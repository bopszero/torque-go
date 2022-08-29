package syswithdrawalmod

import (
	"time"

	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func replaceEthereumTxns(
	ctx comcontext.Context,
	request models.SystemWithdrawalRequest, txns []models.SystemWithdrawalTxn,
	inputFeeInfo blockchainmod.FeeInfo,
) (err error) {
	coin, err := blockchainmod.GetCoin(request.Currency, request.Network)
	if err != nil {
		return
	}
	networkFeeInfo, err := getBalanceLikeFeeInfo(coin)
	if err != nil {
		return
	}
	replaceFeeInfo := networkFeeInfo
	replaceFeeInfo.Price = inputFeeInfo.Price
	replaceFeeInfo.Currency = inputFeeInfo.Currency

	feeBasePrice := replaceFeeInfo.GetBasePrice()
	if feeBasePrice.LessThan(networkFeeInfo.GetBasePriceHigh()) {
		return meta.NewMessageError(
			"System withdrawal cannot replace with this low fee price `%v %v`.",
			replaceFeeInfo.Price, replaceFeeInfo.Currency,
		)
	}
	if feeBasePrice.GreaterThanOrEqual(FeePriceMaxEthereum) {
		return meta.NewMessageError(
			"System withdrawal cannot replace with this high fee price `%v %v`.",
			replaceFeeInfo.Price, replaceFeeInfo.Currency,
		)
	}

	client, err := newCoinClient(coin)
	if err != nil {
		return err
	}

	var (
		replaceCodes = make([]string, 0, len(txns))
		txnCodeMap   = make(map[string]models.SystemWithdrawalTxn, len(txns))
	)
	for _, txn := range txns {
		replaceCodes = append(replaceCodes, txn.RefCode)
		txnCodeMap[txn.RefCode] = txn
	}
	transfers, err := genTransferByCodes(coin, replaceCodes)
	if err != nil {
		return
	}

	var systemAddress models.SystemWithdrawalAddress
	err = database.GetDbSlave().
		First(&systemAddress, request.AddressID).
		Error
	if err != nil {
		err = utils.WrapError(err)
		return
	}
	srcAddressKey, err := systemAddress.Key.GetValue()
	if err != nil {
		return
	}

	txnSigner, err := getBalanceLikeTxnSigner(coin, client, replaceFeeInfo)
	if err != nil {
		return
	}
	if err = txnSigner.SetSrc(srcAddressKey, systemAddress.Address); err != nil {
		return
	}

	var (
		now             = time.Now()
		replaceTxns     = make([]models.SystemWithdrawalTxn, 0, len(txns))
		replaceFeePrice = replaceFeeInfo.GetBasePrice()
	)
	for _, transfer := range transfers {
		var (
			replaceTxn = txnCodeMap[transfer.RefCode]
			nonce      = blockchainmod.NewNumberNonce(replaceTxn.OutputIndex)
		)
		if replaceFeePrice.LessThanOrEqual(replaceTxn.FeePrice) {
			return meta.NewMessageError(
				"System withdrawal expects replacement fee price `%v` must be higher than the old fee price `%v`.",
				replaceFeePrice, replaceTxn.FeePrice,
			)
		}
		if err = txnSigner.SetDst(transfer.ToAddress, transfer.Amount); err != nil {
			return
		}
		if err = txnSigner.SetNonce(nonce); err != nil {
			return
		}
		if err = txnSigner.Sign(true); err != nil {
			return
		}

		feeInfo := txnSigner.GetFeeInfo()
		replaceTxn.ID = 0
		replaceTxn.Status = constants.SystemWithdrawalTxnStatusInit
		replaceTxn.Hash = models.NewString(txnSigner.GetHash())
		replaceTxn.FeePrice = feeInfo.GetBasePrice()
		replaceTxn.FeeMaxQuantity = feeInfo.LimitMaxQuantity
		replaceTxn.SignedBytes = models.NewString(string(txnSigner.GetRaw()))
		replaceTxn.CreateTime = now.Unix()
		replaceTxn.UpdateTime = now.Unix()
		replaceTxns = append(replaceTxns, replaceTxn)
	}

	return database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		for _, txn := range replaceTxns {
			if err = dbTxn.Save(&txn).Error; err != nil {
				return utils.WrapError(err)
			}
		}

		return nil
	})
}
