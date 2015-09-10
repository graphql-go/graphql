package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// FloatValue implements Node, Value
type FloatValue struct {
	Kind  string
	Loc   Location
	Value string
}

func NewFloatValue(v *FloatValue) *FloatValue {
	if v == nil {
		v = &FloatValue{}
	}
	return &FloatValue{
		Kind:  kinds.FloatValue,
		Loc:   v.Loc,
		Value: v.Value,
	}
}

func (v *FloatValue) GetKind() string {
	return v.Kind
}

func (v *FloatValue) GetLoc() Location {
	return v.Loc
}
