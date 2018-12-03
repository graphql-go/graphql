package graphql_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/location"
	"github.com/graphql-go/graphql/testutil"
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

	expected := &graphql.Result{
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
					"",
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
	picResolverFn := func(p graphql.ResolveParams) (interface{}, error) {
		// get and type assert ResolveFn for this field
		picResolver, ok := p.Source.(map[string]interface{})["pic"].(func(size int) string)
		if !ok {
			return nil, nil
		}
		// get and type assert argument
		sizeArg, ok := p.Args["size"].(int)
		if !ok {
			return nil, nil
		}
		return picResolver(sizeArg), nil
	}
	dataType := graphql.NewObject(graphql.ObjectConfig{
		Name: "DataType",
		Fields: graphql.Fields{
			"a": &graphql.Field{
				Type: graphql.String,
			},
			"b": &graphql.Field{
				Type: graphql.String,
			},
			"c": &graphql.Field{
				Type: graphql.String,
			},
			"d": &graphql.Field{
				Type: graphql.String,
			},
			"e": &graphql.Field{
				Type: graphql.String,
			},
			"f": &graphql.Field{
				Type: graphql.String,
			},
			"pic": &graphql.Field{
				Args: graphql.FieldConfigArgument{
					"size": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Type:    graphql.String,
				Resolve: picResolverFn,
			},
		},
	})
	deepDataType := graphql.NewObject(graphql.ObjectConfig{
		Name: "DeepDataType",
		Fields: graphql.Fields{
			"a": &graphql.Field{
				Type: graphql.String,
			},
			"b": &graphql.Field{
				Type: graphql.String,
			},
			"c": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
			"deeper": &graphql.Field{
				Type: graphql.NewList(dataType),
			},
		},
	})

	// Exploring a way to have a Object within itself
	// in this case DataType has DeepDataType has DataType
	dataType.AddFieldConfig("deep", &graphql.Field{
		Type: deepDataType,
	})
	// in this case DataType has DataType
	dataType.AddFieldConfig("promise", &graphql.Field{
		Type: dataType,
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: dataType,
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	astDoc := testutil.TestParse(t, query)

	// execute
	args := map[string]interface{}{
		"size": 100,
	}
	operationName := "Example"
	ep := graphql.ExecuteParams{
		Schema:        schema,
		Root:          data,
		AST:           astDoc,
		OperationName: operationName,
		Args:          args,
	}
	result := testutil.TestExecute(t, ep)
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

	expected := &graphql.Result{
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

	typeObjectType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Type",
		Fields: graphql.Fields{
			"a": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "Apple", nil
				},
			},
			"b": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "Banana", nil
				},
			},
			"c": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "Cherry", nil
				},
			},
		},
	})
	deepTypeFieldConfig := &graphql.Field{
		Type: typeObjectType,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return p.Source, nil
		},
	}
	typeObjectType.AddFieldConfig("deep", deepTypeFieldConfig)

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: typeObjectType,
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.TestParse(t, query)

	// execute
	ep := graphql.ExecuteParams{
		Schema: schema,
		AST:    ast,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestThreadsSourceCorrectly(t *testing.T) {

	query := `
      query Example { a }
    `

	data := map[string]interface{}{
		"key": "value",
	}

	var resolvedSource map[string]interface{}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Type",
			Fields: graphql.Fields{
				"a": &graphql.Field{
					Type: graphql.String,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						resolvedSource = p.Source.(map[string]interface{})
						return resolvedSource, nil
					},
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.TestParse(t, query)

	// execute
	ep := graphql.ExecuteParams{
		Schema: schema,
		Root:   data,
		AST:    ast,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}

	expected := "value"
	if resolvedSource["key"] != expected {
		t.Fatalf("Expected context.key to equal %v, got %v", expected, resolvedSource["key"])
	}
}

func TestCorrectlyThreadsArguments(t *testing.T) {

	query := `
      query Example {
        b(numArg: 123, stringArg: "foo")
      }
    `

	var resolvedArgs map[string]interface{}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Type",
			Fields: graphql.Fields{
				"b": &graphql.Field{
					Args: graphql.FieldConfigArgument{
						"numArg": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
						"stringArg": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Type: graphql.String,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						resolvedArgs = p.Args
						return resolvedArgs, nil
					},
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.TestParse(t, query)

	// execute
	ep := graphql.ExecuteParams{
		Schema: schema,
		AST:    ast,
	}
	result := testutil.TestExecute(t, ep)
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

func TestThreadsRootValueContextCorrectly(t *testing.T) {

	query := `
      query Example { a }
    `

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Type",
			Fields: graphql.Fields{
				"a": &graphql.Field{
					Type: graphql.String,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						val, _ := p.Info.RootValue.(map[string]interface{})["stringKey"].(string)
						return val, nil
					},
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.TestParse(t, query)

	// execute
	ep := graphql.ExecuteParams{
		Schema: schema,
		AST:    ast,
		Root: map[string]interface{}{
			"stringKey": "stringValue",
		},
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}

	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "stringValue",
		},
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestThreadsContextCorrectly(t *testing.T) {

	query := `
      query Example { a }
    `

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Type",
			Fields: graphql.Fields{
				"a": &graphql.Field{
					Type: graphql.String,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return p.Context.Value("foo"), nil
					},
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.TestParse(t, query)

	// execute
	ep := graphql.ExecuteParams{
		Schema:  schema,
		AST:     ast,
		Context: context.WithValue(context.Background(), "foo", "bar"),
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}

	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "bar",
		},
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
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
	originalError := errors.New("Error getting syncError")
	expectedErrors := []gqlerrors.FormattedError{gqlerrors.FormatError(gqlerrors.Error{
		Message: originalError.Error(),
		Locations: []location.SourceLocation{
			{
				Line: 3, Column: 7,
			},
		},
		Path: []interface{}{
			"syncError",
		},
		OriginalError: originalError,
	}),
	}

	data := map[string]interface{}{
		"sync": func() interface{} {
			return "sync"
		},
		"syncError": func() interface{} {
			panic("Error getting syncError")
		},
	}
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Type",
			Fields: graphql.Fields{
				"sync": &graphql.Field{
					Type: graphql.String,
				},
				"syncError": &graphql.Field{
					Type: graphql.String,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.TestParse(t, query)

	// execute
	ep := graphql.ExecuteParams{
		Schema: schema,
		AST:    ast,
		Root:   data,
	}
	result := testutil.TestExecute(t, ep)
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

func TestUsesTheInlineOperationIfNoOperationNameIsProvided(t *testing.T) {

	doc := `{ a }`
	data := map[string]interface{}{
		"a": "b",
	}

	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "b",
		},
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Type",
			Fields: graphql.Fields{
				"a": &graphql.Field{
					Type: graphql.String,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: schema,
		AST:    ast,
		Root:   data,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestUsesTheOnlyOperationIfNoOperationNameIsProvided(t *testing.T) {

	doc := `query Example { a }`
	data := map[string]interface{}{
		"a": "b",
	}

	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "b",
		},
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Type",
			Fields: graphql.Fields{
				"a": &graphql.Field{
					Type: graphql.String,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: schema,
		AST:    ast,
		Root:   data,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestUsesTheNamedOperationIfOperationNameIsProvided(t *testing.T) {

	doc := `query Example { first: a } query OtherExample { second: a }`
	data := map[string]interface{}{
		"a": "b",
	}

	expected := &graphql.Result{
		Data: map[string]interface{}{
			"second": "b",
		},
	}
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Type",
			Fields: graphql.Fields{
				"a": &graphql.Field{
					Type: graphql.String,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema:        schema,
		AST:           ast,
		Root:          data,
		OperationName: "OtherExample",
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestThrowsIfNoOperationIsProvided(t *testing.T) {

	doc := `fragment Example on Type { a }`
	data := map[string]interface{}{
		"a": "b",
	}

	expectedErrors := []gqlerrors.FormattedError{
		{
			Message:   "Must provide an operation.",
			Locations: []location.SourceLocation{},
		},
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Type",
			Fields: graphql.Fields{
				"a": &graphql.Field{
					Type: graphql.String,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: schema,
		AST:    ast,
		Root:   data,
	}
	result := testutil.TestExecute(t, ep)
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
func TestThrowsIfNoOperationNameIsProvidedWithMultipleOperations(t *testing.T) {

	doc := `query Example { a } query OtherExample { a }`
	data := map[string]interface{}{
		"a": "b",
	}

	expectedErrors := []gqlerrors.FormattedError{
		{
			Message:   "Must provide operation name if query contains multiple operations.",
			Locations: []location.SourceLocation{},
		},
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Type",
			Fields: graphql.Fields{
				"a": &graphql.Field{
					Type: graphql.String,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: schema,
		AST:    ast,
		Root:   data,
	}
	result := testutil.TestExecute(t, ep)
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

func TestThrowsIfUnknownOperationNameIsProvided(t *testing.T) {

	doc := `query Example { a } query OtherExample { a }`
	data := map[string]interface{}{
		"a": "b",
	}

	expectedErrors := []gqlerrors.FormattedError{
		{
			Message:   `Unknown operation named "UnknownExample".`,
			Locations: []location.SourceLocation{},
		},
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Type",
			Fields: graphql.Fields{
				"a": &graphql.Field{
					Type: graphql.String,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema:        schema,
		AST:           ast,
		Root:          data,
		OperationName: "UnknownExample",
	}
	result := testutil.TestExecute(t, ep)
	if result.Data != nil {
		t.Fatalf("wrong result, expected nil result.Data, got %v", result.Data)
	}
	if !reflect.DeepEqual(expectedErrors, result.Errors) {
		t.Fatalf("unexpected result, Diff: %v", testutil.Diff(expectedErrors, result.Errors))
	}
}
func TestUsesTheQuerySchemaForQueries(t *testing.T) {

	doc := `query Q { a } mutation M { c } subscription S { a }`
	data := map[string]interface{}{
		"a": "b",
		"c": "d",
	}

	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "b",
		},
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Q",
			Fields: graphql.Fields{
				"a": &graphql.Field{
					Type: graphql.String,
				},
			},
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name: "M",
			Fields: graphql.Fields{
				"c": &graphql.Field{
					Type: graphql.String,
				},
			},
		}),
		Subscription: graphql.NewObject(graphql.ObjectConfig{
			Name: "S",
			Fields: graphql.Fields{
				"a": &graphql.Field{
					Type: graphql.String,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema:        schema,
		AST:           ast,
		Root:          data,
		OperationName: "Q",
	}
	result := testutil.TestExecute(t, ep)
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

	expected := &graphql.Result{
		Data: map[string]interface{}{
			"c": "d",
		},
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Q",
			Fields: graphql.Fields{
				"a": &graphql.Field{
					Type: graphql.String,
				},
			},
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name: "M",
			Fields: graphql.Fields{
				"c": &graphql.Field{
					Type: graphql.String,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema:        schema,
		AST:           ast,
		Root:          data,
		OperationName: "M",
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestUsesTheSubscriptionSchemaForSubscriptions(t *testing.T) {

	doc := `query Q { a } subscription S { a }`
	data := map[string]interface{}{
		"a": "b",
		"c": "d",
	}

	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "b",
		},
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Q",
			Fields: graphql.Fields{
				"a": &graphql.Field{
					Type: graphql.String,
				},
			},
		}),
		Subscription: graphql.NewObject(graphql.ObjectConfig{
			Name: "S",
			Fields: graphql.Fields{
				"a": &graphql.Field{
					Type: graphql.String,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema:        schema,
		AST:           ast,
		Root:          data,
		OperationName: "S",
	}
	result := testutil.TestExecute(t, ep)
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

	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "a",
			"b": "b",
			"c": "c",
			"d": "d",
			"e": "e",
		},
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Type",
			Fields: graphql.Fields{
				"a": &graphql.Field{
					Type: graphql.String,
				},
				"b": &graphql.Field{
					Type: graphql.String,
				},
				"c": &graphql.Field{
					Type: graphql.String,
				},
				"d": &graphql.Field{
					Type: graphql.String,
				},
				"e": &graphql.Field{
					Type: graphql.String,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: schema,
		AST:    ast,
		Root:   data,
	}
	result := testutil.TestExecute(t, ep)
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

	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "b",
		},
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Type",
			Fields: graphql.Fields{
				"a": &graphql.Field{
					Type: graphql.String,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema:        schema,
		AST:           ast,
		Root:          data,
		OperationName: "Q",
	}
	result := testutil.TestExecute(t, ep)
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

	expected := &graphql.Result{
		Data: map[string]interface{}{},
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Q",
			Fields: graphql.Fields{
				"a": &graphql.Field{
					Type: graphql.String,
				},
			},
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name: "M",
			Fields: graphql.Fields{
				"c": &graphql.Field{
					Type: graphql.String,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: schema,
		AST:    ast,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, expected len(%v) errors, got len(%v)", len(expected.Errors), len(result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDoesNotIncludeArgumentsThatWereNotSet(t *testing.T) {

	doc := `{ field(a: true, c: false, e: 0) }`

	expected := &graphql.Result{
		Data: map[string]interface{}{
			"field": `{"a":true,"c":false,"e":0}`,
		},
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Type",
			Fields: graphql.Fields{
				"field": &graphql.Field{
					Type: graphql.String,
					Args: graphql.FieldConfigArgument{
						"a": &graphql.ArgumentConfig{
							Type: graphql.Boolean,
						},
						"b": &graphql.ArgumentConfig{
							Type: graphql.Boolean,
						},
						"c": &graphql.ArgumentConfig{
							Type: graphql.Boolean,
						},
						"d": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
						"e": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						args, _ := json.Marshal(p.Args)
						return string(args), nil
					},
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: schema,
		AST:    ast,
	}
	result := testutil.TestExecute(t, ep)
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

	originalError := gqlerrors.NewFormattedError(`Expected value of type "SpecialType" but got: graphql_test.testNotSpecialType.`)
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"specials": []interface{}{
				map[string]interface{}{
					"value": "foo",
				},
				nil,
			},
		},
		Errors: []gqlerrors.FormattedError{gqlerrors.FormatError(gqlerrors.Error{
			Message: originalError.Message,
			Locations: []location.SourceLocation{
				{
					Line:   1,
					Column: 3,
				},
			},
			Path: []interface{}{
				"specials",
				1,
			},
			OriginalError: originalError,
		})},
	}

	specialType := graphql.NewObject(graphql.ObjectConfig{
		Name: "SpecialType",
		IsTypeOf: func(p graphql.IsTypeOfParams) bool {
			if _, ok := p.Value.(testSpecialType); ok {
				return true
			}
			return false
		},
		Fields: graphql.Fields{
			"value": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(testSpecialType).Value, nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"specials": &graphql.Field{
					Type: graphql.NewList(specialType),
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return p.Source.(map[string]interface{})["specials"], nil
					},
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.TestParse(t, query)

	// execute
	ep := graphql.ExecuteParams{
		Schema: schema,
		AST:    ast,
		Root:   data,
	}
	result := testutil.TestExecute(t, ep)
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
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message:   "GraphQL cannot execute a request containing a ObjectDefinition",
				Locations: []location.SourceLocation{},
			},
		},
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"foo": &graphql.Field{
					Type: graphql.String,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.TestParse(t, query)

	// execute
	ep := graphql.ExecuteParams{
		Schema: schema,
		AST:    ast,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != 1 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestQuery_ExecutionAddsErrorsFromFieldResolveFn(t *testing.T) {
	qError := errors.New("queryError")
	q := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"a": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return nil, qError
				},
			},
			"b": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "ok", nil
				},
			},
		},
	})
	blogSchema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: q,
	})
	if err != nil {
		t.Fatalf("unexpected error, got: %v", err)
	}
	query := "{ a }"
	result := graphql.Do(graphql.Params{
		Schema:        blogSchema,
		RequestString: query,
	})
	if len(result.Errors) == 0 {
		t.Fatal("wrong result, expected errors, got no errors")
	}
	if result.Errors[0].Error() != qError.Error() {
		t.Fatalf("wrong result, unexpected error, got: %v, expected: %v", result.Errors[0], qError)
	}
}

func TestQuery_ExecutionDoesNotAddErrorsFromFieldResolveFn(t *testing.T) {
	qError := errors.New("queryError")
	q := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"a": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return nil, qError
				},
			},
			"b": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "ok", nil
				},
			},
		},
	})
	blogSchema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: q,
	})
	if err != nil {
		t.Fatalf("unexpected error, got: %v", err)
	}
	query := "{ b }"
	result := graphql.Do(graphql.Params{
		Schema:        blogSchema,
		RequestString: query,
	})
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %+v", result.Errors)
	}
}

func TestQuery_InputObjectUsesFieldDefaultValueFn(t *testing.T) {
	inputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "Input",
		Fields: graphql.InputObjectConfigFieldMap{
			"default": &graphql.InputObjectFieldConfig{
				Type:         graphql.String,
				DefaultValue: "bar",
			},
		},
	})
	q := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"a": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"foo": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(inputType),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					val := p.Args["foo"].(map[string]interface{})
					def, ok := val["default"]
					if !ok || def == nil {
						return nil, errors.New("queryError: No 'default' param")
					}
					if def.(string) != "bar" {
						return nil, errors.New("queryError: 'default' param has wrong value")
					}
					return "ok", nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: q,
	})
	if err != nil {
		t.Fatalf("unexpected error, got: %v", err)
	}
	query := `{ a(foo: {}) }`
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %+v", result.Errors)
	}
}

func TestMutation_ExecutionAddsErrorsFromFieldResolveFn(t *testing.T) {
	mError := errors.New("mutationError")
	q := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"a": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
	m := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"foo": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"f": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return nil, mError
				},
			},
			"bar": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"b": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "ok", nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    q,
		Mutation: m,
	})
	if err != nil {
		t.Fatalf("unexpected error, got: %v", err)
	}
	query := "mutation _ { newFoo: foo(f:\"title\") }"
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) == 0 {
		t.Fatal("wrong result, expected errors, got no errors")
	}
	if result.Errors[0].Error() != mError.Error() {
		t.Fatalf("wrong result, unexpected error, got: %v, expected: %v", result.Errors[0], mError)
	}
}

func TestMutation_ExecutionDoesNotAddErrorsFromFieldResolveFn(t *testing.T) {
	mError := errors.New("mutationError")
	q := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"a": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
	m := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"foo": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"f": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return nil, mError
				},
			},
			"bar": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"b": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "ok", nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    q,
		Mutation: m,
	})
	if err != nil {
		t.Fatalf("unexpected error, got: %v", err)
	}
	query := "mutation _ { newBar: bar(b:\"title\") }"
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %+v", result.Errors)
	}
}

func TestGraphqlTag(t *testing.T) {
	typeObjectType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Type",
		Fields: graphql.Fields{
			"fooBar": &graphql.Field{Type: graphql.String},
		},
	})
	var baz = &graphql.Field{
		Type:        typeObjectType,
		Description: "typeObjectType",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			t := struct {
				FooBar string `graphql:"fooBar"`
			}{"foo bar value"}
			return t, nil
		},
	}
	q := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"baz": baz,
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: q,
	})
	if err != nil {
		t.Fatalf("unexpected error, got: %v", err)
	}
	query := "{ baz { fooBar } }"
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %+v", result.Errors)
	}
	expectedData := map[string]interface{}{
		"baz": map[string]interface{}{
			"fooBar": "foo bar value",
		},
	}
	if !reflect.DeepEqual(result.Data, expectedData) {
		t.Fatalf("unexpected result, got: %+v, expected: %+v", expectedData, result.Data)
	}
}

