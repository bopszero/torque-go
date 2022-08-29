package orderchannel

import (
	"fmt"
	"reflect"
	"time"

	"github.com/jinzhu/copier"
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/thirdpartymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/balancemod"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/ordermod"
	"gorm.io/gorm"
)

const DstProfitWithdrawTomorrowOffsetDuration = 18 * time.Hour

func init() {
	comutils.PanicOnError(
		ordermod.ChannelRegister(&DstProfitWithdrawChannel{}),
	)
}

type DstProfitWithdrawChannel struct {
	baseChannel
}

type DstProfitWithdrawMeta struct {
	Currency      meta.Currency          `json:"currency" validate:"required"`
	Network       meta.BlockchainNetwork `json:"network"`
	Address       string                 `json:"address" validate:"required"`
	ExchangeRate  decimal.Decimal        `json:"exchange_rate"`
	ReceiveAmount decimal.Decimal        `json:"receive_amount"`
}

type DstProfitWithdrawCheckoutInfo struct {
	ConversionRate meta.CurrencyConversionRate `json:"conversion_rate"`
	ReceiveFee     meta.CurrencyAmount         `json:"receive_fee"`
}

type DstProfitWithdrawOrderDetail struct {
	Code   string `json:"code"`
	Status string `json:"status"`

	BlockchainHash        string `json:"blockchain_hash"`
	BlockchainExplorerURL string `json:"blockchain_explorer_url"`
}

func (this *DstProfitWithdrawChannel) GetType() meta.ChannelType {
	return constants.ChannelTypeDstProfitWithdraw
}

func (this *DstProfitWithdrawChannel) GetMetaType() reflect.Type {
	return reflect.TypeOf(DstProfitWithdrawMeta{})
}

func (this *DstProfitWithdrawChannel) getMeta(order *models.Order) (*DstProfitWithdrawMeta, error) {
	var metaModel DstProfitWithdrawMeta
	if err := ordermod.GetOrderChannelMetaData(order, this.GetType(), &metaModel); err != nil {
		return nil, err
	}

	return &metaModel, nil
}

func (this *DstProfitWithdrawChannel) getTorqueTxn(ctx comcontext.Context, order *models.Order) (
	*models.TorqueTxn, error,
) {
	var torqueTxn models.TorqueTxn
	err := database.TransactionFromContext(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) error {
		return balancemod.GetBalanceDB(dbTxn).
			FilterProfitWithdraw(order.DstChannelID).
			First(&torqueTxn).
			Error
	})
	if err != nil {
		return nil, err
	}

	return &torqueTxn, nil
}

func (this *DstProfitWithdrawChannel) GetCheckoutInfo(ctx comcontext.Context, order *models.Order) (
	interface{}, error,
) {
	metaModel, err := this.getMeta(order)
	if err != nil {
		return nil, err
	}
	coin, err := blockchainmod.GetCoin(metaModel.Currency, metaModel.Network)
	if err != nil {
		return nil, err
	}

	networkCurrencyInfo := currencymod.GetNetworkCurrencyInfoFastF(
		coin.GetCurrency(), coin.GetNetwork())
	priceUsdt := currencymod.GetCurrencyPriceUsdtFastF(metaModel.Currency)
	priceMarkup := this.getCurrencyPriceMarkup(metaModel.Currency)
	if priceMarkup != nil {
		priceUsdt = priceMarkup.For(priceUsdt)
	}

	checkoutInfo := DstProfitWithdrawCheckoutInfo{
		ReceiveFee: meta.CurrencyAmount{
			Currency: networkCurrencyInfo.Currency,
			Value:    networkCurrencyInfo.WithdrawalFee,
		},
		ConversionRate: meta.CurrencyConversionRate{
			FromCurrency: metaModel.Currency,
			ToCurrency:   constants.CurrencyTorque,
			Value:        currencymod.ConvertUsdtToTorque(priceUsdt),
		},
	}

	return &checkoutInfo, nil
}

func (this *DstProfitWithdrawChannel) Init(ctx comcontext.Context, order *models.Order) (err error) {
	metaModel, err := this.getMeta(order)
	if err != nil {
		return
	}
	coin, err := blockchainmod.GetCoin(metaModel.Currency, metaModel.Network)
	if err != nil {
		return
	}
	if !coin.IsAvailable() {
		return utils.WrapError(constants.ErrorFeatureNotSupport)
	}
	if _, err = coin.NormalizeAddress(metaModel.Address); err != nil {
		return utils.WrapError(constants.ErrorAddress)
	}

	return nil
}

func (this *DstProfitWithdrawChannel) GetOrderDetails(ctx comcontext.Context, order *models.Order) (
	_ interface{}, err error,
) {
	if order.DstChannelID == 0 {
		return
	}

	var torqueTxn models.TorqueTxn
	err = database.GetDbSlave().
		First(&torqueTxn, &models.TorqueTxn{ID: order.DstChannelID}).
		Error
	if database.IsDbError(err) {
		return
	}
	if torqueTxn.ID == 0 {
		return
	}

	var detailModel DstProfitWithdrawOrderDetail
	if err = copier.Copy(&detailModel, &torqueTxn); err != nil {
		return
	}
	if torqueTxn.Status == constants.WithdrawStatusApproved {
		coin := blockchainmod.GetCoinF(torqueTxn.Currency, torqueTxn.Network)
		detailModel.BlockchainExplorerURL = coin.GenTxnExplorerURL(torqueTxn.BlockchainHash)
	}

	return &detailModel, nil
}

