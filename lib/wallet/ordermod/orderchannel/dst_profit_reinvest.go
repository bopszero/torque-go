package orderchannel

import (
	"fmt"
	"reflect"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/trading/depositmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/usermod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/balancemod"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/ordermod"
	"gorm.io/gorm"
)

const (
	DstProfitReinvestTomorrowOffsetDuration = 7 * time.Hour
	DstProfitReinvestRecentTimeGap          = 30 * 24 * time.Hour // ~1 month
	DstProfitReinvestRecentUserIdMaxCount   = 3
)

func init() {
	comutils.PanicOnError(
		ordermod.ChannelRegister(&DstProfitReinvestChannel{}),
	)
}

type DstProfitReinvestChannel struct {
	baseChannel
}

type DstProfitReinvestMeta struct {
	UserIdentity         string          `json:"user_identity" validate:"printascii"`
	Currency             meta.Currency   `json:"currency" validate:"required"`
	FromCoinExchangeRate decimal.Decimal `json:"exchange_rate"`
	ReceiveAmount        decimal.Decimal `json:"receive_amount"`

	// Deprecated
	Address string `json:"address"`
}

type DstProfitReinvestCheckoutInfo struct {
	ConversionRate       meta.CurrencyConversionRate `json:"conversion_rate"`
	RecentUserIdentities []string                    `json:"recent_user_identities"`
}

func (this *DstProfitReinvestChannel) GetType() meta.ChannelType {
	return constants.ChannelTypeDstProfitReinvest
}

func (this *DstProfitReinvestChannel) GetMetaType() reflect.Type {
	return reflect.TypeOf(DstProfitReinvestMeta{})
}

func (this *DstProfitReinvestChannel) getMeta(order *models.Order) (*DstProfitReinvestMeta, error) {
	var metaModel DstProfitReinvestMeta
	if err := ordermod.GetOrderChannelMetaData(order, this.GetType(), &metaModel); err != nil {
		return nil, err
	}
	return &metaModel, nil
}

func (this *DstProfitReinvestChannel) getTargetUser(order *models.Order) (_ models.User, err error) {
	metaModel, err := this.getMeta(order)
	if err != nil {
		return
	}
	coin, err := blockchainmod.GetCoinNative(metaModel.Currency)
	if err != nil {
		return
	}

	var targetUID meta.UID
	switch {
	case metaModel.UserIdentity != "":
		user, getErr := getUserByIdentity(metaModel.UserIdentity)
		if getErr != nil {
			err = getErr
			return
		}
		targetUID = user.ID
		break
	case metaModel.Address != "":
		if _, err = coin.NormalizeAddress(metaModel.Address); err != nil {
			err = utils.WrapError(constants.ErrorAddress)
			return
		}
		targetAddress, addrErr := depositmod.GetDepositAddressByAddress(coin, metaModel.Address)
		if addrErr != nil {
			if utils.IsOurError(addrErr, constants.ErrorCodeDataNotFound) {
				addrErr = utils.WrapError(constants.ErrorAddress)
			}
			err = addrErr
			return
		}
		targetUID = targetAddress.UID
		break
	default:
		targetUID = order.UID
		break
	}

	targetUser, err := usermod.GetUserFast(targetUID)
	if err != nil {
		return
	}
	if targetUser.Status != constants.UserStatusActive {
		err = utils.WrapError(constants.ErrorUserNotFound)
		return
	}
	metaModel.UserIdentity = targetUser.Username
	if err = ordermod.SetOrderChannelMetaData(order, this.GetType(), metaModel); err != nil {
		return
	}
	return targetUser, nil
}

