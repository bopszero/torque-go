package tradingbalance

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type CoinBalanceMap map[meta.Currency]decimal.Decimal

// func NewCoinBalanceMap(currencies []meta.Currency) CoinBalanceMap {
// 	balanceMap := make(CoinBalanceMap, len(currencies))
// 	for _, currency := range currencies {
// 		balanceMap[currency] = decimal.Zero
// 	}
// 	return balanceMap
// }

func (this CoinBalanceMap) CalcValueUSD(infoMap currencymod.CurrencyInfoMap) decimal.Decimal {
	var valueUSD decimal.Decimal
	for currency, value := range this {
		if value.IsZero() {
			continue
		}
		info, ok := infoMap[currency]
		if !ok {
			info = currencymod.GetCurrencyInfoFastF(currency)
		}
		valueUSD = valueUSD.Add(value.Mul(info.PriceUSD))
	}
	return valueUSD
}

func (this CoinBalanceMap) Format(infoMap currencymod.CurrencyInfoMap) CoinBalanceMap {
	formattedMap := make(CoinBalanceMap, len(infoMap))
	for currency := range infoMap {
		if balance, ok := this[currency]; ok {
			formattedMap[currency] = balance
		} else {
			formattedMap[currency] = decimal.Zero
		}
	}
	return formattedMap
}
