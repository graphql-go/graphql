package visitor_test

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/visitor"
	"github.com/graphql-go/graphql/testutil"
)

func parse(t *testing.T, query string) *ast.Document {
	astDoc, err := parser.Parse(parser.ParseParams{
		Source: query,
		Options: parser.ParseOptions{
			NoLocation: true,
		},
	})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	return astDoc
}

func TestVisitor_AllowsForEditingOnEnter(t *testing.T) {

	query := `{ a, b, c { a, b, c } }`
	astDoc := parse(t, query)

	expectedQuery := `{ a,    c { a,    c } }`
	expectedAST := parse(t, expectedQuery)
	v := &visitor.VisitorOptions{
		Enter: func(p visitor.VisitFuncParams) (string, interface{}) {
			switch node := p.Node.(type) {
			case *ast.Field:
				if node.Name != nil && node.Name.Value == "b" {
					return visitor.ActionUpdate, nil
				}
			}
			return visitor.ActionNoChange, nil
		},
	}

	editedAst := visitor.Visit(astDoc, v, nil)
	if !reflect.DeepEqual(expectedAST, editedAst) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedAST, editedAst))
	}

}
func TestVisitor_AllowsForEditingOnLeave(t *testing.T) {

	query := `{ a, b, c { a, b, c } }`
	astDoc := parse(t, query)

	expectedQuery := `{ a,    c { a,    c } }`
	expectedAST := parse(t, expectedQuery)
	v := &visitor.VisitorOptions{
		Leave: func(p visitor.VisitFuncParams) (string, interface{}) {
			switch node := p.Node.(type) {
			case *ast.Field:
				if node.Name != nil && node.Name.Value == "b" {
					return visitor.ActionUpdate, nil
				}
			}
			return visitor.ActionNoChange, nil
		},
	}

	editedAst := visitor.Visit(astDoc, v, nil)
	if !reflect.DeepEqual(expectedAST, editedAst) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedAST, editedAst))
	}
}

func TestVisitor_VisitsEditedNode(t *testing.T) {

	query := `{ a { x } }`
	astDoc := parse(t, query)

	addedField := &ast.Field{
		Kind: "Field",
		Name: &ast.Name{
			Kind:  "Name",
			Value: "__typename",
		},
	}

	didVisitAddedField := false
	v := &visitor.VisitorOptions{
		Enter: func(p visitor.VisitFuncParams) (string, interface{}) {
			switch node := p.Node.(type) {
			case *ast.Field:
				if node.Name != nil && node.Name.Value == "a" {
					s := node.SelectionSet.Selections
					s = append(s, addedField)
					ss := node.SelectionSet
					ss.Selections = s
					return visitor.ActionUpdate, &ast.Field{
						Kind:         "Field",
						SelectionSet: ss,
					}
				}
				if reflect.DeepEqual(node, addedField) {
					didVisitAddedField = true
				}
			}
			return visitor.ActionNoChange, nil
		},
	}

	_ = visitor.Visit(astDoc, v, nil)
	if didVisitAddedField == false {
		t.Fatalf("Unexpected result, expected didVisitAddedField == true")
	}
}
func TestVisitor_AllowsSkippingASubTree(t *testing.T) {

	query := `{ a, b { x }, c }`
	astDoc := parse(t, query)

	visited := []interface{}{}
	expectedVisited := []interface{}{
		[]interface{}{"enter", "Document", nil},
		[]interface{}{"enter", "OperationDefinition", nil},
		[]interface{}{"enter", "SelectionSet", nil},
		[]interface{}{"enter", "Field", nil},
		[]interface{}{"enter", "Name", "a"},
		[]interface{}{"leave", "Name", "a"},
		[]interface{}{"leave", "Field", nil},
		[]interface{}{"enter", "Field", nil},
		[]interface{}{"enter", "Field", nil},
		[]interface{}{"enter", "Name", "c"},
		[]interface{}{"leave", "Name", "c"},
		[]interface{}{"leave", "Field", nil},
		[]interface{}{"leave", "SelectionSet", nil},
		[]interface{}{"leave", "OperationDefinition", nil},
		[]interface{}{"leave", "Document", nil},
	}

	v := &visitor.VisitorOptions{
		Enter: func(p visitor.VisitFuncParams) (string, interface{}) {
			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []interface{}{"enter", node.Kind, node.Value})
			case *ast.Field:
				visited = append(visited, []interface{}{"enter", node.Kind, nil})
				if node.Name != nil && node.Name.Value == "b" {
					return visitor.ActionSkip, nil
				}
			case ast.Node:
				visited = append(visited, []interface{}{"enter", node.GetKind(), nil})
			default:
				visited = append(visited, []interface{}{"enter", nil, nil})
			}
			return visitor.ActionNoChange, nil
		},
		Leave: func(p visitor.VisitFuncParams) (string, interface{}) {
			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []interface{}{"leave", node.Kind, node.Value})
			case ast.Node:
				visited = append(visited, []interface{}{"leave", node.GetKind(), nil})
			default:
				visited = append(visited, []interface{}{"leave", nil, nil})
			}
			return visitor.ActionNoChange, nil
		},
	}

	_ = visitor.Visit(astDoc, v, nil)

	if !reflect.DeepEqual(visited, expectedVisited) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedVisited, visited))
	}
}

