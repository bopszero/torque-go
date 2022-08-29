package ipmod

import (
	"github.com/oschwald/maxminddb-golang"
	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/go-common/comcache"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"time"
)

type MaxMindRecord struct {
	Country MaxMindCountry `maxminddb:"country"`
}

type MaxMindCountry struct {
	ISOCode string `maxminddb:"iso_code"`
}

var timeCache = 30 * time.Second
var whiteListIpCached = comcache.NewCacheObject(
	timeCache,
	func() (interface{}, error) {
		return getListIp(constants.SettingKeyWhiteListIP)
	},
)
var blackListIpCached = comcache.NewCacheObject(
	timeCache,
	func() (interface{}, error) {
		return getListIp(constants.SettingKeyBlackListIP)
	},
)
var maxMindDbReaderCached = comcache.NewCacheObject(
	timeCache,
	func() (interface{}, error) {
		geoIpDBPath := viper.GetString(config.KeyGeoIPDBPath)
		maxMindDbReader, err := maxminddb.Open(geoIpDBPath)
		if err != nil {
			return nil, err
		}
		return maxMindDbReader, nil
	},
)
