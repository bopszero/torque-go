package crontasks

import (
	"fmt"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/ordermod"
)

const OrderAutofixLimit = 250

func OrderAutofixAllPending(workerCount int) {
	defer apiutils.CronRunWithRecovery(
		"OrderAutofixAllPending",
		meta.O{"worker_count": workerCount},
	)

	var wg sync.WaitGroup
	poolFunc := func(seed interface{}) {
		defer wg.Done()
		var (
			ctx    = comcontext.NewContext()
			logger = comlogging.GetLogger()
			orders []models.Order
		)
		comutils.PanicOnError(
			database.GetDbF(database.AliasWalletSlave).
				Where(dbquery.Lte(models.OrderColRetryTime, time.Now().Unix())).
				Where(dbquery.In(models.OrderColStatus, constants.OrderPendingStatuses)).
				Where(fmt.Sprintf("MOD(%s, %d) = %d", models.OrderColID, workerCount, seed)).
				Order(dbquery.OrderAsc(models.OrderColRetryTime)).
				Limit(OrderAutofixLimit).
				Find(&orders).
				Error,
		)
		for _, order := range orders {
			if err := orderAutofix(ctx, order.ID); err != nil {
				logger.
					WithType(constants.LogTypeOrder).
					WithError(err).
					WithField("order_id", order.ID).
					Warnf("autofix an order failed | err=%v", err.Error())
			}
		}
	}
	pool, err := ants.NewPoolWithFunc(workerCount, poolFunc)
	comutils.PanicOnError(err)
	defer pool.Release()

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		comutils.PanicOnError(
			pool.Invoke(i),
		)
	}
	wg.Wait()
}

func orderAutofix(ctx comcontext.Context, orderID uint64) error {
	defer apiutils.CronRunWithRecovery(
		"orderAutofix",
		meta.O{"order_id": orderID},
	)

	order, err := ordermod.ExecuteOrder(ctx, orderID)
	if err != nil {
		return err
	}

	switch order.Status {
	case constants.OrderStatusCompleted:
		if err := ordermod.PushOrderCompletedNotificationAsync(order); err != nil {
			comlogging.GetLogger().
				WithError(err).
				WithField("order_id", order.ID).
				Warnf("push order completed notification failed | err=%v", err.Error())
		}
	case constants.OrderStatusRefunded, constants.OrderStatusFailed:
		if err := ordermod.PushOrderFailedNotificationAsync(order); err != nil {
			comlogging.GetLogger().
				WithError(err).
				WithField("order_id", order.ID).
				Warnf("push order failed notification failed | err=%v", err.Error())
		}
	default:
		break
	}
	return nil
}
