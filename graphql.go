package graphql

import (
	"github.com/sprucehealth/graphql/gqlerrors"
	"github.com/sprucehealth/graphql/language/parser"
	"github.com/sprucehealth/graphql/language/source"
	"golang.org/x/net/context"
)

type Params struct {
	Schema         Schema
	RequestString  string
	RootObject     map[string]interface{}
	VariableValues map[string]interface{}
	OperationName  string

	// Context may be provided to pass application-specific per-request
	// information to resolve functions.
	Context context.Context
}

func Do(p Params) *Result {
	source := source.New("GraphQL request", p.RequestString)
	ast, err := parser.Parse(parser.ParseParams{Source: source})
	if err != nil {
		return &Result{
			Errors: gqlerrors.FormatErrors(err),
		}
	}
	validationResult := ValidateDocument(&p.Schema, ast, nil)

	if !validationResult.IsValid {
		return &Result{
			Errors: validationResult.Errors,
		}
	}

	return Execute(ExecuteParams{
		Schema:        p.Schema,
		Root:          p.RootObject,
		AST:           ast,
		OperationName: p.OperationName,
		Args:          p.VariableValues,
		Context:       p.Context,
	})
}