func TestVisitor_AllowsEarlyExitWhileVisiting(t *testing.T) {

	visited := []interface{}{}

	query := `{ a, b { x }, c }`
	astDoc := parse(t, query)

	expectedVisited := []interface{}{
		[]interface{}{"enter", "Document", nil},
		[]interface{}{"enter", "OperationDefinition", nil},
		[]interface{}{"enter", "SelectionSet", nil},
		[]interface{}{"enter", "Field", nil},
		[]interface{}{"enter", "Name", "a"},
		[]interface{}{"leave", "Name", "a"},
		[]interface{}{"leave", "Field", nil},
		[]interface{}{"enter", "Field", nil},
		[]interface{}{"enter", "Name", "b"},
		[]interface{}{"leave", "Name", "b"},
		[]interface{}{"enter", "SelectionSet", nil},
		[]interface{}{"enter", "Field", nil},
		[]interface{}{"enter", "Name", "x"},
	}

	v := &visitor.VisitorOptions{
		Enter: func(p visitor.VisitFuncParams) (string, interface{}) {
			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []interface{}{"enter", node.Kind, node.Value})
				if node.Value == "x" {
					return visitor.ActionBreak, nil
				}
			case ast.Node:
				visited = append(visited, []interface{}{"enter", node.GetKind(), nil})
			default:
				visited = append(visited, []interface{}{"enter", nil, nil})
			}
			return visitor.ActionNoChange, nil
		},
		Leave: func(p visitor.VisitFuncParams) (string, interface{}) {
			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []interface{}{"leave", node.Kind, node.Value})
			case ast.Node:
				visited = append(visited, []interface{}{"leave", node.GetKind(), nil})
			default:
				visited = append(visited, []interface{}{"leave", nil, nil})
			}
			return visitor.ActionNoChange, nil
		},
	}

	_ = visitor.Visit(astDoc, v, nil)

	if !reflect.DeepEqual(visited, expectedVisited) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedVisited, visited))
	}
}

