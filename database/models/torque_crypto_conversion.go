package models

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/database/dbfields"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	TorqueCryptoConversionTableName = "torque_crypto_conversion"
)

type TorqueCryptoConversion struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	UID            meta.UID           `gorm:"column:uid"`
	Currency       meta.Currency      `gorm:"column:currency"`
	CurrencyAmount decimal.Decimal    `gorm:"column:currency_amount"`
	ExchangeRate   decimal.Decimal    `gorm:"column:exchange_rate"`
	TorqueAmount   decimal.Decimal    `gorm:"column:torque_amount"`
	TxnHash        string             `gorm:"column:txn_hash"`
	SendOrderID    uint64             `gorm:"column:send_order_id"`
	ReceiveOrderID uint64             `gorm:"column:receive_order_id"`
	ExchangeRef    string             `gorm:"column:exchange_ref"`
	ProfitUSDT     decimal.Decimal    `gorm:"column:profit_usdt"`
	ExtraData      dbfields.JsonField `gorm:"column:extra_data"`

	CreateTime int64 `gorm:"column:create_time;autoCreateTime"`
	UpdateTime int64 `gorm:"column:update_time;autoUpdateTime"`
}

func (TorqueCryptoConversion) TableName() string {
	return TorqueCryptoConversionTableName
}
