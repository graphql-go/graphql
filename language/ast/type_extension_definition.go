package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// TypeExtensionDefinition implements Definition
type TypeExtensionDefinition struct {
	Kind       string
	Loc        Location
	Definition ObjectTypeDefinition
}

func NewTypeExtensionDefinition() *TypeExtensionDefinition {
	return &TypeExtensionDefinition{
		Kind: kinds.TypeExtensionDefinition,
	}
}

func (def *TypeExtensionDefinition) GetKind() string {
	return def.Kind
}

func (def *TypeExtensionDefinition) GetLoc() Location {
	return def.Loc
}

func (def *TypeExtensionDefinition) GetName() Name {
	return *NewName() // TODO: maybe update Definition interface to remove GetName()
}
