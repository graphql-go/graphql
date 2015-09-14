package executor_test

import (
	"fmt"
	"github.com/chris-ramon/graphql-go/errors"
	"github.com/chris-ramon/graphql-go/executor"
	"github.com/chris-ramon/graphql-go/language/ast"
	"github.com/chris-ramon/graphql-go/language/location"
	"github.com/chris-ramon/graphql-go/language/parser"
	"github.com/chris-ramon/graphql-go/types"
	"github.com/kr/pretty"
	"reflect"
	"testing"
)

func parse(query string, t *testing.T) *ast.Document {
	astDoc, err := parser.Parse(parser.ParseParams{
		Source: query,
		Options: parser.ParseOptions{
			NoSource: true,
		},
	})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	return astDoc
}

func TestExecutesArbitraryCode(t *testing.T) {
	resultChannel := make(chan *types.GraphQLResult)

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
		"deeper": func() interface{} { return []interface{}{data, nil, nil} },
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

	// TODO make GraphQLResult json friendly
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"b": "Banana",
			"x": "Cookie",
			"d": "Donut",
			"e": "Egg",
			"promise": &types.GraphQLResult{
				Data: map[string]interface{}{
					"a": "Apple",
				},
			},
			"a": "Apple",
			"deep": &types.GraphQLResult{
				Data: map[string]interface{}{
					"a": "Already Been Done",
					"b": "Boring",
					"c": []interface{}{
						"Contrived",
						"",
						"Confusing",
					},
					"deeper": []interface{}{
						&types.GraphQLResult{
							Data: map[string]interface{}{
								"a": "Apple",
								"b": "Banana",
							},
						},
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
		picResolver, ok := p.Source["pic"].(func(size int) string)
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
						Type:         types.GraphQLInt,
						DefaultValue: 100,
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
	astDoc := parse(query, t)

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
	go executor.Execute(ep, resultChannel)
	result := <-resultChannel
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", pretty.Diff(expected, result))
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
			"deep": &types.GraphQLResult{
				Data: map[string]interface{}{
					"c": "Cherry",
					"b": "Banana",
					"deeper": &types.GraphQLResult{
						Data: map[string]interface{}{
							"b": "Banana",
							"c": "Cherry",
						},
						Errors: nil,
					},
				},
				Errors: nil,
			},
			"c": "Cherry",
		},
		Errors: nil,
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
	ast := parse(query, t)

	// execute
	ep := executor.ExecuteParams{
		Schema: schema,
		AST:    ast,
	}
	resultChannel := make(chan *types.GraphQLResult)
	go executor.Execute(ep, resultChannel)
	result := <-resultChannel
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", pretty.Diff(expected, result))
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
						resolvedContext = p.Source
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
	ast := parse(query, t)

	// execute
	ep := executor.ExecuteParams{
		Schema: schema,
		Root:   data,
		AST:    ast,
	}
	resultChannel := make(chan *types.GraphQLResult)
	go executor.Execute(ep, resultChannel)
	result := <-resultChannel
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
							Type: types.GraphQLString,
						},
						"stringArg": &types.GraphQLArgumentConfig{
							Type: types.GraphQLInt,
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
	ast := parse(query, t)

	// execute
	ep := executor.ExecuteParams{
		Schema: schema,
		AST:    ast,
	}
	resultChannel := make(chan *types.GraphQLResult)
	go executor.Execute(ep, resultChannel)
	result := <-resultChannel
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
		"sync": "sync",
	}
	expectedErrors := []graphqlerrors.GraphQLFormattedError{
		graphqlerrors.GraphQLFormattedError{
			Message:   "Error getting syncError",
			Locations: []location.SourceLocation{},
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
	ast := parse(query, t)

	// execute
	ep := executor.ExecuteParams{
		Schema: schema,
		AST:    ast,
		Root:   data,
	}
	resultChannel := make(chan *types.GraphQLResult)
	go executor.Execute(ep, resultChannel)
	result := <-resultChannel
	if len(result.Errors) == 0 {
		t.Fatalf("wrong result, expected errors, got %v", len(result.Errors))
	}
	if !reflect.DeepEqual(expectedData, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", pretty.Diff(expectedData, result.Data))
	}
	if !reflect.DeepEqual(expectedErrors, result.Errors) {
		t.Fatalf("Unexpected result, Diff: %v", pretty.Diff(expectedErrors, result.Errors))
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
	ast := parse(doc, t)

	// execute
	ep := executor.ExecuteParams{
		Schema: schema,
		AST:    ast,
		Root:   data,
	}
	resultChannel := make(chan *types.GraphQLResult)
	go executor.Execute(ep, resultChannel)
	result := <-resultChannel
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", pretty.Diff(expected, result))
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
	ast := parse(doc, t)

	// execute
	ep := executor.ExecuteParams{
		Schema: schema,
		AST:    ast,
		Root:   data,
	}
	resultChannel := make(chan *types.GraphQLResult)
	go executor.Execute(ep, resultChannel)
	result := <-resultChannel
	if len(result.Errors) != 1 {
		t.Fatalf("wrong result, expected len(1) unexpected len: %v", len(result.Errors))
	}
	if result.Data != nil {
		t.Fatalf("wrong result, expected nil result.Data, got %v", result.Data)
	}
	if !reflect.DeepEqual(expectedErrors, result.Errors) {
		t.Fatalf("unexpected result, Diff: %v", pretty.Diff(expectedErrors, result.Errors))
	}
}
