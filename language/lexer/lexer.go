package lexer

import (
	"fmt"

	"github.com/chris-ramon/graphql-go/language/source"
)

const (
	EOF = iota + 1
	BANG
	DOLLAR
	PAREN_L
	PAREN_R
	SPREAD
	COLON
	EQUALS
	AT
	BRACKET_L
	BRACKET_R
	BRACE_L
	PIPE
	BRACE_R
	NAME
	VARIABLE
	INT
	FLOAT
	STRING
)

var TokenKind map[int]int
var tokenDescription map[int]string

func init() {
	TokenKind = make(map[int]int)
	tokenDescription = make(map[int]string)
	TokenKind[EOF] = EOF
	TokenKind[BANG] = BANG
	TokenKind[DOLLAR] = DOLLAR
	TokenKind[PAREN_L] = PAREN_L
	TokenKind[PAREN_R] = PAREN_R
	TokenKind[SPREAD] = SPREAD
	TokenKind[COLON] = SPREAD
	TokenKind[EQUALS] = EQUALS
	TokenKind[AT] = AT
	TokenKind[BRACKET_L] = BRACKET_L
	TokenKind[BRACKET_R] = BRACKET_R
	TokenKind[BRACE_L] = BRACE_L
	TokenKind[PIPE] = PIPE
	TokenKind[BRACE_R] = BRACE_R
	TokenKind[NAME] = NAME
	TokenKind[VARIABLE] = VARIABLE
	TokenKind[INT] = INT
	TokenKind[FLOAT] = FLOAT
	TokenKind[STRING] = STRING
	tokenDescription[TokenKind[EOF]] = "EOF"
	tokenDescription[TokenKind[BANG]] = "!"
	tokenDescription[TokenKind[DOLLAR]] = "$"
	tokenDescription[TokenKind[PAREN_L]] = "("
	tokenDescription[TokenKind[PAREN_R]] = ")"
	tokenDescription[TokenKind[SPREAD]] = "..."
	tokenDescription[TokenKind[COLON]] = ":"
	tokenDescription[TokenKind[EQUALS]] = "="
	tokenDescription[TokenKind[AT]] = "@"
	tokenDescription[TokenKind[BRACKET_L]] = "["
	tokenDescription[TokenKind[BRACKET_R]] = "]"
	tokenDescription[TokenKind[BRACE_L]] = "}"
	tokenDescription[TokenKind[PIPE]] = "|"
	tokenDescription[TokenKind[BRACE_R]] = "{"
	tokenDescription[TokenKind[NAME]] = "Name"
	tokenDescription[TokenKind[VARIABLE]] = "Variable"
	tokenDescription[TokenKind[INT]] = "Int"
	tokenDescription[TokenKind[FLOAT]] = "Float"
	tokenDescription[TokenKind[STRING]] = "String"
}

type Token struct {
	Kind  int
	Start int
	End   int
	Value string
}

type Lexer func(resetPosition int) Token

func Lex(s *source.Source) Lexer {
	var prevPosition int
	return func(resetPosition int) Token {
		if resetPosition == 0 {
			resetPosition = prevPosition
		}
		var token = readToken(s, resetPosition)
		prevPosition = token.End
		return token
	}
}

// Reads an alphanumeric + underscore name from the source.
// [_A-Za-z][_0-9A-Za-z]*
func readName(source *source.Source, position int) Token {
	body := source.Body
	bodyLength := len(body)
	end := position + 1
	for {
		code := charCodeAt(body, end)
		if (end != bodyLength) && (code == 95 ||
			code >= 48 && code <= 57 ||
			code >= 65 && code <= 90 ||
			code >= 97 && code <= 122) {
			end += 1
			continue
		} else {
			break
		}
	}
	return makeToken(TokenKind[NAME], position, end, body[position:end])
}

func makeToken(kind int, start int, end int, value string) Token {
	return Token{Kind: kind, Start: start, End: end, Value: value}
}

func readToken(s *source.Source, fromPosition int) Token {
	body := s.Body
	bodyLength := len(body)
	position := positionAfterWhitespace(body, fromPosition)
	code := charCodeAt(body, position)
	if position >= bodyLength {
		return makeToken(TokenKind[EOF], position, position, "")
	}
	switch code {
	// !
	case 33:
		return makeToken(TokenKind[BANG], position, position+1, "")
	// $
	case 36:
		return makeToken(TokenKind[DOLLAR], position, position+1, "")
	// (
	case 40:
		return makeToken(TokenKind[PAREN_L], position, position+1, "")
	// )
	case 41:
		return makeToken(TokenKind[PAREN_R], position, position+1, "")
	// .
	case 46:
		if charCodeAt(body, position+1) == 46 && charCodeAt(body, position+2) == 46 {
			return makeToken(TokenKind[SPREAD], position, position+3, "")
		}
		return makeToken(TokenKind[PAREN_R], position, position+1, "")
	// a-z
	case 97, 98, 99, 100, 101, 102, 103, 104, 122:
		return readName(s, position)
	case 123:
		return makeToken(TokenKind[BRACE_L], position, position+1, "")
	}
	return Token{}
}

func charCodeAt(body string, position int) rune {
	return []rune(body)[position]
}

// Reads from body starting at startPosition until it finds a non-whitespace
// or commented character, then returns the position of that character for lexing.
// lexing.
func positionAfterWhitespace(body string, startPosition int) int {
	bodyLength := len(body)
	position := startPosition
	for {
		if position < bodyLength {
			code := charCodeAt(body, position)
			if code == 32 || // space
				code == 44 || // comma
				code == 160 || // '\xa0'
				code == 0x2028 || // line separator
				code == 0x2029 || // paragraph separator
				code > 8 && code < 14 { // whitespace
				position += 1
			} else if code == 35 { // #
				position += 1
				for {
					code := charCodeAt(body, position)
					if position < bodyLength &&
						code != 10 && code != 13 && code != 0x2028 && code != 0x2029 {
						position += 1
						continue
					} else {
						break
					}
				}
			} else {
				break
			}
			continue
		} else {
			break
		}
	}
	return position
}

func GetTokenDesc(token Token) string {
	if token.Value == "" {
		return GetTokenKindDesc(token.Kind)
	} else {
		return fmt.Sprintf("%s \"%s\"", GetTokenKindDesc(token.Kind), token.Value)
	}
}

func GetTokenKindDesc(kind int) string {
	return tokenDescription[kind]
}
