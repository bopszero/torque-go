package sysforwardingmod

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	OrderMaxTxnLimit = 100
)

const (
	ImportIndexNetwork = 0
	ImportIndexAddress = 1
	ImportIndexKey     = 2

	ImportColNetwork = "network"
	ImportColAddress = "address"
	ImportColKey     = "key"
)

type ForwardingHandler interface {
	GenerateOrder(ctx comcontext.Context) error
	SignOrder(ctx comcontext.Context, orderID uint64) (*models.SystemForwardingOrder, error)
	ExecuteOrder(ctx comcontext.Context, orderID uint64) (*models.SystemForwardingOrder, error)
}

type ForwardConfig struct {
	ApiProvider          string          `json:"api_provider"`
	Address              string          `json:"address"`
	FeeKey               string          `json:"fee_key"`
	TxnCountMinThreshold int             `json:"txn_count_min_threshold"`
	TxnCountMaxThreshold int             `json:"txn_count_max_threshold"`
	AmountMinThreshold   decimal.Decimal `json:"amount_min_threshold"`
}

type ReportItem struct {
	ID              uint64          `csv:"id"`
	DepositID       uint64          `csv:"deposit_id"`
	DepositAmount   decimal.Decimal `csv:"deposit_amount"`
	Currency        meta.Currency   `csv:"currency"`
	Address         string          `csv:"address"`
	Status          string          `csv:"status"`
	TxnHash         string          `csv:"txn_hash"`
	GrossAmount     decimal.Decimal `csv:"gross_amount"`
	FeeAmount       decimal.Decimal `csv:"fee_amount"`
	NetAmount       decimal.Decimal `csv:"net_amount"`
	SignBalance     decimal.Decimal `csv:"sign_balance"`
	SignBalanceTime string          `csv:"sign_balance_time"`
	Note            string          `csv:"note"`
}
