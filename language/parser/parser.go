package parser

import (
	"fmt"

	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/lexer"
	"github.com/graphql-go/graphql/language/source"
)

type parseFn func(parser *Parser) (interface{}, error)

type ParseOptions struct {
	NoLocation bool
	NoSource   bool
}

type ParseParams struct {
	Source  interface{}
	Options ParseOptions
}

type Parser struct {
	LexToken lexer.Lexer
	Source   *source.Source
	Options  ParseOptions
	PrevEnd  int
	Token    lexer.Token
}

func Parse(p ParseParams) (*ast.Document, error) {
	var sourceObj *source.Source
	switch p.Source.(type) {
	case *source.Source:
		sourceObj = p.Source.(*source.Source)
	default:
		body, _ := p.Source.(string)
		sourceObj = source.NewSource(&source.Source{Body: body})
	}
	parser, err := makeParser(sourceObj, p.Options)
	if err != nil {
		return nil, err
	}
	doc, err := parseDocument(parser)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// TODO: test and expose parseValue as a public
func parseValue(p ParseParams) (ast.Value, error) {
	var value ast.Value
	var sourceObj *source.Source
	switch p.Source.(type) {
	case *source.Source:
		sourceObj = p.Source.(*source.Source)
	default:
		body, _ := p.Source.(string)
		sourceObj = source.NewSource(&source.Source{Body: body})
	}
	parser, err := makeParser(sourceObj, p.Options)
	if err != nil {
		return value, err
	}
	value, err = parseValueLiteral(parser, false)
	if err != nil {
		return value, err
	}
	return value, nil
}

// Converts a name lex token into a name parse node.
func parseName(parser *Parser) (*ast.Name, error) {
	token, err := expect(parser, lexer.TokenKind[lexer.NAME])
	if err != nil {
		return nil, err
	}
	return ast.NewName(&ast.Name{
		Value: token.Value,
		Loc:   loc(parser, token.Start),
	}), nil
}

func makeParser(s *source.Source, opts ParseOptions) (*Parser, error) {
	lexToken := lexer.Lex(s)
	token, err := lexToken(0)
	if err != nil {
		return &Parser{}, err
	}
	return &Parser{
		LexToken: lexToken,
		Source:   s,
		Options:  opts,
		PrevEnd:  0,
		Token:    token,
	}, nil
}

/* Implements the parsing rules in the Document section. */

func parseDocument(parser *Parser) (*ast.Document, error) {
	start := parser.Token.Start
	var nodes []ast.Node
	for {
		if skp, err := skip(parser, lexer.TokenKind[lexer.EOF]); err != nil {
			return nil, err
		} else if skp {
			break
		}
		if peek(parser, lexer.TokenKind[lexer.BRACE_L]) {
			node, err := parseOperationDefinition(parser)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, node)
		} else if peek(parser, lexer.TokenKind[lexer.NAME]) {
			switch parser.Token.Value {
			case "query":
				fallthrough
			case "mutation":
				fallthrough
			case "subscription": // Note: subscription is an experimental non-spec addition.
				node, err := parseOperationDefinition(parser)
				if err != nil {
					return nil, err
				}
				nodes = append(nodes, node)
			case "fragment":
				node, err := parseFragmentDefinition(parser)
				if err != nil {
					return nil, err
				}
				nodes = append(nodes, node)
			case "type":
				node, err := parseObjectTypeDefinition(parser)
				if err != nil {
					return nil, err
				}
				nodes = append(nodes, node)
			case "interface":
				node, err := parseInterfaceTypeDefinition(parser)
				if err != nil {
					return nil, err
				}
				nodes = append(nodes, node)
			case "union":
				node, err := parseUnionTypeDefinition(parser)
				if err != nil {
					return nil, err
				}
				nodes = append(nodes, node)
			case "scalar":
				node, err := parseScalarTypeDefinition(parser)
				if err != nil {
					return nil, err
				}
				nodes = append(nodes, node)
			case "enum":
				node, err := parseEnumTypeDefinition(parser)
				if err != nil {
					return nil, err
				}
				nodes = append(nodes, node)
			case "input":
				node, err := parseInputObjectTypeDefinition(parser)
				if err != nil {
					return nil, err
				}
				nodes = append(nodes, node)
			case "extend":
				node, err := parseTypeExtensionDefinition(parser)
				if err != nil {
					return nil, err
				}
				nodes = append(nodes, node)
			default:
				if err := unexpected(parser, lexer.Token{}); err != nil {
					return nil, err
				}
			}
		} else {
			if err := unexpected(parser, lexer.Token{}); err != nil {
				return nil, err
			}
		}
	}
	return ast.NewDocument(&ast.Document{
		Loc:         loc(parser, start),
		Definitions: nodes,
	}), nil
}

