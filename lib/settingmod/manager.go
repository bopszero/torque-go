package settingmod

import (
	"fmt"
	"time"

	"gitlab.com/snap-clickstaff/go-common/comcache"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func GetSetting(key string, args ...interface{}) (setting models.Setting, err error) {
	fullKey := fmt.Sprintf(key, args...)

	setting.Key = fullKey
	err = database.GetDbSlave().
		First(&setting, models.Setting{Key: fullKey}).
		Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			err = utils.WrapError(err)
			return
		}
		comlogging.GetLogger().
			WithField("key", fullKey).
			Warnf("setting not found | key=%s", fullKey)
	}

	return setting, nil
}

func GetSettingFast(key string, args ...interface{}) (setting models.Setting, err error) {
	fullKey := fmt.Sprintf(key, args...)
	cacheKey := fmt.Sprintf("setting:%v", fullKey)
	err = comcache.GetOrCreate(
		comcache.GetDefaultCache(),
		cacheKey,
		utils.GetEnvCacheDuration(5*time.Minute),
		&setting,
		func() (interface{}, error) {
			return GetSetting(fullKey)
		},
	)
	return
}

func GetSettingValueFast(key string, args ...interface{}) (string, error) {
	setting, err := GetSettingFast(key, args...)
	if err != nil {
		return "", err
	}

	return setting.Value, nil
}