func TestFieldResolver(t *testing.T) {
	typeObjectType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Type",
		Fields: graphql.Fields{
			"fooBar": &graphql.Field{Type: graphql.String},
		},
	})
	var baz = &graphql.Field{
		Type:        typeObjectType,
		Description: "typeObjectType",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return testCustomResolver{}, nil
		},
	}
	var bazPtr = &graphql.Field{
		Type:        typeObjectType,
		Description: "typeObjectType",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return &testCustomResolver{}, nil
		},
	}
	q := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"baz":    baz,
			"bazPtr": bazPtr,
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: q,
	})
	if err != nil {
		t.Fatalf("unexpected error, got: %v", err)
	}
	query := "{ baz { fooBar }, bazPtr { fooBar } }"
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %+v", result.Errors)
	}
	expectedData := map[string]interface{}{
		"baz": map[string]interface{}{
			"fooBar": "foo bar value",
		},
		"bazPtr": map[string]interface{}{
			"fooBar": "foo bar value",
		},
	}
	if !reflect.DeepEqual(result.Data, expectedData) {
		t.Fatalf("unexpected result, got: %+v, expected: %+v", result.Data, expectedData)
	}
}

type testCustomResolver struct{}

func (r testCustomResolver) Resolve(p graphql.ResolveParams) (interface{}, error) {
	if p.Info.FieldName == "fooBar" {
		return "foo bar value", nil
	}
	return "", errors.New("invalid field " + p.Info.FieldName)
}

