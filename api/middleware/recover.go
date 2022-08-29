package middleware

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func Recover(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				recoverErr, ok := err.(error)
				if !ok {
					recoverErr = fmt.Errorf("%v", err)
				}

				c.Error(recoverErr)
			}
		}()

		return next(c)
	}
}
