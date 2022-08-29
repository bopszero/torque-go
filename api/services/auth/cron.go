package auth

import (
	"os"
	"os/signal"

	"gitlab.com/snap-clickstaff/torque-go/api/services/auth/crontasks"
	"gitlab.com/snap-clickstaff/torque-go/config"

	"github.com/robfig/cron/v3"
	"gitlab.com/snap-clickstaff/go-common/comutils"
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

	cronSkip.AddFunc("@every 30s", crontasks.KycUpdateJumioScans)
	cronSkip.AddFunc("@every 30s", crontasks.KycExecuteRequests)

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
