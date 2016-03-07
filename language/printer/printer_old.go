package printer

import (
	"fmt"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/visitor"
	//	"log"
)

var printDocASTReducer11 = map[string]visitor.VisitFunc{
	"Name": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.Name:
			return visitor.ActionUpdate, node.Value
		}
		return visitor.ActionNoChange, nil

	},
	"Variable": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.Variable:
			return visitor.ActionUpdate, fmt.Sprintf("$%v", node.Name)
		}
		return visitor.ActionNoChange, nil
	},
	"Document": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.Document:
			definitions := toSliceString(node.Definitions)
			return visitor.ActionUpdate, join(definitions, "\n\n") + "\n"
		}
		return visitor.ActionNoChange, nil
	},
	"OperationDefinition": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.OperationDefinition:
			op := node.Operation
			name := fmt.Sprintf("%v", node.Name)

			defs := wrap("(", join(toSliceString(node.VariableDefinitions), ", "), ")")
			directives := join(toSliceString(node.Directives), " ")
			selectionSet := fmt.Sprintf("%v", node.SelectionSet)
			str := ""
			if name == "" {
				str = selectionSet
			} else {
				str = join([]string{
					op,
					join([]string{name, defs}, ""),
					directives,
					selectionSet,
				}, " ")
			}
			return visitor.ActionUpdate, str
		}
		return visitor.ActionNoChange, nil
	},
	"VariableDefinition": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.VariableDefinition:
			variable := fmt.Sprintf("%v", node.Variable)
			ttype := fmt.Sprintf("%v", node.Type)
			defaultValue := fmt.Sprintf("%v", node.DefaultValue)

			return visitor.ActionUpdate, variable + ": " + ttype + wrap(" = ", defaultValue, "")

		}
		return visitor.ActionNoChange, nil
	},
	"SelectionSet": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.SelectionSet:
			str := block(node.Selections)
			return visitor.ActionUpdate, str

		}
		return visitor.ActionNoChange, nil
	},
	"Field": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.Field:

			alias := fmt.Sprintf("%v", node.Alias)
			name := fmt.Sprintf("%v", node.Name)
			args := toSliceString(node.Arguments)
			directives := toSliceString(node.Directives)
			selectionSet := fmt.Sprintf("%v", node.SelectionSet)

			str := join(
				[]string{
					wrap("", alias, ": ") + name + wrap("(", join(args, ", "), ")"),
					join(directives, " "),
					selectionSet,
				},
				" ",
			)
			return visitor.ActionUpdate, str
		}
		return visitor.ActionNoChange, nil
	},
	"Argument": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.Argument:
			name := fmt.Sprintf("%v", node.Name)
			value := fmt.Sprintf("%v", node.Value)
			return visitor.ActionUpdate, name + ": " + value
		}
		return visitor.ActionNoChange, nil
	},
	"FragmentSpread": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.FragmentSpread:
			name := fmt.Sprintf("%v", node.Name)
			directives := toSliceString(node.Directives)
			return visitor.ActionUpdate, "..." + name + wrap(" ", join(directives, " "), "")
		}
		return visitor.ActionNoChange, nil
	},
	"InlineFragment": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.InlineFragment:
			typeCondition := fmt.Sprintf("%v", node.TypeCondition)
			directives := toSliceString(node.Directives)
			selectionSet := fmt.Sprintf("%v", node.SelectionSet)
			return visitor.ActionUpdate, "... on " + typeCondition + " " + wrap("", join(directives, " "), " ") + selectionSet
		}
		return visitor.ActionNoChange, nil
	},
	"FragmentDefinition": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.FragmentDefinition:
			name := fmt.Sprintf("%v", node.Name)
			typeCondition := fmt.Sprintf("%v", node.TypeCondition)
			directives := toSliceString(node.Directives)
			selectionSet := fmt.Sprintf("%v", node.SelectionSet)
			return visitor.ActionUpdate, "fragment " + name + " on " + typeCondition + " " + wrap("", join(directives, " "), " ") + selectionSet
		}
		return visitor.ActionNoChange, nil
	},

	"IntValue": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.IntValue:
			return visitor.ActionUpdate, fmt.Sprintf("%v", node.Value)
		}
		return visitor.ActionNoChange, nil
	},
	"FloatValue": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.FloatValue:
			return visitor.ActionUpdate, fmt.Sprintf("%v", node.Value)
		}
		return visitor.ActionNoChange, nil
	},
	"StringValue": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.StringValue:
			return visitor.ActionUpdate, `"` + fmt.Sprintf("%v", node.Value) + `"`
		}
		return visitor.ActionNoChange, nil
	},
	"BooleanValue": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.BooleanValue:
			return visitor.ActionUpdate, fmt.Sprintf("%v", node.Value)
		}
		return visitor.ActionNoChange, nil
	},
	"EnumValue": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.EnumValue:
			return visitor.ActionUpdate, fmt.Sprintf("%v", node.Value)
		}
		return visitor.ActionNoChange, nil
	},
	"ListValue": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.ListValue:
			return visitor.ActionUpdate, "[" + join(toSliceString(node.Values), ", ") + "]"
		}
		return visitor.ActionNoChange, nil
	},
	"ObjectValue": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.ObjectValue:
			return visitor.ActionUpdate, "{" + join(toSliceString(node.Fields), ", ") + "}"
		}
		return visitor.ActionNoChange, nil
	},
	"ObjectField": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.ObjectField:
			name := fmt.Sprintf("%v", node.Name)
			value := fmt.Sprintf("%v", node.Value)
			return visitor.ActionUpdate, name + ": " + value
		}
		return visitor.ActionNoChange, nil
	},

	"Directive": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.Directive:
			name := fmt.Sprintf("%v", node.Name)
			args := toSliceString(node.Arguments)
			return visitor.ActionUpdate, "@" + name + wrap("(", join(args, ", "), ")")
		}
		return visitor.ActionNoChange, nil
	},

	"Named": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.Named:
			return visitor.ActionUpdate, fmt.Sprintf("%v", node.Name)
		}
		return visitor.ActionNoChange, nil
	},
	"List": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.List:
			return visitor.ActionUpdate, "[" + fmt.Sprintf("%v", node.Type) + "]"
		}
		return visitor.ActionNoChange, nil
	},
	"NonNull": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.NonNull:
			return visitor.ActionUpdate, fmt.Sprintf("%v", node.Type) + "!"
		}
		return visitor.ActionNoChange, nil
	},

	"ObjectDefinition": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.ObjectDefinition:
			name := fmt.Sprintf("%v", node.Name)
			interfaces := toSliceString(node.Interfaces)
			fields := node.Fields
			str := "type " + name + " " + wrap("implements ", join(interfaces, ", "), " ") + block(fields)
			return visitor.ActionUpdate, str
		}
		return visitor.ActionNoChange, nil
	},
	"FieldDefinition": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.FieldDefinition:
			name := fmt.Sprintf("%v", node.Name)
			ttype := fmt.Sprintf("%v", node.Type)
			args := toSliceString(node.Arguments)
			str := name + wrap("(", join(args, ", "), ")") + ": " + ttype
			return visitor.ActionUpdate, str
		}
		return visitor.ActionNoChange, nil
	},
	"InputValueDefinition": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.InputValueDefinition:
			name := fmt.Sprintf("%v", node.Name)
			ttype := fmt.Sprintf("%v", node.Type)
			defaultValue := fmt.Sprintf("%v", node.DefaultValue)
			str := name + ": " + ttype + wrap(" = ", defaultValue, "")
			return visitor.ActionUpdate, str
		}
		return visitor.ActionNoChange, nil
	},
	"InterfaceDefinition": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.InterfaceDefinition:
			name := fmt.Sprintf("%v", node.Name)
			fields := node.Fields
			str := "interface " + name + " " + block(fields)
			return visitor.ActionUpdate, str
		}
		return visitor.ActionNoChange, nil
	},
	"UnionDefinition": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.UnionDefinition:
			name := fmt.Sprintf("%v", node.Name)
			types := toSliceString(node.Types)
			str := "union " + name + " = " + join(types, " | ")
			return visitor.ActionUpdate, str
		}
		return visitor.ActionNoChange, nil
	},
	"ScalarDefinition": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.ScalarDefinition:
			name := fmt.Sprintf("%v", node.Name)
			str := "scalar " + name
			return visitor.ActionUpdate, str
		}
		return visitor.ActionNoChange, nil
	},
	"EnumDefinition": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.EnumDefinition:
			name := fmt.Sprintf("%v", node.Name)
			values := node.Values
			str := "enum " + name + " " + block(values)
			return visitor.ActionUpdate, str
		}
		return visitor.ActionNoChange, nil
	},
	"EnumValueDefinition": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.EnumValueDefinition:
			name := fmt.Sprintf("%v", node.Name)
			return visitor.ActionUpdate, name
		}
		return visitor.ActionNoChange, nil
	},
	"InputObjectDefinition": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.InputObjectDefinition:
			name := fmt.Sprintf("%v", node.Name)
			fields := node.Fields
			return visitor.ActionUpdate, "input " + name + " " + block(fields)
		}
		return visitor.ActionNoChange, nil
	},
	"TypeExtensionDefinition": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case *ast.TypeExtensionDefinition:
			definition := fmt.Sprintf("%v", node.Definition)
			str := "extend " + definition
			return visitor.ActionUpdate, str
		}
		return visitor.ActionNoChange, nil
	},
}

func Print11(astNode ast.Node) (printed interface{}) {
	//	defer func() interface{} {
	//		if r := recover(); r != nil {
	//			log.Println("Error: %v", r)
	//			return printed
	//		}
	//		return printed
	//	}()
	printed = visitor.Visit(astNode, &visitor.VisitorOptions{
		LeaveKindMap: printDocASTReducer,
	}, nil)
	return printed
}

//
//func PrintMap(astNodeMap map[string]interface{}) (printed interface{}) {
//	defer func() interface{} {
//		if r := recover(); r != nil {
//			return fmt.Sprintf("%v", astNodeMap)
//		}
//		return printed
//	}()
//	printed = visitor.Visit(astNodeMap, &visitor.VisitorOptions{
//		LeaveKindMap: printDocASTReducer,
//	}, nil)
//	return printed
//}
