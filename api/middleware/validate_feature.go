package middleware

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/auth/kycmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func ValidateFeatureEnable(featureCode string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var (
				ctx             = apiutils.EchoWrapContext(c)
				isFeatureEnable = isFeatureEnable(ctx, featureCode)
			)
			if !isFeatureEnable {
				return utils.WrapError(constants.ErrorFeatureNotSupport)
			}
			return next(c)
		}
	}
}

func isFeatureEnable(ctx comcontext.Context, code string) bool {
	metaFeaturesSetting, err := kycmod.GetMetaFeaturesSetting()
	if err != nil {
		comlogging.GetLogger().
			WithContext(ctx).
			WithError(err).
			Error("cannot get meta features setting")
	}
	for _, feature := range metaFeaturesSetting {
		if feature.Code == code && feature.IsAvailable {
			return true
		}
	}
	return false
}
