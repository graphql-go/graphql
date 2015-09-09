package parser

import (
	"fmt"

	"github.com/chris-ramon/graphql-go/errors"
	"github.com/chris-ramon/graphql-go/language/ast"
	"github.com/chris-ramon/graphql-go/language/fd"
	"github.com/chris-ramon/graphql-go/language/kinds"
	"github.com/chris-ramon/graphql-go/language/lexer"
	"github.com/chris-ramon/graphql-go/language/source"
)

func unexpected(parser *Parser, atToken lexer.Token) error {
	var token lexer.Token
	if (atToken == lexer.Token{}) {
		token = parser.Token
	} else {
		token = parser.Token
	}
	description := fmt.Sprintf("Unexpected %v", lexer.GetTokenDesc(token))
	return graphqlerrors.NewSyntaxError(parser.Source, token.Start, description)
}

type ParseOptions struct {
	NoLocation bool
	NoSource   bool
}

type ParseParams struct {
	Source  interface{}
	Options ParseOptions
}

func Parse(p ParseParams) (ast.Document, error) {
	var doc ast.Document
	var sourceObj *source.Source
	switch p.Source.(type) {
	case *source.Source:
		sourceObj = p.Source.(*source.Source)
	default:
		s, _ := p.Source.(string)
		sourceObj = source.NewSource(s, "")
	}
	parser, err := makeParser(sourceObj, p.Options)
	if err != nil {
		return doc, err
	}
	doc, err = parseDocument(parser)
	if err != nil {
		return doc, err
	}
	return doc, nil
}

