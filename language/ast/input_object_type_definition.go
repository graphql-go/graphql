package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// InputObjectTypeDefinition implements Node, Definition
type InputObjectTypeDefinition struct {
	Kind   string
	Loc    Location
	Name   *Name
	Fields []interface{}
}

func NewInputObjectTypeDefinition(def *InputObjectTypeDefinition) *InputObjectTypeDefinition {
	if def == nil {
		def = &InputObjectTypeDefinition{}
	}
	return &InputObjectTypeDefinition{
		Kind:   kinds.InputObjectTypeDefinition,
		Loc:    def.Loc,
		Name:   def.Name,
		Fields: def.Fields,
	}
}

func (def *InputObjectTypeDefinition) GetKind() string {
	return def.Kind
}

func (def *InputObjectTypeDefinition) GetLoc() Location {
	return def.Loc
}

func (def *InputObjectTypeDefinition) GetName() *Name {
	return def.Name
}
func (def *InputObjectTypeDefinition) GetTypeCondition() NamedType {
	return NamedType{}
}

func (def *InputObjectTypeDefinition) GetVariableDefinitions() []VariableDefinition {
	return []VariableDefinition{}
}

func (def *InputObjectTypeDefinition) GetDirectives() []Directive {
	return []Directive{}
}

func (def *InputObjectTypeDefinition) GetSelectionSet() SelectionSet {
	return SelectionSet{}
}

func (def *InputObjectTypeDefinition) GetOperation() string {
	return ""
}
