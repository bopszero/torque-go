package blockchainmod

import (
	"fmt"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type BlockchainInfoBitcoinBlock struct {
	Hash          string                             `json:"hash"`
	Height        uint64                             `json:"height"`
	IsMainChain   bool                               `json:"main_chain"`
	PrevBlockHash string                             `json:"prev_block"`
	Time          int64                              `json:"time"`
	Txns          []BlockchainInfoBitcoinTransaction `json:"tx"`
}

func (this *BlockchainInfoBitcoinBlock) GetCurrency() meta.Currency {
	return constants.CurrencyBitcoin
}

func (this *BlockchainInfoBitcoinBlock) GetNetwork() meta.BlockchainNetwork {
	if !this.IsMainChain {
		return constants.BlockchainNetworkBitcoinTestnet
	}
	return constants.BlockchainNetworkBitcoin
}

func (this *BlockchainInfoBitcoinBlock) GetHash() string {
	return this.Hash
}

func (this *BlockchainInfoBitcoinBlock) GetHeight() uint64 {
	return this.Height
}

func (this *BlockchainInfoBitcoinBlock) GetTimeUnix() int64 {
	return this.Time
}

func (this *BlockchainInfoBitcoinBlock) GetParentHash() string {
	return this.PrevBlockHash
}

func (this *BlockchainInfoBitcoinBlock) GetTransactions() ([]Transaction, error) {
	txns := make([]Transaction, len(this.Txns))
	for i := range this.Txns {
		txns[i] = &this.Txns[i]
	}
	return txns, nil
}

type BlockchainInfoBitcoinTransaction struct {
	baseTransaction

	cacheDirection   *meta.Direction
	cacheOwnerInputs []*BlockchainInfoBitcoinInput
	cacheFee         *meta.CurrencyAmount
	cacheAmount      *decimal.Decimal

	Hash        string `json:"hash"`
	Fee         uint64 `json:"fee"`
	BlockHeight uint64 `json:"block_height"`
	Time        int64  `json:"time"`

	Inputs  []BlockchainInfoBitcoinInput  `json:"inputs"`
	Outputs []BlockchainInfoBitcoinOutput `json:"out"`
}

func (this *BlockchainInfoBitcoinTransaction) GetCurrency() meta.Currency {
	return constants.CurrencyBitcoin
}

func (this *BlockchainInfoBitcoinTransaction) GetLocalStatus() meta.BlockchainTxnStatus {
	if this.GetConfirmations() > 0 {
		return constants.BlockchainTxnStatusSucceeded
	} else {
		return constants.BlockchainTxnStatusPending
	}
}

func (this *BlockchainInfoBitcoinTransaction) getDirection() meta.Direction {
	if this.ownerAddress == "" {
		return constants.DirectionTypeUnknown
	}

	for _, input := range this.Inputs {
		if input.GetPrevOutAddress() == this.ownerAddress {
			this.cacheOwnerInputs = append(this.cacheOwnerInputs, &input)
			return constants.DirectionTypeSend
		}
	}

	return constants.DirectionTypeReceive
}

func (this *BlockchainInfoBitcoinTransaction) GetDirection() meta.Direction {
	if this.cacheDirection != nil {
		return *this.cacheDirection
	}

	direction := this.getDirection()
	this.cacheDirection = &direction

	return direction
}

func (this *BlockchainInfoBitcoinTransaction) GetFee() meta.CurrencyAmount {
	if this.cacheFee != nil {
		return *this.cacheFee
	}

	feeValue := decimal.Zero
	for _, input := range this.Inputs {
		feeValue = feeValue.Add(input.GetPrevOutAmount())
	}
	for _, output := range this.Outputs {
		feeValue = feeValue.Sub(output.GetAmount())
	}

	fee := meta.CurrencyAmount{
		Currency: constants.CurrencyBitcoin,
		Value:    decimal.Max(feeValue, decimal.Zero),
	}
	this.cacheFee = &fee

	return fee
}

func (this *BlockchainInfoBitcoinTransaction) GetHash() string {
	return this.Hash
}

func (this *BlockchainInfoBitcoinTransaction) GetBlockHeight() uint64 {
	return this.BlockHeight
}

func (this *BlockchainInfoBitcoinTransaction) GetBlockHash() string {
	panic("blockchain.info hasn't support transaction block hash")
}

func (this *BlockchainInfoBitcoinTransaction) GetConfirmations() uint64 {
	if this.BlockHeight == 0 {
		return 0
	}

	networkInfo := GetCoinNativeF(this.GetCurrency()).GetModelNetwork()
	if networkInfo.LatestBlockHeight > this.BlockHeight {
		return networkInfo.LatestBlockHeight - this.BlockHeight + 1
	}

	return 1
}

func (this *BlockchainInfoBitcoinTransaction) GetFromAddress() (string, error) {
	if this.ownerAddress != "" && this.GetDirection() == constants.DirectionTypeSend {
		return this.ownerAddress, nil
	}

	if len(this.Inputs) > 0 {
		return this.Inputs[0].GetPrevOutAddress(), nil
	}

	return "", utils.IssueErrorf("cannot determine bitcoin transaction from address")
}

func (this *BlockchainInfoBitcoinTransaction) GetToAddress() (string, error) {
	if this.ownerAddress != "" && this.GetDirection() == constants.DirectionTypeReceive {
		return this.ownerAddress, nil
	}

	return "", utils.IssueErrorf("cannot determine bitcoin transaction to address")
}

func (this *BlockchainInfoBitcoinTransaction) GetAmount() (decimal.Decimal, error) {
	if this.ownerAddress == "" {
		return decimal.Zero, utils.IssueErrorf("cannot determine bitcoin doesn't have transaction amount")
	}

	if this.cacheAmount != nil {
		return *this.cacheAmount, nil
	}

	amount := decimal.Zero
	for _, input := range this.Inputs {
		if input.GetPrevOutAddress() == this.ownerAddress {
			amount = amount.Sub(input.GetPrevOutAmount())
		}
	}
	for _, output := range this.Outputs {
		if output.GetAddress() == this.ownerAddress {
			amount = amount.Add(output.GetAmount())
		}
	}

	amount = amount.Abs()
	this.cacheAmount = &amount

	return amount, nil
}

func (this *BlockchainInfoBitcoinTransaction) GetInputDataHex() (string, error) {
	return "", nil
}

func (this *BlockchainInfoBitcoinTransaction) GetTimeUnix() int64 {
	return this.Time
}

func (this *BlockchainInfoBitcoinTransaction) GetInputs() ([]Input, error) {
	inputs := make([]Input, len(this.Inputs))

	for i := range this.Inputs {
		inputs[i] = &this.Inputs[i]
	}

	return inputs, nil
}

func (this *BlockchainInfoBitcoinTransaction) GetOutputs() ([]Output, error) {
	outputs := make([]Output, len(this.Outputs))

	for i := range this.Outputs {
		outputs[i] = &this.Outputs[i]
	}

	return outputs, nil
}

type BlockchainInfoBitcoinInput struct {
	PrevOutput BlockchainInfoBitcoinOutput `json:"prev_out"`
}

func (this *BlockchainInfoBitcoinInput) GetPrevOutHash() string {
	panic(fmt.Errorf("blockchain.info hasn't support input prev outpoint hash"))
}

func (this *BlockchainInfoBitcoinInput) GetPrevOutIndex() uint32 {
	return this.PrevOutput.GetIndex()
}

func (this *BlockchainInfoBitcoinInput) GetPrevOutAddress() string {
	return this.PrevOutput.GetAddress()
}

func (this *BlockchainInfoBitcoinInput) GetPrevOutAmount() decimal.Decimal {
	return this.PrevOutput.GetAmount()
}

type BlockchainInfoBitcoinOutput struct {
	Index         uint32 `json:"n"`
	Spent         bool   `json:"spent"`
	AmountSatoshi int64  `json:"value"`
	Script        string `json:"script"`
	Address       string `json:"addr"`
}

func (this *BlockchainInfoBitcoinOutput) GetIndex() uint32 {
	return this.Index
}

func (this *BlockchainInfoBitcoinOutput) GetAddress() string {
	return this.Address
}

func (this *BlockchainInfoBitcoinOutput) GetAmount() decimal.Decimal {
	btcAmount := currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: constants.CurrencySubBitcoinSatoshi,
			Value:    decimal.NewFromInt(this.AmountSatoshi),
		},
		constants.CurrencyBitcoin,
	)

	return btcAmount.Value
}

type BlockchainInfoBitcoinUtxOutput struct {
	address string

	TxnHash       string `json:"tx_hash_big_endian"`
	Index         uint32 `json:"tx_output_n"`
	AmountSatoshi int64  `json:"value"`
	Script        string `json:"script"`
}

func (this *BlockchainInfoBitcoinUtxOutput) GetAddress() string {
	return this.address
}

func (this *BlockchainInfoBitcoinUtxOutput) GetTxnHash() string {
	return this.TxnHash
}

func (this *BlockchainInfoBitcoinUtxOutput) GetIndex() uint32 {
	return this.Index
}

func (this *BlockchainInfoBitcoinUtxOutput) GetAmount() decimal.Decimal {
	btcAmount := currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: constants.CurrencySubBitcoinSatoshi,
			Value:    decimal.NewFromInt(this.AmountSatoshi),
		},
		constants.CurrencyBitcoin,
	)

	return btcAmount.Value
}
