package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// ListValue implements Node, Value
type ListValue struct {
	Kind   string
	Loc    *Location
	Values []Value
}

func NewListValue(v *ListValue) *ListValue {
	if v == nil {
		v = &ListValue{}
	}
	return &ListValue{
		Kind:   kinds.ListValue,
		Loc:    v.Loc,
		Values: v.Values,
	}
}

func (v *ListValue) GetKind() string {
	return v.Kind
}

func (v *ListValue) GetLoc() *Location {
	return v.Loc
}
