package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// FragmentSpread implements Node, Selection
type FragmentSpread struct {
	Kind       string
	Loc        Location
	Name       *Name
	Directives []Directive
}

func NewFragmentSpread(fs *FragmentSpread) *FragmentSpread {
	if fs == nil {
		fs = &FragmentSpread{}
	}
	return &FragmentSpread{
		Kind:       kinds.FragmentSpread,
		Loc:        fs.Loc,
		Name:       fs.Name,
		Directives: fs.Directives,
	}
}

func (fs *FragmentSpread) GetKind() string {
	return fs.Kind
}

func (fs *FragmentSpread) GetLoc() Location {
	return fs.Loc
}
