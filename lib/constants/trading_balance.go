package constants

type TradingBalanceTypeMeta struct {
	ID   uint32
	Code string
}

var (
	TradingBalanceTypeMetaCrSystem                    = TradingBalanceTypeMeta{ID: 1001, Code: "cr_system"}
	TradingBalanceTypeMetaCrTransfer                  = TradingBalanceTypeMeta{ID: 1003, Code: "cr_transfer"}
	TradingBalanceTypeMetaCrDeposit                   = TradingBalanceTypeMeta{ID: 1005, Code: "cr_deposit"}
	TradingBalanceTypeMetaCrWithdrawInvestmentReverse = TradingBalanceTypeMeta{ID: 1007, Code: "cr_withdraw_investment_reverse"}
	TradingBalanceTypeMetaCrDailyProfit               = TradingBalanceTypeMeta{ID: 1009, Code: "cr_daily_profit"}
	TradingBalanceTypeMetaCrAffiliateCommission       = TradingBalanceTypeMeta{ID: 1011, Code: "cr_affiliate_commission"}
	TradingBalanceTypeMetaCrLeaderCommission          = TradingBalanceTypeMeta{ID: 1013, Code: "cr_leader_commission"}
	TradingBalanceTypeMetaCrReinvestSrcReverse        = TradingBalanceTypeMeta{ID: 1015, Code: "cr_reinvest_src_reverse"}
	TradingBalanceTypeMetaCrReinvestDst               = TradingBalanceTypeMeta{ID: 1017, Code: "cr_reinvest_dst"}
	TradingBalanceTypeMetaCrWithdrawProfitReverse     = TradingBalanceTypeMeta{ID: 1019, Code: "cr_withdraw_profit_reverse"}
	TradingBalanceTypeMetaCrPromoCode                 = TradingBalanceTypeMeta{ID: 1021, Code: "cr_promo_code"}
	TradingBalanceTypeMetaCrProductTravelTicketRefund = TradingBalanceTypeMeta{ID: 10003, Code: "cr_product_travel_ticket_refund"}
	TradingBalanceTypeMetaCrProductMallOrderRefund    = TradingBalanceTypeMeta{ID: 10005, Code: "cr_product_mall_order_refund"}
	TradingBalanceTypeMetaCrProductFlightTicketRefund = TradingBalanceTypeMeta{ID: 10007, Code: "cr_product_flight_ticket_refund"}

	TradingBalanceTypeMetaDrSystem                     = TradingBalanceTypeMeta{ID: 1002, Code: "dr_system"}
	TradingBalanceTypeMetaDrTransfer                   = TradingBalanceTypeMeta{ID: 1004, Code: "dr_transfer"}
	TradingBalanceTypeMetaDrDepositReverse             = TradingBalanceTypeMeta{ID: 1006, Code: "dr_deposit_reverse"}
	TradingBalanceTypeMetaDrWithdrawInvestment         = TradingBalanceTypeMeta{ID: 1008, Code: "dr_withdraw_investment"}
	TradingBalanceTypeMetaDrDailyProfitReverse         = TradingBalanceTypeMeta{ID: 1010, Code: "dr_daily_profit_reverse"}
	TradingBalanceTypeMetaDrAffiliateCommissionReverse = TradingBalanceTypeMeta{ID: 1012, Code: "dr_affiliate_commission_reverse"}
	TradingBalanceTypeMetaDrLeaderCommissionReverse    = TradingBalanceTypeMeta{ID: 1014, Code: "dr_leader_commission_reverse"}
	TradingBalanceTypeMetaDrReinvestSrc                = TradingBalanceTypeMeta{ID: 1016, Code: "dr_reinvest_src"}
	TradingBalanceTypeMetaDrWithdrawProfit             = TradingBalanceTypeMeta{ID: 1020, Code: "dr_withdraw_profit"}
	TradingBalanceTypeMetaDrProductEventTicket         = TradingBalanceTypeMeta{ID: 10002, Code: "dr_product_event_ticket"}
	TradingBalanceTypeMetaDrProductTravelTicket        = TradingBalanceTypeMeta{ID: 10004, Code: "dr_product_travel_ticket"}
	TradingBalanceTypeMetaDrProductMallOrder           = TradingBalanceTypeMeta{ID: 10006, Code: "dr_product_mall_order"}
	TradingBalanceTypeMetaDrProductFlightTicket        = TradingBalanceTypeMeta{ID: 10008, Code: "dr_product_flight_ticket"}

	TradingBalanceTypeMetaMap     map[uint32]TradingBalanceTypeMeta
	TradingBalanceTypeMetaCodeMap map[string]TradingBalanceTypeMeta
)

