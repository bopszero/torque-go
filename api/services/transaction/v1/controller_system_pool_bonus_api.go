package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/api"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/api/responses"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/bonuspoolmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func SystemPoolBonusLeaderCheckout(c echo.Context) (err error) {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel PoolBonusLeaderCheckoutRequest
	)
	if err = api.BindAndValidate(c, &reqModel); err != nil {
		return
	}

	poolAmount, err := bonuspoolmod.LeaderCalcPoolAmount(reqModel.FromDate, reqModel.ToDate)
	if err != nil {
		return
	}
	tierInfoList, err := bonuspoolmod.LeaderFetchTierInfoList(ctx)
	if err != nil {
		return
	}
	execHash := bonuspoolmod.LeaderCalcExecutionHash(
		reqModel.FromDate, reqModel.ToDate,
		poolAmount, tierInfoList,
	)

	return responses.Ok(
		ctx,
		PoolBonusLeaderCheckoutResponse{
			ExecutionHash: execHash,
			FromDate:      reqModel.FromDate,
			ToDate:        reqModel.ToDate,
			TotalAmount:   poolAmount,
			TierInfoList:  tierInfoList,
		},
	)
}

func SystemPoolBonusLeaderExecute(c echo.Context) (err error) {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel PoolBonusLeaderExecuteRequest
	)
	if err = api.BindAndValidate(c, &reqModel); err != nil {
		return
	}
	totalTierRate := decimal.Zero
	for _, tierInfo := range reqModel.TierInfoList {
		totalTierRate = totalTierRate.Add(tierInfo.Rate)
	}
	if totalTierRate.GreaterThan(comutils.DecimalOne) {
		return utils.WrapError(constants.ErrorInvalidParams)
	}

	var execution models.LeaderBonusPoolExecution
	err = database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		execution, err = bonuspoolmod.LeaderCreateExecution(
			ctx,
			reqModel.ExecutionHash, reqModel.FromDate, reqModel.ToDate,
			reqModel.TotalAmount, reqModel.TierInfoList)
		if err != nil {
			return
		}
		execution, err = bonuspoolmod.LeaderRunExecution(ctx, execution.ID)
		return
	})
	if err != nil {
		return
	}

	return responses.Ok(
		ctx,
		PoolBonusLeaderExecuteResponse{
			Execution: &execution,
		},
	)
}
