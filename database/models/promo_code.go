package models

import (
	"database/sql"
	"time"
)

const (
	PromoCodeTableName = "promo_code"
)

type PromoCode struct {
	ID        uint32       `gorm:"column:id;primaryKey"`
	Code      string       `gorm:"column:code"`
	Status    string       `gorm:"column:status"`
	IsDeleted sql.NullBool `gorm:"column:deleted"`

	CreateDate time.Time `gorm:"column:date_created"`
	UpdateDate time.Time `gorm:"column:date_modified"`
}

func (PromoCode) TableName() string {
	return PromoCodeTableName
}
