package blockchainmod

import (
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
)

func NewBitcoinCashBitcoreClient() *BitcoreUtxoLikeClient {
	if config.BlockchainUseTestnet {
		return NewBitcoinCashBitcoreTestnetClient()
	} else {
		return NewBitcoinCashBitcoreMainnetClient()
	}
}

func NewBitcoinCashBitcoreMainnetClient() *BitcoreUtxoLikeClient {
	client := NewBitcoreUtxoLikeClient(constants.CurrencyBitcoinCash, BitcoreChainMainnet)
	return &client
}

func NewBitcoinCashBitcoreTestnetClient() *BitcoreUtxoLikeClient {
	client := NewBitcoreUtxoLikeClient(constants.CurrencyBitcoinCash, BitcoreChainTestnet)
	return &client
}
