package graphqlerrors

import (
	"go/ast"
)

func NewLocatedError(err error, nodes []ast.Node) *GraphQLError {
	message := "An unknown error occurred."
	if err != nil {
		message = err.Error()
	}
	stack := message
	return NewGraphQLError(
		message,
		nodes,
		stack,
		nil,
		[]int{},
	)
}
