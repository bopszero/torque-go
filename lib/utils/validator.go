package utils

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comtypes"
)

var ourValidator = comtypes.NewSingleton(func() interface{} { return NewValidator() })

func NewValidator() *validator.Validate {
	v := validator.New()
	v.RegisterCustomTypeFunc(validateDecimal, decimal.Zero)
	return v
}

func validateDecimal(field reflect.Value) interface{} {
	value := field.Interface().(decimal.Decimal)
	if value.Equal(value.Truncate(0)) {
		return value.IntPart()
	}
	valueFloat, _ := value.Float64()
	return valueFloat
}

func ValidateStruct(value interface{}) error {
	return ourValidator.Get().(*validator.Validate).Struct(value)
}
