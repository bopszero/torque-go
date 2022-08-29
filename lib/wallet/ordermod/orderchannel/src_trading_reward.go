package orderchannel

import (
	"reflect"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/ordermod"
	"gorm.io/gorm"
)

func init() {
	comutils.PanicOnError(
		ordermod.ChannelRegister(&SrcTradingRewardChannel{}),
	)
}

type SrcTradingRewardChannel struct {
	baseChannel
}

type SrcTradingRewardMeta struct {
	DailyProfitAmounts        SrcTradingRewardDailyProfitAmounts `json:"daily_profit_amounts" validate:"required"`
	AffiliateCommissionAmount decimal.Decimal                    `json:"affiliate_commission_amount"`
	LeaderCommissionAmount    decimal.Decimal                    `json:"leader_commission_amount"`
}

type SrcTradingRewardDailyProfitAmounts []meta.CurrencyAmount

func (this SrcTradingRewardDailyProfitAmounts) Total() decimal.Decimal {
	totalAmount := decimal.Zero
	for _, amount := range this {
		totalAmount = totalAmount.Add(amount.Value)
	}

	return totalAmount
}

func (this *SrcTradingRewardChannel) GetType() meta.ChannelType {
	return constants.ChannelTypeSrcTradingReward
}

func (this *SrcTradingRewardChannel) GetMetaType() reflect.Type {
	return reflect.TypeOf(SrcTradingRewardMeta{})
}

func (this *SrcTradingRewardChannel) getMeta(order *models.Order) (*SrcTradingRewardMeta, error) {
	var metaModel SrcTradingRewardMeta
	if err := ordermod.GetOrderChannelMetaData(order, this.GetType(), &metaModel); err != nil {
		return nil, err
	}

	return &metaModel, nil
}

func (this *SrcTradingRewardChannel) validateExistingOrder(ctx comcontext.Context, order models.Order) (err error) {
	var existingOrder models.Order
	err = database.TransactionFromContext(ctx, database.AliasWalletMaster, func(dbTxn *gorm.DB) error {
		err := dbTxn.
			Where(dbquery.NotEqual(models.OrderColID, order.ID)).
			Where(dbquery.In(
				models.OrderColStatus,
				[]meta.OrderStatus{
					constants.OrderStatusHandleSrc,
					constants.OrderStatusHandleDst,
					constants.OrderStatusCompleting,
					constants.OrderStatusCompleted,
				},
			)).
			First(
				&existingOrder,
				&models.Order{
					UID:            order.UID,
					Currency:       order.Currency,
					SrcChannelType: order.SrcChannelType,
					DstChannelType: order.DstChannelType,
					SrcChannelRef:  order.SrcChannelRef,
				},
			).
			Error
		if database.IsDbError(err) {
			return utils.WrapError(err)
		}

		return nil
	})
	if err != nil {
		return
	}
	if existingOrder.ID > 0 {
		return utils.WrapError(constants.ErrorOrderDuplicated)
	}

	return nil
}

func (this *SrcTradingRewardChannel) PreValidate(ctx comcontext.Context, order *models.Order) (err error) {
	if order.Currency != constants.CurrencyTorque {
		return utils.IssueErrorf("trading reward only accept Torque currency, not `%v` currency", order.Currency)
	}

	_, err = comutils.TimeParse(constants.DateFormatISO, order.SrcChannelRef)
	if err != nil {
		return utils.IssueErrorf(
			"trading reward need a valid date as source reference | value=%v,err=%v",
			order.SrcChannelRef, err,
		)
	}

	metaModel, err := this.getMeta(order)
	if err != nil {
		return err
	}

	totalMetaAmount := metaModel.AffiliateCommissionAmount.
		Add(metaModel.LeaderCommissionAmount).
		Add(metaModel.DailyProfitAmounts.Total())
	if !totalMetaAmount.Equal(order.SrcChannelAmount) {
		return utils.WrapError(constants.ErrorAmount)
	}

	if err = this.validateExistingOrder(ctx, *order); err != nil {
		return
	}

	return nil
}

func (this *SrcTradingRewardChannel) Prepare(ctx comcontext.Context, order *models.Order) (
	_ meta.OrderStepResultCode, err error,
) {
	err = this.validateExistingOrder(ctx, *order)
	if utils.IsOurError(err, constants.ErrorCodeOrderDuplicated) {
		comutils.PanicOnError(
			ordermod.SetOrderFailStatus(order, constants.OrderStatusCanceled),
		)
		return ordermod.OrderStepResultCodeFail, err
	}
	if err != nil {
		return ordermod.OrderStepResultCodeRetry, err
	}
	return ordermod.OrderStepResultCodeSuccess, nil
}

func (this *SrcTradingRewardChannel) GetNotificationCompleted(ctx comcontext.Context, order models.Order) (
	_ *ordermod.Notification, err error,
) {
	amountUSDT := currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: order.Currency,
			Value:    order.DstChannelAmount,
		},
		constants.CurrencyTetherUSD,
	)
	return this.getSimpleNotification(
		ctx, order,
		constants.TranslationKeyOrderNotifCompletedTitleSrcTradingReward,
		constants.TranslationKeyOrderNotifCompletedMessageSrcTradingReward,
		meta.O{
			"order":           order,
			"estimatedAmount": amountUSDT,
		},
	)
}
