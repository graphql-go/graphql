package graphql_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/graphql-go/graphql/language/visitor"
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

func TestQueryWithCustomRule(t *testing.T) {
	// Test graphql.Do() with custom rule, it extracts query name from each
	// Tests.
	ruleN := len(graphql.SpecifiedRules)
	rules := make([]graphql.ValidationRuleFn, ruleN+1)
	copy(rules[:ruleN], graphql.SpecifiedRules)

	var (
		queryFound bool
		queryName  string
	)
	rules[ruleN] = func(context *graphql.ValidationContext) *graphql.ValidationRuleInstance {
		return &graphql.ValidationRuleInstance{
			VisitorOpts: &visitor.VisitorOptions{
				KindFuncMap: map[string]visitor.NamedVisitFuncs{
					kinds.OperationDefinition: {
						Enter: func(p visitor.VisitFuncParams) (string, interface{}) {
							od, ok := p.Node.(*ast.OperationDefinition)
							if ok && od.Operation == "query" {
								queryFound = true
								if od.Name != nil {
									queryName = od.Name.Value
								}
							}
							return visitor.ActionNoChange, nil
						},
					},
				},
			},
		}
	}

	expectedNames := []string{
		"HeroNameQuery",
		"HeroNameAndFriendsQuery",
		"HumanByIdQuery",
	}

	for i, test := range Tests {
		queryFound, queryName = false, ""
		params := graphql.Params{
			Schema:          test.Schema,
			RequestString:   test.Query,
			VariableValues:  test.Variables,
			ValidationRules: rules,
		}
		testGraphql(test, params, t)
		if !queryFound {
			t.Fatal("can't detect \"query\" operation by validation rule")
		}
		if queryName != expectedNames[i] {
			t.Fatalf("unexpected query name: want=%s got=%s", queryName, expectedNames)
		}
	}
}

// TestCustomRuleWithArgs tests graphql.GetArgumentValues() be able to access
// field's argument values from custom validation rule.
func TestCustomRuleWithArgs(t *testing.T) {
	fieldDef, ok := testutil.StarWarsSchema.QueryType().Fields()["human"]
	if !ok {
		t.Fatal("can't retrieve \"human\" field definition")
	}

	// a custom validation rule to extract argument values of "human" field.
	var actual map[string]interface{}
	enter := func(p visitor.VisitFuncParams) (string, interface{}) {
		// only interested in "human" field.
		fieldNode, ok := p.Node.(*ast.Field)
		if !ok || fieldNode.Name == nil || fieldNode.Name.Value != "human" {
			return visitor.ActionNoChange, nil
		}
		// extract argument values by graphql.GetArgumentValues().
		actual = graphql.GetArgumentValues(fieldDef.Args, fieldNode.Arguments, nil)
		return visitor.ActionNoChange, nil
	}
	checkHumanArgs := func(context *graphql.ValidationContext) *graphql.ValidationRuleInstance {
		return &graphql.ValidationRuleInstance{
			VisitorOpts: &visitor.VisitorOptions{
				KindFuncMap: map[string]visitor.NamedVisitFuncs{
					kinds.Field: {Enter: enter},
				},
			},
		}
	}

	for _, tc := range []struct {
		query      string
		expected   map[string]interface{}
	}{
		{
			`query { human(id: "1000") { name } }`,
			map[string]interface{}{"id": "1000"},
		},
		{
			`query { human(id: "1002") { name } }`,
			map[string]interface{}{"id": "1002"},
		},
		{
			`query { human(id: "9999") { name } }`,
			map[string]interface{}{"id": "9999"},
		},
	} {
		actual = nil
		params := graphql.Params{
			Schema:          testutil.StarWarsSchema,
			RequestString:   tc.query,
			ValidationRules: append(graphql.SpecifiedRules, checkHumanArgs),
		}
		result := graphql.Do(params)
		if len(result.Errors) > 0 {
			t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
		}
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Fatalf("unexpected result: want=%+v got=%+v", tc.expected, actual)
		}
	}
}
