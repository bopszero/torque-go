package ordermod

import (
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
)

type DstPrepareOrderStep struct{}

func (this *DstPrepareOrderStep) GetCode() string {
	return "dp"
}

func (this *DstPrepareOrderStep) Forward(ctx comcontext.Context, order *models.Order) OrderStepResult {
	dstChannel, err := GetOrderDestinationChannel(*order)
	comutils.PanicOnError(err)

	resultCode, err := dstChannel.Prepare(ctx, order)
	return newOrderStepResult(resultCode, err)
}

func (this *DstPrepareOrderStep) Backward(ctx comcontext.Context, order *models.Order) OrderStepResult {
	dstChannel, err := GetOrderDestinationChannel(*order)
	comutils.PanicOnError(err)

	resultCode, err := dstChannel.PrepareReverse(ctx, order)
	return newOrderStepResult(resultCode, err)
}

type DstExecuteOrderStep struct{}

func (this *DstExecuteOrderStep) GetCode() string {
	return "de"
}

func (this *DstExecuteOrderStep) Forward(ctx comcontext.Context, order *models.Order) OrderStepResult {
	dstChannel, err := GetOrderDestinationChannel(*order)
	comutils.PanicOnError(err)

	resultCode, err := dstChannel.Execute(ctx, order)
	return newOrderStepResult(resultCode, err)
}

func (this *DstExecuteOrderStep) Backward(ctx comcontext.Context, order *models.Order) OrderStepResult {
	dstChannel, err := GetOrderDestinationChannel(*order)
	comutils.PanicOnError(err)

	resultCode, err := dstChannel.ExecuteReverse(ctx, order)
	return newOrderStepResult(resultCode, err)
}

type DstCommitOrderStep struct{}

func (this *DstCommitOrderStep) GetCode() string {
	return "dc"
}

func (this *DstCommitOrderStep) Forward(ctx comcontext.Context, order *models.Order) OrderStepResult {
	dstChannel, err := GetOrderDestinationChannel(*order)
	comutils.PanicOnError(err)

	resultCode, err := dstChannel.Commit(ctx, order)
	return newOrderStepResult(resultCode, err)
}

func (this *DstCommitOrderStep) Backward(ctx comcontext.Context, order *models.Order) OrderStepResult {
	dstChannel, err := GetOrderDestinationChannel(*order)
	comutils.PanicOnError(err)

	resultCode, err := dstChannel.CommitReverse(ctx, order)
	return newOrderStepResult(resultCode, err)
}
