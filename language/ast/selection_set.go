package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// SelectionSet implements Node
type SelectionSet struct {
	Kind       string
	Loc        Location
	Selections []Selection
}

func NewSelectionSet(v *SelectionSet) *SelectionSet {
	if v == nil {
		v = &SelectionSet{}
	}
	return &SelectionSet{
		Kind:       kinds.SelectionSet,
		Loc:        v.Loc,
		Selections: v.Selections,
	}
}

func (v *SelectionSet) GetKind() string {
	return v.Kind
}

func (v *SelectionSet) GetLoc() Location {
	return v.Loc
}