func TestContextDeadline(t *testing.T) {
	timeout := time.Millisecond * time.Duration(100)
	acceptableDelay := time.Millisecond * time.Duration(10)
	expectedErrors := []gqlerrors.FormattedError{
		{
			Message:   context.DeadlineExceeded.Error(),
			Locations: []location.SourceLocation{},
		},
	}

	// Query type includes a field that won't resolve within the deadline
	var queryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"hello": &graphql.Field{
					Type: graphql.String,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						time.Sleep(2 * time.Second)
						return "world", nil
					},
				},
			},
		})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
	if err != nil {
		t.Fatalf("unexpected error, got: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	startTime := time.Now()
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: "{hello}",
		Context:       ctx,
	})
	duration := time.Since(startTime)

	if duration > timeout+acceptableDelay {
		t.Fatalf("graphql.Do completed in %s, should have completed in %s", duration, timeout)
	}
	if !result.HasErrors() || len(result.Errors) == 0 {
		t.Fatalf("Result should include errors when deadline is exceeded")
	}
	if !reflect.DeepEqual(expectedErrors, result.Errors) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedErrors, result.Errors))
	}
}

func TestThunkResultsProcessedCorrectly(t *testing.T) {
	barType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Bar",
		Fields: graphql.Fields{
			"bazA": &graphql.Field{
				Type: graphql.String,
			},
			"bazB": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	fooType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Foo",
		Fields: graphql.Fields{
			"bar": &graphql.Field{
				Type: barType,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					var bar struct {
						BazA string
						BazB string
					}
					bar.BazA = "A"
					bar.BazB = "B"

					thunk := func() (interface{}, error) { return &bar, nil }
					return thunk, nil
				},
			},
		},
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"foo": &graphql.Field{
				Type: fooType,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					var foo struct{}
					return foo, nil
				},
			},
		},
	})

	expectNoError := func(err error) {
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
	expectNoError(err)

	query := "{ foo { bar { bazA bazB } } }"
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) != 0 {
		t.Fatalf("expected no errors, got %v", result.Errors)
	}

	foo := result.Data.(map[string]interface{})["foo"].(map[string]interface{})
	bar, ok := foo["bar"].(map[string]interface{})

	if !ok {
		t.Errorf("expected bar to be a map[string]interface{}: actual = %v", reflect.TypeOf(foo["bar"]))
	} else {
		if got, want := bar["bazA"], "A"; got != want {
			t.Errorf("foo.bar.bazA: got=%v, want=%v", got, want)
		}
		if got, want := bar["bazB"], "B"; got != want {
			t.Errorf("foo.bar.bazB: got=%v, want=%v", got, want)
		}
	}

	if t.Failed() {
		b, err := json.Marshal(result.Data)
		expectNoError(err)
		t.Log(string(b))
	}
}

