package printer

import (
	"fmt"
	"testing"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/visitor"
)

func TestGetMapValue_Internal(t *testing.T) {
	// nested map continuation case (map[string]interface{} -> continue)
	m := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{},
		},
	}
	// token "a" maps to {"b": {}} -> maps case, continue to token "b"
	// token "b" maps to {} -> maps case, continue, loop ends -> fallback return valMap = {}
	result := getMapValue(m, "a.b")
	if _, ok := result.(map[string]interface{}); !ok {
		t.Fatalf("expected map, got %T: %v", result, result)
	}
}

func TestGetMapSliceValue_Internal(t *testing.T) {
	// key not found branch
	m := map[string]interface{}{"existing": "value"}
	result := getMapSliceValue(m, "nonexistent")
	if len(result) != 0 {
		t.Fatalf("expected empty slice, got %v", result)
	}

	// fallback return (key exists but value is not a slice)
	m2 := map[string]interface{}{"key": "stringValue"}
	result2 := getMapSliceValue(m2, "key")
	if len(result2) != 0 {
		t.Fatalf("expected empty slice, got %v", result2)
	}
	// fallback return at end of function (multi-token where last is non-slice)
	m3 := map[string]interface{}{
		"a": map[string]interface{}{
			"b": "notslice",
		},
	}
	result3 := getMapSliceValue(m3, "a.b")
	if len(result3) != 0 {
		t.Fatalf("expected empty slice, got %v", result3)
	}
}

func TestGetMapValueString_Internal(t *testing.T) {
	// key not found branch
	m := map[string]interface{}{"existing": "value"}
	result := getMapValueString(m, "nonexistent")
	if result != "" {
		t.Fatalf("expected empty string, got %v", result)
	}

	// nil value branch
	m2 := map[string]interface{}{"key": nil}
	result2 := getMapValueString(m2, "key")
	if result2 != "" {
		t.Fatalf("expected empty string, got %v", result2)
	}

	// default fmt.Sprintf case (int value)
	m3 := map[string]interface{}{"key": 42}
	result3 := getMapValueString(m3, "key")
	if result3 != "42" {
		t.Fatalf("expected '42', got %v", result3)
	}

	// fallback return at end (multi-token where last is unknown)
	m4 := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{},
		},
	}
	result4 := getMapValueString(m4, "a.b")
	if result4 != "" {
		t.Fatalf("expected empty string, got %v", result4)
	}
}

func TestGetDescription_DescribableNode(t *testing.T) {
	// DescribableNode with description
	node := ast.NewScalarDefinition(&ast.ScalarDefinition{
		Name: ast.NewName(&ast.Name{Value: "Foo"}),
		Description: ast.NewStringValue(&ast.StringValue{
			Value: "test description",
		}),
	})
	desc := getDescription(node)
	if desc != `"""test description"""` {
		t.Fatalf("unexpected description: %v", desc)
	}

	// DescribableNode with nil description
	nodeNoDesc := ast.NewScalarDefinition(&ast.ScalarDefinition{
		Name: ast.NewName(&ast.Name{Value: "Foo"}),
	})
	desc2 := getDescription(nodeNoDesc)
	if desc2 != "" {
		t.Fatalf("expected empty description, got %v", desc2)
	}

	// DescribableNode with multiline description
	nodeMultiLine := ast.NewFieldDefinition(&ast.FieldDefinition{
		Name: ast.NewName(&ast.Name{Value: "foo"}),
		Description: ast.NewStringValue(&ast.StringValue{
			Value: "line1\nline2",
		}),
	})
	desc3 := getDescription(nodeMultiLine)
	expected := "\"\"\"\nline1\nline2\n\"\"\""
	if desc3 != expected {
		t.Fatalf("unexpected multiline description:\ngot:  %q\nwant: %q", desc3, expected)
	}
}

