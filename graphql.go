package graphql

import (
	"context"

	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
)

type Params struct {
	// The GraphQL type system to use when validating and executing a query.
	Schema Schema

	// A GraphQL language formatted string representing the requested operation.
	RequestString string

	// The value provided as the first argument to resolver functions on the top
	// level type (e.g. the query object type).
	RootObject map[string]interface{}

	// A mapping of variable name to runtime value to use for all variables
	// defined in the requestString.
	VariableValues map[string]interface{}

	// The name of the operation to use if requestString contains multiple
	// possible operations. Can be omitted if requestString contains only
	// one operation.
	OperationName string

	// Context may be provided to pass application-specific per-request
	// information to resolve functions.
	Context context.Context
}

func Do(p Params) *Result {
	source := source.NewSource(&source.Source{
		Body: []byte(p.RequestString),
		Name: "GraphQL request",
	})
	handleExtensionsInits(&p)

	// parse the source
	handleExtensionsParseDidStart(&p)
	AST, err := parser.Parse(parser.ParseParams{Source: source})
	if err != nil {
		handleExtensionsParseEnded(&p, err)
		return &Result{
			Errors: gqlerrors.FormatErrors(err),
		}
	}
	handleExtensionsParseEnded(&p, err)

	// validate document
	handleExtensionsValidationDidStart(&p)
	validationResult := ValidateDocument(&p.Schema, AST, nil)

	if !validationResult.IsValid {
		handleExtensionsValidationEnded(&p, validationResult.Errors)
		return &Result{
			Errors: validationResult.Errors,
		}
	}
	handleExtensionsValidationEnded(&p, validationResult.Errors)

	return Execute(ExecuteParams{
		Schema:        p.Schema,
		Root:          p.RootObject,
		AST:           AST,
		OperationName: p.OperationName,
		Args:          p.VariableValues,
		Context:       p.Context,
	})
}
