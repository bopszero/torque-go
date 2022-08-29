package v1

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/api"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/api/responses"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/usermod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/ordermod"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/ordermod/orderchannel"
	"gorm.io/gorm"
)

func LegacyTorqueTxnGet(c echo.Context) (err error) {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel = new(LegacyTxnListRequest)
	)
	if err = api.BindAndValidate(c, reqModel); err != nil {
		return
	}

	var orders []models.Order
	ordersQuery := database.
		GetDbF(database.AliasWalletSlave).
		Order(dbquery.OrderDesc(models.OrderColCreateTime)).
		Where(dbquery.In(
			models.OrderColStatus,
			[]meta.OrderStatus{
				constants.OrderStatusHandleSrc,
				constants.OrderStatusHandleDst,
				constants.OrderStatusNeedStaff,
				constants.OrderStatusRefunding,
				constants.OrderStatusCompleting,
				constants.OrderStatusCompleted,
			},
		)).
		Where(&models.Order{UID: reqModel.UID, Currency: reqModel.Currency})

	if reqModel.OrderID == 0 {
		paging := meta.Paging{
			Offset: comutils.MaxUint(0, (reqModel.PagingPage-1)*reqModel.PagingLimit),
			Limit:  comutils.ClampUint(reqModel.PagingLimit, 1, 30),
		}

		ordersQuery = dbquery.Paging(ordersQuery, paging)
	} else {
		ordersQuery = ordersQuery.Where(&models.Order{ID: reqModel.OrderID}).Limit(1)
	}
	if err = ordersQuery.Find(&orders).Error; err != nil {
		return
	}

	db := database.GetDbSlave()
	legacyOrders := make([]*LegacyTxnListItem, 0)
	for _, order := range orders {
		items, err := genLegacyTorqueOrderData(c, db.DB, &order)
		if err != nil {
			return err
		}

		legacyOrders = append(legacyOrders, items...)
	}

	responseBody := meta.O{
		"orders": legacyOrders,
	}
	return responses.Ok(ctx, responseBody)
}

func genLegacyTorqueOrderData(
	c echo.Context, db *gorm.DB, order *models.Order,
) (dataList []*LegacyTxnListItem, err error) {
	data := LegacyTxnListItem{
		ID:          order.ID,
		UID:         order.UID,
		Currency:    order.Currency,
		CreatedTime: time.Unix(order.CreateTime, 0).Format(constants.DateTimeFormatISO),
	}
	dataList = []*LegacyTxnListItem{
		&data,
	}

	switch order.Direction {
	case constants.OrderDirectionPayment:
		data.Ref = comutils.Stringify(order.DstChannelID)
		data.Amount = order.DstChannelAmount.Neg()

		switch order.DstChannelType {
		case constants.ChannelTypeDstTransfer:
			data.TypeCode = "TRANSFER"

			if order.DstChannelRef != "" {
				var p2pTransfer models.P2PTransfer
				err = db.First(&p2pTransfer, fmt.Sprintf("id=%v", order.DstChannelRef)).Error
				if database.IsDbError(err) {
					return
				}
				if p2pTransfer.ID > 0 {
					fromUser := usermod.GetUserFastF(p2pTransfer.FromUID)
					toUser := usermod.GetUserFastF(p2pTransfer.ToUID)
					data.Details = meta.O{
						"from_username": fromUser.Username,
						"to_username":   toUser.Username,
						"from_currency": p2pTransfer.FromCurrency,
						"to_currency":   p2pTransfer.ToCurrency,
						"to_amount":     p2pTransfer.ToAmount,
						"from_amount":   p2pTransfer.FromAmount,
						"fee_amount":    p2pTransfer.FeeAmount,
					}
				}
			}

		case constants.ChannelTypeDstProfitReinvest, constants.ChannelTypeDstProfitWithdraw:
			if order.DstChannelType == constants.ChannelTypeDstProfitReinvest {
				data.TypeCode = "REINVEST"
			} else {
				data.TypeCode = "TORQ_WITHDRAW"
			}

			var torqueTxn models.TorqueTxn
			err = db.First(&torqueTxn, &models.TorqueTxn{ID: order.DstChannelID}).Error
			if database.IsDbError(err) {
				return
			}
			if torqueTxn.ID > 0 {
				data.Details = meta.O{
					"torque_id":       torqueTxn.ID,
					"coin_name":       torqueTxn.Currency,
					"amount":          torqueTxn.Amount,
					"coin_amount":     torqueTxn.CoinAmount,
					"rate":            torqueTxn.ExchangeRate,
					"status":          torqueTxn.Status,
					"txn_hash":        torqueTxn.BlockchainHash,
					"transactionhash": torqueTxn.BlockchainHash,
					"address":         torqueTxn.Address,
				}
			}

		case constants.ChannelTypeDstMerchantGorillaHotel:
			data.TypeCode = "INVOICE_BOOKING"
			data.Details = meta.O{
				"order_id":      order.DstChannelID,
				"product_title": "Hotel Booking",
			}
			data.OrderDetails = meta.O{
				"merchant_ref": order.DstChannelRef,
			}

		case constants.ChannelTypeDstMerchantTorqueMall:
			data.TypeCode = "INVOICE_TMALL"
			data.Details = meta.O{
				"order_id":      order.DstChannelID,
				"product_title": "Torque Mall Order",
			}
			data.OrderDetails = meta.O{
				"merchant_ref": order.DstChannelRef,
			}

		default:
			dataList = dataList[1:]
			return
		}
	default:
		data.Ref = comutils.Stringify(order.SrcChannelID)
		data.Amount = order.SrcChannelAmount

		switch order.SrcChannelType {
		case constants.ChannelTypeSrcTransfer:
			data.TypeCode = "TRANSFER"

			if order.SrcChannelRef != "" {
				var p2pTransfer models.P2PTransfer
				err = db.First(&p2pTransfer, fmt.Sprintf("id=%v", order.SrcChannelRef)).Error
				if database.IsDbError(err) {
					return
				}
				if p2pTransfer.ID > 0 {
					fromUser := usermod.GetUserFastF(p2pTransfer.FromUID)
					toUser := usermod.GetUserFastF(p2pTransfer.ToUID)
					data.Details = meta.O{
						"from_username": fromUser.Username,
						"to_username":   toUser.Username,
						"from_currency": p2pTransfer.FromCurrency,
						"to_currency":   p2pTransfer.ToCurrency,
						"to_amount":     p2pTransfer.ToAmount,
						"from_amount":   p2pTransfer.FromAmount,
						"fee_amount":    p2pTransfer.FeeAmount,
					}
				}
			}

		case constants.ChannelTypeSrcPromoCode:
			data.TypeCode = "REDEEM"

		case constants.ChannelTypeSrcTradingReward:
			data.TypeCode = "PROFIT"
			data.Ref = order.SrcChannelRef
			dataAffiliateCom := data
			dataAffiliateCom.TypeCode = "AFFILIATE"
			dataLeaderCom := data
			dataLeaderCom.TypeCode = "USER_REWARD"

			var rewardMeta orderchannel.SrcTradingRewardMeta
			err = ordermod.GetOrderChannelMetaData(order, constants.ChannelTypeSrcTradingReward, &rewardMeta)
			if err != nil {
				return
			}

			data.Amount = rewardMeta.DailyProfitAmounts.Total()
			dataAffiliateCom.Amount = rewardMeta.AffiliateCommissionAmount
			dataLeaderCom.Amount = rewardMeta.LeaderCommissionAmount

			profitItems := make([]meta.O, len(rewardMeta.DailyProfitAmounts))
			for i, profitAmount := range rewardMeta.DailyProfitAmounts {
				profitItems[i] = meta.O{
					"dailyprofit": profitAmount.Value,
					"coin_name":   profitAmount.Currency,
				}
			}
			data.Details = meta.O{
				"items": profitItems,
			}

			if !dataAffiliateCom.Amount.IsZero() {
				dataList = append(dataList, &dataAffiliateCom)
			}
			if !dataLeaderCom.Amount.IsZero() {
				dataList = append(dataList, &dataLeaderCom)
			}

		default:
			dataList = dataList[1:]
			return
		}
	}

	return dataList, nil
}

