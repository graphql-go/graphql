package graphql

import (
	"github.com/chris-ramon/graphql/gqlerrors"
	"github.com/chris-ramon/graphql/language/ast"
)

type ValidationResult struct {
	IsValid bool
	Errors  []gqlerrors.FormattedError
}

func ValidateDocument(schema Schema, ast *ast.Document) (vr ValidationResult) {
	vr.IsValid = true
	return vr
}
