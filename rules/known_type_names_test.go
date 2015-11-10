package rules_test

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
)

func TestValidate_KnownTypeNames_KnownTypeNamesAreValid(t *testing.T) {
	expectPassesRule(t, graphql.KnownTypeNamesRule, `
      query Foo($var: String, $required: [String!]!) {
        user(id: 4) {
          pets { ... on Pet { name }, ...PetFields }
        }
      }
      fragment PetFields on Pet {
        name
      }
    `)
}
func TestValidate_KnownTypeNames_UnknownTypeNamesAreInValid(t *testing.T) {
	expectFailsRule(t, graphql.KnownTypeNamesRule, `
      query Foo($var: JumbledUpLetters) {
        user(id: 4) {
          name
          pets { ... on Badger { name }, ...PetFields }
        }
      }
      fragment PetFields on Peettt {
        name
      }
    `, []gqlerrors.FormattedError{
		ruleError(`Unknown type "JumbledUpLetters".`, 2, 23),
		ruleError(`Unknown type "Badger".`, 5, 25),
		ruleError(`Unknown type "Peettt".`, 8, 29),
	})
}
