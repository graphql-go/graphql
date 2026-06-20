package ast_test

import (
	"testing"

	"github.com/graphql-go/graphql/language/ast"
)

func TestSchemaDefinition(t *testing.T) {
	loc := &ast.Location{Start: 1, End: 2}
	def := ast.NewSchemaDefinition(&ast.SchemaDefinition{
		Loc: loc,
		Directives: []*ast.Directive{
			{Kind: "Directive", Name: &ast.Name{Value: "include"}},
		},
		OperationTypes: []*ast.OperationTypeDefinition{
			{Operation: "query", Type: &ast.Named{Name: &ast.Name{Value: "Query"}}},
		},
	})
	if def.GetKind() != "SchemaDefinition" {
		t.Fatalf("expected SchemaDefinition, got %s", def.GetKind())
	}
	if def.GetLoc() != loc {
		t.Fatal("Loc mismatch")
	}
	if len(def.GetVariableDefinitions()) != 0 {
		t.Fatal("expected empty VariableDefinitions")
	}
	if def.GetSelectionSet() == nil {
		t.Fatal("expected non-nil SelectionSet")
	}
	if def.GetOperation() != "" {
		t.Fatal("expected empty operation")
	}
}

func TestOperationTypeDefinition(t *testing.T) {
	loc := &ast.Location{Start: 1, End: 2}
	named := &ast.Named{Name: &ast.Name{Value: "Query"}}
	def := ast.NewOperationTypeDefinition(&ast.OperationTypeDefinition{
		Loc:       loc,
		Operation: "query",
		Type:      named,
	})
	if def.GetKind() != "OperationTypeDefinition" {
		t.Fatalf("expected OperationTypeDefinition, got %s", def.GetKind())
	}
	if def.GetLoc() != loc {
		t.Fatal("Loc mismatch")
	}
}

func TestScalarDefinition(t *testing.T) {
	loc := &ast.Location{Start: 1, End: 2}
	name := &ast.Name{Value: "CustomScalar"}
	desc := &ast.StringValue{Value: "A custom scalar"}
	def := ast.NewScalarDefinition(&ast.ScalarDefinition{
		Loc:         loc,
		Name:        name,
		Description: desc,
		Directives:  []*ast.Directive{},
	})
	if def.GetKind() != "ScalarDefinition" {
		t.Fatalf("expected ScalarDefinition, got %s", def.GetKind())
	}
	if def.GetLoc() != loc {
		t.Fatal("Loc mismatch")
	}
	if def.GetName() != name {
		t.Fatal("Name mismatch")
	}
	if def.GetDescription() != desc {
		t.Fatal("Description mismatch")
	}
	if len(def.GetVariableDefinitions()) != 0 {
		t.Fatal("expected empty VariableDefinitions")
	}
	if def.GetSelectionSet() == nil {
		t.Fatal("expected non-nil SelectionSet")
	}
	if def.GetOperation() != "" {
		t.Fatal("expected empty operation")
	}
}

func TestObjectDefinition(t *testing.T) {
	loc := &ast.Location{Start: 1, End: 2}
	name := &ast.Name{Value: "Foo"}
	desc := &ast.StringValue{Value: "A foo type"}
	iface := &ast.Named{Name: &ast.Name{Value: "Bar"}}
	field := &ast.FieldDefinition{Name: &ast.Name{Value: "baz"}}
	def := ast.NewObjectDefinition(&ast.ObjectDefinition{
		Loc:         loc,
		Name:        name,
		Description: desc,
		Interfaces:  []*ast.Named{iface},
		Directives:  []*ast.Directive{},
		Fields:      []*ast.FieldDefinition{field},
	})
	if def.GetKind() != "ObjectDefinition" {
		t.Fatalf("expected ObjectDefinition, got %s", def.GetKind())
	}
	if def.GetLoc() != loc {
		t.Fatal("Loc mismatch")
	}
	if def.GetName() != name {
		t.Fatal("Name mismatch")
	}
	if def.GetDescription() != desc {
		t.Fatal("Description mismatch")
	}
	if len(def.GetVariableDefinitions()) != 0 {
		t.Fatal("expected empty VariableDefinitions")
	}
	if def.GetSelectionSet() == nil {
		t.Fatal("expected non-nil SelectionSet")
	}
	if def.GetOperation() != "" {
		t.Fatal("expected empty operation")
	}
}

