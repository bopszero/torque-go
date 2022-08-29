package orderchannel

import (
	"reflect"

	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/ordermod"
)

func init() {
	comutils.PanicOnError(
		ordermod.ChannelRegister(&SrcTransferChannel{}),
	)
}

type SrcTransferChannel struct {
	baseChannel
}

type SrcTransferMeta struct {
	SenderUID      meta.UID `json:"sender_uid"`
	SenderUsername string   `json:"sender_username" validate:"printascii"`
	Note           string   `json:"note,omitempty"`
}

func (this *SrcTransferChannel) GetType() meta.ChannelType {
	return constants.ChannelTypeSrcTransfer
}

func (this *SrcTransferChannel) GetMetaType() reflect.Type {
	return reflect.TypeOf(SrcTransferMeta{})
}

func (this *SrcTransferChannel) getMeta(order *models.Order) (*SrcTransferMeta, error) {
	var metaModel SrcTransferMeta
	if err := ordermod.GetOrderChannelMetaData(order, this.GetType(), &metaModel); err != nil {
		return nil, err
	}

	return &metaModel, nil
}

func (this *SrcTransferChannel) GetNotificationCompleted(ctx comcontext.Context, order models.Order) (
	_ *ordermod.Notification, err error,
) {
	metaModel, err := this.getMeta(&order)
	if err != nil {
		return
	}

	return this.getSimpleNotification(
		ctx, order,
		constants.TranslationKeyOrderNotifCompletedTitleSrcTransfer,
		constants.TranslationKeyOrderNotifCompletedMessageSrcTransfer,
		meta.O{
			"order": order,
			"meta":  metaModel,
		},
	)
}
