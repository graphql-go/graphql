package types_test

import (
	"github.com/chris-ramon/graphql"
	"github.com/chris-ramon/graphql/errors"
	"github.com/chris-ramon/graphql/testutil"
	"github.com/chris-ramon/graphql/types"
	"reflect"
	"testing"
)

var enumTypeTestColorType = types.NewGraphQLEnumType(types.GraphQLEnumTypeConfig{
	Name: "Color",
	Values: types.GraphQLEnumValueConfigMap{
		"RED": &types.GraphQLEnumValueConfig{
			Value: 0,
		},
		"GREEN": &types.GraphQLEnumValueConfig{
			Value: 1,
		},
		"BLUE": &types.GraphQLEnumValueConfig{
			Value: 2,
		},
	},
})
var enumTypeTestQueryType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
	Name: "Query",
	Fields: types.GraphQLFieldConfigMap{
		"colorEnum": &types.GraphQLFieldConfig{
			Type: enumTypeTestColorType,
			Args: types.GraphQLFieldConfigArgumentMap{
				"fromEnum": &types.GraphQLArgumentConfig{
					Type: enumTypeTestColorType,
				},
				"fromInt": &types.GraphQLArgumentConfig{
					Type: types.GraphQLInt,
				},
				"fromString": &types.GraphQLArgumentConfig{
					Type: types.GraphQLString,
				},
			},
			Resolve: func(p types.GQLFRParams) interface{} {
				if fromInt, ok := p.Args["fromInt"]; ok {
					return fromInt
				}
				if fromString, ok := p.Args["fromString"]; ok {
					return fromString
				}
				if fromEnum, ok := p.Args["fromEnum"]; ok {
					return fromEnum
				}
				return nil
			},
		},
		"colorInt": &types.GraphQLFieldConfig{
			Type: types.GraphQLInt,
			Args: types.GraphQLFieldConfigArgumentMap{
				"fromEnum": &types.GraphQLArgumentConfig{
					Type: enumTypeTestColorType,
				},
				"fromInt": &types.GraphQLArgumentConfig{
					Type: types.GraphQLInt,
				},
			},
			Resolve: func(p types.GQLFRParams) interface{} {
				if fromInt, ok := p.Args["fromInt"]; ok {
					return fromInt
				}
				if fromEnum, ok := p.Args["fromEnum"]; ok {
					return fromEnum
				}
				return nil
			},
		},
	},
})
var enumTypeTestMutationType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
	Name: "Mutation",
	Fields: types.GraphQLFieldConfigMap{
		"favoriteEnum": &types.GraphQLFieldConfig{
			Type: enumTypeTestColorType,
			Args: types.GraphQLFieldConfigArgumentMap{
				"color": &types.GraphQLArgumentConfig{
					Type: enumTypeTestColorType,
				},
			},
			Resolve: func(p types.GQLFRParams) interface{} {
				if color, ok := p.Args["color"]; ok {
					return color
				}
				return nil
			},
		},
	},
})
var enumTypeTestSchema, _ = types.NewGraphQLSchema(types.GraphQLSchemaConfig{
	Query:    enumTypeTestQueryType,
	Mutation: enumTypeTestMutationType,
})

