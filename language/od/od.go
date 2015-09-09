package od

import (
	. "github.com/chris-ramon/graphql-go/language/ast"
	"github.com/chris-ramon/graphql-go/language/kinds"
)

type OperationDefinition struct {
	Kind                string
	Loc                 Location
	Operation           string
	Name                Name
	VariableDefinitions []VariableDefinition
	TypeCondition       NamedType
	Directives          []Directive
	SelectionSet        SelectionSet
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

func (op *OperationDefinition) GetName() Name {
	return op.Name
}

func (op *OperationDefinition) GetTypeCondition() NamedType {
	return op.TypeCondition
}

func (op *OperationDefinition) GetVariableDefinitions() []VariableDefinition {
	return op.VariableDefinitions
}

func (op *OperationDefinition) GetDirectives() []Directive {
	return op.Directives
}

func (op *OperationDefinition) GetSelectionSet() SelectionSet {
	return op.SelectionSet
}

func NewOperationDefinition() *OperationDefinition {
	return &OperationDefinition{
		Kind: kinds.OperationDefinition,
	}
}
