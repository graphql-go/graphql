package ast

import (
	"github.com/sprucehealth/graphql/language/kinds"
)

// Name implements Node
type Name struct {
	Kind  string
	Loc   *Location
	Value string
}

func NewName(node *Name) *Name {
	if node == nil {
		return &Name{Kind: kinds.Name}
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
