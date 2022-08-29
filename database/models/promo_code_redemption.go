package models

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	PromoCodeRedemptionTableName = "redeem_promo"
)

type PromoCodeRedemption struct {
	ID          uint64          `gorm:"column:id;primaryKey"`
	PromoCodeID uint32          `gorm:"column:promo_code_id"`
	UID         meta.UID        `gorm:"column:user_id"`
	CreditValue decimal.Decimal `gorm:"column:torque_credit"`
	IsDeleted   sql.NullBool    `gorm:"column:deleted"`

	CreateDate time.Time `gorm:"column:date_created"`
	UpdateDate time.Time `gorm:"column:date_modified"`
}

func (PromoCodeRedemption) TableName() string {
	return PromoCodeRedemptionTableName
}
