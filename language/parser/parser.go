package parser

import (
	"errors"
	"fmt"

	"github.com/chris-ramon/graphql-go/language/ast"
	"github.com/chris-ramon/graphql-go/language/errors"
	"github.com/chris-ramon/graphql-go/language/fd"
	"github.com/chris-ramon/graphql-go/language/kinds"
	"github.com/chris-ramon/graphql-go/language/lexer"
	"github.com/chris-ramon/graphql-go/language/od"
	"github.com/chris-ramon/graphql-go/language/source"
)

func unexpected(parser Parser, atToken lexer.Token) error {
	var token lexer.Token
	if (atToken == lexer.Token{}) {
		token = parser.Token
	} else {
		token = parser.Token
	}
	return languageerrors.Error(parser.Source, token.Start, lexer.GetTokenDesc(token))
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
	case source.Source:
		sourceObj = p.Source.(*source.Source)
	case string:
		s, _ := p.Source.(string)
		sourceObj = source.NewSource(s, "")
	default:
		return doc, errors.New("Unsupported source type")
	}
	parser := makeParser(sourceObj, p.Options)
	doc, err := parseDocument(parser)
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

func makeParser(s *source.Source, opts ParseOptions) (p Parser) {
	lexToken := lexer.Lex(s)
	return Parser{
		LexToken: lexToken,
		Source:   s,
		Options:  opts,
		PrevEnd:  0,
		Token:    lexToken(0),
	}
}

// Implements the parsing rules in the Document section.
func parseDocument(parser Parser) (ast.Document, error) {
	start := parser.Token.Start
	var iDefinitions []interface{}
	for {
		if skip(parser, lexer.TokenKind[lexer.EOF]) {
			break
		}
		if peek(parser, lexer.TokenKind[lexer.BRACE_L]) {
			oDef, err := parseOperationDefinition(parser)
			if err != nil {
				return ast.Document{}, err
			}
			iDefinitions = append(iDefinitions, oDef)
		} else if peek(parser, lexer.TokenKind[lexer.NAME]) {
			if parser.Token.Value == "query" || parser.Token.Value == "mutation" {
				oDef, err := parseOperationDefinition(parser)
				if err != nil {
					return ast.Document{}, err
				}
				iDefinitions = append(iDefinitions, oDef)
			} else if parser.Token.Value == "fragment" {
				fDef, err := parseFragmentDefinition(parser)
				if err != nil {
					return ast.Document{}, err
				}
				iDefinitions = append(iDefinitions, fDef)
			} else {
				if err := unexpected(parser, lexer.Token{}); err != nil {
					return ast.Document{}, err
				}
			}
		}
	}
	var definitions []ast.Definition
	for i, iDef := range iDefinitions {
		definitions[i] = iDef.(ast.Definition)
	}
	return ast.Document{
		Kind:        kinds.Document,
		Definitions: definitions,
		Loc:         loc(parser, start),
	}, nil
}

// If the next token is of the given kind, return true after advancing
// the parser. Otherwise, do not change the parser state and return false.
func skip(parser Parser, Kind int) bool {
	if parser.Token.Kind == Kind {
		advance(parser)
		return true
	}
	return false
}

// Moves the internal parser object to the next lexed token.
func advance(parser Parser) {
	prevEnd := parser.Token.End
	parser.PrevEnd = prevEnd
	parser.Token = parser.LexToken(prevEnd)
}

// Determines if the next token is of a given kind
func peek(parser Parser, Kind int) bool {
	return parser.Token.Kind == Kind
}

// Implements the parsing rules in the Operations section.
func parseOperationDefinition(parser Parser) (*od.OperationDefinition, error) {
	start := parser.Token.Start
	selectionSet, err := parseSelectionSet(parser)
	if err != nil {
		oDef := od.NewOperationDefinition()
		return oDef, err
	}
	if peek(parser, lexer.TokenKind[lexer.BRACE_L]) {
		oDef := od.NewOperationDefinition()
		oDef.Operation = "query"
		oDef.SelectionSet = selectionSet
		oDef.Loc = loc(parser, start)
		return oDef, err
	}
	operationToken, err := expect(parser, lexer.TokenKind[lexer.NAME])
	if err != nil {
		oDef := od.NewOperationDefinition()
		return oDef, err
	}
	operation := operationToken.Value
	name, err := parseName(parser)
	if err != nil {
		oDef := od.NewOperationDefinition()
		return oDef, err
	}
	variableDefinitions, err := parseVariableDefinitions(parser)
	if err != nil {
		oDef := od.NewOperationDefinition()
		return oDef, err
	}
	directives, err := parseDirectives(parser)
	if err != nil {
		oDef := od.NewOperationDefinition()
		return oDef, err
	}
	oDef := od.NewOperationDefinition()
	oDef.Operation = operation
	oDef.Name = name
	oDef.VariableDefinitions = variableDefinitions
	oDef.Directives = directives
	oDef.SelectionSet = selectionSet
	oDef.Loc = loc(parser, start)
	return oDef, nil
}

