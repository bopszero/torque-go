package affiliate

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

func GetUpperTierMeta(tierCode string) (meta.TierMeta, error) {
	tierIdx, ok := constants.TierOrderCodeMap[tierCode]
	if !ok {
		return constants.TierMetaZero, constants.ErrorInvalidData
	}

	upperTierIdx := tierIdx + 1
	if upperTierIdx >= len(constants.TierMetaOrderedList) {
		return constants.TierMetaZero, constants.ErrorDataNotFound
	}

	upperTierMeta := constants.TierMetaOrderedList[upperTierIdx]
	return upperTierMeta, nil
}

func GetLowerTierMeta(tierCode string) (meta.TierMeta, error) {
	tierIdx, ok := constants.TierOrderCodeMap[tierCode]
	if !ok {
		return constants.TierMetaZero, constants.ErrorInvalidData
	}

	lowerTierIdx := tierIdx - 1
	if lowerTierIdx <= 0 {
		return constants.TierMetaZero, constants.ErrorDataNotFound
	}

	lowerTierMeta := constants.TierMetaOrderedList[lowerTierIdx]
	return lowerTierMeta, nil
}

func ScanNodeDown(
	ctx comcontext.Context,
	node *TreeNode, options ScanOptions, level uint,
	currencyInfoMap currencymod.CurrencyInfoMap,
) ScanNodeInfo {
	var (
		children         []*ScanNodeInfo
		descendantsCount uint
	)
	if options.LimitLevel == 0 || level < options.LimitLevel {
		if level <= MaxScanLevel {
			children = make([]*ScanNodeInfo, 0, len(node.Children))
			for _, childNode := range node.Children {
				childResult := ScanNodeDown(ctx, childNode, options, level+1, currencyInfoMap)
				children = append(children, &childResult)
				descendantsCount += 1 + childResult.DescendantsCount
			}
		} else {
			comlogging.GetLogger().
				WithContext(ctx).
				WithFields(logrus.Fields{
					"uid":     node.User.ID,
					"level":   level,
					"options": options,
				}).
				Error("scan node down too deep")
		}
	}

	var (
		userID   = node.User.ID
		username = node.User.Username
		userTier = "unknown"
	)
	if userTierMeta, ok := constants.TierMetaMap[node.User.TierType]; ok {
		userTier = userTierMeta.Code
	}
	result := ScanNodeInfo{
		User:             node.User,
		DescendantsCount: descendantsCount,
		Level:            level,
		Children:         children,

		// Deprecated
		UserID:   userID,
		Username: username,
		UserTier: userTier,
	}
	if options.GetCoinMap {
		if options.UseRawCoinMap {
			result.CoinBalanceMap = node.CoinBalanceMap
		} else {
			result.CoinBalanceMap = node.CoinBalanceMap.Format(currencyInfoMap)
		}
		balanceUSD := result.CoinBalanceMap.CalcValueUSD(currencyInfoMap)
		result.BalanceUSD = &balanceUSD
	}
	if options.GetChildrenCoinMap {
		if options.UseRawCoinMap {
			result.ChildrenCoinBalanceMap = genChildrenRawCoinMap(node, &result, currencyInfoMap)
		} else {
			result.ChildrenCoinBalanceMap = genChildrenCoinMap(node, &result, currencyInfoMap)
		}
		childrenBalanceUSD := result.ChildrenCoinBalanceMap.CalcValueUSD(currencyInfoMap)
		result.ChildrenBalanceUSD = &childrenBalanceUSD
	}
	return result
}
