package blockchainmod

import (
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
)

func NewLitecoinSoChainClient() *SoChainUtxoLikeClient {
	if config.BlockchainUseTestnet {
		return NewLitecoinSoChainTestnetClient()
	} else {
		return NewLitecoinSoChainMainnetClient()
	}
}

func NewLitecoinSoChainMainnetClient() *SoChainUtxoLikeClient {
	return &SoChainUtxoLikeClient{
		currency:    constants.CurrencyLitecoin,
		networkCode: soChainNetworkLitecoin,

		httpClient: getSoChainClient(),
	}
}

func NewLitecoinSoChainTestnetClient() *SoChainUtxoLikeClient {
	return &SoChainUtxoLikeClient{
		currency:    constants.CurrencyLitecoin,
		networkCode: soChainNetworkLitecoin + soChainNetworkTestSuffix,

		httpClient: getSoChainClient(),
	}
}
