package models

import (
	"database/sql"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	WalletBalanceTxnTableName = "user_balance_txn"

	WalletBalanceTxnColID         = "id"
	WalletBalanceTxnColCurrency   = "currency"
	WalletBalanceTxnColUID        = "uid"
	WalletBalanceTxnColBalance    = "balance"
	WalletBalanceTxnColType       = "type"
	WalletBalanceTxnColCreateTime = "create_time"
)

type WalletBalanceTxn struct {
	ID         uint64          `gorm:"column:id;primaryKey"`
	Currency   meta.Currency   `gorm:"column:currency"`
	UID        meta.UID        `gorm:"column:uid"`
	Amount     decimal.Decimal `gorm:"column:amount"`
	Balance    decimal.Decimal `gorm:"column:balance"`
	Type       uint32          `gorm:"column:type"`
	OrderID    uint64          `gorm:"column:order_id"`
	ParentID   sql.NullInt64   `gorm:"column:parent_id"`
	CreateTime int64           `gorm:"column:create_time;autoCreateTime"`
}

func (WalletBalanceTxn) TableName() string {
	return WalletBalanceTxnTableName
}
