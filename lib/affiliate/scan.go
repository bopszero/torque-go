package affiliate

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/trading/tradingbalance"
)

type ScanOptions struct {
	RootUID            meta.UID `json:"root_uid"`
	LimitLevel         uint     `json:"limit_level"`
	GetCoinMap         bool     `json:"get_coin_map"`
	GetChildrenCoinMap bool     `json:"get_children_coin_map"`
	UseRawCoinMap      bool     `json:"use_raw_coin_map"`
	FetchRootOnly      bool     `json:"fetch_root_only"`
}

type ScanNodeInfo struct {
	User             TreeUser `json:"user"`
	DescendantsCount uint     `json:"descendants_count"`
	Level            uint     `json:"level"`

	CoinBalanceMap         tradingbalance.CoinBalanceMap `json:"coin_balance_map,omitempty"`
	BalanceUSD             *decimal.Decimal              `json:"balance_usd,omitempty"`
	ChildrenCoinBalanceMap tradingbalance.CoinBalanceMap `json:"children_coin_balance_map,omitempty"`
	ChildrenBalanceUSD     *decimal.Decimal              `json:"children_balance_usd,omitempty"`

	Children []*ScanNodeInfo `json:"children"`

	// Deprecated
	UserID   meta.UID `json:"user_id"`
	Username string   `json:"username"`
	UserTier string   `json:"user_role"`
}
