package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/api"
	walletv1 "gitlab.com/snap-clickstaff/torque-go/api/services/wallet/v1"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/settingmod"
)

type InitOrderApiResponse struct {
	Data walletv1.PaymentInitOrderResponse `json:"data"`
}

func TestPaymentProfitReinvest(t *testing.T) {
	e := api.CreateEchoObject()
	rec := httptest.NewRecorder()

	amount, _ := decimal.NewFromString("20")
	checkoutReq := walletv1.PaymentInitOrderRequest{
		Currency: constants.CurrencyTorque,

		SrcChannelType:   constants.ChannelTypeSrcBalance,
		SrcChannelAmount: amount,

		DstChannelType:   constants.ChannelTypeDstProfitReinvest,
		DstChannelAmount: amount,
		DstChannelContext: meta.O{
			"user_identity": "test@gmail.com",
			"currency":      constants.CurrencyTetherUSD,
			"exchange_rate": decimal.NewFromFloat(0.05),
		},

		AmountSubTotal: amount,
		AmountTotal:    amount,

		Note: "Test reinvest USDT",
	}

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(comutils.JsonEncodeF(checkoutReq)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := patchContext(e.NewContext(req, rec))

	if !assert.NoError(t, walletv1.PaymentInitOrder(c)) {
		return
	}
	assertSuccessResponse(t, rec)

	var checkoutResp InitOrderApiResponse
	comutils.JsonDecodeF(readResponseBodyString(rec), &checkoutResp)
	executeReq := walletv1.PaymentExecuteOrderRequest{
		AuthCode: constants.TestTwoFaCode,
		OrderID:  checkoutResp.Data.OrderID,
	}

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(comutils.JsonEncodeF(executeReq)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = patchContext(e.NewContext(req, rec))

	if !assert.NoError(t, walletv1.PaymentExecuteOrder(c)) {
		return
	}

	assertSuccessResponse(t, rec)
}

func TestPaymentProfitWithdraw(t *testing.T) {
	e := api.CreateEchoObject()
	rec := httptest.NewRecorder()

	amount, _ := decimal.NewFromString("25")
	checkoutReq := walletv1.PaymentInitOrderRequest{
		Currency: constants.CurrencyTorque,

		SrcChannelType:   constants.ChannelTypeSrcBalance,
		SrcChannelAmount: amount,

		DstChannelType:   constants.ChannelTypeDstProfitWithdraw,
		DstChannelAmount: amount,
		DstChannelContext: meta.O{
			"address":       "0x83a7663b2b9d6d3f377a41d03b03ba0021e2f831",
			"currency":      constants.CurrencyTetherUSD,
			"exchange_rate": decimal.NewFromFloat(0.05),
		},

		AmountSubTotal: amount,
		AmountTotal:    amount,

		Note: "Test withdraw profit to USDT",
	}

	t.Log(comutils.JsonEncodeF(checkoutReq))
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(comutils.JsonEncodeF(checkoutReq)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := patchContext(e.NewContext(req, rec))

	if !assert.NoError(t, walletv1.PaymentInitOrder(c)) {
		return
	}
	assertSuccessResponse(t, rec)

	var checkoutResp InitOrderApiResponse
	comutils.JsonDecodeF(readResponseBodyString(rec), &checkoutResp)
	executeReq := walletv1.PaymentExecuteOrderRequest{
		AuthCode: constants.TestTwoFaCode,
		OrderID:  checkoutResp.Data.OrderID,
	}

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(comutils.JsonEncodeF(executeReq)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = patchContext(e.NewContext(req, rec))

	if !assert.NoError(t, walletv1.PaymentExecuteOrder(c)) {
		return
	}

	assertSuccessResponse(t, rec)
}

func TestPaymentP2P(t *testing.T) {
	e := api.CreateEchoObject()
	rec := httptest.NewRecorder()

	trasferFeeStr, err := settingmod.GetSettingValueFast(constants.SettingKeyTorqueTransferFee)
	assert.NoError(t, err)
	fee, err := decimal.NewFromString(trasferFeeStr)
	assert.NoError(t, err)

	amount, _ := decimal.NewFromString("12.34567898")
	checkoutReq := walletv1.PaymentInitOrderRequest{
		Currency: constants.CurrencyTorque,

		SrcChannelType:   constants.ChannelTypeSrcBalance,
		SrcChannelAmount: amount,

		DstChannelType:   constants.ChannelTypeDstTransfer,
		DstChannelAmount: amount,
		DstChannelContext: meta.O{
			"user_identity": "torque1@gmail.com",
			"note":          "Hello World!",
		},

		AmountSubTotal: amount,
		AmountTotal:    amount.Add(fee),
		Note:           "Test transfer",
	}

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(comutils.JsonEncodeF(checkoutReq)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := patchContext(e.NewContext(req, rec))

	if !assert.NoError(t, walletv1.PaymentInitOrder(c)) {
		return
	}
	assertSuccessResponse(t, rec)

	var checkoutResp InitOrderApiResponse
	comutils.JsonDecodeF(readResponseBodyString(rec), &checkoutResp)
	executeReq := walletv1.PaymentExecuteOrderRequest{
		AuthCode: constants.TestTwoFaCode,
		OrderID:  checkoutResp.Data.OrderID,
	}

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(comutils.JsonEncodeF(executeReq)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = patchContext(e.NewContext(req, rec))

	if !assert.NoError(t, walletv1.PaymentExecuteOrder(c)) {
		return
	}

	assertSuccessResponse(t, rec)
}

func TestPaymentPushEthTxn(t *testing.T) {
	e := api.CreateEchoObject()
	rec := httptest.NewRecorder()

	amount, _ := decimal.NewFromString("0.001")
	checkoutReq := walletv1.PaymentInitOrderRequest{
		Currency: constants.CurrencyEthereum,

		SrcChannelType:   constants.ChannelTypeSrcBlockchainNetwork,
		SrcChannelAmount: amount,

		DstChannelType:   constants.ChannelTypeDstBlockchainNetwork,
		DstChannelAmount: amount,
		DstChannelContext: meta.O{
			"to_address": "0xf12Db5BeC81DD6f41366EbcC599508212b751479", // PK: 0x0d7c7589fc58ecb2ccc745d68ced03657f9d2bee2bb69d6a14b65798a10d9e66
		},

		AmountSubTotal: amount,
		AmountTotal:    amount,
		Note:           "Test push ETH Txn",
	}

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(comutils.JsonEncodeF(checkoutReq)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := patchContext(e.NewContext(req, rec))

	if !assert.NoError(t, walletv1.PaymentInitOrder(c)) {
		return
	}
	assertSuccessResponse(t, rec)

	var checkoutResp InitOrderApiResponse
	comutils.JsonDecodeF(readResponseBodyString(rec), &checkoutResp)
	executeReq := walletv1.PaymentExecuteOrderRequest{
		AuthCode: constants.TestTwoFaCode,
		OrderID:  checkoutResp.Data.OrderID,
	}

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(comutils.JsonEncodeF(executeReq)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = patchContext(e.NewContext(req, rec))

	if !assert.NoError(t, walletv1.PaymentExecuteOrder(c)) {
		return
	}

	assertSuccessResponse(t, rec)
}

func TestPaymentPushEthTokenTxn(t *testing.T) {
	e := api.CreateEchoObject()
	rec := httptest.NewRecorder()

	amount, _ := decimal.NewFromString("10")
	checkoutReq := walletv1.PaymentInitOrderRequest{
		Currency: constants.CurrencyTetherUSD,

		SrcChannelType:   constants.ChannelTypeSrcBlockchainNetwork,
		SrcChannelAmount: amount,

		DstChannelType:   constants.ChannelTypeDstBlockchainNetwork,
		DstChannelAmount: amount,
		DstChannelContext: meta.O{
			"to_address": "0x946416a1F47c76656cb7517700615501F9940734",
		},

		AmountSubTotal: amount,
		AmountTotal:    amount,
		Note:           "Test push ETH Token Txn",
	}

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(comutils.JsonEncodeF(checkoutReq)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := patchContext(e.NewContext(req, rec))

	if !assert.NoError(t, walletv1.PaymentInitOrder(c)) {
		return
	}
	assertSuccessResponse(t, rec)

	var checkoutResp InitOrderApiResponse
	comutils.JsonDecodeF(readResponseBodyString(rec), &checkoutResp)
	executeReq := walletv1.PaymentExecuteOrderRequest{
		AuthCode: constants.TestTwoFaCode,
		OrderID:  checkoutResp.Data.OrderID,
	}

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(comutils.JsonEncodeF(executeReq)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = patchContext(e.NewContext(req, rec))

	if !assert.NoError(t, walletv1.PaymentExecuteOrder(c)) {
		return
	}

	assertSuccessResponse(t, rec)
}

func TestPaymentPushBtcTxn(t *testing.T) {
	e := api.CreateEchoObject()
	rec := httptest.NewRecorder()

	amount, _ := decimal.NewFromString("0.0001")
	checkoutReq := walletv1.PaymentInitOrderRequest{
		Currency: constants.CurrencyBitcoin,

		SrcChannelType:   constants.ChannelTypeSrcBlockchainNetwork,
		SrcChannelAmount: amount,

		DstChannelType:   constants.ChannelTypeDstBlockchainNetwork,
		DstChannelAmount: amount,
		DstChannelContext: meta.O{
			"to_address": "mvVBfUSjRdhFmzapBDbirEBHtH3AprMPfb",
		},

		AmountSubTotal: amount,
		AmountTotal:    amount,
		Note:           "Test push BTC Txn",
	}

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(comutils.JsonEncodeF(checkoutReq)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := patchContext(e.NewContext(req, rec))

	if !assert.NoError(t, walletv1.PaymentInitOrder(c)) {
		return
	}
	assertSuccessResponse(t, rec)

	var checkoutResp InitOrderApiResponse
	comutils.JsonDecodeF(readResponseBodyString(rec), &checkoutResp)
	executeReq := walletv1.PaymentExecuteOrderRequest{
		AuthCode: constants.TestTwoFaCode,
		OrderID:  checkoutResp.Data.OrderID,
	}

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(comutils.JsonEncodeF(executeReq)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = patchContext(e.NewContext(req, rec))

	if !assert.NoError(t, walletv1.PaymentExecuteOrder(c)) {
		return
	}

	assertSuccessResponse(t, rec)
}

func TestPaymentPushLtcTxn(t *testing.T) {
	e := api.CreateEchoObject()
	rec := httptest.NewRecorder()

	amount, _ := decimal.NewFromString("0.005")
	checkoutReq := walletv1.PaymentInitOrderRequest{
		Currency: constants.CurrencyLitecoin,

		SrcChannelType:   constants.ChannelTypeSrcBlockchainNetwork,
		SrcChannelAmount: amount,

		DstChannelType:   constants.ChannelTypeDstBlockchainNetwork,
		DstChannelAmount: amount,
		DstChannelContext: meta.O{
			"to_address": "mffGhrRYENi7cHzxQbmc33wxqwvMV9rKpF",
		},

		AmountSubTotal: amount,
		AmountTotal:    amount,
		Note:           "Test push LTC Txn",
	}

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(comutils.JsonEncodeF(checkoutReq)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := patchContext(e.NewContext(req, rec))

	if !assert.NoError(t, walletv1.PaymentInitOrder(c)) {
		return
	}
	assertSuccessResponse(t, rec)

	var checkoutResp InitOrderApiResponse
	comutils.JsonDecodeF(readResponseBodyString(rec), &checkoutResp)
	executeReq := walletv1.PaymentExecuteOrderRequest{
		AuthCode: constants.TestTwoFaCode,
		OrderID:  checkoutResp.Data.OrderID,
	}

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(comutils.JsonEncodeF(executeReq)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = patchContext(e.NewContext(req, rec))

	if !assert.NoError(t, walletv1.PaymentExecuteOrder(c)) {
		return
	}

	assertSuccessResponse(t, rec)
}

func TestPaymentPushBchTxn(t *testing.T) {
	e := api.CreateEchoObject()
	rec := httptest.NewRecorder()

	amount, _ := decimal.NewFromString("0.001")
	checkoutReq := walletv1.PaymentInitOrderRequest{
		Currency: constants.CurrencyBitcoinCash,

		SrcChannelType:   constants.ChannelTypeSrcBlockchainNetwork,
		SrcChannelAmount: amount,

		DstChannelType:   constants.ChannelTypeDstBlockchainNetwork,
		DstChannelAmount: amount,
		DstChannelContext: meta.O{
			"to_address": "qztvxcen0vfy6d39mverhp6p03szvvlt2stpgmy65c",
		},

		AmountSubTotal: amount,
		AmountTotal:    amount,
		Note:           "Test push BCH Txn",
	}

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(comutils.JsonEncodeF(checkoutReq)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := patchContext(e.NewContext(req, rec))

	if !assert.NoError(t, walletv1.PaymentInitOrder(c)) {
		return
	}
	assertSuccessResponse(t, rec)

	var checkoutResp InitOrderApiResponse
	comutils.JsonDecodeF(readResponseBodyString(rec), &checkoutResp)
	executeReq := walletv1.PaymentExecuteOrderRequest{
		AuthCode: constants.TestTwoFaCode,
		OrderID:  checkoutResp.Data.OrderID,
	}

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(comutils.JsonEncodeF(executeReq)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = patchContext(e.NewContext(req, rec))

	if !assert.NoError(t, walletv1.PaymentExecuteOrder(c)) {
		return
	}

	assertSuccessResponse(t, rec)
}

func TestPaymentTradingReward(t *testing.T) {
	e := api.CreateEchoObject()
	rec := httptest.NewRecorder()

	amount := decimal.NewFromInt(30)
	checkoutReq := walletv1.PaymentInitOrderRequest{
		Currency: constants.CurrencyTorque,

		SrcChannelType:   constants.ChannelTypeSrcTradingReward,
		SrcChannelRef:    "2020-05-30",
		SrcChannelAmount: amount,
		SrcChannelContext: meta.O{
			"amount_daily_profit":         5,
			"amount_affiliate_commission": 10,
			"amount_leader_commission":    15,
		},

		DstChannelType:   constants.ChannelTypeDstBalance,
		DstChannelAmount: amount,
		DstChannelContext: meta.O{
			"inner_txns": []meta.O{
				{
					"txn_type": constants.WalletBalanceTypeMetaCrDailyProfit.ID,
					"value":    5,
				},
				{
					"txn_type": constants.WalletBalanceTypeMetaCrAffiliateCommission.ID,
					"value":    10,
				},
				{
					"txn_type": constants.WalletBalanceTypeMetaCrLeaderCommission.ID,
					"value":    15,
				},
			},
		},

		AmountSubTotal: amount,
		AmountTotal:    amount,
		Note:           "Test deliver trading reward",
	}

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(comutils.JsonEncodeF(checkoutReq)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := patchContext(e.NewContext(req, rec))

	if !assert.NoError(t, walletv1.PaymentInitOrder(c)) {
		return
	}
	assertSuccessResponse(t, rec)

	var checkoutResp InitOrderApiResponse
	comutils.JsonDecodeF(readResponseBodyString(rec), &checkoutResp)
	executeReq := walletv1.PaymentExecuteOrderRequest{
		AuthCode: constants.TestTwoFaCode,
		OrderID:  checkoutResp.Data.OrderID,
	}

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(comutils.JsonEncodeF(executeReq)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = patchContext(e.NewContext(req, rec))

	if !assert.NoError(t, walletv1.PaymentExecuteOrder(c)) {
		return
	}

	assertSuccessResponse(t, rec)
}

func TestPaymentTorquePurchase(t *testing.T) {
	e := api.CreateEchoObject()
	rec := httptest.NewRecorder()

	amount := decimal.NewFromFloat(1.0 + float64(time.Now().Unix()%100000)/100000)
	checkoutReq := walletv1.PaymentInitOrderRequest{
		Currency: constants.CurrencyTetherUSD,

		SrcChannelType:   constants.ChannelTypeSrcBlockchainNetwork,
		SrcChannelAmount: amount,

		DstChannelType:   constants.ChannelTypeDstTorqueConvert,
		DstChannelAmount: amount,
		DstChannelContext: meta.O{
			"exchange_rate": "20",
		},

		AmountSubTotal: amount,
		AmountTotal:    amount,
		Note:           "Test puchase Torque by USDT",
	}

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(comutils.JsonEncodeF(checkoutReq)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := patchContext(e.NewContext(req, rec))

	if !assert.NoError(t, walletv1.PaymentInitOrder(c)) {
		return
	}
	assertSuccessResponse(t, rec)

	var checkoutResp InitOrderApiResponse
	comutils.JsonDecodeF(readResponseBodyString(rec), &checkoutResp)
	executeReq := walletv1.PaymentExecuteOrderRequest{
		AuthCode: constants.TestTwoFaCode,
		OrderID:  checkoutResp.Data.OrderID,
	}

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(comutils.JsonEncodeF(executeReq)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = patchContext(e.NewContext(req, rec))

	if !assert.NoError(t, walletv1.PaymentExecuteOrder(c)) {
		return
	}

	assertSuccessResponse(t, rec)
}

func TestRedeemPromoCode(t *testing.T) {
	e := api.CreateEchoObject()
	rec := httptest.NewRecorder()

	amount := decimal.NewFromInt(10000)
	checkoutReq := walletv1.PaymentInitOrderRequest{
		Currency: constants.CurrencyTorque,

		SrcChannelType:   constants.ChannelTypeSrcPromoCode,
		SrcChannelAmount: amount,
		SrcChannelContext: meta.O{
			"code": "aaaaaaaaaaa",
		},

		DstChannelType:   constants.ChannelTypeDstBalance,
		DstChannelAmount: amount,

		AmountSubTotal: amount,
		AmountTotal:    amount,
		Note:           "Test redeem Promo code",
	}

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(comutils.JsonEncodeF(checkoutReq)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := patchContext(e.NewContext(req, rec))

	if !assert.NoError(t, walletv1.PaymentInitOrder(c)) {
		return
	}
	assertSuccessResponse(t, rec)

	var checkoutResp InitOrderApiResponse
	comutils.JsonDecodeF(readResponseBodyString(rec), &checkoutResp)
	executeReq := walletv1.PaymentExecuteOrderRequest{
		AuthCode: constants.TestTwoFaCode,
		OrderID:  checkoutResp.Data.OrderID,
	}

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(comutils.JsonEncodeF(executeReq)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = patchContext(e.NewContext(req, rec))

	if !assert.NoError(t, walletv1.PaymentExecuteOrder(c)) {
		return
	}

	assertSuccessResponse(t, rec)
}
