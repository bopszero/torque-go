package blockchainmod

import (
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func NewLitecoinBitcoreClient() (*BitcoreUtxoLikeClient, error) {
	if config.BlockchainUseTestnet {
		return nil, utils.IssueErrorf("bitcore testnet hasn't been supported yet")
	} else {
		return NewLitecoinBitcoreMainnetClient(), nil
	}
}

func NewLitecoinBitcoreMainnetClient() *BitcoreUtxoLikeClient {
	return &BitcoreUtxoLikeClient{
		newBitcoreClient(
			constants.CurrencyLitecoin, constants.CurrencySubBitcoinSatoshi,
			BitcoreChainMainnet,
		),
	}
}
