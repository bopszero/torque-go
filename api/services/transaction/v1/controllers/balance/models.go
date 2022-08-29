package balance

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type GetBalanceRequest struct {
	UID      meta.UID      `json:"uid" validate:"required"`
	Currency meta.Currency `json:"currency"`
}

type GetBalanceResponse struct {
	UID        meta.UID        `json:"uid"`
	Currency   meta.Currency   `json:"currency"`
	Amount     decimal.Decimal `json:"amount"`
	UpdateTime int64           `json:"update_time"`
}

type GetStartingBalanceRequest struct {
	UID          meta.UID      `json:"uid" validate:"required"`
	Currency     meta.Currency `json:"currency" validate:"required"`
	StartingTime int64         `json:"starting_time" validate:"required"`
}

type GetStartingBalanceResponse struct {
	UID          meta.UID        `json:"uid"`
	Currency     meta.Currency   `json:"currency"`
	Amount       decimal.Decimal `json:"amount"`
	StartingTime int64           `json:"starting_time"`
}

type AddTxnRequest struct {
	Currency meta.Currency   `json:"currency" validate:"required"`
	UID      meta.UID        `json:"uid" validate:"required"`
	Amount   decimal.Decimal `json:"amount" validate:"required"`
	TypeCode string          `json:"type_code" validate:"required"`
	Ref      string          `json:"ref" validate:"required"`
}

type AddTxnResponse struct {
	ID       uint64  `json:"id"`
	ParentID *uint64 `json:"parent_id"`

	Currency meta.Currency   `json:"currency"`
	UID      meta.UID        `json:"user_id"`
	Amount   decimal.Decimal `json:"amount"`
	Balance  decimal.Decimal `json:"balance"`
	TypeCode string          `json:"type_code"`
	Ref      string          `json:"ref"`
}
