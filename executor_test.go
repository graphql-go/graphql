package graphql_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

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
	picResolverFn := func(p graphql.ResolveParams) interface{} {
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
				Resolve: func(p graphql.ResolveParams) interface{} {
					return "Apple"
				},
			},
			"b": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) interface{} {
					return "Banana"
				},
			},
			"c": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) interface{} {
					return "Cherry"
				},
			},
		},
	})
	deepTypeFieldConfig := &graphql.Field{
		Type: typeObjectType,
		Resolve: func(p graphql.ResolveParams) interface{} {
			return p.Source
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

func TestThreadsContextCorrectly(t *testing.T) {

	query := `
      query Example { a }
    `

	data := map[string]interface{}{
		"contextThing": "thing",
	}

	var resolvedContext map[string]interface{}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Type",
			Fields: graphql.Fields{
				"a": &graphql.Field{
					Type: graphql.String,
					Resolve: func(p graphql.ResolveParams) interface{} {
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
					Resolve: func(p graphql.ResolveParams) interface{} {
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
	expectedErrors := []gqlerrors.FormattedError{
		gqlerrors.FormattedError{
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

func TestUsesTheInlineOperationIfNoOperationIsProvided(t *testing.T) {

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

func TestUsesTheOnlyOperationIfNoOperationIsProvided(t *testing.T) {

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

func TestThrowsIfNoOperationIsProvidedWithMultipleOperations(t *testing.T) {

	doc := `query Example { a } query OtherExample { a }`
	data := map[string]interface{}{
		"a": "b",
	}

	expectedErrors := []gqlerrors.FormattedError{
		gqlerrors.FormattedError{
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

func TestUsesTheQuerySchemaForQueries(t *testing.T) {

	doc := `query Q { a } mutation M { c }`
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
					Resolve: func(p graphql.ResolveParams) interface{} {
						args, err := json.Marshal(p.Args)
						if err != nil {
							panic(err)
						}
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

	expected := &graphql.Result{
		Data: map[string]interface{}{
			"specials": []interface{}{
				map[string]interface{}{
					"value": "foo",
				},
				nil,
			},
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message:   `Expected value of type "SpecialType" but got: graphql_test.testNotSpecialType.`,
				Locations: []location.SourceLocation{},
			},
		},
	}

	specialType := graphql.NewObject(graphql.ObjectConfig{
		Name: "SpecialType",
		IsTypeOf: func(value interface{}, info graphql.ResolveInfo) bool {
			if _, ok := value.(testSpecialType); ok {
				return true
			}
			return false
		},
		Fields: graphql.Fields{
			"value": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) interface{} {
					return p.Source.(testSpecialType).Value
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
					Resolve: func(p graphql.ResolveParams) interface{} {
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
			gqlerrors.FormattedError{
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
