package blockchainmod

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/shopspring/decimal"
)

const (
	// https://www.buybitcoinworldwide.com/fee-calculator/
	BitcoinTxnLegacySizeInput = 148
	BitcoinTxnSegWitSizeInput = 93.25
	BitcoinTxnSizeOutput      = 34
	BitcoinTxnSizeMeta        = 10

	BitcoinDecimalPlaces       = 8
	BitcoinTxnMinAmountSatoshi = 1000

	BitcoinTxnSingleOutputCount = 2
)

var (
	BitcoinChainConfig        = chaincfg.MainNetParams
	BitcoinChainConfigTestnet = chaincfg.TestNet3Params
)

var (
	BitcoinFeeMinimum = decimal.NewFromInt(1000) // in satoshi
)
