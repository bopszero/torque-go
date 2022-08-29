package blockchainmod

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
)

const (
	TronDecimalPlaces           = 6
	TronMinConfirmations        = 12
	TronTxnLifetime             = 12 * time.Hour
	TronTxnErrorRangeTime       = 1 * time.Second
	TronTxnErrorRangeExpireTime = 5 * time.Minute
)

var (
	TronFeeNormalTxnPrice = decimal.NewFromInt(10)
	TronFeeMaximum        = decimal.NewFromInt(100000) // in sun

	// transfer(address,uint256)
	TronTRC20TransferTokenMethodIdHex = "a9059cbb"
	TronTRC20TransferTokenMethodType  = "TRC20"
	TronTRC20TransferTokenMethodName  = "transfer"
)

var (
	TronTokenMetaMainnetTetherUSD = TokenMeta{
		Currency:      constants.CurrencyTetherUSD,
		Address:       "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
		DecimalPlaces: 6,
	}
	TronTokenMetaTestShastaTetherUSD = TokenMeta{
		Currency:      constants.CurrencyTetherUSD,
		Address:       "TXM3Y5L4hi39NbgnrSL9HqfXNNR38MynVG",
		DecimalPlaces: 6,
	}
)

const (
	TronTxnStatusSuccess             = "success"
	TronTxnTypeTransferContract      = "TransferContract"
	TronTxnTypeTransferAssetContract = "TransferAssetContract"
	TronTxnTypeTriggerSmartContract  = "TriggerSmartContract"
)

// Reference: https://developers.tron.network/docs/tron-grid-intro
const (
	TronHostTestnetRpcMain   = "grpc.trongrid.io:50051"
	TronHostTestnetRpcShasta = "grpc.shasta.trongrid.io:50051"
	TronHostTestnetRpcNile   = "47.252.19.181:50051"
)
