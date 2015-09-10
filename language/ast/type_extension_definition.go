package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// TypeExtensionDefinition implements Node, TypeDefinition
type TypeExtensionDefinition struct {
	Kind       string
	Loc        Location
	Definition *ObjectTypeDefinition
}

func NewTypeExtensionDefinition(def *TypeExtensionDefinition) *TypeExtensionDefinition {
	if def == nil {
		def = &TypeExtensionDefinition{}
	}
	return &TypeExtensionDefinition{
		Kind:       kinds.TypeExtensionDefinition,
		Loc:        def.Loc,
		Definition: def.Definition,
	}
}

func (def *TypeExtensionDefinition) GetKind() string {
	return def.Kind
}

func (def *TypeExtensionDefinition) GetLoc() Location {
	return def.Loc
}

func (def *TypeExtensionDefinition) GetVariableDefinitions() []VariableDefinition {
	return []VariableDefinition{}
}

func (def *TypeExtensionDefinition) GetSelectionSet() SelectionSet {
	return SelectionSet{}
}

func (def *TypeExtensionDefinition) GetOperation() string {
	return ""
}
