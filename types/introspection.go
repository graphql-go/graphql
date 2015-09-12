package types

// TODO: __Schema etc
var __Schema GraphQLObjectType = GraphQLObjectType{
	Name: "__Schema",
	Description: `A GraphQL Schema defines the capabilities of a GraphQL
server. It exposes all available types and directives on
the server, as well as the entry points for query and
mutation operations.`,
	Fields: GraphQLFieldDefinitionMap{
		"types": GraphQLFieldDefinition{
			Description: "A list of all types supported by this server.",
		},
	},
}

// TODO: SchemaMetaFieldDef etc
var SchemaMetaFieldDef GraphQLFieldDefinition = GraphQLFieldDefinition{
	Name: "__schema",
}
var TypeMetaFieldDef GraphQLFieldDefinition = GraphQLFieldDefinition{
	Name: "__type",
}
var TypeNameMetaFieldDef GraphQLFieldDefinition = GraphQLFieldDefinition{
	Name: "__typename",
}