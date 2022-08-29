package tradingbalance

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func GenUserCoinBalanceMapMap(uids []meta.UID, fetchAll bool) (_ map[meta.UID]CoinBalanceMap, err error) {
	query := database.GetDbSlave().
		Where(dbquery.Gt(models.UserBalanceColAmount, 0))
	if !fetchAll {
		query = query.Where(dbquery.In(models.CommonColUID, uids))
	}
	var balances []models.UserBalance
	if err = query.Find(&balances).Error; err != nil {
		err = utils.WrapError(err)
		return
	}

	balanceMapMap := make(map[meta.UID]CoinBalanceMap)
	for _, uid := range uids {
		balanceMapMap[uid] = make(CoinBalanceMap)
	}
	for _, balance := range balances {
		if balanceMap, ok := balanceMapMap[balance.UID]; ok {
			balanceMap[balance.Currency] = balance.Amount
		}
	}
	return balanceMapMap, nil
}

func GenUserCoinBalanceMapMapF(uids []meta.UID, fetchAll bool) map[meta.UID]CoinBalanceMap {
	balanceMapMap, err := GenUserCoinBalanceMapMap(uids, fetchAll)
	comutils.PanicOnError(err)
	return balanceMapMap
}

func GetUserCoinBalanceMap(uid meta.UID) CoinBalanceMap {
	userBalanceMapMap := GenUserCoinBalanceMapMapF([]meta.UID{uid}, false)
	return userBalanceMapMap[uid]
}

func CalcUserCoinBalanceUsdValue(uid meta.UID) decimal.Decimal {
	userBalanceMap := GetUserCoinBalanceMap(uid)
	return userBalanceMap.CalcValueUSD(nil)
}
