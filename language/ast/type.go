package ast

type Type interface {
}

type NamedType struct {
	Kind string
	Loc  Location
	Name *Name
}

type ListType struct {
	Kind string
	Loc  Location
	Type Type
}

type NonNullType struct {
	Kind string
	Loc  Location
	Type Type
}

// Ensure that all value types implements Value interface
var _ Type = (*NamedType)(nil)
var _ Type = (*ListType)(nil)
var _ Type = (*NonNullType)(nil)
