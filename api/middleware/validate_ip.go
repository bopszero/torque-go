package middleware

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"strings"
)

func ValidateIpInWhiteList(whiteList []string, handleError func() error) echo.MiddlewareFunc {
	ipHashSet := convertListToHashSet(whiteList)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userIP := utils.GetRequestUserIP(c.Request())
			if !ipHashSet.Contains(userIP) {
				return handleError()
			}
			return next(c)
		}
	}
}

func convertListToHashSet(whiteList []string) comtypes.HashSet {
	ipHashSet := make(comtypes.HashSet)
	for _, ip := range whiteList{
		ip = strings.TrimSpace(ip)
		ipHashSet.Add(ip)
	}
	return ipHashSet
}
