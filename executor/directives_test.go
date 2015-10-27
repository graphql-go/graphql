package executor_test

import (
	"github.com/chris-ramon/graphql/executor"
	"github.com/chris-ramon/graphql/testutil"
	"github.com/chris-ramon/graphql/types"
	"reflect"
	"testing"
)

var directivesTestSchema, _ = types.NewSchema(types.SchemaConfig{
	Query: types.NewObject(types.ObjectConfig{
		Name: "TestType",
		Fields: types.FieldConfigMap{
			"a": &types.FieldConfig{
				Type: types.String,
			},
			"b": &types.FieldConfig{
				Type: types.String,
			},
		},
	}),
})

var directivesTestData map[string]interface{} = map[string]interface{}{
	"a": func() interface{} { return "a" },
	"b": func() interface{} { return "b" },
}

func executeDirectivesTestQuery(t *testing.T, doc string) *types.Result {
	ast := testutil.Parse(t, doc)
	ep := executor.ExecuteParams{
		Schema: directivesTestSchema,
		AST:    ast,
		Root:   directivesTestData,
	}
	return testutil.Execute(t, ep)
}

func TestDirectivesWorksWithoutDirectives(t *testing.T) {
	query := `{ a, b }`
	expected := &types.Result{
		Data: map[string]interface{}{
			"a": "a",
			"b": "b",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnScalarsIfTrueIncludesScalar(t *testing.T) {
	query := `{ a, b @include(if: true) }`
	expected := &types.Result{
		Data: map[string]interface{}{
			"a": "a",
			"b": "b",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnScalarsIfFalseOmitsOnScalar(t *testing.T) {
	query := `{ a, b @include(if: false) }`
	expected := &types.Result{
		Data: map[string]interface{}{
			"a": "a",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnScalarsUnlessFalseIncludesScalar(t *testing.T) {
	query := `{ a, b @skip(if: false) }`
	expected := &types.Result{
		Data: map[string]interface{}{
			"a": "a",
			"b": "b",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnScalarsUnlessTrueOmitsScalar(t *testing.T) {
	query := `{ a, b @skip(if: true) }`
	expected := &types.Result{
		Data: map[string]interface{}{
			"a": "a",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnFragmentSpreadsIfFalseOmitsFragmentSpread(t *testing.T) {
	query := `
        query Q {
          a
          ...Frag @include(if: false)
        }
        fragment Frag on TestType {
          b
        }
	`
	expected := &types.Result{
		Data: map[string]interface{}{
			"a": "a",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnFragmentSpreadsIfTrueIncludesFragmentSpread(t *testing.T) {
	query := `
        query Q {
          a
          ...Frag @include(if: true)
        }
        fragment Frag on TestType {
          b
        }
	`
	expected := &types.Result{
		Data: map[string]interface{}{
			"a": "a",
			"b": "b",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnFragmentSpreadsUnlessFalseIncludesFragmentSpread(t *testing.T) {
	query := `
        query Q {
          a
          ...Frag @skip(if: false)
        }
        fragment Frag on TestType {
          b
        }
	`
	expected := &types.Result{
		Data: map[string]interface{}{
			"a": "a",
			"b": "b",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnFragmentSpreadsUnlessTrueOmitsFragmentSpread(t *testing.T) {
	query := `
        query Q {
          a
          ...Frag @skip(if: true)
        }
        fragment Frag on TestType {
          b
        }
	`
	expected := &types.Result{
		Data: map[string]interface{}{
			"a": "a",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnInlineFragmentIfFalseOmitsInlineFragment(t *testing.T) {
	query := `
        query Q {
          a
          ... on TestType @include(if: false) {
            b
          }
        }
        fragment Frag on TestType {
          b
        }
	`
	expected := &types.Result{
		Data: map[string]interface{}{
			"a": "a",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnInlineFragmentIfTrueIncludesInlineFragment(t *testing.T) {
	query := `
        query Q {
          a
          ... on TestType @include(if: true) {
            b
          }
        }
        fragment Frag on TestType {
          b
        }
	`
	expected := &types.Result{
		Data: map[string]interface{}{
			"a": "a",
			"b": "b",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnInlineFragmentUnlessFalseIncludesInlineFragment(t *testing.T) {
	query := `
        query Q {
          a
          ... on TestType @skip(if: false) {
            b
          }
        }
        fragment Frag on TestType {
          b
        }
	`
	expected := &types.Result{
		Data: map[string]interface{}{
			"a": "a",
			"b": "b",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnInlineFragmentUnlessTrueIncludesInlineFragment(t *testing.T) {
	query := `
        query Q {
          a
          ... on TestType @skip(if: true) {
            b
          }
        }
        fragment Frag on TestType {
          b
        }
	`
	expected := &types.Result{
		Data: map[string]interface{}{
			"a": "a",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnFragmentIfFalseOmitsFragment(t *testing.T) {
	query := `
        query Q {
          a
          ...Frag
        }
        fragment Frag on TestType @include(if: false) {
          b
        }
	`
	expected := &types.Result{
		Data: map[string]interface{}{
			"a": "a",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnFragmentIfTrueIncludesFragment(t *testing.T) {
	query := `
        query Q {
          a
          ...Frag
        }
        fragment Frag on TestType @include(if: true) {
          b
        }
	`
	expected := &types.Result{
		Data: map[string]interface{}{
			"a": "a",
			"b": "b",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnFragmentUnlessFalseIncludesFragment(t *testing.T) {
	query := `
        query Q {
          a
          ...Frag
        }
        fragment Frag on TestType @skip(if: false) {
          b
        }
	`
	expected := &types.Result{
		Data: map[string]interface{}{
			"a": "a",
			"b": "b",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnFragmentUnlessTrueOmitsFragment(t *testing.T) {
	query := `
        query Q {
          a
          ...Frag
        }
        fragment Frag on TestType @skip(if: true) {
          b
        }
	`
	expected := &types.Result{
		Data: map[string]interface{}{
			"a": "a",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
