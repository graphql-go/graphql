package lexer

import (
	"fmt"

	"github.com/sprucehealth/graphql/gqlerrors"
	"github.com/sprucehealth/graphql/language/source"
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
func readName(s *source.Source, position int) Token {
	end := position + 1
	for {
		code := s.RuneAt(end)
		if !(code != 0 && (code == 95 ||
			code >= 48 && code <= 57 ||
			code >= 65 && code <= 90 ||
			code >= 97 && code <= 122)) {
			break
		}
		end++
	}
	return makeToken(TokenKind[NAME], position, end, s.Body()[position:end])
}

// Reads a number token from the source file, either a float
// or an int depending on whether a decimal point appears.
// Int:   -?(0|[1-9][0-9]*)
// Float: -?(0|[1-9][0-9]*)(\.[0-9]+)?((E|e)(+|-)?[0-9]+)?
func readNumber(s *source.Source, start int, firstCode rune) (Token, error) {
	code := firstCode
	position := start
	isFloat := false
	if code == '-' {
		position++
		code = s.RuneAt(position)
	}
	if code == '0' {
		position++
		code = s.RuneAt(position)
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
		code = s.RuneAt(position)
	}
	if code == '.' {
		isFloat = true
		position++
		code = s.RuneAt(position)
		p, err := readDigits(s, position, code)
		if err != nil {
			return Token{}, err
		}
		position = p
		code = s.RuneAt(position)
	}
	if code == 'E' || code == 'e' {
		isFloat = true
		position++
		code = s.RuneAt(position)
		if code == '+' || code == '-' {
			position++
			code = s.RuneAt(position)
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
	return makeToken(kind, start, position, s.Body()[start:position]), nil
}

// Returns the new position in the source after reading digits.
func readDigits(s *source.Source, start int, firstCode rune) (int, error) {
	if firstCode < '0' || firstCode > '9' {
		var description string
		if firstCode != 0 {
			description = fmt.Sprintf("Invalid number, expected digit but got: \"%c\".", firstCode)
		} else {
			description = "Invalid number, expected digit but got: EOF."
		}
		return start, gqlerrors.NewSyntaxError(s, start, description)
	}

	position := start
	code := firstCode
	for code >= '0' && code <= '9' {
		position++
		code = s.RuneAt(position)
	}
	return position, nil
}

func readString(s *source.Source, start int) (Token, error) {
	body := s.Body()
	position := start + 1
	chunkStart := position
	var code rune
	var value string
	for {
		code = s.RuneAt(position)
		if !(position < len(body) && code != 34 && code != 10 && code != 13 && code != 0x2028 && code != 0x2029) {
			break
		}
		position++
		if code == '\\' {
			value += body[chunkStart : position-1]
			code = s.RuneAt(position)
			switch code {
			case '"':
				value += "\""
			case '/':
				value += "\\/"
			case '\\':
				value += "\\"
			case 'b':
				value += "\b"
			case 'f':
				value += "\f"
			case 'n':
				value += "\n"
			case 'r':
				value += "\r"
			case 't':
				value += "\t"
			case 'u':
				charCode := uniCharCode(
					s.RuneAt(position+1),
					s.RuneAt(position+2),
					s.RuneAt(position+3),
					s.RuneAt(position+4),
				)
				if charCode < 0 {
					return Token{}, gqlerrors.NewSyntaxError(s, position, "Bad character escape sequence.")
				}
				value += fmt.Sprintf("%c", charCode)
				position += 4
			default:
				return Token{}, gqlerrors.NewSyntaxError(s, position, "Bad character escape sequence.")
			}
			position++
			chunkStart = position
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
	switch {
	case a >= '0' && a <= '9': // 0-9
		return int(a) - '0'
	case a >= 'A' && a <= 'F': // A-F
		return int(a) + 10 - 'A'
	case a >= 'a' && a <= 'f': // a-f
		return int(a) + 10 - 'a'
	}
	return -1
}

func makeToken(kind int, start int, end int, value string) Token {
	return Token{Kind: kind, Start: start, End: end, Value: value}
}

func readToken(s *source.Source, fromPosition int) (Token, error) {
	body := s.Body()
	bodyLength := len(body)
	position := positionAfterWhitespace(s, fromPosition)
	code := s.RuneAt(position)
	if position >= bodyLength {
		return makeToken(TokenKind[EOF], position, position, ""), nil
	}
	switch code {
	case '!':
		return makeToken(TokenKind[BANG], position, position+1, ""), nil
	case '$':
		return makeToken(TokenKind[DOLLAR], position, position+1, ""), nil
	case '(':
		return makeToken(TokenKind[PAREN_L], position, position+1, ""), nil
	case ')':
		return makeToken(TokenKind[PAREN_R], position, position+1, ""), nil
	case '.':
		if s.RuneAt(position+1) == '.' && s.RuneAt(position+2) == '.' {
			return makeToken(TokenKind[SPREAD], position, position+3, ""), nil
		}
		break
	case ':':
		return makeToken(TokenKind[COLON], position, position+1, ""), nil
	case '=':
		return makeToken(TokenKind[EQUALS], position, position+1, ""), nil
	case '@':
		return makeToken(TokenKind[AT], position, position+1, ""), nil
	case '[':
		return makeToken(TokenKind[BRACKET_L], position, position+1, ""), nil
	case ']':
		return makeToken(TokenKind[BRACKET_R], position, position+1, ""), nil
	case '{':
		return makeToken(TokenKind[BRACE_L], position, position+1, ""), nil
	case '|':
		return makeToken(TokenKind[PIPE], position, position+1, ""), nil
	case '}':
		return makeToken(TokenKind[BRACE_R], position, position+1, ""), nil
	case '"':
		token, err := readString(s, position)
		if err != nil {
			return token, err
		}
		return token, nil
	// A-Z
	// a-z
	// _
	case 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81,
		82, 83, 84, 85, 86, 87, 88, 89, 90, 95, 97, 98, 99, 100, 101, 102,
		103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116,
		117, 118, 119, 120, 121, 122:
		return readName(s, position), nil
	// -
	// 0-9
	case 45, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57:
		token, err := readNumber(s, position, code)
		if err != nil {
			return token, err
		}
		return token, nil
	}
	description := fmt.Sprintf("Unexpected character \"%c\".", code)
	return Token{}, gqlerrors.NewSyntaxError(s, position, description)
}

// Reads from body starting at startPosition until it finds a non-whitespace
// or commented character, then returns the position of that character for lexing.
// lexing.
func positionAfterWhitespace(s *source.Source, startPosition int) int {
	bodyLength := len(s.Body())
	position := startPosition
	for {
		if position >= bodyLength {
			break
		}
		code := s.RuneAt(position)
		if code == ' ' ||
			code == ',' ||
			code == '\xa0' ||
			code == 0x2028 || // line separator
			code == 0x2029 || // paragraph separator
			code > 8 && code < 14 { // whitespace
			position++
		} else if code == '#' {
			position++
			for {
				code := s.RuneAt(position)
				if !(position < bodyLength &&
					code != 10 && code != 13 && code != 0x2028 && code != 0x2029) {
					break
				}
				position++
			}
		} else {
			break
		}
	}
	return position
}

func GetTokenDesc(token Token) string {
	if token.Value == "" {
		return GetTokenKindDesc(token.Kind)
	}
	return fmt.Sprintf("%s \"%s\"", GetTokenKindDesc(token.Kind), token.Value)
}

func GetTokenKindDesc(kind int) string {
	return tokenDescription[kind]
}
