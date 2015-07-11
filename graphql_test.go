package gql

import (
	"testing"

	"github.com/chris-ramon/graphql-go/types"

	"./testutil"
)

func TestQuery(t *testing.T) {
	var query string
	query = `
		query HeroNameQuery {
			hero {
				name
			}
	}
	`
	expected := testutil.StarWarsChar{
		Name: "R2-D2",
	}
	graphqlParams := GraphqlParams{
		Schema:        testutil.StarWarsSchema,
		RequestString: query,
	}
	resultChannel := make(chan types.GraphQLResult)
	Graphql(graphqlParams, resultChannel)
	graphqlResult := <-resultChannel
	hero := graphqlResult.Data["hero"].(map[string]interface{})
	close(resultChannel)
	if expected.Name != hero["name"] {
		t.Errorf("wrong result, query: %v, graphql result: %v", query, graphqlResult)
	}
}
