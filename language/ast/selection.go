package ast

type Selection interface {
}

type Field struct {
	Kind         string
	Loc          Location
	Alias        *Name
	Name         *Name
	Arguments    []Argument
	Directives   []Directive
	SelectionSet SelectionSet
}

type FragmentSpread struct {
	Kind       string
	Loc        Location
	Name       *Name
	Directives []Directive
}

type InlineFragment struct {
	Kind          string
	Loc           Location
	TypeCondition NamedType
	Directives    []Directive
	SelectionSet  SelectionSet
}

// Ensure that all definition types implements Definition interface
var _ Selection = (*Field)(nil)
var _ Selection = (*FragmentSpread)(nil)
var _ Selection = (*InlineFragment)(nil)
