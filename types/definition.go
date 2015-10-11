package types

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"

	"github.com/chris-ramon/graphql-go/errors"
	"github.com/chris-ramon/graphql-go/language/ast"
)

// These are all of the possible kinds of types.
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

// These types may be used as input types for arguments and directives.
type GraphQLInputType interface {
	GetName() string
	GetDescription() string
	String() string
	GetError() error
}

var _ GraphQLInputType = (*GraphQLScalarType)(nil)
var _ GraphQLInputType = (*GraphQLEnumType)(nil)
var _ GraphQLInputType = (*GraphQLInputObjectType)(nil)
var _ GraphQLInputType = (*GraphQLList)(nil)
var _ GraphQLInputType = (*GraphQLNonNull)(nil)

func IsInputType(ttype GraphQLType) bool {
	namedType := GetNamedType(ttype)
	if _, ok := namedType.(*GraphQLScalarType); ok {
		return true
	}
	if _, ok := namedType.(*GraphQLEnumType); ok {
		return true
	}
	if _, ok := namedType.(*GraphQLInputObjectType); ok {
		return true
	}
	return false
}

func IsOutputType(ttype GraphQLType) bool {
	namedType := GetNamedType(ttype)
	if _, ok := namedType.(*GraphQLScalarType); ok {
		return true
	}
	if _, ok := namedType.(*GraphQLObjectType); ok {
		return true
	}
	if _, ok := namedType.(*GraphQLInterfaceType); ok {
		return true
	}
	if _, ok := namedType.(*GraphQLUnionType); ok {
		return true
	}
	if _, ok := namedType.(*GraphQLEnumType); ok {
		return true
	}
	return false
}

// These types may be used as output types as the result of fields.
type GraphQLOutputType interface {
	GetName() string
	GetDescription() string
	String() string
	GetError() error
}

var _ GraphQLOutputType = (*GraphQLScalarType)(nil)
var _ GraphQLOutputType = (*GraphQLObjectType)(nil)
var _ GraphQLOutputType = (*GraphQLInterfaceType)(nil)
var _ GraphQLOutputType = (*GraphQLUnionType)(nil)
var _ GraphQLOutputType = (*GraphQLEnumType)(nil)
var _ GraphQLOutputType = (*GraphQLList)(nil)
var _ GraphQLOutputType = (*GraphQLNonNull)(nil)

// These types may describe the parent context of a selection set.
type GraphQLCompositeType interface {
	GetName() string
}

var _ GraphQLCompositeType = (*GraphQLObjectType)(nil)
var _ GraphQLCompositeType = (*GraphQLInterfaceType)(nil)
var _ GraphQLCompositeType = (*GraphQLUnionType)(nil)

// These types may describe the parent context of a selection set.
type GraphQLAbstractType interface {
	GetObjectType(value interface{}, info GraphQLResolveInfo) *GraphQLObjectType
	GetPossibleTypes() []*GraphQLObjectType
	IsPossibleType(ttype *GraphQLObjectType) bool
}

var _ GraphQLAbstractType = (*GraphQLInterfaceType)(nil)
var _ GraphQLAbstractType = (*GraphQLUnionType)(nil)

// These named types do not include modifiers like List or NonNull.
type GraphQLNamedType interface {
	String() string
}

var _ GraphQLNamedType = (*GraphQLScalarType)(nil)
var _ GraphQLNamedType = (*GraphQLObjectType)(nil)
var _ GraphQLNamedType = (*GraphQLInterfaceType)(nil)
var _ GraphQLNamedType = (*GraphQLUnionType)(nil)
var _ GraphQLNamedType = (*GraphQLEnumType)(nil)
var _ GraphQLNamedType = (*GraphQLInputObjectType)(nil)

func GetNamedType(ttype GraphQLType) GraphQLNamedType {
	unmodifiedType := ttype
	for {
		if ttype, ok := unmodifiedType.(*GraphQLList); ok {
			unmodifiedType = ttype.OfType
			continue
		}
		if ttype, ok := unmodifiedType.(*GraphQLNonNull); ok {
			unmodifiedType = ttype.OfType
			continue
		}
		break
	}
	return unmodifiedType
}

/**
 * Scalar Type Definition
 *
 * The leaf values of any request and input values to arguments are
 * Scalars (or Enums) and are defined with a name and a series of functions
 * used to parse input from ast or variables and to ensure validity.
 *
 * Example:
 *
 *     var OddType = new GraphQLScalarType({
 *       name: 'Odd',
 *       serialize(value) {
 *         return value % 2 === 1 ? value : null;
 *       }
 *     });
 *
 */
