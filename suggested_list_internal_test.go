package graphql

import (
	"reflect"
	"testing"
)

func TestSuggestionList_ReturnsResultsWhenInputIsEmpty(t *testing.T) {
	expected := []string{"a"}
	result := suggestionList("", []string{"a"})
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Expected %v, got: %v", expected, result)
	}
}
func TestSuggestionList_ReturnsEmptyArrayWhenThereAreNoOptions(t *testing.T) {
	expected := []string{}
	result := suggestionList("input", []string{})
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Expected %v, got: %v", expected, result)
	}
}
func TestSuggestionList_ReturnsOptionsSortedBasedOnSimilarity(t *testing.T) {
	expected := []string{"abc", "ab"}
	result := suggestionList("abc", []string{"a", "ab", "abc"})
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Expected %v, got: %v", expected, result)
	}
}
