package fd

import (
	. "github.com/chris-ramon/graphql-go/language/ast"
	"github.com/chris-ramon/graphql-go/language/kinds"
)

type FragmentDefinition struct {
	Kind                string
	Loc                 Location
	Operation           string
	Name                Name
	VariableDefinitions []VariableDefinition
	TypeCondition       Name
	Directives          []Directive
	SelectionSet        SelectionSet
}

func NewFragmentDefinition() *FragmentDefinition {
	return &FragmentDefinition{
		Kind: kinds.FragmentDefinition,
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

func (fd *FragmentDefinition) GetName() Name {
	return fd.Name
}

func (fd *FragmentDefinition) GetTypeCondition() Name {
	return fd.TypeCondition
}

func (fd *FragmentDefinition) GetVariableDefinitions() []VariableDefinition {
	return fd.VariableDefinitions
}

func (fd *FragmentDefinition) GetDirectives() []Directive {
	return fd.Directives
}

func (fd *FragmentDefinition) GetSelectionSet() SelectionSet {
	return fd.SelectionSet
}
