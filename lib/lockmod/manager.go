package lockmod

import (
	"fmt"
	"time"

	"github.com/go-redsync/redsync/v4"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func GetDefaultRedsync() *OurRedSync {
	return defaultRedsync.Get().(*OurRedSync)
}

func Lock(
	ctx comcontext.Context,
	key string, timeout time.Duration, retryDelay time.Duration,
) (*redsync.Mutex, error) {
	var (
		rs  = GetDefaultRedsync()
		mux = rs.NewMutex(
			key,
			redsync.WithTries(int(timeout/retryDelay)),
			redsync.WithExpiry(timeout),
			redsync.WithRetryDelay(retryDelay),
		)
	)
	if err := mux.LockContext(ctx); err != nil {
		return nil, utils.WrapError(err)
	}
	return mux, nil
}

func LockSimple(key string, params ...interface{}) (*redsync.Mutex, error) {
	var (
		fullKey = fmt.Sprintf(key, params...)
		rs      = GetDefaultRedsync()
	)
	return Lock(
		nil,
		fullKey,
		rs.defaultLockTimeout,
		rs.defaultRetryDelay,
	)
}

func LockNoRetry(key string, timeout time.Duration) (*redsync.Mutex, error) {
	var (
		rs  = GetDefaultRedsync()
		mux = rs.NewMutex(
			key,
			redsync.WithTries(1),
			redsync.WithExpiry(timeout),
		)
	)
	if err := mux.Lock(); err != nil {
		return nil, utils.WrapError(err)
	}
	return mux, nil
}
