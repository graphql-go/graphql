package ast

import (
	"github.com/chris-ramon/graphql-go/language/source"
)

type Location struct {
	Start  int
	End    int
	Source *source.Source
}

type Node interface {
	GetKind() string
	GetLoc()  Location
}

// Name

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

// Document

type Document struct {
	Kind        string
	Loc         Location
	Definitions []Node
}

type Definition interface {
	// Most basic Definition interface
	GetKind() string
	GetLoc() Location

	// Some implementation may or may not have properties for these following getters()
	// But `executor` requires these getters().
	GetName() Name
	GetOperation() string
	GetVariableDefinitions() []VariableDefinition
	GetTypeCondition() NamedType
	GetDirectives() []Directive
	GetSelectionSet() SelectionSet
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

type SelectionSet struct {
	Kind       string
	Loc        Location
	Selections []interface{}
}

type Selection interface{}

func NewSelectionSet() *SelectionSet {
	return &SelectionSet{
		Kind: "SelectionSet",
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

// Fragments

type FragmentSpread struct {
	Kind       string
	Loc        Location
	Name       Name
	Directives []Directive
}

type InlineFragment struct {
	Kind          string
	Loc           Location
	TypeCondition NamedType
	Directives    []Directive
	SelectionSet  SelectionSet
}

// Values

type Value interface {
	//GetKind() string
	//GetLoc() Location
	//GetName() Name
}

type IntValue struct {
	Kind  string
	Loc   Location
	Value string
}

type FloatValue struct {
	Kind  string
	Loc   Location
	Value string
}

type StringValue struct {
	Kind  string
	Loc   Location
	Value string
}

type BooleanValue struct {
	Kind  string
	Loc   Location
	Value bool
}

type EnumValue struct {
	Kind  string
	Loc   Location
	Value string
}

type ListValue struct {
	Kind   string
	Loc    Location
	Values []Value
}

type ObjectValue struct {
	Kind   string
	Loc    Location
	Fields []ObjectField
}

type ObjectField struct {
	Kind  string
	Name  Name
	Loc   Location
	Value Value
}

// Type does not exists in graphql-js ast
type ArrayValue struct {
	Kind   string
	Loc    Location
	Values []Value
}

// Directives

// why?
//type Directive struct {
//	Kind  string
//	Loc   Location
//	Name  Name
//	Value Value
//}

type Directive struct {
	Kind      string
	Loc       Location
	Name      Name
	Arguments []Argument
}

func NewDirective() *Directive {
	return &Directive{
		Kind: "Directive",
	}
}

// Type Reference

type Type interface{}

type NamedType struct {
	Kind string
	Loc  Location
	Name Name
	Type Type
}

type ListType struct {
	Kind string
	Loc  Location
	Type Type
}

type NonNullType struct {
	Kind string
	Loc  Location
	Type Type
}

// TODO: Type Definitions

type TypeDefinition interface{}

type FieldDefinition struct {
	Kind      string
	Loc       Location
	Name      Name
	Arguments []InputValueDefinition
	Type      Type
}

type InputValueDefinition struct {
	Kind         string
	Loc          Location
	Name         Name
	Type         Type
	DefaultValue Value
}

type EnumValueDefinition struct {
	Kind string
	Loc  Location
	Name Name
}
