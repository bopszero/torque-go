package orderchannel

import (
	"reflect"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/settingmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/usermod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/balancemod"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/ordermod"
	"gorm.io/gorm"
)

func init() {
	comutils.PanicOnError(
		ordermod.ChannelRegister(&DstTransferChannel{}),
	)
}

type DstTransferChannel struct {
	baseChannel
}

type DstTransferMeta struct {
	UserIdentity string `json:"user_identity" validate:"required,printascii"`
	Note         string `json:"note,omitempty"`

	ReceiverUsername string   `json:"receiver_username"`
	ReceiverUID      meta.UID `json:"receiver_uid"`
}

type DstTransferCheckoutInfo struct {
	Fee meta.CurrencyAmount `json:"fee"`
}

func (this *DstTransferChannel) GetType() meta.ChannelType {
	return constants.ChannelTypeDstTransfer
}

func (this *DstTransferChannel) GetMetaType() reflect.Type {
	return reflect.TypeOf(DstTransferMeta{})
}

func (this *DstTransferChannel) GetCheckoutInfo(ctx comcontext.Context, order *models.Order) (
	interface{}, error,
) {
	if order.Currency == constants.CurrencyTorque {
		transferFee, err := this.getTransferFee()
		if err != nil {
			return nil, err
		}
		checkoutInfo := DstTransferCheckoutInfo{
			Fee: meta.CurrencyAmount{
				Currency: constants.CurrencyTorque,
				Value:    transferFee,
			},
		}
		return &checkoutInfo, nil
	}

	return nil, nil
}

func (this *DstTransferChannel) getTransferFee() (decimal.Decimal, error) {
	transferFeeStr, err := settingmod.GetSettingValueFast(constants.SettingKeyTorqueTransferFee)
	if err != nil {
		return decimal.Zero, err
	}
	return decimal.NewFromString(transferFeeStr)
}

func (this *DstTransferChannel) Init(ctx comcontext.Context, order *models.Order) error {
	transferFee, err := this.getTransferFee()
	if err != nil {
		return err
	}

	order.AmountFee = order.AmountFee.Add(transferFee)

	return nil
}

func (this *DstTransferChannel) PreValidate(ctx comcontext.Context, order *models.Order) error {
	if blockchainmod.IsNativeCurrency(order.Currency) {
		return utils.IssueErrorf("cannot transfer p2p with blockchain currency `%v`", order.Currency)
	}

	metaModel, err := this.getMeta(order)
	if err != nil {
		return err
	}

	receiveUser, err := getUserByIdentity(metaModel.UserIdentity)
	if err != nil {
		return err
	}
	if receiveUser.ID == order.UID {
		return utils.WrapError(
			constants.ErrorInvalidParams.WithKey(constants.TranslationKeyTransferNoCircle),
		)
	}
	if receiveUser.Status != constants.UserStatusActive {
		return utils.WrapError(constants.ErrorUserNotFound)
	}

	metaModel.ReceiverUID = receiveUser.ID
	metaModel.ReceiverUsername = receiveUser.Username
	if err := ordermod.SetOrderChannelMetaData(order, this.GetType(), metaModel); err != nil {
		return err
	}
	if ordermod.IsUserActionLocked(ctx, order.UID, constants.LockActionTORQTransfer) {
		return utils.WrapError(constants.ErrorUserActionLocked)
	}

	return nil
}

func (this *DstTransferChannel) getMeta(order *models.Order) (*DstTransferMeta, error) {
	var metaModel DstTransferMeta
	if err := ordermod.GetOrderChannelMetaData(order, this.GetType(), &metaModel); err != nil {
		return nil, err
	}

	return &metaModel, nil
}

