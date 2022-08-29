package rewardmod

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

func GetCurrencyDepositMinThresholdMap() map[meta.Currency]decimal.Decimal {
	thresholdMap, err := CurrencyDepositMinThresholdMap.Get()
	comutils.PanicOnError(err)
	return thresholdMap.(map[meta.Currency]decimal.Decimal)
}
