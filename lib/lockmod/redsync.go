package lockmod

import (
	"fmt"
	"time"

	"github.com/go-redsync/redsync/v4"
)

type OurRedSync struct {
	*redsync.Redsync

	keyPrefix          string
	defaultLockTimeout time.Duration
	defaultRetryDelay  time.Duration
}

func (this *OurRedSync) NewMutex(name string, options ...redsync.Option) *redsync.Mutex {
	prefixName := fmt.Sprintf("%s:%s", this.keyPrefix, name)
	return this.Redsync.NewMutex(prefixName, options...)
}
