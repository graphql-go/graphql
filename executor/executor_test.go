package executor_test

import (
	"github.com/chris-ramon/graphql-go/executor"
	"github.com/chris-ramon/graphql-go/language/parser"
	"github.com/chris-ramon/graphql-go/types"
	"github.com/kr/pretty"
	"testing"
)

func TestExecutesArbritraryCode(t *testing.T) {
	resultChannel := make(chan *types.GraphQLResult)

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
	schema := types.GraphQLSchema{
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
				},
			},
		},
	}
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
	pretty.Println("result", result)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
}
