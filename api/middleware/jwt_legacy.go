package middleware

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

var (
	JwtLegacyDefaultConfig = middleware.JWTConfig{
		Skipper:       middleware.DefaultJWTConfig.Skipper,
		SigningMethod: jwt.SigningMethodHS256.Alg(),
		ContextKey:    "jwt",
		AuthScheme:    "Bearer",
		Claims:        JwtNewSimpleMapClaims(),
	}
	JwtLegacyParser = jwt.NewParser(
		jwt.WithJSONNumber(),
	)
)

func JwtLegacyGenConfig(secret []byte) middleware.JWTConfig {
	return middleware.JWTConfig{
		SigningMethod: jwt.SigningMethodHS256.Alg(),
		SigningKey:    secret,
		ContextKey:    "jwt",
		SuccessHandler: func(c echo.Context) {
			token := c.Get("jwt").(*jwt.Token)
			tokenMapClaims := token.Claims.(jwt.MapClaims)

			sub := comutils.Stringify(tokenMapClaims["sub"])
			uid, err := comutils.ParseInt64(sub)
			comutils.PanicOnError(err)
			tokenMapClaims["sub"] = sub
			tokenMapClaims["uid"] = uid
		},
		ErrorHandlerWithContext: func(error, echo.Context) error {
			return constants.ErrorAuth
		},
	}
}

func JwtLegacyNewWithConfig(config middleware.JWTConfig) echo.MiddlewareFunc {
	jwtLegacyPatchConfigDefaults(&config)

	jwtKeyFunc := func(t *jwt.Token) (interface{}, error) {
		// Check the signing method
		if t.Method.Alg() != config.SigningMethod {
			return nil, utils.IssueErrorf("unexpected jwt signing method `%v`", t.Method.Alg())
		}
		return config.SigningKey, nil
	}

	extractor := func(c echo.Context) (string, error) {
		authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 {
			return "", constants.ErrorAuth
		}
		scheme, token := authParts[0], authParts[1]
		if scheme != config.AuthScheme {
			return "", constants.ErrorAuth
		}
		return token, nil
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}
			if config.BeforeFunc != nil {
				config.BeforeFunc(c)
			}

			auth, err := extractor(c)
			if err != nil {
				if config.ErrorHandler != nil {
					return config.ErrorHandler(err)
				}
				if config.ErrorHandlerWithContext != nil {
					return config.ErrorHandlerWithContext(err, c)
				}
				return err
			}

			token, err := JwtLegacyParser.Parse(auth, jwtKeyFunc)
			if err == nil && token.Valid {
				c.Set(config.ContextKey, token)
				if config.SuccessHandler != nil {
					config.SuccessHandler(c)
				}
				return next(c)
			}
			if config.ErrorHandler != nil {
				return config.ErrorHandler(err)
			}
			if config.ErrorHandlerWithContext != nil {
				return config.ErrorHandlerWithContext(err, c)
			}

			return &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Message:  "invalid or expired jwt",
				Internal: err,
			}
		}
	}
}

func jwtLegacyPatchConfigDefaults(config *middleware.JWTConfig) {
	if config.Skipper == nil {
		config.Skipper = JwtLegacyDefaultConfig.Skipper
	}
	if config.SigningKey == nil && len(config.SigningKeys) == 0 {
		panic("jwt middleware requires signing key")
	}
	if config.SigningMethod == "" {
		config.SigningMethod = JwtLegacyDefaultConfig.SigningMethod
	}
	if config.ContextKey == "" {
		config.ContextKey = "jwt"
	}
	if config.Claims == nil {
		config.Claims = JwtLegacyDefaultConfig.Claims
	}
	if config.TokenLookup == "" {
		config.TokenLookup = JwtLegacyDefaultConfig.TokenLookup
	}
	if config.AuthScheme == "" {
		config.AuthScheme = JwtLegacyDefaultConfig.AuthScheme
	}
}
