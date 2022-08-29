package ordermod

import (
	"time"

	"gitlab.com/snap-clickstaff/go-common/comcache"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type (
	ChannelMap      map[meta.ChannelType]Channel
	tChannelsLoader func() ([]Channel, error)
)

var (
	vChannelMapCached comcache.CacheObject
	vChannelsLoaders  []tChannelsLoader
)

func init() {
	var channelMapTimeout time.Duration
	if config.Test {
		channelMapTimeout = 30 * time.Second
	} else {
		channelMapTimeout = 5 * time.Minute
	}
	vChannelMapCached = comcache.NewCacheObject(
		channelMapTimeout,
		func() (value interface{}, err error) {
			return channelLoadMap()
		},
	)
}

func channelLoadMap() (_ ChannelMap, err error) {
	var infoModels []models.ChannelInfo
	if err = database.GetDbF(database.AliasWalletSlave).Find(&infoModels).Error; err != nil {
		err = utils.WrapError(err)
		return
	}
	infoModelMap := make(map[meta.ChannelType]models.ChannelInfo, len(infoModels))
	for _, model := range infoModels {
		infoModelMap[model.Type] = model
	}
	channelMap := make(ChannelMap, len(vChannelsLoaders))
	for _, loader := range vChannelsLoaders {
		channels, channelErr := loader()
		if channelErr != nil {
			err = channelErr
			return
		}
		for _, channel := range channels {
			channelType := channel.GetType()
			if info, ok := infoModelMap[channelType]; ok {
				channel.SetInfoModel(info)
			}
			if existsChannel, ok := channelMap[channelType]; ok && existsChannel != channel {
				err = utils.IssueErrorf(
					"channel loaders have duplicated type on different channel %v",
					channelType,
				)
				return
			}
			channelMap[channelType] = channel
		}
	}
	return channelMap, nil
}

func ChannelRegisterLoader(loader tChannelsLoader) error {
	vChannelsLoaders = append(vChannelsLoaders, loader)
	return nil
}

func ChannelRegister(channel Channel) error {
	return ChannelRegisterLoader(
		func() ([]Channel, error) {
			return []Channel{channel}, nil
		},
	)
}

func GetChannelMap() ChannelMap {
	channelMap, err := vChannelMapCached.Get()
	comutils.PanicOnError(err)

	return channelMap.(ChannelMap)
}
