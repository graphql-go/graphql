package executor_test

import (
	"encoding/json"
	"github.com/chris-ramon/graphql/errors"
	"github.com/chris-ramon/graphql/executor"
	"github.com/chris-ramon/graphql/language/ast"
	"github.com/chris-ramon/graphql/language/location"
	"github.com/chris-ramon/graphql/testutil"
	"github.com/chris-ramon/graphql/types"
	"reflect"
	"testing"
)

var testComplexScalar *types.GraphQLScalarType = types.NewGraphQLScalarType(types.GraphQLScalarTypeConfig{
	Name: "ComplexScalar",
	Serialize: func(value interface{}) interface{} {
		if value == "DeserializedValue" {
			return "SerializedValue"
		}
		return nil
	},
	ParseValue: func(value interface{}) interface{} {
		if value == "SerializedValue" {
			return "DeserializedValue"
		}
		return nil
	},
	ParseLiteral: func(valueAST ast.Value) interface{} {
		astValue := valueAST.GetValue()
		if astValue, ok := astValue.(string); ok && astValue == "SerializedValue" {
			return "DeserializedValue"
		}
		return nil
	},
})

var testInputObject *types.GraphQLInputObjectType = types.NewGraphQLInputObjectType(types.InputObjectConfig{
	Name: "TestInputObject",
	Fields: types.InputObjectConfigFieldMap{
		"a": &types.InputObjectFieldConfig{
			Type: types.GraphQLString,
		},
		"b": &types.InputObjectFieldConfig{
			Type: types.NewGraphQLList(types.GraphQLString),
		},
		"c": &types.InputObjectFieldConfig{
			Type: types.NewGraphQLNonNull(types.GraphQLString),
		},
		"d": &types.InputObjectFieldConfig{
			Type: testComplexScalar,
		},
	},
})

func inputResolved(p types.GQLFRParams) interface{} {
	input, ok := p.Args["input"]
	if !ok {
		return nil
	}
	b, err := json.Marshal(input)
	if err != nil {
		return nil
	}
	return string(b)
}

var testType *types.GraphQLObjectType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
	Name: "TestType",
	Fields: types.GraphQLFieldConfigMap{
		"fieldWithObjectInput": &types.GraphQLFieldConfig{
			Type: types.GraphQLString,
			Args: types.GraphQLFieldConfigArgumentMap{
				"input": &types.GraphQLArgumentConfig{
					Type: testInputObject,
				},
			},
			Resolve: inputResolved,
		},
		"fieldWithNullableStringInput": &types.GraphQLFieldConfig{
			Type: types.GraphQLString,
			Args: types.GraphQLFieldConfigArgumentMap{
				"input": &types.GraphQLArgumentConfig{
					Type: types.GraphQLString,
				},
			},
			Resolve: inputResolved,
		},
		"fieldWithNonNullableStringInput": &types.GraphQLFieldConfig{
			Type: types.GraphQLString,
			Args: types.GraphQLFieldConfigArgumentMap{
				"input": &types.GraphQLArgumentConfig{
					Type: types.NewGraphQLNonNull(types.GraphQLString),
				},
			},
			Resolve: inputResolved,
		},
		"fieldWithDefaultArgumentValue": &types.GraphQLFieldConfig{
			Type: types.GraphQLString,
			Args: types.GraphQLFieldConfigArgumentMap{
				"input": &types.GraphQLArgumentConfig{
					Type:         types.GraphQLString,
					DefaultValue: "Hello World",
				},
			},
			Resolve: inputResolved,
		},
		"list": &types.GraphQLFieldConfig{
			Type: types.GraphQLString,
			Args: types.GraphQLFieldConfigArgumentMap{
				"input": &types.GraphQLArgumentConfig{
					Type: types.NewGraphQLList(types.GraphQLString),
				},
			},
			Resolve: inputResolved,
		},
		"nnList": &types.GraphQLFieldConfig{
			Type: types.GraphQLString,
			Args: types.GraphQLFieldConfigArgumentMap{
				"input": &types.GraphQLArgumentConfig{
					Type: types.NewGraphQLNonNull(types.NewGraphQLList(types.GraphQLString)),
				},
			},
			Resolve: inputResolved,
		},
		"listNN": &types.GraphQLFieldConfig{
			Type: types.GraphQLString,
			Args: types.GraphQLFieldConfigArgumentMap{
				"input": &types.GraphQLArgumentConfig{
					Type: types.NewGraphQLList(types.NewGraphQLNonNull(types.GraphQLString)),
				},
			},
			Resolve: inputResolved,
		},
		"nnListNN": &types.GraphQLFieldConfig{
			Type: types.GraphQLString,
			Args: types.GraphQLFieldConfigArgumentMap{
				"input": &types.GraphQLArgumentConfig{
					Type: types.NewGraphQLNonNull(types.NewGraphQLList(types.NewGraphQLNonNull(types.GraphQLString))),
				},
			},
			Resolve: inputResolved,
		},
	},
})

