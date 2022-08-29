package models

import (
	"database/sql"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	SystemWithdrawalRequestTableName = "system_withdrawal_request"

	SystemWithdrawalRequestColStatus     = "status"
	SystemWithdrawalRequestColCreateTime = "create_time"
)

type SystemWithdrawalRequest struct {
	ID        uint64 `gorm:"column:id;primaryKey"`
	AddressID uint32 `gorm:"column:address_id"`

	Currency            meta.Currency                      `gorm:"column:currency"`
	Network             meta.BlockchainNetwork             `gorm:"column:network"`
	Status              meta.SystemWithdrawalRequestStatus `gorm:"column:status"`
	Amount              decimal.Decimal                    `gorm:"column:amount"`
	AmountEstimatedFee  decimal.Decimal                    `gorm:"column:amount_estimated_fee"`
	CombinedTxnHash     sql.NullString                     `gorm:"column:combined_txn_hash"`
	CombinedSignedBytes sql.NullString                     `gorm:"column:combined_signed_bytes"`

	CreateUID  meta.UID `gorm:"column:create_uid"`
	CreateTime int64    `gorm:"column:create_time;autoCreateTime"`
	UpdateTime int64    `gorm:"column:update_time;autoUpdateTime"`
}

func (SystemWithdrawalRequest) TableName() string {
	return SystemWithdrawalRequestTableName
}
