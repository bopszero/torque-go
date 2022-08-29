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
		// jwtMiddleware = middleware.JwtNewWithConfig(
		// 	middleware.JwtGenConfig(
		// 		authConfig.AccessSecret,
		// 		authmod.GetAccessSigningMethod().Alg(), true,
		// 	),
		// )
	)
	group.Use(corsMiddleware)
	group.Use(middleware.LogRequestDefaultMiddleware)

	metaGroup := group.Group("/meta", jwtLegacyMiddleware)
	addPostWithOptions(metaGroup, "/handshake/", MetaHandshake)
	addPostWithOptions(metaGroup, "/currency-info/", MetaCurrencyInfoGet)

	var (
		portfolioOpenGroup = group.Group("/portfolio")
		portfolioGroup     = portfolioOpenGroup.Group("", jwtLegacyMiddleware)
	)
	addPostWithOptions(portfolioGroup, "/overview/get/", PortfolioGetOverview)
	addPostWithOptions(portfolioGroup, "/currency/get/", PortfolioGetCurrency)
	addPostWithOptions(portfolioGroup, "/currency/order/list/", PortfolioListCurrencyOrder)
	addPostWithOptions(portfolioGroup, "/currency/order/get/", PortfolioGetCurrencyOrder)
	addPostWithOptions(portfolioGroup, "/currency/order/export/", PortfolioExportGenTokenCurrencyOrder)
	portfolioOpenGroup.GET("/currency/order/export/:token/", PortfolioExportDownloadCurrencyOrder)

	paymentGroup := group.Group("/payment", jwtLegacyMiddleware)
	addPostWithOptions(paymentGroup, "/channel/get/", PaymentGetChannel)
	addPostWithOptions(paymentGroup, "/order/checkout/", PaymentCheckoutOrder)
	addPostWithOptions(paymentGroup, "/order/init/", PaymentInitOrder)
	addPostWithOptions(paymentGroup, "/order/execute/", PaymentExecuteOrder)

	helperGroup := group.Group("/helper", jwtLegacyMiddleware)
	addPostWithOptions(helperGroup, "/blockchain/address/validate/", HelperBlockchainValidateAddress)
}

func initS2SGroup(group *echo.Group) {
	macSecret := viper.GetString(config.KeyServiceWalletMacSecret)
	group.Use(middleware.NewValidateMAC(macSecret))

	portfolioGroup := group.Group("/portfolio")
	portfolioGroup.POST("/currency/get/", S2sPortfolioGetCurrency)

	paymentGroup := group.Group("/payment")
	paymentGroup.Use(middleware.LogRequestDefaultMiddleware)
	paymentGroup.POST("/order/init/", S2sPaymentInitOrder)
	paymentGroup.POST("/order/execute/", S2sPaymentExecuteOrder)

	helperGroup := group.Group("/helper")
	helperGroup.POST("/blockchain/address/validate/", HelperBlockchainValidateAddress)

	legacyGroup := group.Group("/legacy")
	legacyGroup.POST("/torque/transaction/", LegacyTorqueTxnGet)
	legacyGroup.POST("/currency/price-map/get/", LegacyCurrencyPriceMapGet)

	metaGroup := group.Group("/meta")
	metaGroup.POST("/handshake/", S2sMetaHandshake)
	metaGroup.POST("/currency-info/", S2sMetaCurrencyInfoGet)

	riskGroup := group.Group("/risk")
	riskGroup.POST("/action-lock/get/", S2sRiskActionLockGet)
}

func addPostWithOptions(
	group *echo.Group, uri string, handler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc,
) {
	group.POST(uri, handler, middlewares...)
	group.OPTIONS(uri, func(c echo.Context) error { return nil }, middlewares...)
}
