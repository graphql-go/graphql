package printer_test

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/chris-ramon/graphql-go"
	"github.com/chris-ramon/graphql-go/language/ast"
	"github.com/chris-ramon/graphql-go/language/printer"
)

func TestSchemaPrinter_PrintsMinimalAST(t *testing.T) {
	astDoc := ast.NewScalarDefinition(&ast.ScalarDefinition{
		Name: ast.NewName(&ast.Name{
			Value: "foo",
		}),
	})
	results := printer.Print(astDoc)
	expected := "scalar foo"
	if !reflect.DeepEqual(results, expected) {
		t.Fatalf("Unexpected result, Diff: %v", graphql.Diff(expected, results))
	}
}

func TestSchemaPrinter_DoesNotAlterAST(t *testing.T) {
	b, err := ioutil.ReadFile("../../schema-kitchen-sink.graphql")
	if err != nil {
		t.Fatalf("unable to load schema-kitchen-sink.graphql")
	}

	query := string(b)
	astDoc := parse(t, query)

	astDocBefore := graphql.ASTToJSON(t, astDoc)

	_ = printer.Print(astDoc)

	astDocAfter := graphql.ASTToJSON(t, astDoc)

	_ = graphql.ASTToJSON(t, astDoc)

	if !reflect.DeepEqual(astDocAfter, astDocBefore) {
		t.Fatalf("Unexpected result, Diff: %v", graphql.Diff(astDocAfter, astDocBefore))
	}
}

func TestSchemaPrinter_PrintsKitchenSink(t *testing.T) {
	b, err := ioutil.ReadFile("../../schema-kitchen-sink.graphql")
	if err != nil {
		t.Fatalf("unable to load schema-kitchen-sink.graphql")
	}

	query := string(b)
	astDoc := parse(t, query)
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
	results := printer.Print(astDoc)
	if !reflect.DeepEqual(expected, results) {
		t.Fatalf("Unexpected result, Diff: %v", graphql.Diff(results, expected))
	}
}
