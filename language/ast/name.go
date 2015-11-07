package ast

import (
	"github.com/graphql-go/graphql/language/kinds"
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
	if node == nil {
		return ""
	}
	return node.Kind
}

func (node *Name) GetLoc() *Location {
	if node == nil {
		return nil
	}
	return node.Loc
}

func (node *Name) GetValue() string {
	if node == nil {
		return ""
	}
	return node.Value
}
