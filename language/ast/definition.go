package ast

type Definition interface {
	// TODO: determine the minimal set of interface for `Definition`
	GetOperation() string
	GetVariableDefinitions() []VariableDefinition
	GetSelectionSet() SelectionSet
}

// TypeDefinition implements Definition
type TypeDefinition interface{
	// TODO: determine the minimal set of interface for `TypeDefinition`
	GetOperation() string
	GetVariableDefinitions() []VariableDefinition
	GetSelectionSet() SelectionSet
}

// Ensure that all definition types implements Definition interface
var _ Definition = (*OperationDefinition)(nil)
var _ Definition = (*FragmentDefinition)(nil)
var _ Definition = (TypeDefinition)(nil)

// Ensure that all typeDefinition types implements TypeDefinition interface
var _ TypeDefinition = (*ObjectTypeDefinition)(nil)
var _ TypeDefinition = (*InterfaceTypeDefinition)(nil)
var _ TypeDefinition = (*UnionTypeDefinition)(nil)
var _ TypeDefinition = (*ScalarTypeDefinition)(nil)
var _ TypeDefinition = (*EnumTypeDefinition)(nil)
var _ TypeDefinition = (*InputObjectTypeDefinition)(nil)
var _ TypeDefinition = (*TypeExtensionDefinition)(nil)
