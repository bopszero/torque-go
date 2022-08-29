package crontasks

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

var systemTradingStartBalanceWalletCurrencies = comtypes.NewHashSetFromListF([]meta.Currency{
	constants.CurrencyTorque,
})

func systemTradingGetStartBalancesCurrencies() []meta.Currency {
	balanceCurrencies := []meta.Currency{
		constants.CurrencyTorque,
	}
	for currency := range currencymod.GetAllLegacyCurrencyInfoMapFastF() {
		balanceCurrencies = append(balanceCurrencies, currency)
	}

	return balanceCurrencies
}

func SystemTradingGenStartBalancesFull(dateStr string, timeBoundaryReach time.Duration) {
	indexDate, err := comutils.TimeParse(constants.DateFormatISO, dateStr)
	comutils.PanicOnError(err)

	var (
		pageSize       = 5000
		lowerBoundTime time.Time
		upperBoundTime time.Time
	)
	if timeBoundaryReach == 0 {
		lowerBoundTime = time.Unix(0, 0)
		upperBoundTime = time.Unix(32503680000, 0) // 3000-01-01
	} else {
		lowerBoundTime = indexDate.Add(-timeBoundaryReach)
		upperBoundTime = indexDate.Add(timeBoundaryReach)
	}

	type (
		RecordKey struct {
			Currency meta.Currency
			UID      meta.UID
		}
		StartingBalanceMap map[RecordKey]decimal.Decimal
	)

	genRecordKey := func(currency meta.Currency, uid meta.UID) RecordKey {
		return RecordKey{Currency: currency, UID: uid}
	}
	getIdsFromRows := func(rows *sql.Rows) (ids []uint64) {
		var idx uint64
		for rows.Next() {
			comutils.PanicOnError(
				rows.Scan(&idx),
			)
			ids = append(ids, idx)
		}
		return ids
	}
	fetchTxnsInChunk := func(db *gorm.DB, ids []uint64) []models.UserBalanceTxn {
		var (
			txns   = make([]models.UserBalanceTxn, 0, len(ids))
			maxIdx = len(ids) - 1
		)
		for i := 0; ; i += pageSize {
			toIdx := i + pageSize
			if toIdx > maxIdx {
				toIdx = maxIdx
			}
			var (
				pageIds  = ids[i:toIdx]
				pageTxns []models.UserBalanceTxn
			)
			comutils.PanicOnError(db.
				Where(dbquery.In(models.CommonColID, pageIds)).
				Find(&pageTxns).
				Error,
			)
			txns = append(txns, pageTxns...)
			if len(pageTxns) < pageSize {
				break
			}
		}
		return txns
	}

	updateBalanceTxnMap := func(mainMap StartingBalanceMap, offerMap StartingBalanceMap) StartingBalanceMap {
		for key, value := range offerMap {
			mainMap[key] = value
		}

		return mainMap
	}
	genStartBalanceMap := func(db *gorm.DB, currencies []meta.Currency) StartingBalanceMap {
		lowerStartingBalanceIdsRows, err := db.
			Model(&models.UserBalanceTxn{}).
			Where(dbquery.In(models.CommonColCurrency, currencies)).
			Where(dbquery.Between(models.UserBalanceTxnColCreateTime, lowerBoundTime.Unix(), indexDate.Unix()-1)).
			Select("MAX(id) AS max_id").
			Group(dbquery.JoinExpr(
				models.UserBalanceTxnColCurrency,
				models.UserBalanceTxnColUID)).
			Rows()
		comutils.PanicOnError(err)
		var (
			lowerStartingBalanceIds = getIdsFromRows(lowerStartingBalanceIdsRows)
			lowerStartingBalanceMap = make(StartingBalanceMap, len(lowerStartingBalanceIds))
			lowerTxns               = fetchTxnsInChunk(db, lowerStartingBalanceIds)
		)
		for _, txn := range lowerTxns {
			key := genRecordKey(txn.Currency, txn.UserID)
			lowerStartingBalanceMap[key] = txn.Balance
		}

		upperStartingBalanceIdsRows, err := db.
			Model(&models.UserBalanceTxn{}).
			Where(dbquery.In(models.CommonColCurrency, currencies)).
			Where(dbquery.Between(models.UserBalanceTxnColCreateTime, indexDate.Unix(), upperBoundTime.Unix())).
			Select("MIN(id) AS min_id").
			Group(dbquery.JoinExpr(
				models.UserBalanceTxnColCurrency,
				models.UserBalanceTxnColUID)).
			Rows()
		comutils.PanicOnError(err)
		var (
			upperStartingBalanceIds = getIdsFromRows(upperStartingBalanceIdsRows)
			upperStartingBalanceMap = make(StartingBalanceMap, len(upperStartingBalanceIds))
			upperTxns               = fetchTxnsInChunk(db, upperStartingBalanceIds)
		)
		for _, txn := range upperTxns {
			key := genRecordKey(txn.Currency, txn.UserID)
			upperStartingBalanceMap[key] = txn.Balance.Sub(txn.Amount)
		}

		var oldBalances []models.UserBalance
		comutils.PanicOnError(db.
			Where(dbquery.In(models.CommonColCurrency, currencies)).
			Where(dbquery.Lt(models.UserBalanceColUpdateTime, indexDate.Unix())).
			Where(dbquery.Gt(models.UserBalanceColAmount, 0)).
			Find(&oldBalances).
			Error,
		)
		oldStartingBalanceMap := make(StartingBalanceMap, len(oldBalances))
		for _, balance := range oldBalances {
			key := genRecordKey(balance.Currency, balance.UID)
			oldStartingBalanceMap[key] = balance.Amount
		}

		totalStartBalanceMap := make(StartingBalanceMap, len(lowerStartingBalanceMap))
		updateBalanceTxnMap(totalStartBalanceMap, oldStartingBalanceMap)
		updateBalanceTxnMap(totalStartBalanceMap, upperStartingBalanceMap)
		updateBalanceTxnMap(totalStartBalanceMap, lowerStartingBalanceMap)

		return totalStartBalanceMap
	}

	var (
		allCurrencies     = systemTradingGetStartBalancesCurrencies()
		tradingCurrencies []meta.Currency
		walletCurrencies  []meta.Currency
	)
	for _, currency := range allCurrencies {
		if systemTradingStartBalanceWalletCurrencies.Contains(currency) {
			walletCurrencies = append(walletCurrencies, currency)
		} else {
			tradingCurrencies = append(tradingCurrencies, currency)
		}
	}

	var (
		dbMain                 = database.GetDbSlave()
		dbWallet               = database.GetDbF(database.AliasWalletSlave)
		tradingStartBalanceMap = genStartBalanceMap(dbMain.DB, tradingCurrencies)
		walletStartBalanceMap  = genStartBalanceMap(dbWallet.DB, walletCurrencies)
		totalStartBalanceMap   = make(StartingBalanceMap)
	)
	updateBalanceTxnMap(totalStartBalanceMap, tradingStartBalanceMap)
	updateBalanceTxnMap(totalStartBalanceMap, walletStartBalanceMap)

	currencyStartingBalanceMap := make(map[meta.Currency]decimal.Decimal)
	for key, startingBalance := range totalStartBalanceMap {
		balance := currencyStartingBalanceMap[key.Currency]
		currencyStartingBalanceMap[key.Currency] = balance.Add(startingBalance)
	}

	for _, currency := range allCurrencies {
		if _, ok := currencyStartingBalanceMap[currency]; !ok {
			currencyStartingBalanceMap[currency] = decimal.Zero
		}
	}

	var (
		systemModels = make([]models.SystemStartDateBalance, 0, len(currencyStartingBalanceMap))
		now          = time.Now()
	)
	for currency, balance := range currencyStartingBalanceMap {
		systemModels = append(
			systemModels,
			models.SystemStartDateBalance{
				Date:       indexDate.Format(constants.DateFormatISO),
				Currency:   currency,
				Amount:     balance,
				CreateTime: now.Unix(),
			},
		)
	}
	comutils.PanicOnError(
		database.GetDbMaster().CreateInBatches(systemModels, 200).Error,
	)
}

