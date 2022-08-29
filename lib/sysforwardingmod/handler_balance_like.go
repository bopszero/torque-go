package sysforwardingmod

import (
	"sort"
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
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

type BalanceLikeForwardingHandler struct {
	baseForwardHandler

	networkCoin blockchainmod.Coin
}

func NewBalanceLikeForwardingHandler(
	currency meta.Currency, network meta.BlockchainNetwork,
	date string,
) (*BalanceLikeForwardingHandler, error) {
	var (
		coin        = blockchainmod.GetCoinF(currency, network)
		networkCoin = blockchainmod.GetCoinF(coin.GetNetworkCurrency(), network)
	)
	dateTime, err := comutils.TimeParse(constants.DateFormatISO, date)
	if err != nil {
		return nil, utils.WrapError(err)
	}
	client, err := newBlockchainClient(coin)
	if err != nil {
		return nil, err
	}

	handler := BalanceLikeForwardingHandler{
		baseForwardHandler: baseForwardHandler{
			coin:   coin,
			date:   dateTime,
			client: client,
		},
		networkCoin: networkCoin,
	}
	return &handler, nil
}

func (this *BalanceLikeForwardingHandler) getTxnSigner(
	order models.SystemForwardingOrder, feeInfo blockchainmod.FeeInfo,
) (_ blockchainmod.SingleTxnSigner, err error) {
	txnSigner, err := this.coin.NewTxnSignerSingle(this.client, &feeInfo)
	if err != nil {
		return
	}
	if err = txnSigner.SetMoveDst(order.Address); err != nil {
		return
	}

	return txnSigner, nil
}

func (this *BalanceLikeForwardingHandler) SignOrder(ctx comcontext.Context, orderID uint64) (
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
			var (
				now                = time.Now()
				validTxnAddressSet = make(comtypes.HashSet, len(txns))
				ignoredTxns        = make([]models.SystemForwardingOrderTxn, 0)
				addressToTxnsMap   = make(map[string][]models.SystemForwardingOrderTxn)
			)
			for i := range txns {
				txn := txns[i]

				if !validTxnAddressSet.Contains(txn.FromAddress) {
					balance, err := this.client.GetBalance(txn.FromAddress)
					if err != nil {
						return err
					}
					txn.SignBalance = balance
					txn.SignBalanceTime = now.Unix()

					if balance.GreaterThanOrEqual(forwardConfig.AmountMinThreshold) {
						validTxnAddressSet.Add(txn.FromAddress)
					}
				} else {
					firstAddrTxn := addressToTxnsMap[txn.FromAddress][0]
					txn.SignBalance = firstAddrTxn.SignBalance
					txn.SignBalanceTime = firstAddrTxn.SignBalanceTime
				}

				if !validTxnAddressSet.Contains(txn.FromAddress) {
					ignoredTxns = append(ignoredTxns, txn)
				}

				addressToTxnsMap[txn.FromAddress] = append(
					addressToTxnsMap[txn.FromAddress],
					txn,
				)
			}

			ignoreMsg := "balance not enough"
			if err = ignoreTxns(ctx, ignoredTxns, ignoreMsg); err != nil {
				return err
			}
			if len(validTxnAddressSet) == 0 {
				order.Status = constants.SystemForwardingOrderStatusFailed
				order.UpdateTime = now.Unix()
				if err := dbTxn.Save(&order).Error; err != nil {
					return utils.WrapError(err)
				}

				return nil
			}

			feeInfo, err := this.coin.GetDefaultFeeInfo()
			if err != nil {
				return
			}
			feeInfo.UsePriceHigh()

			isEnoughFee, err := this.ensureFeeEnoughForToken(
				ctx,
				order, validTxnAddressSet, addressToTxnsMap, feeInfo)
			if err != nil {
				return err
			}
			if !isEnoughFee {
				return nil
			}

			var addressInfoModels []models.SystemForwardingAddress
			err = dbTxn.
				Where(dbquery.In(models.SystemForwardingAddressColAddress, validTxnAddressSet.AsList())).
				Where(&models.SystemForwardingAddress{Network: this.coin.GetNetwork()}).
				Find(&addressInfoModels).
				Error
			if err != nil {
				return utils.WrapError(err)
			}
			if len(addressInfoModels) != len(validTxnAddressSet) {
				return utils.IssueErrorf(
					"system forwarding cannot get user address keys | currency=%v,order_id=%v,address=%v",
					order.Currency, order.ID, validTxnAddressSet.AsList(),
				)
			}

			var lastHash string
			for _, addrModel := range addressInfoModels {
				addressTxns, ok := addressToTxnsMap[addrModel.Address]
				if !ok {
					panic(utils.IssueErrorf(
						"system forwarding address model doesn't exists in txn addresses | currency=%v,address=%v",
						this.coin.GetCurrency(), addrModel.Address,
					))
				}

				txnSigner, err := this.getTxnSigner(order, feeInfo)
				if err != nil {
					return err
				}
				addrKey, err := addrModel.Key.GetValue()
				if err != nil {
					return err
				}
				if err = txnSigner.SetSrc(addrKey, ""); err != nil {
					return err
				}
				if err = txnSigner.Sign(false); err != nil {
					return err
				}

				txnHash := txnSigner.GetHash()
				for _, txn := range addressTxns {
					txn.Hash = models.NewString(txnHash)
					txn.SignedBytes = models.NewString(string(txnSigner.GetRaw()))
					txn.UpdateTime = now.Unix()
					if err = dbTxn.Save(&txn).Error; err != nil {
						return utils.WrapError(err)
					}
				}

				lastHash = txnHash
			}
			order.Status = constants.SystemForwardingOrderStatusSigned
			order.UpdateTime = now.Unix()
			if err := dbTxn.Save(&order).Error; err != nil {
				return utils.WrapError(err)
			}

			comlogging.GetLogger().
				WithType(constants.LogTypeSystemForwarding).
				WithContext(ctx).
				WithFields(logrus.Fields{
					"currency":  order.Currency,
					"order_id":  order.ID,
					"txn_count": len(addressInfoModels),
					"last_hash": lastHash,
				}).
				Infof(
					"signed %v txns for %v order %v",
					len(addressInfoModels), order.Currency, order.ID,
				)

			return nil
		})
	})
	if err != nil {
		return
	}

	return &order, nil
}

