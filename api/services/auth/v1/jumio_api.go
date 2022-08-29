package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/torque-go/api"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/api/responses"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/auth/kycmod"
)

func KycJumioPushScanResult(c echo.Context) (err error) {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel JumioPushScanResultRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return responses.NewError(ctx, http.StatusBadRequest, err)
	}
	if err = kycmod.SubmitRequestByCode(ctx, reqModel.RequestCode, reqModel.Reference); err != nil {
		return responses.NewError(ctx, http.StatusInternalServerError, err)
	}
	return responses.OkEmpty(ctx)
}

func KycJumioRedirectUser(c echo.Context) (err error) {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel KycJumioRedirectRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return responses.NewError(ctx, http.StatusBadRequest, err)
	}
	if err = kycmod.SubmitRequestByCode(ctx, reqModel.RequestCode, reqModel.Reference); err != nil {
		return responses.NewError(ctx, http.StatusInternalServerError, err)
	}
	landingURL := viper.GetString(config.KeyServiceWebBaseURL) + "/user/personal_info/"
	return c.Redirect(http.StatusSeeOther, landingURL)
}
