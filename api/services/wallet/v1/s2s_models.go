package v1

import (
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type S2sPortfolioGetCurrencyRequest struct {
	PortfolioGetCurrencyRequest
	UID meta.UID `json:"uid" validate:"required"`
}

type S2sPaymentInitOrderRequest struct {
	PaymentInitOrderRequest
	UID meta.UID `json:"uid" validate:"required"`
}

type S2sPaymentInitOrderResponse PaymentInitOrderResponse

type S2sPaymentExecuteOrderRequest struct {
	UID     meta.UID `json:"uid" validate:"required"`
	OrderID uint64   `json:"order_id" validate:"required"`
}

type S2sPaymentExecuteOrderResponse PaymentExecuteOrderResponse

type S2sMetaHandshakeResponse MetaHandshakeResponse

type S2sGetUserActionLockRequest struct {
	UID        meta.UID `json:"uid" validate:"required"`
	ActionType uint16   `json:"action_type" validate:"required"`
}

type S2sGetUserActionLockResponse struct {
	IsLocked bool `json:"is_locked"`
}
