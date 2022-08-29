package balancemod

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func getOrCreateBalance(
	ctx comcontext.Context,
	currency meta.Currency, uid meta.UID, lock bool,
) (balance models.UserBalance, err error) {
	currencyInfo, err := currencymod.GetCurrencyInfoFast(currency)
	if err != nil {
		return
	}
	if !currencymod.IsValidWalletInfo(currencyInfo) {
		err = utils.WrapError(constants.ErrorCurrency)
		return
	}
	if !lock {
		err = database.GetDbF(database.AliasWalletSlave).
			First(&balance, &models.UserBalance{Currency: currency, UID: uid}).
			Error
		if database.IsDbError(err) {
			err = utils.WrapError(err)
			return
		}
		if balance.ID != 0 {
			return
		}
	}
	err = database.Atomic(ctx, database.AliasMainMaster, func(dbMainTxn *gorm.DB) error {
		return database.Atomic(ctx, database.AliasWalletMaster, func(dbTxn *gorm.DB) (err error) {
			var dbGetTxn *gorm.DB
			if lock {
				dbGetTxn = dbquery.SelectForUpdate(dbTxn)
			} else {
				dbGetTxn = dbTxn
			}

			err = dbGetTxn.First(&balance, &models.UserBalance{Currency: currency, UID: uid}).Error
			if database.IsDbError(err) {
				return utils.WrapError(err)
			}
			if balance.ID != 0 {
				return nil
			}

			var lockUser models.User
			err = dbquery.SelectForUpdate(dbMainTxn).First(&lockUser, &models.User{ID: uid}).Error
			if database.IsDbError(err) {
				return utils.WrapError(err)
			}
			if lockUser.ID == 0 {
				return utils.IssueErrorf("cannot create balance on a non-existing user | uid=%v", uid)
			}

			newBalance, err := createBalance(ctx, currency, uid)
			if err != nil {
				return utils.WrapError(err)
			}

			balance = *newBalance
			return nil
		})
	})
	return
}

func createBalance(ctx comcontext.Context, currency meta.Currency, uid meta.UID) (*models.UserBalance, error) {
	var balance models.UserBalance
	err := database.Atomic(ctx, database.AliasWalletMaster, func(walletDbTxn *gorm.DB) error {
		err := walletDbTxn.First(&balance, &models.UserBalance{Currency: currency, UID: uid}).Error
		if database.IsDbError(err) {
			return err
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
		return walletDbTxn.Create(&balance).Error
	})
	if err != nil {
		return nil, utils.WrapError(err)
	}

	return &balance, nil
}

func GetAndLockBalance(ctx comcontext.Context, uid meta.UID, currency meta.Currency) (models.UserBalance, error) {
	return getOrCreateBalance(ctx, currency, uid, true)
}

func isCreditType(txnType uint32) bool {
	return txnType&1 == 1
}

func isDebitType(txnType uint32) bool {
	return !isCreditType(txnType)
}

func GetTorqueWithdrawWithDb(db *gorm.DB, ID uint64) (*models.TorqueTxn, error) {
	var torqueTxn models.TorqueTxn
	err := db.First(
		&torqueTxn,
		&models.TorqueTxn{
			ID:         ID,
			IsReinvest: models.NewBool(false),
			IsDeleted:  models.NewBool(false),
		},
	).Error
	if err != nil {
		return nil, err
	}

	return &torqueTxn, nil
}
