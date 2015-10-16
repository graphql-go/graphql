package validator

import (
	"github.com/chris-ramon/graphql-go/errors"
	"github.com/chris-ramon/graphql-go/language/ast"
	"github.com/tallstreet/graphql-go/types"
)

type ValidationResult struct {
	IsValid bool
	Errors  []graphqlerrors.GraphQLFormattedError
}

func ValidateDocument(schema types.GraphQLSchema, ast *ast.Document) (vr ValidationResult) {
	vr.IsValid = true
	return vr
}
