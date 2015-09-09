package ast

import (
	"github.com/chris-ramon/graphql-go/language/kinds"
)

// ScalarTypeDefinition implements Definition
type ScalarTypeDefinition struct {
	Kind string
	Loc  Location
	Name Name
}

func NewScalarTypeDefinition() *ScalarTypeDefinition {
	return &ScalarTypeDefinition{
		Kind: kinds.ScalarTypeDefinition,
	}
}

func (def *ScalarTypeDefinition) GetKind() string {
	return def.Kind
}

func (def *ScalarTypeDefinition) GetLoc() Location {
	return def.Loc
}

func (def *ScalarTypeDefinition) GetName() Name {
	return def.Name
}