type GraphQLScalarType struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	scalarConfig GraphQLScalarTypeConfig
	err          error
}
type SerializeFn func(value interface{}) interface{}
type ParseValueFn func(value interface{}) interface{}
type ParseLiteralFn func(valueAST ast.Value) interface{}
type GraphQLScalarTypeConfig struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Serialize    SerializeFn
	ParseValue   ParseValueFn
	ParseLiteral ParseLiteralFn
}

func NewGraphQLScalarType(config GraphQLScalarTypeConfig) *GraphQLScalarType {
	st := &GraphQLScalarType{}
	err := invariant(config.Name != "", "Type must be named.")
	if err != nil {
		st.err = err
		return st
	}

	err = assertValidName(config.Name)
	if err != nil {
		st.err = err
		return st
	}

	st.Name = config.Name
	st.Description = config.Description

	err = invariant(
		config.Serialize != nil,
		fmt.Sprintf(`%v must provide "serialize" function. If this custom Scalar is `+
			`also used as an input type, ensure "parseValue" and "parseLiteral" `+
			`functions are also provided.`, st),
	)
	if err != nil {
		st.err = err
		return st
	}
	if config.ParseValue != nil || config.ParseLiteral != nil {
		err = invariant(
			config.ParseValue != nil && config.ParseLiteral != nil,
			fmt.Sprintf(`%v must provide both "parseValue" and "parseLiteral" functions.`, st),
		)
		if err != nil {
			st.err = err
			return st
		}
	}

	st.scalarConfig = config
	return st
}
func (st *GraphQLScalarType) Serialize(value interface{}) interface{} {
	if st.scalarConfig.Serialize == nil {
		return value
	}
	return st.scalarConfig.Serialize(value)
}
func (st *GraphQLScalarType) ParseValue(value interface{}) interface{} {
	if st.scalarConfig.ParseValue == nil {
		return value
	}
	return st.scalarConfig.ParseValue(value)
}
func (st *GraphQLScalarType) ParseLiteral(valueAST ast.Value) interface{} {
	if st.scalarConfig.ParseLiteral == nil {
		return nil
	}
	return st.scalarConfig.ParseLiteral(valueAST)
}
func (st *GraphQLScalarType) GetName() string {
	return st.Name
}
func (st *GraphQLScalarType) GetDescription() string {
	return st.Description

}
func (st *GraphQLScalarType) String() string {
	return st.Name
}
func (st *GraphQLScalarType) GetError() error {
	return st.err
}

/**
 * Object Type Definition
 *
 * Almost all of the GraphQL types you define will be object types. Object types
 * have a name, but most importantly describe their fields.
 *
 * Example:
 *
 *     var AddressType = new GraphQLObjectType({
 *       name: 'Address',
 *       fields: {
 *         street: { type: GraphQLString },
 *         number: { type: GraphQLInt },
 *         formatted: {
 *           type: GraphQLString,
 *           resolve(obj) {
 *             return obj.number + ' ' + obj.street
 *           }
 *         }
 *       }
 *     });
 *
 * When two types need to refer to each other, or a type needs to refer to
 * itself in a field, you can use a function expression (aka a closure or a
 * thunk) to supply the fields lazily.
 *
 * Example:
 *
 *     var PersonType = new GraphQLObjectType({
 *       name: 'Person',
 *       fields: () => ({
 *         name: { type: GraphQLString },
 *         bestFriend: { type: PersonType },
 *       })
 *     });
 *
 */
type GraphQLObjectType struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsTypeOf    IsTypeOfFn

	typeConfig GraphQLObjectTypeConfig
	fields     GraphQLFieldDefinitionMap
	interfaces []*GraphQLInterfaceType
	// Interim alternative to throwing an error during schema definition at run-time
	err error
}

type IsTypeOfFn func(value interface{}, info GraphQLResolveInfo) bool

type GraphQLInterfacesThunk func() []*GraphQLInterfaceType

type GraphQLObjectTypeConfig struct {
	Name        string                `json:"description"`
	Interfaces  interface{}           `json:"interfaces"`
	Fields      GraphQLFieldConfigMap `json:"fields"`
	IsTypeOf    IsTypeOfFn            `json:"isTypeOf"`
	Description string                `json:"description"`
}

