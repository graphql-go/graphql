package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// UnionTypeDefinition implements Node, TypeDefinition
type UnionTypeDefinition struct {
	Kind  string
	Loc   Location
	Name  *Name
	Types []NamedType
}

func NewUnionTypeDefinition(def *UnionTypeDefinition) *UnionTypeDefinition {
	if def == nil {
		def = &UnionTypeDefinition{}
	}
	return &UnionTypeDefinition{
		Kind:  kinds.UnionTypeDefinition,
		Loc:   def.Loc,
		Name:  def.Name,
		Types: def.Types,
	}
}

func (def *UnionTypeDefinition) GetKind() string {
	return def.Kind
}

func (def *UnionTypeDefinition) GetLoc() Location {
	return def.Loc
}

func (def *UnionTypeDefinition) GetName() *Name {
	return def.Name
}

func (def *UnionTypeDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}

func (def *UnionTypeDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}

func (def *UnionTypeDefinition) GetOperation() string {
	return ""
}
