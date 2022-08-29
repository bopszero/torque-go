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
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func SubmitProfitWithdraw(
	ctx comcontext.Context,
	uid meta.UID, amount decimal.Decimal, fromCoinExchangeRate decimal.Decimal,
	currency meta.Currency, network meta.BlockchainNetwork, address string,
	priceMarkup *meta.AmountMarkup,
) (torqueTxn models.TorqueTxn, err error) {
	if err = currencymod.ValidateTradingCurrency(currency); err != nil {
		return
	}
	coin, err := blockchainmod.GetCoin(currency, network)
	if err != nil {
		return
	}
	if _, err = coin.NormalizeAddress(address); err != nil {
		err = utils.WrapError(constants.ErrorAddress)
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

	networkCurrencyInfo := coin.GetModelNetworkCurrency()
	err = database.TransactionFromContext(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) error {
		now := time.Now()
		torqueTxn = models.TorqueTxn{
			UserID:    uid,
			CoinID:    coin.GetTradingID(),
			Currency:  coin.GetCurrency(),
			Network:   coin.GetNetwork(),
			Status:    constants.WithdrawStatusPendingConfirm,
			IsDeleted: models.NewBool(false),

			IsReinvest:   models.NewBool(false),
			Amount:       amount,
			Address:      address,
			ExchangeRate: toCoinExchangeRate,
			CoinAmount:   coinAmount,
			CoinFee:      networkCurrencyInfo.WithdrawalFee,

			CreateTime: now,
			UpdateTime: now.Unix(),
		}

		var createErr error
		for i := 0; i < 5; i++ {
			torqueTxn.Code = GenProfitWithdrawCode()
			if createErr = dbTxn.Create(&torqueTxn).Error; createErr == nil {
				break
			}
		}
		if createErr != nil {
			return utils.WrapError(err)
		}

		return nil
	})
	return
}

func markFailProfitWithdraw(
	ctx comcontext.Context, fromStatuses []string, toStatus string, torqueTxnID uint64, note string,
) error {
	return database.TransactionFromContext(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) error {
		now := time.Now()
		updatedDb := GetBalanceDB(dbquery.SelectForUpdate(dbTxn)).
			FilterProfitWithdraw(torqueTxnID).
			Where(dbquery.In(models.TorqueTxnColStatus, fromStatuses)).
			Updates(&models.TorqueTxn{
				Status:     toStatus,
				CloseTime:  now.Unix(),
				Note:       note,
				UpdateTime: now.Unix(),
			})
		if database.IsDbError(updatedDb.Error) {
			return updatedDb.Error
		}
		if updatedDb.RowsAffected < 1 {
			return constants.ErrorDataNotFound
		}

		return nil
	})
}

func CancelProfitWithdraw(ctx comcontext.Context, torqueTxnID uint64, note string) error {
	return markFailProfitWithdraw(
		ctx,
		[]string{constants.WithdrawStatusPendingConfirm}, constants.WithdrawStatusCanceled,
		torqueTxnID, note,
	)
}

func RejectProfitWithdraw(ctx comcontext.Context, torqueTxnID uint64, note string) error {
	return markFailProfitWithdraw(
		ctx,
		[]string{constants.WithdrawStatusPendingConfirm, constants.WithdrawStatusPendingTransfer},
		constants.WithdrawStatusRejected,
		torqueTxnID, note,
	)
}
