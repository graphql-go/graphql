package executor_test

import (
	"fmt"
	"github.com/chris-ramon/graphql-go/executor"
	"github.com/chris-ramon/graphql-go/language/parser"
	"github.com/chris-ramon/graphql-go/types"
	"github.com/kr/pretty"
	"reflect"
	"testing"
)

func TestExecutesArbritraryCode(t *testing.T) {
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

	// TODO make GraphQLResult json marshal friendly
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
	astDoc, err := parser.Parse(parser.ParseParams{
		Source: query,
		Options: parser.ParseOptions{
			NoSource: true,
		},
	})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

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
