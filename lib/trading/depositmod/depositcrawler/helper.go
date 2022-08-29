package depositcrawler

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/trading/depositmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func isUserAddress(coin blockchainmod.Coin, address string) bool {
	if len(address) < 16 {
		return false
	}

	_, err := depositmod.GetDepositAddressByAddress(coin, address)
	if err == nil {
		return true
	}
	if utils.IsOurError(err, constants.ErrorCodeDataNotFound) {
		return false
	}

	comlogging.GetLogger().
		WithError(err).
		WithFields(logrus.Fields{
			"coin":    coin.GetIndexCode(),
			"address": address,
		}).
		Warnf("crawler refuses an address `%v`", address)
	return false
}

func setCrawledDeposit(
	ctx comcontext.Context, coin blockchainmod.Coin,
	fromAddress string, toAddress string, toIdx uint16,
	hash string, amount decimal.Decimal,
	blockHeight uint64, blockTime int64, confirmations uint64,
) (txn models.DepositCryptoTxn, err error) {
	toAddress, err = coin.NormalizeAddress(toAddress)
	if err != nil {
		return
	}

	txn = models.DepositCryptoTxn{
		Currency: coin.GetCurrency(),
		Network:  coin.GetNetwork(),
		ToIndex:  toIdx,
		Hash:     hash,
	}
	err = database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		err = dbquery.SelectForUpdate(dbTxn).
			First(&txn, &txn).
			Error
		if database.IsDbError(err) {
			return utils.WrapError(err)
		}

		now := time.Now()
		if txn.ID == 0 {
			txn.Amount = utils.NormalizeTradingAmount(amount)
			txn.ToAddress = toAddress
			txn.CreateTime = now.Unix()
			txn.IsAccepted = models.NewBool(false)
		}

		txn.FromAddress = fromAddress
		txn.BlockHeight = blockHeight
		txn.BlockTime = blockTime
		txn.Confirmations = comutils.MaxUint64(txn.Confirmations, confirmations)
		txn.UpdateTime = now.Unix()
		if err = dbTxn.Save(&txn).Error; err != nil {
			return utils.WrapError(err)
		}
		return nil
	})
	return
}

func unsetCrawledDeposit(
	ctx comcontext.Context, coin blockchainmod.Coin,
	fromAddress string, toAddress string, toIdx uint16,
	hash string, amount decimal.Decimal, blockHeight uint64,
) (txn models.DepositCryptoTxn, err error) {
	toAddress, err = coin.NormalizeAddress(toAddress)
	if err != nil {
		return
	}

	txn = models.DepositCryptoTxn{
		Currency: coin.GetCurrency(),
		Network:  coin.GetNetwork(),
		ToIndex:  toIdx,
		Hash:     hash,
	}
	err = database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		err = dbquery.SelectForUpdate(dbTxn).
			First(&txn, &txn).
			Error
		if database.IsDbError(err) {
			return utils.WrapError(err)
		}
		if txn.ID == 0 {
			return nil
		}
		if err = dbTxn.Delete(&txn).Error; err != nil {
			return utils.WrapError(err)
		}
		return nil
	})
	return
}

func UpdateCrawledBlock(ctx comcontext.Context, coin blockchainmod.Coin, blockHeight uint64) error {
	return database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) error {
		if err := UpdateCrawledBlockCurrencyInfo(ctx, coin, blockHeight); err != nil {
			return err
		}
		if err := UpdateCrawledBlockCryptoTxns(ctx, coin, blockHeight); err != nil {
			return err
		}
		return nil
	})
}

func UpdateCrawledBlockCurrencyInfo(ctx comcontext.Context, coin blockchainmod.Coin, blockHeight uint64) error {
	return database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		query := dbTxn.Model(&models.LegacyCurrencyInfo{}).
			Where(&models.LegacyCurrencyInfo{Currency: coin.GetCurrency()}).
			Where(dbquery.Lt(models.TradingCurrencyInfoColLatestCrawledBlockHeight, blockHeight))
		err = query.
			Updates(&models.LegacyCurrencyInfo{LatestCrawledBlockHeight: blockHeight}).
			Error
		if err != nil {
			return utils.WrapError(err)
		}
		return nil
	})
}

func UpdateCrawledBlockCryptoTxns(ctx comcontext.Context, coin blockchainmod.Coin, blockHeight uint64) error {
	return database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		var (
			minUpdateDate        = time.Now().Add(-3 * 24 * time.Hour)
			confirmationsExprSQL = fmt.Sprintf(
				"LEAST(%v, %v-%v+1)",
				ConfirmationsUpdateMax,
				blockHeight, models.DepositCryptoTxnColBlockHeight)
			needUpdateTxns []models.DepositCryptoTxn
		)
		err = dbTxn.
			Where(&models.DepositCryptoTxn{
				Network: coin.GetNetwork(),
			}).
			Where(dbquery.Gt(models.DepositCryptoTxnColBlockTime, minUpdateDate)).
			Where(dbquery.Lt(models.DepositCryptoTxnColBlockHeight, blockHeight)).
			Where(dbquery.Lt(models.DepositCryptoTxnColConfirmations, gorm.Expr(confirmationsExprSQL))).
			Find(&needUpdateTxns).
			Error
		if err != nil {
			return utils.WrapError(err)
		}
		if len(needUpdateTxns) == 0 {
			return
		}

		txnIDs := make([]uint64, 0, len(needUpdateTxns))
		for _, txn := range needUpdateTxns {
			txnIDs = append(txnIDs, txn.ID)
		}
		updateConfirmationsQuery := dbTxn.
			Model(&models.DepositCryptoTxn{}).
			Where(dbquery.In(models.CommonColID, txnIDs)).
			Update(models.DepositCryptoTxnColConfirmations, gorm.Expr(confirmationsExprSQL))
		if err = updateConfirmationsQuery.Error; err != nil {
			return utils.WrapError(err)
		}
		return nil
	})
}
