package ast

type Definition interface {
	// TODO: determine the minimal set of interface for `Definition`
	GetOperation() string
	GetVariableDefinitions() []*VariableDefinition
	GetSelectionSet() *SelectionSet
}

// Ensure that all definition types implements Definition interface
var _ Definition = (*OperationDefinition)(nil)
var _ Definition = (*FragmentDefinition)(nil)
var _ Definition = (TypeDefinition)(nil)
