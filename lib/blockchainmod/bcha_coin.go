package blockchainmod

import (
	"fmt"

	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

func init() {
	var (
		currency     = constants.CurrencyBitcoinCashABC
		indexMainnet = meta.BlockchainCurrencyIndex{
			Currency: currency,
			Network:  constants.BlockchainNetworkBitcoinCashABC,
		}
		indexTestnet = meta.BlockchainCurrencyIndex{
			Currency: currency,
			Network:  constants.BlockchainNetworkBitcoinCashAbcTestnet,
		}
		getterMainnet = func() Coin {
			return &BitcoinCashAbcCoin{BitcoinCashCoin{baseCoin{
				currencyIdx: indexMainnet,
				networkMain: indexMainnet.Network,
				networkTest: indexTestnet.Network,
			}}}
		}
		getterTestnet = func() Coin {
			return &BitcoinCashAbcCoin{BitcoinCashCoin{baseCoin{
				currencyIdx: indexTestnet,
				networkMain: indexMainnet.Network,
				networkTest: indexTestnet.Network,
			}}}
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

type BitcoinCashAbcCoin struct {
	BitcoinCashCoin
}

func (this BitcoinCashAbcCoin) GetTradingID() uint16 {
	return 21
}

func (this BitcoinCashAbcCoin) GenTxnExplorerURL(txnHash string) string {
	if this.IsUsingMainnet() {
		return fmt.Sprintf(ExplorerBitcoinCashAbcMainnetTxnUrlPattern, txnHash)
	} else {
		return fmt.Sprintf(ExplorerBitcoinCashAbcTestnetTxnUrlPattern, txnHash)
	}
}
