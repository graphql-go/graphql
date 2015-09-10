package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// ScalarTypeDefinition implements Node, TypeDefinition
type ScalarTypeDefinition struct {
	Kind string
	Loc  Location
	Name *Name
}

func NewScalarTypeDefinition(def *ScalarTypeDefinition) *ScalarTypeDefinition {
	if def == nil {
		def = &ScalarTypeDefinition{}
	}
	return &ScalarTypeDefinition{
		Kind: kinds.ScalarTypeDefinition,
		Loc:  def.Loc,
		Name: def.Name,
	}
}

func (def *ScalarTypeDefinition) GetKind() string {
	return def.Kind
}

func (def *ScalarTypeDefinition) GetLoc() Location {
	return def.Loc
}

func (def *ScalarTypeDefinition) GetName() *Name {
	return def.Name
}

func (def *ScalarTypeDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}

func (def *ScalarTypeDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}

func (def *ScalarTypeDefinition) GetOperation() string {
	return ""
}
