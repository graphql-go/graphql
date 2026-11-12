package ast_test

import (
	"testing"

	"github.com/graphql-go/graphql/language/ast"
)

func TestOperationDefinition(t *testing.T) {
	loc := &ast.Location{Start: 1, End: 2}
	name := &ast.Name{Value: "MyQuery"}
	variable := &ast.VariableDefinition{
		Variable: &ast.Variable{Name: &ast.Name{Value: "id"}},
		Type:     &ast.Named{Name: &ast.Name{Value: "String"}},
	}
	directive := &ast.Directive{Name: &ast.Name{Value: "include"}}
	selection := &ast.Field{Name: &ast.Name{Value: "field1"}}
	selSet := &ast.SelectionSet{Selections: []ast.Selection{selection}}

	def := ast.NewOperationDefinition(&ast.OperationDefinition{
		Loc:                 loc,
		Operation:           "query",
		Name:                name,
		VariableDefinitions: []*ast.VariableDefinition{variable},
		Directives:          []*ast.Directive{directive},
		SelectionSet:        selSet,
	})
	if def.GetKind() != "OperationDefinition" {
		t.Fatalf("expected OperationDefinition, got %s", def.GetKind())
	}
	if def.GetLoc() != loc {
		t.Fatal("Loc mismatch")
	}
	if def.GetOperation() != "query" {
		t.Fatalf("expected query, got %s", def.GetOperation())
	}
	if def.GetName() != name {
		t.Fatal("Name mismatch")
	}
	if len(def.GetVariableDefinitions()) != 1 || def.GetVariableDefinitions()[0] != variable {
		t.Fatal("VariableDefinitions mismatch")
	}
	if len(def.GetDirectives()) != 1 || def.GetDirectives()[0] != directive {
		t.Fatal("Directives mismatch")
	}
	if def.GetSelectionSet() != selSet {
		t.Fatal("SelectionSet mismatch")
	}
}

func TestFragmentDefinition(t *testing.T) {
	loc := &ast.Location{Start: 1, End: 2}
	name := &ast.Name{Value: "MyFragment"}
	typeCond := &ast.Named{Name: &ast.Name{Value: "Foo"}}
	directive := &ast.Directive{Name: &ast.Name{Value: "include"}}
	selection := &ast.Field{Name: &ast.Name{Value: "field1"}}
	selSet := &ast.SelectionSet{Selections: []ast.Selection{selection}}

	def := ast.NewFragmentDefinition(&ast.FragmentDefinition{
		Loc:          loc,
		Name:         name,
		TypeCondition: typeCond,
		Directives:   []*ast.Directive{directive},
		SelectionSet: selSet,
	})
	if def.GetKind() != "FragmentDefinition" {
		t.Fatalf("expected FragmentDefinition, got %s", def.GetKind())
	}
	if def.GetLoc() != loc {
		t.Fatal("Loc mismatch")
	}
	if def.GetName() != name {
		t.Fatal("Name mismatch")
	}
	if len(def.GetVariableDefinitions()) != 0 {
		t.Fatal("expected empty VariableDefinitions")
	}
	if def.GetSelectionSet() != selSet {
		t.Fatal("SelectionSet mismatch")
	}
	if def.GetOperation() != "" {
		t.Fatal("expected empty operation")
	}
}

func TestVariableDefinition(t *testing.T) {
	loc := &ast.Location{Start: 1, End: 2}
	variable := &ast.Variable{Name: &ast.Name{Value: "x"}}
	namedType := &ast.Named{Name: &ast.Name{Value: "String"}}
	defaultVal := &ast.StringValue{Value: "default"}

	def := ast.NewVariableDefinition(&ast.VariableDefinition{
		Loc:          loc,
		Variable:     variable,
		Type:         namedType,
		DefaultValue: defaultVal,
	})
	if def.GetKind() != "VariableDefinition" {
		t.Fatalf("expected VariableDefinition, got %s", def.GetKind())
	}
	if def.GetLoc() != loc {
		t.Fatal("Loc mismatch")
	}
}

func TestTypeExtensionDefinition(t *testing.T) {
	loc := &ast.Location{Start: 1, End: 2}
	innerDef := &ast.ObjectDefinition{
		Kind: "ObjectDefinition",
		Name: &ast.Name{Value: "ExtendedType"},
		Fields: []*ast.FieldDefinition{
			{Name: &ast.Name{Value: "newField"}},
		},
	}

	def := ast.NewTypeExtensionDefinition(&ast.TypeExtensionDefinition{
		Loc:        loc,
		Definition: innerDef,
	})
	if def.GetKind() != "TypeExtensionDefinition" {
		t.Fatalf("expected TypeExtensionDefinition, got %s", def.GetKind())
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

func TestDirectiveDefinition(t *testing.T) {
	loc := &ast.Location{Start: 1, End: 2}
	name := &ast.Name{Value: "skip"}
	desc := &ast.StringValue{Value: "Directive description"}
	arg := &ast.InputValueDefinition{Name: &ast.Name{Value: "if"}}
	location := &ast.Name{Value: "FIELD"}

	def := ast.NewDirectiveDefinition(&ast.DirectiveDefinition{
		Loc:         loc,
		Name:        name,
		Description: desc,
		Arguments:   []*ast.InputValueDefinition{arg},
		Locations:   []*ast.Name{location},
	})
	if def.GetKind() != "DirectiveDefinition" {
		t.Fatalf("expected DirectiveDefinition, got %s", def.GetKind())
	}
	if def.GetLoc() != loc {
		t.Fatal("Loc mismatch")
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

func TestDefinitionInterfaceCompileTimeChecks(t *testing.T) {
	var d ast.Definition
	d = ast.NewOperationDefinition(nil)
	if d.GetKind() != "OperationDefinition" {
		t.Fatal("expected OperationDefinition")
	}
	d = ast.NewFragmentDefinition(nil)
	if d.GetKind() != "FragmentDefinition" {
		t.Fatal("expected FragmentDefinition")
	}

	var tsd ast.TypeSystemDefinition
	tsd = ast.NewTypeExtensionDefinition(nil)
	if tsd.GetKind() != "TypeExtensionDefinition" {
		t.Fatal("expected TypeExtensionDefinition")
	}
	tsd = ast.NewDirectiveDefinition(nil)
	if tsd.GetKind() != "DirectiveDefinition" {
		t.Fatal("expected DirectiveDefinition")
	}

	var n ast.Node
	n = ast.NewOperationDefinition(nil)
	if n.GetKind() != "OperationDefinition" {
		t.Fatal("expected OperationDefinition as Node")
	}
	n = ast.NewFragmentDefinition(nil)
	if n.GetKind() != "FragmentDefinition" {
		t.Fatal("expected FragmentDefinition as Node")
	}
	n = ast.NewVariableDefinition(nil)
	if n.GetKind() != "VariableDefinition" {
		t.Fatal("expected VariableDefinition as Node")
	}
	n = ast.NewTypeExtensionDefinition(nil)
	if n.GetKind() != "TypeExtensionDefinition" {
		t.Fatal("expected TypeExtensionDefinition as Node")
	}
	n = ast.NewDirectiveDefinition(nil)
	if n.GetKind() != "DirectiveDefinition" {
		t.Fatal("expected DirectiveDefinition as Node")
	}
}
