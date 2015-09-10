package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// FieldDefinition implements Node, TypeDefinition
type FieldDefinition struct {
	Kind      string
	Loc       *Location
	Name      *Name
	Arguments []*InputValueDefinition
	Type      Type
}

func NewFieldDefinition(def *FieldDefinition) *FieldDefinition {
	if def == nil {
		def = &FieldDefinition{}
	}
	return &FieldDefinition{
		Kind:      kinds.FieldDefinition,
		Loc:       def.Loc,
		Name:      def.Name,
		Arguments: def.Arguments,
		Type:      def.Type,
	}
}

func (def *FieldDefinition) GetKind() string {
	return def.Kind
}

func (def *FieldDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *FieldDefinition) GetOperation() string {
	return ""
}

func (def *FieldDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}

func (def *FieldDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}
