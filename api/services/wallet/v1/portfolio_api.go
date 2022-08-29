package v1

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/jszwec/csvutil"
	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/go-common/comcache"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/api"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/api/responses"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/settingmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/usermod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/balancemod"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/ordermod"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/ordermod/orderchannel"
)

const (
	PortfolioExportOrdersLimit = 5000
)

func PortfolioGetOverview(c echo.Context) error {
	var (
		ctx = apiutils.EchoWrapContext(c)
		uid = apiutils.GetContextUidF(ctx)
	)
	var (
		allCurrencyInfoMap = currencymod.GetAllCurrencyInfoMapFastF()
		coinMap            = blockchainmod.GetCoinMap()
		responseBalances   = make([]PortfolioGetOverviewBalance, 0, 1+len(coinMap))
	)
	torqueBalance, err := balancemod.GetUserBalance(ctx, uid, constants.CurrencyTorque)
	if err != nil {
		return err
	}
	var (
		torqueCurrencyInfo = allCurrencyInfoMap[constants.CurrencyTorque]
		torqueBalanceInfo  = PortfolioGetOverviewBalance{
			Currency:          constants.CurrencyTorque,
			CurrencyPriceUSD:  torqueCurrencyInfo.PriceUSD,
			CurrencyPriceUSDT: torqueCurrencyInfo.PriceUSDT,
			Amount:            torqueBalance.Amount,
			UpdateTime:        torqueBalance.UpdateTime,
			IsAvailable:       true,
		}
	)
	responseBalances = append(responseBalances, torqueBalanceInfo)

	var (
		logger            = comlogging.GetLogger()
		allNetworkInfoMap = currencymod.GetAllBlockchainNetworkInfoMapFastF()
	)
	for index, coin := range coinMap {
		if _, ok := allNetworkInfoMap[index.Network]; !ok {
			continue
		}
		currencyInfo, ok := allCurrencyInfoMap[index.Currency]
		if !ok || !currencymod.IsDisplayableWalletInfo(currencyInfo) {
			continue
		}
		var (
			currency        = coin.GetCurrency()
			repsonseBalance = PortfolioGetOverviewBalance{
				Currency:          currency,
				Network:           coin.GetNetwork(),
				CurrencyPriceUSD:  currencyInfo.PriceUSD,
				CurrencyPriceUSDT: currencyInfo.PriceUSDT,
				IsAvailable:       true,
			}
		)
		account, err := coin.NewAccountSystem(ctx, uid)
		if err != nil {
			return err
		}
		if balance, err := account.GetBalance(); err == nil {
			repsonseBalance.Amount = balance
		} else {
			repsonseBalance.IsAvailable = false
			logger.
				WithContext(ctx).
				WithError(err).
				WithField("coin", coin.GetIndexCode()).
				Errorf("blockchain account get balance failed | err=%s", err.Error())
		}
		responseBalances = append(responseBalances, repsonseBalance)
	}
	sort.Slice(
		responseBalances,
		func(i, j int) bool {
			var (
				leftInfo  = allCurrencyInfoMap[responseBalances[i].Currency]
				rightInfo = allCurrencyInfoMap[responseBalances[j].Currency]
			)
			return leftInfo.PriorityWallet < rightInfo.PriorityWallet
		},
	)

	return responses.Ok(
		ctx,
		PortfolioGetOverviewResponse{
			Balances: responseBalances,
		},
	)
}

func PortfolioGetCurrency(c echo.Context) error {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		uid      = apiutils.GetContextUidF(ctx)
		reqModel PortfolioGetCurrencyRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}
	return portfolioGetCurrencyByUID(c, uid, reqModel.Currency, reqModel.Network)
}

