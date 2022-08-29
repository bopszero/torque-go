package kycmod

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/trading/tradingbalance"
	"gitlab.com/snap-clickstaff/torque-go/lib/usermod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func ProcessKycApproval(
	ctx comcontext.Context,
	kycRequest *models.KycRequest, idVerification IdentityVerification, additionalChecks AdditionalChecks,
	document Document,
) (err error) {
	switch idVerification.Similarity {
	case JumioSimilarityMatch:
		if !idVerification.Validity {
			setKycValidityFalse(kycRequest)
		} else {
			setKycSimilarityMatch(ctx, kycRequest, additionalChecks, document)
		}
		break
	case JumioSimilarityNoMatch:
		setKycSimilarityNoMatch(kycRequest)
		break
	case JumioSimilarityNotPossible:
		setKycSimilarityNotPossible(kycRequest, idVerification)
		break
	default:
		kycRequest.Status = constants.KycRequestStatusPendingApproval
		break
	}
	return nil
}

func setKycValidityFalse(kycRequest *models.KycRequest) {
	kycRequest.Status = constants.KycRequestStatusRejected
	kycRequest.VerificationStatus = genKycRejectVerificationStatus(
		JumioVerificationStatusApproved,
		KycVerificationStatusValidityFalse,
	)
}

func setKycSimilarityMatch(
	ctx comcontext.Context,
	kycRequest *models.KycRequest, additionalChecks AdditionalChecks, document Document,
) {
	if !isAcceptUserTypeKyc(kycRequest.UID, document.IssuingCountry) {
		kycRequest.Status = constants.KycRequestStatusRejected
		kycRequest.VerificationStatus = genKycRejectVerificationStatus(
			JumioVerificationStatusApproved,
			KycVerificationStatusNotAcceptUserType,
		)
		return
	}

	if document.DOB != "" {
		dobParsed, err := time.Parse(constants.DateFormatISO, document.DOB)
		if err != nil {
			kycRequest.Status = constants.KycRequestStatusPendingApproval
			kycRequest.VerificationStatus = KycVerificationStatusSystemError
			return
		}
		if !IsValidUserAge(dobParsed) {
			kycRequest.Status = constants.KycRequestStatusRejected
			kycRequest.VerificationStatus = genKycRejectVerificationStatus(
				JumioVerificationStatusApproved,
				KycVerificationStatusDobInvalid,
			)
			return
		}
	}

	if !isNeedAdminApproveKyc(ctx, kycRequest, additionalChecks, document) {
		kycRequest.Status = constants.KycRequestStatusApproved
	} else {
		kycRequest.Status = constants.KycRequestStatusPendingApproval
	}
}

func setKycSimilarityNoMatch(kycRequest *models.KycRequest) {
	kycRequest.Status = constants.KycRequestStatusRejected
	kycRequest.VerificationStatus = genKycRejectVerificationStatus(
		JumioVerificationStatusApproved,
		JumioSimilarityNoMatch,
	)
}

func setKycSimilarityNotPossible(kycRequest *models.KycRequest, idVerification IdentityVerification) {
	kycRequest.VerificationStatus = genKycRejectVerificationStatus(
		JumioVerificationStatusApproved,
		JumioSimilarityNotPossible,
	)
	if idVerification.Reason != "" {
		kycRequest.VerificationStatus = genKycRejectVerificationStatus(JumioSimilarityNotPossible, idVerification.Reason)
	}
	kycRequest.Status = constants.KycRequestStatusRejected
}

func isAcceptUserTypeKyc(uid meta.UID, issuingCountry string) bool {
	user, err := usermod.GetUserFast(uid)
	if err != nil {
		return false
	}
	userType, err := IdentifyUserType(user)
	if err != nil {
		return false
	}
	if userType == UserTypeOld {
		return true
	} else if userType == UserTypeStandBy || userType == UserTypeNew {
		var country models.Country
		err = database.GetDbSlave().First(&country, &models.Country{CodeIso3: issuingCountry}).Error
		if err != nil {
			return false
		}
		if !country.IsBanned.Bool {
			return true
		}
		if country.IsBanned.Bool && IsUserInWhiteListKyc(uid) {
			return true
		}
	}
	return false
}

func IsUserInWhiteListKyc(uidKyc meta.UID) bool {
	whiteListUID, err := whiteListUIDCached.Get()
	if err != nil {
		return false
	}
	for _, uid := range whiteListUID.([]meta.UID) {
		if uid == uidKyc {
			return true
		}
	}
	return false
}

