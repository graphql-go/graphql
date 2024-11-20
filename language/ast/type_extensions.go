package ast

import "github.com/graphql-go/graphql/language/kinds"

type TypeExtension interface {
	GetKind() string
	GetLoc() *Location
}

var _ TypeExtension = (*ScalarDefinition)(nil)
var _ TypeExtension = (*ObjectDefinition)(nil)
var _ TypeExtension = (*InterfaceDefinition)(nil)
var _ TypeExtension = (*UnionDefinition)(nil)
var _ TypeExtension = (*EnumDefinition)(nil)
var _ TypeExtension = (*InputObjectDefinition)(nil)

// TypeExtensionDefinition implements Node, Definition
type TypeExtensionDefinition struct {
	Kind       string
	Loc        *Location
	Definition TypeExtension
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

func (def *TypeExtensionDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *TypeExtensionDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}

func (def *TypeExtensionDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}

func (def *TypeExtensionDefinition) GetOperation() string {
	return ""
}