/* Implements the parsing rules in the Operations section. */

func parseOperationDefinition(parser *Parser) (*ast.OperationDefinition, error) {
	start := parser.Token.Start
	if peek(parser, lexer.TokenKind[lexer.BRACE_L]) {
		selectionSet, err := parseSelectionSet(parser)
		if err != nil {
			return nil, err
		}
		return ast.NewOperationDefinition(&ast.OperationDefinition{
			Operation:    "query",
			Directives:   []*ast.Directive{},
			SelectionSet: selectionSet,
			Loc:          loc(parser, start),
		}), nil
	}
	operationToken, err := expect(parser, lexer.TokenKind[lexer.NAME])
	if err != nil {
		return nil, err
	}
	operation := operationToken.Value
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	variableDefinitions, err := parseVariableDefinitions(parser)
	if err != nil {
		return nil, err
	}
	directives, err := parseDirectives(parser)
	if err != nil {
		return nil, err
	}
	selectionSet, err := parseSelectionSet(parser)
	if err != nil {
		return nil, err
	}
	return ast.NewOperationDefinition(&ast.OperationDefinition{
		Operation:           operation,
		Name:                name,
		VariableDefinitions: variableDefinitions,
		Directives:          directives,
		SelectionSet:        selectionSet,
		Loc:                 loc(parser, start),
	}), nil
}

func parseVariableDefinitions(parser *Parser) ([]*ast.VariableDefinition, error) {
	variableDefinitions := []*ast.VariableDefinition{}
	if peek(parser, lexer.TokenKind[lexer.PAREN_L]) {
		vdefs, err := many(parser, lexer.TokenKind[lexer.PAREN_L], parseVariableDefinition, lexer.TokenKind[lexer.PAREN_R])
		for _, vdef := range vdefs {
			if vdef != nil {
				variableDefinitions = append(variableDefinitions, vdef.(*ast.VariableDefinition))
			}
		}
		if err != nil {
			return variableDefinitions, err
		}
		return variableDefinitions, nil
	}
	return variableDefinitions, nil
}

func parseVariableDefinition(parser *Parser) (interface{}, error) {
	start := parser.Token.Start
	variable, err := parseVariable(parser)
	if err != nil {
		return nil, err
	}
	_, err = expect(parser, lexer.TokenKind[lexer.COLON])
	if err != nil {
		return nil, err
	}
	ttype, err := parseType(parser)
	if err != nil {
		return nil, err
	}
	var defaultValue ast.Value
	if skp, err := skip(parser, lexer.TokenKind[lexer.EQUALS]); err != nil {
		return nil, err
	} else if skp {
		dv, err := parseValueLiteral(parser, true)
		if err != nil {
			return nil, err
		}
		defaultValue = dv
	}
	return ast.NewVariableDefinition(&ast.VariableDefinition{
		Variable:     variable,
		Type:         ttype,
		DefaultValue: defaultValue,
		Loc:          loc(parser, start),
	}), nil
}

func parseVariable(parser *Parser) (*ast.Variable, error) {
	start := parser.Token.Start
	_, err := expect(parser, lexer.TokenKind[lexer.DOLLAR])
	if err != nil {
		return nil, err
	}
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	return ast.NewVariable(&ast.Variable{
		Name: name,
		Loc:  loc(parser, start),
	}), nil
}

