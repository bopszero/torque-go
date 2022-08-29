package api

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type OurValidator struct {
	validator *validator.Validate
}

func InitValidator(e *echo.Echo) {
	v := utils.NewValidator()
	v.RegisterTagNameFunc(getJsonTagName)

	e.Validator = OurValidator{validator: v}
}

func getJsonTagName(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

	if name == "-" {
		return ""
	}

	return name
}

func (v OurValidator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}
