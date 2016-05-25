package graphql_test

import (
	"testing"

	"github.com/sprucehealth/graphql"
	"github.com/sprucehealth/graphql/gqlerrors"
	"github.com/sprucehealth/graphql/testutil"
)

func TestValidate_AnonymousOperationMustBeAlone_NoOperations(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.LoneAnonymousOperationRule, `
      fragment fragA on Type {
        field
      }
    `)
}
func TestValidate_AnonymousOperationMustBeAlone_OneAnonOperation(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.LoneAnonymousOperationRule, `
      {
        field
      }
    `)
}
func TestValidate_AnonymousOperationMustBeAlone_MultipleNamedOperations(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.LoneAnonymousOperationRule, `
      query Foo {
        field
      }

      query Bar {
        field
      }
    `)
}
func TestValidate_AnonymousOperationMustBeAlone_AnonOperationWithFragment(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.LoneAnonymousOperationRule, `
      {
        ...Foo
      }
      fragment Foo on Type {
        field
      }
    `)
}
func TestValidate_AnonymousOperationMustBeAlone_MultipleAnonOperations(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.LoneAnonymousOperationRule, `
      {
        fieldA
      }
      {
        fieldB
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`This anonymous operation must be the only defined operation.`, 2, 7),
		testutil.RuleError(`This anonymous operation must be the only defined operation.`, 5, 7),
	})
}
func TestValidate_AnonymousOperationMustBeAlone_AnonOperationWithAnotherOperation(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.LoneAnonymousOperationRule, `
      {
        fieldA
      }
      mutation Foo {
        fieldB
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`This anonymous operation must be the only defined operation.`, 2, 7),
	})
}
