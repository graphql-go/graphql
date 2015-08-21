package lexer

import (
	"reflect"
	"testing"

	"github.com/chris-ramon/graphql-go/language/source"
)

type Expect struct {
	Body     string
	Expected Token
}

func TestLex(t *testing.T) {
	expectations := []Expect{
		Expect{
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
		Expect{
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
	}
	for _, e := range expectations {
		token := Lex(&source.Source{Body: e.Body})(0)
		if !reflect.DeepEqual(token, e.Expected) {
			t.Fatalf("unexpected token, expected: %v, got: %v", e.Expected, token)
		}
	}
}
