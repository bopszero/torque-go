package sysforwardingmod

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/go-common/comutils"
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

type UtxoLikeForwardingHandler struct {
	baseForwardHandler
}

func NewUtxoLikeForwardingHandler(
	currency meta.Currency, network meta.BlockchainNetwork,
	date string,
) (*UtxoLikeForwardingHandler, error) {
	coin := blockchainmod.GetCoinF(currency, network)
	dateTime, err := comutils.TimeParse(constants.DateFormatISO, date)
	if err != nil {
		return nil, utils.WrapError(err)
	}
	client, err := newBlockchainClient(coin)
	if err != nil {
		return nil, err
	}

	handler := UtxoLikeForwardingHandler{
		baseForwardHandler{
			coin:   coin,
			date:   dateTime,
			client: client,
		},
	}
	return &handler, nil
}

func (this *UtxoLikeForwardingHandler) newTxnSigner(order models.SystemForwardingOrder) (
	_ blockchainmod.BulkTxnSigner, err error,
) {
	feeInfo, err := this.coin.GetDefaultFeeInfo()
	if err != nil {
		return
	}
	feeInfo.UsePriceHigh()

	txnSigner, err := this.coin.NewTxnSignerBulk(this.client, &feeInfo)
	if err != nil {
		return
	}
	if err = txnSigner.SetMoveDst(order.Address); err != nil {
		return
	}

	return txnSigner, nil
}

func (this *UtxoLikeForwardingHandler) SignOrder(ctx comcontext.Context, orderID uint64) (
	_ *models.SystemForwardingOrder, err error,
) {
	var (
		forwardConfig = this.getConfig()
		order         models.SystemForwardingOrder
	)
	err = database.Atomic(ctx, database.AliasInternalMaster, func(dbTxn *gorm.DB) error {
		return database.Atomic(ctx, database.AliasMainMaster, func(mainDbTxn *gorm.DB) (err error) {
			err = dbquery.SelectForUpdate(dbTxn).
				Where(dbquery.In(
					models.ForwardingOrderColStatus,
					[]meta.SystemForwardingOrderStatus{
						constants.SystemForwardingOrderStatusInit,
						constants.SystemForwardingOrderStatusSigned,
					},
				)).
				Where(&models.SystemForwardingOrder{
					ID:       orderID,
					Currency: this.coin.GetCurrency(),
				}).
				First(&order).
				Error
			if err != nil {
				return utils.WrapError(err)
			}
			if order.Status == constants.SystemForwardingOrderStatusSigned {
				return nil
			}

			txns, err := this.getOrderTxns(dbTxn, order)
			if err != nil {
				return err
			}
			txnAddressSet := make(comtypes.HashSet, len(txns))
			for _, txn := range txns {
				txnAddressSet.Add(txn.FromAddress)
			}

			var addressInfoModels []models.SystemForwardingAddress
			err = dbTxn.
				Where(dbquery.In(models.SystemForwardingAddressColAddress, txnAddressSet.AsList())).
				Where(&models.SystemForwardingAddress{Network: this.coin.GetNetwork()}).
				Find(&addressInfoModels).
				Error
			if err != nil {
				return utils.WrapError(err)
			}
			if len(addressInfoModels) != len(txnAddressSet) {
				return utils.IssueErrorf(
					"system forwarding cannot get user address keys | currency=%v,order_id=%v,address=%v",
					order.Currency, order.ID, txnAddressSet.AsList(),
				)
			}

			txnSigner, err := this.newTxnSigner(order)
			if err != nil {
				return err
			}

			var (
				addressBalanceMap = make(map[string]decimal.Decimal, len(txnAddressSet))
				signedAddressSet  = make(comtypes.HashSet, len(addressInfoModels))
			)
			for _, addrModel := range addressInfoModels {
				balance, err := this.client.GetBalance(addrModel.Address)
				if err != nil {
					return err
				}
				addressBalanceMap[addrModel.Address] = balance

				if balance.LessThan(forwardConfig.AmountMinThreshold) {
					continue
				}

				addrKey, err := addrModel.Key.GetValue()
				if err != nil {
					return err
				}
				if err = txnSigner.AddSrc(addrKey, addrModel.Address); err != nil {
					return err
				}
				signedAddressSet.Add(addrModel.Address)
			}

			var (
				now         = time.Now()
				signedTxns  = make([]models.SystemForwardingOrderTxn, 0, len(txns))
				ignoredTxns = make([]models.SystemForwardingOrderTxn, 0)
			)
			for _, txn := range txns {
				txn.SignBalance = addressBalanceMap[txn.FromAddress]
				txn.SignBalanceTime = now.Unix()

				if signedAddressSet.Contains(txn.FromAddress) {
					signedTxns = append(signedTxns, txn)
				} else {
					ignoredTxns = append(ignoredTxns, txn)
				}
			}

			ignoreMsg := "balance not enough"
			if err = ignoreTxns(ctx, ignoredTxns, ignoreMsg); err != nil {
				return err
			}
			if len(signedAddressSet) == 0 {
				order.Status = constants.SystemForwardingOrderStatusFailed
				order.UpdateTime = now.Unix()
				if err := dbTxn.Save(&order).Error; err != nil {
					return utils.WrapError(err)
				}

				return nil
			}

			if err = txnSigner.Sign(false); err != nil {
				return err
			}

			txnHash := txnSigner.GetHash()
			order.CombinedTxnHash = models.NewString(txnHash)
			order.CombinedSignedBytes = models.NewString(string(txnSigner.GetRaw()))
			order.Status = constants.SystemForwardingOrderStatusSigned
			order.UpdateTime = now.Unix()
			if err := dbTxn.Save(&order).Error; err != nil {
				return utils.WrapError(err)
			}

			for _, txn := range signedTxns {
				txn.Hash = order.CombinedTxnHash
				txn.UpdateTime = now.Unix()
				if err = dbTxn.Save(&txn).Error; err != nil {
					return utils.WrapError(err)
				}
			}

			comlogging.GetLogger().
				WithType(constants.LogTypeSystemForwarding).
				WithContext(ctx).
				WithFields(logrus.Fields{
					"currency": order.Currency,
					"order_id": order.ID,
					"hash":     order.CombinedTxnHash.String,
				}).
				Infof("signed the %v order %v", order.Currency, order.ID)

			return nil
		})
	})
	if err != nil {
		return
	}

	return &order, nil
}

