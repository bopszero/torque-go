package depositcrawler

import (
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type EthereumCrawler struct {
	BalanceLikeCrawler

	etherscanClient *blockchainmod.EtherscanEthereumClient
}

func NewEthereumCrawler(
	ctx comcontext.Context,
	coin blockchainmod.Coin, options CrawlerOptions,
) (_ *EthereumCrawler, err error) {
	var etherscanClient *blockchainmod.EtherscanEthereumClient
	if coin.IsUsingMainnet() {
		etherscanClient, err = blockchainmod.NewEthereumEtherscanMainnetSystemClient()
	} else {
		etherscanClient, err = blockchainmod.NewEthereumEtherscanTestRopstenSystemClient()
	}
	if err != nil {
		return
	}

	tokenMetaList := []blockchainmod.TokenMeta{
		blockchainmod.GetTokenMetaEthereumTetherUSD(),
	}
	balanceLikeCrawler := NewBalanceLikeCrawler(ctx, coin, tokenMetaList, options)
	client := EthereumCrawler{
		BalanceLikeCrawler: *balanceLikeCrawler,
		etherscanClient:    etherscanClient,
	}
	return &client, nil
}

func (this EthereumCrawler) ConsumeBlock(block blockchainmod.Block) error {
	if err := this.baseConsumeBlock(&this, block); err != nil {
		return err
	}

	internalTxns, err := this.etherscanClient.GetBlockInternalTxns(
		block.GetHeight(),
		meta.Paging{Limit: 10000})
	if err != nil {
		return err
	}
	for _, txn := range internalTxns {
		if err := this.ConsumerInternalTxn(block, txn); err != nil {
			return err
		}
	}

	return nil
}

func (this EthereumCrawler) ConsumerInternalTxn(
	block blockchainmod.Block,
	txn blockchainmod.EtherscanEthereumInternalTransaction,
) (err error) {
	if txn.AmountWei.IsZero() {
		return
	}
	if txn.HasError != 0 {
		return
	}
	if !isUserAddress(this.coin, txn.ToAddress) {
		return
	}

	ethAmount, err := currencymod.ConvertAmount(
		meta.CurrencyAmount{
			Currency: constants.CurrencySubEthereumWei,
			Value:    txn.AmountWei,
		},
		constants.CurrencyEthereum,
	)
	if err != nil {
		return
	}

	_, err = setCrawledDeposit(
		this.ctx, this.coin,
		txn.FromAddress, txn.ToAddress, 0,
		txn.Hash, ethAmount.Value,
		block.GetHeight(), txn.Time, 1,
	)
	return
}
