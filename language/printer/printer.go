package printer

import (
	"fmt"
	"strings"

	"reflect"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/visitor"
)

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
		default:
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
		if v == nil {
			return ""
		}
		switch v := v.(type) {
		case map[string]interface{}:
			valMap = v
			continue
		case string:
			return v
		default:
			return fmt.Sprintf("%v", v)
		}
	}
	return ""
}

func toSliceString(slice interface{}) []string {
	if slice == nil {
		return []string{}
	}
	res := []string{}
	switch reflect.TypeOf(slice).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(slice)
		for i := 0; i < s.Len(); i++ {
			elem := s.Index(i)
			elemInterface := elem.Interface()
			if elem, ok := elemInterface.(string); ok {
				res = append(res, elem)
			}
		}
		return res
	default:
		return res
	}
}

func join(str []string, sep string) string {
	ss := []string{}
	// filter out empty strings
	for _, s := range str {
		if s == "" {
			continue
		}
		ss = append(ss, s)
	}
	return strings.Join(ss, sep)
}

func wrap(start, maybeString, end string) string {
	if maybeString == "" {
		return maybeString
	}
	return start + maybeString + end
}
func block(maybeArray interface{}) string {
	if maybeArray == nil {
		return ""
	}
	s := toSliceString(maybeArray)
	return indent("{\n"+join(s, "\n")) + "\n}"
}

func indent(maybeString interface{}) string {
	if maybeString == nil {
		return ""
	}
	switch str := maybeString.(type) {
	case string:
		return strings.Replace(str, "\n", "\n  ", -1)
	}
	return ""
}

