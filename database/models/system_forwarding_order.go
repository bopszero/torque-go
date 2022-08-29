package models

import (
	"database/sql"

	"gitlab.com/snap-clickstaff/torque-go/database/dbfields"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	ForwardingOrderTableName = "forwarding_order"

	ForwardingOrderColStatus     = "status"
	ForwardingOrderColCreateTime = "create_time"
)

type SystemForwardingOrder struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	Date                dbfields.DateField               `gorm:"column:date"`
	Currency            meta.Currency                    `gorm:"column:currency"`
	Network             meta.BlockchainNetwork           `gorm:"column:network"`
	Address             string                           `gorm:"column:address"`
	Status              meta.SystemForwardingOrderStatus `gorm:"column:status"`
	CombinedTxnHash     sql.NullString                   `gorm:"column:combined_txn_hash"`
	CombinedSignedBytes sql.NullString                   `gorm:"column:combined_signed_bytes"`

	CreateTime int64 `gorm:"column:create_time;autoCreateTime"`
	UpdateTime int64 `gorm:"column:update_time;autoUpdateTime"`
}

func (SystemForwardingOrder) TableName() string {
	return ForwardingOrderTableName
}
