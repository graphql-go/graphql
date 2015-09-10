package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// BooleanValue implements Node, Value
type BooleanValue struct {
	Kind  string
	Loc   Location
	Value bool
}

func NewBooleanValue(v *BooleanValue) *BooleanValue {
	if v == nil {
		v = &BooleanValue{}
	}
	return &BooleanValue{
		Kind:  kinds.BooleanValue,
		Loc:   v.Loc,
		Value: v.Value,
	}
}

func (v *BooleanValue) GetKind() string {
	return v.Kind
}

func (v *BooleanValue) GetLoc() Location {
	return v.Loc
}
