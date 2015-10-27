package validator

import (
	"github.com/chris-ramon/graphql/errors"
	"github.com/chris-ramon/graphql/language/ast"
	"github.com/chris-ramon/graphql/types"
)

type ValidationResult struct {
	IsValid bool
	Errors  []graphqlerrors.FormattedError
}

func ValidateDocument(schema types.Schema, ast *ast.Document) (vr ValidationResult) {
	vr.IsValid = true
	return vr
}
