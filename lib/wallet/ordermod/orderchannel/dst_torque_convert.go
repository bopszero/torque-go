package orderchannel

import (
	"fmt"
	"reflect"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbfields"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/lockmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/thirdpartymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/usermod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/binanceconv"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/ordermod"
	"gorm.io/gorm"
)

type (
	DstTorquePurchaseMeta struct {
		ExchangeRate decimal.Decimal `json:"exchange_rate"`
		TorqueAmount decimal.Decimal `json:"torque_amount"`

		BlockchainFeeInfo blockchainmod.FeeInfo `json:"blockchain_fee_info"`
	}

	DstTorquePurchaseCheckoutInfo struct {
		DstBlockchainNetworkCheckoutInfo
		ConversionRate meta.CurrencyConversionRate `json:"conversion_rate"`
	}
)

func init() {
	comutils.PanicOnError(
		ordermod.ChannelRegister(&DstTorquePurchaseChannel{}),
	)
}

type DstTorquePurchaseChannel struct {
	baseChannel
}

func (this *DstTorquePurchaseChannel) GetType() meta.ChannelType {
	return constants.ChannelTypeDstTorqueConvert
}

func (this *DstTorquePurchaseChannel) GetMetaType() reflect.Type {
	return reflect.TypeOf(DstTorquePurchaseMeta{})
}

func (this *DstTorquePurchaseChannel) getMeta(order *models.Order) (*DstTorquePurchaseMeta, error) {
	var metaModel DstTorquePurchaseMeta
	if err := ordermod.GetOrderChannelMetaData(order, this.GetType(), &metaModel); err != nil {
		return nil, err
	}

	return &metaModel, nil
}

func (this *DstTorquePurchaseChannel) GetCheckoutInfo(ctx comcontext.Context, order *models.Order) (
	interface{}, error,
) {
	coin, err := blockchainmod.GetCoinNative(order.Currency)
	if err != nil {
		return nil, err
	}
	account, err := coin.NewAccountSystem(ctx, order.UID)
	if err != nil {
		return nil, err
	}
	feeInfo, err := account.GetFeeInfoToAddress("")
	if err != nil {
		return nil, err
	}

	priceUsdt := currencymod.GetCurrencyPriceUsdtFastF(order.Currency)
	priceMarkup := this.getCurrencyPriceMarkup(order.Currency)
	if priceMarkup != nil {
		priceUsdt = priceMarkup.For(priceUsdt)
	}

	checkoutInfo := DstTorquePurchaseCheckoutInfo{
		DstBlockchainNetworkCheckoutInfo: DstBlockchainNetworkCheckoutInfo{
			FeeInfo:      feeInfo,
			MinTxnAmount: coin.GetMinTxnAmount(),
		},
		ConversionRate: meta.CurrencyConversionRate{
			FromCurrency: order.Currency,
			ToCurrency:   constants.CurrencyTorque,
			Value:        currencymod.ConvertUsdtToTorque(priceUsdt),
		},
	}
	return &checkoutInfo, nil
}

func (this *DstTorquePurchaseChannel) getCompanyAddress(order *models.Order) string {
	companyAddressMap := viper.GetStringMapString(config.KeyTorquePurchaseCompanyAddressMap)
	if companyAddressMap == nil {
		panic(utils.IssueErrorf("torque purchase address map is empty"))
	}

	companyAddress, ok := companyAddressMap[order.Currency.StringL()]
	if !ok {
		panic(utils.IssueErrorf("torque purchase currency address is missing | currency=%v", order.Currency))
	}

	return companyAddress
}

func (this *DstTorquePurchaseChannel) Init(ctx comcontext.Context, order *models.Order) error {
	metaModel, err := this.getMeta(order)
	if err != nil {
		return err
	}

	systemPriceUsdt, err := currencymod.GetCurrencyPriceUsdtFast(order.Currency)
	if err != nil {
		return err
	}
	priceMarkuper := this.getCurrencyPriceMarkup(order.Currency)
	if priceMarkuper != nil {
		systemPriceUsdt = priceMarkuper.For(systemPriceUsdt)
	}
	systemPriceTorque := currencymod.ConvertUsdtToTorque(systemPriceUsdt)

	systemExchangeRate := systemPriceTorque
	if systemExchangeRate.LessThan(metaModel.ExchangeRate) || metaModel.ExchangeRate.IsZero() {
		metaModel.ExchangeRate = systemExchangeRate
	}
	metaModel.TorqueAmount = currencymod.NormalizeAmount(
		constants.CurrencyTorque,
		order.AmountSubTotal.Mul(metaModel.ExchangeRate))
	if err := ordermod.SetOrderChannelMetaData(order, this.GetType(), metaModel); err != nil {
		return err
	}

	return nil
}