func TestFieldDefinition(t *testing.T) {
	loc := &ast.Location{Start: 1, End: 2}
	name := &ast.Name{Value: "field1"}
	desc := &ast.StringValue{Value: "A field"}
	arg := &ast.InputValueDefinition{Name: &ast.Name{Value: "arg1"}}
	def := ast.NewFieldDefinition(&ast.FieldDefinition{
		Loc:         loc,
		Name:        name,
		Description: desc,
		Arguments:   []*ast.InputValueDefinition{arg},
		Type:        &ast.Named{Name: &ast.Name{Value: "String"}},
		Directives:  []*ast.Directive{},
	})
	if def.GetKind() != "FieldDefinition" {
		t.Fatalf("expected FieldDefinition, got %s", def.GetKind())
	}
	if def.GetLoc() != loc {
		t.Fatal("Loc mismatch")
	}
	if def.GetDescription() != desc {
		t.Fatal("Description mismatch")
	}
}

func TestInputValueDefinition(t *testing.T) {
	loc := &ast.Location{Start: 1, End: 2}
	name := &ast.Name{Value: "input1"}
	desc := &ast.StringValue{Value: "An input"}
	def := ast.NewInputValueDefinition(&ast.InputValueDefinition{
		Loc:          loc,
		Name:         name,
		Description:  desc,
		Type:         &ast.Named{Name: &ast.Name{Value: "String"}},
		DefaultValue: &ast.StringValue{Value: "default"},
		Directives:   []*ast.Directive{},
	})
	if def.GetKind() != "InputValueDefinition" {
		t.Fatalf("expected InputValueDefinition, got %s", def.GetKind())
	}
	if def.GetLoc() != loc {
		t.Fatal("Loc mismatch")
	}
	if def.GetDescription() != desc {
		t.Fatal("Description mismatch")
	}
}

func TestInterfaceDefinition(t *testing.T) {
	loc := &ast.Location{Start: 1, End: 2}
	name := &ast.Name{Value: "MyInterface"}
	desc := &ast.StringValue{Value: "An interface"}
	field := &ast.FieldDefinition{Name: &ast.Name{Value: "field1"}}
	def := ast.NewInterfaceDefinition(&ast.InterfaceDefinition{
		Loc:         loc,
		Name:        name,
		Description: desc,
		Directives:  []*ast.Directive{},
		Fields:      []*ast.FieldDefinition{field},
	})
	if def.GetKind() != "InterfaceDefinition" {
		t.Fatalf("expected InterfaceDefinition, got %s", def.GetKind())
	}
	if def.GetLoc() != loc {
		t.Fatal("Loc mismatch")
	}
	if def.GetName() != name {
		t.Fatal("Name mismatch")
	}
	if def.GetDescription() != desc {
		t.Fatal("Description mismatch")
	}
	if len(def.GetVariableDefinitions()) != 0 {
		t.Fatal("expected empty VariableDefinitions")
	}
	if def.GetSelectionSet() == nil {
		t.Fatal("expected non-nil SelectionSet")
	}
	if def.GetOperation() != "" {
		t.Fatal("expected empty operation")
	}
}

func TestUnionDefinition(t *testing.T) {
	loc := &ast.Location{Start: 1, End: 2}
	name := &ast.Name{Value: "MyUnion"}
	desc := &ast.StringValue{Value: "A union"}
	member := &ast.Named{Name: &ast.Name{Value: "TypeA"}}
	def := ast.NewUnionDefinition(&ast.UnionDefinition{
		Loc:         loc,
		Name:        name,
		Description: desc,
		Directives:  []*ast.Directive{},
		Types:       []*ast.Named{member},
	})
	if def.GetKind() != "UnionDefinition" {
		t.Fatalf("expected UnionDefinition, got %s", def.GetKind())
	}
	if def.GetLoc() != loc {
		t.Fatal("Loc mismatch")
	}
	if def.GetName() != name {
		t.Fatal("Name mismatch")
	}
	if def.GetDescription() != desc {
		t.Fatal("Description mismatch")
	}
	if len(def.GetVariableDefinitions()) != 0 {
		t.Fatal("expected empty VariableDefinitions")
	}
	if def.GetSelectionSet() == nil {
		t.Fatal("expected non-nil SelectionSet")
	}
	if def.GetOperation() != "" {
		t.Fatal("expected empty operation")
	}
}

