package ast

import (
	"github.com/graphql-go/graphql/language/kinds"
)

type Definition interface {
	// TODO: determine the minimal set of interface for `Definition`
	GetOperation() string
	GetVariableDefinitions() []*VariableDefinition
	GetSelectionSet() *SelectionSet
}

// Ensure that all definition types implements Definition interface
var _ Definition = (*OperationDefinition)(nil)
var _ Definition = (*FragmentDefinition)(nil)
var _ Definition = (Definition)(nil)

// OperationDefinition implements Node, Definition
type OperationDefinition struct {
	Kind                string
	Loc                 *Location
	Operation           string
	Name                *Name
	VariableDefinitions []*VariableDefinition
	Directives          []*Directive
	SelectionSet        *SelectionSet
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

func (op *OperationDefinition) GetLoc() *Location {
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

func (op *OperationDefinition) GetDirectives() []*Directive {
	return op.Directives
}

func (op *OperationDefinition) GetSelectionSet() *SelectionSet {
	return op.SelectionSet
}

// FragmentDefinition implements Node, Definition
type FragmentDefinition struct {
	Kind                string
	Loc                 *Location
	Operation           string
	Name                *Name
	VariableDefinitions []*VariableDefinition
	TypeCondition       *Named
	Directives          []*Directive
	SelectionSet        *SelectionSet
}

func NewFragmentDefinition(fd *FragmentDefinition) *FragmentDefinition {
	if fd == nil {
		fd = &FragmentDefinition{}
	}
	return &FragmentDefinition{
		Kind:                kinds.FragmentDefinition,
		Loc:                 fd.Loc,
		Operation:           fd.Operation,
		Name:                fd.Name,
		VariableDefinitions: fd.VariableDefinitions,
		TypeCondition:       fd.TypeCondition,
		Directives:          fd.Directives,
		SelectionSet:        fd.SelectionSet,
	}
}

func (fd *FragmentDefinition) GetKind() string {
	return fd.Kind
}

func (fd *FragmentDefinition) GetLoc() *Location {
	return fd.Loc
}

func (fd *FragmentDefinition) GetOperation() string {
	return fd.Operation
}

func (fd *FragmentDefinition) GetName() *Name {
	return fd.Name
}

func (fd *FragmentDefinition) GetVariableDefinitions() []*VariableDefinition {
	return fd.VariableDefinitions
}

func (fd *FragmentDefinition) GetSelectionSet() *SelectionSet {
	return fd.SelectionSet
}

// VariableDefinition implements Node
type VariableDefinition struct {
	Kind         string
	Loc          *Location
	Variable     *Variable
	Type         Type
	DefaultValue Value
}

func NewVariableDefinition(vd *VariableDefinition) *VariableDefinition {
	if vd == nil {
		vd = &VariableDefinition{}
	}
	return &VariableDefinition{
		Kind:         kinds.VariableDefinition,
		Loc:          vd.Loc,
		Variable:     vd.Variable,
		Type:         vd.Type,
		DefaultValue: vd.DefaultValue,
	}
}

func (vd *VariableDefinition) GetKind() string {
	return vd.Kind
}

func (vd *VariableDefinition) GetLoc() *Location {
	return vd.Loc
}