func SystemTradingGenStartBalancesAccumulate(dateStr string) {
	indexDate, err := comutils.TimeParse(constants.DateFormatISO, dateStr)
	comutils.PanicOnError(err)

	var (
		dbMain  = database.GetDbSlave()
		maxDate string
	)
	err = dbMain.
		Model(&models.SystemStartDateBalance{}).
		Select(dbquery.Max(models.CommonColDate)).
		Where(dbquery.Lt(models.CommonColDate, dateStr)).
		Row().
		Scan(&maxDate)
	if err != nil {
		panic(utils.IssueErrorf(
			"cannot fetch the latest system starting balance | err=%s",
			err.Error()),
		)
	}
	maxDate = maxDate[:10] // Cut the auto parse like `2020-06-01T00:00:00+07:00`

	var prevBalances []models.SystemStartDateBalance
	err = dbMain.
		Where(&models.SystemStartDateBalance{Date: maxDate}).
		Find(&prevBalances).
		Error
	comutils.PanicOnError(err)
	prevBalanceMap := make(map[meta.Currency]models.SystemStartDateBalance, len(prevBalances))
	for _, balance := range prevBalances {
		prevBalanceMap[balance.Currency] = balance
	}

	currencyBalanceMap := make(map[meta.Currency]models.SystemStartDateBalance, len(prevBalances))
	maxDateTime, err := comutils.TimeParse(constants.DateFormatISO, maxDate)
	comutils.PanicOnError(err)

	fillCurrencyBalanceMap := func(db *gorm.DB, currencies []meta.Currency) {
		currencySumRows, err := db.
			Model(&models.UserBalanceTxn{}).
			Select(dbquery.JoinExpr(
				models.UserBalanceTxnColCurrency,
				dbquery.Sum(models.UserBalanceTxnColAmount),
			)).
			Where(dbquery.In(models.CommonColCurrency, currencies)).
			Where(dbquery.Between(models.UserBalanceTxnColCreateTime, maxDateTime.Unix(), indexDate.Unix()-1)).
			Group(models.UserBalanceTxnColCurrency).
			Rows()
		comutils.PanicOnError(err)

		now := time.Now()
		for currencySumRows.Next() {
			var (
				currency      meta.Currency
				inRangeAmount decimal.Decimal
			)
			if err = currencySumRows.Scan(&currency, &inRangeAmount); err != nil {
				comutils.PanicOnError(err)
			}

			var balanceAmount decimal.Decimal
			if prevModel, ok := prevBalanceMap[currency]; ok {
				balanceAmount = prevModel.Amount.Add(inRangeAmount)
			} else {
				balanceAmount = inRangeAmount
			}

			currencyBalanceMap[currency] = models.SystemStartDateBalance{
				Date:       indexDate.Format(constants.DateFormatISO),
				Currency:   currency,
				Amount:     balanceAmount,
				CreateTime: now.Unix(),
			}
		}
	}

	var (
		allCurrencies     = systemTradingGetStartBalancesCurrencies()
		tradingCurrencies []meta.Currency
		walletCurrencies  []meta.Currency
	)
	for _, currency := range allCurrencies {
		if systemTradingStartBalanceWalletCurrencies.Contains(currency) {
			walletCurrencies = append(walletCurrencies, currency)
		} else {
			tradingCurrencies = append(tradingCurrencies, currency)
		}
	}

	dbWallet := database.GetDbF(database.AliasWalletSlave)
	fillCurrencyBalanceMap(dbMain.DB, tradingCurrencies)
	fillCurrencyBalanceMap(dbWallet.DB, walletCurrencies)

	now := time.Now()
	for _, currency := range systemTradingGetStartBalancesCurrencies() {
		if _, ok := currencyBalanceMap[currency]; ok {
			continue
		}

		var missingBalance decimal.Decimal
		if prevModel, ok := prevBalanceMap[currency]; ok {
			missingBalance = prevModel.Amount
		} else {
			missingBalance = decimal.Zero
		}
		currencyBalanceMap[currency] = models.SystemStartDateBalance{
			Date:       indexDate.Format(constants.DateFormatISO),
			Currency:   currency,
			Amount:     missingBalance,
			CreateTime: now.Unix(),
		}
	}

	systemModels := make([]models.SystemStartDateBalance, 0, len(currencyBalanceMap))
	for _, model := range currencyBalanceMap {
		systemModels = append(systemModels, model)
	}
	comutils.PanicOnError(
		database.GetDbMaster().CreateInBatches(systemModels, 200).Error,
	)
}
