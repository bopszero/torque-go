package v1

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/torque-go/api"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/api/responses"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func HelperBlockchainValidateAddress(c echo.Context) error {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel HelperBlockchainValidateAddressRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return utils.WrapError(err)
	}
	coin, err := blockchainmod.GetCoin(reqModel.Currency, reqModel.Network)
	if err != nil {
		return err
	}
	var (
		mainAddress, mainErr     = coin.NormalizeAddress(reqModel.Address)
		legacyAddress, legacyErr = coin.NormalizeAddressLegacy(reqModel.Address)
		isValid                  = mainErr == nil && legacyErr == nil
	)
	response := HelperBlockchainValidateAddressResponse{
		IsValid: isValid,
		Address: mainAddress,
	}
	if isValid && legacyAddress != mainAddress {
		response.AddressLegacy = legacyAddress
	}

	return responses.Ok(ctx, response)
}
