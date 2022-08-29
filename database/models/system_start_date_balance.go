package models

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	SystemStartDateBalanceTableName = "system_start_date_balance"
)

type SystemStartDateBalance struct {
	ID         uint64          `gorm:"column:id;primaryKey"`
	Date       string          `gorm:"column:date"`
	Currency   meta.Currency   `gorm:"column:currency"`
	Amount     decimal.Decimal `gorm:"column:amount"`
	CreateTime int64           `gorm:"column:create_time;autoCreateTime"`
}

func (SystemStartDateBalance) TableName() string {
	return SystemStartDateBalanceTableName
}
