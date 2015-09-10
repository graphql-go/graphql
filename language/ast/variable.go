package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// Variable implements Node, Value
type Variable struct {
	Kind string
	Loc  Location
	Name *Name
}

func NewVariable(v *Variable) *Variable {
	if v == nil {
		v = &Variable{}
	}
	return &Variable{
		Kind: kinds.Variable,
		Loc:  v.Loc,
		Name: v.Name,
	}
}

func (v *Variable) GetKind() string {
	return v.Kind
}

func (v *Variable) GetLoc() Location {
	return v.Loc
}
