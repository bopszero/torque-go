package orderchannel

import (
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/ordermod"
)

func init() {
	for _, channelType := range constants.MerchantChannelTypes {
		merchantChannel := DstMerchantChannel{merchantChannelType: channelType}
		comutils.PanicOnError(
			ordermod.ChannelRegister(&merchantChannel),
		)
	}
}

type DstMerchantChannel struct {
	baseChannel
	merchantChannelType meta.ChannelType
}

func (this *DstMerchantChannel) GetType() meta.ChannelType {
	return this.merchantChannelType
}

func (this *DstMerchantChannel) PreValidate(ctx comcontext.Context, order *models.Order) error {
	if ordermod.IsUserActionLocked(ctx, order.UID, constants.LockActionTORQPurchase) {
		return utils.WrapError(constants.ErrorUserActionLocked)
	}
	return nil
}
