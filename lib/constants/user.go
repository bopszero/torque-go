package constants

import "gitlab.com/snap-clickstaff/torque-go/lib/meta"

const (
	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
	UserStatusPending  = "pending"
	UserStatusRejected = "rejected"
)

const (
	TierTypeStatusDynamic = 1
	TierTypeStatusFixed   = 2

	TierChangeStatusReady   = 1
	TierChangeStatusApplied = 2
)

var (
	TierMetaZero = meta.TierMeta{ID: 0, Code: "unknown"}

	TierMetaAgent           = meta.TierMeta{ID: 1, Code: "agent"}
	TierMetaSeniorPartner   = meta.TierMeta{ID: 2, Code: "market_leader"}
	TierMetaRegionalPartner = meta.TierMeta{ID: 3, Code: "region_leader"}
	TierMetaGlobalPartner   = meta.TierMeta{ID: 6, Code: "global_leader"}
	TierMetaMentor          = meta.TierMeta{ID: 7, Code: "mentor"}

	TierMetaGovernor = meta.TierMeta{ID: 4, Code: "governor"}

	TierMetaOrderedList = []meta.TierMeta{
		TierMetaAgent,
		TierMetaSeniorPartner,
		TierMetaRegionalPartner,
		TierMetaGlobalPartner,
		TierMetaMentor,

		TierMetaGovernor,
	}
	TierMetaLeaderOrderedList []meta.TierMeta
	TierMetaMap               map[int]meta.TierMeta
	TierMetaCodeMap           map[string]meta.TierMeta
	TierOrderCodeMap          map[string]int
)

func init() {
	TierMetaLeaderOrderedList = make([]meta.TierMeta, 0)
	TierMetaMap = make(map[int]meta.TierMeta)
	TierMetaCodeMap = make(map[string]meta.TierMeta)
	TierOrderCodeMap = make(map[string]int)
	for idx, tierMeta := range TierMetaOrderedList {
		TierMetaMap[tierMeta.ID] = tierMeta
		TierMetaCodeMap[tierMeta.Code] = tierMeta
		TierOrderCodeMap[tierMeta.Code] = idx

		if tierMeta.Code != TierMetaAgent.Code && tierMeta.Code != TierMetaGovernor.Code {
			TierMetaLeaderOrderedList = append(TierMetaLeaderOrderedList, tierMeta)
		}
	}
}
