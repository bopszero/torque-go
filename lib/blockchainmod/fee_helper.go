package blockchainmod

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

var feeClient = comtypes.NewSingleton(func() interface{} {
	return utils.NewRestyClient(15 * time.Second)
})

func getFeeClient() *resty.Client {
	return feeClient.Get().(*resty.Client)
}

type bitcoinfeesEarnFeeResponse struct {
	FastestFee  int64 `json:"fastestFee"`
	HourFee     int64 `json:"hourFee"`
	HalfHourFee int64 `json:"halfHourFee"`
}

func GetBitcoinFeeInfoViaBitcoinfeesEarn() (feeInfo FeeInfo, err error) {
	response, err := getFeeClient().R().
		Get("https://bitcoinfees.earn.com/api/v1/fees/recommended")
	if err != nil {
		err = utils.WrapError(err)
		return
	}
	if !response.IsSuccess() {
		err = utils.IssueErrorf(
			"get fee from bitcoinfees.earn.com failed | status_code=%v",
			response.StatusCode(),
		)
		return
	}

	var responseModel bitcoinfeesEarnFeeResponse
	if err = comutils.JsonDecode(response.String(), &responseModel); err != nil {
		err = utils.WrapError(err)
		return
	}

	feeInfo, err = NewFeeInfo(
		constants.CurrencyBitcoin, constants.CurrencySubBitcoinSatoshi,
		decimal.NewFromInt(responseModel.HalfHourFee),
	)
	if err != nil {
		return
	}

	feeInfo.PriceHigh = decimal.NewFromInt(responseModel.FastestFee)
	feeInfo.PriceLow = decimal.NewFromInt(responseModel.HourFee)

	return
}

type blockCypherFeeResponse struct {
	PriceKbHigh     decimal.Decimal `json:"high_fee_per_kb"`
	PriceKbStandard decimal.Decimal `json:"medium_fee_per_kb"`
	PriceKbLow      decimal.Decimal `json:"low_fee_per_kb"`
}

func getBitcoinLikeFeeInfoViaBlockCypher(currency meta.Currency) (feeInfo FeeInfo, err error) {
	apiKey := viper.GetString(config.KeyApiBlockCypherKey)
	if apiKey == "" {
		err = utils.IssueErrorf(
			"missing API key `%v` for https://api.blockcypher.com",
			config.KeyApiBlockCypherKey,
		)
		return
	}

	uri := fmt.Sprintf(
		"https://api.blockcypher.com/v1/%v/main",
		strings.ToLower(currency.String()),
	)
	response, err := getFeeClient().R().
		SetHeader("token", apiKey).
		Get(uri)
	if err != nil {
		err = utils.WrapError(err)
		return
	}
	if !response.IsSuccess() {
		err = utils.IssueErrorf(
			"get fee from api.blockcypher.com failed | status_code=%v",
			response.StatusCode())
		return
	}

	var feeModel blockCypherFeeResponse
	if err = comutils.JsonDecode(response.String(), &feeModel); err != nil {
		err = utils.WrapError(err)
		return
	}

	kbToByteMultipler := decimal.NewFromInt(1024)

	feeInfo, err = NewFeeInfo(
		currency,
		constants.CurrencySubBitcoinSatoshi,
		feeModel.PriceKbStandard.Div(kbToByteMultipler))
	if err != nil {
		return
	}

	feeInfo.PriceLow = feeModel.PriceKbLow.Div(kbToByteMultipler)
	feeInfo.PriceHigh = feeModel.PriceKbHigh.Div(kbToByteMultipler)

	return feeInfo, nil
}

func GetBitcoinFeeInfoViaBlockCypher() (FeeInfo, error) {
	return getBitcoinLikeFeeInfoViaBlockCypher(constants.CurrencyBitcoin)
}

func GetLitecoinFeeInfoViaBlockCypher() (FeeInfo, error) {
	return getBitcoinLikeFeeInfoViaBlockCypher(constants.CurrencyLitecoin)
}

type ethGasStationFeeResponse struct {
	Fast    int64 `json:"fast"`
	Fastest int64 `json:"fastest"`
	Average int64 `json:"average"`
	SafeLow int64 `json:"safeLow"`
}

