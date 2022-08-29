package balancemod

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func GetUserBalance(ctx comcontext.Context, uid meta.UID, currency meta.Currency) (models.UserBalance, error) {
	return getOrCreateBalance(ctx, currency, uid, false)
}

func GetTransaction(
	ctx comcontext.Context,
	currency meta.Currency, uid meta.UID,
	txnType uint32, orderID uint64,
) (*models.WalletBalanceTxn, error) {
	var record models.WalletBalanceTxn
	err := database.Atomic(ctx, database.AliasWalletMaster, func(dbTxn *gorm.DB) error {
		return dbTxn.
			First(
				&record,
				&models.WalletBalanceTxn{
					Currency: currency,
					UID:      uid,
					Type:     txnType,
					OrderID:  orderID,
				},
			).
			Error
	})
	if database.IsDbError(err) {
		return nil, utils.WrapError(err)
	}
	if record.ID == 0 {
		return nil, nil
	} else {
		return &record, nil
	}
}

func AddTransaction(
	ctx comcontext.Context, currency meta.Currency, uid meta.UID, amount decimal.Decimal, txnType uint32, orderID uint64,
) (balanceTxn models.WalletBalanceTxn, err error) {
	currencyInfo, err := currencymod.GetCurrencyInfoFast(currency)
	if err != nil || !currencymod.IsValidWalletInfo(currencyInfo) {
		err = utils.IssueErrorf("currency `%v` is not supported", currency)
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
	if !amount.Equals(currencymod.NormalizeAmount(currency, amount)) {
		err = utils.IssueErrorf("max amount decimal place is %d | amount=%s", currencyInfo.DecimalPlaces, amount)
		return
	}

	logger := comlogging.GetLogger()
	balanceTxn = models.WalletBalanceTxn{
		Currency: currency,
		UID:      uid,
		Amount:   amount,
		Type:     txnType,
		OrderID:  orderID,
	}
	err = database.TransactionFromContext(ctx, database.AliasWalletMaster, func(dbTxn *gorm.DB) error {
		record, err := GetTransaction(ctx, currency, uid, txnType, orderID)
		if err != nil {
			return err
		}
		if record != nil {
			logger.
				WithContext(ctx).
				WithFields(
					logrus.Fields{
						"currency": currency,
						"uid":      uid,
						"amount":   amount,
						"type":     txnType,
						"order_id": orderID,
					}).
				Warn("add balance duplicated")
			if !record.Amount.Equal(amount) {
				return utils.IssueErrorf(
					"balance exists with amount `%s` is mismatched | amount=%s",
					record.Amount, amount,
				)
			}

			balanceTxn = *record
			return nil
		}

		balance, err := GetAndLockBalance(ctx, uid, currency)
		if err != nil {
			return err
		}

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
		comutils.PanicOnError(
			dbTxn.Create(&balanceTxn).Error,
		)

		balance.Amount = afterBalance
		balance.UpdateTime = now.Unix()
		balance.LatestTxnID = models.NewInt64(int64(balanceTxn.ID))
		comutils.PanicOnError(
			dbTxn.Save(&balance).Error,
		)

		return nil
	})
	return
}
