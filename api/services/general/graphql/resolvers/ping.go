package resolvers

import "github.com/graphql-go/graphql"

func PingResolver(p graphql.ResolveParams) (interface{}, error) {
	return "pong", nil
}
