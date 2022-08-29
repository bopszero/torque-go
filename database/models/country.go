package models

import "database/sql"

const (
	CountryTableName = "country"
)

type Country struct {
	ID            uint32       `gorm:"column:id;primaryKey"`
	CodeIso2      string       `gorm:"column:code_iso_2"`
	CodeIso3      string       `gorm:"column:code_iso_3"`
	Name          string       `gorm:"column:name"`
	NameFull      string       `gorm:"column:name_full"`
	PhoneCode     string       `gorm:"column:phone_code"`
	IsBanned      sql.NullBool `gorm:"column:is_banned"`
	FlagUrlSvg    string       `gorm:"column:flag_url_svg"`
	FlagUrlPng32  string       `gorm:"column:flag_url_png_32"`
	FlagUrlPng128 string       `gorm:"column:flag_url_png_128"`
	CreateTime    int64        `gorm:"column:create_time;autoCreateTime"`
	UpdateTime    int64        `gorm:"column:update_time;autoUpdateTime"`
}

func (Country) TableName() string {
	return CountryTableName
}
