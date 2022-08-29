package orderchannel

import (
	"fmt"
	"time"

	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/usermod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/ordermod"
	"gorm.io/gorm"
)

const (
	BlockchainPendingDuration = 15 * time.Minute
)

func executeBlockchainTxn(
	ctx comcontext.Context, order *models.Order, toAddress string, feeInfo *blockchainmod.FeeInfo,
) (cryptoTxn models.CryptoTxn, code meta.OrderStepResultCode, err error) {
	isRetry := order.DstChannelRef != ""
	if isRetry {
		code = ordermod.OrderStepResultCodeRetry
	} else {
		code = ordermod.OrderStepResultCodeFail
	}

	coin, err := blockchainmod.GetCoinNative(order.Currency)
	if err != nil {
		return
	}
	account, err := coin.NewAccountSystem(ctx, order.UID)
	if err != nil {
		return
	}

	var txnBytes []byte
	if order.DstChannelRef != "" {
		cryptoTxn, err = getLocalBlockchainTxn(ctx, order)
		if err != nil {
			return
		}
		var txn blockchainmod.Transaction
		txn, err = account.GetTxn(cryptoTxn.Hash)
		if err != nil && !utils.IsOurError(err, constants.ErrorCodeDataNotFound) {
			return
		}
		if txn != nil && txn.GetHash() != "" {
			return cryptoTxn, ordermod.OrderStepResultCodeSuccess, nil
		}

		txnBytes = []byte(cryptoTxn.SignedBytes.String)
	} else {
		code = ordermod.OrderStepResultCodeRetry
		var finalFeeInfo blockchainmod.FeeInfo
		if feeInfo.IsEmpty() {
			finalFeeInfo, err = account.GetFeeInfoToAddress(toAddress)
			if err != nil {
				return
			}
			feeInfo = &finalFeeInfo
		} else {
			finalFeeInfo = *feeInfo
		}

		var txnSigner blockchainmod.SingleTxnSigner
		txnSigner, err = account.GenerateSignedTxn(
			toAddress,
			order.AmountSubTotal, &finalFeeInfo)
		if err != nil {
			if utils.IsOurError(err, constants.ErrorCodeBalanceNotEnough) {
				code = ordermod.OrderStepResultCodeFail
			}
			return
		}
		txnBytes = txnSigner.GetRaw()

		now := time.Now()
		cryptoTxn = models.CryptoTxn{
			UID:         order.UID,
			Network:     coin.GetNetwork(),
			Hash:        txnSigner.GetHash(),
			SignedBytes: models.NewBytes(txnBytes),

			CreateTime: now.Unix(),
			UpdateTime: now.Unix(),
		}
		err = database.Atomic(ctx, database.AliasWalletMaster, func(dbTxn *gorm.DB) error {
			if err := dbTxn.Save(&cryptoTxn).Error; err != nil {
				return utils.WrapError(err)
			}
			return nil
		})
		if err != nil {
			if database.IsDuplicateEntryError(err) {
				code = ordermod.OrderStepResultCodeFail
				err = utils.WrapError(constants.ErrorOrderDuplicated)
			}
			return
		}
		order.DstChannelRef = cryptoTxn.Hash
	}

	if !config.Debug || config.BlockchainUseTestnet {
		if err = account.PushTxn(txnBytes); err != nil {
			if utils.IsRequestTimeoutError(err) {
				code = ordermod.OrderStepResultCodeRetry
			} else {
				code = ordermod.OrderStepResultCodeFail
			}
			return
		}
	}

	return cryptoTxn, ordermod.OrderStepResultCodeSuccess, nil
}

func getLocalBlockchainTxn(ctx comcontext.Context, order *models.Order) (
	cryptoTxn models.CryptoTxn, err error,
) {
	if order.DstChannelRef == "" {
		err = utils.WrapError(constants.ErrorDataNotFound)
		return
	}

	coin := blockchainmod.GetCoinNativeF(order.Currency)
	err = database.Atomic(ctx, database.AliasWalletMaster, func(dbTxn *gorm.DB) error {
		err := dbTxn.
			First(
				&cryptoTxn,
				&models.CryptoTxn{
					UID:     order.UID,
					Network: coin.GetNetwork(),
					Hash:    order.DstChannelRef,
				},
			).
			Error
		if err != nil {
			return utils.WrapError(err)
		}

		return nil
	})
	return
}

