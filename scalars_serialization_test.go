package graphql

import (
	"reflect"
	"testing"
)

type intSerializationTest struct {
	Value    interface{}
	Expected interface{}
}
type float32SerializationTest struct {
	Value    interface{}
	Expected interface{}
}

type stringSerializationTest struct {
	Value    interface{}
	Expected string
}

type boolSerializationTest struct {
	Value    interface{}
	Expected bool
}

func TestTypeSystem_Scalar_SerializesOutputInt(t *testing.T) {
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
		{float64(1e100), nil},
		{float64(-1e100), nil},
		{"-1.1", -1},
		{"one", nil},
		{false, 0},
		{true, 1},
	}

	for _, test := range tests {
		val := Int.Serialize(test.Value)
		if val != test.Expected {
			reflectedValue := reflect.ValueOf(test.Value)
			t.Fatalf("Failed Int.Serialize(%v(%v)), expected: %v, got %v", reflectedValue.Type(), test.Value, test.Expected, val)
		}
	}
}

func TestTypeSystem_Scalar_SerializesOutputFloat(t *testing.T) {
	tests := []float32SerializationTest{
		{int(1), float32(1.0)},
		{int(0), float32(0.0)},
		{int(-1), float32(-1.0)},
		{float32(0.1), float32(0.1)},
		{float32(1.1), float32(1.1)},
		{float32(-1.1), float32(-1.1)},
		{"-1.1", float32(-1.1)},
		{"one", nil},
		{false, float32(0.0)},
		{true, float32(1.0)},
	}

	for i, test := range tests {
		val := Float.Serialize(test.Value)
		if val != test.Expected {
			reflectedValue := reflect.ValueOf(test.Value)
			t.Fatalf("Failed test #%d - Float.Serialize(%v(%v)), expected: %v, got %v", i, reflectedValue.Type(), test.Value, test.Expected, val)
		}
	}
}

func TestTypeSystem_Scalar_SerializesOutputStrings(t *testing.T) {
	tests := []stringSerializationTest{
		{"string", "string"},
		{int(1), "1"},
		{float32(-1.1), "-1.1"},
		{true, "true"},
		{false, "false"},
	}

	for _, test := range tests {
		val := String.Serialize(test.Value)
		if val != test.Expected {
			reflectedValue := reflect.ValueOf(test.Value)
			t.Fatalf("Failed String.Serialize(%v(%v)), expected: %v, got %v", reflectedValue.Type(), test.Value, test.Expected, val)
		}
	}
}

func TestTypeSystem_Scalar_SerializesOutputBoolean(t *testing.T) {
	tests := []boolSerializationTest{
		{"true", true},
		{"false", false},
		{"string", true},
		{"", false},
		{int(1), true},
		{int(0), false},
		{true, true},
		{false, false},
	}

	for _, test := range tests {
		val := Boolean.Serialize(test.Value)
		if val != test.Expected {
			reflectedValue := reflect.ValueOf(test.Value)
			t.Fatalf("Failed String.Boolean(%v(%v)), expected: %v, got %v", reflectedValue.Type(), test.Value, test.Expected, val)
		}
	}
}
