package graphql

import (
	"github.com/chris-ramon/graphql/kinds"
)

// AstName implements Node
type AstName struct {
	Kind  string
	Loc   *AstLocation
	Value string
}

func NewAstName(node *AstName) *AstName {
	if node == nil {
		node = &AstName{}
	}
	return &AstName{
		Kind:  kinds.Name,
		Value: node.Value,
		Loc:   node.Loc,
	}
}

func (node *AstName) GetKind() string {
	return node.Kind
}

func (node *AstName) GetLoc() *AstLocation {
	return node.Loc
}