func TestThunkErrorsAreHandledCorrectly(t *testing.T) {
	var bazCError = errors.New("barC error")
	barType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Bar",
		Fields: graphql.Fields{
			"bazA": &graphql.Field{
				Type: graphql.String,
			},
			"bazB": &graphql.Field{
				Type: graphql.String,
			},
			"bazC": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					thunk := func() (interface{}, error) {
						return nil, bazCError
					}
					return thunk, nil
				},
			},
		},
	})

	fooType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Foo",
		Fields: graphql.Fields{
			"bar": &graphql.Field{
				Type: barType,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					var bar struct {
						BazA string
						BazB string
					}
					bar.BazA = "A"
					bar.BazB = "B"

					thunk := func() (interface{}, error) {
						return &bar, nil
					}
					return thunk, nil
				},
			},
		},
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"foo": &graphql.Field{
				Type: fooType,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					var foo struct{}
					return foo, nil
				},
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})

	if err != nil {
		t.Fatalf("unexpected error, got: %v", err)
	}

	query := "{ foo { bar { bazA bazB bazC } } }"
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})

	foo := result.Data.(map[string]interface{})["foo"].(map[string]interface{})
	bar, ok := foo["bar"].(map[string]interface{})

	if !ok {
		t.Errorf("expected bar to be a map[string]interface{}: actual = %v", reflect.TypeOf(foo["bar"]))
	} else {
		if got, want := bar["bazA"], "A"; got != want {
			t.Errorf("foo.bar.bazA: got=%v, want=%v", got, want)
		}
		if got, want := bar["bazB"], "B"; got != want {
			t.Errorf("foo.bar.bazB: got=%v, want=%v", got, want)
		}
		if got := bar["bazC"]; got != nil {
			t.Errorf("foo.bar.bazC: got=%v, want=nil", got)
		}
		var errs = result.Errors
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %v", result.Errors)
		}
		if got, want := errs[0].Message, bazCError.Error(); got != want {
			t.Errorf("expected error: got=%v, want=%v", got, want)
		}
	}

	if t.Failed() {
		b, err := json.Marshal(result.Data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		t.Log(string(b))
	}
}

