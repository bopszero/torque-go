package currencymod

import (
	"time"

	"gitlab.com/snap-clickstaff/go-common/comcache"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

var (
	cachedAllLegacyCurrencyInfoMap = comcache.NewCacheObject(
		15*time.Second,
		func() (interface{}, error) {
			return GetAllLegacyCurrencyInfoMap()
		},
	)
	cachedAllLegacyCurrencyInfoIdMap = comcache.NewCacheObject(
		15*time.Second,
		func() (interface{}, error) {
			infoMap, err := GetAllLegacyCurrencyInfoMap()
			if err != nil {
				return nil, err
			}
			infoIdMap := make(map[uint16]models.LegacyCurrencyInfo, len(infoMap))
			for _, info := range infoMap {
				infoIdMap[info.ID] = info
			}
			return infoIdMap, nil
		},
	)
	cachedAllNetworkCurrencyInfoMap = comcache.NewCacheObject(
		15*time.Second,
		func() (interface{}, error) {
			return GetAllNetworkCurrencyInfoMap()
		},
	)
)

func GetAllLegacyCurrencyInfoMap() (infoMap LegacyCurrencyInfoMap, err error) {
	var infoList []models.LegacyCurrencyInfo
	err = database.
		GetDbSlave().
		Find(&infoList).
		Error
	if err != nil {
		err = utils.WrapError(err)
		return
	}

	infoMap = make(LegacyCurrencyInfoMap, len(infoList))
	for _, info := range infoList {
		infoMap[info.Currency] = info
	}
	return
}

func GetAllLegacyCurrencyInfoMapFast() (infoMap LegacyCurrencyInfoMap, err error) {
	infoMapObj, err := cachedAllLegacyCurrencyInfoMap.Get()
	if err != nil {
		return
	}
	infoMap = infoMapObj.(LegacyCurrencyInfoMap)
	return
}

func GetAllLegacyCurrencyInfoMapFastF() LegacyCurrencyInfoMap {
	infoMap, err := GetAllLegacyCurrencyInfoMapFast()
	comutils.PanicOnError(err)
	return infoMap
}

func GetLegacyCurrencyInfoFast(currency meta.Currency) (currencyInfo models.LegacyCurrencyInfo, err error) {
	infoMap, err := GetAllLegacyCurrencyInfoMapFast()
	if err != nil {
		return
	}
	currencyInfo, ok := infoMap[currency]
	if !ok {
		err = utils.WrapError(constants.ErrorDataNotFound)
		return
	}
	return
}

func GetLegacyCurrencyInfoFastF(currency meta.Currency) models.LegacyCurrencyInfo {
	currencyInfo, err := GetLegacyCurrencyInfoFast(currency)
	comutils.PanicOnError(err)

	return currencyInfo
}

func GetAllLegacyCurrencyInfoIdMapFast() (infoIdMap map[uint16]models.LegacyCurrencyInfo, err error) {
	infoIdMapObj, err := cachedAllLegacyCurrencyInfoIdMap.Get()
	if err != nil {
		return
	}
	infoIdMap = infoIdMapObj.(map[uint16]models.LegacyCurrencyInfo)
	return
}

func GetAllLegacyCurrencyInfoIdMapFastF() map[uint16]models.LegacyCurrencyInfo {
	infoIdMap, err := GetAllLegacyCurrencyInfoIdMapFast()
	comutils.PanicOnError(err)
	return infoIdMap
}

func GetLegacyCurrencyInfoByIdFast(infoID uint16) (info models.LegacyCurrencyInfo, err error) {
	infoIdMap, err := GetAllLegacyCurrencyInfoIdMapFast()
	if err != nil {
		return
	}
	info, ok := infoIdMap[infoID]
	if !ok {
		err = utils.WrapError(err)
		return
	}
	return info, nil
}

func GetLegacyCurrencyInfoByIdFastF(infoID uint16) models.LegacyCurrencyInfo {
	info, err := GetLegacyCurrencyInfoByIdFast(infoID)
	comutils.PanicOnError(err)
	return info
}

func GetAllNetworkCurrencyInfoMap() (infoMap map[meta.BlockchainCurrencyIndex]models.NetworkCurrency, err error) {
	var infoList []models.NetworkCurrency
	err = database.GetDbSlave().
		Order(dbquery.OrderAsc(models.CommonColPriority)).
		Find(&infoList).
		Error
	if err != nil {
		err = utils.WrapError(err)
		return
	}

	infoMap = make(map[meta.BlockchainCurrencyIndex]models.NetworkCurrency, len(infoList))
	for _, info := range infoList {
		key := meta.BlockchainCurrencyIndex{
			Currency: info.Currency,
			Network:  info.Network,
		}
		infoMap[key] = info
	}
	return
}

func GetAllNetworkCurrencyInfoMapFast() (infoMap map[meta.BlockchainCurrencyIndex]models.NetworkCurrency, err error) {
	infoMapObj, err := cachedAllNetworkCurrencyInfoMap.Get()
	if err != nil {
		return
	}
	infoMap = infoMapObj.(map[meta.BlockchainCurrencyIndex]models.NetworkCurrency)
	return
}

func GetAllNetworkCurrencyInfoMapFastF() map[meta.BlockchainCurrencyIndex]models.NetworkCurrency {
	infoMap, err := GetAllNetworkCurrencyInfoMapFast()
	comutils.PanicOnError(err)

	return infoMap
}

func GetNetworkCurrencyInfoFast(currency meta.Currency, network meta.BlockchainNetwork) (
	networkCurrency models.NetworkCurrency, err error,
) {
	infoMap, err := GetAllNetworkCurrencyInfoMapFast()
	if err != nil {
		return
	}

	key := meta.BlockchainCurrencyIndex{
		Currency: currency,
		Network:  network,
	}
	networkCurrency, ok := infoMap[key]
	if !ok {
		err = utils.WrapError(constants.ErrorDataNotFound)
		return
	}
	return
}

func GetNetworkCurrencyInfoFastF(
	currency meta.Currency, network meta.BlockchainNetwork,
) models.NetworkCurrency {
	networkCurrency, err := GetNetworkCurrencyInfoFast(currency, network)
	comutils.PanicOnError(err)

	return networkCurrency
}
