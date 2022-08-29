package tradingbalance

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
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

func getBalanceTxn(
	ctx comcontext.Context,
	currency meta.Currency, uid meta.UID, txnType uint32, ref string,
) (*models.UserBalanceTxn, error) {
	var balanceTxn models.UserBalanceTxn
	err := database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) error {
		return dbTxn.
			First(
				&balanceTxn,
				&models.UserBalanceTxn{
					Currency: currency,
					UserID:   uid,
					Type:     txnType,
					Ref:      ref},
			).
			Error
	})
	if database.IsDbError(err) {
		return nil, utils.WrapError(err)
	}
	if balanceTxn.ID == 0 {
		return nil, nil
	} else {
		return &balanceTxn, nil
	}
}

func getOrCreateBalance(
	ctx comcontext.Context,
	currency meta.Currency, uid meta.UID, lock bool,
) (balance models.UserBalance) {
	currencyInfo, err := currencymod.GetCurrencyInfoFast(currency)
	comutils.PanicOnError(err)
	if !currencymod.IsValidTradingInfo(currencyInfo) {
		panic(utils.WrapError(constants.ErrorCurrency))
	}
	if !lock {
		err := database.GetDbSlave().First(&balance, &models.UserBalance{Currency: currency, UID: uid}).Error
		if database.IsDbError(err) {
			panic(err)
		}
		if balance.ID != 0 {
			return balance
		}
	}
	err = database.TransactionFromContext(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		if lock {
			dbTxn = dbquery.SelectForUpdate(dbTxn)
		}
		err = dbTxn.
			First(&balance, &models.UserBalance{Currency: currency, UID: uid}).
			Error
		if database.IsDbError(err) {
			return utils.WrapError(err)
		}
		if balance.ID != 0 {
			return nil
		}
		balance, err = createBalance(ctx, currency, uid)
		return
	})
	comutils.PanicOnError(err)

	return balance
}

func createBalance(ctx comcontext.Context, currency meta.Currency, uid meta.UID) (balance models.UserBalance, err error) {
	err = database.TransactionFromContext(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		var user models.User
		err = dbquery.SelectForUpdate(dbTxn).First(&user, &models.User{ID: uid}).Error
		if database.IsDbError(err) {
			return utils.WrapError(err)
		}
		if user.ID == 0 {
			return utils.IssueErrorf("cannot create balance on a non-existing user | uid=%v", uid)
		}

		err = dbTxn.
			First(&balance, &models.UserBalance{Currency: currency, UID: uid}).
			Error
		if database.IsDbError(err) {
			return utils.WrapError(err)
		}
		if balance.ID != 0 {
			return nil
		}

		balance = models.UserBalance{
			Currency:    currency,
			UID:         uid,
			Amount:      decimal.NewFromInt(0),
			LatestTxnID: sql.NullInt64{},

			UpdateTime: time.Now().Unix(),
		}
		if err = dbTxn.Create(&balance).Error; err != nil {
			return utils.WrapError(err)
		}
		return nil
	})
	return
}

func GetAndLockBalance(ctx comcontext.Context, uid meta.UID, currency meta.Currency) models.UserBalance {
	return getOrCreateBalance(ctx, currency, uid, true)
}

func isCreditType(txnType uint32) bool {
	return txnType&1 == 1
}

func isDebitType(txnType uint32) bool {
	return !isCreditType(txnType)
}
