package models

import (
	"github.com/shopspring/decimal"
)

const (
	PayoutStatsTableName = "payout_stats"
)

type PayoutStats struct {
	ID   uint64 `gorm:"column:id;primaryKey"`
	Date string `gorm:"column:date"`

	TotalProfitAmount           decimal.Decimal `gorm:"column:total_profit_amount"`
	RemainingAmountAffiliate    decimal.Decimal `gorm:"column:remaining_affiliate_amount"`
	RemainingAmountLeaderReward decimal.Decimal `gorm:"column:remaining_leader_reward_amount"`

	CreateTime int64 `gorm:"column:create_time;autoCreateTime"`
}

func (PayoutStats) TableName() string {
	return PayoutStatsTableName
}
