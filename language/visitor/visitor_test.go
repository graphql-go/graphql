package visitor_test

import (
	"github.com/chris-ramon/graphql-go/language/ast"
	"github.com/chris-ramon/graphql-go/language/parser"
	"github.com/chris-ramon/graphql-go/language/visitor"
	"github.com/chris-ramon/graphql-go/testutil"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
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

// Helper functions to get map value by key
// Allow keys to specify dot paths (e.g `Name.Value`)
func getMapValue(m map[string]interface{}, key string) interface{} {
	tokens := strings.Split(key, ".")
	valMap := m
	for _, token := range tokens {
		v, ok := valMap[token]
		if !ok {
			return nil
		}
		switch v := v.(type) {
		case []interface{}:
			return v
		case map[string]interface{}:
			valMap = v
			continue
		case string:
			return v
		}
	}
	return valMap
}
func getMapValueString(m map[string]interface{}, key string) string {
	tokens := strings.Split(key, ".")
	valMap := m
	for _, token := range tokens {
		v, ok := valMap[token]
		if !ok {
			return ""
		}
		switch v := v.(type) {
		case map[string]interface{}:
			valMap = v
			continue
		case string:
			return v
		}
	}
	return ""
}
func TestVisitor_AllowsForEditingOnEnter(t *testing.T) {

	query := `{ a, b, c { a, b, c } }`
	astDoc := parse(t, query)

	expectedQuery := `{ a,    c { a,    c } }`
	expectedAST := parse(t, expectedQuery)
	v := &visitor.VisitorOptions{
		Enter: func(p visitor.VisitFuncParams) (string, interface{}) {
			switch node := p.Node.(type) {
			case map[string]interface{}:
				if getMapValueString(node, "Kind") == "Field" && getMapValueString(node, "Name.Value") == "b" {
					return visitor.ActionUpdate, nil
				}
			}
			return visitor.ActionNoChange, nil
		},
	}

	editedAst := visitor.Visit(astDoc, v, nil)
	if !reflect.DeepEqual(testutil.ASTToJSON(t, expectedAST), editedAst) {
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
			case map[string]interface{}:
				if getMapValueString(node, "Kind") == "Field" && getMapValueString(node, "Name.Value") == "b" {
					return visitor.ActionUpdate, nil
				}
			}
			return visitor.ActionNoChange, nil
		},
	}

	editedAst := visitor.Visit(astDoc, v, nil)
	if !reflect.DeepEqual(testutil.ASTToJSON(t, expectedAST), editedAst) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedAST, editedAst))
	}
}

