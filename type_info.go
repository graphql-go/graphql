package graphql

import (
	"github.com/sprucehealth/graphql/language/ast"
	"github.com/sprucehealth/graphql/language/kinds"
)

// TODO: can move TypeInfo to a utils package if there ever is one
/**
 * TypeInfo is a utility class which, given a GraphQL schema, can keep track
 * of the current field and type definitions at any point in a GraphQL document
 * AST during a recursive descent by calling `enter(node)` and `leave(node)`.
 */
type TypeInfo struct {
	schema          *Schema
	typeStack       []Output
	parentTypeStack []Composite
	inputTypeStack  []Input
	fieldDefStack   []*FieldDefinition
	directive       *Directive
	argument        *Argument
}

func NewTypeInfo(schema *Schema) *TypeInfo {
	return &TypeInfo{
		schema: schema,
	}
}

func (ti *TypeInfo) Type() Output {
	if len(ti.typeStack) > 0 {
		return ti.typeStack[len(ti.typeStack)-1]
	}
	return nil
}

func (ti *TypeInfo) ParentType() Composite {
	if len(ti.parentTypeStack) > 0 {
		return ti.parentTypeStack[len(ti.parentTypeStack)-1]
	}
	return nil
}

func (ti *TypeInfo) InputType() Input {
	if len(ti.inputTypeStack) > 0 {
		return ti.inputTypeStack[len(ti.inputTypeStack)-1]
	}
	return nil
}
func (ti *TypeInfo) FieldDef() *FieldDefinition {
	if len(ti.fieldDefStack) > 0 {
		return ti.fieldDefStack[len(ti.fieldDefStack)-1]
	}
	return nil
}

func (ti *TypeInfo) Directive() *Directive {
	return ti.directive
}

func (ti *TypeInfo) Argument() *Argument {
	return ti.argument
}

