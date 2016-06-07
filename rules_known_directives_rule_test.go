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
      query Foo @onQuery {
        name @include(if: true)
        ...Frag @include(if: true)
        skippedField @skip(if: true)
        ...SkippedFrag @skip(if: true)
      }

      mutation Bar @onMutation {
        someField
      }
    `)
}
func TestValidate_KnownDirectives_WithMisplacedDirectives(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.KnownDirectivesRule, `
      query Foo @include(if: true) {
        name @onQuery
        ...Frag @onQuery
      }

      mutation Bar @onQuery {
        someField
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Directive "include" may not be used on QUERY.`, 2, 17),
		testutil.RuleError(`Directive "onQuery" may not be used on FIELD.`, 3, 14),
		testutil.RuleError(`Directive "onQuery" may not be used on FRAGMENT_SPREAD.`, 4, 17),
		testutil.RuleError(`Directive "onQuery" may not be used on MUTATION.`, 7, 20),
	})
}

func TestValidate_KnownDirectives_WithinSchemaLanguage_WithWellPlacedDirectives(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.KnownDirectivesRule, `
        type MyObj implements MyInterface @onObject {
          myField(myArg: Int @onArgumentDefinition): String @onFieldDefinition
        }

        scalar MyScalar @onScalar

        interface MyInterface @onInterface {
          myField(myArg: Int @onArgumentDefinition): String @onFieldDefinition
        }

        union MyUnion @onUnion = MyObj | Other

        enum MyEnum @onEnum {
          MY_VALUE @onEnumValue
        }

        input MyInput @onInputObject {
          myField: Int @onInputFieldDefinition
        }

        schema @onSchema {
          query: MyQuery
        }
    `)
}

func TestValidate_KnownDirectives_WithinSchemaLanguage_WithMisplacedDirectives(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.KnownDirectivesRule, `
        type MyObj implements MyInterface @onInterface {
          myField(myArg: Int @onInputFieldDefinition): String @onInputFieldDefinition
        }

        scalar MyScalar @onEnum

        interface MyInterface @onObject {
          myField(myArg: Int @onInputFieldDefinition): String @onInputFieldDefinition
        }

        union MyUnion @onEnumValue = MyObj | Other

        enum MyEnum @onScalar {
          MY_VALUE @onUnion
        }

        input MyInput @onEnum {
          myField: Int @onArgumentDefinition
        }

        schema @onObject {
          query: MyQuery
        }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Directive "onInterface" may not be used on OBJECT.`, 2, 43),
		testutil.RuleError(`Directive "onInputFieldDefinition" may not be used on ARGUMENT_DEFINITION.`, 3, 30),
		testutil.RuleError(`Directive "onInputFieldDefinition" may not be used on FIELD_DEFINITION.`, 3, 63),
		testutil.RuleError(`Directive "onEnum" may not be used on SCALAR.`, 6, 25),
		testutil.RuleError(`Directive "onObject" may not be used on INTERFACE.`, 8, 31),
		testutil.RuleError(`Directive "onInputFieldDefinition" may not be used on ARGUMENT_DEFINITION.`, 9, 30),
		testutil.RuleError(`Directive "onInputFieldDefinition" may not be used on FIELD_DEFINITION.`, 9, 63),
		testutil.RuleError(`Directive "onEnumValue" may not be used on UNION.`, 12, 23),
		testutil.RuleError(`Directive "onScalar" may not be used on ENUM.`, 14, 21),
		testutil.RuleError(`Directive "onUnion" may not be used on ENUM_VALUE.`, 15, 20),
		testutil.RuleError(`Directive "onEnum" may not be used on INPUT_OBJECT.`, 18, 23),
		testutil.RuleError(`Directive "onArgumentDefinition" may not be used on INPUT_FIELD_DEFINITION.`, 19, 24),
		testutil.RuleError(`Directive "onObject" may not be used on SCHEMA.`, 22, 16),
	})
}
