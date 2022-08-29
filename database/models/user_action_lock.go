package models

import (
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	UserActionLockTableName = "user_action_lock"
)

type UserActionLock struct {
	ID         uint64   `gorm:"column:id;primaryKey"`
	UID        meta.UID `gorm:"column:uid"`
	ActionType uint16   `gorm:"column:action_type"`
	FromTime   int64    `gorm:"column:from_time"`
	ToTime     int64    `gorm:"column:to_time"`
	Note       string   `gorm:"column:note"`
	CreateTime int64    `gorm:"column:create_time;autoCreateTime"`
	UpdateTime int64    `gorm:"column:update_time;autoUpdateTime"`
}

func (UserActionLock) TableName() string {
	return UserActionLockTableName
}
