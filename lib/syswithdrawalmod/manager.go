package syswithdrawalmod

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/lockmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func GenerateAccount(ctx comcontext.Context, currency meta.Currency, network meta.BlockchainNetwork) (
	sysAddress models.SystemWithdrawalAddress, err error,
) {
	lock, err := lockmod.LockSimple("sys:withdrawal:account:%v-%v", currency, network)
	if err != nil {
		return
	}
	defer lock.Unlock()

	coin, err := blockchainmod.GetCoin(currency, network)
	if err != nil {
		return
	}
	keyHolder, err := coin.NewKey()
	if err != nil {
		return
	}

	now := time.Now()
	sysAddress = models.SystemWithdrawalAddress{
		Network:  coin.GetNetwork(),
		Currency: coin.GetCurrency(),
		Address:  keyHolder.GetAddress(),
		Key:      models.NewSystemWithdrawalAddressKeyEncryptedField(keyHolder.GetPrivateKey()),
		Status:   constants.CommonStatusActive,

		CreateTime: now.Unix(),
		UpdateTime: now.Unix(),
	}
	err = database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		err = dbTxn.
			Model(&models.SystemWithdrawalAddress{}).
			Where(&models.SystemWithdrawalAddress{
				Network:  network,
				Currency: currency,
				Status:   constants.CommonStatusActive,
			}).
			Updates(&models.SystemWithdrawalAddress{
				Status:     constants.CommonStatusInactive,
				UpdateTime: now.Unix(),
			}).
			Error
		if err != nil {
			return err
		}

		if err = dbTxn.Save(&sysAddress).Error; err != nil {
			return utils.WrapError(err)
		}

		return nil
	})
	return
}

func GetAccount(ctx comcontext.Context, coin blockchainmod.Coin) (
	sysAddress models.SystemWithdrawalAddress, err error,
) {
	err = database.GetDbSlave().
		Where(&models.SystemWithdrawalAddress{
			Network:  coin.GetNetwork(),
			Currency: coin.GetCurrency(),
			Status:   constants.CommonStatusActive,
		}).
		Last(&sysAddress).
		Error
	if err != nil {
		err = utils.WrapError(err)
	}
	return
}

func PullAccount(
	ctx comcontext.Context,
	currency meta.Currency, network meta.BlockchainNetwork,
	accountNo, pullAddress string,
) (systemAddress models.SystemWithdrawalAddress, hash string, err error) {
	coin, err := blockchainmod.GetCoin(currency, network)
	if err != nil {
		return
	}
	if accountNo, err = coin.NormalizeAddress(accountNo); err != nil {
		return
	}
	if pullAddress, err = coin.NormalizeAddress(pullAddress); err != nil {
		return
	}

	withdrawalConfig := GetConfig(coin)
	systemPullAddress, err := coin.NormalizeAddress(withdrawalConfig.PullAddress)
	if err != nil {
		return
	}
	if systemPullAddress == "" {
		err = utils.IssueErrorf(
			"system withdrawal needs address to pull balance | currency=%v",
			currency)
		return
	}
	if systemPullAddress, err = coin.NormalizeAddress(withdrawalConfig.PullAddress); err != nil {
		return
	}
	if systemPullAddress != pullAddress {
		err = meta.NewMessageError(
			"System withdrawal pull address `%v` for `%v` is mismatched",
			pullAddress, currency)
		return
	}

	err = database.GetDbSlave().
		Where(&models.SystemWithdrawalAddress{
			Network:  network,
			Currency: currency,
			Address:  accountNo,
			Status:   constants.CommonStatusActive,
		}).
		First(&systemAddress).
		Error
	if err != nil {
		err = utils.WrapError(err)
		return
	}

	client, err := newCoinClient(coin)
	if err != nil {
		return
	}

	srcAddressKey, err := systemAddress.Key.GetValue()
	if err != nil {
		return
	}
	pushTxn := func(txnSigner blockchainmod.SingleTxnSigner) (hash string, err error) {
		if err = txnSigner.SetSrc(srcAddressKey, systemAddress.Address); err != nil {
			return
		}
		if err = txnSigner.SetMoveDst(systemPullAddress); err != nil {
			return
		}
		if err = txnSigner.Sign(false); err != nil {
			return
		}

		hash = txnSigner.GetHash()
		if err = client.PushTxnRaw(txnSigner.GetRaw()); err != nil {
			err = meta.NewMessageError(
				"System withdrawal pushes a txn `%v` for %v has error: %v",
				hash, currency, err.Error(),
			)
			return
		}

		return
	}

	txnSigner, err := coin.NewTxnSignerSingle(client, nil)
	if err != nil {
		return
	}
	hash, err = pushTxn(txnSigner)
	if !utils.IsOurError(err, constants.ErrorCodeAmountTooLow) {
		return
	}

	if !blockchainmod.IsNetworkCoin(coin) {
		networkCurrency := coin.GetNetworkCurrency()
		var (
			networkCoin  = blockchainmod.GetCoinF(networkCurrency, coin.GetNetwork())
			netClient    blockchainmod.Client
			netTxnSigner blockchainmod.SingleTxnSigner
		)
		netClient, err = newCoinClient(networkCoin)
		if err != nil {
			return
		}
		netTxnSigner, err = networkCoin.NewTxnSignerSingle(netClient, nil)
		if err != nil {
			return
		}
		hash, err = pushTxn(netTxnSigner)
	}

	return
}

