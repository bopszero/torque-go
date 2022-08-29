package blockchainmod

import (
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type (
	CoinLoader        func() Coin
	BlockchainCoinMap map[meta.BlockchainCurrencyIndex]Coin
	NetworkCoinMap    map[meta.BlockchainNetwork]Coin
)

var (
	blockchainCoinLoaderMap     = map[meta.BlockchainCurrencyIndex]CoinLoader{}

	blockchainCoinMapProxy = comtypes.NewSingleton(func() interface{} {
		coinMap := make(BlockchainCoinMap, len(blockchainCoinLoaderMap))

		for index, loader := range blockchainCoinLoaderMap {
			coin := loader()
			if index.Currency != coin.GetCurrency() ||
				(index.Network != "" && index.Network != coin.GetNetwork()) {
				panic(utils.IssueErrorf(
					"blockchain init coin index mismatched `%v-%v`",
					index.Currency, index.Network,
				))
			}
			if _, exists := coinMap[index]; exists {
				panic(utils.IssueErrorf("blockchain init coin duplicated index `%v`", index))
			}
			coinMap[index] = coin
		}

		return coinMap
	})
	networkMainCoinMapProxy = comtypes.NewSingleton(func() interface{} {
		var (
			currencyCoinMap = GetCoinMap()
			networkCoinMap  = make(NetworkCoinMap)
		)
		for index, coin := range currencyCoinMap {
			if index.Network == "" || coin.GetCurrency() != coin.GetNetworkCurrency() {
				continue
			}
			network := coin.GetNetwork()
			if _, exists := networkCoinMap[network]; exists {
				panic(utils.IssueErrorf("blockchain init main coin duplicated network `%v`", network))
			}
			networkCoinMap[network] = coin
		}

		return networkCoinMap
	})
)

func RegisterCoinLoader(index meta.BlockchainCurrencyIndex, loader CoinLoader) error {
	if _, exist := blockchainCoinLoaderMap[index]; exist {
		return utils.IssueErrorf("blockchain register coin duplicated index `%v`", index)
	}

	blockchainCoinLoaderMap[index] = loader
	return nil
}

func RegisterNativeCoinLoader(currency meta.Currency, loader CoinLoader) error {
	return RegisterCoinLoader(meta.BlockchainCurrencyIndex{Currency: currency}, loader)
}

func GetCoinMap() BlockchainCoinMap {
	return blockchainCoinMapProxy.Get().(BlockchainCoinMap)
}

func GetNetworkMainCoinMap() NetworkCoinMap {
	return networkMainCoinMapProxy.Get().(NetworkCoinMap)
}

func GetCoinNative(currency meta.Currency) (Coin, error) {
	return GetCoin(currency, "")
}

func GetCoinNativeF(currency meta.Currency) Coin {
	coin, err := GetCoinNative(currency)
	comutils.PanicOnError(err)

	return coin
}

func GetCoin(currency meta.Currency, network meta.BlockchainNetwork) (Coin, error) {
	index := meta.BlockchainCurrencyIndex{
		Currency: currency,
		Network:  network,
	}
	coin, ok := GetCoinMap()[index]
	if !ok {
		return nil, utils.WrapError(constants.ErrorFeatureNotSupport)
	}

	return coin, nil
}

func GetCoinF(currency meta.Currency, network meta.BlockchainNetwork) Coin {
	coin, err := GetCoin(currency, network)
	comutils.PanicOnError(err)

	return coin
}

func GetNetworkMainCoin(network meta.BlockchainNetwork) (Coin, error) {
	coin, ok := GetNetworkMainCoinMap()[network]
	if !ok {
		return nil, utils.IssueErrorf("blockchain network `%v` hasn't been supported yet")
	}

	return coin, nil
}

func IsNativeCurrency(currency meta.Currency) bool {
	index := meta.BlockchainCurrencyIndex{Currency: currency}
	_, ok := GetCoinMap()[index]
	return ok
}

func IsSupportedIndex(currency meta.Currency, network meta.BlockchainNetwork) bool {
	index := meta.BlockchainCurrencyIndex{
		Currency: currency,
		Network:  network,
	}
	_, ok := GetCoinMap()[index]
	return ok
}