func assertJSON(t *testing.T, expected string, actual interface{}) {
	var e interface{}
	if err := json.Unmarshal([]byte(expected), &e); err != nil {
		t.Fatalf(err.Error())
	}
	aJSON, err := json.MarshalIndent(actual, "", "  ")
	if err != nil {
		t.Fatalf(err.Error())
	}
	var a interface{}
	if err := json.Unmarshal(aJSON, &a); err != nil {
		t.Fatalf(err.Error())
	}
	if !reflect.DeepEqual(e, a) {
		eNormalizedJSON, err := json.MarshalIndent(e, "", "  ")
		if err != nil {
			t.Fatalf(err.Error())
		}
		t.Fatalf("Expected JSON:\n\n%v\n\nActual JSON:\n\n%v", string(eNormalizedJSON), string(aJSON))
	}
}

type extendedError struct {
	error
	extensions map[string]interface{}
}

func (err extendedError) Extensions() map[string]interface{} {
	return err.extensions
}

var _ gqlerrors.ExtendedError = &extendedError{}

func testErrors(t *testing.T, nameType graphql.Output, extensions map[string]interface{}, formatErrorFn func(err error) error) *graphql.Result {
	type Hero struct {
		Id      string `graphql:"id"`
		Name    string
		Friends []Hero `graphql:"friends"`
	}

	var heroFields graphql.Fields

	heroType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Hero",
		Fields: graphql.FieldsThunk(func() graphql.Fields {
			return heroFields
		}),
	})

	heroFields = graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
		},
		"name": &graphql.Field{
			Type: nameType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				hero := p.Source.(Hero)
				if hero.Name != "" {
					return hero.Name, nil
				}

				err := fmt.Errorf("Name for character with ID %v could not be fetched.", hero.Id)
				if formatErrorFn != nil {
					err = formatErrorFn(err)
				}

				if extensions != nil {
					return nil, &extendedError{
						error:      err,
						extensions: extensions,
					}
				}
				return nil, err
			},
		},
		"friends": &graphql.Field{
			Type: graphql.NewList(heroType),
		},
	}

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"hero": &graphql.Field{
				Type: heroType,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return Hero{
						Name: "R2-D2",
						Friends: []Hero{
							{Id: "1000", Name: "Luke Skywalker"},
							{Id: "1002"},
							{Id: "1003", Name: "Leia Organa"},
						},
					}, nil
				},
			},
		},
	})

	expectNoError := func(err error) {
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
	expectNoError(err)

	return graphql.Do(graphql.Params{
		Schema: schema,
		RequestString: `{
  hero {
    name
    heroFriends: friends {
      id
      name
    }
  }
}`,
	})
}

