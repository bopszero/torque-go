package ordermod

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/balancemod"
)

type InitOrderStep struct{}

func (this *InitOrderStep) GetCode() string {
	return constants.OrderStepCodeOrderInit
}

func (this *InitOrderStep) Forward(ctx comcontext.Context, order *models.Order) OrderStepResult {
	if order.Status != constants.OrderStatusNew {
		return newOrderStepResult(
			OrderStepResultCodeFail,
			utils.IssueErrorf(
				"cannot init order with invalid status | order_id=%v,expected=%v,actual=%v",
				order.ID, constants.OrderStatusNew, order.Status,
			),
		)
	}

	srcChannel, err := GetOrderSourceChannel(*order)
	if err != nil {
		return newOrderStepResult(OrderStepResultCodeFail, err)
	}
	dstChannel, err := GetOrderDestinationChannel(*order)
	if err != nil {
		return newOrderStepResult(OrderStepResultCodeFail, err)
	}

	if !srcChannel.IsAvailable() {
		return newOrderStepResult(OrderStepResultCodeFail, constants.ErrorChannelNotAvailable)
	}
	if !dstChannel.IsAvailable() {
		return newOrderStepResult(OrderStepResultCodeFail, constants.ErrorChannelNotAvailable)
	}

	if err := srcChannel.Init(ctx, order); err != nil {
		return newOrderStepResult(OrderStepResultCodeFail, err)
	}
	if err := dstChannel.Init(ctx, order); err != nil {
		return newOrderStepResult(OrderStepResultCodeFail, err)
	}

	if err := srcChannel.PreValidate(ctx, order); err != nil {
		return newOrderStepResult(OrderStepResultCodeFail, err)
	}
	if err := dstChannel.PreValidate(ctx, order); err != nil {
		return newOrderStepResult(OrderStepResultCodeFail, err)
	}

	srcInfoModel := srcChannel.GetInfoModel()
	dstInfoModel := dstChannel.GetInfoModel()
	if order.AmountSubTotal.LessThan(srcInfoModel.MinTxnAmount) ||
		order.AmountSubTotal.LessThan(dstInfoModel.MinTxnAmount) {
		return newOrderStepResult(
			OrderStepResultCodeFail,
			constants.ErrorAmountTooLowWithValue.WithData(meta.O{
				"threshold": decimal.Min(srcInfoModel.MinTxnAmount, dstInfoModel.MinTxnAmount),
				"currency":  constants.CurrencyUSD,
			}),
		)
	}
	if (!srcInfoModel.MaxTxnAmount.IsZero() && order.AmountSubTotal.GreaterThan(srcInfoModel.MaxTxnAmount)) ||
		(!dstInfoModel.MaxTxnAmount.IsZero() && order.AmountSubTotal.GreaterThan(dstInfoModel.MaxTxnAmount)) {
		return newOrderStepResult(
			OrderStepResultCodeFail,
			constants.ErrorAmountTooHighWithValue.WithData(meta.O{
				"threshold": decimal.Max(srcInfoModel.MaxTxnAmount, dstInfoModel.MaxTxnAmount),
				"currency":  constants.CurrencyUSD,
			}),
		)
	}

	currencyInfo, err := currencymod.GetCurrencyInfoFast(order.Currency)
	if err != nil || !currencymod.IsValidWalletInfo(currencyInfo) {
		return newOrderStepResult(OrderStepResultCodeFail, constants.ErrorCurrency)
	}
	for _, amount := range []decimal.Decimal{
		order.SrcChannelAmount,
		order.DstChannelAmount,
		order.AmountSubTotal,
		order.AmountFee,
		order.AmountDiscount,
		order.AmountTotal,
	} {
		if !amount.Equal(currencymod.NormalizeAmount(order.Currency, amount)) || amount.IsNegative() {
			return newOrderStepResult(OrderStepResultCodeFail, constants.ErrorAmount)
		}
	}

	// TODO: Validate amounts
	if !order.AmountTotal.Equal(order.AmountSubTotal.Add(order.AmountFee).Sub(order.AmountDiscount)) {
		return newOrderStepResult(OrderStepResultCodeFail, constants.ErrorAmount)
	}

	order.Status = constants.OrderStatusInit

	return OrderStepResultSuccess
}

func (this *InitOrderStep) Backward(ctx comcontext.Context, order *models.Order) OrderStepResult {
	if order.Status == constants.OrderStatusRefunding {
		order.Status = constants.OrderStatusRefunded
	} else if failStatus := getOrderFailStatus(*order); failStatus != constants.OrderStatusUnknown {
		order.Status = failStatus
	} else {
		order.Status = constants.OrderStatusFailed
	}

	return OrderStepResultSuccess
}

