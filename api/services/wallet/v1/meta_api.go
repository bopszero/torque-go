package v1

import (
	"fmt"
	"sort"
	"strings"

	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlocale"
	"gitlab.com/snap-clickstaff/torque-go/api"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/api/responses"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/auth/kycmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/notifmod"
)

func MetaHandshake(c echo.Context) (err error) {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		uid      = apiutils.GetContextUidF(ctx)
		reqModel = new(MetaHandshakeRequest)
	)
	if err = api.BindAndValidate(c, reqModel); err != nil {
		return
	}

	if reqModel.FirebaseToken != "" {
		userAgentHeader := c.Request().Header.Get(config.HttpHeaderUserAgent)
		_, err = notifmod.RegisterFirebaseToken(ctx, uid, reqModel.FirebaseToken, userAgentHeader)
		if err != nil {
			return
		}
	}

	metaFeatures, err := kycmod.GetMetaFeaturesSetting()
	if err != nil {
		return err
	}
	return responses.Ok(
		ctx,
		MetaHandshakeResponse{
			Features:           metaFeatures,
			BlockchainNetworks: metaGetBlockchainNetworks(),
			NetworkCurrencies:  metaGetNetworkCurrencies(),
		},
	)
}

func metaGetBlockchainNetworks() []MetaHandshakeBlockchainNetworkInfo {
	var (
		networkInfoMap = currencymod.GetAllBlockchainNetworkInfoMapFastF()
		infoList       = make([]MetaHandshakeBlockchainNetworkInfo, 0, len(networkInfoMap))
	)
	for network, info := range networkInfoMap {
		infoModel := MetaHandshakeBlockchainNetworkInfo{
			Network:               network,
			Currency:              info.Currency,
			Name:                  info.Name,
			TokenTransferCodeName: info.TokenTransferCodeName,
		}
		infoList = append(infoList, infoModel)
	}
	return infoList
}

func metaGetNetworkCurrencies() []MetaHandshakeNetworkCurrencyInfo {
	var (
		networkCurrencyInfoMap = currencymod.GetAllNetworkCurrencyInfoMapFastF()
		infoList               = make([]MetaHandshakeNetworkCurrencyInfo, 0, len(networkCurrencyInfoMap))
	)
	for _, info := range networkCurrencyInfoMap {
		if info.Priority == 0 {
			continue
		}
		infoModel := MetaHandshakeNetworkCurrencyInfo{
			Currency:      info.Currency,
			Network:       info.Network,
			Priority:      info.Priority,
			WithdrawalFee: info.WithdrawalFee,
		}
		infoList = append(infoList, infoModel)
	}
	return infoList
}

func MetaCurrencyInfoGet(c echo.Context) error {
	var (
		ctx = apiutils.EchoWrapContext(c)
	)
	return responses.Ok(
		ctx,
		MetaCurrencyInfoResponse{
			Currencies: metaGetCurrencyInfo(ctx),
		},
	)
}

func metaGetCurrencyInfo(ctx comcontext.Context) []MetaCurrencyInfo {
	var (
		currencyInfoMap = currencymod.GetAllCurrencyInfoMapFastF()
		currencies      = make([]MetaCurrencyInfo, 0, len(currencyInfoMap))
	)
	for _, info := range currencyInfoMap {
		var (
			msgKey = fmt.Sprintf(
				constants.TranslationKeyCurrencyStatusPattern,
				info.Currency.StringL())
			message, _ = comlocale.TranslateKey(ctx, msgKey)
		)
		currencies = append(
			currencies,
			MetaCurrencyInfo{
				CurrencyInfo:  info,
				StatusMessage: message,
			},
		)
	}
	sort.Slice(
		currencies,
		func(i, j int) bool {
			if currencies[i].PriorityDisplay != currencies[j].PriorityDisplay {
				return currencies[i].PriorityDisplay < currencies[j].PriorityDisplay
			}
			return strings.Compare(currencies[i].Currency.String(), currencies[j].Currency.String()) < 0
		},
	)
	return currencies
}
