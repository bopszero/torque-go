package crontasks

import (
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/sysforwardingmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func ForwardingTradingDeposits(currency meta.Currency, network meta.BlockchainNetwork, date string) {
	defer apiutils.CronRunWithRecovery(
		"ForwardingTradingDeposits",
		meta.O{
			"currency": currency,
			"network":  network,
		},
	)

	var (
		logger = comlogging.GetLogger()
		ctx    = comcontext.NewContext()
	)
	handler, err := sysforwardingmod.NewForwardingHandler(currency, network, date)
	if err != nil {
		logger.
			WithType(constants.LogTypeSystemForwarding).
			WithError(err).
			Errorf("create handler failed | currency=%v,err=%v", currency, err.Error())
		return
	}

	var toExecOrder models.SystemForwardingOrder
	err = database.Atomic(ctx, database.AliasInternalMaster, func(dbTxn *gorm.DB) error {
		err = dbTxn.
			Where(dbquery.In(
				models.ForwardingOrderColStatus,
				[]meta.SystemForwardingOrderStatus{
					constants.SystemForwardingOrderStatusSigned,
					constants.SystemForwardingOrderStatusForwarding,
				})).
			Where(&models.SystemForwardingOrder{
				Currency: currency,
			}).
			First(&toExecOrder).
			Error
		if database.IsDbError(err) {
			return utils.WrapError(err)
		}

		return nil
	})
	if err != nil {
		logger.
			WithType(constants.LogTypeSystemForwarding).
			WithError(err).
			Errorf("query error | err=%v", err.Error())
		return
	}
	if toExecOrder.ID > 0 {
		if _, err := handler.ExecuteOrder(ctx, toExecOrder.ID); err != nil {
			logger.
				WithType(constants.LogTypeSystemForwarding).
				WithError(err).
				Errorf("execute order failed | order_id=%v,err=%v", toExecOrder.ID, err.Error())
		}
		return
	}

	var toSignOrder models.SystemForwardingOrder
	err = database.Atomic(ctx, database.AliasInternalMaster, func(dbTxn *gorm.DB) error {
		err := dbTxn.
			Where(&models.SystemForwardingOrder{
				Currency: currency,
				Status:   constants.SystemForwardingOrderStatusInit,
			}).
			First(&toSignOrder).
			Error
		if database.IsDbError(err) {
			return utils.WrapError(err)
		}

		return nil
	})
	if err != nil {
		logger.
			WithType(constants.LogTypeSystemForwarding).
			WithError(err).
			Errorf("query error | err=%v", err.Error())
		return
	}
	if toSignOrder.ID > 0 {
		if _, err := handler.SignOrder(ctx, toSignOrder.ID); err != nil {
			logger.
				WithType(constants.LogTypeSystemForwarding).
				WithError(err).
				Errorf("sign order failed | order_id=%v,err=%v", toSignOrder.ID, err.Error())
		}
		return
	}

	if err := handler.GenerateOrder(ctx); err != nil {
		logger.
			WithType(constants.LogTypeSystemForwarding).
			WithError(err).
			Errorf("generate order failed | err=%v", err.Error())
	}
	return
}