func TestToSliceString_Internal(t *testing.T) {
	// nil case
	result := toSliceString(nil)
	if len(result) != 0 {
		t.Fatalf("expected empty slice, got %v", result)
	}

	// non-slice default case
	result2 := toSliceString("not a slice")
	if len(result2) != 0 {
		t.Fatalf("expected empty slice, got %v", result2)
	}

	// slice with non-string elements
	result3 := toSliceString([]int{1, 2, 3})
	if len(result3) != 0 {
		t.Fatalf("expected empty slice, got %v", result3)
	}

	// slice with mixed elements
	result4 := toSliceString([]interface{}{"a", 2, "c"})
	if len(result4) != 2 {
		t.Fatalf("expected 2 strings, got %v", result4)
	}
	if result4[0] != "a" || result4[1] != "c" {
		t.Fatalf("unexpected result: %v", result4)
	}
}

func TestIndent_Internal(t *testing.T) {
	// nil case
	result := indent(nil)
	if result != "" {
		t.Fatalf("expected empty string, got %v", result)
	}

	// non-string type
	result2 := indent(42)
	if result2 != "" {
		t.Fatalf("expected empty string, got %v", result2)
	}

	// normal string case
	result3 := indent("hello\nworld")
	if result3 != "hello\n  world" {
		t.Fatalf("unexpected result: %v", result3)
	}
}

func TestPrint_PanicRecover(t *testing.T) {
	// Replace a reducer with one that panics to test Print's recover path
	originalFn := printDocASTReducer["Name"]
	printDocASTReducer["Name"] = func(p visitor.VisitFuncParams) (string, interface{}) {
		panic("simulated panic")
	}
	defer func() {
		printDocASTReducer["Name"] = originalFn
	}()

	doc := ast.NewDocument(&ast.Document{
		Definitions: []ast.Node{
			ast.NewOperationDefinition(&ast.OperationDefinition{
				Operation: ast.OperationTypeQuery,
				SelectionSet: ast.NewSelectionSet(&ast.SelectionSet{
					Selections: []ast.Selection{
						ast.NewField(&ast.Field{
							Name: ast.NewName(&ast.Name{Value: "foo"}),
						}),
					},
				}),
			}),
		},
	})

	_ = Print(doc)
}

func TestPrintDocASTReducer_Name(t *testing.T) {
	fn := printDocASTReducer["Name"]

	// *ast.Name path
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewName(&ast.Name{Value: "hello"}),
	})
	if action != visitor.ActionUpdate || val != "hello" {
		t.Fatalf("unexpected: action=%v val=%v", action, val)
	}

	// map path
	action, val = fn(visitor.VisitFuncParams{
		Node: map[string]interface{}{"Value": "world"},
	})
	if action != visitor.ActionUpdate || val != "world" {
		t.Fatalf("unexpected: action=%v val=%v", action, val)
	}

	// ActionNoChange
	action, val = fn(visitor.VisitFuncParams{Node: 42})
	if action != visitor.ActionNoChange || val != nil {
		t.Fatalf("unexpected: action=%v val=%v", action, val)
	}
}

func TestPrintDocASTReducer_Variable(t *testing.T) {
	fn := printDocASTReducer["Variable"]

	// *ast.Variable path
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewVariable(&ast.Variable{Name: ast.NewName(&ast.Name{Value: "foo"})}),
	})
	if action != visitor.ActionUpdate {
		t.Fatalf("unexpected: action=%v val=%v", action, val)
	}

	// ActionNoChange
	action, val = fn(visitor.VisitFuncParams{Node: 42})
	if action != visitor.ActionNoChange || val != nil {
		t.Fatalf("unexpected: action=%v val=%v", action, val)
	}
}

func TestPrintDocASTReducer_Document(t *testing.T) {
	fn := printDocASTReducer["Document"]

	// *ast.Document path
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewDocument(&ast.Document{}),
	})
	if action != visitor.ActionUpdate || val != "\n" {
		t.Fatalf("unexpected: action=%v val=%v", action, val)
	}

	// ActionNoChange
	action, val = fn(visitor.VisitFuncParams{Node: 42})
	if action != visitor.ActionNoChange || val != nil {
		t.Fatalf("unexpected: action=%v val=%v", action, val)
	}
}

