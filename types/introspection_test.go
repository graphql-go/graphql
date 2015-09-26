package types

import (
	"testing"
)

func TestExecutesAnIntrospectionQuery(t *testing.T) {
	_, err := NewGraphQLSchema(GraphQLSchemaConfig{
		Query: NewGraphQLObjectType(GraphQLObjectTypeConfig{
			Name: "QueryRoot",
			Fields: GraphQLFieldConfigMap{
				"onlyField": &GraphQLFieldConfig{
					Type: GraphQLString,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error creating GraphQLSchema: %v", err.Error())
	}
}
