package models

import "gitlab.com/snap-clickstaff/torque-go/lib/meta"

const (
	UserTierChangeTableName = "user_tier_change"

	UserTierChangeColDate         = "date"
	UserTierChangeColUID          = "uid"
	UserTierChangeColStatus       = "status"
	UserTierChangeColApplyTime    = "apply_time"
	UserTierChangeColFromTierType = "from_tier_type"
	UserTierChangeColToTierType   = "to_tier_type"
	UserTierChangeColExtraData    = "extra_data"
	UserTierChangeColCreateTime   = "create_time"
)

var (
	userTierChangeCreateColumnNames = []string{
		UserTierChangeColDate,
		UserTierChangeColUID,
		UserTierChangeColStatus,
		UserTierChangeColApplyTime,
		UserTierChangeColFromTierType,
		UserTierChangeColToTierType,
		UserTierChangeColExtraData,
		UserTierChangeColCreateTime,
	}
)

type UserTierChange struct {
	ID           uint64   `gorm:"column:id;primaryKey"`
	Date         string   `gorm:"column:date"`
	UID          meta.UID `gorm:"column:uid"`
	Status       int      `gorm:"column:status"`
	ApplyTime    int64    `gorm:"column:apply_time"`
	FromTierType int      `gorm:"column:from_tier_type"`
	ToTierType   int      `gorm:"column:to_tier_type"`

	ExtraData  string `gorm:"column:extra_data"`
	CreateTime int64  `gorm:"column:create_time;autoCreateTime"`
}

func (UserTierChange) TableName() string {
	return UserTierChangeTableName
}

func (UserTierChange) GetCreateColumnNames() []string {
	return userTierChangeCreateColumnNames
}

func (m UserTierChange) GetCreateValues() []interface{} {
	return []interface{}{
		m.Date,
		m.UID,
		m.Status,
		m.ApplyTime,
		m.FromTierType,
		m.ToTierType,
		m.ExtraData,
		m.CreateTime,
	}
}
