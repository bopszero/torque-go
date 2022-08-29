package test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

func blockchainScanBitcoinTxn(t *testing.T, txn blockchainmod.Transaction) {
	txn.GetConfirmations()

	txn.GetHash()
	// txn.GetBlockHash()
	txn.GetLocalStatus()

	txn.GetTimeUnix()

	inputs, err := txn.GetInputs()
	assert.Nil(t, err)
	for _, input := range inputs {
		// input.GetPrevOutHash()
		input.GetPrevOutIndex()

		// input.GetPrevOutAddress()
		// input.GetPrevOutAmount()
	}

	outputs, err := txn.GetOutputs()
	assert.Nil(t, err)
	for _, output := range outputs {
		output.GetIndex()
		output.GetAddress()
		output.GetAmount()
	}
}

func testBlockchainBitcoinLikeClient(t *testing.T, client blockchainmod.Client, blockHash string, address string) {
	var err error

	balance, err := client.GetBalance(address)
	assert.Nil(t, err)
	balance.Add(decimal.Zero)

	block, err := client.GetBlock(blockHash)
	assert.Nil(t, err)
	block.GetCurrency()
	block.GetNetwork()
	block.GetHash()
	block.GetHeight()
	block.GetParentHash()
	blockTxns, err := block.GetTransactions()
	assert.Nil(t, err)
	for _, txn := range blockTxns {
		blockchainScanBitcoinTxn(t, txn)
	}

	txns, err := client.GetTxns(address, meta.Paging{Limit: 10})
	assert.Nil(t, err)
	for _, txn := range txns {
		blockchainScanBitcoinTxn(t, txn)
	}

	txn, err := client.GetTxn(txns[0].GetHash())
	assert.Nil(t, err)
	blockchainScanBitcoinTxn(t, txn)
}

func TestBlockchainBitcoinBitcoreMainnetClient(t *testing.T) {
	client := blockchainmod.NewBitcoinBitcoreMainnetClient()

	address := "1AWLm7LqhfUzWfdyhEMLiAFavgB1LRXvm5"

	_, err := client.GetBalance(address)
	assert.Nil(t, err)

	block, err := client.GetLatestBlock()
	assert.Nil(t, err)
	assert.True(t, block.GetHeight() > 0)

	txns, err := client.GetTxns(address, meta.Paging{Limit: 10})
	assert.Nil(t, err)
	assert.Greater(t, len(txns), 0)

	utxOutputs, err := client.GetUtxOutputs(address, comutils.DecimalTen)
	assert.Nil(t, err)
	assert.Greater(t, len(utxOutputs), 0)
}

func TestBlockchainBitcoinTestnetTxn(t *testing.T) {
	var (
		keyFrom     = "cTsFgBzQhZd1HqS7MiU1SwaRBNjNajMZ2P2t9sDe9ifzecoRHuZs"
		addressFrom = "n2LbMGPvpXoKYhEbhp8ciE9XzwPHyhcKzs"
		keyTo       = "cNRKedbtiuqrkxAFRZjgkJMauYTBgYbdxV7qM3DErYWzpHxZWHba"
		addressTo   = "tb1qaczy4s8ayaqs5jph5xc0gfaegly7s9c7uyr77h"
		amount      = decimal.NewFromFloat(0.001)
	)

	coin := blockchainmod.GetCoinNativeF(constants.CurrencyBitcoin)

	fromAcc, err := coin.NewAccountOwner(keyFrom, addressFrom)
	assert.Nil(t, err)
	toAcc, err := coin.NewAccountOwner(keyTo, addressTo)
	assert.Nil(t, err)

	client, err := coin.NewClientDefault()
	assert.Nil(t, err)

	feeInfo, err := fromAcc.GetFeeInfoToAddress(toAcc.GetAddress())
	assert.Nil(t, err)
	feeInfo.Price = comutils.DecimalOne
	feeInfo.SetLimitMaxQuantity(feeInfo.LimitMaxQuantity)

	txn, err := fromAcc.GenerateSignedTxn(toAcc.GetAddress(), amount, &feeInfo)
	assert.Nil(t, err)

	t.Logf("txn hash: %s", txn.GetHash())
	txnBytes := txn.GetRaw()
	t.Logf("txn hex: %s", comutils.HexEncode(txn.GetRaw()))
	err = client.PushTxnRaw(txnBytes)
	assert.Nil(t, err)
}

func TestBlockchainTronMainnetListTokenTxns(t *testing.T) {
	var (
		fromAddress = "TYKXLBFQjwFDYNGoNs9wzpgxiS1YtL3bFx"
	)

	coin := blockchainmod.GetCoinF(constants.CurrencyTetherUSD, constants.BlockchainNetworkTron)
	client, err := coin.NewClientDefault()
	assert.Nil(t, err)

	txns, err := client.GetTxns(fromAddress, meta.Paging{Limit: 3})
	assert.Nil(t, err)

	for _, txn := range txns {
		txnAmount, err := txn.GetAmount()
		assert.Nil(t, err)
		assert.True(t, txnAmount.GreaterThan(decimal.Zero))
	}
}