func (this *BalanceLikeForwardingHandler) ensureFeeEnoughForToken(
	ctx comcontext.Context,
	order models.SystemForwardingOrder, validAddresses comtypes.HashSet,
	addressToTxnsMap map[string][]models.SystemForwardingOrderTxn,
	feeInfo blockchainmod.FeeInfo,
) (_ bool, err error) {
	if order.Currency == this.networkCoin.GetCurrency() {
		return true, nil
	}

	feeAmount := feeInfo.GetBaseValue()
	insufficientAddresses, hasPendingFee, err := this.filterInsufficientAddresses(
		validAddresses, feeAmount,
		addressToTxnsMap)
	if err != nil {
		return
	}

	if len(insufficientAddresses) == 0 {
		return !hasPendingFee, nil
	} else if !hasPendingFee {
		err := this.signFeeTxns(ctx, order, insufficientAddresses, addressToTxnsMap, feeInfo)
		return false, err
	} else {
		err := this.pushFeeTxns(ctx, order, validAddresses, addressToTxnsMap)
		return false, err
	}
}

func (this *BalanceLikeForwardingHandler) filterInsufficientAddresses(
	addressSet comtypes.HashSet, amount decimal.Decimal,
	addressToTxnsMap map[string][]models.SystemForwardingOrderTxn,
) (insufficientAddresses []string, hasPendingFee bool, err error) {
	networkClient, err := newBlockchainClient(this.networkCoin)
	if err != nil {
		return
	}

	collectAddressFunc := func(address string) error {
		ethBalance, err := networkClient.GetBalance(address)
		if err == nil {
			if ethBalance.LessThan(amount) {
				insufficientAddresses = append(insufficientAddresses, address)
			}
		}
		return err
	}
	for addressObj := range addressSet {
		address := addressObj.(string)
		anyTxn := addressToTxnsMap[address][0]
		switch anyTxn.Status {

		case constants.SystemForwardingOrderTxnStatusInit:
			if err = collectAddressFunc(address); err != nil {
				return
			}
			break

		case constants.SystemForwardingOrderTxnStatusFeeComing:
			txn, getTxnErr := networkClient.GetTxn(anyTxn.FeeTxnHash.String)
			if getTxnErr != nil && !utils.IsOurError(getTxnErr, constants.ErrorCodeDataNotFound) {
				err = getTxnErr
				return
			}
			if txn == nil || txn.GetHash() == "" || txn.GetConfirmations() < 1 {
				insufficientAddresses = append(insufficientAddresses, address)
				hasPendingFee = true
				break
			}
			if err = collectAddressFunc(address); err != nil {
				return
			}
			break

		default:
			panic(utils.IssueErrorf("not supported"))
		}
	}

	return
}

