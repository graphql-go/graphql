package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

type Type interface {
	GetKind() string
	GetLoc() *Location
	String() string
}

// Ensure that all value types implements Value interface
var _ Type = (*NamedType)(nil)
var _ Type = (*ListType)(nil)
var _ Type = (*NonNullType)(nil)

// NamedType implements Node, Type
type NamedType struct {
	Kind string
	Loc  *Location
	Name *Name
}

func NewNamedType(t *NamedType) *NamedType {
	if t == nil {
		t = &NamedType{}
	}
	return &NamedType{
		Kind: kinds.NamedType,
		Loc:  t.Loc,
		Name: t.Name,
	}
}

func (t *NamedType) GetKind() string {
	return t.Kind
}

func (t *NamedType) GetLoc() *Location {
	return t.Loc
}

func (t *NamedType) String() string {
	return t.GetKind()
}

// ListType implements Node, Type
type ListType struct {
	Kind string
	Loc  *Location
	Type Type
}

func NewListType(t *ListType) *ListType {
	if t == nil {
		t = &ListType{}
	}
	return &ListType{
		Kind: kinds.ListType,
		Loc:  t.Loc,
		Type: t.Type,
	}
}

func (t *ListType) GetKind() string {
	return t.Kind
}

func (t *ListType) GetLoc() *Location {
	return t.Loc
}

func (t *ListType) String() string {
	return t.GetKind()
}

// NonNullType implements Node, Type
type NonNullType struct {
	Kind string
	Loc  *Location
	Type Type
}

func NewNonNullType(t *NonNullType) *NonNullType {
	if t == nil {
		t = &NonNullType{}
	}
	return &NonNullType{
		Kind: kinds.NonNullType,
		Loc:  t.Loc,
		Type: t.Type,
	}
}

func (t *NonNullType) GetKind() string {
	return t.Kind
}

func (t *NonNullType) GetLoc() *Location {
	return t.Loc
}

func (t *NonNullType) String() string {
	return t.GetKind()
}