func NewGraphQLObjectType(config GraphQLObjectTypeConfig) *GraphQLObjectType {
	objectType := &GraphQLObjectType{}

	err := invariant(config.Name != "", "Type must be named.")
	if err != nil {
		objectType.err = err
		return objectType
	}
	err = assertValidName(config.Name)
	if err != nil {
		objectType.err = err
		return objectType
	}

	objectType.Name = config.Name
	objectType.Description = config.Description
	objectType.IsTypeOf = config.IsTypeOf
	objectType.typeConfig = config

	/*
			addImplementationToInterfaces()
			Update the interfaces to know about this implementation.
			This is an rare and unfortunate use of mutation in the type definition
		 	implementations, but avoids an expensive "getPossibleTypes"
		 	implementation for Interface types.
	*/
	interfaces := objectType.GetInterfaces()
	if interfaces == nil {
		return objectType
	}
	for _, iface := range interfaces {
		iface.implementations = append(iface.implementations, objectType)
	}

	return objectType
}
func (gt *GraphQLObjectType) AddFieldConfig(fieldName string, fieldConfig *GraphQLFieldConfig) {
	if fieldName == "" || fieldConfig == nil {
		return
	}
	gt.typeConfig.Fields[fieldName] = fieldConfig

}
func (gt *GraphQLObjectType) GetName() string {
	return gt.Name
}
func (gt *GraphQLObjectType) GetDescription() string {
	return ""
}
func (gt *GraphQLObjectType) String() string {
	return gt.Name
}
func (gt *GraphQLObjectType) GetFields() GraphQLFieldDefinitionMap {
	fields, err := defineFieldMap(gt, gt.typeConfig.Fields)
	gt.err = err
	gt.fields = fields
	return gt.fields
}
func (gt *GraphQLObjectType) GetInterfaces() []*GraphQLInterfaceType {
	var configInterfaces []*GraphQLInterfaceType
	switch gt.typeConfig.Interfaces.(type) {
	case GraphQLInterfacesThunk:
		configInterfaces = gt.typeConfig.Interfaces.(GraphQLInterfacesThunk)()
	case []*GraphQLInterfaceType:
		configInterfaces = gt.typeConfig.Interfaces.([]*GraphQLInterfaceType)
	case nil:
	default:
		gt.err = errors.New(fmt.Sprintf("Unknown GraphQLObjectType.Interfaces type: %v", reflect.TypeOf(gt.typeConfig.Interfaces)))
		return nil
	}
	interfaces, err := defineInterfaces(gt, configInterfaces)
	gt.err = err
	gt.interfaces = interfaces
	return gt.interfaces
}
func (gt *GraphQLObjectType) GetError() error {
	return gt.err
}

func defineInterfaces(ttype *GraphQLObjectType, interfaces []*GraphQLInterfaceType) ([]*GraphQLInterfaceType, error) {
	ifaces := []*GraphQLInterfaceType{}

	if len(interfaces) == 0 {
		return ifaces, nil
	}
	for _, iface := range interfaces {
		err := invariant(
			iface != nil,
			fmt.Sprintf(`%v may only implement Interface types, it cannot implement: %v.`, ttype, iface),
		)
		if err != nil {
			return ifaces, err
		}
		if iface.ResolveType != nil {
			err = invariant(
				iface.ResolveType != nil,
				fmt.Sprintf(`Interface Type %v does not provide a "resolveType" function `+
					`and implementing Type %v does not provide a "isTypeOf" `+
					`function. There is no way to resolve this implementing type `+
					`during execution.`, iface, ttype),
			)
			if err != nil {
				return ifaces, err
			}
		}
		ifaces = append(ifaces, iface)
	}

	return ifaces, nil
}