func parseSelectionSet(parser *Parser) (*ast.SelectionSet, error) {
	start := parser.Token.Start
	iSelections, err := many(parser, lexer.TokenKind[lexer.BRACE_L], parseSelection, lexer.TokenKind[lexer.BRACE_R])
	if err != nil {
		return nil, err
	}
	selections := []ast.Selection{}
	for _, iSelection := range iSelections {
		if iSelection != nil {
			// type assert interface{} into Selection interface
			selections = append(selections, iSelection.(ast.Selection))
		}
	}

	return ast.NewSelectionSet(&ast.SelectionSet{
		Selections: selections,
		Loc:        loc(parser, start),
	}), nil
}

func parseSelection(parser *Parser) (interface{}, error) {
	if peek(parser, lexer.TokenKind[lexer.SPREAD]) {
		r, err := parseFragment(parser)
		return r, err
	} else {
		return parseField(parser)
	}
}

func parseField(parser *Parser) (*ast.Field, error) {
	start := parser.Token.Start
	nameOrAlias, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	var (
		name  *ast.Name
		alias *ast.Name
	)
	if skp, err := skip(parser, lexer.TokenKind[lexer.COLON]); err != nil {
		return nil, err
	} else if skp {
		alias = nameOrAlias
		n, err := parseName(parser)
		if err != nil {
			return nil, err
		}
		name = n
	} else {
		name = nameOrAlias
	}
	arguments, err := parseArguments(parser)
	if err != nil {
		return nil, err
	}
	directives, err := parseDirectives(parser)
	if err != nil {
		return nil, err
	}
	var selectionSet *ast.SelectionSet
	if peek(parser, lexer.TokenKind[lexer.BRACE_L]) {
		sSet, err := parseSelectionSet(parser)
		if err != nil {
			return nil, err
		}
		selectionSet = sSet
	}
	return ast.NewField(&ast.Field{
		Alias:        alias,
		Name:         name,
		Arguments:    arguments,
		Directives:   directives,
		SelectionSet: selectionSet,
		Loc:          loc(parser, start),
	}), nil
}

func parseArguments(parser *Parser) ([]*ast.Argument, error) {
	arguments := []*ast.Argument{}
	if peek(parser, lexer.TokenKind[lexer.PAREN_L]) {
		iArguments, err := many(parser, lexer.TokenKind[lexer.PAREN_L], parseArgument, lexer.TokenKind[lexer.PAREN_R])
		if err != nil {
			return arguments, err
		}
		for _, iArgument := range iArguments {
			if iArgument != nil {
				arguments = append(arguments, iArgument.(*ast.Argument))
			}
		}
		return arguments, nil
	}
	return arguments, nil
}

func parseArgument(parser *Parser) (interface{}, error) {
	start := parser.Token.Start
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	_, err = expect(parser, lexer.TokenKind[lexer.COLON])
	if err != nil {
		return nil, err
	}
	value, err := parseValueLiteral(parser, false)
	if err != nil {
		return nil, err
	}
	return ast.NewArgument(&ast.Argument{
		Name:  name,
		Value: value,
		Loc:   loc(parser, start),
	}), nil
}

/* Implements the parsing rules in the Fragments section. */

func parseFragment(parser *Parser) (interface{}, error) {
	start := parser.Token.Start
	_, err := expect(parser, lexer.TokenKind[lexer.SPREAD])
	if err != nil {
		return nil, err
	}
	if parser.Token.Value == "on" {
		if err := advance(parser); err != nil {
			return nil, err
		}
		name, err := parseNamed(parser)
		if err != nil {
			return nil, err
		}
		directives, err := parseDirectives(parser)
		if err != nil {
			return nil, err
		}
		selectionSet, err := parseSelectionSet(parser)
		if err != nil {
			return nil, err
		}
		return ast.NewInlineFragment(&ast.InlineFragment{
			TypeCondition: name,
			Directives:    directives,
			SelectionSet:  selectionSet,
			Loc:           loc(parser, start),
		}), nil
	}
	name, err := parseFragmentName(parser)
	if err != nil {
		return nil, err
	}
	directives, err := parseDirectives(parser)
	if err != nil {
		return nil, err
	}
	return ast.NewFragmentSpread(&ast.FragmentSpread{
		Name:       name,
		Directives: directives,
		Loc:        loc(parser, start),
	}), nil
}

