package orderchannel

import (
	"reflect"

	"github.com/shopspring/decimal"
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
		ordermod.ChannelRegister(&DstBalanceChannel{}),
	)
}

type DstBalanceChannel struct {
	baseChannel
}

type DstBalanceMeta struct {
	InnerTxns []DstBalanceMetaInnerTxn `json:"inner_txns"`
}

type DstBalanceMetaInnerTxn struct {
	TxnType uint32          `json:"txn_type"`
	Value   decimal.Decimal `json:"value"`
}

func (this *DstBalanceChannel) GetType() meta.ChannelType {
	return constants.ChannelTypeDstBalance
}

func (this *DstBalanceChannel) GetMetaType() reflect.Type {
	return reflect.TypeOf(DstBalanceMeta{})
}

func (this *DstBalanceChannel) getMeta(order *models.Order) (*DstBalanceMeta, error) {
	var metaModel DstBalanceMeta

	err := ordermod.GetOrderChannelMetaData(order, this.GetType(), &metaModel)
	if err != nil && !utils.IsOurError(err, constants.ErrorCodeDataNotFound) {
		return nil, err
	}

	return &metaModel, nil
}

func (this *DstBalanceChannel) PreValidate(ctx comcontext.Context, order *models.Order) error {
	metaModel, err := this.getMeta(order)
	if err != nil {
		return err
	}

	if len(metaModel.InnerTxns) > 0 {
		totalInnerAmount := decimal.Zero
		for _, txn := range metaModel.InnerTxns {
			totalInnerAmount = totalInnerAmount.Add(txn.Value)
		}

		if !totalInnerAmount.Equal(order.AmountSubTotal) {
			return constants.ErrorAmount
		}
	}

	return nil
}

func (this *DstBalanceChannel) Execute(ctx comcontext.Context, order *models.Order) (
	meta.OrderStepResultCode, error,
) {
	metaModel, err := this.getMeta(order)
	if err != nil {
		return ordermod.OrderStepResultCodeFail, err
	}

	if len(metaModel.InnerTxns) == 0 {
		return this.executeSingle(ctx, order)
	} else {
		return this.executeMultiple(ctx, order)
	}
}

func (this *DstBalanceChannel) executeSingle(ctx comcontext.Context, order *models.Order) (
	meta.OrderStepResultCode, error,
) {
	balanceTypeMeta, err := ordermod.GetBalanceTypeMetaCrByChannel(order.SrcChannelType)
	if err != nil {
		return ordermod.OrderStepResultCodeRetry, err
	}

	_, err = balancemod.AddTransaction(
		ctx,
		order.Currency, order.UID, order.AmountSubTotal,
		balanceTypeMeta.ID, order.ID)
	if err != nil {
		return ordermod.OrderStepResultCodeRetry, err
	}

	return ordermod.OrderStepResultCodeSuccess, nil
}

func (this *DstBalanceChannel) executeMultiple(ctx comcontext.Context, order *models.Order) (
	meta.OrderStepResultCode, error,
) {
	metaModel, err := this.getMeta(order)
	if err != nil {
		return ordermod.OrderStepResultCodeFail, err
	}

	for _, txn := range metaModel.InnerTxns {
		if txn.Value.IsZero() {
			continue
		}

		_, err := balancemod.AddTransaction(
			ctx,
			order.Currency, order.UID, txn.Value,
			txn.TxnType, order.ID)
		if err != nil {
			return ordermod.OrderStepResultCodeRetry, err
		}
	}

	return ordermod.OrderStepResultCodeSuccess, nil
}

func (this *DstBalanceChannel) CommitReverse(ctx comcontext.Context, order *models.Order) (
	meta.OrderStepResultCode, error,
) {
	err := utils.IssueErrorf("cannot reverse a top-up order, kindly use adjust balance for it | order_id=%v", order.ID)
	return ordermod.OrderStepResultCodeFail, err
}
