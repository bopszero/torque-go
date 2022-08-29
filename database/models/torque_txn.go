package models

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	TorqueTxnTableName = "torque"

	TorqueTxnColID            = "torque_id"
	TorqueTxnColCode          = "code"
	TorqueTxnColStatus        = "status"
	TorqueTxnColExecuteStatus = "execute_status"
)

type TorqueTxn struct {
	ID        uint64                 `gorm:"column:torque_id;primaryKey"`
	Code      string                 `gorm:"column:code"`
	UserID    meta.UID               `gorm:"column:user_id"`
	CoinID    uint16                 `gorm:"column:coin_id"`
	Currency  meta.Currency          `gorm:"column:currency"`
	Network   meta.BlockchainNetwork `gorm:"column:network"`
	Status    string                 `gorm:"column:status"`
	Note      string                 `gorm:"column:remark"`
	IsDeleted sql.NullBool           `gorm:"column:deleted"`

	IsReinvest     sql.NullBool    `gorm:"column:is_reinvest"`
	Amount         decimal.Decimal `gorm:"column:amount"`
	Address        string          `gorm:"column:address"`
	ExchangeRate   decimal.Decimal `gorm:"column:rate"`
	CoinAmount     decimal.Decimal `gorm:"column:coin_amount"`
	CoinFee        decimal.Decimal `gorm:"column:coin_fee"`
	Balance        decimal.Decimal `gorm:"column:balance"`
	BlockchainHash string          `gorm:"column:transactionhash"`

	ExecuteStatus meta.SystemWithdrawalTxnStatus `gorm:"column:execute_status"`

	CreateTime time.Time `gorm:"column:date_created"`
	UpdateTime int64     `gorm:"column:update_time;autoUpdateTime"`
	CloseTime  int64     `gorm:"column:close_time"`
}

func (TorqueTxn) TableName() string {
	return TorqueTxnTableName
}
