package v1

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type PaymentGetChannelRequest struct {
	Type meta.ChannelType `json:"type" validate:"required"`
}

type PaymentGetChannelResponse Channel

type PaymentCheckoutOrderRequest struct {
	Currency meta.Currency `json:"currency" validate:"required"`

	SrcChannelType    meta.ChannelType `json:"src_channel_type" validate:"required"`
	SrcChannelID      uint64           `json:"src_channel_id"`
	SrcChannelRef     string           `json:"src_channel_ref"`
	SrcChannelAmount  decimal.Decimal  `json:"src_channel_amount"`
	SrcChannelContext meta.O           `json:"src_channel_context"`

	DstChannelType    meta.ChannelType `json:"dst_channel_type" validate:"required"`
	DstChannelID      uint64           `json:"dst_channel_id"`
	DstChannelRef     string           `json:"dst_channel_ref"`
	DstChannelAmount  decimal.Decimal  `json:"dst_channel_amount"`
	DstChannelContext meta.O           `json:"dst_channel_context"`
}

type PaymentCheckoutOrderResponse struct {
	// ChannelSrcMeta Channel `json:"channel_src_meta"`
	// ChannelDstMeta Channel `json:"channel_dst_meta"`

	ChannelSrcInfo interface{} `json:"channel_src_info"`
	ChannelDstInfo interface{} `json:"channel_dst_info"`
}

type PaymentInitOrderRequest struct {
	Currency meta.Currency `json:"currency" validate:"required"`

	SrcChannelType    meta.ChannelType `json:"src_channel_type" validate:"required"`
	SrcChannelID      uint64           `json:"src_channel_id"`
	SrcChannelRef     string           `json:"src_channel_ref"`
	SrcChannelAmount  decimal.Decimal  `json:"src_channel_amount" validate:"required"`
	SrcChannelContext meta.O           `json:"src_channel_context"`

	DstChannelType    meta.ChannelType `json:"dst_channel_type" validate:"required"`
	DstChannelID      uint64           `json:"dst_channel_id"`
	DstChannelRef     string           `json:"dst_channel_ref"`
	DstChannelAmount  decimal.Decimal  `json:"dst_channel_amount" validate:"required"`
	DstChannelContext meta.O           `json:"dst_channel_context"`

	AmountSubTotal decimal.Decimal `json:"amount_sub_total" validate:"required"`
	AmountTotal    decimal.Decimal `json:"amount_total" validate:"required"`
	Note           string          `json:"note"`
}

type PaymentInitOrderResponse struct {
	OrderID uint64 `json:"order_id"`
}

type PaymentExecuteOrderRequest struct {
	OrderID  uint64 `json:"order_id" validate:"required"`
	AuthCode string `json:"auth_code"`
}

type PaymentExecuteOrderResponse Order
