package models

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	LeaderRewardTableName = "user_reward"

	LeaderRewardColUID    = "leader_id"
	LeaderRewardColDate   = "date"
	LeaderRewardColAmount = "amount"
)

type LeaderReward struct {
	ID        uint64          `gorm:"column:id;primaryKey"`
	UID       meta.UID        `gorm:"column:leader_id"`
	Date      string          `gorm:"column:date"`
	Amount    decimal.Decimal `gorm:"column:amount"`
	IsDeleted sql.NullBool    `gorm:"column:deleted"`

	CreateTime time.Time `gorm:"column:date_created"`
}

func (LeaderReward) TableName() string {
	return LeaderRewardTableName
}
