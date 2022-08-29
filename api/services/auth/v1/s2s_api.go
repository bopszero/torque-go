package v1

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/api"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/api/responses"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/auth/kycmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
)

func S2sKycGet(c echo.Context) (err error) {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel S2sKycGetRequest
	)
	if err = api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}
	var kycGetResponse KycGetResponse
	switch {
	case reqModel.UID != 0:
		kycGetResponse, err = kycGetByUID(ctx, reqModel.UID)
		if err != nil {
			return err
		}
		break
	case reqModel.Email != "":
		apiKycRequest, err := kycGetByEmail(ctx, reqModel.Email)
		if err != nil {
			return err
		}
		kycGetResponse.Request = apiKycRequest
	}
	return responses.Ok(ctx, S2sKycGetResponse{kycGetResponse})
}

func kycGetByEmail(ctx apiutils.EchoWrappedContext, email string) (*KycRequest, error) {
	modelKycRequest, err := kycmod.GetKycRequestByUserEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if modelKycRequest == nil {
		return nil, nil
	}
	var address kycmod.KycAddress
	if modelKycRequest.Address != "" {
		err = comutils.JsonDecode(modelKycRequest.Address, &address)
		if err != nil {
			return nil, err
		}
	}
	apiKycRequest := KycRequest{
		ID:                 modelKycRequest.ID,
		UID:                modelKycRequest.UID,
		Status:             modelKycRequest.Status,
		FullName:           modelKycRequest.FullName,
		DOB:                modelKycRequest.DOB.Format(constants.DateFormatISO),
		Nationality:        modelKycRequest.Nationality,
		ResidentialAddress: address.ResidentialAddress,
		PostalCode:         address.PostalCode,
		City:               address.City,
		Country:            address.Country,
	}
	return &apiKycRequest, nil
}

func S2sKycSendEmail(c echo.Context) (err error) {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel KycSendEmailRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}
	var kycRequest models.KycRequest
	err = database.GetDbF(database.AliasWalletSlave).First(&kycRequest, reqModel.RequestId).Error
	if err != nil {
		return err
	}
	if reqModel.RequestStatus != kycRequest.Status {
		return constants.ErrorKycRequestInvalidStatus
	}
	if reqModel.RequestStatus != constants.KycRequestStatusApproved &&
		reqModel.RequestStatus != constants.KycRequestStatusRejected {
		return constants.ErrorKycRequestInvalidStatus
	}
	isSent, err := kycmod.SendRequestEmail(ctx, kycRequest)
	if err != nil {
		return err
	}
	if !isSent {
		return constants.ErrorEmailNotSent
	}
	return responses.OkEmpty(ctx)
}
