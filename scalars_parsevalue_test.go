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
