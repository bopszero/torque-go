package sysforwardingmod

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func getConfig(coin blockchainmod.Coin) (forwardConfig ForwardConfig) {
	configMap := viper.GetStringMap(config.KeySystemForwardingConfigMap)

	value, ok := configMap[strings.ToLower(coin.GetIndexCode())]
	if !ok {
		return
	}

	if err := utils.DumpDataByJSON(&value, &forwardConfig); err != nil {
		comlogging.GetLogger().
			WithType(constants.LogTypeSystemForwarding).
			WithError(err).
			Warnf("load config failed | err=%s", err.Error())
	}

	return
}

func newBlockchainClient(coin blockchainmod.Coin) (_ blockchainmod.Client, err error) {
	conf := getConfig(coin)
	if conf.ApiProvider == "spare" {
		return coin.NewClientSpare()
	} else {
		return coin.NewClientDefault()
	}
}

func countDepositsLeft(ctx comcontext.Context, excludeIDs []uint64, fromTime, toTime time.Time) int64 {
	var depositLeftCount int64
	err := database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) error {
		query := dbTxn.
			Model(&models.Deposit{}).
			Where(dbquery.Gte(models.DepositColCreateTime, fromTime)).
			Where(dbquery.Lt(models.DepositColCreateTime, toTime)).
			Where(
				fmt.Sprintf(
					"%v=? OR (%v=? AND %v=0)",
					models.DepositColStatus,
					models.DepositColStatus,
					models.DepositColForwardStatus,
				),
				constants.DepositStatusPendingReinvest,
				constants.DepositStatusApproved,
			)
		if len(excludeIDs) > 0 {
			query = query.Where(dbquery.NotIn(models.DepositColID, excludeIDs))
		}

		return query.Count(&depositLeftCount).Error
	})
	if err != nil {
		comlogging.GetLogger().
			WithType(constants.LogTypeSystemForwarding).
			WithContext(ctx).
			WithError(err).
			Errorf("unexpected database query error | err=%s", err.Error())
		return -1
	}

	return depositLeftCount
}

func ignoreTxns(ctx comcontext.Context, txns []models.SystemForwardingOrderTxn, note string) error {
	if len(txns) == 0 {
		return nil
	}

	depositIDs := make([]uint64, 0)
	for _, txn := range txns {
		depositIDs = append(depositIDs, txn.DepositID)
	}

	return database.Atomic(ctx, database.AliasInternalMaster, func(dbTxn *gorm.DB) error {
		return database.Atomic(ctx, database.AliasMainMaster, func(mainDbTxn *gorm.DB) (err error) {
			now := time.Now()

			for _, txn := range txns {
				txn.Status = constants.SystemForwardingOrderTxnStatusIgnored
				txn.Note = note
				txn.UpdateTime = now.Unix()
				if err = dbTxn.Save(&txn).Error; err != nil {
					return utils.WrapError(err)
				}
			}

			err = mainDbTxn.
				Model(&models.Deposit{}).
				Where(dbquery.In(models.DepositColID, depositIDs)).
				Updates(&models.Deposit{
					ForwardStatus: constants.SystemForwardingOrderTxnStatusIgnored,
					UpdateTime:    now.Unix(),
				}).
				Error
			if err != nil {
				return utils.WrapError(err)
			}

			return nil
		})
	})
}

func dumpTxnToReportItem(
	txn models.SystemForwardingOrderTxn, deposit *models.Deposit,
) ReportItem {
	var grossAmount decimal.Decimal
	if blockchainmod.IsNetworkCurrency(txn.Currency) {
		grossAmount = txn.Amount.Add(txn.Fee)
	} else {
		grossAmount = txn.Amount
	}

	item := ReportItem{
		ID:          txn.ID,
		DepositID:   txn.DepositID,
		Currency:    txn.Currency,
		Address:     txn.FromAddress,
		Status:      constants.SystemForwardingOrderTxnStatusNameMap[txn.Status],
		TxnHash:     txn.Hash.String,
		GrossAmount: grossAmount,
		FeeAmount:   txn.Fee,
		NetAmount:   txn.Amount,
		SignBalance: txn.SignBalance,
		Note:        txn.Note,
	}
	if deposit != nil {
		item.DepositAmount = deposit.Amount
	}
	if txn.SignBalanceTime > 0 {
		item.SignBalanceTime = time.
			Unix(txn.SignBalanceTime, 0).
			Format(constants.DateTimeFormatISO)
	}

	return item
}
