package v1

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type Channel struct {
	Type        meta.ChannelType `json:"type"`
	Name        string           `json:"name"`
	Description string           `json:"description"`

	MinTxnAmount decimal.Decimal `json:"min_txn_amount"`
	MaxTxnAmount decimal.Decimal `json:"max_txn_amount"`
}

type Order struct {
	ID        uint64           `json:"id"`
	Code      string           `json:"code"`
	UID       meta.UID         `json:"uid"`
	Direction meta.Direction   `json:"direction_type"`
	Currency  meta.Currency    `json:"currency"`
	Status    meta.OrderStatus `json:"status"`

	SrcChannelType    meta.ChannelType    `json:"src_channel_type"`
	SrcChannelID      uint64              `json:"src_channel_id"`
	SrcChannelRef     string              `json:"src_channel_ref"`
	SrcChannelAmount  decimal.Decimal     `json:"src_channel_amount"`
	SrcChannelContext OrderChannelContext `json:"src_channel_context,omitempty"`

	DstChannelType    meta.ChannelType    `json:"dst_channel_type"`
	DstChannelID      uint64              `json:"dst_channel_id"`
	DstChannelRef     string              `json:"dst_channel_ref"`
	DstChannelAmount  decimal.Decimal     `json:"dst_channel_amount"`
	DstChannelContext OrderChannelContext `json:"dst_channel_context,omitempty"`

	AmountSubTotal decimal.Decimal `json:"amount_sub_total"`
	AmountFee      decimal.Decimal `json:"amount_fee"`
	AmountDiscount decimal.Decimal `json:"amount_discount"`
	AmountTotal    decimal.Decimal `json:"amount_total"`

	Note       string `json:"note"`
	ExtraData  meta.O `json:"extra_data"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}

type OrderExportItem struct {
	Code     string          `csv:"Code"`
	Currency meta.Currency   `csv:"Currency"`
	Type     string          `csv:"Type"`
	Status   string          `csv:"Status"`
	Amount   decimal.Decimal `csv:"Amount"`
	Time     string          `csv:"Time"`
	Note     string          `csv:"Note"`
}

type OrderChannelContext struct {
	Meta    interface{} `json:"meta,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

type MetaCurrencyInfo struct {
	models.CurrencyInfo
	StatusMessage string `json:"status_message"`
}
