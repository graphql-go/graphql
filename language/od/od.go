package od

import (
	. "github.com/chris-ramon/graphql-go/language/ast"
)

type OperationDefinition struct {
	Kind                string
	Loc                 Location
	Operation           string
	Name                Name
	VariableDefinitions []VariableDefinition
	TypeCondition       Name
	Directives          []Directive
	SelectionSet        SelectionSet
}

func GetKind(op *OperationDefinition) string {
	return op.Kind
}

func GetLoc(op *OperationDefinition) Location {
	return op.Loc
}

func GetOperation(op *OperationDefinition) string {
	return op.Operation
}

func GetName(op *OperationDefinition) Name {
	return op.Name
}

func GetTypeCondition(op *OperationDefinition) Name {
	return op.TypeCondition
}

func GetVariableDefinitions(op *OperationDefinition) []VariableDefinition {
	return op.VariableDefinitions
}

func GetDirectives(op *OperationDefinition) []Directive {
	return op.Directives
}

func GetSelectionSet(op *OperationDefinition) SelectionSet {
	return op.SelectionSet
}

func NewOperationDefinition() *OperationDefinition {
	return &OperationDefinition{
		Kind: "OperationDefinition",
	}
}
