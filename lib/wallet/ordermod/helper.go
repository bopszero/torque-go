package ordermod

import (
	"time"

	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func confirmOrderStep(order *models.Order, step OrderStep, direction meta.Direction, result OrderStepResult) {
	currentStep := models.OrderStep{
		Direction: direction,
		Code:      step.GetCode(),
		Time:      time.Now().Unix(),
	}
	if result.GetError() != nil {
		errMsg := result.GetError().Error()
		if len(errMsg) > StepsDataErrorMessageMaxLength {
			errMsg = errMsg[:StepsDataErrorMessageMaxLength] + "..."
		}
		currentStep.Error = errMsg
	}

	order.StepsData.Current = &currentStep

	if result.IsSuccess() || result.IsFail() {
		order.StepsData.History = append(order.StepsData.History, currentStep)
		order.StepsData.Current = nil
	}
}

func getOrderExtraData(order *models.Order, section string, key string) (interface{}, bool) {
	sectionData, ok := order.ExtraData[section]
	if !ok {
		return nil, false
	}
	value, ok := sectionData.(map[string]interface{})[key]
	if !ok {
		return nil, false
	}
	return value, true
}

func setOrderExtraData(order *models.Order, section string, key string, value interface{}) {
	sectionData, ok := order.ExtraData[section]
	if !ok {
		sectionData = make(map[string]interface{})
		order.ExtraData[section] = sectionData
	}
	sectionData.(map[string]interface{})[key] = value
}

func SetOrderFailStatus(order *models.Order, status meta.OrderStatus) error {
	if status > constants.OrderStatusUnknown {
		return utils.IssueErrorf("invalid order fail status | order_id=%v,status=%v", order.ID, status)
	}

	var exMeta map[string]interface{}
	if data, ok := order.ExtraData[ExtraDataSectionMeta]; ok {
		exMeta = data.(map[string]interface{})
	} else {
		exMeta = make(map[string]interface{})
	}
	exMeta[ExtraDataMetaFailStatus] = status

	order.ExtraData[ExtraDataSectionMeta] = exMeta
	return nil
}

func getOrderFailStatus(order models.Order) meta.OrderStatus {
	exMetaVal, ok := order.ExtraData[ExtraDataSectionMeta]
	if !ok {
		return constants.OrderStatusUnknown
	}
	failStatusVal, ok := exMetaVal.(map[string]interface{})[ExtraDataMetaFailStatus]
	if !ok {
		return constants.OrderStatusUnknown
	}

	return meta.OrderStatus(
		comutils.ParseIntF(comutils.Stringify(failStatusVal)),
	)
}
