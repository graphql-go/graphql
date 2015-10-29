package graphql

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/chris-ramon/graphql/gqlerrors"
	"github.com/chris-ramon/graphql/language/ast"
	"github.com/chris-ramon/graphql/language/location"
)

var testComplexScalar *Scalar = NewScalar(ScalarConfig{
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

var testInputObject *InputObject = NewInputObject(InputObjectConfig{
	Name: "TestInputObject",
	Fields: InputObjectConfigFieldMap{
		"a": &InputObjectFieldConfig{
			Type: String,
		},
		"b": &InputObjectFieldConfig{
			Type: NewList(String),
		},
		"c": &InputObjectFieldConfig{
			Type: NewNonNull(String),
		},
		"d": &InputObjectFieldConfig{
			Type: testComplexScalar,
		},
	},
})

func inputResolved(p GQLFRParams) interface{} {
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

var testType *Object = NewObject(ObjectConfig{
	Name: "TestType",
	Fields: FieldConfigMap{
		"fieldWithObjectInput": &FieldConfig{
			Type: String,
			Args: FieldConfigArgument{
				"input": &ArgumentConfig{
					Type: testInputObject,
				},
			},
			Resolve: inputResolved,
		},
		"fieldWithNullableStringInput": &FieldConfig{
			Type: String,
			Args: FieldConfigArgument{
				"input": &ArgumentConfig{
					Type: String,
				},
			},
			Resolve: inputResolved,
		},
		"fieldWithNonNullableStringInput": &FieldConfig{
			Type: String,
			Args: FieldConfigArgument{
				"input": &ArgumentConfig{
					Type: NewNonNull(String),
				},
			},
			Resolve: inputResolved,
		},
		"fieldWithDefaultArgumentValue": &FieldConfig{
			Type: String,
			Args: FieldConfigArgument{
				"input": &ArgumentConfig{
					Type:         String,
					DefaultValue: "Hello World",
				},
			},
			Resolve: inputResolved,
		},
		"list": &FieldConfig{
			Type: String,
			Args: FieldConfigArgument{
				"input": &ArgumentConfig{
					Type: NewList(String),
				},
			},
			Resolve: inputResolved,
		},
		"nnList": &FieldConfig{
			Type: String,
			Args: FieldConfigArgument{
				"input": &ArgumentConfig{
					Type: NewNonNull(NewList(String)),
				},
			},
			Resolve: inputResolved,
		},
		"listNN": &FieldConfig{
			Type: String,
			Args: FieldConfigArgument{
				"input": &ArgumentConfig{
					Type: NewList(NewNonNull(String)),
				},
			},
			Resolve: inputResolved,
		},
		"nnListNN": &FieldConfig{
			Type: String,
			Args: FieldConfigArgument{
				"input": &ArgumentConfig{
					Type: NewNonNull(NewList(NewNonNull(String))),
				},
			},
			Resolve: inputResolved,
		},
	},
})

var variablesTestSchema, _ = NewSchema(SchemaConfig{
	Query: testType,
})

func TestVariables_ObjectsAndNullability_UsingInlineStructs_ExecutesWithComplexInput(t *testing.T) {
	doc := `
        {
          fieldWithObjectInput(input: {a: "foo", b: ["bar"], c: "baz"})
        }
	`
	expected := &Result{
		Data: map[string]interface{}{
			"fieldWithObjectInput": `{"a":"foo","b":["bar"],"c":"baz"}`,
		},
	}
	// parse query
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestVariables_ObjectsAndNullability_UsingInlineStructs_ProperlyParsesSingleValueToList(t *testing.T) {
	doc := `
        {
          fieldWithObjectInput(input: {a: "foo", b: "bar", c: "baz"})
        }
	`
	expected := &Result{
		Data: map[string]interface{}{
			"fieldWithObjectInput": `{"a":"foo","b":["bar"],"c":"baz"}`,
		},
	}
	// parse query
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestVariables_ObjectsAndNullability_UsingInlineStructs_DoesNotUseIncorrectValue(t *testing.T) {
	doc := `
        {
          fieldWithObjectInput(input: ["foo", "bar", "baz"])
        }
	`
	expected := &Result{
		Data: map[string]interface{}{
			"fieldWithObjectInput": nil,
		},
	}
	// parse query
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}

func testVariables_ObjectsAndNullability_UsingVariables_GetAST(t *testing.T) *ast.Document {
	doc := `
        query q($input: TestInputObject) {
          fieldWithObjectInput(input: $input)
        }
	`
	return TestParse(t, doc)
}
func TestVariables_ObjectsAndNullability_UsingVariables_ExecutesWithComplexInput(t *testing.T) {

	params := map[string]interface{}{
		"input": map[string]interface{}{
			"a": "foo",
			"b": []interface{}{"bar"},
			"c": "baz",
		},
	}
	expected := &Result{
		Data: map[string]interface{}{
			"fieldWithObjectInput": `{"a":"foo","b":["bar"],"c":"baz"}`,
		},
	}

	ast := testVariables_ObjectsAndNullability_UsingVariables_GetAST(t)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}

func TestVariables_ObjectsAndNullability_UsingVariables_UsesDefaultValueWhenNotProvided(t *testing.T) {

	doc := `
	  query q($input: TestInputObject = {a: "foo", b: ["bar"], c: "baz"}) {
		fieldWithObjectInput(input: $input)
	  }
	`
	expected := &Result{
		Data: map[string]interface{}{
			"fieldWithObjectInput": `{"a":"foo","b":["bar"],"c":"baz"}`,
		},
	}

	withDefaultsAST := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    withDefaultsAST,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: map[string]interface{}{
			"fieldWithObjectInput": `{"a":"foo","b":["bar"],"c":"baz"}`,
		},
	}

	ast := testVariables_ObjectsAndNullability_UsingVariables_GetAST(t)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestVariables_ObjectsAndNullability_UsingVariables_ExecutesWithComplexScalarInput(t *testing.T) {
	params := map[string]interface{}{
		"input": map[string]interface{}{
			"c": "foo",
			"d": "SerializedValue",
		},
	}
	expected := &Result{
		Data: map[string]interface{}{
			"fieldWithObjectInput": `{"c":"foo","d":"DeserializedValue"}`,
		},
	}

	ast := testVariables_ObjectsAndNullability_UsingVariables_GetAST(t)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
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
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestVariables_ObjectsAndNullability_UsingVariables_ErrorsOnIncorrectType(t *testing.T) {
	params := map[string]interface{}{
		"input": "foo bar",
	}
	expected := &Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
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
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestVariables_ObjectsAndNullability_UsingVariables_ErrorsOnOmissionOfNestedNonNull(t *testing.T) {
	params := map[string]interface{}{
		"input": map[string]interface{}{
			"a": "foo",
			"b": "bar",
		},
	}
	expected := &Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
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
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
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
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}

func TestVariables_NullableScalars_AllowsNullableInputsToBeOmitted(t *testing.T) {
	doc := `
      {
        fieldWithNullableStringInput
      }
	`
	expected := &Result{
		Data: map[string]interface{}{
			"fieldWithNullableStringInput": nil,
		},
	}

	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestVariables_NullableScalars_AllowsNullableInputsToBeOmittedInAVariable(t *testing.T) {
	doc := `
      query SetsNullable($value: String) {
        fieldWithNullableStringInput(input: $value)
      }
	`
	expected := &Result{
		Data: map[string]interface{}{
			"fieldWithNullableStringInput": nil,
		},
	}

	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestVariables_NullableScalars_AllowsNullableInputsToBeOmittedInAnUnlistedVariable(t *testing.T) {
	doc := `
      query SetsNullable {
        fieldWithNullableStringInput(input: $value)
      }
	`
	expected := &Result{
		Data: map[string]interface{}{
			"fieldWithNullableStringInput": nil,
		},
	}

	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: map[string]interface{}{
			"fieldWithNullableStringInput": nil,
		},
	}

	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: map[string]interface{}{
			"fieldWithNullableStringInput": `"a"`,
		},
	}

	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestVariables_NullableScalars_AllowsNullableInputsToBeSetToAValueDirectly(t *testing.T) {
	doc := `
      {
        fieldWithNullableStringInput(input: "a")
      }
	`
	expected := &Result{
		Data: map[string]interface{}{
			"fieldWithNullableStringInput": `"a"`,
		},
	}

	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}

func TestVariables_NonNullableScalars_DoesNotAllowNonNullableInputsToBeOmittedInAVariable(t *testing.T) {

	doc := `
        query SetsNonNullable($value: String!) {
          fieldWithNonNullableStringInput(input: $value)
        }
	`

	expected := &Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: `Variable "$value" of required type "String!" was not provided.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line: 2, Column: 31,
					},
				},
			},
		},
	}

	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: `Variable "$value" of required type "String!" was not provided.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line: 2, Column: 31,
					},
				},
			},
		},
	}

	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: map[string]interface{}{
			"fieldWithNonNullableStringInput": `"a"`,
		},
	}

	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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

	expected := &Result{
		Data: map[string]interface{}{
			"fieldWithNonNullableStringInput": `"a"`,
		},
	}

	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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

	expected := &Result{
		Data: map[string]interface{}{
			"fieldWithNonNullableStringInput": nil,
		},
	}

	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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

	expected := &Result{
		Data: map[string]interface{}{
			"list": nil,
		},
	}
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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

	expected := &Result{
		Data: map[string]interface{}{
			"list": `["A"]`,
		},
	}
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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

	expected := &Result{
		Data: map[string]interface{}{
			"list": `["A",null,"B"]`,
		},
	}
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestVariables_ListsAndNullability_DoesNotAllowNonNullListsToBeNull(t *testing.T) {
	doc := `
        query q($input: [String]!) {
          nnList(input: $input)
        }
	`
	expected := &Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: `Variable "$input" of required type "[String]!" was not provided.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line: 2, Column: 17,
					},
				},
			},
		},
	}
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: map[string]interface{}{
			"nnList": `["A"]`,
		},
	}
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: map[string]interface{}{
			"nnList": `["A",null,"B"]`,
		},
	}
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: map[string]interface{}{
			"listNN": nil,
		},
	}
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: map[string]interface{}{
			"listNN": `["A"]`,
		},
	}
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
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
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: `Variable "$input" of required type "[String!]!" was not provided.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line: 2, Column: 17,
					},
				},
			},
		},
	}
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: map[string]interface{}{
			"nnListNN": `["A"]`,
		},
	}
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
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
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: `Variable "$input" expected value of type "TestType!" which cannot be used as an input type.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line: 2, Column: 17,
					},
				},
			},
		},
	}
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: `Variable "$input" expected value of type "UnknownType!" which cannot be used as an input type.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line: 2, Column: 17,
					},
				},
			},
		},
	}
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
		Args:   params,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}

func TestVariables_UsesArgumentDefaultValues_WhenNoArgumentProvided(t *testing.T) {
	doc := `
	{
      fieldWithDefaultArgumentValue
    }
	`
	expected := &Result{
		Data: map[string]interface{}{
			"fieldWithDefaultArgumentValue": `"Hello World"`,
		},
	}
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestVariables_UsesArgumentDefaultValues_WhenNullableVariableProvided(t *testing.T) {
	doc := `
	query optionalVariable($optional: String) {
        fieldWithDefaultArgumentValue(input: $optional)
    }
	`
	expected := &Result{
		Data: map[string]interface{}{
			"fieldWithDefaultArgumentValue": `"Hello World"`,
		},
	}
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestVariables_UsesArgumentDefaultValues_WhenArgumentProvidedCannotBeParsed(t *testing.T) {
	doc := `
	{
		fieldWithDefaultArgumentValue(input: WRONG_TYPE)
	}
	`
	expected := &Result{
		Data: map[string]interface{}{
			"fieldWithDefaultArgumentValue": `"Hello World"`,
		},
	}
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: variablesTestSchema,
		AST:    ast,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
