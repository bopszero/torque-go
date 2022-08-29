package models

import (
	"database/sql"

	"github.com/shopspring/decimal"
)

const (
	DailyPercentageTableName = "daily_percentage"
)

type DailyPercentage struct {
	ID           uint64          `gorm:"column:daily_percentage_id;primaryKey"`
	CoinID       uint16          `gorm:"column:coin_id"`
	Date         string          `gorm:"column:date"`
	Percentage   decimal.Decimal `gorm:"column:percentage"`
	RateToCoin   decimal.Decimal `gorm:"column:to_coin"`
	RateFromCoin decimal.Decimal `gorm:"column:from_coin"`
	IsDeleted    sql.NullBool    `gorm:"column:deleted"`
}

func (DailyPercentage) TableName() string {
	return DailyPercentageTableName
}
