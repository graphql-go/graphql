package lexer

import (
	"reflect"
	"testing"

	"github.com/chris-ramon/graphql-go/language/source"
)

type Test struct {
	Body     string
	Expected interface{}
}

func TestSkipsWhiteSpace(t *testing.T) {
	tests := []Test{
		Test{
			Body: `

    foo

`,
			Expected: Token{
				Kind:  TokenKind[NAME],
				Start: 6,
				End:   9,
				Value: "foo",
			},
		},
		Test{
			Body: `
    #comment
    foo#comment
`,
			Expected: Token{
				Kind:  TokenKind[NAME],
				Start: 18,
				End:   21,
				Value: "foo",
			},
		},
		Test{
			Body: `,,,foo,,,`,
			Expected: Token{
				Kind:  TokenKind[NAME],
				Start: 3,
				End:   6,
				Value: "foo",
			},
		},
	}
	for _, test := range tests {
		token, err := Lex(&source.Source{Body: test.Body})(0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(token, test.Expected) {
			t.Fatalf("unexpected token, expected: %v, got: %v, body: %s", test.Expected, token, test.Body)
		}
	}
}

func TestErrorsRespectWhitespace(t *testing.T) {
	body := `

    ?

`
	source := source.NewSource(body, "")
	_, err := Lex(source)(0)
	expected := "Syntax Error GraphQL (3:5) Unexpected character \"?\".\n\n2: \n3:     ?\n       ^\n4: \n"
	if err.Error() != expected {
		t.Fatalf("unexpected error.\nexpected:\n%v\n\ngot:\n%v", expected, err.Error())
	}
}

func TestLexesStrings(t *testing.T) {
	tests := []Test{
		Test{
			Body: "\"simple\"",
			Expected: Token{
				Kind:  TokenKind[STRING],
				Start: 0,
				End:   8,
				Value: "simple",
			},
		},
		Test{
			Body: "\" white space \"",
			Expected: Token{
				Kind:  TokenKind[STRING],
				Start: 0,
				End:   15,
				Value: " white space ",
			},
		},
		Test{
			Body: "\"quote \\\"\"",
			Expected: Token{
				Kind:  TokenKind[STRING],
				Start: 0,
				End:   10,
				Value: `quote "`,
			},
		},
		Test{
			Body: "\"escaped \\n\\r\\b\\t\\f\"",
			Expected: Token{
				Kind:  TokenKind[STRING],
				Start: 0,
				End:   20,
				Value: "escaped \n\r\b\t\f",
			},
		},
		Test{
			Body: "\"slashes \\\\ \\/\"",
			Expected: Token{
				Kind:  TokenKind[STRING],
				Start: 0,
				End:   15,
				Value: "slashes \\ \\/",
			},
		},
		Test{
			Body: "\"unicode \\u1234\\u5678\\u90AB\\uCDEF\"",
			Expected: Token{
				Kind:  TokenKind[STRING],
				Start: 0,
				End:   34,
				Value: "unicode \u1234\u5678\u90AB\uCDEF",
			},
		},
	}
	for _, test := range tests {
		token, err := Lex(&source.Source{Body: test.Body})(0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(token, test.Expected) {
			t.Fatalf("unexpected token, expected: %v, got: %v", test.Expected, token)
		}
	}
}

func TestLexReportsUsefulStringErrors(t *testing.T) {
	tests := []Test{
		Test{
			Body: "\"no end quote",
			Expected: `Syntax Error GraphQL (1:14) Unterminated string.

1: "no end quote
                ^
`,
		},
		Test{
			Body: "\"multi\nline\"",
			Expected: `Syntax Error GraphQL (1:7) Unterminated string.

1: "multi
         ^
2: line"
`,
		},
		Test{
			Body: "\"multi\rline\"",
			Expected: `Syntax Error GraphQL (1:7) Unterminated string.

1: "multi
         ^
2: line"
`,
		},
		Test{
			Body: "\"multi\u2028line\"",
			Expected: `Syntax Error GraphQL (1:7) Unterminated string.

1: "multi
         ^
2: line"
`,
		},
		Test{
			Body: "\"multi\u2029line\"",
			Expected: `Syntax Error GraphQL (1:7) Unterminated string.

1: "multi
         ^
2: line"
`,
		},
		Test{
			Body: "\"bad \\z esc\"",
			Expected: `Syntax Error GraphQL (1:7) Bad character escape sequence.

1: "bad \z esc"
         ^
`,
		},
		Test{
			Body: "\"bad \\x esc\"",
			Expected: `Syntax Error GraphQL (1:7) Bad character escape sequence.

1: "bad \x esc"
         ^
`,
		},
		Test{
			Body: "\"bad \\u1 esc\"",
			Expected: `Syntax Error GraphQL (1:7) Bad character escape sequence.

1: "bad \u1 esc"
         ^
`,
		},
		Test{
			Body: "\"bad \\u0XX1 esc\"",
			Expected: `Syntax Error GraphQL (1:7) Bad character escape sequence.

1: "bad \u0XX1 esc"
         ^
`,
		},
		Test{
			Body: "\"bad \\uXXXX esc\"",
			Expected: `Syntax Error GraphQL (1:7) Bad character escape sequence.

1: "bad \uXXXX esc"
         ^
`,
		},
		Test{
			Body: "\"bad \\uFXXX esc\"",
			Expected: `Syntax Error GraphQL (1:7) Bad character escape sequence.

1: "bad \uFXXX esc"
         ^
`,
		},
		Test{
			Body: "\"bad \\uXXXF esc\"",
			Expected: `Syntax Error GraphQL (1:7) Bad character escape sequence.

1: "bad \uXXXF esc"
         ^
`,
		},
	}
	for _, test := range tests {
		_, err := Lex(source.NewSource(test.Body, ""))(0)
		if err.Error() != test.Expected {
			t.Fatalf("unexpected error.\nexpected:\n%v\n\ngot:\n%v", test.Expected, err.Error())
		}
	}
}

func TestLexesNumbers(t *testing.T) {
	tests := []Test{
		Test{
			Body: "4",
			Expected: Token{
				Kind:  TokenKind[INT],
				Start: 0,
				End:   1,
				Value: "4",
			},
		},
		Test{
			Body: "4.123",
			Expected: Token{
				Kind:  TokenKind[FLOAT],
				Start: 0,
				End:   5,
				Value: "4.123",
			},
		},
		Test{
			Body: "-4",
			Expected: Token{
				Kind:  TokenKind[INT],
				Start: 0,
				End:   2,
				Value: "-4",
			},
		},
		Test{
			Body: "9",
			Expected: Token{
				Kind:  TokenKind[INT],
				Start: 0,
				End:   1,
				Value: "9",
			},
		},
		Test{
			Body: "0",
			Expected: Token{
				Kind:  TokenKind[INT],
				Start: 0,
				End:   1,
				Value: "0",
			},
		},
		Test{
			Body: "-4.123",
			Expected: Token{
				Kind:  TokenKind[FLOAT],
				Start: 0,
				End:   6,
				Value: "-4.123",
			},
		},
		Test{
			Body: "0.123",
			Expected: Token{
				Kind:  TokenKind[FLOAT],
				Start: 0,
				End:   5,
				Value: "0.123",
			},
		},
		Test{
			Body: "123e4",
			Expected: Token{
				Kind:  TokenKind[FLOAT],
				Start: 0,
				End:   5,
				Value: "123e4",
			},
		},
		Test{
			Body: "123E4",
			Expected: Token{
				Kind:  TokenKind[FLOAT],
				Start: 0,
				End:   5,
				Value: "123E4",
			},
		},
		Test{
			Body: "123e-4",
			Expected: Token{
				Kind:  TokenKind[FLOAT],
				Start: 0,
				End:   6,
				Value: "123e-4",
			},
		},
		Test{
			Body: "123e+4",
			Expected: Token{
				Kind:  TokenKind[FLOAT],
				Start: 0,
				End:   6,
				Value: "123e+4",
			},
		},
		Test{
			Body: "-1.123e4",
			Expected: Token{
				Kind:  TokenKind[FLOAT],
				Start: 0,
				End:   8,
				Value: "-1.123e4",
			},
		},
		Test{
			Body: "-1.123E4",
			Expected: Token{
				Kind:  TokenKind[FLOAT],
				Start: 0,
				End:   8,
				Value: "-1.123E4",
			},
		},
		Test{
			Body: "-1.123e-4",
			Expected: Token{
				Kind:  TokenKind[FLOAT],
				Start: 0,
				End:   9,
				Value: "-1.123e-4",
			},
		},
		Test{
			Body: "-1.123e+4",
			Expected: Token{
				Kind:  TokenKind[FLOAT],
				Start: 0,
				End:   9,
				Value: "-1.123e+4",
			},
		},
		Test{
			Body: "-1.123e4567",
			Expected: Token{
				Kind:  TokenKind[FLOAT],
				Start: 0,
				End:   11,
				Value: "-1.123e4567",
			},
		},
	}
	for _, test := range tests {
		token, err := Lex(source.NewSource(test.Body, ""))(0)
		if err != nil {
			t.Fatalf("unexpected error: %v, test: %s", err, test)
		}
		if !reflect.DeepEqual(token, test.Expected) {
			t.Fatalf("unexpected token, expected: %v, got: %v, test: %v", test.Expected, token, test)
		}
	}
}
