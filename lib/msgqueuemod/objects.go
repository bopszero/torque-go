package msgqueuemod

import (
	"fmt"
	"sync"

	"github.com/adjust/rmq/v4"
	"github.com/go-redis/redis/v8"
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

var (
	rmqConnectionDefaultSingleton = comtypes.NewSingleton(func() interface{} {
		var (
			conf      = GetConfig()
			queueOpts = conf.RedisOptions

			connTag     = fmt.Sprintf("%s-default", conf.KeyPrefix)
			redisClient = redis.NewClient(&redis.Options{
				Addr:     conf.DSN.Host,
				Username: queueOpts.Username,
				Password: queueOpts.Password,
				DB:       queueOpts.DB,

				DialTimeout:  queueOpts.DialTimeout.Duration,
				ReadTimeout:  queueOpts.ReadTimeout.Duration,
				WriteTimeout: queueOpts.WriteTimeout.Duration,

				PoolSize:     queueOpts.PoolSize,
				MinIdleConns: queueOpts.MinIdleConns,
				MaxConnAge:   queueOpts.MaxConnAge.Duration,
				PoolTimeout:  queueOpts.PoolTimeout.Duration,
				IdleTimeout:  queueOpts.IdleTimeout.Duration,
			})
		)
		conn, err := rmq.OpenConnectionWithRedisClient(connTag, redisClient, nil)
		comutils.PanicOnError(err)
		return conn
	})
	queueMap     = make(map[string]rmq.Queue)
	queueMapLock sync.Mutex
)

func getQueue(key string) (queue rmq.Queue, err error) {
	conf := GetConfig()
	key = fmt.Sprintf("%s-%s", conf.KeyPrefix, key)

	queue, ok := queueMap[key]
	if ok {
		return queue, nil
	}

	queueMapLock.Lock()
	defer queueMapLock.Unlock()

	queue, ok = queueMap[key]
	if !ok {
		queue, err = GetConnectionDefault().OpenQueue(key)
		if err != nil {
			return nil, utils.WrapError(err)
		}
		queueMap[key] = queue
	}
	return queue, nil
}
