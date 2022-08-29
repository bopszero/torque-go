package blockchainmod

import (
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func NewBitcoinBitcoreClient() (*BitcoreUtxoLikeClient, error) {
	if config.BlockchainUseTestnet {
		return nil, utils.IssueErrorf("bitcore testnet hasn't been supported yet")
	} else {
		return NewBitcoinBitcoreMainnetClient(), nil
	}
}

func NewBitcoinBitcoreMainnetClient() *BitcoreUtxoLikeClient {
	return &BitcoreUtxoLikeClient{
		newBitcoreClient(
			constants.CurrencyBitcoin, constants.CurrencySubBitcoinSatoshi,
			BitcoreChainMainnet,
		),
	}
}
