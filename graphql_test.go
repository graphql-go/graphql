package gql

import (
	"reflect"
	"testing"

	"github.com/chris-ramon/graphql-go/types"

	"./testutil"
)

type T struct {
	Query    string
	Schema   types.GraphQLSchema
	Expected interface{}
	Result   interface{}
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
			Expected: &testutil.StarWarsChar{
				Name: "R2-D2",
			},
			Result: &testutil.StarWarsChar{},
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
			Expected: &testutil.StarWarsChar{
				Id:   "2001",
				Name: "R2-D2",
				Friends: []testutil.StarWarsChar{
					testutil.StarWarsChar{
						Name: "Luke Skywalker",
					},
					testutil.StarWarsChar{
						Name: "Han Solo",
					},
					testutil.StarWarsChar{
						Name: "Leia Organa",
					},
				},
			},
			Result: &testutil.StarWarsChar{},
		},
	}
)

func TestQuery(t *testing.T) {
	for _, test := range Tests {
		graphqlParams := GraphqlParams{
			Schema:        test.Schema,
			RequestString: test.Query,
			Result:        test.Result,
		}
		testGraphql(test, graphqlParams, t)
	}
}

func testGraphql(test T, p GraphqlParams, t *testing.T) {
	resultChannel := make(chan types.GraphQLResult)
	go Graphql(p, resultChannel)
	graphqlResult := <-resultChannel
	if len(graphqlResult.Errors) > 0 {
		t.Errorf("wrong result, unexpected errors: %v", graphqlResult.Errors)
	}
	if !reflect.DeepEqual(test.Result, test.Expected) {
		t.Errorf("wrong result, query: %v, graphql result: %v, expected: %v", test.Query, test.Result, test.Expected)
	}
}
