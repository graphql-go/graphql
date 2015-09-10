package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// EnumValue implements Node, Value
type EnumValue struct {
	Kind  string
	Loc   Location
	Value string
}

func NewEnumValue(v *EnumValue) *EnumValue {
	if v == nil {
		v = &EnumValue{}
	}
	return &EnumValue{
		Kind:  kinds.EnumValue,
		Loc:   v.Loc,
		Value: v.Value,
	}
}

func (v *EnumValue) GetKind() string {
	return v.Kind
}

func (v *EnumValue) GetLoc() Location {
	return v.Loc
}
