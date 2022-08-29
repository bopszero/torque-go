package ordermod

import (
	"fmt"

	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

func getStepByCode(code string) OrderStep {
	step, ok := StepMap[code]
	if !ok {
		panic(fmt.Errorf("unknown step code `%v`", code))
	}

	return step
}

func GetOrderLastStep(order models.Order) *models.OrderStep {
	if len(order.StepsData.History) == 0 {
		return nil
	}

	return &order.StepsData.History[len(order.StepsData.History)-1]
}

func GetOrderNextStep(order models.Order, direction meta.Direction) OrderStep {
	lastStep := GetOrderLastStep(order)
	if lastStep == nil {
		if direction == constants.OrderStepDirectionForward {
			return OrderedSteps[0]
		} else {
			return nil
		}
	}

	if lastStep.Direction != direction {
		if order.StepsData.Current != nil {
			return getStepByCode(order.StepsData.Current.Code)
		} else {
			return getStepByCode(lastStep.Code)
		}
	}

	switch direction {
	case constants.OrderStepDirectionForward:
		if lastStep.Code == OrderedSteps[len(OrderedSteps)-1].GetCode() {
			return nil
		}
		for i := 0; i < len(OrderedSteps)-1; i++ {
			step := OrderedSteps[i]

			if step.GetCode() == lastStep.Code {
				return OrderedSteps[i+1]
			}
		}

	case constants.OrderStepDirectionBackward:
		if lastStep.Code == OrderedSteps[0].GetCode() {
			return nil
		}
		for i := len(OrderedSteps) - 1; i > 0; i-- {
			step := OrderedSteps[i]

			if step.GetCode() == lastStep.Code {
				return OrderedSteps[i-1]
			}
		}

	default:
		panic(fmt.Errorf("unknown order step direction `%v`", lastStep.Direction))
	}

	return nil
}
