package middleware

import (
	"strings"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
)

func SentryPrepare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cloneHub := sentry.CurrentHub().Clone()
		c.Set(comlogging.ContextKeySentryHub, cloneHub)

		scope := cloneHub.Scope()
		scope.SetRequest(c.Request())
		scope.AddEventProcessor(sentryUpdateEventUserProcessor)

		return next(c)
	}
}

func sentryUpdateEventUserProcessor(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
	reqHeaders := event.Request.Headers
	for _, header := range constants.RequestIpHeaders {
		if value, ok := reqHeaders[header]; ok && value != "" {
			event.User.IPAddress = strings.Split(value, ",")[0]
			break
		}
	}

	return event
}
