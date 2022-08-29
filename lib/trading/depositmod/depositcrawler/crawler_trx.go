package depositcrawler

import (
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
)

var TronAcceptedTxnTypes = comtypes.NewHashSetFromListF([]string{
	blockchainmod.TronTxnTypeTransferContract,
	blockchainmod.TronTxnTypeTriggerSmartContract,
})

type TronCrawler struct {
	BalanceLikeCrawler
}

func NewTronCrawler(
	ctx comcontext.Context,
	coin blockchainmod.Coin, options CrawlerOptions,
) *TronCrawler {
	tokenMetaList := []blockchainmod.TokenMeta{
		// blockchainmod.GetTokenMetaTronTetherUSD(),  // TODO: Enable when needed
	}
	balanceLikeCrawler := NewBalanceLikeCrawler(ctx, coin, tokenMetaList, options)
	return &TronCrawler{*balanceLikeCrawler}
}

func (this TronCrawler) ConsumeBlock(block blockchainmod.Block) error {
	return this.baseConsumeBlock(&this, block)
}

func (this TronCrawler) ConsumeTxn(txn blockchainmod.Transaction) (err error) {
	if !TronAcceptedTxnTypes.Contains(txn.GetTypeCode()) {
		return
	}
	return this.BalanceLikeCrawler.ConsumeTxn(txn)
}
