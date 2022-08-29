package binanceconv

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcache"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/thirdpartymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func GetBalance(currency meta.Currency) (decimal.Decimal, error) {
	client := thirdpartymod.GetBinanceSystemClient()
	account, err := client.GetAccountInfo()
	if err != nil {
		return decimal.Zero, utils.WrapError(err)
	}

	currencyStr := currency.String()
	for _, balance := range account.Balances {
		if balance.Asset == currencyStr {
			freeBalance, err := decimal.NewFromString(balance.Free)
			if err != nil {
				return decimal.Zero, utils.WrapError(err)
			}

			return freeBalance, nil
		}
	}

	return decimal.Zero, nil
}

func GetBalanceFast(currency meta.Currency) (decimal.Decimal, error) {
	cacheKey := fmt.Sprintf("binance_info:balance:%v", currency)

	var balance decimal.Decimal
	err := comcache.GetOrCreate(
		comcache.GetRemoteCache(),
		cacheKey,
		10*time.Second,
		&balance,
		func() (interface{}, error) {
			return GetBalance(currency)
		},
	)
	if err != nil {
		return decimal.Zero, err
	}

	return balance, nil
}
