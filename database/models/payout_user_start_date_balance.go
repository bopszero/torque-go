package models

import (
	"database/sql"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	PayoutUserStartDateBalanceTableName = "user_start_date_balance"

	PayoutUserStartDateBalanceColCurrency = "currency"
	PayoutUserStartDateBalanceColUID      = "uid"
)

type PayoutUserStartDateBalance struct {
	ID         uint64          `gorm:"column:id;primaryKey"`
	Date       string          `gorm:"column:date"`
	Currency   meta.Currency   `gorm:"column:currency"`
	UID        meta.UID        `gorm:"column:uid"`
	Amount     decimal.Decimal `gorm:"column:amount"`
	TxnID      sql.NullInt64   `gorm:"column:txn_id"`
	CreateTime int64           `gorm:"column:create_time;autoCreateTime"`
}

func (PayoutUserStartDateBalance) TableName() string {
	return PayoutUserStartDateBalanceTableName
}
