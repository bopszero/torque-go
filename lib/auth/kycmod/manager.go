package kycmod

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/settingmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/usermod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func isValidKycStatus(ctx comcontext.Context, uid meta.UID) bool {
	user, err := usermod.GetUserFast(uid)
	if err != nil {
		return false
	}
	kycRequest, err := GetKycRequestByUserEmail(ctx, user.Email)
	if err != nil {
		return false
	}
	if kycRequest == nil {
		return true
	}
	if kycRequest.Status == constants.KycRequestStatusPendingAnalysis ||
		kycRequest.Status == constants.KycRequestStatusApproved ||
		kycRequest.Status == constants.KycRequestStatusPendingApproval {
		return false
	}
	return true
}

func IdentifyUserType(user models.User) (meta.KycUserType, error) {
	kycLaunchTime, err := settingmod.GetSettingValueFast(constants.SettingKeyKycLaunchTime)
	if err != nil {
		return 0, err
	}
	userCreateTime := user.CreateDate.Unix()
	if userCreateTime < InitTime {
		return UserTypeOld, nil
	} else if InitTime <= userCreateTime &&
		userCreateTime < comutils.ParseInt64F(kycLaunchTime) {
		return UserTypeStandBy, nil
	}
	return UserTypeNew, nil
}

func ValidateUserSubmission(ctx comcontext.Context, uid meta.UID) (err error) {
	if !isValidKycStatus(ctx, uid) {
		return utils.WrapError(constants.ErrorKycRequestInvalidStatus)
	}
	if isUserSubmittedOverLimit(ctx, uid) {
		return utils.WrapError(constants.ErrorUserSubmittedOverLimit)
	}
	if err = validateSpecialInfo(uid); err != nil {
		return err
	}
	return nil
}

func isUserSubmittedOverLimit(ctx comcontext.Context, uid meta.UID) bool {
	err, weightSubmittedKyc := getWeightUserSubmittedKyc(ctx, uid)
	if err != nil {
		return false
	}
	err, kycSubmitLimit := getKycSubmitLimit()
	if err != nil {
		return false
	}
	return weightSubmittedKyc >= kycSubmitLimit
}

func validateSpecialInfo(uid meta.UID) error {
	specialUserInfo, err := GetKycSpecialUserInfo(uid)
	if err != nil && !utils.IsOurError(err, constants.ErrorCodeDataNotFound) {
		return err
	}
	if specialUserInfo.ID == 0 {
		return nil
	}
	if specialUserInfo.Type == KycSpecialUserTypeBlacklist {
		return utils.WrapError(constants.ErrorKycUserBlacklist)
	}
	if specialUserInfo.IsPending.Bool {
		return utils.WrapError(constants.ErrorKycRequestInvalidStatus)
	}
	return nil
}

func IsValidUserAge(userDOB time.Time) bool {
	// -1 for don't need care timezone
	return userDOB.AddDate(AllowAge, 0, -1).Before(time.Now())
}

func getWeightUserSubmittedKyc(ctx comcontext.Context, uid meta.UID) (error, int64) {
	user, err := usermod.GetUserFast(uid)
	if err != nil {
		return err, 0
	}
	var kycRequestList []models.KycRequest
	err = database.GetDbF(database.AliasWalletSlave).
		Model(&models.KycRequest{}).
		Where(&models.KycRequest{
			EmailOriginal: ParseEmailOriginal(ctx, user.Email),
		}).
		Find(&kycRequestList).
		Error
	if err != nil {
		return err, 0
	}
	var weightSubmittedKyc int64
	for _, kycRequest := range kycRequestList {
		weightSubmittedKyc += kycRequest.AttemptWeight
	}
	return nil, weightSubmittedKyc
}

func getKycSubmitLimit() (error, int64) {
	kycSubmitLimitStr, err := settingmod.GetSettingValueFast(constants.SettingKeyKycSubmitLimit)
	if err != nil {
		return err, 0
	}
	kycSubmitLimit, err := comutils.ParseInt64(kycSubmitLimitStr)
	if err != nil {
		return err, 0
	}
	return nil, kycSubmitLimit
}

