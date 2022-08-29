package constants

import "gitlab.com/snap-clickstaff/torque-go/lib/meta"

const (
	SystemWithdrawalRequestStatusInit         = meta.SystemWithdrawalRequestStatus(1)
	SystemWithdrawalRequestStatusConfirmed    = meta.SystemWithdrawalRequestStatus(3)
	SystemWithdrawalRequestStatusTransferring = meta.SystemWithdrawalRequestStatus(4)
	SystemWithdrawalRequestStatusCancelled    = meta.SystemWithdrawalRequestStatus(9)
	SystemWithdrawalRequestStatusCompleted    = meta.SystemWithdrawalRequestStatus(10)
)

const (
	SystemWithdrawalTxnStatusReplaced  = meta.SystemWithdrawalTxnStatus(-2)
	SystemWithdrawalTxnStatusFailed    = meta.SystemWithdrawalTxnStatus(-1)
	SystemWithdrawalTxnStatusNew       = meta.SystemWithdrawalTxnStatus(0)
	SystemWithdrawalTxnStatusInit      = meta.SystemWithdrawalTxnStatus(1)
	SystemWithdrawalTxnStatusPushed    = meta.SystemWithdrawalTxnStatus(2)
	SystemWithdrawalTxnStatusCancelled = meta.SystemWithdrawalTxnStatus(3)
	SystemWithdrawalTxnStatusSuccess   = meta.SystemWithdrawalTxnStatus(10)
)

const (
	SystemForwardingOrderStatusFailed     = meta.SystemForwardingOrderStatus(-1)
	SystemForwardingOrderStatusInit       = meta.SystemForwardingOrderStatus(1)
	SystemForwardingOrderStatusSigned     = meta.SystemForwardingOrderStatus(2)
	SystemForwardingOrderStatusForwarding = meta.SystemForwardingOrderStatus(3)
	SystemForwardingOrderStatusCompleted  = meta.SystemForwardingOrderStatus(10)
)

const (
	SystemForwardingOrderTxnStatusIgnored   = meta.SystemForwardingOrderTxnStatus(-2)
	SystemForwardingOrderTxnStatusFailed    = meta.SystemForwardingOrderTxnStatus(-1)
	SystemForwardingOrderTxnStatusUnknown   = meta.SystemForwardingOrderTxnStatus(0)
	SystemForwardingOrderTxnStatusInit      = meta.SystemForwardingOrderTxnStatus(1)
	SystemForwardingOrderTxnStatusFeeComing = meta.SystemForwardingOrderTxnStatus(2)
	SystemForwardingOrderTxnStatusPushed    = meta.SystemForwardingOrderTxnStatus(3)
	SystemForwardingOrderTxnStatusSuccess   = meta.SystemForwardingOrderTxnStatus(10)
)

var (
	SystemForwardingOrderTxnStatusNameMap = map[meta.SystemForwardingOrderTxnStatus]string{
		SystemForwardingOrderTxnStatusIgnored:   "Ignored",
		SystemForwardingOrderTxnStatusFailed:    "Failed",
		SystemForwardingOrderTxnStatusUnknown:   "Unknown",
		SystemForwardingOrderTxnStatusInit:      "Init",
		SystemForwardingOrderTxnStatusFeeComing: "Fee Coming",
		SystemForwardingOrderTxnStatusPushed:    "Pushed",
		SystemForwardingOrderTxnStatusSuccess:   "Success",
	}
)
