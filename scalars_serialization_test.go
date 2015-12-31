package graphql_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/graphql-go/graphql"
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
		// Bigger than 2^32, but still representable as an Int
		{float32(1e5), 100000},
		{float32(math.MaxFloat32), nil},
		{9876504321, 9876504321},
		{-9876504321, -9876504321},
		{float64(1e100), nil},
		{float64(-1e100), nil},
		{"-1.1", -1},
		{"one", nil},
		{false, 0},
		{true, 1},
		{int8(1), 1},
		{int16(1), 1},
		{int32(1), 1},
		{int64(1), 1},
		{uint(1), 1},
		{uint8(1), 1},
		{uint16(1), 1},
		{uint32(1), 1},
		{uint32(math.MaxUint32), nil},
		{uint64(1), 1},
		{uint64(math.MaxInt32), math.MaxInt32},
		{int64(math.MaxInt32) + int64(1), nil},
		{int64(math.MinInt32) - int64(1), nil},
		{uint64(math.MaxInt64) + uint64(1), nil},
		{byte(127), 127},
		{'世', int('世')},
		// testing types that don't match a value in the array.
		{[]int{}, nil},
	}

	for _, test := range tests {
		val := graphql.Int.Serialize(test.Value)
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
		val := graphql.Float.Serialize(test.Value)
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
		val := graphql.String.Serialize(test.Value)
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
		val := graphql.Boolean.Serialize(test.Value)
		if val != test.Expected {
			reflectedValue := reflect.ValueOf(test.Value)
			t.Fatalf("Failed String.Boolean(%v(%v)), expected: %v, got %v", reflectedValue.Type(), test.Value, test.Expected, val)
		}
	}
}
