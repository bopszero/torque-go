package v1

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type WithdrawalAccountGenerateRequest struct {
	Currency meta.Currency          `json:"currency" validate:"required"`
	Network  meta.BlockchainNetwork `json:"network"`
}

type WithdrawalAccountGenerateResponse struct {
	Currency  meta.Currency          `json:"currency"`
	Network   meta.BlockchainNetwork `json:"network"`
	AccountNo string                 `json:"account_no"`
}

type WithdrawalAccountGetRequest struct {
	Currency meta.Currency          `json:"currency" validate:"required"`
	Network  meta.BlockchainNetwork `json:"network"`
}

type WithdrawalAccountGetResponse struct {
	Currency    meta.Currency          `json:"currency"`
	Network     meta.BlockchainNetwork `json:"network"`
	AccountNo   string                 `json:"account_no"`
	PullAddress string                 `json:"pull_address"`
}

type WithdrawalAccountPullRequest struct {
	Currency    meta.Currency          `json:"currency" validate:"required"`
	Network     meta.BlockchainNetwork `json:"network"`
	AccountNo   string                 `json:"account_no" validate:"required"`
	PullAddress string                 `json:"pull_address" validate:"required"`
}

type WithdrawalAccountPullResponse struct {
	Currency meta.Currency          `json:"currency"`
	Network  meta.BlockchainNetwork `json:"network"`
	Hash     string                 `json:"hash"`
}

type WithdrawalTransferMeta struct {
	Currency meta.Currency          `json:"currency"`
	Network  meta.BlockchainNetwork `json:"network"`
	FeeInfo  blockchainmod.FeeInfo  `json:"fee_info"`
}

type WithdrawalTransferMetaResponse struct {
	Currencies []WithdrawalTransferMeta `json:"currencies"`
}

type WithdrawalTransferSubmitRequest struct {
	RequestUID meta.UID `json:"request_uid" validate:"required"`

	Currency    meta.Currency          `json:"currency" validate:"required"`
	Network     meta.BlockchainNetwork `json:"network"`
	SrcAddress  string                 `json:"src_address" validate:"required"`
	Codes       []string               `json:"codes" validate:"required,min=1,max=1800"`
	TotalAmount decimal.Decimal        `json:"total_amount" validate:"required"`
}

type WithdrawalTransferSubmitResponse struct {
	RequestID uint64 `json:"request_id"`
	AddressID uint32 `json:"address_id"`

	Status             meta.SystemWithdrawalRequestStatus `json:"status"`
	Currency           meta.Currency                      `json:"currency"`
	Network            meta.BlockchainNetwork             `json:"network"`
	Amount             decimal.Decimal                    `json:"amount"`
	AmountEstimatedFee meta.CurrencyAmount                `json:"amount_estimated_fee"`
	CombinedTxnHash    string                             `json:"combined_txn_hash"`

	CreateUID  meta.UID `json:"create_uid"`
	CreateTime int64    `json:"create_time"`
}

type WithdrawalTransferConfirmRequest struct {
	RequestID uint64 `json:"request_id" validate:"required"`
}

type WithdrawalTransferReplaceRequest struct {
	RequestID uint64                `json:"request_id" validate:"required"`
	FeeInfo   blockchainmod.FeeInfo `json:"fee_info" validate:"required"`
}
