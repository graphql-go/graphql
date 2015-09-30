package gql

import (
	"github.com/chris-ramon/graphql-go/errors"
	"github.com/chris-ramon/graphql-go/executor"
	"github.com/chris-ramon/graphql-go/gqltypes"
	"github.com/chris-ramon/graphql-go/language/parser"
	"github.com/chris-ramon/graphql-go/language/source"
	"github.com/chris-ramon/graphql-go/validator"
)

type GraphqlParams struct {
	Schema         gqltypes.GraphQLSchema
	RequestString  string
	RootObject     map[string]interface{}
	VariableValues map[string]interface{}
	OperationName  string
}

func Graphql(p GraphqlParams, resultChannel chan *gqltypes.GraphQLResult) {
	source := source.NewSource(&source.Source{
		Body: p.RequestString,
		Name: "GraphQL request",
	})
	AST, err := parser.Parse(parser.ParseParams{Source: source})
	if err != nil {
		result := gqltypes.GraphQLResult{
			Errors: graphqlerrors.FormatErrors(err),
		}
		resultChannel <- &result
		return
	}
	validationResult := validator.ValidateDocument(p.Schema, AST)

	if !validationResult.IsValid {
		result := gqltypes.GraphQLResult{
			Errors: validationResult.Errors,
		}
		resultChannel <- &result
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
