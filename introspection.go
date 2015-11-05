package graphql

import (
	"fmt"
	"reflect"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/printer"
)

const (
	TypeKindScalar      = "SCALAR"
	TypeKindObject      = "OBJECT"
	TypeKindInterface   = "INTERFACE"
	TypeKindUnion       = "UNION"
	TypeKindEnum        = "ENUM"
	TypeKindInputObject = "INPUT_OBJECT"
	TypeKindList        = "LIST"
	TypeKindNonNull     = "NON_NULL"
)

var __Directive *Object
var __Schema *Object
var __Type *Object
var __Field *Object
var __InputValue *Object
var __EnumValue *Object

var __TypeKind *Enum

var SchemaMetaFieldDef *FieldDefinition
var TypeMetaFieldDef *FieldDefinition
var TypeNameMetaFieldDef *FieldDefinition

func init() {

	__TypeKind = NewEnum(EnumConfig{
		Name:        "__TypeKind",
		Description: "An enum describing what kind of type a given __Type is",
		Values: EnumValueConfigMap{
			"SCALAR": &EnumValueConfig{
				Value:       TypeKindScalar,
				Description: "Indicates this type is a scalar.",
			},
			"OBJECT": &EnumValueConfig{
				Value: TypeKindObject,
				Description: "Indicates this type is an object. " +
					"`fields` and `interfaces` are valid fields.",
			},
			"INTERFACE": &EnumValueConfig{
				Value: TypeKindInterface,
				Description: "Indicates this type is an interface. " +
					"`fields` and `possibleTypes` are valid fields.",
			},
			"UNION": &EnumValueConfig{
				Value: TypeKindUnion,
				Description: "Indicates this type is a union. " +
					"`possibleTypes` is a valid field.",
			},
			"ENUM": &EnumValueConfig{
				Value: TypeKindEnum,
				Description: "Indicates this type is an enum. " +
					"`enumValues` is a valid field.",
			},
			"INPUT_OBJECT": &EnumValueConfig{
				Value: TypeKindInputObject,
				Description: "Indicates this type is an input object. " +
					"`inputFields` is a valid field.",
			},
			"LIST": &EnumValueConfig{
				Value: TypeKindList,
				Description: "Indicates this type is a list. " +
					"`ofType` is a valid field.",
			},
			"NON_NULL": &EnumValueConfig{
				Value: TypeKindNonNull,
				Description: "Indicates this type is a non-null. " +
					"`ofType` is a valid field.",
			},
		},
	})

	// Note: some fields (for e.g "fields", "interfaces") are defined later due to cyclic reference
	__Type = NewObject(ObjectConfig{
		Name: "__Type",
		Fields: FieldConfigMap{
			"kind": &FieldConfig{
				Type: NewNonNull(__TypeKind),
				Resolve: func(p GQLFRParams) interface{} {
					switch p.Source.(type) {
					case *Scalar:
						return TypeKindScalar
					case *Object:
						return TypeKindObject
					case *Interface:
						return TypeKindInterface
					case *Union:
						return TypeKindUnion
					case *Enum:
						return TypeKindEnum
					case *InputObject:
						return TypeKindInputObject
					case *List:
						return TypeKindList
					case *NonNull:
						return TypeKindNonNull
					}
					panic(fmt.Sprintf("Unknown kind of type: %v", p.Source))
				},
			},
			"name": &FieldConfig{
				Type: String,
			},
			"description": &FieldConfig{
				Type: String,
			},
			"fields":        &FieldConfig{},
			"interfaces":    &FieldConfig{},
			"possibleTypes": &FieldConfig{},
			"enumValues":    &FieldConfig{},
			"inputFields":   &FieldConfig{},
			"ofType":        &FieldConfig{},
		},
	})

	__InputValue = NewObject(ObjectConfig{
		Name: "__InputValue",
		Fields: FieldConfigMap{
			"name": &FieldConfig{
				Type: NewNonNull(String),
			},
			"description": &FieldConfig{
				Type: String,
			},
			"type": &FieldConfig{
				Type: NewNonNull(__Type),
			},
			"defaultValue": &FieldConfig{
				Type: String,
				Resolve: func(p GQLFRParams) interface{} {
					if inputVal, ok := p.Source.(*Argument); ok {
						if inputVal.DefaultValue == nil {
							return nil
						}
						astVal := astFromValue(inputVal.DefaultValue, inputVal)
						return printer.Print(astVal)
					}
					if inputVal, ok := p.Source.(*InputObjectField); ok {
						if inputVal.DefaultValue == nil {
							return nil
						}
						astVal := astFromValue(inputVal.DefaultValue, inputVal)
						return printer.Print(astVal)
					}
					return nil
				},
			},
		},
	})

	__Field = NewObject(ObjectConfig{
		Name: "__Field",
		Fields: FieldConfigMap{
			"name": &FieldConfig{
				Type: NewNonNull(String),
			},
			"description": &FieldConfig{
				Type: String,
			},
			"args": &FieldConfig{
				Type: NewNonNull(NewList(NewNonNull(__InputValue))),
				Resolve: func(p GQLFRParams) interface{} {
					if field, ok := p.Source.(*FieldDefinition); ok {
						return field.Args
					}
					return []interface{}{}
				},
			},
			"type": &FieldConfig{
				Type: NewNonNull(__Type),
			},
			"isDeprecated": &FieldConfig{
				Type: NewNonNull(Boolean),
				Resolve: func(p GQLFRParams) interface{} {
					if field, ok := p.Source.(*FieldDefinition); ok {
						return (field.DeprecationReason != "")
					}
					return false
				},
			},
			"deprecationReason": &FieldConfig{
				Type: String,
			},
		},
	})

	__Directive = NewObject(ObjectConfig{
		Name: "__Directive",
		Fields: FieldConfigMap{
			"name": &FieldConfig{
				Type: NewNonNull(String),
			},
			"description": &FieldConfig{
				Type: String,
			},
			"args": &FieldConfig{
				Type: NewNonNull(NewList(
					NewNonNull(__InputValue),
				)),
			},
			"onOperation": &FieldConfig{
				Type: NewNonNull(Boolean),
			},
			"onFragment": &FieldConfig{
				Type: NewNonNull(Boolean),
			},
			"onField": &FieldConfig{
				Type: NewNonNull(Boolean),
			},
		},
	})

	__Schema = NewObject(ObjectConfig{
		Name: "__Schema",
		Description: `A GraphQL Schema defines the capabilities of a GraphQL
server. It exposes all available types and directives on
the server, as well as the entry points for query and
mutation operations.`,
		Fields: FieldConfigMap{
			"types": &FieldConfig{
				Description: "A list of all types supported by this server.",
				Type: NewNonNull(NewList(
					NewNonNull(__Type),
				)),
				Resolve: func(p GQLFRParams) interface{} {
					if schema, ok := p.Source.(Schema); ok {
						results := []Type{}
						for _, ttype := range schema.GetTypeMap() {
							results = append(results, ttype)
						}
						return results
					}
					return []Type{}
				},
			},
			"queryType": &FieldConfig{
				Description: "The type that query operations will be rooted at.",
				Type:        NewNonNull(__Type),
				Resolve: func(p GQLFRParams) interface{} {
					if schema, ok := p.Source.(Schema); ok {
						return schema.GetQueryType()
					}
					return nil
				},
			},
			"mutationType": &FieldConfig{
				Description: `If this server supports mutation, the type that ` +
					`mutation operations will be rooted at.`,
				Type: __Type,
				Resolve: func(p GQLFRParams) interface{} {
					if schema, ok := p.Source.(Schema); ok {
						if schema.GetMutationType() != nil {
							return schema.GetMutationType()
						}
					}
					return nil
				},
			},
			"directives": &FieldConfig{
				Description: `A list of all directives supported by this server.`,
				Type: NewNonNull(NewList(
					NewNonNull(__Directive),
				)),
				Resolve: func(p GQLFRParams) interface{} {
					if schema, ok := p.Source.(Schema); ok {
						return schema.GetDirectives()
					}
					return nil
				},
			},
		},
	})

	__EnumValue = NewObject(ObjectConfig{
		Name: "__EnumValue",
		Fields: FieldConfigMap{
			"name": &FieldConfig{
				Type: NewNonNull(String),
			},
			"description": &FieldConfig{
				Type: String,
			},
			"isDeprecated": &FieldConfig{
				Type: NewNonNull(Boolean),
				Resolve: func(p GQLFRParams) interface{} {
					if field, ok := p.Source.(*EnumValueDefinition); ok {
						return (field.DeprecationReason != "")
					}
					return false
				},
			},
			"deprecationReason": &FieldConfig{
				Type: String,
			},
		},
	})

	// Again, adding field configs to __Type that have cyclic reference here
	// because golang don't like them too much during init/compile-time
	__Type.AddFieldConfig("fields", &FieldConfig{
		Type: NewList(NewNonNull(__Field)),
		Args: FieldConfigArgument{
			"includeDeprecated": &ArgumentConfig{
				Type:         Boolean,
				DefaultValue: false,
			},
		},
		Resolve: func(p GQLFRParams) interface{} {
			includeDeprecated, _ := p.Args["includeDeprecated"].(bool)
			switch ttype := p.Source.(type) {
			case *Object:
				if ttype == nil {
					return nil
				}
				fields := []*FieldDefinition{}
				for _, field := range ttype.GetFields() {
					if !includeDeprecated && field.DeprecationReason != "" {
						continue
					}
					fields = append(fields, field)
				}
				return fields
			case *Interface:
				if ttype == nil {
					return nil
				}
				fields := []*FieldDefinition{}
				for _, field := range ttype.GetFields() {
					if !includeDeprecated && field.DeprecationReason != "" {
						continue
					}
					fields = append(fields, field)
				}
				return fields
			}
			return nil
		},
	})
	__Type.AddFieldConfig("interfaces", &FieldConfig{
		Type: NewList(NewNonNull(__Type)),
		Resolve: func(p GQLFRParams) interface{} {
			switch ttype := p.Source.(type) {
			case *Object:
				return ttype.GetInterfaces()
			}
			return nil
		},
	})
	__Type.AddFieldConfig("possibleTypes", &FieldConfig{
		Type: NewList(NewNonNull(__Type)),
		Resolve: func(p GQLFRParams) interface{} {
			switch ttype := p.Source.(type) {
			case *Interface:
				return ttype.GetPossibleTypes()
			case *Union:
				return ttype.GetPossibleTypes()
			}
			return nil
		},
	})
	__Type.AddFieldConfig("enumValues", &FieldConfig{
		Type: NewList(NewNonNull(__EnumValue)),
		Args: FieldConfigArgument{
			"includeDeprecated": &ArgumentConfig{
				Type:         Boolean,
				DefaultValue: false,
			},
		},
		Resolve: func(p GQLFRParams) interface{} {
			includeDeprecated, _ := p.Args["includeDeprecated"].(bool)
			switch ttype := p.Source.(type) {
			case *Enum:
				if includeDeprecated {
					return ttype.GetValues()
				}
				values := []*EnumValueDefinition{}
				for _, value := range ttype.GetValues() {
					if value.DeprecationReason != "" {
						continue
					}
					values = append(values, value)
				}
				return values
			}
			return nil
		},
	})
	__Type.AddFieldConfig("inputFields", &FieldConfig{
		Type: NewList(NewNonNull(__InputValue)),
		Resolve: func(p GQLFRParams) interface{} {
			switch ttype := p.Source.(type) {
			case *InputObject:
				fields := []*InputObjectField{}
				for _, field := range ttype.GetFields() {
					fields = append(fields, field)
				}
				return fields
			}
			return nil
		},
	})
	__Type.AddFieldConfig("ofType", &FieldConfig{
		Type: __Type,
	})

	/**
	 * Note that these are FieldDefinition and not FieldConfig,
	 * so the format for args is different.
	 */

	SchemaMetaFieldDef = &FieldDefinition{
		Name:        "__schema",
		Type:        NewNonNull(__Schema),
		Description: "Access the current type schema of this server.",
		Args:        []*Argument{},
		Resolve: func(p GQLFRParams) interface{} {
			return p.Info.Schema
		},
	}
	TypeMetaFieldDef = &FieldDefinition{
		Name:        "__type",
		Type:        __Type,
		Description: "Request the type information of a single type.",
		Args: []*Argument{
			&Argument{
				Name: "name",
				Type: NewNonNull(String),
			},
		},
		Resolve: func(p GQLFRParams) interface{} {
			name, ok := p.Args["name"].(string)
			if !ok {
				return nil
			}
			return p.Info.Schema.GetType(name)
		},
	}

	TypeNameMetaFieldDef = &FieldDefinition{
		Name:        "__typename",
		Type:        NewNonNull(String),
		Description: "The name of the current Object type at runtime.",
		Args:        []*Argument{},
		Resolve: func(p GQLFRParams) interface{} {
			return p.Info.ParentType.GetName()
		},
	}

}

