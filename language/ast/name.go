package ast

import "github.com/chris-ramon/graphql-go/language/kinds"

// Name implements Node
type Name struct {
	Kind  string
	Loc   Location
	Value string
}

func NewName(n *Name) *Name {
	if n == nil {
		n = &Name{}
	}
	return &Name{
		Kind:  kinds.Name,
		Value: n.Value,
		Loc:   n.Loc,
	}
}

func (node *Name) GetKind() string {
	return node.Kind
}

func (node *Name) GetLoc() Location {
	return node.Loc
}
