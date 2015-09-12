package types

import (
	"fmt"
	"github.com/chris-ramon/graphql-go/errors"
	"github.com/chris-ramon/graphql-go/language/ast"
	"github.com/kr/pretty"
)

type Schema interface{}

type GraphQLResult struct {
	Data   interface{}
	Errors []graphqlerrors.GraphQLFormattedError
}

func (gqR *GraphQLResult) HasErrors() bool {
	return (len(gqR.Errors) > 0)
}

type GraphQLEnumType struct {
	Name string
}
func (gt *GraphQLEnumType) GetName() string {
	return gt.Name
}
func (gt *GraphQLEnumType) SetName(name string) {
	gt.Name = name
}
func (gt *GraphQLEnumType) GetDescription() string {
	return ""
}
func (gt *GraphQLEnumType) Coerce(value interface{}) interface{} {
	return value
}
func (gt *GraphQLEnumType) CoerceLiteral(value interface{}) interface{} {
	return value
}
func (gt *GraphQLEnumType) ToString() string {
	return fmt.Sprint("%v", gt)
}

type GraphQLInterfaceTypeConfig struct {
	Name        string
	Fields      interface{}
	ResolveType GraphQLObjectType
	Description string
}

type GraphQLInterfaceType struct {
	Name              string
	Description       string
	TypeConfig        GraphQLInterfaceTypeConfig
	Fields            GraphQLFieldDefinitionMap
	Implementations   []GraphQLObjectType
	PossibleTypeNames map[string]bool
}

func (it *GraphQLInterfaceType) GetName() string {
	return it.Name
}
func (it *GraphQLInterfaceType) SetName(name string) {
	it.Name = name
}
func (it *GraphQLInterfaceType) GetDescription() string {
	return it.Description
}
func (it *GraphQLInterfaceType) Coerce(value interface{}) (r interface{}) {
	return r
}
func (it *GraphQLInterfaceType) CoerceLiteral(value interface{}) (r interface{}) {
	return r
}

func (it *GraphQLInterfaceType) Constructor(config GraphQLInterfaceTypeConfig) {
	//TODO: figure out how to make next line work
	//invariant(config.name, 'Type must be named.');
	it.Name = config.Name
	it.Description = config.Description
	it.TypeConfig = config
	//it.Implementations = []GraphQLObjectType;
}

func (it *GraphQLInterfaceType) GetFields() (fields GraphQLFieldDefinitionMap) {
	//return this._fields ||
	//(this._fields = defineFieldMap(this._typeConfig.fields));
	return fields
}

func (it *GraphQLInterfaceType) GetPossibleTypes() []GraphQLObjectType {
	return it.Implementations
}

func (it *GraphQLInterfaceType) IsPossibleType(objType GraphQLObjectType) bool {
	//var possibleTypeNames = this._possibleTypeNames;
	//if (!possibleTypeNames) {
	//this._possibleTypeNames = possibleTypeNames =
	//this.getPossibleTypes().reduce(
	//(map, possibleType) => ((map[possibleType.name] = true), map),
	//{}
	//);
	//}
	//return possibleTypeNames[type.name] === true;
	return true
}

func (it *GraphQLInterfaceType) ResolveType(any interface{}) (r GraphQLObjectType) {
	//var resolver = this._typeConfig.resolveType;
	//return resolver ? resolver(value) : getTypeOf(value, this);
	return r
}

func (it *GraphQLInterfaceType) ToString() string {
	return it.Name
}

type GQLFRParams struct {
	Source     interface{}
	Args       map[string]interface{}
	Context    interface{}
	FieldAST   interface{}
	FieldType  interface{}
	ParentType interface{}
	Schema     GraphQLSchema

	Info       GraphQLResolveInfo
}

type GraphQLFieldResolveFn func(p GQLFRParams) interface{}

type GraphQLScalarTypeConfig struct {
	Name        string
	Description string
}

func (stC *GraphQLScalarTypeConfig) Coerce(value interface{}) (r interface{}) {
	return r
}

func (stC *GraphQLScalarTypeConfig) CoerceLiteral(value ast.Value) (r interface{}) {
	return r
}

type GraphQLScalarType struct {
	Name         string
	Description  string
	ScalarConfig GraphQLScalarTypeConfig
}

