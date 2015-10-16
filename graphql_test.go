package gql

import (
	"reflect"
	"testing"

	"github.com/chris-ramon/graphql-go/types"
	"golang.org/x/net/context"

	"./testutil"
)

type T struct {
	Query    string
	Schema   types.GraphQLSchema
	Expected interface{}
}

var (
	Tests = []T{
		T{
			Query: `
				query HeroNameQuery {
					hero {
						name
					}
				}
			`,
			Schema: testutil.StarWarsSchema,
			Expected: &types.GraphQLResult{
				Data: map[string]interface{}{
					"hero": map[string]interface{}{
						"name": "R2-D2",
					},
				},
			},
		},
		T{
			Query: `
				query HeroNameAndFriendsQuery {
					hero {
						id
						name
						friends {
							name
						}
					}
				}
				`,
			Schema: testutil.StarWarsSchema,
			Expected: &types.GraphQLResult{
				Data: map[string]interface{}{
					"hero": map[string]interface{}{
						"id":   "2001",
						"name": "R2-D2",
						"friends": []interface{}{
							map[string]interface{}{
								"name": "Luke Skywalker",
							},
							map[string]interface{}{
								"name": "Han Solo",
							},
							map[string]interface{}{
								"name": "Leia Organa",
							},
						},
					},
				},
			},
		},
	}
)

func TestQuery(t *testing.T) {
	for _, test := range Tests {
		graphqlParams := GraphqlParams{
			Schema:        test.Schema,
			RequestString: test.Query,
		}
		testGraphql(test, graphqlParams, t)
	}
}

func testGraphql(test T, p GraphqlParams, t *testing.T) {
	resultChannel := make(chan *types.GraphQLResult)
	go Graphql(p, resultChannel)
	result := <-resultChannel
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(result, test.Expected) {
		t.Fatalf("wrong result, query: %v, graphql result diff: %v", test.Query, testutil.Diff(test.Expected, result))
	}
}

func TestBasicGraphQLExample(t *testing.T) {
	// taken from `graphql-js` README

	helloFieldResolved := func(ctx context.Context, p types.GQLFRParams) interface{} {
		return "world"
	}

	schema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "RootQueryType",
			Fields: types.GraphQLFieldConfigMap{
				"hello": &types.GraphQLFieldConfig{
					Description: "Returns `world`",
					Type:        types.GraphQLString,
					Resolve:     helloFieldResolved,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("wrong result, unexpected errors: %v", err.Error())
	}
	query := "{ hello }"
	var expected interface{}
	expected = map[string]interface{}{
		"hello": "world",
	}

	resultChannel := make(chan *types.GraphQLResult)
	go Graphql(GraphqlParams{
		Schema:        schema,
		RequestString: query,
	}, resultChannel)
	result := <-resultChannel
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(result.Data, expected) {
		t.Fatalf("wrong result, query: %v, graphql result diff: %v", query, testutil.Diff(expected, result))
	}

}
