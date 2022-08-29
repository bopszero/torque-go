package trading

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/api/middleware"
	apiV1 "gitlab.com/snap-clickstaff/torque-go/api/services/wallet/v1"
	"gitlab.com/snap-clickstaff/torque-go/config"

	echoSwagger "github.com/swaggo/echo-swagger"
	_ "gitlab.com/snap-clickstaff/torque-go/api/services/wallet/docs"
)

func InitRouter(e *echo.Echo) {
	initMetaRoutes(e)

	apiV1.InitGroup(e.Group("/v1", middleware.Locale))
}

func initMetaRoutes(e *echo.Echo) {
	e.GET("/health/", apiutils.HealthCheck)

	if config.Test {
		e.GET("/swagger/*any", echoSwagger.WrapHandler)
	}
}
