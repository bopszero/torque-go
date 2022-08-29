package models

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	NetworkCurrencyTableName = "network_currency"
)

type NetworkCurrency struct {
	ID uint16 `gorm:"column:id;primaryKey"`

	Currency      meta.Currency          `gorm:"column:currency"`
	Network       meta.BlockchainNetwork `gorm:"column:network"`
	Priority      uint8                  `gorm:"column:priority"`
	WithdrawalFee decimal.Decimal        `gorm:"column:withdrawal_fee"`

	CreateTime int64 `gorm:"column:create_time;autoCreateTime"`
	UpdateTime int64 `gorm:"column:update_time;autoUpdateTime"`
}

func (NetworkCurrency) TableName() string {
	return NetworkCurrencyTableName
}
