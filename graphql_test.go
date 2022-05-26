package graphql_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/testutil"
)

type T struct {
	Query     string
	Schema    graphql.Schema
	Expected  interface{}
	Variables map[string]interface{}
}

var Tests = []T{}

func init() {
	Tests = []T{
		{
			Query: `
				query HeroNameQuery {
					hero {
						name
					}
				}
			`,
			Schema: testutil.StarWarsSchema,
			Expected: &graphql.Result{
				Data: map[string]interface{}{
					"hero": map[string]interface{}{
						"name": "R2-D2",
					},
				},
			},
		},
		{
			Query: `
				query HeroNameAndFriendsQuery {
					hero {
						id
						name
						friends {
							name
						}
					}
				}
			`,
			Schema: testutil.StarWarsSchema,
			Expected: &graphql.Result{
				Data: map[string]interface{}{
					"hero": map[string]interface{}{
						"id":   "2001",
						"name": "R2-D2",
						"friends": []interface{}{
							map[string]interface{}{
								"name": "Luke Skywalker",
							},
							map[string]interface{}{
								"name": "Han Solo",
							},
							map[string]interface{}{
								"name": "Leia Organa",
							},
						},
					},
				},
			},
		},
		{
			Query: `
				query HumanByIdQuery($id: String!) {
					human(id: $id) {
						name
					}
				}
			`,
			Schema: testutil.StarWarsSchema,
			Expected: &graphql.Result{
				Data: map[string]interface{}{
					"human": map[string]interface{}{
						"name": "Darth Vader",
					},
				},
			},
			Variables: map[string]interface{}{
				"id": "1001",
			},
		},
	}
}

func TestQuery(t *testing.T) {
	for _, test := range Tests {
		params := graphql.Params{
			Schema:         test.Schema,
			RequestString:  test.Query,
			VariableValues: test.Variables,
		}
		testGraphql(test, params, t)
	}
}

func testGraphql(test T, p graphql.Params, t *testing.T) {
	result := graphql.Do(p)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(result, test.Expected) {
		t.Fatalf("wrong result, query: %v, graphql result diff: %v", test.Query, testutil.Diff(test.Expected, result))
	}
}

func TestBasicGraphQLExample(t *testing.T) {
	// taken from `graphql-js` README

	helloFieldResolved := func(p graphql.ResolveParams) (interface{}, error) {
		return "world", nil
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "RootQueryType",
			Fields: graphql.Fields{
				"hello": &graphql.Field{
					Description: "Returns `world`",
					Type:        graphql.String,
					Resolve:     helloFieldResolved,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("wrong result, unexpected errors: %v", err.Error())
	}
	query := "{ hello }"
	var expected interface{}
	expected = map[string]interface{}{
		"hello": "world",
	}

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(result.Data, expected) {
		t.Fatalf("wrong result, query: %v, graphql result diff: %v", query, testutil.Diff(expected, result))
	}

}

func TestThreadsContextFromParamsThrough(t *testing.T) {
	extractFieldFromContextFn := func(p graphql.ResolveParams) (interface{}, error) {
		return p.Context.Value(p.Args["key"]), nil
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"value": &graphql.Field{
					Type: graphql.String,
					Args: graphql.FieldConfigArgument{
						"key": &graphql.ArgumentConfig{Type: graphql.String},
					},
					Resolve: extractFieldFromContextFn,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("wrong result, unexpected errors: %v", err.Error())
	}
	query := `{ value(key:"a") }`

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       context.WithValue(context.TODO(), "a", "xyz"),
	})
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	expected := map[string]interface{}{"value": "xyz"}
	if !reflect.DeepEqual(result.Data, expected) {
		t.Fatalf("wrong result, query: %v, graphql result diff: %v", query, testutil.Diff(expected, result))
	}

}

func TestNewErrorChecksNilNodes(t *testing.T) {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"graphql_is": &graphql.Field{
					Type: graphql.String,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return "", nil
					},
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("unexpected errors: %v", err.Error())
	}
	query := `{graphql_is:great(sort:ByPopularity)}{stars}`
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) == 0 {
		t.Fatalf("expected errors, got: %v", result)
	}
}

