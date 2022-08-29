package models

import (
	"database/sql"

	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	KycSpecialUserInfoTableName = "kyc_special_user"
)

type KycSpecialUserInfo struct {
	ID         uint64                  `gorm:"column:id;primaryKey"`
	UID        meta.UID                `gorm:"column:uid"`
	Type       meta.KycSpecialUserType `gorm:"column:type"`
	IsPending  sql.NullBool            `gorm:"column:is_pending"`
	Note       string                  `gorm:"column:note"`
	CreateTime int64                   `gorm:"column:create_time;autoCreateTime"`
	UpdateTime int64                   `gorm:"column:update_time;autoUpdateTime"`
}

func (KycSpecialUserInfo) TableName() string {
	return KycSpecialUserInfoTableName
}
