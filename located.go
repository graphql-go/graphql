package graphql

import (
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/ast"
)

func NewLocatedError(err interface{}, nodes []ast.Node) *gqlerrors.Error {
	message := "An unknown error occurred."
	if err, ok := err.(error); ok {
		message = err.Error()
	}
	if err, ok := err.(string); ok {
		message = err
	}
	stack := message
	return gqlerrors.NewError(
		message,
		nodes,
		stack,
		nil,
		[]int{},
	)
}

func FieldASTsToNodeASTs(fieldASTs []*ast.Field) []ast.Node {
	nodes := []ast.Node{}
	for _, fieldAST := range fieldASTs {
		nodes = append(nodes, fieldAST)
	}
	return nodes
}
