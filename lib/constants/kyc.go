package constants

import (
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	KycRequestStatusInit            = meta.KycRequestStatus(1)
	KycRequestStatusPendingAnalysis = meta.KycRequestStatus(2)
	KycRequestStatusPendingApproval = meta.KycRequestStatus(5)
	KycRequestStatusApproved        = meta.KycRequestStatus(3)
	KycRequestStatusRejected        = meta.KycRequestStatus(6)
	KycRequestStatusFailed          = meta.KycRequestStatus(7)
)
