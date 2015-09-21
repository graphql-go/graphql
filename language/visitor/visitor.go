package visitor
import (
	"github.com/chris-ramon/graphql-go/language/ast"
	"reflect"
	"fmt"
)

const (
	ActionNoChange = "NOCHANGE"
	ActionBreak    = "BREAK"
	ActionRemove   = "REMOVE"
	ActionUpdate   = ""
)
type KeyMap map[string][]string

// note that the keys are in Capital letters, equivalent to the ast.Node field Names
var QueryDocumentKeys KeyMap = KeyMap{
	"Name": []string{},
	"Document": []string{"Definitions"},
	"OperationDefinition": []string{
		"Name",
		"VariableDefinitions",
		"Directives",
		"SelectionSet",
	},
	"VariableDefinition": []string{
		"Variable",
		"Type",
		"DefaultValue",
	},
	"Variable": []string{"Name"},
	"SelectionSet": []string{"Selections"},
	"Field": []string{
		"Alias",
		"Name",
		"Arguments",
		"Directives",
		"SelectionSet",
	},
	"Argument": []string{
		"Name",
		"Value",
	},

	"FragmentSpread": []string{
		"Name",
		"Directives",
	},
	"InlineFragment": []string{
		"TypeCondition",
		"Directives",
		"SelectionSet",
	},
	"FragmentDefinition": []string{
		"Name",
		"TypeCondition",
		"Directives",
		"SelectionSet",
	},

	"IntValue": []string{},
	"FloatValue": []string{},
	"StringValue": []string{},
	"BooleanValue": []string{},
	"EnumValue": []string{},
	"ListValue": []string{"Values" },
	"ObjectValue": []string{"Fields" },
	"ObjectField": []string{
		"Name",
		"Value",
	},

	"Directive": []string{
		"Name",
		"Arguments",
	},

	"NamedType": []string{"Name" },
	"ListType": []string{"Type" },
	"NonNullType": []string{"Type" },

	"ObjectTypeDefinition": []string{
		"Name",
		"Interfaces",
		"Fields",
	},
	"FieldDefinition": []string{
		"Name",
		"Arguments",
		"Type",
	},
	"InputValueDefinition": []string{
		"Name",
		"Type",
		"DefaultValue",
	},
	"InterfaceTypeDefinition": []string{
		"Name",
		"Fields",
	},
	"UnionTypeDefinition": []string{
		"Name",
		"Types",
	},
	"ScalarTypeDefinition": []string{"Name" },
	"EnumTypeDefinition": []string{
		"Name",
		"Values",
	},
	"EnumValueDefinition": []string{"Name" },
	"InputObjectTypeDefinition": []string{
		"Name",
		"Fields",
	},
	"TypeExtensionDefinition": []string{"Definition" },

}

type stack struct {
	Index   int
	Keys    []interface{}
	Edits   []*edit
	InArray bool
	Prev    *stack

}

type VisitFuncParams struct {
	Node      interface{}
	Key       interface{}
	Parent    interface{}
	Path      []interface{}
	Ancestors []interface{}
}
type VisitFuncResults struct {
	Break      bool        // set to true to stop traversal, default false
	Skip       bool        // set to true to skip over sub-tree, default false
	Remove     bool        // set to true to remove node, default false
	Edit       bool        // set to true to edit node, default false
	EditedNode interface{} // default nil
}
type VisitFunc func(p VisitFuncParams) (string, interface{})

type NamedVisitFuncs struct {
	Kind  VisitFunc // 1) Named visitors triggered when entering a node a specific kind.
	Leave VisitFunc // 2) Named visitors that trigger upon entering and leaving a node of
	Enter VisitFunc // 2) Named visitors that trigger upon entering and leaving a node of
}

type VisitorOptions struct {
	KindFuncMap  map[string]NamedVisitFuncs
	Enter        VisitFunc            // 3) Generic visitors that trigger upon entering and leaving any node.
	Leave        VisitFunc            // 3) Generic visitors that trigger upon entering and leaving any node.

	EnterKindMap map[string]VisitFunc // 4) Parallel visitors for entering and leaving nodes of a specific kind
	LeaveKindMap map[string]VisitFunc // 4) Parallel visitors for entering and leaving nodes of a specific kind
}

func isSlice(Value interface{}) bool {
	val := reflect.ValueOf(Value)
	if val.IsValid() && val.Type().Kind() == reflect.Slice {
		return true
	}
	return false
}

