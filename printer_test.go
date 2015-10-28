package graphql

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func parse(t *testing.T, query string) *AstDocument {
	astDoc, err := Parse(ParseParams{
		Source: query,
		Options: ParseOptions{
			NoLocation: true,
		},
	})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	return astDoc
}

func TestPrinter_DoesNotAlterAST(t *testing.T) {
	b, err := ioutil.ReadFile("./kitchen-sink.graphql")
	if err != nil {
		t.Fatalf("unable to load kitchen-sink.graphql")
	}

	query := string(b)
	astDoc := parse(t, query)

	astDocBefore := ASTToJSON(t, astDoc)

	_ = Print(astDoc)

	astDocAfter := ASTToJSON(t, astDoc)

	_ = ASTToJSON(t, astDoc)

	if !reflect.DeepEqual(astDocAfter, astDocBefore) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(astDocAfter, astDocBefore))
	}
}

func TestPrinter_PrintsMinimalAST(t *testing.T) {
	astDoc := NewField(&AstField{
		Name: NewAstName(&AstName{
			Value: "foo",
		}),
	})
	results := Print(astDoc)
	expected := "foo"
	if !reflect.DeepEqual(results, expected) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, results))
	}
}

func TestPrinter_PrintsKitchenSink(t *testing.T) {
	b, err := ioutil.ReadFile("./kitchen-sink.graphql")
	if err != nil {
		t.Fatalf("unable to load kitchen-sink.graphql")
	}

	query := string(b)
	astDoc := parse(t, query)
	expected := `query namedQuery($foo: ComplexFooType, $bar: Bar = DefaultBarValue) {
  customUser: user(id: [987, 654]) {
    id
    ... on User @defer {
      field2 {
        id
        alias: field1(first: 10, after: $foo) @include(if: $foo) {
          id
          ...frag
        }
      }
    }
  }
}

mutation favPost {
  fav(post: 123) @defer {
    post {
      id
    }
  }
}

fragment frag on Follower {
  foo(size: $size, bar: $b, obj: {key: "value"})
}

{
  unnamed(truthyVal: true, falseyVal: false)
  query
}
`
	results := Print(astDoc)

	if !reflect.DeepEqual(expected, results) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(results, expected))
	}
}