func parseFragmentDefinition(parser *Parser) (*ast.FragmentDefinition, error) {
	start := parser.Token.Start
	_, err := expectKeyWord(parser, "fragment")
	if err != nil {
		return nil, err
	}
	name, err := parseFragmentName(parser)
	if err != nil {
		return nil, err
	}
	_, err = expectKeyWord(parser, "on")
	if err != nil {
		return nil, err
	}
	typeCondition, err := parseNamed(parser)
	if err != nil {
		return nil, err
	}
	directives, err := parseDirectives(parser)
	if err != nil {
		return nil, err
	}
	selectionSet, err := parseSelectionSet(parser)
	if err != nil {
		return nil, err
	}
	return ast.NewFragmentDefinition(&ast.FragmentDefinition{
		Name:          name,
		TypeCondition: typeCondition,
		Directives:    directives,
		SelectionSet:  selectionSet,
		Loc:           loc(parser, start),
	}), nil
}

func parseFragmentName(parser *Parser) (*ast.Name, error) {
	if parser.Token.Value == "on" {
		return nil, unexpected(parser, lexer.Token{})
	}
	return parseName(parser)
}

/* Implements the parsing rules in the Values section. */

func parseValueLiteral(parser *Parser, isConst bool) (ast.Value, error) {
	token := parser.Token
	switch token.Kind {
	case lexer.TokenKind[lexer.BRACKET_L]:
		return parseList(parser, isConst)
	case lexer.TokenKind[lexer.BRACE_L]:
		return parseObject(parser, isConst)
	case lexer.TokenKind[lexer.INT]:
		if err := advance(parser); err != nil {
			return nil, err
		}
		return ast.NewIntValue(&ast.IntValue{
			Value: token.Value,
			Loc:   loc(parser, token.Start),
		}), nil
	case lexer.TokenKind[lexer.FLOAT]:
		if err := advance(parser); err != nil {
			return nil, err
		}
		return ast.NewFloatValue(&ast.FloatValue{
			Value: token.Value,
			Loc:   loc(parser, token.Start),
		}), nil
	case lexer.TokenKind[lexer.STRING]:
		if err := advance(parser); err != nil {
			return nil, err
		}
		return ast.NewStringValue(&ast.StringValue{
			Value: token.Value,
			Loc:   loc(parser, token.Start),
		}), nil
	case lexer.TokenKind[lexer.NAME]:
		if token.Value == "true" || token.Value == "false" {
			if err := advance(parser); err != nil {
				return nil, err
			}
			value := true
			if token.Value == "false" {
				value = false
			}
			return ast.NewBooleanValue(&ast.BooleanValue{
				Value: value,
				Loc:   loc(parser, token.Start),
			}), nil
		} else if token.Value != "null" {
			if err := advance(parser); err != nil {
				return nil, err
			}
			return ast.NewEnumValue(&ast.EnumValue{
				Value: token.Value,
				Loc:   loc(parser, token.Start),
			}), nil
		}
	case lexer.TokenKind[lexer.DOLLAR]:
		if !isConst {
			return parseVariable(parser)
		}
	}
	if err := unexpected(parser, lexer.Token{}); err != nil {
		return nil, err
	}
	return nil, nil
}

func parseConstValue(parser *Parser) (interface{}, error) {
	value, err := parseValueLiteral(parser, true)
	if err != nil {
		return value, err
	}
	return value, nil
}

func parseValueValue(parser *Parser) (interface{}, error) {
	return parseValueLiteral(parser, false)
}

func parseList(parser *Parser, isConst bool) (*ast.ListValue, error) {
	start := parser.Token.Start
	var item parseFn
	if isConst {
		item = parseConstValue
	} else {
		item = parseValueValue
	}
	iValues, err := any(parser, lexer.TokenKind[lexer.BRACKET_L], item, lexer.TokenKind[lexer.BRACKET_R])
	if err != nil {
		return nil, err
	}
	values := []ast.Value{}
	for _, iValue := range iValues {
		values = append(values, iValue.(ast.Value))
	}
	return ast.NewListValue(&ast.ListValue{
		Values: values,
		Loc:    loc(parser, start),
	}), nil
}

