package types

import (
	"fmt"
	"github.com/chris-ramon/graphql-go/language/ast"
)

type GraphQLType interface {
	GetName() string
	GetDescription() string
	Coerce(value interface{}) interface{}
	CoerceLiteral(value interface{}) interface{}
	String() string
}

var _ GraphQLType = (*GraphQLScalarType)(nil)
var _ GraphQLType = (*GraphQLObjectType)(nil)
var _ GraphQLType = (*GraphQLInterfaceType)(nil)
var _ GraphQLType = (*GraphQLUnionType)(nil)
var _ GraphQLType = (*GraphQLEnumType)(nil)
var _ GraphQLType = (*GraphQLInputObjectType)(nil)
var _ GraphQLType = (*GraphQLList)(nil)
var _ GraphQLType = (*GraphQLNonNull)(nil)

type GraphQLArgument struct {
	Name         string
	Type         GraphQLInputType
	DefaultValue interface{}
	Description  string
}

//type GraphQLNonNull interface {
//	GetName() string
//	GetDescription() string
//	Coerce(value interface{}) interface{}
//	CoerceLiteral(value interface{}) interface{}
//	ToString() string
//}

type GraphQLNonNull struct {
	OfType GraphQLType
}

func NewGraphQLNonNull(ofType GraphQLType) *GraphQLNonNull {
	return &GraphQLNonNull{
		OfType: ofType,
	}
}
func (gl *GraphQLNonNull) GetName() string {
	return fmt.Sprintf("%v", gl.OfType)
}
func (gl *GraphQLNonNull) GetDescription() string {
	return ""
}
func (gl *GraphQLNonNull) Coerce(value interface{}) interface{} {
	return value
}
func (gl *GraphQLNonNull) CoerceLiteral(value interface{}) interface{} {
	return value
}
func (gl *GraphQLNonNull) String() string {
	if gl.OfType != nil {
		return gl.OfType.GetName()
	}
	return ""
}

type GraphQLResolveInfo struct {
	FieldName      string
	FieldASTs      []*ast.Field
	ReturnType     GraphQLOutputType
	ParentType     GraphQLCompositeType
	Schema         GraphQLSchema
	Fragments      map[string]ast.Definition
	RootValue      interface{}
	Operation      ast.Definition
	VariableValues map[string]interface{}
}

type GraphQLCompositeType interface {
	GetName() string
}

var _ GraphQLCompositeType = (*GraphQLObjectType)(nil)
var _ GraphQLCompositeType = (*GraphQLInterfaceType)(nil)
var _ GraphQLCompositeType = (*GraphQLUnionType)(nil)

type GraphQLAbstractType interface {
	GetObjectType(value interface{}, info GraphQLResolveInfo) *GraphQLObjectType
	GetPossibleTypes() []*GraphQLObjectType
	IsPossibleType(ttype *GraphQLObjectType) bool
}

var _ GraphQLAbstractType = (*GraphQLInterfaceType)(nil)
var _ GraphQLAbstractType = (*GraphQLUnionType)(nil)