func (sT *GraphQLScalarType) GetName() string {
	return sT.Name

}
func (sT *GraphQLScalarType) SetName(name string) {
	sT.Name = name
}
func (sT *GraphQLScalarType) GetDescription() string {
	return sT.Description

}
func (sT *GraphQLScalarType) Coerce(value interface{}) (r interface{}) {
	return r

}
func (sT *GraphQLScalarType) CoerceLiteral(value interface{}) (r interface{}) {
	return r

}
func (sT *GraphQLScalarType) ToString() string {
	return sT.Name
}

type GraphQLOutputType interface {
	GetName() string
	GetDescription() string
	Coerce(value interface{}) (r interface{})
	CoerceLiteral(value interface{}) (r interface{})
	ToString() string
}

type GraphQLInputType interface {
	GetName() string
	GetDescription() string
	Coerce(value interface{}) interface{}
	CoerceLiteral(value interface{}) interface{}
	ToString() string
}
var _ GraphQLInputType = (*GraphQLScalarType)(nil)
var _ GraphQLInputType = (*GraphQLEnumType)(nil)
//var _ GraphQLInputType = (*GraphQLInputObjectType)(nil)
var _ GraphQLInputType = (*GraphQLList)(nil)

type GraphQLFieldArgument struct {
	Name         string
	Type         GraphQLType
	DefaultValue interface{}
	Description  string
}

type GraphQLFieldDefinition struct {
	Name              string
	Description       string
	Type              GraphQLType
	Args              []GraphQLFieldArgument
	Resolve           GraphQLFieldResolveFn
	DeprecationReason string
}

type GraphQLFieldDefinitionMap map[string]GraphQLFieldDefinition

type GraphQLObjectType struct {
	Name        string
	Description string
	//	TypeConfig GraphQLObjectTypeConfig
	Fields      GraphQLFieldDefinitionMap
	Interfaces  []GraphQLInterfaceType
}
func (gt *GraphQLObjectType) GetName() string {
	return gt.Name
}
func (gt *GraphQLObjectType) SetName(name string) {
	gt.Name = name
}
func (gt *GraphQLObjectType) GetDescription() string {
	return ""
}
func (gt *GraphQLObjectType) Coerce(value interface{}) interface{} {
	return value
}
func (gt *GraphQLObjectType) CoerceLiteral(value interface{}) interface{} {
	return value
}
func (gt *GraphQLObjectType) ToString() string {
	return fmt.Sprint("%v", gt.Name)
}

type GraphQLList struct {
	OfType GraphQLType
}

func NewGraphQLList(ofType GraphQLType) *GraphQLList {
	return &GraphQLList{
		OfType: ofType,
	}
}
func (gl *GraphQLList) GetName() string {
	return gl.OfType.GetName()
}
func (gl *GraphQLList) SetName(name string) {
	gl.OfType.SetName(name)
}
func (gl *GraphQLList) GetDescription() string {
	return ""
}
func (gl *GraphQLList) Coerce(value interface{}) interface{} {
	return value
}
func (gl *GraphQLList) CoerceLiteral(value interface{}) interface{} {
	return value
}
func (gl *GraphQLList) ToString() string {
	return fmt.Sprint("%v", gl)
}


type GraphQLSchemaConfig struct {
	Query    GraphQLObjectType
	Mutation GraphQLObjectType
}

// chose to name as GraphQLTypeMap instead of TypeMap
type GraphQLTypeMap map[string]GraphQLType

type GraphQLSchema struct {
	Query        GraphQLObjectType
	SchemaConfig GraphQLSchemaConfig

	typeMap      GraphQLTypeMap
}

