package ast

import (
	"github.com/chris-ramon/graphql-go/language/source"
)

type Location struct {
	Start  int
	End    int
	Source *source.Source
}

type SelectionSet struct {
	Kind       string
	Loc        Location
	Selections []interface{}
}

func NewSelectionSet() *SelectionSet {
	return &SelectionSet{
		Kind: "SelectionSet",
	}
}

type Argument struct {
	Kind  string
	Loc   Location
	Name  *Name
	Value Value
}

func NewArgument() *Name {
	return &Name{
		Kind: "Argument",
	}
}

type Directive struct {
	Kind      string
	Loc       Location
	Name      *Name
	Arguments []Argument
}

func NewDirective() *Directive {
	return &Directive{
		Kind: "Directive",
	}
}
