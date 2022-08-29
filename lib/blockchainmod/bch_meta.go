package blockchainmod

import (
	"github.com/gcash/bchd/chaincfg"
	"github.com/shopspring/decimal"
)

var (
	BitcoinCashChainConfig        = chaincfg.MainNetParams
	BitcoinCashChainConfigTestnet = chaincfg.TestNet3Params
)

var (
	BitcoinCashFeeMinimum = decimal.NewFromInt(1000) // in satoshi
)
