package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// NamedType implements Node, Type
type NamedType struct {
	Kind string
	Loc  *Location
	Name *Name
}

func NewNamedType(v *NamedType) *NamedType {
	if v == nil {
		v = &NamedType{}
	}
	return &NamedType{
		Kind: kinds.NamedType,
		Loc:  v.Loc,
		Name: v.Name,
	}
}

func (v *NamedType) GetKind() string {
	return v.Kind
}

func (v *NamedType) GetLoc() *Location {
	return v.Loc
}
