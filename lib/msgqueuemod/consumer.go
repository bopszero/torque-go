package msgqueuemod

import (
	"github.com/adjust/rmq/v4"
	"github.com/sirupsen/logrus"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comutils"
)

func DefaultConsume(delivery rmq.Delivery) {
	var message Message
	if err := comutils.JsonDecode(delivery.Payload(), &message); err != nil {
		comlogging.GetLogger().
			WithError(err).
			WithField("payload", delivery.Payload()).
			Errorf("message queue consumer cannot JSON parse message")
		handleDeliveryActionError(delivery, "reject", delivery.Reject())
		return
	}

	logEntry := comlogging.GetLogger().WithFields(logrus.Fields{
		"id":   message.ID,
		"type": message.Type,
		"data": message.Data,
	})

	handler, ok := MessageHandlerMap[message.Type]
	if !ok {
		logEntry.Errorf("message queue consumer doesn't support type `%v`", message.Type)
		handleDeliveryActionError(delivery, "reject", delivery.Reject())
		return
	}

	logEntry.Infof("message queue is consuming message")
	if err := handler(message); err != nil {
		logEntry.
			WithError(err).
			Errorf("message queue consume message failed")
		handleDeliveryActionError(delivery, "push", delivery.Ack())
		return
	}
	logEntry.Infof("message queue consumed message successfully")

	handleDeliveryActionError(delivery, "ack", delivery.Reject())
}
