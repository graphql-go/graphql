package types

import (
	"fmt"
	"github.com/chris-ramon/graphql-go/errors"
	"github.com/chris-ramon/graphql-go/language/ast"
	"regexp"
)

type Schema interface{}

type GraphQLResult struct {
	Data   interface{}                           `json:"data"`
	Errors []graphqlerrors.GraphQLFormattedError `json:"errors"`
}

func (gqR *GraphQLResult) HasErrors() bool {
	return (len(gqR.Errors) > 0)
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

type GraphQLEnumType struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	enumConfig   GraphQLEnumTypeConfig
	values       []*GraphQLEnumValueDefinition
	valuesLookup map[interface{}]*GraphQLEnumValueDefinition
	nameLookup   map[string]*GraphQLEnumValueDefinition

	err error
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
func (gt *GraphQLEnumType) Coerce(value interface{}) interface{} {
	return value
}
func (gt *GraphQLEnumType) CoerceLiteral(value interface{}) interface{} {
	return value
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

type GraphQLInterfaceTypeConfig struct {
	Name        string                `json:"name"`
	Fields      GraphQLFieldConfigMap `json:"fields"`
	ResolveType ResolveTypeFn
	Description string `json:"description"`
}

type ResolveTypeFn func(value interface{}, info GraphQLResolveInfo) *GraphQLObjectType
type GraphQLInterfaceType struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ResolveType ResolveTypeFn

	typeConfig      GraphQLInterfaceTypeConfig
	fields          GraphQLFieldDefinitionMap
	implementations []*GraphQLObjectType
	possibleTypes   map[string]bool

	// Interim alternative to throwing an error during schema definition at run-time
	err error
}

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
func (it *GraphQLInterfaceType) Coerce(value interface{}) (r interface{}) {
	return value
}
func (it *GraphQLInterfaceType) CoerceLiteral(value interface{}) (r interface{}) {
	return value
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
func (it *GraphQLInterfaceType) String() string {
	return it.Name
}
func (it *GraphQLInterfaceType) GetError() error {
	return it.err
}

// TODO: clean up GQLFRParams fields
type GQLFRParams struct {
	Source interface{}
	Args   map[string]interface{}
	Info   GraphQLResolveInfo

	//	Context    interface{}
	//	FieldAST   interface{}
	//	FieldType  interface{}
	//	ParentType interface{}
	//
	Schema GraphQLSchema
	//	Directive GraphQLDirective

}

// TODO: relook at GraphQLFieldResolveFn params
type GraphQLFieldResolveFn func(p GQLFRParams) interface{}

type GraphQLOutputType interface {
	GetName() string
	GetDescription() string
	Coerce(value interface{}) (r interface{})
	CoerceLiteral(value interface{}) (r interface{})
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

type GraphQLInputType interface {
	GetName() string
	GetDescription() string
	Coerce(value interface{}) interface{}
	CoerceLiteral(value interface{}) interface{}
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

type GraphQLFieldArgument struct {
	Name         string      `json:"name"`
	Type         GraphQLType `json:"type"`
	DefaultValue interface{} `json:"defaultValue"`
	Description  string      `json:"description"`
}

type GraphQLFieldDefinition struct {
	Name              string                `json:"name"`
	Description       string                `json:"description"`
	Type              GraphQLOutputType     `json:"type"`
	Args              []*GraphQLArgument    `json:"args"`
	Resolve           GraphQLFieldResolveFn `json:"-"`
	DeprecationReason string                `json:"deprecationReason"`
}

type GraphQLFieldDefinitionMap map[string]*GraphQLFieldDefinition

type IsTypeOfFn func(value interface{}, info GraphQLResolveInfo) bool

type GraphQLObjectTypeConfig struct {
	Name        string                  `json:"description"`
	Interfaces  []*GraphQLInterfaceType `json:"interfaces"`
	Fields      GraphQLFieldConfigMap   `json:"fields"`
	IsTypeOf    IsTypeOfFn              `json:"isTypeOf"`
	Description string                  `json:"description"`
}
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
	for _, iface := range objectType.GetInterfaces() {
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
func (gt *GraphQLObjectType) Coerce(value interface{}) interface{} {
	return value
}
func (gt *GraphQLObjectType) CoerceLiteral(value interface{}) interface{} {
	return value
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
	interfaces, err := defineInterfaces(gt, gt.typeConfig.Interfaces)
	gt.err = err
	gt.interfaces = interfaces
	return gt.interfaces
}
func (gt *GraphQLObjectType) GetError() error {
	return gt.err
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
func (gl *GraphQLList) Coerce(value interface{}) interface{} {
	return value
}
func (gl *GraphQLList) CoerceLiteral(value interface{}) interface{} {
	return value
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

type GraphQLUnionTypeConfig struct {
	Name        string               `json:"name"`
	Types       []*GraphQLObjectType `json:"types"`
	ResolveType ResolveTypeFn
	Description string `json:"description"`
}
type GraphQLUnionType struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ResolveType ResolveTypeFn

	typeConfig    GraphQLUnionTypeConfig
	types         []*GraphQLObjectType
	possibleTypes map[string]bool

	err error
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
func (ut *GraphQLUnionType) Coerce(value interface{}) (r interface{}) {
	return value
}
func (ut *GraphQLUnionType) CoerceLiteral(value interface{}) (r interface{}) {
	return value
}
func (ut *GraphQLUnionType) GetError() error {
	return ut.err
}

// These named types do not include modifiers like List or NonNull.
type GraphQLNamedType interface {
	String() string
}

var (
	_ GraphQLNamedType = (*GraphQLScalarType)(nil)
	_ GraphQLNamedType = (*GraphQLObjectType)(nil)
	_ GraphQLNamedType = (*GraphQLInterfaceType)(nil)
	_ GraphQLNamedType = (*GraphQLUnionType)(nil)
	_ GraphQLNamedType = (*GraphQLEnumType)(nil)
	_ GraphQLNamedType = (*GraphQLInputObjectType)(nil)
)

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

// TODO: there is another invariant() func in `executor`
func invariant(condition bool, message string) error {
	if !condition {
		return graphqlerrors.NewGraphQLFormattedError(message)
	}
	return nil
}

var NAME_REGEXP, _ = regexp.Compile("^[_a-zA-Z][_a-zA-Z0-9]*$")

func assertValidName(name string) error {
	return invariant(
		NAME_REGEXP.MatchString(name),
		fmt.Sprintf(`Names must match /^[_a-zA-Z][_a-zA-Z0-9]*$/ but "%v" does not.`, name),
	)
}

func defineFieldMap(ttype GraphQLNamedType, fields GraphQLFieldConfigMap) (GraphQLFieldDefinitionMap, error) {

	resultFieldMap := GraphQLFieldDefinitionMap{}

	err := invariant(
		len(fields) > 0,
		fmt.Sprintf(`%v fields must be an object with field names as keys.`, ttype),
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
type InputObjectConfigFieldMap map[string]*InputObjectFieldConfig
type InputObjectFieldMap map[string]*InputObjectField
type InputObjectConfig struct {
	Name        string                    `json:"name"`
	Fields      InputObjectConfigFieldMap `json:"fields"`
	Description string                    `json:"description"`
}

// TODO: GraphQLInputObjectType is more of a map in Golang
type GraphQLInputObjectType struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	typeConfig InputObjectConfig
	fields     InputObjectFieldMap

	err error
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
	fieldMap := gt.typeConfig.Fields
	resultFieldMap := InputObjectFieldMap{}

	err := invariant(
		len(fieldMap) > 0,
		fmt.Sprintf(`%v fields must be an object with field names as keys.`, gt),
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
		field.DefaultValue = field.DefaultValue
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
func (gt *GraphQLInputObjectType) Coerce(value interface{}) interface{} {
	return value
}
func (gt *GraphQLInputObjectType) CoerceLiteral(value interface{}) interface{} {
	return value
}
func (gt *GraphQLInputObjectType) String() string {
	return gt.Name
}

func (gt *GraphQLInputObjectType) GetError() error {
	return gt.err
}
