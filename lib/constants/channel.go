package constants

import "gitlab.com/snap-clickstaff/torque-go/lib/meta"

const (
	ChannelTypeSrcBalance           = 10000
	ChannelTypeSrcSystem            = 10001
	ChannelTypeSrcTransfer          = 10002
	ChannelTypeSrcBlockchainNetwork = 10003
	ChannelTypeSrcTorqueConvert     = 10004
	ChannelTypeSrcPromoCode         = 10100
	ChannelTypeSrcTradingReward     = 10101
	ChannelTypeSrcBonusPoolLeader   = 10102
)

const (
	ChannelTypeDstBalance           = 100000
	ChannelTypeDstSystem            = 100001
	ChannelTypeDstTransfer          = 100002
	ChannelTypeDstBlockchainNetwork = 100003
	ChannelTypeDstTorqueConvert     = 100004
	ChannelTypeDstProfitReinvest    = 100100
	ChannelTypeDstProfitWithdraw    = 100101

	ChannelTypeDstMerchantGorillaHotel  = 110000
	ChannelTypeDstMerchantGorillaFlight = 110001
	ChannelTypeDstMerchantTorqueMall    = 110002
)

var (
	ChannelDstToBalanceTypeMetaDrMap = map[meta.ChannelType]meta.WalletBalanceTypeMeta{
		ChannelTypeDstSystem:         WalletBalanceTypeMetaDrSystem,
		ChannelTypeDstTransfer:       WalletBalanceTypeMetaDrTransfer,
		ChannelTypeDstTorqueConvert:  WalletBalanceTypeMetaDrConversion,
		ChannelTypeDstProfitReinvest: WalletBalanceTypeMetaDrProfitReinvest,
		ChannelTypeDstProfitWithdraw: WalletBalanceTypeMetaDrProfitWithdraw,

		ChannelTypeDstMerchantGorillaHotel:  WalletBalanceTypeMetaDrMerchant,
		ChannelTypeDstMerchantGorillaFlight: WalletBalanceTypeMetaDrMerchant,
		ChannelTypeDstMerchantTorqueMall:    WalletBalanceTypeMetaDrMerchant,
	}

	ChannelSrcToBalanceTypeMetaCrMap = map[meta.ChannelType]meta.WalletBalanceTypeMeta{
		ChannelTypeSrcSystem:          WalletBalanceTypeMetaCrSystem,
		ChannelTypeSrcTransfer:        WalletBalanceTypeMetaCrTransfer,
		ChannelTypeSrcTorqueConvert:   WalletBalanceTypeMetaCrConversion,
		ChannelTypeSrcPromoCode:       WalletBalanceTypeMetaCrPromoCode,
		ChannelTypeSrcBonusPoolLeader: WalletBalanceTypeMetaCrBonusPool,
	}
)
