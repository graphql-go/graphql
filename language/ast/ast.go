package ast

import "github.com/chris-ramon/graphql-go/language/source"

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

type Definition interface {
	GetKind() string
	GetLoc() Location
	GetOperation() string
	GetName() Name
	GetVariableDefinitions() []VariableDefinition
	GetTypeCondition() Name
	GetDirectives() []Directive
	GetSelectionSet() SelectionSet
}

type Argument struct {
	Kind  string
	Loc   Location
	Name  Name
	Value Value
}

func NewArgument() *Name {
	return &Name{
		Kind: "Argument",
	}
}

type Field struct {
	Kind         string
	Loc          Location
	Alias        Name
	Name         Name
	Arguments    []Argument
	Directives   []Directive
	SelectionSet SelectionSet
}

func NewField() *Name {
	return &Name{
		Kind: "Field",
	}
}

type Value interface {
	GetKind() string
	GetLoc() Location
	GetName() Name
}

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

type Variable struct {
	Kind string
	Loc  Location
	Name Name
}

func NewVariable() *Variable {
	return &Variable{
		Kind: "Variable",
	}
}

type VariableDefinition struct {
	Kind         string
	Loc          Location
	Variable     Variable
	Type         interface{}
	DefaultValue Value
}

func NewVariableDefinition() *VariableDefinition {
	return &VariableDefinition{
		Kind: "VariableDefinition",
	}
}
