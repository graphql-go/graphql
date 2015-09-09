package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// ObjectTypeDefinition implements Definition
type ObjectTypeDefinition struct {
	Kind       string
	Loc        Location
	Name       Name
	Interfaces []NamedType
	Fields     []interface{}
}

func NewObjectTypeDefinition() *ObjectTypeDefinition {
	return &ObjectTypeDefinition{
		Kind: kinds.ObjectTypeDefinition,
	}
}

func (def *ObjectTypeDefinition) GetKind() string {
	return def.Kind
}

func (def *ObjectTypeDefinition) GetLoc() Location {
	return def.Loc
}

func (def *ObjectTypeDefinition) GetName() Name {
	return def.Name
}