func removeLocalBlockchainTxn(ctx comcontext.Context, order *models.Order) error {
	if order.DstChannelRef == "" {
		return nil
	}

	coin := blockchainmod.GetCoinNativeF(order.Currency)
	cryptoTxn := models.CryptoTxn{
		UID:     order.UID,
		Network: coin.GetNetwork(),
		Hash:    order.DstChannelRef,
	}
	return database.Atomic(ctx, database.AliasWalletMaster, func(dbTxn *gorm.DB) error {
		if err := dbTxn.Delete(models.CryptoTxn{}, &cryptoTxn).Error; err != nil {
			return utils.WrapError(err)
		}
		return nil
	})
}

func watchBlockchainTxnConfirmations(ctx comcontext.Context, order *models.Order) (
	meta.OrderStepResultCode, error,
) {
	coin := blockchainmod.GetCoinNativeF(order.Currency)
	blockchainAccount, err := coin.NewAccountSystem(ctx, order.UID)
	if err != nil {
		return ordermod.OrderStepResultCodeRetry, err
	}

	txn, err := blockchainAccount.GetTxn(order.DstChannelRef)
	if err != nil {
		return ordermod.OrderStepResultCodeRetry, err
	}

	if txn.GetHash() == "" {
		return ordermod.OrderStepResultCodeRetry, nil
	}

	switch txn.GetLocalStatus() {
	case constants.BlockchainTxnStatusSucceeded:
		return ordermod.OrderStepResultCodeSuccess, nil
	case constants.BlockchainTxnStatusPending:
		return ordermod.OrderStepResultCodeRetry, nil
	case constants.BlockchainTxnStatusFailed:
		return ordermod.OrderStepResultCodeFail, utils.IssueErrorf("blockchain transaction has failed | hash=%v", txn.GetHash())
	default:
		panic(fmt.Errorf("blockchain transcation has an unknown local status `%v`", txn.GetLocalStatus()))
	}
}

func validateBlockchainOrder(ctx comcontext.Context, channel ordermod.Channel, order models.Order) error {
	coin, err := blockchainmod.GetCoinNative(order.Currency)
	if err != nil {
		return err
	}

	minTxnAmount := coin.GetMinTxnAmount()
	if order.AmountSubTotal.LessThan(minTxnAmount) {
		return constants.ErrorAmountTooLowWithValue.WithData(meta.O{
			"threshold": minTxnAmount,
			"currency":  order.Currency,
		})
	}

	infoModel := channel.GetInfoModel()
	if infoModel.BlockchainNetworkConfig.CurrencyAvailabilityMap != nil {
		if ok, exists := infoModel.BlockchainNetworkConfig.CurrencyAvailabilityMap[order.Currency]; exists && !ok {
			return constants.ErrorChannelNotAvailable
		}
	}
	if infoModel.BlockchainNetworkConfig.CurrencyAmountThresholdMap != nil {
		threshold, exists := infoModel.BlockchainNetworkConfig.CurrencyAmountThresholdMap[order.Currency]
		if exists {
			if order.AmountSubTotal.LessThan(threshold.Min) {
				return constants.ErrorAmountTooLowWithValue.WithData(meta.O{
					"threshold": threshold.Min,
					"currency":  order.Currency,
				})
			} else if !threshold.Max.IsZero() && order.AmountSubTotal.GreaterThan(threshold.Max) {
				return constants.ErrorAmountTooHighWithValue.WithData(meta.O{
					"threshold": threshold.Max,
					"currency":  order.Currency,
				})
			}
		}
	}

	return nil
}

func getUserByIdentity(userIdentity string) (user models.User, err error) {
	if user, err = usermod.GetUserByUsername(userIdentity); err == nil {
		return
	}
	if !utils.IsOurError(err, constants.ErrorCodeUserNotFound) {
		return
	}
	return usermod.GetUserByReferralCode(userIdentity)
}

func blockchainGetPendingOrders(
	ctx comcontext.Context, order *models.Order,
) (orders []models.Order, err error) {
	err = database.Atomic(ctx, database.AliasWalletMaster, func(dbTxn *gorm.DB) error {
		latestPendingTime := order.CreateTime - int64(BlockchainPendingDuration.Seconds())
		err := dbTxn.
			Where(dbquery.NotEqual(models.CommonColID, order.ID)).
			Where(&models.Order{
				Currency: order.Currency,
				UID:      order.UID,
			}).
			Where(dbquery.In(
				models.OrderColStatus,
				[]meta.OrderStatus{
					constants.OrderStatusHandleSrc,
					constants.OrderStatusHandleDst,
				},
			)).
			Where(dbquery.In(
				models.OrderColDstChannelType,
				constants.BlockchainChannelTypes,
			)).
			Where(dbquery.Gte(models.OrderColCreateTime, latestPendingTime)).
			Find(&orders).
			Error
		if err != nil {
			return utils.WrapError(err)
		}
		return nil
	})
	return
}
