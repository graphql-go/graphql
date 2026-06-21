package ast

import (
	"testing"

	"github.com/graphql-go/graphql/language/kinds"
)

func TestDefConstructors_NilInput(t *testing.T) {
	tests := []struct {
		name string
		got  Node
	}{
		{"NewOperationDefinition", NewOperationDefinition(nil)},
		{"NewFragmentDefinition", NewFragmentDefinition(nil)},
		{"NewVariableDefinition", NewVariableDefinition(nil)},
		{"NewTypeExtensionDefinition", NewTypeExtensionDefinition(nil)},
		{"NewDirectiveDefinition", NewDirectiveDefinition(nil)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got == nil {
				t.Fatal("expected non-nil result")
			}
		})
	}
}

func TestDefOperationDefinition_NewWithInput(t *testing.T) {
	loc := &Location{Start: 1, End: 2}
	name := &Name{Value: "MyOp"}
	vd := &VariableDefinition{}
	dir := &Directive{}
	ss := &SelectionSet{}
	op := NewOperationDefinition(&OperationDefinition{
		Loc:                 loc,
		Operation:           "query",
		Name:                name,
		VariableDefinitions: []*VariableDefinition{vd},
		Directives:          []*Directive{dir},
		SelectionSet:        ss,
	})
	if op.Kind != kinds.OperationDefinition {
		t.Fatalf("expected kind %s, got %s", kinds.OperationDefinition, op.Kind)
	}
	if op.Loc != loc {
		t.Fatal("Loc mismatch")
	}
	if op.Operation != "query" {
		t.Fatal("Operation mismatch")
	}
	if op.Name != name {
		t.Fatal("Name mismatch")
	}
	if len(op.VariableDefinitions) != 1 || op.VariableDefinitions[0] != vd {
		t.Fatal("VariableDefinitions mismatch")
	}
	if len(op.Directives) != 1 || op.Directives[0] != dir {
		t.Fatal("Directives mismatch")
	}
	if op.SelectionSet != ss {
		t.Fatal("SelectionSet mismatch")
	}
}

func TestDefVariableDefinition_NewWithInput(t *testing.T) {
	loc := &Location{Start: 1, End: 2}
	variable := &Variable{Name: &Name{Value: "x"}}
	typ := &Named{Name: &Name{Value: "String"}}
	defaultVal := &StringValue{Value: "hello"}
	vd := NewVariableDefinition(&VariableDefinition{
		Loc:          loc,
		Variable:     variable,
		Type:         typ,
		DefaultValue: defaultVal,
	})
	if vd.Kind != kinds.VariableDefinition {
		t.Fatalf("expected kind %s, got %s", kinds.VariableDefinition, vd.Kind)
	}
	if vd.Loc != loc {
		t.Fatal("Loc mismatch")
	}
	if vd.Variable != variable {
		t.Fatal("Variable mismatch")
	}
	if vd.Type != typ {
		t.Fatal("Type mismatch")
	}
	if vd.DefaultValue != defaultVal {
		t.Fatal("DefaultValue mismatch")
	}
}

func TestDefTypeExtensionDefinition_NewWithInput(t *testing.T) {
	loc := &Location{Start: 1, End: 2}
	def := &ObjectDefinition{Name: &Name{Value: "ExtendedType"}}
	ted := NewTypeExtensionDefinition(&TypeExtensionDefinition{
		Loc:        loc,
		Definition: def,
	})
	if ted.Kind != kinds.TypeExtensionDefinition {
		t.Fatalf("expected kind %s, got %s", kinds.TypeExtensionDefinition, ted.Kind)
	}
	if ted.Loc != loc {
		t.Fatal("Loc mismatch")
	}
	if ted.Definition != def {
		t.Fatal("Definition mismatch")
	}
}

func TestDefDirectiveDefinition_NewWithInput(t *testing.T) {
	loc := &Location{Start: 1, End: 2}
	name := &Name{Value: "skip"}
	desc := &StringValue{Value: "Directive description"}
	arg := &InputValueDefinition{Name: &Name{Value: "if"}}
	location := &Name{Value: "FIELD"}
	dd := NewDirectiveDefinition(&DirectiveDefinition{
		Loc:         loc,
		Name:        name,
		Description: desc,
		Arguments:   []*InputValueDefinition{arg},
		Locations:   []*Name{location},
	})
	if dd.Kind != kinds.DirectiveDefinition {
		t.Fatalf("expected kind %s, got %s", kinds.DirectiveDefinition, dd.Kind)
	}
	if dd.Loc != loc {
		t.Fatal("Loc mismatch")
	}
	if dd.Name != name {
		t.Fatal("Name mismatch")
	}
	if dd.Description != desc {
		t.Fatal("Description mismatch")
	}
	if len(dd.Arguments) != 1 || dd.Arguments[0] != arg {
		t.Fatal("Arguments mismatch")
	}
	if len(dd.Locations) != 1 || dd.Locations[0] != location {
		t.Fatal("Locations mismatch")
	}
}
