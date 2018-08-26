package graphql_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/graphql-go/graphql"
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

func TestTypeSystem_Scalar_ParseValueOutputDuration(t *testing.T) {
	str := "1h30m10s"
	d, _ := time.ParseDuration(str)

	tests := []struct {
		in  interface{}
		out interface{}
	}{
		{nil, nil},
		{"", nil},
		{(*string)(nil), nil},
		{"2017-07-23", nil},
		{str, d},
		{&str, d},
		{int64(5410000000000), d},
		{d.Seconds(), d},
	}

	for _, test := range tests {
		val := graphql.Duration.ParseValue(test.in)
		if val != test.out {
			reflectedValue := reflect.ValueOf(test.in)
			t.Fatalf("failed Duration.ParseValue(%v(%v)), expected: %v, got %v", reflectedValue.Type(), test.in, test.out, val)
		}
	}
}