func parseObject(parser *Parser, isConst bool) (*ast.ObjectValue, error) {
	start := parser.Token.Start
	_, err := expect(parser, lexer.TokenKind[lexer.BRACE_L])
	if err != nil {
		return nil, err
	}
	fields := []*ast.ObjectField{}
	fieldNames := map[string]bool{}
	for {
		if skp, err := skip(parser, lexer.TokenKind[lexer.BRACE_R]); err != nil {
			return nil, err
		} else if skp {
			break
		}
		field, fieldName, err := parseObjectField(parser, isConst, fieldNames)
		if err != nil {
			return nil, err
		}
		fieldNames[fieldName] = true
		fields = append(fields, field)
	}
	return ast.NewObjectValue(&ast.ObjectValue{
		Fields: fields,
		Loc:    loc(parser, start),
	}), nil
}

func parseObjectField(parser *Parser, isConst bool, fieldNames map[string]bool) (*ast.ObjectField, string, error) {
	start := parser.Token.Start
	name, err := parseName(parser)
	if err != nil {
		return nil, "", err
	}
	fieldName := name.Value
	if _, ok := fieldNames[fieldName]; ok {
		descp := fmt.Sprintf("Duplicate input object field %v.", fieldName)
		return nil, "", gqlerrors.NewSyntaxError(parser.Source, start, descp)
	}
	_, err = expect(parser, lexer.TokenKind[lexer.COLON])
	if err != nil {
		return nil, "", err
	}
	value, err := parseValueLiteral(parser, isConst)
	if err != nil {
		return nil, "", err
	}
	return ast.NewObjectField(&ast.ObjectField{
		Name:  name,
		Value: value,
		Loc:   loc(parser, start),
	}), fieldName, nil
}

/* Implements the parsing rules in the Directives section. */

func parseDirectives(parser *Parser) ([]*ast.Directive, error) {
	directives := []*ast.Directive{}
	for {
		if !peek(parser, lexer.TokenKind[lexer.AT]) {
			break
		}
		directive, err := parseDirective(parser)
		if err != nil {
			return directives, err
		}
		directives = append(directives, directive)
	}
	return directives, nil
}

func parseDirective(parser *Parser) (*ast.Directive, error) {
	start := parser.Token.Start
	_, err := expect(parser, lexer.TokenKind[lexer.AT])
	if err != nil {
		return nil, err
	}
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	args, err := parseArguments(parser)
	if err != nil {
		return nil, err
	}
	return ast.NewDirective(&ast.Directive{
		Name:      name,
		Arguments: args,
		Loc:       loc(parser, start),
	}), nil
}

/* Implements the parsing rules in the Types section. */

func parseType(parser *Parser) (ast.Type, error) {
	start := parser.Token.Start
	var ttype ast.Type
	if skp, err := skip(parser, lexer.TokenKind[lexer.BRACKET_L]); err != nil {
		return nil, err
	} else if skp {
		t, err := parseType(parser)
		if err != nil {
			return t, err
		}
		ttype = t
		_, err = expect(parser, lexer.TokenKind[lexer.BRACKET_R])
		if err != nil {
			return ttype, err
		}
		ttype = ast.NewList(&ast.List{
			Type: ttype,
			Loc:  loc(parser, start),
		})
	} else {
		name, err := parseNamed(parser)
		if err != nil {
			return ttype, err
		}
		ttype = name
	}
	if skp, err := skip(parser, lexer.TokenKind[lexer.BANG]); err != nil {
		return nil, err
	} else if skp {
		ttype = ast.NewNonNull(&ast.NonNull{
			Type: ttype,
			Loc:  loc(parser, start),
		})
		return ttype, nil
	}
	return ttype, nil
}

func parseNamed(parser *Parser) (*ast.Named, error) {
	start := parser.Token.Start
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	return ast.NewNamed(&ast.Named{
		Name: name,
		Loc:  loc(parser, start),
	}), nil
}

/* Implements the parsing rules in the Type Definition section. */

