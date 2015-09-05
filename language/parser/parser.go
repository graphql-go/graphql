package parser

import (
	"fmt"

	"github.com/chris-ramon/graphql-go/errors"
	"github.com/chris-ramon/graphql-go/language/ast"
	"github.com/chris-ramon/graphql-go/language/errors"
	"github.com/chris-ramon/graphql-go/language/fd"
	"github.com/chris-ramon/graphql-go/language/kinds"
	"github.com/chris-ramon/graphql-go/language/lexer"
	"github.com/chris-ramon/graphql-go/language/od"
	"github.com/chris-ramon/graphql-go/language/source"
)

func unexpected(parser *Parser, atToken lexer.Token) graphqlerrors.GraphQLError {
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

func Parse(p ParseParams) (ast.Document, graphqlerrors.GraphQLError) {
	var doc ast.Document
	var sourceObj *source.Source
	switch p.Source.(type) {
	case source.Source:
		sourceObj = p.Source.(*source.Source)
	default:
		s, _ := p.Source.(string)
		sourceObj = source.NewSource(s, "")
	}
	parser, errMakeParser := makeParser(sourceObj, p.Options)
	if errMakeParser.Error != nil {
		return doc, errMakeParser
	}
	doc, errParseDocument := parseDocument(parser)
	if errParseDocument.Error != nil {
		return doc, errParseDocument
	}
	return doc, graphqlerrors.GraphQLError{}
}

type Parser struct {
	LexToken lexer.Lexer
	Source   *source.Source
	Options  ParseOptions
	PrevEnd  int
	Token    lexer.Token
}

func makeParser(s *source.Source, opts ParseOptions) (*Parser, graphqlerrors.GraphQLError) {
	lexToken := lexer.Lex(s)
	token, err := lexToken(0)
	if err.Error != nil {
		return &Parser{}, err
	}
	return &Parser{
		LexToken: lexToken,
		Source:   s,
		Options:  opts,
		PrevEnd:  0,
		Token:    token,
	}, graphqlerrors.GraphQLError{}
}

// Implements the parsing rules in the Document section.
func parseDocument(parser *Parser) (ast.Document, graphqlerrors.GraphQLError) {
	start := parser.Token.Start
	var definitions []ast.Definition
	for {
		if skip(parser, lexer.TokenKind[lexer.EOF]) {
			break
		}
		if peek(parser, lexer.TokenKind[lexer.BRACE_L]) {
			oDef, err := parseOperationDefinition(parser)
			if err.Error != nil {
				return ast.Document{}, err
			}
			definitions = append(definitions, oDef)
		} else if peek(parser, lexer.TokenKind[lexer.NAME]) {
			if parser.Token.Value == "query" || parser.Token.Value == "mutation" {
				oDef, err := parseOperationDefinition(parser)
				if err.Error != nil {
					return ast.Document{}, err
				}
				definitions = append(definitions, oDef)
			} else if parser.Token.Value == "fragment" {
				fDef, err := parseFragmentDefinition(parser)
				if err.Error != nil {
					return ast.Document{}, err
				}
				definitions = append(definitions, fDef)
			} else {
				if err := unexpected(parser, lexer.Token{}); err.Error != nil {
					return ast.Document{}, err
				}
			}
		}
	}
	return ast.Document{
		Kind:        kinds.Document,
		Loc:         loc(parser, start),
		Definitions: definitions,
	}, graphqlerrors.GraphQLError{}
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
func advance(parser *Parser) graphqlerrors.GraphQLError {
	prevEnd := parser.Token.End
	parser.PrevEnd = prevEnd
	token, err := parser.LexToken(prevEnd)
	if err.Error != nil {
		return err
	}
	parser.Token = token
	return graphqlerrors.GraphQLError{}
}

// Determines if the next token is of a given kind
func peek(parser *Parser, Kind int) bool {
	return parser.Token.Kind == Kind
}

// Implements the parsing rules in the Operations section.
func parseOperationDefinition(parser *Parser) (*od.OperationDefinition, graphqlerrors.GraphQLError) {
	start := parser.Token.Start
	if peek(parser, lexer.TokenKind[lexer.BRACE_L]) {
		selectionSet, err := parseSelectionSet(parser)
		if err.Error != nil {
			oDef := od.NewOperationDefinition()
			return oDef, err
		}
		oDef := od.NewOperationDefinition()
		oDef.Operation = "query"
		oDef.Directives = []ast.Directive{}
		oDef.SelectionSet = selectionSet
		oDef.Loc = loc(parser, start)
		return oDef, err
	}
	operationToken, err := expect(parser, lexer.TokenKind[lexer.NAME])
	if err.Error != nil {
		oDef := od.NewOperationDefinition()
		return oDef, err
	}
	operation := operationToken.Value
	name, err := parseName(parser)
	if err.Error != nil {
		oDef := od.NewOperationDefinition()
		return oDef, err
	}
	variableDefinitions, err := parseVariableDefinitions(parser)
	if err.Error != nil {
		oDef := od.NewOperationDefinition()
		return oDef, err
	}
	directives, err := parseDirectives(parser)
	if err.Error != nil {
		oDef := od.NewOperationDefinition()
		return oDef, err
	}
	selectionSet, err := parseSelectionSet(parser)
	if err.Error != nil {
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
	return oDef, graphqlerrors.GraphQLError{}
}

func parseFragmentDefinition(parser *Parser) (*fd.FragmentDefinition, graphqlerrors.GraphQLError) {
	start := parser.Token.Start
	_, errFragment := expectKeyWord(parser, "fragment")
	if errFragment.Error != nil {
		fDef := fd.NewFragmentDefinition()
		return fDef, errFragment
	}
	name, errName := parseName(parser)
	if errName.Error != nil {
		fDef := fd.NewFragmentDefinition()
		return fDef, errName
	}
	_, errOn := expectKeyWord(parser, "on")
	if errOn.Error != nil {
		fDef := fd.NewFragmentDefinition()
		return fDef, errOn
	}
	typeCondition, errTypeCondition := parseName(parser)
	if errTypeCondition.Error != nil {
		fDef := fd.NewFragmentDefinition()
		return fDef, errTypeCondition
	}
	selectionSet, errSelectionSet := parseSelectionSet(parser)
	if errSelectionSet.Error != nil {
		fDef := fd.NewFragmentDefinition()
		return fDef, errSelectionSet
	}
	directives, errDirectives := parseDirectives(parser)
	if errDirectives.Error != nil {
		fDef := fd.NewFragmentDefinition()
		return fDef, errDirectives
	}
	fDef := fd.NewFragmentDefinition()
	fDef.Name = name
	fDef.TypeCondition = typeCondition
	fDef.Directives = directives
	fDef.SelectionSet = selectionSet
	fDef.Loc = loc(parser, start)
	return fDef, graphqlerrors.GraphQLError{}
}

func expectKeyWord(parser *Parser, value string) (lexer.Token, graphqlerrors.GraphQLError) {
	token := parser.Token
	if token.Kind == lexer.TokenKind[lexer.NAME] && token.Value == value {
		advance(parser)
		return token, graphqlerrors.GraphQLError{}
	}
	descp := fmt.Sprintf("Expected \"%s\", found %s", value, lexer.GetTokenDesc(token))
	return token, languageerrors.Error(parser.Source, token.Start, descp)
}

func parseSelectionSet(parser *Parser) (ast.SelectionSet, graphqlerrors.GraphQLError) {
	start := parser.Token.Start
	selections, err := many(parser, lexer.TokenKind[lexer.BRACE_L], parseSelection, lexer.TokenKind[lexer.BRACE_R])
	if err.Error != nil {
		return ast.SelectionSet{}, err
	}
	return ast.SelectionSet{
		Kind:       kinds.SelectionSet,
		Selections: selections,
		Loc:        loc(parser, start),
	}, graphqlerrors.GraphQLError{}
}

func parseSelection(parser *Parser) (interface{}, graphqlerrors.GraphQLError) {
	if peek(parser, lexer.TokenKind[lexer.SPREAD]) {
		r, err := parseFragment(parser)
		if err.Error != nil {
			return r, err
		}
		return r, graphqlerrors.GraphQLError{}
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

func expect(parser *Parser, kind int) (lexer.Token, graphqlerrors.GraphQLError) {
	token := parser.Token
	if token.Kind == kind {
		advance(parser)
		return token, graphqlerrors.GraphQLError{}
	}
	descp := fmt.Sprintf("Expected %s, found %s", lexer.GetTokenKindDesc(kind), lexer.GetTokenDesc(token))
	return token, languageerrors.Error(parser.Source, token.Start, descp)
}

// Converts a name lex token into a name parse node.
func parseName(parser *Parser) (ast.Name, graphqlerrors.GraphQLError) {
	token, err := expect(parser, lexer.TokenKind[lexer.NAME])
	if err.Error != nil {
		return ast.Name{}, err
	}
	return ast.Name{
		Kind:  kinds.Name,
		Value: token.Value,
		Loc:   loc(parser, token.Start),
	}, graphqlerrors.GraphQLError{}
}

func parseVariableDefinitions(parser *Parser) ([]ast.VariableDefinition, graphqlerrors.GraphQLError) {
	if peek(parser, lexer.TokenKind[lexer.PAREN_L]) {
		vdefs, err := many(parser, lexer.TokenKind[lexer.PAREN_L], parseVariableDefinition, lexer.TokenKind[lexer.PAREN_R])
		var variableDefinitions []ast.VariableDefinition
		for i, vdef := range vdefs {
			variableDefinitions[i] = vdef.(ast.VariableDefinition)
		}
		if err.Error != nil {
			return variableDefinitions, err
		}
		return variableDefinitions, graphqlerrors.GraphQLError{}
	} else {
		var vd []ast.VariableDefinition
		return vd, graphqlerrors.GraphQLError{}
	}
}

func parseDirectives(parser *Parser) ([]ast.Directive, graphqlerrors.GraphQLError) {
	directives := []ast.Directive{}
	for {
		if !peek(parser, lexer.TokenKind[lexer.AT]) {
			break
		}
		directive, err := parseDirective(parser)
		if err.Error != nil {
			return directives, err
		}
		directives = append(directives, directive)
	}
	return directives, graphqlerrors.GraphQLError{}
}

func parseDirective(parser *Parser) (ast.Directive, graphqlerrors.GraphQLError) {
	start := parser.Token.Start
	_, err := expect(parser, lexer.TokenKind[lexer.AT])
	if err.Error != nil {
		return ast.Directive{}, err
	}
	name, err := parseName(parser)
	if err.Error != nil {
		return ast.Directive{}, err
	}
	var value ast.Value
	if skip(parser, lexer.TokenKind[lexer.COLON]) {
		v, err := parseValue(parser, false)
		if err.Error != nil {
			return ast.Directive{}, err
		}
		value = v
	}
	return ast.Directive{
		Kind:  kinds.Directive,
		Name:  name,
		Value: value,
		Loc:   loc(parser, start),
	}, graphqlerrors.GraphQLError{}
}

func parseVariableDefinition(parser *Parser) (interface{}, graphqlerrors.GraphQLError) {
	start := parser.Token.Start
	var defaultValue ast.Value
	if skip(parser, lexer.TokenKind[lexer.EQUALS]) {
		dv, err := parseValue(parser, true)
		if err.Error != nil {
			return dv, err
		}
		defaultValue = dv
	}
	_, err := expect(parser, lexer.TokenKind[lexer.COLON])
	if err.Error != nil {
		return ast.VariableDefinition{}, err
	}
	variable, err := parseVariable(parser)
	if err.Error != nil {
		return ast.VariableDefinition{}, err
	}
	ttype, err := parseType(parser)
	if err.Error != nil {
		return ast.VariableDefinition{}, err
	}
	return ast.VariableDefinition{
		Kind:         kinds.VariableDefinition,
		Variable:     variable,
		Type:         ttype,
		DefaultValue: defaultValue,
		Loc:          loc(parser, start),
	}, graphqlerrors.GraphQLError{}
}

func parseVariable(parser *Parser) (ast.Variable, graphqlerrors.GraphQLError) {
	start := parser.Token.Start
	_, err := expect(parser, lexer.TokenKind[lexer.DOLLAR])
	if err.Error != nil {
		return ast.Variable{}, err
	}
	name, err := parseName(parser)
	if err.Error != nil {
		return ast.Variable{}, err
	}
	return ast.Variable{
		Kind: kinds.Variable,
		Name: name,
		Loc:  loc(parser, start),
	}, graphqlerrors.GraphQLError{}
}

func parseType(parser *Parser) (ast.Type, graphqlerrors.GraphQLError) {
	start := parser.Token.Start
	var ttype ast.Type
	if skip(parser, lexer.TokenKind[lexer.BRACE_L]) {
		t, errParseType := parseType(parser)
		if errParseType.Error != nil {
			return t, errParseType
		}
		ttype = t
		_, errExpect := expect(parser, lexer.TokenKind[lexer.BRACKET_R])
		if errExpect.Error != nil {
			return ttype, errExpect
		}
		ttype = ast.ListType{
			Kind: kinds.ListType,
			Type: ttype,
			Loc:  loc(parser, start),
		}
	} else {
		name, err := parseName(parser)
		if err.Error != nil {
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
		return ttype, graphqlerrors.GraphQLError{}
	}
	return ttype, graphqlerrors.GraphQLError{}
}

func parseValue(parser *Parser, isConst bool) (ast.Value, graphqlerrors.GraphQLError) {
	token := parser.Token
	switch token.Kind {
	case lexer.TokenKind[lexer.BRACE_L]:
		value, err := parseArray(parser, isConst)
		if err.Error != nil {
			return value, err
		}
		return value, graphqlerrors.GraphQLError{}
	}
	if err := unexpected(parser, lexer.Token{}); err.Error != nil {
		return nil, err
	}
	return nil, graphqlerrors.GraphQLError{}
}

type parseFn func(parser *Parser) (interface{}, graphqlerrors.GraphQLError)

func many(parser *Parser, openKind int, parseFn parseFn, closeKind int) ([]interface{}, graphqlerrors.GraphQLError) {
	_, err := expect(parser, openKind)
	if err.Error != nil {
		return nil, err
	}
	var nodes []interface{}
	node, err := parseFn(parser)
	if err.Error != nil {
		return nodes, err
	}
	nodes = append(nodes, node)
	for {
		if skip(parser, closeKind) {
			break
		}
		node, err := parseFn(parser)
		if err.Error != nil {
			return nodes, err
		}
		nodes = append(nodes, node)
	}
	return nodes, graphqlerrors.GraphQLError{}
}

func parseFragment(parser *Parser) (interface{}, graphqlerrors.GraphQLError) {
	start := parser.Token.Start
	_, err := expect(parser, lexer.TokenKind[lexer.SPREAD])
	if err.Error != nil {
		return nil, err
	}
	name, err := parseName(parser)
	if err.Error != nil {
		return nil, err
	}
	if parser.Token.Value == "on" {
		advance(parser)
		selectionSet, err := parseSelectionSet(parser)
		if err.Error != nil {
			return ast.InlineFragment{}, err
		}
		directives, err := parseDirectives(parser)
		if err.Error != nil {
			return ast.InlineFragment{}, err
		}
		return ast.InlineFragment{
			Kind:          kinds.InlineFragment,
			TypeCondition: name,
			Directives:    directives,
			SelectionSet:  selectionSet,
			Loc:           loc(parser, start),
		}, graphqlerrors.GraphQLError{}
	}
	directives, err := parseDirectives(parser)
	if err.Error != nil {
		return ast.InlineFragment{}, err
	}
	return ast.FragmentSpread{
		Kind:       kinds.FragmentSpread,
		Name:       name,
		Directives: directives,
		Loc:        loc(parser, start),
	}, graphqlerrors.GraphQLError{}
}

func parseField(parser *Parser) (ast.Field, graphqlerrors.GraphQLError) {
	start := parser.Token.Start
	nameOrAlias, err := parseName(parser)
	if err.Error != nil {
		return ast.Field{}, err
	}
	var (
		name  ast.Name
		alias ast.Name
	)
	if skip(parser, lexer.TokenKind[lexer.COLON]) {
		alias = nameOrAlias
		n, err := parseName(parser)
		if err.Error != nil {
			return ast.Field{}, err
		}
		name = n
	} else {
		name = nameOrAlias
	}
	var selectionSet ast.SelectionSet
	if peek(parser, lexer.TokenKind[lexer.BRACE_L]) {
		sSet, err := parseSelectionSet(parser)
		if err.Error != nil {
			return ast.Field{}, err
		}
		selectionSet = sSet
	}
	arguments, err := parseArguments(parser)
	if err.Error != nil {
		return ast.Field{}, err
	}
	directives, err := parseDirectives(parser)
	if err.Error != nil {
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
	}, graphqlerrors.GraphQLError{}
}

func parseArray(parser *Parser, isConst bool) (ast.ArrayValue, graphqlerrors.GraphQLError) {
	start := parser.Token.Start
	var item parseFn
	if isConst {
		item = parseConstValue
	} else {
		item = parseVariableValue
	}
	iValues, err := any(parser, lexer.TokenKind[lexer.BRACE_L], item, lexer.TokenKind[lexer.BRACKET_R])
	if err.Error != nil {
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
	}, graphqlerrors.GraphQLError{}
}

func any(parser *Parser, openKind int, parseFn parseFn, closeKind int) ([]interface{}, graphqlerrors.GraphQLError) {
	var nodes []interface{}
	_, err := expect(parser, openKind)
	if err.Error != nil {
		return nodes, graphqlerrors.GraphQLError{}
	}
	for {
		if skip(parser, closeKind) {
			break
		}
		n, err := parseFn(parser)
		if err.Error != nil {
			return nodes, err
		}
		nodes = append(nodes, n)
	}
	return nodes, graphqlerrors.GraphQLError{}
}

func parseArguments(parser *Parser) ([]ast.Argument, graphqlerrors.GraphQLError) {
	if peek(parser, lexer.TokenKind[lexer.PAREN_L]) {
		iArguments, err := many(parser, lexer.TokenKind[lexer.PAREN_L], parseArgument, lexer.TokenKind[lexer.PAREN_R])
		var arguments []ast.Argument
		if err.Error != nil {
			return arguments, err
		}
		for i, iArgument := range iArguments {
			arguments[i] = iArgument.(ast.Argument)
		}
		return arguments, graphqlerrors.GraphQLError{}
	} else {
		return []ast.Argument{}, graphqlerrors.GraphQLError{}
	}
}

func parseArgument(parser *Parser) (interface{}, graphqlerrors.GraphQLError) {
	start := parser.Token.Start
	name, err := parseName(parser)
	if err.Error != nil {
		return ast.Argument{}, err
	}
	_, errExpect := expect(parser, lexer.TokenKind[lexer.COLON])
	if errExpect.Error != nil {
		return ast.Argument{}, errExpect
	}
	value, err := parseValue(parser, false)
	if err.Error != nil {
		return ast.Argument{}, err
	}
	return ast.Argument{
		Kind:  kinds.Argument,
		Name:  name,
		Value: value,
		Loc:   loc(parser, start),
	}, graphqlerrors.GraphQLError{}
}

func parseConstValue(parser *Parser) (interface{}, graphqlerrors.GraphQLError) {
	value, err := parseValue(parser, true)
	if err.Error != nil {
		return value, err
	}
	return value, graphqlerrors.GraphQLError{}
}

func parseVariableValue(parser *Parser) (interface{}, graphqlerrors.GraphQLError) {
	value, err := parseValue(parser, false)
	if err.Error != nil {
		return value, err
	}
	return value, graphqlerrors.GraphQLError{}
}
