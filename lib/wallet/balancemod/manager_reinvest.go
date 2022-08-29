package balancemod

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/trading/depositmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/trading/tradingbalance"
	"gitlab.com/snap-clickstaff/torque-go/lib/usermod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func SubmitProfitReinvest(
	ctx comcontext.Context,
	uid meta.UID, amount decimal.Decimal, fromCoinExchangeRate decimal.Decimal,
	currency meta.Currency, targetUID meta.UID,
	priceMarkup *meta.AmountMarkup,
) (torqueTxn models.TorqueTxn, err error) {
	if err = currencymod.ValidateTradingCurrency(currency); err != nil {
		return
	}
	coin, err := blockchainmod.GetCoinNative(currency)
	if err != nil {
		return
	}
	targetAddress, err := depositmod.GetUserDepositAddress(targetUID, coin)
	if err != nil {
		return
	}

	systemPriceUsdt, err := currencymod.GetCurrencyPriceUsdtFast(currency)
	if err != nil {
		return
	}
	if priceMarkup != nil {
		systemPriceUsdt = priceMarkup.For(systemPriceUsdt)
	}
	systemFromCoinExchangeRate := currencymod.ConvertUsdtToTorque(systemPriceUsdt)
	if fromCoinExchangeRate.IsZero() ||
		systemFromCoinExchangeRate.GreaterThan(fromCoinExchangeRate) {
		fromCoinExchangeRate = systemFromCoinExchangeRate
	}

	var (
		toCoinExchangeRate = comutils.DecimalOneDivide(fromCoinExchangeRate)
		coinAmount         = utils.NormalizeTradingAmount(amount.Mul(toCoinExchangeRate))
	)
	if coinAmount.IsZero() {
		err = utils.WrapError(constants.ErrorAmountTooLow)
		return
	}

	err = database.TransactionFromContext(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) error {
		now := time.Now()

		torqueTxn = models.TorqueTxn{
			UserID:    uid,
			CoinID:    coin.GetTradingID(),
			Currency:  coin.GetCurrency(),
			Network:   constants.BlockchainNetworkTorque,
			Status:    constants.JudgeStatusPending,
			IsDeleted: models.NewBool(false),

			IsReinvest:   models.NewBool(true),
			Amount:       amount,
			Address:      targetAddress.Address,
			ExchangeRate: toCoinExchangeRate,
			CoinAmount:   coinAmount,

			CreateTime: now,
			UpdateTime: now.Unix(),
		}
		var createErr error
		for i := 0; i < 5; i++ {
			torqueTxn.Code = GenProfitReinvestCode()
			if createErr = dbTxn.Create(&torqueTxn).Error; createErr == nil {
				break
			}
		}
		if createErr != nil {
			return utils.WrapError(err)
		}

		deposit := models.Deposit{
			UID:       targetAddress.UID,
			CoinID:    coin.GetTradingID(),
			Currency:  coin.GetCurrency(),
			Network:   constants.BlockchainNetworkTorque,
			Status:    constants.JudgeStatusPending,
			IsDeleted: models.NewBool(false),

			TxnHash:     torqueTxn.Code,
			Address:     targetAddress.Address,
			Amount:      coinAmount,
			TorqueTxnID: models.NewUInt64(torqueTxn.ID),
			IsReinvest:  models.NewBool(true),

			CreateTime: now,
			UpdateTime: now.Unix(),
		}
		if err := dbTxn.Create(&deposit).Error; err != nil {
			return utils.WrapError(err)
		}

		return nil
	})
	return
}

func ApproveProfitReinvest(ctx comcontext.Context, torqueTxnID uint64, note string) (
	torqueTxn models.TorqueTxn, err error,
) {
	err = database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		err = GetBalanceDB(dbquery.SelectForUpdate(dbTxn)).
			FilterProfitReinvest(torqueTxnID).
			First(
				&torqueTxn,
				&models.TorqueTxn{
					Status: constants.JudgeStatusPending,
				},
			).
			Error
		if database.IsDbError(err) {
			err = utils.WrapError(err)
			return
		}
		if torqueTxn.ID == 0 {
			return utils.WrapError(constants.ErrorDataNotFound)
		}

		var deposit models.Deposit
		err = dbTxn.
			First(
				&deposit,
				&models.Deposit{
					TorqueTxnID: models.NewUInt64(torqueTxn.ID),
					Status:      constants.DepositStatusPendingReinvest,
					IsReinvest:  models.NewBool(true),
					IsDeleted:   models.NewBool(false),
				},
			).
			Error
		if database.IsDbError(err) {
			err = utils.WrapError(err)
			return
		}
		if torqueTxn.ID == 0 {
			return utils.WrapError(constants.ErrorDataNotFound)
		}
		targetUser, err := usermod.GetUser(deposit.UID)
		if err != nil {
			return
		}
		if targetUser.Status != constants.UserStatusActive {
			return utils.WrapError(constants.ErrorUserNotFound)
		}

		now := time.Now()

		torqueTxn.Status = constants.JudgeStatusApproved
		torqueTxn.CloseTime = now.Unix()
		torqueTxn.Note = note
		torqueTxn.UpdateTime = now.Unix()
		if err = dbTxn.Save(&torqueTxn).Error; err != nil {
			err = utils.WrapError(err)
			return
		}

		deposit.Status = constants.DepositStatusApproved
		deposit.CloseTime = now.Unix()
		deposit.Note = note
		deposit.UpdateTime = now.Unix()
		if err = dbTxn.Save(&deposit).Error; err != nil {
			err = utils.WrapError(err)
			return
		}

		_, err = tradingbalance.AddTransaction(
			ctx,
			deposit.Currency,
			deposit.UID,
			deposit.Amount,
			constants.TradingBalanceTypeMetaCrReinvestDst.ID,
			comutils.Stringify(torqueTxn.ID),
		)
		return
	})
	return
}

func RejectProfitReinvest(ctx comcontext.Context, torqueTxnID uint64, note string) (
	torqueTxn models.TorqueTxn, err error,
) {
	err = database.TransactionFromContext(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		err = GetBalanceDB(dbquery.SelectForUpdate(dbTxn)).
			FilterProfitReinvest(torqueTxnID).
			First(
				&torqueTxn,
				&models.TorqueTxn{
					Status: constants.JudgeStatusPending,
				},
			).
			Error
		if database.IsDbError(err) {
			err = utils.WrapError(err)
			return
		}
		if torqueTxn.ID == 0 {
			return utils.WrapError(constants.ErrorDataNotFound)
		}

		now := time.Now()

		updatedDepositDb := dbTxn.
			Where(&models.Deposit{
				TorqueTxnID: models.NewUInt64(torqueTxn.ID),
				Status:      constants.DepositStatusPendingReinvest,
				IsReinvest:  models.NewBool(true),
				IsDeleted:   models.NewBool(false),
			}).
			Model(&models.Deposit{}).
			Updates(&models.Deposit{
				Status:     constants.DepositStatusRejected,
				CloseTime:  now.Unix(),
				Note:       note,
				UpdateTime: now.Unix(),
			})
		if updatedDepositDb.RowsAffected != 1 {
			return utils.WrapError(constants.ErrorDataNotFound)
		}

		torqueTxn.Status = constants.JudgeStatusRejected
		torqueTxn.CloseTime = now.Unix()
		torqueTxn.Note = note
		torqueTxn.UpdateTime = now.Unix()
		if err = dbTxn.Save(&torqueTxn).Error; err != nil {
			err = utils.WrapError(err)
			return
		}

		return nil
	})
	return
}
