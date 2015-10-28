package graphql

import (
	"fmt"
)

type Error struct {
	Message   string
	Stack     string
	Nodes     []Node
	Source    *Source
	Positions []int
	Locations []SourceLocation
}

// implements Golang's built-in `error` interface
func (g Error) Error() string {
	return fmt.Sprintf("%v", g.Message)
}

func NewError(message string, nodes []Node, stack string, source *Source, positions []int) *Error {
	if stack == "" && message != "" {
		stack = message
	}
	if source == nil {
		for _, node := range nodes {
			// get source from first node
			if node.GetLoc() != nil {
				source = node.GetLoc().Source
			}
			break
		}
	}
	if len(positions) == 0 && len(nodes) > 0 {
		for _, node := range nodes {
			if node.GetLoc() == nil {
				continue
			}
			positions = append(positions, node.GetLoc().Start)
		}
	}
	locations := []SourceLocation{}
	for _, pos := range positions {
		loc := GetLocation(source, pos)
		locations = append(locations, loc)
	}
	return &Error{
		Message:   message,
		Stack:     stack,
		Nodes:     nodes,
		Source:    source,
		Positions: positions,
		Locations: locations,
	}
}