func (this *BalanceLikeForwardingHandler) signFeeTxns(
	ctx comcontext.Context,
	order models.SystemForwardingOrder, addresses []string,
	addressToTxnsMap map[string][]models.SystemForwardingOrderTxn,
	feeInfo blockchainmod.FeeInfo,
) (err error) {
	networkClient, err := newBlockchainClient(this.networkCoin)
	if err != nil {
		return
	}
	feeTxnSigner, err := this.genFeeDistributionTxnSigner(order.Currency, networkClient, feeInfo)
	if err != nil {
		return
	}
	nonce, err := networkClient.GetNextNonce(feeTxnSigner.GetSrcAddress())
	if err != nil {
		return
	}

	feeAmount := feeInfo.GetBaseValue()
	return database.Atomic(ctx, database.AliasInternalMaster, func(dbTxn *gorm.DB) (err error) {
		var (
			now      = time.Now()
			lastHash string
		)
		for _, address := range addresses {
			if err = feeTxnSigner.SetDst(address, feeAmount); err != nil {
				return
			}
			if err = feeTxnSigner.SetNonce(nonce); err != nil {
				return
			}
			if err = feeTxnSigner.Sign(true); err != nil {
				return
			}

			var (
				txnHash        = feeTxnSigner.GetHash()
				txnSignedBytes = feeTxnSigner.GetRaw()
			)
			for _, txn := range addressToTxnsMap[address] {
				if txn.FeeTxnIndex, err = nonce.GetNumber(); err != nil {
					return
				}
				txn.Status = constants.SystemForwardingOrderTxnStatusFeeComing
				txn.FeeTxnHash = models.NewString(txnHash)
				txn.FeeTxnSignedBytes = models.NewBytes(txnSignedBytes)
				txn.UpdateTime = now.Unix()
				if err = dbTxn.Save(&txn).Error; err != nil {
					return utils.WrapError(err)
				}
			}

			lastHash = txnHash
			if nonce, err = nonce.Next(); err != nil {
				return
			}
		}

		comlogging.GetLogger().
			WithType(constants.LogTypeSystemForwarding).
			WithContext(ctx).
			WithFields(logrus.Fields{
				"currency":  order.Currency,
				"order_id":  order.ID,
				"txn_count": len(addresses),
				"last_hash": lastHash,
			}).
			Infof(
				"signed %v fee txns for %v order %v",
				len(addresses), order.Currency, order.ID,
			)

		return nil
	})
}

func (this *BalanceLikeForwardingHandler) pushFeeTxns(
	ctx comcontext.Context,
	order models.SystemForwardingOrder, addressSet comtypes.HashSet,
	addressToTxnsMap map[string][]models.SystemForwardingOrderTxn,
) (err error) {
	networkClient, err := newBlockchainClient(this.networkCoin)
	if err != nil {
		return
	}

	validTxns := make([]models.SystemForwardingOrderTxn, len(addressSet))
	for addressObj := range addressSet {
		txns := addressToTxnsMap[addressObj.(string)]
		validTxns = append(validTxns, txns...)
	}
	sort.Slice(
		validTxns,
		func(i, j int) bool {
			return validTxns[i].FeeTxnIndex < validTxns[j].FeeTxnIndex
		},
	)

	var (
		checkFeeTxnHashSet = make(comtypes.HashSet, len(validTxns))
		lastPushedHash     string
	)
	for _, txn := range validTxns {
		if txn.FeeTxnHash.String == "" || !checkFeeTxnHashSet.Add(txn.FeeTxnHash.String) {
			continue
		}

		feeTxn, err := networkClient.GetTxn(txn.FeeTxnHash.String)
		if err != nil && !utils.IsOurError(err, constants.ErrorCodeDataNotFound) {
			return err
		}
		if feeTxn != nil && feeTxn.GetHash() != "" {
			continue
		}

		feeTxnBytes := []byte(txn.FeeTxnSignedBytes.String)
		if err := networkClient.PushTxnRaw(feeTxnBytes); err != nil {
			return utils.IssueErrorf(
				"system forwarding push ETH fee txn failed | hash=%v",
				txn.FeeTxnHash.String)
		}
		lastPushedHash = txn.FeeTxnHash.String
	}

	if lastPushedHash != "" {
		comlogging.GetLogger().
			WithType(constants.LogTypeSystemForwarding).
			WithContext(ctx).
			WithFields(logrus.Fields{
				"currency":  order.Currency,
				"order_id":  order.ID,
				"last_hash": lastPushedHash,
			}).
			Infof(
				"pushed %v fee txns for %v order %v",
				len(checkFeeTxnHashSet), order.Currency, order.ID,
			)
	}

	return nil
}

func (this *BalanceLikeForwardingHandler) genFeeDistributionTxnSigner(
	currency meta.Currency,
	client blockchainmod.Client, feeInfo blockchainmod.FeeInfo,
) (_ blockchainmod.SingleTxnSigner, err error) {
	forwardConfig := this.getConfig()
	if forwardConfig.FeeKey == "" {
		err = utils.IssueErrorf(
			"system forwarding need Fee Key to distribute fee to forward `%v`",
			currency)
		return
	}

	txnSigner, err := this.networkCoin.NewTxnSignerSingle(client, &feeInfo)
	if err != nil {
		return
	}
	if err = txnSigner.SetSrc(forwardConfig.FeeKey, forwardConfig.Address); err != nil {
		return
	}

	return txnSigner, nil
}

