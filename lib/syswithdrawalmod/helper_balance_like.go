package syswithdrawalmod

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func getBalanceLikeTxnSigner(
	coin blockchainmod.Coin, client blockchainmod.Client, feeInfo blockchainmod.FeeInfo,
) (blockchainmod.SingleTxnSigner, error) {
	txnSigner, err := coin.NewTxnSignerSingle(client, &feeInfo)
	if err != nil {
		return nil, err
	}

	txnSigner.MarkAsPreferOffline()
	return txnSigner, nil
}

func getBalanceLikeFeeInfo(coin blockchainmod.Coin) (feeInfo blockchainmod.FeeInfo, err error) {
	if feeInfo, err = coin.GetDefaultFeeInfo(); err != nil {
		return
	}
	switch coin.GetCurrency() {
	case constants.CurrencyEthereum:
		feeInfo.SetLimitMaxQuantity(blockchainmod.GetConfig().EthereumContractDefaultGasLimit)
		break
	default:
		break
	}

	return feeInfo, nil
}

func submitBalanceLikeBulk(
	ctx comcontext.Context, requestUID meta.UID,
	coin blockchainmod.Coin, systemAddress models.SystemWithdrawalAddress, transfers []Transfer,
) (request models.SystemWithdrawalRequest, err error) {
	client, err := newCoinClient(coin)
	if err != nil {
		return
	}

	feeInfo, err := getBalanceLikeFeeInfo(coin)
	if err != nil {
		return
	}
	txnSigner, err := getBalanceLikeTxnSigner(coin, client, feeInfo)
	if err != nil {
		return
	}
	srcAddressKey, err := systemAddress.Key.GetValue()
	if err != nil {
		return
	}
	if err = txnSigner.SetSrc(srcAddressKey, systemAddress.Address); err != nil {
		return
	}

	nonce, err := client.GetNextNonce(systemAddress.Address)
	if err != nil {
		return
	}
	var (
		txns                    = make([]models.SystemWithdrawalTxn, 0, len(transfers))
		totalTransferAmount     = decimal.Zero
		totalEstimatedFeeAmount = decimal.Zero
	)
	for _, transfer := range transfers {
		if err = txnSigner.SetDst(transfer.ToAddress, transfer.Amount); err != nil {
			return
		}
		if err = txnSigner.SetNonce(nonce); err != nil {
			return
		}
		if err = txnSigner.Sign(true); err != nil {
			return
		}

		estimatedFeeAmount, feeErr := txnSigner.GetEstimatedFee()
		if feeErr != nil {
			err = feeErr
			return
		}
		nonceValue, nonceErr := nonce.GetNumber()
		if nonceErr != nil {
			err = nonceErr
			return
		}
		feeInfo := txnSigner.GetFeeInfo()
		txn := models.SystemWithdrawalTxn{
			RefCode:        transfer.RefCode,
			Status:         constants.SystemWithdrawalTxnStatusInit,
			Currency:       coin.GetCurrency(),
			Hash:           models.NewString(txnSigner.GetHash()),
			ToAddress:      transfer.ToAddress,
			OutputIndex:    nonceValue,
			FeePrice:       feeInfo.GetBasePrice(),
			FeeMaxQuantity: feeInfo.LimitMaxQuantity,
			SignedBytes:    models.NewBytes(txnSigner.GetRaw()),
		}
		txns = append(txns, txn)

		if nonce, err = nonce.Next(); err != nil {
			return
		}

		totalTransferAmount = totalTransferAmount.Add(transfer.Amount)
		totalEstimatedFeeAmount = totalEstimatedFeeAmount.Add(estimatedFeeAmount.Value)
	}

	request = models.SystemWithdrawalRequest{
		AddressID: systemAddress.ID,
		CreateUID: requestUID,

		Status:             constants.SystemWithdrawalRequestStatusInit,
		Currency:           coin.GetCurrency(),
		Network:            coin.GetNetwork(),
		Amount:             totalTransferAmount,
		AmountEstimatedFee: totalEstimatedFeeAmount,
	}
	systemBalance, err := client.GetBalance(systemAddress.Address)
	if err != nil {
		return
	}
	if coin.GetCurrency() == coin.GetNetworkCurrency() {
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
	} else {
		if systemBalance.LessThan(request.Amount) {
			err = utils.WrapError(
				constants.ErrorAmountTooHigh.WithData(meta.O{
					"threshold": request.Amount,
					"currency":  request.Currency,
				}),
			)
			return
		}

		var (
			networkCoin = blockchainmod.GetCoinF(coin.GetNetworkCurrency(), systemAddress.Network)
			networkErr  error
		)
		networkClient, networkErr := newCoinClient(networkCoin)
		if networkErr != nil {
			err = networkErr
			return
		}
		networkBalance, networkErr := networkClient.GetBalance(systemAddress.Address)
		if networkErr != nil {
			err = networkErr
			return
		}
		if networkBalance.LessThan(request.AmountEstimatedFee) {
			err = constants.ErrorBlockchainBalanceNotEnoughForFee.WithData(meta.O{
				"threshold": request.AmountEstimatedFee,
				"currency":  networkCoin.GetCurrency(),
			})
			return
		}
	}

	now := time.Now()
	err = database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		request.CreateTime = now.Unix()
		request.UpdateTime = now.Unix()
		if err = dbTxn.Save(&request).Error; err != nil {
			return utils.WrapError(err)
		}

		for i := range txns {
			txn := &txns[i]
			txn.RequestID = request.ID
			txn.CreateTime = now.Unix()
			txn.UpdateTime = now.Unix()
		}
		if err = dbTxn.CreateInBatches(txns, 200).Error; err != nil {
			return utils.WrapError(err)
		}
		return nil
	})
	return
}