func (this *DstProfitReinvestChannel) getRecentUserIdentities(order *models.Order) (usernames []string, err error) {
	var (
		orders         []models.Order
		now            = time.Now()
		recentFromTime = now.Add(-DstProfitReinvestRecentTimeGap).Unix()
	)
	err = database.GetDbF(database.AliasWalletSlave).
		Where(dbquery.Gte(models.OrderColCreateTime, recentFromTime)).
		Order(dbquery.OrderDesc(models.OrderColCreateTime)).
		Find(
			&orders,
			&models.Order{
				Currency:       order.Currency,
				UID:            order.UID,
				Status:         constants.OrderStatusCompleted,
				SrcChannelType: order.SrcChannelType,
				DstChannelType: order.DstChannelType,
			},
		).
		Error
	if err != nil {
		err = utils.WrapError(err)
		return
	}

	var (
		user        = usermod.GetUserFastF(order.UID)
		usernameSet = make(comtypes.HashSet, DstProfitReinvestRecentUserIdMaxCount)
	)
	for i := range orders {
		var (
			order        = orders[i]
			reinvestMeta DstProfitReinvestMeta
		)
		if err = ordermod.GetOrderChannelMetaData(&order, order.DstChannelType, &reinvestMeta); err != nil {
			return
		}
		if reinvestMeta.UserIdentity == "" || utils.IsSameStringCI(reinvestMeta.UserIdentity, user.Username) {
			continue
		}
		if !usernameSet.Add(reinvestMeta.UserIdentity) {
			continue
		}

		usernames = append(usernames, reinvestMeta.UserIdentity)
		if len(usernameSet) >= DstProfitReinvestRecentUserIdMaxCount {
			break
		}
	}
	return
}

func (this *DstProfitReinvestChannel) GetCheckoutInfo(ctx comcontext.Context, order *models.Order) (
	interface{}, error,
) {
	metaModel, err := this.getMeta(order)
	if err != nil {
		return nil, err
	}
	priceUsdt := currencymod.GetCurrencyPriceUsdtFastF(metaModel.Currency)
	priceMarkup := this.getCurrencyPriceMarkup(metaModel.Currency)
	if priceMarkup != nil {
		priceUsdt = priceMarkup.For(priceUsdt)
	}

	recentUsernames, err := this.getRecentUserIdentities(order)
	if err != nil {
		return nil, err
	}

	checkoutInfo := DstProfitReinvestCheckoutInfo{
		ConversionRate: meta.CurrencyConversionRate{
			FromCurrency: metaModel.Currency,
			ToCurrency:   constants.CurrencyTorque,
			Value:        currencymod.ConvertUsdtToTorque(priceUsdt),
		},
		RecentUserIdentities: recentUsernames,
	}

	return &checkoutInfo, nil
}

func (this *DstProfitReinvestChannel) Init(ctx comcontext.Context, order *models.Order) (err error) {
	metaModel, err := this.getMeta(order)
	if err != nil {
		return
	}
	coin, err := blockchainmod.GetCoinNative(metaModel.Currency)
	if err != nil {
		return
	}
	targetUser, err := this.getTargetUser(order)
	if err != nil {
		return
	}
	if _, err = depositmod.GetOrCreateDepositUserAddress(ctx, coin, targetUser.ID); err != nil {
		return
	}

	return nil
}

func (this *DstProfitReinvestChannel) PreValidate(ctx comcontext.Context, order *models.Order) error {
	if order.Currency != constants.CurrencyTorque {
		return utils.IssueErrorf("profit reinvest only accept Torque currency, not `%v` currency", order.Currency)
	}
	if ordermod.IsUserActionLocked(ctx, order.UID, constants.LockActionTORQReallocate) {
		return utils.WrapError(constants.ErrorUserActionLocked)
	}
	return nil
}

func (this *DstProfitReinvestChannel) Execute(ctx comcontext.Context, order *models.Order) (
	meta.OrderStepResultCode, error,
) {
	metaModel, err := this.getMeta(order)
	if err != nil {
		return ordermod.OrderStepResultCodeFail, err
	}
	targetUser, err := this.getTargetUser(order)
	if err != nil {
		return ordermod.OrderStepResultCodeFail, err
	}

	var torqueTxn models.TorqueTxn
	err = database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		torqueTxn, err = balancemod.SubmitProfitReinvest(
			ctx,
			order.UID, order.AmountSubTotal, metaModel.FromCoinExchangeRate,
			metaModel.Currency, targetUser.ID,
			this.getCurrencyPriceMarkup(metaModel.Currency))
		if err != nil {
			return err
		}

		metaModel.FromCoinExchangeRate = torqueTxn.ExchangeRate
		metaModel.ReceiveAmount = torqueTxn.CoinAmount
		return ordermod.SetOrderChannelMetaData(order, this.GetType(), metaModel)
	})
	if err != nil {
		return ordermod.OrderStepResultCodeFail, err
	}

	order.DstChannelID = torqueTxn.ID

	return ordermod.OrderStepResultCodeSuccess, err
}

