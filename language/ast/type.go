package ast

type Type interface {
	GetKind() string
}

// Ensure that all value types implements Value interface
var _ Type = (*NamedType)(nil)
var _ Type = (*ListType)(nil)
var _ Type = (*NonNullType)(nil)
