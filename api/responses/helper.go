package responses

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gorm.io/gorm"
)

func errorToResponseInfo(ctx comcontext.Context, err error) ApiResponse {
	if err == nil {
		return ApiResponse{
			Code: constants.ErrorCodeSuccess,
		}
	}

	var (
		code          meta.ErrorCode
		message       string
		errorMessages []string
	)

	switch errT := err.(type) {
	case *comlogging.SentryError:
		err = errT.GetError()
	}

	switch err {
	case gorm.ErrRecordNotFound:
		code = constants.ErrorDataNotFound.Code()
	default:
		switch errT := err.(type) {
		case validator.ValidationErrors:
			code = constants.ErrorCodeInvalidParams
			message = constants.ErrorInvalidParams.Message(ctx)

			for _, fieldError := range errT {
				errorMsg := fmt.Sprintf(
					"Field validation for '%s' failed on the '%s' tag.",
					fieldError.Field(), fieldError.Tag(),
				)
				errorMessages = append(errorMessages, errorMsg)
			}
		case *meta.GeneralError:
			code = errT.Code()
			message = errT.Message(ctx)
		case *meta.MessageError:
			code = constants.ErrorUnknown.Code()
			message = errT.Error()
		default:
			code = constants.ErrorUnknown.Code()
			message = constants.ErrorUnknown.Message(ctx)
		}
	}

	if len(errorMessages) == 0 {
		errorMessages = append(errorMessages, err.Error())
	}

	response := ApiResponse{
		Code:    code,
		Message: message,
	}
	if config.Test {
		response.Errors = errorMessages
	}

	return response
}