func SubmitBulk(
	ctx comcontext.Context, requestUID meta.UID,
	currency meta.Currency, network meta.BlockchainNetwork, srcAddress string,
	codes []string, totalAmount decimal.Decimal,
) (_ models.SystemWithdrawalRequest, err error) {
	var (
		db   = database.GetDbSlave()
		coin = blockchainmod.GetCoinF(currency, network)
	)
	if srcAddress, err = coin.NormalizeAddress(srcAddress); err != nil {
		return
	}

	var systemAddress models.SystemWithdrawalAddress
	err = db.
		Where(&models.SystemWithdrawalAddress{
			Network:  coin.GetNetwork(),
			Currency: coin.GetCurrency(),
			Address:  srcAddress,
			Status:   constants.CommonStatusActive,
		}).
		First(&systemAddress).
		Error
	if err != nil {
		err = utils.WrapError(err)
		return
	}

	lockKey := fmt.Sprintf("sys:withdrawal:submit:%v", systemAddress.Address)
	lock, err := lockmod.Lock(ctx, lockKey, LockDurationSubmit, time.Second)
	if err != nil {
		return
	}
	defer lock.Unlock()

	var addressPendingRequest models.SystemWithdrawalRequest
	err = db.
		Where(models.SystemWithdrawalRequest{AddressID: systemAddress.ID}).
		Where(dbquery.NotIn(
			models.CommonColStatus,
			[]meta.SystemWithdrawalRequestStatus{
				constants.SystemWithdrawalRequestStatusCancelled,
				constants.SystemWithdrawalRequestStatusCompleted,
			})).
		First(&addressPendingRequest).
		Error
	if database.IsDbError(err) {
		err = utils.WrapError(err)
		return
	}
	if addressPendingRequest.ID > 0 {
		var latestPendingTxn models.SystemWithdrawalTxn
		err = db.
			Where(&models.SystemWithdrawalTxn{RequestID: addressPendingRequest.ID}).
			Where(dbquery.In(
				models.SystemWithdrawalTxnColStatus,
				[]meta.SystemWithdrawalTxnStatus{
					constants.SystemWithdrawalTxnStatusInit,
					constants.SystemWithdrawalTxnStatusPushed,
				}),
			).
			Order(dbquery.OrderAsc(models.SystemWithdrawalTxnColOutputIndex)).
			Take(&latestPendingTxn).
			Error
		if err != nil {
			err = utils.WrapError(err)
			return
		}
		if latestPendingTxn.Hash.String != "" {
			err = meta.NewMessageError(
				"Source address `%s` is being used, kindly wait until the txn `%s` is finished.",
				systemAddress.Address, latestPendingTxn.Hash.String,
			)
		} else {
			err = meta.NewMessageError(
				"Source address `%s` is being used, kindly wait until the request `%v` is finished.",
				systemAddress.Address, latestPendingTxn.Hash.String,
			)
		}
		return
	}

	transfers, err := genTransferByCodes(coin, codes)
	if err != nil {
		return
	}
	var (
		invalidAddressSet   = make(comtypes.HashSet)
		totalTransferAmount = decimal.Zero
	)
	for _, transfer := range transfers {
		if _, err := coin.NormalizeAddress(transfer.ToAddress); err != nil {
			invalidAddressSet.Add(transfer.ToAddress)
		}

		if !transfer.Amount.IsPositive() {
			err = utils.WrapError(constants.ErrorAmount)
			return
		}
		totalTransferAmount = totalTransferAmount.Add(transfer.Amount)
	}
	if !totalTransferAmount.Equal(totalAmount) {
		err = utils.WrapError(constants.ErrorAmount)
		return
	}
	if len(invalidAddressSet) > 0 {
		err = meta.NewMessageError(
			"System Withdrawal detects some invalid addresses %v.",
			invalidAddressSet,
		)
		return
	}

	if constants.BlockchainChannelUtxoCurrencySet.Contains(coin.GetCurrency()) {
		return submitUtxoLikeBulk(ctx, requestUID, coin, systemAddress, transfers)
	} else {
		return submitBalanceLikeBulk(ctx, requestUID, coin, systemAddress, transfers)
	}
}

