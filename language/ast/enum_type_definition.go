package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// EnumTypeDefinition implements Node, Definition
type EnumTypeDefinition struct {
	Kind   string
	Loc    Location
	Name   *Name
	Values []interface{}
}

func NewEnumTypeDefinition() *EnumTypeDefinition {
	return &EnumTypeDefinition{
		Kind: kinds.EnumTypeDefinition,
	}
}

func (def *EnumTypeDefinition) GetKind() string {
	return def.Kind
}

func (def *EnumTypeDefinition) GetLoc() Location {
	return def.Loc
}

func (def *EnumTypeDefinition) GetName() *Name {
	return def.Name
}

func (def *EnumTypeDefinition) GetTypeCondition() NamedType {
	return NamedType{}
}

func (def *EnumTypeDefinition) GetVariableDefinitions() []VariableDefinition {
	return []VariableDefinition{}
}

func (def *EnumTypeDefinition) GetDirectives() []Directive {
	return []Directive{}
}

func (def *EnumTypeDefinition) GetSelectionSet() SelectionSet {
	return SelectionSet{}
}

func (def *EnumTypeDefinition) GetOperation() string {
	return ""
}
