package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// ListType implements Node, Type
type ListType struct {
	Kind string
	Loc  *Location
	Type Type
}

func NewListType(v *ListType) *ListType {
	if v == nil {
		v = &ListType{}
	}
	return &ListType{
		Kind: kinds.ListType,
		Loc:  v.Loc,
		Type: v.Type,
	}
}

func (v *ListType) GetKind() string {
	return v.Kind
}

func (v *ListType) GetLoc() *Location {
	return v.Loc
}