func defineFieldMap(ttype GraphQLNamedType, fields GraphQLFieldConfigMap) (GraphQLFieldDefinitionMap, error) {

	resultFieldMap := GraphQLFieldDefinitionMap{}

	err := invariant(
		len(fields) > 0,
		fmt.Sprintf(`%v fields must be an object with field names as keys or a function which return such an object.`, ttype),
	)
	if err != nil {
		return resultFieldMap, err
	}

	for fieldName, field := range fields {
		if field == nil {
			continue
		}
		err = invariant(
			field.Type != nil,
			fmt.Sprintf(`%v.%v field type must be Output Type but got: %v.`, ttype, fieldName, field.Type),
		)
		if err != nil {
			return resultFieldMap, err
		}
		if field.Type.GetError() != nil {
			return resultFieldMap, field.Type.GetError()
		}
		err = assertValidName(fieldName)
		if err != nil {
			return resultFieldMap, err
		}
		fieldDef := &GraphQLFieldDefinition{
			Name:              fieldName,
			Description:       field.Description,
			Type:              field.Type,
			Resolve:           field.Resolve,
			DeprecationReason: field.DeprecationReason,
		}

		fieldDef.Args = []*GraphQLArgument{}
		for argName, arg := range field.Args {
			err := assertValidName(argName)
			if err != nil {
				return resultFieldMap, err
			}
			err = invariant(
				arg != nil,
				fmt.Sprintf(`%v.%v args must be an object with argument names as keys.`, ttype, fieldName),
			)
			if err != nil {
				return resultFieldMap, err
			}
			err = invariant(
				arg.Type != nil,
				fmt.Sprintf(`%v.%v(%v:) argument type must be Input Type but got: %v.`, ttype, fieldName, argName, arg.Type),
			)
			if err != nil {
				return resultFieldMap, err
			}
			fieldArg := &GraphQLArgument{
				Name:         argName,
				Description:  arg.Description,
				Type:         arg.Type,
				DefaultValue: arg.DefaultValue,
			}
			fieldDef.Args = append(fieldDef.Args, fieldArg)
		}
		resultFieldMap[fieldName] = fieldDef
	}
	return resultFieldMap, nil
}

// TODO: clean up GQLFRParams fields
type GQLFRParams struct {
	Source interface{}
	Args   map[string]interface{}
	Info   GraphQLResolveInfo
	Schema GraphQLSchema
}

// TODO: relook at GraphQLFieldResolveFn params
type GraphQLFieldResolveFn func(p GQLFRParams) interface{}

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

type GraphQLFieldConfigMap map[string]*GraphQLFieldConfig

type GraphQLFieldConfig struct {
	Name              string                        `json:"name"` // used by graphlql-relay
	Type              GraphQLOutputType             `json:"type"`
	Args              GraphQLFieldConfigArgumentMap `json:"args"`
	Resolve           GraphQLFieldResolveFn
	DeprecationReason string `json:"deprecationReason"`
	Description       string `json:"description"`
}

type GraphQLFieldConfigArgumentMap map[string]*GraphQLArgumentConfig

type GraphQLArgumentConfig struct {
	Type         GraphQLInputType `json:"type"`
	DefaultValue interface{}      `json:"defaultValue"`
	Description  string           `json:"description"`
}

type GraphQLFieldDefinitionMap map[string]*GraphQLFieldDefinition
type GraphQLFieldDefinition struct {
	Name              string                `json:"name"`
	Description       string                `json:"description"`
	Type              GraphQLOutputType     `json:"type"`
	Args              []*GraphQLArgument    `json:"args"`
	Resolve           GraphQLFieldResolveFn `json:"-"`
	DeprecationReason string                `json:"deprecationReason"`
}

type GraphQLFieldArgument struct {
	Name         string      `json:"name"`
	Type         GraphQLType `json:"type"`
	DefaultValue interface{} `json:"defaultValue"`
	Description  string      `json:"description"`
}

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

/**
 * Interface Type Definition
 *
 * When a field can return one of a heterogeneous set of types, a Interface type
 * is used to describe what types are possible, what fields are in common across
 * all types, as well as a function to determine which type is actually used
 * when the field is resolved.
 *
 * Example:
 *
 *     var EntityType = new GraphQLInterfaceType({
 *       name: 'Entity',
 *       fields: {
 *         name: { type: GraphQLString }
 *       }
 *     });
 *
 */
type GraphQLInterfaceType struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ResolveType ResolveTypeFn

	typeConfig      GraphQLInterfaceTypeConfig
	fields          GraphQLFieldDefinitionMap
	implementations []*GraphQLObjectType
	possibleTypes   map[string]bool

	err error
}
type GraphQLInterfaceTypeConfig struct {
	Name        string                `json:"name"`
	Fields      GraphQLFieldConfigMap `json:"fields"`
	ResolveType ResolveTypeFn
	Description string `json:"description"`
}
type ResolveTypeFn func(value interface{}, info GraphQLResolveInfo) *GraphQLObjectType

func NewGraphQLInterfaceType(config GraphQLInterfaceTypeConfig) *GraphQLInterfaceType {
	it := &GraphQLInterfaceType{}

	err := invariant(config.Name != "", "Type must be named.")
	if err != nil {
		it.err = err
		return it
	}
	err = assertValidName(config.Name)
	if err != nil {
		it.err = err
		return it
	}
	it.Name = config.Name
	it.Description = config.Description
	it.ResolveType = config.ResolveType
	it.typeConfig = config
	it.implementations = []*GraphQLObjectType{}

	return it
}

