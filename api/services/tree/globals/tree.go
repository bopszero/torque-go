package globals

import (
	"time"

	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/go-common/comcache"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/affiliate"
)

var globalTree comcache.CacheObject

func GetTree() *affiliate.Tree {
	if globalTree == nil {
		refreshIntervalDesc := viper.GetString(config.KeyTreeRefreshDuration)
		if refreshIntervalDesc == "" {
			refreshIntervalDesc = "5m"
		}
		refreshInterval, err := time.ParseDuration(refreshIntervalDesc)
		comutils.PanicOnError(err)

		globalTree = comcache.NewCacheObjectAsync(
			refreshInterval,
			func() (interface{}, error) {
				return affiliate.GenerateTreeAtNow(), nil
			},
		)
	}

	treeObj, err := globalTree.Get()
	comutils.PanicOnError(err)

	return treeObj.(*affiliate.Tree)
}
