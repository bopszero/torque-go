package blockchainmod

import (
	"math"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type FeeInfo struct {
	Currency         meta.Currency   `json:"currency"`
	Price            decimal.Decimal `json:"price"`
	PriceHigh        decimal.Decimal `json:"price_high,omitempty"`
	PriceLow         decimal.Decimal `json:"price_low,omitempty"`
	LimitMaxQuantity uint32          `json:"limit_max_quantity"`
	LimitMinValue    decimal.Decimal `json:"limit_min_value,omitempty"`
	LimitMaxValue    decimal.Decimal `json:"limit_max_value"`

	BaseCurrency   meta.Currency   `json:"base_currency,omitempty"`
	ToBaseMultiple decimal.Decimal `json:"to_base_multiple,omitempty"`
}

func NewFeeInfo(baseCurrency meta.Currency, currency meta.Currency, price decimal.Decimal) (
	feeInfo FeeInfo, err error,
) {
	rate, err := currencymod.GetConversionRate(currency, baseCurrency)
	if err != nil {
		return
	}

	maxValue := decimal.NewFromInt(MaxBalance)
	feeInfo = FeeInfo{
		BaseCurrency:   baseCurrency,
		ToBaseMultiple: rate.Value,

		Currency: currency,
		Price:    price,

		LimitMaxValue: maxValue,
		LimitMaxQuantity: uint32(comutils.MinInt64(
			maxValue.Div(price).IntPart(),
			math.MaxUint32,
		)),
	}
	return
}

func (this *FeeInfo) IsEmpty() bool {
	return this.Currency == ""
}

func (this *FeeInfo) UsePriceHigh() {
	if !this.PriceHigh.IsZero() {
		this.Price = this.PriceHigh
	}
}

func (this *FeeInfo) UsePriceLow() {
	if !this.PriceLow.IsZero() {
		this.Price = this.PriceLow
	}
}

func (this *FeeInfo) GetBasePrice() decimal.Decimal {
	currencyAmount := currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: this.Currency,
			Value:    this.Price,
		},
		this.BaseCurrency,
	)
	return currencyAmount.Value
}

func (this *FeeInfo) GetBaseValue() decimal.Decimal {
	return this.GetBasePrice().Mul(decimal.NewFromInt(int64(this.LimitMaxQuantity)))
}

func (this *FeeInfo) GetBasePriceHigh() decimal.Decimal {
	priceHigh := this.PriceHigh
	if priceHigh.IsZero() {
		priceHigh = this.Price
	}

	currencyAmount := currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: this.Currency,
			Value:    priceHigh,
		},
		this.BaseCurrency,
	)
	return currencyAmount.Value
}

func (this *FeeInfo) GetBaseLimitMinValue() decimal.Decimal {
	currencyAmount := currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: this.Currency,
			Value:    this.LimitMinValue,
		},
		this.BaseCurrency,
	)
	return currencyAmount.Value
}

func (this *FeeInfo) GetBaseLimitMaxValue() decimal.Decimal {
	currencyAmount := currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: this.Currency,
			Value:    this.LimitMaxValue,
		},
		this.BaseCurrency,
	)
	return currencyAmount.Value
}

func (this *FeeInfo) truncateAmount(
	amount decimal.Decimal,
	unitCurrency meta.Currency,
) meta.CurrencyAmount {
	unitAmount := currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: this.Currency,
			Value:    amount,
		},
		unitCurrency,
	)
	unitAmount.Value = unitAmount.Value.Truncate(0)

	return currencymod.ConvertAmountF(unitAmount, this.Currency)
}

func (this *FeeInfo) SetLimitMaxValue(amount decimal.Decimal, unitCurrency meta.Currency) {
	truncatedAmount := this.truncateAmount(amount, unitCurrency)
	truncatedAmount.Value = decimal.Max(truncatedAmount.Value, this.LimitMinValue)

	this.LimitMaxValue = truncatedAmount.Value
	this.LimitMaxQuantity = uint32(truncatedAmount.Value.Div(this.Price).IntPart())
}

func (this *FeeInfo) SetLimitMaxQuantity(maxQuantity uint32) {
	truncatedAmount := this.truncateAmount(
		this.Price.Mul(decimal.NewFromInt(int64(maxQuantity))),
		this.Currency)
	truncatedAmount.Value = decimal.Max(truncatedAmount.Value, this.LimitMinValue)

	this.LimitMaxQuantity = maxQuantity
	this.LimitMaxValue = comutils.DecimalClamp(
		truncatedAmount.Value,
		this.LimitMinValue, this.LimitMaxValue,
	)
}