type Parser struct {
	LexToken lexer.Lexer
	Source   *source.Source
	Options  ParseOptions
	PrevEnd  int
	Token    lexer.Token
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

// Implements the parsing rules in the Document section.
func parseDocument(parser *Parser) (ast.Document, error) {
	start := parser.Token.Start
	var definitions []ast.Definition
	for {
		if skip(parser, lexer.TokenKind[lexer.EOF]) {
			break
		}
		if peek(parser, lexer.TokenKind[lexer.BRACE_L]) {
			oDef, err := parseOperationDefinition(parser)
			if err != nil {
				return ast.Document{}, err
			}
			definitions = append(definitions, oDef)
		} else if peek(parser, lexer.TokenKind[lexer.NAME]) {
			switch parser.Token.Value {
			case "query":
				fallthrough
			case "mutation":
				fallthrough
			case "subscription": // Note: subscription is an experimental non-spec addition.
				oDef, err := parseOperationDefinition(parser)
				if err != nil {
					return ast.Document{}, err
				}
				definitions = append(definitions, oDef)
			case "fragment":
				fDef, err := parseFragmentDefinition(parser)
				if err != nil {
					return ast.Document{}, err
				}
				definitions = append(definitions, fDef)
			case "type":
				def, err := parseObjectTypeDefinition(parser)
				if err != nil {
					return ast.Document{}, err
				}
				definitions = append(definitions, def)
			case "interface":
				def, err := parseInterfaceTypeDefinition(parser)
				if err != nil {
					return ast.Document{}, err
				}
				definitions = append(definitions, def)
			case "union":
				def, err := parseUnionTypeDefinition(parser)
				if err != nil {
					return ast.Document{}, err
				}
				definitions = append(definitions, def)
			case "scalar":
				def, err := parseScalarTypeDefinition(parser)
				if err != nil {
					return ast.Document{}, err
				}
				definitions = append(definitions, def)
			case "enum":
				def, err := parseEnumTypeDefinition(parser)
				if err != nil {
					return ast.Document{}, err
				}
				definitions = append(definitions, def)
			case "input":
				def, err := parseInputObjectTypeDefinition(parser)
				if err != nil {
					return ast.Document{}, err
				}
				definitions = append(definitions, def)
			case "extend":
				def, err := parseTypeExtensionDefinition(parser)
				if err != nil {
					return ast.Document{}, err
				}
				definitions = append(definitions, def)
			default:
				if err := unexpected(parser, lexer.Token{}); err != nil {
					return ast.Document{}, err
				}
			}
		} else {
			if err := unexpected(parser, lexer.Token{}); err != nil {
				return ast.Document{}, err
			}
		}
	}
	return ast.Document{
		Kind:        kinds.Document,
		Loc:         loc(parser, start),
		Definitions: definitions,
	}, nil
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

// Implements the parsing rules in the Operations section.
func parseOperationDefinition(parser *Parser) (*ast.OperationDefinition, error) {
	start := parser.Token.Start
	if peek(parser, lexer.TokenKind[lexer.BRACE_L]) {
		selectionSet, err := parseSelectionSet(parser)
		if err != nil {
			oDef := ast.NewOperationDefinition()
			return oDef, err
		}
		oDef := ast.NewOperationDefinition()
		oDef.Operation = "query"
		oDef.Directives = []ast.Directive{}
		oDef.SelectionSet = selectionSet
		oDef.Loc = loc(parser, start)
		return oDef, err
	}
	operationToken, err := expect(parser, lexer.TokenKind[lexer.NAME])
	if err != nil {
		oDef := ast.NewOperationDefinition()
		return oDef, err
	}
	operation := operationToken.Value
	name, err := parseName(parser)
	if err != nil {
		oDef := ast.NewOperationDefinition()
		return oDef, err
	}
	variableDefinitions, err := parseVariableDefinitions(parser)
	if err != nil {
		oDef := ast.NewOperationDefinition()
		return oDef, err
	}
	directives, err := parseDirectives(parser)
	if err != nil {
		oDef := ast.NewOperationDefinition()
		return oDef, err
	}
	selectionSet, err := parseSelectionSet(parser)
	if err != nil {
		oDef := ast.NewOperationDefinition()
		return oDef, err
	}
	oDef := ast.NewOperationDefinition()
	oDef.Operation = operation
	oDef.Name = name
	oDef.VariableDefinitions = variableDefinitions
	oDef.Directives = directives
	oDef.SelectionSet = selectionSet
	oDef.Loc = loc(parser, start)
	return oDef, nil
}

func parseFragmentDefinition(parser *Parser) (*fd.FragmentDefinition, error) {
	start := parser.Token.Start
	_, err := expectKeyWord(parser, "fragment")
	if err != nil {
		fDef := fd.NewFragmentDefinition()
		return fDef, err
	}
	name, err := parseFragmentName(parser)
	if err != nil {
		fDef := fd.NewFragmentDefinition()
		return fDef, err
	}
	_, err = expectKeyWord(parser, "on")
	if err != nil {
		fDef := fd.NewFragmentDefinition()
		return fDef, err
	}
	typeCondition, err := parseNamedType(parser)
	if err != nil {
		fDef := fd.NewFragmentDefinition()
		return fDef, err
	}
	selectionSet, err := parseSelectionSet(parser)
	if err != nil {
		fDef := fd.NewFragmentDefinition()
		return fDef, err
	}
	directives, err := parseDirectives(parser)
	if err != nil {
		fDef := fd.NewFragmentDefinition()
		return fDef, err
	}
	fDef := fd.NewFragmentDefinition()
	fDef.Name = name
	fDef.TypeCondition = typeCondition
	fDef.Directives = directives
	fDef.SelectionSet = selectionSet
	fDef.Loc = loc(parser, start)
	return fDef, nil
}

func expectKeyWord(parser *Parser, value string) (lexer.Token, error) {
	token := parser.Token
	if token.Kind == lexer.TokenKind[lexer.NAME] && token.Value == value {
		advance(parser)
		return token, nil
	}
	descp := fmt.Sprintf("Expected \"%s\", found %s", value, lexer.GetTokenDesc(token))
	return token, graphqlerrors.NewSyntaxError(parser.Source, token.Start, descp)
}

func parseSelectionSet(parser *Parser) (ast.SelectionSet, error) {
	start := parser.Token.Start
	selections, err := many(parser, lexer.TokenKind[lexer.BRACE_L], parseSelection, lexer.TokenKind[lexer.BRACE_R])
	if err != nil {
		return ast.SelectionSet{}, err
	}
	return ast.SelectionSet{
		Kind:       kinds.SelectionSet,
		Selections: selections,
		Loc:        loc(parser, start),
	}, nil
}

func parseSelection(parser *Parser) (interface{}, error) {
	if peek(parser, lexer.TokenKind[lexer.SPREAD]) {
		r, err := parseFragment(parser)
		return r, err
	} else {
		return parseField(parser)
	}
}

func loc(parser *Parser, start int) ast.Location {
	if parser.Options.NoLocation {
		return ast.Location{}
	}
	if parser.Options.NoSource {
		return ast.Location{
			Start: start,
			End:   parser.PrevEnd,
		}
	}
	return ast.Location{
		Start:  start,
		End:    parser.PrevEnd,
		Source: parser.Source,
	}
}

func expect(parser *Parser, kind int) (lexer.Token, error) {
	token := parser.Token
	if token.Kind == kind {
		advance(parser)
		return token, nil
	}
	descp := fmt.Sprintf("Expected %s, found %s", lexer.GetTokenKindDesc(kind), lexer.GetTokenDesc(token))
	return token, graphqlerrors.NewSyntaxError(parser.Source, token.Start, descp)
}

// Converts a name lex token into a name parse node.
func parseName(parser *Parser) (ast.Name, error) {
	token, err := expect(parser, lexer.TokenKind[lexer.NAME])
	if err != nil {
		return ast.Name{}, err
	}
	return ast.Name{
		Kind:  kinds.Name,
		Value: token.Value,
		Loc:   loc(parser, token.Start),
	}, nil
}

func parseNamedType(parser *Parser) (ast.NamedType, error) {
	start := parser.Token.Start
	name, err := parseName(parser)
	if err != nil {
		return ast.NamedType{}, err
	}
	return ast.NamedType{
		Kind: kinds.NamedType,
		Name: name,
		Loc:  loc(parser, start),
	}, nil
}

func parseFragmentName(parser *Parser) (ast.Name, error) {
	if parser.Token.Value == "on" {
		return ast.Name{}, unexpected(parser, lexer.Token{})
	}
	return parseName(parser)
}

func parseVariableDefinitions(parser *Parser) ([]ast.VariableDefinition, error) {
	if peek(parser, lexer.TokenKind[lexer.PAREN_L]) {
		vdefs, err := many(parser, lexer.TokenKind[lexer.PAREN_L], parseVariableDefinition, lexer.TokenKind[lexer.PAREN_R])
		variableDefinitions := []ast.VariableDefinition{}
		for _, vdef := range vdefs {
			variableDefinitions = append(variableDefinitions, vdef.(ast.VariableDefinition))
		}
		if err != nil {
			return variableDefinitions, err
		}
		return variableDefinitions, nil
	} else {
		return []ast.VariableDefinition{}, nil
	}
}

func parseDirectives(parser *Parser) ([]ast.Directive, error) {
	directives := []ast.Directive{}
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

func parseDirective(parser *Parser) (ast.Directive, error) {
	start := parser.Token.Start
	_, err := expect(parser, lexer.TokenKind[lexer.AT])
	if err != nil {
		return ast.Directive{}, err
	}
	name, err := parseName(parser)
	if err != nil {
		return ast.Directive{}, err
	}
	args, err := parseArguments(parser)
	if err != nil {
		return ast.Directive{}, err
	}
	return ast.Directive{
		Kind:      kinds.Directive,
		Name:      name,
		Arguments: args,
		Loc:       loc(parser, start),
	}, nil
}

func parseVariableDefinition(parser *Parser) (interface{}, error) {
	start := parser.Token.Start
	variable, err := parseVariable(parser)
	if err != nil {
		return ast.VariableDefinition{}, err
	}
	_, err = expect(parser, lexer.TokenKind[lexer.COLON])
	if err != nil {
		return ast.VariableDefinition{}, err
	}
	ttype, err := parseType(parser)
	if err != nil {
		return ast.VariableDefinition{}, err
	}
	var defaultValue ast.Value
	if skip(parser, lexer.TokenKind[lexer.EQUALS]) {
		dv, err := parseValueLiteral(parser, true)
		if err != nil {
			return dv, err
		}
		defaultValue = dv
	}
	return ast.VariableDefinition{
		Kind:         kinds.VariableDefinition,
		Variable:     variable,
		Type:         ttype,
		DefaultValue: defaultValue,
		Loc:          loc(parser, start),
	}, nil
}

func parseVariable(parser *Parser) (ast.Variable, error) {
	start := parser.Token.Start
	_, err := expect(parser, lexer.TokenKind[lexer.DOLLAR])
	if err != nil {
		return ast.Variable{}, err
	}
	name, err := parseName(parser)
	if err != nil {
		return ast.Variable{}, err
	}
	return ast.Variable{
		Kind: kinds.Variable,
		Name: name,
		Loc:  loc(parser, start),
	}, nil
}

func parseType(parser *Parser) (ast.Type, error) {
	start := parser.Token.Start
	var ttype ast.Type
	if skip(parser, lexer.TokenKind[lexer.BRACE_L]) {
		t, err := parseType(parser)
		if err != nil {
			return t, err
		}
		ttype = t
		_, err = expect(parser, lexer.TokenKind[lexer.BRACKET_R])
		if err != nil {
			return ttype, err
		}
		ttype = ast.ListType{
			Kind: kinds.ListType,
			Type: ttype,
			Loc:  loc(parser, start),
		}
	} else {
		name, err := parseNamedType(parser)
		if err != nil {
			return ttype, err
		}
		ttype = name
	}
	if skip(parser, lexer.TokenKind[lexer.BANG]) {
		ttype = ast.NonNullType{
			Kind: kinds.NonNullType,
			Type: ttype,
			Loc:  loc(parser, start),
		}
		return ttype, nil
	}
	return ttype, nil
}

func parseValueLiteral(parser *Parser, isConst bool) (ast.Value, error) {
	token := parser.Token
	switch token.Kind {
	case lexer.TokenKind[lexer.BRACKET_L]:
		return parseList(parser, isConst)
	case lexer.TokenKind[lexer.BRACE_L]:
		return parseObject(parser, isConst)
	case lexer.TokenKind[lexer.INT]:
		advance(parser)
		return ast.IntValue{
			Kind:  kinds.IntValue,
			Value: token.Value,
			Loc:   loc(parser, token.Start),
		}, nil
	case lexer.TokenKind[lexer.FLOAT]:
		advance(parser)
		return ast.FloatValue{
			Kind:  kinds.FloatValue,
			Value: token.Value,
			Loc:   loc(parser, token.Start),
		}, nil
	case lexer.TokenKind[lexer.STRING]:
		advance(parser)
		return ast.StringValue{
			Kind:  kinds.StringValue,
			Value: token.Value,
			Loc:   loc(parser, token.Start),
		}, nil
	case lexer.TokenKind[lexer.NAME]:
		if token.Value == "true" || token.Value == "false" {
			advance(parser)
			value := true
			if token.Value == "false" {
				value = false
			}
			return ast.BooleanValue{
				Kind:  kinds.BooleanValue,
				Value: value,
				Loc:   loc(parser, token.Start),
			}, nil
		} else if token.Value != "null" {
			advance(parser)
			return ast.EnumValue{
				Kind:  kinds.EnumValue,
				Value: token.Value,
				Loc:   loc(parser, token.Start),
			}, nil
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

type parseFn func(parser *Parser) (interface{}, error)

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

func parseFragment(parser *Parser) (interface{}, error) {
	start := parser.Token.Start
	_, err := expect(parser, lexer.TokenKind[lexer.SPREAD])
	if err != nil {
		return nil, err
	}
	if parser.Token.Value == "on" {
		advance(parser)
		name, err := parseNamedType(parser)
		if err != nil {
			return ast.InlineFragment{}, err
		}
		directives, err := parseDirectives(parser)
		if err != nil {
			return ast.InlineFragment{}, err
		}
		selectionSet, err := parseSelectionSet(parser)
		if err != nil {
			return ast.InlineFragment{}, err
		}
		return ast.InlineFragment{
			Kind:          kinds.InlineFragment,
			TypeCondition: name,
			Directives:    directives,
			SelectionSet:  selectionSet,
			Loc:           loc(parser, start),
		}, nil
	}
	name, err := parseFragmentName(parser)
	if err != nil {
		return ast.FragmentSpread{}, err
	}
	directives, err := parseDirectives(parser)
	if err != nil {
		return ast.FragmentSpread{}, err
	}
	return ast.FragmentSpread{
		Kind:       kinds.FragmentSpread,
		Name:       name,
		Directives: directives,
		Loc:        loc(parser, start),
	}, nil
}

func parseField(parser *Parser) (ast.Field, error) {
	start := parser.Token.Start
	nameOrAlias, err := parseName(parser)
	if err != nil {
		return ast.Field{}, err
	}
	var (
		name  ast.Name
		alias ast.Name
	)
	if skip(parser, lexer.TokenKind[lexer.COLON]) {
		alias = nameOrAlias
		n, err := parseName(parser)
		if err != nil {
			return ast.Field{}, err
		}
		name = n
	} else {
		name = nameOrAlias
	}
	arguments, err := parseArguments(parser)
	if err != nil {
		return ast.Field{}, err
	}
	directives, err := parseDirectives(parser)
	if err != nil {
		return ast.Field{}, err
	}
	var selectionSet ast.SelectionSet
	if peek(parser, lexer.TokenKind[lexer.BRACE_L]) {
		sSet, err := parseSelectionSet(parser)
		if err != nil {
			return ast.Field{}, err
		}
		selectionSet = sSet
	}
	return ast.Field{
		Kind:         kinds.Field,
		Alias:        alias,
		Name:         name,
		Arguments:    arguments,
		Directives:   directives,
		SelectionSet: selectionSet,
		Loc:          loc(parser, start),
	}, nil
}

func parseList(parser *Parser, isConst bool) (ast.ListValue, error) {
	start := parser.Token.Start
	var item parseFn
	if isConst {
		item = parseConstValue
	} else {
		item = parseValueValue
	}
	iValues, err := any(parser, lexer.TokenKind[lexer.BRACKET_L], item, lexer.TokenKind[lexer.BRACKET_R])
	if err != nil {
		return ast.ListValue{}, err
	}
	values := []ast.Value{}
	for _, iValue := range iValues {
		values = append(values, iValue.(ast.Value))
	}
	return ast.ListValue{
		Kind:   kinds.ListValue,
		Values: values,
		Loc:    loc(parser, start),
	}, nil
}

func parseObject(parser *Parser, isConst bool) (ast.ObjectValue, error) {
	start := parser.Token.Start
	_, err := expect(parser, lexer.TokenKind[lexer.BRACE_L])
	if err != nil {
		return ast.ObjectValue{}, err
	}
	fields := []ast.ObjectField{}
	fieldNames := map[string]bool{}
	for {
		if skip(parser, lexer.TokenKind[lexer.BRACE_R]) {
			break
		}
		field, fieldName, err := parseObjectField(parser, isConst, fieldNames)
		if err != nil {
			return ast.ObjectValue{}, err
		}
		fieldNames[fieldName] = true
		fields = append(fields, field)
	}
	return ast.ObjectValue{
		Kind:   kinds.ObjectValue,
		Fields: fields,
		Loc:    loc(parser, start),
	}, nil
}
func parseObjectField(parser *Parser, isConst bool, fieldNames map[string]bool) (ast.ObjectField, string, error) {
	start := parser.Token.Start
	name, err := parseName(parser)
	if err != nil {
		return ast.ObjectField{}, "", err
	}
	fieldName := name.Value
	if _, ok := fieldNames[fieldName]; ok {
		descp := fmt.Sprintf("Duplicate input object field %v.", fieldName)
		return ast.ObjectField{}, "", graphqlerrors.NewSyntaxError(parser.Source, start, descp)
	}
	_, err = expect(parser, lexer.TokenKind[lexer.COLON])
	if err != nil {
		return ast.ObjectField{}, "", err
	}
	value, err := parseValueLiteral(parser, isConst)
	if err != nil {
		return ast.ObjectField{}, "", err
	}
	return ast.ObjectField{
		Kind:  kinds.ObjectField,
		Name:  name,
		Value: value,
		Loc:   loc(parser, start),
	}, fieldName, nil
}

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

func parseArguments(parser *Parser) ([]ast.Argument, error) {
	if peek(parser, lexer.TokenKind[lexer.PAREN_L]) {
		iArguments, err := many(parser, lexer.TokenKind[lexer.PAREN_L], parseArgument, lexer.TokenKind[lexer.PAREN_R])
		arguments := []ast.Argument{}
		if err != nil {
			return arguments, err
		}
		for _, iArgument := range iArguments {
			arguments = append(arguments, iArgument.(ast.Argument))
		}
		return arguments, nil
	} else {
		return []ast.Argument{}, nil
	}
}

func parseArgument(parser *Parser) (interface{}, error) {
	start := parser.Token.Start
	name, err := parseName(parser)
	if err != nil {
		return ast.Argument{}, err
	}
	_, err = expect(parser, lexer.TokenKind[lexer.COLON])
	if err != nil {
		return ast.Argument{}, err
	}
	value, err := parseValueLiteral(parser, false)
	if err != nil {
		return ast.Argument{}, err
	}
	return ast.Argument{
		Kind:  kinds.Argument,
		Name:  name,
		Value: value,
		Loc:   loc(parser, start),
	}, nil
}

func parseConstValue(parser *Parser) (interface{}, error) {
	value, err := parseValueLiteral(parser, true)
	if err != nil {
		return value, err
	}
	return value, nil
}

func parseVariableValue(parser *Parser) (interface{}, error) {
	value, err := parseValueLiteral(parser, false)
	if err != nil {
		return value, err
	}
	return value, nil
}

func parseValueValue(parser *Parser) (interface{}, error) {
	return parseValueLiteral(parser, false)
}

func parseObjectTypeDefinition(parser *Parser) (*ast.ObjectTypeDefinition, error) {
	start := parser.Token.Start
	_, err := expectKeyWord(parser, "type")
	if err != nil {
		return ast.NewObjectTypeDefinition(), err
	}
	name, err := parseName(parser)
	if err != nil {
		return ast.NewObjectTypeDefinition(), err
	}
	fields, err := any(parser, lexer.TokenKind[lexer.BRACE_L], parseFieldDefinition, lexer.TokenKind[lexer.BRACE_R])
	def := ast.NewObjectTypeDefinition()
	def.Name = name
	def.Loc = loc(parser, start)
	def.Fields = fields
	return def, nil
}

func parseInterfaceTypeDefinition(parser *Parser) (*ast.InterfaceTypeDefinition, error) {
	start := parser.Token.Start
	_, err := expectKeyWord(parser, "interface")
	if err != nil {
		return ast.NewInterfaceTypeDefinition(), err
	}
	name, err := parseName(parser)
	if err != nil {
		return ast.NewInterfaceTypeDefinition(), err
	}
	fields, err := any(parser, lexer.TokenKind[lexer.BRACE_L], parseFieldDefinition, lexer.TokenKind[lexer.BRACE_R])
	def := ast.NewInterfaceTypeDefinition()
	def.Name = name
	def.Loc = loc(parser, start)
	def.Fields = fields
	return def, nil
}

func parseUnionTypeDefinition(parser *Parser) (*ast.UnionTypeDefinition, error) {
	start := parser.Token.Start
	_, err := expectKeyWord(parser, "union")
	if err != nil {
		return ast.NewUnionTypeDefinition(), err
	}
	name, err := parseName(parser)
	if err != nil {
		return ast.NewUnionTypeDefinition(), err
	}
	_, err = expect(parser, lexer.TokenKind[lexer.EQUALS])
	if err != nil {
		return ast.NewUnionTypeDefinition(), err
	}
	types, err := parseUnionMembers(parser)
	if err != nil {
		return ast.NewUnionTypeDefinition(), err
	}
	def := ast.NewUnionTypeDefinition()
	def.Name = name
	def.Loc = loc(parser, start)
	def.Types = types
	return def, nil
}

func parseScalarTypeDefinition(parser *Parser) (*ast.ScalarTypeDefinition, error) {
	start := parser.Token.Start
	_, err := expectKeyWord(parser, "scalar")
	if err != nil {
		return ast.NewScalarTypeDefinition(), err
	}
	name, err := parseName(parser)
	if err != nil {
		return ast.NewScalarTypeDefinition(), err
	}
	def := ast.NewScalarTypeDefinition()
	def.Name = name
	def.Loc = loc(parser, start)
	return def, nil
}

func parseEnumTypeDefinition(parser *Parser) (*ast.EnumTypeDefinition, error) {
	start := parser.Token.Start
	_, err := expectKeyWord(parser, "enum")
	if err != nil {
		return ast.NewEnumTypeDefinition(), err
	}
	name, err := parseName(parser)
	if err != nil {
		return ast.NewEnumTypeDefinition(), err
	}
	values, err := any(parser, lexer.TokenKind[lexer.BRACE_L], parseEnumValueDefinition, lexer.TokenKind[lexer.BRACE_R])

	def := ast.NewEnumTypeDefinition()
	def.Name = name
	def.Loc = loc(parser, start)
	def.Values = values
	return def, nil
}

func parseInputObjectTypeDefinition(parser *Parser) (*ast.InputObjectTypeDefinition, error) {
	start := parser.Token.Start
	_, err := expectKeyWord(parser, "input")
	if err != nil {
		return ast.NewInputObjectTypeDefinition(), err
	}
	name, err := parseName(parser)
	if err != nil {
		return ast.NewInputObjectTypeDefinition(), err
	}
	fields, err := any(parser, lexer.TokenKind[lexer.BRACE_L], parseInputValueDef, lexer.TokenKind[lexer.BRACE_R])

	def := ast.NewInputObjectTypeDefinition()
	def.Name = name
	def.Loc = loc(parser, start)
	def.Fields = fields
	return def, nil
}

func parseTypeExtensionDefinition(parser *Parser) (*ast.TypeExtensionDefinition, error) {
	start := parser.Token.Start
	_, err := expectKeyWord(parser, "extend")
	if err != nil {
		return ast.NewTypeExtensionDefinition(), err
	}

	definition, err := parseObjectTypeDefinition(parser)
	if err != nil {
		return ast.NewTypeExtensionDefinition(), err
	}

	def := ast.NewTypeExtensionDefinition()
	def.Loc = loc(parser, start)
	if definition != nil {
		def.Definition = *definition
	}
	return def, nil
}

func parseFieldDefinition(parser *Parser) (interface{}, error) {
	start := parser.Token.Start
	name, err := parseName(parser)
	if err != nil {
		return ast.FieldDefinition{}, err
	}
	args, err := parseArgumentDefs(parser)
	if err != nil {
		return ast.FieldDefinition{}, err
	}
	_, err = expect(parser, lexer.TokenKind[lexer.COLON])
	if err != nil {
		return ast.FieldDefinition{}, err
	}
	ttype, err := parseType(parser)
	if err != nil {
		return ast.FieldDefinition{}, err
	}
	return ast.FieldDefinition{
		Kind:      kinds.FieldDefinition,
		Name:      name,
		Arguments: args,
		Type:      ttype,
		Loc:       loc(parser, start),
	}, nil
}
func parseEnumValueDefinition(parser *Parser) (interface{}, error) {
	start := parser.Token.Start
	name, err := parseName(parser)
	if err != nil {
		return ast.EnumValueDefinition{}, err
	}
	return ast.EnumValueDefinition{
		Kind: kinds.EnumValueDefinition,
		Name: name,
		Loc:  loc(parser, start),
	}, nil
}

func parseUnionMembers(parser *Parser) ([]ast.NamedType, error) {
	members := []ast.NamedType{}
	for {
		member, err := parseNamedType(parser)
		if err != nil {
			return members, err
		}
		members = append(members, member)
		if !skip(parser, lexer.TokenKind[lexer.PIPE]) {
			break
		}
	}
	return members, nil
}

func parseArgumentDefs(parser *Parser) ([]ast.InputValueDefinition, error) {
	if !peek(parser, lexer.TokenKind[lexer.PAREN_L]) {
		return []ast.InputValueDefinition{}, nil
	}
	iInputValueDefinitions, err := many(parser, lexer.TokenKind[lexer.PAREN_L], parseInputValueDef, lexer.TokenKind[lexer.PAREN_R])
	inputValueDefinitions := []ast.InputValueDefinition{}
	if err != nil {
		return inputValueDefinitions, err
	}
	for _, iInputValueDefinition := range iInputValueDefinitions {
		inputValueDefinitions = append(inputValueDefinitions, iInputValueDefinition.(ast.InputValueDefinition))
	}

	return inputValueDefinitions, err
}

func parseInputValueDef(parser *Parser) (interface{}, error) {
	start := parser.Token.Start
	name, err := parseName(parser)
	if err != nil {
		return ast.InputValueDefinition{}, err
	}
	_, err = expect(parser, lexer.TokenKind[lexer.COLON])
	if err != nil {
		return ast.InputValueDefinition{}, err
	}
	ttype, err := parseType(parser)
	if err != nil {
		return ast.InputValueDefinition{}, err
	}
	var defaultValue ast.Value
	if skip(parser, lexer.TokenKind[lexer.EQUALS]) {
		defaultValue, err = parseConstValue(parser)
		if err != nil {
			return ast.InputValueDefinition{}, err
		}
	}
	return ast.InputValueDefinition{
		Kind:         kinds.InputValueDefinition,
		Name:         name,
		Type:         ttype,
		DefaultValue: defaultValue,
		Loc:          loc(parser, start),
	}, nil

}
