package graphql

import (
	"github.com/chris-ramon/graphql/gqlerrors"
	"github.com/chris-ramon/graphql/language/source"
)

type Params struct {
	Schema         Schema
	RequestString  string
	RootObject     map[string]interface{}
	VariableValues map[string]interface{}
	OperationName  string
}

func Graphql(p Params, resultChannel chan *Result) {
	source := source.NewSource(&source.Source{
		Body: p.RequestString,
		Name: "GraphQL request",
	})
	AST, err := Parse(ParseParams{Source: source})
	if err != nil {
		result := Result{
			Errors: gqlerrors.FormatErrors(err),
		}
		resultChannel <- &result
		return
	}
	validationResult := ValidateDocument(p.Schema, AST)

	if !validationResult.IsValid {
		result := Result{
			Errors: validationResult.Errors,
		}
		resultChannel <- &result
		return
	} else {
		ep := ExecuteParams{
			Schema:        p.Schema,
			Root:          p.RootObject,
			AST:           AST,
			OperationName: p.OperationName,
			Args:          p.VariableValues,
		}
		Execute(ep, resultChannel)
		return
	}
}
