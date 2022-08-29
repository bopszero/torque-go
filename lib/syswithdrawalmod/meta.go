package syswithdrawalmod

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	LockDurationSubmit = 60 * time.Second
	InitTimeout        = 5 * time.Minute
)

var (
	// Round up to 600 from 546 for 1000 relay fee
	// Reference: https://github.com/btcsuite/btcd/blob/56cc42fe07c206e76812fc57a216b59c41189f04/mempool/policy.go#L180
	UtxoLikeDustAmountThreshold = comutils.NewDecimalF("0.000006")

	// 1000 Gwei
	FeePriceMaxEthereum = comutils.NewDecimalF("0.000001")

	ConfirmAcceptStatuses = []meta.SystemWithdrawalTxnStatus{
		constants.SystemWithdrawalTxnStatusFailed,
		constants.SystemWithdrawalTxnStatusNew,
		constants.SystemWithdrawalTxnStatusCancelled,
	}
)

type Transfer struct {
	RefCode   string
	Currency  meta.Currency
	ToAddress string
	Amount    decimal.Decimal
}

type WithdrawalConfig struct {
	PullAddress string `json:"pull_address"`
}
