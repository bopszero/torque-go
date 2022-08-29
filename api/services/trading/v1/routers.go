package v1

import (
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/torque-go/api/middleware"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/authmod"
)

func InitGroup(group *echo.Group) {
	initS2CGroup(group.Group(""))
	initS2SGroup(group.Group("/s2s"))
}

func initS2CGroup(group *echo.Group) {
	var (
		authConfig  = authmod.GetAuthConfig()
		corsOrigins = []string{
			viper.GetString(config.KeyServiceWebBaseURL),
		}
	)
	if config.Test {
		corsOrigins = []string{"*"}
	}
	var (
		corsMiddleware = echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
			Skipper:      echoMiddleware.DefaultSkipper,
			AllowOrigins: corsOrigins,
			AllowMethods: []string{"POST"},
			MaxAge:       int(time.Hour),
		})
		jwtLegacyMiddleware = middleware.JwtLegacyNewWithConfig(
			middleware.JwtLegacyGenConfig(authConfig.LegacySecret),
		)
	)
	group.Use(corsMiddleware)
	group.Use(middleware.LogRequestDefaultMiddleware)

	var (
		txnOpenGroup = group.Group("/txn")
		txnGroup     = txnOpenGroup.Group("", jwtLegacyMiddleware)
	)
	addPostWithOptions(txnGroup, "/deposit/export/", TxnDepositExportGenToken)
	txnOpenGroup.GET("/deposit/export/:token/", TxnDepositExportDownload)
	addPostWithOptions(txnGroup, "/withdrawal/export/", TxnWithdrawalExportGenToken)
	txnOpenGroup.GET("/withdrawal/export/:token/", TxnWithdrawalExportDownload)
}

func initS2SGroup(group *echo.Group) {
	macSecret := viper.GetString(config.KeyServiceTradingSecret)
	group.Use(middleware.NewValidateMAC(macSecret))

	// TODO: Should use gRPC
}

func addPostWithOptions(
	group *echo.Group, uri string, handler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc,
) {
	group.POST(uri, handler, middlewares...)
	group.OPTIONS(uri, func(c echo.Context) error { return nil }, middlewares...)
}
