package rules_test

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
)

func TestValidate_UniqueOperationNames_NoOperations(t *testing.T) {
	expectPassesRule(t, graphql.UniqueOperationNamesRule, `
      fragment fragA on Type {
        field
      }
    `)
}
func TestValidate_UniqueOperationNames_OneAnonOperation(t *testing.T) {
	expectPassesRule(t, graphql.UniqueOperationNamesRule, `
      {
        field
      }
    `)
}
func TestValidate_UniqueOperationNames_OneNamedOperation(t *testing.T) {
	expectPassesRule(t, graphql.UniqueOperationNamesRule, `
      query Foo {
        field
      }
    `)
}
func TestValidate_UniqueOperationNames_MultipleOperations(t *testing.T) {
	expectPassesRule(t, graphql.UniqueOperationNamesRule, `
      query Foo {
        field
      }

      query Bar {
        field
      }
    `)
}
func TestValidate_UniqueOperationNames_MultipleOperationsOfDifferentTypes(t *testing.T) {
	expectPassesRule(t, graphql.UniqueOperationNamesRule, `
      query Foo {
        field
      }

      mutation Bar {
        field
      }
    `)
}
func TestValidate_UniqueOperationNames_FragmentAndOperationNamedTheSame(t *testing.T) {
	expectPassesRule(t, graphql.UniqueOperationNamesRule, `
      query Foo {
        field
      }

      mutation Bar {
        field
      }
    `)
}
func TestValidate_UniqueOperationNames_MultipleOperationsOfSameName(t *testing.T) {
	expectFailsRule(t, graphql.UniqueOperationNamesRule, `
      query Foo {
        fieldA
      }
      query Foo {
        fieldB
      }
    `, []gqlerrors.FormattedError{
		ruleError(`There can only be one operation named "Foo".`, 2, 13, 5, 13),
	})
}
func TestValidate_UniqueOperationNames_MultipleOperationsOfSameNameOfDifferentTypes(t *testing.T) {
	expectFailsRule(t, graphql.UniqueOperationNamesRule, `
      query Foo {
        fieldA
      }
      mutation Foo {
        fieldB
      }
    `, []gqlerrors.FormattedError{
		ruleError(`There can only be one operation named "Foo".`, 2, 13, 5, 16),
	})
}
