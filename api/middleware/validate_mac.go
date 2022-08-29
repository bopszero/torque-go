package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/api/responses"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
)

const (
	MacHeaderPrefix   = "HMAC "
	MacCurrentVersion = "1.0"
)

func NewValidateMAC(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var (
				ctx     = apiutils.EchoWrapContext(c)
				request = c.Request()
			)
			authHeader := request.Header.Get(echo.HeaderAuthorization)
			if authHeader == "" {
				if config.Test {
					return next(c)
				}

				return c.NoContent(http.StatusForbidden)
			}

			if !strings.HasPrefix(authHeader, MacHeaderPrefix) {
				return responses.AutoErrorCode(ctx, constants.ErrorAuth)
			}

			macValue := authHeader[len(MacHeaderPrefix):]
			macParts := strings.Split(macValue, ":")
			if len(macParts) != 3 {
				return responses.AutoErrorCode(ctx, constants.ErrorAuth)
			}

			version, nonce, macHex := macParts[0], macParts[1], macParts[2]
			if version != MacCurrentVersion {
				return responses.AutoErrorCode(ctx, constants.ErrorAuth)
			}

			requestBody := readContextRequestBody(c)
			if !isValidMAC(secret, requestBody, nonce, macHex) {
				comlogging.GetLogger().
					WithContext(ctx).
					WithFields(logrus.Fields{
						"body":  requestBody,
						"nonce": nonce,
						"mac":   macHex,
					}).
					Warning("request MAC mismatched")
				return responses.AutoErrorCode(ctx, constants.ErrorAuth)
			}

			return next(c)
		}
	}
}

func isValidMAC(secret string, body string, nonce string, actualMACHex string) bool {
	actualMAC, err := hex.DecodeString(actualMACHex)
	if err != nil {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(body + nonce))
	expectedMAC := mac.Sum(nil)

	return hmac.Equal(expectedMAC, actualMAC)
}
