package blockchainmod

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type BitcoreUtxoLikeBlock struct {
	tBitcoreBaseBlock
}

type BitcoreUtxoLikeTransaction struct {
	tBitcoreBaseTransaction
	cacheDirection *meta.Direction
	cacheAmount    *decimal.Decimal

	Inputs  []BitcoreUtxoLikeTxnCoinFlow `json:"inputs"`
	Outputs []BitcoreUtxoLikeTxnCoinFlow `json:"outputs"`
}

func (this *BitcoreUtxoLikeTransaction) GetDirection() meta.Direction {
	if this.cacheDirection != nil {
		return *this.cacheDirection
	}

	direction := this.getDirection()
	this.cacheDirection = &direction

	return direction
}

func (this *BitcoreUtxoLikeTransaction) getDirection() meta.Direction {
	if this.ownerAddress == "" {
		return constants.DirectionTypeUnknown
	}

	for _, input := range this.Inputs {
		if input.GetPrevOutAddress() == this.ownerAddress {
			return constants.DirectionTypeSend
		}
	}

	return constants.DirectionTypeReceive
}

func (this *BitcoreUtxoLikeTransaction) GetFromAddress() (string, error) {
	if this.GetDirection() == constants.DirectionTypeSend {
		return this.ownerAddress, nil
	}

	return "", utils.IssueErrorf("cannot determine bitcoin transaction from address")
}

func (this *BitcoreUtxoLikeTransaction) GetToAddress() (string, error) {
	if this.GetDirection() == constants.DirectionTypeReceive {
		return this.ownerAddress, nil
	}

	return "", utils.IssueErrorf("cannot determine bitcoin transaction to address")
}

func (this *BitcoreUtxoLikeTransaction) GetAmount() (decimal.Decimal, error) {
	if this.cacheAmount != nil {
		return *this.cacheAmount, nil
	}
	if this.ownerAddress == "" {
		return decimal.Zero, utils.IssueErrorf("cannot determine bitcoin transaction amount without owner address")
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

func (this *BitcoreUtxoLikeTransaction) GetInputDataHex() (string, error) {
	return "", nil
}

func (this *BitcoreUtxoLikeTransaction) GetInputs() ([]Input, error) {
	inputs := make([]Input, len(this.Inputs))
	for i := range this.Inputs {
		inputs[i] = &this.Inputs[i]
	}

	return inputs, nil
}

func (this *BitcoreUtxoLikeTransaction) GetOutputs() ([]Output, error) {
	outputs := make([]Output, len(this.Outputs))
	for i := range this.Outputs {
		outputs[i] = &this.Outputs[i]
	}

	return outputs, nil
}

type BitcoreUtxoLikeTxnCoinFlow struct {
	Currency       meta.Currency   `json:"chain"`
	Network        string          `json:"network"`
	Address        string          `json:"address"`
	Script         string          `json:"script"`
	Value          decimal.Decimal `json:"value"`
	Confirmations  int64           `json:"confirmations"`
	InIndex        uint32          `json:"mintIndex"`
	InTxnHash      string          `json:"mintTxid"`
	InBlockHeight  int64           `json:"mintHeight"`
	OutTxnHash     string          `json:"spentTxid"`
	OutBlockHeight int64           `json:"spentHeight"`
}

func (this *BitcoreUtxoLikeTxnCoinFlow) GetInTxnHash() string {
	return this.InTxnHash
}

func (this *BitcoreUtxoLikeTxnCoinFlow) GetOutTxnHash() string {
	return this.OutTxnHash
}

func (this *BitcoreUtxoLikeTxnCoinFlow) GetIndex() uint32 {
	return this.InIndex
}

func (this *BitcoreUtxoLikeTxnCoinFlow) GetAddress() string {
	return this.Address
}

func (this *BitcoreUtxoLikeTxnCoinFlow) GetAmount() decimal.Decimal {
	return currencymod.
		ConvertAmountF(
			meta.CurrencyAmount{
				Currency: constants.CurrencySubBitcoinSatoshi,
				Value:    this.Value,
			},
			this.Currency,
		).
		Value
}

func (this *BitcoreUtxoLikeTxnCoinFlow) GetPrevOutHash() string {
	return this.InTxnHash
}

func (this *BitcoreUtxoLikeTxnCoinFlow) GetPrevOutIndex() uint32 {
	return this.InIndex
}

func (this *BitcoreUtxoLikeTxnCoinFlow) GetPrevOutAddress() string {
	return this.Address
}

func (this *BitcoreUtxoLikeTxnCoinFlow) GetPrevOutAmount() decimal.Decimal {
	return this.GetAmount()
}

type BitcoreUtxoLikeTransactionUtxOutput struct {
	BitcoreUtxoLikeTxnCoinFlow
}

func (this *BitcoreUtxoLikeTransactionUtxOutput) GetTxnHash() string {
	return this.GetInTxnHash()
}