func (it *GraphQLInterfaceType) AddFieldConfig(fieldName string, fieldConfig *GraphQLFieldConfig) {
	if fieldName == "" || fieldConfig == nil {
		return
	}
	it.typeConfig.Fields[fieldName] = fieldConfig
}
func (it *GraphQLInterfaceType) GetName() string {
	return it.Name
}
func (it *GraphQLInterfaceType) GetDescription() string {
	return it.Description
}
func (it *GraphQLInterfaceType) GetFields() (fields GraphQLFieldDefinitionMap) {
	it.fields, it.err = defineFieldMap(it, it.typeConfig.Fields)
	return it.fields
}
func (it *GraphQLInterfaceType) GetPossibleTypes() []*GraphQLObjectType {
	return it.implementations
}
func (it *GraphQLInterfaceType) IsPossibleType(ttype *GraphQLObjectType) bool {
	if ttype == nil {
		return false
	}
	if len(it.possibleTypes) == 0 {
		possibleTypes := map[string]bool{}
		for _, possibleType := range it.GetPossibleTypes() {
			if possibleType == nil {
				continue
			}
			possibleTypes[possibleType.Name] = true
		}
		it.possibleTypes = possibleTypes
	}
	if val, ok := it.possibleTypes[ttype.Name]; ok {
		return val
	}
	return false
}
func (it *GraphQLInterfaceType) GetObjectType(value interface{}, info GraphQLResolveInfo) *GraphQLObjectType {
	if it.ResolveType != nil {
		return it.ResolveType(value, info)
	}
	return getTypeOf(value, info, it)
}
func (it *GraphQLInterfaceType) String() string {
	return it.Name
}
func (it *GraphQLInterfaceType) GetError() error {
	return it.err
}

func getTypeOf(value interface{}, info GraphQLResolveInfo, abstractType GraphQLAbstractType) *GraphQLObjectType {
	possibleTypes := abstractType.GetPossibleTypes()
	for _, possibleType := range possibleTypes {
		if possibleType.IsTypeOf == nil {
			continue
		}
		if res := possibleType.IsTypeOf(value, info); res {
			return possibleType
		}
	}
	return nil
}

/**
 * Union Type Definition
 *
 * When a field can return one of a heterogeneous set of types, a Union type
 * is used to describe what types are possible as well as providing a function
 * to determine which type is actually used when the field is resolved.
 *
 * Example:
 *
 *     var PetType = new GraphQLUnionType({
 *       name: 'Pet',
 *       types: [ DogType, CatType ],
 *       resolveType(value) {
 *         if (value instanceof Dog) {
 *           return DogType;
 *         }
 *         if (value instanceof Cat) {
 *           return CatType;
 *         }
 *       }
 *     });
 *
 */
type GraphQLUnionType struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ResolveType ResolveTypeFn

	typeConfig    GraphQLUnionTypeConfig
	types         []*GraphQLObjectType
	possibleTypes map[string]bool

	err error
}
type GraphQLUnionTypeConfig struct {
	Name        string               `json:"name"`
	Types       []*GraphQLObjectType `json:"types"`
	ResolveType ResolveTypeFn
	Description string `json:"description"`
}