func (this *DstProfitReinvestChannel) getCurrencyPriceMarkup(currency meta.Currency) *meta.AmountMarkup {
	if this.infoModel.BlockchainNetworkConfig.CurrencyMarkupPriceMap == nil {
		return nil
	}

	priceMarkup, ok := this.infoModel.BlockchainNetworkConfig.CurrencyMarkupPriceMap[currency]
	if !ok {
		return nil
	}

	return &priceMarkup
}

func (this *DstProfitReinvestChannel) Commit(ctx comcontext.Context, order *models.Order) (
	meta.OrderStepResultCode, error,
) {
	torqueTxn, err := this.getTorqueTxn(ctx, order)
	if err != nil {
		return ordermod.OrderStepResultCodeRetry, err
	}

	switch torqueTxn.Status {
	case constants.JudgeStatusApproved:
		return ordermod.OrderStepResultCodeSuccess, nil
	case constants.JudgeStatusRejected:
		return ordermod.OrderStepResultCodeFail, fmt.Errorf("profit reinvest fail with reason '%s'", torqueTxn.Note)
	default:
		return this.autoApprove(ctx, order, torqueTxn)
	}
}

func (this *DstProfitReinvestChannel) autoApprove(ctx comcontext.Context, order *models.Order, torqueTxn models.TorqueTxn) (
	meta.OrderStepResultCode, error,
) {
	if order.DstChannelID != torqueTxn.ID {
		panic(fmt.Errorf("reinvest cannot approve a mismatch ID Torque transaction"))
	}

	var (
		now               = time.Now()
		orderTomorrowTime = comutils.TimeRoundDate(time.Unix(order.CreateTime, 0)).Add(24 * time.Hour)
		orderApprovalTime = orderTomorrowTime.Add(DstProfitReinvestTomorrowOffsetDuration)
	)
	if now.Before(orderApprovalTime) {
		order.RetryTime = now.Add(3 * time.Minute).Unix()
		return ordermod.OrderStepResultCodeRetry, nil
	}

	if _, err := balancemod.ApproveProfitReinvest(ctx, torqueTxn.ID, ""); err != nil {
		if utils.IsOurError(err, constants.ErrorCodeUserNotFound) {
			order.Note = constants.ErrorUserNotFound.Message(ctx)
			if _, err := balancemod.RejectProfitReinvest(ctx, torqueTxn.ID, err.Error()); err != nil {
				return ordermod.OrderStepResultCodeRetry, err
			}
			return ordermod.OrderStepResultCodeFail, err
		}
		return ordermod.OrderStepResultCodeRetry, err
	}

	return ordermod.OrderStepResultCodeSuccess, nil
}

func (this *DstProfitReinvestChannel) getTorqueTxn(ctx comcontext.Context, order *models.Order) (
	txn models.TorqueTxn, err error,
) {
	err = database.TransactionFromContext(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) error {
		return balancemod.GetBalanceDB(dbTxn).
			FilterProfitReinvest(order.DstChannelID).
			First(&txn).
			Error
	})
	if err != nil {
		err = utils.WrapError(err)
	}
	return
}

func (this *DstProfitReinvestChannel) getNotification(
	ctx comcontext.Context, order models.Order, titleKey string, messageKey string,
) (*ordermod.Notification, error) {
	metaModel, err := this.getMeta(&order)
	if err != nil {
		return nil, err
	}

	return this.getSimpleNotification(
		ctx, order,
		titleKey, messageKey,
		meta.O{
			"order": order,
			"meta":  metaModel,
		},
	)
}

func (this *DstProfitReinvestChannel) GetNotificationCompleted(ctx comcontext.Context, order models.Order) (
	*ordermod.Notification, error,
) {
	return this.getNotification(
		ctx, order,
		constants.TranslationKeyOrderNotifCompletedTitleDstProfitReinvest,
		constants.TranslationKeyOrderNotifCompletedMessageDstProfitReinvest,
	)
}

func (this *DstProfitReinvestChannel) GetNotificationFailed(ctx comcontext.Context, order models.Order) (
	*ordermod.Notification, error,
) {
	return this.getNotification(
		ctx, order,
		constants.TranslationKeyOrderNotifFailedTitleDstProfitReinvest,
		constants.TranslationKeyOrderNotifFailedMessageDstProfitReinvest,
	)
}