func IdentifyUserPermissionType(user models.User) (meta.KycUserType, error) {
	kycLaunchTime, err := settingmod.GetSettingValueFast(constants.SettingKeyKycLaunchTime)
	if err != nil {
		return 0, err
	}
	kycReminderDeadlineStr, err := settingmod.GetSettingValueFast(constants.SettingKeyKycReminderDeadlineTime)
	if err != nil {
		return 0, err
	}
	userCreateTime := user.CreateDate.Unix()
	switch true {
	case userCreateTime < InitTime:
		return UserTypeOld, nil
	case InitTime <= userCreateTime &&
		userCreateTime < comutils.ParseInt64F(kycLaunchTime) &&
		time.Now().Unix() < comutils.ParseInt64F(kycReminderDeadlineStr):
		return UserTypeOld, nil
	case InitTime <= userCreateTime &&
		userCreateTime < comutils.ParseInt64F(kycLaunchTime):
		return UserTypeStandBy, nil
	}
	return UserTypeNew, nil
}

func GetMetaFeaturesSetting() ([]FeatureDetails, error) {
	metaFeaturesStr, err := settingmod.GetSettingValueFast(constants.SettingKeyMetaFeatures)
	if err != nil {
		return nil, err
	}
	if metaFeaturesStr == "" {
		return []FeatureDetails{}, nil
	}

	var metaFeatures []FeatureDetails
	if err := comutils.JsonDecode(metaFeaturesStr, &metaFeatures); err != nil {
		return nil, utils.WrapError(err)
	}
	return metaFeatures, nil
}

func SubmitRequestByCode(ctx comcontext.Context, requestCode, scanRef string) (err error) {
	return database.TransactionFromContext(ctx, database.AliasWalletMaster, func(dbTxn *gorm.DB) (err error) {
		var request models.KycRequest
		err = dbquery.SelectForUpdate(dbTxn).
			First(&request, &models.KycRequest{Code: requestCode}).
			Error
		if err != nil {
			return utils.WrapError(err)
		}
		if request.Status != constants.KycRequestStatusInit {
			return nil
		}

		request.Status = constants.KycRequestStatusPendingAnalysis
		request.AttemptWeight = SubmitAttemptWeight
		request.UpdateTime = time.Now().Unix()
		if err = dbTxn.Save(&request).Error; err != nil {
			return utils.WrapError(err)
		}

		_, err = initJumioScan(ctx, scanRef)
		return
	})
}

func ExecuteRequest(ctx comcontext.Context, requestID uint64) (request models.KycRequest, err error) {
	err = database.Atomic(ctx, database.AliasWalletMaster, func(dbTxn *gorm.DB) (err error) {
		err = dbquery.SelectForUpdate(dbTxn).
			First(&request, &models.KycRequest{ID: requestID}).
			Error
		if err != nil {
			return utils.WrapError(err)
		}

		now := time.Now()
		switch request.Status {
		case constants.KycRequestStatusInit:
			if request.CreateTime < now.Add(-RequestInitLifeTime).Unix() {
				request.Status = constants.KycRequestStatusFailed
				request.NoteInternal = KycVerificationStatusExpired
				request.CloseUID = SystemUID
				request.CloseTime = now.Unix()
				request.UpdateTime = now.Unix()
				if err = dbTxn.Save(&request).Error; err != nil {
					err = utils.WrapError(err)
				}
				return
			}
			break
		case constants.KycRequestStatusPendingAnalysis:
			break
		default:
			return nil
		}

		var jumioScan models.JumioScan
		err = dbTxn.
			Where(dbquery.NotEqual(models.JumioScanColStatus, JumioScanStatusInit)).
			First(&jumioScan, &models.JumioScan{RequestCode: request.Code}).
			Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil
			}
			return utils.WrapError(err)
		}

		request.Reference = jumioScan.Reference

		switch {
		case request.UserCode != jumioScan.UserCode:
			request.Status = constants.KycRequestStatusFailed
			request.NoteInternal = KycVerificationStatusDifferenceUserCode
			break
		case jumioScan.Status == JumioScanStatusExpired:
			request.Status = constants.KycRequestStatusFailed
			request.NoteInternal = KycVerificationStatusExpired
			break
		case jumioScan.Status == JumioScanStatusUpdatedData:
			if err = judgeRequestStatus(ctx, &request); err != nil {
				return
			}
			break
		default:
			comlogging.GetLogger().
				WithType(constants.LogTypeKYC).
				WithContext(ctx).
				WithFields(logrus.Fields{
					"id":     jumioScan.ID,
					"status": jumioScan.Status,
				}).
				Warnf("jumio scan has an invalid status | status=%v", jumioScan.Status)
			return nil
		}

		switch request.Status {
		case constants.KycRequestStatusApproved, constants.KycRequestStatusRejected, constants.KycRequestStatusFailed:
			request.CloseUID = SystemUID
			request.CloseTime = now.Unix()
			break
		default:
			break
		}

		request.UpdateTime = now.Unix()
		if err = dbTxn.Save(&request).Error; err != nil {
			return utils.WrapError(err)
		}

		_, err = SendRequestEmail(ctx, request)
		if err != nil {
			return err
		}
		return nil
	})
	return
}

