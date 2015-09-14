package types

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
			"SCALAR": GraphQLEnumValueConfig{
				Value:       TypeKindScalar,
				Description: "Indicates this type is a scalar.",
			},
			"OBJECT": GraphQLEnumValueConfig{
				Value: TypeKindObject,
				Description: "Indicates this type is an object. " +
					"`fields` and `interfaces` are valid fields.",
			},
			"INTERFACE": GraphQLEnumValueConfig{
				Value: TypeKindInterface,
				Description: "Indicates this type is an interface. " +
					"`fields` and `possibleTypes` are valid fields.",
			},
			"UNION": GraphQLEnumValueConfig{
				Value: TypeKindUnion,
				Description: "Indicates this type is an interface. " +
					"`fields` and `possibleTypes` are valid fields.",
			},
			"ENUM": GraphQLEnumValueConfig{
				Value: TypeKindEnum,
				Description: "Indicates this type is an enum. " +
					"`enumValues` are valid fields.",
			},
			"INPUT_OBJECT": GraphQLEnumValueConfig{
				Value: TypeKindInputObject,
				Description: "Indicates this type is an input object. " +
					"`inputFields` is a valid field.",
			},
			"LIST": GraphQLEnumValueConfig{
				Value: TypeKindList,
				Description: "Indicates this type is a list. " +
					"`ofType` is a valid field.",
			},
			"NON_NULL": GraphQLEnumValueConfig{
				Value: TypeKindNonNull,
				Description: "Indicates this type is a non-null. " +
					"`ofType` is a valid field.",
			},
		},
	})

	// define __Type without field.Type definitions first
	// we'll point them back to itself `__Type` after init
	__Type = NewGraphQLObjectType(GraphQLObjectTypeConfig{
		Name: "__Type",
		Fields: GraphQLFieldConfigMap{
			"kind": &GraphQLFieldConfig{
				Type: NewGraphQLNonNull(__TypeKind),
				Resolve: func(p GQLFRParams) interface{} {
					// TODO: resolveFn for __Type
					return nil
				},
			},
			"name": &GraphQLFieldConfig{
				Type: GraphQLString,
			},
			"description": &GraphQLFieldConfig{
				Type: GraphQLString,
			},
			"fields": &GraphQLFieldConfig{
				//				Type: NewGraphQLList(NewGraphQLNonNull(__Field)),
				Args: GraphQLFieldConfigArgumentMap{
					"includeDeprecated": &GraphQLArgumentConfig{
						Type:         GraphQLBoolean,
						DefaultValue: false,
					},
				},
				Resolve: func(p GQLFRParams) interface{} {
					// TODO: resolveFn for __Type
					return nil
				},
			},
			"interfaces": &GraphQLFieldConfig{
				//				Type: NewGraphQLList(NewGraphQLNonNull(__Type)),
				Resolve: func(p GQLFRParams) interface{} {
					// TODO: resolveFn for __Type
					return nil
				},
			},
			"possibleTypes": &GraphQLFieldConfig{
				//				Type: NewGraphQLList(NewGraphQLNonNull(__Type)),
				Resolve: func(p GQLFRParams) interface{} {
					// TODO: resolveFn for __Type
					return nil
				},
			},
			"enumValues": &GraphQLFieldConfig{
				//				Type: NewGraphQLList(NewGraphQLNonNull(__EnumValue)),
				Args: GraphQLFieldConfigArgumentMap{
					"includeDeprecated": &GraphQLArgumentConfig{
						Type:         GraphQLBoolean,
						DefaultValue: false,
					},
				},
				Resolve: func(p GQLFRParams) interface{} {
					// TODO: resolveFn for __Type
					return nil
				},
			},
			"inputFields": &GraphQLFieldConfig{
				//				Type: NewGraphQLList(NewGraphQLNonNull(__InputValue)),
				Resolve: func(p GQLFRParams) interface{} {
					// TODO: resolveFn for __Type
					return nil
				},
			},
			"ofType": &GraphQLFieldConfig{
			//				Type: __Type,
			},
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
					// TODO: resolveFn for __InputValue
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
					// TODO: resolveFn for __Field
					return nil
				},
			},
			"type": &GraphQLFieldConfig{
				Type: NewGraphQLNonNull(__Type),
			},
			"isDeprecated": &GraphQLFieldConfig{
				Type: NewGraphQLNonNull(GraphQLBoolean),
				Resolve: func(p GQLFRParams) interface{} {
					// TODO: resolveFn for __Field
					return nil
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
				Resolve: func(p GQLFRParams) interface{} {
					return p.Directive.Args
				},
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
					typeMap := p.Schema.GetTypeMap()
					results := []GraphQLType{}
					for _, ttype := range typeMap {
						results = append(results, ttype)
					}
					return results
				},
			},
			"queryType": &GraphQLFieldConfig{
				Description: "The type that query operations will be rooted at.",
				Type:        NewGraphQLNonNull(__Type),
				Resolve: func(p GQLFRParams) interface{} {
					return p.Schema.GetQueryType()
				},
			},
			"mutationType": &GraphQLFieldConfig{
				Description: `If this server supports mutation, the type that ` +
					`mutation operations will be rooted at.`,
				Type: __Type,
				Resolve: func(p GQLFRParams) interface{} {
					return p.Schema.GetMutationType()
				},
			},
			"directives": &GraphQLFieldConfig{
				Description: `A list of all directives supported by this server.`,
				Type: NewGraphQLNonNull(NewGraphQLList(
					NewGraphQLNonNull(__Directive),
				)),
				Resolve: func(p GQLFRParams) interface{} {
					return p.Schema.GetDirectives()
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
					// TODO: resolveFn for __EnumValue
					return nil
				},
			},
			"deprecationReason": &GraphQLFieldConfig{
				Type: GraphQLString,
			},
		},
	})

	__Type.__setTypesFor__Type(
		__Type,
		__Field,
		__EnumValue,
		__InputValue,
	)

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
			// TODO: resolveFn for SchemaMetaFieldDef
			return nil
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
			// TODO: resolveFn for TypeMetaFieldDef
			return nil
		},
	}

	TypeNameMetaFieldDef = &GraphQLFieldDefinition{
		Name:        "__typename",
		Type:        NewGraphQLNonNull(GraphQLString),
		Description: "The name of the current Object type at runtime.",
		Args:        []*GraphQLArgument{},
		Resolve: func(p GQLFRParams) interface{} {
			// TODO: resolveFn for TypeNameMetaFieldDef
			return nil
		},
	}

}
