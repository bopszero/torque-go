package constants

import "gitlab.com/snap-clickstaff/torque-go/lib/meta"

const (
	RewardTypeNormal    = meta.RewardType(0)
	RewardTypeWhitelist = meta.RewardType(1)
	RewardTypeBlacklist = meta.RewardType(-1)
)
