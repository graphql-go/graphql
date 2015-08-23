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

func TestSkipsWhiteSpace(t *testing.T) {
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
		Expect{
			Body: `,,,foo,,,`,
			Expected: Token{
				Kind:  TokenKind[NAME],
				Start: 3,
				End:   6,
				Value: "foo",
			},
		},
	}
	for _, e := range expectations {
		token, err := Lex(&source.Source{Body: e.Body})(0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(token, e.Expected) {
			t.Fatalf("unexpected token, expected: %v, got: %v", e.Expected, token)
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
