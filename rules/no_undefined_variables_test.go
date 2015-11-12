package rules_test

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
)

func TestValidate_NoUndefinedVariables_AllVariablesDefined(t *testing.T) {
	expectPassesRule(t, graphql.NoUndefinedVariablesRule, `
      query Foo($a: String, $b: String, $c: String) {
        field(a: $a, b: $b, c: $c)
      }
    `)
}
func TestValidate_NoUndefinedVariables_AllVariablesDeeplyDefined(t *testing.T) {
	expectPassesRule(t, graphql.NoUndefinedVariablesRule, `
      query Foo($a: String, $b: String, $c: String) {
        field(a: $a) {
          field(b: $b) {
            field(c: $c)
          }
        }
      }
    `)
}
func TestValidate_NoUndefinedVariables_AllVariablesDeeplyDefinedInInlineFragmentsDefined(t *testing.T) {
	expectPassesRule(t, graphql.NoUndefinedVariablesRule, `
      query Foo($a: String, $b: String, $c: String) {
        ... on Type {
          field(a: $a) {
            field(b: $b) {
              ... on Type {
                field(c: $c)
              }
            }
          }
        }
      }
    `)
}
func TestValidate_NoUndefinedVariables_AllVariablesInFragmentsDeeplyDefined(t *testing.T) {
	expectPassesRule(t, graphql.NoUndefinedVariablesRule, `
      query Foo($a: String, $b: String, $c: String) {
        ...FragA
      }
      fragment FragA on Type {
        field(a: $a) {
          ...FragB
        }
      }
      fragment FragB on Type {
        field(b: $b) {
          ...FragC
        }
      }
      fragment FragC on Type {
        field(c: $c)
      }
    `)
}
func TestValidate_NoUndefinedVariables_VariablesWithinSingleFragmentDefinedInMultipleOperations(t *testing.T) {
	expectPassesRule(t, graphql.NoUndefinedVariablesRule, `
      query Foo($a: String) {
        ...FragA
      }
      query Bar($a: String) {
        ...FragA
      }
      fragment FragA on Type {
        field(a: $a)
      }
    `)
}
func TestValidate_NoUndefinedVariables_VariableWithinFragmentsDefinedInOperations(t *testing.T) {
	expectPassesRule(t, graphql.NoUndefinedVariablesRule, `
      query Foo($a: String) {
        ...FragA
      }
      query Bar($b: String) {
        ...FragB
      }
      fragment FragA on Type {
        field(a: $a)
      }
      fragment FragB on Type {
        field(b: $b)
      }
    `)
}
func TestValidate_NoUndefinedVariables_VariableWithinRecursiveFragmentDefined(t *testing.T) {
	expectPassesRule(t, graphql.NoUndefinedVariablesRule, `
      query Foo($a: String) {
        ...FragA
      }
      fragment FragA on Type {
        field(a: $a) {
          ...FragA
        }
      }
    `)
}
func TestValidate_NoUndefinedVariables_VariableNotDefined(t *testing.T) {
	expectFailsRule(t, graphql.NoUndefinedVariablesRule, `
      query Foo($a: String, $b: String, $c: String) {
        field(a: $a, b: $b, c: $c, d: $d)
      }
    `, []gqlerrors.FormattedError{
		ruleError(`Variable "$d" is not defined.`, 3, 39),
	})
}
func TestValidate_NoUndefinedVariables_VariableNotDefinedByUnnamedQuery(t *testing.T) {
	expectFailsRule(t, graphql.NoUndefinedVariablesRule, `
      {
        field(a: $a)
      }
    `, []gqlerrors.FormattedError{
		ruleError(`Variable "$a" is not defined.`, 3, 18),
	})
}
func TestValidate_NoUndefinedVariables_MultipleVariablesNotDefined(t *testing.T) {
	expectFailsRule(t, graphql.NoUndefinedVariablesRule, `
      query Foo($b: String) {
        field(a: $a, b: $b, c: $c)
      }
    `, []gqlerrors.FormattedError{
		ruleError(`Variable "$a" is not defined.`, 3, 18),
		ruleError(`Variable "$c" is not defined.`, 3, 32),
	})
}

func TestValidate_NoUndefinedVariables_VariableInFragmentNotDefinedByUnnamedQuery(t *testing.T) {
	expectFailsRule(t, graphql.NoUndefinedVariablesRule, `
      {
        ...FragA
      }
      fragment FragA on Type {
        field(a: $a)
      }
    `, []gqlerrors.FormattedError{
		ruleError(`Variable "$a" is not defined.`, 6, 18),
	})
}