func isNeedAdminApproveKyc(
	ctx comcontext.Context,
	kycRequest *models.KycRequest, additionalChecks AdditionalChecks, document Document,
) bool {
	isUserHaveBannedAccount, err := isUserHasBannedAccountAfterInitTime(kycRequest, document)
	if err != nil {
		kycRequest.VerificationStatus = KycVerificationStatusSystemError
		return true
	}
	if isUserHaveBannedAccount {
		kycRequest.VerificationStatus = KycVerificationStatusBannedAccountAfterInitTime
		return true
	}
	//if kycRequest.DOB.Format(constants.DateFormatISO) != document.DOB {
	//	kycRequest.VerificationStatus = KycVerificationStatusDOBNotMatch
	//	return true
	//}
	//if !isMatchDocumentNationality(ctx, kycRequest.Nationality, document) {
	//	kycRequest.VerificationStatus = KycVerificationStatusNationalityNotMatch
	//	return true
	//}
	if additionalChecks.WatchlistScreening.SearchResults > 0 {
		kycRequest.VerificationStatus = KycVerificationStatusJumioComplianceIssue
		return true
	}
	if document.DOB == "" {
		kycRequest.VerificationStatus = KycVerificationStatusJumioDobMissing
		return true
	}
	//if !isMatchFullName(kycRequest, document) {
	//	kycRequest.VerificationStatus = KycVerificationStatusNotMatchName
	//	return true
	//}
	return false
}

func isMatchFullName(kycRequest *models.KycRequest, document Document) bool {
	requestFullName := TruncateName(kycRequest.FullName)
	if requestFullName == "" {
		return false
	}

	if document.FirstName == JumioValueNotAvailable {
		document.FirstName = ""
	}
	if document.LastName == JumioValueNotAvailable {
		document.LastName = ""
	}
	var (
		fullNameFist = TruncateName(fmt.Sprintf("%s %s", document.FirstName, document.LastName))
		fullNameLast = TruncateName(fmt.Sprintf("%s %s", document.LastName, document.FirstName))
	)
	return utils.IsSameStringCI(fullNameFist, requestFullName) ||
		utils.IsSameStringCI(fullNameLast, requestFullName)
}

func isUserHasBannedAccountAfterInitTime(kycRequest *models.KycRequest, document Document) (bool, error) {
	var country models.Country
	err := database.GetDbSlave().
		First(&country, &models.Country{CodeIso3: document.IssuingCountry}).
		Error
	if err != nil {
		return false, err
	}
	if !country.IsBanned.Bool {
		return false, nil
	}
	var users []models.User
	err = database.GetDbSlave().
		Where(&models.User{OriginalEmail: kycRequest.EmailOriginal}).
		Find(&users).
		Error
	if err != nil {
		return false, err
	}
	for _, user := range users {
		if user.CreateDate.Unix() >= InitTime {
			return true, nil
		}
	}
	return false, nil
}

func genVerificationStatusNotReadable(verificationStatus string, rejectReason RejectReason) string {
	if rejectReason.Code == "" {
		return fmt.Sprintf("%v:%v", KeyVerificationStatusJumio, strings.ToLower(verificationStatus))
	}
	rejectDetailsList := rejectReason.GetDetails()
	if len(rejectDetailsList) > 0 {
		return fmt.Sprintf(
			"%v:%v:%v_%v",
			KeyVerificationStatusJumio,
			strings.ToLower(verificationStatus),
			rejectReason.Code,
			rejectDetailsList[0].Code,
		)
	}
	return genKycRejectVerificationStatus(verificationStatus, rejectReason.Code)
}

func SendRequestEmail(ctx comcontext.Context, request models.KycRequest) (isSent bool, err error) {
	submittedUser, err := usermod.GetUserFast(request.UID)
	if err != nil {
		return false, err
	}
	switch request.Status {
	case constants.KycRequestStatusApproved:
		return sendEmailKycApproved(ctx, submittedUser, request)
	case constants.KycRequestStatusRejected:
		return sendEmailKycRejected(ctx, submittedUser, request)
	case constants.KycRequestStatusPendingApproval:
		return sendEmailKycPendingApproval(ctx, submittedUser, request)
	default:
		return false, nil
	}
}

