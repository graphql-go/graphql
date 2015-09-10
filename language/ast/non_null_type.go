package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// NonNullType implements Node, Type
type NonNullType struct {
	Kind string
	Loc  *Location
	Type Type
}

func NewNonNullType(v *NonNullType) *NonNullType {
	if v == nil {
		v = &NonNullType{}
	}
	return &NonNullType{
		Kind: kinds.NonNullType,
		Loc:  v.Loc,
		Type: v.Type,
	}
}

func (v *NonNullType) GetKind() string {
	return v.Kind
}

func (v *NonNullType) GetLoc() *Location {
	return v.Loc
}
