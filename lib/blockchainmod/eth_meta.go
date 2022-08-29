package blockchainmod

import (
	"math/big"

	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
)

const (
	EthereumDecimalPlaces    = 18
	EthereumMinConfirmations = 12
)

var (
	ChainIdEthereum                  = big.NewInt(1)
	ChainIdEthereumTestRopsten       = big.NewInt(3)
	ChainIdEthereumTestRinkeby       = big.NewInt(4)
	ChainIdEthereumTestKovan         = big.NewInt(42)
	ChainIdEthereumClassic           = big.NewInt(61)
	ChainIdEthereumClassicTestMorden = big.NewInt(62)

	NetworkIdEthereum                  = big.NewInt(1)
	NetworkIdEthereumTestRopsten       = big.NewInt(3)
	NetworkIdEthereumTestRinkeby       = big.NewInt(4)
	NetworkIdEthereumTestKovan         = big.NewInt(42)
	NetworkIdEthereumClassic           = big.NewInt(1)
	NetworkIdEthereumClassicTestMorden = big.NewInt(2)
)

const (
	// transfer(address,uint256)
	EthereumERC20TransferTokenMethodIdHex = "0xa9059cbb"
	// Transfer(address,address,uint256)
	EthereumERC20TrasnferTokenLogTopicHex = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
)

const (
	EthereumStandardGasLimit     = 21000
	EthereumTokenUsdtGasLimitMin = 60000
)

var (
	EthereumTokenMetaMainnetTetherUSD = TokenMeta{
		Currency:      constants.CurrencyTetherUSD,
		Address:       "0xdac17f958d2ee523a2206206994597c13d831ec7",
		DecimalPlaces: 6,
	}
	EthereumTokenMetaTestRopstenTetherUSD = TokenMeta{
		Currency:      constants.CurrencyTetherUSD,
		Address:       "0x6ab27024eaa44896740311bb663928de056709c2",
		DecimalPlaces: 6,
	}
)

const (
	EtherscanModuleAccount = "account"
	EtherscanModuleProxy   = "proxy"

	EtherscanActionGetBalance           = "balance"
	EtherscanActionGetTokenBalance      = "tokenbalance"
	EtherscanActionListTxns             = "txlist"
	EtherscanActionListInternalTxns     = "txlistinternal"
	EtherscanActionListErc20TokenTxns   = "tokentx"
	EtherscanActionGetBlockHeight       = "eth_blockNumber"
	EtherscanActionGetBlockByHeight     = "eth_getBlockByNumber"
	EtherscanActionGetCode              = "eth_getCode"
	EtherscanActionGetLatestBlockHeight = "eth_blockNumber"
	EtherscanActionGetTxn               = "eth_getTransactionByHash"
	EtherscanActionGetTxnReceipt        = "eth_getTransactionReceipt"
	EtherscanActionGetTxnsCount         = "eth_getTransactionCount"
	EtherscanActionPushTxn              = "eth_sendRawTransaction"

	EtherscanStatusSuccess = "1"
)