func TestEnumDefinition(t *testing.T) {
	loc := &ast.Location{Start: 1, End: 2}
	name := &ast.Name{Value: "MyEnum"}
	desc := &ast.StringValue{Value: "An enum"}
	val := &ast.EnumValueDefinition{Name: &ast.Name{Value: "VAL_A"}}
	def := ast.NewEnumDefinition(&ast.EnumDefinition{
		Loc:         loc,
		Name:        name,
		Description: desc,
		Directives:  []*ast.Directive{},
		Values:      []*ast.EnumValueDefinition{val},
	})
	if def.GetKind() != "EnumDefinition" {
		t.Fatalf("expected EnumDefinition, got %s", def.GetKind())
	}
	if def.GetLoc() != loc {
		t.Fatal("Loc mismatch")
	}
	if def.GetName() != name {
		t.Fatal("Name mismatch")
	}
	if def.GetDescription() != desc {
		t.Fatal("Description mismatch")
	}
	if len(def.GetVariableDefinitions()) != 0 {
		t.Fatal("expected empty VariableDefinitions")
	}
	if def.GetSelectionSet() == nil {
		t.Fatal("expected non-nil SelectionSet")
	}
	if def.GetOperation() != "" {
		t.Fatal("expected empty operation")
	}
}

func TestEnumValueDefinition(t *testing.T) {
	loc := &ast.Location{Start: 1, End: 2}
	name := &ast.Name{Value: "VAL_A"}
	desc := &ast.StringValue{Value: "An enum value"}
	def := ast.NewEnumValueDefinition(&ast.EnumValueDefinition{
		Loc:         loc,
		Name:        name,
		Description: desc,
		Directives:  []*ast.Directive{},
	})
	if def.GetKind() != "EnumValueDefinition" {
		t.Fatalf("expected EnumValueDefinition, got %s", def.GetKind())
	}
	if def.GetLoc() != loc {
		t.Fatal("Loc mismatch")
	}
	if def.GetDescription() != desc {
		t.Fatal("Description mismatch")
	}
}

func TestInputObjectDefinition(t *testing.T) {
	loc := &ast.Location{Start: 1, End: 2}
	name := &ast.Name{Value: "MyInput"}
	desc := &ast.StringValue{Value: "An input object"}
	field := &ast.InputValueDefinition{Name: &ast.Name{Value: "field1"}}
	def := ast.NewInputObjectDefinition(&ast.InputObjectDefinition{
		Loc:         loc,
		Name:        name,
		Description: desc,
		Directives:  []*ast.Directive{},
		Fields:      []*ast.InputValueDefinition{field},
	})
	if def.GetKind() != "InputObjectDefinition" {
		t.Fatalf("expected InputObjectDefinition, got %s", def.GetKind())
	}
	if def.GetLoc() != loc {
		t.Fatal("Loc mismatch")
	}
	if def.GetName() != name {
		t.Fatal("Name mismatch")
	}
	if def.GetDescription() != desc {
		t.Fatal("Description mismatch")
	}
	if len(def.GetVariableDefinitions()) != 0 {
		t.Fatal("expected empty VariableDefinitions")
	}
	if def.GetSelectionSet() == nil {
		t.Fatal("expected non-nil SelectionSet")
	}
	if def.GetOperation() != "" {
		t.Fatal("expected empty operation")
	}
}

func TestTypeDefInterfaceCompileTimeChecks(t *testing.T) {
	var td ast.TypeDefinition
	td = ast.NewScalarDefinition(nil)
	if td.GetKind() != "ScalarDefinition" {
		t.Fatal("expected ScalarDefinition")
	}
	td = ast.NewObjectDefinition(nil)
	if td.GetKind() != "ObjectDefinition" {
		t.Fatal("expected ObjectDefinition")
	}
	td = ast.NewInterfaceDefinition(nil)
	if td.GetKind() != "InterfaceDefinition" {
		t.Fatal("expected InterfaceDefinition")
	}
	td = ast.NewUnionDefinition(nil)
	if td.GetKind() != "UnionDefinition" {
		t.Fatal("expected UnionDefinition")
	}
	td = ast.NewEnumDefinition(nil)
	if td.GetKind() != "EnumDefinition" {
		t.Fatal("expected EnumDefinition")
	}
	td = ast.NewInputObjectDefinition(nil)
	if td.GetKind() != "InputObjectDefinition" {
		t.Fatal("expected InputObjectDefinition")
	}

	var tsd ast.TypeSystemDefinition
	tsd = ast.NewSchemaDefinition(nil)
	if tsd.GetKind() != "SchemaDefinition" {
		t.Fatal("expected SchemaDefinition")
	}
	tsd = ast.NewTypeExtensionDefinition(nil)
	if tsd.GetKind() != "TypeExtensionDefinition" {
		t.Fatal("expected TypeExtensionDefinition")
	}
	tsd = ast.NewDirectiveDefinition(nil)
	if tsd.GetKind() != "DirectiveDefinition" {
		t.Fatal("expected DirectiveDefinition")
	}

	var n ast.Node
	n = ast.NewScalarDefinition(nil)
	if n.GetKind() != "ScalarDefinition" {
		t.Fatal("expected ScalarDefinition as Node")
	}
	_ = n
}
