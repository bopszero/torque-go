package models

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	AffiliateCommissionTableName = "affiliate_com"

	AffiliateCommissionColUID    = "user_id"
	AffiliateCommissionColDate   = "date"
	AffiliateCommissionColAmount = "amount"
)

type AffiliateCommission struct {
	ID        uint64          `gorm:"column:id;primaryKey"`
	CoinID    uint16          `gorm:"column:coin_id"`
	UID       meta.UID        `gorm:"column:user_id"`
	Date      string          `gorm:"column:date"`
	Amount    decimal.Decimal `gorm:"column:amount"`
	IsDeleted sql.NullBool    `gorm:"column:deleted"`

	CreateTime time.Time `gorm:"column:created_date"`
	UpdateTime time.Time `gorm:"column:last_modified"`
}

func (AffiliateCommission) TableName() string {
	return AffiliateCommissionTableName
}
