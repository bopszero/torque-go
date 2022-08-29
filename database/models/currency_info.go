package models

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/database/dbfields"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	CurrencyTableName = "currency"
)

type CurrencyInfo struct {
	ID uint16 `gorm:"column:id;primaryKey" json:"-"`

	Currency        meta.Currency      `gorm:"column:code" json:"currency"`
	PriceUSD        decimal.Decimal    `gorm:"column:price_usd" json:"price_usd"`
	PriceUSDT       decimal.Decimal    `gorm:"column:price_usdt" json:"price_usdt"`
	PriorityDisplay uint16             `gorm:"column:priority_display" json:"priority_display"`
	PriorityTrading uint16             `gorm:"column:priority_trading" json:"priority_trading"`
	PriorityWallet  uint16             `gorm:"column:priority_wallet" json:"priority_wallet"`
	IsFiat          dbfields.BoolField `gorm:"column:is_fiat" json:"is_fiat"`
	DecimalPlaces   uint8              `gorm:"column:decimal_places" json:"decimal_places"`
	Symbol          string             `gorm:"column:symbol" json:"symbol"`
	IconURL         string             `gorm:"column:icon_url" json:"icon_url"`
	BannerURL       string             `gorm:"column:banner_url" json:"banner_url"`
	ColorHex        string             `gorm:"column:color_hex" json:"color_hex"`

	UpdateTime int64 `gorm:"column:update_time;autoUpdateTime" json:"-"`
}

func (CurrencyInfo) TableName() string {
	return CurrencyTableName
}
