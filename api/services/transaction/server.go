package transaction

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/torque-go/api"
)

func genEchoObject() *echo.Echo {
	e := api.CreateEchoObject()

	InitRouter(e)

	return e
}

func StartServer(host string, port int) {
	e := genEchoObject()
	api.StartServer(host, port, e)
}
