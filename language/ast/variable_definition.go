package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// VariableDefinition implements Node
type VariableDefinition struct {
	Kind         string
	Loc          Location
	Variable     *Variable
	Type         Type
	DefaultValue Value
}

func NewVariableDefinition(vdef *VariableDefinition) *VariableDefinition {
	if vdef == nil {
		vdef = &VariableDefinition{}
	}
	return &VariableDefinition{
		Kind:         kinds.VariableDefinition,
		Loc:          vdef.Loc,
		Variable:     vdef.Variable,
		Type:         vdef.Type,
		DefaultValue: vdef.DefaultValue,
	}
}

func (vdef *VariableDefinition) GetKind() string {
	return vdef.Kind
}

func (vdef *VariableDefinition) GetLoc() Location {
	return vdef.Loc
}
