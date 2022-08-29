package orderchannel

import (
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/ordermod"
)

func init() {
	comutils.PanicOnError(
		ordermod.ChannelRegister(&DstSystemChannel{}),
	)
}

type DstSystemChannel struct {
	baseChannel
}

func (this *DstSystemChannel) GetType() meta.ChannelType {
	return constants.ChannelTypeDstSystem
}

func (this *DstSystemChannel) GetNotificationCompleted(ctx comcontext.Context, order models.Order) (
	*ordermod.Notification, error,
) {
	return this.getSimpleNotification(
		ctx, order,
		constants.TranslationKeyOrderNotifCompletedTitleDstSystem,
		constants.TranslationKeyOrderNotifCompletedMessageDstSystem,
		meta.O{"order": order},
	)
}
