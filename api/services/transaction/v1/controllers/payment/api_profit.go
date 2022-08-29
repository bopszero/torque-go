package payment

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/torque-go/api"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/api/responses"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/trading/withdrawalmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/balancemod"
)

func RejectProfitWithdraw(c echo.Context) error {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel ProfitWithdrawRejectRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}

	_, err := withdrawalmod.RejectProfitWithdraw(ctx, reqModel.ID, reqModel.Note)
	if err != nil {
		return err
	}

	return responses.Ok(ctx, meta.O{"ok": true})
}

func CancelProfitWithdraw(c echo.Context) error {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel ProfitWithdrawCancelRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}

	_, err := withdrawalmod.CancelProfitWithdraw(ctx, reqModel.ID, reqModel.Note)
	if err != nil {
		return err
	}

	return responses.Ok(ctx, meta.O{"ok": true})
}

func ApproveProfitReinvest(c echo.Context) error {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel ProfitReinvestApproveRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}

	_, err := balancemod.ApproveProfitReinvest(ctx, reqModel.ID, reqModel.Note)
	if err != nil {
		return err
	}

	return responses.Ok(ctx, meta.O{"ok": true})
}

func RejectProfitReinvest(c echo.Context) error {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel ProfitReinvestRejectRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}

	_, err := balancemod.RejectProfitReinvest(ctx, reqModel.ID, reqModel.Note)
	if err != nil {
		return err
	}

	return responses.Ok(ctx, meta.O{"ok": true})
}