func ConfirmRequest(ctx comcontext.Context, requestID uint64) (
	request models.SystemWithdrawalRequest, err error,
) {
	err = database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		if err = dbquery.SelectForUpdate(dbTxn).First(&request, requestID).Error; err != nil {
			return
		}

		if request.Status != constants.SystemWithdrawalRequestStatusInit {
			if request.Status == constants.SystemWithdrawalRequestStatusConfirmed {
				return
			}
			return utils.WrapError(constants.ErrorInvalidParams)
		}

		now := time.Now()

		txns, err := getRequestTxns(dbTxn, request.ID)
		if err != nil {
			return
		}
		refCodes := make([]string, len(txns))
		for i, txn := range txns {
			refCodes[i] = txn.RefCode
		}

		investmentWithdrawalQuery := dbTxn.
			Model(&models.Withdraw{}).
			Where(dbquery.In(models.WithdrawColCode, refCodes)).
			Where(dbquery.In(models.WithdrawColExecuteStatus, ConfirmAcceptStatuses)).
			Updates(&models.Withdraw{
				ExecuteStatus: constants.SystemWithdrawalTxnStatusInit,
				UpdateTime:    now.Unix(),
			})
		if investmentWithdrawalQuery.Error != nil {
			return utils.WrapError(err)
		}
		profitWithdrawalQuery := dbTxn.
			Model(&models.TorqueTxn{}).
			Where(dbquery.In(models.TorqueTxnColCode, refCodes)).
			Where(dbquery.In(models.TorqueTxnColExecuteStatus, ConfirmAcceptStatuses)).
			Updates(&models.TorqueTxn{
				ExecuteStatus: constants.SystemWithdrawalTxnStatusInit,
				UpdateTime:    now.Unix(),
			})
		if profitWithdrawalQuery.Error != nil {
			return utils.WrapError(err)
		}
		lockedRowCount := int(investmentWithdrawalQuery.RowsAffected + profitWithdrawalQuery.RowsAffected)
		if lockedRowCount != len(txns) {
			return meta.NewMessageError(
				"Request's withdrawals have been processed before, can only lock %v/%v records",
				lockedRowCount, len(txns),
			)
		}

		request.Status = constants.SystemWithdrawalRequestStatusConfirmed
		request.UpdateTime = now.Unix()
		if err = dbTxn.Save(&request).Error; err != nil {
			return utils.WrapError(err)
		}

		return nil
	})
	return
}

func CancelRequest(ctx comcontext.Context, requestID uint64) (
	request models.SystemWithdrawalRequest, err error,
) {
	err = database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		if err = dbquery.SelectForUpdate(dbTxn).First(&request, requestID).Error; err != nil {
			return utils.WrapError(err)
		}

		switch request.Status {
		case constants.SystemWithdrawalRequestStatusCancelled:
			return nil
		case constants.SystemWithdrawalRequestStatusInit,
			constants.SystemWithdrawalRequestStatusConfirmed,
			constants.SystemWithdrawalRequestStatusTransferring:
			request.Status = constants.SystemWithdrawalRequestStatusCancelled
			request.UpdateTime = time.Now().Unix()
		default:
			return utils.WrapError(constants.ErrorInvalidParams)
		}

		return dbTxn.Save(&request).Error
	})
	return
}

func ExecuteRequest(ctx comcontext.Context, requestID uint64) (
	request models.SystemWithdrawalRequest, err error,
) {
	err = database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		err = dbquery.SelectForUpdate(dbTxn).
			Where(dbquery.In(
				models.SystemWithdrawalRequestColStatus,
				[]meta.SystemWithdrawalRequestStatus{
					constants.SystemWithdrawalRequestStatusConfirmed,
					constants.SystemWithdrawalRequestStatusTransferring,
				})).
			First(&request, requestID).
			Error
		if err != nil {
			return utils.WrapError(err)
		}

		coin, err := blockchainmod.GetCoin(request.Currency, request.Network)
		if err != nil {
			return
		}
		txns, err := getRequestTxns(dbTxn, request.ID)
		if err != nil {
			return
		}

		if request.Status != constants.SystemWithdrawalRequestStatusTransferring {
			request.Status = constants.SystemWithdrawalRequestStatusTransferring
			request.UpdateTime = time.Now().Unix()
			if err = dbTxn.Save(&request).Error; err != nil {
				return utils.WrapError(err)
			}
		}

		if constants.BlockchainChannelUtxoCurrencySet.Contains(request.Currency) {
			return executeUtxoLikeTxns(ctx, coin, request, txns)
		} else {
			return executeBalanceLikeTxns(ctx, coin, request, txns)
		}
	})
	return
}

