package graphql_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

func TestTypeSystem_Scalar_ParseValueOutputDateTime(t *testing.T) {
	t1, _ := time.Parse(time.RFC3339, "2017-07-23T03:46:56.647Z")
	tests := []dateTimeSerializationTest{
		{nil, nil},
		{"", nil},
		{(*string)(nil), nil},
		{"2017-07-23", nil},
		{"2017-07-23T03:46:56.647Z", t1},
	}
	for _, test := range tests {
		val := graphql.DateTime.ParseValue(test.Value)
		if val != test.Expected {
			reflectedValue := reflect.ValueOf(test.Value)
			t.Fatalf("failed DateTime.ParseValue(%v(%v)), expected: %v, got %v", reflectedValue.Type(), test.Value, test.Expected, val)
		}
	}
}

func TestTypeSystem_Scalar_ParseLiteralOutputDateTime(t *testing.T) {
	t1, _ := time.Parse(time.RFC3339, "2017-07-23T03:46:56.647Z")
	for name, testCase := range map[string]struct {
		Literal  ast.Value
		Expected interface{}
	}{
		"String": {
			Literal: &ast.StringValue{
				Value: "2017-07-23T03:46:56.647Z",
			},
			Expected: t1,
		},
		"NotAString": {
			Literal:  &ast.IntValue{},
			Expected: nil,
		},
	} {
		t.Run(name, func(t *testing.T) {
			parsed := graphql.DateTime.ParseLiteral(testCase.Literal)
			if parsed != testCase.Expected {
				t.Fatalf("failed DateTime.ParseLiteral(%T(%v)), expected: %v, got %v", testCase.Literal, testCase.Literal, parsed, testCase.Expected)
			}
		})
	}
}
