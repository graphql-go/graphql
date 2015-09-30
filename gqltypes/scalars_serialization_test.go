package gqltypes

import (
	"math"
	"reflect"
	"testing"
)

type intSerializationTest struct {
	Value    interface{}
	Expected int
}
type float32SerializationTest struct {
	Value    interface{}
	Expected float32
}

type stringSerializationTest struct {
	Value    interface{}
	Expected string
}

type boolSerializationTest struct {
	Value    interface{}
	Expected bool
}

func TestSerializesOutputInt(t *testing.T) {
	tests := []intSerializationTest{
		{1, 1},
		{0, 0},
		{-1, -1},
		{float32(0.1), 0},
		{float32(1.1), 1},
		{float32(-1.1), -1},
		{float32(1e5), 100000}, // Bigger than 2^32, but still representable as an Int
		{9876504321, 9876504321},
		{-9876504321, -9876504321},
		{float64(1e100), math.MinInt64},
		{float64(-1e100), math.MinInt64},
		{"-1.1", -1},
		{"one", 0},
		{false, 0},
		{true, 1},
	}

	for _, test := range tests {
		val := GraphQLInt.Serialize(test.Value)
		if val != test.Expected {
			reflectedValue := reflect.ValueOf(test.Value)
			t.Fatalf("Failed GraphQLInt.Serialize(%v(%v)), expected: %v, got %v", reflectedValue.Type(), test.Value, test.Expected, val)
		}
	}
}

func TestSerializesOutputFloat(t *testing.T) {
	tests := []float32SerializationTest{
		{int(1), 1.0},
		{int(0), 0.0},
		{int(-1), -1.0},
		{float32(0.1), 0.1},
		{float32(1.1), 1.1},
		{float32(-1.1), -1.1},
		{"-1.1", -1.1},
		{"one", 0.0},
		{false, 0.0},
		{true, 1.0},
	}

	for _, test := range tests {
		val := GraphQLFloat.Serialize(test.Value)
		if val != test.Expected {
			reflectedValue := reflect.ValueOf(test.Value)
			t.Fatalf("Failed GraphQLFloat.Serialize(%v(%v)), expected: %v, got %v", reflectedValue.Type(), test.Value, test.Expected, val)
		}
	}
}

func TestSerializesOutputStrings(t *testing.T) {
	tests := []stringSerializationTest{
		{"string", "string"},
		{int(1), "1"},
		{float32(-1.1), "-1.1"},
		{true, "true"},
		{false, "false"},
	}

	for _, test := range tests {
		val := GraphQLString.Serialize(test.Value)
		if val != test.Expected {
			reflectedValue := reflect.ValueOf(test.Value)
			t.Fatalf("Failed GraphQLString.Serialize(%v(%v)), expected: %v, got %v", reflectedValue.Type(), test.Value, test.Expected, val)
		}
	}
}

func TestSerializesOutputBoolean(t *testing.T) {
	tests := []boolSerializationTest{
		{"true", true},
		{"false", false},
		{"string", false},
		{"", false},
		{int(1), true},
		{int(0), false},
		{true, true},
		{false, false},
	}

	for _, test := range tests {
		val := GraphQLBoolean.Serialize(test.Value)
		if val != test.Expected {
			reflectedValue := reflect.ValueOf(test.Value)
			t.Fatalf("Failed GraphQLString.GraphQLBoolean(%v(%v)), expected: %v, got %v", reflectedValue.Type(), test.Value, test.Expected, val)
		}
	}
}
