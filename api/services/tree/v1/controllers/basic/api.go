package basic

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/torque-go/api"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/api/responses"
	"gitlab.com/snap-clickstaff/torque-go/api/services/tree/globals"
	"gitlab.com/snap-clickstaff/torque-go/lib/affiliate"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func GetNodeDown(c echo.Context) error {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel GetNodeDownRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}

	userNode := globals.GetTree().GetNode(reqModel.UID)
	if userNode == nil {
		return utils.WrapError(constants.ErrorDataNotFound)
	}
	if reqModel.Options.RootUID != 0 && reqModel.Options.RootUID != reqModel.UID {
		upperNode := userNode.Parent
		for upperNode != nil && upperNode.User.ID != reqModel.Options.RootUID {
			upperNode = upperNode.Parent
		}
		if upperNode == nil || upperNode.User.ID != reqModel.Options.RootUID {
			return utils.WrapError(constants.ErrorDataNotFound)
		}
	}
	var currencyInfoMap currencymod.CurrencyInfoMap
	if reqModel.Options.UseRawCoinMap {
		currencyInfoMap = currencymod.GetNonFiatCurrencyInfoMapFastF()
	} else {
		currencyInfoMap = currencymod.GetTradingCurrencyInfoMapFastF()
	}

	nodeInfo := affiliate.ScanNodeDown(ctx, userNode, reqModel.Options, 0, currencyInfoMap)
	if reqModel.Options.FetchRootOnly {
		nodeInfo.Children = nodeInfo.Children[:0]
	}

	return responses.Ok(
		ctx,
		GetNodeDownResponse{
			Node: nodeInfo,
		},
	)
}
