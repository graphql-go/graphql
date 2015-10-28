package graphql

import (
	"github.com/chris-ramon/graphql/kinds"
)

// Ensure that all typeDefinition types implements Definition interface
var _ Definition = (*AstObjectDefinition)(nil)
var _ Definition = (*AstInterfaceDefinition)(nil)
var _ Definition = (*AstUnionDefinition)(nil)
var _ Definition = (*AstScalarDefinition)(nil)
var _ Definition = (*AstEnumDefinition)(nil)
var _ Definition = (*AstInputObjectDefinition)(nil)
var _ Definition = (*AstTypeExtensionDefinition)(nil)

// AstObjectDefinition implements Node, Definition
type AstObjectDefinition struct {
	Kind       string
	Loc        *AstLocation
	Name       *AstName
	Interfaces []*AstNamed
	Fields     []*AstFieldDefinition
}

func NewAstObjectDefinition(def *AstObjectDefinition) *AstObjectDefinition {
	if def == nil {
		def = &AstObjectDefinition{}
	}
	return &AstObjectDefinition{
		Kind:       kinds.ObjectDefinition,
		Loc:        def.Loc,
		Name:       def.Name,
		Interfaces: def.Interfaces,
		Fields:     def.Fields,
	}
}

func (def *AstObjectDefinition) GetKind() string {
	return def.Kind
}

func (def *AstObjectDefinition) GetLoc() *AstLocation {
	return def.Loc
}

func (def *AstObjectDefinition) GetName() *AstName {
	return def.Name
}

func (def *AstObjectDefinition) GetVariableDefinitions() []*AstVariableDefinition {
	return []*AstVariableDefinition{}
}

func (def *AstObjectDefinition) GetSelectionSet() *AstSelectionSet {
	return &AstSelectionSet{}
}

func (def *AstObjectDefinition) GetOperation() string {
	return ""
}

// AstFieldDefinition implements Node
type AstFieldDefinition struct {
	Kind      string
	Loc       *AstLocation
	Name      *AstName
	Arguments []*AstInputValueDefinition
	Type      AstType
}

func NewAstFieldDefinition(def *AstFieldDefinition) *AstFieldDefinition {
	if def == nil {
		def = &AstFieldDefinition{}
	}
	return &AstFieldDefinition{
		Kind:      kinds.FieldDefinition,
		Loc:       def.Loc,
		Name:      def.Name,
		Arguments: def.Arguments,
		Type:      def.Type,
	}
}

func (def *AstFieldDefinition) GetKind() string {
	return def.Kind
}

func (def *AstFieldDefinition) GetLoc() *AstLocation {
	return def.Loc
}

// AstInputValueDefinition implements Node
type AstInputValueDefinition struct {
	Kind         string
	Loc          *AstLocation
	Name         *AstName
	Type         AstType
	DefaultValue Value
}

func NewAstInputValueDefinition(def *AstInputValueDefinition) *AstInputValueDefinition {
	if def == nil {
		def = &AstInputValueDefinition{}
	}
	return &AstInputValueDefinition{
		Kind:         kinds.InputValueDefinition,
		Loc:          def.Loc,
		Name:         def.Name,
		Type:         def.Type,
		DefaultValue: def.DefaultValue,
	}
}

func (def *AstInputValueDefinition) GetKind() string {
	return def.Kind
}

func (def *AstInputValueDefinition) GetLoc() *AstLocation {
	return def.Loc
}

// AstInterfaceDefinition implements Node, Definition
type AstInterfaceDefinition struct {
	Kind   string
	Loc    *AstLocation
	Name   *AstName
	Fields []*AstFieldDefinition
}

func NewAstInterfaceDefinition(def *AstInterfaceDefinition) *AstInterfaceDefinition {
	if def == nil {
		def = &AstInterfaceDefinition{}
	}
	return &AstInterfaceDefinition{
		Kind:   kinds.InterfaceDefinition,
		Loc:    def.Loc,
		Name:   def.Name,
		Fields: def.Fields,
	}
}

func (def *AstInterfaceDefinition) GetKind() string {
	return def.Kind
}

func (def *AstInterfaceDefinition) GetLoc() *AstLocation {
	return def.Loc
}

func (def *AstInterfaceDefinition) GetName() *AstName {
	return def.Name
}

func (def *AstInterfaceDefinition) GetVariableDefinitions() []*AstVariableDefinition {
	return []*AstVariableDefinition{}
}

func (def *AstInterfaceDefinition) GetSelectionSet() *AstSelectionSet {
	return &AstSelectionSet{}
}

func (def *AstInterfaceDefinition) GetOperation() string {
	return ""
}

// AstUnionDefinition implements Node, Definition
type AstUnionDefinition struct {
	Kind  string
	Loc   *AstLocation
	Name  *AstName
	Types []*AstNamed
}

func NewAstUnionDefinition(def *AstUnionDefinition) *AstUnionDefinition {
	if def == nil {
		def = &AstUnionDefinition{}
	}
	return &AstUnionDefinition{
		Kind:  kinds.UnionDefinition,
		Loc:   def.Loc,
		Name:  def.Name,
		Types: def.Types,
	}
}

func (def *AstUnionDefinition) GetKind() string {
	return def.Kind
}

func (def *AstUnionDefinition) GetLoc() *AstLocation {
	return def.Loc
}

func (def *AstUnionDefinition) GetName() *AstName {
	return def.Name
}

func (def *AstUnionDefinition) GetVariableDefinitions() []*AstVariableDefinition {
	return []*AstVariableDefinition{}
}

