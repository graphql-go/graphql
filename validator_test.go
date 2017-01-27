package graphql_test

import (
	"reflect"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/location"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"github.com/graphql-go/graphql/testutil"
)

func expectValid(t *testing.T, schema *graphql.Schema, queryString string) {
	source := source.NewSource(&source.Source{
		Body: []byte(queryString),
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

	expectValid(t, testutil.TestSchema, `
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

// NOTE: experimental
func TestValidator_SupportsFullValidation_ValidatesUsingACustomTypeInfo(t *testing.T) {

	// This TypeInfo will never return a valid field.
	typeInfo := graphql.NewTypeInfo(&graphql.TypeInfoConfig{
		Schema: testutil.TestSchema,
		FieldDefFn: func(schema *graphql.Schema, parentType graphql.Type, fieldAST *ast.Field) *graphql.FieldDefinition {
			return nil
		},
	})

	ast := testutil.TestParse(t, `
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

	errors := graphql.VisitUsingRules(testutil.TestSchema, typeInfo, ast, graphql.SpecifiedRules)

	expectedErrors := []gqlerrors.FormattedError{
		{
			Message: `Cannot query field "catOrDog" on type "QueryRoot". Did you mean "catOrDog"?`,
			Locations: []location.SourceLocation{
				{Line: 3, Column: 9},
			},
		},
		{
			Message: `Cannot query field "furColor" on type "Cat". Did you mean "furColor"?`,
			Locations: []location.SourceLocation{
				{Line: 5, Column: 13},
			},
		},
		{
			Message: `Cannot query field "isHousetrained" on type "Dog". Did you mean "isHousetrained"?`,
			Locations: []location.SourceLocation{
				{Line: 8, Column: 13},
			},
		},
	}
	if !reflect.DeepEqual(expectedErrors, errors) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedErrors, errors))
	}
}
