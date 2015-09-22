package printer

import (
	"fmt"
	"github.com/chris-ramon/graphql-go/language/ast"
	"github.com/chris-ramon/graphql-go/language/visitor"
	"strings"
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
	for _, s := range slice.([]interface{}) {
		switch s := s.(type) {
		case string:
			res = append(res, s)
		}
	}
	return res
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
	"Name": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			return visitor.ActionUpdate, getMapValue(node, "Value")
		}
		return visitor.ActionNoChange, nil
	},
	"Variable": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			return visitor.ActionUpdate, "$" + getMapValueString(node, "Name")
		}
		return visitor.ActionNoChange, nil
	},
	"Document": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			definitions := toSliceString(getMapValue(node, "Definitions"))
			return visitor.ActionUpdate, join(definitions, "\n\n") + "\n"
		}
		return visitor.ActionNoChange, nil
	},
	"OperationDefinition": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
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
			return visitor.ActionUpdate, str
		}
		return visitor.ActionNoChange, nil
	},
	"VariableDefinition": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:

			variable := getMapValueString(node, "Variable")
			ttype := getMapValueString(node, "Type")
			defaultValue := getMapValueString(node, "DefaultValue")

			return visitor.ActionUpdate, variable + ": " + ttype + wrap(" = ", defaultValue, "")

		}
		return visitor.ActionNoChange, nil
	},
	"SelectionSet": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			selections := getMapValue(node, "Selections")
			str := block(selections)
			return visitor.ActionUpdate, str

		}
		return visitor.ActionNoChange, nil
	},
	"Field": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
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
			return visitor.ActionUpdate, str
		}
		return visitor.ActionNoChange, nil
	},
	"Argument": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			value := getMapValueString(node, "Value")
			return visitor.ActionUpdate, name + ": " + value
		}
		return visitor.ActionNoChange, nil
	},
	"FragmentSpread": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			directives := toSliceString(getMapValue(node, "Directives"))
			return visitor.ActionUpdate, "..." + name + wrap(" ", join(directives, " "), "")
		}
		return visitor.ActionNoChange, nil
	},
	"InlineFragment": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			typeCondition := getMapValueString(node, "TypeCondition")
			directives := toSliceString(getMapValue(node, "Directives"))
			selectionSet := getMapValueString(node, "SelectionSet")
			return visitor.ActionUpdate, "... on " + typeCondition + " " + wrap("", join(directives, " "), " ") + selectionSet
		}
		return visitor.ActionNoChange, nil
	},
	"FragmentDefinition": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			typeCondition := getMapValueString(node, "TypeCondition")
			directives := toSliceString(getMapValue(node, "Directives"))
			selectionSet := getMapValueString(node, "SelectionSet")
			return visitor.ActionUpdate, "fragment " + name + " on " + typeCondition + " " + wrap("", join(directives, " "), " ") + selectionSet
		}
		return visitor.ActionNoChange, nil
	},

	"IntValue": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			return visitor.ActionUpdate, getMapValueString(node, "Value")
		}
		return visitor.ActionNoChange, nil
	},
	"FloatValue": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			return visitor.ActionUpdate, getMapValueString(node, "Value")
		}
		return visitor.ActionNoChange, nil
	},
	"StringValue": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			return visitor.ActionUpdate, `"` + getMapValueString(node, "Value") + `"`
		}
		return visitor.ActionNoChange, nil
	},
	"BooleanValue": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			return visitor.ActionUpdate, getMapValueString(node, "Value")
		}
		return visitor.ActionNoChange, nil
	},
	"EnumValue": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			return visitor.ActionUpdate, getMapValueString(node, "Value")
		}
		return visitor.ActionNoChange, nil
	},
	"ListValue": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			return visitor.ActionUpdate, "[" + join(toSliceString(getMapValue(node, "Values")), ", ") + "]"
		}
		return visitor.ActionNoChange, nil
	},
	"ObjectValue": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			return visitor.ActionUpdate, "{" + join(toSliceString(getMapValue(node, "Fields")), ", ") + "}"
		}
		return visitor.ActionNoChange, nil
	},
	"ObjectField": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			value := getMapValueString(node, "Value")
			return visitor.ActionUpdate, name + ": " + value
		}
		return visitor.ActionNoChange, nil
	},

	"Directive": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			args := toSliceString(getMapValue(node, "Arguments"))
			return visitor.ActionUpdate, "@" + name + wrap("(", join(args, ", "), ")")
		}
		return visitor.ActionNoChange, nil
	},

	"NamedType": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			return visitor.ActionUpdate, getMapValueString(node, "Name")
		}
		return visitor.ActionNoChange, nil
	},
	"ListType": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			return visitor.ActionUpdate, "[" + getMapValueString(node, "Type") + "]"
		}
		return visitor.ActionNoChange, nil
	},
	"NonNullType": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			return visitor.ActionUpdate, getMapValueString(node, "Type") + "!"
		}
		return visitor.ActionNoChange, nil
	},

	"ObjectTypeDefinition": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			interfaces := toSliceString(getMapValue(node, "Interfaces"))
			fields := toSliceString(getMapValue(node, "Fields"))
			str := "type " + name + " " + wrap("implements ", join(interfaces, ", "), " ") + block(fields)
			return visitor.ActionUpdate, str
		}
		return visitor.ActionNoChange, nil
	},
	"FieldDefinition": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			ttype := getMapValueString(node, "Type")
			args := toSliceString(getMapValue(node, "Arguments"))
			str := name + wrap("(", join(args, ", "), ")") + ": " + ttype
			return visitor.ActionUpdate, str
		}
		return visitor.ActionNoChange, nil
	},
	"InputValueDefinition": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			ttype := getMapValueString(node, "Type")
			defaultValue := getMapValueString(node, "DefaultValue")
			str := name + ": " + ttype + wrap(" = ", defaultValue, "")
			return visitor.ActionUpdate, str
		}
		return visitor.ActionNoChange, nil
	},
	"InterfaceTypeDefinition": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			fields := toSliceString(getMapValue(node, "Fields"))
			str := "interface " + name + " " + block(fields)
			return visitor.ActionUpdate, str
		}
		return visitor.ActionNoChange, nil
	},
	"UnionTypeDefinition": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			types := toSliceString(getMapValue(node, "Types"))
			str := "union " + name + " = " + join(types, " | ")
			return visitor.ActionUpdate, str
		}
		return visitor.ActionNoChange, nil
	},
	"ScalarTypeDefinition": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			str := "scalar " + name
			return visitor.ActionUpdate, str
		}
		return visitor.ActionNoChange, nil
	},
	"EnumTypeDefinition": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			values := toSliceString(getMapValue(node, "Values"))
			str := "enum " + name + " " + block(values)
			return visitor.ActionUpdate, str
		}
		return visitor.ActionNoChange, nil
	},
	"EnumValueDefinition": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			return visitor.ActionUpdate, name
		}
		return visitor.ActionNoChange, nil
	},
	"InputObjectTypeDefinition": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			fields := toSliceString(getMapValue(node, "Fields"))
			return visitor.ActionUpdate, "input " + name + " " + block(fields)
		}
		return visitor.ActionNoChange, nil
	},
	"TypeExtensionDefinition": func(p visitor.VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			definition := getMapValueString(node, "definition")
			str := "extend " + definition
			return visitor.ActionUpdate, str
		}
		return visitor.ActionNoChange, nil
	},
}

func Print(ast ast.Node) interface{} {
	return visitor.Visit(ast, &visitor.VisitorOptions{
		LeaveKindMap: printDocASTReducer,
	}, nil)
}
