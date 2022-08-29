package models

import (
	"database/sql"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	SystemForwardingOrderTxnTableName = "forwarding_order_txn"

	SystemForwardingOrderTxnColOrderID = "order_id"
	SystemForwardingOrderTxnColStatus  = "status"
)

type SystemForwardingOrderTxn struct {
	ID        uint64 `gorm:"column:id;primaryKey"`
	OrderID   uint64 `gorm:"column:order_id"`
	DepositID uint64 `gorm:"column:deposit_id"`

	Currency          meta.Currency                       `gorm:"column:currency"`
	Status            meta.SystemForwardingOrderTxnStatus `gorm:"column:status"`
	FromAddress       string                              `gorm:"column:from_address"`
	Amount            decimal.Decimal                     `gorm:"column:amount"`
	Fee               decimal.Decimal                     `gorm:"column:fee"`
	Note              string                              `gorm:"column:note"`
	Hash              sql.NullString                      `gorm:"column:hash"`
	SignedBytes       sql.NullString                      `gorm:"column:signed_bytes"`
	SignBalance       decimal.Decimal                     `gorm:"column:sign_balance"`
	SignBalanceTime   int64                               `gorm:"column:sign_balance_time"`
	FeeTxnIndex       uint64                              `gorm:"column:fee_txn_index"`
	FeeTxnHash        sql.NullString                      `gorm:"column:fee_txn_hash"`
	FeeTxnSignedBytes sql.NullString                      `gorm:"column:fee_txn_signed_bytes"`

	CreateTime int64 `gorm:"column:create_time;autoCreateTime"`
	UpdateTime int64 `gorm:"column:update_time;autoUpdateTime"`
}

func (SystemForwardingOrderTxn) TableName() string {
	return SystemForwardingOrderTxnTableName
}
