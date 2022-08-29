package models

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	DailyProfitTableName = "daily_profit"

	DailyProfitColCoinID = "coin_id"
	DailyProfitColUID    = "user_id"
	DailyProfitColDate   = "date"
	DailyProfitColAmount = "dailyprofit"
)

type DailyProfit struct {
	ID        uint64          `gorm:"column:id;primaryKey"`
	CoinID    uint16          `gorm:"column:coin_id"`
	UID       meta.UID        `gorm:"column:user_id"`
	Date      string          `gorm:"column:date"`
	Amount    decimal.Decimal `gorm:"column:dailyprofit"`
	IsDeleted sql.NullBool    `gorm:"column:deleted"`

	CreateTime time.Time `gorm:"column:date_created"`
	UpdateTime time.Time `gorm:"column:date_modified"`
}

func (DailyProfit) TableName() string {
	return DailyProfitTableName
}
