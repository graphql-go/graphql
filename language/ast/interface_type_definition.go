package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// InterfaceTypeDefinition implements Node, TypeDefinition
type InterfaceTypeDefinition struct {
	Kind   string
	Loc    *Location
	Name   *Name
	Fields []*FieldDefinition
}

func NewInterfaceTypeDefinition(def *InterfaceTypeDefinition) *InterfaceTypeDefinition {
	if def == nil {
		def = &InterfaceTypeDefinition{}
	}
	return &InterfaceTypeDefinition{
		Kind:   kinds.InterfaceTypeDefinition,
		Loc:    def.Loc,
		Name:   def.Name,
		Fields: def.Fields,
	}
}

func (def *InterfaceTypeDefinition) GetKind() string {
	return def.Kind
}

func (def *InterfaceTypeDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *InterfaceTypeDefinition) GetName() *Name {
	return def.Name
}

func (def *InterfaceTypeDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}

func (def *InterfaceTypeDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}

func (def *InterfaceTypeDefinition) GetOperation() string {
	return ""
}
