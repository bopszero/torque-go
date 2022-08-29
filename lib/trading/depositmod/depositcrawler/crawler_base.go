package depositcrawler

import (
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type baseCrawler struct {
	ctx           comcontext.Context
	coin          blockchainmod.Coin
	tokenMetaList []blockchainmod.TokenMeta

	options CrawlerOptions
}

func newBaseCrawler(
	ctx comcontext.Context,
	coin blockchainmod.Coin, tokenMetaList []blockchainmod.TokenMeta,
	options CrawlerOptions,
) baseCrawler {
	return baseCrawler{
		ctx:           ctx,
		coin:          coin,
		tokenMetaList: tokenMetaList,
		options:       options,
	}
}

func (this *baseCrawler) GetSystemCrawledBlockHeight() (_ uint64, err error) {
	var currencyInfo models.LegacyCurrencyInfo
	err = database.GetDbSlave().
		First(
			&currencyInfo,
			&models.LegacyCurrencyInfo{Currency: this.coin.GetCurrency()},
		).
		Error
	if err != nil {
		err = utils.WrapError(err)
		return
	}

	return currencyInfo.LatestCrawledBlockHeight, nil
}

func (this *baseCrawler) GetScanBlockHeights() (heights []uint64, err error) {
	systemBlockHeight, err := this.GetSystemCrawledBlockHeight()
	if err != nil {
		return
	}

	return this.options.GenBlockHeights(systemBlockHeight), nil
}