func NewGraphQLUnionType(config GraphQLUnionTypeConfig) *GraphQLUnionType {
	objectType := &GraphQLUnionType{}

	err := invariant(config.Name != "", "Type must be named.")
	if err != nil {
		objectType.err = err
		return objectType
	}
	err = assertValidName(config.Name)
	if err != nil {
		objectType.err = err
		return objectType
	}
	objectType.Name = config.Name
	objectType.Description = config.Description
	objectType.ResolveType = config.ResolveType

	err = invariant(
		len(config.Types) > 0,
		fmt.Sprintf(`Must provide Array of types for Union %v.`, config.Name),
	)
	if err != nil {
		objectType.err = err
		return objectType
	}
	for _, ttype := range config.Types {
		err := invariant(
			ttype != nil,
			fmt.Sprintf(`%v may only contain Object types, it cannot contain: %v.`, objectType, ttype),
		)
		if err != nil {
			objectType.err = err
			return objectType
		}
		if objectType.ResolveType == nil {
			err = invariant(
				ttype.IsTypeOf != nil,
				fmt.Sprintf(`Union Type %v does not provide a "resolveType" function `+
					`and possible Type %v does not provide a "isTypeOf" `+
					`function. There is no way to resolve this possible type `+
					`during execution.`, objectType, ttype),
			)
			if err != nil {
				objectType.err = err
				return objectType
			}
		}
	}
	objectType.types = config.Types
	objectType.typeConfig = config

	return objectType
}
func (ut *GraphQLUnionType) GetPossibleTypes() []*GraphQLObjectType {
	return ut.types
}
func (ut *GraphQLUnionType) IsPossibleType(ttype *GraphQLObjectType) bool {

	if ttype == nil {
		return false
	}
	if len(ut.possibleTypes) == 0 {
		possibleTypes := map[string]bool{}
		for _, possibleType := range ut.GetPossibleTypes() {
			if possibleType == nil {
				continue
			}
			possibleTypes[possibleType.Name] = true
		}
		ut.possibleTypes = possibleTypes
	}

	if val, ok := ut.possibleTypes[ttype.Name]; ok {
		return val
	}
	return false
}
func (ut *GraphQLUnionType) GetObjectType(value interface{}, info GraphQLResolveInfo) *GraphQLObjectType {
	if ut.ResolveType != nil {
		return ut.ResolveType(value, info)
	}
	return getTypeOf(value, info, ut)
}
func (ut *GraphQLUnionType) String() string {
	return ut.Name
}
func (ut *GraphQLUnionType) GetName() string {
	return ut.Name
}
func (ut *GraphQLUnionType) GetDescription() string {
	return ut.Description
}
func (ut *GraphQLUnionType) GetError() error {
	return ut.err
}

/**
 * Enum Type Definition
 *
 * Some leaf values of requests and input values are Enums. GraphQL serializes
 * Enum values as strings, however internally Enums can be represented by any
 * kind of type, often integers.
 *
 * Example:
 *
 *     var RGBType = new GraphQLEnumType({
 *       name: 'RGB',
 *       values: {
 *         RED: { value: 0 },
 *         GREEN: { value: 1 },
 *         BLUE: { value: 2 }
 *       }
 *     });
 *
 * Note: If a value is not provided in a definition, the name of the enum value
 * will be used as it's internal value.
 */
type GraphQLEnumType struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	enumConfig   GraphQLEnumTypeConfig
	values       []*GraphQLEnumValueDefinition
	valuesLookup map[interface{}]*GraphQLEnumValueDefinition
	nameLookup   map[string]*GraphQLEnumValueDefinition

	err error
}
type GraphQLEnumValueConfigMap map[string]*GraphQLEnumValueConfig
type GraphQLEnumValueConfig struct {
	Value             interface{} `json:"value"`
	DeprecationReason string      `json:"deprecationReason"`
	Description       string      `json:"description"`
}
type GraphQLEnumTypeConfig struct {
	Name        string                    `json:"name"`
	Values      GraphQLEnumValueConfigMap `json:"values"`
	Description string                    `json:"description"`
}
type GraphQLEnumValueDefinition struct {
	Name              string      `json:"name"`
	Value             interface{} `json:"value"`
	DeprecationReason string      `json:"deprecationReason"`
	Description       string      `json:"description"`
}