func (ti *TypeInfo) Enter(node ast.Node) {
	schema := ti.schema
	var ttype Type
	switch node := node.(type) {
	case *ast.SelectionSet:
		namedType := GetNamed(ti.Type())
		var compositeType Composite = nil
		if IsCompositeType(namedType) {
			compositeType, _ = namedType.(Composite)
		}
		ti.parentTypeStack = append(ti.parentTypeStack, compositeType)
	case *ast.Field:
		parentType := ti.ParentType()
		var fieldDef *FieldDefinition
		if parentType != nil {
			fieldDef = TypeInfoFieldDef(*schema, parentType.(Type), node)
		}
		ti.fieldDefStack = append(ti.fieldDefStack, fieldDef)
		if fieldDef != nil {
			ti.typeStack = append(ti.typeStack, fieldDef.Type)
		} else {
			ti.typeStack = append(ti.typeStack, nil)
		}
	case *ast.Directive:
		nameVal := ""
		if node.Name != nil {
			nameVal = node.Name.Value
		}
		ti.directive = schema.Directive(nameVal)
	case *ast.OperationDefinition:
		if node.Operation == "query" {
			ttype = schema.QueryType()
		} else if node.Operation == "mutation" {
			ttype = schema.MutationType()
		}
		ti.typeStack = append(ti.typeStack, ttype)
	case *ast.InlineFragment:
		ttype, _ = typeFromAST(*schema, node.TypeCondition)
		ti.typeStack = append(ti.typeStack, ttype)
	case *ast.FragmentDefinition:
		ttype, _ = typeFromAST(*schema, node.TypeCondition)
		ti.typeStack = append(ti.typeStack, ttype)
	case *ast.VariableDefinition:
		ttype, _ = typeFromAST(*schema, node.Type)
		ti.inputTypeStack = append(ti.inputTypeStack, ttype)
	case *ast.Argument:
		nameVal := ""
		if node.Name != nil {
			nameVal = node.Name.Value
		}
		var argType Input
		var argDef *Argument
		directive := ti.Directive()
		fieldDef := ti.FieldDef()
		if directive != nil {
			for _, arg := range directive.Args {
				if arg.Name() == nameVal {
					argDef = arg
				}
			}
		} else if fieldDef != nil {
			for _, arg := range fieldDef.Args {
				if arg.Name() == nameVal {
					argDef = arg
				}
			}
		}
		if argDef != nil {
			argType = argDef.Type
		}
		ti.argument = argDef
		ti.inputTypeStack = append(ti.inputTypeStack, argType)
	case *ast.ListValue:
		listType := GetNullable(ti.InputType())
		if list, ok := listType.(*List); ok {
			ti.inputTypeStack = append(ti.inputTypeStack, list.OfType)
		} else {
			ti.inputTypeStack = append(ti.inputTypeStack, nil)
		}
	case *ast.ObjectField:
		var fieldType Input
		objectType := GetNamed(ti.InputType())

		if objectType, ok := objectType.(*InputObject); ok {
			nameVal := ""
			if node.Name != nil {
				nameVal = node.Name.Value
			}
			if inputField, ok := objectType.Fields()[nameVal]; ok {
				fieldType = inputField.Type
			}
		}
		ti.inputTypeStack = append(ti.inputTypeStack, fieldType)
	}
}
func (ti *TypeInfo) Leave(node ast.Node) {
	kind := node.GetKind()
	switch kind {
	case kinds.SelectionSet:
		// pop ti.parentTypeStack
		_, ti.parentTypeStack = ti.parentTypeStack[len(ti.parentTypeStack)-1], ti.parentTypeStack[:len(ti.parentTypeStack)-1]
	case kinds.Field:
		// pop ti.fieldDefStack
		_, ti.fieldDefStack = ti.fieldDefStack[len(ti.fieldDefStack)-1], ti.fieldDefStack[:len(ti.fieldDefStack)-1]
		// pop ti.typeStack
		_, ti.typeStack = ti.typeStack[len(ti.typeStack)-1], ti.typeStack[:len(ti.typeStack)-1]
	case kinds.Directive:
		ti.directive = nil
	case kinds.OperationDefinition:
		fallthrough
	case kinds.InlineFragment:
		fallthrough
	case kinds.FragmentDefinition:
		// pop ti.typeStack
		_, ti.typeStack = ti.typeStack[len(ti.typeStack)-1], ti.typeStack[:len(ti.typeStack)-1]
	case kinds.VariableDefinition:
		// pop ti.inputTypeStack
		_, ti.inputTypeStack = ti.inputTypeStack[len(ti.inputTypeStack)-1], ti.inputTypeStack[:len(ti.inputTypeStack)-1]
	case kinds.Argument:
		ti.argument = nil
		// pop ti.inputTypeStack
		_, ti.inputTypeStack = ti.inputTypeStack[len(ti.inputTypeStack)-1], ti.inputTypeStack[:len(ti.inputTypeStack)-1]
	case kinds.ListValue:
		fallthrough
	case kinds.ObjectField:
		// pop ti.inputTypeStack
		_, ti.inputTypeStack = ti.inputTypeStack[len(ti.inputTypeStack)-1], ti.inputTypeStack[:len(ti.inputTypeStack)-1]
	}
}

/**
 * Not exactly the same as the executor's definition of FieldDef, in this
 * statically evaluated environment we do not always have an Object type,
 * and need to handle Interface and Union types.
 */
func TypeInfoFieldDef(schema Schema, parentType Type, fieldAST *ast.Field) *FieldDefinition {
	name := ""
	if fieldAST.Name != nil {
		name = fieldAST.Name.Value
	}
	if name == SchemaMetaFieldDef.Name &&
		schema.QueryType() == parentType {
		return SchemaMetaFieldDef
	}
	if name == TypeMetaFieldDef.Name &&
		schema.QueryType() == parentType {
		return TypeMetaFieldDef
	}
	if name == TypeNameMetaFieldDef.Name {
		if _, ok := parentType.(*Object); ok && parentType != nil {
			return TypeNameMetaFieldDef
		}
		if _, ok := parentType.(*Interface); ok && parentType != nil {
			return TypeNameMetaFieldDef
		}
		if _, ok := parentType.(*Union); ok && parentType != nil {
			return TypeNameMetaFieldDef
		}
	}

	if parentType, ok := parentType.(*Object); ok && parentType != nil {
		field, _ := parentType.Fields()[name]
		return field
	}
	if parentType, ok := parentType.(*Interface); ok && parentType != nil {
		field, _ := parentType.Fields()[name]
		return field
	}
	return nil
}
