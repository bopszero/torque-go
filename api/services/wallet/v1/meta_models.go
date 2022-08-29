package v1

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/auth/kycmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type MetaHandshakeRequest struct {
	FirebaseToken string `json:"firebase_token"`
}

type MetaHandshakeResponse struct {
	Features           []kycmod.FeatureDetails              `json:"features"`
	BlockchainNetworks []MetaHandshakeBlockchainNetworkInfo `json:"blockchain_networks"`
	NetworkCurrencies  []MetaHandshakeNetworkCurrencyInfo   `json:"network_currencies"`
}

type MetaHandshakeBlockchainNetworkInfo struct {
	Network               meta.BlockchainNetwork `json:"code"`
	Currency              meta.Currency          `json:"currency"`
	Name                  string                 `json:"name"`
	TokenTransferCodeName string                 `json:"token_transfer_code_name"`
}

type MetaHandshakeNetworkCurrencyInfo struct {
	Currency      meta.Currency          `json:"currency"`
	Network       meta.BlockchainNetwork `json:"network"`
	Priority      uint8                  `json:"priority"`
	WithdrawalFee decimal.Decimal        `json:"withdrawal_fee"`
}

type MetaHandshakeCurrencyNetworkInfo struct {
	Currency        meta.Currency              `json:"currency"`
	NetworkInfoList []MetaHandshakeNetworkInfo `json:"network_info_list"`
}

type MetaHandshakeNetworkInfo struct {
	Priority              uint32                 `json:"priority"`
	Network               meta.BlockchainNetwork `json:"network"`
	Name                  string                 `json:"name"`
	TokenTransferCodeName string                 `json:"token_transfer_code_name"`
	WithdrawalFee         decimal.Decimal        `json:"withdrawal_fee"`
}

type MetaCurrencyInfoResponse struct {
	Currencies []MetaCurrencyInfo `json:"currencies"`
}
