package graphqlerrors

import (
	"go/ast"

	"github.com/chris-ramon/graphql-go/language/source"
)

type GraphQLFormattedError struct {
	Message   string
	Locations []struct {
		Line   int
		Column int
	}
}

type GraphQLError struct {
	Error     error
	Stack     string
	Nodes     []ast.Node
	Source    *source.Source
	Positions []int
	Locations interface{}
}
