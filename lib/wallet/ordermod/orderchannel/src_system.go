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
		ordermod.ChannelRegister(&SrcSystemChannel{}),
	)
}

type SrcSystemChannel struct {
	baseChannel
}

func (this *SrcSystemChannel) GetType() meta.ChannelType {
	return constants.ChannelTypeSrcSystem
}

func (this *SrcSystemChannel) GetNotificationCompleted(ctx comcontext.Context, order models.Order) (
	*ordermod.Notification, error,
) {
	return this.getSimpleNotification(
		ctx, order,
		constants.TranslationKeyOrderNotifCompletedTitleSrcSystem,
		constants.TranslationKeyOrderNotifCompletedMessageSrcSystem,
		meta.O{"order": order},
	)
}
