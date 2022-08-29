package models

import (
	"database/sql"

	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	CryptoTxnTableName = "crypto_txn"
)

type CryptoTxn struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	UID         meta.UID               `gorm:"column:uid"`
	Network     meta.BlockchainNetwork `gorm:"column:network"`
	Hash        string                 `gorm:"column:hash"`
	SignedBytes sql.NullString         `gorm:"column:signed_bytes"`

	CreateTime int64 `gorm:"column:create_time;autoCreateTime"`
	UpdateTime int64 `gorm:"column:update_time;autoUpdateTime"`
}

func (CryptoTxn) TableName() string {
	return CryptoTxnTableName
}
