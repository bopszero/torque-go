package crontasks

import (
	"fmt"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type CurrencyPriceMap map[meta.Currency]decimal.Decimal

func CurrencyUpdatePrices() {
	defer apiutils.CronRunWithRecovery("CurrencyUpdatePrices", nil)

	apiKey := viper.GetString(config.KeyApiCryptoCompareKey)
	if apiKey == "" {
		panic(utils.IssueErrorf("missing config key `%s`", config.KeyApiCryptoCompareKey))
	}

	var currencyInfoList []models.CurrencyInfo
	comutils.PanicOnError(
		database.
			GetDbF(database.AliasWalletSlave).
			Find(&currencyInfoList).
			Error,
	)
	var (
		targetCurrencies    = make([]meta.Currency, 0, len(currencyInfoList))
		targetCurrencyCodes = make([]string, 0, len(currencyInfoList))
	)
	for _, currencyInfo := range currencyInfoList {
		if currencyInfo.Currency == constants.CurrencyTorque {
			continue
		}
		targetCurrencies = append(targetCurrencies, currencyInfo.Currency)
		targetCurrencyCodes = append(targetCurrencyCodes, currencyInfo.Currency.String())
	}

	client := getHttpClient(10 * time.Second)
	request := client.R().
		SetHeader(echo.HeaderAuthorization, fmt.Sprintf("Apikey %s", apiKey)).
		SetQueryParams(map[string]string{
			"fsyms": fmt.Sprintf("%v,%v", constants.CurrencyUSD, constants.CurrencyTetherUSD),
			"tsyms": strings.Join(targetCurrencyCodes, ","),
		})
	response, err := request.Get("https://min-api.cryptocompare.com/data/pricemulti")
	comutils.PanicOnError(err)
	if !response.IsSuccess() {
		panic(fmt.Errorf(
			"cryptocompare unexpected failed response | status=%v,body=%v",
			response.StatusCode(), response.String(),
		))
	}

	var priceTableMap map[meta.Currency]CurrencyPriceMap
	if err := comutils.JsonDecode(response.String(), &priceTableMap); err != nil {
		panic(err)
	}

	var (
		toPrice = func(usdRate decimal.Decimal) decimal.Decimal {
			return comutils.DecimalDivide(comutils.DecimalOne, usdRate)
		}
		usdPriceMap  = priceTableMap[constants.CurrencyUSD]
		usdtPriceMap = priceTableMap[constants.CurrencyTetherUSD]

		now            = time.Now()
		logger         = comlogging.GetLogger()
		walletMasterDB = database.GetDbF(database.AliasWalletMaster)
	)
	for _, currency := range targetCurrencies {
		rateUSD, hasUSD := usdPriceMap[currency]
		rateUSDT, hasUSDT := usdtPriceMap[currency]
		if !hasUSD || !hasUSDT {
			logger.
				WithField("currency", currency).
				Warnf("currency `%v` is missing from cryptocompare response", currency)
			continue
		}

		err := walletMasterDB.
			Model(&models.CurrencyInfo{}).
			Where(&models.CurrencyInfo{Currency: currency}).
			Updates(&models.CurrencyInfo{
				PriceUSD:   currencymod.NormalizeAmount(constants.CurrencyUSD, toPrice(rateUSD)),
				PriceUSDT:  currencymod.NormalizeAmount(constants.CurrencyTetherUSD, toPrice(rateUSDT)),
				UpdateTime: now.Unix(),
			}).
			Error
		if err != nil {
			panic(fmt.Errorf("update price for currency `%v` failed | err=%v", currency, err.Error()))
		}
	}

	usdtRateUSD := comutils.DecimalDivide(
		usdPriceMap[constants.CurrencyBitcoin],
		usdtPriceMap[constants.CurrencyBitcoin])
	err = walletMasterDB.
		Model(&models.CurrencyInfo{}).
		Where(&models.CurrencyInfo{Currency: constants.CurrencyTorque}).
		Updates(&models.CurrencyInfo{
			PriceUSD: currencymod.NormalizeAmount(
				constants.CurrencyUSD,
				toPrice(usdtRateUSD).Mul(constants.CurrencyTorquePriceUSDT),
			),
			UpdateTime: now.Unix(),
		}).
		Error
	if err != nil {
		panic(fmt.Errorf("update price for currency `%v` failed | err=%v", constants.CurrencyTorque, err.Error()))
	}
}