func (this *DstTorquePurchaseChannel) PreValidate(ctx comcontext.Context, order *models.Order) error {
	if ordermod.IsUserActionLocked(ctx, order.UID, constants.LockActionConvertCoinToTORQ) {
		return utils.WrapError(constants.ErrorUserActionLocked)
	}
	if err := validateBlockchainOrder(ctx, this, *order); err != nil {
		return err
	}
	if err := this.validateExchangeValues(*order); err != nil {
		return err
	}

	metaModel, err := this.getMeta(order)
	if err != nil {
		return err
	}
	if metaModel.TorqueAmount.IsZero() {
		return constants.ErrorAmountTooLow
	}

	pendingOrders, err := blockchainGetPendingOrders(ctx, order)
	if err != nil {
		return err
	}
	if len(pendingOrders) > 0 {
		return utils.WrapError(constants.ErrorOrderConcurrent)
	}

	return nil
}

func (this *DstTorquePurchaseChannel) getCurrencyPriceMarkup(currency meta.Currency) *meta.AmountMarkup {
	if this.infoModel.BlockchainNetworkConfig.CurrencyMarkupPriceMap == nil {
		return nil
	}
	priceMarkup, ok := this.infoModel.BlockchainNetworkConfig.CurrencyMarkupPriceMap[currency]
	if !ok {
		return nil
	}
	return &priceMarkup
}

func (this *DstTorquePurchaseChannel) Prepare(ctx comcontext.Context, order *models.Order) (
	meta.OrderStepResultCode, error,
) {
	metaModel, err := this.getMeta(order)
	if err != nil {
		return ordermod.OrderStepResultCodeFail, err
	}

	if !this.isDestinationCurrency(*order) {
		lock, err := lockmod.LockSimple("order:torque-convert:currency:%v", order.Currency)
		if err != nil {
			return ordermod.OrderStepResultCodeRetry, err
		}
		defer lock.Unlock()

		binanceBalance, err := binanceconv.GetBalance(order.Currency)
		if err != nil {
			return ordermod.OrderStepResultCodeRetry, err
		}
		if order.AmountSubTotal.GreaterThan(binanceBalance) {
			comlogging.GetLogger().
				WithContext(ctx).
				WithFields(logrus.Fields{
					"currency":       order.Currency,
					"balance":        binanceBalance,
					"request_amount": order.AmountSubTotal,
				}).
				Errorf("binance `%v` balance insufficient", order.Currency)
			return ordermod.OrderStepResultCodeFail, utils.WrapError(constants.ErrorChannelConversionOutOfCoin)
		}
	}

	companyCurrencyAddress := this.getCompanyAddress(order)
	cryptoTxn, resultCode, err := executeBlockchainTxn(
		ctx,
		order, companyCurrencyAddress, &metaModel.BlockchainFeeInfo)
	if err == nil {
		if err := ordermod.SetOrderChannelMetaData(order, this.GetType(), metaModel); err != nil {
			comlogging.GetLogger().
				WithContext(ctx).
				WithField("order_id", order.ID).
				Warn("cannot set meta info after submit blockchain txn")
		}
	}
	if cryptoTxn.ID > 0 {
		order.DstChannelRef = cryptoTxn.Hash
	}
	if err != nil || resultCode != ordermod.OrderStepResultCodeSuccess {
		return resultCode, err
	}

	conversionModel, err := this.recordPurchasing(ctx, order)
	if err != nil {
		return ordermod.OrderStepResultCodeRetry, err
	}
	order.DstChannelID = conversionModel.ID

	return ordermod.OrderStepResultCodeSuccess, nil
}

func (this *DstTorquePurchaseChannel) isDestinationCurrency(order models.Order) bool {
	return order.Currency == constants.CurrencyTetherUSD
}

