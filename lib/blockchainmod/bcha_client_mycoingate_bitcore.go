package blockchainmod

import (
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
)

func NewBitcoinCashAbcBitcoreClient() *BitcoreUtxoLikeClient {
	if config.BlockchainUseTestnet {
		return NewBitcoinCashAbcBitcoreTestnetClient()
	} else {
		return NewBitcoinCashAbcBitcoreMainnetClient()
	}
}

func NewBitcoinCashAbcBitcoreMainnetClient() *BitcoreUtxoLikeClient {
	client := NewBitcoreUtxoLikeClient(constants.CurrencyBitcoinCashABC, BitcoreChainMainnet)
	return &client
}

func NewBitcoinCashAbcBitcoreTestnetClient() *BitcoreUtxoLikeClient {
	client := NewBitcoreUtxoLikeClient(constants.CurrencyBitcoinCashABC, BitcoreChainTestnet)
	return &client
}
