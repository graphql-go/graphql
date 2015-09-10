package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// OperationDefinition implements Node, Definition
type OperationDefinition struct {
	Kind                string
	Loc                 Location
	Operation           string
	Name                *Name
	VariableDefinitions []*VariableDefinition
	Directives          []Directive
	SelectionSet        SelectionSet
}

func NewOperationDefinition(op *OperationDefinition) *OperationDefinition {
	if op == nil {
		op = &OperationDefinition{}
	}
	return &OperationDefinition{
		Kind:                kinds.OperationDefinition,
		Loc:                 op.Loc,
		Operation:           op.Operation,
		Name:                op.Name,
		VariableDefinitions: op.VariableDefinitions,
		Directives:          op.Directives,
		SelectionSet:        op.SelectionSet,
	}
}

func (op *OperationDefinition) GetKind() string {
	return op.Kind
}

func (op *OperationDefinition) GetLoc() Location {
	return op.Loc
}

func (op *OperationDefinition) GetOperation() string {
	return op.Operation
}

func (op *OperationDefinition) GetName() *Name {
	return op.Name
}

func (op *OperationDefinition) GetVariableDefinitions() []*VariableDefinition {
	return op.VariableDefinitions
}

func (op *OperationDefinition) GetDirectives() []Directive {
	return op.Directives
}

func (op *OperationDefinition) GetSelectionSet() SelectionSet {
	return op.SelectionSet
}
