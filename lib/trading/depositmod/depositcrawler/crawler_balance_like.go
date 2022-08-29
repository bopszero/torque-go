package depositcrawler

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
)

type BalanceLikeCrawler struct {
	baseCrawler
}

func NewBalanceLikeCrawler(
	ctx comcontext.Context,
	coin blockchainmod.Coin, tokenMetaList []blockchainmod.TokenMeta,
	options CrawlerOptions,
) *BalanceLikeCrawler {
	return &BalanceLikeCrawler{newBaseCrawler(ctx, coin, tokenMetaList, options)}
}

func (BalanceLikeCrawler) baseConsumeBlock(that Crawler, block blockchainmod.Block) error {
	txns, err := block.GetTransactions()
	if err != nil {
		return err
	}
	for _, txn := range txns {
		if err := that.ConsumeTxn(txn); err != nil {
			return err
		}
	}

	return nil
}

func (this BalanceLikeCrawler) ConsumeBlock(block blockchainmod.Block) error {
	return this.baseConsumeBlock(&this, block)
}

func (this BalanceLikeCrawler) ConsumeTxn(txn blockchainmod.Transaction) (err error) {
	if err = this.consumeTxnNative(txn); err != nil {
		return
	}
	for _, tokenMeta := range this.tokenMetaList {
		rc20Transfers, err := txn.GetRC20Transfers(tokenMeta)
		if err != nil {
			return err
		}
		for _, rc20Transfer := range rc20Transfers {
			if err := this.consumeRC20Transfer(txn, tokenMeta, rc20Transfer); err != nil {
				return err
			}
		}
	}

	return nil
}

func (this BalanceLikeCrawler) consumeTxnNative(txn blockchainmod.Transaction) (err error) {
	fromAddress, toAddress, amount, ok, err := this.getTxnValues(txn)
	if err != nil || !ok {
		return
	}
	if !isUserAddress(this.coin, toAddress) {
		return
	}

	switch txn.GetLocalStatus() {
	case constants.BlockchainTxnStatusPending, constants.BlockchainTxnStatusSucceeded:
		_, err = setCrawledDeposit(
			this.ctx, this.coin,
			fromAddress, toAddress, 0,
			txn.GetHash(), amount, txn.GetBlockHeight(),
			txn.GetTimeUnix(), txn.GetConfirmations(),
		)
		break
	case constants.BlockchainTxnStatusFailed:
		_, err = unsetCrawledDeposit(
			this.ctx, this.coin,
			fromAddress, toAddress, 0,
			txn.GetHash(), amount, txn.GetBlockHeight(),
		)
		break
	default:
		break
	}
	return
}

func (this BalanceLikeCrawler) getTxnValues(txn blockchainmod.Transaction) (
	fromAddress, toAddress string, amount decimal.Decimal, ok bool, err error,
) {
	if amount, err = txn.GetAmount(); err != nil {
		return
	}
	if fromAddress, err = txn.GetFromAddress(); err != nil {
		return
	}
	if toAddress, err = txn.GetToAddress(); err != nil {
		return
	}
	if amount.IsZero() {
		return
	}
	if toAddress == "" {
		return
	}

	ok = true
	return
}

func (this BalanceLikeCrawler) consumeRC20Transfer(
	txn blockchainmod.Transaction,
	tokenMeta blockchainmod.TokenMeta, transfer blockchainmod.RC20Transfer,
) (err error) {
	var (
		fromAddress = transfer.GetFromAddress()
		toAddress   = transfer.GetToAddress()
		amount      = transfer.GetAmount()
	)
	if amount.IsZero() {
		return
	}
	if toAddress == "" {
		return
	}

	tokenCoin := blockchainmod.GetCoinF(tokenMeta.Currency, this.coin.GetNetwork())
	if !isUserAddress(tokenCoin, toAddress) {
		return
	}

	switch txn.GetLocalStatus() {
	case constants.BlockchainTxnStatusPending, constants.BlockchainTxnStatusSucceeded:
		_, err = setCrawledDeposit(
			this.ctx, tokenCoin,
			fromAddress, toAddress, 0,
			txn.GetHash(), amount, txn.GetBlockHeight(),
			txn.GetTimeUnix(), txn.GetConfirmations(),
		)
		break
	case constants.BlockchainTxnStatusFailed:
		_, err = unsetCrawledDeposit(
			this.ctx, tokenCoin,
			fromAddress, toAddress, 0,
			txn.GetHash(), amount, txn.GetBlockHeight(),
		)
		break
	default:
		break
	}
	return
}
