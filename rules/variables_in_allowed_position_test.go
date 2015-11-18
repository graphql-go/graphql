package rules_test

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
)

func TestValidate_VariablesInAllowedPosition_BooleanToBoolean(t *testing.T) {
	expectPassesRule(t, graphql.VariablesInAllowedPositionRule, `
      query Query($booleanArg: Boolean)
      {
        complicatedArgs {
          booleanArgField(booleanArg: $booleanArg)
        }
      }
    `)
}
func TestValidate_VariablesInAllowedPosition_BooleanToBooleanWithinFragment(t *testing.T) {
	expectPassesRule(t, graphql.VariablesInAllowedPositionRule, `
      fragment booleanArgFrag on ComplicatedArgs {
        booleanArgField(booleanArg: $booleanArg)
      }
      query Query($booleanArg: Boolean)
      {
        complicatedArgs {
          ...booleanArgFrag
        }
      }
    `)
	expectPassesRule(t, graphql.VariablesInAllowedPositionRule, `
      query Query($booleanArg: Boolean)
      {
        complicatedArgs {
          ...booleanArgFrag
        }
      }
      fragment booleanArgFrag on ComplicatedArgs {
        booleanArgField(booleanArg: $booleanArg)
      }
    `)
}
func TestValidate_VariablesInAllowedPosition_NonNullableBooleanToBoolean(t *testing.T) {
	expectPassesRule(t, graphql.VariablesInAllowedPositionRule, `
      query Query($nonNullBooleanArg: Boolean!)
      {
        complicatedArgs {
          booleanArgField(booleanArg: $nonNullBooleanArg)
        }
      }
    `)
}
func TestValidate_VariablesInAllowedPosition_NonNullableBooleanToBooleanWithinFragment(t *testing.T) {
	expectPassesRule(t, graphql.VariablesInAllowedPositionRule, `
      fragment booleanArgFrag on ComplicatedArgs {
        booleanArgField(booleanArg: $nonNullBooleanArg)
      }

      query Query($nonNullBooleanArg: Boolean!)
      {
        complicatedArgs {
          ...booleanArgFrag
        }
      }
    `)
}
func TestValidate_VariablesInAllowedPosition_IntToNonNullableIntWithDefault(t *testing.T) {
	expectPassesRule(t, graphql.VariablesInAllowedPositionRule, `
      query Query($intArg: Int = 1)
      {
        complicatedArgs {
          nonNullIntArgField(nonNullIntArg: $intArg)
        }
      }
    `)
}
func TestValidate_VariablesInAllowedPosition_ListOfStringToListOfString(t *testing.T) {
	expectPassesRule(t, graphql.VariablesInAllowedPositionRule, `
      query Query($stringListVar: [String])
      {
        complicatedArgs {
          stringListArgField(stringListArg: $stringListVar)
        }
      }
    `)
}
func TestValidate_VariablesInAllowedPosition_ListOfNonNullableStringToListOfString(t *testing.T) {
	expectPassesRule(t, graphql.VariablesInAllowedPositionRule, `
      query Query($stringListVar: [String!])
      {
        complicatedArgs {
          stringListArgField(stringListArg: $stringListVar)
        }
      }
    `)
}
func TestValidate_VariablesInAllowedPosition_StringToListOfStringInItemPosition(t *testing.T) {
	expectPassesRule(t, graphql.VariablesInAllowedPositionRule, `
      query Query($stringVar: String)
      {
        complicatedArgs {
          stringListArgField(stringListArg: [$stringVar])
        }
      }
    `)
}
func TestValidate_VariablesInAllowedPosition_NonNullableStringToListOfStringInItemPosition(t *testing.T) {
	expectPassesRule(t, graphql.VariablesInAllowedPositionRule, `
      query Query($stringVar: String!)
      {
        complicatedArgs {
          stringListArgField(stringListArg: [$stringVar])
        }
      }
    `)
}
func TestValidate_VariablesInAllowedPosition_ComplexInputToComplexInput(t *testing.T) {
	expectPassesRule(t, graphql.VariablesInAllowedPositionRule, `
      query Query($complexVar: ComplexInput)
      {
        complicatedArgs {
          complexArgField(complexArg: $ComplexInput)
        }
      }
    `)
}
func TestValidate_VariablesInAllowedPosition_ComplexInputToComplexInputInFieldPosition(t *testing.T) {
	expectPassesRule(t, graphql.VariablesInAllowedPositionRule, `
      query Query($boolVar: Boolean = false)
      {
        complicatedArgs {
          complexArgField(complexArg: {requiredArg: $boolVar})
        }
      }
    `)
}
func TestValidate_VariablesInAllowedPosition_NonNullableBooleanToNonNullableBooleanInDirective(t *testing.T) {
	expectPassesRule(t, graphql.VariablesInAllowedPositionRule, `
      query Query($boolVar: Boolean!)
      {
        dog @include(if: $boolVar)
      }
    `)
}
func TestValidate_VariablesInAllowedPosition_NonNullableBooleanToNonNullableBooleanInDirectiveInDirectiveWithDefault(t *testing.T) {
	expectPassesRule(t, graphql.VariablesInAllowedPositionRule, `
      query Query($boolVar: Boolean = false)
      {
        dog @include(if: $boolVar)
      }
    `)
}
func TestValidate_VariablesInAllowedPosition_IntToNonNullableInt(t *testing.T) {
	expectFailsRule(t, graphql.VariablesInAllowedPositionRule, `
      query Query($intArg: Int)
      {
        complicatedArgs {
          nonNullIntArgField(nonNullIntArg: $intArg)
        }
      }
    `, []gqlerrors.FormattedError{
		ruleError(`Variable "$intArg" of type "Int" used in position `+
			`expecting type "Int!".`, 5, 45),
	})
}
func TestValidate_VariablesInAllowedPosition_IntToNonNullableIntWithinFragment(t *testing.T) {
	expectFailsRule(t, graphql.VariablesInAllowedPositionRule, `
      fragment nonNullIntArgFieldFrag on ComplicatedArgs {
        nonNullIntArgField(nonNullIntArg: $intArg)
      }

      query Query($intArg: Int)
      {
        complicatedArgs {
          ...nonNullIntArgFieldFrag
        }
      }
    `, []gqlerrors.FormattedError{
		ruleError(`Variable "$intArg" of type "Int" used in position `+
			`expecting type "Int!".`, 3, 43),
	})
}
func TestValidate_VariablesInAllowedPosition_IntToNonNullableIntWithinNestedFragment(t *testing.T) {
	expectFailsRule(t, graphql.VariablesInAllowedPositionRule, `
      fragment outerFrag on ComplicatedArgs {
        ...nonNullIntArgFieldFrag
      }

      fragment nonNullIntArgFieldFrag on ComplicatedArgs {
        nonNullIntArgField(nonNullIntArg: $intArg)
      }

      query Query($intArg: Int)
      {
        complicatedArgs {
          ...outerFrag
        }
      }
    `, []gqlerrors.FormattedError{
		ruleError(`Variable "$intArg" of type "Int" used in position `+
			`expecting type "Int!".`, 7, 43),
	})
}
func TestValidate_VariablesInAllowedPosition_StringOverBoolean(t *testing.T) {
	expectFailsRule(t, graphql.VariablesInAllowedPositionRule, `
      query Query($stringVar: String)
      {
        complicatedArgs {
          booleanArgField(booleanArg: $stringVar)
        }
      }
    `, []gqlerrors.FormattedError{
		ruleError(`Variable "$stringVar" of type "String" used in position `+
			`expecting type "Boolean".`, 5, 39),
	})
}
func TestValidate_VariablesInAllowedPosition_StringToListOfString(t *testing.T) {
	expectFailsRule(t, graphql.VariablesInAllowedPositionRule, `
      query Query($stringVar: String)
      {
        complicatedArgs {
          stringListArgField(stringListArg: $stringVar)
        }
      }
    `, []gqlerrors.FormattedError{
		ruleError(`Variable "$stringVar" of type "String" used in position `+
			`expecting type "[String]".`, 5, 45),
	})
}
func TestValidate_VariablesInAllowedPosition_BooleanToNonNullableBooleanInDirective(t *testing.T) {
	expectFailsRule(t, graphql.VariablesInAllowedPositionRule, `
      query Query($boolVar: Boolean)
      {
        dog @include(if: $boolVar)
      }
    `, []gqlerrors.FormattedError{
		ruleError(`Variable "$boolVar" of type "Boolean" used in position `+
			`expecting type "Boolean!".`, 4, 26),
	})
}
func TestValidate_VariablesInAllowedPosition_StringToNonNullableBooleanInDirective(t *testing.T) {
	expectFailsRule(t, graphql.VariablesInAllowedPositionRule, `
      query Query($stringVar: String)
      {
        dog @include(if: $stringVar)
      }
    `, []gqlerrors.FormattedError{
		ruleError(`Variable "$stringVar" of type "String" used in position `+
			`expecting type "Boolean!".`, 4, 26),
	})
}
