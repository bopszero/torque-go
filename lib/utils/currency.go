package utils

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
)

func NormalizeTradingAmount(amount decimal.Decimal) decimal.Decimal {
	return amount.Truncate(constants.AmountTradingMaxDecimalPlaces)
}
