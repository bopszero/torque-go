package currencymod

import (
	"fmt"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func NormalizeAmount(currency meta.Currency, amount decimal.Decimal) decimal.Decimal {
	if amount.IsZero() {
		return amount
	}
	currencyInfo, ok := GetAllCurrencyInfoMapFastF()[currency]
	if !ok {
		return amount
	}
	for dp := int32(currencyInfo.DecimalPlaces); dp < comutils.DecimalDivisionPrecision; dp++ {
		var (
			truncatedAmount = amount.Truncate(dp)
			errorRate       = comutils.DecimalDivide(truncatedAmount, amount)
		)
		if errorRate.GreaterThanOrEqual(constants.AmountNormalizeErrorThreshold) {
			return truncatedAmount
		}
	}

	return amount.Truncate(comutils.DecimalDivisionPrecision)
}

func NormalizeCurrencyAmount(amount meta.CurrencyAmount) meta.CurrencyAmount {
	amount.Value = NormalizeAmount(amount.Currency, amount.Value)
	return amount
}

func GetConversionRate(fromCurrency meta.Currency, toCurrency meta.Currency) (
	rate meta.CurrencyConversionRate, err error,
) {
	if fromCurrency == toCurrency {
		rate = meta.CurrencyConversionRate{
			FromCurrency: fromCurrency,
			ToCurrency:   toCurrency,
			Value:        constants.DecimalOne,
		}
		return
	}

	key := fmt.Sprintf("%v-%v", fromCurrency, toCurrency)
	rate, ok := constants.CurrencyConversionRateMap[key]
	if !ok {
		err = utils.IssueErrorf("cannot find the conversion rate for currency pair `%v`", key)
	}
	return
}

func ConvertAmount(amount meta.CurrencyAmount, toCurrency meta.Currency) (meta.CurrencyAmount, error) {
	if amount.Currency == "" || toCurrency == "" {
		return amount, utils.WrapError(constants.ErrorCurrency)
	}
	if amount.Currency == toCurrency {
		return amount, nil
	}

	toAmount := meta.CurrencyAmount{
		Currency: toCurrency,
		Value:    decimal.Zero,
	}

	if amount.Value.IsZero() {
		return toAmount, nil
	}
	conversionRate, err := GetConversionRate(amount.Currency, toCurrency)
	if err != nil {
		return toAmount, err
	}
	toAmount.Value = amount.Value.Mul(conversionRate.Value)

	return NormalizeCurrencyAmount(toAmount), nil
}

func ConvertAmountF(amount meta.CurrencyAmount, toCurrency meta.Currency) meta.CurrencyAmount {
	convertedAmount, err := ConvertAmount(amount, toCurrency)
	comutils.PanicOnError(err)

	return convertedAmount
}

func IsValidWalletInfo(currencyInfo models.CurrencyInfo) bool {
	return !currencyInfo.IsFiat.Bool &&
		0 < currencyInfo.PriorityWallet &&
		currencyInfo.PriorityWallet < PriorityCommonSoonThreshold

}

func IsDisplayableWalletInfo(currencyInfo models.CurrencyInfo) bool {
	return !currencyInfo.IsFiat.Bool && currencyInfo.PriorityWallet > 0
}

func IsValidTradingInfo(currencyInfo models.CurrencyInfo) bool {
	return !currencyInfo.IsFiat.Bool &&
		0 < currencyInfo.PriorityTrading &&
		currencyInfo.PriorityTrading < PriorityCommonSoonThreshold
}

func IsDisplayableTradingInfo(currencyInfo models.CurrencyInfo) bool {
	return !currencyInfo.IsFiat.Bool && currencyInfo.PriorityTrading > 0
}

func IsLocalCurrency(currency meta.Currency) bool {
	return constants.CurrencyLocalSet.Contains(currency)
}
