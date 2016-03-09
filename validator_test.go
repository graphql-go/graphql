package graphql_test

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"github.com/graphql-go/graphql/testutil"
)

func expectValid(t *testing.T, schema *graphql.Schema, queryString string) {
	source := source.NewSource(&source.Source{
		Body: queryString,
		Name: "GraphQL request",
	})
	AST, err := parser.Parse(parser.ParseParams{Source: source})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	validationResult := graphql.ValidateDocument(schema, AST, nil)

	if !validationResult.IsValid || len(validationResult.Errors) > 0 {
		t.Fatalf("Unexpected error: %v", validationResult.Errors)
	}

}

func TestValidator_SupportsFullValidation_ValidatesQueries(t *testing.T) {

	expectValid(t, testutil.DefaultRulesTestSchema, `
      query {
        catOrDog {
          ... on Cat {
            furColor
          }
          ... on Dog {
            isHousetrained
          }
        }
      }
    `)
}
