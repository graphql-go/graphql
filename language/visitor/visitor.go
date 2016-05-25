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
	Parent    ast.Node
	Ancestors []ast.Node
}

type VisitFunc func(p VisitFuncParams) (string, interface{})

type VisitorOptions struct {
	Enter VisitFunc
	Leave VisitFunc
}

type actionBreak struct{}

func visit(root ast.Node, visitorOpts *VisitorOptions, ancestors []ast.Node, parent ast.Node) {
	if root == nil || reflect.ValueOf(root).IsNil() {
		return
	}

	p := VisitFuncParams{
		Node:      root,
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
		visit(root.Name, visitorOpts, p.Ancestors, root)
	case *ast.Document:
		for _, n := range root.Definitions {
			visit(n, visitorOpts, p.Ancestors, root)
		}
	case *ast.OperationDefinition:
		visit(root.Name, visitorOpts, p.Ancestors, root)
		for _, n := range root.VariableDefinitions {
			visit(n, visitorOpts, p.Ancestors, root)
		}
		for _, n := range root.Directives {
			visit(n, visitorOpts, p.Ancestors, root)
		}
		visit(root.SelectionSet, visitorOpts, p.Ancestors, root)
	case *ast.VariableDefinition:
		visit(root.Variable, visitorOpts, p.Ancestors, root)
		visit(root.Type, visitorOpts, p.Ancestors, root)
		visit(root.DefaultValue, visitorOpts, p.Ancestors, root)
	case *ast.SelectionSet:
		for _, n := range root.Selections {
			visit(n, visitorOpts, p.Ancestors, root)
		}
	case *ast.Field:
		visit(root.Alias, visitorOpts, p.Ancestors, root)
		visit(root.Name, visitorOpts, p.Ancestors, root)
		for _, n := range root.Arguments {
			visit(n, visitorOpts, p.Ancestors, root)
		}
		for _, n := range root.Directives {
			visit(n, visitorOpts, p.Ancestors, root)
		}
		visit(root.SelectionSet, visitorOpts, p.Ancestors, root)
	case *ast.Argument:
		visit(root.Name, visitorOpts, p.Ancestors, root)
		visit(root.Value, visitorOpts, p.Ancestors, root)
	case *ast.FragmentSpread:
		visit(root.Name, visitorOpts, p.Ancestors, root)
		for _, n := range root.Directives {
			visit(n, visitorOpts, p.Ancestors, root)
		}
	case *ast.InlineFragment:
		visit(root.TypeCondition, visitorOpts, p.Ancestors, root)
		for _, n := range root.Directives {
			visit(n, visitorOpts, p.Ancestors, root)
		}
		visit(root.SelectionSet, visitorOpts, p.Ancestors, root)
	case *ast.FragmentDefinition:
		visit(root.Name, visitorOpts, p.Ancestors, root)
		visit(root.TypeCondition, visitorOpts, p.Ancestors, root)
		for _, n := range root.Directives {
			visit(n, visitorOpts, p.Ancestors, root)
		}
		visit(root.SelectionSet, visitorOpts, p.Ancestors, root)
	case *ast.IntValue:
	case *ast.FloatValue:
	case *ast.StringValue:
	case *ast.BooleanValue:
	case *ast.EnumValue:
	case *ast.ListValue:
		for _, n := range root.Values {
			visit(n, visitorOpts, p.Ancestors, root)
		}
	case *ast.ObjectValue:
		for _, n := range root.Fields {
			visit(n, visitorOpts, p.Ancestors, root)
		}
	case *ast.ObjectField:
		visit(root.Name, visitorOpts, p.Ancestors, root)
		visit(root.Value, visitorOpts, p.Ancestors, root)
	case *ast.Directive:
		visit(root.Name, visitorOpts, p.Ancestors, root)
		for _, n := range root.Arguments {
			visit(n, visitorOpts, p.Ancestors, root)
		}
	case *ast.Named:
		visit(root.Name, visitorOpts, p.Ancestors, root)
	case *ast.List:
		visit(root.Type, visitorOpts, p.Ancestors, root)
	case *ast.NonNull:
		visit(root.Type, visitorOpts, p.Ancestors, root)
	case *ast.ObjectDefinition:
		visit(root.Name, visitorOpts, p.Ancestors, root)
		for _, n := range root.Interfaces {
			visit(n, visitorOpts, p.Ancestors, root)
		}
		for _, n := range root.Fields {
			visit(n, visitorOpts, p.Ancestors, root)
		}
	case *ast.FieldDefinition:
		visit(root.Name, visitorOpts, p.Ancestors, root)
		for _, n := range root.Arguments {
			visit(n, visitorOpts, p.Ancestors, root)
		}
		visit(root.Type, visitorOpts, p.Ancestors, root)
	case *ast.InputValueDefinition:
		visit(root.Name, visitorOpts, p.Ancestors, root)
		visit(root.Type, visitorOpts, p.Ancestors, root)
		visit(root.DefaultValue, visitorOpts, p.Ancestors, root)
	case *ast.InterfaceDefinition:
		visit(root.Name, visitorOpts, p.Ancestors, root)
		for _, n := range root.Fields {
			visit(n, visitorOpts, p.Ancestors, root)
		}
	case *ast.UnionDefinition:
		visit(root.Name, visitorOpts, p.Ancestors, root)
		for _, n := range root.Types {
			visit(n, visitorOpts, p.Ancestors, root)
		}
	case *ast.ScalarDefinition:
		visit(root.Name, visitorOpts, p.Ancestors, root)
	case *ast.EnumDefinition:
		visit(root.Name, visitorOpts, p.Ancestors, root)
		for _, n := range root.Values {
			visit(n, visitorOpts, p.Ancestors, root)
		}
	case *ast.EnumValueDefinition:
		visit(root.Name, visitorOpts, p.Ancestors, root)
	case *ast.InputObjectDefinition:
		visit(root.Name, visitorOpts, p.Ancestors, root)
		for _, n := range root.Fields {
			visit(n, visitorOpts, p.Ancestors, root)
		}
	case *ast.TypeExtensionDefinition:
		visit(root.Definition, visitorOpts, p.Ancestors, root)
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
	visit(root, visitorOpts, make([]ast.Node, 0, 64), nil)
	return nil
}
