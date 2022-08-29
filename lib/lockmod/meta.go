package lockmod

import (
	"fmt"
	"net/url"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
)

var (
	// Supported query params
	// - key_prefix: Key prefix
	// - lock_timeout: Default lock timeout in ms
	// - retry_delay: Default lock retry delay in ms
	defaultRedsync = comtypes.NewSingleton(func() interface{} {
		var err error
		lockURL, err := url.Parse(viper.GetString(config.KeyRedLockDSN))
		comutils.PanicOnError(err)

		if lockURL.Scheme != "redis" {
			comutils.PanicOnError(fmt.Errorf("redlock only accept Redis (not %v)", lockURL.Scheme))
		}

		var (
			dsnQuery    = lockURL.Query()
			keyPrefix   string
			lockTimeout = 8 * time.Second
			retryDelay  = 500 * time.Millisecond
		)
		keyPrefix = dsnQuery.Get("key_prefix")
		if lockTimeoutStr := dsnQuery.Get("lock_timeout"); lockTimeoutStr != "" {
			lockTimeout, err = time.ParseDuration(lockTimeoutStr)
			comutils.PanicOnError(err)
		}
		if retryDelayStr := dsnQuery.Get("retry_delay"); retryDelayStr != "" {
			retryDelay, err = time.ParseDuration(retryDelayStr)
			comutils.PanicOnError(err)
		}

		systemPool := createRedisPool(lockURL)
		return &OurRedSync{
			Redsync: redsync.New(systemPool),

			keyPrefix:          keyPrefix,
			defaultLockTimeout: lockTimeout,
			defaultRetryDelay:  retryDelay,
		}
	})
)
