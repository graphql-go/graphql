package ast

import (
	"github.com/bigdrum/graphql/language/kinds"
)

// Name implements Node
type Name struct {
	Kind  string
	Loc   *Location
	Value string
}

func NewName(node *Name) *Name {
	if node == nil {
		node = &Name{}
	}
	return &Name{
		Kind:  kinds.Name,
		Value: node.Value,
		Loc:   node.Loc,
	}
}

func (node *Name) GetKind() string {
	return node.Kind
}

func (node *Name) GetLoc() *Location {
	return node.Loc
}
