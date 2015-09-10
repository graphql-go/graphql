package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// InlineFragment implements Node, Selection
type InlineFragment struct {
	Kind          string
	Loc           Location
	TypeCondition NamedType
	Directives    []Directive
	SelectionSet  *SelectionSet
}

func NewInlineFragment(f *InlineFragment) *InlineFragment {
	if f == nil {
		f = &InlineFragment{}
	}
	return &InlineFragment{
		Kind:          kinds.InlineFragment,
		Loc:           f.Loc,
		TypeCondition: f.TypeCondition,
		Directives:    f.Directives,
		SelectionSet:  f.SelectionSet,
	}
}

func (f *InlineFragment) GetKind() string {
	return f.Kind
}

func (f *InlineFragment) GetLoc() Location {
	return f.Loc
}
