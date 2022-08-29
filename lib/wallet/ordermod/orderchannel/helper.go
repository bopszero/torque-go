package orderchannel

import (
	"fmt"
	"time"

	"gitlab.com/snap-clickstaff/go-common/comcache"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func GetChannelInfoMap() (channelMap map[meta.ChannelType]models.ChannelInfo, err error) {
	var channelInfoList []models.ChannelInfo
	err = database.
		GetDbF(database.AliasWalletSlave).
		Find(&channelInfoList).
		Error
	if err != nil {
		err = utils.WrapError(err)
		return
	}

	channelMap = make(map[meta.ChannelType]models.ChannelInfo, len(channelInfoList))
	for _, channel := range channelInfoList {
		channelMap[channel.Type] = channel
	}

	return
}

func GetChannelInfoMapFast() (channelMap map[meta.ChannelType]models.ChannelInfo, err error) {
	err = comcache.GetOrCreate(
		comcache.GetDefaultCache(),
		"channel_info:all_map",
		5*time.Minute,
		&channelMap,
		func() (interface{}, error) {
			return GetChannelInfoMap()
		},
	)
	return
}

func GetChannelInfo(channelType meta.ChannelType) (channel models.ChannelInfo, err error) {
	err = database.
		GetDbF(database.AliasWalletSlave).
		First(&channel, &models.ChannelInfo{Type: channelType}).
		Error
	return
}

func GetChannelInfoFast(channelType meta.ChannelType) (channel models.ChannelInfo, err error) {
	cacheKey := fmt.Sprintf("channel_info:%v", channelType)
	err = comcache.GetOrCreate(
		comcache.GetDefaultCache(),
		cacheKey,
		5*time.Minute,
		&channel,
		func() (interface{}, error) {
			return GetChannelInfo(channelType)
		},
	)
	return
}
