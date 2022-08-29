package affiliate

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/trading/tradingbalance"
)

func genChildrenRawCoinMap(
	node *TreeNode, result *ScanNodeInfo,
	infoMap currencymod.CurrencyInfoMap,
) tradingbalance.CoinBalanceMap {
	if len(result.Children) == 0 {
		return nil
	}
	totalCoinMap := make(tradingbalance.CoinBalanceMap, len(infoMap))
	for _, childResult := range result.Children {
		for currency, balance := range childResult.CoinBalanceMap {
			totalCoinMap[currency] = totalCoinMap[currency].Add(balance)
		}
		for currency, balance := range childResult.ChildrenCoinBalanceMap {
			totalCoinMap[currency] = totalCoinMap[currency].Add(balance)
		}
	}
	return totalCoinMap
}

func genChildrenCoinMap(
	node *TreeNode, result *ScanNodeInfo,
	infoMap currencymod.CurrencyInfoMap,
) tradingbalance.CoinBalanceMap {
	totalCoinMap := make(tradingbalance.CoinBalanceMap, len(infoMap))
	for currency := range infoMap {
		var totalBalance decimal.Decimal
		for _, childResult := range result.Children {
			totalBalance = totalBalance.
				Add(childResult.CoinBalanceMap[currency]).
				Add(childResult.ChildrenCoinBalanceMap[currency])
		}
		totalCoinMap[currency] = totalBalance
	}
	return totalCoinMap
}