func (this *DstTorquePurchaseChannel) recordPurchasing(ctx comcontext.Context, order *models.Order) (
	conversionModel models.TorqueCryptoConversion, err error,
) {
	metaModel, err := this.getMeta(order)
	if err != nil {
		return
	}

	now := time.Now()
	conversionModel = models.TorqueCryptoConversion{
		UID:      order.UID,
		Currency: order.Currency,
		TxnHash:  order.DstChannelRef,

		SendOrderID: order.ID,
	}
	err = database.Atomic(ctx, database.AliasWalletMaster, func(dbTxn *gorm.DB) (err error) {
		err = dbTxn.First(&conversionModel, &conversionModel).Error
		if database.IsDbError(err) {
			return utils.WrapError(err)
		}
		if conversionModel.ID > 0 {
			return nil
		}
		conversionModel.CurrencyAmount = order.AmountSubTotal
		conversionModel.ExchangeRate = metaModel.ExchangeRate
		conversionModel.TorqueAmount = metaModel.TorqueAmount
		conversionModel.ExtraData = make(dbfields.JsonField)
		conversionModel.CreateTime = now.Unix()
		conversionModel.UpdateTime = now.Unix()
		if err = dbTxn.Save(&conversionModel).Error; err != nil {
			return utils.WrapError(err)
		}
		return nil
	})
	if err != nil {
		return
	}
	err = database.Atomic(ctx, database.AliasWalletMaster, func(dbTxn *gorm.DB) (err error) {
		if err = this.requestExchangeSell(ctx, *order, &conversionModel); err != nil {
			return
		}
		if err = dbTxn.Save(&conversionModel).Error; err != nil {
			return utils.WrapError(err)
		}
		return nil
	})
	if err != nil {
		return
	}

	return
}

func (this *DstTorquePurchaseChannel) requestExchangeSell(
	ctx comcontext.Context,
	order models.Order, model *models.TorqueCryptoConversion,
) error {
	if this.isDestinationCurrency(order) {
		return nil
	}
	if model.ExchangeRef != "" {
		return nil
	}

	var (
		client    = thirdpartymod.GetBinanceSystemClient()
		reference = fmt.Sprintf("torq-%v", model.ID)
	)
	respOrder, err := client.GetOrderByRef(order.Currency, reference)
	if err != nil && !thirdpartymod.IsBinanceErrorCode(err, thirdpartymod.BinanceErrorCodeOrderFound) {
		return utils.WrapError(err)
	}
	if respOrder != nil {
		model.ExtraData["binance_response"] = respOrder
		return nil
	}

	convAmount, err := this.fetchExchangeAmount(order)
	if err != nil {
		return err
	}
	response, err := client.SubmitOrderUsdtMarketSell(order.Currency, convAmount, reference)
	logEntry := comlogging.GetLogger().
		WithContext(ctx).
		WithFields(logrus.Fields{
			"order_id": order.ID,
			"model_id": model.ID,
			"response": comutils.JsonEncodeF(response),
		})
	if err != nil {
		logEntry.
			WithError(err).
			Error("torque purchase to binance conversion failed")
		model.ExtraData["binance_response"] = meta.O{"error": err.Error()}
		return nil
	}

	commitAmountUSDT := currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: constants.CurrencyTorque,
			Value:    model.TorqueAmount,
		},
		constants.CurrencyTetherUSD,
	)
	binanceReceiveAmountUSDT := comutils.NewDecimalF(response.CummulativeQuoteQuantity)

	model.ProfitUSDT = binanceReceiveAmountUSDT.Sub(commitAmountUSDT.Value)
	model.ExchangeRef = comutils.Stringify(response.OrderID)
	model.ExtraData["binance_response"] = response

	logEntry.Info("torque purchase to binance conversion success")

	return nil
}

func (this *DstTorquePurchaseChannel) PrepareReverse(ctx comcontext.Context, order *models.Order) (
	meta.OrderStepResultCode, error,
) {
	if err := removeLocalBlockchainTxn(ctx, order); err != nil {
		return ordermod.OrderStepResultCodeRetry, err
	}

	return ordermod.OrderStepResultCodeSuccess, nil
}

func (this *DstTorquePurchaseChannel) Execute(ctx comcontext.Context, order *models.Order) (
	meta.OrderStepResultCode, error,
) {
	if !config.Debug && usermod.IsUserContext(ctx) {
		return ordermod.OrderStepResultCodeRetry, nil
	}

	if code, err := watchBlockchainTxnConfirmations(ctx, order); code != ordermod.OrderStepResultCodeSuccess {
		return code, err
	}

	err := database.Atomic(ctx, database.AliasWalletMaster, func(dbTxn *gorm.DB) (err error) {
		topupOrder, err := this.genTopupOrder(*order)
		if err != nil {
			return
		}
		if err = dbTxn.Save(&topupOrder).Error; err != nil {
			return
		}

		return dbTxn.
			Model(&models.TorqueCryptoConversion{
				ID: order.DstChannelID,
			}).
			Updates(&models.TorqueCryptoConversion{
				ReceiveOrderID: topupOrder.ID,
				UpdateTime:     time.Now().Unix(),
			}).
			Error
	})
	if err != nil {
		return ordermod.OrderStepResultCodeRetry, err
	}

	return ordermod.OrderStepResultCodeSuccess, nil
}

