package orderchannel

import (
	"reflect"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/usermod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/ordermod"
)

func init() {
	comutils.PanicOnError(
		ordermod.ChannelRegister(&DstBlockchainNetworkChannel{}),
	)
}

type DstBlockchainNetworkChannel struct {
	baseChannel
}

type DstBlockchainNetworkMeta struct {
	ToAddress    string                `json:"to_address" validate:"required"`
	FeeInfo      blockchainmod.FeeInfo `json:"fee_info"`
	InputDataHex string                `json:"input_data_hex,omitempty"`
}

type DstBlockchainNetworkCheckoutInfo struct {
	FeeInfo      blockchainmod.FeeInfo `json:"fee_info"`
	MinTxnAmount decimal.Decimal       `json:"min_txn_amount"`
}

func (this *DstBlockchainNetworkChannel) GetType() meta.ChannelType {
	return constants.ChannelTypeDstBlockchainNetwork
}

func (this *DstBlockchainNetworkChannel) GetMetaType() reflect.Type {
	return reflect.TypeOf(DstBlockchainNetworkMeta{})
}

func (this *DstBlockchainNetworkChannel) getMeta(order *models.Order) (*DstBlockchainNetworkMeta, error) {
	var metaModel DstBlockchainNetworkMeta
	if err := ordermod.GetOrderChannelMetaData(order, this.GetType(), &metaModel); err != nil {
		return nil, err
	}

	return &metaModel, nil
}

func (this *DstBlockchainNetworkChannel) GetCheckoutInfo(ctx comcontext.Context, order *models.Order) (
	interface{}, error,
) {
	metaModel, err := this.getMeta(order)
	if err != nil {
		return nil, err
	}

	coin, err := blockchainmod.GetCoinNative(order.Currency)
	if err != nil {
		return nil, err
	}
	account, err := coin.NewAccountSystem(ctx, order.UID)
	if err != nil {
		return nil, err
	}
	feeInfo, err := account.GetFeeInfoToAddress(metaModel.ToAddress)
	if err != nil {
		return nil, err
	}

	checkoutInfo := DstBlockchainNetworkCheckoutInfo{
		FeeInfo:      feeInfo,
		MinTxnAmount: coin.GetMinTxnAmount(),
	}
	return &checkoutInfo, nil
}

func (this *DstBlockchainNetworkChannel) PreValidate(ctx comcontext.Context, order *models.Order) error {
	if err := validateBlockchainOrder(ctx, this, *order); err != nil {
		return err
	}

	metaModel, err := this.getMeta(order)
	if err != nil {
		return err
	}

	coin, err := blockchainmod.GetCoinNative(order.Currency)
	if err != nil {
		return err
	}
	if _, err := coin.NormalizeAddress(metaModel.ToAddress); err != nil {
		return utils.WrapError(constants.ErrorAddress)
	}

	if ordermod.IsUserActionLocked(ctx, order.UID, constants.LockActionSendPersonal) {
		return utils.WrapError(constants.ErrorUserActionLocked)
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

func (this *DstBlockchainNetworkChannel) Execute(ctx comcontext.Context, order *models.Order) (
	meta.OrderStepResultCode, error,
) {
	metaModel, err := this.getMeta(order)
	if err != nil {
		return ordermod.OrderStepResultCodeFail, err
	}

	cryptoTxn, code, err := executeBlockchainTxn(ctx, order, metaModel.ToAddress, &metaModel.FeeInfo)
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
		order.DstChannelID = cryptoTxn.ID
	}

	return code, err
}

func (this *DstBlockchainNetworkChannel) ExecuteReverse(ctx comcontext.Context, order *models.Order) (
	meta.OrderStepResultCode, error,
) {
	if err := removeLocalBlockchainTxn(ctx, order); err != nil {
		return ordermod.OrderStepResultCodeRetry, err
	}

	return ordermod.OrderStepResultCodeSuccess, nil
}

func (this *DstBlockchainNetworkChannel) Commit(ctx comcontext.Context, order *models.Order) (
	meta.OrderStepResultCode, error,
) {
	if usermod.IsUserContext(ctx) {
		return ordermod.OrderStepResultCodeRetry, nil
	}

	return watchBlockchainTxnConfirmations(ctx, order)
}

func (this *DstBlockchainNetworkChannel) getNotification(
	ctx comcontext.Context, order models.Order, titleKey string, messageKey string,
) (*ordermod.Notification, error) {
	metaModel, err := this.getMeta(&order)
	if err != nil {
		return nil, err
	}
	if order.DstChannelRef == "" {
		return nil, nil
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

func (this *DstBlockchainNetworkChannel) GetNotificationCompleted(ctx comcontext.Context, order models.Order) (
	*ordermod.Notification, error,
) {
	return this.getNotification(
		ctx, order,
		constants.TranslationKeyOrderNotifCompletedTitleDstBlockchain,
		constants.TranslationKeyOrderNotifCompletedMessageDstBlockchain,
	)
}

func (this *DstBlockchainNetworkChannel) GetNotificationFailed(ctx comcontext.Context, order models.Order) (
	*ordermod.Notification, error,
) {
	return this.getNotification(
		ctx, order,
		constants.TranslationKeyOrderNotifFailedTitleDstBlockchain,
		constants.TranslationKeyOrderNotifFailedMessageDstBlockchain,
	)
}
