package ordermod

import (
	"reflect"

	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type Channel interface {
	SetInfoModel(infoModel models.ChannelInfo)
	IsAvailable() bool

	GetType() meta.ChannelType
	GetInfoModel() models.ChannelInfo
	GetMetaType() reflect.Type
	GetCheckoutInfo(comcontext.Context, *models.Order) (interface{}, error)
	GetOrderDetails(comcontext.Context, *models.Order) (interface{}, error)

	Init(comcontext.Context, *models.Order) error
	PreValidate(comcontext.Context, *models.Order) error

	Prepare(comcontext.Context, *models.Order) (meta.OrderStepResultCode, error)
	Execute(comcontext.Context, *models.Order) (meta.OrderStepResultCode, error)
	Commit(comcontext.Context, *models.Order) (meta.OrderStepResultCode, error)
	PrepareReverse(comcontext.Context, *models.Order) (meta.OrderStepResultCode, error)
	ExecuteReverse(comcontext.Context, *models.Order) (meta.OrderStepResultCode, error)
	CommitReverse(comcontext.Context, *models.Order) (meta.OrderStepResultCode, error)

	GetNotificationCompleted(ctx comcontext.Context, order models.Order) (*Notification, error)
	GetNotificationFailed(ctx comcontext.Context, order models.Order) (*Notification, error)
}
