package blockchainmod

import (
	"math/rand"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

func tronToSun(value decimal.Decimal) int64 {
	amount := currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: constants.CurrencyTron,
			Value:    value,
		},
		constants.CurrencySubTronSun,
	)
	return amount.Value.IntPart()
}

func tronMakeErrorTimeMs(timeMs int64, errRange time.Duration) int64 {
	rangeMs := errRange.Milliseconds()
	return timeMs - rangeMs/2 + (rand.Int63() % rangeMs)
}
