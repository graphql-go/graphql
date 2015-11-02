package graphql

import (
	"reflect"
	"testing"
)

var directivesTestSchema, _ = NewSchema(SchemaConfig{
	Query: NewObject(ObjectConfig{
		Name: "TestType",
		Fields: FieldConfigMap{
			"a": &FieldConfig{
				Type: String,
			},
			"b": &FieldConfig{
				Type: String,
			},
		},
	}),
})

var directivesTestData map[string]interface{} = map[string]interface{}{
	"a": func() interface{} { return "a" },
	"b": func() interface{} { return "b" },
}

func executeDirectivesTestQuery(t *testing.T, doc string) *Result {
	ast := TestParse(t, doc)
	ep := ExecuteParams{
		Schema: directivesTestSchema,
		AST:    ast,
		Root:   directivesTestData,
	}
	return TestExecute(t, ep)
}

func TestDirectivesWorksWithoutDirectives(t *testing.T) {
	query := `{ a, b }`
	expected := &Result{
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
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}

func TestDirectivesWorksOnScalarsIfTrueIncludesScalar(t *testing.T) {
	query := `{ a, b @include(if: true) }`
	expected := &Result{
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
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}

func TestDirectivesWorksOnScalarsIfFalseOmitsOnScalar(t *testing.T) {
	query := `{ a, b @include(if: false) }`
	expected := &Result{
		Data: map[string]interface{}{
			"a": "a",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}

func TestDirectivesWorksOnScalarsUnlessFalseIncludesScalar(t *testing.T) {
	query := `{ a, b @skip(if: false) }`
	expected := &Result{
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
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}

func TestDirectivesWorksOnScalarsUnlessTrueOmitsScalar(t *testing.T) {
	query := `{ a, b @skip(if: true) }`
	expected := &Result{
		Data: map[string]interface{}{
			"a": "a",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: map[string]interface{}{
			"a": "a",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
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
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
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
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: map[string]interface{}{
			"a": "a",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: map[string]interface{}{
			"a": "a",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
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
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
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
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: map[string]interface{}{
			"a": "a",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: map[string]interface{}{
			"a": "a",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
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
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
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
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: map[string]interface{}{
			"a": "a",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
