package constants

import (
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	OrderDirectionTopup   = meta.Direction(1)
	OrderDirectionPayment = meta.Direction(-1)
)

const (
	OrderStepDirectionForward  = meta.Direction(1)
	OrderStepDirectionBackward = meta.Direction(-1)
)

const (
	OrderStatusExpired  = meta.OrderStatus(-3)
	OrderStatusCanceled = meta.OrderStatus(-2)
	OrderStatusFailed   = meta.OrderStatus(-1)

	OrderStatusUnknown = meta.OrderStatus(0)

	OrderStatusNew       = meta.OrderStatus(1)
	OrderStatusInit      = meta.OrderStatus(2)
	OrderStatusHandleSrc = meta.OrderStatus(3)
	OrderStatusHandleDst = meta.OrderStatus(4)

	OrderStatusNeedStaff = meta.OrderStatus(50)

	OrderStatusFailing    = meta.OrderStatus(97)
	OrderStatusRefunding  = meta.OrderStatus(98)
	OrderStatusCompleting = meta.OrderStatus(99)
	OrderStatusCompleted  = meta.OrderStatus(100)
	OrderStatusRefunded   = meta.OrderStatus(101)
)

var (
	OrderPendingStatuses = []meta.OrderStatus{
		OrderStatusHandleSrc,
		OrderStatusHandleDst,
		OrderStatusFailing,
		OrderStatusRefunding,
		OrderStatusCompleting,
	}
	OrderVisibleStatuses = []meta.OrderStatus{
		OrderStatusExpired,
		OrderStatusFailed,
		OrderStatusHandleSrc,
		OrderStatusHandleDst,
		OrderStatusNeedStaff,
		OrderStatusFailing,
		OrderStatusRefunding,
		OrderStatusCompleting,
		OrderStatusCompleted,
		OrderStatusRefunded,
	}
	OrderFailingStatusSet = comtypes.NewHashSetFromListF([]meta.OrderStatus{
		OrderStatusFailing,
		OrderStatusRefunding,
	})
)

const (
	OrderStepCodeEmpty = ""

	OrderStepCodeOrderInit      = "oi"
	OrderStepCodeOrderStartSrc  = "oss"
	OrderStepCodeOrderFinishSrc = "ofs"
	OrderStepCodeOrderStartDst  = "osd"
	OrderStepCodeOrderFinishDst = "ofd"
	OrderStepCodeOrderDone      = "od"
)

const (
	OrderExtraDataBlockchainInfo = "blockchain_info"
)
