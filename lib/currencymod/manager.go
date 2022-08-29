package currencymod

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcache"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

var (
	cachedAllCurrencyInfoMap = comcache.NewCacheObject(
		15*time.Second,
		func() (interface{}, error) {
			return GetAllCurrencyInfoMap()
		},
	)
	cachedAllBlockchainNetworkInfoMap = comcache.NewCacheObject(
		5*time.Minute,
		func() (interface{}, error) {
			return GetAllBlockchainNetworkInfoMap()
		},
	)
)

func GetAllCurrencyInfoMap() (infoMap CurrencyInfoMap, err error) {
	var infoList []models.CurrencyInfo
	err = database.
		GetDbF(database.AliasWalletSlave).
		Find(&infoList).
		Error
	if err != nil {
		err = utils.WrapError(err)
		return
	}
	infoMap = make(CurrencyInfoMap, len(infoList))
	for _, info := range infoList {
		infoMap[info.Currency] = info
	}
	return
}

func GetAllCurrencyInfoMapFast() (infoMap CurrencyInfoMap, err error) {
	infoMapObj, err := cachedAllCurrencyInfoMap.Get()
	if err == nil {
		infoMap = infoMapObj.(CurrencyInfoMap)
	}
	return
}

func GetAllCurrencyInfoMapFastF() CurrencyInfoMap {
	infoMap, err := GetAllCurrencyInfoMapFast()
	comutils.PanicOnError(err)
	return infoMap
}

func GetNonFiatCurrencyInfoMapFastF() CurrencyInfoMap {
	var (
		allInfoMap = GetAllCurrencyInfoMapFastF()
		infoMap    = make(CurrencyInfoMap)
	)
	for currency, info := range allInfoMap {
		if !info.IsFiat.Bool {
			infoMap[currency] = info
		}
	}
	return infoMap
}

func GetWalletCurrencyInfoMapFastF() CurrencyInfoMap {
	var (
		allInfoMap = GetAllCurrencyInfoMapFastF()
		infoMap    = make(CurrencyInfoMap)
	)
	for currency, info := range allInfoMap {
		if IsValidWalletInfo(info) {
			infoMap[currency] = info
		}
	}
	return infoMap
}

func GetTradingCurrencyInfoMapFastF() CurrencyInfoMap {
	var (
		allInfoMap = GetAllCurrencyInfoMapFastF()
		infoMap    = make(CurrencyInfoMap)
	)
	for currency, info := range allInfoMap {
		if IsValidTradingInfo(info) {
			infoMap[currency] = info
		}
	}
	return infoMap
}

func GetCurrencyInfoFast(currency meta.Currency) (currencyInfo models.CurrencyInfo, err error) {
	infoMap, err := GetAllCurrencyInfoMapFast()
	if err != nil {
		return
	}
	currencyInfo, ok := infoMap[currency]
	if !ok {
		err = utils.WrapError(constants.ErrorCurrency)
		return
	}
	return
}

func GetCurrencyInfoFastF(currency meta.Currency) models.CurrencyInfo {
	currencyInfo, err := GetCurrencyInfoFast(currency)
	comutils.PanicOnError(err)

	return currencyInfo
}

func GetCurrencyPriceUsdtFast(currency meta.Currency) (decimal.Decimal, error) {
	currencyInfo, err := GetCurrencyInfoFast(currency)
	if err != nil {
		return decimal.Zero, err
	}

	return currencyInfo.PriceUSDT, nil
}

func GetCurrencyPriceUsdtFastF(currency meta.Currency) decimal.Decimal {
	priceUSDT, err := GetCurrencyPriceUsdtFast(currency)
	comutils.PanicOnError(err)

	return priceUSDT
}

func GetCurrencyPriceTorqueFast(currency meta.Currency) (decimal.Decimal, error) {
	currencyInfo, err := GetCurrencyInfoFast(currency)
	if err != nil {
		return decimal.Zero, err
	}

	priceTorque := comutils.DecimalDivide(currencyInfo.PriceUSDT, constants.CurrencyTorquePriceUSDT)
	return priceTorque, nil
}

func GetCurrencyPriceTorqueFastF(currency meta.Currency) decimal.Decimal {
	priceTorque, err := GetCurrencyPriceTorqueFast(currency)
	comutils.PanicOnError(err)

	return priceTorque
}

func ConvertUsdtToTorque(usdtValue decimal.Decimal) decimal.Decimal {
	amount := ConvertAmountF(
		meta.CurrencyAmount{
			Currency: constants.CurrencyTetherUSD,
			Value:    usdtValue,
		},
		constants.CurrencyTorque,
	)
	return amount.Value
}

func GetAllBlockchainNetworkInfoMap() (infoMap map[meta.BlockchainNetwork]models.BlockchainNetworkInfo, err error) {
	var infoList []models.BlockchainNetworkInfo
	err = database.
		GetDbF(database.AliasWalletSlave).
		Find(&infoList).
		Error
	if err != nil {
		err = utils.WrapError(err)
		return
	}

	infoMap = make(map[meta.BlockchainNetwork]models.BlockchainNetworkInfo, len(infoList))
	for _, info := range infoList {
		infoMap[info.Network] = info
	}
	return
}

func GetAllBlockchainNetworkInfoMapFast() (infoMap map[meta.BlockchainNetwork]models.BlockchainNetworkInfo, err error) {
	infoMapObj, err := cachedAllBlockchainNetworkInfoMap.Get()
	if err == nil {
		infoMap = infoMapObj.(map[meta.BlockchainNetwork]models.BlockchainNetworkInfo)
	}
	return
}

func GetAllBlockchainNetworkInfoMapFastF() map[meta.BlockchainNetwork]models.BlockchainNetworkInfo {
	infoList, err := GetAllBlockchainNetworkInfoMapFast()
	comutils.PanicOnError(err)

	return infoList
}

func GetBlockchainNetworkInfoFast(network meta.BlockchainNetwork) (
	networkInfo models.BlockchainNetworkInfo, err error,
) {
	infoMap, err := GetAllBlockchainNetworkInfoMapFast()
	if err != nil {
		return
	}
	networkInfo, ok := infoMap[network]
	if !ok {
		err = utils.WrapError(constants.ErrorDataNotFound)
		return
	}
	return
}

func ValidateTradingCurrency(currency meta.Currency) error {
	info, err := GetCurrencyInfoFast(currency)
	if err != nil {
		return err
	}
	if !IsValidTradingInfo(info) {
		return utils.WrapError(constants.ErrorCurrency)
	}
	return nil
}
