package api

import (
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/torque-go/config"
)

func InitGraphQLRoutes(e *echo.Echo, schema *graphql.Schema) {
	config := handler.Config{
		Schema:     schema,
		Pretty:     true,
		GraphiQL:   true,
		Playground: config.Debug,
	}

	graphqlHandler := handler.New(&config)

	getHandler := func(c echo.Context) error {
		request := c.Request()

		if config.Playground && request.Method == "GET" {
			acceptHeader := request.Header.Get("Accept")
			acceptJson := strings.Contains(acceptHeader, echo.MIMEApplicationJSON)
			acceptHtml := strings.Contains(acceptHeader, echo.MIMETextHTML)
			_, enableRaw := request.URL.Query()["raw"]

			if !enableRaw && acceptHtml && !acceptJson {
				return c.File("templates/graphql_playground.html")
			}
		}

		graphqlHandler.ServeHTTP(c.Response(), request)
		return nil
	}

	e.GET("/graphql/", getHandler)
	e.POST("/graphql/", echo.WrapHandler(graphqlHandler))
}
