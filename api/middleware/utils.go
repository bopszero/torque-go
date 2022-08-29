package middleware

import (
	"bytes"
	"io/ioutil"

	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/go-common/comutils"
)

const contextKeyRequestBody = "__request:body"

func readContextRequestBody(ctx echo.Context) string {
	contextBody := ctx.Get(contextKeyRequestBody)
	if contextBody != nil {
		return contextBody.(string)
	}

	request := ctx.Request()
	if request.Body == nil {
		return ""
	}

	requestBody, err := ioutil.ReadAll(request.Body)
	comutils.PanicOnError(err)
	request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))

	requestBodyStr := string(requestBody)
	ctx.Set(contextKeyRequestBody, requestBodyStr)

	return requestBodyStr
}

func truncateStringDisplay(s string, limit int) string {
	if limit > 0 && len(s) <= limit {
		return s
	}

	return s[:limit] + "..."
}
