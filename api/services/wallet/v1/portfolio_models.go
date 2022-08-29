package v1

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type PortfolioGetOverviewResponse struct {
	Balances []PortfolioGetOverviewBalance `json:"balances"`
}

type PortfolioGetOverviewBalance struct {
	Currency          meta.Currency          `json:"currency"`
	Network           meta.BlockchainNetwork `json:"network"`
	CurrencyPriceUSD  decimal.Decimal        `json:"currency_price_usd"`
	CurrencyPriceUSDT decimal.Decimal        `json:"currency_price_usdt"`
	Amount            decimal.Decimal        `json:"amount"`
	UpdateTime        int64                  `json:"update_time"`

	IsAvailable bool `json:"is_available"`
}

type PortfolioGetCurrencyRequest struct {
	Currency meta.Currency          `json:"currency" validate:"required"`
	Network  meta.BlockchainNetwork `json:"network"`
}

type PortfolioGetCurrencyResponse struct {
	Currency  meta.Currency          `json:"currency"`
	Network   meta.BlockchainNetwork `json:"network"`
	PriceUSD  decimal.Decimal        `json:"price_usd"`
	PriceUSDT decimal.Decimal        `json:"price_usdt"`
	AccountNo string                 `json:"account_no"`
	Balance   decimal.Decimal        `json:"balance"`
	Notice    string                 `json:"notice"`
}

type PortfolioListOrdersRequest struct {
	Currency meta.Currency          `json:"currency" validate:"required"`
	Network  meta.BlockchainNetwork `json:"network"`
	Paging   meta.Paging            `json:"paging"`
}

type PortfolioListOrdersResponse struct {
	Items  []Order     `json:"items"`
	Paging meta.Paging `json:"paging"`
}

type PortfolioExportOrdersRequest struct {
	Currency meta.Currency          `json:"currency" validate:"required"`
	Network  meta.BlockchainNetwork `json:"network"`
	FromTime int64                  `json:"from_time" validate:"required"`
	ToTime   int64                  `json:"to_time" validate:"required"`
}

type PortfolioExportOrdersResponse struct {
	Token string `json:"token"`
}

type PortfolioExportOrdersTokenModel struct {
	PortfolioExportOrdersRequest
	UID meta.UID `json:"uid"`
}

type PortfolioGetCurrencyOrdersRequest struct {
	Currency  meta.Currency          `json:"currency" validate:"required"`
	Network   meta.BlockchainNetwork `json:"network"`
	Reference string                 `json:"ref" validate:"required"`
}

type PortfolioGetCurrencyOrderResponse Order
