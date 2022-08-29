package orderchannel

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcache"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

const (
	DstTorquePurchaseBinanceSymbolStatusTrading = "TRADING"
	DstTorquePurchaseBinanceFilterLotSize       = "LOT_SIZE"
	DstTorquePurchaseBinanceFilterMinNotional   = "MIN_NOTIONAL"
)

type (
	DstTorquePurchaseBinancePairFilter struct {
		FilterType string `json:"filterType"`

		// -- LOT_SIZE --
		MinQuantity decimal.Decimal `json:"minQty"`
		MaxQuantity decimal.Decimal `json:"maxQty"`
		StepSize    decimal.Decimal `json:"stepSize"`

		// -- MIN_NOTIONAL --
		MinNotional decimal.Decimal `json:"minNotional"`
	}
	DstTorquePurchaseBinancePairInfo struct {
		Symbol        string                               `json:"symbol"`
		Status        string                               `json:"status"`
		CaseCurrency  meta.Currency                        `json:"baseAsset"`
		QuoteCurrency meta.Currency                        `json:"quoteAsset"`
		Filters       []DstTorquePurchaseBinancePairFilter `json:"filters"`
		FilterMap     map[string]DstTorquePurchaseBinancePairFilter
	}
	DstTorquePurchaseBinanceExchangeInfo struct {
		Symbols []DstTorquePurchaseBinancePairInfo `json:"symbols"`
	}
)

var (
	DstTorquePurchaseBinancePairInfoMapProxy = comcache.NewCacheObject(
		15*time.Minute,
		func() (_ interface{}, err error) {
			client := utils.NewRestyClient(5 * time.Second)
			response, err := client.R().Get("https://api.binance.com/api/v3/exchangeInfo")
			if err != nil {
				return
			}
			var exchangeInfo DstTorquePurchaseBinanceExchangeInfo
			if err = comutils.JsonDecode(response.String(), &exchangeInfo); err != nil {
				return
			}
			pairInfoMap := make(map[string]DstTorquePurchaseBinancePairInfo)
			for _, sym := range exchangeInfo.Symbols {
				if sym.Status != DstTorquePurchaseBinanceSymbolStatusTrading {
					continue
				}
				if sym.QuoteCurrency != constants.CurrencyTetherUSD {
					continue
				}

				filterMap := make(map[string]DstTorquePurchaseBinancePairFilter, len(sym.Filters))
				for _, filter := range sym.Filters {
					filterMap[filter.FilterType] = filter
				}
				sym.FilterMap = filterMap

				pairInfoMap[sym.Symbol] = sym
			}
			return pairInfoMap, nil
		},
	)
)

func DstTorquePurchaseGetBinancePairInfoMap() (map[string]DstTorquePurchaseBinancePairInfo, error) {
	mapObj, err := DstTorquePurchaseBinancePairInfoMapProxy.Get()
	if err != nil {
		return nil, err
	}
	return mapObj.(map[string]DstTorquePurchaseBinancePairInfo), nil
}

func DstTorquePurchaseGetBinancePairInfo(currency meta.Currency) (info DstTorquePurchaseBinancePairInfo, err error) {
	infoMap, err := DstTorquePurchaseGetBinancePairInfoMap()
	if err != nil {
		return
	}
	pairCode := fmt.Sprintf("%v%v", currency, constants.CurrencyTetherUSD)
	info, ok := infoMap[pairCode]
	if !ok {
		err = utils.WrapError(constants.ErrorCurrency)
		return
	}
	return
}

func (this *DstTorquePurchaseChannel) validateExchangeValues(order models.Order) (err error) {
	if order.Currency == constants.CurrencyTetherUSD {
		return
	}
	pairInfo, err := DstTorquePurchaseGetBinancePairInfo(order.Currency)
	if err != nil {
		return
	}
	amount := order.AmountSubTotal
	if filter, ok := pairInfo.FilterMap[DstTorquePurchaseBinanceFilterLotSize]; ok {
		if amount.LessThan(filter.MinQuantity) {
			return utils.WrapError(constants.ErrorAmountTooLow)
		}
		if amount.GreaterThan(filter.MaxQuantity) {
			return utils.WrapError(constants.ErrorAmountTooHigh)
		}
	}
	if filter, ok := pairInfo.FilterMap[DstTorquePurchaseBinanceFilterMinNotional]; ok {
		var (
			currencyInfo = currencymod.GetCurrencyInfoFastF(order.Currency)
			usdtPrice    = amount.Mul(currencyInfo.PriceUSDT)
		)
		if usdtPrice.LessThan(filter.MinNotional) {
			return utils.WrapError(constants.ErrorAmountTooLow)
		}
	}
	return nil
}

func (this *DstTorquePurchaseChannel) fetchExchangeAmount(order models.Order) (amount decimal.Decimal, err error) {
	if this.isDestinationCurrency(order) {
		return decimal.Zero, nil
	}

	pairInfo, err := DstTorquePurchaseGetBinancePairInfo(order.Currency)
	if err != nil {
		return
	}
	amount = order.AmountSubTotal
	if filter, ok := pairInfo.FilterMap[DstTorquePurchaseBinanceFilterLotSize]; ok && !filter.StepSize.IsZero() {
		amount = comutils.DecimalDivide(amount, filter.StepSize).Truncate(0).Mul(filter.StepSize)
	}
	return
}
