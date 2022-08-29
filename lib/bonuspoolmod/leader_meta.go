package bonuspoolmod

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	LeaderDirectSameLevelDownlinePoolMinThreshold = 3
)

var (
	LeaderTierMetaList = []meta.TierMeta{
		constants.TierMetaSeniorPartner,
		constants.TierMetaRegionalPartner,
		constants.TierMetaGlobalPartner,
	}
)

type (
	LeaderExecutionHashMeta struct {
		Secret       string           `msgpack:"secret"`
		FromDate     string           `msgpack:"from_date"`
		ToDate       string           `msgpack:"to_date"`
		TotalAmount  decimal.Decimal  `msgpack:"total_amount"`
		TierInfoList []LeaderTierInfo `msgpack:"tier_info_list"`
	}
)
