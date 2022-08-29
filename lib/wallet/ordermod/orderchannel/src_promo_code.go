package orderchannel

import (
	"reflect"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/trading/tradingbalance"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/balancemod"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/ordermod"
)

func init() {
	comutils.PanicOnError(
		ordermod.ChannelRegister(&SrcPromoCodeChannel{}),
	)
}

type SrcPromoCodeChannel struct {
	baseChannel
}

type SrcPromoCodeMeta struct {
	Code string `json:"code" validate:"required,alphanum"`
}

func (this *SrcPromoCodeChannel) GetType() meta.ChannelType {
	return constants.ChannelTypeSrcPromoCode
}

func (this *SrcPromoCodeChannel) GetMetaType() reflect.Type {
	return reflect.TypeOf(SrcPromoCodeMeta{})
}

func (this *SrcPromoCodeChannel) getMeta(order *models.Order) (*SrcPromoCodeMeta, error) {
	var metaModel SrcPromoCodeMeta
	if err := ordermod.GetOrderChannelMetaData(order, this.GetType(), &metaModel); err != nil {
		return nil, err
	}

	return &metaModel, nil
}

func (this *SrcPromoCodeChannel) Init(ctx comcontext.Context, order *models.Order) error {
	order.SrcChannelAmount = decimal.Zero
	order.DstChannelAmount = balancemod.PromoCodeCreditValue

	order.AmountSubTotal = balancemod.PromoCodeCreditValue
	order.AmountFee = decimal.Zero
	order.AmountDiscount = decimal.Zero
	order.AmountTotal = balancemod.PromoCodeCreditValue

	return nil
}

func (this *SrcPromoCodeChannel) PreValidate(ctx comcontext.Context, order *models.Order) (err error) {
	if order.Currency != constants.CurrencyTorque {
		return utils.IssueErrorf("promo code redemption only accept Torque currency, not `%v` currency", order.Currency)
	}

	metaModel, err := this.getMeta(order)
	if err != nil {
		return
	}

	db := database.GetDbSlave()

	var promoCode models.PromoCode
	err = db.
		First(
			&promoCode,
			&models.PromoCode{
				Code:      metaModel.Code,
				Status:    constants.PromoCodeStatusAvailable,
				IsDeleted: models.NewBool(false),
			},
		).
		Error
	if database.IsDbError(err) {
		return
	}
	if promoCode.ID == 0 {
		return constants.ErrorDataNotFound
	}

	var redemption models.PromoCodeRedemption
	err = db.
		First(
			&redemption,
			&models.PromoCodeRedemption{UID: order.UID, IsDeleted: models.NewBool(false)},
		).
		Error
	if database.IsDbError(err) {
		return err
	}
	if redemption.ID != 0 {
		return constants.ErrorDataExists
	}

	portfolioValueUsd := tradingbalance.CalcUserCoinBalanceUsdValue(order.UID)
	if portfolioValueUsd.LessThan(balancemod.PromoCodeMinBalanceThreshold) {
		return constants.ErrorAmountTooLowWithValue.WithData(meta.O{
			"threshold": balancemod.PromoCodeMinBalanceThreshold,
			"currency":  constants.CurrencyUSD,
		})
	}

	return nil
}

func (this *SrcPromoCodeChannel) Execute(ctx comcontext.Context, order *models.Order) (
	_ meta.OrderStepResultCode, err error,
) {
	metaModel, err := this.getMeta(order)
	if err != nil {
		return ordermod.OrderStepResultCodeFail, err
	}

	codeRedemption, err := balancemod.RedeemPromoCode(ctx, order.UID, metaModel.Code)
	if err != nil {
		return ordermod.OrderStepResultCodeFail, err
	}

	order.SrcChannelID = codeRedemption.ID
	order.SrcChannelRef = metaModel.Code

	return ordermod.OrderStepResultCodeSuccess, nil
}
