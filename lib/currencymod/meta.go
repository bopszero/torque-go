package currencymod

import (
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type (
	CurrencyInfoMap       map[meta.Currency]models.CurrencyInfo
	LegacyCurrencyInfoMap map[meta.Currency]models.LegacyCurrencyInfo
)

const (
	PriorityCommonSoonThreshold = 10000
)
