package test

import (
	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/go-common/comcache"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/database"
)

func init() {
	config.Init()
	comutils.PanicOnError(
		comlogging.Init(viper.GetString(config.KeyLogFile), config.Debug),
	)
	comutils.PanicOnError(
		comcache.Init(config.KeyCacheMap, &comcache.MsgPackEncoder{}),
	)

	database.Init()
}
