package ordermod

import (
	"fmt"
	"reflect"

	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func GetChannelByType(channelType meta.ChannelType) (Channel, error) {
	channel, ok := GetChannelMap()[channelType]
	if !ok {
		return nil, utils.IssueErrorf("there is no channel of type `%v`", channelType)
	}

	return channel, nil
}

func genChannelMetaKey(channelType meta.ChannelType) string {
	return fmt.Sprintf("channel:%v", channelType)
}

func SetOrderChannelMetaData(order *models.Order, channelType meta.ChannelType, data interface{}) error {
	channel, err := GetChannelByType(channelType)
	if err != nil {
		return err
	}
	metaType := channel.GetMetaType()
	if metaType == nil {
		return utils.IssueErrorf("channel `%v` doesn't accept any meta", channelType)
	}

	var channelMeta interface{}
	switch reflect.Indirect(reflect.ValueOf(data)).Kind() {
	case reflect.Map:
		channelMeta = comutils.NewByType(metaType)
		if err = utils.DumpDataByJSON(data, channelMeta); err != nil {
			return err
		}
	case reflect.Struct:
		channelMeta = data
	default:
		panic(fmt.Errorf("order channel meta data type `%v` is invalid", reflect.TypeOf(data).Kind()))
	}
	if err = utils.ValidateStruct(channelMeta); err != nil {
		return err
	}

	contextData, exist := order.ExtraData[ExtraDataSectionMeta]
	if !exist {
		contextData = make(map[string]interface{})
	}

	key := genChannelMetaKey(channelType)
	contextData.(map[string]interface{})[key] = channelMeta

	order.ExtraData[ExtraDataSectionMeta] = contextData
	return nil
}

func GetOrderChannelMetaData(order *models.Order, channelType meta.ChannelType, value interface{}) error {
	key := genChannelMetaKey(channelType)
	data, ok := getOrderExtraData(order, ExtraDataSectionMeta, key)
	if !ok || data == nil {
		return utils.WrapError(constants.ErrorDataNotFound)
	}
	orderDataValue := reflect.Indirect(reflect.ValueOf(data))
	switch orderDataValue.Kind() {
	case reflect.Map:
		if err := utils.DumpDataByJSON(data, value); err != nil {
			return err
		}
	case reflect.Struct:
		if err := utils.DumpData(data, value); err != nil {
			return err
		}
	default:
		panic(fmt.Errorf("order channel meta data type `%v` is invalid", reflect.TypeOf(data).Kind()))
	}
	var (
		receiverValue = reflect.Indirect(reflect.ValueOf(value))
		needSetBack   = orderDataValue.Kind() != reflect.Struct &&
			receiverValue.Kind() == reflect.Struct
	)
	if needSetBack {
		setOrderExtraData(order, ExtraDataSectionMeta, key, value)
	}
	return nil
}

func GetOrderSrcChannelMetaData(order *models.Order, value interface{}) error {
	return GetOrderChannelMetaData(order, order.SrcChannelType, value)
}

func GetOrderDstChannelMetaData(order *models.Order, value interface{}) error {
	return GetOrderChannelMetaData(order, order.DstChannelType, value)
}

func GetBalanceTypeMetaDrByChannel(channelType meta.ChannelType) (*meta.WalletBalanceTypeMeta, error) {
	balanceMeta, ok := constants.ChannelDstToBalanceTypeMetaDrMap[channelType]
	if !ok {
		err := utils.IssueErrorf("channel `%v` doesn't have a debit balance type", channelType)
		return nil, err
	}

	return &balanceMeta, nil
}

func GetBalanceTypeMetaCrByChannel(channelType meta.ChannelType) (*meta.WalletBalanceTypeMeta, error) {
	balanceMeta, ok := constants.ChannelSrcToBalanceTypeMetaCrMap[channelType]
	if !ok {
		err := utils.IssueErrorf("channel `%v` doesn't have a credit balance type", channelType)
		return nil, err
	}

	return &balanceMeta, nil
}

func GetOrderDstChannelDetails(ctx comcontext.Context, order *models.Order) (interface{}, error) {
	channel, err := GetChannelByType(order.DstChannelType)
	if err != nil {
		return nil, err
	}

	return channel.GetOrderDetails(ctx, order)
}

func GetOrderMainChannelType(order models.Order) meta.ChannelType {
	if order.Direction == constants.OrderDirectionPayment {
		return order.DstChannelType
	} else {
		return order.SrcChannelType
	}
}

func GetOrderSourceChannel(order models.Order) (Channel, error) {
	return GetChannelByType(order.SrcChannelType)
}

func GetOrderDestinationChannel(order models.Order) (Channel, error) {
	return GetChannelByType(order.DstChannelType)
}
