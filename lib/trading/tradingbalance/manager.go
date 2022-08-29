package tradingbalance

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func GetUserBalance(ctx comcontext.Context, uid meta.UID, currency meta.Currency) models.UserBalance {
	return getOrCreateBalance(ctx, currency, uid, false)
}

func GetUserBalances(ctx comcontext.Context, uid meta.UID) (balances []models.UserBalance) {
	comutils.PanicOnError(
		database.GetDbSlave().
			Find(&balances, &models.UserBalance{UID: uid}).
			Error,
	)

	existCurrencies := make(comtypes.HashSet, len(balances))
	for _, balance := range balances {
		existCurrencies.Add(balance.Currency)
	}
	currencyInfoMap := currencymod.GetAllCurrencyInfoMapFastF()
	for currency, info := range currencyInfoMap {
		if existCurrencies.Contains(currency) || !currencymod.IsValidTradingInfo(info) {
			continue
		}
		newBalance := GetUserBalance(ctx, uid, currency)
		balances = append(balances, newBalance)
	}

	return balances
}

func AddTransaction(
	ctx comcontext.Context,
	currency meta.Currency, uid meta.UID, amount decimal.Decimal, txnType uint32, ref string,
) (balanceTxn models.UserBalanceTxn, err error) {
	currencyInfo, err := currencymod.GetCurrencyInfoFast(currency)
	if err != nil || !currencymod.IsValidTradingInfo(currencyInfo) {
		err = utils.WrapError(constants.ErrorCurrency)
		return
	}
	if amount.IsZero() {
		err = utils.IssueErrorf("amount cannot be Zero")
		return
	} else if isCreditType(txnType) && !amount.IsPositive() {
		err = utils.IssueErrorf("credit type requires positive amount | amount=%s", amount)
		return
	} else if isDebitType(txnType) && !amount.IsNegative() {
		err = utils.IssueErrorf("debit type requires negative amount | amount=%s", amount)
		return
	}
	if !amount.Equals(utils.NormalizeTradingAmount(amount)) {
		err = utils.IssueErrorf(
			"max amount decimal place is %d | amount=%s",
			constants.AmountTradingMaxDecimalPlaces, amount)
		return
	}

	logger := comlogging.GetLogger()
	balanceTxn = models.UserBalanceTxn{
		Currency: currency,
		UserID:   uid,
		Amount:   amount,
		Type:     txnType,
		Ref:      ref,
	}
	err = database.TransactionFromContext(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		existsBalanceTxn, err := getBalanceTxn(ctx, currency, uid, txnType, ref)
		if err != nil {
			return
		}
		if existsBalanceTxn != nil {
			logger.
				WithContext(ctx).
				WithFields(
					logrus.Fields{
						"currency": currency,
						"uid":      uid,
						"amount":   amount,
						"type":     txnType,
						"ref":      ref,
					}).
				Warn("add balance duplicated")
			if !existsBalanceTxn.Amount.Equal(amount) {
				return utils.IssueErrorf(
					"balance exists with amount `%s` is mismatched | amount=%s",
					existsBalanceTxn.Amount, amount,
				)
			}

			balanceTxn = *existsBalanceTxn
			return
		}

		balance := GetAndLockBalance(ctx, uid, currency)

		afterBalance := balance.Amount.Add(amount)
		if afterBalance.LessThan(decimal.Zero) {
			logger.
				WithContext(ctx).
				WithFields(logrus.Fields{
					"currency": currency,
					"uid":      uid,
					"balance":  balance.Amount,
					"amount":   amount,
				}).
				Error("balance not enough")
			return utils.IssueErrorf("balance not enough | balance=%s,amount=%s", balance.Amount, amount)
		}

		now := time.Now()
		balanceTxn.Balance = afterBalance
		balanceTxn.CreateTime = now.Unix()
		balanceTxn.ParentID = balance.LatestTxnID
		if err = dbTxn.Create(&balanceTxn).Error; err != nil {
			err = utils.WrapError(err)
			return
		}

		balance.Amount = afterBalance
		balance.UpdateTime = now.Unix()
		balance.LatestTxnID = models.NewInt64(int64(balanceTxn.ID))
		if err = dbTxn.Save(&balance).Error; err != nil {
			err = utils.WrapError(err)
			return
		}

		return nil
	})
	return
}

func RemoveUserBalanceTxnRange(ctx comcontext.Context, uid meta.UID, currency meta.Currency, fromTxnID uint64, toTxnID uint64) error {
	var fromTxn models.UserBalanceTxn
	var txns []models.UserBalanceTxn

	err := database.TransactionFromContext(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) error {
		userBalance := GetAndLockBalance(ctx, uid, currency)

		comutils.PanicOnError(dbTxn.
			First(&fromTxn, &models.UserBalanceTxn{ID: fromTxnID, UserID: uid, Currency: currency}).
			Error,
		)

		comutils.PanicOnError(dbTxn.
			Where(dbquery.Gt(models.UserBalanceTxnColID, fromTxnID)).
			Order(dbquery.OrderAsc(models.UserBalanceTxnColID)).
			Find(&txns, &models.UserBalanceTxn{UserID: uid, Currency: currency}).
			Error,
		)

		userBalance.Amount = fromTxn.Balance.Sub(fromTxn.Amount)
		userBalance.LatestTxnID = fromTxn.ParentID
		userBalance.UpdateTime = time.Now().Unix()
		comutils.PanicOnError(
			dbTxn.Save(&userBalance).Error,
		)
		comutils.PanicOnError(
			dbTxn.Delete(&fromTxn).Error,
		)

		for _, txn := range txns {
			dbTxn.Delete(&txn)

			if txn.ID > toTxnID {
				_, err := AddTransaction(ctx, currency, uid, txn.Amount, txn.Type, txn.Ref)
				comutils.PanicOnError(err)
			}
		}

		return nil
	})

	return err
}
