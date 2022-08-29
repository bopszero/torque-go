package v1

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/torque-go/api"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/api/responses"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/auth/kycmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/ordermod"
	"gorm.io/gorm"
)

func S2sPortfolioGetCurrency(c echo.Context) error {
	var reqModel S2sPortfolioGetCurrencyRequest
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}
	return portfolioGetCurrencyByUID(c, reqModel.UID, reqModel.Currency, reqModel.Network)
}

func S2sPaymentInitOrder(c echo.Context) (err error) {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel S2sPaymentInitOrderRequest
	)
	if err = api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}
	uid := reqModel.UID

	order := ordermod.NewUserOrder(
		uid, reqModel.Currency,
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
		if err = ordermod.SetOrderChannelMetaData(&order, order.SrcChannelType, reqModel.SrcChannelContext); err != nil {
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
		if err = ordermod.SetOrderChannelMetaData(&order, order.DstChannelType, reqModel.DstChannelContext); err != nil {
			return err
		}
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
		S2sPaymentInitOrderResponse{
			OrderID: order.ID,
		},
	)
}

func S2sPaymentExecuteOrder(c echo.Context) (err error) {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel S2sPaymentExecuteOrderRequest
	)
	if err = api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}
	uid := reqModel.UID

	var (
		order     models.Order
		execError error
	)
	err = database.TransactionFromContext(ctx, database.AliasWalletMaster, func(dbTxn *gorm.DB) error {
		err := dbTxn.
			Select(models.OrderColID).
			First(&order, &models.Order{ID: reqModel.OrderID, UID: uid}).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				err = constants.ErrorOrderNotFound
			}
			return err
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

	var responseOrder S2sPaymentExecuteOrderResponse
	if err = dumpOrder(&order, (*Order)(&responseOrder)); err != nil {
		return
	}

	return responses.Ok(ctx, responseOrder)
}

func S2sMetaHandshake(c echo.Context) (err error) {
	var (
		ctx = apiutils.EchoWrapContext(c)
	)
	metaFeatures, err := kycmod.GetMetaFeaturesSetting()
	if err != nil {
		return err
	}
	return responses.Ok(
		ctx,
		MetaHandshakeResponse{
			Features:           metaFeatures,
			BlockchainNetworks: metaGetBlockchainNetworks(),
			NetworkCurrencies:  metaGetNetworkCurrencies(),
		},
	)
}

func S2sMetaCurrencyInfoGet(c echo.Context) error {
	var (
		ctx = apiutils.EchoWrapContext(c)
	)
	return responses.Ok(
		ctx,
		MetaCurrencyInfoResponse{
			Currencies: metaGetCurrencyInfo(ctx),
		},
	)
}

func S2sRiskActionLockGet(c echo.Context) (err error) {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel S2sGetUserActionLockRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}
	isCheck := ordermod.IsUserActionLocked(ctx, reqModel.UID, reqModel.ActionType)
	return responses.Ok(
		ctx,
		S2sGetUserActionLockResponse{
			IsLocked: isCheck,
		},
	)
}
