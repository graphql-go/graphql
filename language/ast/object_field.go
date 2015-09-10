package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// ObjectField implements Node, Value
type ObjectField struct {
	Kind  string
	Name  *Name
	Loc   *Location
	Value Value
}

func NewObjectField(f *ObjectField) *ObjectField {
	if f == nil {
		f = &ObjectField{}
	}
	return &ObjectField{
		Kind:  kinds.ObjectField,
		Loc:   f.Loc,
		Name:  f.Name,
		Value: f.Value,
	}
}

func (f *ObjectField) GetKind() string {
	return f.Kind
}

func (f *ObjectField) GetLoc() *Location {
	return f.Loc
}