func TestPrintDocASTReducer_OperationDefinition(t *testing.T) {
	fn := printDocASTReducer["OperationDefinition"]

	// *ast.OperationDefinition path (query, no name, no vars, no directives)
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewOperationDefinition(&ast.OperationDefinition{
			Operation: ast.OperationTypeQuery,
			SelectionSet: ast.NewSelectionSet(&ast.SelectionSet{
				Selections: []ast.Selection{
					ast.NewField(&ast.Field{
						Name: ast.NewName(&ast.Name{Value: "id"}),
					}),
				},
			}),
		}),
	})
	if action != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action)
	}
	_ = val

	// *ast.OperationDefinition path (mutation with name and directives)
	mutFnDef := ast.NewField(&ast.Field{
		Name: ast.NewName(&ast.Name{Value: "doStuff"}),
		SelectionSet: ast.NewSelectionSet(&ast.SelectionSet{
			Selections: []ast.Selection{
				ast.NewField(&ast.Field{
					Name: ast.NewName(&ast.Name{Value: "result"}),
				}),
			},
		}),
	})
	selSet := ast.NewSelectionSet(&ast.SelectionSet{
		Selections: []ast.Selection{mutFnDef},
	})
	action2, val2 := fn(visitor.VisitFuncParams{
		Node: ast.NewOperationDefinition(&ast.OperationDefinition{
			Operation:     ast.OperationTypeMutation,
			Name:          ast.NewName(&ast.Name{Value: "myMutation"}),
			SelectionSet:  selSet,
		}),
	})
	if action2 != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action2)
	}
	_ = val2

	// *ast.OperationDefinition path (subscription)
	subSel := ast.NewSelectionSet(&ast.SelectionSet{
		Selections: []ast.Selection{
			ast.NewField(&ast.Field{
				Name: ast.NewName(&ast.Name{Value: "x"}),
			}),
		},
	})
	action3, val3 := fn(visitor.VisitFuncParams{
		Node: ast.NewOperationDefinition(&ast.OperationDefinition{
			Operation:    ast.OperationTypeSubscription,
			SelectionSet: subSel,
		}),
	})
	if action3 != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action3)
	}
	_ = val3

	// ActionNoChange
	action4, val4 := fn(visitor.VisitFuncParams{Node: 42})
	if action4 != visitor.ActionNoChange || val4 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action4, val4)
	}
}