func parseObjectTypeDefinition(parser *Parser) (*ast.ObjectDefinition, error) {
	start := parser.Token.Start
	_, err := expectKeyWord(parser, "type")
	if err != nil {
		return nil, err
	}
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	interfaces, err := parseImplementsInterfaces(parser)
	if err != nil {
		return nil, err
	}
	iFields, err := any(parser, lexer.TokenKind[lexer.BRACE_L], parseFieldDefinition, lexer.TokenKind[lexer.BRACE_R])
	if err != nil {
		return nil, err
	}
	fields := []*ast.FieldDefinition{}
	for _, iField := range iFields {
		if iField != nil {
			fields = append(fields, iField.(*ast.FieldDefinition))
		}
	}
	return ast.NewObjectDefinition(&ast.ObjectDefinition{
		Name:       name,
		Loc:        loc(parser, start),
		Interfaces: interfaces,
		Fields:     fields,
	}), nil
}

func parseImplementsInterfaces(parser *Parser) ([]*ast.Named, error) {
	types := []*ast.Named{}
	if parser.Token.Value == "implements" {
		if err := advance(parser); err != nil {
			return nil, err
		}
		for {
			ttype, err := parseNamed(parser)
			if err != nil {
				return types, err
			}
			types = append(types, ttype)
			if peek(parser, lexer.TokenKind[lexer.BRACE_L]) {
				break
			}
		}
	}
	return types, nil
}

func parseFieldDefinition(parser *Parser) (interface{}, error) {
	start := parser.Token.Start
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	args, err := parseArgumentDefs(parser)
	if err != nil {
		return nil, err
	}
	_, err = expect(parser, lexer.TokenKind[lexer.COLON])
	if err != nil {
		return nil, err
	}
	ttype, err := parseType(parser)
	if err != nil {
		return nil, err
	}
	return ast.NewFieldDefinition(&ast.FieldDefinition{
		Name:      name,
		Arguments: args,
		Type:      ttype,
		Loc:       loc(parser, start),
	}), nil
}

func parseArgumentDefs(parser *Parser) ([]*ast.InputValueDefinition, error) {
	inputValueDefinitions := []*ast.InputValueDefinition{}

	if !peek(parser, lexer.TokenKind[lexer.PAREN_L]) {
		return inputValueDefinitions, nil
	}
	iInputValueDefinitions, err := many(parser, lexer.TokenKind[lexer.PAREN_L], parseInputValueDef, lexer.TokenKind[lexer.PAREN_R])
	if err != nil {
		return inputValueDefinitions, err
	}
	for _, iInputValueDefinition := range iInputValueDefinitions {
		if iInputValueDefinition != nil {
			inputValueDefinitions = append(inputValueDefinitions, iInputValueDefinition.(*ast.InputValueDefinition))
		}
	}
	return inputValueDefinitions, err
}

func parseInputValueDef(parser *Parser) (interface{}, error) {
	start := parser.Token.Start
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	_, err = expect(parser, lexer.TokenKind[lexer.COLON])
	if err != nil {
		return nil, err
	}
	ttype, err := parseType(parser)
	if err != nil {
		return nil, err
	}
	var defaultValue ast.Value
	if skp, err := skip(parser, lexer.TokenKind[lexer.EQUALS]); err != nil {
		return nil, err
	} else if skp {
		val, err := parseConstValue(parser)
		if err != nil {
			return nil, err
		}
		if val, ok := val.(ast.Value); ok {
			defaultValue = val
		}
	}
	return ast.NewInputValueDefinition(&ast.InputValueDefinition{
		Name:         name,
		Type:         ttype,
		DefaultValue: defaultValue,
		Loc:          loc(parser, start),
	}), nil
}

func parseInterfaceTypeDefinition(parser *Parser) (*ast.InterfaceDefinition, error) {
	start := parser.Token.Start
	_, err := expectKeyWord(parser, "interface")
	if err != nil {
		return nil, err
	}
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	iFields, err := any(parser, lexer.TokenKind[lexer.BRACE_L], parseFieldDefinition, lexer.TokenKind[lexer.BRACE_R])
	if err != nil {
		return nil, err
	}
	fields := []*ast.FieldDefinition{}
	for _, iField := range iFields {
		if iField != nil {
			fields = append(fields, iField.(*ast.FieldDefinition))
		}
	}
	return ast.NewInterfaceDefinition(&ast.InterfaceDefinition{
		Name:   name,
		Loc:    loc(parser, start),
		Fields: fields,
	}), nil
}