func judgeRequestStatus(ctx comcontext.Context, request *models.KycRequest) (err error) {
	scanData, err := GetJumioScanDetails(ctx, request.Reference)
	if err != nil {
		return err
	}
	verificationData, err := GetJumioVerificationData(ctx, request.Reference)
	if err != nil {
		return err
	}

	switch scanData.Document.Status {
	case JumioVerificationStatusApproved:
		return ProcessKycApproval(
			ctx,
			request,
			scanData.Verification.IdentityVerification,
			verificationData.AdditionalChecks,
			scanData.Document,
		)
	case JumioVerificationStatusNotReadable:
		verificationStatus := genVerificationStatusNotReadable(
			scanData.Document.Status,
			scanData.Verification.RejectReason,
		)
		request.Status = constants.KycRequestStatusRejected
		request.VerificationStatus = verificationStatus
		return nil
	default:
		request.Status = constants.KycRequestStatusRejected
		request.VerificationStatus = fmt.Sprintf(
			"%v:%v",
			KeyVerificationStatusJumio,
			strings.ToLower(scanData.Document.Status),
		)
		if scanData.Document.Status == JumioVerificationStatusNoIdUploaded {
			request.AttemptWeight = 0
		}
		return nil
	}
}

func GetKycRequestByUserEmail(ctx comcontext.Context, userEmail string) (*models.KycRequest, error) {
	var requests []models.KycRequest
	err := database.GetDbF(database.AliasWalletSlave).
		Where(&models.KycRequest{
			EmailOriginal: ParseEmailOriginal(ctx, userEmail),
		}).
		Find(&requests).
		Error
	if err != nil {
		return nil, utils.WrapError(err)
	}
	if len(requests) == 0 {
		return nil, nil
	}
	bestRequest := requests[0]
	for _, request := range requests[1:] {
		if isPreferLeftRequest(request, bestRequest) {
			bestRequest = request
		}
	}
	return &bestRequest, nil
}

func GetUserKycRequest(ctx comcontext.Context, uid meta.UID) (*models.KycRequest, error) {
	user, err := usermod.GetUserFast(uid)
	if err != nil {
		return nil, err
	}
	return GetKycRequestByUserEmail(ctx, user.OriginalEmail)
}

func ValidateUserKYC(ctx comcontext.Context, uid meta.UID) error {
	req, err := GetUserKycRequest(ctx, uid)
	if err != nil {
		return err
	}
	if req == nil || req.Status != constants.KycRequestStatusApproved {
		return utils.WrapError(constants.ErrorKycRequired)
	}
	return nil
}

func GetUserRemainingAttempts(ctx comcontext.Context, uid meta.UID) (int64, error) {
	err, weightSubmittedKyc := getWeightUserSubmittedKyc(ctx, uid)
	if err != nil {
		return 0, err
	}
	err, kycSubmitLimit := getKycSubmitLimit()
	if err != nil {
		return 0, err
	}
	remainingAttempts := kycSubmitLimit - weightSubmittedKyc
	if remainingAttempts < 0 {
		remainingAttempts = 0
	}
	return remainingAttempts, nil
}

func GetKycSpecialUserInfo(uid meta.UID) (specialUser models.KycSpecialUserInfo, err error) {
	err = database.GetDbF(database.AliasWalletSlave).
		Take(&specialUser, &models.KycSpecialUserInfo{UID: uid}).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = constants.ErrorDataNotFound
		}
		err = utils.WrapError(err)
	}
	return
}
