package rules_test

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
)

func TestValidate_KnownFragmentNames_KnownFragmentNamesAreValid(t *testing.T) {
	expectPassesRule(t, graphql.KnownFragmentNamesRule, `
      {
        human(id: 4) {
          ...HumanFields1
          ... on Human {
            ...HumanFields2
          }
        }
      }
      fragment HumanFields1 on Human {
        name
        ...HumanFields3
      }
      fragment HumanFields2 on Human {
        name
      }
      fragment HumanFields3 on Human {
        name
      }
    `)
}
func TestValidate_KnownFragmentNames_UnknownFragmentNamesAreInvalid(t *testing.T) {
	expectFailsRule(t, graphql.KnownFragmentNamesRule, `
      {
        human(id: 4) {
          ...UnknownFragment1
          ... on Human {
            ...UnknownFragment2
          }
        }
      }
      fragment HumanFields on Human {
        name
        ...UnknownFragment3
      }
    `, []gqlerrors.FormattedError{
		ruleError(`Unknown fragment "UnknownFragment1".`, 4, 14),
		ruleError(`Unknown fragment "UnknownFragment2".`, 6, 16),
		ruleError(`Unknown fragment "UnknownFragment3".`, 12, 12),
	})
}