// http://facebook.github.io/graphql/June2018/#example-bc485
func TestQuery_ErrorPath(t *testing.T) {
	result := testErrors(t, graphql.String, nil, nil)

	assertJSON(t, `{
	  "errors": [
		{
		  "message": "Name for character with ID 1002 could not be fetched.",
		  "locations": [ { "line": 6, "column": 7 } ],
		  "path": [ "hero", "heroFriends", 1, "name" ]
		}
	  ],
	  "data": {
		"hero": {
		  "name": "R2-D2",
		  "heroFriends": [
			{
			  "id": "1000",
			  "name": "Luke Skywalker"
			},
			{
			  "id": "1002",
			  "name": null
			},
			{
			  "id": "1003",
			  "name": "Leia Organa"
			}
		  ]
		}
	  }
	}`, result)
}

// http://facebook.github.io/graphql/June2018/#example-08b62
func TestQuery_ErrorPathForNonNullField(t *testing.T) {
	result := testErrors(t, graphql.NewNonNull(graphql.String), nil, nil)

	assertJSON(t, `{
	  "errors": [
		{
		  "message": "Name for character with ID 1002 could not be fetched.",
		  "locations": [ { "line": 6, "column": 7 } ],
		  "path": [ "hero", "heroFriends", 1, "name" ]
		}
	  ],
	  "data": {
		"hero": {
		  "name": "R2-D2",
		  "heroFriends": [
			{
			  "id": "1000",
			  "name": "Luke Skywalker"
			},
			null,
			{
			  "id": "1003",
			  "name": "Leia Organa"
			}
		  ]
		}
	  }
	}`, result)
}