func portfolioGetCurrencyByUID(
	c echo.Context,
	uid meta.UID, currency meta.Currency, network meta.BlockchainNetwork,
) error {
	var (
		ctx          = apiutils.EchoWrapContext(c)
		currencyInfo = currencymod.GetCurrencyInfoFastF(currency)
	)
	if !currencymod.IsValidWalletInfo(currencyInfo) {
		return utils.WrapError(constants.ErrorCurrency)
	}

	currencyNotice, err := settingmod.GetSettingValueFast(
		constants.SettingKeyCurrencyNoticePattern, currency)
	if err != nil {
		return err
	}
	responseModel := PortfolioGetCurrencyResponse{
		Currency:  currency,
		PriceUSD:  currencyInfo.PriceUSD,
		PriceUSDT: currencyInfo.PriceUSDT,
		Notice:    currencyNotice,
	}
	if blockchainmod.IsSupportedIndex(currency, network) {
		// TODO: Temporary disabled
		return meta.NewMessageError("Deposit is temporarily suspended. Please try again later.")

		coin := blockchainmod.GetCoinNativeF(currency)
		responseModel.Network = coin.GetNetwork()

		account, err := coin.NewAccountSystem(ctx, uid)
		if err != nil {
			return err
		}
		responseModel.AccountNo = account.GetAddress()
		if responseModel.Balance, err = account.GetBalance(); err != nil {
			return err
		}
	} else if currencymod.IsLocalCurrency(currencyInfo.Currency) {
		user := usermod.GetUserFastF(uid)
		responseModel.AccountNo = user.Code

		balanceModel, err := balancemod.GetUserBalance(ctx, uid, currency)
		if err != nil {
			return err
		}
		responseModel.Balance = balanceModel.Amount
	} else {
		return utils.WrapError(constants.ErrorCurrency)
	}

	return responses.Ok(ctx, responseModel)
}

func PortfolioListCurrencyOrder(c echo.Context) (err error) {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		uid      = apiutils.GetContextUidF(ctx)
		reqModel PortfolioListOrdersRequest
	)
	if err = api.BindAndValidate(c, &reqModel); err != nil {
		return
	}
	currencyInfo := currencymod.GetCurrencyInfoFastF(reqModel.Currency)
	if !currencymod.IsValidWalletInfo(currencyInfo) {
		return utils.WrapError(constants.ErrorCurrency)
	}

	var orders []models.Order
	if blockchainmod.IsSupportedIndex(currencyInfo.Currency, reqModel.Network) {
		orders, err = listOrderBlockchain(ctx, uid, reqModel)
	} else if currencymod.IsLocalCurrency(currencyInfo.Currency) {
		orders, err = listOrderLocal(ctx, uid, reqModel)
	} else {
		err = utils.WrapError(constants.ErrorCurrency)
	}
	if err != nil {
		return
	}

	responseOrders := make([]Order, len(orders))
	for i, order := range orders {
		responseOrder := &responseOrders[i]
		if err = dumpOrder(&order, responseOrder); err != nil {
			return
		}

		var srcMeta interface{}
		err = ordermod.GetOrderSrcChannelMetaData(&order, &srcMeta)
		if err != nil && !utils.IsOurError(err, constants.ErrorCodeDataNotFound) {
			return
		}
		responseOrder.SrcChannelContext = OrderChannelContext{
			Meta: srcMeta,
		}

		var dstMeta interface{}
		err = ordermod.GetOrderDstChannelMetaData(&order, &dstMeta)
		if err != nil && !utils.IsOurError(err, constants.ErrorCodeDataNotFound) {
			return
		}
		responseOrder.DstChannelContext = OrderChannelContext{
			Meta: dstMeta,
		}
	}

	return responses.Ok(
		ctx,
		PortfolioListOrdersResponse{
			Items:  responseOrders,
			Paging: reqModel.Paging,
		},
	)
}

