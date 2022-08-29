package v1

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/torque-go/api/middleware"
	"gitlab.com/snap-clickstaff/torque-go/api/services/tree/v1/controllers/basic"
	"gitlab.com/snap-clickstaff/torque-go/api/services/tree/v1/controllers/promotion"
)

func InitGroup(group *echo.Group) {
	logRequestMiddleware := middleware.NewRequestLogger(middleware.LogRequestOptions{
		RequestLimit: 500,
		LogResponse:  false,
	})

	group.POST("/basic/get_node_down/", basic.GetNodeDown, logRequestMiddleware)
	group.POST("/promotion/stats/get/", promotion.GetStats, middleware.LogRequestDefaultMiddleware)
}