func parseUnionTypeDefinition(parser *Parser) (*ast.UnionDefinition, error) {
	start := parser.Token.Start
	_, err := expectKeyWord(parser, "union")
	if err != nil {
		return nil, err
	}
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	_, err = expect(parser, lexer.TokenKind[lexer.EQUALS])
	if err != nil {
		return nil, err
	}
	types, err := parseUnionMembers(parser)
	if err != nil {
		return nil, err
	}
	return ast.NewUnionDefinition(&ast.UnionDefinition{
		Name:  name,
		Loc:   loc(parser, start),
		Types: types,
	}), nil
}

func parseUnionMembers(parser *Parser) ([]*ast.Named, error) {
	members := []*ast.Named{}
	for {
		member, err := parseNamed(parser)
		if err != nil {
			return members, err
		}
		members = append(members, member)
		if skp, err := skip(parser, lexer.TokenKind[lexer.PIPE]); err != nil {
			return nil, err
		} else if !skp {
			break
		}
	}
	return members, nil
}

func parseScalarTypeDefinition(parser *Parser) (*ast.ScalarDefinition, error) {
	start := parser.Token.Start
	_, err := expectKeyWord(parser, "scalar")
	if err != nil {
		return nil, err
	}
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	def := ast.NewScalarDefinition(&ast.ScalarDefinition{
		Name: name,
		Loc:  loc(parser, start),
	})
	return def, nil
}

func parseEnumTypeDefinition(parser *Parser) (*ast.EnumDefinition, error) {
	start := parser.Token.Start
	_, err := expectKeyWord(parser, "enum")
	if err != nil {
		return nil, err
	}
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	iEnumValueDefs, err := any(parser, lexer.TokenKind[lexer.BRACE_L], parseEnumValueDefinition, lexer.TokenKind[lexer.BRACE_R])
	if err != nil {
		return nil, err
	}
	values := []*ast.EnumValueDefinition{}
	for _, iEnumValueDef := range iEnumValueDefs {
		if iEnumValueDef != nil {
			values = append(values, iEnumValueDef.(*ast.EnumValueDefinition))
		}
	}
	return ast.NewEnumDefinition(&ast.EnumDefinition{
		Name:   name,
		Loc:    loc(parser, start),
		Values: values,
	}), nil
}

func parseEnumValueDefinition(parser *Parser) (interface{}, error) {
	start := parser.Token.Start
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	return ast.NewEnumValueDefinition(&ast.EnumValueDefinition{
		Name: name,
		Loc:  loc(parser, start),
	}), nil
}

func parseInputObjectTypeDefinition(parser *Parser) (*ast.InputObjectDefinition, error) {
	start := parser.Token.Start
	_, err := expectKeyWord(parser, "input")
	if err != nil {
		return nil, err
	}
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	iInputValueDefinitions, err := any(parser, lexer.TokenKind[lexer.BRACE_L], parseInputValueDef, lexer.TokenKind[lexer.BRACE_R])
	if err != nil {
		return nil, err
	}
	fields := []*ast.InputValueDefinition{}
	for _, iInputValueDefinition := range iInputValueDefinitions {
		if iInputValueDefinition != nil {
			fields = append(fields, iInputValueDefinition.(*ast.InputValueDefinition))
		}
	}
	return ast.NewInputObjectDefinition(&ast.InputObjectDefinition{
		Name:   name,
		Loc:    loc(parser, start),
		Fields: fields,
	}), nil
}

func parseTypeExtensionDefinition(parser *Parser) (*ast.TypeExtensionDefinition, error) {
	start := parser.Token.Start
	_, err := expectKeyWord(parser, "extend")
	if err != nil {
		return nil, err
	}

	definition, err := parseObjectTypeDefinition(parser)
	if err != nil {
		return nil, err
	}
	return ast.NewTypeExtensionDefinition(&ast.TypeExtensionDefinition{
		Loc:        loc(parser, start),
		Definition: definition,
	}), nil
}

/* Core parsing utility functions */