func executeEnumTypeTest(t *testing.T, query string) *types.GraphQLResult {
	result := g(t, graphql.GraphqlParams{
		Schema:        enumTypeTestSchema,
		RequestString: query,
	})
	return result
}
func executeEnumTypeTestWithParams(t *testing.T, query string, params map[string]interface{}) *types.GraphQLResult {
	result := g(t, graphql.GraphqlParams{
		Schema:         enumTypeTestSchema,
		RequestString:  query,
		VariableValues: params,
	})
	return result
}
func TestTypeSystem_EnumValues_AcceptsEnumLiteralsAsInput(t *testing.T) {
	query := "{ colorInt(fromEnum: GREEN) }"
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"colorInt": 1,
		},
	}
	result := executeEnumTypeTest(t, query)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestTypeSystem_EnumValues_EnumMayBeOutputType(t *testing.T) {
	query := "{ colorEnum(fromInt: 1) }"
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"colorEnum": "GREEN",
		},
	}
	result := executeEnumTypeTest(t, query)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_EnumMayBeBothInputAndOutputType(t *testing.T) {
	query := "{ colorEnum(fromEnum: GREEN) }"
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"colorEnum": "GREEN",
		},
	}
	result := executeEnumTypeTest(t, query)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_DoesNotAcceptStringLiterals(t *testing.T) {
	query := `{ colorEnum(fromEnum: "GREEN") }`
	expected := &types.GraphQLResult{
		Data: nil,
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: `Argument "fromEnum" expected type "Color" but got: "GREEN".`,
			},
		},
	}
	result := executeEnumTypeTest(t, query)
	t.Skipf("Pending `validator` implementation")
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_DoesNotAcceptIncorrectInternalValue(t *testing.T) {
	query := `{ colorEnum(fromString: "GREEN") }`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"colorEnum": nil,
		},
	}
	result := executeEnumTypeTest(t, query)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_DoesNotAcceptInternalValueInPlaceOfEnumLiteral(t *testing.T) {
	query := `{ colorEnum(fromEnum: 1) }`
	expected := &types.GraphQLResult{
		Data: nil,
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: `Argument "fromEnum" expected type "Color" but got: 1.`,
			},
		},
	}
	result := executeEnumTypeTest(t, query)
	t.Skipf("Pending `validator` implementation")
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestTypeSystem_EnumValues_DoesNotAcceptEnumLiteralInPlaceOfInt(t *testing.T) {
	query := `{ colorEnum(fromInt: GREEN) }`
	expected := &types.GraphQLResult{
		Data: nil,
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: `Argument "fromInt" expected type "Int" but got: GREEN.`,
			},
		},
	}
	result := executeEnumTypeTest(t, query)
	t.Skipf("Pending `validator` implementation")
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestTypeSystem_EnumValues_AcceptsJSONStringAsEnumVariable(t *testing.T) {
	query := `query test($color: Color!) { colorEnum(fromEnum: $color) }`
	params := map[string]interface{}{
		"color": "BLUE",
	}
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"colorEnum": "BLUE",
		},
	}
	result := executeEnumTypeTestWithParams(t, query, params)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestTypeSystem_EnumValues_AcceptsEnumLiteralsAsInputArgumentsToMutations(t *testing.T) {
	query := `mutation x($color: Color!) { favoriteEnum(color: $color) }`
	params := map[string]interface{}{
		"color": "GREEN",
	}
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"favoriteEnum": "GREEN",
		},
	}
	result := executeEnumTypeTestWithParams(t, query, params)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_DoesNotAcceptInternalValueAsEnumVariable(t *testing.T) {
	query := `query test($color: Color!) { colorEnum(fromEnum: $color) }`
	params := map[string]interface{}{
		"color": 2,
	}
	expected := &types.GraphQLResult{
		Data: nil,
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: `Variable "$color" expected value of type "Color!" but got: 2.`,
			},
		},
	}
	result := executeEnumTypeTestWithParams(t, query, params)
	t.Skipf("Pending `validator` implementation")
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_DoesNotAcceptStringVariablesAsEnumInput(t *testing.T) {
	query := `query test($color: String!) { colorEnum(fromEnum: $color) }`
	params := map[string]interface{}{
		"color": "BLUE",
	}
	expected := &types.GraphQLResult{
		Data: nil,
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: `Variable "$color" of type "String!" used in position expecting type "Color".`,
			},
		},
	}
	result := executeEnumTypeTestWithParams(t, query, params)
	t.Skipf("Pending `validator` implementation")
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_DoesNotAcceptInternalValueVariableAsEnumInput(t *testing.T) {
	query := `query test($color: Int!) { colorEnum(fromEnum: $color) }`
	params := map[string]interface{}{
		"color": 2,
	}
	expected := &types.GraphQLResult{
		Data: nil,
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: `Variable "$color" of type "Int!" used in position expecting type "Color".`,
			},
		},
	}
	result := executeEnumTypeTestWithParams(t, query, params)
	t.Skipf("Pending `validator` implementation")
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_EnumValueMayHaveAnInternalValueOfZero(t *testing.T) {
	query := `{
        colorEnum(fromEnum: RED)
        colorInt(fromEnum: RED)
      }`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"colorEnum": "RED",
			"colorInt":  0,
		},
	}
	result := executeEnumTypeTest(t, query)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_EnumValueMayBeNullable(t *testing.T) {
	query := `{
        colorEnum
        colorInt
      }`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"colorEnum": nil,
			"colorInt":  nil,
		},
	}
	result := executeEnumTypeTest(t, query)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
