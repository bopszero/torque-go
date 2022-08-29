package depositcrawler

import (
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
)

type RippleCrawler struct {
	BalanceLikeCrawler
}

func NewRippleCrawler(
	ctx comcontext.Context,
	coin blockchainmod.Coin, options CrawlerOptions,
) *RippleCrawler {
	balanceLikeCrawler := NewBalanceLikeCrawler(ctx, coin, nil, options)
	return &RippleCrawler{*balanceLikeCrawler}
}

func (this RippleCrawler) ConsumeBlock(block blockchainmod.Block) error {
	return this.baseConsumeBlock(&this, block)
}

func (this RippleCrawler) consumeCoinTxn(
	coin blockchainmod.Coin, txn blockchainmod.Transaction,
) (err error) {
	fromAddress, toAddress, amount, ok, err := this.getTxnValues(txn)
	if err != nil || !ok {
		return
	}
	toXAddress, err := blockchainmod.RippleParseXAddress(toAddress, coin.IsUsingMainnet())
	if err != nil {
		return
	}
	if !isUserAddress(coin, toXAddress.GetAddress()) {
		return
	}

	_, err = setCrawledDeposit(
		this.ctx, coin,
		fromAddress, toXAddress.GetAddress(), 0,
		txn.GetHash(), amount, txn.GetBlockHeight(),
		txn.GetTimeUnix(), txn.GetConfirmations(),
	)
	return
}

func (this RippleCrawler) ConsumeTxn(txn blockchainmod.Transaction) (err error) {
	return this.consumeCoinTxn(this.coin, txn)
}
