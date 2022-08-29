package v1

import (
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type HelperBlockchainValidateAddressRequest struct {
	Currency meta.Currency          `json:"currency" validate:"required"`
	Network  meta.BlockchainNetwork `json:"network"`
	Address  string                 `json:"address" validate:"required"`
}

type HelperBlockchainValidateAddressResponse struct {
	IsValid       bool   `json:"is_valid"`
	Address       string `json:"address"`
	AddressLegacy string `json:"address_legacy"`
}
