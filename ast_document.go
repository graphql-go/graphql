package graphql

import (
	"github.com/chris-ramon/graphql/kinds"
)

// Document implements Node
type AstDocument struct {
	Kind        string
	Loc         *AstLocation
	Definitions []Node
}

func NewAstDocument(d *AstDocument) *AstDocument {
	if d == nil {
		d = &AstDocument{}
	}
	return &AstDocument{
		Kind:        kinds.Document,
		Loc:         d.Loc,
		Definitions: d.Definitions,
	}
}

func (node *AstDocument) GetKind() string {
	return node.Kind
}

func (node *AstDocument) GetLoc() *AstLocation {
	return node.Loc
}
