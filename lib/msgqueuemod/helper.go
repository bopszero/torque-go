package msgqueuemod

import (
	"github.com/adjust/rmq/v4"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
)

func handleDeliveryActionError(delivery rmq.Delivery, action string, err error) {
	if err == nil {
		return
	}
	comlogging.GetLogger().
		WithType(constants.LogTypeMessageQueue).
		WithError(err).
		WithField("payload", delivery.Payload()).
		Errorf("delivery %s failed", action)
}
