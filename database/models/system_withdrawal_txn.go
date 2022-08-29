package models

import (
	"database/sql"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	SystemWithdrawalTxnTableName = "system_withdrawal_txn"

	SystemWithdrawalTxnColRequestID   = "request_id"
	SystemWithdrawalTxnColStatus      = "status"
	SystemWithdrawalTxnColOutputIndex = "output_index"
)

type SystemWithdrawalTxn struct {
	ID        uint64 `gorm:"column:id;primaryKey"`
	RequestID uint64 `gorm:"column:request_id"`
	RefCode   string `gorm:"column:ref_code"`

	Currency       meta.Currency                  `gorm:"column:currency"`
	Status         meta.SystemWithdrawalTxnStatus `gorm:"column:status"`
	Hash           sql.NullString                 `gorm:"column:hash"`
	ToAddress      string                         `gorm:"column:to_address"`
	OutputIndex    uint64                         `gorm:"column:output_index"`
	FeeAmount      decimal.Decimal                `gorm:"column:fee_amount"`
	FeePrice       decimal.Decimal                `gorm:"column:fee_price"`
	FeeMaxQuantity uint32                         `gorm:"column:fee_max_quantity"`
	SignedBytes    sql.NullString                 `gorm:"column:signed_bytes"`

	CreateTime int64 `gorm:"column:create_time;autoCreateTime"`
	UpdateTime int64 `gorm:"column:update_time;autoUpdateTime"`
}

func (SystemWithdrawalTxn) TableName() string {
	return SystemWithdrawalTxnTableName
}