func parseFragmentDefinition(parser Parser) (*fd.FragmentDefinition, error) {
	start := parser.Token.Start
	_, errFragment := expectKeyWord(parser, "fragment")
	if errFragment != nil {
		fDef := fd.NewFragmentDefinition()
		return fDef, errFragment
	}
	name, errName := parseName(parser)
	if errName != nil {
		fDef := fd.NewFragmentDefinition()
		return fDef, errName
	}
	_, errOn := expectKeyWord(parser, "on")
	if errOn != nil {
		fDef := fd.NewFragmentDefinition()
		return fDef, errOn
	}
	typeCondition, errTypeCondition := parseName(parser)
	if errTypeCondition != nil {
		fDef := fd.NewFragmentDefinition()
		return fDef, errTypeCondition
	}
	selectionSet, errSelectionSet := parseSelectionSet(parser)
	if errSelectionSet != nil {
		fDef := fd.NewFragmentDefinition()
		return fDef, errSelectionSet
	}
	directives, errDirectives := parseDirectives(parser)
	if errDirectives != nil {
		fDef := fd.NewFragmentDefinition()
		return fDef, errDirectives
	}
	fDef := fd.NewFragmentDefinition()
	fDef.Name = name
	fDef.TypeCondition = typeCondition
	fDef.Directives = directives
	fDef.SelectionSet = selectionSet
	fDef.Loc = loc(parser, start)
	return fDef, nil
}

func expectKeyWord(parser Parser, value string) (lexer.Token, error) {
	token := parser.Token
	if token.Kind == lexer.TokenKind[lexer.NAME] && token.Value == value {
		advance(parser)
		return token, nil
	}
	descp := fmt.Sprintf("Expected \"%s\", found %s", value, lexer.GetTokenDesc(token))
	return token, languageerrors.Error(parser.Source, token.Start, descp)
}

