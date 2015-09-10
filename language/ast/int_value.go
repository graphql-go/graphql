package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// IntValue implements Node, Value
type IntValue struct {
	Kind  string
	Loc   *Location
	Value string
}

func NewIntValue(v *IntValue) *IntValue {
	if v == nil {
		v = &IntValue{}
	}
	return &IntValue{
		Kind:  kinds.IntValue,
		Loc:   v.Loc,
		Value: v.Value,
	}
}

func (v *IntValue) GetKind() string {
	return v.Kind
}

func (v *IntValue) GetLoc() *Location {
	return v.Loc
}
