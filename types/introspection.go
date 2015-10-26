package types

import (
	"fmt"
	"github.com/chris-ramon/graphql/language/ast"
	"github.com/chris-ramon/graphql/language/printer"
	"math"
	"reflect"
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

var __Directive *GraphQLObjectType
var __Schema *GraphQLObjectType
var __Type *GraphQLObjectType
var __Field *GraphQLObjectType
var __InputValue *GraphQLObjectType
var __EnumValue *GraphQLObjectType

var __TypeKind *GraphQLEnumType

var SchemaMetaFieldDef *GraphQLFieldDefinition
var TypeMetaFieldDef *GraphQLFieldDefinition
var TypeNameMetaFieldDef *GraphQLFieldDefinition

func init() {

	__TypeKind = NewGraphQLEnumType(GraphQLEnumTypeConfig{
		Name:        "__TypeKind",
		Description: "An enum describing what kind of type a given __Type is",
		Values: GraphQLEnumValueConfigMap{
			"SCALAR": &GraphQLEnumValueConfig{
				Value:       TypeKindScalar,
				Description: "Indicates this type is a scalar.",
			},
			"OBJECT": &GraphQLEnumValueConfig{
				Value: TypeKindObject,
				Description: "Indicates this type is an object. " +
					"`fields` and `interfaces` are valid fields.",
			},
			"INTERFACE": &GraphQLEnumValueConfig{
				Value: TypeKindInterface,
				Description: "Indicates this type is an interface. " +
					"`fields` and `possibleTypes` are valid fields.",
			},
			"UNION": &GraphQLEnumValueConfig{
				Value: TypeKindUnion,
				Description: "Indicates this type is a union. " +
					"`possibleTypes` is a valid field.",
			},
			"ENUM": &GraphQLEnumValueConfig{
				Value: TypeKindEnum,
				Description: "Indicates this type is an enum. " +
					"`enumValues` is a valid field.",
			},
			"INPUT_OBJECT": &GraphQLEnumValueConfig{
				Value: TypeKindInputObject,
				Description: "Indicates this type is an input object. " +
					"`inputFields` is a valid field.",
			},
			"LIST": &GraphQLEnumValueConfig{
				Value: TypeKindList,
				Description: "Indicates this type is a list. " +
					"`ofType` is a valid field.",
			},
			"NON_NULL": &GraphQLEnumValueConfig{
				Value: TypeKindNonNull,
				Description: "Indicates this type is a non-null. " +
					"`ofType` is a valid field.",
			},
		},
	})

	// Note: some fields (for e.g "fields", "interfaces") are defined later due to cyclic reference
	__Type = NewGraphQLObjectType(GraphQLObjectTypeConfig{
		Name: "__Type",
		Fields: GraphQLFieldConfigMap{
			"kind": &GraphQLFieldConfig{
				Type: NewGraphQLNonNull(__TypeKind),
				Resolve: func(p GQLFRParams) interface{} {
					switch p.Source.(type) {
					case *GraphQLScalarType:
						return TypeKindScalar
					case *GraphQLObjectType:
						return TypeKindObject
					case *GraphQLInterfaceType:
						return TypeKindInterface
					case *GraphQLUnionType:
						return TypeKindUnion
					case *GraphQLEnumType:
						return TypeKindEnum
					case *GraphQLInputObjectType:
						return TypeKindInputObject
					case *GraphQLList:
						return TypeKindList
					case *GraphQLNonNull:
						return TypeKindNonNull
					}
					panic(fmt.Sprintf("Unknown kind of type: %v", p.Source))
				},
			},
			"name": &GraphQLFieldConfig{
				Type: GraphQLString,
			},
			"description": &GraphQLFieldConfig{
				Type: GraphQLString,
			},
			"fields":        &GraphQLFieldConfig{},
			"interfaces":    &GraphQLFieldConfig{},
			"possibleTypes": &GraphQLFieldConfig{},
			"enumValues":    &GraphQLFieldConfig{},
			"inputFields":   &GraphQLFieldConfig{},
			"ofType":        &GraphQLFieldConfig{},
		},
	})

	__InputValue = NewGraphQLObjectType(GraphQLObjectTypeConfig{
		Name: "__InputValue",
		Fields: GraphQLFieldConfigMap{
			"name": &GraphQLFieldConfig{
				Type: NewGraphQLNonNull(GraphQLString),
			},
			"description": &GraphQLFieldConfig{
				Type: GraphQLString,
			},
			"type": &GraphQLFieldConfig{
				Type: NewGraphQLNonNull(__Type),
			},
			"defaultValue": &GraphQLFieldConfig{
				Type: GraphQLString,
				Resolve: func(p GQLFRParams) interface{} {
					if inputVal, ok := p.Source.(*GraphQLArgument); ok {
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

	__Field = NewGraphQLObjectType(GraphQLObjectTypeConfig{
		Name: "__Field",
		Fields: GraphQLFieldConfigMap{
			"name": &GraphQLFieldConfig{
				Type: NewGraphQLNonNull(GraphQLString),
			},
			"description": &GraphQLFieldConfig{
				Type: GraphQLString,
			},
			"args": &GraphQLFieldConfig{
				Type: NewGraphQLNonNull(NewGraphQLList(NewGraphQLNonNull(__InputValue))),
				Resolve: func(p GQLFRParams) interface{} {
					if field, ok := p.Source.(*GraphQLFieldDefinition); ok {
						return field.Args
					}
					return []interface{}{}
				},
			},
			"type": &GraphQLFieldConfig{
				Type: NewGraphQLNonNull(__Type),
			},
			"isDeprecated": &GraphQLFieldConfig{
				Type: NewGraphQLNonNull(GraphQLBoolean),
				Resolve: func(p GQLFRParams) interface{} {
					if field, ok := p.Source.(*GraphQLFieldDefinition); ok {
						return (field.DeprecationReason != "")
					}
					return false
				},
			},
			"deprecationReason": &GraphQLFieldConfig{
				Type: GraphQLString,
			},
		},
	})

	__Directive = NewGraphQLObjectType(GraphQLObjectTypeConfig{
		Name: "__Directive",
		Fields: GraphQLFieldConfigMap{
			"name": &GraphQLFieldConfig{
				Type: NewGraphQLNonNull(GraphQLString),
			},
			"description": &GraphQLFieldConfig{
				Type: GraphQLString,
			},
			"args": &GraphQLFieldConfig{
				Type: NewGraphQLNonNull(NewGraphQLList(
					NewGraphQLNonNull(__InputValue),
				)),
			},
			"onOperation": &GraphQLFieldConfig{
				Type: NewGraphQLNonNull(GraphQLBoolean),
			},
			"onFragment": &GraphQLFieldConfig{
				Type: NewGraphQLNonNull(GraphQLBoolean),
			},
			"onField": &GraphQLFieldConfig{
				Type: NewGraphQLNonNull(GraphQLBoolean),
			},
		},
	})

	__Schema = NewGraphQLObjectType(GraphQLObjectTypeConfig{
		Name: "__Schema",
		Description: `A GraphQL Schema defines the capabilities of a GraphQL
server. It exposes all available types and directives on
the server, as well as the entry points for query and
mutation operations.`,
		Fields: GraphQLFieldConfigMap{
			"types": &GraphQLFieldConfig{
				Description: "A list of all types supported by this server.",
				Type: NewGraphQLNonNull(NewGraphQLList(
					NewGraphQLNonNull(__Type),
				)),
				Resolve: func(p GQLFRParams) interface{} {
					if schema, ok := p.Source.(GraphQLSchema); ok {
						results := []GraphQLType{}
						for _, ttype := range schema.GetTypeMap() {
							results = append(results, ttype)
						}
						return results
					}
					return []GraphQLType{}
				},
			},
			"queryType": &GraphQLFieldConfig{
				Description: "The type that query operations will be rooted at.",
				Type:        NewGraphQLNonNull(__Type),
				Resolve: func(p GQLFRParams) interface{} {
					if schema, ok := p.Source.(GraphQLSchema); ok {
						return schema.GetQueryType()
					}
					return nil
				},
			},
			"mutationType": &GraphQLFieldConfig{
				Description: `If this server supports mutation, the type that ` +
					`mutation operations will be rooted at.`,
				Type: __Type,
				Resolve: func(p GQLFRParams) interface{} {
					if schema, ok := p.Source.(GraphQLSchema); ok {
						if schema.GetMutationType() != nil {
							return schema.GetMutationType()
						}
					}
					return nil
				},
			},
			"directives": &GraphQLFieldConfig{
				Description: `A list of all directives supported by this server.`,
				Type: NewGraphQLNonNull(NewGraphQLList(
					NewGraphQLNonNull(__Directive),
				)),
				Resolve: func(p GQLFRParams) interface{} {
					if schema, ok := p.Source.(GraphQLSchema); ok {
						return schema.GetDirectives()
					}
					return nil
				},
			},
		},
	})

	__EnumValue = NewGraphQLObjectType(GraphQLObjectTypeConfig{
		Name: "__EnumValue",
		Fields: GraphQLFieldConfigMap{
			"name": &GraphQLFieldConfig{
				Type: NewGraphQLNonNull(GraphQLString),
			},
			"description": &GraphQLFieldConfig{
				Type: GraphQLString,
			},
			"isDeprecated": &GraphQLFieldConfig{
				Type: NewGraphQLNonNull(GraphQLBoolean),
				Resolve: func(p GQLFRParams) interface{} {
					if field, ok := p.Source.(*GraphQLEnumValueDefinition); ok {
						return (field.DeprecationReason != "")
					}
					return false
				},
			},
			"deprecationReason": &GraphQLFieldConfig{
				Type: GraphQLString,
			},
		},
	})

	// Again, adding field configs to __Type that have cyclic reference here
	// because golang don't like them too much during init/compile-time
	__Type.AddFieldConfig("fields", &GraphQLFieldConfig{
		Type: NewGraphQLList(NewGraphQLNonNull(__Field)),
		Args: GraphQLFieldConfigArgumentMap{
			"includeDeprecated": &GraphQLArgumentConfig{
				Type:         GraphQLBoolean,
				DefaultValue: false,
			},
		},
		Resolve: func(p GQLFRParams) interface{} {
			includeDeprecated, _ := p.Args["includeDeprecated"].(bool)
			switch ttype := p.Source.(type) {
			case *GraphQLObjectType:
				if ttype == nil {
					return nil
				}
				fields := []*GraphQLFieldDefinition{}
				for _, field := range ttype.GetFields() {
					if !includeDeprecated && field.DeprecationReason != "" {
						continue
					}
					fields = append(fields, field)
				}
				return fields
			case *GraphQLInterfaceType:
				if ttype == nil {
					return nil
				}
				fields := []*GraphQLFieldDefinition{}
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
	__Type.AddFieldConfig("interfaces", &GraphQLFieldConfig{
		Type: NewGraphQLList(NewGraphQLNonNull(__Type)),
		Resolve: func(p GQLFRParams) interface{} {
			switch ttype := p.Source.(type) {
			case *GraphQLObjectType:
				return ttype.GetInterfaces()
			}
			return nil
		},
	})
	__Type.AddFieldConfig("possibleTypes", &GraphQLFieldConfig{
		Type: NewGraphQLList(NewGraphQLNonNull(__Type)),
		Resolve: func(p GQLFRParams) interface{} {
			switch ttype := p.Source.(type) {
			case *GraphQLInterfaceType:
				return ttype.GetPossibleTypes()
			case *GraphQLUnionType:
				return ttype.GetPossibleTypes()
			}
			return nil
		},
	})
	__Type.AddFieldConfig("enumValues", &GraphQLFieldConfig{
		Type: NewGraphQLList(NewGraphQLNonNull(__EnumValue)),
		Args: GraphQLFieldConfigArgumentMap{
			"includeDeprecated": &GraphQLArgumentConfig{
				Type:         GraphQLBoolean,
				DefaultValue: false,
			},
		},
		Resolve: func(p GQLFRParams) interface{} {
			includeDeprecated, _ := p.Args["includeDeprecated"].(bool)
			switch ttype := p.Source.(type) {
			case *GraphQLEnumType:
				if includeDeprecated {
					return ttype.GetValues()
				}
				values := []*GraphQLEnumValueDefinition{}
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
	__Type.AddFieldConfig("inputFields", &GraphQLFieldConfig{
		Type: NewGraphQLList(NewGraphQLNonNull(__InputValue)),
		Resolve: func(p GQLFRParams) interface{} {
			switch ttype := p.Source.(type) {
			case *GraphQLInputObjectType:
				fields := []*InputObjectField{}
				for _, field := range ttype.GetFields() {
					fields = append(fields, field)
				}
				return fields
			}
			return nil
		},
	})
	__Type.AddFieldConfig("ofType", &GraphQLFieldConfig{
		Type: __Type,
	})

	/**
	 * Note that these are GraphQLFieldDefinition and not GraphQLFieldConfig,
	 * so the format for args is different.
	 */

	SchemaMetaFieldDef = &GraphQLFieldDefinition{
		Name:        "__schema",
		Type:        NewGraphQLNonNull(__Schema),
		Description: "Access the current type schema of this server.",
		Args:        []*GraphQLArgument{},
		Resolve: func(p GQLFRParams) interface{} {
			return p.Info.Schema
		},
	}
	TypeMetaFieldDef = &GraphQLFieldDefinition{
		Name:        "__type",
		Type:        __Type,
		Description: "Request the type information of a single type.",
		Args: []*GraphQLArgument{
			&GraphQLArgument{
				Name: "name",
				Type: NewGraphQLNonNull(GraphQLString),
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

	TypeNameMetaFieldDef = &GraphQLFieldDefinition{
		Name:        "__typename",
		Type:        NewGraphQLNonNull(GraphQLString),
		Description: "The name of the current Object type at runtime.",
		Args:        []*GraphQLArgument{},
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
func astFromValue(value interface{}, ttype GraphQLType) ast.Value {

	if ttype, ok := ttype.(*GraphQLNonNull); ok {
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

	// Convert Golang slice to GraphQL list. If the GraphQLType is a list, but
	// the value is not an array, convert the value using the list's item type.
	if ttype, ok := ttype.(*GraphQLList); ok {
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
		// TODO: implement astFromValue from Map to ast.Value
	}

	if value, ok := value.(bool); ok {
		return ast.NewBooleanValue(&ast.BooleanValue{
			Value: value,
		})
	}
	if value, ok := value.(int); ok {
		if ttype == GraphQLFloat {
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
		if _, ok := ttype.(*GraphQLEnumType); ok {
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

// Returns true if a value is null, undefined, or NaN.
// TODO: move isNullish to utilities. This is a copy of isNullish in `executor`
func isNullish(value interface{}) bool {
	if value, ok := value.(string); ok {
		return value == ""
	}
	if value, ok := value.(int); ok {
		return math.IsNaN(float64(value))
	}
	if value, ok := value.(float32); ok {
		return math.IsNaN(float64(value))
	}
	if value, ok := value.(float64); ok {
		return math.IsNaN(value)
	}
	return value == nil
}
