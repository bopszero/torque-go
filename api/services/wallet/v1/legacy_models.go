package v1

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type LegacyTxnListRequest struct {
	UID      meta.UID      `json:"uid" validate:"required"`
	Currency meta.Currency `json:"currency" validate:"required,eq=TORQ"`

	PagingPage  uint   `json:"page"`
	PagingLimit uint   `json:"limit"`
	OrderID     uint64 `json:"order_id"`
}

type LegacyTxnListItem struct {
	ID          uint64          `json:"id"`
	UID         meta.UID        `json:"uid"`
	Amount      decimal.Decimal `json:"amount"`
	Ref         string          `json:"ref"`
	Currency    meta.Currency   `json:"coin_name"`
	TypeCode    string          `json:"type_code"`
	Reversed    bool            `json:"reversed"`
	CreatedTime string          `json:"created_at"`

	Details      meta.O `json:"transaction,omitempty"`
	OrderDetails meta.O `json:"order,omitempty"`
}

type LegacyCurrencyPriceMapGetResponse map[meta.Currency]map[meta.Currency]float64
