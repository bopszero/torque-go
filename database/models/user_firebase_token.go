package models

import (
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	UserFirebaseTokenTableName = "user_firebase_token"

	UserFirebaseTokenColUID = "uid"
)

type UserFirebaseToken struct {
	ID         uint64   `gorm:"column:id;primaryKey"`
	UID        meta.UID `gorm:"column:uid"`
	Token      string   `gorm:"column:token"`
	DeviceOS   string   `gorm:"column:device_os"`
	DeviceDesc string   `gorm:"column:device_desc"`
	CreateTime int64    `gorm:"column:create_time;autoCreateTime"`
}

func (UserFirebaseToken) TableName() string {
	return UserFirebaseTokenTableName
}
