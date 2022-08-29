package constants

import (
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	BlockchainNetworkTorque                = meta.BlockchainNetwork("TORQ")
	BlockchainNetworkBitcoin               = meta.BlockchainNetwork("BTC")
	BlockchainNetworkBitcoinTestnet        = meta.BlockchainNetwork("BTC:TESTNET")
	BlockchainNetworkBitcoinCash           = meta.BlockchainNetwork("BCH")
	BlockchainNetworkBitcoinCashTestnet    = meta.BlockchainNetwork("BCH:TESTNET")
	BlockchainNetworkBitcoinCashABC        = meta.BlockchainNetwork("BCHA")
	BlockchainNetworkBitcoinCashAbcTestnet = meta.BlockchainNetwork("BCHA:TESTNET")
	BlockchainNetworkEthereum              = meta.BlockchainNetwork("ETH")
	BlockchainNetworkEthereumTestRopsten   = meta.BlockchainNetwork("ETH:TEST_ROPSTEN")
	BlockchainNetworkLitecoin              = meta.BlockchainNetwork("LTC")
	BlockchainNetworkLitecoinTestnet       = meta.BlockchainNetwork("LTC:TESTNET")
	BlockchainNetworkTron                  = meta.BlockchainNetwork("TRX")
	BlockchainNetworkTronTestShasta        = meta.BlockchainNetwork("TRX:TEST_SHASTA")
	BlockchainNetworkRipple                = meta.BlockchainNetwork("XRP")
	BlockchainNetworkRippleTestnet         = meta.BlockchainNetwork("XRP:TESTNET")
)

const (
	BlockchainTxnStatusFailed    = meta.BlockchainTxnStatus(-1)
	BlockchainTxnStatusPending   = meta.BlockchainTxnStatus(1)
	BlockchainTxnStatusSucceeded = meta.BlockchainTxnStatus(10)
)

var (
	BlockchainChannelTypes = []meta.ChannelType{
		ChannelTypeDstBlockchainNetwork,
		ChannelTypeDstTorqueConvert,
	}
	BlockchainChannelUtxoCurrencySet = comtypes.NewHashSetFromListF([]meta.Currency{
		CurrencyBitcoin,
		CurrencyBitcoinCash,
		CurrencyBitcoinCashABC,
		CurrencyLitecoin,
	})
)
