package graphql

import (
	"reflect"
	"testing"

	"github.com/chris-ramon/graphql-go/gqlerrors"
)

var enumTypeTestColorType = NewEnum(EnumConfig{
	Name: "Color",
	Values: EnumValueConfigMap{
		"RED": &EnumValueConfig{
			Value: 0,
		},
		"GREEN": &EnumValueConfig{
			Value: 1,
		},
		"BLUE": &EnumValueConfig{
			Value: 2,
		},
	},
})
var enumTypeTestQueryType = NewObject(ObjectConfig{
	Name: "Query",
	Fields: FieldConfigMap{
		"colorEnum": &FieldConfig{
			Type: enumTypeTestColorType,
			Args: FieldConfigArgument{
				"fromEnum": &ArgumentConfig{
					Type: enumTypeTestColorType,
				},
				"fromInt": &ArgumentConfig{
					Type: Int,
				},
				"fromString": &ArgumentConfig{
					Type: String,
				},
			},
			Resolve: func(p GQLFRParams) interface{} {
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
		"colorInt": &FieldConfig{
			Type: Int,
			Args: FieldConfigArgument{
				"fromEnum": &ArgumentConfig{
					Type: enumTypeTestColorType,
				},
				"fromInt": &ArgumentConfig{
					Type: Int,
				},
			},
			Resolve: func(p GQLFRParams) interface{} {
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
var enumTypeTestMutationType = NewObject(ObjectConfig{
	Name: "Mutation",
	Fields: FieldConfigMap{
		"favoriteEnum": &FieldConfig{
			Type: enumTypeTestColorType,
			Args: FieldConfigArgument{
				"color": &ArgumentConfig{
					Type: enumTypeTestColorType,
				},
			},
			Resolve: func(p GQLFRParams) interface{} {
				if color, ok := p.Args["color"]; ok {
					return color
				}
				return nil
			},
		},
	},
})
var enumTypeTestSchema, _ = NewSchema(SchemaConfig{
	Query:    enumTypeTestQueryType,
	Mutation: enumTypeTestMutationType,
})

func executeEnumTypeTest(t *testing.T, query string) *Result {
	result := g(t, Params{
		Schema:        enumTypeTestSchema,
		RequestString: query,
	})
	return result
}
func executeEnumTypeTestWithParams(t *testing.T, query string, params map[string]interface{}) *Result {
	result := g(t, Params{
		Schema:         enumTypeTestSchema,
		RequestString:  query,
		VariableValues: params,
	})
	return result
}
func TestTypeSystem_EnumValues_AcceptsEnumLiteralsAsInput(t *testing.T) {
	query := "{ colorInt(fromEnum: GREEN) }"
	expected := &Result{
		Data: map[string]interface{}{
			"colorInt": 1,
		},
	}
	result := executeEnumTypeTest(t, query)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}

func TestTypeSystem_EnumValues_EnumMayBeOutputType(t *testing.T) {
	query := "{ colorEnum(fromInt: 1) }"
	expected := &Result{
		Data: map[string]interface{}{
			"colorEnum": "GREEN",
		},
	}
	result := executeEnumTypeTest(t, query)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_EnumMayBeBothInputAndOutputType(t *testing.T) {
	query := "{ colorEnum(fromEnum: GREEN) }"
	expected := &Result{
		Data: map[string]interface{}{
			"colorEnum": "GREEN",
		},
	}
	result := executeEnumTypeTest(t, query)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_DoesNotAcceptStringLiterals(t *testing.T) {
	query := `{ colorEnum(fromEnum: "GREEN") }`
	expected := &Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: `Argument "fromEnum" expected type "Color" but got: "GREEN".`,
			},
		},
	}
	result := executeEnumTypeTest(t, query)
	t.Skipf("Pending `validator` implementation")
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_DoesNotAcceptIncorrectInternalValue(t *testing.T) {
	query := `{ colorEnum(fromString: "GREEN") }`
	expected := &Result{
		Data: map[string]interface{}{
			"colorEnum": nil,
		},
	}
	result := executeEnumTypeTest(t, query)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_DoesNotAcceptInternalValueInPlaceOfEnumLiteral(t *testing.T) {
	query := `{ colorEnum(fromEnum: 1) }`
	expected := &Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: `Argument "fromEnum" expected type "Color" but got: 1.`,
			},
		},
	}
	result := executeEnumTypeTest(t, query)
	t.Skipf("Pending `validator` implementation")
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}

func TestTypeSystem_EnumValues_DoesNotAcceptEnumLiteralInPlaceOfInt(t *testing.T) {
	query := `{ colorEnum(fromInt: GREEN) }`
	expected := &Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: `Argument "fromInt" expected type "Int" but got: GREEN.`,
			},
		},
	}
	result := executeEnumTypeTest(t, query)
	t.Skipf("Pending `validator` implementation")
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}

func TestTypeSystem_EnumValues_AcceptsJSONStringAsEnumVariable(t *testing.T) {
	query := `query test($color: Color!) { colorEnum(fromEnum: $color) }`
	params := map[string]interface{}{
		"color": "BLUE",
	}
	expected := &Result{
		Data: map[string]interface{}{
			"colorEnum": "BLUE",
		},
	}
	result := executeEnumTypeTestWithParams(t, query, params)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}

func TestTypeSystem_EnumValues_AcceptsEnumLiteralsAsInputArgumentsToMutations(t *testing.T) {
	query := `mutation x($color: Color!) { favoriteEnum(color: $color) }`
	params := map[string]interface{}{
		"color": "GREEN",
	}
	expected := &Result{
		Data: map[string]interface{}{
			"favoriteEnum": "GREEN",
		},
	}
	result := executeEnumTypeTestWithParams(t, query, params)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_DoesNotAcceptInternalValueAsEnumVariable(t *testing.T) {
	query := `query test($color: Color!) { colorEnum(fromEnum: $color) }`
	params := map[string]interface{}{
		"color": 2,
	}
	expected := &Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: `Variable "$color" expected value of type "Color!" but got: 2.`,
			},
		},
	}
	result := executeEnumTypeTestWithParams(t, query, params)
	t.Skipf("Pending `validator` implementation")
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_DoesNotAcceptStringVariablesAsEnumInput(t *testing.T) {
	query := `query test($color: String!) { colorEnum(fromEnum: $color) }`
	params := map[string]interface{}{
		"color": "BLUE",
	}
	expected := &Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: `Variable "$color" of type "String!" used in position expecting type "Color".`,
			},
		},
	}
	result := executeEnumTypeTestWithParams(t, query, params)
	t.Skipf("Pending `validator` implementation")
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_DoesNotAcceptInternalValueVariableAsEnumInput(t *testing.T) {
	query := `query test($color: Int!) { colorEnum(fromEnum: $color) }`
	params := map[string]interface{}{
		"color": 2,
	}
	expected := &Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: `Variable "$color" of type "Int!" used in position expecting type "Color".`,
			},
		},
	}
	result := executeEnumTypeTestWithParams(t, query, params)
	t.Skipf("Pending `validator` implementation")
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_EnumValueMayHaveAnInternalValueOfZero(t *testing.T) {
	query := `{
        colorEnum(fromEnum: RED)
        colorInt(fromEnum: RED)
      }`
	expected := &Result{
		Data: map[string]interface{}{
			"colorEnum": "RED",
			"colorInt":  0,
		},
	}
	result := executeEnumTypeTest(t, query)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_EnumValueMayBeNullable(t *testing.T) {
	query := `{
        colorEnum
        colorInt
      }`
	expected := &Result{
		Data: map[string]interface{}{
			"colorEnum": nil,
			"colorInt":  nil,
		},
	}
	result := executeEnumTypeTest(t, query)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
