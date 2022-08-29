package models

import "gitlab.com/snap-clickstaff/torque-go/lib/meta"

const UserConfigTableName = "user_config"

type UserConfig struct {
	UID        meta.UID        `gorm:"column:uid;primaryKey"`
	RewardType meta.RewardType `gorm:"column:reward_type"`
}

func (UserConfig) TableName() string {
	return UserConfigTableName
}
