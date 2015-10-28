package graphql

import (
	"fmt"
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
	LexToken Lexer
	Source   *Source
	Options  ParseOptions
	PrevEnd  int
	Token    Token
}

func Parse(p ParseParams) (*AstDocument, error) {
	var sourceObj *Source
	switch p.Source.(type) {
	case *Source:
		sourceObj = p.Source.(*Source)
	default:
		body, _ := p.Source.(string)
		sourceObj = NewSource(&Source{Body: body})
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
func parseValue(p ParseParams) (Value, error) {
	var value Value
	var sourceObj *Source
	switch p.Source.(type) {
	case *Source:
		sourceObj = p.Source.(*Source)
	default:
		body, _ := p.Source.(string)
		sourceObj = NewSource(&Source{Body: body})
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
func parseName(parser *Parser) (*AstName, error) {
	token, err := expect(parser, TokenKind[NAME])
	if err != nil {
		return nil, err
	}
	return NewAstName(&AstName{
		Value: token.Value,
		Loc:   loc(parser, token.Start),
	}), nil
}

func makeParser(s *Source, opts ParseOptions) (*Parser, error) {
	lexToken := Lex(s)
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

/* Implements the parsing rules in the AstDocument section. */

func parseDocument(parser *Parser) (*AstDocument, error) {
	start := parser.Token.Start
	var nodes []Node
	for {
		if skip(parser, TokenKind[EOF]) {
			break
		}
		if peek(parser, TokenKind[BRACE_L]) {
			node, err := parseOperationDefinition(parser)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, node)
		} else if peek(parser, TokenKind[NAME]) {
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
				if err := unexpected(parser, Token{}); err != nil {
					return nil, err
				}
			}
		} else {
			if err := unexpected(parser, Token{}); err != nil {
				return nil, err
			}
		}
	}
	return NewAstDocument(&AstDocument{
		Loc:         loc(parser, start),
		Definitions: nodes,
	}), nil
}

/* Implements the parsing rules in the Operations section. */

func parseOperationDefinition(parser *Parser) (*AstOperationDefinition, error) {
	start := parser.Token.Start
	if peek(parser, TokenKind[BRACE_L]) {
		selectionSet, err := parseSelectionSet(parser)
		if err != nil {
			return nil, err
		}
		return NewAstOperationDefinition(&AstOperationDefinition{
			Operation:    "query",
			Directives:   []*AstDirective{},
			SelectionSet: selectionSet,
			Loc:          loc(parser, start),
		}), nil
	}
	operationToken, err := expect(parser, TokenKind[NAME])
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
	return NewAstOperationDefinition(&AstOperationDefinition{
		Operation:           operation,
		Name:                name,
		VariableDefinitions: variableDefinitions,
		Directives:          directives,
		SelectionSet:        selectionSet,
		Loc:                 loc(parser, start),
	}), nil
}

func parseVariableDefinitions(parser *Parser) ([]*AstVariableDefinition, error) {
	variableDefinitions := []*AstVariableDefinition{}
	if peek(parser, TokenKind[PAREN_L]) {
		vdefs, err := many(parser, TokenKind[PAREN_L], parseVariableDefinition, TokenKind[PAREN_R])
		for _, vdef := range vdefs {
			if vdef != nil {
				variableDefinitions = append(variableDefinitions, vdef.(*AstVariableDefinition))
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
	_, err = expect(parser, TokenKind[COLON])
	if err != nil {
		return nil, err
	}
	ttype, err := parseType(parser)
	if err != nil {
		return nil, err
	}
	var defaultValue Value
	if skip(parser, TokenKind[EQUALS]) {
		dv, err := parseValueLiteral(parser, true)
		if err != nil {
			return nil, err
		}
		defaultValue = dv
	}
	return NewAstVariableDefinition(&AstVariableDefinition{
		Variable:     variable,
		Type:         ttype,
		DefaultValue: defaultValue,
		Loc:          loc(parser, start),
	}), nil
}

func parseVariable(parser *Parser) (*AstVariable, error) {
	start := parser.Token.Start
	_, err := expect(parser, TokenKind[DOLLAR])
	if err != nil {
		return nil, err
	}
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	return NewAstVariable(&AstVariable{
		Name: name,
		Loc:  loc(parser, start),
	}), nil
}

func parseSelectionSet(parser *Parser) (*AstSelectionSet, error) {
	start := parser.Token.Start
	iSelections, err := many(parser, TokenKind[BRACE_L], parseSelection, TokenKind[BRACE_R])
	if err != nil {
		return nil, err
	}
	selections := []Selection{}
	for _, iSelection := range iSelections {
		if iSelection != nil {
			// type assert interface{} into Selection interface
			selections = append(selections, iSelection.(Selection))
		}
	}

	return NewAstSelectionSet(&AstSelectionSet{
		Selections: selections,
		Loc:        loc(parser, start),
	}), nil
}

func parseSelection(parser *Parser) (interface{}, error) {
	if peek(parser, TokenKind[SPREAD]) {
		r, err := parseFragment(parser)
		return r, err
	} else {
		return parseField(parser)
	}
}

func parseField(parser *Parser) (*AstField, error) {
	start := parser.Token.Start
	nameOrAlias, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	var (
		name  *AstName
		alias *AstName
	)
	if skip(parser, TokenKind[COLON]) {
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
	var selectionSet *AstSelectionSet
	if peek(parser, TokenKind[BRACE_L]) {
		sSet, err := parseSelectionSet(parser)
		if err != nil {
			return nil, err
		}
		selectionSet = sSet
	}
	return NewField(&AstField{
		Alias:        alias,
		Name:         name,
		Arguments:    arguments,
		Directives:   directives,
		SelectionSet: selectionSet,
		Loc:          loc(parser, start),
	}), nil
}

func parseArguments(parser *Parser) ([]*AstArgument, error) {
	arguments := []*AstArgument{}
	if peek(parser, TokenKind[PAREN_L]) {
		iArguments, err := many(parser, TokenKind[PAREN_L], parseArgument, TokenKind[PAREN_R])
		if err != nil {
			return arguments, err
		}
		for _, iArgument := range iArguments {
			if iArgument != nil {
				arguments = append(arguments, iArgument.(*AstArgument))
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
	_, err = expect(parser, TokenKind[COLON])
	if err != nil {
		return nil, err
	}
	value, err := parseValueLiteral(parser, false)
	if err != nil {
		return nil, err
	}
	return NewAstArgument(&AstArgument{
		Name:  name,
		Value: value,
		Loc:   loc(parser, start),
	}), nil
}

/* Implements the parsing rules in the Fragments section. */

func parseFragment(parser *Parser) (interface{}, error) {
	start := parser.Token.Start
	_, err := expect(parser, TokenKind[SPREAD])
	if err != nil {
		return nil, err
	}
	if parser.Token.Value == "on" {
		advance(parser)
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
		return NewAstInlineFragment(&AstInlineFragment{
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
	return NewAstFragmentSpread(&AstFragmentSpread{
		Name:       name,
		Directives: directives,
		Loc:        loc(parser, start),
	}), nil
}

func parseFragmentDefinition(parser *Parser) (*AstFragmentDefinition, error) {
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
	return NewAstFragmentDefinition(&AstFragmentDefinition{
		Name:          name,
		TypeCondition: typeCondition,
		Directives:    directives,
		SelectionSet:  selectionSet,
		Loc:           loc(parser, start),
	}), nil
}

func parseFragmentName(parser *Parser) (*AstName, error) {
	if parser.Token.Value == "on" {
		return nil, unexpected(parser, Token{})
	}
	return parseName(parser)
}

/* Implements the parsing rules in the Values section. */

func parseValueLiteral(parser *Parser, isConst bool) (Value, error) {
	token := parser.Token
	switch token.Kind {
	case TokenKind[BRACKET_L]:
		return parseList(parser, isConst)
	case TokenKind[BRACE_L]:
		return parseObject(parser, isConst)
	case TokenKind[INT]:
		advance(parser)
		return NewAstIntValue(&AstIntValue{
			Value: token.Value,
			Loc:   loc(parser, token.Start),
		}), nil
	case TokenKind[FLOAT]:
		advance(parser)
		return NewAstFloatValue(&AstFloatValue{
			Value: token.Value,
			Loc:   loc(parser, token.Start),
		}), nil
	case TokenKind[STRING]:
		advance(parser)
		return NewAstStringValue(&AstStringValue{
			Value: token.Value,
			Loc:   loc(parser, token.Start),
		}), nil
	case TokenKind[NAME]:
		if token.Value == "true" || token.Value == "false" {
			advance(parser)
			value := true
			if token.Value == "false" {
				value = false
			}
			return NewAstBooleanValue(&AstBooleanValue{
				Value: value,
				Loc:   loc(parser, token.Start),
			}), nil
		} else if token.Value != "null" {
			advance(parser)
			return NewAstEnumValue(&AstEnumValue{
				Value: token.Value,
				Loc:   loc(parser, token.Start),
			}), nil
		}
	case TokenKind[DOLLAR]:
		if !isConst {
			return parseVariable(parser)
		}
	}
	if err := unexpected(parser, Token{}); err != nil {
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

func parseList(parser *Parser, isConst bool) (*AstListValue, error) {
	start := parser.Token.Start
	var item parseFn
	if isConst {
		item = parseConstValue
	} else {
		item = parseValueValue
	}
	iValues, err := any(parser, TokenKind[BRACKET_L], item, TokenKind[BRACKET_R])
	if err != nil {
		return nil, err
	}
	values := []Value{}
	for _, iValue := range iValues {
		values = append(values, iValue.(Value))
	}
	return NewAstListValue(&AstListValue{
		Values: values,
		Loc:    loc(parser, start),
	}), nil
}

func parseObject(parser *Parser, isConst bool) (*AstObjectValue, error) {
	start := parser.Token.Start
	_, err := expect(parser, TokenKind[BRACE_L])
	if err != nil {
		return nil, err
	}
	fields := []*AstObjectField{}
	fieldNames := map[string]bool{}
	for {
		if skip(parser, TokenKind[BRACE_R]) {
			break
		}
		field, fieldName, err := parseObjectField(parser, isConst, fieldNames)
		if err != nil {
			return nil, err
		}
		fieldNames[fieldName] = true
		fields = append(fields, field)
	}
	return NewAstObjectValue(&AstObjectValue{
		Fields: fields,
		Loc:    loc(parser, start),
	}), nil
}

func parseObjectField(parser *Parser, isConst bool, fieldNames map[string]bool) (*AstObjectField, string, error) {
	start := parser.Token.Start
	name, err := parseName(parser)
	if err != nil {
		return nil, "", err
	}
	fieldName := name.Value
	if _, ok := fieldNames[fieldName]; ok {
		descp := fmt.Sprintf("Duplicate input object field %v.", fieldName)
		return nil, "", NewSyntaxError(parser.Source, start, descp)
	}
	_, err = expect(parser, TokenKind[COLON])
	if err != nil {
		return nil, "", err
	}
	value, err := parseValueLiteral(parser, isConst)
	if err != nil {
		return nil, "", err
	}
	return NewAstObjectField(&AstObjectField{
		Name:  name,
		Value: value,
		Loc:   loc(parser, start),
	}), fieldName, nil
}

/* Implements the parsing rules in the Directives section. */

func parseDirectives(parser *Parser) ([]*AstDirective, error) {
	directives := []*AstDirective{}
	for {
		if !peek(parser, TokenKind[AT]) {
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

func parseDirective(parser *Parser) (*AstDirective, error) {
	start := parser.Token.Start
	_, err := expect(parser, TokenKind[AT])
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
	return NewAstDirective(&AstDirective{
		Name:      name,
		Arguments: args,
		Loc:       loc(parser, start),
	}), nil
}

/* Implements the parsing rules in the Types section. */

func parseType(parser *Parser) (AstType, error) {
	start := parser.Token.Start
	var ttype AstType
	if skip(parser, TokenKind[BRACKET_L]) {
		t, err := parseType(parser)
		if err != nil {
			return t, err
		}
		ttype = t
		_, err = expect(parser, TokenKind[BRACKET_R])
		if err != nil {
			return ttype, err
		}
		ttype = NewAstList(&AstList{
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
	if skip(parser, TokenKind[BANG]) {
		ttype = NewAstNonNull(&AstNonNull{
			Type: ttype,
			Loc:  loc(parser, start),
		})
		return ttype, nil
	}
	return ttype, nil
}

func parseNamed(parser *Parser) (*AstNamed, error) {
	start := parser.Token.Start
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	return NewAstNamed(&AstNamed{
		Name: name,
		Loc:  loc(parser, start),
	}), nil
}

/* Implements the parsing rules in the Type Definition section. */

func parseObjectTypeDefinition(parser *Parser) (*AstObjectDefinition, error) {
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
	iFields, err := any(parser, TokenKind[BRACE_L], parseFieldDefinition, TokenKind[BRACE_R])
	if err != nil {
		return nil, err
	}
	fields := []*AstFieldDefinition{}
	for _, iField := range iFields {
		if iField != nil {
			fields = append(fields, iField.(*AstFieldDefinition))
		}
	}
	return NewAstObjectDefinition(&AstObjectDefinition{
		Name:       name,
		Loc:        loc(parser, start),
		Interfaces: interfaces,
		Fields:     fields,
	}), nil
}

func parseImplementsInterfaces(parser *Parser) ([]*AstNamed, error) {
	types := []*AstNamed{}
	if parser.Token.Value == "implements" {
		advance(parser)
		for {
			ttype, err := parseNamed(parser)
			if err != nil {
				return types, err
			}
			types = append(types, ttype)
			if peek(parser, TokenKind[BRACE_L]) {
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
	_, err = expect(parser, TokenKind[COLON])
	if err != nil {
		return nil, err
	}
	ttype, err := parseType(parser)
	if err != nil {
		return nil, err
	}
	return NewAstFieldDefinition(&AstFieldDefinition{
		Name:      name,
		Arguments: args,
		Type:      ttype,
		Loc:       loc(parser, start),
	}), nil
}

func parseArgumentDefs(parser *Parser) ([]*AstInputValueDefinition, error) {
	inputValueDefinitions := []*AstInputValueDefinition{}

	if !peek(parser, TokenKind[PAREN_L]) {
		return inputValueDefinitions, nil
	}
	iInputValueDefinitions, err := many(parser, TokenKind[PAREN_L], parseInputValueDef, TokenKind[PAREN_R])
	if err != nil {
		return inputValueDefinitions, err
	}
	for _, iInputValueDefinition := range iInputValueDefinitions {
		if iInputValueDefinition != nil {
			inputValueDefinitions = append(inputValueDefinitions, iInputValueDefinition.(*AstInputValueDefinition))
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
	_, err = expect(parser, TokenKind[COLON])
	if err != nil {
		return nil, err
	}
	ttype, err := parseType(parser)
	if err != nil {
		return nil, err
	}
	var defaultValue Value
	if skip(parser, TokenKind[EQUALS]) {
		val, err := parseConstValue(parser)
		if err != nil {
			return nil, err
		}
		if val, ok := val.(Value); ok {
			defaultValue = val
		}
	}
	return NewAstInputValueDefinition(&AstInputValueDefinition{
		Name:         name,
		Type:         ttype,
		DefaultValue: defaultValue,
		Loc:          loc(parser, start),
	}), nil
}

func parseInterfaceTypeDefinition(parser *Parser) (*AstInterfaceDefinition, error) {
	start := parser.Token.Start
	_, err := expectKeyWord(parser, "interface")
	if err != nil {
		return nil, err
	}
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	iFields, err := any(parser, TokenKind[BRACE_L], parseFieldDefinition, TokenKind[BRACE_R])
	if err != nil {
		return nil, err
	}
	fields := []*AstFieldDefinition{}
	for _, iField := range iFields {
		if iField != nil {
			fields = append(fields, iField.(*AstFieldDefinition))
		}
	}
	return NewAstInterfaceDefinition(&AstInterfaceDefinition{
		Name:   name,
		Loc:    loc(parser, start),
		Fields: fields,
	}), nil
}

func parseUnionTypeDefinition(parser *Parser) (*AstUnionDefinition, error) {
	start := parser.Token.Start
	_, err := expectKeyWord(parser, "union")
	if err != nil {
		return nil, err
	}
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	_, err = expect(parser, TokenKind[EQUALS])
	if err != nil {
		return nil, err
	}
	types, err := parseUnionMembers(parser)
	if err != nil {
		return nil, err
	}
	return NewAstUnionDefinition(&AstUnionDefinition{
		Name:  name,
		Loc:   loc(parser, start),
		Types: types,
	}), nil
}

func parseUnionMembers(parser *Parser) ([]*AstNamed, error) {
	members := []*AstNamed{}
	for {
		member, err := parseNamed(parser)
		if err != nil {
			return members, err
		}
		members = append(members, member)
		if !skip(parser, TokenKind[PIPE]) {
			break
		}
	}
	return members, nil
}

func parseScalarTypeDefinition(parser *Parser) (*AstScalarDefinition, error) {
	start := parser.Token.Start
	_, err := expectKeyWord(parser, "scalar")
	if err != nil {
		return nil, err
	}
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	def := NewAstScalarDefinition(&AstScalarDefinition{
		Name: name,
		Loc:  loc(parser, start),
	})
	return def, nil
}

func parseEnumTypeDefinition(parser *Parser) (*AstEnumDefinition, error) {
	start := parser.Token.Start
	_, err := expectKeyWord(parser, "enum")
	if err != nil {
		return nil, err
	}
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	iEnumValueDefs, err := any(parser, TokenKind[BRACE_L], parseEnumValueDefinition, TokenKind[BRACE_R])
	if err != nil {
		return nil, err
	}
	values := []*AstEnumValueDefinition{}
	for _, iEnumValueDef := range iEnumValueDefs {
		if iEnumValueDef != nil {
			values = append(values, iEnumValueDef.(*AstEnumValueDefinition))
		}
	}
	return NewAstEnumDefinition(&AstEnumDefinition{
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
	return NewAstEnumValueDefinition(&AstEnumValueDefinition{
		Name: name,
		Loc:  loc(parser, start),
	}), nil
}

func parseInputObjectTypeDefinition(parser *Parser) (*AstInputObjectDefinition, error) {
	start := parser.Token.Start
	_, err := expectKeyWord(parser, "input")
	if err != nil {
		return nil, err
	}
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	iInputValueDefinitions, err := any(parser, TokenKind[BRACE_L], parseInputValueDef, TokenKind[BRACE_R])
	if err != nil {
		return nil, err
	}
	fields := []*AstInputValueDefinition{}
	for _, iInputValueDefinition := range iInputValueDefinitions {
		if iInputValueDefinition != nil {
			fields = append(fields, iInputValueDefinition.(*AstInputValueDefinition))
		}
	}
	return NewAstInputObjectDefinition(&AstInputObjectDefinition{
		Name:   name,
		Loc:    loc(parser, start),
		Fields: fields,
	}), nil
}

func parseTypeExtensionDefinition(parser *Parser) (*AstTypeExtensionDefinition, error) {
	start := parser.Token.Start
	_, err := expectKeyWord(parser, "extend")
	if err != nil {
		return nil, err
	}

	definition, err := parseObjectTypeDefinition(parser)
	if err != nil {
		return nil, err
	}
	return NewAstTypeExtensionDefinition(&AstTypeExtensionDefinition{
		Loc:        loc(parser, start),
		Definition: definition,
	}), nil
}

/* Core parsing utility functions */

// Returns a location object, used to identify the place in
// the source that created a given parsed object.
func loc(parser *Parser, start int) *AstLocation {
	if parser.Options.NoLocation {
		return nil
	}
	if parser.Options.NoSource {
		return NewAstLocation(&AstLocation{
			Start: start,
			End:   parser.PrevEnd,
		})
	}
	return NewAstLocation(&AstLocation{
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
func skip(parser *Parser, Kind int) bool {
	if parser.Token.Kind == Kind {
		advance(parser)
		return true
	} else {
		return false
	}
}

// If the next token is of the given kind, return that token after advancing
// the parser. Otherwise, do not change the parser state and return false.
func expect(parser *Parser, kind int) (Token, error) {
	token := parser.Token
	if token.Kind == kind {
		advance(parser)
		return token, nil
	}
	descp := fmt.Sprintf("Expected %s, found %s", GetTokenKindDesc(kind), GetTokenDesc(token))
	return token, NewSyntaxError(parser.Source, token.Start, descp)
}

// If the next token is a keyword with the given value, return that token after
// advancing the parser. Otherwise, do not change the parser state and return false.
func expectKeyWord(parser *Parser, value string) (Token, error) {
	token := parser.Token
	if token.Kind == TokenKind[NAME] && token.Value == value {
		advance(parser)
		return token, nil
	}
	descp := fmt.Sprintf("Expected \"%s\", found %s", value, GetTokenDesc(token))
	return token, NewSyntaxError(parser.Source, token.Start, descp)
}

// Helper function for creating an error when an unexpected lexed token
// is encountered.
func unexpected(parser *Parser, atToken Token) error {
	var token Token
	if (atToken == Token{}) {
		token = parser.Token
	} else {
		token = parser.Token
	}
	description := fmt.Sprintf("Unexpected %v", GetTokenDesc(token))
	return NewSyntaxError(parser.Source, token.Start, description)
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
		if skip(parser, closeKind) {
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
		if skip(parser, closeKind) {
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
