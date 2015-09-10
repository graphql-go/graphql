package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// StringValue implements Node, Value
type StringValue struct {
	Kind  string
	Loc   Location
	Value string
}

func NewStringValue(v *StringValue) *StringValue {
	if v == nil {
		v = &StringValue{}
	}
	return &StringValue{
		Kind:  kinds.StringValue,
		Loc:   v.Loc,
		Value: v.Value,
	}
}

func (v *StringValue) GetKind() string {
	return v.Kind
}

func (v *StringValue) GetLoc() Location {
	return v.Loc
}
