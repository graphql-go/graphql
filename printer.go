package graphql

import (
	"fmt"
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

var printDocASTReducer = map[string]VisitFunc{
	"Name": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			return ActionUpdate, getMapValue(node, "Value")
		}
		return ActionNoChange, nil
	},
	"Variable": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			return ActionUpdate, "$" + getMapValueString(node, "Name")
		}
		return ActionNoChange, nil
	},
	"Document": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			definitions := toSliceString(getMapValue(node, "Definitions"))
			return ActionUpdate, join(definitions, "\n\n") + "\n"
		}
		return ActionNoChange, nil
	},
	"OperationDefinition": func(p VisitFuncParams) (string, interface{}) {
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
			return ActionUpdate, str
		}
		return ActionNoChange, nil
	},
	"VariableDefinition": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:

			variable := getMapValueString(node, "Variable")
			ttype := getMapValueString(node, "Type")
			defaultValue := getMapValueString(node, "DefaultValue")

			return ActionUpdate, variable + ": " + ttype + wrap(" = ", defaultValue, "")

		}
		return ActionNoChange, nil
	},
	"SelectionSet": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			selections := getMapValue(node, "Selections")
			str := block(selections)
			return ActionUpdate, str

		}
		return ActionNoChange, nil
	},
	"Field": func(p VisitFuncParams) (string, interface{}) {
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
			return ActionUpdate, str
		}
		return ActionNoChange, nil
	},
	"Argument": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			value := getMapValueString(node, "Value")
			return ActionUpdate, name + ": " + value
		}
		return ActionNoChange, nil
	},
	"FragmentSpread": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			directives := toSliceString(getMapValue(node, "Directives"))
			return ActionUpdate, "..." + name + wrap(" ", join(directives, " "), "")
		}
		return ActionNoChange, nil
	},
	"InlineFragment": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			typeCondition := getMapValueString(node, "TypeCondition")
			directives := toSliceString(getMapValue(node, "Directives"))
			selectionSet := getMapValueString(node, "SelectionSet")
			return ActionUpdate, "... on " + typeCondition + " " + wrap("", join(directives, " "), " ") + selectionSet
		}
		return ActionNoChange, nil
	},
	"FragmentDefinition": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			typeCondition := getMapValueString(node, "TypeCondition")
			directives := toSliceString(getMapValue(node, "Directives"))
			selectionSet := getMapValueString(node, "SelectionSet")
			return ActionUpdate, "fragment " + name + " on " + typeCondition + " " + wrap("", join(directives, " "), " ") + selectionSet
		}
		return ActionNoChange, nil
	},

	"IntValue": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			return ActionUpdate, getMapValueString(node, "Value")
		}
		return ActionNoChange, nil
	},
	"FloatValue": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			return ActionUpdate, getMapValueString(node, "Value")
		}
		return ActionNoChange, nil
	},
	"StringValue": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			return ActionUpdate, `"` + getMapValueString(node, "Value") + `"`
		}
		return ActionNoChange, nil
	},
	"BooleanValue": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			return ActionUpdate, getMapValueString(node, "Value")
		}
		return ActionNoChange, nil
	},
	"EnumValue": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			return ActionUpdate, getMapValueString(node, "Value")
		}
		return ActionNoChange, nil
	},
	"ListValue": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			return ActionUpdate, "[" + join(toSliceString(getMapValue(node, "Values")), ", ") + "]"
		}
		return ActionNoChange, nil
	},
	"ObjectValue": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			return ActionUpdate, "{" + join(toSliceString(getMapValue(node, "Fields")), ", ") + "}"
		}
		return ActionNoChange, nil
	},
	"ObjectField": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			value := getMapValueString(node, "Value")
			return ActionUpdate, name + ": " + value
		}
		return ActionNoChange, nil
	},

	"Directive": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			args := toSliceString(getMapValue(node, "Arguments"))
			return ActionUpdate, "@" + name + wrap("(", join(args, ", "), ")")
		}
		return ActionNoChange, nil
	},

	"Named": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			return ActionUpdate, getMapValueString(node, "Name")
		}
		return ActionNoChange, nil
	},
	"List": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			return ActionUpdate, "[" + getMapValueString(node, "Type") + "]"
		}
		return ActionNoChange, nil
	},
	"NonNull": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			return ActionUpdate, getMapValueString(node, "Type") + "!"
		}
		return ActionNoChange, nil
	},

	"ObjectDefinition": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			interfaces := toSliceString(getMapValue(node, "Interfaces"))
			fields := getMapValue(node, "Fields")
			str := "type " + name + " " + wrap("implements ", join(interfaces, ", "), " ") + block(fields)
			return ActionUpdate, str
		}
		return ActionNoChange, nil
	},
	"FieldDefinition": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			ttype := getMapValueString(node, "Type")
			args := toSliceString(getMapValue(node, "Arguments"))
			str := name + wrap("(", join(args, ", "), ")") + ": " + ttype
			return ActionUpdate, str
		}
		return ActionNoChange, nil
	},
	"InputValueDefinition": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			ttype := getMapValueString(node, "Type")
			defaultValue := getMapValueString(node, "DefaultValue")
			str := name + ": " + ttype + wrap(" = ", defaultValue, "")
			return ActionUpdate, str
		}
		return ActionNoChange, nil
	},
	"InterfaceDefinition": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			fields := getMapValue(node, "Fields")
			str := "interface " + name + " " + block(fields)
			return ActionUpdate, str
		}
		return ActionNoChange, nil
	},
	"UnionDefinition": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			types := toSliceString(getMapValue(node, "Types"))
			str := "union " + name + " = " + join(types, " | ")
			return ActionUpdate, str
		}
		return ActionNoChange, nil
	},
	"ScalarDefinition": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			str := "scalar " + name
			return ActionUpdate, str
		}
		return ActionNoChange, nil
	},
	"EnumDefinition": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			values := getMapValue(node, "Values")
			str := "enum " + name + " " + block(values)
			return ActionUpdate, str
		}
		return ActionNoChange, nil
	},
	"EnumValueDefinition": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			return ActionUpdate, name
		}
		return ActionNoChange, nil
	},
	"InputObjectDefinition": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			name := getMapValueString(node, "Name")
			fields := getMapValue(node, "Fields")
			return ActionUpdate, "input " + name + " " + block(fields)
		}
		return ActionNoChange, nil
	},
	"TypeExtensionDefinition": func(p VisitFuncParams) (string, interface{}) {
		switch node := p.Node.(type) {
		case map[string]interface{}:
			definition := getMapValueString(node, "Definition")
			str := "extend " + definition
			return ActionUpdate, str
		}
		return ActionNoChange, nil
	},
}

func Print(astNode Node) (printed interface{}) {
	defer func() interface{} {
		if r := recover(); r != nil {
			return fmt.Sprintf("%v", astNode)
		}
		return printed
	}()
	printed = Visit(astNode, &VisitorOptions{
		LeaveKindMap: printDocASTReducer,
	}, nil)
	return printed
}
