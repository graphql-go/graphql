package graphql_test

import (
	"reflect"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/testutil"
)

var enumTypeTestColorType = graphql.NewEnum(graphql.EnumConfig{
	Name: "Color",
	Values: graphql.EnumValueConfigMap{
		"RED": &graphql.EnumValueConfig{
			Value: 0,
		},
		"GREEN": &graphql.EnumValueConfig{
			Value: 1,
		},
		"BLUE": &graphql.EnumValueConfig{
			Value: 2,
		},
	},
})
var enumTypeTestQueryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"colorEnum": &graphql.Field{
			Type: enumTypeTestColorType,
			Args: graphql.FieldConfigArgument{
				"fromEnum": &graphql.ArgumentConfig{
					Type: enumTypeTestColorType,
				},
				"fromInt": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"fromString": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if fromInt, ok := p.Args["fromInt"]; ok {
					return fromInt, nil
				}
				if fromString, ok := p.Args["fromString"]; ok {
					return fromString, nil
				}
				if fromEnum, ok := p.Args["fromEnum"]; ok {
					return fromEnum, nil
				}
				return nil, nil
			},
		},
		"colorInt": &graphql.Field{
			Type: graphql.Int,
			Args: graphql.FieldConfigArgument{
				"fromEnum": &graphql.ArgumentConfig{
					Type: enumTypeTestColorType,
				},
				"fromInt": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if fromInt, ok := p.Args["fromInt"]; ok {
					return fromInt, nil
				}
				if fromEnum, ok := p.Args["fromEnum"]; ok {
					return fromEnum, nil
				}
				return nil, nil
			},
		},
	},
})
var enumTypeTestMutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"favoriteEnum": &graphql.Field{
			Type: enumTypeTestColorType,
			Args: graphql.FieldConfigArgument{
				"color": &graphql.ArgumentConfig{
					Type: enumTypeTestColorType,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if color, ok := p.Args["color"]; ok {
					return color, nil
				}
				return nil, nil
			},
		},
	},
})
var enumTypeTestSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    enumTypeTestQueryType,
	Mutation: enumTypeTestMutationType,
})

func executeEnumTypeTest(t *testing.T, query string) *graphql.Result {
	result := g(t, graphql.Params{
		Schema:        enumTypeTestSchema,
		RequestString: query,
	})
	return result
}
func executeEnumTypeTestWithParams(t *testing.T, query string, params map[string]interface{}) *graphql.Result {
	result := g(t, graphql.Params{
		Schema:         enumTypeTestSchema,
		RequestString:  query,
		VariableValues: params,
	})
	return result
}
func TestTypeSystem_EnumValues_AcceptsEnumLiteralsAsInput(t *testing.T) {
	query := "{ colorInt(fromEnum: GREEN) }"
	expected := &graphql.Result{
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
	expected := &graphql.Result{
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
	expected := &graphql.Result{
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
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: `Argument "fromEnum" expected type "Color" but got: "GREEN".`,
			},
		},
	}
	result := executeEnumTypeTest(t, query)
	if !testutil.EqualErrorMessage(expected, result, 0) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_DoesNotAcceptIncorrectInternalValue(t *testing.T) {
	query := `{ colorEnum(fromString: "GREEN") }`
	expected := &graphql.Result{
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
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: `Argument "fromEnum" expected type "Color" but got: 1.`,
			},
		},
	}
	result := executeEnumTypeTest(t, query)
	if !testutil.EqualErrorMessage(expected, result, 0) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestTypeSystem_EnumValues_DoesNotAcceptEnumLiteralInPlaceOfInt(t *testing.T) {
	query := `{ colorEnum(fromInt: GREEN) }`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: `Argument "fromInt" expected type "Int" but got: GREEN.`,
			},
		},
	}
	result := executeEnumTypeTest(t, query)
	if !testutil.EqualErrorMessage(expected, result, 0) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestTypeSystem_EnumValues_AcceptsJSONStringAsEnumVariable(t *testing.T) {
	query := `query test($color: Color!) { colorEnum(fromEnum: $color) }`
	params := map[string]interface{}{
		"color": "BLUE",
	}
	expected := &graphql.Result{
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
	expected := &graphql.Result{
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
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: `Variable "$color" expected value of type "Color!" but got: 2.`,
			},
		},
	}
	result := executeEnumTypeTestWithParams(t, query, params)
	if !testutil.EqualErrorMessage(expected, result, 0) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_DoesNotAcceptStringVariablesAsEnumInput(t *testing.T) {
	query := `query test($color: String!) { colorEnum(fromEnum: $color) }`
	params := map[string]interface{}{
		"color": "BLUE",
	}
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: `Variable "$color" of type "String!" used in position expecting type "Color".`,
			},
		},
	}
	result := executeEnumTypeTestWithParams(t, query, params)
	if !testutil.EqualErrorMessage(expected, result, 0) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_DoesNotAcceptInternalValueVariableAsEnumInput(t *testing.T) {
	query := `query test($color: Int!) { colorEnum(fromEnum: $color) }`
	params := map[string]interface{}{
		"color": 2,
	}
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: `Variable "$color" of type "Int!" used in position expecting type "Color".`,
			},
		},
	}
	result := executeEnumTypeTestWithParams(t, query, params)
	if !testutil.EqualErrorMessage(expected, result, 0) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_EnumValueMayHaveAnInternalValueOfZero(t *testing.T) {
	query := `{
        colorEnum(fromEnum: RED)
        colorInt(fromEnum: RED)
      }`
	expected := &graphql.Result{
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
	expected := &graphql.Result{
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
