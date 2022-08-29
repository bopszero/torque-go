package promotion

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/api"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/api/responses"
	"gitlab.com/snap-clickstaff/torque-go/api/services/tree/globals"
	"gitlab.com/snap-clickstaff/torque-go/lib/affiliate"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func GetStats(c echo.Context) error {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel = new(PromotionStatsGetRequest)
	)
	if err := api.BindAndValidate(c, reqModel); err != nil {
		return err
	}

	userNode := globals.GetTree().GetNode(reqModel.UID)
	if userNode == nil {
		return utils.WrapError(constants.ErrorDataNotFound)
	}

	downLineStatsList := make([]*PromotionStats, 0)

	seniorPartnerStats := NewPromotionStats(
		constants.TierMetaAgent.Code,
		constants.TierMetaSeniorPartner.Code)
	for _, child := range userNode.Children {
		if child.GetStats().GetMaxChainLength() >= affiliate.TierPromotionMarketLeaderDepth {
			seniorPartnerStats.UIDs = append(seniorPartnerStats.UIDs, child.User.ID)
		}
	}
	downLineStatsList = append(downLineStatsList, seniorPartnerStats)

	for _, tierMeta := range constants.TierMetaLeaderOrderedList {
		if tierMeta.Code == constants.TierMetaMentor.Code {
			continue
		}

		toTier, err := affiliate.GetUpperTierMeta(tierMeta.Code)
		comutils.PanicOnError(err)
		tierStats := NewPromotionStats(tierMeta.Code, toTier.Code)

		for _, child := range userNode.Children {
			tierNode := child.GetStats().GetNodeByTier(tierMeta.Code)
			if tierNode != nil {
				tierStats.UIDs = append(tierStats.UIDs, tierNode.User.ID)
			}
		}

		downLineStatsList = append(downLineStatsList, tierStats)
	}

	return responses.Ok(
		ctx,
		meta.O{
			"stats": PromotionStatsGetResponse{
				UID:               reqModel.UID,
				Tier:              userNode.GetTierCode(),
				DownLineStatsList: downLineStatsList,
			},
		},
	)
}