func ReplaceRequest(ctx comcontext.Context, requestID uint64, feeInfo blockchainmod.FeeInfo) (
	request models.SystemWithdrawalRequest, err error,
) {
	err = database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		err = dbquery.SelectForUpdate(dbTxn).
			Where(dbquery.In(
				models.SystemWithdrawalRequestColStatus,
				[]meta.SystemWithdrawalRequestStatus{
					constants.SystemWithdrawalRequestStatusConfirmed,
					constants.SystemWithdrawalRequestStatusTransferring,
				})).
			First(&request, requestID).
			Error
		if err != nil {
			return utils.WrapError(err)
		}

		coin, err := blockchainmod.GetCoin(request.Currency, request.Network)
		if err != nil {
			return
		}
		if coin.GetNetworkCurrency() != constants.CurrencyEthereum {
			return meta.NewMessageError(
				"System withdrawal only support replacing ETH transactions (not %v).",
				coin.GetCurrency(),
			)
		}

		txns, err := getRequestTxns(dbTxn, request.ID)
		if err != nil {
			return
		}
		client, err := newCoinClient(coin)
		if err != nil {
			return
		}

		var replaceTxns []models.SystemWithdrawalTxn
		for i, txn := range txns {
			if txn.Status != constants.SystemWithdrawalTxnStatusInit &&
				txn.Status != constants.SystemWithdrawalTxnStatusPushed {
				continue
			}

			networkTxn, err := client.GetTxn(txn.Hash.String)
			if err != nil && !utils.IsOurError(err, constants.ErrorCodeDataNotFound) {
				return err
			}

			if networkTxn != nil && networkTxn.GetHash() != "" && networkTxn.GetConfirmations() > 0 {
				continue
			}

			replaceTxns = txns[i:]
			break
		}

		return replaceEthereumTxns(ctx, request, replaceTxns, feeInfo)
	})
	return
}

func CancelExpiredRequests(ctx comcontext.Context, coin blockchainmod.Coin) error {
	now := time.Now()
	return database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		var toBeCancelledRequests []models.SystemWithdrawalRequest
		err = dbquery.SelectForUpdate(dbTxn).
			Model(&models.SystemWithdrawalRequest{}).
			Where(&models.SystemWithdrawalRequest{
				Status:   constants.SystemWithdrawalRequestStatusInit,
				Currency: coin.GetCurrency(),
				Network:  coin.GetNetwork(),
			}).
			Where(dbquery.Lte(
				models.SystemWithdrawalRequestColCreateTime,
				now.Add(-InitTimeout).Unix(),
			)).
			Find(&toBeCancelledRequests).
			Error
		if err != nil {
			return utils.WrapError(err)
		}
		if len(toBeCancelledRequests) == 0 {
			return
		}

		requestIDs := make([]uint64, len(toBeCancelledRequests))
		for i, r := range toBeCancelledRequests {
			requestIDs[i] = r.ID
		}
		txnsQuery := dbTxn.
			Model(&models.SystemWithdrawalTxn{}).
			Where(dbquery.In(models.SystemWithdrawalTxnColRequestID, requestIDs))
		err = txnsQuery.
			Updates(&models.SystemWithdrawalTxn{
				Status:      constants.SystemWithdrawalTxnStatusCancelled,
				Hash:        models.NewStringNull(),
				SignedBytes: models.NewStringNull(),
				UpdateTime:  now.Unix(),
			}).
			Error
		if err != nil {
			return utils.WrapError(err)
		}

		var txns []models.SystemWithdrawalTxn
		err = txnsQuery.Find(&txns).Error
		if err != nil {
			return utils.WrapError(err)
		}
		refCodes := make([]string, len(txns))
		for i, txn := range txns {
			refCodes[i] = txn.RefCode
		}
		err = updateTradingWithdrawals(ctx, refCodes, constants.SystemWithdrawalTxnStatusCancelled, "")
		if err != nil {
			return
		}

		err = dbTxn.
			Model(&models.SystemWithdrawalRequest{}).
			Where(dbquery.In(models.CommonColID, requestIDs)).
			Updates(&models.SystemWithdrawalRequest{
				Status:              constants.SystemWithdrawalRequestStatusCancelled,
				CombinedTxnHash:     models.NewStringNull(),
				CombinedSignedBytes: models.NewStringNull(),
				UpdateTime:          now.Unix(),
			}).
			Error
		if err != nil {
			return utils.WrapError(err)
		}

		return nil
	})
}
