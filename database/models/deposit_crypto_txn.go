package models

import (
	"database/sql"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	DepositCryptoTxnTableName = "deposit_crypto_txn"

	DepositCryptoTxnColConfirmations = "confirmations"
	DepositCryptoTxnColBlockHeight   = "block_height"
	DepositCryptoTxnColBlockTime     = "block_time"
)

type DepositCryptoTxn struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	Currency      meta.Currency          `gorm:"column:currency"`
	Network       meta.BlockchainNetwork `gorm:"column:network"`
	FromAddress   string                 `gorm:"column:from_address"`
	ToAddress     string                 `gorm:"column:to_address"`
	ToIndex       uint16                 `gorm:"column:to_index"`
	Hash          string                 `gorm:"column:hash"`
	Amount        decimal.Decimal        `gorm:"column:amount"`
	Confirmations uint64                 `gorm:"column:confirmations"`
	BlockHeight   uint64                 `gorm:"column:block_height"`
	BlockTime     int64                  `gorm:"column:block_time"`
	IsAccepted    sql.NullBool           `gorm:"column:is_accepted"`

	CreateTime int64 `gorm:"column:create_time;autoCreateTime"`
	UpdateTime int64 `gorm:"column:update_time;autoUpdateTime"`
}

func (DepositCryptoTxn) TableName() string {
	return DepositCryptoTxnTableName
}
