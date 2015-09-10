package ast

// Interfaces
type Node interface {
	GetKind() string
	GetLoc() Location
}

// The list of all possible AST node types.
// Ensure that all node types implements Node interface
var _ Node = (*Name)(nil)
var _ Node = (*Document)(nil)
var _ Node = (*OperationDefinition)(nil)
var _ Node = (*VariableDefinition)(nil)
var _ Node = (*Variable)(nil)
var _ Node = (*SelectionSet)(nil)
var _ Node = (*Field)(nil)

//var _ Node = (*Argument)(nil)
var _ Node = (*FragmentSpread)(nil)
var _ Node = (*InlineFragment)(nil)
var _ Node = (*FragmentDefinition)(nil)
var _ Node = (*IntValue)(nil)
var _ Node = (*FloatValue)(nil)
var _ Node = (*StringValue)(nil)
var _ Node = (*BooleanValue)(nil)
var _ Node = (*EnumValue)(nil)
var _ Node = (*ListValue)(nil)
var _ Node = (*ObjectValue)(nil)
var _ Node = (*ObjectField)(nil)

//var _ Node = (*Directive)(nil)
//var _ Node = (*ListType)(nil)
//var _ Node = (*NonNullType)(nil)
var _ Node = (*ObjectTypeDefinition)(nil)

//var _ Node = (*FieldDefinition)(nil)
//var _ Node = (*InputValueDefinition)(nil)
var _ Node = (*InterfaceTypeDefinition)(nil)
var _ Node = (*UnionTypeDefinition)(nil)
var _ Node = (*ScalarTypeDefinition)(nil)
var _ Node = (*EnumTypeDefinition)(nil)

//var _ Node = (*EnumValueDefinition)(nil)
var _ Node = (*InputObjectTypeDefinition)(nil)
var _ Node = (*TypeExtensionDefinition)(nil)
