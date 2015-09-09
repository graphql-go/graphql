package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// InputObjectTypeDefinition implements Definition
type InputObjectTypeDefinition struct {
	Kind   string
	Loc    Location
	Name   Name
	Fields []interface{}
}

func NewInputObjectTypeDefinition() *InputObjectTypeDefinition {
	return &InputObjectTypeDefinition{
		Kind: kinds.InputObjectTypeDefinition,
	}
}

func (def *InputObjectTypeDefinition) GetKind() string {
	return def.Kind
}

func (def *InputObjectTypeDefinition) GetLoc() Location {
	return def.Loc
}

func (def *InputObjectTypeDefinition) GetName() Name {
	return def.Name
}
