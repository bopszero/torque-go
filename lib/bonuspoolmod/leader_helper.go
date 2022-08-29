package bonuspoolmod

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/affiliate"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/rewardmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/thirdpartymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/ordermod"
)

func leaderValidateDateRange(fromDate, toDate string) (fromTime time.Time, toTime time.Time, err error) {
	if fromTime, err = utils.TimeParseDate(fromDate); err != nil {
		return
	}
	if toTime, err = utils.TimeParseDate(toDate); err != nil {
		return
	}
	if fromTime.After(toTime) {
		err = utils.IssueErrorf(
			"from date cannot be after to date | from_date=%v,to_date=%v",
			fromDate, toDate)
		return
	}
	return fromTime, toTime, nil
}

func leaderGetTierDefaultRate(tierMeta meta.TierMeta) (rate decimal.Decimal, err error) {
	rate, ok := LeaderGetConfig().DefaultTierRateMap[tierMeta.Code]
	if !ok {
		err = utils.IssueErrorf("leader bonus pool doesn't contains rate for Tier %v", tierMeta.Code)
		return
	}
	return rate, nil
}

func leaderGenOrder(
	execution models.LeaderBonusPoolExecution, detail models.LeaderBonusPoolDetail,
) (order models.Order, err error) {
	var (
		now           = time.Now()
		tierRate      decimal.Decimal
		tierUserCount uint16
	)
	switch detail.TierType {
	case constants.TierMetaSeniorPartner.ID:
		tierRate = execution.TierSeniorRate
		tierUserCount = execution.TierSeniorReceiverCount
		break
	case constants.TierMetaRegionalPartner.ID:
		tierRate = execution.TierRegionalRate
		tierUserCount = execution.TierRegionalReceiverCount
		break
	case constants.TierMetaGlobalPartner.ID:
		tierRate = execution.TierGlobalRate
		tierUserCount = execution.TierGlobalReceiverCount
		break
	default:
		err = utils.WrapError(constants.ErrorInvalidData)
		return
	}

	order = ordermod.NewUserOrder(
		detail.UID, constants.CurrencyTorque,
		constants.ChannelTypeSrcBonusPoolLeader, constants.ChannelTypeDstBalance,
	)
	orderAmount := currencymod.NormalizeAmount(
		order.Currency,
		comutils.DecimalDivide(execution.TotalAmount.Mul(tierRate), decimal.NewFromInt(int64(tierUserCount))),
	)
	if orderAmount.IsZero() {
		err = utils.WrapError(constants.ErrorAmount)
		return
	}

	order.SrcChannelID = detail.ID
	order.SrcChannelAmount = orderAmount
	order.DstChannelAmount = orderAmount
	order.AmountSubTotal = orderAmount
	order.AmountTotal = orderAmount
	order.Status = constants.OrderStatusHandleSrc
	order.StepsData = models.OrderStepsData{
		History: []models.OrderStep{
			{
				Direction: constants.OrderStepDirectionForward,
				Code:      constants.OrderStepCodeOrderInit,
				Time:      now.Unix(),
			},
			{
				Direction: constants.OrderStepDirectionForward,
				Code:      constants.OrderStepCodeOrderStartSrc,
				Time:      now.Unix(),
			},
		},
	}
	return order, nil
}

func filterValidTierUIDs(
	ctx comcontext.Context,
	tierMeta meta.TierMeta, uids []meta.UID,
) (_ []meta.UID, err error) {
	var (
		treeClient      = thirdpartymod.GetTreeServiceSystemClient()
		treeScanOptions = affiliate.ScanOptions{
			LimitLevel: 1, // Direct children
		}
	)
	canRewardUIDs, err := rewardmod.FilterValidUIDs(uids)
	if err != nil {
		return
	}
	validUIDs := make([]meta.UID, 0, len(canRewardUIDs))
	for _, uid := range canRewardUIDs {
		nodeInfo, nodeErr := treeClient.GetNodeDown(ctx, uid, treeScanOptions)
		if nodeErr != nil {
			err = nodeErr
			return
		}
		directSameLevelCount := 0
		for _, childNode := range nodeInfo.Children {
			if childNode.User.TierType == nodeInfo.User.TierType {
				directSameLevelCount++
			}
		}
		if directSameLevelCount < LeaderDirectSameLevelDownlinePoolMinThreshold {
			continue
		}
		validUIDs = append(validUIDs, uid)
	}
	return validUIDs, nil
}
