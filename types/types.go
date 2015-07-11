package types

import "github.com/chris-ramon/graphql-go/errors"

type Schema interface{}

type GraphQLResult struct {
	Data   interface{}
	Errors []errors.GraphQLFormattedError
}

type GraphQLEnumType struct {
}

type GraphQLInterfaceType struct {
}

type GraphQLObjectTypeFields func()

type GraphQLObjectType struct {
	Name   string
	Fields GraphQLObjectTypeFields
}

type GraphQLList struct {
}

type GraphQLNonNull struct {
}

type GraphQLSchemaConfig struct {
	query    GraphQLObjectType
	mutation GraphQLObjectType
}

type GraphQLSchema struct {
	Query        GraphQLObjectType
	schemaConfig GraphQLSchemaConfig
}

func (gq *GraphQLSchema) Constructor(config GraphQLSchemaConfig) {
	gq.schemaConfig = config
}

type GraphQLString struct {
}
