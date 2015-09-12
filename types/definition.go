package types
import (
	"github.com/chris-ramon/graphql-go/language/ast"
	"fmt"
)

type GraphQLType interface {
	GetName() string
	SetName(string)
	GetDescription() string
	Coerce(value interface{}) interface{}
	CoerceLiteral(value interface{}) interface{}
	ToString() string
}

var _ GraphQLType = (*GraphQLScalarType)(nil)
var _ GraphQLType = (*GraphQLObjectType)(nil)
var _ GraphQLType = (*GraphQLInterfaceType)(nil)
//var _ GraphQLType = (*GraphQLUnionType)(nil)
var _ GraphQLType = (*GraphQLEnumType)(nil)
//var _ GraphQLType = (*GraphQLInputObjectType)(nil)
var _ GraphQLType = (*GraphQLList)(nil)
var _ GraphQLType = (*GraphQLNonNull)(nil)

type GraphQLArgument struct {
	Name        string
	Type        GraphQLInputType
	DefaultValue interface{}
	Description string
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
func (gl *GraphQLNonNull) GetName() string {
	return gl.OfType.GetName()
}
func (gl *GraphQLNonNull) SetName(name string) {
	gl.OfType.SetName(name)
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
func (gl *GraphQLNonNull) ToString() string {
	return fmt.Sprint("%v", gl)
}

type GraphQLResolveInfo struct {
	FieldName string
	FieldASTs []*ast.Field
	ReturnType GraphQLOutputType
	ParentType GraphQLCompositeType
	Schema GraphQLSchema
	Fragments map[string]ast.Definition
	RootValue interface{}
	Operation ast.Definition
	VariableValues map[string]interface{}
}

type GraphQLCompositeType interface {

}
var _ GraphQLCompositeType = (*GraphQLObjectType)(nil)
var _ GraphQLCompositeType = (*GraphQLInterfaceType)(nil)
//var _ GraphQLCompositeType = (*GraphQLUnionType)(nil)