package models

import (
	"database/sql"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	UserBalanceTxnTableName = "user_balance_txn"

	UserBalanceTxnColID         = "id"
	UserBalanceTxnColCurrency   = "currency"
	UserBalanceTxnColUID        = "uid"
	UserBalanceTxnColAmount     = "amount"
	UserBalanceTxnColBalance    = "balance"
	UserBalanceTxnColType       = "type"
	UserBalanceTxnColCreateTime = "create_time"
)

type UserBalanceTxn struct {
	ID         uint64          `gorm:"column:id;primaryKey"`
	UserID     meta.UID        `gorm:"column:uid"`
	Currency   meta.Currency   `gorm:"column:currency"`
	Amount     decimal.Decimal `gorm:"column:amount"`
	Balance    decimal.Decimal `gorm:"column:balance"`
	Type       uint32          `gorm:"column:type"`
	Ref        string          `gorm:"column:ref"`
	ParentID   sql.NullInt64   `gorm:"column:parent_id"`
	CreateTime int64           `gorm:"column:create_time;autoCreateTime"`

	// User   User     `gorm:"column:uid,foreignkey:UID"`
	// Parent   *UserBalanceTxn `gorm:"column:parent_id,foreignKey:ParentID"`
}

func (UserBalanceTxn) TableName() string {
	return UserBalanceTxnTableName
}
