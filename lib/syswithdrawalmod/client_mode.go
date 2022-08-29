package syswithdrawalmod

import (
	"time"

	"gitlab.com/snap-clickstaff/go-common/comcache"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type tClientModeMap map[string]string

var clientModeMapSettingProxy = comcache.NewCacheObject(
	15*time.Second,
	func() (interface{}, error) {
		var setting models.Setting
		err := database.GetDbSlave().
			First(
				&setting,
				&models.Setting{Key: constants.SettingKeySystemWithdrawalClientModeMap},
			).
			Error
		if database.IsDbError(err) {
			return nil, utils.WrapError(err)
		}
		var clientModeMap tClientModeMap
		if setting.Value != "" {
			if err := comutils.JsonDecode(setting.Value, &clientModeMap); err != nil {
				return nil, utils.WrapError(err)
			}
		}
		return clientModeMap, nil
	},
)
