package graphql

func NewLocatedError(err interface{}, nodes []Node) *Error {
	message := "An unknown error occurred."
	if err, ok := err.(error); ok {
		message = err.Error()
	}
	if err, ok := err.(string); ok {
		message = err
	}
	stack := message
	return NewError(
		message,
		nodes,
		stack,
		nil,
		[]int{},
	)
}

func FieldASTsToNodeASTs(fieldASTs []*AstField) []Node {
	nodes := []Node{}
	for _, fieldAST := range fieldASTs {
		nodes = append(nodes, fieldAST)
	}
	return nodes
}
