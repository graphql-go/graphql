package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// InterfaceTypeDefinition implements Definition
type InterfaceTypeDefinition struct {
	Kind   string
	Loc    Location
	Name   Name
	Fields []interface{}
}

func NewInterfaceTypeDefinition() *InterfaceTypeDefinition {
	return &InterfaceTypeDefinition{
		Kind: kinds.InterfaceTypeDefinition,
	}
}

func (def *InterfaceTypeDefinition) GetKind() string {
	return def.Kind
}

func (def *InterfaceTypeDefinition) GetLoc() Location {
	return def.Loc
}

func (def *InterfaceTypeDefinition) GetName() Name {
	return def.Name
}

func (def *InterfaceTypeDefinition) GetTypeCondition() NamedType {
	return NamedType{}
}

func (def *InterfaceTypeDefinition) GetVariableDefinitions() []VariableDefinition {
	return []VariableDefinition{}
}

func (def *InterfaceTypeDefinition) GetDirectives() []Directive {
	return []Directive{}
}

func (def *InterfaceTypeDefinition) GetSelectionSet() SelectionSet {
	return SelectionSet{}
}

func (def *InterfaceTypeDefinition) GetOperation() string {
	return ""
}

