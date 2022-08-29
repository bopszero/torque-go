package sysforwardingmod

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbfields"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

type baseForwardHandler struct {
	coin   blockchainmod.Coin
	date   time.Time
	client blockchainmod.Client
}

func (this *baseForwardHandler) getConfig() ForwardConfig {
	return getConfig(this.coin)
}

func (this *baseForwardHandler) getOrderTxns(dbTxn *gorm.DB, order models.SystemForwardingOrder) (
	txns []models.SystemForwardingOrderTxn, err error,
) {
	err = dbTxn.
		Where(dbquery.Gte(models.SystemForwardingOrderTxnColStatus, constants.SystemForwardingOrderTxnStatusInit)).
		Find(&txns, &models.SystemForwardingOrderTxn{OrderID: order.ID}).
		Error
	if err != nil {
		err = utils.WrapError(err)
	}

	return
}

func (this *baseForwardHandler) GenerateOrder(ctx comcontext.Context) error {
	forwardingConfig := this.getConfig()
	if forwardingConfig.Address == "" {
		return utils.IssueErrorf(
			"system forwarding needs address to forward balances | currency=%v",
			this.coin.GetCurrency(),
		)
	}

	return database.Atomic(ctx, database.AliasInternalMaster, func(dbTxn *gorm.DB) error {
		return database.Atomic(ctx, database.AliasMainMaster, func(mainDbTxn *gorm.DB) (err error) {
			var (
				fromTime = this.date.Unix()
				toTime   = this.date.Add(24 * time.Hour).Unix()
				deposits []models.Deposit
			)
			err = mainDbTxn.
				Where(dbquery.Gte(models.DepositColCloseTime, fromTime)).
				Where(dbquery.Lt(models.DepositColCloseTime, toTime)).
				Where(dbquery.In(
					models.DepositColForwardStatus,
					[]meta.SystemForwardingOrderTxnStatus{
						constants.SystemForwardingOrderTxnStatusIgnored,
						constants.SystemForwardingOrderTxnStatusFailed,
						constants.SystemForwardingOrderTxnStatusUnknown,
					},
				)).
				Where(&models.Deposit{
					IsReinvest: models.NewBool(false),
					Currency:   this.coin.GetCurrency(),
					Network:    this.coin.GetNetwork(),
					Status:     constants.DepositStatusApproved,
				}).
				Find(&deposits).
				Error
			if err != nil {
				return utils.WrapError(err)
			}

			validDeposits := this.filterValidDeposits(ctx, deposits)
			if len(validDeposits) == 0 {
				return nil
			}

			depositIDs := make([]uint64, len(validDeposits))
			for i, deposit := range validDeposits {
				depositIDs[i] = deposit.ID
			}
			err = dbquery.SelectForUpdate(mainDbTxn).
				Where(dbquery.In(models.DepositColID, depositIDs)).
				Find(&validDeposits).
				Error
			if err != nil {
				return utils.WrapError(err)
			}

			now := time.Now()
			order := models.SystemForwardingOrder{
				Date:       dbfields.NewDateFieldFromTime(this.date),
				Currency:   this.coin.GetCurrency(),
				Address:    forwardingConfig.Address,
				Status:     constants.SystemForwardingOrderStatusInit,
				CreateTime: now.Unix(),
				UpdateTime: now.Unix(),
			}
			if err = dbTxn.Save(&order).Error; err != nil {
				return utils.WrapError(err)
			}

			txns := make([]models.SystemForwardingOrderTxn, len(validDeposits))
			for i, deposit := range validDeposits {
				normalizedAddress, err := this.coin.NormalizeAddress(deposit.Address)
				comutils.PanicOnError(err)

				txns[i] = models.SystemForwardingOrderTxn{
					OrderID:   order.ID,
					DepositID: deposit.ID,

					Currency:    order.Currency,
					Status:      constants.SystemForwardingOrderTxnStatusInit,
					FromAddress: normalizedAddress,

					CreateTime: now.Unix(),
					UpdateTime: now.Unix(),
				}
			}
			if err = dbTxn.CreateInBatches(txns, 200).Error; err != nil {
				return utils.WrapError(err)
			}

			err = mainDbTxn.
				Model(&models.Deposit{}).
				Where(dbquery.In(models.DepositColID, depositIDs)).
				Updates(&models.Deposit{
					ForwardStatus: constants.SystemForwardingOrderTxnStatusInit,
				}).
				Error
			if err != nil {
				return utils.WrapError(err)
			}

			comlogging.GetLogger().
				WithType(constants.LogTypeSystemForwarding).
				WithContext(ctx).
				WithFields(logrus.Fields{
					"currency":  order.Currency,
					"order_id":  order.ID,
					"txn_count": len(validDeposits),
				}).
				Infof(
					"generated a new %v order %v with %v txns",
					order.Currency, order.ID, len(validDeposits),
				)

			return nil
		})
	})
}

func (this *baseForwardHandler) filterValidDeposits(
	ctx comcontext.Context, deposits []models.Deposit,
) []models.Deposit {
	if config.Debug {
		return deposits
	}

	var (
		logger                 = comlogging.GetLogger()
		addressTotalDepositMap = make(map[string]decimal.Decimal, len(deposits))
		notIgnoredAddressSet   = make(comtypes.HashSet)
	)
	for _, d := range deposits {
		addrAmount := addressTotalDepositMap[d.Address]
		addressTotalDepositMap[d.Address] = addrAmount.Add(d.Amount)

		if d.ForwardStatus != constants.SystemForwardingOrderTxnStatusIgnored {
			notIgnoredAddressSet.Add(d.Address)
		}
	}

	var (
		now              = time.Now()
		nextDateTime     = this.date.Add(24 * time.Hour)
		isBeforeTomorrow = now.Before(nextDateTime)

		forwardConfig   = this.getConfig()
		validAddressSet = make(comtypes.HashSet, len(addressTotalDepositMap))
	)
	for address, totalAmount := range addressTotalDepositMap {
		if _, err := this.coin.NormalizeAddress(address); err != nil {
			logger.
				WithType(constants.LogTypeSystemForwarding).
				WithContext(ctx).
				WithFields(logrus.Fields{
					"currency": this.coin.GetCurrency(),
					"address":  address,
				}).
				Warnf("invalid %v deposit address `%v`", this.coin.GetCurrency(), address)
			continue
		}

		if totalAmount.LessThan(forwardConfig.AmountMinThreshold) && isBeforeTomorrow {
			continue
		}

		validAddressSet.Add(address)
		if forwardConfig.TxnCountMaxThreshold > 0 &&
			len(validAddressSet) >= forwardConfig.TxnCountMaxThreshold {
			break
		}
	}

	validDeposits := make([]models.Deposit, 0, len(deposits))
	for _, d := range deposits {
		if d.ForwardStatus == constants.SystemForwardingOrderTxnStatusIgnored {
			if _, accept := notIgnoredAddressSet[d.Address]; !accept {
				continue
			}
		}

		if _, err := this.coin.NormalizeAddress(d.Address); err != nil {
			logger.
				WithType(constants.LogTypeSystemForwarding).
				WithContext(ctx).
				WithFields(logrus.Fields{
					"currency":   d.Currency,
					"network":    d.Network,
					"address":    d.Address,
					"deposit_id": d.ID,
				}).
				Errorf("A deposit has an invalid address")
			continue
		}

		if !validAddressSet.Contains(d.Address) {
			continue
		}

		validDeposits = append(validDeposits, d)
	}

	if len(validDeposits) < forwardConfig.TxnCountMinThreshold && isBeforeTomorrow {
		return nil
	}

	return validDeposits
}