func listOrderBlockchain(ctx comcontext.Context, uid meta.UID, reqModel PortfolioListOrdersRequest) (
	orders []models.Order, err error,
) {
	coin := blockchainmod.GetCoinNativeF(reqModel.Currency)
	account, err := coin.NewAccountSystem(ctx, uid)
	if err != nil {
		return
	}
	txns, err := account.GetTxns(reqModel.Paging)
	if err != nil {
		return
	}
	if len(txns) == 0 {
		return
	}
	orders = make([]models.Order, 0, len(txns))

	blockchainSendHashSet := make(comtypes.HashSet, len(txns))
	for _, txn := range txns {
		order := ordermod.NewOrder(txn.GetCurrency())
		if err = orderchannel.DumpBlockchainTxnToOrder(uid, txn, &order, false); err != nil {
			return
		}
		if order.DstChannelRef != "" {
			blockchainSendHashSet.Add(order.DstChannelRef)
		}
		orders = append(orders, order)
	}

	appendLocalSendOrders := func() (err error) {
		fromTime := txns[len(txns)-1].GetTimeUnix()
		var toTime int64
		if reqModel.Paging.Offset == 0 {
			toTime = time.Now().Unix()
			if len(txns) < int(reqModel.Paging.Limit) {
				fromTime = 0
			}
		} else {
			toTime = txns[0].GetTimeUnix()
		}

		var localOrders []models.Order
		err = ordermod.GetOrderDB(database.GetDbF(database.AliasWalletSlave).DB).
			FilterUserCurrency(uid, reqModel.Currency).
			FilterBlockchainTxns().
			Where(dbquery.Between(models.OrderColCreateTime, fromTime, toTime)).
			Find(&localOrders).
			Error
		if err != nil {
			return
		}
		localOrderMap := make(map[string]models.Order)
		for _, order := range localOrders {
			sendOrder := order
			if !blockchainSendHashSet.Contains(sendOrder.DstChannelRef) {
				orders = append(orders, order)
			}
			localOrderMap[sendOrder.DstChannelRef] = sendOrder
		}

		for i := range orders {
			order := &orders[i]
			if order.CreateTime == 0 {
				if localOrder, ok := localOrderMap[order.DstChannelRef]; ok {
					order.CreateTime = localOrder.CreateTime
					order.UpdateTime = localOrder.CreateTime
				}
			}
		}

		return nil
	}
	replaceLocalBusinessOrders := func() (err error) {
		var localOrders []models.Order
		err = ordermod.GetOrderDB(database.GetDbF(database.AliasWalletSlave).DB).
			FilterUserCurrency(uid, reqModel.Currency).
			FilterBlockchainTxns().
			Where(dbquery.In(models.OrderColDstChannelRef, blockchainSendHashSet.AsList())).
			Find(&localOrders).
			Error
		if err != nil {
			return
		}
		localOrderMap := make(map[string]models.Order)
		for _, order := range localOrders {
			localOrderMap[order.DstChannelRef] = order
		}

		for i, order := range orders {
			if nonSendOrder, ok := localOrderMap[order.DstChannelRef]; ok {
				for key, value := range order.ExtraData {
					if strings.HasPrefix(key, "_") {
						continue
					}
					nonSendOrder.ExtraData[key] = value
				}
				orders[i] = nonSendOrder
			}
		}

		return nil
	}

	if err = appendLocalSendOrders(); err != nil {
		return
	}
	if err = replaceLocalBusinessOrders(); err != nil {
		return
	}

	for i := range orders {
		if err = orderchannel.DumpBlockchainInfoToOrder(&orders[i], nil); err != nil {
			return
		}
	}

	sort.Slice(
		orders,
		func(i, j int) bool {
			return orders[j].CreateTime < orders[i].CreateTime
		},
	)

	return orders, nil
}

func listOrderLocal(ctx comcontext.Context, uid meta.UID, reqModel PortfolioListOrdersRequest) (
	orders []models.Order, err error,
) {
	query := ordermod.
		GetOrderDB(database.GetDbF(database.AliasWalletSlave).DB).
		FilterUserCurrency(uid, reqModel.Currency).
		Order(dbquery.OrderDesc(models.OrderColCreateTime))
	err = dbquery.Paging(query, reqModel.Paging).
		Find(&orders).
		Error
	if err != nil {
		err = utils.WrapError(err)
		return
	}

	return
}

func PortfolioExportGenTokenCurrencyOrder(c echo.Context) (err error) {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		uid      = apiutils.GetContextUidF(ctx)
		reqModel PortfolioExportOrdersRequest
	)
	if err = api.BindAndValidate(c, &reqModel); err != nil {
		return
	}
	if _, err = blockchainmod.GetCoin(reqModel.Currency, reqModel.Network); err != nil {
		return
	}

	var (
		token      = comutils.NewUuid4Code()
		tokenModel = PortfolioExportOrdersTokenModel{
			PortfolioExportOrdersRequest: reqModel,
			UID:                          uid,
		}
	)
	cacheKey := fmt.Sprintf("PortfolioExportCurrencyOrder:%s", token)
	err = comcache.GetRemoteCache().Set(cacheKey, tokenModel, 30*time.Second)
	if err != nil {
		return
	}

	return responses.Ok(ctx, PortfolioExportOrdersResponse{Token: token})
}

