package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// InputValueDefinition implements Node, TypeDefinition
type InputValueDefinition struct {
	Kind         string
	Loc          *Location
	Name         *Name
	Type         Type
	DefaultValue Value
}

func NewInputValueDefinition(def *InputValueDefinition) *InputValueDefinition {
	if def == nil {
		def = &InputValueDefinition{}
	}
	return &InputValueDefinition{
		Kind:         kinds.InputValueDefinition,
		Loc:          def.Loc,
		Name:         def.Name,
		Type:         def.Type,
		DefaultValue: def.DefaultValue,
	}
}

func (def *InputValueDefinition) GetKind() string {
	return def.Kind
}

func (def *InputValueDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *InputValueDefinition) GetOperation() string {
	return ""
}

func (def *InputValueDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}

func (def *InputValueDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}
