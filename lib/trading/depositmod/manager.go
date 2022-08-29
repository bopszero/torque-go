package depositmod

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/auth/kycmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/lockmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/trading/tradingbalance"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func GetOrCreateDepositUserAddress(
	ctx comcontext.Context,
	coin blockchainmod.Coin, uid meta.UID,
) (address models.DepositUserAddress, err error) {
	err = database.GetDbSlave().
		First(
			&address,
			&models.DepositUserAddress{
				UID:      uid,
				Currency: coin.GetCurrency(),
				Network:  coin.GetNetwork(),
			},
		).
		Error
	if database.IsDbError(err) {
		err = utils.WrapError(err)
		return
	}

	if address.ID == 0 {
		address, err = CreateDepositUserAddress(ctx, coin, uid)
	}

	return
}

func CreateDepositUserAddress(
	ctx comcontext.Context, coin blockchainmod.Coin, uid meta.UID,
) (address models.DepositUserAddress, err error) {
	if err = currencymod.ValidateTradingCurrency(coin.GetCurrency()); err != nil {
		return
	}
	if err = kycmod.ValidateUserKYC(ctx, uid); err != nil {
		return
	}

	lock, err := lockmod.LockSimple("deposit:user_address")
	if err != nil {
		return
	}
	defer lock.Unlock()

	err = database.TransactionFromContext(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		err = dbTxn.
			First(
				&address,
				&models.DepositUserAddress{
					UID:      uid,
					Currency: coin.GetCurrency(),
					Network:  coin.GetNetwork(),
				},
			).
			Error
		if database.IsDbError(err) {
			return utils.WrapError(err)
		}
		if address.ID > 0 {
			return nil
		}

		var randomAvailAddress models.DepositAddressStock
		err = dbTxn.
			Order("RAND()").
			Limit(1).
			First(
				&randomAvailAddress,
				&models.DepositAddressStock{
					Currency: coin.GetCurrency(),
					Network:  coin.GetNetwork(),
					IsUsed:   models.NewBool(false),
				},
			).
			Error
		if database.IsDbError(err) {
			return utils.WrapError(err)
		}
		if randomAvailAddress.ID == 0 {
			return utils.IssueErrorf(
				"system ran out of address for coin %v:%v",
				coin.GetCurrency(), coin.GetNetwork(),
			)
		}

		normalizedAddress, err := coin.NormalizeAddress(randomAvailAddress.Address)
		if err != nil {
			return
		}

		now := time.Now()
		address = models.DepositUserAddress{
			UID:      uid,
			Currency: coin.GetCurrency(),
			Network:  coin.GetNetwork(),
			Address:  normalizedAddress,

			CreateTime: now.Unix(),
		}
		if err = dbTxn.Create(&address).Error; err != nil {
			err = utils.WrapError(err)
			return
		}

		randomAvailAddress.IsUsed = models.NewBool(true)
		if err = dbTxn.Save(&randomAvailAddress).Error; err != nil {
			err = utils.WrapError(err)
			return
		}

		return nil
	})
	return
}

func GetDepositAddressByAddress(coin blockchainmod.Coin, address string) (
	depositAddress models.DepositUserAddress, err error,
) {
	address, err = coin.NormalizeAddress(address)
	if err != nil {
		return
	}
	err = database.GetDbSlave().
		First(
			&depositAddress,
			&models.DepositUserAddress{
				Currency: coin.GetCurrency(),
				Network:  coin.GetNetwork(),
				Address:  address,
			},
		).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = constants.ErrorDataNotFound
		}
		err = utils.WrapError(err)
		return
	}

	return depositAddress, nil
}

func GetUserDepositAddress(uid meta.UID, coin blockchainmod.Coin) (address models.DepositUserAddress, err error) {
	err = database.GetDbSlave().
		First(
			&address,
			&models.DepositUserAddress{
				UID:      uid,
				Currency: coin.GetCurrency(),
				Network:  coin.GetNetwork(),
			},
		).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = constants.ErrorDataNotFound
		}
		err = utils.WrapError(err)
		return
	}

	return address, nil
}