func executeBalanceLikeTxns(
	ctx comcontext.Context, coin blockchainmod.Coin,
	request models.SystemWithdrawalRequest, txns []models.SystemWithdrawalTxn,
) (err error) {
	client, err := newCoinClient(coin)
	if err != nil {
		return err
	}

	refCodeToTxnsMap := make(map[string][]models.SystemWithdrawalTxn, len(txns))
	for _, txn := range txns {
		refCodeToTxnsMap[txn.RefCode] = append(refCodeToTxnsMap[txn.RefCode], txn)
	}
	var (
		updatedTxnMap = make(map[uint64]models.SystemWithdrawalTxn, len(txns))
		doneCount     int
	)
	for refCode, txns := range refCodeToTxnsMap {
		if executeBalanceLikeRefCode(coin, request, refCode, txns, client, updatedTxnMap) {
			doneCount++
		}
	}

	return database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		if len(updatedTxnMap) > 0 {
			now := time.Now()
			for _, txn := range updatedTxnMap {
				txn.UpdateTime = now.Unix()
				if err = dbTxn.Save(&txn).Error; err != nil {
					return utils.WrapError(err)
				}
				if txn.Status != constants.SystemWithdrawalTxnStatusReplaced {
					err = updateTradingWithdrawals(
						ctx,
						[]string{txn.RefCode}, txn.Status, txn.Hash.String)
					if err != nil {
						return err
					}
				}
			}
		}

		if doneCount < len(refCodeToTxnsMap) {
			return nil
		}

		request.Status = constants.SystemWithdrawalRequestStatusCompleted
		request.UpdateTime = time.Now().Unix()
		if err = dbTxn.Save(&request).Error; err != nil {
			return utils.WrapError(err)
		}

		return nil
	})
}

func executeBalanceLikeRefCode(
	coin blockchainmod.Coin, request models.SystemWithdrawalRequest,
	refCode string, txns []models.SystemWithdrawalTxn,
	client blockchainmod.Client, updatedTxnMap map[uint64]models.SystemWithdrawalTxn,
) (isDone bool) {
	logger := comlogging.GetLogger()
	for _, txn := range txns {
		isTxnDone, err := executeBalanceLikeTxn(request, txn, coin, client, updatedTxnMap)
		if err != nil {
			logger.
				WithType(constants.LogTypeSystemWithdrawal).
				WithError(err).
				WithField("txn_id", request.ID).
				Errorf("execute %v txn error | err=%s", coin.GetIndexCode(), err.Error())
		}
		if !isTxnDone {
			continue
		}
		isDone = isTxnDone

		updatedTxn, ok := updatedTxnMap[txn.ID]
		if !ok {
			updatedTxn = txn
		}
		switch updatedTxn.Status {
		case constants.SystemWithdrawalTxnStatusSuccess, constants.SystemWithdrawalTxnStatusFailed:
			for _, replacedTxn := range txns {
				if replacedTxn.ID != txn.ID {
					replacedTxn.Status = constants.SystemWithdrawalTxnStatusReplaced
					updatedTxnMap[replacedTxn.ID] = replacedTxn
				}
			}
			return true
		default:
			break
		}
	}
	return
}

func executeBalanceLikeTxn(
	request models.SystemWithdrawalRequest, txn models.SystemWithdrawalTxn,
	coin blockchainmod.Coin, client blockchainmod.Client,
	updatedTxnMap map[uint64]models.SystemWithdrawalTxn,
) (isDone bool, err error) {
	switch txn.Status {
	case constants.SystemWithdrawalTxnStatusInit,
		constants.SystemWithdrawalTxnStatusPushed:
		networkTxn, err := client.GetTxn(txn.Hash.String)
		if err != nil && !utils.IsOurError(err, constants.ErrorCodeDataNotFound) {
			return false, err
		}

		if networkTxn != nil && networkTxn.GetHash() != "" {
			if networkTxn.GetConfirmations() < 1 {
				return false, nil
			}

			switch networkTxn.GetLocalStatus() {
			case constants.BlockchainTxnStatusSucceeded:
				txn.Status = constants.SystemWithdrawalTxnStatusSuccess
				break
			case constants.BlockchainTxnStatusFailed:
				txn.Status = constants.SystemWithdrawalTxnStatusFailed
				break
			default:
				return false, nil
			}
			txn.FeeAmount = currencymod.ConvertAmountF(networkTxn.GetFee(), coin.GetNetworkCurrency()).Value

			updatedTxnMap[txn.ID] = txn
			return executeBalanceLikeTxn(request, txn, coin, client, updatedTxnMap)
		}

		txnBytes := []byte(txn.SignedBytes.String)
		if err := client.PushTxnRaw(txnBytes); err != nil {
			return false, err
		}

		txn.Status = constants.SystemWithdrawalTxnStatusPushed
		updatedTxnMap[txn.ID] = txn
		return false, nil

	case constants.SystemWithdrawalTxnStatusFailed:
		// TODO: Support replace
		comlogging.GetLogger().
			WithType(constants.LogTypeSystemWithdrawal).
			WithFields(logrus.Fields{
				"request_id": request.ID,
				"txn_id":     txn.ID,
				"hash":       txn.Hash.String,
			}).
			Errorf("execute %v txn failed | hash=%s", coin.GetIndexCode(), txn.Hash.String)
		return true, nil

	case constants.SystemWithdrawalTxnStatusSuccess:
		return true, nil

	default:
		panic(utils.IssueErrorf("not supported"))
	}
}