// http://facebook.github.io/graphql/June2018/#example-fce18
func TestQuery_ErrorExtensions(t *testing.T) {
	result := testErrors(t, graphql.NewNonNull(graphql.String), map[string]interface{}{
		"code":      "CAN_NOT_FETCH_BY_ID",
		"timestamp": "Fri Feb 9 14:33:09 UTC 2018",
	}, nil)

	assertJSON(t, `{
	  "errors": [
		{
		  "message": "Name for character with ID 1002 could not be fetched.",
		  "locations": [ { "line": 6, "column": 7 } ],
		  "path": [ "hero", "heroFriends", 1, "name" ],
		  "extensions": {
			  "code": "CAN_NOT_FETCH_BY_ID",
			  "timestamp": "Fri Feb 9 14:33:09 UTC 2018"
		  }}
	  ],
	  "data": {
		"hero": {
		  "name": "R2-D2",
		  "heroFriends": [
			{
			  "id": "1000",
			  "name": "Luke Skywalker"
			},
			null,
			{
			  "id": "1003",
			  "name": "Leia Organa"
			}
		  ]
		}
	  }
	}`, result)
}

func TestQuery_OriginalErrorBuiltin(t *testing.T) {
	result := testErrors(t, graphql.String, nil, nil)
	originalError := result.Errors[0].OriginalError()
	switch originalError.(type) {
	case error:
	default:
		t.Fatalf("unexpected error: %v", reflect.TypeOf(originalError))
	}
}

