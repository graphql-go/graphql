package rules_test

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
)

func TestValidate_AnonymousOperationMustBeAlone_NoOperations(t *testing.T) {
	expectPassesRule(t, graphql.LoneAnonymousOperationRule, `
      fragment fragA on Type {
        field
      }
    `)
}
func TestValidate_AnonymousOperationMustBeAlone_OneAnonOperation(t *testing.T) {
	expectPassesRule(t, graphql.LoneAnonymousOperationRule, `
      {
        field
      }
    `)
}
func TestValidate_AnonymousOperationMustBeAlone_MultipleNamedOperations(t *testing.T) {
	expectPassesRule(t, graphql.LoneAnonymousOperationRule, `
      query Foo {
        field
      }

      query Bar {
        field
      }
    `)
}
func TestValidate_AnonymousOperationMustBeAlone_AnonOperationWithFragment(t *testing.T) {
	expectPassesRule(t, graphql.LoneAnonymousOperationRule, `
      {
        ...Foo
      }
      fragment Foo on Type {
        field
      }
    `)
}

func TestValidate_AnonymousOperationMustBeAlone_MultipleAnonOperations(t *testing.T) {
	expectFailsRule(t, graphql.LoneAnonymousOperationRule, `
      {
        fieldA
      }
      {
        fieldB
      }
    `, []gqlerrors.FormattedError{
		ruleError(`This anonymous operation must be the only defined operation.`, 2, 7),
		ruleError(`This anonymous operation must be the only defined operation.`, 5, 7),
	})
}
func TestValidate_AnonymousOperationMustBeAlone_AnonOperationWithAnotherOperation(t *testing.T) {
	expectFailsRule(t, graphql.LoneAnonymousOperationRule, `
      {
        fieldA
      }
      mutation Foo {
        fieldB
      }
    `, []gqlerrors.FormattedError{
		ruleError(`This anonymous operation must be the only defined operation.`, 2, 7),
	})
}
