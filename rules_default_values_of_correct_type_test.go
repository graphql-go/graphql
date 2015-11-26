package graphql_test

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/testutil"
)

func TestValidate_VariableDefaultValuesOfCorrectType_VariablesWithNoDefaultValues(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.DefaultValuesOfCorrectTypeRule, `
      query NullableValues($a: Int, $b: String, $c: ComplexInput) {
        dog { name }
      }
    `)
}
func TestValidate_VariableDefaultValuesOfCorrectType_RequiredVariablesWithoutDefaultValues(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.DefaultValuesOfCorrectTypeRule, `
      query RequiredValues($a: Int!, $b: String!) {
        dog { name }
      }
    `)
}
func TestValidate_VariableDefaultValuesOfCorrectType_VariablesWithValidDefaultValues(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.DefaultValuesOfCorrectTypeRule, `
      query WithDefaultValues(
        $a: Int = 1,
        $b: String = "ok",
        $c: ComplexInput = { requiredField: true, intField: 3 }
      ) {
        dog { name }
      }
    `)
}
func TestValidate_VariableDefaultValuesOfCorrectType_NoRequiredVariablesWithDefaultValues(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.DefaultValuesOfCorrectTypeRule, `
      query UnreachableDefaultValues($a: Int! = 3, $b: String! = "default") {
        dog { name }
      }
    `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				`Variable "$a" of type "Int!" is required and will not `+
					`use the default value. Perhaps you meant to use type "Int".`,
				2, 49,
			),
			testutil.RuleError(
				`Variable "$b" of type "String!" is required and will not `+
					`use the default value. Perhaps you meant to use type "String".`,
				2, 66,
			),
		})
}
func TestValidate_VariableDefaultValuesOfCorrectType_VariablesWithInvalidDefaultValues(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.DefaultValuesOfCorrectTypeRule, `
      query InvalidDefaultValues(
        $a: Int = "one",
        $b: String = 4,
        $c: ComplexInput = "notverycomplex"
      ) {
        dog { name }
      }
    `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(`Variable "$a" of type "Int" has invalid default value: "one".`, 3, 19),
			testutil.RuleError(`Variable "$b" of type "String" has invalid default value: 4.`, 4, 22),
			testutil.RuleError(`Variable "$c" of type "ComplexInput" has invalid default value: "notverycomplex".`, 5, 28),
		})
}
func TestValidate_VariableDefaultValuesOfCorrectType_ComplexVariablesMissingRequiredField(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.DefaultValuesOfCorrectTypeRule, `
      query MissingRequiredField($a: ComplexInput = {intField: 3}) {
        dog { name }
      }
    `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(`Variable "$a" of type "ComplexInput" has invalid default value: {intField: 3}.`, 2, 53),
		})
}
func TestValidate_VariableDefaultValuesOfCorrectType_ListVariablesWithInvalidItem(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.DefaultValuesOfCorrectTypeRule, `
      query InvalidItem($a: [String] = ["one", 2]) {
        dog { name }
      }
    `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(`Variable "$a" of type "[String]" has invalid default value: ["one", 2].`, 2, 40),
		})
}
