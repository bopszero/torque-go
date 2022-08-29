package ordermod

import (
	"fmt"
	"time"

	"github.com/jinzhu/copier"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbfields"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/thirdpartymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func NewOrder(currency meta.Currency) models.Order {
	now := time.Now()

	return models.Order{
		Currency:  currency,
		Direction: constants.DirectionTypeUnknown,
		Status:    constants.OrderStatusUnknown,

		ExtraData: make(dbfields.JsonField),

		CreateTime: now.Unix(),
		UpdateTime: now.Unix(),
	}
}

func NewUserOrder(
	uid meta.UID, currency meta.Currency,
	srcChannelType meta.ChannelType, dstChannelType meta.ChannelType,
) models.Order {
	var direction meta.Direction
	if dstChannelType == constants.ChannelTypeDstBalance {
		direction = constants.OrderDirectionTopup
	} else {
		direction = constants.OrderDirectionPayment
	}

	order := NewOrder(currency)
	order.UID = uid
	order.Code = comutils.NewUuidCode()
	order.Direction = direction
	order.Status = constants.OrderStatusNew

	order.SrcChannelType = srcChannelType
	order.DstChannelType = dstChannelType

	order.AmountSubTotal = decimal.Zero
	order.AmountFee = decimal.Zero
	order.AmountDiscount = decimal.Zero
	order.AmountTotal = decimal.Zero

	return order
}

func ValidateOrder(order models.Order) error {
	channelPair := ChannelPair{
		SourceType:      order.SrcChannelType,
		DestinationType: order.DstChannelType,
	}
	if _, ok := ValidChannelPairSet[channelPair]; !ok {
		return constants.ErrorOrderInvalid
	}

	return nil
}

func ValidateUserRequestOrder(order models.Order) error {
	if _, ok := ValidUserSourceChannelSet[order.SrcChannelType]; !ok {
		return constants.ErrorOrderInvalid
	}

	return ValidateOrder(order)
}

func InitOrder(ctx comcontext.Context, orderID uint64) (models.Order, error) {
	return ProcessOrderSteps(
		ctx, orderID,
		[]meta.OrderStatus{
			constants.OrderStatusNew,
		},
		constants.OrderStepCodeOrderInit,
	)
}

func ExecuteOrder(ctx comcontext.Context, orderID uint64) (models.Order, error) {
	return ProcessOrderSteps(
		ctx, orderID,
		[]meta.OrderStatus{
			constants.OrderStatusNew,
			constants.OrderStatusInit,
			constants.OrderStatusHandleSrc,
			constants.OrderStatusHandleDst,
			constants.OrderStatusFailing,
			constants.OrderStatusRefunding,
			constants.OrderStatusCompleting,
			constants.OrderStatusCompleted,
		},
		constants.OrderStepCodeEmpty,
	)
}

func MarkOrderAsFailing(ctx comcontext.Context, orderID uint64) (order models.Order, err error) {
	err = database.TransactionFromContext(ctx, database.AliasWalletMaster, func(dbTxn *gorm.DB) error {
		lockOrderQueryDB := dbquery.SelectForUpdate(dbTxn).
			Where(dbquery.Gt(models.OrderColStatus, constants.OrderStatusInit)).
			Where(dbquery.NotIn(
				models.OrderColStatus,
				[]meta.OrderStatus{
					constants.OrderStatusNeedStaff,
					constants.OrderStatusRefunded,
					constants.OrderStatusFailing,
					constants.OrderStatusRefunding,
				},
			)).
			First(&order, &models.Order{ID: orderID})
		if err := lockOrderQueryDB.Error; err != nil {
			return err
		}

		order.Status = constants.OrderStatusFailing

		comlogging.GetLogger().
			WithContext(ctx).
			WithField("order_id", order.ID).
			Info("order mark fail request")

		if err := dbTxn.Save(&order).Error; err != nil {
			return utils.WrapError(err)
		}

		return nil
	})
	return
}

