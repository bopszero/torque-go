package blockchainmod

import (
	"time"

	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type TokenMeta struct {
	Currency      meta.Currency
	Address       string
	DecimalPlaces uint8
}

const (
	Base58BitcoinAlphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	Base58IpfsAlphabet    = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	Base58FlickrAlphabet  = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
	Base58RippleAlphabet  = "rpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65jkm8oFqi1tuvAxyz"
)

const (
	MaxBalance = 1000000000000
)

const (
	ApiDefaultListPaging = 50

	ClientRequestTimeout = 5 * time.Second
)

const (
	HostSnapNodeProduction          = "http://node-api.myblockchain.live"
	HostSnapNodeTest                = "https://103.231.191.13/"
	HostMyCoinGateBitcoreProduction = "https://mycoingate.com"
	HostMyCoinGateBitcoreTest       = "https://staging.mycoingate.com"

	HostBlockchainInfoProduction          = "https://blockchain.info"
	HostEtherscanProduction               = "https://api.etherscan.io"
	HostEtherscanSandboxRopsten           = "https://api-ropsten.etherscan.io"
	HostInsightBitcoreBitcoinProduction   = "https://insight.bitpay.com"
	HostInsightLitecoreLitecoinProduction = "https://insight.litecore.io"
	HostSoChainProduction                 = "https://sochain.com"
)

const (
	ExplorerBitcoinMainnetTxnUrlPattern        = "https://blockchair.com/bitcoin/transaction/%s/"
	ExplorerBitcoinTestnetTxnUrlPattern        = "https://blockchair.com/bitcoin/testnet/transaction/%s/"
	ExplorerBitcoinCashMainnetTxnUrlPattern    = "https://explorer.bitcoin.com/bch/tx/%s/"
	ExplorerBitcoinCashTestnetTxnUrlPattern    = "https://explorer.bitcoin.com/tbch/tx/%s/"
	ExplorerBitcoinCashAbcMainnetTxnUrlPattern = "https://explorer.bitcoinabc.org/tx/%s/"
	ExplorerBitcoinCashAbcTestnetTxnUrlPattern = "https://texplorer.bitcoinabc.org//tx/%s/"
	ExplorerLitecoinMainnetTxnUrlPattern       = "https://sochain.com/tx/LTC/%s/"
	ExplorerLitecoinTestnetTxnUrlPattern       = "https://sochain.com/tx/LTCTEST/%s/"
	ExplorerEthereumMainnetTxnUrlPattern       = "https://etherscan.io/tx/%s/"
	ExplorerEthereumTestRopstenTxnUrlPattern   = "https://ropsten.etherscan.io/tx/%s/"
	ExplorerTronMainnetTxnUrlPattern           = "https://tronscan.org/#/transaction/%s/"
	ExplorerTronTestShastaTxnUrlPattern        = "https://shasta.tronscan.org/#/transaction/%s/"
	ExplorerTronTestNileTxnUrlPattern          = "https://nile.tronscan.org/#/transaction/%s/"
	ExplorerRippleMainnetTxnUrlPattern         = "https://xrpscan.com/tx/%s/"
	ExplorerRippleTestnetTxnUrlPattern         = "https://test.bithomp.com/explorer/%s"
)

const (
	ClientModeDefault = "default"
	ClientModeSpare   = "spare"
)

const (
	SoChainStatusSuccess = "success"
)

const (
	MyCoinGateAuthUsername = "mycoingate"
	MyCoinGateAuthPassword = "bYvnXehz7Vlg1Bl5NbrGxUXo3wb2U5Fl"
)

const (
	BitcoreBlockHeightAsPending  = -1
	BitcoreBlockHeightAsUnspent  = -2
	BitcoreBlockHeightAsConflict = -3
	BitcoreBlockHeightAsError    = -4

	BitcoreChainMainnet     = "mainnet"
	BitcoreChainTestnet     = "testnet"
	BitcoreChainTestRopsten = "ropsten"
	BitcoreChainTestShasta  = "shasta"
	BitcoreChainTestNile    = "nile"
)
