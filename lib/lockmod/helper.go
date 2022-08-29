package lockmod

import (
	"net/url"

	redislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"gitlab.com/snap-clickstaff/go-common/comcache"
	"gitlab.com/snap-clickstaff/go-common/comutils"
)

func createRedisPool(lockURL *url.URL) redis.Pool {
	var (
		queryMap = comutils.UrlQueryToMap(lockURL)
		connOpts comcache.RedisQueryOptions
	)
	comutils.JsonDecodeF(comutils.JsonEncodeF(queryMap), &connOpts)

	return goredis.NewPool(
		redislib.NewClient(&redislib.Options{
			Addr:     lockURL.Host,
			Username: connOpts.Username,
			Password: connOpts.Password,
			DB:       connOpts.DB,

			DialTimeout:  connOpts.DialTimeout.Duration,
			ReadTimeout:  connOpts.ReadTimeout.Duration,
			WriteTimeout: connOpts.WriteTimeout.Duration,

			PoolSize:     connOpts.PoolSize,
			MinIdleConns: connOpts.MinIdleConns,
			MaxConnAge:   connOpts.MaxConnAge.Duration,
			PoolTimeout:  connOpts.PoolTimeout.Duration,
			IdleTimeout:  connOpts.IdleTimeout.Duration,
		}),
	)

}
