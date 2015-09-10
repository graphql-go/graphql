package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// EnumValueDefinition implements Node, TypeDefinition
type EnumValueDefinition struct {
	Kind string
	Loc  *Location
	Name *Name
}

func NewEnumValueDefinition(def *EnumValueDefinition) *EnumValueDefinition {
	if def == nil {
		def = &EnumValueDefinition{}
	}
	return &EnumValueDefinition{
		Kind: kinds.EnumValueDefinition,
		Loc:  def.Loc,
		Name: def.Name,
	}
}

func (def *EnumValueDefinition) GetKind() string {
	return def.Kind
}

func (def *EnumValueDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *EnumValueDefinition) GetOperation() string {
	return ""
}

func (def *EnumValueDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}

func (def *EnumValueDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}
