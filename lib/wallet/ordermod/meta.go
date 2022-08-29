package ordermod

import (
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	ExtraDataSectionMeta    = "__meta__"
	ExtraDataMetaFailStatus = "fail_status"

	StepsDataErrorMessageMaxLength = 1500
)

type ChannelPair struct {
	SourceType      meta.ChannelType
	DestinationType meta.ChannelType
}

type ChannelTypeSet map[meta.ChannelType]bool

var (
	ValidDestinationChannelMap = map[meta.ChannelType]ChannelTypeSet{
		constants.ChannelTypeSrcBalance: {
			constants.ChannelTypeDstSystem:         true,
			constants.ChannelTypeDstTransfer:       true,
			constants.ChannelTypeDstTorqueConvert:  true,
			constants.ChannelTypeDstProfitReinvest: true,
			constants.ChannelTypeDstProfitWithdraw: true,

			constants.ChannelTypeDstMerchantGorillaHotel:  true,
			constants.ChannelTypeDstMerchantGorillaFlight: true,
			constants.ChannelTypeDstMerchantTorqueMall:    true,
		},
		constants.ChannelTypeSrcBlockchainNetwork: {
			constants.ChannelTypeDstBlockchainNetwork: true,
			constants.ChannelTypeDstTorqueConvert:     true,
		},
	}
	ValidSourceChannelMap = map[meta.ChannelType]ChannelTypeSet{
		constants.ChannelTypeDstBalance: {
			constants.ChannelTypeSrcSystem:          true,
			constants.ChannelTypeSrcTransfer:        true,
			constants.ChannelTypeSrcTorqueConvert:   true,
			constants.ChannelTypeSrcPromoCode:       true,
			constants.ChannelTypeSrcTradingReward:   true,
			constants.ChannelTypeSrcBonusPoolLeader: true,
		},
	}
	ValidChannelPairSet map[ChannelPair]bool

	ValidUserSourceChannelSet = ChannelTypeSet{
		constants.ChannelTypeSrcBalance:           true,
		constants.ChannelTypeSrcBlockchainNetwork: true,
		constants.ChannelTypeSrcPromoCode:         true,
	}
)

func init() {
	ValidChannelPairSet = make(map[ChannelPair]bool)
	for srcChannel, dstChannelSet := range ValidDestinationChannelMap {
		for dstChannel := range dstChannelSet {
			pair := ChannelPair{
				SourceType:      srcChannel,
				DestinationType: dstChannel,
			}
			ValidChannelPairSet[pair] = true
		}
	}
	for dstChannel, srcChannelSet := range ValidSourceChannelMap {
		for srcChannel := range srcChannelSet {
			pair := ChannelPair{
				SourceType:      srcChannel,
				DestinationType: dstChannel,
			}
			ValidChannelPairSet[pair] = true
		}
	}
}
