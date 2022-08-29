package blockchainmod

import (
	chaincfgBTC "github.com/btcsuite/btcd/chaincfg"
	"github.com/jinzhu/copier"
	chaincfgLTC "github.com/ltcsuite/ltcd/chaincfg"
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comutils"
)

const (
	LitecoinTxnMinAmountSatoshi = 100000
)

var (
	LitecoinChainConfig        chaincfgBTC.Params
	LitecoinChainConfigTestnet chaincfgBTC.Params
)

var (
	LitecoinFeeMinimum = decimal.NewFromInt(10000)    // in satoshi
	LitecoinFeeMaximum = decimal.NewFromInt(10000000) // in satoshi
)

func init() {
	comutils.PanicOnError(
		copier.Copy(&LitecoinChainConfig, &chaincfgLTC.MainNetParams),
	)
	if err := chaincfgBTC.Register(&LitecoinChainConfig); err != nil {
		panic(err)
	}

	comutils.PanicOnError(
		copier.Copy(&LitecoinChainConfigTestnet, &chaincfgLTC.TestNet4Params),
	)
	if err := chaincfgBTC.Register(&LitecoinChainConfigTestnet); err != nil {
		panic(err)
	}
}
