package blockchainmod

import (
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

func init() {
	var (
		currency     = constants.CurrencyTetherUSD
		indexMainnet = meta.BlockchainCurrencyIndex{
			Currency: currency,
			Network:  constants.BlockchainNetworkEthereum,
		}
		indexTestnet = meta.BlockchainCurrencyIndex{
			Currency: currency,
			Network:  constants.BlockchainNetworkEthereumTestRopsten,
		}
		getterMainnet = func() Coin {
			return &EthereumTokenCoin{
				EthereumCoin: newEthereumCoin(indexMainnet),
				tokenMeta:    EthereumTokenMetaMainnetTetherUSD,
			}
		}
		getterTestnet = func() Coin {
			return &EthereumTokenCoin{
				EthereumCoin: newEthereumCoin(indexTestnet),
				tokenMeta:    EthereumTokenMetaTestRopstenTetherUSD,
			}
		}
	)
	comutils.PanicOnError(
		RegisterNativeCoinLoader(
			currency,
			func() Coin {
				if config.BlockchainUseTestnet {
					return getterTestnet()
				} else {
					return getterMainnet()
				}
			},
		),
	)
	comutils.PanicOnError(
		RegisterCoinLoader(indexMainnet, getterMainnet),
	)
	comutils.PanicOnError(
		RegisterCoinLoader(indexTestnet, getterTestnet),
	)
}
