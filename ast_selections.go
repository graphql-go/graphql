package graphql

import (
	"github.com/chris-ramon/graphql/kinds"
)

type Selection interface {
}

// Ensure that all definition types implements Selection interface
var _ Selection = (*AstField)(nil)
var _ Selection = (*AstFragmentSpread)(nil)
var _ Selection = (*AstInlineFragment)(nil)

// AstField implements Node, Selection
type AstField struct {
	Kind         string
	Loc          *AstLocation
	Alias        *AstName
	Name         *AstName
	Arguments    []*AstArgument
	Directives   []*AstDirective
	SelectionSet *AstSelectionSet
}

func NewField(f *AstField) *AstField {
	if f == nil {
		f = &AstField{}
	}
	return &AstField{
		Kind:         kinds.Field,
		Loc:          f.Loc,
		Alias:        f.Alias,
		Name:         f.Name,
		Arguments:    f.Arguments,
		Directives:   f.Directives,
		SelectionSet: f.SelectionSet,
	}
}

func (f *AstField) GetKind() string {
	return f.Kind
}

func (f *AstField) GetLoc() *AstLocation {
	return f.Loc
}

// FragmentSpread implements Node, Selection
type AstFragmentSpread struct {
	Kind       string
	Loc        *AstLocation
	Name       *AstName
	Directives []*AstDirective
}

func NewAstFragmentSpread(fs *AstFragmentSpread) *AstFragmentSpread {
	if fs == nil {
		fs = &AstFragmentSpread{}
	}
	return &AstFragmentSpread{
		Kind:       kinds.FragmentSpread,
		Loc:        fs.Loc,
		Name:       fs.Name,
		Directives: fs.Directives,
	}
}

func (fs *AstFragmentSpread) GetKind() string {
	return fs.Kind
}

func (fs *AstFragmentSpread) GetLoc() *AstLocation {
	return fs.Loc
}

// InlineFragment implements Node, Selection
type AstInlineFragment struct {
	Kind          string
	Loc           *AstLocation
	TypeCondition *AstNamed
	Directives    []*AstDirective
	SelectionSet  *AstSelectionSet
}

func NewAstInlineFragment(f *AstInlineFragment) *AstInlineFragment {
	if f == nil {
		f = &AstInlineFragment{}
	}
	return &AstInlineFragment{
		Kind:          kinds.InlineFragment,
		Loc:           f.Loc,
		TypeCondition: f.TypeCondition,
		Directives:    f.Directives,
		SelectionSet:  f.SelectionSet,
	}
}

func (f *AstInlineFragment) GetKind() string {
	return f.Kind
}

func (f *AstInlineFragment) GetLoc() *AstLocation {
	return f.Loc
}

// AstSelectionSet implements Node
type AstSelectionSet struct {
	Kind       string
	Loc        *AstLocation
	Selections []Selection
}

func NewAstSelectionSet(ss *AstSelectionSet) *AstSelectionSet {
	if ss == nil {
		ss = &AstSelectionSet{}
	}
	return &AstSelectionSet{
		Kind:       kinds.SelectionSet,
		Loc:        ss.Loc,
		Selections: ss.Selections,
	}
}

func (ss *AstSelectionSet) GetKind() string {
	return ss.Kind
}

func (ss *AstSelectionSet) GetLoc() *AstLocation {
	return ss.Loc
}