func TestEmptyStringIsNotNull(t *testing.T) {
	checkForEmptyString := func(p graphql.ResolveParams) (interface{}, error) {
		arg := p.Args["arg"]
		if arg == nil || arg.(string) != "" {
			t.Errorf("Expected empty string for input arg, got %#v", arg)
		}
		return "yay", nil
	}
	returnEmptyString := func(p graphql.ResolveParams) (interface{}, error) {
		return "", nil
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"checkEmptyArg": &graphql.Field{
					Type: graphql.String,
					Args: graphql.FieldConfigArgument{
						"arg": &graphql.ArgumentConfig{Type: graphql.String},
					},
					Resolve: checkForEmptyString,
				},
				"checkEmptyResult": &graphql.Field{
					Type:    graphql.String,
					Resolve: returnEmptyString,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("wrong result, unexpected errors: %v", err.Error())
	}
	query := `{ checkEmptyArg(arg:"") checkEmptyResult }`

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	expected := map[string]interface{}{"checkEmptyArg": "yay", "checkEmptyResult": ""}
	if !reflect.DeepEqual(result.Data, expected) {
		t.Errorf("wrong result, query: %v, graphql result diff: %v", query, testutil.Diff(expected, result))
	}
}

func TestNullLiteralArguments(t *testing.T) {
	checkForNull := func(p graphql.ResolveParams) (interface{}, error) {
		arg, ok := p.Args["arg"]
		if !ok || arg != nil {
			t.Errorf("expected null for input arg, got %#v", arg)
		}
		return "yay", nil
	}
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"checkNullStringArg": &graphql.Field{
					Type: graphql.String,
					Args: graphql.FieldConfigArgument{
						"arg": &graphql.ArgumentConfig{Type: graphql.String},
					},
					Resolve: checkForNull,
				},
				"checkNullIntArg": &graphql.Field{
					Type: graphql.String,
					Args: graphql.FieldConfigArgument{
						"arg": &graphql.ArgumentConfig{Type: graphql.Int},
					},
					Resolve: checkForNull,
				},
				"checkNullBooleanArg": &graphql.Field{
					Type: graphql.String,
					Args: graphql.FieldConfigArgument{
						"arg": &graphql.ArgumentConfig{Type: graphql.Boolean},
					},
					Resolve: checkForNull,
				},
				"checkNullListArg": &graphql.Field{
					Type: graphql.String,
					Args: graphql.FieldConfigArgument{
						"arg": &graphql.ArgumentConfig{Type: graphql.NewList(graphql.String)},
					},
					Resolve: checkForNull,
				},
				"checkNullInputObjectArg": &graphql.Field{
					Type: graphql.String,
					Args: graphql.FieldConfigArgument{
						"arg": &graphql.ArgumentConfig{Type: graphql.NewInputObject(
							graphql.InputObjectConfig{
								Name: "InputType",
								Fields: graphql.InputObjectConfigFieldMap{
									"field1": {Type: graphql.String},
									"field2": {Type: graphql.Int},
								},
							})},
					},
					Resolve: checkForNull,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("wrong result, unexpected errors: %v", err.Error())
	}
	query := `{ checkNullStringArg(arg:null) checkNullIntArg(arg:null) checkNullBooleanArg(arg:null) checkNullListArg(arg:null) checkNullInputObjectArg(arg:null) }`

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	expected := map[string]interface{}{
		"checkNullStringArg": "yay", "checkNullIntArg": "yay",
		"checkNullBooleanArg": "yay", "checkNullListArg": "yay",
		"checkNullInputObjectArg": "yay"}
	if !reflect.DeepEqual(result.Data, expected) {
		t.Errorf("wrong result, query: %v, graphql result diff: %v", query, testutil.Diff(expected, result))
	}
}

func TestNullLiteralDefaultVariableValue(t *testing.T) {
	checkForNull := func(p graphql.ResolveParams) (interface{}, error) {
		arg, ok := p.Args["arg"]
		if !ok || arg != nil {
			t.Errorf("expected null for input arg, got %#v", arg)
		}
		return "yay", nil
	}
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"checkNullStringArg": &graphql.Field{
					Type: graphql.String,
					Args: graphql.FieldConfigArgument{
						"arg": &graphql.ArgumentConfig{Type: graphql.String},
					},
					Resolve: checkForNull,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("wrong result, unexpected errors: %v", err.Error())
	}
	query := `query Test($value: String = null) { checkNullStringArg(arg: $value) }`

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		VariableValues: map[string]interface{}{"value2": nil},
	})
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	expected := map[string]interface{}{ "checkNullStringArg": "yay", }
	if !reflect.DeepEqual(result.Data, expected) {
		t.Errorf("wrong result, query: %v, graphql result diff: %v", query, testutil.Diff(expected, result))
	}
}

