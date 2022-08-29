package withdrawalmod

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/trading/tradingbalance"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/balancemod"
	"gorm.io/gorm"
)

func SubmitInvestmentWithdraw(
	ctx comcontext.Context, uid meta.UID, amount decimal.Decimal,
	currency meta.Currency, network meta.BlockchainNetwork, address string,
) (withdraw models.Withdraw, err error) {
	if err = currencymod.ValidateTradingCurrency(currency); err != nil {
		return
	}
	coin, err := blockchainmod.GetCoin(currency, network)
	if err != nil {
		return
	}
	if !coin.IsAvailable() {
		err = utils.WrapError(constants.ErrorFeatureNotSupport)
	}
	if _, err = coin.NormalizeAddress(address); err != nil {
		err = utils.WrapError(constants.ErrorAddress)
		return
	}

	var (
		now                 = time.Now()
		networkCurrencyInfo = coin.GetModelNetworkCurrency()
	)
	withdraw = models.Withdraw{
		CoinID:    coin.GetTradingID(),
		Currency:  coin.GetCurrency(),
		Network:   coin.GetNetwork(),
		UserID:    uid,
		Status:    constants.WithdrawStatusPendingConfirm,
		Amount:    amount,
		Fee:       networkCurrencyInfo.WithdrawalFee,
		Address:   address,
		IsDeleted: models.NewBool(false),

		CreateTime: now,
		UpdateTime: now.Unix(),
	}
	err = database.TransactionFromContext(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) error {
		var createErr error
		for i := 0; i < 5; i++ {
			withdraw.Code = balancemod.GenInvestmentWithdrawCode()
			if createErr = dbTxn.Create(&withdraw).Error; createErr == nil {
				break
			}
		}
		if createErr != nil {
			return utils.WrapError(createErr)
		}

		negativeAmount := withdraw.Amount.Mul(constants.DecimalOneNegative)
		_, err := tradingbalance.AddTransaction(
			ctx,
			coin.GetCurrency(),
			withdraw.UserID,
			negativeAmount,
			constants.TradingBalanceTypeMetaDrWithdrawInvestment.ID,
			comutils.Stringify(withdraw.ID),
		)
		return err
	})
	return
}

func markFailInvestmentWithdraw(
	ctx comcontext.Context, fromStatus string, toStatus string, withdrawID uint64, note string,
) (withdraw models.Withdraw, err error) {
	var balanceTxn models.UserBalanceTxn
	err = database.TransactionFromContext(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		err = dbquery.SelectForUpdate(dbTxn).
			First(
				&withdraw,
				&models.Withdraw{
					ID:        withdrawID,
					Status:    fromStatus,
					IsDeleted: models.NewBool(false),
				},
			).
			Error
		if database.IsDbError(err) {
			err = utils.WrapError(err)
			return err
		}
		if withdraw.ID == 0 {
			return utils.WrapError(constants.ErrorDataNotFound)
		}

		now := time.Now()

		withdraw.Status = toStatus
		withdraw.CloseTime = now.Unix()
		withdraw.Note = note
		withdraw.UpdateTime = now.Unix()
		comutils.PanicOnError(
			dbTxn.Save(&withdraw).Error,
		)

		balanceTxn, err = tradingbalance.AddTransaction(
			ctx,
			withdraw.Currency,
			withdraw.UserID,
			withdraw.Amount,
			constants.TradingBalanceTypeMetaCrWithdrawInvestmentReverse.ID,
			comutils.Stringify(withdraw.ID),
		)

		return err
	})
	if err != nil {
		return
	}

	if toStatus == constants.WithdrawStatusRejected {
		if err := PushWithdrawRejectedNotificationAsync(balanceTxn); err != nil {
			comlogging.GetLogger().
				WithContext(ctx).
				WithError(err).
				WithField("withdraw_id", withdraw.ID).
				Warn("push investment withdraw notification failed | err=%v", err.Error())
		}
	}
	return
}

func RejectInvestmentWithdraw(ctx comcontext.Context, withdrawID uint64, note string) (models.Withdraw, error) {
	return markFailInvestmentWithdraw(
		ctx,
		constants.WithdrawStatusPendingTransfer, constants.WithdrawStatusRejected,
		withdrawID, note,
	)
}

func CancelInvestmentWithdraw(ctx comcontext.Context, withdrawID uint64, note string) (models.Withdraw, error) {
	return markFailInvestmentWithdraw(
		ctx,
		constants.WithdrawStatusPendingConfirm, constants.WithdrawStatusCanceled,
		withdrawID, note,
	)
}

func markFailProfitWithdraw(
	ctx comcontext.Context, fromStatus string, toStatus string, torqueTxnID uint64, note string,
) (*models.TorqueTxn, error) {
	var (
		now       = time.Now()
		torqueTxn models.TorqueTxn
	)
	err := database.TransactionFromContext(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		err = dbquery.SelectForUpdate(dbTxn).
			First(
				&torqueTxn,
				&models.TorqueTxn{
					ID:         torqueTxnID,
					Status:     fromStatus,
					IsReinvest: models.NewBool(false),
					IsDeleted:  models.NewBool(false),
				},
			).
			Error
		if database.IsDbError(err) {
			return err
		}
		if torqueTxn.ID == 0 {
			return utils.WrapError(constants.ErrorDataNotFound)
		}

		torqueTxn.Status = toStatus
		torqueTxn.CloseTime = now.Unix()
		torqueTxn.Note = note
		torqueTxn.UpdateTime = now.Unix()
		comutils.PanicOnError(
			dbTxn.Save(&torqueTxn).Error,
		)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &torqueTxn, nil
}

func RejectProfitWithdraw(ctx comcontext.Context, torqueTxnID uint64, note string) (*models.TorqueTxn, error) {
	return markFailProfitWithdraw(
		ctx,
		constants.WithdrawStatusPendingTransfer, constants.WithdrawStatusRejected,
		torqueTxnID, note,
	)
}

func CancelProfitWithdraw(ctx comcontext.Context, torqueTxnID uint64, note string) (*models.TorqueTxn, error) {
	return markFailProfitWithdraw(
		ctx,
		constants.WithdrawStatusPendingConfirm, constants.WithdrawStatusCanceled,
		torqueTxnID, note,
	)
}