func TestBlockchainTronTestShastaTxn(t *testing.T) {
	// From Address: TAWdc5J6ANwpETh5f43E2ZEdrtrQSZsAFE
	// To Address:   TT5sMopof2Dd2gXdPSfehSHMvp5k6MuaZA
	var (
		keyFrom = "add80c549e2e25af1e2ac0900b9c3ad759c2a9eaea9d6b2a4cf03a265790b6a8"
		// toKey   = "84c374f14ac6aef769be05f75ddcc8fb2348e64954323b73c12267b2511fcb9b"
		toAddress = "TT5sMopof2Dd2gXdPSfehSHMvp5k6MuaZA"
		amount    = decimal.NewFromFloat(1)
	)

	coin := blockchainmod.GetCoinNativeF(constants.CurrencyTron)

	fromAcc, err := coin.NewAccountOwner(keyFrom, "")
	assert.Nil(t, err)
	toAcc, err := coin.NewAccountGuest(toAddress)
	assert.Nil(t, err)

	client, err := coin.NewClientDefault()
	assert.Nil(t, err)

	txn, err := fromAcc.GenerateSignedTxn(toAcc.GetAddress(), amount, nil)
	assert.Nil(t, err)

	t.Logf("txn hash: %s", txn.GetHash())
	txnBytes := txn.GetRaw()
	t.Logf("txn hex: %s", comutils.HexEncode(txn.GetRaw()))
	err = client.PushTxnRaw(txnBytes)
	assert.Nil(t, err)
}

func TestBlockchainTronTestShastaTokenTxn(t *testing.T) {
	// From Address: TAWdc5J6ANwpETh5f43E2ZEdrtrQSZsAFE
	// To Address:   TT5sMopof2Dd2gXdPSfehSHMvp5k6MuaZA
	var (
		keyFrom = "add80c549e2e25af1e2ac0900b9c3ad759c2a9eaea9d6b2a4cf03a265790b6a8"
		keyTo   = "84c374f14ac6aef769be05f75ddcc8fb2348e64954323b73c12267b2511fcb9b"
		amount  = decimal.NewFromFloat(1.765432)
	)

	coin := blockchainmod.GetCoinF(constants.CurrencyTetherUSD, constants.BlockchainNetworkTronTestShasta)

	fromAcc, err := coin.NewAccountOwner(keyFrom, "")
	assert.Nil(t, err)
	toAcc, err := coin.NewAccountOwner(keyTo, "")
	assert.Nil(t, err)

	client, err := coin.NewClientDefault()
	assert.Nil(t, err)

	txn, err := fromAcc.GenerateSignedTxn(toAcc.GetAddress(), amount, nil)
	assert.Nil(t, err)

	t.Logf("txn hash: %s", txn.GetHash())
	txnBytes := txn.GetRaw()
	t.Logf("txn hex: %s", comutils.HexEncode(txn.GetRaw()))
	err = client.PushTxnRaw(txnBytes)
	assert.Nil(t, err)
}

func TestBlockchainRippleTestnetTxn(t *testing.T) {
	var (
		fromSecret = "sEd7tch9nBLos4rXxpAA3arsqbnfabj"
		// fromLegacyAddress = "rUC3dM9dZdhd1nUzQMZu4jEDoZop8Dcc1f"
		fromXAddress = "TVshSSsodvqkAGT7yUkZgppKWbmsmYfteQUhNh2EgUzadeu"

		// toSecret = "sEd76f3MRGyXTkWkM3hFpaCXU27PR1d"
		// toLegacyAddress = "rJkX8FQoUdtTkoibS3cfdZgjyXt6dxzVAc"
		toXAddress = "TV4vr4JTDiwDncukQTm3NxsUTpHCC1oM8ikcyc11QJrLVBh"

		amount = decimal.NewFromFloat(1.1)
	)

	coin := blockchainmod.GetCoinF(constants.CurrencyRipple, constants.BlockchainNetworkRippleTestnet)

	fromAcc, err := coin.NewAccountOwner(fromSecret, fromXAddress)
	assert.Nil(t, err)
	toAcc, err := coin.NewAccountGuest(toXAddress)
	assert.Nil(t, err)

	client, err := coin.NewClientDefault()
	assert.Nil(t, err)

	txn, err := fromAcc.GenerateSignedTxn(toAcc.GetAddress(), amount, nil)
	assert.Nil(t, err)

	t.Logf("txn hash: %s", txn.GetHash())
	txnBytes := txn.GetRaw()
	t.Logf("txn hex: %s", comutils.HexEncode(txn.GetRaw()))
	err = client.PushTxnRaw(txnBytes)
	assert.Nil(t, err)
}
