package payment

import (
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/api"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/api/responses"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/trading/depositmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/trading/depositmod/depositcrawler"
	"gitlab.com/snap-clickstaff/torque-go/lib/trading/withdrawalmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func GetDepositAccount(c echo.Context) error {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel DepositGetAccountRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return utils.WrapError(err)
	}

	coin, err := blockchainmod.GetCoin(reqModel.Currency, reqModel.Network)
	if err != nil {
		return err
	}
	depositAddress, err := depositmod.GetOrCreateDepositUserAddress(ctx, coin, reqModel.UID)
	if err != nil {
		return err
	}
	address, err := coin.NormalizeAddress(depositAddress.Address)
	if err != nil {
		return err
	}

	var responseModel DepositGetAccountResponse
	if err := copier.Copy(&responseModel, &depositAddress); err != nil {
		return utils.WrapError(err)
	}

	legacyAddress, err := coin.NormalizeAddressLegacy(depositAddress.Address)
	if err == nil && legacyAddress != address {
		responseModel.AddressLegacy = legacyAddress
	}

	return responses.Ok(ctx, responseModel)
}

func CrawlDeposit(c echo.Context) (err error) {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel DepositCrawlRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}

	coin, err := blockchainmod.GetCoin(reqModel.Currency, reqModel.Network)
	if err != nil {
		return
	}
	client, err := coin.NewClientDefault()
	if err != nil {
		return
	}

	emptyOptions := depositcrawler.CrawlerOptions{}
	crawler, err := depositcrawler.NewCrawler(ctx, coin, emptyOptions)
	if err != nil {
		return
	}

	block, err := client.GetBlockByHeight(reqModel.BlockHeight)
	if err != nil {
		return
	}
	if err = crawler.ConsumeBlock(block); err != nil {
		return
	}

	return responses.OkEmpty(ctx)
}

func SubmitDeposit(c echo.Context) error {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel DepositSubmitRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}

	deposit, err := depositmod.SubmitDeposit(
		ctx,
		reqModel.UID, reqModel.Currency, reqModel.Network,
		reqModel.TxnHash, reqModel.TxnIndex, reqModel.Address, reqModel.Amount,
	)
	if err != nil {
		return err
	}

	return responses.Ok(
		ctx,
		DepositSubmitResponse{
			ID:       deposit.ID,
			CoinID:   deposit.CoinID,
			Currency: deposit.Currency,
			Network:  deposit.Network,
			UID:      deposit.UID,
			Status:   deposit.Status,

			TxnHash:  deposit.TxnHash,
			TxnIndex: deposit.TxnIndex,
			Address:  deposit.Address,
			Amount:   deposit.Amount,

			CreateTime: deposit.CreateTime.Unix(),
		},
	)
}

func ApproveDeposit(c echo.Context) error {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel DepositApproveRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}

	_, err := depositmod.ApproveDeposit(ctx, reqModel.ID, reqModel.Note)
	if err != nil {
		return err
	}

	return responses.OkEmpty(ctx)
}

func SubmitInvestmentWithdraw(c echo.Context) error {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel InvestmentWithdrawSubmitRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}

	withdrawModel, err := withdrawalmod.SubmitInvestmentWithdraw(
		ctx,
		reqModel.UID, reqModel.Amount,
		reqModel.Currency, reqModel.Network, reqModel.Address)
	if err != nil {
		return err
	}

	return responses.Ok(
		ctx,
		meta.O{
			"withdraw_request": InvestmentWithdrawSubmitResponse{
				ID:         withdrawModel.ID,
				Code:       withdrawModel.Code,
				UID:        withdrawModel.UserID,
				Amount:     withdrawModel.Amount,
				Currency:   withdrawModel.Currency,
				Network:    withdrawModel.Network,
				Address:    withdrawModel.Address,
				Status:     withdrawModel.Status,
				CreateTime: withdrawModel.CreateTime.Unix(),
			},
		},
	)
}

func RejectInvestmentWithdraw(c echo.Context) error {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel InvestmentWithdrawRejectRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}

	_, err := withdrawalmod.RejectInvestmentWithdraw(ctx, reqModel.ID, reqModel.Note)
	comutils.PanicOnError(err)

	return responses.Ok(ctx, meta.O{"ok": true})
}

func CancelInvestmentWithdraw(c echo.Context) error {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel InvestmentWithdrawCancelRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}

	_, err := withdrawalmod.CancelInvestmentWithdraw(ctx, reqModel.ID, reqModel.Note)
	comutils.PanicOnError(err)

	return responses.Ok(ctx, meta.O{"ok": true})
}
