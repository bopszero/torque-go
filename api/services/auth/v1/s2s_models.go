package v1

import "gitlab.com/snap-clickstaff/torque-go/lib/meta"

type S2sKycGetRequest struct {
	UID   meta.UID `json:"uid"`
	Email string   `json:"email" validate:"omitempty,email"`
}

type S2sKycGetResponse struct {
	KycGetResponse
}
