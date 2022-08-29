package apiutils

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
)

type EchoWrappedContext interface {
	context.Context
	echo.Context
}

type echoWrappedContext struct {
	echo.Context
}

func (*echoWrappedContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (*echoWrappedContext) Done() <-chan struct{} {
	return nil
}

func (*echoWrappedContext) Err() error {
	return nil
}

func (this *echoWrappedContext) Value(key interface{}) interface{} {
	return this.Get(fmt.Sprintf("%v", key))
}

func EchoWrapContext(c echo.Context) EchoWrappedContext {
	return &echoWrappedContext{c}
}
