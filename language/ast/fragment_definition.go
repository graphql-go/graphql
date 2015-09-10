package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// FragmentDefinition implements Node, Definition
type FragmentDefinition struct {
	Kind                string
	Loc                 Location
	Operation           string
	Name                *Name
	VariableDefinitions []*VariableDefinition
	TypeCondition       NamedType
	Directives          []Directive
	SelectionSet        SelectionSet
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

func (fd *FragmentDefinition) GetLoc() Location {
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

func (fd *FragmentDefinition) GetSelectionSet() SelectionSet {
	return fd.SelectionSet
}
