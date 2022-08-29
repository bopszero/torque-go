package models

import (
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	SettingTableName = "setting"
)

type Setting struct {
	ID    meta.UID `gorm:"column:setting_id;primaryKey"`
	Key   string   `gorm:"column:field"`
	Value string   `gorm:"column:value"`
}

func (Setting) TableName() string {
	return SettingTableName
}
