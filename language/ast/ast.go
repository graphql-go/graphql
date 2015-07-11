package ast

import "github.com/chris-ramon/graphql-go/language/source"

type Definition interface {
}

type Name struct {
	Kind  string
	Loc   Location
	Value string
}

func NewName() *Name {
	return &Name{
		Kind: "Name",
	}
}

type SelectionSet struct {
	Kind       string
	Location   Location
	Selections []interface{}
}

func NewSelectionSet() *SelectionSet {
	return &SelectionSet{
		Kind: "SelectionSet",
	}
}

type Value interface{}

type Directive struct {
	Kind  string
	Loc   Location
	Name  Name
	Value Value
}

func NewDirective() *Directive {
	return &Directive{
		Kind: "Directive",
	}
}

type FragmentDefinition struct {
	Kind          string
	Loc           Location
	Name          Name
	TypeCondition Name
	Directives    []Directive
	SelectionSet  SelectionSet
}

func NewFragmentDefinition() *FragmentDefinition {
	return &FragmentDefinition{
		Kind: "FragmentDefinition",
	}
}

type Location struct {
	Start  int
	End    int
	Source source.Source
}

type Document struct {
	Kind        string
	Loc         Location
	Definitions []Definition
}

type OperationDefinition interface{}