func TestVisitor_VisitsEditedNode(t *testing.T) {

	query := `{ a { x } }`
	astDoc := parse(t, query)

	addedField := map[string]interface{}{
		"Kind": "Field",
		"Name": map[string]interface{}{
			"Kind":  "Name",
			"Value": "__typename",
		},
	}

	didVisitAddedField := false
	v := &visitor.VisitorOptions{
		Enter: func(p visitor.VisitFuncParams) (string, interface{}) {
			switch node := p.Node.(type) {
			case map[string]interface{}:
				if getMapValueString(node, "Kind") == "Field" && getMapValueString(node, "Name.Value") == "a" {
					s := getMapValue(node, "SelectionSet.Selections").([]interface{})
					s = append(s, addedField)
					return visitor.ActionUpdate, map[string]interface{}{
						"Kind":         "Field",
						"SelectionSet": s,
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
			case map[string]interface{}:
				visited = append(visited, []interface{}{"enter", getMapValue(node, "Kind"), getMapValue(node, "Value")})
				if getMapValueString(node, "Kind") == "Field" && getMapValueString(node, "Name.Value") == "b" {
					return visitor.ActionSkip, nil
				}
			}
			return visitor.ActionNoChange, nil
		},
		Leave: func(p visitor.VisitFuncParams) (string, interface{}) {
			switch node := p.Node.(type) {
			case map[string]interface{}:
				visited = append(visited, []interface{}{"leave", getMapValue(node, "Kind"), getMapValue(node, "Value")})
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
					case map[string]interface{}:
						visited = append(visited, []interface{}{"enter", getMapValue(node, "Kind"), getMapValue(node, "Value")})
					}
					return visitor.ActionNoChange, nil
				},
			},
			"SelectionSet": visitor.NamedVisitFuncs{
				Enter: func(p visitor.VisitFuncParams) (string, interface{}) {
					switch node := p.Node.(type) {
					case map[string]interface{}:
						visited = append(visited, []interface{}{"enter", getMapValue(node, "Kind"), getMapValue(node, "Value")})
					}
					return visitor.ActionNoChange, nil
				},
				Leave: func(p visitor.VisitFuncParams) (string, interface{}) {
					switch node := p.Node.(type) {
					case map[string]interface{}:
						visited = append(visited, []interface{}{"leave", getMapValue(node, "Kind"), getMapValue(node, "Value")})
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
	b, err := ioutil.ReadFile("./../parser/kitchen-sink.graphql")
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
		[]interface{}{"enter", "NamedType", "Type", "VariableDefinition"},
		[]interface{}{"enter", "Name", "Name", "NamedType"},
		[]interface{}{"leave", "Name", "Name", "NamedType"},
		[]interface{}{"leave", "NamedType", "Type", "VariableDefinition"},
		[]interface{}{"leave", "VariableDefinition", 0, nil},
		[]interface{}{"enter", "VariableDefinition", 1, nil},
		[]interface{}{"enter", "Variable", "Variable", "VariableDefinition"},
		[]interface{}{"enter", "Name", "Name", "Variable"},
		[]interface{}{"leave", "Name", "Name", "Variable"},
		[]interface{}{"leave", "Variable", "Variable", "VariableDefinition"},
		[]interface{}{"enter", "NamedType", "Type", "VariableDefinition"},
		[]interface{}{"enter", "Name", "Name", "NamedType"},
		[]interface{}{"leave", "Name", "Name", "NamedType"},
		[]interface{}{"leave", "NamedType", "Type", "VariableDefinition"},
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
		[]interface{}{"enter", "NamedType", "TypeCondition", "InlineFragment"},
		[]interface{}{"enter", "Name", "Name", "NamedType"},
		[]interface{}{"leave", "Name", "Name", "NamedType"},
		[]interface{}{"leave", "NamedType", "TypeCondition", "InlineFragment"},
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
		[]interface{}{"enter", "NamedType", "TypeCondition", "FragmentDefinition"},
		[]interface{}{"enter", "Name", "Name", "NamedType"},
		[]interface{}{"leave", "Name", "Name", "NamedType"},
		[]interface{}{"leave", "NamedType", "TypeCondition", "FragmentDefinition"},
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
			case map[string]interface{}:

				parent, ok := p.Parent.(map[string]interface{})
				if !ok {
					parent = map[string]interface{}{}
				}

				visited = append(visited, []interface{}{"enter", getMapValue(node, "Kind"), p.Key, getMapValue(parent, "Kind")})
			}
			return visitor.ActionNoChange, nil
		},
		Leave: func(p visitor.VisitFuncParams) (string, interface{}) {
			switch node := p.Node.(type) {
			case map[string]interface{}:

				parent, ok := p.Parent.(map[string]interface{})
				if !ok {
					parent = map[string]interface{}{}
				}

				visited = append(visited, []interface{}{"leave", getMapValue(node, "Kind"), p.Key, getMapValue(parent, "Kind")})
			}
			return visitor.ActionNoChange, nil
		},
	}

	_ = visitor.Visit(astDoc, v, nil)

	if !reflect.DeepEqual(visited, expectedVisited) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedVisited, visited))
	}
}

func TestVisitor_ProducesHelpfulErrorMessages(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			err := r.(string)
			expectedErr := `Invalid AST Node (4): map[random:Data]`
			if err != expectedErr {
				t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(err, expectedErr))
			}
			return
		}
		t.Fatalf("expected to panic")
	}()
	astDoc := map[string]interface{}{
		"random": "Data",
	}
	_ = visitor.Visit(astDoc, nil, nil)
}
