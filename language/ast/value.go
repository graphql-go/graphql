package ast

type Value interface {
}

type Variable struct {
	Kind string
	Loc  Location
	Name *Name
}

type IntValue struct {
	Kind  string
	Loc   Location
	Value string
}

type FloatValue struct {
	Kind  string
	Loc   Location
	Value string
}

type StringValue struct {
	Kind  string
	Loc   Location
	Value string
}

type BooleanValue struct {
	Kind  string
	Loc   Location
	Value bool
}

type EnumValue struct {
	Kind  string
	Loc   Location
	Value string
}

type ListValue struct {
	Kind   string
	Loc    Location
	Values []Value
}

type ObjectValue struct {
	Kind   string
	Loc    Location
	Fields []ObjectField
}

type ObjectField struct {
	Kind  string
	Name  *Name
	Loc   Location
	Value Value
}

// Type does not exists in graphql-js ast
type ArrayValue struct {
	Kind   string
	Loc    Location
	Values []Value
}

// Ensure that all value types implements Value interface
var _ Value = (*Variable)(nil)
var _ Value = (*IntValue)(nil)
var _ Value = (*FloatValue)(nil)
var _ Value = (*StringValue)(nil)
var _ Value = (*BooleanValue)(nil)
var _ Value = (*EnumValue)(nil)
var _ Value = (*ListValue)(nil)
var _ Value = (*ObjectValue)(nil)