func (this *DstTransferChannel) Prepare(ctx comcontext.Context, order *models.Order) (
	_ meta.OrderStepResultCode, err error,
) {
	receiveOrder, err := this.genReceiveOrder(*order)
	if err != nil {
		return ordermod.OrderStepResultCodeFail, err
	}

	err = database.TransactionFromContext(ctx, database.AliasWalletMaster, func(dbTxn *gorm.DB) error {
		err := dbTxn.Save(&receiveOrder).Error
		if err != nil {
			return err
		}

		order.DstChannelID = receiveOrder.ID
		return nil
	})
	if err != nil {
		return ordermod.OrderStepResultCodeFail, err
	}

	return ordermod.OrderStepResultCodeSuccess, nil
}

func (this *DstTransferChannel) genReceiveOrder(order models.Order) (receiveOrder models.Order, err error) {
	metaModel, err := this.getMeta(&order)
	if err != nil {
		return
	}
	senderUser, err := usermod.GetUserFast(order.UID)
	if err != nil {
		return
	}

	receiveOrder = ordermod.NewUserOrder(
		metaModel.ReceiverUID, order.Currency,
		constants.ChannelTypeSrcTransfer, constants.ChannelTypeDstBalance)
	receiveOrder.SrcChannelID = order.ID
	receiveOrder.SrcChannelAmount = order.AmountSubTotal
	receiveOrder.DstChannelAmount = order.AmountSubTotal
	receiveOrder.AmountSubTotal = order.AmountSubTotal
	receiveOrder.AmountTotal = order.AmountSubTotal

	receiveMeta := SrcTransferMeta{
		SenderUID:      senderUser.ID,
		SenderUsername: senderUser.Username,
		Note:           metaModel.Note,
	}
	err = ordermod.SetOrderChannelMetaData(&receiveOrder, constants.ChannelTypeSrcTransfer, receiveMeta)
	if err != nil {
		return
	}

	return receiveOrder, nil
}

func (this *DstTransferChannel) Execute(ctx comcontext.Context, order *models.Order) (
	meta.OrderStepResultCode, error,
) {
	metaModel, err := this.getMeta(order)
	if err != nil {
		return ordermod.OrderStepResultCodeRetry, err
	}

	receiveOrder, err := ordermod.ExecuteOrder(ctx, order.DstChannelID)
	if err != nil {
		return ordermod.OrderStepResultCodeRetry, err
	}
	if receiveOrder.Status != constants.OrderStatusCompleted {
		return ordermod.OrderStepResultCodeRetry, constants.ErrorOrderStatus
	}

	err = database.TransactionFromContext(ctx, database.AliasWalletMaster, func(dbTxn *gorm.DB) error {
		transferTxn, err := balancemod.AddP2pTransfer(
			ctx,
			order.Currency, order.AmountSubTotal,
			order.UID, receiveOrder.UID,
			order.AmountFee, metaModel.Note)
		if err != nil {
			return err
		}

		order.DstChannelRef = comutils.Stringify(transferTxn.ID)
		receiveOrder.SrcChannelRef = comutils.Stringify(transferTxn.ID)

		return dbTxn.Save(&receiveOrder).Error
	})
	if err != nil {
		return ordermod.OrderStepResultCodeRetry, err
	}

	if usermod.IsUserContext(ctx) {
		database.OnCommit(ctx, database.AliasWalletMaster, func() error {
			if err := ordermod.PushOrderCompletedNotificationAsync(receiveOrder); err != nil {
				comlogging.GetLogger().
					WithContext(ctx).
					WithError(err).
					WithField("order_id", receiveOrder.ID).
					Warn("p2p transfer push receiver notification failed")
			}

			return nil
		})
	}

	return ordermod.OrderStepResultCodeSuccess, nil
}

func (this *DstTransferChannel) ExecuteReverse(ctx comcontext.Context, order *models.Order) (
	meta.OrderStepResultCode, error,
) {
	err := utils.IssueErrorf("cannot reverse a transfer order, kindly use adjust balance for it | order_id=%v", order.ID)
	return ordermod.OrderStepResultCodeFail, err
}
