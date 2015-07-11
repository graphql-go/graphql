package gql

import (
	"github.com/chris-ramon/graphql-go/executor"
	"github.com/chris-ramon/graphql-go/language/parser"
	"github.com/chris-ramon/graphql-go/language/source"
	"github.com/chris-ramon/graphql-go/types"
	"github.com/chris-ramon/graphql-go/validator"
)

type GraphqlParams struct {
	Schema         types.GraphQLSchema
	RequestString  string
	RootObject     map[string]interface{}
	VariableValues map[string]string
	OperationName  string
}

func Graphql(p GraphqlParams, resultChannel chan types.GraphQLResult) {
	source := source.NewSource(p.RequestString, "GraphQL request")
	AST := parser.Parse(source, parser.ParseOptions{})
	validationResult := validator.ValidateDocument(p.Schema, AST)
	if !validationResult.IsValid {
		result := types.GraphQLResult{
			Errors: validationResult.Errors,
		}
		resultChannel <- result
		return
	} else {
		ep := executor.ExecuteParams{
			Schema:        p.Schema,
			Root:          p.RootObject,
			AST:           AST,
			OperationName: p.OperationName,
			Args:          p.VariableValues,
		}
		executor.Execute(ep, resultChannel)
		return
	}
}
