package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// Field implements Node, Selection
type Field struct {
	Kind         string
	Loc          *Location
	Alias        *Name
	Name         *Name
	Arguments    []*Argument
	Directives   []*Directive
	SelectionSet *SelectionSet
}

func NewField(f *Field) *Field {
	if f == nil {
		f = &Field{}
	}
	return &Field{
		Kind:         kinds.Field,
		Loc:          f.Loc,
		Alias:        f.Alias,
		Name:         f.Name,
		Arguments:    f.Arguments,
		Directives:   f.Directives,
		SelectionSet: f.SelectionSet,
	}
}

func (f *Field) GetKind() string {
	return f.Kind
}

func (f *Field) GetLoc() *Location {
	return f.Loc
}
