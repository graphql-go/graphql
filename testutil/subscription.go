package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"testing"

	"github.com/graphql-go/graphql"
)

// TestResponse models the expected response
type TestResponse struct {
	Data   string
	Errors []string
}

// TestSubscription is a GraphQL test case to be used with RunSubscribe.
type TestSubscription struct {
	Name            string
	Schema          graphql.Schema
	Query           string
	OperationName   string
	Variables       map[string]interface{}
	ExpectedResults []TestResponse
}

// RunSubscribes runs the given GraphQL subscription test cases as subtests.
func RunSubscribes(t *testing.T, tests []*TestSubscription) {
	for i, test := range tests {
		if test.Name == "" {
			test.Name = strconv.Itoa(i + 1)
		}

		t.Run(test.Name, func(t *testing.T) {
			RunSubscribe(t, test)
		})
	}
}

// RunSubscribe runs a single GraphQL subscription test case.
func RunSubscribe(t *testing.T, test *TestSubscription) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := graphql.Subscribe(graphql.Params{
		Context:        ctx,
		OperationName:  test.OperationName,
		RequestString:  test.Query,
		VariableValues: test.Variables,
		Schema:         test.Schema,
	})
	// if err != nil {
	// 	if err.Error() != test.ExpectedErr.Error() {
	// 		t.Fatalf("unexpected error: got %+v, want %+v", err, test.ExpectedErr)
	// 	}

	// 	return
	// }

	var results []*graphql.Result
	for res := range c {
		t.Log(pretty(res))
		results = append(results, res)
	}

	for i, expected := range test.ExpectedResults {
		if len(results)-1 < i {
			t.Error(errors.New("not enough results, expected results are more than actual results"))
			return
		}
		res := results[i]

		var errs []string
		for _, err := range res.Errors {
			errs = append(errs, err.Message)
		}
		checkErrorStrings(t, expected.Errors, errs)
		if expected.Data == "" {
			continue
		}

		got, err := json.MarshalIndent(res.Data, "", "  ")
		if err != nil {
			t.Fatalf("got: invalid JSON: %s; raw: %s", err, got)
		}

		if err != nil {
			t.Fatal(err)
		}
		want, err := formatJSON(expected.Data)
		if err != nil {
			t.Fatalf("got: invalid JSON: %s; raw: %s", err, res.Data)
		}

		if !bytes.Equal(got, want) {
			t.Logf("got:  %s", got)
			t.Logf("want: %s", want)
			t.Fail()
		}
	}
}

func checkErrorStrings(t *testing.T, expected, actual []string) {
	expectedCount, actualCount := len(expected), len(actual)

	if expectedCount != actualCount {
		t.Fatalf("unexpected number of errors: want `%d`, got `%d`", expectedCount, actualCount)
	}

	if expectedCount > 0 {
		for i, want := range expected {
			got := actual[i]

			if got != want {
				t.Fatalf("unexpected error: got `%+v`, want `%+v`", got, want)
			}
		}

		// Return because we're done checking.
		return
	}

	for _, err := range actual {
		t.Errorf("unexpected error: '%s'", err)
	}
}

func formatJSON(data string) ([]byte, error) {
	var v interface{}
	if err := json.Unmarshal([]byte(data), &v); err != nil {
		return nil, err
	}
	formatted, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, err
	}
	return formatted, nil
}

func pretty(x interface{}) string {
	got, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(got)
}
