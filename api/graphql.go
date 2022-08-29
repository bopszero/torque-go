package api

import "github.com/graphql-go/graphql"

func GraphQLUpdateFields(mainFields graphql.Fields, offerFields graphql.Fields) graphql.Fields {
	for key, value := range offerFields {
		mainFields[key] = value
	}

	return mainFields
}
