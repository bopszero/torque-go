package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gitlab.com/snap-clickstaff/go-common/comlocale"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
)

func Locale(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		reqHeader := c.Request().Header
		acceptLanguage := reqHeader.Get(config.HttpHeaderTorqueLanguage)
		if acceptLanguage == "" {
			acceptLanguage = reqHeader.Get(config.HttpHeaderAcceptLanguage)
		}

		if acceptLanguage != "" && acceptLanguage != config.LanguageCode {
			bundle, err := comlocale.GetBundle()
			if err != nil {
				comlogging.GetLogger().
					WithContext(apiutils.EchoWrapContext(c)).
					WithError(err).
					Error("comlocale get bundle failed | err=%s", err.Error())
			} else {
				localizer := i18n.NewLocalizer(bundle, acceptLanguage)
				c.Set(comlocale.ContextKeyLocalizer, localizer)
			}
		}

		return next(c)
	}
}
