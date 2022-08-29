package bonuspoolmod

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type LeaderTierInfo struct {
	TierType int             `json:"tier_type" validate:"required" msgpack:"tier_type"`
	UIDs     []meta.UID      `json:"uids" validate:"required,min=1" msgpack:"uids"`
	Rate     decimal.Decimal `json:"rate" validate:"required" msgpack:"-"`

	AdditionalUserInfoList []LeaderTierAdditionalUserInfo `json:"additional_user_info_list,omitempty" msgpack:"-"`
}

type LeaderTierAdditionalUserInfo struct {
	UID  meta.UID `json:"uid" validate:"required"`
	Note string   `json:"note" validate:"required"`
}