func pop(a []interface{}) (x interface{}, aa []interface{}) {
	if len(a) == 0 {
		return x, aa
	}
	x, aa = a[len(a) - 1], a[:len(a) - 1]
	return x, aa
}
func copy(a []interface{}) ([]interface{}) {
	return append([]interface{}(nil), a...)
}
func spliceSelections(a []ast.Selection, i int) ([]ast.Selection) {
	if i >= len(a) {
		return a
	}
	if i < 0 {
		return []ast.Selection{}
	}
	return append(a[:i], a[i + 1:]...)
}
func spliceNodes(a []ast.Node, i int) ([]ast.Node) {
	if i >= len(a) {
		return a
	}
	if i < 0 {
		return []ast.Node{}
	}
	return append(a[:i], a[i + 1:]...)
}
type edit struct {
	Key          interface{}
	Value        interface{}
	Change       VisitFuncResults
	UpdateParent bool
	ChildNode    interface{}
}
func Visit(root ast.Node, visitorOpts *VisitorOptions, keyMap KeyMap) ast.Node {
	visitorKeys := keyMap
	if visitorKeys == nil {
		visitorKeys = QueryDocumentKeys
	}

	var sstack *stack
	var parent interface{}
	inArray := isSlice(root)
	keys := []interface{}{root }
	index := -1
	edits := []*edit{}
	path := []interface{}{}
	ancestors := []interface{}{}
	newRoot := root
	Loop:
	for {
		index = index + 1

		isLeaving := (len(keys) == index)
		var key interface{}
		var node interface{}
		isEdited := (isLeaving && len(edits) != 0 )

		if isLeaving {
			if len(ancestors) == 0 {
				key = nil
			} else {
				key, path = pop(path)
			}

			node = parent
			parent, ancestors = pop(ancestors)
			if isEdited {
				editOffset := 0
				for _, edit := range edits {
					arrayEditKey := 0
					if inArray {
						keyInt := edit.Key.(int)
						edit.Key = keyInt - editOffset
						arrayEditKey = edit.Key.(int)
					}
					if inArray && isNilNode(edit.Value) {
						if n, ok := node.([]ast.Selection); ok {
							node = spliceSelections(n, arrayEditKey)
						} else if n, ok := node.([]ast.Node); ok {
							node = spliceNodes(n, arrayEditKey)
						} else {
							panic(fmt.Sprintf("1 Invalid AST Node: %v", node))
						}
						editOffset = editOffset + 1
					} else {
						if inArray {
							if n, ok := node.([]ast.Selection); ok {
								n[arrayEditKey] = edit.Value.(ast.Selection)
								node = n
							} else if n, ok := node.([]ast.Node); ok {
								n[arrayEditKey] = edit.Value.(ast.Node)
								node = n
							} else {
								panic(fmt.Sprintf("2 Invalid AST Node: %v", node))
							}
						} else {
							setField(node, edit.Key, edit.Value)
						}
					}
				}
			}
			index = sstack.Index
			keys = sstack.Keys
			edits = sstack.Edits
			inArray = sstack.InArray
			sstack = sstack.Prev
		} else {
			// get key
			if !isNilNode(parent) {
				if inArray {
					key = index
				} else {
					key = getField(keys, index)
				}
			} else {
				// initial conditions
				key = nil
			}
			// get node
			if !isNilNode(parent) {
				node = getField(parent, key)
			} else {
				// initial conditions
				node = newRoot
			}

			if isNilNode(node) {
				continue
			}
			if !isNilNode(parent) {
				path = append(path, key)
			}
		}

		// get result from visitFn for a node if set
		var result interface{}
		resultIsUndefined := true
		if !isSlice(node) && !isNilNode(node) {
			astNode, ok := node.(ast.Node)
			if !ok {
				panic(fmt.Sprintf("3 Invalid AST Node: %v", node))
			}
			visitFn := getVisitFn(visitorOpts, isLeaving, astNode.GetKind())
			if visitFn != nil {
				p := VisitFuncParams{
					Node: node,
					Key: key,
					Parent: parent,
					Path: path,
					Ancestors: ancestors,
				}
				action := ActionUpdate
				action, result = visitFn(p)
				if action == ActionBreak {
					break Loop
				}
				if action == ActionRemove {
					if isLeaving {
						_, path = pop(path)
						continue
					}
				}
				if action != ActionNoChange {
					resultIsUndefined = false
					edits = append(edits, &edit{
						Key: key,
						Value: result,
					})
					if !isLeaving {
						if n, ok := result.(ast.Selection); ok {
							node = n
						} else if n, ok := result.(ast.Node); ok {
							node = n
						} else {
							_, path = pop(path)
							continue
						}
					}
				} else {
					resultIsUndefined = true

				}
			}

		}
		if resultIsUndefined && isEdited {
			edits = append(edits, &edit{
				Key: key,
				Value: node,
			})
		}

		if !isLeaving {

			// add to stack
			prevStack := sstack
			sstack = &stack{
				InArray: inArray,
				Index: index,
				Keys: keys,
				Edits: edits,
				Prev: prevStack,
			}

			// replace keys
			inArray = isSlice(node)
			keys = []interface{}{}
			if !isNilNode(node) {
				if inArray {
					// get keys
					if n, ok := node.([]ast.Node); ok {
						for _, m := range n {
							keys = append(keys, m)
						}
					} else if n, ok := node.([]ast.Selection); ok {
						for _, m := range n {
							keys = append(keys, m)
						}
					} else if n, ok := node.([]*ast.VariableDefinition); ok {
						for _, m := range n {
							keys = append(keys, m)
						}
					} else if n, ok := node.([]*ast.Argument); ok {
						for _, m := range n {
							keys = append(keys, m)
						}
					} else if n, ok := node.([]*ast.Directive); ok {
						for _, m := range n {
							keys = append(keys, m)
						}
					} else if n, ok := node.([]*ast.ObjectField); ok {
						for _, m := range n {
							keys = append(keys, m)
						}
					} else if n, ok := node.([]ast.Value); ok {
						for _, m := range n {
							keys = append(keys, m)
						}
					} else {
						panic(fmt.Sprintf("4 Invalid AST Node: %v", node))
					}

				} else {
					if n, ok := node.(ast.Node); ok {
						if n, ok := visitorKeys[n.GetKind()]; ok {
							for _, m := range n {
								keys = append(keys, m)
							}
						}
					} else {
						panic(fmt.Sprintf("5 Invalid AST Node: %v", node))
					}
				}
			}

			index = -1
			edits = []*edit{}
			if !isNilNode(parent) {
				ancestors = append(ancestors, parent)
			}
			parent = node
		}

		// loop guard
		if sstack == nil {
			break Loop
		}
	}
	if len(edits) != 0 {
		newRoot = edits[0].Value.(ast.Node)
	}
	return newRoot
}

