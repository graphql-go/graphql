package executor_test

import (
	"fmt"
	"github.com/chris-ramon/graphql-go/executor"
	"github.com/chris-ramon/graphql-go/language/parser"
	"github.com/chris-ramon/graphql-go/types"
	"github.com/kr/pretty"
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
		"promise": func() interface{} {
			// TODO: instead of promise, let's try go-routines
			return "TODO  go-routines"
		},
	}
	deepData = map[string]interface{}{
		"a":      func() string { return "Already Been Done" },
		"b":      func() string { return "Boring" },
		"c":      func() []string { return []string{"Contrived", "", "Confusing"} },
		"deeper": func() []interface{} { return []interface{}{data, nil, nil} },
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
	pparams := parser.ParseParams{
		Source: query,
		Options: parser.ParseOptions{
			NoSource: true,
		},
	}
	astDoc, err := parser.Parse(pparams)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	picResolverFn := func(p types.GQLFRParams) interface{} {
		return p.Source["pic"].(func(size int) string)(1000)
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

	schema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: dataType,
	})

	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

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
	pretty.Println(result)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
}
