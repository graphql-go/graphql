package graphql

import (
	"github.com/chris-ramon/graphql/kinds"
)

// Directive implements Node
type AstDirective struct {
	Kind      string
	Loc       *AstLocation
	Name      *AstName
	Arguments []*AstArgument
}

func NewAstDirective(dir *AstDirective) *AstDirective {
	if dir == nil {
		dir = &AstDirective{}
	}
	return &AstDirective{
		Kind:      kinds.Directive,
		Loc:       dir.Loc,
		Name:      dir.Name,
		Arguments: dir.Arguments,
	}
}

func (dir *AstDirective) GetKind() string {
	return dir.Kind
}

func (dir *AstDirective) GetLoc() *AstLocation {
	return dir.Loc
}
