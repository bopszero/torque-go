package orderchannel

import (
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/bonuspoolmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/ordermod"
	"gorm.io/gorm"
)

func init() {
	comutils.PanicOnError(
		ordermod.ChannelRegister(&SrcBonusPoolLeaderChannel{}),
	)
}

type SrcBonusPoolLeaderChannel struct {
	baseChannel
}

func (this *SrcBonusPoolLeaderChannel) getDetail(
	ctx comcontext.Context, order *models.Order,
) (detail models.LeaderBonusPoolDetail, err error) {
	err = database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		detail, err = bonuspoolmod.LeaderGetDB(dbTxn).GetDetail(order.SrcChannelID)
		if err != nil {
			return
		}
		if detail.OrderID > 0 && detail.OrderID != order.ID {
			return utils.IssueErrorf(
				"leader bonus pool order has mismatched Detail | order_id=%v,detail_order_id=%v",
				order.ID, detail.OrderID,
			)
		}
		return nil
	})
	return
}

func (this *SrcBonusPoolLeaderChannel) GetType() meta.ChannelType {
	return constants.ChannelTypeSrcBonusPoolLeader
}

func (this *SrcBonusPoolLeaderChannel) PreValidate(ctx comcontext.Context, order *models.Order) (err error) {
	if order.Currency != constants.CurrencyTorque {
		return utils.IssueErrorf("bonus pool leader only accept Torque currency, not `%v` currency", order.Currency)
	}
	_, err = this.getDetail(ctx, order)
	return
}

func (this *SrcBonusPoolLeaderChannel) Execute(ctx comcontext.Context, order *models.Order) (
	_ meta.OrderStepResultCode, err error,
) {
	detail, err := this.getDetail(ctx, order)
	if err != nil {
		return ordermod.OrderStepResultCodeRetry, err
	}
	err = database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) error {
		detail.OrderID = order.ID
		if err = dbTxn.Save(&detail).Error; err != nil {
			return utils.WrapError(err)
		}
		return nil
	})
	if err != nil {
		return ordermod.OrderStepResultCodeRetry, err
	}
	return ordermod.OrderStepResultCodeSuccess, nil
}

func (this *SrcBonusPoolLeaderChannel) Commit(ctx comcontext.Context, order *models.Order) (
	_ meta.OrderStepResultCode, err error,
) {
	detail, err := this.getDetail(ctx, order)
	if err != nil {
		return ordermod.OrderStepResultCodeRetry, err
	}
	_, err = bonuspoolmod.LeaderTryCompleteExecution(ctx, detail.ExecutionID)
	if err != nil {
		return ordermod.OrderStepResultCodeRetry, err
	}
	return ordermod.OrderStepResultCodeSuccess, nil
}

func (this *SrcBonusPoolLeaderChannel) GetNotificationCompleted(ctx comcontext.Context, order models.Order) (
	_ *ordermod.Notification, err error,
) {
	return this.getSimpleNotification(
		ctx, order,
		constants.TranslationKeyOrderNotifCompletedTitleSrcBonusPoolLeader,
		constants.TranslationKeyOrderNotifCompletedMessageSrcBonusPoolLeader,
		meta.O{
			"order": order,
		},
	)
}
