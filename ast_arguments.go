package graphql

import (
	"github.com/chris-ramon/graphql/kinds"
)

// Argument implements Node
type AstArgument struct {
	Kind  string
	Loc   *AstLocation
	Name  *AstName
	Value Value
}

func NewAstArgument(arg *AstArgument) *AstArgument {
	if arg == nil {
		arg = &AstArgument{}
	}
	return &AstArgument{
		Kind:  kinds.Argument,
		Loc:   arg.Loc,
		Name:  arg.Name,
		Value: arg.Value,
	}
}

func (arg *AstArgument) GetKind() string {
	return arg.Kind
}

func (arg *AstArgument) GetLoc() *AstLocation {
	return arg.Loc
}
