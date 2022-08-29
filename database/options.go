package database

import (
	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type DbPoolOptions struct {
	PoolSize    int                   `json:"pool_size"`
	MaxLifetime comtypes.TimeDuration `json:"max_lifetime"`
	MaxSize     int                   `json:"max_size"`
}

func GetDbPoolOptions() (opts DbPoolOptions, err error) {
	opts = DbPoolOptions{
		PoolSize:    DefaultPoolSize,
		MaxLifetime: comtypes.TimeDuration{DefaultPoolMaxLifetime},
		MaxSize:     DefaultPoolMaxSize,
	}
	rawOptions := viper.GetStringMap(config.KeyDbPoolOptions)
	err = utils.DumpDataByJSON(rawOptions, &opts)
	return
}
