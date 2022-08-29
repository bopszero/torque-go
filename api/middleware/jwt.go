package middleware

import (
	"crypto/ecdsa"
	"strings"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/authmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

const (
	JwtHeader           = echo.HeaderAuthorization
	JwtCookieAccessKey  = "jwta"
	JwtCookieRefreshKey = "jwtr"
)

type JwtExConfig struct {
	middleware.JWTConfig

	TokenLookups     []string
	RequireAuthorize bool
}

var (
	JwtDefaultConfig = JwtExConfig{
		JWTConfig: middleware.JWTConfig{
			Skipper:       middleware.DefaultJWTConfig.Skipper,
			SigningMethod: authmod.GetAccessSigningMethod().Alg(),
			ContextKey:    config.ContextKeyJWT,
			AuthScheme:    "Bearer",
			Claims:        JwtNewSimpleMapClaims(),
		},
		TokenLookups: []string{
			"cookie:" + JwtCookieAccessKey,
			"header:" + JwtHeader,
		},
	}
)

func JwtGenConfig(secret interface{}, signingMethod string, requireAuthorize bool) JwtExConfig {
	return JwtExConfig{
		JWTConfig: middleware.JWTConfig{
			SigningMethod: signingMethod,
			SigningKey:    secret,
			ContextKey:    config.ContextKeyJWT,
			ErrorHandlerWithContext: func(error, echo.Context) error {
				return constants.ErrorAuth
			},
		},
		RequireAuthorize: requireAuthorize,
	}
}

func JwtRefreshGenConfig(secret interface{}, signingMethod string) JwtExConfig {
	return JwtExConfig{
		JWTConfig: middleware.JWTConfig{
			SigningMethod: signingMethod,
			SigningKey:    secret,
			ContextKey:    config.ContextKeyJWT,
			ErrorHandlerWithContext: func(error, echo.Context) error {
				return constants.ErrorAuth
			},
		},
		RequireAuthorize: true,
		TokenLookups: []string{
			"cookie:" + JwtCookieRefreshKey,
			"header:" + JwtHeader,
		},
	}
}

func JwtNewWithConfig(exConf JwtExConfig) echo.MiddlewareFunc {
	exConf = jwtPatchConfigDefaultValues(exConf)

	jwtKeyFunc := func(t *jwt.Token) (key interface{}, err error) {
		// Check the signing method
		if t.Method.Alg() != exConf.SigningMethod {
			err = utils.IssueErrorf("unexpected jwt signing method `%v`", t.Method.Alg())
			return
		}
		if len(exConf.SigningKeys) > 0 {
			if kid, ok := t.Header["kid"].(string); ok {
				if idKey, ok := exConf.SigningKeys[kid]; ok {
					key = idKey
				}
			}
			if key == nil {
				err = utils.IssueErrorf("unexpected jwt key id `%v`", t.Header["kid"])
				return
			}
		} else {
			key = exConf.SigningKey
		}

		switch keyT := key.(type) {
		case string:
			key = []byte(keyT)
			break
		case *ecdsa.PrivateKey:
			key = keyT.Public()
			break
		default:
			break
		}
		return
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if exConf.Skipper(c) {
				return next(c)
			}
			if exConf.BeforeFunc != nil {
				exConf.BeforeFunc(c)
			}

			token := jwtExtractToken(c, exConf)
			if token == "" {
				return jwtHandleError(c, exConf, constants.ErrorAuth)
			}

			tokenInfo, err := authmod.JwtParser.ParseWithClaims(token, &authmod.JwtExClaims{}, jwtKeyFunc)
			if err != nil || !tokenInfo.Valid {
				return jwtHandleError(c, exConf, err)
			}
			if exConf.RequireAuthorize {
				if err := jwtAuthorizeClaims(tokenInfo.Claims); err != nil {
					return jwtHandleError(c, exConf, err)
				}
			}

			c.Set(exConf.ContextKey, tokenInfo)
			if exConf.SuccessHandler != nil {
				exConf.SuccessHandler(c)
			}
			return next(c)
		}
	}
}

func jwtPatchConfigDefaultValues(exConf JwtExConfig) JwtExConfig {
	if exConf.Skipper == nil {
		exConf.Skipper = JwtDefaultConfig.Skipper
	}
	if exConf.SigningKey == nil && len(exConf.SigningKeys) == 0 {
		panic("jwt middleware requires signing key")
	}
	if exConf.SigningMethod == "" {
		exConf.SigningMethod = JwtDefaultConfig.SigningMethod
	}
	if exConf.ContextKey == "" {
		exConf.ContextKey = config.ContextKeyJWT
	}
	if exConf.Claims == nil {
		exConf.Claims = JwtDefaultConfig.Claims
	}
	if len(exConf.TokenLookups) == 0 {
		exConf.TokenLookups = JwtDefaultConfig.TokenLookups
	}
	if exConf.AuthScheme == "" {
		exConf.AuthScheme = JwtDefaultConfig.AuthScheme
	}

	return exConf
}

func jwtHandleError(c echo.Context, exConf JwtExConfig, err error) error {
	if exConf.ErrorHandler != nil {
		return exConf.ErrorHandler(err)
	}
	if exConf.ErrorHandlerWithContext != nil {
		return exConf.ErrorHandlerWithContext(err, c)
	}
	return utils.WrapError(constants.ErrorAuth)
}

func jwtAuthorizeClaims(claims jwt.Claims) error {
	exClaims, err := authmod.ParseExClaim(claims)
	if err != nil {
		return err
	}
	if !exClaims.IsAuthorized() {
		return utils.WrapError(constants.ErrorAuth)
	}
	return nil
}

func jwtExtractToken(c echo.Context, config JwtExConfig) (token string) {
	for _, lookup := range config.TokenLookups {
		lookupParts := strings.Split(lookup, ":")
		if len(lookupParts) != 2 {
			continue
		}
		var (
			location = lookupParts[0]
			index    = lookupParts[1]
		)
		switch location {
		case "header":
			token = jwtExtractTokenInHeader(c, config, index)
			break
		case "cookie":
			token = jwtExtractTokenInCookie(c, config, index)
			break
		default:
			break
		}
		if token != "" {
			return token
		}
	}
	return ""
}

func jwtExtractTokenInHeader(c echo.Context, config JwtExConfig, index string) string {
	authHeader := c.Request().Header.Get(index)
	if authHeader == "" {
		return ""
	}

	authParts := strings.Split(authHeader, " ")
	if len(authParts) != 2 {
		return ""
	}

	var (
		scheme = authParts[0]
		token  = authParts[1]
	)
	if scheme != config.AuthScheme {
		return ""
	}

	return token
}

func jwtExtractTokenInCookie(c echo.Context, config JwtExConfig, index string) string {
	cookie, err := c.Cookie(index)
	if err != nil {
		return ""
	} else {
		return cookie.Value
	}
}
