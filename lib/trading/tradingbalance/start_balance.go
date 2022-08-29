package tradingbalance

import (
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

func GenUserCoinStartBalanceMapMap(uids []meta.UID, fetchAll bool) map[meta.UID]CoinBalanceMap {
	query := database.GetDbSlave().
		Where(dbquery.Gt(models.UserBalanceColAmount, 0))
	if !fetchAll {
		query = query.Where(dbquery.In(models.CommonColUID, uids))
	}
	var balances []models.UserBalance
	comutils.PanicOnError(
		query.Find(&balances).Error,
	)

	userStartBalancesMap := make(map[meta.UID]CoinBalanceMap)
	for _, uid := range uids {
		userStartBalancesMap[uid] = make(CoinBalanceMap)
	}
	for _, balance := range balances {
		if coinBalanceMap, ok := userStartBalancesMap[balance.UID]; ok {
			coinBalanceMap[balance.Currency] = balance.Amount
		}
	}
	return userStartBalancesMap
}
