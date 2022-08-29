package msgqueuemod

import (
	"fmt"
	"net/url"

	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/go-common/comcache"
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
)

var (
	globalConfig = comtypes.NewSingleton(func() interface{} {
		queueURL, err := url.Parse(viper.GetString(config.KeyMsgQueueDSN))
		comutils.PanicOnError(err)
		if queueURL.Scheme != "redis" {
			comutils.PanicOnError(fmt.Errorf("rmq only accept Redis as broker (not %v)", queueURL.Scheme))
		}
		var (
			queueQuery = queueURL.Query()
			queryMap   = comutils.UrlQueryToMap(queueURL)
			connOpts   comcache.RedisQueryOptions
		)
		comutils.JsonDecodeF(comutils.JsonEncodeF(queryMap), &connOpts)

		return &Config{
			DSN:          queueURL,
			KeyPrefix:    queueQuery.Get("key_prefix"),
			RedisOptions: connOpts,
		}
	})
)

type Config struct {
	DSN          *url.URL
	KeyPrefix    string
	RedisOptions comcache.RedisQueryOptions
}

func GetConfig() *Config {
	return globalConfig.Get().(*Config)
}