func TestPrintDocASTReducer_VariableDefinition(t *testing.T) {
	fn := printDocASTReducer["VariableDefinition"]

	// *ast.VariableDefinition path
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewVariableDefinition(&ast.VariableDefinition{
			Variable: ast.NewVariable(&ast.Variable{Name: ast.NewName(&ast.Name{Value: "x"})}),
			Type:     ast.NewNamed(&ast.Named{Name: ast.NewName(&ast.Name{Value: "Int"})}),
		}),
	})
	if action != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action)
	}
	_ = val

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_SelectionSet(t *testing.T) {
	fn := printDocASTReducer["SelectionSet"]

	// *ast.SelectionSet path with empty selections
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewSelectionSet(&ast.SelectionSet{}),
	})
	if action != visitor.ActionUpdate || val != "{}" {
		t.Fatalf("unexpected: action=%v val=%v", action, val)
	}

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_Field(t *testing.T) {
	fn := printDocASTReducer["Field"]

	// *ast.Argument path (dead code path under "Field" key)
	arg := ast.NewArgument(&ast.Argument{
		Name:  ast.NewName(&ast.Name{Value: "arg1"}),
		Value: ast.NewStringValue(&ast.StringValue{Value: "val1"}),
	})
	action, val := fn(visitor.VisitFuncParams{
		Node: arg,
	})
	if action != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action)
	}
	_ = val

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_Argument(t *testing.T) {
	fn := printDocASTReducer["Argument"]

	// *ast.FragmentSpread path (dead code path under "Argument" key)
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewFragmentSpread(&ast.FragmentSpread{
			Name: ast.NewName(&ast.Name{Value: "myFrag"}),
		}),
	})
	if action != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action)
	}
	_ = val

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_FragmentSpread(t *testing.T) {
	fn := printDocASTReducer["FragmentSpread"]

	// *ast.InlineFragment path (dead code path under "FragmentSpread" key)
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewInlineFragment(&ast.InlineFragment{
			TypeCondition: ast.NewNamed(&ast.Named{Name: ast.NewName(&ast.Name{Value: "User"})}),
			SelectionSet: ast.NewSelectionSet(&ast.SelectionSet{
				Selections: []ast.Selection{
					ast.NewField(&ast.Field{
						Name: ast.NewName(&ast.Name{Value: "id"}),
					}),
				},
			}),
		}),
	})
	if action != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action)
	}
	_ = val

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_InlineFragment(t *testing.T) {
	fn := printDocASTReducer["InlineFragment"]

	// ActionNoChange (map path already covered by kitchen-sink)
	action, val := fn(visitor.VisitFuncParams{Node: 42})
	if action != visitor.ActionNoChange || val != nil {
		t.Fatalf("unexpected: action=%v val=%v", action, val)
	}
}

func TestPrintDocASTReducer_FragmentDefinition(t *testing.T) {
	fn := printDocASTReducer["FragmentDefinition"]

	// *ast.FragmentDefinition path
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewFragmentDefinition(&ast.FragmentDefinition{
			Name: ast.NewName(&ast.Name{Value: "myFrag"}),
			TypeCondition: ast.NewNamed(&ast.Named{Name: ast.NewName(&ast.Name{Value: "User"})}),
			SelectionSet: ast.NewSelectionSet(&ast.SelectionSet{
				Selections: []ast.Selection{
					ast.NewField(&ast.Field{
						Name: ast.NewName(&ast.Name{Value: "id"}),
					}),
				},
			}),
		}),
	})
	if action != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action)
	}
	_ = val

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_IntValue(t *testing.T) {
	fn := printDocASTReducer["IntValue"]

	// map path
	action, val := fn(visitor.VisitFuncParams{
		Node: map[string]interface{}{"Value": "42"},
	})
	if action != visitor.ActionUpdate || val != "42" {
		t.Fatalf("unexpected: action=%v val=%v", action, val)
	}

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_FloatValue(t *testing.T) {
	fn := printDocASTReducer["FloatValue"]

	// *ast.FloatValue path
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewFloatValue(&ast.FloatValue{Value: "3.14"}),
	})
	if action != visitor.ActionUpdate || val != "3.14" {
		t.Fatalf("unexpected: action=%v val=%v", action, val)
	}

	// map path
	action2, val2 := fn(visitor.VisitFuncParams{
		Node: map[string]interface{}{"Value": "2.718"},
	})
	if action2 != visitor.ActionUpdate || val2 != "2.718" {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}

	// ActionNoChange
	action3, val3 := fn(visitor.VisitFuncParams{Node: 42})
	if action3 != visitor.ActionNoChange || val3 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action3, val3)
	}
}

