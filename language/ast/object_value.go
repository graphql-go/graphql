package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// ObjectValue implements Node, Value
type ObjectValue struct {
	Kind   string
	Loc    *Location
	Fields []*ObjectField
}

func NewObjectValue(v *ObjectValue) *ObjectValue {
	if v == nil {
		v = &ObjectValue{}
	}
	return &ObjectValue{
		Kind:   kinds.ObjectValue,
		Loc:    v.Loc,
		Fields: v.Fields,
	}
}

func (v *ObjectValue) GetKind() string {
	return v.Kind
}

func (v *ObjectValue) GetLoc() *Location {
	return v.Loc
}