func TestQuery_OriginalErrorExtended(t *testing.T) {
	result := testErrors(t, graphql.String, map[string]interface{}{
		"code": "CAN_NOT_FETCH_BY_ID",
	}, nil)
	originalError := result.Errors[0].OriginalError()
	switch originalError.(type) {
	case *extendedError:
	case extendedError:
	default:
		t.Fatalf("unexpected error: %v", reflect.TypeOf(originalError))
	}
}

type customError struct {
	error
}

func (e customError) Error() string {
	return e.error.Error()
}

func TestQuery_OriginalErrorCustom(t *testing.T) {
	result := testErrors(t, graphql.String, nil, func(err error) error {
		return customError{error: err}
	})
	originalError := result.Errors[0].OriginalError()
	switch originalError.(type) {
	case customError:
	default:
		t.Fatalf("unexpected error: %v", reflect.TypeOf(originalError))
	}
}

func TestQuery_OriginalErrorCustomPtr(t *testing.T) {
	result := testErrors(t, graphql.String, nil, func(err error) error {
		return &customError{error: err}
	})
	originalError := result.Errors[0].OriginalError()
	switch originalError.(type) {
	case *customError:
	default:
		t.Fatalf("unexpected error: %v", reflect.TypeOf(originalError))
	}
}

func TestQuery_OriginalErrorPanic(t *testing.T) {
	result := testErrors(t, graphql.String, nil, func(err error) error {
		panic(errors.New("panic error"))
	})
	originalError := result.Errors[0].OriginalError()
	switch originalError.(type) {
	case error:
	default:
		t.Fatalf("unexpected error: %v", reflect.TypeOf(originalError))
	}
}