func TestValidate_NoUndefinedVariables_VariableInFragmentNotDefinedByOperation(t *testing.T) {
	expectFailsRule(t, graphql.NoUndefinedVariablesRule, `
      query Foo($a: String, $b: String) {
        ...FragA
      }
      fragment FragA on Type {
        field(a: $a) {
          ...FragB
        }
      }
      fragment FragB on Type {
        field(b: $b) {
          ...FragC
        }
      }
      fragment FragC on Type {
        field(c: $c)
      }
    `, []gqlerrors.FormattedError{
		ruleError(`Variable "$c" is not defined by operation "Foo".`, 16, 18, 2, 7),
	})
}

func TestValidate_NoUndefinedVariables_MultipleVariablesInFragmentsNotDefined(t *testing.T) {
	expectFailsRule(t, graphql.NoUndefinedVariablesRule, `
      query Foo($b: String) {
        ...FragA
      }
      fragment FragA on Type {
        field(a: $a) {
          ...FragB
        }
      }
      fragment FragB on Type {
        field(b: $b) {
          ...FragC
        }
      }
      fragment FragC on Type {
        field(c: $c)
      }
    `, []gqlerrors.FormattedError{
		ruleError(`Variable "$a" is not defined by operation "Foo".`, 6, 18, 2, 7),
		ruleError(`Variable "$c" is not defined by operation "Foo".`, 16, 18, 2, 7),
	})
}

func TestValidate_NoUndefinedVariables_SingleVariableInFragmentNotDefinedByMultipleOperations(t *testing.T) {
	expectFailsRule(t, graphql.NoUndefinedVariablesRule, `
      query Foo($a: String) {
        ...FragAB
      }
      query Bar($a: String) {
        ...FragAB
      }
      fragment FragAB on Type {
        field(a: $a, b: $b)
      }
    `, []gqlerrors.FormattedError{
		ruleError(`Variable "$b" is not defined by operation "Foo".`, 9, 25, 2, 7),
		ruleError(`Variable "$b" is not defined by operation "Bar".`, 9, 25, 5, 7),
	})
}

func TestValidate_NoUndefinedVariables_VariablesInFragmentNotDefinedByMultipleOperations(t *testing.T) {
	expectFailsRule(t, graphql.NoUndefinedVariablesRule, `
      query Foo($b: String) {
        ...FragAB
      }
      query Bar($a: String) {
        ...FragAB
      }
      fragment FragAB on Type {
        field(a: $a, b: $b)
      }
    `, []gqlerrors.FormattedError{
		ruleError(`Variable "$a" is not defined by operation "Foo".`, 9, 18, 2, 7),
		ruleError(`Variable "$b" is not defined by operation "Bar".`, 9, 25, 5, 7),
	})
}
func TestValidate_NoUndefinedVariables_VariableInFragmentUsedByOtherOperation(t *testing.T) {
	expectFailsRule(t, graphql.NoUndefinedVariablesRule, `
      query Foo($b: String) {
        ...FragA
      }
      query Bar($a: String) {
        ...FragB
      }
      fragment FragA on Type {
        field(a: $a)
      }
      fragment FragB on Type {
        field(b: $b)
      }
    `, []gqlerrors.FormattedError{
		ruleError(`Variable "$a" is not defined by operation "Foo".`, 9, 18, 2, 7),
		ruleError(`Variable "$b" is not defined by operation "Bar".`, 12, 18, 5, 7),
	})
}

func TestValidate_NoUndefinedVariables_VaMultipleUndefinedVariablesProduceMultipleErrors(t *testing.T) {
	expectFailsRule(t, graphql.NoUndefinedVariablesRule, `
      query Foo($b: String) {
        ...FragAB
      }
      query Bar($a: String) {
        ...FragAB
      }
      fragment FragAB on Type {
        field1(a: $a, b: $b)
        ...FragC
        field3(a: $a, b: $b)
      }
      fragment FragC on Type {
        field2(c: $c)
      }
    `, []gqlerrors.FormattedError{
		ruleError(`Variable "$a" is not defined by operation "Foo".`, 9, 19, 2, 7),
		ruleError(`Variable "$c" is not defined by operation "Foo".`, 14, 19, 2, 7),
		ruleError(`Variable "$a" is not defined by operation "Foo".`, 11, 19, 2, 7),
		ruleError(`Variable "$b" is not defined by operation "Bar".`, 9, 26, 5, 7),
		ruleError(`Variable "$c" is not defined by operation "Bar".`, 14, 19, 5, 7),
		ruleError(`Variable "$b" is not defined by operation "Bar".`, 11, 26, 5, 7),
	})
}