func parseSelectionSet(parser Parser) (ast.SelectionSet, error) {
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

func parseSelection(parser Parser) (interface{}, error) {
	if peek(parser, lexer.TokenKind[lexer.SPREAD]) {
		r, err := parseFragment(parser)
		if err != nil {
			return r, err
		}
		return r, nil
	} else {
		return parseField(parser)
	}
}

func loc(parser Parser, start int) ast.Location {
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

func expect(parser Parser, kind int) (lexer.Token, error) {
	token := parser.Token
	if token.Kind == kind {
		advance(parser)
		return token, nil
	}
	descp := fmt.Sprintf("Expected %s, found %s", lexer.GetTokenKindDesc(kind), lexer.GetTokenDesc(token))
	return token, languageerrors.Error(parser.Source, token.Start, descp)
}

// Converts a name lex token into a name parse node.
func parseName(parser Parser) (ast.Name, error) {
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

func parseVariableDefinitions(parser Parser) ([]ast.VariableDefinition, error) {
	if peek(parser, lexer.TokenKind[lexer.PAREN_L]) {
		vdefs, err := many(parser, lexer.TokenKind[lexer.PAREN_L], parseVariableDefinition, lexer.TokenKind[lexer.PAREN_R])
		var variableDefinitions []ast.VariableDefinition
		for i, vdef := range vdefs {
			variableDefinitions[i] = vdef.(ast.VariableDefinition)
		}
		if err != nil {
			return variableDefinitions, err
		}
		return variableDefinitions, nil
	} else {
		var vd []ast.VariableDefinition
		return vd, nil
	}
}

func parseDirectives(parser Parser) ([]ast.Directive, error) {
	var directives []ast.Directive
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

func parseDirective(parser Parser) (ast.Directive, error) {
	start := parser.Token.Start
	_, err := expect(parser, lexer.TokenKind[lexer.AT])
	if err != nil {
		return ast.Directive{}, err
	}
	name, err := parseName(parser)
	if err != nil {
		return ast.Directive{}, err
	}
	var value ast.Value
	if skip(parser, lexer.TokenKind[lexer.COLON]) {
		v, err := parseValue(parser, false)
		if err != nil {
			return ast.Directive{}, err
		}
		value = v
	}
	return ast.Directive{
		Kind:  kinds.Directive,
		Name:  name,
		Value: value,
		Loc:   loc(parser, start),
	}, nil
}

func parseVariableDefinition(parser Parser) (interface{}, error) {
	start := parser.Token.Start
	var defaultValue ast.Value
	if skip(parser, lexer.TokenKind[lexer.EQUALS]) {
		dv, err := parseValue(parser, true)
		if err != nil {
			return dv, err
		}
		defaultValue = dv
	}
	_, err := expect(parser, lexer.TokenKind[lexer.COLON])
	if err != nil {
		return ast.VariableDefinition{}, err
	}
	variable, err := parseVariable(parser)
	if err != nil {
		return ast.VariableDefinition{}, err
	}
	ttype, err := parseType(parser)
	if err != nil {
		return ast.VariableDefinition{}, err
	}
	return ast.VariableDefinition{
		Kind:         kinds.VariableDefinition,
		Variable:     variable,
		Type:         ttype,
		DefaultValue: defaultValue,
		Loc:          loc(parser, start),
	}, nil
}

func parseVariable(parser Parser) (ast.Variable, error) {
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

func parseType(parser Parser) (ast.Type, error) {
	start := parser.Token.Start
	var ttype ast.Type
	if skip(parser, lexer.TokenKind[lexer.BRACE_L]) {
		t, errParseType := parseType(parser)
		if errParseType != nil {
			return t, errParseType
		}
		ttype = t
		_, errExpect := expect(parser, lexer.TokenKind[lexer.BRACKET_R])
		if errExpect != nil {
			return ttype, errExpect
		}
		ttype = ast.ListType{
			Kind: kinds.ListType,
			Type: ttype,
			Loc:  loc(parser, start),
		}
	} else {
		name, err := parseName(parser)
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

func parseValue(parser Parser, isConst bool) (ast.Value, error) {
	token := parser.Token
	switch token.Kind {
	case lexer.TokenKind[lexer.BRACE_L]:
		value, err := parseArray(parser, isConst)
		if err != nil {
			return value, err
		}
		return value, nil
	}
	if err := unexpected(parser, lexer.Token{}); err != nil {
		return nil, err
	}
	return nil, nil
}

type parseFn func(parser Parser) (interface{}, error)

func many(parser Parser, openKind int, parseFn parseFn, closeKind int) ([]interface{}, error) {
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

func parseFragment(parser Parser) (interface{}, error) {
	start := parser.Token.Start
	_, err := expect(parser, lexer.TokenKind[lexer.SPREAD])
	if err != nil {
		return nil, err
	}
	name, err := parseName(parser)
	if err != nil {
		return nil, err
	}
	if parser.Token.Value == "on" {
		advance(parser)
		selectionSet, err := parseSelectionSet(parser)
		if err != nil {
			return ast.InlineFragment{}, err
		}
		directives, err := parseDirectives(parser)
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
	directives, err := parseDirectives(parser)
	if err != nil {
		return ast.InlineFragment{}, err
	}
	return ast.FragmentSpread{
		Kind:       kinds.FragmentSpread,
		Name:       name,
		Directives: directives,
		Loc:        loc(parser, start),
	}, nil
}

func parseField(parser Parser) (ast.Field, error) {
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
	var selectionSet ast.SelectionSet
	if peek(parser, lexer.TokenKind[lexer.BRACE_L]) {
		sSet, err := parseSelectionSet(parser)
		if err != nil {
			return ast.Field{}, err
		}
		selectionSet = sSet
	}
	arguments, err := parseArguments(parser)
	if err != nil {
		return ast.Field{}, err
	}
	directives, err := parseDirectives(parser)
	if err != nil {
		return ast.Field{}, err
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

func parseArray(parser Parser, isConst bool) (ast.ArrayValue, error) {
	start := parser.Token.Start
	var item parseFn
	if isConst {
		item = parseConstValue
	} else {
		item = parseVariableValue
	}
	iValues, err := any(parser, lexer.TokenKind[lexer.BRACE_L], item, lexer.TokenKind[lexer.BRACKET_R])
	if err != nil {
		return ast.ArrayValue{}, err
	}
	var values []ast.Value
	for i, iValue := range iValues {
		values[i] = iValue.(ast.Value)
	}
	return ast.ArrayValue{
		Kind:   kinds.Array,
		Values: values,
		Loc:    loc(parser, start),
	}, nil
}

func any(parser Parser, openKind int, parseFn parseFn, closeKind int) ([]interface{}, error) {
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

func parseArguments(parser Parser) ([]ast.Argument, error) {
	if peek(parser, lexer.TokenKind[lexer.PAREN_L]) {
		iArguments, err := many(parser, lexer.TokenKind[lexer.PAREN_L], parseArgument, lexer.TokenKind[lexer.PAREN_R])
		var arguments []ast.Argument
		if err != nil {
			return arguments, err
		}
		for i, iArgument := range iArguments {
			arguments[i] = iArgument.(ast.Argument)
		}
		return arguments, nil
	} else {
		return []ast.Argument{}, nil
	}
}

func parseArgument(parser Parser) (interface{}, error) {
	start := parser.Token.Start
	name, err := parseName(parser)
	if err != nil {
		return ast.Argument{}, err
	}
	_, errExpect := expect(parser, lexer.TokenKind[lexer.COLON])
	if errExpect != nil {
		return ast.Argument{}, errExpect
	}
	value, err := parseValue(parser, false)
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

func parseConstValue(parser Parser) (interface{}, error) {
	value, err := parseValue(parser, true)
	if err != nil {
		return value, err
	}
	return value, nil
}

func parseVariableValue(parser Parser) (interface{}, error) {
	value, err := parseValue(parser, false)
	if err != nil {
		return value, err
	}
	return value, nil
}