func (this *DstProfitWithdrawChannel) PreValidate(ctx comcontext.Context, order *models.Order) error {
	if order.Currency != constants.CurrencyTorque {
		return utils.IssueErrorf(
			"profit withdrawal only accept Torque currency, not `%v` currency",
			order.Currency,
		)
	}
	if ordermod.IsUserActionLocked(ctx, order.UID, constants.LockActionSendPersonal) {
		return utils.WrapError(constants.ErrorUserActionLocked)
	}
	return nil
}

func (this *DstProfitWithdrawChannel) Execute(ctx comcontext.Context, order *models.Order) (
	meta.OrderStepResultCode, error,
) {
	metaModel, err := this.getMeta(order)
	if err != nil {
		return ordermod.OrderStepResultCodeFail, err
	}

	var torqueTxn models.TorqueTxn
	err = database.TransactionFromContext(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		torqueTxn, err = balancemod.SubmitProfitWithdraw(
			ctx,
			order.UID, order.AmountSubTotal, metaModel.ExchangeRate,
			metaModel.Currency, metaModel.Network, metaModel.Address,
			this.getCurrencyPriceMarkup(metaModel.Currency))
		if err != nil {
			return err
		}

		err = database.OnCommit(ctx, database.AliasWalletMaster, func() error {
			webClient := thirdpartymod.GetWebServiceSystemClient()
			return webClient.SendEmailProfitWithdraw(ctx, torqueTxn.ID)
		})
		if err != nil {
			return err
		}

		metaModel.ExchangeRate = torqueTxn.ExchangeRate
		metaModel.ReceiveAmount = torqueTxn.CoinAmount
		return ordermod.SetOrderChannelMetaData(order, this.GetType(), metaModel)
	})
	if err != nil {
		return ordermod.OrderStepResultCodeFail, err
	}

	order.DstChannelID = torqueTxn.ID

	return ordermod.OrderStepResultCodeSuccess, err
}

func (this *DstProfitWithdrawChannel) getCurrencyPriceMarkup(currency meta.Currency) *meta.AmountMarkup {
	if this.infoModel.BlockchainNetworkConfig.CurrencyMarkupPriceMap == nil {
		return nil
	}

	priceMarkup, ok := this.infoModel.BlockchainNetworkConfig.CurrencyMarkupPriceMap[currency]
	if !ok {
		return nil
	}

	return &priceMarkup
}

func (this *DstProfitWithdrawChannel) ExecuteReverse(ctx comcontext.Context, order *models.Order) (
	meta.OrderStepResultCode, error,
) {
	torqueTxn, err := this.getTorqueTxn(ctx, order)
	if err != nil {
		return ordermod.OrderStepResultCodeRetry, err
	}

	switch torqueTxn.Status {
	case constants.WithdrawStatusPendingConfirm:
		err = balancemod.CancelProfitWithdraw(ctx, order.DstChannelID, "")
	case constants.WithdrawStatusPendingTransfer:
		err = balancemod.RejectProfitWithdraw(ctx, order.DstChannelID, "")
	case constants.WithdrawStatusCanceled, constants.WithdrawStatusRejected:
		err = nil
	default:
		err = utils.IssueErrorf("cannot rollback profit withdraw torque with status `%v`", torqueTxn.Status)
	}
	if err != nil {
		return ordermod.OrderStepResultCodeRetry, err
	}

	return ordermod.OrderStepResultCodeSuccess, nil
}

func (this *DstProfitWithdrawChannel) Commit(ctx comcontext.Context, order *models.Order) (
	meta.OrderStepResultCode, error,
) {
	torqueTxn, err := this.getTorqueTxn(ctx, order)
	if err != nil {
		return ordermod.OrderStepResultCodeRetry, err
	}

	switch torqueTxn.Status {
	case constants.WithdrawStatusApproved:
		return ordermod.OrderStepResultCodeSuccess, nil
	case constants.WithdrawStatusCanceled, constants.WithdrawStatusRejected:
		return ordermod.OrderStepResultCodeFail, fmt.Errorf("profit withdraw fail with reason '%s'", torqueTxn.Note)
	default:
		now := time.Now()
		orderTomorrowTime := comutils.TimeRoundDate(time.Unix(order.CreateTime, 0)).Add(24 * time.Hour)
		orderApprovalTime := orderTomorrowTime.Add(DstProfitWithdrawTomorrowOffsetDuration)
		if now.Before(orderApprovalTime) {
			order.RetryTime = now.Add(3 * time.Minute).Unix()
		}

		return ordermod.OrderStepResultCodeRetry, nil
	}
}

func (this *DstProfitWithdrawChannel) getNotification(
	ctx comcontext.Context, order models.Order, titleKey string, messageKey string,
) (*ordermod.Notification, error) {
	metaModel, err := this.getMeta(&order)
	if err != nil {
		return nil, err
	}

	return this.getSimpleNotification(
		ctx, order,
		titleKey, messageKey,
		meta.O{
			"order": order,
			"meta":  metaModel,
		},
	)
}

func (this *DstProfitWithdrawChannel) GetNotificationCompleted(ctx comcontext.Context, order models.Order) (
	*ordermod.Notification, error,
) {
	return this.getNotification(
		ctx, order,
		constants.TranslationKeyOrderNotifCompletedTitleDstProfitWithdraw,
		constants.TranslationKeyOrderNotifCompletedMessageDstProfitWithdraw,
	)
}

func (this *DstProfitWithdrawChannel) GetNotificationFailed(ctx comcontext.Context, order models.Order) (
	*ordermod.Notification, error,
) {
	return this.getNotification(
		ctx, order,
		constants.TranslationKeyOrderNotifFailedTitleDstProfitWithdraw,
		constants.TranslationKeyOrderNotifFailedMessageDstProfitWithdraw,
	)
}
