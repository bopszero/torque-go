package trading

import (
	"os"
	"os/signal"

	"github.com/robfig/cron/v3"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
)

func StartCron() {
	cronSkip := cron.New(
		cron.WithLogger(cron.DefaultLogger),
		cron.WithChain(
			cron.SkipIfStillRunning(cron.DefaultLogger),
		),
	)
	cronFree := cron.New(
		cron.WithLogger(cron.DefaultLogger),
	)

	// TODO: Add jobs here
	// cronSkip.AddFunc("@every 1m", crontasks.DoSomething)

	cronSkip.Start()
	cronFree.Start()
	comutils.EchoWithTime("Started the Cron Service.")

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, config.SignalInterrupt)

	osSignal := <-signalChannel
	comutils.EchoWithTime("Received signal `%s`, stopping the service.", osSignal)
	<-cronSkip.Stop().Done()
	<-cronFree.Stop().Done()
	comutils.EchoWithTime("Stopped the Cron Service.")
}
