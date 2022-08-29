package v1

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/api"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/api/responses"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/syswithdrawalmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func SystemWithdrawalAccountMeta(c echo.Context) error {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel WithdrawalAccountGenerateRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}

	sysAddress, err := syswithdrawalmod.GenerateAccount(ctx, reqModel.Currency, reqModel.Network)
	if err != nil {
		return err
	}

	return responses.Ok(
		ctx,
		WithdrawalAccountGenerateResponse{
			Currency:  reqModel.Currency,
			Network:   sysAddress.Network,
			AccountNo: sysAddress.Address,
		},
	)
}

func SystemWithdrawalAccountGenerate(c echo.Context) error {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel WithdrawalAccountGenerateRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}

	sysAddress, err := syswithdrawalmod.GenerateAccount(ctx, reqModel.Currency, reqModel.Network)
	if err != nil {
		return err
	}

	return responses.Ok(
		ctx,
		WithdrawalAccountGenerateResponse{
			Currency:  reqModel.Currency,
			Network:   sysAddress.Network,
			AccountNo: sysAddress.Address,
		},
	)
}

func SystemWithdrawalAccountGet(c echo.Context) error {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel WithdrawalAccountGetRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}

	coin, err := blockchainmod.GetCoin(reqModel.Currency, reqModel.Network)
	if err != nil {
		return err
	}
	withdrawalConfig := syswithdrawalmod.GetConfig(coin)
	sysAddress, err := syswithdrawalmod.GetAccount(ctx, coin)
	if err != nil {
		return err
	}

	return responses.Ok(
		ctx,
		WithdrawalAccountGetResponse{
			Currency:    reqModel.Currency,
			Network:     sysAddress.Network,
			AccountNo:   sysAddress.Address,
			PullAddress: withdrawalConfig.PullAddress,
		},
	)
}

func SystemWithdrawalAccountPull(c echo.Context) error {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel WithdrawalAccountPullRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}

	sysAddress, hash, err := syswithdrawalmod.PullAccount(
		ctx,
		reqModel.Currency, reqModel.Network,
		reqModel.AccountNo, reqModel.PullAddress)
	if err != nil {
		return err
	}

	return responses.Ok(
		ctx,
		WithdrawalAccountPullResponse{
			Currency: reqModel.Currency,
			Network:  sysAddress.Network,
			Hash:     hash,
		},
	)
}

func SystemWithdrawalTransferMeta(c echo.Context) error {
	var (
		ctx        = apiutils.EchoWrapContext(c)
		coinMap    = blockchainmod.GetCoinMap()
		currencies = make([]WithdrawalTransferMeta, 0, len(coinMap))
	)
	for index, coin := range coinMap {
		feeInfo, err := coin.GetDefaultFeeInfo()
		comutils.PanicOnError(err)

		currencies = append(
			currencies,
			WithdrawalTransferMeta{
				Currency: index.Currency,
				Network:  index.Network,
				FeeInfo:  feeInfo,
			},
		)
	}

	return responses.Ok(
		ctx,
		WithdrawalTransferMetaResponse{
			Currencies: currencies,
		},
	)
}

func SystemWithdrawalTransferSubmit(c echo.Context) error {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel WithdrawalTransferSubmitRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}

	coin, err := blockchainmod.GetCoin(reqModel.Currency, reqModel.Network)
	if err != nil {
		return err
	}
	request, err := syswithdrawalmod.SubmitBulk(
		ctx, reqModel.RequestUID,
		coin.GetCurrency(), coin.GetNetwork(), reqModel.SrcAddress,
		reqModel.Codes, reqModel.TotalAmount)
	if err != nil {
		return err
	}

	return responses.Ok(
		ctx,
		WithdrawalTransferSubmitResponse{
			RequestID: request.ID,
			AddressID: request.AddressID,

			Status:   request.Status,
			Currency: request.Currency,
			Network:  request.Network,
			Amount:   request.Amount,
			AmountEstimatedFee: meta.CurrencyAmount{
				Currency: coin.GetNetworkCurrency(),
				Value:    request.AmountEstimatedFee,
			},
			CombinedTxnHash: request.CombinedTxnHash.String,

			CreateUID:  request.CreateUID,
			CreateTime: request.CreateTime,
		},
	)
}

func SystemWithdrawalTransferConfirm(c echo.Context) error {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel WithdrawalTransferConfirmRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}

	request, err := syswithdrawalmod.ConfirmRequest(ctx, reqModel.RequestID)
	if err != nil {
		return err
	}

	execFunc := func() {
		var (
			err error
			ctx = comcontext.NewContext()
		)
		defer func() {
			utils.ErrorCatchWithLog(ctx, "system withdraw execute on confirm", err)
		}()
		_, err = syswithdrawalmod.ExecuteRequest(ctx, request.ID)
	}
	if config.Debug {
		execFunc()
	} else {
		go execFunc()
	}

	return responses.OkEmpty(ctx)
}

func SystemWithdrawalTransferReplace(c echo.Context) error {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel WithdrawalTransferReplaceRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}

	_, err := syswithdrawalmod.ReplaceRequest(ctx, reqModel.RequestID, reqModel.FeeInfo)
	if err != nil {
		return err
	}

	return responses.OkEmpty(ctx)
}
