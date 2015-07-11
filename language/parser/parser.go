package parser

import "github.com/chris-ramon/graphql-go/language/ast"

type ParseOptions struct {
	NoLocation bool
	NoSource   bool
}

func Parse(source interface{}, options ParseOptions) (doc ast.Document) {
	return doc
}