func TestVisitor_AllowsANamedFunctionsVisitorAPI(t *testing.T) {

	query := `{ a, b { x }, c }`
	astDoc := parse(t, query)

	visited := []interface{}{}
	expectedVisited := []interface{}{
		[]interface{}{"enter", "SelectionSet", nil},
		[]interface{}{"enter", "Name", "a"},
		[]interface{}{"enter", "Name", "b"},
		[]interface{}{"enter", "SelectionSet", nil},
		[]interface{}{"enter", "Name", "x"},
		[]interface{}{"leave", "SelectionSet", nil},
		[]interface{}{"enter", "Name", "c"},
		[]interface{}{"leave", "SelectionSet", nil},
	}

	v := &visitor.VisitorOptions{
		KindFuncMap: map[string]visitor.NamedVisitFuncs{
			"Name": visitor.NamedVisitFuncs{
				Kind: func(p visitor.VisitFuncParams) (string, interface{}) {
					switch node := p.Node.(type) {
					case *ast.Name:
						visited = append(visited, []interface{}{"enter", node.Kind, node.Value})
					}
					return visitor.ActionNoChange, nil
				},
			},
			"SelectionSet": visitor.NamedVisitFuncs{
				Enter: func(p visitor.VisitFuncParams) (string, interface{}) {
					switch node := p.Node.(type) {
					case *ast.SelectionSet:
						visited = append(visited, []interface{}{"enter", node.Kind, nil})
					}
					return visitor.ActionNoChange, nil
				},
				Leave: func(p visitor.VisitFuncParams) (string, interface{}) {
					switch node := p.Node.(type) {
					case *ast.SelectionSet:
						visited = append(visited, []interface{}{"leave", node.Kind, nil})
					}
					return visitor.ActionNoChange, nil
				},
			},
		},
	}

	_ = visitor.Visit(astDoc, v, nil)

	if !reflect.DeepEqual(visited, expectedVisited) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedVisited, visited))
	}
}
func TestVisitor_VisitsKitchenSink(t *testing.T) {
	b, err := ioutil.ReadFile("../../kitchen-sink.graphql")
	if err != nil {
		t.Fatalf("unable to load kitchen-sink.graphql")
	}

	query := string(b)
	astDoc := parse(t, query)

	visited := []interface{}{}
	expectedVisited := []interface{}{
		[]interface{}{"enter", "Document", nil, nil},
		[]interface{}{"enter", "OperationDefinition", 0, nil},
		[]interface{}{"enter", "Name", "Name", "OperationDefinition"},
		[]interface{}{"leave", "Name", "Name", "OperationDefinition"},
		[]interface{}{"enter", "VariableDefinition", 0, nil},
		[]interface{}{"enter", "Variable", "Variable", "VariableDefinition"},
		[]interface{}{"enter", "Name", "Name", "Variable"},
		[]interface{}{"leave", "Name", "Name", "Variable"},
		[]interface{}{"leave", "Variable", "Variable", "VariableDefinition"},
		[]interface{}{"enter", "Named", "Type", "VariableDefinition"},
		[]interface{}{"enter", "Name", "Name", "Named"},
		[]interface{}{"leave", "Name", "Name", "Named"},
		[]interface{}{"leave", "Named", "Type", "VariableDefinition"},
		[]interface{}{"leave", "VariableDefinition", 0, nil},
		[]interface{}{"enter", "VariableDefinition", 1, nil},
		[]interface{}{"enter", "Variable", "Variable", "VariableDefinition"},
		[]interface{}{"enter", "Name", "Name", "Variable"},
		[]interface{}{"leave", "Name", "Name", "Variable"},
		[]interface{}{"leave", "Variable", "Variable", "VariableDefinition"},
		[]interface{}{"enter", "Named", "Type", "VariableDefinition"},
		[]interface{}{"enter", "Name", "Name", "Named"},
		[]interface{}{"leave", "Name", "Name", "Named"},
		[]interface{}{"leave", "Named", "Type", "VariableDefinition"},
		[]interface{}{"enter", "EnumValue", "DefaultValue", "VariableDefinition"},
		[]interface{}{"leave", "EnumValue", "DefaultValue", "VariableDefinition"},
		[]interface{}{"leave", "VariableDefinition", 1, nil},
		[]interface{}{"enter", "SelectionSet", "SelectionSet", "OperationDefinition"},
		[]interface{}{"enter", "Field", 0, nil},
		[]interface{}{"enter", "Name", "Alias", "Field"},
		[]interface{}{"leave", "Name", "Alias", "Field"},
		[]interface{}{"enter", "Name", "Name", "Field"},
		[]interface{}{"leave", "Name", "Name", "Field"},
		[]interface{}{"enter", "Argument", 0, nil},
		[]interface{}{"enter", "Name", "Name", "Argument"},
		[]interface{}{"leave", "Name", "Name", "Argument"},
		[]interface{}{"enter", "ListValue", "Value", "Argument"},
		[]interface{}{"enter", "IntValue", 0, nil},
		[]interface{}{"leave", "IntValue", 0, nil},
		[]interface{}{"enter", "IntValue", 1, nil},
		[]interface{}{"leave", "IntValue", 1, nil},
		[]interface{}{"leave", "ListValue", "Value", "Argument"},
		[]interface{}{"leave", "Argument", 0, nil},
		[]interface{}{"enter", "SelectionSet", "SelectionSet", "Field"},
		[]interface{}{"enter", "Field", 0, nil},
		[]interface{}{"enter", "Name", "Name", "Field"},
		[]interface{}{"leave", "Name", "Name", "Field"},
		[]interface{}{"leave", "Field", 0, nil},
		[]interface{}{"enter", "InlineFragment", 1, nil},
		[]interface{}{"enter", "Named", "TypeCondition", "InlineFragment"},
		[]interface{}{"enter", "Name", "Name", "Named"},
		[]interface{}{"leave", "Name", "Name", "Named"},
		[]interface{}{"leave", "Named", "TypeCondition", "InlineFragment"},
		[]interface{}{"enter", "Directive", 0, nil},
		[]interface{}{"enter", "Name", "Name", "Directive"},
		[]interface{}{"leave", "Name", "Name", "Directive"},
		[]interface{}{"leave", "Directive", 0, nil},
		[]interface{}{"enter", "SelectionSet", "SelectionSet", "InlineFragment"},
		[]interface{}{"enter", "Field", 0, nil},
		[]interface{}{"enter", "Name", "Name", "Field"},
		[]interface{}{"leave", "Name", "Name", "Field"},
		[]interface{}{"enter", "SelectionSet", "SelectionSet", "Field"},
		[]interface{}{"enter", "Field", 0, nil},
		[]interface{}{"enter", "Name", "Name", "Field"},
		[]interface{}{"leave", "Name", "Name", "Field"},
		[]interface{}{"leave", "Field", 0, nil},
		[]interface{}{"enter", "Field", 1, nil},
		[]interface{}{"enter", "Name", "Alias", "Field"},
		[]interface{}{"leave", "Name", "Alias", "Field"},
		[]interface{}{"enter", "Name", "Name", "Field"},
		[]interface{}{"leave", "Name", "Name", "Field"},
		[]interface{}{"enter", "Argument", 0, nil},
		[]interface{}{"enter", "Name", "Name", "Argument"},
		[]interface{}{"leave", "Name", "Name", "Argument"},
		[]interface{}{"enter", "IntValue", "Value", "Argument"},
		[]interface{}{"leave", "IntValue", "Value", "Argument"},
		[]interface{}{"leave", "Argument", 0, nil},
		[]interface{}{"enter", "Argument", 1, nil},
		[]interface{}{"enter", "Name", "Name", "Argument"},
		[]interface{}{"leave", "Name", "Name", "Argument"},
		[]interface{}{"enter", "Variable", "Value", "Argument"},
		[]interface{}{"enter", "Name", "Name", "Variable"},
		[]interface{}{"leave", "Name", "Name", "Variable"},
		[]interface{}{"leave", "Variable", "Value", "Argument"},
		[]interface{}{"leave", "Argument", 1, nil},
		[]interface{}{"enter", "Directive", 0, nil},
		[]interface{}{"enter", "Name", "Name", "Directive"},
		[]interface{}{"leave", "Name", "Name", "Directive"},
		[]interface{}{"enter", "Argument", 0, nil},
		[]interface{}{"enter", "Name", "Name", "Argument"},
		[]interface{}{"leave", "Name", "Name", "Argument"},
		[]interface{}{"enter", "Variable", "Value", "Argument"},
		[]interface{}{"enter", "Name", "Name", "Variable"},
		[]interface{}{"leave", "Name", "Name", "Variable"},
		[]interface{}{"leave", "Variable", "Value", "Argument"},
		[]interface{}{"leave", "Argument", 0, nil},
		[]interface{}{"leave", "Directive", 0, nil},
		[]interface{}{"enter", "SelectionSet", "SelectionSet", "Field"},
		[]interface{}{"enter", "Field", 0, nil},
		[]interface{}{"enter", "Name", "Name", "Field"},
		[]interface{}{"leave", "Name", "Name", "Field"},
		[]interface{}{"leave", "Field", 0, nil},
		[]interface{}{"enter", "FragmentSpread", 1, nil},
		[]interface{}{"enter", "Name", "Name", "FragmentSpread"},
		[]interface{}{"leave", "Name", "Name", "FragmentSpread"},
		[]interface{}{"leave", "FragmentSpread", 1, nil},
		[]interface{}{"leave", "SelectionSet", "SelectionSet", "Field"},
		[]interface{}{"leave", "Field", 1, nil},
		[]interface{}{"leave", "SelectionSet", "SelectionSet", "Field"},
		[]interface{}{"leave", "Field", 0, nil},
		[]interface{}{"leave", "SelectionSet", "SelectionSet", "InlineFragment"},
		[]interface{}{"leave", "InlineFragment", 1, nil},
		[]interface{}{"leave", "SelectionSet", "SelectionSet", "Field"},
		[]interface{}{"leave", "Field", 0, nil},
		[]interface{}{"leave", "SelectionSet", "SelectionSet", "OperationDefinition"},
		[]interface{}{"leave", "OperationDefinition", 0, nil},
		[]interface{}{"enter", "OperationDefinition", 1, nil},
		[]interface{}{"enter", "Name", "Name", "OperationDefinition"},
		[]interface{}{"leave", "Name", "Name", "OperationDefinition"},
		[]interface{}{"enter", "SelectionSet", "SelectionSet", "OperationDefinition"},
		[]interface{}{"enter", "Field", 0, nil},
		[]interface{}{"enter", "Name", "Name", "Field"},
		[]interface{}{"leave", "Name", "Name", "Field"},
		[]interface{}{"enter", "Argument", 0, nil},
		[]interface{}{"enter", "Name", "Name", "Argument"},
		[]interface{}{"leave", "Name", "Name", "Argument"},
		[]interface{}{"enter", "IntValue", "Value", "Argument"},
		[]interface{}{"leave", "IntValue", "Value", "Argument"},
		[]interface{}{"leave", "Argument", 0, nil},
		[]interface{}{"enter", "Directive", 0, nil},
		[]interface{}{"enter", "Name", "Name", "Directive"},
		[]interface{}{"leave", "Name", "Name", "Directive"},
		[]interface{}{"leave", "Directive", 0, nil},
		[]interface{}{"enter", "SelectionSet", "SelectionSet", "Field"},
		[]interface{}{"enter", "Field", 0, nil},
		[]interface{}{"enter", "Name", "Name", "Field"},
		[]interface{}{"leave", "Name", "Name", "Field"},
		[]interface{}{"enter", "SelectionSet", "SelectionSet", "Field"},
		[]interface{}{"enter", "Field", 0, nil},
		[]interface{}{"enter", "Name", "Name", "Field"},
		[]interface{}{"leave", "Name", "Name", "Field"},
		[]interface{}{"leave", "Field", 0, nil},
		[]interface{}{"leave", "SelectionSet", "SelectionSet", "Field"},
		[]interface{}{"leave", "Field", 0, nil},
		[]interface{}{"leave", "SelectionSet", "SelectionSet", "Field"},
		[]interface{}{"leave", "Field", 0, nil},
		[]interface{}{"leave", "SelectionSet", "SelectionSet", "OperationDefinition"},
		[]interface{}{"leave", "OperationDefinition", 1, nil},
		[]interface{}{"enter", "FragmentDefinition", 2, nil},
		[]interface{}{"enter", "Name", "Name", "FragmentDefinition"},
		[]interface{}{"leave", "Name", "Name", "FragmentDefinition"},
		[]interface{}{"enter", "Named", "TypeCondition", "FragmentDefinition"},
		[]interface{}{"enter", "Name", "Name", "Named"},
		[]interface{}{"leave", "Name", "Name", "Named"},
		[]interface{}{"leave", "Named", "TypeCondition", "FragmentDefinition"},
		[]interface{}{"enter", "SelectionSet", "SelectionSet", "FragmentDefinition"},
		[]interface{}{"enter", "Field", 0, nil},
		[]interface{}{"enter", "Name", "Name", "Field"},
		[]interface{}{"leave", "Name", "Name", "Field"},
		[]interface{}{"enter", "Argument", 0, nil},
		[]interface{}{"enter", "Name", "Name", "Argument"},
		[]interface{}{"leave", "Name", "Name", "Argument"},
		[]interface{}{"enter", "Variable", "Value", "Argument"},
		[]interface{}{"enter", "Name", "Name", "Variable"},
		[]interface{}{"leave", "Name", "Name", "Variable"},
		[]interface{}{"leave", "Variable", "Value", "Argument"},
		[]interface{}{"leave", "Argument", 0, nil},
		[]interface{}{"enter", "Argument", 1, nil},
		[]interface{}{"enter", "Name", "Name", "Argument"},
		[]interface{}{"leave", "Name", "Name", "Argument"},
		[]interface{}{"enter", "Variable", "Value", "Argument"},
		[]interface{}{"enter", "Name", "Name", "Variable"},
		[]interface{}{"leave", "Name", "Name", "Variable"},
		[]interface{}{"leave", "Variable", "Value", "Argument"},
		[]interface{}{"leave", "Argument", 1, nil},
		[]interface{}{"enter", "Argument", 2, nil},
		[]interface{}{"enter", "Name", "Name", "Argument"},
		[]interface{}{"leave", "Name", "Name", "Argument"},
		[]interface{}{"enter", "ObjectValue", "Value", "Argument"},
		[]interface{}{"enter", "ObjectField", 0, nil},
		[]interface{}{"enter", "Name", "Name", "ObjectField"},
		[]interface{}{"leave", "Name", "Name", "ObjectField"},
		[]interface{}{"enter", "StringValue", "Value", "ObjectField"},
		[]interface{}{"leave", "StringValue", "Value", "ObjectField"},
		[]interface{}{"leave", "ObjectField", 0, nil},
		[]interface{}{"leave", "ObjectValue", "Value", "Argument"},
		[]interface{}{"leave", "Argument", 2, nil},
		[]interface{}{"leave", "Field", 0, nil},
		[]interface{}{"leave", "SelectionSet", "SelectionSet", "FragmentDefinition"},
		[]interface{}{"leave", "FragmentDefinition", 2, nil},
		[]interface{}{"enter", "OperationDefinition", 3, nil},
		[]interface{}{"enter", "SelectionSet", "SelectionSet", "OperationDefinition"},
		[]interface{}{"enter", "Field", 0, nil},
		[]interface{}{"enter", "Name", "Name", "Field"},
		[]interface{}{"leave", "Name", "Name", "Field"},
		[]interface{}{"enter", "Argument", 0, nil},
		[]interface{}{"enter", "Name", "Name", "Argument"},
		[]interface{}{"leave", "Name", "Name", "Argument"},
		[]interface{}{"enter", "BooleanValue", "Value", "Argument"},
		[]interface{}{"leave", "BooleanValue", "Value", "Argument"},
		[]interface{}{"leave", "Argument", 0, nil},
		[]interface{}{"enter", "Argument", 1, nil},
		[]interface{}{"enter", "Name", "Name", "Argument"},
		[]interface{}{"leave", "Name", "Name", "Argument"},
		[]interface{}{"enter", "BooleanValue", "Value", "Argument"},
		[]interface{}{"leave", "BooleanValue", "Value", "Argument"},
		[]interface{}{"leave", "Argument", 1, nil},
		[]interface{}{"leave", "Field", 0, nil},
		[]interface{}{"enter", "Field", 1, nil},
		[]interface{}{"enter", "Name", "Name", "Field"},
		[]interface{}{"leave", "Name", "Name", "Field"},
		[]interface{}{"leave", "Field", 1, nil},
		[]interface{}{"leave", "SelectionSet", "SelectionSet", "OperationDefinition"},
		[]interface{}{"leave", "OperationDefinition", 3, nil},
		[]interface{}{"leave", "Document", nil, nil},
	}

	v := &visitor.VisitorOptions{
		Enter: func(p visitor.VisitFuncParams) (string, interface{}) {
			switch node := p.Node.(type) {
			case ast.Node:
				if p.Parent != nil {
					visited = append(visited, []interface{}{"enter", node.GetKind(), p.Key, p.Parent.GetKind()})
				} else {
					visited = append(visited, []interface{}{"enter", node.GetKind(), p.Key, nil})
				}
			}
			return visitor.ActionNoChange, nil
		},
		Leave: func(p visitor.VisitFuncParams) (string, interface{}) {
			switch node := p.Node.(type) {
			case ast.Node:
				if p.Parent != nil {
					visited = append(visited, []interface{}{"leave", node.GetKind(), p.Key, p.Parent.GetKind()})
				} else {
					visited = append(visited, []interface{}{"leave", node.GetKind(), p.Key, nil})
				}
			}
			return visitor.ActionNoChange, nil
		},
	}

	_ = visitor.Visit(astDoc, v, nil)

	if !reflect.DeepEqual(visited, expectedVisited) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedVisited, visited))
	}
}
