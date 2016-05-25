package ast

import (
	"github.com/sprucehealth/graphql/language/kinds"
)

// Argument implements Node
type Argument struct {
	Kind  string
	Loc   *Location
	Name  *Name
	Value Value
}

func NewArgument(arg *Argument) *Argument {
	if arg == nil {
		return &Argument{Kind: kinds.Argument}
	}
	return &Argument{
		Kind:  kinds.Argument,
		Loc:   arg.Loc,
		Name:  arg.Name,
		Value: arg.Value,
	}
}

func (arg *Argument) GetKind() string {
	return arg.Kind
}

func (arg *Argument) GetLoc() *Location {
	return arg.Loc
}