func LegacyCurrencyPriceMapGet(c echo.Context) (err error) {
	var (
		ctx = apiutils.EchoWrapContext(c)
	)
	priceTable := map[meta.Currency]map[meta.Currency]decimal.Decimal{
		constants.CurrencyUSD:       make(map[meta.Currency]decimal.Decimal),
		constants.CurrencyTetherUSD: make(map[meta.Currency]decimal.Decimal),
	}

	allCurrencyInfoMap := currencymod.GetTradingCurrencyInfoMapFastF()
	for currency, info := range allCurrencyInfoMap {
		priceTable[currency] = map[meta.Currency]decimal.Decimal{
			currency:                    comutils.DecimalOne,
			constants.CurrencyUSD:       info.PriceUSD,
			constants.CurrencyTetherUSD: info.PriceUSDT,
		}
		priceTable[constants.CurrencyUSD][currency] = info.PriceUSD
		priceTable[constants.CurrencyTetherUSD][currency] = info.PriceUSDT
	}

	for currency, priceMap := range priceTable {
		for currencyInner, priceMapInner := range priceTable {
			if currencyInner == currency {
				continue
			}
			if _, ok := priceMap[currencyInner]; ok {
				continue
			}
			var (
				priceUSD      = priceMapInner[constants.CurrencyUSD]
				priceCurrency decimal.Decimal
			)
			if priceUSD.IsZero() {
				priceCurrency = decimal.Zero
			} else {
				priceCurrency = priceMap[constants.CurrencyUSD].Div(priceUSD)
			}
			priceMap[currencyInner] = utils.NormalizeTradingAmount(priceCurrency)
		}
	}

	responsePriceTable := make(LegacyCurrencyPriceMapGetResponse)
	for currency, priceMap := range priceTable {
		responsePriceMap := make(map[meta.Currency]float64)
		for currency, price := range priceMap {
			responsePriceMap[currency], _ = price.Float64()
		}
		responsePriceTable[currency] = responsePriceMap
	}

	return responses.Ok(ctx, responsePriceTable)
}