func PortfolioExportDownloadCurrencyOrder(c echo.Context) (err error) {
	var (
		ctx         = apiutils.EchoWrapContext(c)
		token       = c.Param("token")
		cacheKey    = fmt.Sprintf("PortfolioExportCurrencyOrder:%s", token)
		cacheClient = comcache.GetRemoteCache()
		tokenModel  PortfolioExportOrdersTokenModel
	)
	err = cacheClient.Get(cacheKey, &tokenModel)
	if err != nil {
		c.String(http.StatusInternalServerError, constants.ErrorUnknown.Message(ctx))
		return utils.WrapError(err)
	}
	cacheClient.Delete(cacheKey)

	channelInfoMap, err := orderchannel.GetChannelInfoMapFast()
	if err != nil {
		c.String(http.StatusInternalServerError, constants.ErrorUnknown.Message(ctx))
		return utils.WrapError(err)
	}
	var (
		orders      []models.Order
		ordersCount int64
		toTime      = tokenModel.ToTime - 1
		query       = ordermod.
			GetOrderDB(database.GetDbF(database.AliasWalletSlave).DB).
			FilterUserCurrency(tokenModel.UID, tokenModel.Currency).
			Model(&models.Order{}).
			Where(dbquery.Between(models.OrderColCreateTime, tokenModel.FromTime, toTime)).
			Order(dbquery.OrderDesc(models.OrderColCreateTime))
	)
	if err = query.Count(&ordersCount).Error; err != nil {
		c.String(http.StatusInternalServerError, constants.ErrorUnknown.Message(ctx))
		return utils.WrapError(err)
	}
	err = query.
		Limit(PortfolioExportOrdersLimit).
		Find(&orders).
		Error
	if err != nil {
		c.String(http.StatusInternalServerError, constants.ErrorUnknown.Message(ctx))
		return utils.WrapError(err)
	}

	orderToExportItem := func(order models.Order) OrderExportItem {
		item := OrderExportItem{
			Code:     order.Code,
			Currency: order.Currency,
			Amount:   order.AmountTotal,
			Time:     time.Unix(order.CreateTime, 0).Format(constants.DateTimeFormatISO),
			Note:     order.Note,
		}
		if order.Status < constants.OrderStatusNew {
			item.Status = "Fail"
		} else if order.Status == constants.OrderStatusCompleting ||
			order.Status == constants.OrderStatusCompleted {
			item.Status = "Success"
		} else {
			item.Status = "Pending"
		}
		channelInfo, ok := channelInfoMap[ordermod.GetOrderMainChannelType(order)]
		if ok {
			item.Type = channelInfo.Name
		} else {
			item.Type = "<Unknown>"
		}
		return item
	}
	var orderFromIdx = 0
	readerGen := utils.NewReaderGenerator(func() ([]byte, error) {
		var (
			pageItemsCount = 500
			pageToIdx      = comutils.MinInt(orderFromIdx+pageItemsCount, len(orders))
			pageOrders     = orders[orderFromIdx:pageToIdx]
			exportOrders   = make([]OrderExportItem, len(pageOrders))
		)
		for i, order := range pageOrders {
			exportOrders[i] = orderToExportItem(order)
		}

		var (
			buf       bytes.Buffer
			csvWriter = csv.NewWriter(&buf)
		)
		csvEncoder := csvutil.NewEncoder(csvWriter)
		if err := csvEncoder.Encode(exportOrders); err != nil {
			return nil, utils.WrapError(err)
		}
		csvWriter.Flush()

		orderFromIdx = pageToIdx
		return buf.Bytes(), nil
	})

	exportFileName := fmt.Sprintf(
		"[%v] %s -- %s",
		tokenModel.Currency,
		time.Unix(tokenModel.FromTime, 0).Format(constants.DateFormatISO),
		time.Unix(toTime, 0).Format(constants.DateFormatISO),
	)
	if ordersCount > PortfolioExportOrdersLimit {
		exportFileName += fmt.Sprintf(" (%v of %v)", PortfolioExportOrdersLimit, ordersCount)
	}
	return apiutils.EchoResponseDownloadSteam(
		ctx,
		exportFileName+".csv", constants.ContentTypeCSV, readerGen,
	)
}