func NewGraphQLEnumType(config GraphQLEnumTypeConfig) *GraphQLEnumType {
	gt := &GraphQLEnumType{}
	gt.enumConfig = config

	err := assertValidName(config.Name)
	if err != nil {
		gt.err = err
		return gt
	}

	gt.Name = config.Name
	gt.Description = config.Description
	gt.values, err = gt.defineEnumValues(config.Values)
	if err != nil {
		gt.err = err
		return gt
	}

	return gt
}
func (gt *GraphQLEnumType) defineEnumValues(valueMap GraphQLEnumValueConfigMap) ([]*GraphQLEnumValueDefinition, error) {
	values := []*GraphQLEnumValueDefinition{}

	err := invariant(
		len(valueMap) > 0,
		fmt.Sprintf(`%v values must be an object with value names as keys.`, gt),
	)
	if err != nil {
		return values, err
	}

	for valueName, valueConfig := range valueMap {
		err := invariant(
			valueConfig != nil,
			fmt.Sprintf(`%v.%v must refer to an object with a "value" key `+
				`representing an internal value but got: %v.`, gt, valueName, valueConfig),
		)
		if err != nil {
			return values, err
		}
		err = assertValidName(valueName)
		if err != nil {
			return values, err
		}
		value := &GraphQLEnumValueDefinition{
			Name:              valueName,
			Value:             valueConfig.Value,
			DeprecationReason: valueConfig.DeprecationReason,
			Description:       valueConfig.Description,
		}
		if value.Value == nil {
			value.Value = valueName
		}
		values = append(values, value)
	}
	return values, nil
}
func (gt *GraphQLEnumType) GetValues() []*GraphQLEnumValueDefinition {
	return gt.values
}
func (gt *GraphQLEnumType) Serialize(value interface{}) interface{} {
	if enumValue, ok := gt.getValueLookup()[value]; ok {
		return enumValue.Name
	}
	return nil
}
func (gt *GraphQLEnumType) ParseValue(value interface{}) interface{} {
	valueStr, ok := value.(string)
	if !ok {
		return nil
	}
	if enumValue, ok := gt.getNameLookup()[valueStr]; ok {
		return enumValue.Value
	}
	return nil
}
func (gt *GraphQLEnumType) ParseLiteral(valueAST ast.Value) interface{} {
	if valueAST, ok := valueAST.(*ast.EnumValue); ok {
		if enumValue, ok := gt.getNameLookup()[valueAST.Value]; ok {
			return enumValue.Value
		}
	}
	return nil
}
func (gt *GraphQLEnumType) GetName() string {
	return gt.Name
}
func (gt *GraphQLEnumType) GetDescription() string {
	return ""
}
func (gt *GraphQLEnumType) String() string {
	return gt.Name
}
func (gt *GraphQLEnumType) GetError() error {
	return gt.err
}
func (gt *GraphQLEnumType) getValueLookup() map[interface{}]*GraphQLEnumValueDefinition {
	if len(gt.valuesLookup) > 0 {
		return gt.valuesLookup
	}
	valuesLookup := map[interface{}]*GraphQLEnumValueDefinition{}
	for _, value := range gt.GetValues() {
		valuesLookup[value.Value] = value
	}
	gt.valuesLookup = valuesLookup
	return gt.valuesLookup
}

func (gt *GraphQLEnumType) getNameLookup() map[string]*GraphQLEnumValueDefinition {
	if len(gt.nameLookup) > 0 {
		return gt.nameLookup
	}
	nameLookup := map[string]*GraphQLEnumValueDefinition{}
	for _, value := range gt.GetValues() {
		nameLookup[value.Name] = value
	}
	gt.nameLookup = nameLookup
	return gt.nameLookup
}

/**
 * Input Object Type Definition
 *
 * An input object defines a structured collection of fields which may be
 * supplied to a field argument.
 *
 * Using `NonNull` will ensure that a value must be provided by the query
 *
 * Example:
 *
 *     var GeoPoint = new GraphQLInputObjectType({
 *       name: 'GeoPoint',
 *       fields: {
 *         lat: { type: new GraphQLNonNull(GraphQLFloat) },
 *         lon: { type: new GraphQLNonNull(GraphQLFloat) },
 *         alt: { type: GraphQLFloat, defaultValue: 0 },
 *       }
 *     });
 *
 */
type GraphQLInputObjectType struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	typeConfig InputObjectConfig
	fields     InputObjectFieldMap

	err error
}
type InputObjectFieldConfig struct {
	Type         GraphQLInputType `json:"type"`
	DefaultValue interface{}      `json:"defaultValue"`
	Description  string           `json:"description"`
}
type InputObjectField struct {
	Name         string           `json:"name"`
	Type         GraphQLInputType `json:"type"`
	DefaultValue interface{}      `json:"defaultValue"`
	Description  string           `json:"description"`
}

func (st *InputObjectField) GetName() string {
	return st.Name
}
func (st *InputObjectField) GetDescription() string {
	return st.Description

}
func (st *InputObjectField) String() string {
	return st.Name
}
func (st *InputObjectField) GetError() error {
	return nil
}

type InputObjectConfigFieldMap map[string]*InputObjectFieldConfig
type InputObjectFieldMap map[string]*InputObjectField
type InputObjectConfigFieldMapThunk func() InputObjectConfigFieldMap
type InputObjectConfig struct {
	Name        string      `json:"name"`
	Fields      interface{} `json:"fields"`
	Description string      `json:"description"`
}

// TODO: rename InputObjectConfig to GraphQLInputObjecTypeConfig for consistency?
func NewGraphQLInputObjectType(config InputObjectConfig) *GraphQLInputObjectType {
	gt := &GraphQLInputObjectType{}
	err := invariant(config.Name != "", "Type must be named.")
	if err != nil {
		gt.err = err
		return gt
	}

	gt.Name = config.Name
	gt.Description = config.Description
	gt.typeConfig = config
	gt.fields = gt.defineFieldMap()
	return gt
}

