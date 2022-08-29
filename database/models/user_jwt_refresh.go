package models

import "gitlab.com/snap-clickstaff/torque-go/lib/meta"

const (
	UserJwtRefreshTableName = "user_jwt_refresh"
)

type UserJwtRefresh struct {
	ID        string   `gorm:"PRIMARY_KEY;Column:id"`
	UID       meta.UID `gorm:"Column:uid"`
	DeviceUID string   `gorm:"Column:device_uid"`

	ExpireTime int64 `gorm:"Column:expire_time"`
	RotateTime int64 `gorm:"Column:rotate_time"`
	CreateTime int64 `gorm:"Column:create_time;AutoCreateTime"`
}

func (UserJwtRefresh) TableName() string {
	return UserJwtRefreshTableName
}
