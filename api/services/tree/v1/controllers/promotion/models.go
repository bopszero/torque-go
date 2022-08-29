package promotion

import "gitlab.com/snap-clickstaff/torque-go/lib/meta"

type PromotionStatsGetRequest struct {
	UID meta.UID `json:"uid" validate:"required"`
}

type PromotionStats struct {
	FromTier string     `json:"from_tier"`
	ToTier   string     `json:"to_tier"`
	UIDs     []meta.UID `json:"uids"`
}

func NewPromotionStats(fromTier string, toTier string) *PromotionStats {
	return &PromotionStats{
		FromTier: fromTier,
		ToTier:   toTier,
		UIDs:     make([]meta.UID, 0),
	}
}

type PromotionStatsGetResponse struct {
	UID               meta.UID          `json:"uid"`
	Tier              string            `json:"tier"`
	DownLineStatsList []*PromotionStats `json:"down_line_stats_list"`
}
