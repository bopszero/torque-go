package config

import (
	"fmt"

	"gitlab.com/snap-clickstaff/go-common/comlogging"
)

var cmdRootDeferFuncs []func()

func CmdRegisterRootDefer(deferFunc func()) {
	cmdRootDeferFuncs = append(cmdRootDeferFuncs, deferFunc)
}

func CmdExecuteRootDefers() {
	for _, deferFunc := range cmdRootDeferFuncs {
		defer cmdRecoverRootDefer()
		deferFunc()
	}
}

func cmdRecoverRootDefer() {
	errObj := recover()
	if errObj == nil {
		return
	}
	err, ok := errObj.(error)
	if !ok {
		err = fmt.Errorf("%v", err)
	}
	comlogging.GetLogger().
		WithError(err).
		Errorf("execute root defer failed | err=%s", err.Error())
}
