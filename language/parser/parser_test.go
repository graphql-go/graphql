package parser

import (
	"reflect"
	"testing"

	"github.com/chris-ramon/graphql-go/language/ast"
	"github.com/chris-ramon/graphql-go/language/od"
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
	if err.Error != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	oDef := od.OperationDefinition{
		Kind: "OperationDefinition",
		Loc: ast.Location{
			Start: 0, End: 9,
		},
		Operation:  "query",
		Directives: []ast.Directive{},
		SelectionSet: ast.SelectionSet{
			Kind: "SelectionSet",
			Loc: ast.Location{
				Start: 0, End: 9,
			},
			Selections: []interface{}{
				ast.Field{
					Kind: "Field",
					Loc: ast.Location{
						Start: 2, End: 7,
					},
					Name: ast.Name{
						Kind: "Name",
						Loc: ast.Location{
							Start: 2, End: 7,
						},
						Value: "field",
					},
					Arguments:  []ast.Argument{},
					Directives: []ast.Directive{},
				},
			},
		},
	}
	expectedDocument := ast.Document{
		Kind: "Document",
		Loc: ast.Location{
			Start: 0, End: 9,
		},
		Definitions: []ast.Definition{&oDef},
	}
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
	expectedError := `Syntax Error GraphQL (1:2) Expected Name, found EOF

1: {
    ^
`
	if err.Error.Error() != expectedError {
		t.Fatalf("unexpected error. \n\n expected: \n %v \n\n got: \n %v", expectedError, err)
	}
}
