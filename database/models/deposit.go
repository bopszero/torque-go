package models

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	DepositTableName = "deposit"

	DepositColID            = "deposit_id"
	DepositColStatus        = "status"
	DepositColUserID        = "user_id"
	DepositColTorqueTxnID   = "torque_pair_id"
	DepositColForwardStatus = "forward_status"
	DepositColCreateTime    = "date_created"
	DepositColCloseTime     = "close_time"
)

type Deposit struct {
	ID        uint64                 `gorm:"column:deposit_id;primaryKey"`
	CoinID    uint16                 `gorm:"column:coin_id"`
	Currency  meta.Currency          `gorm:"column:currency"`
	Network   meta.BlockchainNetwork `gorm:"column:network"`
	UID       meta.UID               `gorm:"column:user_id"`
	Status    string                 `gorm:"column:status"`
	IsDeleted sql.NullBool           `gorm:"column:deleted"`

	TxnHash     string          `gorm:"column:txn_hash"`
	TxnIndex    uint16          `gorm:"column:txn_to_index"`
	Address     string          `gorm:"column:address"`
	Amount      decimal.Decimal `gorm:"column:amount"`
	Note        string          `gorm:"column:notes"`
	IsReinvest  sql.NullBool    `gorm:"column:is_reinvest"`
	TorqueTxnID sql.NullInt64   `gorm:"column:torque_pair_id"`

	ForwardStatus meta.SystemForwardingOrderTxnStatus `gorm:"column:forward_status"`

	CreateTime time.Time `gorm:"column:date_created"`
	UpdateTime int64     `gorm:"column:update_time;autoUpdateTime"`
	CloseTime  int64     `gorm:"column:close_time"`
}

func (Deposit) TableName() string {
	return DepositTableName
}