func TestPrintDocASTReducer_StringValue(t *testing.T) {
	fn := printDocASTReducer["StringValue"]

	// map path
	action, val := fn(visitor.VisitFuncParams{
		Node: map[string]interface{}{"Value": "hello"},
	})
	if action != visitor.ActionUpdate || val != `"hello"` {
		t.Fatalf("unexpected: action=%v val=%v", action, val)
	}

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_BooleanValue(t *testing.T) {
	fn := printDocASTReducer["BooleanValue"]

	// map path
	action, val := fn(visitor.VisitFuncParams{
		Node: map[string]interface{}{"Value": "true"},
	})
	if action != visitor.ActionUpdate || val != "true" {
		t.Fatalf("unexpected: action=%v val=%v", action, val)
	}

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_EnumValue(t *testing.T) {
	fn := printDocASTReducer["EnumValue"]

	// map path
	action, val := fn(visitor.VisitFuncParams{
		Node: map[string]interface{}{"Value": "SOMETHING"},
	})
	if action != visitor.ActionUpdate || val != "SOMETHING" {
		t.Fatalf("unexpected: action=%v val=%v", action, val)
	}

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_ListValue(t *testing.T) {
	fn := printDocASTReducer["ListValue"]

	// *ast.ListValue path
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewListValue(&ast.ListValue{}),
	})
	if action != visitor.ActionUpdate || val != "[]" {
		t.Fatalf("unexpected: action=%v val=%v", action, val)
	}

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_ObjectValue(t *testing.T) {
	fn := printDocASTReducer["ObjectValue"]

	// *ast.ObjectValue path
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewObjectValue(&ast.ObjectValue{}),
	})
	if action != visitor.ActionUpdate || val != "{}" {
		t.Fatalf("unexpected: action=%v val=%v", action, val)
	}

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_ObjectField(t *testing.T) {
	fn := printDocASTReducer["ObjectField"]

	// *ast.ObjectField path
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewObjectField(&ast.ObjectField{
			Name:  ast.NewName(&ast.Name{Value: "key"}),
			Value: ast.NewStringValue(&ast.StringValue{Value: "val"}),
		}),
	})
	if action != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action)
	}
	_ = val

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_Directive(t *testing.T) {
	fn := printDocASTReducer["Directive"]

	// *ast.Directive path
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewDirective(&ast.Directive{
			Name: ast.NewName(&ast.Name{Value: "skip"}),
		}),
	})
	if action != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action)
	}
	_ = val

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_Named(t *testing.T) {
	fn := printDocASTReducer["Named"]

	// *ast.Named path
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewNamed(&ast.Named{Name: ast.NewName(&ast.Name{Value: "String"})}),
	})
	if action != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action)
	}
	_ = val

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_List(t *testing.T) {
	fn := printDocASTReducer["List"]

	// *ast.List path
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewList(&ast.List{Type: ast.NewNamed(&ast.Named{Name: ast.NewName(&ast.Name{Value: "Int"})})}),
	})
	if action != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action)
	}
	_ = val

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_NonNull(t *testing.T) {
	fn := printDocASTReducer["NonNull"]

	// *ast.NonNull path
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewNonNull(&ast.NonNull{Type: ast.NewNamed(&ast.Named{Name: ast.NewName(&ast.Name{Value: "String"})})}),
	})
	if action != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action)
	}
	_ = val

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_SchemaDefinition(t *testing.T) {
	fn := printDocASTReducer["SchemaDefinition"]

	// *ast.SchemaDefinition path (with directives to cover directives loop)
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewSchemaDefinition(&ast.SchemaDefinition{
			Directives: []*ast.Directive{
				ast.NewDirective(&ast.Directive{Name: ast.NewName(&ast.Name{Value: "myDirective"})}),
			},
			OperationTypes: []*ast.OperationTypeDefinition{
				ast.NewOperationTypeDefinition(&ast.OperationTypeDefinition{
					Operation: "QUERY",
					Type:      ast.NewNamed(&ast.Named{Name: ast.NewName(&ast.Name{Value: "QueryType"})}),
				}),
			},
		}),
	})
	if action != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action)
	}
	_ = val

	// map path with directives (to cover map directives loop)
	action3, val3 := fn(visitor.VisitFuncParams{
		Node: map[string]interface{}{
			"OperationTypes": []interface{}{"query: QueryType"},
			"Directives":     []interface{}{"@myDirective"},
		},
	})
	if action3 != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action3)
	}
	_ = val3

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_OperationTypeDefinition(t *testing.T) {
	fn := printDocASTReducer["OperationTypeDefinition"]

	// *ast.OperationTypeDefinition path
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewOperationTypeDefinition(&ast.OperationTypeDefinition{
			Operation: "QUERY",
			Type:      ast.NewNamed(&ast.Named{Name: ast.NewName(&ast.Name{Value: "QueryRoot"})}),
		}),
	})
	if action != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action)
	}
	_ = val

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_ScalarDefinition(t *testing.T) {
	fn := printDocASTReducer["ScalarDefinition"]

	// *ast.ScalarDefinition path with directives and description
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewScalarDefinition(&ast.ScalarDefinition{
			Name: ast.NewName(&ast.Name{Value: "CustomScalar"}),
			Directives: []*ast.Directive{
				ast.NewDirective(&ast.Directive{Name: ast.NewName(&ast.Name{Value: "myDir"})}),
			},
			Description: ast.NewStringValue(&ast.StringValue{Value: "a scalar"}),
		}),
	})
	if action != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action)
	}
	_ = val

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_ObjectDefinition(t *testing.T) {
	fn := printDocASTReducer["ObjectDefinition"]

	// *ast.ObjectDefinition path with directives and description
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewObjectDefinition(&ast.ObjectDefinition{
			Name: ast.NewName(&ast.Name{Value: "Foo"}),
			Directives: []*ast.Directive{
				ast.NewDirective(&ast.Directive{Name: ast.NewName(&ast.Name{Value: "myDir"})}),
			},
			Description: ast.NewStringValue(&ast.StringValue{Value: "an object"}),
		}),
	})
	if action != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action)
	}
	_ = val

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_FieldDefinition(t *testing.T) {
	fn := printDocASTReducer["FieldDefinition"]

	// *ast.FieldDefinition path with directives, description, and args with descriptions (hasArgDesc=true)
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewFieldDefinition(&ast.FieldDefinition{
			Name: ast.NewName(&ast.Name{Value: "foo"}),
			Type: ast.NewNamed(&ast.Named{Name: ast.NewName(&ast.Name{Value: "Int"})}),
			Directives: []*ast.Directive{
				ast.NewDirective(&ast.Directive{Name: ast.NewName(&ast.Name{Value: "myDir"})}),
			},
			Description: ast.NewStringValue(&ast.StringValue{Value: "a field"}),
			Arguments: []*ast.InputValueDefinition{
				ast.NewInputValueDefinition(&ast.InputValueDefinition{
					Name:        ast.NewName(&ast.Name{Value: "arg1"}),
					Type:        ast.NewNamed(&ast.Named{Name: ast.NewName(&ast.Name{Value: "String"})}),
					Description: ast.NewStringValue(&ast.StringValue{Value: "an arg"}),
				}),
			},
		}),
	})
	if action != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action)
	}
	_ = val

	// *ast.FieldDefinition path with args but NO descriptions (hasArgDesc=false, else branch)
	action3, val3 := fn(visitor.VisitFuncParams{
		Node: ast.NewFieldDefinition(&ast.FieldDefinition{
			Name: ast.NewName(&ast.Name{Value: "bar"}),
			Type: ast.NewNamed(&ast.Named{Name: ast.NewName(&ast.Name{Value: "String"})}),
			Arguments: []*ast.InputValueDefinition{
				ast.NewInputValueDefinition(&ast.InputValueDefinition{
					Name: ast.NewName(&ast.Name{Value: "x"}),
					Type: ast.NewNamed(&ast.Named{Name: ast.NewName(&ast.Name{Value: "Int"})}),
				}),
			},
		}),
	})
	if action3 != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action3)
	}
	_ = val3

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_InputValueDefinition(t *testing.T) {
	fn := printDocASTReducer["InputValueDefinition"]

	// *ast.InputValueDefinition path with directives and description
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewInputValueDefinition(&ast.InputValueDefinition{
			Name: ast.NewName(&ast.Name{Value: "arg1"}),
			Type: ast.NewNamed(&ast.Named{Name: ast.NewName(&ast.Name{Value: "String"})}),
			Directives: []*ast.Directive{
				ast.NewDirective(&ast.Directive{Name: ast.NewName(&ast.Name{Value: "myDir"})}),
			},
			Description: ast.NewStringValue(&ast.StringValue{Value: "an input"}),
		}),
	})
	if action != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action)
	}
	_ = val

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_InterfaceDefinition(t *testing.T) {
	fn := printDocASTReducer["InterfaceDefinition"]

	// *ast.InterfaceDefinition path with directives and description
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewInterfaceDefinition(&ast.InterfaceDefinition{
			Name: ast.NewName(&ast.Name{Value: "Bar"}),
			Directives: []*ast.Directive{
				ast.NewDirective(&ast.Directive{Name: ast.NewName(&ast.Name{Value: "myDir"})}),
			},
			Description: ast.NewStringValue(&ast.StringValue{Value: "an interface"}),
		}),
	})
	if action != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action)
	}
	_ = val

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_UnionDefinition(t *testing.T) {
	fn := printDocASTReducer["UnionDefinition"]

	// *ast.UnionDefinition path with directives and description
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewUnionDefinition(&ast.UnionDefinition{
			Name: ast.NewName(&ast.Name{Value: "SearchResult"}),
			Types: []*ast.Named{
				ast.NewNamed(&ast.Named{Name: ast.NewName(&ast.Name{Value: "Article"})}),
			},
			Directives: []*ast.Directive{
				ast.NewDirective(&ast.Directive{Name: ast.NewName(&ast.Name{Value: "myDir"})}),
			},
			Description: ast.NewStringValue(&ast.StringValue{Value: "a union"}),
		}),
	})
	if action != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action)
	}
	_ = val

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_EnumDefinition(t *testing.T) {
	fn := printDocASTReducer["EnumDefinition"]

	// *ast.EnumDefinition path with directives and description
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewEnumDefinition(&ast.EnumDefinition{
			Name: ast.NewName(&ast.Name{Value: "Color"}),
			Directives: []*ast.Directive{
				ast.NewDirective(&ast.Directive{Name: ast.NewName(&ast.Name{Value: "myDir"})}),
			},
			Description: ast.NewStringValue(&ast.StringValue{Value: "an enum"}),
		}),
	})
	if action != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action)
	}
	_ = val

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_EnumValueDefinition(t *testing.T) {
	fn := printDocASTReducer["EnumValueDefinition"]

	// *ast.EnumValueDefinition path with directives and description
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewEnumValueDefinition(&ast.EnumValueDefinition{
			Name: ast.NewName(&ast.Name{Value: "RED"}),
			Directives: []*ast.Directive{
				ast.NewDirective(&ast.Directive{Name: ast.NewName(&ast.Name{Value: "myDir"})}),
			},
			Description: ast.NewStringValue(&ast.StringValue{Value: "the color red"}),
		}),
	})
	if action != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action)
	}
	_ = val

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_InputObjectDefinition(t *testing.T) {
	fn := printDocASTReducer["InputObjectDefinition"]

	// *ast.InputObjectDefinition path with directives and description
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewInputObjectDefinition(&ast.InputObjectDefinition{
			Name: ast.NewName(&ast.Name{Value: "InputType"}),
			Directives: []*ast.Directive{
				ast.NewDirective(&ast.Directive{Name: ast.NewName(&ast.Name{Value: "myDir"})}),
			},
			Description: ast.NewStringValue(&ast.StringValue{Value: "an input"}),
		}),
	})
	if action != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action)
	}
	_ = val

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_TypeExtensionDefinition(t *testing.T) {
	fn := printDocASTReducer["TypeExtensionDefinition"]

	// *ast.TypeExtensionDefinition path
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewTypeExtensionDefinition(&ast.TypeExtensionDefinition{
			Definition: ast.NewObjectDefinition(&ast.ObjectDefinition{
				Name: ast.NewName(&ast.Name{Value: "Foo"}),
			}),
		}),
	})
	if action != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action)
	}
	_ = val

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