func fillInFieldNames(field *GraphQLFieldDefinition, defaultName string) {
	fieldName := field.Name
	if fieldName == "" {
		fieldName = defaultName
		field.Name = fieldName
		if field.Type.GetName() == "" {
			field.Type.SetName(fieldName)
		}
	}
	for _, arg := range field.Args {
		argName := arg.Name
		if argName == "" {
			argName = defaultName
			arg.Name = argName
		}
		if arg.Type.GetName() == "" {
			arg.Type.SetName(argName)
		}
	}
}
func NewGraphQLSchema(config GraphQLSchemaConfig) (GraphQLSchema, error) {

	// allow user to optionally not explicitly specify `Name` in `GraphQLObjectType.Fields`
	// if `Name` is not specified, use FieldDefinitionMap keys
	for key, field := range config.Query.Fields {
		fillInFieldNames(&field, key)
		config.Query.Fields[key] = field

	}
	for key, field := range config.Mutation.Fields {
		fillInFieldNames(&field, key)
		config.Mutation.Fields[key] = field
	}

	schema := GraphQLSchema{
		SchemaConfig: config,
	}

	var err error
	typeMap := GraphQLTypeMap{}
	objectTypes := []GraphQLObjectType{
		schema.GetQueryType(),
		schema.GetMutationType(),
		__Schema,
	}
	for _, objectType := range objectTypes {
		typeMap, err = typeMapReducer(typeMap, &objectType)
		if err != nil {
			return schema, err
		}

	}
	schema.typeMap = typeMap

	return schema, nil
}
func typeMapReducer(typeMap GraphQLTypeMap, objectType GraphQLType) (GraphQLTypeMap, error) {
	var err error
	if objectType == nil {
		return typeMap, nil
	}

	switch objectType := objectType.(type) {
	case *GraphQLList:
		return typeMapReducer(typeMap, objectType.OfType)
	case *GraphQLNonNull:
		return typeMapReducer(typeMap, objectType.OfType)
	}

	if _, ok := typeMap[objectType.GetName()]; ok {
		return typeMap, graphqlerrors.NewGraphQLFormattedError(
			fmt.Sprintf(`Schema must contain unique named types but contains multiple types named "%v".`, objectType.GetName()),
		)
	}
	if objectType.GetName() == "" {
		pretty.Println("-----> EMPTY NAME", objectType.GetName(), objectType)
		return typeMap, nil
	}

	typeMap[objectType.GetName()] = objectType

	switch objectType := objectType.(type) {
	//	case *GraphQLUnionType:
	//	fallthrough
	case *GraphQLInterfaceType:
		for _, innerObjectType := range objectType.GetPossibleTypes() {
			typeMap, err = typeMapReducer(typeMap, &innerObjectType)
			if err != nil {
				return typeMap, err
			}
		}
	case *GraphQLObjectType:
		for _, innerObjectType := range objectType.Interfaces {
			typeMap, err = typeMapReducer(typeMap, &innerObjectType)
			if err != nil {
				return typeMap, err
			}
		}
	}

	switch objectType := objectType.(type) {
	case *GraphQLObjectType:
		fieldMap := objectType.Fields
		for _, field := range fieldMap {
			for _, arg := range field.Args {
				typeMap, err = typeMapReducer(typeMap, arg.Type)
				if err != nil {
					return typeMap, err
				}
			}
			typeMap, err = typeMapReducer(typeMap, field.Type)
			if err != nil {
				return typeMap, err
			}
		}
	case *GraphQLInterfaceType:
		fieldMap := objectType.Fields
		for _, field := range fieldMap {
			for _, arg := range field.Args {
				typeMap, err = typeMapReducer(typeMap, arg.Type)
				if err != nil {
					return typeMap, err
				}
			}
			typeMap, err = typeMapReducer(typeMap, field.Type)
			if err != nil {
				return typeMap, err
			}
		}
	//	case *GraphQLInputObjectType:
	}
	return typeMap, nil
}

func (gq *GraphQLSchema) GetQueryType() GraphQLObjectType {
	return gq.SchemaConfig.Query
}

func (gq *GraphQLSchema) GetMutationType() GraphQLObjectType {
	return gq.SchemaConfig.Mutation
}


func (gq *GraphQLSchema) GetTypeMap() GraphQLTypeMap {
	return gq.typeMap
}

func (gq *GraphQLSchema) GetType(name string) GraphQLType {
	return gq.GetTypeMap()[name]
}

type GraphQLString struct {
	Name        string
	Description string
}

func (gs *GraphQLString) GetName() string {
	return gs.Name
}
func (gs *GraphQLString) SetName(name string) {
	gs.Name = name
}
func (gs *GraphQLString) GetDescription() string {
	return gs.Description
}
func (gs *GraphQLString) Coerce(value interface{}) (r interface{}) {
	return r
}
func (gs *GraphQLString) CoerceLiteral(value interface{}) (r interface{}) {
	return r
}
func (gs *GraphQLString) ToString() string {
	return fmt.Sprint("%v", gs)
}

type GraphQLInt struct {
	Name        string
	Description string
}
func (gi *GraphQLInt) GetName() string {
	return gi.Name
}
func (gi *GraphQLInt) SetName(name string) {
	gi.Name = name
}
func (gi *GraphQLInt) GetDescription() string {
	return gi.Description
}
func (gs *GraphQLInt) Coerce(value interface{}) (r interface{}) {
	return r
}
func (gi *GraphQLInt) CoerceLiteral(value interface{}) (r interface{}) {
	return r
}
func (gi *GraphQLInt) ToString() string {
	return fmt.Sprint("%v", gi)
}