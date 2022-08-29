package payment

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type DepositGetAccountRequest struct {
	UID      meta.UID               `json:"uid" validate:"required"`
	Currency meta.Currency          `json:"currency" validate:"required"`
	Network  meta.BlockchainNetwork `json:"network"`
}

type DepositGetAccountResponse struct {
	UID           meta.UID               `json:"uid"`
	Currency      meta.Currency          `json:"currency"`
	Network       meta.BlockchainNetwork `json:"network"`
	Address       string                 `json:"address"`
	AddressLegacy string                 `json:"address_legacy"`
	CreateTime    int64                  `json:"create_time"`
}

type DepositCrawlRequest struct {
	Currency    meta.Currency          `json:"currency" validate:"required"`
	Network     meta.BlockchainNetwork `json:"network"`
	BlockHeight uint64                 `json:"block_height"`
}

type DepositSubmitRequest struct {
	UID      meta.UID               `json:"uid" validate:"required"`
	Currency meta.Currency          `json:"currency" validate:"required"`
	Network  meta.BlockchainNetwork `json:"network"`

	TxnHash  string          `json:"txn_hash" validate:"required"`
	TxnIndex uint16          `json:"txn_index"`
	Address  string          `json:"address" validate:"required"`
	Amount   decimal.Decimal `json:"amount" validate:"required"`
}

type DepositSubmitResponse struct {
	ID       uint64                 `json:"id"`
	CoinID   uint16                 `json:"coin_id"`
	Currency meta.Currency          `json:"currency"`
	Network  meta.BlockchainNetwork `json:"network"`
	UID      meta.UID               `json:"uid"`
	Status   string                 `json:"status"`

	TxnHash  string          `json:"txn_hash"`
	TxnIndex uint16          `json:"txn_index"`
	Address  string          `json:"address"`
	Amount   decimal.Decimal `json:"amount"`

	CreateTime int64 `json:"create_time"`
}

type DepositApproveRequest struct {
	ID   uint64 `json:"id" validate:"required"`
	Note string `json:"note" validate:"max=255"`
}

type InvestmentWithdrawSubmitRequest struct {
	UID      meta.UID               `json:"uid" validate:"required"`
	Amount   decimal.Decimal        `json:"amount" validate:"required"`
	Currency meta.Currency          `json:"currency" validate:"required"`
	Network  meta.BlockchainNetwork `json:"network"`
	Address  string                 `json:"address" validate:"required"`
}

type InvestmentWithdrawSubmitResponse struct {
	ID         uint64                 `json:"id"`
	Code       string                 `json:"code"`
	UID        meta.UID               `json:"uid"`
	Amount     decimal.Decimal        `json:"amount"`
	Currency   meta.Currency          `json:"currency"`
	Network    meta.BlockchainNetwork `json:"network"`
	Address    string                 `json:"address"`
	Status     string                 `json:"status"`
	CreateTime int64                  `json:"create_time"`
}

type InvestmentWithdrawRejectRequest struct {
	ID   uint64 `json:"id" validate:"required"`
	Note string `json:"note"`
}

type InvestmentWithdrawCancelRequest InvestmentWithdrawRejectRequest

type ProfitWithdrawSubmitRequest struct {
	UID          meta.UID        `json:"uid" validate:"required"`
	Amount       decimal.Decimal `json:"amount" validate:"required"`
	ExchangeRate decimal.Decimal `json:"exchange_rate" validate:"required"`
	Currency     meta.Currency   `json:"currency" validate:"required"`
	Address      string          `json:"address" validate:"required"`
}

type ProfitWithdrawSubmitResponse struct {
	ID             uint64          `json:"id"`
	Code           string          `json:"code"`
	UID            meta.UID        `json:"uid"`
	Amount         decimal.Decimal `json:"amount"`
	ExchangeRate   decimal.Decimal `json:"exchange_rate"`
	CurrencyAmount decimal.Decimal `json:"currency_amount"`
	Currency       meta.Currency   `json:"currency"`
	Address        string          `json:"address"`
	Status         string          `json:"status"`
	CreateTime     int64           `json:"create_time"`
}

type ProfitWithdrawRejectRequest struct {
	ID   uint64 `json:"id" validate:"required"`
	Note string `json:"note"`
}

type ProfitWithdrawCancelRequest ProfitWithdrawRejectRequest

type ProfitReinvestSubmitRequest struct {
	UID          meta.UID        `json:"uid" validate:"required"`
	Amount       decimal.Decimal `json:"amount" validate:"required"`
	ExchangeRate decimal.Decimal `json:"exchange_rate" validate:"required"`
	Currency     meta.Currency   `json:"currency" validate:"required"`
	Address      string          `json:"address" validate:"required"`
}

type ProfitReinvestSubmitResponse struct {
	ID             uint64          `json:"id"`
	UID            meta.UID        `json:"uid"`
	Amount         decimal.Decimal `json:"amount"`
	ExchangeRate   decimal.Decimal `json:"exchange_rate"`
	CurrencyAmount decimal.Decimal `json:"currency_amount"`
	Currency       meta.Currency   `json:"currency"`
	Address        string          `json:"address"`
	Status         string          `json:"status"`
	CreateTime     int64           `json:"create_time"`
}
type ProfitReinvestApproveRequest struct {
	ID   uint64 `json:"id" validate:"required"`
	Note string `json:"note"`
}

type ProfitReinvestRejectRequest struct {
	ID   uint64 `json:"id" validate:"required"`
	Note string `json:"note"`
}

type P2PTransferRequest struct {
	FromUID      meta.UID        `json:"from_uid" validate:"required"`
	FromCurrency meta.Currency   `json:"from_currency" validate:"required"`
	FromAmount   decimal.Decimal `json:"from_amount" validate:"required"`
	ToUID        meta.UID        `json:"to_uid" validate:"required"`
	ToCurrency   meta.Currency   `json:"to_currency" validate:"required"`
	ExchangeRate decimal.Decimal `json:"exchange_rate" validate:"required"`
	FeeAmount    decimal.Decimal `json:"fee_amount"`
	Note         string          `json:"note"`
}

type P2PTransferResponse struct {
	ID           uint64          `json:"id"`
	FromUID      meta.UID        `json:"from_uid"`
	FromCurrency meta.Currency   `json:"from_currency"`
	FromAmount   decimal.Decimal `json:"from_amount"`
	ToUID        meta.UID        `json:"to_uid"`
	ToCurrency   meta.Currency   `json:"to_currency"`
	ToAmount     decimal.Decimal `json:"to_amount"`
	ExchangeRate decimal.Decimal `json:"exchange_rate"`
	FeeAmount    decimal.Decimal `json:"fee_amount"`
	Note         string          `json:"note"`
	Status       int8            `json:"status"`
	CreateTime   int64           `json:"create_time"`
}

type PromoCodeRedeemRequest struct {
	UID  meta.UID `json:"uid"`
	Code string   `json:"code"`
}

type PromoCodeRedeemResponse struct {
	ID         uint64   `json:"id"`
	UID        meta.UID `json:"uid"`
	Code       string   `json:"code"`
	CreateTime int64    `json:"create_time"`
}
