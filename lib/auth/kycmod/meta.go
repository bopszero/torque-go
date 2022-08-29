package kycmod

import (
	"regexp"
	"time"

	"gitlab.com/snap-clickstaff/go-common/comcache"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/settingmod"
)

type FeatureDetails struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	IsAvailable bool   `json:"is_available"`
}

type KycEmailTemplateData struct {
	SubmittedUserFirstName string
	LogoUrl                string
}

type KycApproveEmailTemplateData struct {
	TemplateData  KycEmailTemplateData
	OriginalEmail string
}

type KycRejectEmailTemplateData struct {
	TemplateData KycEmailTemplateData
	Reason       string
}

type KycPendingApproveEmailTemplateData struct {
	TemplateData KycEmailTemplateData
	Content      string
}

type KycAddress struct {
	ResidentialAddress string `json:"residential_address"`
	PostalCode         string `json:"postal_code"`
	City               string `json:"city"`
	Country            string `json:"country"`
}

const (
	SystemUID                     = -1
	SubmitAttemptWeight           = 1
	TorqueLogoImageUrl            = "https://torquebot.net/assests/images/logo2.png"
	FeatureCodeName               = "kyc"
	JumioValueNotAvailable        = "N/A"
	InitTime                      = 1599325200 // 2020-09-06
	FullNameDirtyCharacterPattern = `[^\s\p{L}]`
	RequestInitLifeTime           = 24 * time.Hour
	JumioScanInitLifeTime         = 32 * time.Hour
	AllowAge                      = 18
)

const (
	KycEmailSubjectSucceeded                    = "Torque - Identity Verification Successful"
	KycEmailSubjectFailed                       = "Torque - Identity Verification Failed"
	KycEmailSubjectPendingApprovalHaveFund      = "Torque - Unable to process your KYC verification - Restricted Accounts"
	KycEmailSubjectPendingApprovalNoFund        = "Torque - Unable to process your KYC verification - Minimum Balance"
	KycEmailSubjectPendingApprovalHaveAndNoFund = "Torque - Account verification pending"
)

const (
	KycVerificationStatusDOBNotMatch                = "dob_not_match"
	KycVerificationStatusNationalityNotMatch        = "nationality_not_match"
	KycVerificationStatusNotMatchName               = "full_name_not_match"
	KycVerificationStatusJumioComplianceIssue       = "jumio_compliance_issue"
	KycVerificationStatusBannedAccountAfterInitTime = "banned_account_after_init_kyc_time"
	KycVerificationStatusJumioDobMissing            = "jumio_dob_missing"
	KycVerificationStatusSystemError                = "system_error"
	KycVerificationStatusNotAcceptUserType          = "not_accept_user_type"
	KycVerificationStatusDobInvalid                 = "dob_invalid"
	KycVerificationStatusValidityFalse              = "validity_false"
	KycVerificationStatusDifferenceUserCode         = "difference_user_code"
	KycVerificationStatusExpired                    = "expired"
	KycVerificationStatusNoIdUploaded               = "jumio:no_id_uploaded"
)

const (
	EmailDomainGmailPattern      = `^gmail\..*$`
	EmailDomainYahooPattern      = `^yahoo\..*$`
	EmailDomainProtonMailPattern = `^protonmail\..*$`
)

const (
	UserTypeOld     = meta.KycUserType(1)
	UserTypeStandBy = meta.KycUserType(2)
	UserTypeNew     = meta.KycUserType(3)
)

const KeyVerificationStatusJumio = "jumio"

const (
	JumioVerificationStatusApproved        = "APPROVED_VERIFIED"
	JumioVerificationStatusNotReadable     = "ERROR_NOT_READABLE_ID"
	JumioVerificationStatusDeniedIDType    = "DENIED_UNSUPPORTED_ID_TYPE"
	JumioVerificationStatusDeniedIDCountry = "DENIED_UNSUPPORTED_ID_COUNTRY"
	JumioVerificationStatusFraud           = "DENIED_FRAUD"
	JumioVerificationStatusNoIdUploaded    = "NO_ID_UPLOADED"
)

const (
	JumioSimilarityMatch       = "MATCH"
	JumioSimilarityNoMatch     = "NO_MATCH"
	JumioSimilarityNotPossible = "NOT_POSSIBLE"
)

const (
	JumioScanStatusInit        = 1
	JumioScanStatusUpdatedData = 2
	JumioScanStatusExpired     = 3
)

const (
	KycTemplatePathApproved        = "./resources/kyc/email/kyc_approve.html"
	KycTemplatePathRejected        = "./resources/kyc/email/kyc_reject.html"
	KycTemplatePathPendingApproval = "./resources/kyc/email/kyc_pending_approve.html"
)

const (
	KycSpecialUserTypeBlacklist = meta.KycSpecialUserType(-1)
	KycSpecialUserTypeWhitelist = meta.KycSpecialUserType(1)
)

var (
	RequestStatusPriorityMap = map[meta.KycRequestStatus]int{
		constants.KycRequestStatusApproved:        1,
		constants.KycRequestStatusPendingApproval: 10,
		constants.KycRequestStatusPendingAnalysis: 20,
		constants.KycRequestStatusRejected:        30,
		constants.KycRequestStatusInit:            40,
		constants.KycRequestStatusFailed:          100,
	}

	whiteListUIDCached = comcache.NewCacheObject(
		30*time.Second,
		func() (interface{}, error) {
			kycWhiteListUID, err := settingmod.GetSetting(constants.SettingKeyKycWhiteListUID)
			if err != nil {
				return nil, err
			}
			var whiteListUID []meta.UID
			err = comutils.JsonDecode(kycWhiteListUID.Value, &whiteListUID)
			if err != nil {
				return nil, err
			}
			return whiteListUID, nil
		},
	)

	FullNameDirtyCharacterRegex *regexp.Regexp
)

func init() {
	var err error
	FullNameDirtyCharacterRegex, err = regexp.Compile(FullNameDirtyCharacterPattern)
	comutils.PanicOnError(err)
}
