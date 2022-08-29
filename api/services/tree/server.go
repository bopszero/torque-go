package tree

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/torque-go/api"
	"gitlab.com/snap-clickstaff/torque-go/api/services/tree/globals"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func genEchoObject() *echo.Echo {
	e := api.CreateEchoObject()

	InitRouter(e)

	return e
}

func StartServer(host string, port int) {
	e := genEchoObject()

	go func() {
		defer utils.ErrorCatchWithLog(nil, "tree service pre-init tree", nil)
		globals.GetTree()
	}()

	api.StartServer(host, port, e)
}