func TestNullLiteralVariables(t *testing.T) {
	checkForNull := func(p graphql.ResolveParams) (interface{}, error) {
		arg, ok := p.Args["arg"]
		if !ok || arg != nil {
			t.Errorf("expected null for input arg, got %#v", arg)
		}
		return "yay", nil
	}
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"checkNullStringArg": &graphql.Field{
					Type: graphql.String,
					Args: graphql.FieldConfigArgument{
						"arg": &graphql.ArgumentConfig{Type: graphql.String},
					},
					Resolve: checkForNull,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("wrong result, unexpected errors: %v", err.Error())
	}
	query := `query Test($value: String) { checkNullStringArg(arg: $value) }`

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		VariableValues: map[string]interface{}{"value": nil},
	})
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	expected := map[string]interface{}{ "checkNullStringArg": "yay", }
	if !reflect.DeepEqual(result.Data, expected) {
		t.Errorf("wrong result, query: %v, graphql result diff: %v", query, testutil.Diff(expected, result))
	}
}

func TestErrorNullLiteralForNotNullArgument(t *testing.T) {
	checkNotCalled := func(p graphql.ResolveParams) (interface{}, error) {
		t.Error("shouldn't have been called")
		return nil, nil
	}
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"checkNotNullArg": &graphql.Field{
					Type: graphql.String,
					Args: graphql.FieldConfigArgument{
						"arg": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String) },
					},
					Resolve: checkNotCalled,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("wrong result, unexpected errors: %v", err.Error())
	}
	query := `{ checkNotNullArg(arg:null) }`

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})

	if len(result.Errors) == 0 {
		t.Fatalf("expected errors, got: %v", result)
	}

	expectedMessage := `Argument "arg" has invalid value <nil>.
Expected "String!", found null.`;

	if result.Errors[0].Message != expectedMessage {
		t.Fatalf("unexpected error.\nexpected:\n%s\ngot:\n%s\n", expectedMessage, result.Errors[0].Message)
	}
}

func TestNullInputObjectFields(t *testing.T) {
	checkForNull := func(p graphql.ResolveParams) (interface{}, error) {
		arg := p.Args["arg"]
		expectedValue := map[string]interface{}{ "field1": nil, "field2": nil, "field3": nil, "field4" : "abc", "field5": 42, "field6": true}
		if value, ok := arg.(map[string]interface{}); !ok  {
			t.Errorf("expected map[string]interface{} for input arg, got %#v", arg)
		} else if !reflect.DeepEqual(expectedValue, value) {
			t.Errorf("unexpected input object, diff: %v", testutil.Diff(expectedValue, value))
		}
		return "yay", nil
	}
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"checkNullInputObjectFields": &graphql.Field{
					Type: graphql.String,
					Args: graphql.FieldConfigArgument{
						"arg": &graphql.ArgumentConfig{Type: graphql.NewInputObject(
							graphql.InputObjectConfig{
								Name: "InputType",
								Fields: graphql.InputObjectConfigFieldMap{
									"field1": {Type: graphql.String},
									"field2": {Type: graphql.Int},
									"field3": {Type: graphql.Boolean},
									"field4": {Type: graphql.String},
									"field5": {Type: graphql.Int},
									"field6": {Type: graphql.Boolean},
								},
							})},
					},
					Resolve: checkForNull,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("wrong result, unexpected errors: %v", err.Error())
	}
	query := `{ checkNullInputObjectFields(arg: {field1: null, field2: null, field3: null, field4: "abc", field5: 42, field6: true }) }`

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	expected := map[string]interface{}{ "checkNullInputObjectFields": "yay" }
	if !reflect.DeepEqual(result.Data, expected) {
		t.Errorf("wrong result, query: %v, graphql result diff: %v", query, testutil.Diff(expected, result))
	}
}

func TestErrorNullInList(t *testing.T) {
	checkNotCalled := func(p graphql.ResolveParams) (interface{}, error) {
		t.Error("shouldn't have been called")
		return nil, nil
	}
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"checkNotNullInListArg": &graphql.Field{
					Type: graphql.String,
					Args: graphql.FieldConfigArgument{
						"arg": &graphql.ArgumentConfig{Type: graphql.NewList(graphql.String) },
					},
					Resolve: checkNotCalled,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("wrong result, unexpected errors: %v", err.Error())
	}
	query := `{ checkNotNullInListArg(arg: [null, null]) }`

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})

	if len(result.Errors) == 0 {
		t.Fatalf("expected errors, got: %v", result)
	}

	expectedMessage := `Argument "arg" has invalid value [<nil>, <nil>].
In element #1: Unexpected null literal.
In element #2: Unexpected null literal.`

	if result.Errors[0].Message != expectedMessage {
		t.Fatalf("unexpected error.\nexpected:\n%s\ngot:\n%s\n", expectedMessage, result.Errors[0].Message)
	}
}
