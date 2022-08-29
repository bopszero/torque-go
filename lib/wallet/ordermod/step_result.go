package ordermod

import (
	"fmt"

	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	OrderStepResultCodeFail    = meta.OrderStepResultCode(-1)
	OrderStepResultCodeIgnore  = meta.OrderStepResultCode(0)
	OrderStepResultCodeSuccess = meta.OrderStepResultCode(10)

	OrderStepResultCodeRetry     = meta.OrderStepResultCode(1)
	OrderStepResultCodeNeedStaff = meta.OrderStepResultCode(2)
)

type (
	OrderStepResult interface {
		String() string

		GetError() error
		GetCode() meta.OrderStepResultCode

		IsSuccess() bool
		IsFail() bool
		IsRetry() bool

		IsPause() bool
	}

	orderStepResult struct {
		Code  meta.OrderStepResultCode
		Error error
	}
)

func newOrderStepResult(code meta.OrderStepResultCode, err error) OrderStepResult {
	return &orderStepResult{
		Code:  code,
		Error: err,
	}
}

func (this *orderStepResult) String() string {
	if this.Error == nil {
		return fmt.Sprintf("%v", this.Code)
	} else {
		return fmt.Sprintf("%v(%v)", this.Code, this.Error.Error())
	}
}

func (this *orderStepResult) GetError() error {
	return this.Error
}

func (this *orderStepResult) GetCode() meta.OrderStepResultCode {
	return this.Code
}

func (this *orderStepResult) IsSuccess() bool {
	return this.Code == OrderStepResultCodeSuccess || this.Code == OrderStepResultCodeIgnore
}

func (this *orderStepResult) IsFail() bool {
	return this.Code == OrderStepResultCodeFail
}

func (this *orderStepResult) IsRetry() bool {
	return this.Code == OrderStepResultCodeRetry
}

func (this *orderStepResult) IsPause() bool {
	return this.Code == OrderStepResultCodeNeedStaff
}

var (
	OrderStepResultSuccess    = newOrderStepResult(OrderStepResultCodeSuccess, nil)
	OrderStepResultIgnore     = newOrderStepResult(OrderStepResultCodeIgnore, nil)
	OrderStepResultRetryEmpty = newOrderStepResult(OrderStepResultCodeRetry, nil)
)
