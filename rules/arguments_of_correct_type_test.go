package rules_test

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
)

func TestValidate_ArgValuesOfCorrectType_ValidValue_GoodIntValue(t *testing.T) {
	expectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            intArgField(intArg: 2)
          }
        }
    `)
}
func TestValidate_ArgValuesOfCorrectType_ValidValue_GoodBooleanValue(t *testing.T) {
	expectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            booleanArgField(booleanArg: true)
          }
        }
    `)
}
func TestValidate_ArgValuesOfCorrectType_ValidValue_GoodStringValue(t *testing.T) {
	expectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            stringArgField(stringArg: "foo")
          }
        }
    `)
}
func TestValidate_ArgValuesOfCorrectType_ValidValue_GoodFloatValue(t *testing.T) {
	expectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            floatArgField(floatArg: 1.1)
          }
        }
    `)
}
func TestValidate_ArgValuesOfCorrectType_ValidValue_IntIntoFloat(t *testing.T) {
	expectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            floatArgField(floatArg: 1)
          }
        }
    `)
}
func TestValidate_ArgValuesOfCorrectType_ValidValue_IntIntoID(t *testing.T) {
	expectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            idArgField(idArg: 1)
          }
        }
    `)
}
func TestValidate_ArgValuesOfCorrectType_ValidValue_StringIntoID(t *testing.T) {
	expectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            idArgField(idArg: "someIdString")
          }
        }
    `)
}
func TestValidate_ArgValuesOfCorrectType_ValidValue_GoodEnumValue(t *testing.T) {
	expectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          dog {
            doesKnowCommand(dogCommand: SIT)
          }
        }
    `)
}

func TestValidate_ArgValuesOfCorrectType_InvalidStringValues_IntIntoString(t *testing.T) {
	expectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            stringArgField(stringArg: 1)
          }
        }
        `,
		[]gqlerrors.FormattedError{
			ruleError(
				`Argument "stringArg" expected type "String" but got: 1.`,
				4, 39,
			),
		})
}
func TestValidate_ArgValuesOfCorrectType_InvalidStringValues_FloatIntoString(t *testing.T) {
	expectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            stringArgField(stringArg: 1.0)
          }
        }
        `,
		[]gqlerrors.FormattedError{
			ruleError(
				`Argument "stringArg" expected type "String" but got: 1.0.`,
				4, 39,
			),
		})
}