type SrcStartOrderStep struct{}

func (this *SrcStartOrderStep) GetCode() string {
	return constants.OrderStepCodeOrderStartSrc
}

func (this *SrcStartOrderStep) Forward(ctx comcontext.Context, order *models.Order) OrderStepResult {
	if order.Status != constants.OrderStatusInit {
		return newOrderStepResult(
			OrderStepResultCodeFail,
			fmt.Errorf(
				"cannot start order source channel with invalid status | order_id=%v,expected=%v,actual=%v",
				order.ID, constants.OrderStatusInit, order.Status,
			),
		)
	}

	order.Status = constants.OrderStatusHandleSrc

	return OrderStepResultSuccess
}

func (this *SrcStartOrderStep) Backward(ctx comcontext.Context, order *models.Order) OrderStepResult {
	return OrderStepResultIgnore
}

type SrcFinishOrderStep struct{}

func (this *SrcFinishOrderStep) GetCode() string {
	return constants.OrderStepCodeOrderFinishSrc
}

func (this *SrcFinishOrderStep) Forward(ctx comcontext.Context, order *models.Order) OrderStepResult {
	return OrderStepResultIgnore
}

func (this *SrcFinishOrderStep) Backward(ctx comcontext.Context, order *models.Order) OrderStepResult {
	order.Status = constants.OrderStatusRefunding

	return OrderStepResultSuccess
}

type DstStartOrderStep struct{}

func (this *DstStartOrderStep) GetCode() string {
	return constants.OrderStepCodeOrderStartDst
}

func (this *DstStartOrderStep) Forward(ctx comcontext.Context, order *models.Order) OrderStepResult {
	if order.Status != constants.OrderStatusHandleSrc {
		return newOrderStepResult(
			OrderStepResultCodeFail,
			fmt.Errorf(
				"cannot start order destination channel with invalid status | order_id=%v,expected=%v,actual=%v",
				order.ID, constants.OrderStatusHandleSrc, order.Status,
			),
		)
	}

	order.Status = constants.OrderStatusHandleDst

	return OrderStepResultSuccess
}

func (this *DstStartOrderStep) Backward(ctx comcontext.Context, order *models.Order) OrderStepResult {
	return OrderStepResultIgnore
}

type DstFinishOrderStep struct{}

func (this *DstFinishOrderStep) GetCode() string {
	return constants.OrderStepCodeOrderFinishDst
}

func (this *DstFinishOrderStep) Forward(ctx comcontext.Context, order *models.Order) OrderStepResult {
	if order.AmountFee.GreaterThan(decimal.Zero) {
		_, err := balancemod.AddTransaction(
			ctx,
			order.Currency, order.UID, order.AmountFee.Mul(constants.DecimalOneNegative),
			constants.WalletBalanceTypeMetaDrFee.ID, order.ID,
		)
		if err != nil {
			return newOrderStepResult(OrderStepResultCodeFail, err)
		}
	}

	order.Status = constants.OrderStatusCompleting

	return OrderStepResultSuccess
}

func (this *DstFinishOrderStep) Backward(ctx comcontext.Context, order *models.Order) OrderStepResult {
	if order.AmountFee.GreaterThan(decimal.Zero) {
		_, err := balancemod.AddTransaction(
			ctx,
			order.Currency, order.UID, order.AmountFee,
			constants.WalletBalanceTypeMetaCrFeeReverse.ID, order.ID,
		)
		if err != nil {
			return newOrderStepResult(OrderStepResultCodeRetry, err)
		}
	}

	return OrderStepResultSuccess
}

type DoneOrderStep struct{}

func (this *DoneOrderStep) GetCode() string {
	return constants.OrderStepCodeOrderDone
}

func (this *DoneOrderStep) Forward(ctx comcontext.Context, order *models.Order) OrderStepResult {
	if order.Status != constants.OrderStatusCompleting {
		return newOrderStepResult(
			OrderStepResultCodeFail,
			fmt.Errorf(
				"cannot finish order with invalid status | order_id=%v,expected=%v,actual=%v",
				order.ID, constants.OrderStatusCompleting, order.Status,
			),
		)
	}

	order.SucceedTime = time.Now().Unix()
	order.Status = constants.OrderStatusCompleted

	return OrderStepResultSuccess
}

func (this *DoneOrderStep) Backward(ctx comcontext.Context, order *models.Order) OrderStepResult {
	return OrderStepResultIgnore
}
