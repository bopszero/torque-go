package authmod

import (
	"fmt"

	"gitlab.com/snap-clickstaff/go-common/comcache"
	"gitlab.com/snap-clickstaff/torque-go/config"
)

func passwordSetNonce(info PasswordNonce) error {
	var (
		cache    = comcache.GetRemoteCache()
		cacheKey = fmt.Sprintf("%s:%s", info.Username, info.ID)
		timeout  = PasswordNonceTimeout
	)
	if config.Test {
		timeout *= 100
	}
	return cache.Set(cacheKey, info.Value, timeout)
}

func passwordGetNonce(username, nonceID string) (value []byte, err error) {
	var (
		cache    = comcache.GetRemoteCache()
		cacheKey = fmt.Sprintf("%s:%s", username, nonceID)
	)
	err = cache.Get(cacheKey, &value)
	return
}

func passwordDeleteNonce(username, nonceID string) error {
	var (
		cache    = comcache.GetRemoteCache()
		cacheKey = fmt.Sprintf("%s:%s", username, nonceID)
	)
	return cache.Delete(cacheKey)
}
