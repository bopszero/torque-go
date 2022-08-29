package models

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/database/dbfields"
)

const (
	LeaderBonusPoolExecutionTableName   = "leader_bonus_pool_execution"
	LeaderBonusPoolExecutionColFromDate = "from_date"
	LeaderBonusPoolExecutionColToDate   = "to_date"
)

type LeaderBonusPoolExecution struct {
	ID          uint32             `gorm:"column:id;primaryKey" json:"id"`
	Hash        string             `gorm:"column:hash" json:"hash"`
	FromDate    dbfields.DateField `gorm:"column:from_date" json:"from_date"`
	ToDate      dbfields.DateField `gorm:"column:to_date" json:"to_date"`
	TotalAmount decimal.Decimal    `gorm:"column:total_amount" json:"total_amount"`
	Status      int8               `gorm:"column:status" json:"status"`

	TierSeniorRate            decimal.Decimal `gorm:"column:tier_senior_rate" json:"-"`
	TierSeniorReceiverCount   uint16          `gorm:"column:tier_senior_receiver_count" json:"-"`
	TierRegionalRate          decimal.Decimal `gorm:"column:tier_regional_rate" json:"-"`
	TierRegionalReceiverCount uint16          `gorm:"column:tier_regional_receiver_count" json:"-"`
	TierGlobalRate            decimal.Decimal `gorm:"column:tier_global_rate" json:"-"`
	TierGlobalReceiverCount   uint16          `gorm:"column:tier_global_receiver_count" json:"-"`

	CreateTime int64 `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime int64 `gorm:"column:update_time;autoUpdateTime" json:"-"`
}

func (LeaderBonusPoolExecution) TableName() string {
	return LeaderBonusPoolExecutionTableName
}