/**
 * Produces a GraphQL Value AST given a Golang value.
 *
 * Optionally, a GraphQL type may be provided, which will be used to
 * disambiguate between value primitives.
 *
 * | JSON Value    | GraphQL Value        |
 * | ------------- | -------------------- |
 * | Object        | Input Object         |
 * | Array         | List                 |
 * | Boolean       | Boolean              |
 * | String        | String / Enum Value  |
 * | Number        | Int / Float          |
 *
 */
func astFromValue(value interface{}, ttype Type) ast.Value {

	if ttype, ok := ttype.(*NonNull); ok {
		// Note: we're not checking that the result is non-null.
		// This function is not responsible for validating the input value.
		val := astFromValue(value, ttype.OfType)
		return val
	}
	if isNullish(value) {
		return nil
	}
	valueVal := reflect.ValueOf(value)
	if !valueVal.IsValid() {
		return nil
	}
	if valueVal.Type().Kind() == reflect.Ptr {
		valueVal = valueVal.Elem()
	}
	if !valueVal.IsValid() {
		return nil
	}

	// Convert Golang slice to GraphQL list. If the Type is a list, but
	// the value is not an array, convert the value using the list's item type.
	if ttype, ok := ttype.(*List); ok {
		if valueVal.Type().Kind() == reflect.Slice {
			itemType := ttype.OfType
			values := []ast.Value{}
			for i := 0; i < valueVal.Len(); i++ {
				item := valueVal.Index(i).Interface()
				itemAST := astFromValue(item, itemType)
				if itemAST != nil {
					values = append(values, itemAST)
				}
			}
			return ast.NewListValue(&ast.ListValue{
				Values: values,
			})
		} else {
			// Because GraphQL will accept single values as a "list of one" when
			// expecting a list, if there's a non-array value and an expected list type,
			// create an AST using the list's item type.
			val := astFromValue(value, ttype.OfType)
			return val
		}
	}

	if valueVal.Type().Kind() == reflect.Map {
		// TODO: implement astFromValue from Map to Value
	}

	if value, ok := value.(bool); ok {
		return ast.NewBooleanValue(&ast.BooleanValue{
			Value: value,
		})
	}
	if value, ok := value.(int); ok {
		if ttype == Float {
			return ast.NewIntValue(&ast.IntValue{
				Value: fmt.Sprintf("%v.0", value),
			})
		}
		return ast.NewIntValue(&ast.IntValue{
			Value: fmt.Sprintf("%v", value),
		})
	}
	if value, ok := value.(float32); ok {
		return ast.NewFloatValue(&ast.FloatValue{
			Value: fmt.Sprintf("%v", value),
		})
	}
	if value, ok := value.(float64); ok {
		return ast.NewFloatValue(&ast.FloatValue{
			Value: fmt.Sprintf("%v", value),
		})
	}

	if value, ok := value.(string); ok {
		if _, ok := ttype.(*Enum); ok {
			return ast.NewEnumValue(&ast.EnumValue{
				Value: fmt.Sprintf("%v", value),
			})
		}
		return ast.NewStringValue(&ast.StringValue{
			Value: fmt.Sprintf("%v", value),
		})
	}

	// fallback, treat as string
	return ast.NewStringValue(&ast.StringValue{
		Value: fmt.Sprintf("%v", value),
	})
}
