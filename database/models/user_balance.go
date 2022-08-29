package models

import (
	"database/sql"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	UserBalanceTableName = "user_balance"

	UserBalanceColCurrency   = "currency"
	UserBalanceColUID        = "uid"
	UserBalanceColAmount     = "amount"
	UserBalanceColUpdateTime = "update_time"
)

type UserBalance struct {
	ID         uint64          `gorm:"column:id;primaryKey"`
	Currency   meta.Currency   `gorm:"column:currency"`
	Amount     decimal.Decimal `gorm:"column:amount"`
	UpdateTime int64           `gorm:"column:update_time;autoUpdateTime"`

	UID         meta.UID      `gorm:"column:uid"`
	LatestTxnID sql.NullInt64 `gorm:"column:latest_txn_id"`
	// User   User     `gorm:"column:uid,foreignkey:UID"`
	// LatestTxn   *UserBalanceTxn `gorm:"column:latest_txn_id,foreignKey:LatestHistoryID"`
}

func (UserBalance) TableName() string {
	return UserBalanceTableName
}
