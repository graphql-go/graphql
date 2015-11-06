package graphql

import "github.com/graphql-go/graphql/language/ast"

func ValidateDocument(schema *Schema, ast *ast.Document) (bool, []error) {
	return true, nil
}
