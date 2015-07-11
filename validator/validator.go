package validator

import (
	"github.com/chris-ramon/graphql-go/errors"
	"github.com/chris-ramon/graphql-go/language/ast"
	"github.com/chris-ramon/graphql-go/types"
)

type ValidationResult struct {
	IsValid bool
	Errors  []errors.GraphQLFormattedError
}

func ValidateDocument(schema types.GraphQLSchema, ast ast.Document) (vr ValidationResult) {
	return vr
}
