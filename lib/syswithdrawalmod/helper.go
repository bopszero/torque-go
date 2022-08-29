package syswithdrawalmod

import (
	"sort"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func GetConfig(coin blockchainmod.Coin) (withdrawalConfig WithdrawalConfig) {
	configMap := viper.GetStringMap(config.KeySystemWithdrawalConfigMap)
	value, ok := configMap[strings.ToLower(coin.GetIndexCode())]
	if !ok {
		return
	}

	if err := utils.DumpDataByJSON(&value, &withdrawalConfig); err != nil {
		comlogging.GetLogger().
			WithType(constants.LogTypeSystemForwarding).
			WithError(err).
			WithField("coin", coin.GetIndexCode()).
			Warnf("load config failed | err=%s", err.Error())
	}

	return
}

func getClientModMapSetting() (tClientModeMap, error) {
	clientModeMapObj, err := clientModeMapSettingProxy.Get()
	if err != nil {
		return nil, err
	}
	return clientModeMapObj.(tClientModeMap), nil
}

func newCoinClient(coin blockchainmod.Coin) (blockchainmod.Client, error) {
	getClientModMap, err := getClientModMapSetting()
	if err != nil {
		return nil, err
	}
	clientMode, ok := getClientModMap[coin.GetIndexCode()]
	if ok && strings.ToLower(clientMode) == blockchainmod.ClientModeSpare {
		return coin.NewClientSpare()
	} else {
		return coin.NewClientDefault()
	}
}

func updateTradingWithdrawals(
	ctx comcontext.Context,
	refCodes []string, status meta.SystemWithdrawalTxnStatus, txnHash string,
) error {
	now := time.Now()
	return database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		err = dbTxn.Model(&models.Withdraw{}).
			Where(dbquery.In(models.WithdrawColCode, refCodes)).
			Updates(&models.Withdraw{
				TxnHash:       txnHash,
				ExecuteStatus: status,
				UpdateTime:    now.Unix(),
			}).
			Error
		if err != nil {
			return utils.WrapError(err)
		}

		err = dbTxn.Model(&models.TorqueTxn{}).
			Where(dbquery.In(models.TorqueTxnColCode, refCodes)).
			Updates(&models.TorqueTxn{
				BlockchainHash: txnHash,
				ExecuteStatus:  status,
				UpdateTime:     now.Unix(),
			}).
			Error
		if err != nil {
			return utils.WrapError(err)
		}

		return nil
	})
}

func getRequestTxns(db *gorm.DB, requestID uint64) (txns []models.SystemWithdrawalTxn, err error) {
	err = db.
		Where(&models.SystemWithdrawalTxn{RequestID: requestID}).
		Where(dbquery.NotIn(
			models.SystemWithdrawalTxnColStatus,
			[]meta.SystemWithdrawalTxnStatus{
				constants.SystemWithdrawalTxnStatusReplaced,
				constants.SystemWithdrawalTxnStatusCancelled,
			},
		)).
		Find(&txns).
		Error
	if err != nil {
		err = utils.WrapError(err)
	}

	sort.Slice(
		txns,
		func(i, j int) bool {
			if txns[i].OutputIndex < txns[j].OutputIndex {
				return true
			}
			if txns[i].OutputIndex == txns[j].OutputIndex {
				return txns[i].CreateTime < txns[j].CreateTime
			}
			return false
		},
	)
	return
}

func genTransferByCodes(coin blockchainmod.Coin, codes []string) (
	transfers []Transfer, err error,
) {
	var (
		db                    = database.GetDbSlave()
		investmentWithdrawals []models.Withdraw
		profitWithdrawals     []models.TorqueTxn
	)
	err = db.
		Where(dbquery.In(models.WithdrawColCode, codes)).
		Where(&models.Withdraw{
			Currency:  coin.GetCurrency(),
			Network:   coin.GetNetwork(),
			Status:    constants.WithdrawStatusPendingTransfer,
			IsDeleted: models.NewBool(false),
		}).
		Find(&investmentWithdrawals).
		Error
	if err != nil {
		err = utils.WrapError(err)
		return
	}
	err = db.
		Where(dbquery.In(models.TorqueTxnColCode, codes)).
		Where(&models.TorqueTxn{
			Currency:   coin.GetCurrency(),
			Network:    coin.GetNetwork(),
			Status:     constants.WithdrawStatusPendingTransfer,
			IsReinvest: models.NewBool(false),
			IsDeleted:  models.NewBool(false),
		}).
		Find(&profitWithdrawals).
		Error
	if len(investmentWithdrawals)+len(profitWithdrawals) != len(codes) {
		err = utils.WrapError(constants.ErrorDataNotFound)
		return
	}

	transfers = make([]Transfer, 0, len(investmentWithdrawals)+len(profitWithdrawals))
	for _, w := range investmentWithdrawals {
		address, addrErr := coin.NormalizeAddress(w.Address)
		if addrErr != nil {
			err = addrErr
			return
		}
		transfer := Transfer{
			RefCode:   w.Code,
			Currency:  w.Currency,
			ToAddress: address,
			Amount:    w.Amount.Sub(w.Fee),
		}
		transfers = append(transfers, transfer)
	}
	for _, w := range profitWithdrawals {
		address, addrErr := coin.NormalizeAddress(w.Address)
		if addrErr != nil {
			err = addrErr
			return
		}
		transfer := Transfer{
			RefCode:   w.Code,
			Currency:  w.Currency,
			ToAddress: address,
			Amount:    w.CoinAmount.Sub(w.CoinFee),
		}
		transfers = append(transfers, transfer)
	}

	return transfers, nil
}
