package parser

import "testing"

//func TestErrors(t *testing.T) {
//s := `{ ...MissingOn }
//fragment MissingOn Type`
//_, err := Parse(ParseParams{Source: s})
//expected := "Syntax Error GraphQL (2:20) Expected \"on\", found Name \"Type\""
//if err == nil {
//t.Fatalf("expected an error, got nil")
//}
//if err.Error() != expected {
//t.Errorf("wrong result, expected: %v, got: %v", expected, err.Error())
//}
//}
