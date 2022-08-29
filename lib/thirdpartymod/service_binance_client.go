package thirdpartymod

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

var (
	binanceSystemClient *BinanceClient
)

type BinanceClient struct {
	*binance.Client

	RequestTimeout time.Duration
}

func NewBinanceClient(key string, secret string, timeout time.Duration) *BinanceClient {
	return &BinanceClient{
		Client:         binance.NewClient(key, secret),
		RequestTimeout: timeout,
	}
}

func GetBinanceSystemClient() *BinanceClient {
	if binanceSystemClient == nil {
		binanceSystemClient = NewBinanceClient(
			viper.GetString(config.KeyApiBinanceKey),
			viper.GetString(config.KeyApiBinanceSecret),
			5*time.Second,
		)
	}

	return binanceSystemClient
}

func (this *BinanceClient) getTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), this.RequestTimeout)
}

func (this *BinanceClient) getPairSymbolUsdt(currency meta.Currency) string {
	return fmt.Sprintf("%v%v", currency, constants.CurrencyTetherUSD)
}

func (this *BinanceClient) GetExchangeInfo() (*binance.ExchangeInfo, error) {
	ctx, cancel := this.getTimeoutContext()
	defer cancel()
	return this.NewExchangeInfoService().Do(ctx)
}

func (this *BinanceClient) GetAccountInfo() (*binance.Account, error) {
	ctx, cancel := this.getTimeoutContext()
	defer cancel()
	return this.NewGetAccountService().Do(ctx)
}

func (this *BinanceClient) GetAvgPriceUSDT(currency meta.Currency) (*binance.AvgPrice, error) {
	avgPriceService := this.NewAveragePriceService().
		Symbol(this.getPairSymbolUsdt(currency))
	ctx, cancel := this.getTimeoutContext()
	defer cancel()
	return avgPriceService.Do(ctx)
}

func (this *BinanceClient) GetOrderByRef(currency meta.Currency, ref string) (*binance.Order, error) {
	getOrderService := this.NewGetOrderService().
		Symbol(this.getPairSymbolUsdt(currency)).
		OrigClientOrderID(ref)
	ctx, cancel := this.getTimeoutContext()
	defer cancel()
	return getOrderService.Do(ctx)
}

func (this *BinanceClient) SubmitOrderUsdtMarket(
	currency meta.Currency, sideType binance.SideType,
	amount decimal.Decimal, reference string,
) (*binance.CreateOrderResponse, error) {
	orderService := this.NewCreateOrderService().
		Side(sideType).
		Type(binance.OrderTypeMarket).
		Symbol(this.getPairSymbolUsdt(currency)).
		Quantity(amount.String()).
		NewOrderRespType(binance.NewOrderRespTypeRESULT)
	if reference != "" {
		orderService.NewClientOrderID(reference)
	}

	ctx, cancel := this.getTimeoutContext()
	defer cancel()

	if config.Debug {
		testErr := orderService.Test(
			ctx,
			binance.WithRecvWindow(10000),
		)
		return nil, testErr
	}

	return orderService.Do(
		ctx,
		binance.WithRecvWindow(10000),
	)
}

func (this *BinanceClient) SubmitOrderUsdtMarketBuy(
	currency meta.Currency, amount decimal.Decimal, reference string,
) (*binance.CreateOrderResponse, error) {
	return this.SubmitOrderUsdtMarket(currency, binance.SideTypeBuy, amount, reference)
}

func (this *BinanceClient) SubmitOrderUsdtMarketSell(
	currency meta.Currency, amount decimal.Decimal, reference string,
) (*binance.CreateOrderResponse, error) {
	return this.SubmitOrderUsdtMarket(currency, binance.SideTypeSell, amount, reference)
}
