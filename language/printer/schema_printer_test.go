package printer

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func parseee(t *testing.T, query string) *AstDocument {
	astDoc, err := Parse(ParseParams{
		Source: query,
		Options: ParseOptions{
			NoLocation: false,
			NoSource:   true,
		},
	})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	return astDoc
}

func TestSchemaPrinter_PrintsMinimalAST(t *testing.T) {
	astDoc := NewAstScalarDefinition(&AstScalarDefinition{
		Name: NewAstName(&AstName{
			Value: "foo",
		}),
	})
	results := Print(astDoc)
	expected := "scalar foo"
	if !reflect.DeepEqual(results, expected) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, results))
	}
}

func TestSchemaPrinter_DoesNotAlterAST(t *testing.T) {
	b, err := ioutil.ReadFile("./schema-kitchen-sink.graphql")
	if err != nil {
		t.Fatalf("unable to load schema-kitchen-sink.graphql")
	}

	query := string(b)
	astDoc := parseee(t, query)

	astDocBefore := ASTToJSON(t, astDoc)

	_ = Print(astDoc)

	astDocAfter := ASTToJSON(t, astDoc)

	_ = ASTToJSON(t, astDoc)

	if !reflect.DeepEqual(astDocAfter, astDocBefore) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(astDocAfter, astDocBefore))
	}
}

func TestSchemaPrinter_PrintsKitchenSink(t *testing.T) {
	b, err := ioutil.ReadFile("./schema-kitchen-sink.graphql")
	if err != nil {
		t.Fatalf("unable to load schema-kitchen-sink.graphql")
	}

	query := string(b)
	astDoc := parseee(t, query)
	expected := `type Foo implements Bar {
  one: Type
  two(argument: InputType!): Type
  three(argument: InputType, other: String): Int
  four(argument: String = "string"): String
  five(argument: [String] = ["string", "string"]): String
  six(argument: InputType = {key: "value"}): Type
}

interface Bar {
  one: Type
  four(argument: String = "string"): String
}

union Feed = Story | Article | Advert

scalar CustomScalar

enum Site {
  DESKTOP
  MOBILE
}

input InputType {
  key: String!
  answer: Int = 42
}

extend type Foo {
  seven(argument: [String]): Type
}
`
	results := Print(astDoc)
	if !reflect.DeepEqual(expected, results) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(results, expected))
	}
}