var printDocASTReducer = map[string]visitor.VisitFunc{
	"Name": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.Name:
			return visitor.ActionUpdate, node.Value, nil
		case map[string]interface{}:
			return visitor.ActionUpdate, getMapValue(node, "Value"), nil
		}
		return visitor.ActionNoChange, nil, nil
	},
	"Variable": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.Variable:
			return visitor.ActionUpdate, fmt.Sprintf("$%v", node.Name), nil
		case map[string]interface{}:
			return visitor.ActionUpdate, "$" + getMapValueString(node, "Name"), nil
		}
		return visitor.ActionNoChange, nil, nil
	},
	"Document": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.Document:
			definitions := toSliceString(node.Definitions)
			return visitor.ActionUpdate, join(definitions, "\n\n") + "\n", nil
		case map[string]interface{}:
			definitions := toSliceString(getMapValue(node, "Definitions"))
			return visitor.ActionUpdate, join(definitions, "\n\n") + "\n", nil
		}
		return visitor.ActionNoChange, nil, nil
	},
	"OperationDefinition": func(p visitor.VisitFuncParams) (string, interface{}, error) {
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
			return visitor.ActionUpdate, str, nil
		case map[string]interface{}:
			op := getMapValueString(node, "Operation")
			name := getMapValueString(node, "Name")

			defs := wrap("(", join(toSliceString(getMapValue(node, "VariableDefinitions")), ", "), ")")
			directives := join(toSliceString(getMapValue(node, "Directives")), " ")
			selectionSet := getMapValueString(node, "SelectionSet")
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
			return visitor.ActionUpdate, str, nil
		}
		return visitor.ActionNoChange, nil, nil
	},
	"VariableDefinition": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.VariableDefinition:
			variable := fmt.Sprintf("%v", node.Variable)
			ttype := fmt.Sprintf("%v", node.Type)
			defaultValue := fmt.Sprintf("%v", node.DefaultValue)

			return visitor.ActionUpdate, variable + ": " + ttype + wrap(" = ", defaultValue, ""), nil
		case map[string]interface{}:

			variable := getMapValueString(node, "Variable")
			ttype := getMapValueString(node, "Type")
			defaultValue := getMapValueString(node, "DefaultValue")

			return visitor.ActionUpdate, variable + ": " + ttype + wrap(" = ", defaultValue, ""), nil

		}
		return visitor.ActionNoChange, nil, nil
	},
	"SelectionSet": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.SelectionSet:
			str := block(node.Selections)
			return visitor.ActionUpdate, str, nil
		case map[string]interface{}:
			selections := getMapValue(node, "Selections")
			str := block(selections)
			return visitor.ActionUpdate, str, nil

		}
		return visitor.ActionNoChange, nil, nil
	},
	"Field": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.Argument:
			name := fmt.Sprintf("%v", node.Name)
			value := fmt.Sprintf("%v", node.Value)
			return visitor.ActionUpdate, name + ": " + value, nil
		case map[string]interface{}:

			alias := getMapValueString(node, "Alias")
			name := getMapValueString(node, "Name")
			args := toSliceString(getMapValue(node, "Arguments"))
			directives := toSliceString(getMapValue(node, "Directives"))
			selectionSet := getMapValueString(node, "SelectionSet")

			str := join(
				[]string{
					wrap("", alias, ": ") + name + wrap("(", join(args, ", "), ")"),
					join(directives, " "),
					selectionSet,
				},
				" ",
			)
			return visitor.ActionUpdate, str, nil
		}
		return visitor.ActionNoChange, nil, nil
	},
	"Argument": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.FragmentSpread:
			name := fmt.Sprintf("%v", node.Name)
			directives := toSliceString(node.Directives)
			return visitor.ActionUpdate, "..." + name + wrap(" ", join(directives, " "), ""), nil
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			value := getMapValueString(node, "Value")
			return visitor.ActionUpdate, name + ": " + value, nil
		}
		return visitor.ActionNoChange, nil, nil
	},
	"FragmentSpread": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.InlineFragment:
			typeCondition := fmt.Sprintf("%v", node.TypeCondition)
			directives := toSliceString(node.Directives)
			selectionSet := fmt.Sprintf("%v", node.SelectionSet)
			return visitor.ActionUpdate, "... on " + typeCondition + " " + wrap("", join(directives, " "), " ") + selectionSet, nil
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			directives := toSliceString(getMapValue(node, "Directives"))
			return visitor.ActionUpdate, "..." + name + wrap(" ", join(directives, " "), ""), nil
		}
		return visitor.ActionNoChange, nil, nil
	},
	"InlineFragment": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			typeCondition := getMapValueString(node, "TypeCondition")
			directives := toSliceString(getMapValue(node, "Directives"))
			selectionSet := getMapValueString(node, "SelectionSet")
			return visitor.ActionUpdate, "... on " + typeCondition + " " + wrap("", join(directives, " "), " ") + selectionSet, nil
		}
		return visitor.ActionNoChange, nil, nil
	},
	"FragmentDefinition": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.FragmentDefinition:
			name := fmt.Sprintf("%v", node.Name)
			typeCondition := fmt.Sprintf("%v", node.TypeCondition)
			directives := toSliceString(node.Directives)
			selectionSet := fmt.Sprintf("%v", node.SelectionSet)
			return visitor.ActionUpdate, "fragment " + name + " on " + typeCondition + " " + wrap("", join(directives, " "), " ") + selectionSet, nil
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			typeCondition := getMapValueString(node, "TypeCondition")
			directives := toSliceString(getMapValue(node, "Directives"))
			selectionSet := getMapValueString(node, "SelectionSet")
			return visitor.ActionUpdate, "fragment " + name + " on " + typeCondition + " " + wrap("", join(directives, " "), " ") + selectionSet, nil
		}
		return visitor.ActionNoChange, nil, nil
	},

	"IntValue": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.IntValue:
			return visitor.ActionUpdate, fmt.Sprintf("%v", node.Value), nil
		case map[string]interface{}:
			return visitor.ActionUpdate, getMapValueString(node, "Value"), nil
		}
		return visitor.ActionNoChange, nil, nil
	},
	"FloatValue": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.FloatValue:
			return visitor.ActionUpdate, fmt.Sprintf("%v", node.Value), nil
		case map[string]interface{}:
			return visitor.ActionUpdate, getMapValueString(node, "Value"), nil
		}
		return visitor.ActionNoChange, nil, nil
	},
	"StringValue": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.StringValue:
			return visitor.ActionUpdate, `"` + fmt.Sprintf("%v", node.Value) + `"`, nil
		case map[string]interface{}:
			return visitor.ActionUpdate, `"` + getMapValueString(node, "Value") + `"`, nil
		}
		return visitor.ActionNoChange, nil, nil
	},
	"BooleanValue": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.BooleanValue:
			return visitor.ActionUpdate, fmt.Sprintf("%v", node.Value), nil
		case map[string]interface{}:
			return visitor.ActionUpdate, getMapValueString(node, "Value"), nil
		}
		return visitor.ActionNoChange, nil, nil
	},
	"EnumValue": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.EnumValue:
			return visitor.ActionUpdate, fmt.Sprintf("%v", node.Value), nil
		case map[string]interface{}:
			return visitor.ActionUpdate, getMapValueString(node, "Value"), nil
		}
		return visitor.ActionNoChange, nil, nil
	},
	"ListValue": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.ListValue:
			return visitor.ActionUpdate, "[" + join(toSliceString(node.Values), ", ") + "]", nil
		case map[string]interface{}:
			return visitor.ActionUpdate, "[" + join(toSliceString(getMapValue(node, "Values")), ", ") + "]", nil
		}
		return visitor.ActionNoChange, nil, nil
	},
	"ObjectValue": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.ObjectValue:
			return visitor.ActionUpdate, "{" + join(toSliceString(node.Fields), ", ") + "}", nil
		case map[string]interface{}:
			return visitor.ActionUpdate, "{" + join(toSliceString(getMapValue(node, "Fields")), ", ") + "}", nil
		}
		return visitor.ActionNoChange, nil, nil
	},
	"ObjectField": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.ObjectField:
			name := fmt.Sprintf("%v", node.Name)
			value := fmt.Sprintf("%v", node.Value)
			return visitor.ActionUpdate, name + ": " + value, nil
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			value := getMapValueString(node, "Value")
			return visitor.ActionUpdate, name + ": " + value, nil
		}
		return visitor.ActionNoChange, nil, nil
	},

	"Directive": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.Directive:
			name := fmt.Sprintf("%v", node.Name)
			args := toSliceString(node.Arguments)
			return visitor.ActionUpdate, "@" + name + wrap("(", join(args, ", "), ")"), nil
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			args := toSliceString(getMapValue(node, "Arguments"))
			return visitor.ActionUpdate, "@" + name + wrap("(", join(args, ", "), ")"), nil
		}
		return visitor.ActionNoChange, nil, nil
	},

	"Named": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.Named:
			return visitor.ActionUpdate, fmt.Sprintf("%v", node.Name), nil
		case map[string]interface{}:
			return visitor.ActionUpdate, getMapValueString(node, "Name"), nil
		}
		return visitor.ActionNoChange, nil, nil
	},
	"List": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.List:
			return visitor.ActionUpdate, "[" + fmt.Sprintf("%v", node.Type) + "]", nil
		case map[string]interface{}:
			return visitor.ActionUpdate, "[" + getMapValueString(node, "Type") + "]", nil
		}
		return visitor.ActionNoChange, nil, nil
	},
	"NonNull": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.NonNull:
			return visitor.ActionUpdate, fmt.Sprintf("%v", node.Type) + "!", nil
		case map[string]interface{}:
			return visitor.ActionUpdate, getMapValueString(node, "Type") + "!", nil
		}
		return visitor.ActionNoChange, nil, nil
	},

	"ObjectDefinition": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.ObjectDefinition:
			name := fmt.Sprintf("%v", node.Name)
			interfaces := toSliceString(node.Interfaces)
			fields := node.Fields
			str := "type " + name + " " + wrap("implements ", join(interfaces, ", "), " ") + block(fields)
			return visitor.ActionUpdate, str, nil
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			interfaces := toSliceString(getMapValue(node, "Interfaces"))
			fields := getMapValue(node, "Fields")
			str := "type " + name + " " + wrap("implements ", join(interfaces, ", "), " ") + block(fields)
			return visitor.ActionUpdate, str, nil
		}
		return visitor.ActionNoChange, nil, nil
	},
	"FieldDefinition": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.FieldDefinition:
			name := fmt.Sprintf("%v", node.Name)
			ttype := fmt.Sprintf("%v", node.Type)
			args := toSliceString(node.Arguments)
			str := name + wrap("(", join(args, ", "), ")") + ": " + ttype
			return visitor.ActionUpdate, str, nil
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			ttype := getMapValueString(node, "Type")
			args := toSliceString(getMapValue(node, "Arguments"))
			str := name + wrap("(", join(args, ", "), ")") + ": " + ttype
			return visitor.ActionUpdate, str, nil
		}
		return visitor.ActionNoChange, nil, nil
	},
	"InputValueDefinition": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.InputValueDefinition:
			name := fmt.Sprintf("%v", node.Name)
			ttype := fmt.Sprintf("%v", node.Type)
			defaultValue := fmt.Sprintf("%v", node.DefaultValue)
			str := name + ": " + ttype + wrap(" = ", defaultValue, "")
			return visitor.ActionUpdate, str, nil
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			ttype := getMapValueString(node, "Type")
			defaultValue := getMapValueString(node, "DefaultValue")
			str := name + ": " + ttype + wrap(" = ", defaultValue, "")
			return visitor.ActionUpdate, str, nil
		}
		return visitor.ActionNoChange, nil, nil
	},
	"InterfaceDefinition": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.InterfaceDefinition:
			name := fmt.Sprintf("%v", node.Name)
			fields := node.Fields
			str := "interface " + name + " " + block(fields)
			return visitor.ActionUpdate, str, nil
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			fields := getMapValue(node, "Fields")
			str := "interface " + name + " " + block(fields)
			return visitor.ActionUpdate, str, nil
		}
		return visitor.ActionNoChange, nil, nil
	},
	"UnionDefinition": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.UnionDefinition:
			name := fmt.Sprintf("%v", node.Name)
			types := toSliceString(node.Types)
			str := "union " + name + " = " + join(types, " | ")
			return visitor.ActionUpdate, str, nil
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			types := toSliceString(getMapValue(node, "Types"))
			str := "union " + name + " = " + join(types, " | ")
			return visitor.ActionUpdate, str, nil
		}
		return visitor.ActionNoChange, nil, nil
	},
	"ScalarDefinition": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.ScalarDefinition:
			name := fmt.Sprintf("%v", node.Name)
			str := "scalar " + name
			return visitor.ActionUpdate, str, nil
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			str := "scalar " + name
			return visitor.ActionUpdate, str, nil
		}
		return visitor.ActionNoChange, nil, nil
	},
	"EnumDefinition": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.EnumDefinition:
			name := fmt.Sprintf("%v", node.Name)
			values := node.Values
			str := "enum " + name + " " + block(values)
			return visitor.ActionUpdate, str, nil
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			values := getMapValue(node, "Values")
			str := "enum " + name + " " + block(values)
			return visitor.ActionUpdate, str, nil
		}
		return visitor.ActionNoChange, nil, nil
	},
	"EnumValueDefinition": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.EnumValueDefinition:
			name := fmt.Sprintf("%v", node.Name)
			return visitor.ActionUpdate, name, nil
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			return visitor.ActionUpdate, name, nil
		}
		return visitor.ActionNoChange, nil, nil
	},
	"InputObjectDefinition": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.InputObjectDefinition:
			name := fmt.Sprintf("%v", node.Name)
			fields := node.Fields
			return visitor.ActionUpdate, "input " + name + " " + block(fields), nil
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			fields := getMapValue(node, "Fields")
			return visitor.ActionUpdate, "input " + name + " " + block(fields), nil
		}
		return visitor.ActionNoChange, nil, nil
	},
	"TypeExtensionDefinition": func(p visitor.VisitFuncParams) (string, interface{}, error) {
		switch node := p.Node.(type) {
		case *ast.TypeExtensionDefinition:
			definition := fmt.Sprintf("%v", node.Definition)
			str := "extend " + definition
			return visitor.ActionUpdate, str, nil
		case map[string]interface{}:
			definition := getMapValueString(node, "Definition")
			str := "extend " + definition
			return visitor.ActionUpdate, str, nil
		}
		return visitor.ActionNoChange, nil, nil
	},
}

func Print(astNode ast.Node) (interface{}, error) {
	var printed interface{}
	defer func() (interface{}, error) {
		if r := recover(); r != nil {
			return fmt.Sprintf("%v", astNode), nil
		}
		return printed, nil
	}()
	printed, err := visitor.Visit(astNode, &visitor.VisitorOptions{
		LeaveKindMap: printDocASTReducer,
	}, nil)
	if err != nil {
		return printed, err
	}
	return printed, nil
}
