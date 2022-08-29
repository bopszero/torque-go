package ordermod

import (
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
)

type SrcPrepareOrderStep struct{}

func (this *SrcPrepareOrderStep) GetCode() string {
	return "sp"
}

func (this *SrcPrepareOrderStep) Forward(ctx comcontext.Context, order *models.Order) OrderStepResult {
	srcChannel, err := GetOrderSourceChannel(*order)
	comutils.PanicOnError(err)

	resultCode, err := srcChannel.Prepare(ctx, order)
	return newOrderStepResult(resultCode, err)
}

func (this *SrcPrepareOrderStep) Backward(ctx comcontext.Context, order *models.Order) OrderStepResult {
	srcChannel, err := GetOrderSourceChannel(*order)
	comutils.PanicOnError(err)

	resultCode, err := srcChannel.PrepareReverse(ctx, order)
	return newOrderStepResult(resultCode, err)
}

type SrcExecuteOrderStep struct{}

func (this *SrcExecuteOrderStep) GetCode() string {
	return "se"
}

func (this *SrcExecuteOrderStep) Forward(ctx comcontext.Context, order *models.Order) OrderStepResult {
	srcChannel, err := GetOrderSourceChannel(*order)
	comutils.PanicOnError(err)

	resultCode, err := srcChannel.Execute(ctx, order)
	return newOrderStepResult(resultCode, err)
}

func (this *SrcExecuteOrderStep) Backward(ctx comcontext.Context, order *models.Order) OrderStepResult {
	srcChannel, err := GetOrderSourceChannel(*order)
	comutils.PanicOnError(err)

	resultCode, err := srcChannel.ExecuteReverse(ctx, order)
	return newOrderStepResult(resultCode, err)
}

type SrcCommitOrderStep struct{}

func (this *SrcCommitOrderStep) GetCode() string {
	return "sc"
}

func (this *SrcCommitOrderStep) Forward(ctx comcontext.Context, order *models.Order) OrderStepResult {
	srcChannel, err := GetOrderSourceChannel(*order)
	comutils.PanicOnError(err)

	resultCode, err := srcChannel.Commit(ctx, order)
	return newOrderStepResult(resultCode, err)
}

func (this *SrcCommitOrderStep) Backward(ctx comcontext.Context, order *models.Order) OrderStepResult {
	srcChannel, err := GetOrderSourceChannel(*order)
	comutils.PanicOnError(err)

	resultCode, err := srcChannel.CommitReverse(ctx, order)
	return newOrderStepResult(resultCode, err)
}
