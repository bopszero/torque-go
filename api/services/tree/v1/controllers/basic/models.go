package basic

import (
	"gitlab.com/snap-clickstaff/torque-go/lib/affiliate"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type GetNodeDownRequest struct {
	UID     meta.UID              `json:"user_id" validate:"required"`
	Options affiliate.ScanOptions `json:"options"`
}

type GetNodeDownResponse struct {
	Node affiliate.ScanNodeInfo `json:"node"`
}
