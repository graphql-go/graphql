package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// EnumTypeDefinition implements Definition
type EnumTypeDefinition struct {
	Kind   string
	Loc    Location
	Name   Name
	Values []interface{}
}

func NewEnumTypeDefinition() *EnumTypeDefinition {
	return &EnumTypeDefinition{
		Kind: kinds.EnumTypeDefinition,
	}
}

func (def *EnumTypeDefinition) GetKind() string {
	return def.Kind
}

func (def *EnumTypeDefinition) GetLoc() Location {
	return def.Loc
}

func (def *EnumTypeDefinition) GetName() Name {
	return def.Name
}
