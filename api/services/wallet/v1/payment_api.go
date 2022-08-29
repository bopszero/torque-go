package v1

import (
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/torque-go/api"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/api/responses"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/usermod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/ordermod"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/ordermod/orderchannel"
	"gorm.io/gorm"
)

func PaymentGetChannel(c echo.Context) (err error) {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel PaymentGetChannelRequest
	)
	if err = api.BindAndValidate(c, &reqModel); err != nil {
		return
	}

	channel, err := orderchannel.GetChannelInfoFast(reqModel.Type)
	if err != nil {
		return
	}

	var response PaymentGetChannelResponse
	if err = copier.Copy(&response, channel); err != nil {
		return
	}

	return responses.AutoErrorCodeData(ctx, err, &response)
}

func PaymentCheckoutOrder(c echo.Context) (err error) {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		uid      = apiutils.GetContextUidF(ctx)
		reqModel PaymentCheckoutOrderRequest
	)
	if err = api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}
	currencyInfo := currencymod.GetCurrencyInfoFastF(reqModel.Currency)
	if !currencymod.IsValidWalletInfo(currencyInfo) {
		return utils.WrapError(constants.ErrorCurrency)
	}

	order := ordermod.NewUserOrder(
		uid, currencyInfo.Currency,
		reqModel.SrcChannelType, reqModel.DstChannelType,
	)

	order.SrcChannelType = reqModel.SrcChannelType
	order.SrcChannelID = reqModel.SrcChannelID
	order.SrcChannelRef = reqModel.SrcChannelRef
	order.SrcChannelAmount = reqModel.SrcChannelAmount

	order.DstChannelType = reqModel.DstChannelType
	order.DstChannelID = reqModel.DstChannelID
	order.DstChannelRef = reqModel.DstChannelRef
	order.DstChannelAmount = reqModel.DstChannelAmount

	if order.DstChannelType == constants.ChannelTypeDstBalance {
		order.AmountSubTotal = reqModel.DstChannelAmount
		order.AmountTotal = reqModel.DstChannelAmount
	} else {
		order.AmountSubTotal = reqModel.SrcChannelAmount
		order.AmountTotal = reqModel.SrcChannelAmount
	}

	srcChannel, err := ordermod.GetChannelByType(order.SrcChannelType)
	if err != nil {
		return
	}
	if srcChannel.GetMetaType() != nil {
		err = ordermod.SetOrderChannelMetaData(&order, order.SrcChannelType, reqModel.SrcChannelContext)
		if err != nil {
			return
		}
	}
	dstChannel, err := ordermod.GetChannelByType(order.DstChannelType)
	if err != nil {
		return
	}
	if dstChannel.GetMetaType() != nil {
		err = ordermod.SetOrderChannelMetaData(&order, order.DstChannelType, reqModel.DstChannelContext)
		if err != nil {
			return
		}
	}

	if err = ordermod.ValidateUserRequestOrder(order); err != nil {
		return
	}

	// srcChannelInfo, err := ordermod.GetChannelInfoFast(order.SrcChannelType)
	// if err != nil {
	// 	return
	// }
	// dstChannelInfo, err := ordermod.GetChannelInfoFast(order.DstChannelType)
	// if err != nil {
	// 	return
	// }
	// var srcChannelInfoResponse, dstChannelInfoResponse Channel
	// if err = copier.Copy(&srcChannelInfoResponse, srcChannelInfo); err != nil {
	// 	return
	// }
	// if err = copier.Copy(&dstChannelInfoResponse, dstChannelInfo); err != nil {
	// 	return
	// }

	var srcInfo, dstInfo interface{}
	if srcInfo, err = srcChannel.GetCheckoutInfo(ctx, &order); err != nil {
		return
	}
	if dstInfo, err = dstChannel.GetCheckoutInfo(ctx, &order); err != nil {
		return
	}

	validationFuncs := []func(ctx comcontext.Context, order *models.Order) error{
		srcChannel.Init,
		srcChannel.PreValidate,
		dstChannel.Init,
		dstChannel.PreValidate,
	}
	for _, validateionFunc := range validationFuncs {
		if err = validateionFunc(ctx, &order); err != nil {
			break
		}
	}

	response := PaymentCheckoutOrderResponse{
		ChannelSrcInfo: srcInfo,
		ChannelDstInfo: dstInfo,
	}
	return responses.AutoErrorCodeData(ctx, err, &response)
}

