package depositcrawler

import (
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
)

func NewCrawler(
	ctx comcontext.Context,
	coin blockchainmod.Coin, options CrawlerOptions,
) (Crawler, error) {
	currency := coin.GetNetworkCurrency()
	switch currency {
	case constants.CurrencyEthereum:
		return NewEthereumCrawler(ctx, coin, options)
	case constants.CurrencyTron:
		return NewTronCrawler(ctx, coin, options), nil
	case constants.CurrencyRipple:
		return NewRippleCrawler(ctx, coin, options), nil
	default:
		if constants.BlockchainChannelUtxoCurrencySet.Contains(currency) {
			return NewUtxoLikeLikeCrawler(ctx, coin, options), nil
		} else {
			return NewBalanceLikeCrawler(ctx, coin, nil, options), nil
		}
	}
}

func NewCrawlerDefaultOptions(ctx comcontext.Context, coin blockchainmod.Coin) (Crawler, error) {
	var (
		options  CrawlerOptions
		currency = coin.GetNetworkCurrency()
	)
	switch currency {
	case constants.CurrencyEthereum:
		options = CrawlerOptions{
			BlockScanBackward: 30,
			BlockScanForward:  300,
		}
		break
	case constants.CurrencyTron:
		options = CrawlerOptions{
			BlockScanBackward: 50,
			BlockScanForward:  500,
		}
		break
	case constants.CurrencyRipple:
		options = CrawlerOptions{
			BlockScanBackward: 50,
			BlockScanForward:  500,
		}
		break
	default:
		if constants.BlockchainChannelUtxoCurrencySet.Contains(currency) {
			options = CrawlerOptions{
				BlockScanBackward: 2,
				BlockScanForward:  50,
			}
		} else {
			options = CrawlerOptions{
				BlockScanBackward: 20,
				BlockScanForward:  200,
			}
		}
	}

	return NewCrawler(ctx, coin, options)
}
