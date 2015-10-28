package graphql

import (
	"github.com/chris-ramon/graphql/kinds"
)

type Definition interface {
	// TODO: determine the minimal set of interface for `Definition`
	GetOperation() string
	GetVariableDefinitions() []*AstVariableDefinition
	GetSelectionSet() *AstSelectionSet
}

// Ensure that all definition types implements Definition interface
var _ Definition = (*AstOperationDefinition)(nil)
var _ Definition = (*AstFragmentDefinition)(nil)
var _ Definition = (Definition)(nil)

// AstOperationDefinition implements Node, Definition
type AstOperationDefinition struct {
	Kind                string
	Loc                 *AstLocation
	Operation           string
	Name                *AstName
	VariableDefinitions []*AstVariableDefinition
	Directives          []*AstDirective
	SelectionSet        *AstSelectionSet
}

func NewAstOperationDefinition(op *AstOperationDefinition) *AstOperationDefinition {
	if op == nil {
		op = &AstOperationDefinition{}
	}
	return &AstOperationDefinition{
		Kind:                kinds.OperationDefinition,
		Loc:                 op.Loc,
		Operation:           op.Operation,
		Name:                op.Name,
		VariableDefinitions: op.VariableDefinitions,
		Directives:          op.Directives,
		SelectionSet:        op.SelectionSet,
	}
}

func (op *AstOperationDefinition) GetKind() string {
	return op.Kind
}

func (op *AstOperationDefinition) GetLoc() *AstLocation {
	return op.Loc
}

func (op *AstOperationDefinition) GetOperation() string {
	return op.Operation
}

func (op *AstOperationDefinition) GetName() *AstName {
	return op.Name
}

func (op *AstOperationDefinition) GetVariableDefinitions() []*AstVariableDefinition {
	return op.VariableDefinitions
}

func (op *AstOperationDefinition) GetDirectives() []*AstDirective {
	return op.Directives
}

func (op *AstOperationDefinition) GetSelectionSet() *AstSelectionSet {
	return op.SelectionSet
}

// AstFragmentDefinition implements Node, Definition
type AstFragmentDefinition struct {
	Kind                string
	Loc                 *AstLocation
	Operation           string
	Name                *AstName
	VariableDefinitions []*AstVariableDefinition
	TypeCondition       *AstNamed
	Directives          []*AstDirective
	SelectionSet        *AstSelectionSet
}

func NewAstFragmentDefinition(fd *AstFragmentDefinition) *AstFragmentDefinition {
	if fd == nil {
		fd = &AstFragmentDefinition{}
	}
	return &AstFragmentDefinition{
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

func (fd *AstFragmentDefinition) GetKind() string {
	return fd.Kind
}

func (fd *AstFragmentDefinition) GetLoc() *AstLocation {
	return fd.Loc
}

func (fd *AstFragmentDefinition) GetOperation() string {
	return fd.Operation
}

func (fd *AstFragmentDefinition) GetName() *AstName {
	return fd.Name
}

func (fd *AstFragmentDefinition) GetVariableDefinitions() []*AstVariableDefinition {
	return fd.VariableDefinitions
}

func (fd *AstFragmentDefinition) GetSelectionSet() *AstSelectionSet {
	return fd.SelectionSet
}

// AstVariableDefinition implements Node
type AstVariableDefinition struct {
	Kind         string
	Loc          *AstLocation
	Variable     *AstVariable
	Type         AstType
	DefaultValue Value
}

func NewAstVariableDefinition(vd *AstVariableDefinition) *AstVariableDefinition {
	if vd == nil {
		vd = &AstVariableDefinition{}
	}
	return &AstVariableDefinition{
		Kind:         kinds.VariableDefinition,
		Loc:          vd.Loc,
		Variable:     vd.Variable,
		Type:         vd.Type,
		DefaultValue: vd.DefaultValue,
	}
}

func (vd *AstVariableDefinition) GetKind() string {
	return vd.Kind
}

func (vd *AstVariableDefinition) GetLoc() *AstLocation {
	return vd.Loc
}
