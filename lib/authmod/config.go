package authmod

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
)

type AuthConfig struct {
	AccessSecret   []byte
	AccessTimeout  time.Duration
	RefreshSecret  []byte
	RefreshTimeout time.Duration

	// Deprecated
	LegacySecret []byte
}

var (
	AuthConfigProxy = comtypes.NewSingleton(func() interface{} {
		panicConfigMissing := func(key string) {
			panic(fmt.Errorf("jwt config `%s` is missing", key))
		}

		accessTimeoutDesc := viper.GetString(config.KeyAuthAccessTimeout)
		if accessTimeoutDesc == "" {
			panicConfigMissing(config.KeyAuthAccessTimeout)
		}
		accessTimeout, err := time.ParseDuration(accessTimeoutDesc)
		comutils.PanicOnError(err)

		refreshTimeoutDesc := viper.GetString(config.KeyAuthRefreshTimeout)
		if accessTimeoutDesc == "" {
			panicConfigMissing(config.KeyAuthRefreshTimeout)
		}
		refreshTimeout, err := time.ParseDuration(refreshTimeoutDesc)
		comutils.PanicOnError(err)

		accessSecret := viper.GetString(config.KeyAuthAccessSecret)
		if accessSecret == "" {
			panicConfigMissing(config.KeyAuthAccessSecret)
		}
		refreshSecret := viper.GetString(config.KeyAuthRefreshSecret)
		if refreshSecret == "" {
			panicConfigMissing(config.KeyAuthRefreshSecret)
		}
		legacySecret := viper.GetString(config.KeyAuthLegacySecret)
		if legacySecret == "" {
			panicConfigMissing(config.KeyAuthLegacySecret)
		}

		accessKey, err := comutils.Base64DecodeNoPadding(accessSecret)
		comutils.PanicOnError(err)
		refreshKey, err := comutils.Base64DecodeNoPadding(refreshSecret)
		comutils.PanicOnError(err)

		return AuthConfig{
			AccessSecret:   accessKey,
			AccessTimeout:  accessTimeout,
			RefreshSecret:  refreshKey,
			RefreshTimeout: refreshTimeout,

			// Deprecated
			LegacySecret: []byte(legacySecret),
		}
	})
)

func GetAuthConfig() AuthConfig {
	return AuthConfigProxy.Get().(AuthConfig)
}
