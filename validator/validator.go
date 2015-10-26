package validator

import (
	"github.com/chris-ramon/graphql/errors"
	"github.com/chris-ramon/graphql/language/ast"
	"github.com/chris-ramon/graphql/types"
)

type ValidationResult struct {
	IsValid bool
	Errors  []graphqlerrors.GraphQLFormattedError
}

func ValidateDocument(schema types.GraphQLSchema, ast *ast.Document) (vr ValidationResult) {
	vr.IsValid = true
	return vr
}