func (def *AstUnionDefinition) GetSelectionSet() *AstSelectionSet {
	return &AstSelectionSet{}
}

func (def *AstUnionDefinition) GetOperation() string {
	return ""
}

// AstScalarDefinition implements Node, Definition
type AstScalarDefinition struct {
	Kind string
	Loc  *AstLocation
	Name *AstName
}

func NewAstScalarDefinition(def *AstScalarDefinition) *AstScalarDefinition {
	if def == nil {
		def = &AstScalarDefinition{}
	}
	return &AstScalarDefinition{
		Kind: kinds.ScalarDefinition,
		Loc:  def.Loc,
		Name: def.Name,
	}
}

func (def *AstScalarDefinition) GetKind() string {
	return def.Kind
}

func (def *AstScalarDefinition) GetLoc() *AstLocation {
	return def.Loc
}

func (def *AstScalarDefinition) GetName() *AstName {
	return def.Name
}

func (def *AstScalarDefinition) GetVariableDefinitions() []*AstVariableDefinition {
	return []*AstVariableDefinition{}
}

func (def *AstScalarDefinition) GetSelectionSet() *AstSelectionSet {
	return &AstSelectionSet{}
}

func (def *AstScalarDefinition) GetOperation() string {
	return ""
}

// AstEnumDefinition implements Node, Definition
type AstEnumDefinition struct {
	Kind   string
	Loc    *AstLocation
	Name   *AstName
	Values []*AstEnumValueDefinition
}

func NewAstEnumDefinition(def *AstEnumDefinition) *AstEnumDefinition {
	if def == nil {
		def = &AstEnumDefinition{}
	}
	return &AstEnumDefinition{
		Kind:   kinds.EnumDefinition,
		Loc:    def.Loc,
		Name:   def.Name,
		Values: def.Values,
	}
}

func (def *AstEnumDefinition) GetKind() string {
	return def.Kind
}

func (def *AstEnumDefinition) GetLoc() *AstLocation {
	return def.Loc
}

func (def *AstEnumDefinition) GetName() *AstName {
	return def.Name
}

func (def *AstEnumDefinition) GetVariableDefinitions() []*AstVariableDefinition {
	return []*AstVariableDefinition{}
}

func (def *AstEnumDefinition) GetSelectionSet() *AstSelectionSet {
	return &AstSelectionSet{}
}

func (def *AstEnumDefinition) GetOperation() string {
	return ""
}

// EnumValueDefinition implements Node, Definition
type AstEnumValueDefinition struct {
	Kind string
	Loc  *AstLocation
	Name *AstName
}

func NewAstEnumValueDefinition(def *AstEnumValueDefinition) *AstEnumValueDefinition {
	if def == nil {
		def = &AstEnumValueDefinition{}
	}
	return &AstEnumValueDefinition{
		Kind: kinds.EnumValueDefinition,
		Loc:  def.Loc,
		Name: def.Name,
	}
}

func (def *AstEnumValueDefinition) GetKind() string {
	return def.Kind
}

func (def *AstEnumValueDefinition) GetLoc() *AstLocation {
	return def.Loc
}

// AstInputObjectDefinition implements Node, Definition
type AstInputObjectDefinition struct {
	Kind   string
	Loc    *AstLocation
	Name   *AstName
	Fields []*AstInputValueDefinition
}

func NewAstInputObjectDefinition(def *AstInputObjectDefinition) *AstInputObjectDefinition {
	if def == nil {
		def = &AstInputObjectDefinition{}
	}
	return &AstInputObjectDefinition{
		Kind:   kinds.InputObjectDefinition,
		Loc:    def.Loc,
		Name:   def.Name,
		Fields: def.Fields,
	}
}

func (def *AstInputObjectDefinition) GetKind() string {
	return def.Kind
}

func (def *AstInputObjectDefinition) GetLoc() *AstLocation {
	return def.Loc
}

func (def *AstInputObjectDefinition) GetName() *AstName {
	return def.Name
}

func (def *AstInputObjectDefinition) GetVariableDefinitions() []*AstVariableDefinition {
	return []*AstVariableDefinition{}
}

func (def *AstInputObjectDefinition) GetSelectionSet() *AstSelectionSet {
	return &AstSelectionSet{}
}

func (def *AstInputObjectDefinition) GetOperation() string {
	return ""
}

// TypeExtensionDefinition implements Node, Definition
type AstTypeExtensionDefinition struct {
	Kind       string
	Loc        *AstLocation
	Definition *AstObjectDefinition
}

func NewAstTypeExtensionDefinition(def *AstTypeExtensionDefinition) *AstTypeExtensionDefinition {
	if def == nil {
		def = &AstTypeExtensionDefinition{}
	}
	return &AstTypeExtensionDefinition{
		Kind:       kinds.TypeExtensionDefinition,
		Loc:        def.Loc,
		Definition: def.Definition,
	}
}

func (def *AstTypeExtensionDefinition) GetKind() string {
	return def.Kind
}

func (def *AstTypeExtensionDefinition) GetLoc() *AstLocation {
	return def.Loc
}

func (def *AstTypeExtensionDefinition) GetVariableDefinitions() []*AstVariableDefinition {
	return []*AstVariableDefinition{}
}

func (def *AstTypeExtensionDefinition) GetSelectionSet() *AstSelectionSet {
	return &AstSelectionSet{}
}

func (def *AstTypeExtensionDefinition) GetOperation() string {
	return ""
}
