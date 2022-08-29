package v1

import (
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type KycInitRequest struct {
	FullName           string `json:"full_name" validate:"required"`
	DOB                string `json:"dob" validate:"required,datetime=2006-01-02"`
	Nationality        string `json:"nationality" validate:"required"`
	ResidentialAddress string `json:"residential_address" validate:"required"`
	PostalCode         string `json:"postal_code" validate:"required"`
	City               string `json:"city" validate:"required"`
	Country            string `json:"country" validate:"required"`
}

type KycInitResponse struct {
	KycCode        string `json:"kyc_code"`
	UserCode       string `json:"user_code"`
	JumioApiToken  string `json:"jumio_api_token"`
	JumioApiSecret string `json:"jumio_api_secret"`
}

type KycInitUrlRequest struct {
	KycInitRequest
	Locale string `json:"locale"`
}

type KycInitUrlResponse struct {
	RedirectUrl string `json:"redirect_url"`
}

type KycSubmitRequest struct {
	KycCode   string `json:"kyc_code" validate:"required"`
	Reference string `json:"reference" validate:"required"`
}

type KycGetResponse struct {
	Request           *KycRequest      `json:"request"`
	UserType          meta.KycUserType `json:"user_type"`
	RemainingAttempts int64            `json:"remaining_attempts"`
	BlacklistInfo     KycBlacklistInfo `json:"blacklist_info"`
}

type KycBlacklistInfo struct {
	InBlacklist bool   `json:"in_blacklist"`
	Reason      string `json:"reason"`
}

type KycRequest struct {
	ID                 uint64                `json:"id"`
	UID                meta.UID              `json:"uid"`
	Status             meta.KycRequestStatus `json:"status"`
	FullName           string                `json:"full_name"`
	Username           string                `json:"username"`
	Note               string                `json:"note"`
	DOB                string                `json:"dob" `
	Nationality        string                `json:"nationality" `
	ResidentialAddress string                `json:"residential_address" `
	PostalCode         string                `json:"postal_code" `
	City               string                `json:"city" `
	Country            string                `json:"country" `
}

type JumioPushScanResultRequest struct {
	RequestCode          string `form:"merchantIdScanReference" validate:"required"`
	Reference            string `form:"jumioIdScanReference" validate:"required"`
	IdentityVerification string `form:"identityVerification"`
	VerificationStatus   string `form:"verificationStatus"`
	RejectReason         string `form:"rejectReason"`
	AdditionalChecks     string `form:"additionalChecks"`
}

type KycReminder struct {
	DeadlineTime int64 `json:"deadline_time"`
	IsEnabled    bool  `json:"is_enabled"`
}

type KycMetaResponse struct {
	Reminder KycReminder `json:"reminder"`
}

type KycJumioTransactionInitRequest struct {
	CustomerInternalReference string `json:"customerInternalReference" validate:"required"`
	UserReference             string `json:"userReference" validate:"required"`
}

type KycJumioTransactionInitResponse struct {
	RedirectUrl          string `json:"redirectUrl"`
	Timestamp            string `json:"timestamp"`
	TransactionReference string `json:"transactionReference"`
}

type KycJumioRedirectRequest struct {
	RequestCode       string `query:"customerInternalReference" validate:"required"`
	Reference         string `query:"transactionReference" validate:"required"`
	TransactionStatus string `query:"transactionStatus"`
}

type KycValidateEmailRequest struct {
	Email string `json:"email" validate:"required"`
}

type KycValidateEmailResponse struct {
	IsValidEmail bool   `json:"is_valid_email"`
	Message      string `json:"message"`
}

type KycSendEmailRequest struct {
	RequestId     uint64                `json:"request_id" validate:"required"`
	RequestStatus meta.KycRequestStatus `json:"request_status" validate:"required"`
}
