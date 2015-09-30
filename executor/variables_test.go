package executor_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/chris-ramon/graphql-go/errors"
	"github.com/chris-ramon/graphql-go/executor"
	"github.com/chris-ramon/graphql-go/gqltypes"
	"github.com/chris-ramon/graphql-go/language/ast"
	"github.com/chris-ramon/graphql-go/language/location"
	"github.com/chris-ramon/graphql-go/testutil"
)

var testComplexScalar *gqltypes.GraphQLScalarType = gqltypes.NewGraphQLScalarType(gqltypes.GraphQLScalarTypeConfig{
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

var testInputObject *gqltypes.GraphQLInputObjectType = gqltypes.NewGraphQLInputObjectType(gqltypes.InputObjectConfig{
	Name: "TestInputObject",
	Fields: gqltypes.InputObjectConfigFieldMap{
		"a": &gqltypes.InputObjectFieldConfig{
			Type: gqltypes.GraphQLString,
		},
		"b": &gqltypes.InputObjectFieldConfig{
			Type: gqltypes.NewGraphQLList(gqltypes.GraphQLString),
		},
		"c": &gqltypes.InputObjectFieldConfig{
			Type: gqltypes.NewGraphQLNonNull(gqltypes.GraphQLString),
		},
		"d": &gqltypes.InputObjectFieldConfig{
			Type: testComplexScalar,
		},
	},
})

func inputResolved(p gqltypes.GQLFRParams) interface{} {
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

var testType *gqltypes.GraphQLObjectType = gqltypes.NewGraphQLObjectType(gqltypes.GraphQLObjectTypeConfig{
	Name: "TestType",
	Fields: gqltypes.GraphQLFieldConfigMap{
		"fieldWithObjectInput": &gqltypes.GraphQLFieldConfig{
			Type: gqltypes.GraphQLString,
			Args: gqltypes.GraphQLFieldConfigArgumentMap{
				"input": &gqltypes.GraphQLArgumentConfig{
					Type: testInputObject,
				},
			},
			Resolve: inputResolved,
		},
		"fieldWithNullableStringInput": &gqltypes.GraphQLFieldConfig{
			Type: gqltypes.GraphQLString,
			Args: gqltypes.GraphQLFieldConfigArgumentMap{
				"input": &gqltypes.GraphQLArgumentConfig{
					Type: gqltypes.GraphQLString,
				},
			},
			Resolve: inputResolved,
		},
		"fieldWithNonNullableStringInput": &gqltypes.GraphQLFieldConfig{
			Type: gqltypes.GraphQLString,
			Args: gqltypes.GraphQLFieldConfigArgumentMap{
				"input": &gqltypes.GraphQLArgumentConfig{
					Type: gqltypes.NewGraphQLNonNull(gqltypes.GraphQLString),
				},
			},
			Resolve: inputResolved,
		},
		"fieldWithDefaultArgumentValue": &gqltypes.GraphQLFieldConfig{
			Type: gqltypes.GraphQLString,
			Args: gqltypes.GraphQLFieldConfigArgumentMap{
				"input": &gqltypes.GraphQLArgumentConfig{
					Type:         gqltypes.GraphQLString,
					DefaultValue: "Hello World",
				},
			},
			Resolve: inputResolved,
		},
		"list": &gqltypes.GraphQLFieldConfig{
			Type: gqltypes.GraphQLString,
			Args: gqltypes.GraphQLFieldConfigArgumentMap{
				"input": &gqltypes.GraphQLArgumentConfig{
					Type: gqltypes.NewGraphQLList(gqltypes.GraphQLString),
				},
			},
			Resolve: inputResolved,
		},
		"nnList": &gqltypes.GraphQLFieldConfig{
			Type: gqltypes.GraphQLString,
			Args: gqltypes.GraphQLFieldConfigArgumentMap{
				"input": &gqltypes.GraphQLArgumentConfig{
					Type: gqltypes.NewGraphQLNonNull(gqltypes.NewGraphQLList(gqltypes.GraphQLString)),
				},
			},
			Resolve: inputResolved,
		},
		"listNN": &gqltypes.GraphQLFieldConfig{
			Type: gqltypes.GraphQLString,
			Args: gqltypes.GraphQLFieldConfigArgumentMap{
				"input": &gqltypes.GraphQLArgumentConfig{
					Type: gqltypes.NewGraphQLList(gqltypes.NewGraphQLNonNull(gqltypes.GraphQLString)),
				},
			},
			Resolve: inputResolved,
		},
		"nnListNN": &gqltypes.GraphQLFieldConfig{
			Type: gqltypes.GraphQLString,
			Args: gqltypes.GraphQLFieldConfigArgumentMap{
				"input": &gqltypes.GraphQLArgumentConfig{
					Type: gqltypes.NewGraphQLNonNull(gqltypes.NewGraphQLList(gqltypes.NewGraphQLNonNull(gqltypes.GraphQLString))),
				},
			},
			Resolve: inputResolved,
		},
	},
})

var variablesTestSchema, _ = gqltypes.NewGraphQLSchema(gqltypes.GraphQLSchemaConfig{
	Query: testType,
})

func TestVariables_ObjectsAndNullability_UsingInlineStructs_ExecutesWithComplexInput(t *testing.T) {
	doc := `
        {
          fieldWithObjectInput(input: {a: "foo", b: ["bar"], c: "baz"})
        }
	`
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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

	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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

	expected := &gqltypes.GraphQLResult{
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

	expected := &gqltypes.GraphQLResult{
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

	expected := &gqltypes.GraphQLResult{
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

	expected := &gqltypes.GraphQLResult{
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

	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
	expected := &gqltypes.GraphQLResult{
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
