package graphql

import (
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
)

func Do(schema *Schema, request string) *Result {
	source := source.NewSource("GraphQL request", request)
	document, err := parser.Parse(source, nil)
	if err != nil {
		return &Result{
			Errors: gqlerrors.FormatErrors(err),
		}
	}

	valid, errs := ValidateDocument(schema, document)
	if !valid {
		return &Result{
			Errors: gqlerrors.FormatErrors(errs...),
		}
	}

	return Execute(schema, document)
}
