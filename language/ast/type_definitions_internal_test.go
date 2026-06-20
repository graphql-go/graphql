package ast

import (
	"testing"
)

func TestTypeDefConstructors_NilInput(t *testing.T) {
	tests := []struct {
		name string
		got  Node
	}{
		{"NewSchemaDefinition", NewSchemaDefinition(nil)},
		{"NewOperationTypeDefinition", NewOperationTypeDefinition(nil)},
		{"NewScalarDefinition", NewScalarDefinition(nil)},
		{"NewObjectDefinition", NewObjectDefinition(nil)},
		{"NewFieldDefinition", NewFieldDefinition(nil)},
		{"NewInputValueDefinition", NewInputValueDefinition(nil)},
		{"NewInterfaceDefinition", NewInterfaceDefinition(nil)},
		{"NewUnionDefinition", NewUnionDefinition(nil)},
		{"NewEnumDefinition", NewEnumDefinition(nil)},
		{"NewEnumValueDefinition", NewEnumValueDefinition(nil)},
		{"NewInputObjectDefinition", NewInputObjectDefinition(nil)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got == nil {
				t.Fatal("expected non-nil result")
			}
		})
	}
}

func TestTypeDefStubReturns(t *testing.T) {
	tests := []struct {
		name string
		desc DescribableNode
	}{
		{"ScalarDefinition", NewScalarDefinition(&ScalarDefinition{Description: &StringValue{Value: "desc"}})},
		{"ObjectDefinition", NewObjectDefinition(&ObjectDefinition{Description: &StringValue{Value: "desc"}})},
		{"FieldDefinition", NewFieldDefinition(&FieldDefinition{Description: &StringValue{Value: "desc"}})},
		{"InputValueDefinition", NewInputValueDefinition(&InputValueDefinition{Description: &StringValue{Value: "desc"}})},
		{"InterfaceDefinition", NewInterfaceDefinition(&InterfaceDefinition{Description: &StringValue{Value: "desc"}})},
		{"UnionDefinition", NewUnionDefinition(&UnionDefinition{Description: &StringValue{Value: "desc"}})},
		{"EnumDefinition", NewEnumDefinition(&EnumDefinition{Description: &StringValue{Value: "desc"}})},
		{"EnumValueDefinition", NewEnumValueDefinition(&EnumValueDefinition{Description: &StringValue{Value: "desc"}})},
		{"InputObjectDefinition", NewInputObjectDefinition(&InputObjectDefinition{Description: &StringValue{Value: "desc"}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.desc.GetDescription() == nil {
				t.Fatal("expected non-nil description")
			}
		})
	}
}
