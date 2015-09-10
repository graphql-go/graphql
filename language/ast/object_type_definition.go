package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// ObjectTypeDefinition implements Node, TypeDefinition
type ObjectTypeDefinition struct {
	Kind       string
	Loc        Location
	Name       *Name
	Interfaces []NamedType
	Fields     []interface{}
}

func NewObjectTypeDefinition(def *ObjectTypeDefinition) *ObjectTypeDefinition {
	if def == nil {
		def = &ObjectTypeDefinition{}
	}
	return &ObjectTypeDefinition{
		Kind:       kinds.ObjectTypeDefinition,
		Loc:        def.Loc,
		Name:       def.Name,
		Interfaces: def.Interfaces,
		Fields:     def.Fields,
	}
}

func (def *ObjectTypeDefinition) GetKind() string {
	return def.Kind
}

func (def *ObjectTypeDefinition) GetLoc() Location {
	return def.Loc
}

func (def *ObjectTypeDefinition) GetName() *Name {
	return def.Name
}

func (def *ObjectTypeDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}

func (def *ObjectTypeDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}

func (def *ObjectTypeDefinition) GetOperation() string {
	return ""
}
