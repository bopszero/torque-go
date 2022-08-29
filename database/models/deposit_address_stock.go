package models

import (
	"database/sql"
	"time"

	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	DepositAddressStockTableName = "deposit_address_stock"
)

type DepositAddressStock struct {
	ID         uint64                 `gorm:"column:id;primaryKey"`
	CoinID     uint16                 `gorm:"column:coin_id"`
	Currency   meta.Currency          `gorm:"column:currency"`
	Network    meta.BlockchainNetwork `gorm:"column:network"`
	Address    string                 `gorm:"column:address"`
	IsUsed     sql.NullBool           `gorm:"column:is_used"`
	CreateDate time.Time              `gorm:"column:create_date"`
}

func (DepositAddressStock) TableName() string {
	return DepositAddressStockTableName
}
