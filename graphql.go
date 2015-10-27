package graphql

import (
	"github.com/chris-ramon/graphql/errors"
	"github.com/chris-ramon/graphql/executor"
	"github.com/chris-ramon/graphql/language/parser"
	"github.com/chris-ramon/graphql/language/source"
	"github.com/chris-ramon/graphql/types"
	"github.com/chris-ramon/graphql/validator"
)

type Params struct {
	Schema         types.Schema
	RequestString  string
	RootObject     map[string]interface{}
	VariableValues map[string]interface{}
	OperationName  string
}

func Graphql(p Params, resultChannel chan *types.Result) {
	source := source.NewSource(&source.Source{
		Body: p.RequestString,
		Name: "GraphQL request",
	})
	AST, err := parser.Parse(parser.ParseParams{Source: source})
	if err != nil {
		result := types.Result{
			Errors: graphqlerrors.FormatErrors(err),
		}
		resultChannel <- &result
		return
	}
	validationResult := validator.ValidateDocument(p.Schema, AST)

	if !validationResult.IsValid {
		result := types.Result{
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
