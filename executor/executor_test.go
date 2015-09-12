package executor_test

import (
	"github.com/chris-ramon/graphql-go/executor"
	"github.com/chris-ramon/graphql-go/language/parser"
	"github.com/chris-ramon/graphql-go/types"
	"testing"
	"github.com/kr/pretty"
	"fmt"
)

func TestExecutesArbritraryCode(t *testing.T) {
	resultChannel := make(chan *types.GraphQLResult)

	deepData := map[string]interface{} {}
	data := map[string]interface{} {
		"a": func () string { return "Apple" },
		"b": func () string { return "Banana" },
		"c": func () string { return "Cookie" },
		"d": func () string { return "Donut" },
		"e": func () string { return "Egg" },
		"f": "Fish",
		"pic": func(size int) string {
			return fmt.Sprintf("Pic of size: %v", size)
		},
		"deep": func() interface{} { return deepData },
		"promise": func() interface{} {
			// instead of promise, let's try go-routines
			return "TODO  go-routines"
		},
	}
	deepData = map[string]interface{} {
		"a": func () string { return "Already Been Done" },
		"b": func () string { return "Boring" },
		"c": func () []string { return []string{"Contrived", "", "Confusing"} },
		"deeper": func () []interface{} { return []interface{}{data, nil, nil} },
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
		pretty.Println("----==> picResolverFn p.ARG", p.Args)
		return fmt.Sprintf("$$$$$Pic of size: %v", 100000)
	}
	schema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.GraphQLObjectType{
			Name: "DataType",
			Fields: types.GraphQLFieldDefinitionMap{
				"a": types.GraphQLFieldDefinition{
					Type: &types.GraphQLString{},
				},
				"b": types.GraphQLFieldDefinition{
					Type: &types.GraphQLString{},
				},
				"c": types.GraphQLFieldDefinition{
					Type: &types.GraphQLString{},
				},
				"d": types.GraphQLFieldDefinition{
					Type: &types.GraphQLString{},
				},
				"e": types.GraphQLFieldDefinition{
					Type: &types.GraphQLString{},
				},
				"f": types.GraphQLFieldDefinition{
					Type: &types.GraphQLString{},
				},
				"pic": types.GraphQLFieldDefinition{
					Args: []types.GraphQLFieldArgument{
						types.GraphQLFieldArgument{
							Name: "size",
							Type: &types.GraphQLInt{},
						},
					},
					Type: &types.GraphQLString{},
					Resolve: picResolverFn,
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("Error in schema", err.Error())
	}

	root := map[string]interface{}{}
	args := map[string]interface{}{
		"size": 100,
	}
	operationName := "Example"
	ep := executor.ExecuteParams{
		Schema:        schema,
		Root:          root,
		AST:           astDoc,
		OperationName: operationName,
		Args:          args,
	}
	go executor.Execute(ep, resultChannel)
	result := <-resultChannel
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
}
