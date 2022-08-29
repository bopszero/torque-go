package auth

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/api/middleware"
	_ "gitlab.com/snap-clickstaff/torque-go/api/services/auth/docs"
	apiV1 "gitlab.com/snap-clickstaff/torque-go/api/services/auth/v1"
	"gitlab.com/snap-clickstaff/torque-go/config"
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