func SubmitDeposit(
	ctx comcontext.Context,
	uid meta.UID, currency meta.Currency, network meta.BlockchainNetwork,
	txnHash string, txnIndex uint16, address string, amount decimal.Decimal,
) (deposit models.Deposit, err error) {
	if err = currencymod.ValidateTradingCurrency(currency); err != nil {
		return
	}
	coin, err := blockchainmod.GetCoin(currency, network)
	if err != nil {
		return
	}

	now := time.Now()
	deposit = models.Deposit{
		CoinID:   coin.GetTradingID(),
		Currency: coin.GetCurrency(),
		Network:  coin.GetNetwork(),
		UID:      uid,
		Status:   constants.DepositStatusPendingConfirmations,

		TxnHash:  txnHash,
		TxnIndex: txnIndex,
		Address:  address,
		Amount:   amount,

		CreateTime: now,
		UpdateTime: now.Unix(),
		IsReinvest: models.NewBool(false),
		IsDeleted:  models.NewBool(false),
	}
	err = database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		err = dbTxn.Create(&deposit).Error
		if err != nil {
			if database.IsDuplicateEntryError(err) {
				err = constants.ErrorDataExists
			}
			err = utils.WrapError(err)
		}
		return
	})
	return
}

func closeDeposit(ctx comcontext.Context, depositID uint64, isApproved bool, note string) (
	deposit models.Deposit, err error,
) {
	var balanceTxn models.UserBalanceTxn
	err = database.TransactionFromContext(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		err = dbquery.SelectForUpdate(dbTxn).
			First(
				&deposit,
				&models.Deposit{
					ID:         depositID,
					Status:     constants.DepositStatusPendingConfirmations,
					IsReinvest: models.NewBool(false),
					IsDeleted:  models.NewBool(false),
				},
			).
			Error
		if database.IsDbError(err) {
			err = utils.WrapError(err)
			return
		}
		if deposit.ID == 0 {
			err = utils.WrapError(constants.ErrorDataNotFound)
			return
		}

		now := time.Now()

		if isApproved {
			deposit.Status = constants.DepositStatusApproved
		} else {
			deposit.Status = constants.DepositStatusRejected
		}
		deposit.CloseTime = now.Unix()
		deposit.Note = note
		deposit.UpdateTime = now.Unix()
		if err = dbTxn.Save(&deposit).Error; err != nil {
			err = utils.WrapError(err)
			return
		}

		if isApproved {
			balanceTxn, err = tradingbalance.AddTransaction(
				ctx,
				deposit.Currency,
				deposit.UID,
				deposit.Amount,
				constants.TradingBalanceTypeMetaCrDeposit.ID,
				comutils.Stringify(deposit.ID),
			)
		}
		return
	})
	if err != nil {
		return
	}

	if isApproved {
		err := PushDepositApprovedNotificationAsync(balanceTxn)
		if err != nil {
			comlogging.GetLogger().
				WithContext(ctx).
				WithError(err).
				WithFields(logrus.Fields{
					"deposit_id":     depositID,
					"balance_txn_id": balanceTxn.ID,
				}).
				Warn("push deposit approved notification failed | err=%v", err.Error())
		}
	}

	return
}

func ApproveDeposit(ctx comcontext.Context, depositID uint64, note string) (
	deposit models.Deposit, err error,
) {
	var balanceTxn models.UserBalanceTxn
	err = database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		err = dbquery.SelectForUpdate(dbTxn).
			First(
				&deposit,
				&models.Deposit{
					ID:         depositID,
					Status:     constants.DepositStatusPendingConfirmations,
					IsReinvest: models.NewBool(false),
					IsDeleted:  models.NewBool(false),
				},
			).
			Error
		if database.IsDbError(err) {
			err = utils.WrapError(err)
			return
		}
		if deposit.ID == 0 {
			err = utils.WrapError(constants.ErrorDataNotFound)
			return
		}

		now := time.Now()

		deposit.Status = constants.DepositStatusApproved
		deposit.CloseTime = now.Unix()
		deposit.Note = note
		deposit.UpdateTime = now.Unix()
		if err = dbTxn.Save(&deposit).Error; err != nil {
			err = utils.WrapError(err)
			return
		}

		balanceTxn, err = tradingbalance.AddTransaction(
			ctx,
			deposit.Currency,
			deposit.UID,
			deposit.Amount,
			constants.TradingBalanceTypeMetaCrDeposit.ID,
			comutils.Stringify(deposit.ID),
		)
		return
	})
	if err != nil {
		return
	}

	if err := PushDepositApprovedNotificationAsync(balanceTxn); err != nil {
		comlogging.GetLogger().
			WithError(err).
			WithFields(logrus.Fields{
				"deposit_id":     depositID,
				"balance_txn_id": balanceTxn.ID,
			}).
			Warn("push deposit approved notification failed | err=%v", err.Error())
	}

	return
}
