package ast

type Value interface {
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
