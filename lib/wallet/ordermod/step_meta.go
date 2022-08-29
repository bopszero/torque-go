package ordermod

import (
	"fmt"

	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
)

type (
	OrderStep interface {
		GetCode() string

		Forward(ctx comcontext.Context, order *models.Order) OrderStepResult
		Backward(ctx comcontext.Context, order *models.Order) OrderStepResult
	}
)

var (
	OrderedSteps = []OrderStep{
		&InitOrderStep{},

		&SrcStartOrderStep{},
		&SrcPrepareOrderStep{},
		&SrcExecuteOrderStep{},
		&SrcCommitOrderStep{},
		&SrcFinishOrderStep{},

		&DstStartOrderStep{},
		&DstPrepareOrderStep{},
		&DstExecuteOrderStep{},
		&DstCommitOrderStep{},
		&DstFinishOrderStep{},

		&DoneOrderStep{},
	}
	StepMap map[string]OrderStep
)

func init() {
	StepMap = map[string]OrderStep{}
	for _, step := range OrderedSteps {
		stepCode := step.GetCode()

		if _, exist := StepMap[stepCode]; exist {
			panic(fmt.Errorf("order steps has duplicated at step `%v`", stepCode))
		}

		StepMap[stepCode] = step
	}
}
