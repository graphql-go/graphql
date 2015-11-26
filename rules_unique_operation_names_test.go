package graphql_test

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/testutil"
)

func TestValidate_UniqueOperationNames_NoOperations(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.UniqueOperationNamesRule, `
      fragment fragA on Type {
        field
      }
    `)
}
func TestValidate_UniqueOperationNames_OneAnonOperation(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.UniqueOperationNamesRule, `
      {
        field
      }
    `)
}
func TestValidate_UniqueOperationNames_OneNamedOperation(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.UniqueOperationNamesRule, `
      query Foo {
        field
      }
    `)
}
func TestValidate_UniqueOperationNames_MultipleOperations(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.UniqueOperationNamesRule, `
      query Foo {
        field
      }

      query Bar {
        field
      }
    `)
}
func TestValidate_UniqueOperationNames_MultipleOperationsOfDifferentTypes(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.UniqueOperationNamesRule, `
      query Foo {
        field
      }

      mutation Bar {
        field
      }
    `)
}
func TestValidate_UniqueOperationNames_FragmentAndOperationNamedTheSame(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.UniqueOperationNamesRule, `
      query Foo {
        ...Foo
      }
      fragment Foo on Type {
        field
      }
    `)
}
func TestValidate_UniqueOperationNames_MultipleOperationsOfSameName(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.UniqueOperationNamesRule, `
      query Foo {
        fieldA
      }
      query Foo {
        fieldB
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`There can only be one operation named "Foo".`, 2, 13, 5, 13),
	})
}
func TestValidate_UniqueOperationNames_MultipleOperationsOfSameNameOfDifferentTypes(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.UniqueOperationNamesRule, `
      query Foo {
        fieldA
      }
      mutation Foo {
        fieldB
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`There can only be one operation named "Foo".`, 2, 13, 5, 16),
	})
}
