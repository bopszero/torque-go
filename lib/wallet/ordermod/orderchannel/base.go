package orderchannel

import (
	"reflect"

	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlocale"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/ordermod"
)

type baseChannel struct {
	infoModel models.ChannelInfo
}

func (this *baseChannel) SetInfoModel(infoModel models.ChannelInfo) {
	this.infoModel = infoModel
}

func (this *baseChannel) IsAvailable() bool {
	return this.infoModel.IsAvailable.Bool
}

func (this *baseChannel) GetInfoModel() models.ChannelInfo {
	return this.infoModel
}

func (this *baseChannel) GetMetaType() reflect.Type {
	return nil
}

func (this *baseChannel) GetCheckoutInfo(ctx comcontext.Context, order *models.Order) (
	interface{}, error,
) {
	return nil, nil
}

func (this *baseChannel) GetOrderDetails(ctx comcontext.Context, order *models.Order) (
	interface{}, error,
) {
	return nil, nil
}

func (this *baseChannel) Init(ctx comcontext.Context, order *models.Order) error {
	return nil
}

func (this *baseChannel) PreValidate(ctx comcontext.Context, order *models.Order) error {
	return nil
}

func (this *baseChannel) Prepare(ctx comcontext.Context, order *models.Order) (
	meta.OrderStepResultCode, error,
) {
	return ordermod.OrderStepResultCodeIgnore, nil
}

func (this *baseChannel) Execute(ctx comcontext.Context, order *models.Order) (
	meta.OrderStepResultCode, error,
) {
	return ordermod.OrderStepResultCodeIgnore, nil
}

func (this *baseChannel) Commit(ctx comcontext.Context, order *models.Order) (
	meta.OrderStepResultCode, error,
) {
	return ordermod.OrderStepResultCodeIgnore, nil
}

func (this *baseChannel) PrepareReverse(ctx comcontext.Context, order *models.Order) (
	meta.OrderStepResultCode, error,
) {
	return ordermod.OrderStepResultCodeIgnore, nil
}

func (this *baseChannel) ExecuteReverse(ctx comcontext.Context, order *models.Order) (
	meta.OrderStepResultCode, error,
) {
	return ordermod.OrderStepResultCodeIgnore, nil
}

func (this *baseChannel) CommitReverse(ctx comcontext.Context, order *models.Order) (
	meta.OrderStepResultCode, error,
) {
	return ordermod.OrderStepResultCodeIgnore, nil
}

func (this *baseChannel) getSimpleNotification(
	ctx comcontext.Context, order models.Order,
	translationTitleKey string, translationMessageKey string, translationData interface{},
) (_ *ordermod.Notification, err error) {
	notif := new(ordermod.Notification)
	notif.Title, err = comlocale.TranslateKeyData(ctx, translationTitleKey, translationData)
	if err != nil {
		return
	}
	notif.Message, err = comlocale.TranslateKeyData(ctx, translationMessageKey, translationData)
	if err != nil {
		return
	}
	return notif, nil
}

func (this *baseChannel) GetNotificationCompleted(ctx comcontext.Context, order models.Order) (
	*ordermod.Notification, error,
) {
	return this.getSimpleNotification(
		ctx, order,
		constants.TranslationKeyOrderNotifCompletedTitleDefault,
		constants.TranslationKeyOrderNotifCompletedMessageDefault,
		meta.O{"order": order},
	)
}

func (this *baseChannel) GetNotificationFailed(ctx comcontext.Context, order models.Order) (
	*ordermod.Notification, error,
) {
	return this.getSimpleNotification(
		ctx, order,
		constants.TranslationKeyOrderNotifFailedTitleDefault,
		constants.TranslationKeyOrderNotifFailedMessageDefault,
		meta.O{"order": order},
	)
}
