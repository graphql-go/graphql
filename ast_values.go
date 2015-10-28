package graphql

import (
	"github.com/chris-ramon/graphql/kinds"
)

type Value interface {
	GetValue() interface{}
	GetKind() string
	GetLoc() *AstLocation
}

// Ensure that all value types implements Value interface
var _ Value = (*AstVariable)(nil)
var _ Value = (*AstIntValue)(nil)
var _ Value = (*AstFloatValue)(nil)
var _ Value = (*AstStringValue)(nil)
var _ Value = (*AstBooleanValue)(nil)
var _ Value = (*AstEnumValue)(nil)
var _ Value = (*AstListValue)(nil)
var _ Value = (*AstObjectValue)(nil)

// AstVariable implements Node, Value
type AstVariable struct {
	Kind string
	Loc  *AstLocation
	Name *AstName
}

func NewAstVariable(v *AstVariable) *AstVariable {
	if v == nil {
		v = &AstVariable{}
	}
	return &AstVariable{
		Kind: kinds.Variable,
		Loc:  v.Loc,
		Name: v.Name,
	}
}

func (v *AstVariable) GetKind() string {
	return v.Kind
}

func (v *AstVariable) GetLoc() *AstLocation {
	return v.Loc
}

// GetValue alias to AstVariable.GetName()
func (v *AstVariable) GetValue() interface{} {
	return v.GetName()
}

func (v *AstVariable) GetName() interface{} {
	return v.Name
}

// AstIntValue implements Node, Value
type AstIntValue struct {
	Kind  string
	Loc   *AstLocation
	Value string
}

func NewAstIntValue(v *AstIntValue) *AstIntValue {
	if v == nil {
		v = &AstIntValue{}
	}
	return &AstIntValue{
		Kind:  kinds.IntValue,
		Loc:   v.Loc,
		Value: v.Value,
	}
}

func (v *AstIntValue) GetKind() string {
	return v.Kind
}

func (v *AstIntValue) GetLoc() *AstLocation {
	return v.Loc
}

func (v *AstIntValue) GetValue() interface{} {
	return v.Value
}

// AstFloatValue implements Node, Value
type AstFloatValue struct {
	Kind  string
	Loc   *AstLocation
	Value string
}

func NewAstFloatValue(v *AstFloatValue) *AstFloatValue {
	if v == nil {
		v = &AstFloatValue{}
	}
	return &AstFloatValue{
		Kind:  kinds.FloatValue,
		Loc:   v.Loc,
		Value: v.Value,
	}
}

func (v *AstFloatValue) GetKind() string {
	return v.Kind
}

func (v *AstFloatValue) GetLoc() *AstLocation {
	return v.Loc
}

func (v *AstFloatValue) GetValue() interface{} {
	return v.Value
}

// AstStringValue implements Node, Value
type AstStringValue struct {
	Kind  string
	Loc   *AstLocation
	Value string
}

func NewAstStringValue(v *AstStringValue) *AstStringValue {
	if v == nil {
		v = &AstStringValue{}
	}
	return &AstStringValue{
		Kind:  kinds.StringValue,
		Loc:   v.Loc,
		Value: v.Value,
	}
}

func (v *AstStringValue) GetKind() string {
	return v.Kind
}

func (v *AstStringValue) GetLoc() *AstLocation {
	return v.Loc
}

func (v *AstStringValue) GetValue() interface{} {
	return v.Value
}

// AstBooleanValue implements Node, Value
type AstBooleanValue struct {
	Kind  string
	Loc   *AstLocation
	Value bool
}

func NewAstBooleanValue(v *AstBooleanValue) *AstBooleanValue {
	if v == nil {
		v = &AstBooleanValue{}
	}
	return &AstBooleanValue{
		Kind:  kinds.BooleanValue,
		Loc:   v.Loc,
		Value: v.Value,
	}
}

func (v *AstBooleanValue) GetKind() string {
	return v.Kind
}

func (v *AstBooleanValue) GetLoc() *AstLocation {
	return v.Loc
}

func (v *AstBooleanValue) GetValue() interface{} {
	return v.Value
}

// AstEnumValue implements Node, Value
type AstEnumValue struct {
	Kind  string
	Loc   *AstLocation
	Value string
}

func NewAstEnumValue(v *AstEnumValue) *AstEnumValue {
	if v == nil {
		v = &AstEnumValue{}
	}
	return &AstEnumValue{
		Kind:  kinds.EnumValue,
		Loc:   v.Loc,
		Value: v.Value,
	}
}

func (v *AstEnumValue) GetKind() string {
	return v.Kind
}

func (v *AstEnumValue) GetLoc() *AstLocation {
	return v.Loc
}

func (v *AstEnumValue) GetValue() interface{} {
	return v.Value
}

// AstListValue implements Node, Value
type AstListValue struct {
	Kind   string
	Loc    *AstLocation
	Values []Value
}

func NewAstListValue(v *AstListValue) *AstListValue {
	if v == nil {
		v = &AstListValue{}
	}
	return &AstListValue{
		Kind:   kinds.ListValue,
		Loc:    v.Loc,
		Values: v.Values,
	}
}

func (v *AstListValue) GetKind() string {
	return v.Kind
}

func (v *AstListValue) GetLoc() *AstLocation {
	return v.Loc
}

// GetValue alias to AstListValue.GetValues()
func (v *AstListValue) GetValue() interface{} {
	return v.GetValues()
}

func (v *AstListValue) GetValues() interface{} {
	// TODO: verify AstObjectValue.GetValue()
	return v.Values
}

// AstObjectValue implements Node, Value
type AstObjectValue struct {
	Kind   string
	Loc    *AstLocation
	Fields []*AstObjectField
}

func NewAstObjectValue(v *AstObjectValue) *AstObjectValue {
	if v == nil {
		v = &AstObjectValue{}
	}
	return &AstObjectValue{
		Kind:   kinds.ObjectValue,
		Loc:    v.Loc,
		Fields: v.Fields,
	}
}

func (v *AstObjectValue) GetKind() string {
	return v.Kind
}

func (v *AstObjectValue) GetLoc() *AstLocation {
	return v.Loc
}

func (v *AstObjectValue) GetValue() interface{} {
	// TODO: verify AstObjectValue.GetValue()
	return v.Fields
}

// AstObjectField implements Node, Value
type AstObjectField struct {
	Kind  string
	Name  *AstName
	Loc   *AstLocation
	Value Value
}

func NewAstObjectField(f *AstObjectField) *AstObjectField {
	if f == nil {
		f = &AstObjectField{}
	}
	return &AstObjectField{
		Kind:  kinds.ObjectField,
		Loc:   f.Loc,
		Name:  f.Name,
		Value: f.Value,
	}
}

func (f *AstObjectField) GetKind() string {
	return f.Kind
}

func (f *AstObjectField) GetLoc() *AstLocation {
	return f.Loc
}

func (f *AstObjectField) GetValue() interface{} {
	return f.Value
}
