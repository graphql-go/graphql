package validator

import (
	"github.com/chris-ramon/graphql-go/errors"
	"github.com/chris-ramon/graphql-go/gqltypes"
	"github.com/chris-ramon/graphql-go/language/ast"
)

type ValidationResult struct {
	IsValid bool
	Errors  []graphqlerrors.GraphQLFormattedError
}

func ValidateDocument(schema gqltypes.GraphQLSchema, ast *ast.Document) (vr ValidationResult) {
	vr.IsValid = true
	return vr
}
