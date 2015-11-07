package ast

import (
	"github.com/graphql-go/graphql/language/kinds"
)

// Ensure that all typeDefinition types implements Definition interface
var _ Definition = (*ObjectDefinition)(nil)
var _ Definition = (*InterfaceDefinition)(nil)
var _ Definition = (*UnionDefinition)(nil)
var _ Definition = (*ScalarDefinition)(nil)
var _ Definition = (*EnumDefinition)(nil)
var _ Definition = (*InputObjectDefinition)(nil)
var _ Definition = (*TypeExtensionDefinition)(nil)

// ObjectDefinition implements Node, Definition
type ObjectDefinition struct {
	Kind       string
	Loc        *Location
	Name       *Name
	Interfaces []*Named
	Fields     []*FieldDefinition
}

func NewObjectDefinition(def *ObjectDefinition) *ObjectDefinition {
	if def == nil {
		def = &ObjectDefinition{}
	}
	return &ObjectDefinition{
		Kind:       kinds.ObjectDefinition,
		Loc:        def.Loc,
		Name:       def.Name,
		Interfaces: def.Interfaces,
		Fields:     def.Fields,
	}
}

func (def *ObjectDefinition) GetKind() string {
	return def.Kind
}

func (def *ObjectDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *ObjectDefinition) GetName() *Name {
	return def.Name
}

func (def *ObjectDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}

func (def *ObjectDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}

func (def *ObjectDefinition) GetOperation() string {
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

// InterfaceDefinition implements Node, Definition
type InterfaceDefinition struct {
	Kind   string
	Loc    *Location
	Name   *Name
	Fields []*FieldDefinition
}

func NewInterfaceDefinition(def *InterfaceDefinition) *InterfaceDefinition {
	if def == nil {
		def = &InterfaceDefinition{}
	}
	return &InterfaceDefinition{
		Kind:   kinds.InterfaceDefinition,
		Loc:    def.Loc,
		Name:   def.Name,
		Fields: def.Fields,
	}
}

func (def *InterfaceDefinition) GetKind() string {
	return def.Kind
}

func (def *InterfaceDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *InterfaceDefinition) GetName() *Name {
	return def.Name
}

func (def *InterfaceDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}

func (def *InterfaceDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}

func (def *InterfaceDefinition) GetOperation() string {
	return ""
}

// UnionDefinition implements Node, Definition
type UnionDefinition struct {
	Kind  string
	Loc   *Location
	Name  *Name
	Types []*Named
}

func NewUnionDefinition(def *UnionDefinition) *UnionDefinition {
	if def == nil {
		def = &UnionDefinition{}
	}
	return &UnionDefinition{
		Kind:  kinds.UnionDefinition,
		Loc:   def.Loc,
		Name:  def.Name,
		Types: def.Types,
	}
}

func (def *UnionDefinition) GetKind() string {
	return def.Kind
}

func (def *UnionDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *UnionDefinition) GetName() *Name {
	return def.Name
}

func (def *UnionDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}

func (def *UnionDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}

func (def *UnionDefinition) GetOperation() string {
	return ""
}

// ScalarDefinition implements Node, Definition
type ScalarDefinition struct {
	Kind string
	Loc  *Location
	Name *Name
}

func NewScalarDefinition(def *ScalarDefinition) *ScalarDefinition {
	if def == nil {
		def = &ScalarDefinition{}
	}
	return &ScalarDefinition{
		Kind: kinds.ScalarDefinition,
		Loc:  def.Loc,
		Name: def.Name,
	}
}

func (def *ScalarDefinition) GetKind() string {
	return def.Kind
}

func (def *ScalarDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *ScalarDefinition) GetName() *Name {
	return def.Name
}

func (def *ScalarDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}

func (def *ScalarDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}

func (def *ScalarDefinition) GetOperation() string {
	return ""
}

// EnumDefinition implements Node, Definition
type EnumDefinition struct {
	Kind   string
	Loc    *Location
	Name   *Name
	Values []*EnumValueDefinition
}

func NewEnumDefinition(def *EnumDefinition) *EnumDefinition {
	if def == nil {
		def = &EnumDefinition{}
	}
	return &EnumDefinition{
		Kind:   kinds.EnumDefinition,
		Loc:    def.Loc,
		Name:   def.Name,
		Values: def.Values,
	}
}

func (def *EnumDefinition) GetKind() string {
	return def.Kind
}

func (def *EnumDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *EnumDefinition) GetName() *Name {
	return def.Name
}

func (def *EnumDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}

func (def *EnumDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}

func (def *EnumDefinition) GetOperation() string {
	return ""
}

// EnumValueDefinition implements Node, Definition
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

// InputObjectDefinition implements Node, Definition
type InputObjectDefinition struct {
	Kind   string
	Loc    *Location
	Name   *Name
	Fields []*InputValueDefinition
}

func NewInputObjectDefinition(def *InputObjectDefinition) *InputObjectDefinition {
	if def == nil {
		def = &InputObjectDefinition{}
	}
	return &InputObjectDefinition{
		Kind:   kinds.InputObjectDefinition,
		Loc:    def.Loc,
		Name:   def.Name,
		Fields: def.Fields,
	}
}

func (def *InputObjectDefinition) GetKind() string {
	return def.Kind
}

func (def *InputObjectDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *InputObjectDefinition) GetName() *Name {
	return def.Name
}

func (def *InputObjectDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}

func (def *InputObjectDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}

func (def *InputObjectDefinition) GetOperation() string {
	return ""
}

// TypeExtensionDefinition implements Node, Definition
type TypeExtensionDefinition struct {
	Kind       string
	Loc        *Location
	Definition *ObjectDefinition
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
