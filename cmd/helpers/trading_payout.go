package helpers

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func TradingPayoutGenUserStartCoinBalances(dateStr string) {
	dateTime, err := comutils.TimeParse(constants.DateFormatISO, dateStr)
	comutils.PanicOnError(err)

	db := database.GetDbSlave()

	weekDuration := 7 * 24 * time.Hour
	lowerBoundTime := dateTime.Add(-weekDuration)
	upperBoundTime := dateTime.Add(weekDuration * 4)

	type StartingBalanceMap map[string]*models.PayoutUserStartDateBalance

	genTxnKey := func(currency meta.Currency, uid meta.UID) string {
		return fmt.Sprintf("%v:%v", currency, uid)
	}
	getIdsFromRows := func(rows *sql.Rows) []uint64 {
		ids := make([]uint64, 0)
		for rows.Next() {
			var idx uint64
			comutils.PanicOnError(
				rows.Scan(&idx),
			)
			ids = append(ids, idx)
		}
		return ids
	}
	genBalanceTxnMapByIds := func(ids []uint64) StartingBalanceMap {
		var txns []models.UserBalanceTxn
		comutils.PanicOnError(
			db.
				Where(dbquery.In(models.UserBalanceTxnColID, ids)).
				Where(dbquery.Gt(models.UserBalanceTxnColBalance, 0)).
				Find(&txns).
				Error,
		)
		txnMap := make(StartingBalanceMap, len(txns))
		createTimeUnix := time.Now().Unix()
		for _, txn := range txns {
			key := genTxnKey(txn.Currency, txn.UserID)
			txnMap[key] = &models.PayoutUserStartDateBalance{
				Date:       dateStr,
				Currency:   txn.Currency,
				UID:        txn.UserID,
				Amount:     txn.Balance,
				TxnID:      models.NewUInt64(txn.ID),
				CreateTime: createTimeUnix,
			}
		}

		return txnMap
	}

	lowerStartingBalanceIdsRows, err := db.
		Model(&models.UserBalanceTxn{}).
		Where(dbquery.Between(models.UserBalanceTxnColCreateTime, lowerBoundTime.Unix(), dateTime.Unix()-1)).
		Select("MAX(id) AS max_id").
		Group(dbquery.JoinExpr(
			models.UserBalanceTxnColCurrency,
			models.UserBalanceTxnColUID)).
		Rows()
	comutils.PanicOnError(err)
	lowerStartingBalanceIds := getIdsFromRows(lowerStartingBalanceIdsRows)
	lowerStartingBalanceMap := genBalanceTxnMapByIds(lowerStartingBalanceIds)

	upperStartingBalanceIdsRows, err := db.
		Model(&models.UserBalanceTxn{}).
		Where(dbquery.Between(models.UserBalanceTxnColCreateTime, dateTime.Unix(), upperBoundTime.Unix())).
		Select("MIN(id) AS min_id").
		Group(dbquery.JoinExpr(
			models.UserBalanceTxnColCurrency,
			models.UserBalanceTxnColUID)).
		Rows()
	comutils.PanicOnError(err)
	upperStartingBalanceIds := getIdsFromRows(upperStartingBalanceIdsRows)
	upperStartingBalanceMap := genBalanceTxnMapByIds(upperStartingBalanceIds)

	var oldBalances []models.UserBalance
	comutils.PanicOnError(db.
		Where(dbquery.Lt(models.UserBalanceColUpdateTime, dateTime.Unix())).
		Where(dbquery.Gt(models.UserBalanceColAmount, 0)).
		Find(&oldBalances).
		Error,
	)
	oldStartingBalanceMap := make(StartingBalanceMap, len(oldBalances))
	oldStartingBalanceCreateTimeUnix := time.Now().Unix()
	for _, balance := range oldBalances {
		key := genTxnKey(balance.Currency, balance.UID)
		oldStartingBalanceMap[key] = &models.PayoutUserStartDateBalance{
			Date:       dateStr,
			Currency:   balance.Currency,
			UID:        balance.UID,
			Amount:     balance.Amount,
			TxnID:      balance.LatestTxnID,
			CreateTime: oldStartingBalanceCreateTimeUnix,
		}
	}

	updateBalanceTxnMap := func(mainMap StartingBalanceMap, offerMap StartingBalanceMap) StartingBalanceMap {
		for key, value := range offerMap {
			mainMap[key] = value
		}

		return mainMap
	}

	totalStartingBalanceMap := make(StartingBalanceMap, 0)
	updateBalanceTxnMap(totalStartingBalanceMap, oldStartingBalanceMap)
	updateBalanceTxnMap(totalStartingBalanceMap, upperStartingBalanceMap)
	updateBalanceTxnMap(totalStartingBalanceMap, lowerStartingBalanceMap)

	holdWithdrawalAggRowsSQL := `
    SELECT
        user_id
        , coin_id
        , SUM(amount) AS sum_amount
    FROM
        withdraw
    WHERE
        (
            (
                status IN ('Pending', 'Processing', 'Confirmed', 'Waiting Confirmation')
                OR (
                    status = 'Approved'
                    AND close_time >= ?
                )
                OR (
                    status IN ('Canceled', 'Rejected')
                    AND close_time >= ?
                )
            )
            AND date_created < ?
        )
    GROUP BY
        user_id
        , coin_id
    `
	holdWithdrawalAggRows, err := db.
		Raw(
			holdWithdrawalAggRowsSQL,
			dateTime.Add(24*time.Hour).Unix(),
			dateTime.Unix(),
			dateTime.Format(constants.DateTimeFormatISO),
		).
		Rows()
	comutils.PanicOnError(err)

	curerncyInfoMap := currencymod.GetAllLegacyCurrencyInfoIdMapFastF()
	for holdWithdrawalAggRows.Next() {
		var (
			uid        meta.UID
			coinID     uint16
			holdAmount decimal.Decimal
		)
		comutils.PanicOnError(
			holdWithdrawalAggRows.Scan(&uid, &coinID, &holdAmount),
		)

		currencyInfo, ok := curerncyInfoMap[coinID]
		if !ok {
			panic(utils.IssueErrorf("invalid coin id `%v`", coinID))
		}

		mapKey := genTxnKey(currencyInfo.Currency, uid)
		startBalanceTxn, ok := totalStartingBalanceMap[mapKey]
		if !ok {
			startBalanceTxn = &models.PayoutUserStartDateBalance{
				Date:     dateStr,
				Currency: currencyInfo.Currency,
				UID:      uid,
				Amount:   decimal.Zero,
			}
		}

		startBalanceTxn.Amount = startBalanceTxn.Amount.Add(holdAmount)
		totalStartingBalanceMap[mapKey] = startBalanceTxn
	}

	balanceModels := make([]models.PayoutUserStartDateBalance, 0, len(totalStartingBalanceMap))
	for _, startingBalance := range totalStartingBalanceMap {
		balanceModels = append(balanceModels, *startingBalance)
	}
	comutils.PanicOnError(
		database.GetDbMaster().CreateInBatches(balanceModels, 200).Error,
	)
}
