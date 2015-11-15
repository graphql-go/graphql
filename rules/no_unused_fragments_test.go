package rules_test

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
)

func TestValidate_NoUnusedFragments_AllFragmentNamesAreUsed(t *testing.T) {
	expectPassesRule(t, graphql.NoUnusedFragmentsRule, `
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
func TestValidate_NoUnusedFragments_AllFragmentNamesAreUsedByMultipleOperations(t *testing.T) {
	expectPassesRule(t, graphql.NoUnusedFragmentsRule, `
      query Foo {
        human(id: 4) {
          ...HumanFields1
        }
      }
      query Bar {
        human(id: 4) {
          ...HumanFields2
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
func TestValidate_NoUnusedFragments_ContainsUnknownFragments(t *testing.T) {
	expectFailsRule(t, graphql.NoUnusedFragmentsRule, `
      query Foo {
        human(id: 4) {
          ...HumanFields1
        }
      }
      query Bar {
        human(id: 4) {
          ...HumanFields2
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
      fragment Unused1 on Human {
        name
      }
      fragment Unused2 on Human {
        name
      }
    `, []gqlerrors.FormattedError{
		ruleError(`Fragment "Unused1" is never used.`, 22, 7),
		ruleError(`Fragment "Unused2" is never used.`, 25, 7),
	})
}

func TestValidate_NoUnusedFragments_ContainsUnknownFragmentsWithRefCycle(t *testing.T) {
	expectFailsRule(t, graphql.NoUnusedFragmentsRule, `
      query Foo {
        human(id: 4) {
          ...HumanFields1
        }
      }
      query Bar {
        human(id: 4) {
          ...HumanFields2
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
      fragment Unused1 on Human {
        name
        ...Unused2
      }
      fragment Unused2 on Human {
        name
        ...Unused1
      }
    `, []gqlerrors.FormattedError{
		ruleError(`Fragment "Unused1" is never used.`, 22, 7),
		ruleError(`Fragment "Unused2" is never used.`, 26, 7),
	})
}

func TestValidate_NoUnusedFragments_ContainsUnknownAndUndefFragments(t *testing.T) {
	expectFailsRule(t, graphql.NoUnusedFragmentsRule, `
      query Foo {
        human(id: 4) {
          ...bar
        }
      }
      fragment foo on Human {
        name
      }
    `, []gqlerrors.FormattedError{
		ruleError(`Fragment "foo" is never used.`, 7, 7),
	})
}
