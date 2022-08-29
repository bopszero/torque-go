package graphql

import (
	"github.com/graphql-go/graphql"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/api"
	"gitlab.com/snap-clickstaff/torque-go/api/services/general/graphql/resolvers"
)

func GenSchema() *graphql.Schema {
	rootQuery := genRootQuery()
	rootMutation := genRootMutation()

	schemaConfig := graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	}

	schema, err := graphql.NewSchema(schemaConfig)
	comutils.PanicOnError(err)

	return &schema
}

func genRootQuery() *graphql.Object {
	rootFields := graphql.Fields{}
	rootFields = api.GraphQLUpdateFields(rootFields, PingQueryFields)

	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "RootQuery",
			Fields: rootFields,
		},
	)
}

func genRootMutation() *graphql.Object {
	rootFields := graphql.Fields{}
	rootFields = api.GraphQLUpdateFields(rootFields, PingMutationFields)

	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "RootMutation",
			Fields: rootFields,
		},
	)
}

var PingQueryFields = graphql.Fields{
	"ping": &graphql.Field{
		Type:        graphql.String,
		Resolve:     resolvers.PingResolver,
		Description: "Simply for ping/pong.",
	},
}

var PingMutationFields = graphql.Fields{
	"ping": &graphql.Field{
		Type:        graphql.String,
		Resolve:     resolvers.PingResolver,
		Description: "Simply for ping/pong.",
	},
}