// Returns a location object, used to identify the place in
// the source that created a given parsed object.
func loc(parser *Parser, start int) *ast.Location {
	if parser.Options.NoLocation {
		return nil
	}
	if parser.Options.NoSource {
		return ast.NewLocation(&ast.Location{
			Start: start,
			End:   parser.PrevEnd,
		})
	}
	return ast.NewLocation(&ast.Location{
		Start:  start,
		End:    parser.PrevEnd,
		Source: parser.Source,
	})
}

// Moves the internal parser object to the next lexed token.
func advance(parser *Parser) error {
	prevEnd := parser.Token.End
	parser.PrevEnd = prevEnd
	token, err := parser.LexToken(prevEnd)
	if err != nil {
		return err
	}
	parser.Token = token
	return nil
}

// Determines if the next token is of a given kind
func peek(parser *Parser, Kind int) bool {
	return parser.Token.Kind == Kind
}

// If the next token is of the given kind, return true after advancing
// the parser. Otherwise, do not change the parser state and return false.
func skip(parser *Parser, Kind int) (bool, error) {
	if parser.Token.Kind == Kind {
		err := advance(parser)
		return true, err
	} else {
		return false, nil
	}
}

// If the next token is of the given kind, return that token after advancing
// the parser. Otherwise, do not change the parser state and return false.
func expect(parser *Parser, kind int) (lexer.Token, error) {
	token := parser.Token
	if token.Kind == kind {
		err := advance(parser)
		return token, err
	}
	descp := fmt.Sprintf("Expected %s, found %s", lexer.GetTokenKindDesc(kind), lexer.GetTokenDesc(token))
	return token, gqlerrors.NewSyntaxError(parser.Source, token.Start, descp)
}

// If the next token is a keyword with the given value, return that token after
// advancing the parser. Otherwise, do not change the parser state and return false.
func expectKeyWord(parser *Parser, value string) (lexer.Token, error) {
	token := parser.Token
	if token.Kind == lexer.TokenKind[lexer.NAME] && token.Value == value {
		err := advance(parser)
		return token, err
	}
	descp := fmt.Sprintf("Expected \"%s\", found %s", value, lexer.GetTokenDesc(token))
	return token, gqlerrors.NewSyntaxError(parser.Source, token.Start, descp)
}

// Helper function for creating an error when an unexpected lexed token
// is encountered.
func unexpected(parser *Parser, atToken lexer.Token) error {
	var token lexer.Token
	if (atToken == lexer.Token{}) {
		token = parser.Token
	} else {
		token = parser.Token
	}
	description := fmt.Sprintf("Unexpected %v", lexer.GetTokenDesc(token))
	return gqlerrors.NewSyntaxError(parser.Source, token.Start, description)
}

//  Returns a possibly empty list of parse nodes, determined by
// the parseFn. This list begins with a lex token of openKind
// and ends with a lex token of closeKind. Advances the parser
// to the next lex token after the closing token.
func any(parser *Parser, openKind int, parseFn parseFn, closeKind int) ([]interface{}, error) {
	var nodes []interface{}
	_, err := expect(parser, openKind)
	if err != nil {
		return nodes, nil
	}
	for {
		if skp, err := skip(parser, closeKind); err != nil {
			return nil, err
		} else if skp {
			break
		}
		n, err := parseFn(parser)
		if err != nil {
			return nodes, err
		}
		nodes = append(nodes, n)
	}
	return nodes, nil
}

//  Returns a non-empty list of parse nodes, determined by
// the parseFn. This list begins with a lex token of openKind
// and ends with a lex token of closeKind. Advances the parser
// to the next lex token after the closing token.
func many(parser *Parser, openKind int, parseFn parseFn, closeKind int) ([]interface{}, error) {
	_, err := expect(parser, openKind)
	if err != nil {
		return nil, err
	}
	var nodes []interface{}
	node, err := parseFn(parser)
	if err != nil {
		return nodes, err
	}
	nodes = append(nodes, node)
	for {
		if skp, err := skip(parser, closeKind); err != nil {
			return nil, err
		} else if skp {
			break
		}
		node, err := parseFn(parser)
		if err != nil {
			return nodes, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}
