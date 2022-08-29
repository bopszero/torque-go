package syswithdrawalmod

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
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

func submitUtxoLikeBulk(
	ctx comcontext.Context, requestUID meta.UID,
	coin blockchainmod.Coin, systemAddress models.SystemWithdrawalAddress, transfers []Transfer,
) (request models.SystemWithdrawalRequest, err error) {
	client, err := newCoinClient(coin)
	if err != nil {
		return
	}
	txnSigner, err := coin.NewTxnSignerBulk(client, nil)
	if err != nil {
		return
	}
	if err = txnSigner.SetLeftoverAddress(systemAddress.Address); err != nil {
		return
	}

	srcAddressKey, err := systemAddress.Key.GetValue()
	if err != nil {
		return
	}
	if err = txnSigner.AddSrc(srcAddressKey, systemAddress.Address); err != nil {
		return
	}

	totalTransferAmount := decimal.Zero
	for _, transfer := range transfers {
		if transfer.Amount.LessThan(UtxoLikeDustAmountThreshold) {
			err = meta.NewMessageError(
				"The transfer `%v` is considered dust.",
				transfer.RefCode,
			)
			return
		}
		if err = txnSigner.AddDst(transfer.ToAddress, transfer.Amount); err != nil {
			return
		}
		totalTransferAmount = totalTransferAmount.Add(transfer.Amount)
	}
	if err = txnSigner.Sign(false); err != nil {
		return
	}

	feeAmount, err := txnSigner.GetEstimatedFee()
	if err != nil {
		return
	}
	now := time.Now()
	request = models.SystemWithdrawalRequest{
		AddressID: systemAddress.ID,

		Status:             constants.SystemWithdrawalRequestStatusInit,
		Currency:           coin.GetCurrency(),
		Network:            coin.GetNetwork(),
		Amount:             totalTransferAmount,
		AmountEstimatedFee: feeAmount.Value,

		CombinedTxnHash:     models.NewString(txnSigner.GetHash()),
		CombinedSignedBytes: models.NewBytes(txnSigner.GetRaw()),

		CreateUID:  requestUID,
		CreateTime: now.Unix(),
		UpdateTime: now.Unix(),
	}
	systemBalance, err := client.GetBalance(systemAddress.Address)
	if err != nil {
		return
	}
	requestTotalAmount := request.Amount.Add(request.AmountEstimatedFee)
	if systemBalance.LessThan(requestTotalAmount) {
		err = utils.WrapError(
			constants.ErrorAmountTooHigh.WithData(meta.O{
				"threshold": requestTotalAmount,
				"currency":  request.Currency,
			}),
		)
		return
	}

	err = database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		if err = dbTxn.Save(&request).Error; err != nil {
			return utils.WrapError(err)
		}

		now := time.Now()
		txns := make([]models.SystemWithdrawalTxn, len(transfers))
		for i, transfer := range transfers {
			txns[i] = models.SystemWithdrawalTxn{
				RequestID: request.ID,
				RefCode:   transfer.RefCode,

				Status:      constants.SystemWithdrawalTxnStatusInit,
				Currency:    coin.GetCurrency(),
				Hash:        models.NewString(txnSigner.GetHash()),
				ToAddress:   transfer.ToAddress,
				OutputIndex: uint64(i),

				CreateTime: now.Unix(),
				UpdateTime: now.Unix(),
			}
		}
		if err = dbTxn.CreateInBatches(txns, 200).Error; err != nil {
			return utils.WrapError(err)
		}
		return nil
	})
	return
}

func executeUtxoLikeTxns(
	ctx comcontext.Context, coin blockchainmod.Coin,
	request models.SystemWithdrawalRequest, txns []models.SystemWithdrawalTxn,
) (err error) {
	if !request.CombinedTxnHash.Valid {
		panic(utils.IssueErrorf(
			"system withdrawal request for `%v` doesn't have txn hash | request_id=%v",
			request.Currency, request.ID,
		))
	}

	client, err := newCoinClient(coin)
	if err != nil {
		return err
	}

	networkTxn, err := client.GetTxn(request.CombinedTxnHash.String)
	if err != nil && !utils.IsOurError(err, constants.ErrorCodeDataNotFound) {
		return err
	}

	if networkTxn != nil && networkTxn.GetHash() != "" {
		if networkTxn.GetConfirmations() < 1 {
			return
		}
		if networkTxn.GetLocalStatus() != constants.BlockchainTxnStatusSucceeded {
			return
		}
		return completeUtxoLikeTxns(ctx, request, txns, networkTxn.GetFee())
	}

	return pushUtxoLikeTxns(ctx, coin, request, txns)
}

func pushUtxoLikeTxns(
	ctx comcontext.Context, coin blockchainmod.Coin,
	request models.SystemWithdrawalRequest, txns []models.SystemWithdrawalTxn,
) error {
	client, err := newCoinClient(coin)
	if err != nil {
		return err
	}

	if !request.CombinedSignedBytes.Valid {
		panic(utils.IssueErrorf(
			"system withdrawal request for `%v` is not signed | request_id=%v",
			request.Currency, request.ID,
		))
	}

	txnBytes := []byte(request.CombinedSignedBytes.String)
	if err := client.PushTxnRaw(txnBytes); err != nil {
		return err
	}

	return database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		var (
			now = time.Now()

			txnIDs   = make([]uint64, len(txns))
			refCodes = make([]string, len(txns))
		)
		for i, txn := range txns {
			txnIDs[i] = txn.ID
			refCodes[i] = txn.RefCode
		}
		err = dbTxn.
			Model(&models.SystemWithdrawalTxn{}).
			Where(dbquery.In(models.CommonColID, txnIDs)).
			Updates(&models.SystemWithdrawalTxn{
				Status:     constants.SystemWithdrawalTxnStatusPushed,
				UpdateTime: now.Unix(),
			}).
			Error
		if err != nil {
			return utils.WrapError(err)
		}

		return updateTradingWithdrawals(
			ctx,
			refCodes,
			constants.SystemWithdrawalTxnStatusPushed, request.CombinedTxnHash.String,
		)
	})
}

func completeUtxoLikeTxns(
	ctx comcontext.Context,
	request models.SystemWithdrawalRequest, txns []models.SystemWithdrawalTxn,
	feeAmount meta.CurrencyAmount,
) error {
	baseFeeAmount := currencymod.ConvertAmountF(feeAmount, request.Currency)
	return database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		now := time.Now()

		request.Status = constants.SystemWithdrawalRequestStatusCompleted
		request.UpdateTime = now.Unix()
		if err = dbTxn.Save(&request).Error; err != nil {
			return utils.WrapError(err)
		}

		var (
			txnIDs   = make([]uint64, len(txns))
			refCodes = make([]string, len(txns))
		)
		for i, txn := range txns {
			txnIDs[i] = txn.ID
			refCodes[i] = txn.RefCode
		}
		err = dbTxn.
			Model(&models.SystemWithdrawalTxn{}).
			Where(dbquery.In(models.CommonColID, txnIDs)).
			Updates(&models.SystemWithdrawalTxn{
				FeeAmount:  baseFeeAmount.Value,
				Status:     constants.SystemWithdrawalTxnStatusSuccess,
				UpdateTime: now.Unix(),
			}).
			Error
		if err != nil {
			return utils.WrapError(err)
		}

		return updateTradingWithdrawals(
			ctx,
			refCodes,
			constants.SystemWithdrawalTxnStatusSuccess, request.CombinedTxnHash.String,
		)
	})
}