func PaymentInitOrder(c echo.Context) (err error) {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		uid      = apiutils.GetContextUidF(ctx)
		reqModel PaymentInitOrderRequest
	)
	if err = api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}
	currencyInfo := currencymod.GetCurrencyInfoFastF(reqModel.Currency)
	if !currencymod.IsValidWalletInfo(currencyInfo) {
		return utils.WrapError(constants.ErrorCurrency)
	}

	order := ordermod.NewUserOrder(
		uid, currencyInfo.Currency,
		reqModel.SrcChannelType, reqModel.DstChannelType,
	)

	// TODO: Validate Amounts
	order.AmountSubTotal = reqModel.AmountSubTotal
	order.AmountTotal = reqModel.AmountTotal

	order.SrcChannelType = reqModel.SrcChannelType
	order.SrcChannelID = reqModel.SrcChannelID
	order.SrcChannelRef = reqModel.SrcChannelRef
	order.SrcChannelAmount = reqModel.SrcChannelAmount

	srcChannel, err := ordermod.GetChannelByType(order.SrcChannelType)
	if err != nil {
		return err
	}
	if srcChannel.GetMetaType() != nil {
		err = ordermod.SetOrderChannelMetaData(&order, order.SrcChannelType, reqModel.SrcChannelContext)
		if err != nil {
			return err
		}
	}

	order.DstChannelType = reqModel.DstChannelType
	order.DstChannelID = reqModel.DstChannelID
	order.DstChannelRef = reqModel.DstChannelRef
	order.DstChannelAmount = reqModel.DstChannelAmount

	dstChannel, err := ordermod.GetChannelByType(order.DstChannelType)
	if err != nil {
		return err
	}
	if dstChannel.GetMetaType() != nil {
		err = ordermod.SetOrderChannelMetaData(&order, order.DstChannelType, reqModel.DstChannelContext)
		if err != nil {
			return err
		}
	}

	if err := ordermod.ValidateUserRequestOrder(order); err != nil {
		return err
	}

	err = database.TransactionFromContext(ctx, database.AliasWalletMaster, func(dbTxn *gorm.DB) error {
		if err = dbTxn.Create(&order).Error; err != nil {
			return err
		}
		if order, err = ordermod.InitOrder(ctx, order.ID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return responses.Ok(
		ctx,
		PaymentInitOrderResponse{
			OrderID: order.ID,
		},
	)
}

func PaymentExecuteOrder(c echo.Context) (err error) {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		uid      = apiutils.GetContextUidF(ctx)
		reqModel PaymentExecuteOrderRequest
	)
	if err = api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}

	var (
		order     models.Order
		execError error
	)
	err = database.TransactionFromContext(ctx, database.AliasWalletMaster, func(dbTxn *gorm.DB) error {
		err := dbTxn.
			First(&order, &models.Order{ID: reqModel.OrderID, UID: uid}).
			Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				err = constants.ErrorOrderNotFound
			}
			return err
		}
		// TODO: Enhance this logic to handle TOTP check more robust
		if order.DstChannelType != constants.ChannelTypeDstBalance {
			if !usermod.IsValidUserTOTP(uid, reqModel.AuthCode, true) {
				return constants.ErrorAuthInput
			}
		}
		if order.Status != constants.OrderStatusInit {
			return constants.ErrorOrderStatus
		}

		order, execError = ordermod.ExecuteOrder(ctx, order.ID)
		return nil
	})
	if err != nil {
		return err
	}
	if execError != nil {
		return execError
	}

	var responseOrder PaymentExecuteOrderResponse
	if err = dumpOrder(&order, (*Order)(&responseOrder)); err != nil {
		return
	}

	return responses.Ok(ctx, responseOrder)
}
