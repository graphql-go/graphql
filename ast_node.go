package graphql

type Node interface {
	GetKind() string
	GetLoc() *AstLocation
}

// The list of all possible AST node graphql.
// Ensure that all node types implements Node interface
var _ Node = (*AstName)(nil)
var _ Node = (*AstDocument)(nil)
var _ Node = (*AstOperationDefinition)(nil)
var _ Node = (*AstVariableDefinition)(nil)
var _ Node = (*AstVariable)(nil)
var _ Node = (*AstSelectionSet)(nil)
var _ Node = (*AstField)(nil)
var _ Node = (*AstArgument)(nil)
var _ Node = (*AstFragmentSpread)(nil)
var _ Node = (*AstInlineFragment)(nil)
var _ Node = (*AstFragmentDefinition)(nil)
var _ Node = (*AstIntValue)(nil)
var _ Node = (*AstFloatValue)(nil)
var _ Node = (*AstStringValue)(nil)
var _ Node = (*AstBooleanValue)(nil)
var _ Node = (*AstEnumValue)(nil)
var _ Node = (*AstListValue)(nil)
var _ Node = (*AstObjectValue)(nil)
var _ Node = (*AstObjectField)(nil)
var _ Node = (*AstDirective)(nil)
var _ Node = (*AstList)(nil)
var _ Node = (*AstNonNull)(nil)
var _ Node = (*AstObjectDefinition)(nil)
var _ Node = (*AstFieldDefinition)(nil)
var _ Node = (*AstInputValueDefinition)(nil)
var _ Node = (*AstInterfaceDefinition)(nil)
var _ Node = (*AstUnionDefinition)(nil)
var _ Node = (*AstScalarDefinition)(nil)
var _ Node = (*AstEnumDefinition)(nil)
var _ Node = (*AstEnumValueDefinition)(nil)
var _ Node = (*AstInputObjectDefinition)(nil)
var _ Node = (*AstTypeExtensionDefinition)(nil)

// TODO: File issue in `graphql-js` where AstNamed is not
// defined as a Node. This might be a mistake in `graphql-js`?
var _ Node = (*AstNamed)(nil)
