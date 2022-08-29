package wallet

import (
	"os"
	"os/signal"
	"time"

	"github.com/adjust/rmq/v4"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/msgqueuemod"
)

func StartConsumers() {
	var (
		logger      = comlogging.GetLogger()
		queues      = []rmq.Queue{}
		defaultConn = msgqueuemod.GetConnectionDefault()
	)
	if _, err := rmq.NewCleaner(defaultConn).Clean(); err != nil {
		key := "default"
		comutils.EchoWithTime("Clean connection failed | key=%s,err=%s", key, err.Error())
		logger.
			WithError(err).
			WithField("key", key).
			Warn("consumers clean connection failed")
	}

	targetQueue, err := msgqueuemod.GetQueueWallet()
	comutils.PanicOnError(err)
	comutils.PanicOnError(
		targetQueue.StartConsuming(500, 2*time.Second),
	)
	_, err = targetQueue.AddConsumerFunc("default", msgqueuemod.DefaultConsume)
	comutils.PanicOnError(err)
	queues = append(queues, targetQueue)

	comutils.EchoWithTime("Started the Consumers.")
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, config.SignalInterrupt)

	osSignal := <-signalChannel
	comutils.EchoWithTime("Received signal `%s`, stopping the service...", osSignal)
	for _, queue := range queues {
		<-queue.StopConsuming()
	}
	comutils.EchoWithTime("Stopped the Consumers.")
}
