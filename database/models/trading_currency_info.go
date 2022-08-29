package models

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	TradingCurrencyInfoTableName = "coins"

	TradingCurrencyInfoColLatestCrawledBlockHeight = "latest_crawled_block_height"
)

type LegacyCurrencyInfo struct {
	ID       uint16        `gorm:"column:coin_id;primaryKey"`
	Currency meta.Currency `gorm:"column:coin_name"`

	DepositMinThreshold       decimal.Decimal `gorm:"column:min_deposit_display"`
	DepositMinThresholdPayout decimal.Decimal `gorm:"column:min_deposit"`
	LatestCrawledBlockHeight  uint64          `gorm:"column:latest_crawled_block_height"`
}

func (LegacyCurrencyInfo) TableName() string {
	return TradingCurrencyInfoTableName
}
