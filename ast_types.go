package graphql

import (
	"github.com/chris-ramon/graphql/kinds"
)

type AstType interface {
	GetKind() string
	GetLoc() *AstLocation
	String() string
}

// Ensure that all value types implements Value interface
var _ AstType = (*AstNamed)(nil)
var _ AstType = (*AstList)(nil)
var _ AstType = (*AstNonNull)(nil)

// AstNamed implements Node, Type
type AstNamed struct {
	Kind string
	Loc  *AstLocation
	Name *AstName
}

func NewAstNamed(t *AstNamed) *AstNamed {
	if t == nil {
		t = &AstNamed{}
	}
	return &AstNamed{
		Kind: kinds.Named,
		Loc:  t.Loc,
		Name: t.Name,
	}
}

func (t *AstNamed) GetKind() string {
	return t.Kind
}

func (t *AstNamed) GetLoc() *AstLocation {
	return t.Loc
}

func (t *AstNamed) String() string {
	return t.GetKind()
}

// List implements Node, Type
type AstList struct {
	Kind string
	Loc  *AstLocation
	Type AstType
}

func NewAstList(t *AstList) *AstList {
	if t == nil {
		t = &AstList{}
	}
	return &AstList{
		Kind: kinds.List,
		Loc:  t.Loc,
		Type: t.Type,
	}
}

func (t *AstList) GetKind() string {
	return t.Kind
}

func (t *AstList) GetLoc() *AstLocation {
	return t.Loc
}

func (t *AstList) String() string {
	return t.GetKind()
}

// AstNonNull implements Node, Type
type AstNonNull struct {
	Kind string
	Loc  *AstLocation
	Type AstType
}

func NewAstNonNull(t *AstNonNull) *AstNonNull {
	if t == nil {
		t = &AstNonNull{}
	}
	return &AstNonNull{
		Kind: kinds.NonNull,
		Loc:  t.Loc,
		Type: t.Type,
	}
}

func (t *AstNonNull) GetKind() string {
	return t.Kind
}

func (t *AstNonNull) GetLoc() *AstLocation {
	return t.Loc
}

func (t *AstNonNull) String() string {
	return t.GetKind()
}