func (this *DstTorquePurchaseChannel) genTopupOrder(order models.Order) (
	topupOrder models.Order, err error,
) {
	metaModel, err := this.getMeta(&order)
	if err != nil {
		return
	}

	topupOrder = ordermod.NewUserOrder(
		order.UID, constants.CurrencyTorque,
		constants.ChannelTypeSrcTorqueConvert, constants.ChannelTypeDstBalance)
	topupOrder.SrcChannelID = order.ID
	topupOrder.SrcChannelAmount = metaModel.TorqueAmount
	topupOrder.DstChannelAmount = metaModel.TorqueAmount
	topupOrder.AmountSubTotal = metaModel.TorqueAmount
	topupOrder.AmountTotal = metaModel.TorqueAmount

	srcMeta := SrcTorquePurchaseMeta{
		ExchangeRate: metaModel.ExchangeRate,
		CurrencyAmount: meta.CurrencyAmount{
			Currency: order.Currency,
			Value:    order.AmountSubTotal,
		},
	}
	err = ordermod.SetOrderChannelMetaData(&topupOrder, constants.ChannelTypeSrcTorqueConvert, &srcMeta)
	if err != nil {
		return
	}

	return topupOrder, nil
}

func (this *DstTorquePurchaseChannel) Commit(ctx comcontext.Context, order *models.Order) (
	meta.OrderStepResultCode, error,
) {
	conversionModel, err := this.getConversionModel(ctx, *order)
	if err != nil {
		return ordermod.OrderStepResultCodeRetry, err
	}
	topupOrder, err := ordermod.ExecuteOrder(ctx, conversionModel.ReceiveOrderID)
	if err != nil {
		return ordermod.OrderStepResultCodeRetry, err
	}
	if topupOrder.Status != constants.OrderStatusCompleted {
		panic(fmt.Errorf("torque purchase topup order has unxepceted status | status=%v", topupOrder.Status))
	}

	return ordermod.OrderStepResultCodeSuccess, nil
}

func (this *DstTorquePurchaseChannel) getConversionModel(ctx comcontext.Context, order models.Order) (
	model models.TorqueCryptoConversion, err error,
) {
	err = database.Atomic(ctx, database.AliasWalletMaster, func(dbTxn *gorm.DB) error {
		err := dbTxn.
			First(&model, &models.TorqueCryptoConversion{ID: order.DstChannelID}).
			Error
		if err != nil {
			return utils.WrapError(err)
		}
		return nil
	})
	return
}

func (this *DstTorquePurchaseChannel) getNotification(
	ctx comcontext.Context, order models.Order, titleKey string, messageKey string,
) (_ *ordermod.Notification, err error) {
	if order.DstChannelRef == "" {
		return nil, nil
	}

	var receiveOrder models.Order
	if order.DstChannelID > 0 {
		conversionModel, err := this.getConversionModel(ctx, order)
		if err != nil {
			return nil, err
		}
		err = database.Atomic(ctx, database.AliasWalletMaster, func(dbTxn *gorm.DB) error {
			err := dbTxn.
				First(&receiveOrder, &models.Order{ID: conversionModel.ReceiveOrderID}).
				Error
			if database.IsDbError(err) {
				return utils.WrapError(err)
			}
			return nil
		})
	}
	if receiveOrder.ID == 0 {
		if receiveOrder, err = this.genTopupOrder(order); err != nil {
			return nil, err
		}
	}

	return this.getSimpleNotification(
		ctx, order,
		titleKey, messageKey,
		meta.O{
			"order":        order,
			"receiveOrder": receiveOrder,
		},
	)
}

func (this *DstTorquePurchaseChannel) GetNotificationCompleted(ctx comcontext.Context, order models.Order) (
	*ordermod.Notification, error,
) {
	return this.getNotification(
		ctx,
		order,
		constants.TranslationKeyOrderNotifCompletedTitleDstTorquePurchase,
		constants.TranslationKeyOrderNotifCompletedMessageDstTorquePurchase,
	)
}

func (this *DstTorquePurchaseChannel) GetNotificationFailed(ctx comcontext.Context, order models.Order) (
	*ordermod.Notification, error,
) {
	return this.getNotification(
		ctx,
		order,
		constants.TranslationKeyOrderNotifFailedTitleDstTorquePurchase,
		constants.TranslationKeyOrderNotifFailedMessageDstTorquePurchase,
	)
}