func (gt *GraphQLInputObjectType) defineFieldMap() InputObjectFieldMap {
	var fieldMap InputObjectConfigFieldMap
	switch gt.typeConfig.Fields.(type) {
	case InputObjectConfigFieldMap:
		fieldMap = gt.typeConfig.Fields.(InputObjectConfigFieldMap)
	case InputObjectConfigFieldMapThunk:
		fieldMap = gt.typeConfig.Fields.(InputObjectConfigFieldMapThunk)()
	}
	resultFieldMap := InputObjectFieldMap{}

	err := invariant(
		len(fieldMap) > 0,
		fmt.Sprintf(`%v fields must be an object with field names as keys or a function which return such an object.`, gt),
	)
	if err != nil {
		gt.err = err
		return resultFieldMap
	}

	for fieldName, fieldConfig := range fieldMap {
		if fieldConfig == nil {
			continue
		}
		err := assertValidName(fieldName)
		if err != nil {
			continue
		}
		err = invariant(
			fieldConfig.Type != nil,
			fmt.Sprintf(`%v.%v field type must be Input Type but got: %v.`, gt, fieldName, fieldConfig.Type),
		)
		if err != nil {
			gt.err = err
			return resultFieldMap
		}
		field := &InputObjectField{}
		field.Name = fieldName
		field.Type = fieldConfig.Type
		field.Description = fieldConfig.Description
		field.DefaultValue = fieldConfig.DefaultValue
		resultFieldMap[fieldName] = field
	}
	return resultFieldMap
}
func (gt *GraphQLInputObjectType) GetFields() InputObjectFieldMap {
	return gt.fields
}
func (gt *GraphQLInputObjectType) GetName() string {
	return gt.Name
}
func (gt *GraphQLInputObjectType) GetDescription() string {
	return gt.Description
}
func (gt *GraphQLInputObjectType) String() string {
	return gt.Name
}
func (gt *GraphQLInputObjectType) GetError() error {
	return gt.err
}

/**
 * List Modifier
 *
 * A list is a kind of type marker, a wrapping type which points to another
 * type. Lists are often created within the context of defining the fields of
 * an object type.
 *
 * Example:
 *
 *     var PersonType = new GraphQLObjectType({
 *       name: 'Person',
 *       fields: () => ({
 *         parents: { type: new GraphQLList(Person) },
 *         children: { type: new GraphQLList(Person) },
 *       })
 *     })
 *
 */
type GraphQLList struct {
	OfType GraphQLType `json:"ofType"`

	err error
}

func NewGraphQLList(ofType GraphQLType) *GraphQLList {
	gl := &GraphQLList{}

	err := invariant(ofType != nil, fmt.Sprintf(`Can only create List of a GraphQLType but got: %v.`, ofType))
	if err != nil {
		gl.err = err
		return gl
	}

	gl.OfType = ofType
	return gl
}
func (gl *GraphQLList) GetName() string {
	return fmt.Sprintf("%v", gl.OfType)
}
func (gl *GraphQLList) GetDescription() string {
	return ""
}
func (gl *GraphQLList) String() string {
	if gl.OfType != nil {
		return fmt.Sprintf("[%v]", gl.OfType)
	}
	return ""
}
func (gl *GraphQLList) GetError() error {
	return gl.err
}

/**
 * Non-Null Modifier
 *
 * A non-null is a kind of type marker, a wrapping type which points to another
 * type. Non-null types enforce that their values are never null and can ensure
 * an error is raised if this ever occurs during a request. It is useful for
 * fields which you can make a strong guarantee on non-nullability, for example
 * usually the id field of a database row will never be null.
 *
 * Example:
 *
 *     var RowType = new GraphQLObjectType({
 *       name: 'Row',
 *       fields: () => ({
 *         id: { type: new GraphQLNonNull(GraphQLString) },
 *       })
 *     })
 *
 * Note: the enforcement of non-nullability occurs within the executor.
 */
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

var NAME_REGEXP, _ = regexp.Compile("^[_a-zA-Z][_a-zA-Z0-9]*$")

func assertValidName(name string) error {
	return invariant(
		NAME_REGEXP.MatchString(name),
		fmt.Sprintf(`Names must match /^[_a-zA-Z][_a-zA-Z0-9]*$/ but "%v" does not.`, name),
	)
}

// TODO: there is another invariant() func in `executor`
func invariant(condition bool, message string) error {
	if !condition {
		return graphqlerrors.NewGraphQLFormattedError(message)
	}
	return nil
}
