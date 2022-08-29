package general

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/torque-go/api"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/api/responses"
	"gitlab.com/snap-clickstaff/torque-go/api/services/general/graphql"
	apiV1 "gitlab.com/snap-clickstaff/torque-go/api/services/general/v1"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"

	echoSwagger "github.com/swaggo/echo-swagger"
)

// @BasePath /api

func InitRouter(e *echo.Echo) {
	initMetaRoutes(e)

	apiGroup := e.Group("/api")
	apiGroup.GET("/ping/", ping)

	apiV1.InitGroup(apiGroup.Group("/v1"))
}

func ping(c echo.Context) error {
	responseData := meta.O{
		"message": "pong",
	}

	return responses.Ok(apiutils.EchoWrapContext(c), responseData)
}

func initMetaRoutes(e *echo.Echo) {
	e.GET("/health/", apiutils.HealthCheck)
	e.GET("/swagger/*any", echoSwagger.WrapHandler)

	api.InitGraphQLRoutes(e, graphql.GenSchema())
}
