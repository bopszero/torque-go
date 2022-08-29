package ipmod

import (
	"net"
	"strings"

	"github.com/oschwald/maxminddb-golang"
	"github.com/sirupsen/logrus"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/settingmod"
)

func getListIp(settingKeyListIp string) (comtypes.HashSet, error) {
	ipHashSet := make(comtypes.HashSet)
	listIp, err := settingmod.GetSetting(settingKeyListIp)
	if err != nil {
		return nil, err
	}
	ipSlice := strings.Split(listIp.Value, ",")
	for _, ip := range ipSlice {
		ipHashSet.Add(strings.Trim(ip, " "))
	}
	return ipHashSet, nil
}

func GetWhiteListIp() (comtypes.HashSet, error) {
	whiteListIp, err := whiteListIpCached.Get()
	if err != nil {
		return nil, err
	}
	return whiteListIp.(comtypes.HashSet), nil
}

func GetBlackListIp() (comtypes.HashSet, error) {
	blackListIp, err := blackListIpCached.Get()
	if err != nil {
		return nil, err
	}
	return blackListIp.(comtypes.HashSet), nil
}

func IsEnableBanIp(ctx comcontext.Context) bool {
	isEnableBanIpStr, err := settingmod.GetSettingValueFast(constants.SettingKeyIsEnableBanIP)
	if err != nil {
		comlogging.GetLogger().
			WithContext(ctx).
			WithError(err).
			WithFields(logrus.Fields{
				"key": constants.SettingKeyIsEnableBanIP,
			}).
			Errorf("get setting '%s' error | err=%s", constants.SettingKeyIsEnableBanIP, err.Error())
		return false
	}
	return isEnableBanIpStr == constants.ToggleCodeOn
}

func GetIpCountry(ip string) (country string, err error) {
	maxMindDbReaderCached, err := maxMindDbReaderCached.Get()
	maxMindDbReader := maxMindDbReaderCached.(*maxminddb.Reader)
	if err != nil {
		return "", err
	}
	var maxMindRecord MaxMindRecord
	err = maxMindDbReader.Lookup(net.ParseIP(ip), &maxMindRecord)
	if err != nil {
		return "", err
	}
	return maxMindRecord.Country.ISOCode, nil
}
