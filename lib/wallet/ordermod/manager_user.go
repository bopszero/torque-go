package ordermod

import (
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func GetUserOrder(uid meta.UID, orderID uint64) (order models.Order, err error) {
	order, err = GetOrderDB(database.GetDbF(database.AliasWalletSlave).DB).
		GetOrder(orderID)
	if err != nil {
		return
	}
	if order.UID != uid {
		err = utils.IssueErrorf(
			"cannot get order of another user | order_id=%v,uid=%v",
			orderID, uid,
		)
		return
	}

	return order, nil
}

func GetUserCurrencyOrder(currency meta.Currency, uid meta.UID, orderID uint64) (
	order models.Order, err error,
) {
	order, err = GetOrderDB(database.GetDbF(database.AliasWalletSlave).DB).
		GetCurrencyOrder(currency, orderID)
	if err != nil {
		return
	}
	if order.UID != uid {
		err = utils.IssueErrorf(
			"cannot get order of another user | order_id=%v,uid=%v",
			orderID, uid,
		)
		return
	}

	return order, nil
}
