package crontasks

import (
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/syswithdrawalmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func SystemWithdrawalExecuteCurrencyRequest(coin blockchainmod.Coin) {
	defer apiutils.CronRunWithRecovery(
		"SystemWithdrawExecuteCurrencyRequest",
		meta.O{"coin": coin.GetIndexCode()},
	)

	logError := func(err error) {
		if err == nil {
			return
		}

		comlogging.GetLogger().
			WithError(err).
			Errorf("`SystemWithdrawExecuteCurrencyRequest` has error | err=%s", err.Error())
	}

	ctx := comcontext.NewContext()
	if err := syswithdrawalmod.CancelExpiredRequests(ctx, coin); err != nil {
		logError(err)
		return
	}
	err := database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		var (
			request       models.SystemWithdrawalRequest
			coinQueryCond = models.SystemWithdrawalRequest{
				Currency: coin.GetCurrency(),
				Network:  coin.GetNetwork(),
			}
		)
		err = dbTxn.
			Where(&coinQueryCond).
			Where(&models.SystemWithdrawalRequest{
				Status: constants.SystemWithdrawalRequestStatusTransferring,
			}).
			First(&request).
			Error
		if database.IsDbError(err) {
			return utils.WrapError(err)
		}
		if request.ID == 0 {
			err = dbTxn.
				Where(&coinQueryCond).
				Where(&models.SystemWithdrawalRequest{
					Status: constants.SystemWithdrawalRequestStatusConfirmed,
				}).
				First(&request).
				Error
			if database.IsDbError(err) {
				return utils.WrapError(err)
			}
		}
		if request.ID == 0 {
			return nil
		}
		if _, err = syswithdrawalmod.ExecuteRequest(ctx, request.ID); err != nil {
			return utils.WrapError(err)
		}
		return nil
	})
	if err != nil {
		logError(err)
		return
	}
}
