package blockchainmod

// import (
// 	"gitlab.com/snap-clickstaff/go-common/comutils"
// 	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
// 	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
// )

// func init() {
// 	var (
// 		currency     = constants.CurrencyTetherUSD
// 		indexMainnet = meta.BlockchainCurrencyIndex{
// 			Currency: currency,
// 			Network:  constants.BlockchainNetworkTron,
// 		}
// 		indexTestnet = meta.BlockchainCurrencyIndex{
// 			Currency: currency,
// 			Network:  constants.BlockchainNetworkTronTestShasta,
// 		}
// 		getterMainnet = func() Coin {
// 			return &TronTokenCoin{
// 				TronCoin:  newTronCoin(indexMainnet),
// 				tokenMeta: TronTokenMetaMainnetTetherUSD,
// 			}
// 		}
// 		getterTestnet = func() Coin {
// 			return &TronTokenCoin{
// 				TronCoin:  newTronCoin(indexTestnet),
// 				tokenMeta: TronTokenMetaTestShastaTetherUSD,
// 			}
// 		}
// 	)
// 	comutils.PanicOnError(
// 		RegisterCoinLoader(indexMainnet, getterMainnet),
// 	)
// 	comutils.PanicOnError(
// 		RegisterCoinLoader(indexTestnet, getterTestnet),
// 	)
// }