func getField(obj interface{}, key interface{}) interface{} {
	val := reflect.ValueOf(obj)
	if val.Type().Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Type().Kind() == reflect.Struct {
		key, ok := key.(string)
		if !ok {
			return nil
		}
		valField := val.FieldByName(key)
		if valField.IsValid() {
			return valField.Interface()
		}
		return nil
	}
	if val.Type().Kind() == reflect.Slice {
		key, ok := key.(int)
		if !ok {
			return nil
		}
		if key >= val.Len() {
			return nil
		}
		valField := val.Index(key)
		if valField.IsValid() {
			return valField.Interface()
		}
		return nil
	}
	if val.Type().Kind() == reflect.Map {
		keyVal := reflect.ValueOf(key)
		valField := val.MapIndex(keyVal)
		if valField.IsValid() {
			return valField.Interface()
		}
		return nil
	}
	return nil
}

func setField(obj interface{}, key interface{}, value interface{}) {
	val := reflect.ValueOf(obj)
	if !val.IsValid() {
		return
	}
	if val.Type().Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Type().Kind() == reflect.Struct {
		keyStr, ok := key.(string)
		if !ok {
			return
		}
		valField := val.FieldByName(keyStr)
		if valField.CanSet() {

			valueVal := reflect.ValueOf(value)
			if valueVal.IsValid() {
				valField.Set(valueVal)
			} else {
				// set to zero
				valueVal = reflect.New(val.Type().Elem())
				valField.Set(valueVal)
			}

		}
		return
	}
	return
}

func isNilNode(node interface{}) bool {
	val := reflect.ValueOf(node)
	if !val.IsValid() {
		return true
	}
	if val.Type().Kind() == reflect.Ptr {
		return val.IsNil()
	}
	if val.Type().Kind() == reflect.Slice {
		return val.Len() == 0
	}
	if val.Type().Kind() == reflect.Map {
		return val.Len() == 0
	}
	if val.Type().Kind() == reflect.Bool {
		return val.Interface().(bool)
	}
	return true
}
func copyNode(node interface{}) ast.Node {
	val := reflect.ValueOf(node)
	if !val.IsValid() {
		return nil
	}
	switch node := node.(type) {
	case *ast.Document:
		n := *node
		return &n
	default:
		return node.(ast.Node)
	}
}

func getVisitFn(visitorOpts *VisitorOptions, isLeaving bool, kind string) VisitFunc {
	kindVisitor, ok := visitorOpts.KindFuncMap[kind]
	if ok {
		if !isLeaving && kindVisitor.Kind != nil {
			// { Kind() {} }
			return kindVisitor.Kind
		}
		if isLeaving {
			// { Kind: { leave() {} } }
			return kindVisitor.Leave
		} else {
			// { Kind: { enter() {} } }
			return kindVisitor.Enter
		}
		return nil
	}

	if isLeaving {
		// { enter() {} }
		specificVisitor := visitorOpts.Leave
		if specificVisitor != nil {
			return specificVisitor
		}
		if specificKindVisitor, ok := visitorOpts.LeaveKindMap[kind]; ok {
			// { leave: { Kind() {} } }
			return specificKindVisitor
		}
		return nil

	} else {
		// { leave() {} }
		specificVisitor := visitorOpts.Enter
		if specificVisitor != nil {
			return specificVisitor
		}
		if specificKindVisitor, ok := visitorOpts.EnterKindMap[kind]; ok {
			// { enter: { Kind() {} } }
			return specificKindVisitor
		}
		return nil
	}
	return nil
}
