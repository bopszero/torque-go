package msgqueuemod

import (
	"github.com/adjust/rmq/v4"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func GetConnectionDefault() rmq.Connection {
	return rmqConnectionDefaultSingleton.Get().(rmq.Connection)
}

func GetQueueAuth() (rmq.Queue, error) {
	return getQueue(QueueKeyAuth)
}

func GetQueueWallet() (rmq.Queue, error) {
	return getQueue(QueueKeyWallet)
}

func PublishMessage(queue rmq.Queue, msg Message) error {
	msgJSON, err := comutils.JsonEncode(msg)
	if err != nil {
		return utils.WrapError(err)
	}
	if err := queue.Publish(msgJSON); err != nil {
		return utils.WrapError(err)
	}
	return nil
}