var variablesTestSchema, _ = types.NewGraphQLSchema(types.GraphQLSchemaConfig{
	Query: testType,
})

func TestVariables_ObjectsAndNullability_UsingInlineStructs_ExecutesWithComplexInput(t *testing.T) {
	doc := `
        {
          fieldWithObjectInput(input: {a: "foo", b: ["bar"], c: "baz"})
        }
	`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"fieldWithObjectInput": `{"a":"foo","b":["bar"],"c":"baz"}`,
		},
	}
	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_ObjectsAndNullability_UsingInlineStructs_ProperlyParsesSingleValueToList(t *testing.T) {
	doc := `
        {
          fieldWithObjectInput(input: {a: "foo", b: "bar", c: "baz"})
        }
	`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"fieldWithObjectInput": `{"a":"foo","b":["bar"],"c":"baz"}`,
		},
	}
	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_ObjectsAndNullability_UsingInlineStructs_DoesNotUseIncorrectValue(t *testing.T) {
	doc := `
        {
          fieldWithObjectInput(input: ["foo", "bar", "baz"])
        }
	`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"fieldWithObjectInput": nil,
		},
	}
	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func testVariables_ObjectsAndNullability_UsingVariables_GetAST(t *testing.T) *ast.Document {
	doc := `
        query q($input: TestInputObject) {
          fieldWithObjectInput(input: $input)
        }
	`
	return testutil.Parse(t, doc)
}
func TestVariables_ObjectsAndNullability_UsingVariables_ExecutesWithComplexInput(t *testing.T) {

	params := map[string]interface{}{
		"input": map[string]interface{}{
			"a": "foo",
			"b": []interface{}{"bar"},
			"c": "baz",
		},
	}
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"fieldWithObjectInput": `{"a":"foo","b":["bar"],"c":"baz"}`,
		},
	}

	ast := testVariables_ObjectsAndNullability_UsingVariables_GetAST(t)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestVariables_ObjectsAndNullability_UsingVariables_UsesDefaultValueWhenNotProvided(t *testing.T) {

	doc := `
	  query q($input: TestInputObject = {a: "foo", b: ["bar"], c: "baz"}) {
		fieldWithObjectInput(input: $input)
	  }
	`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"fieldWithObjectInput": `{"a":"foo","b":["bar"],"c":"baz"}`,
		},
	}

	withDefaultsAST := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    withDefaultsAST,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_ObjectsAndNullability_UsingVariables_ProperlyParsesSingleValueToList(t *testing.T) {
	params := map[string]interface{}{
		"input": map[string]interface{}{
			"a": "foo",
			"b": "bar",
			"c": "baz",
		},
	}
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"fieldWithObjectInput": `{"a":"foo","b":["bar"],"c":"baz"}`,
		},
	}

	ast := testVariables_ObjectsAndNullability_UsingVariables_GetAST(t)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_ObjectsAndNullability_UsingVariables_ExecutesWithComplexScalarInput(t *testing.T) {
	params := map[string]interface{}{
		"input": map[string]interface{}{
			"c": "foo",
			"d": "SerializedValue",
		},
	}
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"fieldWithObjectInput": `{"c":"foo","d":"DeserializedValue"}`,
		},
	}

	ast := testVariables_ObjectsAndNullability_UsingVariables_GetAST(t)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_ObjectsAndNullability_UsingVariables_ErrorsOnNullForNestedNonNull(t *testing.T) {
	params := map[string]interface{}{
		"input": map[string]interface{}{
			"a": "foo",
			"b": "bar",
			"c": nil,
		},
	}
	expected := &types.GraphQLResult{
		Data: nil,
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: `Variable "$input" expected value of type "TestInputObject" but ` +
					`got: {"a":"foo","b":"bar","c":null}.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line: 2, Column: 17,
					},
				},
			},
		},
	}

	ast := testVariables_ObjectsAndNullability_UsingVariables_GetAST(t)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_ObjectsAndNullability_UsingVariables_ErrorsOnIncorrectType(t *testing.T) {
	params := map[string]interface{}{
		"input": "foo bar",
	}
	expected := &types.GraphQLResult{
		Data: nil,
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: `Variable "$input" expected value of type "TestInputObject" but ` +
					`got: "foo bar".`,
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line: 2, Column: 17,
					},
				},
			},
		},
	}

	ast := testVariables_ObjectsAndNullability_UsingVariables_GetAST(t)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_ObjectsAndNullability_UsingVariables_ErrorsOnOmissionOfNestedNonNull(t *testing.T) {
	params := map[string]interface{}{
		"input": map[string]interface{}{
			"a": "foo",
			"b": "bar",
		},
	}
	expected := &types.GraphQLResult{
		Data: nil,
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: `Variable "$input" expected value of type "TestInputObject" but ` +
					`got: {"a":"foo","b":"bar"}.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line: 2, Column: 17,
					},
				},
			},
		},
	}

	ast := testVariables_ObjectsAndNullability_UsingVariables_GetAST(t)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_ObjectsAndNullability_UsingVariables_ErrorsOnAdditionOfUnknownInputField(t *testing.T) {
	params := map[string]interface{}{
		"input": map[string]interface{}{
			"a": "foo",
			"b": "bar",
			"c": "baz",
			"d": "dog",
		},
	}
	expected := &types.GraphQLResult{
		Data: nil,
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: `Variable "$input" expected value of type "TestInputObject" but ` +
					`got: {"a":"foo","b":"bar","c":"baz","d":"dog"}.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line: 2, Column: 17,
					},
				},
			},
		},
	}

	ast := testVariables_ObjectsAndNullability_UsingVariables_GetAST(t)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestVariables_NullableScalars_AllowsNullableInputsToBeOmitted(t *testing.T) {
	doc := `
      {
        fieldWithNullableStringInput
      }
	`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"fieldWithNullableStringInput": nil,
		},
	}

	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_NullableScalars_AllowsNullableInputsToBeOmittedInAVariable(t *testing.T) {
	doc := `
      query SetsNullable($value: String) {
        fieldWithNullableStringInput(input: $value)
      }
	`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"fieldWithNullableStringInput": nil,
		},
	}

	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_NullableScalars_AllowsNullableInputsToBeOmittedInAnUnlistedVariable(t *testing.T) {
	doc := `
      query SetsNullable {
        fieldWithNullableStringInput(input: $value)
      }
	`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"fieldWithNullableStringInput": nil,
		},
	}

	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_NullableScalars_AllowsNullableInputsToBeSetToNullInAVariable(t *testing.T) {
	doc := `
      query SetsNullable($value: String) {
        fieldWithNullableStringInput(input: $value)
      }
	`
	params := map[string]interface{}{
		"value": nil,
	}
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"fieldWithNullableStringInput": nil,
		},
	}

	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_NullableScalars_AllowsNullableInputsToBeSetToAValueInAVariable(t *testing.T) {
	doc := `
      query SetsNullable($value: String) {
        fieldWithNullableStringInput(input: $value)
      }
	`
	params := map[string]interface{}{
		"value": "a",
	}
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"fieldWithNullableStringInput": `"a"`,
		},
	}

	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_NullableScalars_AllowsNullableInputsToBeSetToAValueDirectly(t *testing.T) {
	doc := `
      {
        fieldWithNullableStringInput(input: "a")
      }
	`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"fieldWithNullableStringInput": `"a"`,
		},
	}

	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestVariables_NonNullableScalars_DoesNotAllowNonNullableInputsToBeOmittedInAVariable(t *testing.T) {

	doc := `
        query SetsNonNullable($value: String!) {
          fieldWithNonNullableStringInput(input: $value)
        }
	`

	expected := &types.GraphQLResult{
		Data: nil,
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: `Variable "$value" of required type "String!" was not provided.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line: 2, Column: 31,
					},
				},
			},
		},
	}

	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_NonNullableScalars_DoesNotAllowNonNullableInputsToBeSetToNullInAVariable(t *testing.T) {
	doc := `
        query SetsNonNullable($value: String!) {
          fieldWithNonNullableStringInput(input: $value)
        }
	`

	params := map[string]interface{}{
		"value": nil,
	}
	expected := &types.GraphQLResult{
		Data: nil,
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: `Variable "$value" of required type "String!" was not provided.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line: 2, Column: 31,
					},
				},
			},
		},
	}

	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_NonNullableScalars_AllowsNonNullableInputsToBeSetToAValueInAVariable(t *testing.T) {
	doc := `
        query SetsNonNullable($value: String!) {
          fieldWithNonNullableStringInput(input: $value)
        }
	`

	params := map[string]interface{}{
		"value": "a",
	}
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"fieldWithNonNullableStringInput": `"a"`,
		},
	}

	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_NonNullableScalars_AllowsNonNullableInputsToBeSetToAValueDirectly(t *testing.T) {
	doc := `
      {
        fieldWithNonNullableStringInput(input: "a")
      }
	`

	params := map[string]interface{}{
		"value": "a",
	}

	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"fieldWithNonNullableStringInput": `"a"`,
		},
	}

	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_NonNullableScalars_PassesAlongNullForNonNullableInputsIfExplicitlySetInTheQuery(t *testing.T) {
	doc := `
      {
        fieldWithNonNullableStringInput
      }
	`

	params := map[string]interface{}{
		"value": "a",
	}

	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"fieldWithNonNullableStringInput": nil,
		},
	}

	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestVariables_ListsAndNullability_AllowsListsToBeNull(t *testing.T) {
	doc := `
        query q($input: [String]) {
          list(input: $input)
        }
	`
	params := map[string]interface{}{
		"input": nil,
	}

	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"list": nil,
		},
	}
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_ListsAndNullability_AllowsListsToContainValues(t *testing.T) {
	doc := `
        query q($input: [String]) {
          list(input: $input)
        }
	`
	params := map[string]interface{}{
		"input": []interface{}{"A"},
	}

	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"list": `["A"]`,
		},
	}
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_ListsAndNullability_AllowsListsToContainNull(t *testing.T) {
	doc := `
        query q($input: [String]) {
          list(input: $input)
        }
	`
	params := map[string]interface{}{
		"input": []interface{}{"A", nil, "B"},
	}

	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"list": `["A",null,"B"]`,
		},
	}
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_ListsAndNullability_DoesNotAllowNonNullListsToBeNull(t *testing.T) {
	doc := `
        query q($input: [String]!) {
          nnList(input: $input)
        }
	`
	expected := &types.GraphQLResult{
		Data: nil,
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: `Variable "$input" of required type "[String]!" was not provided.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line: 2, Column: 17,
					},
				},
			},
		},
	}
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_ListsAndNullability_AllowsNonNullListsToContainValues(t *testing.T) {
	doc := `
        query q($input: [String]!) {
          nnList(input: $input)
        }
	`
	params := map[string]interface{}{
		"input": []interface{}{"A"},
	}
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"nnList": `["A"]`,
		},
	}
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_ListsAndNullability_AllowsNonNullListsToContainNull(t *testing.T) {
	doc := `
        query q($input: [String]!) {
          nnList(input: $input)
        }
	`
	params := map[string]interface{}{
		"input": []interface{}{"A", nil, "B"},
	}
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"nnList": `["A",null,"B"]`,
		},
	}
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_ListsAndNullability_AllowsListsOfNonNullsToBeNull(t *testing.T) {
	doc := `
        query q($input: [String!]) {
          listNN(input: $input)
        }
	`
	params := map[string]interface{}{
		"input": nil,
	}
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"listNN": nil,
		},
	}
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_ListsAndNullability_AllowsListsOfNonNullsToContainValues(t *testing.T) {
	doc := `
        query q($input: [String!]) {
          listNN(input: $input)
        }
	`
	params := map[string]interface{}{
		"input": []interface{}{"A"},
	}
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"listNN": `["A"]`,
		},
	}
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_ListsAndNullability_DoesNotAllowListOfNonNullsToContainNull(t *testing.T) {
	doc := `
        query q($input: [String!]) {
          listNN(input: $input)
        }
	`
	params := map[string]interface{}{
		"input": []interface{}{"A", nil, "B"},
	}
	expected := &types.GraphQLResult{
		Data: nil,
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: `Variable "$input" expected value of type "[String!]" but got: ` +
					`["A",null,"B"].`,
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line: 2, Column: 17,
					},
				},
			},
		},
	}
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_ListsAndNullability_DoesNotAllowNonNullListOfNonNullsToBeNull(t *testing.T) {
	doc := `
        query q($input: [String!]!) {
          nnListNN(input: $input)
        }
	`
	params := map[string]interface{}{
		"input": nil,
	}
	expected := &types.GraphQLResult{
		Data: nil,
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: `Variable "$input" of required type "[String!]!" was not provided.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line: 2, Column: 17,
					},
				},
			},
		},
	}
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_ListsAndNullability_AllowsNonNullListsOfNonNulsToContainValues(t *testing.T) {
	doc := `
        query q($input: [String!]!) {
          nnListNN(input: $input)
        }
	`
	params := map[string]interface{}{
		"input": []interface{}{"A"},
	}
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"nnListNN": `["A"]`,
		},
	}
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_ListsAndNullability_DoesNotAllowNonNullListOfNonNullsToContainNull(t *testing.T) {
	doc := `
        query q($input: [String!]!) {
          nnListNN(input: $input)
        }
	`
	params := map[string]interface{}{
		"input": []interface{}{"A", nil, "B"},
	}
	expected := &types.GraphQLResult{
		Data: nil,
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: `Variable "$input" expected value of type "[String!]!" but got: ` +
					`["A",null,"B"].`,
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line: 2, Column: 17,
					},
				},
			},
		},
	}
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_ListsAndNullability_DoesNotAllowInvalidTypesToBeUsedAsValues(t *testing.T) {
	doc := `
        query q($input: TestType!) {
          fieldWithObjectInput(input: $input)
        }
	`
	params := map[string]interface{}{
		"input": map[string]interface{}{
			"list": []interface{}{"A", "B"},
		},
	}
	expected := &types.GraphQLResult{
		Data: nil,
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: `Variable "$input" expected value of type "TestType!" which cannot be used as an input type.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line: 2, Column: 17,
					},
				},
			},
		},
	}
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_ListsAndNullability_DoesNotAllowUnknownTypesToBeUsedAsValues(t *testing.T) {
	doc := `
        query q($input: UnknownType!) {
          fieldWithObjectInput(input: $input)
        }
	`
	params := map[string]interface{}{
		"input": "whoknows",
	}
	expected := &types.GraphQLResult{
		Data: nil,
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: `Variable "$input" expected value of type "UnknownType!" which cannot be used as an input type.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line: 2, Column: 17,
					},
				},
			},
		},
	}
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestVariables_UsesArgumentDefaultValues_WhenNoArgumentProvided(t *testing.T) {
	doc := `
	{
      fieldWithDefaultArgumentValue
    }
	`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"fieldWithDefaultArgumentValue": `"Hello World"`,
		},
	}
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_UsesArgumentDefaultValues_WhenNullableVariableProvided(t *testing.T) {
	doc := `
	query optionalVariable($optional: String) {
        fieldWithDefaultArgumentValue(input: $optional)
    }
	`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"fieldWithDefaultArgumentValue": `"Hello World"`,
		},
	}
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestVariables_UsesArgumentDefaultValues_WhenArgumentProvidedCannotBeParsed(t *testing.T) {
	doc := `
	{
		fieldWithDefaultArgumentValue(input: WRONG_TYPE)
	}
	`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"fieldWithDefaultArgumentValue": `"Hello World"`,
		},
	}
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
