package executor_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/chris-ramon/graphql-go/errors"
	"github.com/chris-ramon/graphql-go/executor"
	"github.com/chris-ramon/graphql-go/language/location"
	"github.com/chris-ramon/graphql-go/testutil"
	"github.com/chris-ramon/graphql-go/types"
)

func TestExecutesArbitraryCode(t *testing.T) {

	deepData := map[string]interface{}{}
	data := map[string]interface{}{
		"a": func() interface{} { return "Apple" },
		"b": func() interface{} { return "Banana" },
		"c": func() interface{} { return "Cookie" },
		"d": func() interface{} { return "Donut" },
		"e": func() interface{} { return "Egg" },
		"f": "Fish",
		"pic": func(size int) string {
			return fmt.Sprintf("Pic of size: %v", size)
		},
		"deep": func() interface{} { return deepData },
	}
	data["promise"] = func() interface{} {
		return data
	}
	deepData = map[string]interface{}{
		"a":      func() interface{} { return "Already Been Done" },
		"b":      func() interface{} { return "Boring" },
		"c":      func() interface{} { return []string{"Contrived", "", "Confusing"} },
		"deeper": func() interface{} { return []interface{}{data, nil, data} },
	}

	query := `
      query Example($size: Int) {
        a,
        b,
        x: c
        ...c
        f
        ...on DataType {
          pic(size: $size)
          promise {
            a
          }
        }
        deep {
          a
          b
          c
          deeper {
            a
            b
          }
        }
      }

      fragment c on DataType {
        d
        e
      }
    `

	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"b": "Banana",
			"x": "Cookie",
			"d": "Donut",
			"e": "Egg",
			"promise": map[string]interface{}{
				"a": "Apple",
			},
			"a": "Apple",
			"deep": map[string]interface{}{
				"a": "Already Been Done",
				"b": "Boring",
				"c": []interface{}{
					"Contrived",
					nil,
					"Confusing",
				},
				"deeper": []interface{}{
					map[string]interface{}{
						"a": "Apple",
						"b": "Banana",
					},
					nil,
					map[string]interface{}{
						"a": "Apple",
						"b": "Banana",
					},
				},
			},
			"f":   "Fish",
			"pic": "Pic of size: 100",
		},
	}

	// Schema Definitions
	picResolverFn := func(p types.GQLFRParams) interface{} {
		// get and type assert ResolveFn for this field
		picResolver, ok := p.Source.(map[string]interface{})["pic"].(func(size int) string)
		if !ok {
			return nil
		}
		// get and type assert argument
		sizeArg, ok := p.Args["size"].(int)
		if !ok {
			return nil
		}
		return picResolver(sizeArg)
	}
	dataType := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "DataType",
		Fields: types.GraphQLFieldConfigMap{
			"a": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
			"b": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
			"c": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
			"d": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
			"e": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
			"f": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
			"pic": &types.GraphQLFieldConfig{
				Args: types.GraphQLFieldConfigArgumentMap{
					"size": &types.GraphQLArgumentConfig{
						Type: types.GraphQLInt,
					},
				},
				Type:    types.GraphQLString,
				Resolve: picResolverFn,
			},
		},
	})
	deepDataType := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "DeepDataType",
		Fields: types.GraphQLFieldConfigMap{
			"a": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
			"b": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
			"c": &types.GraphQLFieldConfig{
				Type: types.NewGraphQLList(types.GraphQLString),
			},
			"deeper": &types.GraphQLFieldConfig{
				Type: types.NewGraphQLList(dataType),
			},
		},
	})

	// Exploring a way to have a GraphQLObjectType within itself
	// in this case DataType has DeepDataType has DataType
	dataType.AddFieldConfig("deep", &types.GraphQLFieldConfig{
		Type: deepDataType,
	})
	// in this case DataType has DataType
	dataType.AddFieldConfig("promise", &types.GraphQLFieldConfig{
		Type: dataType,
	})

	schema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: dataType,
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	astDoc := testutil.Parse(t, query)

	// execute
	args := map[string]interface{}{
		"size": 100,
	}
	operationName := "Example"
	ep := executor.ExecuteParams{
		Schema:        schema,
		Root:          data,
		AST:           astDoc,
		OperationName: operationName,
		Args:          args,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestMergesParallelFragments(t *testing.T) {

	query := `
      { a, ...FragOne, ...FragTwo }

      fragment FragOne on Type {
        b
        deep { b, deeper: deep { b } }
      }

      fragment FragTwo on Type {
        c
        deep { c, deeper: deep { c } }
      }
    `

	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"a": "Apple",
			"b": "Banana",
			"deep": map[string]interface{}{
				"c": "Cherry",
				"b": "Banana",
				"deeper": map[string]interface{}{
					"b": "Banana",
					"c": "Cherry",
				},
			},
			"c": "Cherry",
		},
	}

	typeObjectType := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "Type",
		Fields: types.GraphQLFieldConfigMap{
			"a": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Resolve: func(p types.GQLFRParams) interface{} {
					return "Apple"
				},
			},
			"b": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Resolve: func(p types.GQLFRParams) interface{} {
					return "Banana"
				},
			},
			"c": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Resolve: func(p types.GQLFRParams) interface{} {
					return "Cherry"
				},
			},
		},
	})
	deepTypeFieldConfig := &types.GraphQLFieldConfig{
		Type: typeObjectType,
		Resolve: func(p types.GQLFRParams) interface{} {
			return p.Source
		},
	}
	typeObjectType.AddFieldConfig("deep", deepTypeFieldConfig)

	schema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: typeObjectType,
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.Parse(t, query)

	// execute
	ep := executor.ExecuteParams{
		Schema: schema,
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

func TestThreadsContextCorrectly(t *testing.T) {

	query := `
      query Example { a }
    `

	data := map[string]interface{}{
		"contextThing": "thing",
	}

	var resolvedContext map[string]interface{}

	schema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "Type",
			Fields: types.GraphQLFieldConfigMap{
				"a": &types.GraphQLFieldConfig{
					Type: types.GraphQLString,
					Resolve: func(p types.GQLFRParams) interface{} {
						resolvedContext = p.Source.(map[string]interface{})
						return resolvedContext
					},
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.Parse(t, query)

	// execute
	ep := executor.ExecuteParams{
		Schema: schema,
		Root:   data,
		AST:    ast,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}

	expected := "thing"
	if resolvedContext["contextThing"] != expected {
		t.Fatalf("Expected context.contextThing to equal %v, got %v", expected, resolvedContext["contextThing"])
	}
}

func TestCorrectlyThreadsArguments(t *testing.T) {

	query := `
      query Example {
        b(numArg: 123, stringArg: "foo")
      }
    `

	var resolvedArgs map[string]interface{}

	schema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "Type",
			Fields: types.GraphQLFieldConfigMap{
				"b": &types.GraphQLFieldConfig{
					Args: types.GraphQLFieldConfigArgumentMap{
						"numArg": &types.GraphQLArgumentConfig{
							Type: types.GraphQLInt,
						},
						"stringArg": &types.GraphQLArgumentConfig{
							Type: types.GraphQLString,
						},
					},
					Type: types.GraphQLString,
					Resolve: func(p types.GQLFRParams) interface{} {
						resolvedArgs = p.Args
						return resolvedArgs
					},
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.Parse(t, query)

	// execute
	ep := executor.ExecuteParams{
		Schema: schema,
		AST:    ast,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}

	expectedNum := 123
	expectedString := "foo"
	if resolvedArgs["numArg"] != expectedNum {
		t.Fatalf("Expected args.numArg to equal `%v`, got `%v`", expectedNum, resolvedArgs["numArg"])
	}
	if resolvedArgs["stringArg"] != expectedString {
		t.Fatalf("Expected args.stringArg to equal `%v`, got `%v`", expectedNum, resolvedArgs["stringArg"])
	}
}

func TestNullsOutErrorSubtrees(t *testing.T) {

	// TODO: TestNullsOutErrorSubtrees test for go-routines if implemented
	query := `{
      sync,
      syncError,
    }`

	expectedData := map[string]interface{}{
		"sync":      "sync",
		"syncError": nil,
	}
	expectedErrors := []graphqlerrors.GraphQLFormattedError{
		graphqlerrors.GraphQLFormattedError{
			Message: "Error getting syncError",
			Locations: []location.SourceLocation{
				location.SourceLocation{
					Line: 3, Column: 7,
				},
			},
		},
	}

	data := map[string]interface{}{
		"sync": func() interface{} {
			return "sync"
		},
		"syncError": func() interface{} {
			panic("Error getting syncError")
		},
	}
	schema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "Type",
			Fields: types.GraphQLFieldConfigMap{
				"sync": &types.GraphQLFieldConfig{
					Type: types.GraphQLString,
				},
				"syncError": &types.GraphQLFieldConfig{
					Type: types.GraphQLString,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.Parse(t, query)

	// execute
	ep := executor.ExecuteParams{
		Schema: schema,
		AST:    ast,
		Root:   data,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) == 0 {
		t.Fatalf("wrong result, expected errors, got %v", len(result.Errors))
	}
	if !reflect.DeepEqual(expectedData, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedData, result.Data))
	}
	if !reflect.DeepEqual(expectedErrors, result.Errors) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedErrors, result.Errors))
	}
}

func TestUsesTheInlineOperationIfNoOperationIsProvided(t *testing.T) {

	doc := `{ a }`
	data := map[string]interface{}{
		"a": "b",
	}

	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"a": "b",
		},
	}

	schema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "Type",
			Fields: types.GraphQLFieldConfigMap{
				"a": &types.GraphQLFieldConfig{
					Type: types.GraphQLString,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: schema,
		AST:    ast,
		Root:   data,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestUsesTheOnlyOperationIfNoOperationIsProvided(t *testing.T) {

	doc := `query Example { a }`
	data := map[string]interface{}{
		"a": "b",
	}

	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"a": "b",
		},
	}

	schema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "Type",
			Fields: types.GraphQLFieldConfigMap{
				"a": &types.GraphQLFieldConfig{
					Type: types.GraphQLString,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: schema,
		AST:    ast,
		Root:   data,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestThrowsIfNoOperationIsProvidedWithMultipleOperations(t *testing.T) {

	doc := `query Example { a } query OtherExample { a }`
	data := map[string]interface{}{
		"a": "b",
	}

	expectedErrors := []graphqlerrors.GraphQLFormattedError{
		graphqlerrors.GraphQLFormattedError{
			Message:   "Must provide operation name if query contains multiple operations.",
			Locations: []location.SourceLocation{},
		},
	}

	schema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "Type",
			Fields: types.GraphQLFieldConfigMap{
				"a": &types.GraphQLFieldConfig{
					Type: types.GraphQLString,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: schema,
		AST:    ast,
		Root:   data,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != 1 {
		t.Fatalf("wrong result, expected len(1) unexpected len: %v", len(result.Errors))
	}
	if result.Data != nil {
		t.Fatalf("wrong result, expected nil result.Data, got %v", result.Data)
	}
	if !reflect.DeepEqual(expectedErrors, result.Errors) {
		t.Fatalf("unexpected result, Diff: %v", testutil.Diff(expectedErrors, result.Errors))
	}
}

func TestUsesTheQuerySchemaForQueries(t *testing.T) {

	doc := `query Q { a } mutation M { c }`
	data := map[string]interface{}{
		"a": "b",
		"c": "d",
	}

	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"a": "b",
		},
	}

	schema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "Q",
			Fields: types.GraphQLFieldConfigMap{
				"a": &types.GraphQLFieldConfig{
					Type: types.GraphQLString,
				},
			},
		}),
		Mutation: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "M",
			Fields: types.GraphQLFieldConfigMap{
				"c": &types.GraphQLFieldConfig{
					Type: types.GraphQLString,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema:        schema,
		AST:           ast,
		Root:          data,
		OperationName: "Q",
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestUsesTheMutationSchemaForMutations(t *testing.T) {

	doc := `query Q { a } mutation M { c }`
	data := map[string]interface{}{
		"a": "b",
		"c": "d",
	}

	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"c": "d",
		},
	}

	schema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "Q",
			Fields: types.GraphQLFieldConfigMap{
				"a": &types.GraphQLFieldConfig{
					Type: types.GraphQLString,
				},
			},
		}),
		Mutation: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "M",
			Fields: types.GraphQLFieldConfigMap{
				"c": &types.GraphQLFieldConfig{
					Type: types.GraphQLString,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema:        schema,
		AST:           ast,
		Root:          data,
		OperationName: "M",
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestCorrectFieldOrderingDespiteExecutionOrder(t *testing.T) {

	doc := `
	{
      b,
      a,
      c,
      d,
      e
    }
	`
	data := map[string]interface{}{
		"a": func() interface{} { return "a" },
		"b": func() interface{} { return "b" },
		"c": func() interface{} { return "c" },
		"d": func() interface{} { return "d" },
		"e": func() interface{} { return "e" },
	}

	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"a": "a",
			"b": "b",
			"c": "c",
			"d": "d",
			"e": "e",
		},
	}

	schema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "Type",
			Fields: types.GraphQLFieldConfigMap{
				"a": &types.GraphQLFieldConfig{
					Type: types.GraphQLString,
				},
				"b": &types.GraphQLFieldConfig{
					Type: types.GraphQLString,
				},
				"c": &types.GraphQLFieldConfig{
					Type: types.GraphQLString,
				},
				"d": &types.GraphQLFieldConfig{
					Type: types.GraphQLString,
				},
				"e": &types.GraphQLFieldConfig{
					Type: types.GraphQLString,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: schema,
		AST:    ast,
		Root:   data,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}

	// TODO: test to ensure key ordering
	// The following does not work
	// - iterating over result.Data map
	//   Note that golang's map iteration order is randomized
	//   So, iterating over result.Data won't do it for a test
	// - Marshal the result.Data to json string and assert it
	//   json.Marshal seems to re-sort the keys automatically
	//
	t.Skipf("TODO: Ensure key ordering")
}

func TestAvoidsRecursion(t *testing.T) {

	doc := `
      query Q {
        a
        ...Frag
        ...Frag
      }

      fragment Frag on Type {
        a,
        ...Frag
      }
    `
	data := map[string]interface{}{
		"a": "b",
	}

	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"a": "b",
		},
	}

	schema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "Type",
			Fields: types.GraphQLFieldConfigMap{
				"a": &types.GraphQLFieldConfig{
					Type: types.GraphQLString,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema:        schema,
		AST:           ast,
		Root:          data,
		OperationName: "Q",
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}

}

func TestDoesNotIncludeIllegalFieldsInOutput(t *testing.T) {

	doc := `mutation M {
      thisIsIllegalDontIncludeMe
    }`

	expected := &types.GraphQLResult{
		Data: map[string]interface{}{},
	}

	schema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "Q",
			Fields: types.GraphQLFieldConfigMap{
				"a": &types.GraphQLFieldConfig{
					Type: types.GraphQLString,
				},
			},
		}),
		Mutation: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "M",
			Fields: types.GraphQLFieldConfigMap{
				"c": &types.GraphQLFieldConfig{
					Type: types.GraphQLString,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: schema,
		AST:    ast,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, expected len(%v) errors, got len(%v)", len(expected.Errors), len(result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDoesNotIncludeArgumentsThatWereNotSet(t *testing.T) {

	doc := `{ field(a: true, c: false, e: 0) }`

	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"field": `{"a":true,"c":false,"e":0}`,
		},
	}

	schema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "Type",
			Fields: types.GraphQLFieldConfigMap{
				"field": &types.GraphQLFieldConfig{
					Type: types.GraphQLString,
					Args: types.GraphQLFieldConfigArgumentMap{
						"a": &types.GraphQLArgumentConfig{
							Type: types.GraphQLBoolean,
						},
						"b": &types.GraphQLArgumentConfig{
							Type: types.GraphQLBoolean,
						},
						"c": &types.GraphQLArgumentConfig{
							Type: types.GraphQLBoolean,
						},
						"d": &types.GraphQLArgumentConfig{
							Type: types.GraphQLInt,
						},
						"e": &types.GraphQLArgumentConfig{
							Type: types.GraphQLInt,
						},
					},
					Resolve: func(p types.GQLFRParams) interface{} {
						args, _ := json.Marshal(p.Args)
						return string(args)
					},
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: schema,
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

type testSpecialType struct {
	Value string
}
type testNotSpecialType struct {
	Value string
}

func TestFailsWhenAnIsTypeOfCheckIsNotMet(t *testing.T) {

	query := `{ specials { value } }`

	data := map[string]interface{}{
		"specials": []interface{}{
			testSpecialType{"foo"},
			testNotSpecialType{"bar"},
		},
	}

	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"specials": []interface{}{
				map[string]interface{}{
					"value": "foo",
				},
				nil,
			},
		},
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message:   `Expected value of type "SpecialType" but got: executor_test.testNotSpecialType.`,
				Locations: []location.SourceLocation{},
			},
		},
	}

	specialType := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "SpecialType",
		IsTypeOf: func(value interface{}, info types.GraphQLResolveInfo) bool {
			if _, ok := value.(testSpecialType); ok {
				return true
			}
			return false
		},
		Fields: types.GraphQLFieldConfigMap{
			"value": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Resolve: func(p types.GQLFRParams) interface{} {
					return p.Source.(testSpecialType).Value
				},
			},
		},
	})
	schema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "Query",
			Fields: types.GraphQLFieldConfigMap{
				"specials": &types.GraphQLFieldConfig{
					Type: types.NewGraphQLList(specialType),
					Resolve: func(p types.GQLFRParams) interface{} {
						return p.Source.(map[string]interface{})["specials"]
					},
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.Parse(t, query)

	// execute
	ep := executor.ExecuteParams{
		Schema: schema,
		AST:    ast,
		Root:   data,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) == 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestFailsToExecuteQueryContainingATypeDefinition(t *testing.T) {

	query := `
      { foo }

      type Query { foo: String }
	`
	expected := &types.GraphQLResult{
		Data: nil,
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message:   "GraphQL cannot execute a request containing a ObjectTypeDefinition",
				Locations: []location.SourceLocation{},
			},
		},
	}

	schema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "Query",
			Fields: types.GraphQLFieldConfigMap{
				"foo": &types.GraphQLFieldConfig{
					Type: types.GraphQLString,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.Parse(t, query)

	// execute
	ep := executor.ExecuteParams{
		Schema: schema,
		AST:    ast,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != 1 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
