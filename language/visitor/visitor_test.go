package visitor_test

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/sprucehealth/graphql/language/ast"
	"github.com/sprucehealth/graphql/language/parser"
	"github.com/sprucehealth/graphql/language/visitor"
	"github.com/sprucehealth/graphql/testutil"
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

	_ = visitor.Visit(astDoc, v)

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

	_ = visitor.Visit(astDoc, v)

	if !reflect.DeepEqual(visited, expectedVisited) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedVisited, visited))
	}
}

func TestVisitor_VisitsKitchenSink(t *testing.T) {
	t.Skip("This test seems bad")

	b, err := ioutil.ReadFile("../../kitchen-sink.graphql")
	if err != nil {
		t.Fatalf("unable to load kitchen-sink.graphql")
	}

	query := string(b)
	astDoc := parse(t, query)

	visited := []interface{}{}
	expectedVisited := []interface{}{
		[]interface{}{"enter", "Document", nil},
		[]interface{}{"enter", "OperationDefinition", nil},
		[]interface{}{"enter", "Name", "OperationDefinition"},
		[]interface{}{"leave", "Name", "OperationDefinition"},
		[]interface{}{"enter", "VariableDefinition", nil},
		[]interface{}{"enter", "Variable", "VariableDefinition"},
		[]interface{}{"enter", "Name", "Variable"},
		[]interface{}{"leave", "Name", "Variable"},
		[]interface{}{"leave", "Variable", "VariableDefinition"},
		[]interface{}{"enter", "Named", "VariableDefinition"},
		[]interface{}{"enter", "Name", "Named"},
		[]interface{}{"leave", "Name", "Named"},
		[]interface{}{"leave", "Named", "VariableDefinition"},
		[]interface{}{"leave", "VariableDefinition", nil},
		[]interface{}{"enter", "VariableDefinition", nil},
		[]interface{}{"enter", "Variable", "VariableDefinition"},
		[]interface{}{"enter", "Name", "Variable"},
		[]interface{}{"leave", "Name", "Variable"},
		[]interface{}{"leave", "Variable", "VariableDefinition"},
		[]interface{}{"enter", "Named", "VariableDefinition"},
		[]interface{}{"enter", "Name", "Named"},
		[]interface{}{"leave", "Name", "Named"},
		[]interface{}{"leave", "Named", "VariableDefinition"},
		[]interface{}{"enter", "EnumValue", "VariableDefinition"},
		[]interface{}{"leave", "EnumValue", "VariableDefinition"},
		[]interface{}{"leave", "VariableDefinition", nil},
		[]interface{}{"enter", "SelectionSet", "OperationDefinition"},
		[]interface{}{"enter", "Field", nil},
		[]interface{}{"enter", "Name", "Field"},
		[]interface{}{"leave", "Name", "Field"},
		[]interface{}{"enter", "Name", "Field"},
		[]interface{}{"leave", "Name", "Field"},
		[]interface{}{"enter", "Argument", nil},
		[]interface{}{"enter", "Name", "Argument"},
		[]interface{}{"leave", "Name", "Argument"},
		[]interface{}{"enter", "ListValue", "Argument"},
		[]interface{}{"enter", "IntValue", nil},
		[]interface{}{"leave", "IntValue", nil},
		[]interface{}{"enter", "IntValue", nil},
		[]interface{}{"leave", "IntValue", nil},
		[]interface{}{"leave", "ListValue", "Argument"},
		[]interface{}{"leave", "Argument", nil},
		[]interface{}{"enter", "SelectionSet", "Field"},
		[]interface{}{"enter", "Field", nil},
		[]interface{}{"enter", "Name", "Field"},
		[]interface{}{"leave", "Name", "Field"},
		[]interface{}{"leave", "Field", nil},
		[]interface{}{"enter", "InlineFragment", nil},
		[]interface{}{"enter", "Named", "InlineFragment"},
		[]interface{}{"enter", "Name", "Named"},
		[]interface{}{"leave", "Name", "Named"},
		[]interface{}{"leave", "Named", "InlineFragment"},
		[]interface{}{"enter", "Directive", nil},
		[]interface{}{"enter", "Name", "Directive"},
		[]interface{}{"leave", "Name", "Directive"},
		[]interface{}{"leave", "Directive", nil},
		[]interface{}{"enter", "SelectionSet", "InlineFragment"},
		[]interface{}{"enter", "Field", nil},
		[]interface{}{"enter", "Name", "Field"},
		[]interface{}{"leave", "Name", "Field"},
		[]interface{}{"enter", "SelectionSet", "Field"},
		[]interface{}{"enter", "Field", nil},
		[]interface{}{"enter", "Name", "Field"},
		[]interface{}{"leave", "Name", "Field"},
		[]interface{}{"leave", "Field", nil},
		[]interface{}{"enter", "Field", nil},
		[]interface{}{"enter", "Name", "Field"},
		[]interface{}{"leave", "Name", "Field"},
		[]interface{}{"enter", "Name", "Field"},
		[]interface{}{"leave", "Name", "Field"},
		[]interface{}{"enter", "Argument", nil},
		[]interface{}{"enter", "Name", "Argument"},
		[]interface{}{"leave", "Name", "Argument"},
		[]interface{}{"enter", "IntValue", "Argument"},
		[]interface{}{"leave", "IntValue", "Argument"},
		[]interface{}{"leave", "Argument", nil},
		[]interface{}{"enter", "Argument", nil},
		[]interface{}{"enter", "Name", "Argument"},
		[]interface{}{"leave", "Name", "Argument"},
		[]interface{}{"enter", "Variable", "Argument"},
		[]interface{}{"enter", "Name", "Variable"},
		[]interface{}{"leave", "Name", "Variable"},
		[]interface{}{"leave", "Variable", "Argument"},
		[]interface{}{"leave", "Argument", nil},
		[]interface{}{"enter", "Directive", nil},
		[]interface{}{"enter", "Name", "Directive"},
		[]interface{}{"leave", "Name", "Directive"},
		[]interface{}{"enter", "Argument", nil},
		[]interface{}{"enter", "Name", "Argument"},
		[]interface{}{"leave", "Name", "Argument"},
		[]interface{}{"enter", "Variable", "Argument"},
		[]interface{}{"enter", "Name", "Variable"},
		[]interface{}{"leave", "Name", "Variable"},
		[]interface{}{"leave", "Variable", "Argument"},
		[]interface{}{"leave", "Argument", nil},
		[]interface{}{"leave", "Directive", nil},
		[]interface{}{"enter", "SelectionSet", "Field"},
		[]interface{}{"enter", "Field", nil},
		[]interface{}{"enter", "Name", "Field"},
		[]interface{}{"leave", "Name", "Field"},
		[]interface{}{"leave", "Field", nil},
		[]interface{}{"enter", "FragmentSpread", nil},
		[]interface{}{"enter", "Name", "FragmentSpread"},
		[]interface{}{"leave", "Name", "FragmentSpread"},
		[]interface{}{"leave", "FragmentSpread", nil},
		[]interface{}{"leave", "SelectionSet", "Field"},
		[]interface{}{"leave", "Field", nil},
		[]interface{}{"leave", "SelectionSet", "Field"},
		[]interface{}{"leave", "Field", nil},
		[]interface{}{"leave", "SelectionSet", "InlineFragment"},
		[]interface{}{"leave", "InlineFragment", nil},
		[]interface{}{"leave", "SelectionSet", "Field"},
		[]interface{}{"leave", "Field", nil},
		[]interface{}{"leave", "SelectionSet", "OperationDefinition"},
		[]interface{}{"leave", "OperationDefinition", nil},
		[]interface{}{"enter", "OperationDefinition", nil},
		[]interface{}{"enter", "Name", "OperationDefinition"},
		[]interface{}{"leave", "Name", "OperationDefinition"},
		[]interface{}{"enter", "SelectionSet", "OperationDefinition"},
		[]interface{}{"enter", "Field", nil},
		[]interface{}{"enter", "Name", "Field"},
		[]interface{}{"leave", "Name", "Field"},
		[]interface{}{"enter", "Argument", nil},
		[]interface{}{"enter", "Name", "Argument"},
		[]interface{}{"leave", "Name", "Argument"},
		[]interface{}{"enter", "IntValue", "Argument"},
		[]interface{}{"leave", "IntValue", "Argument"},
		[]interface{}{"leave", "Argument", nil},
		[]interface{}{"enter", "Directive", nil},
		[]interface{}{"enter", "Name", "Directive"},
		[]interface{}{"leave", "Name", "Directive"},
		[]interface{}{"leave", "Directive", nil},
		[]interface{}{"enter", "SelectionSet", "Field"},
		[]interface{}{"enter", "Field", nil},
		[]interface{}{"enter", "Name", "Field"},
		[]interface{}{"leave", "Name", "Field"},
		[]interface{}{"enter", "SelectionSet", "Field"},
		[]interface{}{"enter", "Field", nil},
		[]interface{}{"enter", "Name", "Field"},
		[]interface{}{"leave", "Name", "Field"},
		[]interface{}{"leave", "Field", nil},
		[]interface{}{"leave", "SelectionSet", "Field"},
		[]interface{}{"leave", "Field", nil},
		[]interface{}{"leave", "SelectionSet", "Field"},
		[]interface{}{"leave", "Field", nil},
		[]interface{}{"leave", "SelectionSet", "OperationDefinition"},
		[]interface{}{"leave", "OperationDefinition", nil},
		[]interface{}{"enter", "FragmentDefinition", nil},
		[]interface{}{"enter", "Name", "FragmentDefinition"},
		[]interface{}{"leave", "Name", "FragmentDefinition"},
		[]interface{}{"enter", "Named", "FragmentDefinition"},
		[]interface{}{"enter", "Name", "Named"},
		[]interface{}{"leave", "Name", "Named"},
		[]interface{}{"leave", "Named", "FragmentDefinition"},
		[]interface{}{"enter", "SelectionSet", "FragmentDefinition"},
		[]interface{}{"enter", "Field", nil},
		[]interface{}{"enter", "Name", "Field"},
		[]interface{}{"leave", "Name", "Field"},
		[]interface{}{"enter", "Argument", nil},
		[]interface{}{"enter", "Name", "Argument"},
		[]interface{}{"leave", "Name", "Argument"},
		[]interface{}{"enter", "Variable", "Argument"},
		[]interface{}{"enter", "Name", "Variable"},
		[]interface{}{"leave", "Name", "Variable"},
		[]interface{}{"leave", "Variable", "Argument"},
		[]interface{}{"leave", "Argument", nil},
		[]interface{}{"enter", "Argument", nil},
		[]interface{}{"enter", "Name", "Argument"},
		[]interface{}{"leave", "Name", "Argument"},
		[]interface{}{"enter", "Variable", "Argument"},
		[]interface{}{"enter", "Name", "Variable"},
		[]interface{}{"leave", "Name", "Variable"},
		[]interface{}{"leave", "Variable", "Argument"},
		[]interface{}{"leave", "Argument", nil},
		[]interface{}{"enter", "Argument", nil},
		[]interface{}{"enter", "Name", "Argument"},
		[]interface{}{"leave", "Name", "Argument"},
		[]interface{}{"enter", "ObjectValue", "Argument"},
		[]interface{}{"enter", "ObjectField", nil},
		[]interface{}{"enter", "Name", "ObjectField"},
		[]interface{}{"leave", "Name", "ObjectField"},
		[]interface{}{"enter", "StringValue", "ObjectField"},
		[]interface{}{"leave", "StringValue", "ObjectField"},
		[]interface{}{"leave", "ObjectField", nil},
		[]interface{}{"leave", "ObjectValue", "Argument"},
		[]interface{}{"leave", "Argument", nil},
		[]interface{}{"leave", "Field", nil},
		[]interface{}{"leave", "SelectionSet", "FragmentDefinition"},
		[]interface{}{"leave", "FragmentDefinition", nil},
		[]interface{}{"enter", "OperationDefinition", nil},
		[]interface{}{"enter", "SelectionSet", "OperationDefinition"},
		[]interface{}{"enter", "Field", nil},
		[]interface{}{"enter", "Name", "Field"},
		[]interface{}{"leave", "Name", "Field"},
		[]interface{}{"enter", "Argument", nil},
		[]interface{}{"enter", "Name", "Argument"},
		[]interface{}{"leave", "Name", "Argument"},
		[]interface{}{"enter", "BooleanValue", "Argument"},
		[]interface{}{"leave", "BooleanValue", "Argument"},
		[]interface{}{"leave", "Argument", nil},
		[]interface{}{"enter", "Argument", nil},
		[]interface{}{"enter", "Name", "Argument"},
		[]interface{}{"leave", "Name", "Argument"},
		[]interface{}{"enter", "BooleanValue", "Argument"},
		[]interface{}{"leave", "BooleanValue", "Argument"},
		[]interface{}{"leave", "Argument", nil},
		[]interface{}{"leave", "Field", nil},
		[]interface{}{"enter", "Field", nil},
		[]interface{}{"enter", "Name", "Field"},
		[]interface{}{"leave", "Name", "Field"},
		[]interface{}{"leave", "Field", nil},
		[]interface{}{"leave", "SelectionSet", "OperationDefinition"},
		[]interface{}{"leave", "OperationDefinition", nil},
		[]interface{}{"leave", "Document", nil},
	}

	v := &visitor.VisitorOptions{
		Enter: func(p visitor.VisitFuncParams) (string, interface{}) {
			switch node := p.Node.(type) {
			case ast.Node:
				if p.Parent != nil {
					visited = append(visited, []interface{}{"enter", node.GetKind(), p.Parent.GetKind()})
				} else {
					visited = append(visited, []interface{}{"enter", node.GetKind(), nil})
				}
			}
			return visitor.ActionNoChange, nil
		},
		Leave: func(p visitor.VisitFuncParams) (string, interface{}) {
			switch node := p.Node.(type) {
			case ast.Node:
				if p.Parent != nil {
					visited = append(visited, []interface{}{"leave", node.GetKind(), p.Parent.GetKind()})
				} else {
					visited = append(visited, []interface{}{"leave", node.GetKind(), nil})
				}
			}
			return visitor.ActionNoChange, nil
		},
	}

	_ = visitor.Visit(astDoc, v)

	if !reflect.DeepEqual(visited, expectedVisited) {
		for i, v := range visited {
			if !reflect.DeepEqual(v, expectedVisited[i]) {
				t.Logf("%d    %v != %v", i, v, expectedVisited[i])
				break
			}
		}
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedVisited, visited))
	}
}
