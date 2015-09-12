package types

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