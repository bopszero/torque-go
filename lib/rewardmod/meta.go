package rewardmod

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcache"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

var CurrencyDepositMinThresholdMap = comcache.NewCacheObject(15*time.Minute, func() (interface{}, error) {
	currencyInfoMap, err := currencymod.GetAllLegacyCurrencyInfoMapFast()
	if err != nil {
		return nil, err
	}
	minThresholdMap := make(map[meta.Currency]decimal.Decimal, len(currencyInfoMap))
	for currency, currencyInfo := range currencyInfoMap {
		minThresholdMap[currency] = currencyInfo.DepositMinThreshold
	}
	return minThresholdMap, nil
})
