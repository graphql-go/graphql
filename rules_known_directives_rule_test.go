package graphql_test

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/testutil"
)

func TestValidate_KnownDirectives_WithNoDirectives(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.KnownDirectivesRule, `
      query Foo {
        name
        ...Frag
      }

      fragment Frag on Dog {
        name
      }
    `)
}
func TestValidate_KnownDirectives_WithKnownDirective(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.KnownDirectivesRule, `
      {
        dog @include(if: true) {
          name
        }
        human @skip(if: false) {
          name
        }
      }
    `)
}
func TestValidate_KnownDirectives_WithUnknownDirective(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.KnownDirectivesRule, `
      {
        dog @unknown(directive: "value") {
          name
        }
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Unknown directive "unknown".`, 3, 13),
	})
}
func TestValidate_KnownDirectives_WithManyUnknownDirectives(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.KnownDirectivesRule, `
      {
        dog @unknown(directive: "value") {
          name
        }
        human @unknown(directive: "value") {
          name
          pets @unknown(directive: "value") {
            name
          }
        }
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Unknown directive "unknown".`, 3, 13),
		testutil.RuleError(`Unknown directive "unknown".`, 6, 15),
		testutil.RuleError(`Unknown directive "unknown".`, 8, 16),
	})
}
func TestValidate_KnownDirectives_WithWellPlacedDirectives(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.KnownDirectivesRule, `
      query Foo {
        name @include(if: true)
        ...Frag @include(if: true)
        skippedField @skip(if: true)
        ...SkippedFrag @skip(if: true)
      }
    `)
}
func TestValidate_KnownDirectives_WithMisplacedDirectives(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.KnownDirectivesRule, `
      query Foo @include(if: true) {
        name
        ...Frag
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Directive "include" may not be used on "operation".`, 2, 17),
	})
}
