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

func GetKind(fd *FragmentDefinition) string {
	return fd.Kind
}

func GetLoc(fd *FragmentDefinition) Location {
	return fd.Loc
}

func GetOperation(fd *FragmentDefinition) string {
	return fd.Operation
}

func GetName(fd *FragmentDefinition) Name {
	return fd.Name
}

func GetTypeCondition(fd *FragmentDefinition) Name {
	return fd.TypeCondition
}

func GetVariableDefinitions(fd *FragmentDefinition) []VariableDefinition {
	return fd.VariableDefinitions
}

func GetDirectives(fd *FragmentDefinition) []Directive {
	return fd.Directives
}

func GetSelectionSet(fd *FragmentDefinition) SelectionSet {
	return fd.SelectionSet
}
