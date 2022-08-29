package affiliate

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/rewardmod"
)

type TreeNodeTier struct {
	Meta      meta.TierMeta
	IsFixed   bool
	BaseNodes []*TreeNode
}

type TreeNodeStats struct {
	node              *TreeNode
	profitableBalance *meta.CurrencyAmount
	maxChainLength    int
	firstTierNodeMap  map[string]*TreeNode
}

func (tns *TreeNodeStats) GetNodeByTier(tier string) *TreeNode {
	if tns.firstTierNodeMap == nil {
		tns.firstTierNodeMap = make(map[string]*TreeNode)
	}

	node, ok := tns.firstTierNodeMap[tier]
	if !ok {
		node = tns.findFirstNodeByTier(tier)
		tns.firstTierNodeMap[tier] = node
	}

	return node
}

func (tns *TreeNodeStats) findFirstNodeByTier(tier string) *TreeNode {
	if tns.node.GetTierCode() == tier {
		return tns.node
	}

	var node *TreeNode
	for _, child := range tns.node.Children {
		node = child.GetStats().GetNodeByTier(tier)
		if node != nil {
			return node
		}
	}

	return nil
}

func (tns *TreeNodeStats) GetMaxChainLength() int {
	if !tns.HasProfitableBalance() {
		return 0
	}

	if tns.maxChainLength == 0 {
		maxChildLength := 0
		for _, child := range tns.node.Children {
			childLength := child.GetStats().GetMaxChainLength()
			if childLength > maxChildLength {
				maxChildLength = childLength
			}
		}

		tns.maxChainLength = maxChildLength + 1
	}

	return tns.maxChainLength
}

func (tns *TreeNodeStats) HasProfitableBalance() bool {
	if tns.profitableBalance == nil {
		minThresholdMap := rewardmod.GetCurrencyDepositMinThresholdMap()
		for currency, balance := range tns.node.CoinBalanceMap {
			threshold, ok := minThresholdMap[currency]
			if !ok {
				panic(constants.ErrorInvalidData)
			}

			if balance.GreaterThanOrEqual(threshold) {
				balanceAmount := meta.CurrencyAmount{Currency: currency, Value: balance}
				tns.profitableBalance = &balanceAmount
				break
			}
		}
		if tns.profitableBalance == nil {
			anyCoinCode := constants.CurrencyEthereum
			dummyAmount := meta.CurrencyAmount{Currency: anyCoinCode, Value: decimal.Zero}
			tns.profitableBalance = &dummyAmount
		}
	}
	return tns.profitableBalance.Value.GreaterThan(decimal.Zero)
}
