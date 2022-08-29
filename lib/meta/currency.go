package meta

import (
	"strings"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comutils"
)

type Currency string

func (this Currency) String() string {
	return string(this)
}

func (this Currency) StringU() string {
	return strings.ToUpper(this.String())
}

func (this Currency) StringL() string {
	return strings.ToLower(this.String())
}

func NewCurrency(code string) Currency {
	return Currency(code)
}

func NewCurrencyU(code string) Currency {
	return NewCurrency(strings.ToUpper(code))
}

type CurrencyMeta struct {
	ID            uint16
	Code          Currency
	DecimalPlaces uint8
}

type CurrencyConversionRate struct {
	FromCurrency Currency        `json:"from_currency" validate:"required"`
	ToCurrency   Currency        `json:"to_currency" validate:"required"`
	Value        decimal.Decimal `json:"value" validate:"required"`
}

type CurrencyAmount struct {
	Currency Currency        `json:"currency" validate:"required"`
	Value    decimal.Decimal `json:"value" validate:"required"`
}

type AmountMarkup struct {
	isPercentage bool
	value        decimal.Decimal
}

func NewAmountModifier(strValue string) (*AmountMarkup, error) {
	handler := AmountMarkup{}
	err := handler.UnmarshalText([]byte(strValue))
	if err != nil {
		return nil, err
	}

	return &handler, nil
}

func (this *AmountMarkup) String() string {
	strValue := this.value.String()
	if this.isPercentage {
		strValue += "%"
	}

	return strValue
}

func (this *AmountMarkup) For(value decimal.Decimal) decimal.Decimal {
	if this.value.IsZero() {
		return value
	}

	if this.isPercentage {
		rate := comutils.DecimalDivide(this.value, decimal.NewFromInt(100))
		return value.Add(value.Mul(rate))
	}

	return value.Add(this.value)
}

func (this AmountMarkup) MarshalText() ([]byte, error) {
	return []byte(this.String()), nil
}

func (this *AmountMarkup) UnmarshalText(text []byte) (err error) {
	strVal := strings.TrimSpace(string(text))

	if strings.HasSuffix(strVal, "%") {
		this.isPercentage = true
		strVal = strVal[:len(strVal)-1]
	}

	this.value, err = decimal.NewFromString(strVal)
	return
}

func (this AmountMarkup) MarshalBinary() ([]byte, error) {
	return this.MarshalText()
}

func (this *AmountMarkup) UnmarshalBinary(data []byte) (err error) {
	return this.UnmarshalText(data)
}
