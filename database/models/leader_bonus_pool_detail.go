package models

import (
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	LeaderBonusPoolDetailTableName = "leader_bonus_pool_detail"

	LeaderBonusPoolDetailColExecutionID = "exec_id"
	LeaderBonusPoolDetailColOrderID = "order_id"
)

type LeaderBonusPoolDetail struct {
	ID          uint64   `gorm:"column:id;primaryKey"`
	ExecutionID uint32   `gorm:"column:exec_id"`
	TierType    int      `gorm:"column:tier_type"`
	UID         meta.UID `gorm:"column:uid"`
	OrderID     uint64   `gorm:"column:order_id"`
	Note        string   `gorm:"column:note"`
	CreateTime  int64    `gorm:"column:create_time;autoCreateTime"`
	UpdateTime  int64    `gorm:"column:update_time;autoUpdateTime"`
}

func (LeaderBonusPoolDetail) TableName() string {
	return LeaderBonusPoolDetailTableName
}
