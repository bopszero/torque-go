package models

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	P2PTransferTableName = "p2p_transfer"
	P2PTransferColStatus = "status"
)

type P2PTransfer struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	FromUID      meta.UID        `gorm:"column:from_uid"`
	FromCurrency meta.Currency   `gorm:"column:from_currency"`
	FromAmount   decimal.Decimal `gorm:"column:from_amount"`
	ToUID        meta.UID        `gorm:"column:to_uid"`
	ToCurrency   meta.Currency   `gorm:"column:to_currency"`
	ToAmount     decimal.Decimal `gorm:"column:to_amount"`
	ExchangeRate decimal.Decimal `gorm:"column:exchange_rate"`

	FeeAmount decimal.Decimal `gorm:"column:fee_amount"`
	Note      string          `gorm:"column:note"`
	Status    int8            `gorm:"column:status"`

	CreateTime int64 `gorm:"column:create_time;autoCreateTime"`
	UpdateTime int64 `gorm:"column:update_time;autoUpdateTime"`
}

func (P2PTransfer) TableName() string {
	return P2PTransferTableName
}