func (this *BalanceLikeForwardingHandler) ExecuteOrder(ctx comcontext.Context, orderID uint64) (
	_ *models.SystemForwardingOrder, err error,
) {
	var (
		order    models.SystemForwardingOrder
		lastHash string
	)
	err = database.Atomic(ctx, database.AliasInternalMaster, func(dbTxn *gorm.DB) error {
		return database.Atomic(ctx, database.AliasMainMaster, func(mainDbTxn *gorm.DB) (err error) {
			err = dbquery.SelectForUpdate(dbTxn).
				Where(dbquery.In(
					models.ForwardingOrderColStatus,
					[]meta.SystemForwardingOrderStatus{
						constants.SystemForwardingOrderStatusSigned,
						constants.SystemForwardingOrderStatusForwarding,
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
			txns, err := this.getOrderTxns(dbTxn, order)
			if err != nil {
				return err
			}

			var (
				updatedTxnMap = make(map[uint64]models.SystemForwardingOrderTxn, len(txns))
				doneCount     = 0
			)

			logger := comlogging.GetLogger()
			for _, txn := range txns {
				isDone, err := this.executeTxn(order, txn, updatedTxnMap)
				if err != nil {
					logger.
						WithType(constants.LogTypeSystemForwarding).
						WithContext(ctx).
						WithError(err).
						WithFields(logrus.Fields{
							"txn_id": txn.ID,
							"coin":   this.coin.GetIndexCode(),
						}).
						Errorf("execute %v txn error | err=%s", this.coin.GetIndexCode(), err.Error())
					break
				}
				if isDone {
					doneCount++
					lastHash = txn.Hash.String
				}
			}

			if len(updatedTxnMap) > 0 {
				now := time.Now()
				for _, txn := range updatedTxnMap {
					txn.UpdateTime = now.Unix()
					if err = dbTxn.Save(&txn).Error; err != nil {
						return utils.WrapError(err)
					}
					err = updateTradingDeposits(ctx, []uint64{txn.DepositID}, txn.Status)
					if err != nil {
						return err
					}
				}
			}

			if doneCount < len(txns) {
				order.Status = constants.SystemForwardingOrderStatusForwarding
				order.UpdateTime = time.Now().Unix()
				if err = dbTxn.Save(&order).Error; err != nil {
					return utils.WrapError(err)
				}
			} else {
				order.Status = constants.SystemForwardingOrderStatusCompleted
				order.UpdateTime = time.Now().Unix()
				if err = dbTxn.Save(&order).Error; err != nil {
					return utils.WrapError(err)
				}

				logger.
					WithType(constants.LogTypeSystemForwarding).
					WithContext(ctx).
					WithFields(logrus.Fields{
						"currency":  order.Currency,
						"order_id":  order.ID,
						"last_hash": lastHash,
					}).
					Infof("completed the %v order %v", order.Currency, order.ID)
			}

			return nil
		})
	})
	if err != nil {
		return
	}

	return &order, nil
}

func (this *BalanceLikeForwardingHandler) executeTxn(
	order models.SystemForwardingOrder, txn models.SystemForwardingOrderTxn,
	updatedTxnMap map[uint64]models.SystemForwardingOrderTxn,
) (bool, error) {
	switch txn.Status {
	case constants.SystemForwardingOrderTxnStatusInit,
		constants.SystemForwardingOrderTxnStatusFeeComing,
		constants.SystemForwardingOrderTxnStatusPushed:

		networkTxn, err := this.client.GetTxn(txn.Hash.String)
		if err != nil && !utils.IsOurError(err, constants.ErrorCodeDataNotFound) {
			return false, err
		}

		if networkTxn != nil && networkTxn.GetHash() != "" {
			if networkTxn.GetConfirmations() == 0 {
				return false, nil
			}

			txnAmount, err := networkTxn.GetAmount()
			if err != nil {
				return false, err
			}
			txn.Amount = txnAmount
			txn.Fee = networkTxn.GetFee().Value

			if networkTxn.GetLocalStatus() == constants.BlockchainTxnStatusFailed {
				txn.Status = constants.SystemForwardingOrderTxnStatusFailed
			} else {
				txn.Status = constants.SystemForwardingOrderTxnStatusSuccess
			}
			updatedTxnMap[txn.ID] = txn

			return true, nil
		}

		txnBytes := []byte(txn.SignedBytes.String)
		if err := this.client.PushTxnRaw(txnBytes); err != nil {
			return false, err
		}

		updatedTxnMap[txn.ID] = txn
		return false, nil

	case constants.SystemForwardingOrderTxnStatusFailed,
		constants.SystemForwardingOrderTxnStatusSuccess:
		return true, nil

	default:
		panic(utils.IssueErrorf("not supported"))
	}
}