var LegacyBalanceTxnTypeCodeNameMap = map[uint32]string{
	1001: "ADJUST_ADD",
	1003: "TRANSFER",
	1005: "DEPOSIT",
	1009: "PROFIT",
	1011: "AFFILIATE",
	1013: "USER_REWARD",
	1017: "REINVEST",
	1021: "REDEEM",

	1002:  "ADJUST_DEDUCT",
	1004:  "TRANSFER",
	1008:  "WITHDRAW",
	1020:  "TORQ_WITHDRAW",
	1016:  "REINVEST",
	10002: "INVOICE_EVENT",
	10004: "INVOICE_BOOKING",
	10006: "INVOICE_TMALL",
	10008: "INVOICE_FLIGHT",

	1007: "WITHDRAW::REVERSE",
	1015: "REINVEST::REVERSE",
	1019: "TORQ_WITHDRAW::REVERSE",
	1006: "DEPOSIT::REVERSE",
	1010: "PROFIT::REVERSE",
	1012: "AFFILIATE::REVERSE",
	1014: "USER_REWARD::REVERSE",

	10003: "INVOICE_BOOKING::REVERSE",
	10005: "INVOICE_TMALL::REVERSE",
	10007: "INVOICE_FLIGHT::REVERSE",
}

func init() {
	TradingBalanceTypeMetaMap = map[uint32]TradingBalanceTypeMeta{
		TradingBalanceTypeMetaCrSystem.ID:                    TradingBalanceTypeMetaCrSystem,
		TradingBalanceTypeMetaCrDeposit.ID:                   TradingBalanceTypeMetaCrDeposit,
		TradingBalanceTypeMetaCrWithdrawInvestmentReverse.ID: TradingBalanceTypeMetaCrWithdrawInvestmentReverse,
		TradingBalanceTypeMetaCrDailyProfit.ID:               TradingBalanceTypeMetaCrDailyProfit,
		TradingBalanceTypeMetaCrAffiliateCommission.ID:       TradingBalanceTypeMetaCrAffiliateCommission,
		TradingBalanceTypeMetaCrLeaderCommission.ID:          TradingBalanceTypeMetaCrLeaderCommission,
		TradingBalanceTypeMetaCrReinvestSrcReverse.ID:        TradingBalanceTypeMetaCrReinvestSrcReverse,
		TradingBalanceTypeMetaCrReinvestDst.ID:               TradingBalanceTypeMetaCrReinvestDst,
		TradingBalanceTypeMetaCrWithdrawProfitReverse.ID:     TradingBalanceTypeMetaCrWithdrawProfitReverse,
		TradingBalanceTypeMetaCrPromoCode.ID:                 TradingBalanceTypeMetaCrPromoCode,
		TradingBalanceTypeMetaCrProductTravelTicketRefund.ID: TradingBalanceTypeMetaCrProductTravelTicketRefund,
		TradingBalanceTypeMetaCrProductMallOrderRefund.ID:    TradingBalanceTypeMetaCrProductMallOrderRefund,
		TradingBalanceTypeMetaCrProductFlightTicketRefund.ID: TradingBalanceTypeMetaCrProductFlightTicketRefund,

		TradingBalanceTypeMetaDrSystem.ID:                     TradingBalanceTypeMetaDrSystem,
		TradingBalanceTypeMetaDrDepositReverse.ID:             TradingBalanceTypeMetaDrDepositReverse,
		TradingBalanceTypeMetaDrWithdrawInvestment.ID:         TradingBalanceTypeMetaDrWithdrawInvestment,
		TradingBalanceTypeMetaDrDailyProfitReverse.ID:         TradingBalanceTypeMetaDrDailyProfitReverse,
		TradingBalanceTypeMetaDrAffiliateCommissionReverse.ID: TradingBalanceTypeMetaDrAffiliateCommissionReverse,
		TradingBalanceTypeMetaDrLeaderCommissionReverse.ID:    TradingBalanceTypeMetaDrLeaderCommissionReverse,
		TradingBalanceTypeMetaDrReinvestSrc.ID:                TradingBalanceTypeMetaDrReinvestSrc,
		TradingBalanceTypeMetaDrWithdrawProfit.ID:             TradingBalanceTypeMetaDrWithdrawProfit,
		TradingBalanceTypeMetaDrProductEventTicket.ID:         TradingBalanceTypeMetaDrProductEventTicket,
		TradingBalanceTypeMetaDrProductTravelTicket.ID:        TradingBalanceTypeMetaDrProductTravelTicket,
		TradingBalanceTypeMetaDrProductMallOrder.ID:           TradingBalanceTypeMetaDrProductMallOrder,
		TradingBalanceTypeMetaDrProductFlightTicket.ID:        TradingBalanceTypeMetaDrProductFlightTicket,
	}

	TradingBalanceTypeMetaCodeMap = map[string]TradingBalanceTypeMeta{}
	for _, meta := range TradingBalanceTypeMetaMap {
		TradingBalanceTypeMetaCodeMap[meta.Code] = meta
	}
}
