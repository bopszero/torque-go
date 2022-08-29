package api

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/api/middleware"
	"gitlab.com/snap-clickstaff/torque-go/api/responses"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func CreateEchoObject() *echo.Echo {
	e := echo.New()
	logger := comlogging.GetLogger()

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		var (
			responseErr error
			ctx         = apiutils.EchoWrapContext(c)
		)
		if !c.Response().Committed {
			responseErr = responses.AutoErrorCode(ctx, err)
		}

		logEntry := logger.
			WithContext(ctx).
			WithError(err)
		if responseErr != nil {
			logEntry = logEntry.WithField("response_error", responseErr)
		}
		if uid, uidErr := apiutils.GetContextUID(ctx); uidErr == nil {
			logEntry = logEntry.WithField(comlogging.FieldKeyUID, uid)
		}
		logEntry.Errorf("http error response | err=%s", err.Error())
	}

	if config.Debug {
		e.Debug = true
		e.Use(echoMiddleware.Logger(), echoMiddleware.Recover())
	} else {
		e.HideBanner = true
		e.HidePort = true
		e.Use(middleware.Recover, middleware.SentryPrepare)
	}

	InitValidator(e)

	return e
}

func StartServer(host string, port int, e *echo.Echo) {
	logger := comlogging.GetLogger()

	go func() {
		bindAddr := fmt.Sprintf("%s:%d", host, port)
		err := e.Start(bindAddr)

		logger.
			WithError(err).
			WithField("address", bindAddr).
			Info("Server stopped.")
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, config.SignalInterrupt)
	signal.Notify(quit, config.SignalTerminate)
	<-quit
	logger.Info("Server shutting down...")

	gracefulDuration := time.Duration(viper.GetInt(config.KeyServerGracefulTimeout)) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), gracefulDuration)
	defer cancel()

	err := e.Shutdown(ctx)
	logger.
		WithError(err).
		WithField("server_timeout", int(gracefulDuration.Seconds())).
		Info("Server exited.")
}

func BindAndValidate(c echo.Context, model interface{}) (err error) {
	if err = c.Bind(model); err != nil {
		return utils.WrapError(err)
	}
	if err = c.Validate(model); err != nil {
		return utils.WrapError(err)
	}
	return err
}
