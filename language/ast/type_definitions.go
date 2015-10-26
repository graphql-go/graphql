package ast

import (
	"github.com/chris-ramon/graphql/language/kinds"
)

// TypeDefinition implements Definition
type TypeDefinition interface {
	// TODO: determine the minimal set of interface for `TypeDefinition`
	GetOperation() string
	GetVariableDefinitions() []*VariableDefinition
	GetSelectionSet() *SelectionSet
}

// Ensure that all typeDefinition types implements TypeDefinition interface
var _ TypeDefinition = (*ObjectTypeDefinition)(nil)
var _ TypeDefinition = (*InterfaceTypeDefinition)(nil)
var _ TypeDefinition = (*UnionTypeDefinition)(nil)
var _ TypeDefinition = (*ScalarTypeDefinition)(nil)
var _ TypeDefinition = (*EnumTypeDefinition)(nil)
var _ TypeDefinition = (*InputObjectTypeDefinition)(nil)
var _ TypeDefinition = (*TypeExtensionDefinition)(nil)

// ObjectTypeDefinition implements Node, TypeDefinition
type ObjectTypeDefinition struct {
	Kind       string
	Loc        *Location
	Name       *Name
	Interfaces []*NamedType
	Fields     []*FieldDefinition
}

func NewObjectTypeDefinition(def *ObjectTypeDefinition) *ObjectTypeDefinition {
	if def == nil {
		def = &ObjectTypeDefinition{}
	}
	return &ObjectTypeDefinition{
		Kind:       kinds.ObjectTypeDefinition,
		Loc:        def.Loc,
		Name:       def.Name,
		Interfaces: def.Interfaces,
		Fields:     def.Fields,
	}
}

func (def *ObjectTypeDefinition) GetKind() string {
	return def.Kind
}

func (def *ObjectTypeDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *ObjectTypeDefinition) GetName() *Name {
	return def.Name
}

func (def *ObjectTypeDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}

func (def *ObjectTypeDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}

func (def *ObjectTypeDefinition) GetOperation() string {
	return ""
}

// FieldDefinition implements Node
type FieldDefinition struct {
	Kind      string
	Loc       *Location
	Name      *Name
	Arguments []*InputValueDefinition
	Type      Type
}

func NewFieldDefinition(def *FieldDefinition) *FieldDefinition {
	if def == nil {
		def = &FieldDefinition{}
	}
	return &FieldDefinition{
		Kind:      kinds.FieldDefinition,
		Loc:       def.Loc,
		Name:      def.Name,
		Arguments: def.Arguments,
		Type:      def.Type,
	}
}

func (def *FieldDefinition) GetKind() string {
	return def.Kind
}

func (def *FieldDefinition) GetLoc() *Location {
	return def.Loc
}

// InputValueDefinition implements Node
type InputValueDefinition struct {
	Kind         string
	Loc          *Location
	Name         *Name
	Type         Type
	DefaultValue Value
}

func NewInputValueDefinition(def *InputValueDefinition) *InputValueDefinition {
	if def == nil {
		def = &InputValueDefinition{}
	}
	return &InputValueDefinition{
		Kind:         kinds.InputValueDefinition,
		Loc:          def.Loc,
		Name:         def.Name,
		Type:         def.Type,
		DefaultValue: def.DefaultValue,
	}
}

func (def *InputValueDefinition) GetKind() string {
	return def.Kind
}

func (def *InputValueDefinition) GetLoc() *Location {
	return def.Loc
}

// InterfaceTypeDefinition implements Node, TypeDefinition
type InterfaceTypeDefinition struct {
	Kind   string
	Loc    *Location
	Name   *Name
	Fields []*FieldDefinition
}

func NewInterfaceTypeDefinition(def *InterfaceTypeDefinition) *InterfaceTypeDefinition {
	if def == nil {
		def = &InterfaceTypeDefinition{}
	}
	return &InterfaceTypeDefinition{
		Kind:   kinds.InterfaceTypeDefinition,
		Loc:    def.Loc,
		Name:   def.Name,
		Fields: def.Fields,
	}
}

func (def *InterfaceTypeDefinition) GetKind() string {
	return def.Kind
}

func (def *InterfaceTypeDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *InterfaceTypeDefinition) GetName() *Name {
	return def.Name
}

func (def *InterfaceTypeDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}

func (def *InterfaceTypeDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}

func (def *InterfaceTypeDefinition) GetOperation() string {
	return ""
}

// UnionTypeDefinition implements Node, TypeDefinition
type UnionTypeDefinition struct {
	Kind  string
	Loc   *Location
	Name  *Name
	Types []*NamedType
}

