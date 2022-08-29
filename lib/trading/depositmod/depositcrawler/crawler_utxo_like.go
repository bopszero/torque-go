package depositcrawler

import (
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
)

type UtxoLikeCrawler struct {
	baseCrawler
}

func NewUtxoLikeLikeCrawler(
	ctx comcontext.Context,
	coin blockchainmod.Coin, options CrawlerOptions,
) *UtxoLikeCrawler {
	return &UtxoLikeCrawler{newBaseCrawler(ctx, coin, nil, options)}
}

func (this *UtxoLikeCrawler) ConsumeBlock(block blockchainmod.Block) error {
	txns, err := block.GetTransactions()
	if err != nil {
		return err
	}

	for _, txn := range txns {
		if err := this.ConsumeTxn(txn); err != nil {
			return err
		}
	}

	return nil
}

func (this *UtxoLikeCrawler) ConsumeTxn(txn blockchainmod.Transaction) error {
	inputs, err := txn.GetInputs()
	if err != nil {
		return err
	}
	outputs, err := txn.GetOutputs()
	if err != nil {
		return err
	}

	var (
		selectedFromAddress string
		inputAddress        = make(comtypes.HashSet, len(inputs))
	)
	for _, input := range inputs {
		address := input.GetPrevOutAddress()
		inputAddress.Add(address)
		if selectedFromAddress == "" && address != "" {
			selectedFromAddress = address
		}
	}
	if selectedFromAddress == "" {
		selectedFromAddress = "<Unknown>"
	}

	for i, output := range outputs {
		outputAddress := output.GetAddress()
		if inputAddress.Contains(outputAddress) {
			continue
		}

		if !isUserAddress(this.coin, outputAddress) {
			continue
		}

		_, err := setCrawledDeposit(
			this.ctx, this.coin,
			selectedFromAddress, outputAddress, uint16(i),
			txn.GetHash(), output.GetAmount(), txn.GetBlockHeight(),
			txn.GetTimeUnix(), txn.GetConfirmations(),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