func (this *UtxoLikeForwardingHandler) ExecuteOrder(ctx comcontext.Context, orderID uint64) (
	_ *models.SystemForwardingOrder, err error,
) {
	var order models.SystemForwardingOrder
	err = database.Atomic(ctx, database.AliasInternalMaster, func(dbTxn *gorm.DB) error {
		return database.Atomic(ctx, database.AliasMainMaster, func(mainDbTxn *gorm.DB) (err error) {
			err = dbquery.SelectForUpdate(dbTxn).
				Where(&models.SystemForwardingOrder{
					ID:       orderID,
					Currency: this.coin.GetCurrency(),
				}).
				First(&order).
				Error
			if err != nil {
				return utils.WrapError(err)
			}

			switch order.Status {
			case constants.SystemForwardingOrderStatusSigned, constants.SystemForwardingOrderStatusForwarding:
				break
			default:
				return utils.IssueErrorf(
					"system forwarding cannot sign an invalid status order | order_id=%v,status=%v",
					order.ID, order.Status,
				)
			}

			txn, err := this.client.GetTxn(order.CombinedTxnHash.String)
			if err != nil && !utils.IsOurError(err, constants.ErrorCodeDataNotFound) {
				return err
			}
			if txn == nil || txn.GetHash() == "" {
				return this.forwardOrder(ctx, order)
			}
			if txn.GetConfirmations() < 1 {
				return nil
			}

			return this.completeOrder(ctx, order, txn)
		})
	})
	if err != nil {
		return
	}

	return &order, nil
}

