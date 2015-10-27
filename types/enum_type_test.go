package types_test

import (
	"github.com/chris-ramon/graphql"
	"github.com/chris-ramon/graphql/errors"
	"github.com/chris-ramon/graphql/testutil"
	"github.com/chris-ramon/graphql/types"
	"reflect"
	"testing"
)

var enumTypeTestColorType = types.NewEnum(types.EnumConfig{
	Name: "Color",
	Values: types.EnumValueConfigMap{
		"RED": &types.EnumValueConfig{
			Value: 0,
		},
		"GREEN": &types.EnumValueConfig{
			Value: 1,
		},
		"BLUE": &types.EnumValueConfig{
			Value: 2,
		},
	},
})
var enumTypeTestQueryType = types.NewObject(types.ObjectConfig{
	Name: "Query",
	Fields: types.FieldConfigMap{
		"colorEnum": &types.FieldConfig{
			Type: enumTypeTestColorType,
			Args: types.FieldConfigArgument{
				"fromEnum": &types.ArgumentConfig{
					Type: enumTypeTestColorType,
				},
				"fromInt": &types.ArgumentConfig{
					Type: types.Int,
				},
				"fromString": &types.ArgumentConfig{
					Type: types.String,
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
		"colorInt": &types.FieldConfig{
			Type: types.Int,
			Args: types.FieldConfigArgument{
				"fromEnum": &types.ArgumentConfig{
					Type: enumTypeTestColorType,
				},
				"fromInt": &types.ArgumentConfig{
					Type: types.Int,
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
var enumTypeTestMutationType = types.NewObject(types.ObjectConfig{
	Name: "Mutation",
	Fields: types.FieldConfigMap{
		"favoriteEnum": &types.FieldConfig{
			Type: enumTypeTestColorType,
			Args: types.FieldConfigArgument{
				"color": &types.ArgumentConfig{
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
var enumTypeTestSchema, _ = types.NewSchema(types.SchemaConfig{
	Query:    enumTypeTestQueryType,
	Mutation: enumTypeTestMutationType,
})

func executeEnumTypeTest(t *testing.T, query string) *types.Result {
	result := g(t, graphql.Params{
		Schema:        enumTypeTestSchema,
		RequestString: query,
	})
	return result
}
func executeEnumTypeTestWithParams(t *testing.T, query string, params map[string]interface{}) *types.Result {
	result := g(t, graphql.Params{
		Schema:         enumTypeTestSchema,
		RequestString:  query,
		VariableValues: params,
	})
	return result
}
func TestTypeSystem_EnumValues_AcceptsEnumLiteralsAsInput(t *testing.T) {
	query := "{ colorInt(fromEnum: GREEN) }"
	expected := &types.Result{
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
	expected := &types.Result{
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
	expected := &types.Result{
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
	expected := &types.Result{
		Data: nil,
		Errors: []graphqlerrors.FormattedError{
			graphqlerrors.FormattedError{
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
	expected := &types.Result{
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
	expected := &types.Result{
		Data: nil,
		Errors: []graphqlerrors.FormattedError{
			graphqlerrors.FormattedError{
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
	expected := &types.Result{
		Data: nil,
		Errors: []graphqlerrors.FormattedError{
			graphqlerrors.FormattedError{
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
	expected := &types.Result{
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
	expected := &types.Result{
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
	expected := &types.Result{
		Data: nil,
		Errors: []graphqlerrors.FormattedError{
			graphqlerrors.FormattedError{
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
	expected := &types.Result{
		Data: nil,
		Errors: []graphqlerrors.FormattedError{
			graphqlerrors.FormattedError{
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
	expected := &types.Result{
		Data: nil,
		Errors: []graphqlerrors.FormattedError{
			graphqlerrors.FormattedError{
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
	expected := &types.Result{
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
	expected := &types.Result{
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
