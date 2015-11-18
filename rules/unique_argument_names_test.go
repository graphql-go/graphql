package rules_test

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
)

func TestValidate_UniqueArgumentNames_NoArgumentsOnField(t *testing.T) {
	expectPassesRule(t, graphql.UniqueArgumentNamesRule, `
      {
        field
      }
    `)
}
func TestValidate_UniqueArgumentNames_NoArgumentsOnDirective(t *testing.T) {
	expectPassesRule(t, graphql.UniqueArgumentNamesRule, `
      {
        field
      }
    `)
}
func TestValidate_UniqueArgumentNames_ArgumentOnField(t *testing.T) {
	expectPassesRule(t, graphql.UniqueArgumentNamesRule, `
      {
        field(arg: "value")
      }
    `)
}
func TestValidate_UniqueArgumentNames_ArgumentOnDirective(t *testing.T) {
	expectPassesRule(t, graphql.UniqueArgumentNamesRule, `
      {
        field @directive(arg: "value")
      }
    `)
}
func TestValidate_UniqueArgumentNames_SameArgumentOnTwoFields(t *testing.T) {
	expectPassesRule(t, graphql.UniqueArgumentNamesRule, `
      {
        one: field(arg: "value")
        two: field(arg: "value")
      }
    `)
}
func TestValidate_UniqueArgumentNames_SameArgumentOnFieldAndDirective(t *testing.T) {
	expectPassesRule(t, graphql.UniqueArgumentNamesRule, `
      {
        field(arg: "value") @directive(arg: "value")
      }
    `)
}
func TestValidate_UniqueArgumentNames_SameArgumentOnTwoDirectives(t *testing.T) {
	expectPassesRule(t, graphql.UniqueArgumentNamesRule, `
      {
        field @directive1(arg: "value") @directive2(arg: "value")
      }
    `)
}
func TestValidate_UniqueArgumentNames_MultipleFieldArguments(t *testing.T) {
	expectPassesRule(t, graphql.UniqueArgumentNamesRule, `
      {
        field(arg1: "value", arg2: "value", arg3: "value")
      }
    `)
}
func TestValidate_UniqueArgumentNames_MultipleDirectiveArguments(t *testing.T) {
	expectPassesRule(t, graphql.UniqueArgumentNamesRule, `
      {
        field @directive(arg1: "value", arg2: "value", arg3: "value")
      }
    `)
}
func TestValidate_UniqueArgumentNames_DuplicateFieldArguments(t *testing.T) {
	expectFailsRule(t, graphql.UniqueArgumentNamesRule, `
      {
        field(arg1: "value", arg1: "value")
      }
    `, []gqlerrors.FormattedError{
		ruleError(`There can be only one argument named "arg1".`, 3, 15, 3, 30),
	})
}
func TestValidate_UniqueArgumentNames_ManyDuplicateFieldArguments(t *testing.T) {
	expectFailsRule(t, graphql.UniqueArgumentNamesRule, `
      {
        field(arg1: "value", arg1: "value", arg1: "value")
      }
    `, []gqlerrors.FormattedError{
		ruleError(`There can be only one argument named "arg1".`, 3, 15, 3, 30),
		ruleError(`There can be only one argument named "arg1".`, 3, 15, 3, 45),
	})
}
func TestValidate_UniqueArgumentNames_DuplicateDirectiveArguments(t *testing.T) {
	expectFailsRule(t, graphql.UniqueArgumentNamesRule, `
      {
        field @directive(arg1: "value", arg1: "value")
      }
    `, []gqlerrors.FormattedError{
		ruleError(`There can be only one argument named "arg1".`, 3, 26, 3, 41),
	})
}
func TestValidate_UniqueArgumentNames_ManyDuplicateDirectiveArguments(t *testing.T) {
	expectFailsRule(t, graphql.UniqueArgumentNamesRule, `
      {
        field @directive(arg1: "value", arg1: "value", arg1: "value")
      }
    `, []gqlerrors.FormattedError{
		ruleError(`There can be only one argument named "arg1".`, 3, 26, 3, 41),
		ruleError(`There can be only one argument named "arg1".`, 3, 26, 3, 56),
	})
}