func (this *UtxoLikeForwardingHandler) forwardOrder(ctx comcontext.Context, order models.SystemForwardingOrder) error {
	return database.Atomic(ctx, database.AliasInternalMaster, func(dbTxn *gorm.DB) error {
		return database.Atomic(ctx, database.AliasMainMaster, func(mainDbTxn *gorm.DB) (err error) {
			txns, err := this.getOrderTxns(dbTxn, order)
			if err != nil {
				return err
			}

			var (
				txnIDs     = make([]uint64, len(txns))
				depositIDs = make([]uint64, len(txns))
			)
			for i, txn := range txns {
				txnIDs[i] = txn.ID
				depositIDs[i] = txn.DepositID
			}

			txnBytes := []byte(order.CombinedSignedBytes.String)
			if err := this.client.PushTxnRaw(txnBytes); err != nil {
				return err
			}

			comlogging.GetLogger().
				WithType(constants.LogTypeSystemForwarding).
				WithContext(ctx).
				WithFields(logrus.Fields{
					"currency": order.Currency,
					"order_id": order.ID,
					"hash":     order.CombinedTxnHash.String,
				}).
				Infof("pushed the %v order %v", order.Currency, order.ID)

			now := time.Now()

			order.Status = constants.SystemForwardingOrderStatusForwarding
			order.UpdateTime = now.Unix()
			if err := dbTxn.Save(&order).Error; err != nil {
				return utils.WrapError(err)
			}

			err = dbTxn.
				Model(&models.SystemForwardingOrderTxn{}).
				Where(dbquery.In(models.CommonColID, txnIDs)).
				Updates(&models.SystemForwardingOrderTxn{
					Status:     constants.SystemForwardingOrderTxnStatusPushed,
					UpdateTime: now.Unix(),
				}).
				Error
			if err != nil {
				return utils.WrapError(err)
			}

			return updateTradingDeposits(
				ctx,
				depositIDs, constants.SystemForwardingOrderTxnStatusPushed,
			)
		})
	})
}

func (this *UtxoLikeForwardingHandler) completeOrder(
	ctx comcontext.Context, order models.SystemForwardingOrder, txn blockchainmod.Transaction,
) (err error) {
	var (
		txnAmount decimal.Decimal
		txnFee    meta.CurrencyAmount
	)
	txnOutputs, err := txn.GetOutputs()
	if err != nil {
		return err
	}
	txnAmount = txnOutputs[0].GetAmount()
	txnFee, err = currencymod.ConvertAmount(txn.GetFee(), this.coin.GetCurrency())
	if err != nil {
		return err
	}

	return database.Atomic(ctx, database.AliasInternalMaster, func(dbTxn *gorm.DB) error {
		return database.Atomic(ctx, database.AliasMainMaster, func(mainDbTxn *gorm.DB) (err error) {
			now := time.Now()

			txns, err := this.getOrderTxns(dbTxn, order)
			if err != nil {
				return err
			}

			var (
				txnIDs     = make([]uint64, len(txns))
				depositIDs = make([]uint64, len(txns))
			)
			for i, txn := range txns {
				txnIDs[i] = txn.ID
				depositIDs[i] = txn.DepositID
			}

			order.Status = constants.SystemForwardingOrderStatusCompleted
			order.UpdateTime = now.Unix()
			if err := dbTxn.Save(&order).Error; err != nil {
				return utils.WrapError(err)
			}

			err = dbTxn.
				Model(&models.SystemForwardingOrderTxn{}).
				Where(dbquery.In(models.CommonColID, txnIDs)).
				Updates(&models.SystemForwardingOrderTxn{
					Amount:     txnAmount,
					Fee:        txnFee.Value,
					Status:     constants.SystemForwardingOrderTxnStatusSuccess,
					UpdateTime: now.Unix(),
				}).
				Error
			if err != nil {
				return utils.WrapError(err)
			}

			comlogging.GetLogger().
				WithType(constants.LogTypeSystemForwarding).
				WithContext(ctx).
				WithFields(logrus.Fields{
					"currency": order.Currency,
					"order_id": order.ID,
					"hash":     order.CombinedTxnHash.String,
				}).
				Infof("completed the %v order %v", order.Currency, order.ID)

			return updateTradingDeposits(
				ctx,
				depositIDs, constants.SystemForwardingOrderTxnStatusSuccess,
			)
		})
	})
}