func sendEmailKycApproved(
	ctx comcontext.Context,
	submittedUser models.User, kycRequest models.KycRequest,
) (isSent bool, err error) {
	templateData := KycApproveEmailTemplateData{
		TemplateData: KycEmailTemplateData{
			SubmittedUserFirstName: submittedUser.FirstName,
			LogoUrl:                TorqueLogoImageUrl,
		},
		OriginalEmail: kycRequest.EmailOriginal,
	}
	err = sendEmailTemplate(
		ctx,
		submittedUser.Email, KycEmailSubjectSucceeded,
		KycTemplatePathApproved, templateData,
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

func sendEmailKycRejected(
	ctx comcontext.Context,
	submittedUser models.User, kycRequest models.KycRequest,
) (isSent bool, err error) {
	if kycRequest.VerificationStatus == KycVerificationStatusNoIdUploaded {
		return false, nil
	}
	var reason string
	if kycRequest.CloseUID == SystemUID {
		ctx := comcontext.NewContext()
		reason = meta.NewGeneralError(meta.ErrorCode(kycRequest.VerificationStatus)).Message(ctx)
	} else {
		reason = kycRequest.Note
	}
	templateData := KycRejectEmailTemplateData{
		TemplateData: KycEmailTemplateData{
			SubmittedUserFirstName: submittedUser.FirstName,
			LogoUrl:                TorqueLogoImageUrl,
		},
		Reason: reason,
	}
	err = sendEmailTemplate(
		ctx,
		submittedUser.Email, KycEmailSubjectFailed,
		KycTemplatePathRejected, templateData,
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

func sendEmailKycPendingApproval(
	ctx comcontext.Context,
	submittedUser models.User, kycRequest models.KycRequest,
) (isSent bool, err error) {
	if kycRequest.VerificationStatus != KycVerificationStatusBannedAccountAfterInitTime {
		return false, nil
	}
	subject, err := getEmailSubjectPendingApproval(kycRequest)
	if err != nil {
		return false, err
	}
	var content string
	switch subject {
	case KycEmailSubjectPendingApprovalHaveFund:
		content = kycContentHaveFund
	case KycEmailSubjectPendingApprovalNoFund:
		content = kycContentNoFund
	case KycEmailSubjectPendingApprovalHaveAndNoFund:
		content = kycContentHaveAndNoFund
	}
	templateData := KycPendingApproveEmailTemplateData{
		TemplateData: KycEmailTemplateData{
			SubmittedUserFirstName: submittedUser.FirstName,
			LogoUrl:                TorqueLogoImageUrl,
		},
		Content: content,
	}
	err = sendEmailTemplate(
		ctx,
		submittedUser.Email, subject,
		KycTemplatePathPendingApproval, templateData,
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

func getEmailSubjectPendingApproval(kycRequest models.KycRequest) (string, error) {
	ctx := comcontext.NewContext()
	var usersCreatedAfterInitTime []models.User
	err := database.GetDbSlave().
		Where(&models.User{OriginalEmail: kycRequest.EmailOriginal}).
		Where(dbquery.Gt(models.UserColCreateDate, time.Unix(InitTime, 0))).
		Find(&usersCreatedAfterInitTime).
		Error
	if err != nil {
		return "", err
	}
	countUsersHaveFund := 0
	for _, user := range usersCreatedAfterInitTime {
		userHaveFundTrading := 0
		userBalancesTrading := tradingbalance.GetUserBalances(ctx, user.ID)
		for _, userBalanceTrading := range userBalancesTrading {
			if userBalanceTrading.Amount.GreaterThan(decimal.Zero) {
				countUsersHaveFund++
				userHaveFundTrading++
				break
			}
		}
		// if userHaveFundTrading == 0 {
		// 	userBalances, err := balancemod.GetUserBalances(ctx, user.ID)
		// 	if err != nil {
		// 		return "", err
		// 	}
		// 	for _, userBalance := range userBalances {
		// 		if userBalance.Amount.GreaterThan(decimal.Zero) {
		// 			countUsersHaveFund++
		// 			break
		// 		}
		// 	}
		// }
	}
	if countUsersHaveFund == len(usersCreatedAfterInitTime) {
		return KycEmailSubjectPendingApprovalHaveFund, nil
	} else if countUsersHaveFund == 0 {
		return KycEmailSubjectPendingApprovalNoFund, nil
	} else {
		return KycEmailSubjectPendingApprovalHaveAndNoFund, nil
	}
}
