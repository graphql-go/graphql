package graphql

import "testing"

func TestIsIterable(t *testing.T) {
	if !isIterable([]int{}) {
		t.Fatal("expected isIterable to return true for a slice, got false")
	}
	if !isIterable([]int{}) {
		t.Fatal("expected isIterable to return true for an array, got false")
	}
	if isIterable(1) {
		t.Fatal("expected isIterable to return false for an int, got true")
	}
	if isIterable(nil) {
		t.Fatal("expected isIterable to return false for nil, got true")
	}
}
