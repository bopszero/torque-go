package payment

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/torque-go/api"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/api/responses"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/trading/tradingtxn"
)

func TransferP2P(c echo.Context) error {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel P2PTransferRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}

	p2pTransferModel, err := tradingtxn.TransferP2P(
		ctx,
		reqModel.FromUID, reqModel.FromCurrency, reqModel.FromAmount,
		reqModel.ToUID, reqModel.ToCurrency, reqModel.ExchangeRate, reqModel.FeeAmount, reqModel.Note)
	if err != nil {
		return err
	}

	return responses.Ok(
		ctx,
		meta.O{
			"transfer": P2PTransferResponse{
				ID:           p2pTransferModel.ID,
				FromUID:      p2pTransferModel.FromUID,
				FromCurrency: p2pTransferModel.FromCurrency,
				FromAmount:   p2pTransferModel.FromAmount,
				ToUID:        p2pTransferModel.ToUID,
				ToCurrency:   p2pTransferModel.ToCurrency,
				ToAmount:     p2pTransferModel.ToAmount,
				ExchangeRate: p2pTransferModel.ExchangeRate,
				FeeAmount:    p2pTransferModel.FeeAmount,
				Note:         p2pTransferModel.Note,
				Status:       p2pTransferModel.Status,
				CreateTime:   p2pTransferModel.CreateTime,
			},
		},
	)
}
