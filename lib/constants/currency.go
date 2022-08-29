package constants

import (
	"fmt"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	CurrencyUSD    = meta.Currency("USD")
	CurrencyTorque = meta.Currency("TORQ")

	CurrencyBitcoin        = meta.Currency("BTC")
	CurrencyBitcoinCash    = meta.Currency("BCH")
	CurrencyBitcoinCashABC = meta.Currency("BCHA")
	CurrencyEthereum       = meta.Currency("ETH")
	CurrencyLitecoin       = meta.Currency("LTC")
	CurrencyTetherUSD      = meta.Currency("USDT")
	CurrencyTron           = meta.Currency("TRX")
	CurrencyRipple         = meta.Currency("XRP")

	CurrencySubBitcoinSatoshi = meta.Currency("satoshi")
	CurrencySubEthereumGwei   = meta.Currency("gwei")
	CurrencySubEthereumWei    = meta.Currency("wei")
	CurrencySubTronSun        = meta.Currency("sun")
	CurrencySubTronBandwidth  = meta.Currency("bandwidth")
	CurrencySubTronEnergy     = meta.Currency("energy")
	CurrencySubRippleDrop     = meta.Currency("drop")
)

var (
	CurrencyTorquePriceUSDT decimal.Decimal
	CurrencyLocalList       = []meta.Currency{
		CurrencyTorque,
	}
	CurrencyLocalSet = comtypes.NewHashSetFromListF(CurrencyLocalList)
)

func init() {
	CurrencyTorquePriceUSDT, _ = decimal.NewFromString("0.05")

	initConversion()
}

var CurrencyConversionRateMap map[string]meta.CurrencyConversionRate

func initConversion() {
	conversionRates := []meta.CurrencyConversionRate{
		{
			FromCurrency: CurrencyTetherUSD,
			ToCurrency:   CurrencyTorque,
			Value:        decimal.NewFromInt(20),
		},

		{
			FromCurrency: CurrencyBitcoin,
			ToCurrency:   CurrencySubBitcoinSatoshi,
			Value:        decimal.NewFromInt(100000000),
		},
		{
			FromCurrency: CurrencyBitcoinCash,
			ToCurrency:   CurrencySubBitcoinSatoshi,
			Value:        decimal.NewFromInt(100000000),
		},
		{
			FromCurrency: CurrencyBitcoinCashABC,
			ToCurrency:   CurrencySubBitcoinSatoshi,
			Value:        decimal.NewFromInt(100000000),
		},
		{
			FromCurrency: CurrencyLitecoin,
			ToCurrency:   CurrencySubBitcoinSatoshi,
			Value:        decimal.NewFromInt(100000000),
		},
		{
			FromCurrency: CurrencyEthereum,
			ToCurrency:   CurrencySubEthereumGwei,
			Value:        decimal.NewFromInt(1000000000),
		},
		{
			FromCurrency: CurrencyEthereum,
			ToCurrency:   CurrencySubEthereumWei,
			Value:        decimal.NewFromInt(1000000000000000000),
		},
		{
			FromCurrency: CurrencySubEthereumGwei,
			ToCurrency:   CurrencySubEthereumWei,
			Value:        decimal.NewFromInt(1000000000),
		},
		{
			FromCurrency: CurrencyTron,
			ToCurrency:   CurrencySubTronSun,
			Value:        decimal.NewFromInt(1000000),
		},
		{
			FromCurrency: CurrencyTron,
			ToCurrency:   CurrencySubTronBandwidth,
			Value:        decimal.NewFromInt(100000),
		},
		{
			FromCurrency: CurrencySubTronBandwidth,
			ToCurrency:   CurrencySubTronSun,
			Value:        decimal.NewFromInt(10),
		},
		{
			FromCurrency: CurrencyTron,
			ToCurrency:   CurrencySubTronEnergy,
			Value:        decimal.NewFromInt(100000),
		},
		{
			FromCurrency: CurrencySubTronEnergy,
			ToCurrency:   CurrencySubTronSun,
			Value:        decimal.NewFromInt(10),
		},
		{
			FromCurrency: CurrencyRipple,
			ToCurrency:   CurrencySubRippleDrop,
			Value:        decimal.NewFromInt(1000000),
		},
	}
	CurrencyConversionRateMap = make(map[string]meta.CurrencyConversionRate, len(conversionRates))
	for _, rate := range conversionRates {
		key := fmt.Sprintf("%v-%v", rate.FromCurrency, rate.ToCurrency)
		CurrencyConversionRateMap[key] = rate

		reversedKey := fmt.Sprintf("%v-%v", rate.ToCurrency, rate.FromCurrency)
		reversedRate := meta.CurrencyConversionRate{
			FromCurrency: rate.ToCurrency,
			ToCurrency:   rate.FromCurrency,
			Value:        DecimalOne.DivRound(rate.Value, 32),
		}
		CurrencyConversionRateMap[reversedKey] = reversedRate
	}
}
