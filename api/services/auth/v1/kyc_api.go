package v1

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/api"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/api/responses"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/auth/kycmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/settingmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/thirdpartymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/usermod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func KycInit(c echo.Context) (err error) {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel KycInitRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}
	kycRequest, err := createKycInit(ctx, reqModel)
	if err != nil {
		return err
	}
	jumioCredentials, err := thirdpartymod.GetJumioCredentials()
	if err != nil {
		return err
	}
	return responses.Ok(
		ctx,
		KycInitResponse{
			KycCode:        kycRequest.Code,
			UserCode:       kycRequest.UserCode,
			JumioApiToken:  jumioCredentials.ApiToken,
			JumioApiSecret: jumioCredentials.ApiSecret,
		},
	)
}

func KycSubmit(c echo.Context) (err error) {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		uid      = apiutils.GetContextUidF(ctx)
		reqModel KycSubmitRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}
	if err = kycmod.ValidateUserSubmission(ctx, uid); err != nil {
		return err
	}

	var kycRequest models.KycRequest
	err = database.GetDbF(database.AliasWalletSlave).
		First(&kycRequest, &models.KycRequest{UID: uid, Code: reqModel.KycCode}).
		Error
	if err != nil {
		return utils.WrapError(err)
	}
	if err = kycmod.SubmitRequestByCode(ctx, reqModel.KycCode, reqModel.Reference); err != nil {
		return
	}

	return responses.OkEmpty(ctx)
}

func KycGet(c echo.Context) (err error) {
	var (
		ctx = apiutils.EchoWrapContext(c)
		uid = apiutils.GetContextUidF(ctx)
	)
	kycGetResponse, err := kycGetByUID(ctx, uid)
	if err != nil {
		return err
	}
	return responses.Ok(ctx, kycGetResponse)
}

func kycGetByUID(ctx apiutils.EchoWrappedContext, uid meta.UID) (kycGetResponse KycGetResponse, err error) {
	var blacklistInfo KycBlacklistInfo
	user, err := usermod.GetUserFast(uid)
	if err != nil {
		return
	}
	userPermissionType, err := kycmod.IdentifyUserPermissionType(user)
	if err != nil {
		return
	}
	specialUserInfo, err := kycmod.GetKycSpecialUserInfo(user.ID)
	if err != nil && !utils.IsOurError(err, constants.ErrorCodeDataNotFound) {
		return
	}
	if specialUserInfo.ID > 0 {
		if specialUserInfo.Type == kycmod.KycSpecialUserTypeBlacklist {
			blacklistInfo = KycBlacklistInfo{
				InBlacklist: true,
				Reason:      constants.ErrorKycUserBlacklist.Message(ctx),
			}
		}
		if specialUserInfo.IsPending.Bool {
			kycGetResponse = KycGetResponse{
				Request: &KycRequest{
					Status: constants.KycRequestStatusPendingApproval,
				},
				UserType:      userPermissionType,
				BlacklistInfo: blacklistInfo,
			}
			return
		}
	}

	remainingAttempts, err := kycmod.GetUserRemainingAttempts(ctx, uid)
	if err != nil {
		return
	}
	kycRequest, err := kycmod.GetKycRequestByUserEmail(ctx, user.Email)
	if err != nil {
		return
	}
	if kycRequest == nil {
		kycGetResponse = KycGetResponse{
			Request:           nil,
			UserType:          userPermissionType,
			RemainingAttempts: remainingAttempts,
			BlacklistInfo:     blacklistInfo,
		}
		return
	}
	submittedUser, err := usermod.GetUserFast(kycRequest.UID)
	if err != nil {
		return
	}
	kycRequestNote := kycRequest.Note
	if kycRequest.Status == constants.KycRequestStatusRejected &&
		kycRequest.CloseUID == kycmod.SystemUID {
		verificationStatus := meta.ErrorCode(kycRequest.VerificationStatus)
		kycRequestNote = meta.NewGeneralError(verificationStatus).Message(ctx)
	}
	var address kycmod.KycAddress
	if kycRequest.Address != "" {
		err = comutils.JsonDecode(kycRequest.Address, &address)
		if err != nil {
			return
		}
	}
	kycGetResponse = KycGetResponse{
		Request: &KycRequest{
			ID:                 kycRequest.ID,
			UID:                kycRequest.UID,
			Status:             kycRequest.Status,
			FullName:           kycRequest.FullName,
			Username:           submittedUser.Username,
			Note:               kycRequestNote,
			DOB:                kycRequest.DOB.Format(constants.DateFormatISO),
			Nationality:        kycRequest.Nationality,
			ResidentialAddress: address.ResidentialAddress,
			PostalCode:         address.PostalCode,
			City:               address.City,
			Country:            address.Country,
		},
		UserType:          userPermissionType,
		RemainingAttempts: remainingAttempts,
		BlacklistInfo:     blacklistInfo,
	}
	return kycGetResponse, nil
}

