package graphql_test

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/testutil"
)

func TestValidate_FieldsOnCorrectType_ObjectFieldSelection(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment objectFieldSelection on Dog {
        __typename
        name
      }
    `)
}
func TestValidate_FieldsOnCorrectType_AliasedObjectFieldSelection(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment aliasedObjectFieldSelection on Dog {
        tn : __typename
        otherName : name
      }
    `)
}
func TestValidate_FieldsOnCorrectType_InterfaceFieldSelection(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment interfaceFieldSelection on Pet {
        __typename
        name
      }
    `)
}
func TestValidate_FieldsOnCorrectType_AliasedInterfaceFieldSelection(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment interfaceFieldSelection on Pet {
        otherName : name
      }
    `)
}
func TestValidate_FieldsOnCorrectType_LyingAliasSelection(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment lyingAliasSelection on Dog {
        name : nickname
      }
    `)
}
func TestValidate_FieldsOnCorrectType_IgnoresFieldsOnUnknownType(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment unknownSelection on UnknownType {
        unknownField
      }
    `)
}
func TestValidate_FieldsOnCorrectType_FieldNotDefinedOnFragment(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment fieldNotDefined on Dog {
        meowVolume
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Cannot query field "meowVolume" on "Dog".`, 3, 9),
	})
}
func TestValidate_FieldsOnCorrectType_FieldNotDefinedDeeplyOnlyReportsFirst(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment deepFieldNotDefined on Dog {
        unknown_field {
          deeper_unknown_field
        }
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Cannot query field "unknown_field" on "Dog".`, 3, 9),
	})
}
func TestValidate_FieldsOnCorrectType_SubFieldNotDefined(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment subFieldNotDefined on Human {
        pets {
          unknown_field
        }
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Cannot query field "unknown_field" on "Pet".`, 4, 11),
	})
}
func TestValidate_FieldsOnCorrectType_FieldNotDefinedOnInlineFragment(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment fieldNotDefined on Pet {
        ... on Dog {
          meowVolume
        }
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Cannot query field "meowVolume" on "Dog".`, 4, 11),
	})
}
func TestValidate_FieldsOnCorrectType_AliasedFieldTargetNotDefined(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment aliasedFieldTargetNotDefined on Dog {
        volume : mooVolume
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Cannot query field "mooVolume" on "Dog".`, 3, 9),
	})
}
func TestValidate_FieldsOnCorrectType_AliasedLyingFieldTargetNotDefined(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment aliasedLyingFieldTargetNotDefined on Dog {
        barkVolume : kawVolume
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Cannot query field "kawVolume" on "Dog".`, 3, 9),
	})
}
func TestValidate_FieldsOnCorrectType_NotDefinedOnInterface(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment notDefinedOnInterface on Pet {
        tailLength
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Cannot query field "tailLength" on "Pet".`, 3, 9),
	})
}
func TestValidate_FieldsOnCorrectType_DefinedOnImplementorsButNotOnInterface(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment definedOnImplementorsButNotInterface on Pet {
        nickname
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Cannot query field "nickname" on "Pet".`, 3, 9),
	})
}
func TestValidate_FieldsOnCorrectType_MetaFieldSelectionOnUnion(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment directFieldSelectionOnUnion on CatOrDog {
        __typename
      }
    `)
}
func TestValidate_FieldsOnCorrectType_DirectFieldSelectionOnUnion(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment directFieldSelectionOnUnion on CatOrDog {
        directField
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Cannot query field "directField" on "CatOrDog".`, 3, 9),
	})
}
func TestValidate_FieldsOnCorrectType_DirectImplementorsQueriedOnUnion(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment definedOnImplementorsQueriedOnUnion on CatOrDog {
        name
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Cannot query field "name" on "CatOrDog".`, 3, 9),
	})
}
func TestValidate_FieldsOnCorrectType_ValidFieldInInlineFragment(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment objectFieldSelection on Pet {
        ... on Dog {
          name
        }
      }
    `)
}