func NewUnionTypeDefinition(def *UnionTypeDefinition) *UnionTypeDefinition {
	if def == nil {
		def = &UnionTypeDefinition{}
	}
	return &UnionTypeDefinition{
		Kind:  kinds.UnionTypeDefinition,
		Loc:   def.Loc,
		Name:  def.Name,
		Types: def.Types,
	}
}

func (def *UnionTypeDefinition) GetKind() string {
	return def.Kind
}

func (def *UnionTypeDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *UnionTypeDefinition) GetName() *Name {
	return def.Name
}

func (def *UnionTypeDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}

func (def *UnionTypeDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}

func (def *UnionTypeDefinition) GetOperation() string {
	return ""
}

// ScalarTypeDefinition implements Node, TypeDefinition
type ScalarTypeDefinition struct {
	Kind string
	Loc  *Location
	Name *Name
}

func NewScalarTypeDefinition(def *ScalarTypeDefinition) *ScalarTypeDefinition {
	if def == nil {
		def = &ScalarTypeDefinition{}
	}
	return &ScalarTypeDefinition{
		Kind: kinds.ScalarTypeDefinition,
		Loc:  def.Loc,
		Name: def.Name,
	}
}

func (def *ScalarTypeDefinition) GetKind() string {
	return def.Kind
}

func (def *ScalarTypeDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *ScalarTypeDefinition) GetName() *Name {
	return def.Name
}

func (def *ScalarTypeDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}

func (def *ScalarTypeDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}

func (def *ScalarTypeDefinition) GetOperation() string {
	return ""
}

// EnumTypeDefinition implements Node, TypeDefinition
type EnumTypeDefinition struct {
	Kind   string
	Loc    *Location
	Name   *Name
	Values []*EnumValueDefinition
}

func NewEnumTypeDefinition(def *EnumTypeDefinition) *EnumTypeDefinition {
	if def == nil {
		def = &EnumTypeDefinition{}
	}
	return &EnumTypeDefinition{
		Kind:   kinds.EnumTypeDefinition,
		Loc:    def.Loc,
		Name:   def.Name,
		Values: def.Values,
	}
}

func (def *EnumTypeDefinition) GetKind() string {
	return def.Kind
}

func (def *EnumTypeDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *EnumTypeDefinition) GetName() *Name {
	return def.Name
}

func (def *EnumTypeDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}

func (def *EnumTypeDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}

func (def *EnumTypeDefinition) GetOperation() string {
	return ""
}

// EnumValueDefinition implements Node, TypeDefinition
type EnumValueDefinition struct {
	Kind string
	Loc  *Location
	Name *Name
}

func NewEnumValueDefinition(def *EnumValueDefinition) *EnumValueDefinition {
	if def == nil {
		def = &EnumValueDefinition{}
	}
	return &EnumValueDefinition{
		Kind: kinds.EnumValueDefinition,
		Loc:  def.Loc,
		Name: def.Name,
	}
}

func (def *EnumValueDefinition) GetKind() string {
	return def.Kind
}

func (def *EnumValueDefinition) GetLoc() *Location {
	return def.Loc
}

// InputObjectTypeDefinition implements Node, TypeDefinition
type InputObjectTypeDefinition struct {
	Kind   string
	Loc    *Location
	Name   *Name
	Fields []*InputValueDefinition
}

func NewInputObjectTypeDefinition(def *InputObjectTypeDefinition) *InputObjectTypeDefinition {
	if def == nil {
		def = &InputObjectTypeDefinition{}
	}
	return &InputObjectTypeDefinition{
		Kind:   kinds.InputObjectTypeDefinition,
		Loc:    def.Loc,
		Name:   def.Name,
		Fields: def.Fields,
	}
}

func (def *InputObjectTypeDefinition) GetKind() string {
	return def.Kind
}

func (def *InputObjectTypeDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *InputObjectTypeDefinition) GetName() *Name {
	return def.Name
}

func (def *InputObjectTypeDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}

func (def *InputObjectTypeDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}

func (def *InputObjectTypeDefinition) GetOperation() string {
	return ""
}

// TypeExtensionDefinition implements Node, TypeDefinition
type TypeExtensionDefinition struct {
	Kind       string
	Loc        *Location
	Definition *ObjectTypeDefinition
}

func NewTypeExtensionDefinition(def *TypeExtensionDefinition) *TypeExtensionDefinition {
	if def == nil {
		def = &TypeExtensionDefinition{}
	}
	return &TypeExtensionDefinition{
		Kind:       kinds.TypeExtensionDefinition,
		Loc:        def.Loc,
		Definition: def.Definition,
	}
}

func (def *TypeExtensionDefinition) GetKind() string {
	return def.Kind
}

func (def *TypeExtensionDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *TypeExtensionDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}

func (def *TypeExtensionDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}

func (def *TypeExtensionDefinition) GetOperation() string {
	return ""
}