func KycMeta(c echo.Context) (err error) {
	ctx := apiutils.EchoWrapContext(c)
	kycReminderDeadlineStr, err := settingmod.GetSettingValueFast(constants.SettingKeyKycReminderDeadlineTime)
	if err != nil {
		return err
	}
	kycReminderIsEnabledStr, err := settingmod.GetSettingValueFast(constants.SettingKeyKycReminderIsEnabled)
	if err != nil {
		return err
	}
	reminderDeadlineTime, err := comutils.ParseInt64(kycReminderDeadlineStr)
	if err != nil {
		return err
	}
	reminderIsEnabled, err := strconv.ParseBool(kycReminderIsEnabledStr)
	if err != nil {
		return err
	}
	return responses.Ok(
		ctx,
		KycMetaResponse{
			Reminder: KycReminder{
				DeadlineTime: reminderDeadlineTime,
				IsEnabled:    reminderIsEnabled,
			},
		},
	)
}

func createKycInit(ctx comcontext.Context, reqModel KycInitRequest) (
	kycRequest *models.KycRequest, err error,
) {
	uid := apiutils.GetContextUidF(ctx)
	if err = kycmod.ValidateUserSubmission(ctx, uid); err != nil {
		return nil, err
	}
	user, err := usermod.GetUserFast(uid)
	if err != nil {
		return nil, err
	}
	dobParsed, err := time.Parse(constants.DateFormatISO, reqModel.DOB)
	if err != nil {
		return nil, utils.WrapError(err)
	}
	if !kycmod.IsValidUserAge(dobParsed) {
		return nil, utils.WrapError(constants.ErrorDobInvalid)
	}
	address := kycmod.KycAddress{
		ResidentialAddress: reqModel.ResidentialAddress,
		PostalCode:         reqModel.PostalCode,
		City:               reqModel.City,
		Country:            reqModel.Country,
	}
	addressJSON, err := comutils.JsonEncode(address)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	kycRequest = &models.KycRequest{
		Code:          comutils.NewUuidCode(),
		UID:           uid,
		UserCode:      user.Code,
		EmailOriginal: kycmod.ParseEmailOriginal(ctx, user.Email),
		FullName:      utils.StringTrim(reqModel.FullName),
		DOB:           dobParsed,
		Status:        constants.KycRequestStatusInit,
		Nationality:   reqModel.Nationality,
		Address:       addressJSON,
		CreateTime:    now.Unix(),
		UpdateTime:    now.Unix(),
	}
	err = database.TransactionFromContext(ctx, database.AliasWalletMaster, func(dbTxn *gorm.DB) error {
		return dbTxn.Create(&kycRequest).Error
	})
	if err != nil {
		return nil, err
	}
	return kycRequest, nil
}

func KycInitUrl(c echo.Context) (err error) {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel KycInitUrlRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}
	kycRequest, err := createKycInit(ctx, reqModel.KycInitRequest)
	if err != nil {
		return err
	}
	jumioClient, err := thirdpartymod.GetJumioServiceSystemClient()
	if err != nil {
		return err
	}
	resp, err := jumioClient.InitTransaction(ctx, kycRequest.Code, kycRequest.UserCode, reqModel.Locale)
	if err != nil {
		return err
	}
	var transactionInitResponse KycJumioTransactionInitResponse
	if err = comutils.JsonDecode(resp.String(), &transactionInitResponse); err != nil {
		return err
	}
	return responses.Ok(
		ctx,
		KycInitUrlResponse{
			RedirectUrl: transactionInitResponse.RedirectUrl,
		},
	)
}

func KycValidateEmail(c echo.Context) (err error) {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel KycValidateEmailRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}
	kycRequest, err := kycmod.GetKycRequestByUserEmail(ctx, reqModel.Email)
	if err != nil {
		return err
	}
	if kycRequest == nil {
		return responses.Ok(ctx, KycValidateEmailResponse{IsValidEmail: true})
	}
	if kycRequest.Status == constants.KycRequestStatusApproved {
		scanDetailsResp, err := kycmod.GetJumioScanDetails(ctx, kycRequest.Reference)
		if err != nil {
			return err
		}
		var country models.Country
		err = database.GetDbSlave().
			First(&country, &models.Country{CodeIso3: scanDetailsResp.Document.IssuingCountry}).
			Error
		if err != nil {
			return err
		}
		if country.IsBanned.Bool && !kycmod.IsUserInWhiteListKyc(kycRequest.UID) {
			return responses.Ok(
				ctx,
				KycValidateEmailResponse{
					IsValidEmail: false,
					Message:      constants.ErrorEmailInvalid.Message(ctx),
				},
			)
		}
	}
	return responses.Ok(ctx, KycValidateEmailResponse{IsValidEmail: true})
}
