package orderchannel

import (
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/balancemod"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/ordermod"
)

func init() {
	comutils.PanicOnError(
		ordermod.ChannelRegister(&SrcBalanceChannel{}),
	)
}

type SrcBalanceChannel struct {
	baseChannel
}

func (this *SrcBalanceChannel) GetType() meta.ChannelType {
	return constants.ChannelTypeSrcBalance
}

func (this *SrcBalanceChannel) PreValidate(ctx comcontext.Context, order *models.Order) error {
	if order.Currency != constants.CurrencyTorque {
		return utils.IssueErrorf("source balance only accept Torque currency, not `%v` currency", order.Currency)
	}

	userBalance, err := balancemod.GetUserBalance(ctx, order.UID, order.Currency)
	if err != nil {
		return err
	}
	if userBalance.Amount.LessThan(order.AmountTotal) {
		return constants.ErrorBalanceNotEnough
	}

	return nil
}

func (this *SrcBalanceChannel) Execute(ctx comcontext.Context, order *models.Order) (
	_ meta.OrderStepResultCode, err error,
) {
	balanceTypeMeta, err := ordermod.GetBalanceTypeMetaDrByChannel(order.DstChannelType)
	if err != nil {
		return ordermod.OrderStepResultCodeFail, err
	}
	_, err = balancemod.AddTransaction(
		ctx,
		order.Currency, order.UID, order.AmountSubTotal.Mul(constants.DecimalOneNegative),
		balanceTypeMeta.ID, order.ID)
	if err != nil {
		return ordermod.OrderStepResultCodeFail, err
	}

	return ordermod.OrderStepResultCodeSuccess, nil
}

func (this *SrcBalanceChannel) ExecuteReverse(ctx comcontext.Context, order *models.Order) (
	_ meta.OrderStepResultCode, err error,
) {
	balanceTypeMeta, err := ordermod.GetBalanceTypeMetaDrByChannel(order.DstChannelType)
	if err != nil {
		return ordermod.OrderStepResultCodeRetry, err
	}
	txn, err := balancemod.GetTransaction(
		ctx,
		order.Currency, order.UID,
		balanceTypeMeta.ID, order.ID)
	if err != nil {
		return ordermod.OrderStepResultCodeRetry, err
	}
	if txn == nil {
		return ordermod.OrderStepResultCodeIgnore, nil
	}

	refundTxn, err := balancemod.GetTransaction(
		ctx,
		order.Currency, order.UID,
		constants.WalletBalanceTypeMetaCrRefund.ID, order.ID)
	if err != nil {
		return ordermod.OrderStepResultCodeRetry, err
	}
	if refundTxn != nil {
		return ordermod.OrderStepResultCodeIgnore, nil
	}

	_, err = balancemod.AddTransaction(
		ctx,
		order.Currency, order.UID, order.AmountSubTotal,
		constants.WalletBalanceTypeMetaCrRefund.ID, order.ID)
	if err != nil {
		return ordermod.OrderStepResultCodeRetry, err
	}

	return ordermod.OrderStepResultCodeSuccess, nil
}
