package v1

import (
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/torque-go/api/middleware"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/auth/kycmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/authmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
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
			Skipper:          echoMiddleware.DefaultSkipper,
			AllowCredentials: true,
			AllowOrigins:     corsOrigins,
			AllowMethods:     []string{"POST"},
			MaxAge:           int(time.Hour),
		})
		jwtLegacyMiddleware = middleware.JwtLegacyNewWithConfig(
			middleware.JwtLegacyGenConfig(authConfig.LegacySecret),
		)
		jwtAccessNonAuthMiddleware = middleware.JwtNewWithConfig(
			middleware.JwtGenConfig(
				authConfig.AccessSecret,
				authmod.GetAccessSigningMethod().Alg(), false,
			),
		)
		// jwtAccessMiddleware = middleware.JwtNewWithConfig(
		// 	middleware.JwtGenConfig(
		// 		authConfig.AccessSecret,
		// 		authmod.GetAccessSigningMethod().Alg(), true,
		// 	),
		// )
		jwtRefreshMiddleware = middleware.JwtNewWithConfig(
			middleware.JwtRefreshGenConfig(
				authConfig.RefreshSecret,
				authmod.GetRefreshSigningMethod().Alg(),
			),
		)
	)
	group.Use(corsMiddleware)
	group.Use(middleware.LogRequestDefaultMiddleware)

	loginGroup := group.Group("/login")
	addPostOpts(loginGroup, "/input/prepare/", LoginInputPrepare)
	addPostOpts(loginGroup, "/input/execute/", LoginInputExecute)
	addPostNameOpts(loginGroup, "/input/commit/", LoginInputCommit, RouteNameLoginInputCommit, jwtAccessNonAuthMiddleware)
	addPostNameOpts(loginGroup, "/refresh/", LoginRefresh, RouteNameLoginRefresh, jwtRefreshMiddleware)
	addPostOpts(loginGroup, "/logout/", LoginLogout)
	if config.Debug {
		loginGroup.GET("/input/test/", LoginInputTest)
	}

	addPostOpts(group, "/meta/status/", MetaStatus)
	addPostOpts(group, "/kyc/validate/email/", KycValidateEmail)

	kycJumioGroup := group.Group("/kyc/jumio")
	middlewareIpInWhiteList := middleware.ValidateIpInWhiteList(
		kycmod.GetJumioIpWhiteList(),
		func() error {
			return constants.ErrorAuth
		},
	)
	addPostOpts(kycJumioGroup, "/scan-result/", KycJumioPushScanResult, middlewareIpInWhiteList)
	kycJumioGroup.GET("/redirect-user/", KycJumioRedirectUser)

	var (
		kycGroup                = group.Group("/kyc", jwtLegacyMiddleware)
		middlewareFeatureEnable = middleware.ValidateFeatureEnable(kycmod.FeatureCodeName)
	)
	addPostOpts(kycGroup, "/init/", KycInit, middlewareFeatureEnable)
	addPostOpts(kycGroup, "/init/url/", KycInitUrl, middlewareFeatureEnable)
	addPostOpts(kycGroup, "/submit/", KycSubmit, middlewareFeatureEnable)
	addPostOpts(kycGroup, "/get/", KycGet)
	addPostOpts(kycGroup, "/meta/", KycMeta)
}

func initS2SGroup(group *echo.Group) {
	macSecret := viper.GetString(config.KeyServiceAuthMacSecret)
	group.Use(middleware.NewValidateMAC(macSecret))

	kycGroup := group.Group("/kyc")
	kycGroup.POST("/get/", S2sKycGet)
	kycGroup.POST("/meta/", KycMeta)
	kycGroup.POST("/validate/email/", KycValidateEmail)
	kycGroup.POST("/send/email/", S2sKycSendEmail)
}

func addPostOpts(
	group *echo.Group, uri string, handler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc,
) *echo.Route {
	group.OPTIONS(uri, emptyOptionsHandler, middlewares...)
	return group.POST(uri, handler, middlewares...)
}

func addPostNameOpts(
	group *echo.Group,
	uri string, handler echo.HandlerFunc,
	name string, middlewares ...echo.MiddlewareFunc,
) *echo.Route {
	route := addPostOpts(group, uri, handler, middlewares...)
	if name != "" {
		route.Name = name
	}
	return route
}

func emptyOptionsHandler(c echo.Context) error {
	return nil
}
