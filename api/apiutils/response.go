package apiutils

import (
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

func EchoSetDownloadResponse(c echo.Context, name string) {
	c.Response().Header().Set(
		echo.HeaderContentDisposition,
		fmt.Sprintf("attachment; filename=%q", name),
	)
}

func EchoResponseDownloadBytes(c echo.Context, name, contentType string, data []byte) error {
	EchoSetDownloadResponse(c, name)
	return c.Blob(http.StatusOK, contentType, data)
}

func EchoResponseDownloadSteam(c echo.Context, name, contentType string, reader io.Reader) error {
	EchoSetDownloadResponse(c, name)
	return c.Stream(http.StatusOK, contentType, reader)
}
