package graphql

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

func TestAcceptsOptionToNotIncludeSource(t *testing.T) {
	opts := ParseOptions{
		NoSource: true,
	}
	params := ParseParams{
		Source:  "{ field }",
		Options: opts,
	}
	document, err := Parse(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	oDef := AstOperationDefinition{
		Kind: "OperationDefinition",
		Loc: &AstLocation{
			Start: 0, End: 9,
		},
		Operation:  "query",
		Directives: []*AstDirective{},
		SelectionSet: &AstSelectionSet{
			Kind: "SelectionSet",
			Loc: &AstLocation{
				Start: 0, End: 9,
			},
			Selections: []Selection{
				&AstField{
					Kind: "Field",
					Loc: &AstLocation{
						Start: 2, End: 7,
					},
					Name: &AstName{
						Kind: "Name",
						Loc: &AstLocation{
							Start: 2, End: 7,
						},
						Value: "field",
					},
					Arguments:  []*AstArgument{},
					Directives: []*AstDirective{},
				},
			},
		},
	}
	expectedDocument := NewAstDocument(&AstDocument{
		Loc: &AstLocation{
			Start: 0, End: 9,
		},
		Definitions: []Node{&oDef},
	})
	if !reflect.DeepEqual(document, expectedDocument) {
		t.Fatalf("unexpected document, expected: %v, got: %v", expectedDocument, document)
	}
}

func TestParseProvidesUsefulErrors(t *testing.T) {
	opts := ParseOptions{
		NoSource: true,
	}
	params := ParseParams{
		Source:  "{",
		Options: opts,
	}
	_, err := Parse(params)

	expectedError := &Error{
		Message: `Syntax Error GraphQL (1:2) Expected Name, found EOF

1: {
    ^
`,
		Positions: []int{1},
		Locations: []SourceLocation{{1, 2}},
	}
	checkError(t, err, expectedError)

	testErrorMessagesTable := []errorMessageTest{
		{
			`{ ...MissingOn }
fragment MissingOn Type
`,
			`Syntax Error GraphQL (2:20) Expected "on", found Name "Type"`,
			false,
		},
		{
			`{ field: {} }`,
			`Syntax Error GraphQL (1:10) Expected Name, found {`,
			false,
		},
		{
			`notanoperation Foo { field }`,
			`Syntax Error GraphQL (1:1) Unexpected Name "notanoperation"`,
			false,
		},
		{
			"...",
			`Syntax Error GraphQL (1:1) Unexpected ...`,
			false,
		},
	}
	for _, test := range testErrorMessagesTable {
		if test.skipped != false {
			t.Skipf("Skipped test: %v", test.source)
		}
		_, err := Parse(ParseParams{Source: test.source})
		checkErrorMessage(t, err, test.expectedMessage)
	}

}

func TestParseProvidesUsefulErrorsWhenUsingSource(t *testing.T) {
	test := errorMessageTest{
		NewSource(&Source{Body: "query", Name: "MyQuery.graphql"}),
		`Syntax Error MyQuery.graphql (1:6) Expected Name, found EOF`,
		false,
	}
	testErrorMessage(t, test)
}

func TestParsesVariableInlineValues(t *testing.T) {
	source := `{ field(complex: { a: { b: [ $var ] } }) }`
	// should not return error
	_, err := Parse(ParseParams{Source: source})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParsesConstantDefaultValues(t *testing.T) {
	test := errorMessageTest{
		`query Foo($x: Complex = { a: { b: [ $var ] } }) { field }`,
		`Syntax Error GraphQL (1:37) Unexpected $`,
		false,
	}
	testErrorMessage(t, test)
}

func TestDuplicatedKeysInInputObject(t *testing.T) {
	test := errorMessageTest{
		`{ field(arg: { a: 1, a: 2 }) }'`,
		`Syntax Error GraphQL (1:22) Duplicate input object field a.`,
		false,
	}
	testErrorMessage(t, test)
}

func TestDoesNotAcceptFragmentsNameOn(t *testing.T) {
	test := errorMessageTest{
		`fragment on on on { on }`,
		`Syntax Error GraphQL (1:10) Unexpected Name "on"`,
		false,
	}
	testErrorMessage(t, test)
}

func TestDoesNotAcceptFragmentsSpreadOfOn(t *testing.T) {
	test := errorMessageTest{
		`{ ...on }'`,
		`Syntax Error GraphQL (1:9) Expected Name, found }`,
		false,
	}
	testErrorMessage(t, test)
}

func TestDoesNotAllowNullAsValue(t *testing.T) {
	test := errorMessageTest{
		`{ fieldWithNullableStringInput(input: null) }'`,
		`Syntax Error GraphQL (1:39) Unexpected Name "null"`,
		false,
	}
	testErrorMessage(t, test)
}

func TestParsesKitchenSink(t *testing.T) {
	b, err := ioutil.ReadFile("./kitchen-sink.graphql")
	if err != nil {
		t.Fatalf("unable to load kitchen-sink.graphql")
	}
	source := string(b)
	_, err = Parse(ParseParams{Source: source})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAllowsNonKeywordsAnywhereNameIsAllowed(t *testing.T) {
	nonKeywords := []string{
		"on",
		"fragment",
		"query",
		"mutation",
		"true",
		"false",
	}
	for _, keyword := range nonKeywords {
		fragmentName := keyword
		// You can't define or reference a fragment named `on`.
		if keyword == "on" {
			fragmentName = "a"
		}
		source := fmt.Sprintf(`query %v {
			... %v
			... on %v { field }
		}
		fragment %v on Type {
		%v(%v: $%v) @%v(%v: $%v)
		}
		`, keyword, fragmentName, keyword, fragmentName, keyword, keyword, keyword, keyword, keyword, keyword)
		_, err := Parse(ParseParams{Source: source})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
}

func TestParsesExperimentalSubscriptionFeature(t *testing.T) {
	source := `
      subscription Foo {
        subscriptionField
      }
    `
	_, err := Parse(ParseParams{Source: source})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseCreatesAst(t *testing.T) {
	body := `{
  node(id: 4) {
    id,
    name
  }
}
`
	source := NewSource(&Source{Body: body})
	document, err := Parse(
		ParseParams{
			Source: source,
			Options: ParseOptions{
				NoSource: true,
			},
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	oDef := AstOperationDefinition{
		Kind: "OperationDefinition",
		Loc: &AstLocation{
			Start: 0, End: 40,
		},
		Operation:  "query",
		Directives: []*AstDirective{},
		SelectionSet: &AstSelectionSet{
			Kind: "SelectionSet",
			Loc: &AstLocation{
				Start: 0, End: 40,
			},
			Selections: []Selection{
				&AstField{
					Kind: "Field",
					Loc: &AstLocation{
						Start: 4, End: 38,
					},
					Name: &AstName{
						Kind: "Name",
						Loc: &AstLocation{
							Start: 4, End: 8,
						},
						Value: "node",
					},
					Arguments: []*AstArgument{
						{
							Kind: "Argument",
							Name: &AstName{
								Kind: "Name",
								Loc: &AstLocation{
									Start: 9, End: 11,
								},
								Value: "id",
							},
							Value: &AstIntValue{
								Kind: "IntValue",
								Loc: &AstLocation{
									Start: 13, End: 14,
								},
								Value: "4",
							},
							Loc: &AstLocation{
								Start: 9, End: 14,
							},
						},
					},
					Directives: []*AstDirective{},
					SelectionSet: &AstSelectionSet{
						Kind: "SelectionSet",
						Loc: &AstLocation{
							Start: 16, End: 38,
						},
						Selections: []Selection{
							&AstField{
								Kind: "Field",
								Loc: &AstLocation{
									Start: 22, End: 24,
								},
								Name: &AstName{
									Kind: "Name",
									Loc: &AstLocation{
										Start: 22, End: 24,
									},
									Value: "id",
								},
								Arguments:    []*AstArgument{},
								Directives:   []*AstDirective{},
								SelectionSet: nil,
							},
							&AstField{
								Kind: "Field",
								Loc: &AstLocation{
									Start: 30, End: 34,
								},
								Name: &AstName{
									Kind: "Name",
									Loc: &AstLocation{
										Start: 30, End: 34,
									},
									Value: "name",
								},
								Arguments:    []*AstArgument{},
								Directives:   []*AstDirective{},
								SelectionSet: nil,
							},
						},
					},
				},
			},
		},
	}
	expectedDocument := NewAstDocument(&AstDocument{
		Loc: &AstLocation{
			Start: 0, End: 41,
		},
		Definitions: []Node{&oDef},
	})
	if !reflect.DeepEqual(document, expectedDocument) {
		t.Fatalf("unexpected document, expected: %v, got: %v", expectedDocument, document.Definitions)
	}

}

type errorMessageTest struct {
	source          interface{}
	expectedMessage string
	skipped         bool
}

func testErrorMessage(t *testing.T, test errorMessageTest) {
	if test.skipped != false {
		t.Skipf("Skipped test: %v", test.source)
	}
	_, err := Parse(ParseParams{Source: test.source})
	checkErrorMessage(t, err, test.expectedMessage)
}

func checkError(t *testing.T, err error, expectedError *Error) {
	if expectedError == nil {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		return // ok
	}
	// else expectedError != nil
	if err == nil {
		t.Fatalf("unexpected nil error\nexpected:\n%v\n\ngot:\n%v", expectedError, err)
	}
	if err.Error() != expectedError.Message {
		t.Fatalf("unexpected error.\nexpected:\n%v\n\ngot:\n%v", expectedError, err.Error())
	}
	gErr := toError(err)
	if gErr == nil {
		t.Fatalf("unexpected nil Error")
	}
	if len(expectedError.Positions) > 0 && !reflect.DeepEqual(gErr.Positions, expectedError.Positions) {
		t.Fatalf("unexpected Error.Positions.\nexpected:\n%v\n\ngot:\n%v", expectedError.Positions, gErr.Positions)
	}
	if len(expectedError.Locations) > 0 && !reflect.DeepEqual(gErr.Locations, expectedError.Locations) {
		t.Fatalf("unexpected Error.Locations.\nexpected:\n%v\n\ngot:\n%v", expectedError.Locations, gErr.Locations)
	}
}

func checkErrorMessage(t *testing.T, err error, expectedMessage string) {
	if err == nil {
		t.Fatalf("unexpected nil error\nexpected:\n%v\n\ngot:\n%v", expectedMessage, err)
	}
	if err.Error() != expectedMessage {
		// only check first line of error message
		lines := strings.Split(err.Error(), "\n")
		if lines[0] != expectedMessage {
			t.Fatalf("unexpected error.\nexpected:\n%v\n\ngot:\n%v", expectedMessage, lines[0])
		}
	}
}

func toError(err error) *Error {
	if err == nil {
		return nil
	}
	switch err := err.(type) {
	case *Error:
		return err
	default:
		return nil
	}
}
