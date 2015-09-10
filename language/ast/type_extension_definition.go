package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// TypeExtensionDefinition implements Node, Definition
type TypeExtensionDefinition struct {
	Kind       string
	Loc        Location
	Definition ObjectTypeDefinition
}

func NewTypeExtensionDefinition() *TypeExtensionDefinition {
	return &TypeExtensionDefinition{
		Kind: kinds.TypeExtensionDefinition,
	}
}

func (def *TypeExtensionDefinition) GetKind() string {
	return def.Kind
}

func (def *TypeExtensionDefinition) GetLoc() Location {
	return def.Loc
}

func (def *TypeExtensionDefinition) GetName() *Name {
	return NewName(nil)
}

func (def *TypeExtensionDefinition) GetTypeCondition() NamedType {
	return NamedType{}
}

func (def *TypeExtensionDefinition) GetVariableDefinitions() []VariableDefinition {
	return []VariableDefinition{}
}

func (def *TypeExtensionDefinition) GetDirectives() []Directive {
	return []Directive{}
}

func (def *TypeExtensionDefinition) GetSelectionSet() SelectionSet {
	return SelectionSet{}
}

func (def *TypeExtensionDefinition) GetOperation() string {
	return ""
}
