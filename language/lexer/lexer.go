package lexer

import (
	"fmt"

	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/source"
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
	TokenKind[COLON] = COLON
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
	tokenDescription[TokenKind[BRACE_L]] = "{"
	tokenDescription[TokenKind[PIPE]] = "|"
	tokenDescription[TokenKind[BRACE_R]] = "}"
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

func (t *Token) String() string {
	return fmt.Sprintf("%s", tokenDescription[t.Kind])
}

type Lexer func(resetPosition int) (Token, error)

func Lex(s *source.Source) Lexer {
	var prevPosition int
	return func(resetPosition int) (Token, error) {
		if resetPosition == 0 {
			resetPosition = prevPosition
		}
		token, err := readToken(s, resetPosition)
		if err != nil {
			return token, err
		}
		prevPosition = token.End
		return token, nil
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

// Reads a number token from the source file, either a float
// or an int depending on whether a decimal point appears.
// Int:   -?(0|[1-9][0-9]*)
// Float: -?(0|[1-9][0-9]*)(\.[0-9]+)?((E|e)(+|-)?[0-9]+)?
func readNumber(s *source.Source, start int, firstCode rune) (Token, error) {
	code := firstCode
	body := s.Body
	position := start
	isFloat := false
	if code == 45 { // -
		position += 1
		code = charCodeAt(body, position)
	}
	if code == 48 { // 0
		position += 1
		code = charCodeAt(body, position)
		if code >= 48 && code <= 57 {
			description := fmt.Sprintf("Invalid number, unexpected digit after 0: \"%c\".", code)
			return Token{}, gqlerrors.NewSyntaxError(s, position, description)
		}
	} else {
		p, err := readDigits(s, position, code)
		if err != nil {
			return Token{}, err
		}
		position = p
		code = charCodeAt(body, position)
	}
	if code == 46 { // .
		isFloat = true
		position += 1
		code = charCodeAt(body, position)
		p, err := readDigits(s, position, code)
		if err != nil {
			return Token{}, err
		}
		position = p
		code = charCodeAt(body, position)
	}
	if code == 69 || code == 101 { // E e
		isFloat = true
		position += 1
		code = charCodeAt(body, position)
		if code == 43 || code == 45 { // + -
			position += 1
			code = charCodeAt(body, position)
		}
		p, err := readDigits(s, position, code)
		if err != nil {
			return Token{}, err
		}
		position = p
	}
	kind := TokenKind[INT]
	if isFloat {
		kind = TokenKind[FLOAT]
	}
	return makeToken(kind, start, position, body[start:position]), nil
}

// Returns the new position in the source after reading digits.
func readDigits(s *source.Source, start int, firstCode rune) (int, error) {
	body := s.Body
	position := start
	code := firstCode
	if code >= 48 && code <= 57 { // 0 - 9
		for {
			if code >= 48 && code <= 57 { // 0 - 9
				position += 1
				code = charCodeAt(body, position)
				continue
			} else {
				break
			}
		}
		return position, nil
	}
	var description string
	if code != 0 {
		description = fmt.Sprintf("Invalid number, expected digit but got: \"%c\".", code)
	} else {
		description = fmt.Sprintf("Invalid number, expected digit but got: EOF.")
	}
	return position, gqlerrors.NewSyntaxError(s, position, description)
}

func readString(s *source.Source, start int) (Token, error) {
	body := s.Body
	position := start + 1
	chunkStart := position
	var code rune
	var value string
	for {
		code = charCodeAt(body, position)
		if position < len(body) && code != 34 && code != 10 && code != 13 && code != 0x2028 && code != 0x2029 {
			position += 1
			if code == 92 { // \
				value += body[chunkStart : position-1]
				code = charCodeAt(body, position)
				switch code {
				case 34:
					value += "\""
					break
				case 47:
					value += "\\/"
					break
				case 92:
					value += "\\"
					break
				case 98:
					value += "\b"
					break
				case 102:
					value += "\f"
					break
				case 110:
					value += "\n"
					break
				case 114:
					value += "\r"
					break
				case 116:
					value += "\t"
					break
				case 117:
					charCode := uniCharCode(
						charCodeAt(body, position+1),
						charCodeAt(body, position+2),
						charCodeAt(body, position+3),
						charCodeAt(body, position+4),
					)
					if charCode < 0 {
						return Token{}, gqlerrors.NewSyntaxError(s, position, "Bad character escape sequence.")
					}
					value += fmt.Sprintf("%c", charCode)
					position += 4
					break
				default:
					return Token{}, gqlerrors.NewSyntaxError(s, position, "Bad character escape sequence.")
				}
				position += 1
				chunkStart = position
			}
			continue
		} else {
			break
		}
	}
	if code != 34 {
		return Token{}, gqlerrors.NewSyntaxError(s, position, "Unterminated string.")
	}
	value += body[chunkStart:position]
	return makeToken(TokenKind[STRING], start, position+1, value), nil
}

// Converts four hexidecimal chars to the integer that the
// string represents. For example, uniCharCode('0','0','0','f')
// will return 15, and uniCharCode('0','0','f','f') returns 255.
// Returns a negative number on error, if a char was invalid.
// This is implemented by noting that char2hex() returns -1 on error,
// which means the result of ORing the char2hex() will also be negative.
func uniCharCode(a, b, c, d rune) rune {
	return rune(char2hex(a)<<12 | char2hex(b)<<8 | char2hex(c)<<4 | char2hex(d))
}

// Converts a hex character to its integer value.
// '0' becomes 0, '9' becomes 9
// 'A' becomes 10, 'F' becomes 15
// 'a' becomes 10, 'f' becomes 15
// Returns -1 on error.
func char2hex(a rune) int {
	if a >= 48 && a <= 57 { // 0-9
		return int(a) - 48
	} else if a >= 65 && a <= 70 { // A-F
		return int(a) - 55
	} else if a >= 97 && a <= 102 { // a-f
		return int(a) - 87
	} else {
		return -1
	}
}

func makeToken(kind int, start int, end int, value string) Token {
	return Token{Kind: kind, Start: start, End: end, Value: value}
}

func readToken(s *source.Source, fromPosition int) (Token, error) {
	body := s.Body
	bodyLength := len(body)
	position := positionAfterWhitespace(body, fromPosition)
	code := charCodeAt(body, position)
	if position >= bodyLength {
		return makeToken(TokenKind[EOF], position, position, ""), nil
	}
	switch code {
	// !
	case 33:
		return makeToken(TokenKind[BANG], position, position+1, ""), nil
	// $
	case 36:
		return makeToken(TokenKind[DOLLAR], position, position+1, ""), nil
	// (
	case 40:
		return makeToken(TokenKind[PAREN_L], position, position+1, ""), nil
	// )
	case 41:
		return makeToken(TokenKind[PAREN_R], position, position+1, ""), nil
	// .
	case 46:
		if charCodeAt(body, position+1) == 46 && charCodeAt(body, position+2) == 46 {
			return makeToken(TokenKind[SPREAD], position, position+3, ""), nil
		}
		break
	// :
	case 58:
		return makeToken(TokenKind[COLON], position, position+1, ""), nil
	// =
	case 61:
		return makeToken(TokenKind[EQUALS], position, position+1, ""), nil
	// @
	case 64:
		return makeToken(TokenKind[AT], position, position+1, ""), nil
	// [
	case 91:
		return makeToken(TokenKind[BRACKET_L], position, position+1, ""), nil
	// ]
	case 93:
		return makeToken(TokenKind[BRACKET_R], position, position+1, ""), nil
	// {
	case 123:
		return makeToken(TokenKind[BRACE_L], position, position+1, ""), nil
	// |
	case 124:
		return makeToken(TokenKind[PIPE], position, position+1, ""), nil
	// }
	case 125:
		return makeToken(TokenKind[BRACE_R], position, position+1, ""), nil
	// A-Z
	case 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81,
		82, 83, 84, 85, 86, 87, 88, 89, 90:
		return readName(s, position), nil
	// _
	// a-z
	case 95, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110,
		111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122:
		return readName(s, position), nil
	// -
	// 0-9
	case 45, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57:
		token, err := readNumber(s, position, code)
		if err != nil {
			return token, err
		} else {
			return token, nil
		}
	// "
	case 34:
		token, err := readString(s, position)
		if err != nil {
			return token, err
		}
		return token, nil
	}
	description := fmt.Sprintf("Unexpected character \"%c\".", code)
	return Token{}, gqlerrors.NewSyntaxError(s, position, description)
}

func charCodeAt(body string, position int) rune {
	r := []rune(body)
	if len(r) > position {
		return r[position]
	} else {
		return 0
	}
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
