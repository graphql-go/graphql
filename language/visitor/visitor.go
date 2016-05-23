package visitor

import (
	"fmt"
	"github.com/sprucehealth/graphql/language/ast"
	"reflect"
)

const (
	ActionNoChange = ""
	ActionBreak    = "BREAK"
	ActionSkip     = "SKIP"
)

type VisitFuncParams struct {
	Node      interface{}
	Key       interface{} // The name of this node's field in its parent node
	Parent    ast.Node
	Ancestors []ast.Node
}

type VisitFunc func(p VisitFuncParams) (string, interface{})

type VisitorOptions struct {
	Enter VisitFunc
	Leave VisitFunc
}

type actionBreak struct{}

func visit(root ast.Node, visitorOpts *VisitorOptions, ancestors []ast.Node, parent ast.Node, key interface{}) {
	if root == nil || reflect.ValueOf(root).IsNil() {
		return
	}

	p := VisitFuncParams{
		Node:      root,
		Key:       key,
		Parent:    parent,
		Ancestors: ancestors,
	}
	if parent != nil {
		p.Ancestors = append(p.Ancestors, parent)
	}

	if visitorOpts.Enter != nil {
		// TODO: ignoring result (i.e. error) for now
		action, _ := visitorOpts.Enter(p)
		switch action {
		case ActionSkip:
			return
		case ActionBreak:
			panic(actionBreak{})
		}
	}

	switch root := root.(type) {
	case *ast.Name:
	case *ast.Variable:
		visit(root.Name, visitorOpts, p.Ancestors, root, "Name")
	case *ast.Document:
		for i, n := range root.Definitions {
			visit(n, visitorOpts, p.Ancestors, root, i)
		}
	case *ast.OperationDefinition:
		visit(root.Name, visitorOpts, p.Ancestors, root, "Name")
		for i, n := range root.VariableDefinitions {
			visit(n, visitorOpts, p.Ancestors, root, i)
		}
		for i, n := range root.Directives {
			visit(n, visitorOpts, p.Ancestors, root, i)
		}
		visit(root.SelectionSet, visitorOpts, p.Ancestors, root, "SelectionSet")
	case *ast.VariableDefinition:
		visit(root.Variable, visitorOpts, p.Ancestors, root, "Variable")
		visit(root.Type, visitorOpts, p.Ancestors, root, "Type")
		visit(root.DefaultValue, visitorOpts, p.Ancestors, root, "DefaultValue")
	case *ast.SelectionSet:
		for i, n := range root.Selections {
			visit(n, visitorOpts, p.Ancestors, root, i)
		}
	case *ast.Field:
		visit(root.Alias, visitorOpts, p.Ancestors, root, "Alias")
		visit(root.Name, visitorOpts, p.Ancestors, root, "Name")
		for i, n := range root.Arguments {
			visit(n, visitorOpts, p.Ancestors, root, i)
		}
		for i, n := range root.Directives {
			visit(n, visitorOpts, p.Ancestors, root, i)
		}
		visit(root.SelectionSet, visitorOpts, p.Ancestors, root, "SelectionSet")
	case *ast.Argument:
		visit(root.Name, visitorOpts, p.Ancestors, root, "Name")
		visit(root.Value, visitorOpts, p.Ancestors, root, "Value")
	case *ast.FragmentSpread:
		visit(root.Name, visitorOpts, p.Ancestors, root, "Name")
		for i, n := range root.Directives {
			visit(n, visitorOpts, p.Ancestors, root, i)
		}
	case *ast.InlineFragment:
		visit(root.TypeCondition, visitorOpts, p.Ancestors, root, "TypeCondition")
		for i, n := range root.Directives {
			visit(n, visitorOpts, p.Ancestors, root, i)
		}
		visit(root.SelectionSet, visitorOpts, p.Ancestors, root, "SelectionSet")
	case *ast.FragmentDefinition:
		visit(root.Name, visitorOpts, p.Ancestors, root, "Name")
		visit(root.TypeCondition, visitorOpts, p.Ancestors, root, "TypeCondition")
		for i, n := range root.Directives {
			visit(n, visitorOpts, p.Ancestors, root, i)
		}
		visit(root.SelectionSet, visitorOpts, p.Ancestors, root, "SelectionSet")
	case *ast.IntValue:
	case *ast.FloatValue:
	case *ast.StringValue:
	case *ast.BooleanValue:
	case *ast.EnumValue:
	case *ast.ListValue:
		for i, n := range root.Values {
			visit(n, visitorOpts, p.Ancestors, root, i)
		}
	case *ast.ObjectValue:
		for i, n := range root.Fields {
			visit(n, visitorOpts, p.Ancestors, root, i)
		}
	case *ast.ObjectField:
		visit(root.Name, visitorOpts, p.Ancestors, root, "Name")
		visit(root.Value, visitorOpts, p.Ancestors, root, "Value")
	case *ast.Directive:
		visit(root.Name, visitorOpts, p.Ancestors, root, "Name")
		for i, n := range root.Arguments {
			visit(n, visitorOpts, p.Ancestors, root, i)
		}
	case *ast.Named:
		visit(root.Name, visitorOpts, p.Ancestors, root, "Name")
	case *ast.List:
		visit(root.Type, visitorOpts, p.Ancestors, root, "Type")
	case *ast.NonNull:
		visit(root.Type, visitorOpts, p.Ancestors, root, "Type")
	case *ast.ObjectDefinition:
		visit(root.Name, visitorOpts, p.Ancestors, root, "Name")
		for i, n := range root.Interfaces {
			visit(n, visitorOpts, p.Ancestors, root, i)
		}
		for i, n := range root.Fields {
			visit(n, visitorOpts, p.Ancestors, root, i)
		}
	case *ast.FieldDefinition:
		visit(root.Name, visitorOpts, p.Ancestors, root, "Name")
		for i, n := range root.Arguments {
			visit(n, visitorOpts, p.Ancestors, root, i)
		}
		visit(root.Type, visitorOpts, p.Ancestors, root, "Type")
	case *ast.InputValueDefinition:
		visit(root.Name, visitorOpts, p.Ancestors, root, "Name")
		visit(root.Type, visitorOpts, p.Ancestors, root, "Type")
		visit(root.DefaultValue, visitorOpts, p.Ancestors, root, "DefaultValue")
	case *ast.InterfaceDefinition:
		visit(root.Name, visitorOpts, p.Ancestors, root, "Name")
		for i, n := range root.Fields {
			visit(n, visitorOpts, p.Ancestors, root, i)
		}
	case *ast.UnionDefinition:
		visit(root.Name, visitorOpts, p.Ancestors, root, "Name")
		for i, n := range root.Types {
			visit(n, visitorOpts, p.Ancestors, root, i)
		}
	case *ast.ScalarDefinition:
		visit(root.Name, visitorOpts, p.Ancestors, root, "Name")
	case *ast.EnumDefinition:
		visit(root.Name, visitorOpts, p.Ancestors, root, "Name")
		for i, n := range root.Values {
			visit(n, visitorOpts, p.Ancestors, root, i)
		}
	case *ast.EnumValueDefinition:
		visit(root.Name, visitorOpts, p.Ancestors, root, "Name")
	case *ast.InputObjectDefinition:
		visit(root.Name, visitorOpts, p.Ancestors, root, "Name")
		for i, n := range root.Fields {
			visit(n, visitorOpts, p.Ancestors, root, i)
		}
	case *ast.TypeExtensionDefinition:
		visit(root.Definition, visitorOpts, p.Ancestors, root, "Definition")
	default:
		panic("unknown node type")
	}

	if visitorOpts.Leave != nil {
		// TODO: ignoring result (i.e. error) for now
		action, _ := visitorOpts.Leave(p)
		switch action {
		case ActionBreak:
			panic(actionBreak{})
		}
	}
}

func Visit(root ast.Node, visitorOpts *VisitorOptions) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(actionBreak); ok {
				err = nil
			} else if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("runtime error: %v", r)
			}
		}
	}()
	visit(root, visitorOpts, make([]ast.Node, 0, 64), nil, nil)
	return nil
}
