package rules_test

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
)

func TestValidate_KnownDirectives_WithNoDirectives(t *testing.T) {
	expectPassesRule(t, graphql.KnownDirectivesRule, `
      query Foo {
        name
        ...Frag
      }

      fragment Frag on Dog {
        name
      }
    `)
}
func TestValidate_KnownDirectives_WithUnknownDirective(t *testing.T) {
	expectFailsRule(t, graphql.KnownDirectivesRule, `
      {
        dog @unknown(directive: "value") {
          name
        }
      }
    `, []gqlerrors.FormattedError{
		ruleError(`Unknown directive "unknown".`, 3, 13),
	})
}
func TestValidate_KnownDirectives_WithManyUnknownDirectives(t *testing.T) {
	expectFailsRule(t, graphql.KnownDirectivesRule, `
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
		ruleError(`Unknown directive "unknown".`, 3, 13),
		ruleError(`Unknown directive "unknown".`, 6, 15),
		ruleError(`Unknown directive "unknown".`, 8, 16),
	})
}
func TestValidate_KnownDirectives_WithWellPlacedDirectives(t *testing.T) {
	expectPassesRule(t, graphql.KnownDirectivesRule, `
      query Foo {
        name @include(if: true)
        ...Frag @include(if: true)
        skippedField @skip(if: true)
        ...SkippedFrag @skip(if: true)
      }
    `)
}
func TestValidate_KnownDirectives_WithMisplacedDirectives(t *testing.T) {
	expectFailsRule(t, graphql.KnownDirectivesRule, `
      query Foo @include(if: true) {
        name
        ...Frag
      }
    `, []gqlerrors.FormattedError{
		ruleError(`Directive "include" may not be used on "operation".`, 2, 17),
	})
}