func TestPrintDocASTReducer_DirectiveDefinition(t *testing.T) {
	fn := printDocASTReducer["DirectiveDefinition"]

	// *ast.DirectiveDefinition path with description and args with descriptions (hasArgDesc=true)
	action, val := fn(visitor.VisitFuncParams{
		Node: ast.NewDirectiveDefinition(&ast.DirectiveDefinition{
			Name: ast.NewName(&ast.Name{Value: "myDirective"}),
			Locations: []*ast.Name{
				ast.NewName(&ast.Name{Value: "FIELD"}),
			},
			Description: ast.NewStringValue(&ast.StringValue{Value: "a directive"}),
			Arguments: []*ast.InputValueDefinition{
				ast.NewInputValueDefinition(&ast.InputValueDefinition{
					Name:        ast.NewName(&ast.Name{Value: "arg1"}),
					Type:        ast.NewNamed(&ast.Named{Name: ast.NewName(&ast.Name{Value: "String"})}),
					Description: ast.NewStringValue(&ast.StringValue{Value: "an arg"}),
				}),
			},
		}),
	})
	if action != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action)
	}
	_ = val

	// *ast.DirectiveDefinition path with NO descriptions (hasArgDesc=false, else branch)
	action3, val3 := fn(visitor.VisitFuncParams{
		Node: ast.NewDirectiveDefinition(&ast.DirectiveDefinition{
			Name: ast.NewName(&ast.Name{Value: "simpleDirective"}),
			Locations: []*ast.Name{
				ast.NewName(&ast.Name{Value: "FIELD"}),
			},
			Arguments: []*ast.InputValueDefinition{
				ast.NewInputValueDefinition(&ast.InputValueDefinition{
					Name: ast.NewName(&ast.Name{Value: "x"}),
					Type: ast.NewNamed(&ast.Named{Name: ast.NewName(&ast.Name{Value: "Int"})}),
				}),
			},
		}),
	})
	if action3 != visitor.ActionUpdate {
		t.Fatalf("expected ActionUpdate, got %v", action3)
	}
	_ = val3

	// ActionNoChange
	action2, val2 := fn(visitor.VisitFuncParams{Node: 42})
	if action2 != visitor.ActionNoChange || val2 != nil {
		t.Fatalf("unexpected: action=%v val=%v", action2, val2)
	}
}

