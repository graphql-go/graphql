package graphqlerrors

import (

	"fmt"
	"github.com/chris-ramon/graphql-go/language/ast"
	"github.com/chris-ramon/graphql-go/language/location"
	"github.com/chris-ramon/graphql-go/language/source"
)

type GraphQLError struct {
	Message   string
	Stack     string
	Nodes     []ast.Node
	Source    *source.Source
	Positions []int
	Locations []location.SourceLocation
}

// implements Golang's built-in `error` interface
func (g GraphQLError) Error() string {
	return fmt.Sprintf("%v", g.Message)
}

func NewGraphQLError(message string, nodes []ast.Node, stack string, source *source.Source, positions []int) *GraphQLError {
	if stack == "" && message != "" {
		stack = message
	}
	locations := []location.SourceLocation{}
	for _, pos := range positions {
		loc := location.GetLocation(source, pos)
		locations = append(locations, loc)
	}
	if len(positions) == 0 && len(nodes) > 0 {
		for _, node := range nodes {
			if node.GetLoc() == nil {
				continue
			}
			positions = append(positions, node.GetLoc().Start)
		}
	}
	return &GraphQLError{
		Message:   message,
		Stack:     stack,
		Nodes:     nodes,
		Source:    source,
		Positions: positions,
		Locations: locations,
	}
}
