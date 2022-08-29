package kycmod

import (
	"reflect"

	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type JumioTransaction struct {
	UserCode string `json:"customerId"`
	KycCode  string `json:"merchantScanReference"`
}
type Document struct {
	Status         string `json:"status"`
	IssuingCountry string `json:"issuingCountry"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	DOB            string `json:"dob"`
	Nationality    string `json:"nationality"`
}

type Verification struct {
	IdentityVerification IdentityVerification `json:"identityVerification"`
	RejectReason         RejectReason         `json:"rejectReason"`
}

type IdentityVerification struct {
	Similarity string `json:"similarity"`
	Validity   bool   `json:"validity,string"`
	Reason     string `json:"reason"`
}

type RejectReason struct {
	Code    string      `json:"rejectReasonCode"`
	Details interface{} `json:"rejectReasonDetails"`
}

func (this RejectReason) GetDetails() (detailsList []RejectReasonDetails) {
	if this.Details == nil {
		return
	}
	switch reflect.TypeOf(this.Details).Kind() {
	case reflect.Array, reflect.Slice:
		if err := utils.DumpDataByJSON(this.Details, &detailsList); err != nil {
			panic(err)
		}
		break
	default:
		var detail RejectReasonDetails
		if err := utils.DumpDataByJSON(this.Details, &detail); err != nil {
			panic(err)
		}
		detailsList = append(detailsList, detail)
		break
	}
	return detailsList
}

type RejectReasonDetails struct {
	Code string `json:"detailsCode"`
}

type JumioScanDetailsResponse struct {
	Transaction   JumioTransaction `json:"transaction"`
	ScanReference string           `json:"scanReference"`
	Document      Document         `json:"document"`
	Verification  Verification     `json:"verification"`
}

type JumioDataVerificationResponse struct {
	ScanReference    string           `json:"scanReference"`
	AdditionalChecks AdditionalChecks `json:"additionalChecks"`
}

type AdditionalChecks struct {
	WatchlistScreening WatchlistScreening `json:"watchlistScreening"`
}

type WatchlistScreening struct {
	SearchResults int `json:"searchResults,string"`
}