func GetEthereumFeeInfoViaEthGasStation() (feeInfo FeeInfo, err error) {
	response, err := getFeeClient().R().Get("https://ethgasstation.info/api/ethgasAPI.json")
	if err != nil {
		err = utils.WrapError(err)
		return
	}
	if !response.IsSuccess() {
		err = utils.IssueErrorf(
			"get fee from ethgasstation.info failed | status_code=%v",
			response.StatusCode(),
		)
		return
	}

	var responseModel ethGasStationFeeResponse
	if err = comutils.JsonDecode(response.String(), &responseModel); err != nil {
		err = utils.WrapError(err)
		return
	}

	return NewFeeInfo(
		constants.CurrencyEthereum, constants.CurrencySubEthereumGwei,
		decimal.NewFromInt(responseModel.Average).
			Div(decimal.NewFromInt(10)),
	)
}

type etherscanGasOracleResponse struct {
	Status  int                       `json:"status,string"`
	Message string                    `json:"message"`
	FeeData etherscanGasOracleFeeData `json:"result"`
}

type etherscanGasOracleFeeData struct {
	LastBlock       int64           `json:"LastBlock,string"`
	SafeGasPrice    decimal.Decimal `json:"SafeGasPrice"`
	ProposeGasPrice decimal.Decimal `json:"ProposeGasPrice"`
	FastGasPrice    decimal.Decimal `json:"FastGasPrice"`
}

func GetEthereumFeeInfoViaEtherscan() (feeInfo FeeInfo, err error) {
	response, err := getFeeClient().R().
		SetQueryParams(map[string]string{
			"module": "gastracker",
			"action": "gasoracle",
			"apikey": viper.GetString(config.KeyApiEtherscanKey),
		}).
		Get("https://api.etherscan.io/api")
	if err != nil {
		err = utils.WrapError(err)
		return
	}
	if !response.IsSuccess() {
		err = utils.IssueErrorf(
			"get fee from api.etherscan.io failed | status_code=%v,response=%s",
			response.StatusCode(), response.String(),
		)
		return
	}

	var responseModel etherscanGasOracleResponse
	if err = comutils.JsonDecode(response.String(), &responseModel); err != nil {
		err = utils.WrapError(err)
		return
	}

	feeInfo, err = NewFeeInfo(
		constants.CurrencyEthereum,
		constants.CurrencySubEthereumGwei,
		responseModel.FeeData.ProposeGasPrice,
	)
	if err != nil {
		return
	}

	feeInfo.PriceLow = responseModel.FeeData.SafeGasPrice
	feeInfo.PriceHigh = responseModel.FeeData.FastGasPrice

	return
}

type myCoinGateBitCoreFeeResponse struct {
	Chain        string          `json:"chain"`
	Network      string          `json:"network"`
	Price        decimal.Decimal `json:"price"`
	Currency     meta.Currency   `json:"currency"`
	InBlockCount uint64          `json:"blocks"`
}

func GetFeeResponseFromMyCoinGateBitcore(currency meta.Currency, backBlockCount uint8) (
	feeInfo FeeInfo, err error,
) {
	var (
		client = getBitcoreClient()
		uri    = fmt.Sprintf("/api/%v/mainnet/fee/%v/", currency.StringU(), backBlockCount)
	)
	response, err := client.R().Get(uri)
	if err != nil {
		err = utils.WrapError(err)
		return
	}
	if !response.IsSuccess() {
		err = utils.IssueErrorf(
			"get fee from %v failed | status_code=%v,currency=%v",
			client.HostURL, response.StatusCode(), currency)
		return
	}

	var feeResponse myCoinGateBitCoreFeeResponse
	if err = comutils.JsonDecode(response.String(), &feeResponse); err != nil {
		err = utils.WrapError(err)
		return
	}

	feeInfo, err = NewFeeInfo(
		currency,
		feeResponse.Currency,
		decimal.Max(feeResponse.Price, comutils.DecimalOne).Truncate(0),
	)
	return
}

func GetBitcoinFeeInfoViaMyCoinGateBitcore() (feeInfo FeeInfo, err error) {
	return GetFeeResponseFromMyCoinGateBitcore(constants.CurrencyBitcoin, 3)
}

func GetLitecoinFeeInfoViaMyCoinGateBitcore() (feeInfo FeeInfo, err error) {
	return GetFeeResponseFromMyCoinGateBitcore(constants.CurrencyLitecoin, 3)
}

func GetTronFeeInfoViaMyCoinGateBitcore() (feeInfo FeeInfo, err error) {
	return NewFeeInfo(
		constants.CurrencyTron,
		constants.CurrencySubTronSun,
		TronFeeNormalTxnPrice,
	)
}

func GetRippleFeeInfoViaMyCoinGateBitcore() (feeInfo FeeInfo, err error) {
	return NewFeeInfo(
		constants.CurrencyRipple,
		constants.CurrencySubRippleDrop,
		decimal.NewFromInt(RippleNormalTxnFee),
	)
}
