package printer_test

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/sprucehealth/graphql/language/ast"
	"github.com/sprucehealth/graphql/language/printer"
	"github.com/sprucehealth/graphql/testutil"
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
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, results))
	}
}

func TestSchemaPrinter_DoesNotAlterAST(t *testing.T) {
	b, err := ioutil.ReadFile("../../schema-kitchen-sink.graphql")
	if err != nil {
		t.Fatalf("unable to load schema-kitchen-sink.graphql")
	}

	query := string(b)
	astDoc := parse(t, query)

	astDocBefore := testutil.ASTToJSON(t, astDoc)

	_ = printer.Print(astDoc)

	astDocAfter := testutil.ASTToJSON(t, astDoc)

	_ = testutil.ASTToJSON(t, astDoc)

	if !reflect.DeepEqual(astDocAfter, astDocBefore) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(astDocAfter, astDocBefore))
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
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(results, expected))
	}
}