func ProcessOrderSteps(
	ctx comcontext.Context,
	orderID uint64, acceptedStatuses []meta.OrderStatus,
	lastStepCode string,
) (
	order models.Order, execErr error,
) {
	err := database.TransactionFromContext(ctx, database.AliasWalletMaster, func(dbTxn *gorm.DB) error {
		if err := dbquery.SelectForUpdate(dbTxn).First(&order, &models.Order{ID: orderID}).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = constants.ErrorOrderNotFound
			}
			return err
		}

		if err := ValidateOrder(order); err != nil {
			return err
		}

		isValidStatus := false
		for _, status := range acceptedStatuses {
			if order.Status == status {
				isValidStatus = true
			}
		}
		if !isValidStatus {
			return constants.ErrorOrderStatus
		}

		lastStep := GetOrderLastStep(order)
		var direction meta.Direction
		if constants.OrderFailingStatusSet.Contains(order.Status) {
			direction = constants.OrderStepDirectionBackward
		} else if lastStep != nil {
			direction = lastStep.Direction
		} else {
			direction = constants.OrderStepDirectionForward
		}

		logger := comlogging.GetLogger()
		var (
			result  OrderStepResult
			isSaved bool
		)
		for {
			step := GetOrderNextStep(order, direction)
			if step == nil {
				break
			}
			logEntry := logger.
				WithContext(ctx).
				WithFields(logrus.Fields{
					"order_id": order.ID,
					"step":     step.GetCode(),
				})

			err := database.TransactionFromContext(ctx, database.AliasWalletMaster, func(_ *gorm.DB) (err error) {
				defer func() {
					err = comutils.ToError(recover())
				}()

				isSaved = false
				switch direction {
				case constants.OrderStepDirectionForward:
					result = step.Forward(ctx, &order)
					logEntry = logEntry.WithField("direction", "forward")
				case constants.OrderStepDirectionBackward:
					result = step.Backward(ctx, &order)
					logEntry = logEntry.WithField("direction", "backward")
				default:
					panic(fmt.Errorf("unknown order step direciton `%v`", direction))
				}

				return
			})
			if err != nil {
				now := time.Now()
				order.RetryTime = now.Add(time.Minute).Unix()
				order.UpdateTime = now.Unix()
				comutils.PanicOnError(
					dbTxn.Save(&order).Error,
				)

				execErr = utils.WrapError(err)
				return nil
			}

			confirmOrderStep(&order, step, direction, result)
			if result.IsFail() && direction != constants.OrderStepDirectionBackward {
				order.Status = constants.OrderStatusFailing
				direction = constants.OrderStepDirectionBackward
			}

			logEntry = logEntry.WithField("step_result", result.GetCode())
			if result.GetError() != nil {
				logEntry.WithError(result.GetError()).Error("order execution error")
			} else {
				logEntry.Info("order execution")
			}

			now := time.Now()
			isNeedUpdate := true

			switch result.GetCode() {

			case OrderStepResultCodeIgnore:
				isNeedUpdate = false

			case OrderStepResultCodeNeedStaff:
				order.Status = constants.OrderStatusNeedStaff

			case OrderStepResultCodeFail, OrderStepResultCodeRetry:
				if order.RetryTime <= now.Unix() {
					var retryWait time.Duration
					if result.IsFail() {
						retryWait = 5 * time.Second
					} else {
						retryWait = time.Minute
					}
					order.RetryTime = now.Add(retryWait).Unix()
				}

			default:
				break
			}

			order.UpdateTime = now.Unix()
			if isNeedUpdate {
				comutils.PanicOnError(
					dbTxn.Save(&order).Error,
				)
				isSaved = true
			}

			execErr = result.GetError()
			if execErr != nil || step == nil || !result.IsSuccess() {
				break
			}
			if step.GetCode() == lastStepCode {
				break
			}
		}

		if !isSaved {
			comutils.PanicOnError(
				dbTxn.Save(&order).Error,
			)
		}

		return nil
	})
	if err != nil {
		return order, err
	} else {
		return order, execErr
	}
}

type OrderNotifData struct {
	Currency meta.Currency `json:"currency"`
	Ref      string        `json:"ref"`

	Direction      meta.Direction   `json:"direction_type"`
	SrcChannelType meta.ChannelType `json:"src_channel_type"`
	DstChannelType meta.ChannelType `json:"dst_channel_type"`
	AmountSubTotal decimal.Decimal  `json:"amount_sub_total"`
	AmountTotal    decimal.Decimal  `json:"amount_total"`
}

func PushOrderNotification(ctx comcontext.Context, order models.Order, notif *Notification) (err error) {
	var orderNotifData OrderNotifData
	if err = copier.Copy(&orderNotifData, &order); err != nil {
		return
	}
	if blockchainmod.IsNativeCurrency(order.Currency) {
		orderNotifData.Ref = order.DstChannelRef
	} else {
		orderNotifData.Ref = comutils.Stringify(order.ID)
	}

	client := thirdpartymod.GetPushServiceSystemClient()
	messageData := thirdpartymod.PushServiceMessageData{
		Title:   notif.Title,
		Message: notif.Message,

		Action:            thirdpartymod.ServicePushActionActivity,
		ActionDestination: thirdpartymod.ServicePushActionDestinationWalletTxn,
		ActionData:        comutils.JsonEncodeF(orderNotifData),
	}
	err = client.Push(ctx, order.UID, messageData)
	if err != nil {
		return
	}

	return nil
}

func IsUserActionLocked(ctx comcontext.Context, uid meta.UID, actionType uint16) bool {
	var userAction models.UserActionLock
	err := database.GetDbF(database.AliasWalletSlave).
		Last(&userAction, &models.UserActionLock{UID: uid, ActionType: actionType}).
		Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			comlogging.GetLogger().
				WithContext(ctx).
				WithError(err).
				WithFields(logrus.Fields{
					"uid":         uid,
					"action_type": actionType,
				}).
				Errorf("get user action lock failed | err=%v", err.Error())
		}
		return false
	}
	now := time.Now().Unix()
	return userAction.FromTime <= now && (now <= userAction.ToTime || userAction.ToTime == 0)
}
