package graphql

import (
	"reflect"
	"testing"
)

func TestQuotedOrList_DoesNoAcceptAnEmptyList(t *testing.T) {
	expected := ""
	result := quotedOrList([]string{})
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Expected %v, got: %v", expected, result)
	}
}
func TestQuotedOrList_ReturnsSingleQuotedItem(t *testing.T) {
	expected := `"A"`
	result := quotedOrList([]string{"A"})
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Expected %v, got: %v", expected, result)
	}
}
func TestQuotedOrList_ReturnsTwoItems(t *testing.T) {
	expected := `"A" or "B"`
	result := quotedOrList([]string{"A", "B"})
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Expected %v, got: %v", expected, result)
	}
}
func TestQuotedOrList_ReturnsCommaSeparatedManyItemList(t *testing.T) {
	expected := `"A", "B", "C", "D", or "E"`
	result := quotedOrList([]string{"A", "B", "C", "D", "E", "F"})
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Expected %v, got: %v", expected, result)
	}
}
