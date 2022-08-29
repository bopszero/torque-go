package ordermod

import (
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/msgqueuemod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type OrderNotificationAsyncParams struct {
	OrderID uint64 `json:"order_id"`
}

func init() {
	comutils.PanicOnError(
		msgqueuemod.RegisterHandler(
			msgqueuemod.MessageTypeOrderNotifCompleted,
			PushOrderCompletedNotificationHandler,
		),
	)
	comutils.PanicOnError(
		msgqueuemod.RegisterHandler(
			msgqueuemod.MessageTypeOrderNotifFailed,
			PushOrderFailedNotificationHandler,
		),
	)
}

func PushOrderCompletedNotificationAsync(order models.Order) error {
	msg := msgqueuemod.NewMessageJsonF(
		msgqueuemod.MessageTypeOrderNotifCompleted,
		OrderNotificationAsyncParams{OrderID: order.ID},
	)
	queue, err := msgqueuemod.GetQueueWallet()
	if err != nil {
		return err
	}
	return msgqueuemod.PublishMessage(queue, msg)
}

func PushOrderCompletedNotificationHandler(msg msgqueuemod.Message) (err error) {
	var params OrderNotificationAsyncParams
	if err = comutils.JsonDecode(msg.Data.(string), &params); err != nil {
		err = utils.WrapError(err)
		return
	}

	var (
		ctx   = comcontext.NewContext()
		db    = database.GetDbF(database.AliasWalletMaster)
		order models.Order
	)
	if err = db.First(&order, &models.Order{ID: params.OrderID}).Error; err != nil {
		err = utils.WrapError(err)
		return
	}

	channel, err := GetChannelByType(GetOrderMainChannelType(order))
	if err != nil {
		return
	}
	notif, err := channel.GetNotificationCompleted(ctx, order)
	if err != nil || notif == nil {
		return
	}

	return PushOrderNotification(ctx, order, notif)
}

func PushOrderFailedNotificationAsync(order models.Order) (err error) {
	msg := msgqueuemod.NewMessageJsonF(
		msgqueuemod.MessageTypeOrderNotifFailed,
		OrderNotificationAsyncParams{OrderID: order.ID},
	)
	queue, err := msgqueuemod.GetQueueWallet()
	if err != nil {
		return err
	}
	return msgqueuemod.PublishMessage(queue, msg)
}

func PushOrderFailedNotificationHandler(msg msgqueuemod.Message) (err error) {
	var params OrderNotificationAsyncParams
	if err := comutils.JsonDecode(msg.Data.(string), &params); err != nil {
		return utils.WrapError(err)
	}

	ctx := comcontext.NewContext()
	db := database.GetDbF(database.AliasWalletMaster)

	var order models.Order
	if err := db.First(&order, &models.Order{ID: params.OrderID}).Error; err != nil {
		return err
	}

	channel, err := GetChannelByType(GetOrderMainChannelType(order))
	if err != nil {
		return
	}
	notif, err := channel.GetNotificationFailed(ctx, order)
	if err != nil || notif == nil {
		return
	}

	return PushOrderNotification(ctx, order, notif)
}