func PortfolioGetCurrencyOrder(c echo.Context) (err error) {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		uid      = apiutils.GetContextUidF(ctx)
		reqModel PortfolioGetCurrencyOrdersRequest
	)
	if err = api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}
	currencyInfo := currencymod.GetCurrencyInfoFastF(reqModel.Currency)
	if !currencymod.IsValidWalletInfo(currencyInfo) {
		return utils.WrapError(constants.ErrorCurrency)
	}

	var order *models.Order
	if blockchainmod.IsSupportedIndex(currencyInfo.Currency, reqModel.Network) {
		order, err = getOrderBlockchain(ctx, uid, reqModel)
	} else if currencymod.IsLocalCurrency(currencyInfo.Currency) {
		order, err = getOrderLocal(ctx, uid, reqModel)
	} else {
		err = utils.WrapError(constants.ErrorCurrency)
	}
	if err != nil {
		return err
	}

	var responseModel PortfolioGetCurrencyOrderResponse
	if err = dumpOrder(order, (*Order)(&responseModel)); err != nil {
		return
	}

	var (
		srcMeta    interface{}
		srcDetails interface{}
	)
	err = ordermod.GetOrderSrcChannelMetaData(order, &srcMeta)
	if err != nil && !utils.IsOurError(err, constants.ErrorCodeDataNotFound) {
		return err
	}
	// TODO: Update
	// if srcDetails, err = ordermod.GetOrderSrcChannelDetails(ctx, order); err != nil {
	// 	return err
	// }
	responseModel.SrcChannelContext = OrderChannelContext{
		Meta:    srcMeta,
		Details: srcDetails,
	}

	var (
		dstMeta    interface{}
		dstDetails interface{}
	)
	err = ordermod.GetOrderDstChannelMetaData(order, &dstMeta)
	if err != nil && !utils.IsOurError(err, constants.ErrorCodeDataNotFound) {
		return err
	}
	if dstDetails, err = ordermod.GetOrderDstChannelDetails(ctx, order); err != nil {
		return err
	}
	responseModel.DstChannelContext = OrderChannelContext{
		Meta:    dstMeta,
		Details: dstDetails,
	}

	return responses.Ok(ctx, responseModel)
}

func getOrderBlockchain(ctx comcontext.Context, uid meta.UID, reqModel PortfolioGetCurrencyOrdersRequest) (
	order *models.Order, err error,
) {
	coin := blockchainmod.GetCoinNativeF(reqModel.Currency)
	account, err := coin.NewAccountSystem(ctx, uid)
	if err != nil {
		return
	}

	txn, err := account.GetTxn(reqModel.Reference)
	if err != nil && !utils.IsOurError(err, constants.ErrorCodeDataNotFound) {
		return
	}

	localOrder := ordermod.NewOrder(reqModel.Currency)
	err = ordermod.GetOrderDB(database.GetDbF(database.AliasWalletSlave).DB).
		FilterUserCurrency(uid, reqModel.Currency).
		FilterBlockchainTxn(reqModel.Currency, reqModel.Reference).
		First(&localOrder).
		Error
	if database.IsDbError(err) {
		err = utils.WrapError(err)
		return
	}

	if txn == nil || txn.GetHash() == "" {
		if localOrder.ID == 0 {
			return nil, constants.ErrorOrderNotFound
		}
		if err = orderchannel.DumpBlockchainInfoToOrder(&localOrder, nil); err != nil {
			return
		}
		order = &localOrder
		return
	}

	txnOrder := ordermod.NewOrder(txn.GetCurrency())
	if err = orderchannel.DumpBlockchainTxnToOrder(uid, txn, &txnOrder, true); err != nil {
		return
	}
	if localOrder.ID == 0 {
		order = &txnOrder
		return
	}

	for key, value := range txnOrder.ExtraData {
		if strings.HasPrefix(key, "_") {
			continue
		}
		localOrder.ExtraData[key] = value
	}

	order = &localOrder
	return
}

func getOrderLocal(ctx comcontext.Context, uid meta.UID, reqModel PortfolioGetCurrencyOrdersRequest) (
	_ *models.Order, err error,
) {
	orderID, err := comutils.ParseUint64(reqModel.Reference)
	if err != nil {
		return
	}

	order, err := ordermod.GetUserCurrencyOrder(reqModel.Currency, uid, orderID)
	if err != nil {
		return
	}

	return &order, nil
}
