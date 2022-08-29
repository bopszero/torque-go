package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/torque-go/api/middleware"
	"gitlab.com/snap-clickstaff/torque-go/api/services/transaction/v1/controllers/balance"
	"gitlab.com/snap-clickstaff/torque-go/api/services/transaction/v1/controllers/payment"
	"gitlab.com/snap-clickstaff/torque-go/config"
)

func InitGroup(group *echo.Group) {
	macSecret := viper.GetString(config.KeyServiceTransactionMacSecret)
	group.Use(middleware.NewValidateMAC(macSecret))

	group.POST("/balance/user/get/", balance.GetUserBalance)
	group.POST("/balance/txn/add/", balance.AddTxn, middleware.LogRequestDefaultMiddleware)

	paymentGroup := group.Group("/payment")
	paymentGroup.Use(middleware.LogRequestDefaultMiddleware)
	paymentGroup.POST("/investment/deposit/account/get/", payment.GetDepositAccount)
	paymentGroup.POST("/investment/deposit/crawl/", payment.CrawlDeposit)
	paymentGroup.POST("/investment/deposit/submit/", payment.SubmitDeposit)
	paymentGroup.POST("/investment/deposit/approve/", payment.ApproveDeposit)
	paymentGroup.POST("/investment/withdraw/submit/", payment.SubmitInvestmentWithdraw)
	paymentGroup.POST("/investment/withdraw/reject/", payment.RejectInvestmentWithdraw)
	paymentGroup.POST("/investment/withdraw/cancel/", payment.CancelInvestmentWithdraw)
	paymentGroup.POST("/profit/withdraw/reject/", payment.RejectProfitWithdraw)
	paymentGroup.POST("/profit/withdraw/cancel/", payment.CancelProfitWithdraw)
	paymentGroup.POST("/profit/reinvest/approve/", payment.ApproveProfitReinvest)
	paymentGroup.POST("/profit/reinvest/reject/", payment.RejectProfitReinvest)
	paymentGroup.POST("/p2p/transfer/", payment.TransferP2P)

	companyGroup := group.Group("/system")
	companyGroup.Use(middleware.LogRequestDefaultMiddleware)
	companyGroup.POST("/withdrawal/account/generate/", SystemWithdrawalAccountGenerate)
	companyGroup.POST("/withdrawal/account/get/", SystemWithdrawalAccountGet)
	companyGroup.POST("/withdrawal/account/pull/", SystemWithdrawalAccountPull)
	companyGroup.POST("/withdrawal/transfer/meta/", SystemWithdrawalTransferMeta)
	companyGroup.POST("/withdrawal/transfer/submit/", SystemWithdrawalTransferSubmit)
	companyGroup.POST("/withdrawal/transfer/confirm/", SystemWithdrawalTransferConfirm)
	companyGroup.POST("/withdrawal/transfer/replace/", SystemWithdrawalTransferReplace)

	companyGroup.POST("/pool/bonus/leader/checkout/", SystemPoolBonusLeaderCheckout)
	companyGroup.POST("/pool/bonus/leader/execute/", SystemPoolBonusLeaderExecute)
}
