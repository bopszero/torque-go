package models

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	WithdrawTableName = "withdraw"

	WithdrawColCode          = "code"
	WithdrawColStatus        = "status"
	WithdrawColExecuteStatus = "execute_status"
)

type Withdraw struct {
	ID        uint64                 `gorm:"column:withdraw_id;primaryKey"`
	Code      string                 `gorm:"column:code"`
	CoinID    uint16                 `gorm:"column:coin_id"`
	Currency  meta.Currency          `gorm:"column:currency"`
	Network   meta.BlockchainNetwork `gorm:"column:network"`
	UserID    meta.UID               `gorm:"column:user_id"`
	Status    string                 `gorm:"column:status"`
	Amount    decimal.Decimal        `gorm:"column:amount"`
	Fee       decimal.Decimal        `gorm:"column:fee"`
	Address   string                 `gorm:"column:address"`
	Note      string                 `gorm:"column:remark"`
	IsDeleted sql.NullBool           `gorm:"column:deleted"`

	TxnHash       string                         `gorm:"column:txn_hash"`
	ExecuteStatus meta.SystemWithdrawalTxnStatus `gorm:"column:execute_status"`

	CreateTime time.Time `gorm:"column:date_created"`
	UpdateTime int64     `gorm:"column:update_time;autoUpdateTime"`
	CloseTime  int64     `gorm:"column:close_time"`
}

func (Withdraw) TableName() string {
	return WithdrawTableName
}
