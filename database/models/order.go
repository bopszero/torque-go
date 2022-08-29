package models

import (
	"database/sql/driver"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/database/dbfields"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	OrderTableName = "order"

	OrderColID             = "id"
	OrderColCurrency       = "currency"
	OrderColStatus         = "status"
	OrderColSrcChannelType = "src_channel_type"
	OrderColDstChannelType = "dst_channel_type"
	OrderColDstChannelID   = "dst_channel_id"
	OrderColDstChannelRef  = "dst_channel_ref"
	OrderColCreateTime     = "create_time"
	OrderColRetryTime      = "retry_time"
)

type Order struct {
	ID        uint64           `gorm:"column:id;primaryKey"`
	Code      string           `gorm:"column:code"`
	UID       meta.UID         `gorm:"column:uid"`
	Direction meta.Direction   `gorm:"column:direction_type"`
	Currency  meta.Currency    `gorm:"column:currency"`
	Status    meta.OrderStatus `gorm:"column:status"`

	SrcChannelType   meta.ChannelType `gorm:"column:src_channel_type"`
	SrcChannelID     uint64           `gorm:"column:src_channel_id"`
	SrcChannelRef    string           `gorm:"column:src_channel_ref"`
	SrcChannelAmount decimal.Decimal  `gorm:"column:src_channel_amount"`

	DstChannelType   meta.ChannelType `gorm:"column:dst_channel_type"`
	DstChannelID     uint64           `gorm:"column:dst_channel_id"`
	DstChannelRef    string           `gorm:"column:dst_channel_ref"`
	DstChannelAmount decimal.Decimal  `gorm:"column:dst_channel_amount"`

	AmountSubTotal decimal.Decimal `gorm:"column:amount_sub_total"`
	AmountFee      decimal.Decimal `gorm:"column:amount_fee"`
	AmountDiscount decimal.Decimal `gorm:"column:amount_discount"`
	AmountTotal    decimal.Decimal `gorm:"column:amount_total"`

	Note      string             `gorm:"column:note"`
	StepsData OrderStepsData     `gorm:"column:steps_data"`
	ExtraData dbfields.JsonField `gorm:"column:extra_data"`

	SucceedTime int64 `gorm:"column:succeed_time"`
	RetryTime   int64 `gorm:"column:retry_time"`
	CreateTime  int64 `gorm:"column:create_time;autoCreateTime"`
	UpdateTime  int64 `gorm:"column:update_time;autoUpdateTime"`
}

func (Order) TableName() string {
	return OrderTableName
}

type OrderStep struct {
	Direction meta.Direction `json:"d"`
	Code      string         `json:"c"`
	Time      int64          `json:"t"`
	Error     string         `json:"e,omitempty"`
}

type OrderStepsData struct {
	Current *OrderStep  `json:"current"`
	History []OrderStep `json:"history"`
}

func (osd OrderStepsData) Value() (driver.Value, error) {
	return toJson(osd)
}

func (osd *OrderStepsData) Scan(input interface{}) error {
	return fromJson(input, osd)
}
