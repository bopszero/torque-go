package constants

import "gitlab.com/snap-clickstaff/torque-go/lib/meta"

var (
	WalletBalanceTypeMetaCrSystem              = meta.WalletBalanceTypeMeta{ID: 1001, Code: "cr_system"}
	WalletBalanceTypeMetaCrTransfer            = meta.WalletBalanceTypeMeta{ID: 1003, Code: "cr_transfer"}
	WalletBalanceTypeMetaCrRefund              = meta.WalletBalanceTypeMeta{ID: 1005, Code: "cr_refund"}
	WalletBalanceTypeMetaCrFeeReverse          = meta.WalletBalanceTypeMeta{ID: 1007, Code: "cr_fee_reverse"}
	WalletBalanceTypeMetaCrPromoCode           = meta.WalletBalanceTypeMeta{ID: 1009, Code: "cr_promo_code"}
	WalletBalanceTypeMetaCrConversion          = meta.WalletBalanceTypeMeta{ID: 1011, Code: "cr_conversion"}
	WalletBalanceTypeMetaCrBonusPool           = meta.WalletBalanceTypeMeta{ID: 1013, Code: "cr_bonus_pool"}
	WalletBalanceTypeMetaCrDailyProfit         = meta.WalletBalanceTypeMeta{ID: 2005, Code: "cr_daily_profit"}
	WalletBalanceTypeMetaCrAffiliateCommission = meta.WalletBalanceTypeMeta{ID: 2007, Code: "cr_affiliate_commission"}
	WalletBalanceTypeMetaCrLeaderCommission    = meta.WalletBalanceTypeMeta{ID: 2009, Code: "cr_leader_commission"}

	WalletBalanceTypeMetaDrSystem         = meta.WalletBalanceTypeMeta{ID: 1002, Code: "dr_system"}
	WalletBalanceTypeMetaDrTransfer       = meta.WalletBalanceTypeMeta{ID: 1004, Code: "dr_transfer"}
	WalletBalanceTypeMetaDrReverse        = meta.WalletBalanceTypeMeta{ID: 1006, Code: "dr_reverse"} // TODO: Recheck usage
	WalletBalanceTypeMetaDrFee            = meta.WalletBalanceTypeMeta{ID: 1008, Code: "dr_fee"}
	WalletBalanceTypeMetaDrConversion     = meta.WalletBalanceTypeMeta{ID: 1012, Code: "dr_conversion"}
	WalletBalanceTypeMetaDrProfitReinvest = meta.WalletBalanceTypeMeta{ID: 2002, Code: "dr_profit_reinvest"}
	WalletBalanceTypeMetaDrProfitWithdraw = meta.WalletBalanceTypeMeta{ID: 2004, Code: "dr_profit_withdraw"}
	WalletBalanceTypeMetaDrMerchant       = meta.WalletBalanceTypeMeta{ID: 3002, Code: "dr_merchant"}
)
