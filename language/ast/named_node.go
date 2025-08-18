package ast

type NamedNode interface {
	GetName() *Name
}

// Ensure that the following types implement the NamedNode interface
var _ NamedNode = (*Argument)(nil)
var _ NamedNode = (*DirectiveDefinition)(nil)
var _ NamedNode = (*FragmentDefinition)(nil)
var _ NamedNode = (*OperationDefinition)(nil)
var _ NamedNode = (*Field)(nil)
var _ NamedNode = (*FragmentSpread)(nil)
var _ NamedNode = (*EnumDefinition)(nil)
var _ NamedNode = (*EnumValueDefinition)(nil)
var _ NamedNode = (*FieldDefinition)(nil)
var _ NamedNode = (*InputObjectDefinition)(nil)
var _ NamedNode = (*InputValueDefinition)(nil)
var _ NamedNode = (*InterfaceDefinition)(nil)
var _ NamedNode = (*ObjectDefinition)(nil)
var _ NamedNode = (*ScalarDefinition)(nil)
var _ NamedNode = (*UnionDefinition)(nil)
var _ NamedNode = (*ObjectField)(nil)
var _ NamedNode = (*Variable)(nil)
