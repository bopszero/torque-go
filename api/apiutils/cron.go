package apiutils

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

func CronRunWithRecovery(funcName string, data meta.O) {
	errLike := recover()
	if errLike == nil {
		return
	}
	fmt.Printf("cron task panic | exc=%v,data=%v\n", errLike, data)
	logEntry := comlogging.GetLogger().
		WithError(comutils.ToError(errLike)).
		WithField("func", funcName)
	if len(data) > 0 {
		logEntry = logEntry.WithFields(logrus.Fields(data))
	}
	logEntry.Error("cron task panic")
}
