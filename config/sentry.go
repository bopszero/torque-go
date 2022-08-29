package config

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/torque-go/buildmeta"
)

func InitSentry() {
	if Debug {
		return
	}

	dsn := viper.GetString(KeySentryDSN)
	if dsn == "" {
		return
	}

	releaseDesc := fmt.Sprintf(
		"%s@v%s-%s",
		AppName,
		buildmeta.Version, buildmeta.GitCommitID[:8],
	)
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		AttachStacktrace: true,
		Environment:      Env,
		MaxBreadcrumbs:   30,
		Release:          releaseDesc,
	})

	if err != nil {
		panic(err)
	}

	CmdRegisterRootDefer(CloseSentry)
}

func CloseSentry() {
	sentry.Flush(time.Second * 5)
}
