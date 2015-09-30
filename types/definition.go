package types

import (
	"fmt"
	"github.com/chris-ramon/graphql-go/language/ast"
)

type GraphQLType interface {
	GetName() string
	GetDescription() string
	String() string
	GetError() error
}

var _ GraphQLType = (*GraphQLScalarType)(nil)
var _ GraphQLType = (*GraphQLObjectType)(nil)
var _ GraphQLType = (*GraphQLInterfaceType)(nil)
var _ GraphQLType = (*GraphQLUnionType)(nil)
var _ GraphQLType = (*GraphQLEnumType)(nil)
var _ GraphQLType = (*GraphQLInputObjectType)(nil)
var _ GraphQLType = (*GraphQLList)(nil)
var _ GraphQLType = (*GraphQLNonNull)(nil)
var _ GraphQLType = (*GraphQLArgument)(nil)

type GraphQLArgument struct {
	Name         string           `json:"name"`
	Type         GraphQLInputType `json:"type"`
	DefaultValue interface{}      `json:"defaultValue"`
	Description  string           `json:"description"`
}

func (st *GraphQLArgument) GetName() string {
	return st.Name
}
func (st *GraphQLArgument) GetDescription() string {
	return st.Description

}
func (st *GraphQLArgument) String() string {
	return st.Name
}
func (st *GraphQLArgument) GetError() error {
	return nil
}

type GraphQLNonNull struct {
	Name   string      `json:"name"` // added to conform with introspection for NonNull.Name = nil
	OfType GraphQLType `json:"ofType"`

	err error
}

func NewGraphQLNonNull(ofType GraphQLType) *GraphQLNonNull {
	gl := &GraphQLNonNull{}

	_, isOfTypeNonNull := ofType.(*GraphQLNonNull)
	err := invariant(ofType != nil && !isOfTypeNonNull, fmt.Sprintf(`Can only create NonNull of a Nullable GraphQLType but got: %v.`, ofType))
	if err != nil {
		gl.err = err
		return gl
	}
	gl.OfType = ofType
	return gl
}
func (gl *GraphQLNonNull) GetName() string {
	return fmt.Sprintf("%v!", gl.OfType)
}
func (gl *GraphQLNonNull) GetDescription() string {
	return ""
}
func (gl *GraphQLNonNull) String() string {
	if gl.OfType != nil {
		return gl.GetName()
	}
	return ""
}
func (gl *GraphQLNonNull) GetError() error {
	return gl.err
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
