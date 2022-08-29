package responses

import (
	"net/http"

	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

func setResponseCode(ctx apiutils.EchoWrappedContext, code meta.ErrorCode) {
	responseHeader := ctx.Response().Header()
	responseHeader.Set(config.HttpHeaderApiResponseCode, string(code))
	responseHeader.Set("Torque-Response-Code", string(code)) // TODO: Deprecated
}

func New(ctx apiutils.EchoWrappedContext, code meta.ErrorCode, data interface{}) error {
	setResponseCode(ctx, code)

	return ctx.JSON(
		http.StatusOK,
		ApiResponse{
			Code: code,
			Data: data,
		},
	)
}

func NewError(ctx apiutils.EchoWrappedContext, statusCode int, err error) error {
	if statusCode < http.StatusBadRequest {
		statusCode = http.StatusNotImplemented
	}
	responseInfo := errorToResponseInfo(ctx, err)

	return ctx.JSON(statusCode, responseInfo)
}

func Ok(ctx apiutils.EchoWrappedContext, data interface{}) error {
	return New(ctx, constants.ErrorCodeSuccess, data)
}

func OkEmpty(ctx apiutils.EchoWrappedContext) error {
	return New(ctx, constants.ErrorCodeSuccess, nil)
}

func AutoErrorCode(ctx apiutils.EchoWrappedContext, err error) error {
	responseInfo := errorToResponseInfo(ctx, err)
	setResponseCode(ctx, responseInfo.Code)

	return ctx.JSON(http.StatusOK, responseInfo)
}

func AutoErrorCodeData(ctx apiutils.EchoWrappedContext, err error, data interface{}) error {
	responseInfo := errorToResponseInfo(ctx, err)
	responseInfo.Data = data
	setResponseCode(ctx, responseInfo.Code)

	return ctx.JSON(http.StatusOK, responseInfo)
}
