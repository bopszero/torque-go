package helpers

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/ordermod"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/ordermod/orderchannel"
)

func WalletCollectRewards(date string) {
	comutils.EchoWithTime("Wallet collect rewards for date %s started.", date)

	var (
		db  = database.GetDbMaster()
		err error
	)
	currencyInfoIdMap, err := currencymod.GetAllLegacyCurrencyInfoIdMapFast()
	comutils.PanicOnError(err)

	updateBalanceMap := func(balanceMap map[meta.UID]decimal.Decimal, rows *sql.Rows) {
		var (
			userID meta.UID
			amount decimal.Decimal
		)
		err := rows.Scan(&userID, &amount)
		comutils.PanicOnError(err)

		balance, ok := balanceMap[userID]
		if !ok {
			balance = decimal.NewFromInt(0)
		}

		balanceMap[userID] = balance.Add(amount)
	}

	comutils.EchoWithTime("Wallet collect rewards for date %s collecting records...", date)

	affComRows, err := db.
		Table(models.AffiliateCommissionTableName).
		Select(fmt.Sprintf("`%s`, SUM(`%s`) AS `amount`", models.AffiliateCommissionColUID, models.AffiliateCommissionColAmount)).
		Where(&models.AffiliateCommission{Date: date, IsDeleted: models.NewBool(false)}).
		Where(dbquery.Gt(models.AffiliateCommissionColAmount, 0)).
		Group(models.AffiliateCommissionColUID).
		Rows()
	comutils.PanicOnError(err)

	userAffComBalanceMap := map[meta.UID]decimal.Decimal{}
	for affComRows.Next() {
		updateBalanceMap(userAffComBalanceMap, affComRows)
	}

	dailyProfitRows, err := db.
		Table(models.DailyProfitTableName).
		Select(fmt.Sprintf(
			"`%s`, `%s`, SUM(`%s`) AS `amount`",
			models.DailyProfitColCoinID,
			models.DailyProfitColUID,
			models.DailyProfitColAmount,
		)).
		Where(&models.DailyProfit{Date: date, IsDeleted: models.NewBool(false)}).
		Where(dbquery.Gt(models.DailyProfitColAmount, 0)).
		Group(fmt.Sprintf("%s, %s", models.DailyProfitColUID, models.DailyProfitColCoinID)).
		Rows()
	comutils.PanicOnError(err)

	type ProfitValues map[meta.Currency]decimal.Decimal
	userDailyProfitValuesMap := map[meta.UID]ProfitValues{}
	for dailyProfitRows.Next() {
		var (
			coinID uint16
			uid    meta.UID
			amount decimal.Decimal
		)
		err := dailyProfitRows.Scan(&coinID, &uid, &amount)
		comutils.PanicOnError(err)

		currencyInfo, ok := currencyInfoIdMap[coinID]
		if !ok {
			continue
		}

		profitValues, ok := userDailyProfitValuesMap[uid]
		if !ok {
			profitValues = make(ProfitValues)
			userDailyProfitValuesMap[uid] = profitValues
		}

		balance, ok := profitValues[currencyInfo.Currency]
		if !ok {
			balance = decimal.Zero
		}

		profitValues[currencyInfo.Currency] = balance.Add(amount)
	}

	leaderRewardRows, err := db.
		Table(models.LeaderRewardTableName).
		Select(fmt.Sprintf("`%s`, SUM(`%s`) AS `amount`", models.LeaderRewardColUID, models.LeaderRewardColAmount)).
		Where(&models.LeaderReward{Date: date, IsDeleted: models.NewBool(false)}).
		Where(dbquery.Gt(models.LeaderRewardColAmount, 0)).
		Group(models.LeaderRewardColUID).
		Rows()
	comutils.PanicOnError(err)

	userLeaderRewardBalanceMap := map[meta.UID]decimal.Decimal{}
	for leaderRewardRows.Next() {
		updateBalanceMap(userLeaderRewardBalanceMap, leaderRewardRows)
	}

	newRewardMeta := func() *orderchannel.SrcTradingRewardMeta {
		return &orderchannel.SrcTradingRewardMeta{
			AffiliateCommissionAmount: decimal.Zero,
			LeaderCommissionAmount:    decimal.Zero,
		}
	}

	userTradingRewardMetaMap := make(map[meta.UID]*orderchannel.SrcTradingRewardMeta)
	for uid, affComAmount := range userAffComBalanceMap {
		var rewardMeta *orderchannel.SrcTradingRewardMeta

		rewardMeta, ok := userTradingRewardMetaMap[uid]
		if !ok {
			rewardMeta = newRewardMeta()
		}

		rewardMeta.AffiliateCommissionAmount = affComAmount
		userTradingRewardMetaMap[uid] = rewardMeta
	}

	for uid, leaderCommAmount := range userLeaderRewardBalanceMap {
		var rewardMeta *orderchannel.SrcTradingRewardMeta

		rewardMeta, ok := userTradingRewardMetaMap[uid]
		if !ok {
			rewardMeta = newRewardMeta()
		}

		rewardMeta.LeaderCommissionAmount = leaderCommAmount
		userTradingRewardMetaMap[uid] = rewardMeta
	}

	for uid, dailyProfitValues := range userDailyProfitValuesMap {
		var rewardMeta *orderchannel.SrcTradingRewardMeta

		rewardMeta, ok := userTradingRewardMetaMap[uid]
		if !ok {
			rewardMeta = newRewardMeta()
		}

		rewardMeta.DailyProfitAmounts = make(
			orderchannel.SrcTradingRewardDailyProfitAmounts,
			0, len(dailyProfitValues))
		for currency, amount := range dailyProfitValues {
			currencyAmount := meta.CurrencyAmount{
				Currency: currency,
				Value:    amount,
			}
			rewardMeta.DailyProfitAmounts = append(rewardMeta.DailyProfitAmounts, currencyAmount)
		}
		userTradingRewardMetaMap[uid] = rewardMeta
	}

	comutils.EchoWithTime("Wallet collect rewards for date %s executing orders...", date)

	genOrder := func(uid meta.UID, rewardMeta *orderchannel.SrcTradingRewardMeta) (order models.Order, err error) {
		var (
			totalDailyProfit = rewardMeta.DailyProfitAmounts.Total()
			totalAmount      = rewardMeta.AffiliateCommissionAmount.
				Add(rewardMeta.LeaderCommissionAmount).
				Add(totalDailyProfit)
		)

		order = ordermod.NewUserOrder(
			uid, constants.CurrencyTorque,
			constants.ChannelTypeSrcTradingReward, constants.ChannelTypeDstBalance)
		order.SrcChannelRef = date
		order.SrcChannelAmount = totalAmount
		order.DstChannelAmount = totalAmount
		order.AmountSubTotal = totalAmount
		order.AmountTotal = totalAmount

		err = ordermod.SetOrderChannelMetaData(&order, order.SrcChannelType, rewardMeta)
		if err != nil {
			return
		}

		dstMeta := orderchannel.DstBalanceMeta{
			InnerTxns: []orderchannel.DstBalanceMetaInnerTxn{
				{
					TxnType: constants.WalletBalanceTypeMetaCrAffiliateCommission.ID,
					Value:   rewardMeta.AffiliateCommissionAmount,
				},
				{
					TxnType: constants.WalletBalanceTypeMetaCrDailyProfit.ID,
					Value:   totalDailyProfit,
				},
				{
					TxnType: constants.WalletBalanceTypeMetaCrLeaderCommission.ID,
					Value:   rewardMeta.LeaderCommissionAmount,
				},
			},
		}
		if err = ordermod.SetOrderChannelMetaData(&order, order.DstChannelType, dstMeta); err != nil {
			return
		}

		now := time.Now()
		order.Status = constants.OrderStatusHandleSrc
		order.StepsData = models.OrderStepsData{
			History: []models.OrderStep{
				{
					Direction: constants.OrderStepDirectionForward,
					Code:      constants.OrderStepCodeOrderInit,
					Time:      now.Unix(),
				},
				{
					Direction: constants.OrderStepDirectionForward,
					Code:      constants.OrderStepCodeOrderStartSrc,
					Time:      now.Unix(),
				},
			},
		}

		return order, nil
	}

	var (
		userErrorMap = make(map[meta.UID]error)
		rewardOrders = make([]models.Order, 0, len(userTradingRewardMetaMap))
	)
	for uid, rewardMeta := range userTradingRewardMetaMap {
		order, err := genOrder(uid, rewardMeta)
		if err != nil {
			userErrorMap[uid] = err
			continue
		}

		rewardOrders = append(rewardOrders, order)
	}
	comutils.PanicOnError(
		database.GetDbF(database.AliasWalletMaster).
			CreateInBatches(rewardOrders, 1000).
			Error,
	)
	logger := comlogging.GetLogger()
	for uid, err := range userErrorMap {
		logger.Errorf("user collect payout reward failed | uid=%v,error=%v", uid, err.Error())
	}
	comutils.EchoWithTime(
		"Wallet collect rewards for date %s finished with %v orders.",
		date, len(rewardOrders),
	)
}
