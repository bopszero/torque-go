package orderchannel

import (
	"reflect"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/ordermod"
)

func init() {
	comutils.PanicOnError(
		ordermod.ChannelRegister(&SrcTorquePurchaseChannel{}),
	)
}

type SrcTorquePurchaseChannel struct {
	baseChannel
}

type SrcTorquePurchaseMeta struct {
	ExchangeRate   decimal.Decimal     `json:"exchange_rate" validate:"required"`
	CurrencyAmount meta.CurrencyAmount `json:"currency_amount" validate:"required"`
}

func (this *SrcTorquePurchaseChannel) GetType() meta.ChannelType {
	return constants.ChannelTypeSrcTorqueConvert
}

func (this *SrcTorquePurchaseChannel) GetMetaType() reflect.Type {
	return reflect.TypeOf(SrcTorquePurchaseMeta{})
}

func (this *SrcTorquePurchaseChannel) GetNotificationCompleted(ctx comcontext.Context, order models.Order) (
	*ordermod.Notification, error,
) {
	return nil, nil
}