// Test helper function call patterns to cover edge cases
func TestJoin_EmptyStrings(t *testing.T) {
	result := join([]string{"a", "", "b", "", "c"}, " ")
	if result != "a b c" {
		t.Fatalf("unexpected: %v", result)
	}
}

func TestWrap_EmptyString(t *testing.T) {
	result := wrap("(", "", ")")
	if result != "" {
		t.Fatalf("expected empty string, got %v", result)
	}

	result2 := wrap("(", "content", ")")
	if result2 != "(content)" {
		t.Fatalf("unexpected: %v", result2)
	}
}

func TestBlock_Empty(t *testing.T) {
	result := block(nil)
	if result != "{}" {
		t.Fatalf("unexpected: %v", result)
	}
}

func TestGetMapValue_Coverage(t *testing.T) {
	// Ensure the !ok branch is tested
	m := map[string]interface{}{}
	result := getMapValue(m, "nonexistent")
	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}

	// Ensure default return (non-map, non-slice, non-nil value)
	m2 := map[string]interface{}{"k": "v"}
	result2 := getMapValue(m2, "k")
	if result2 != "v" {
		t.Fatalf("expected 'v', got %v", result2)
	}

	// slice return
	m3 := map[string]interface{}{"k": []interface{}{"a", "b"}}
	result3 := getMapValue(m3, "k")
	if fmt.Sprintf("%v", result3) != "[a b]" {
		t.Fatalf("unexpected: %v", result3)
	}
}
