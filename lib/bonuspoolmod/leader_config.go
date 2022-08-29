package bonuspoolmod

import (
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

var leaderConfig = comtypes.NewSingleton(func() interface{} {
	var (
		configMap = viper.GetStringMap(config.KeyBonusPoolLeaderConfig)
		conf      LeaderConfig
	)
	comutils.PanicOnError(
		utils.DumpDataByJSON(&configMap, &conf),
	)
	return &conf
})

type LeaderConfig struct {
	DefaultTierRateMap map[string]decimal.Decimal `json:"default_tier_rate_map"`
}

func LeaderGetConfig() *LeaderConfig {
	return leaderConfig.Get().(*LeaderConfig)
}
