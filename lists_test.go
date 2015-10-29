package graphql

import (
	"reflect"
	"testing"

	"github.com/chris-ramon/graphql/gqlerrors"
	"github.com/chris-ramon/graphql/language/location"
)

func checkList(t *testing.T, testType Type, testData interface{}, expected *Result) {
	data := map[string]interface{}{
		"test": testData,
	}

	dataType := NewObject(ObjectConfig{
		Name: "DataType",
		Fields: FieldConfigMap{
			"test": &FieldConfig{
				Type: testType,
			},
		},
	})
	dataType.AddFieldConfig("nest", &FieldConfig{
		Type: dataType,
		Resolve: func(p GQLFRParams) interface{} {
			return data
		},
	})

	schema, err := NewSchema(SchemaConfig{
		Query: dataType,
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := TestParse(t, `{ nest { test } }`)

	// execute
	ep := ExecuteParams{
		Schema: schema,
		AST:    ast,
		Root:   data,
	}
	result := TestExecute(t, ep)
	if len(expected.Errors) != len(result.Errors) {
		t.Fatalf("wrong result, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}

}

// Describe [T] Array<T>
func TestLists_ListOfNullableObjects_ContainsValues(t *testing.T) {
	ttype := NewList(Int)
	data := []interface{}{
		1, 2,
	}
	expected := &Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": []interface{}{
					1, 2,
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_ListOfNullableObjects_ContainsNull(t *testing.T) {
	ttype := NewList(Int)
	data := []interface{}{
		1, nil, 2,
	}
	expected := &Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": []interface{}{
					1, nil, 2,
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_ListOfNullableObjects_ReturnsNull(t *testing.T) {
	ttype := NewList(Int)
	expected := &Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": nil,
			},
		},
	}
	checkList(t, ttype, nil, expected)
}

// Describe [T] Func()Array<T> // equivalent to Promise<Array<T>>
func TestLists_ListOfNullableFunc_ContainsValues(t *testing.T) {
	ttype := NewList(Int)

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return []interface{}{
			1, 2,
		}
	}
	expected := &Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": []interface{}{
					1, 2,
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_ListOfNullableFunc_ContainsNull(t *testing.T) {
	ttype := NewList(Int)

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return []interface{}{
			1, nil, 2,
		}
	}
	expected := &Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": []interface{}{
					1, nil, 2,
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_ListOfNullableFunc_ReturnsNull(t *testing.T) {
	ttype := NewList(Int)

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return nil
	}
	expected := &Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": nil,
			},
		},
	}
	checkList(t, ttype, data, expected)
}

// Describe [T] Array<Func()<T>> // equivalent to Array<Promise<T>>
func TestLists_ListOfNullableArrayOfFuncContainsValues(t *testing.T) {
	ttype := NewList(Int)

	// `data` is a slice of functions that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := []interface{}{
		func() interface{} {
			return 1
		},
		func() interface{} {
			return 2
		},
	}
	expected := &Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": []interface{}{
					1, 2,
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_ListOfNullableArrayOfFuncContainsNulls(t *testing.T) {
	ttype := NewList(Int)

	// `data` is a slice of functions that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := []interface{}{
		func() interface{} {
			return 1
		},
		func() interface{} {
			return nil
		},
		func() interface{} {
			return 2
		},
	}
	expected := &Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": []interface{}{
					1, nil, 2,
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}

// Describe [T]! Array<T>
func TestLists_NonNullListOfNullableObjectsContainsValues(t *testing.T) {
	ttype := NewNonNull(NewList(Int))
	data := []interface{}{
		1, 2,
	}
	expected := &Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": []interface{}{
					1, 2,
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NonNullListOfNullableObjectsContainsNull(t *testing.T) {
	ttype := NewNonNull(NewList(Int))
	data := []interface{}{
		1, nil, 2,
	}
	expected := &Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": []interface{}{
					1, nil, 2,
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NonNullListOfNullableObjectsReturnsNull(t *testing.T) {
	ttype := NewNonNull(NewList(Int))
	expected := &Result{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: "Cannot return null for non-nullable field DataType.test.",
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line:   1,
						Column: 10,
					},
				},
			},
		},
	}
	checkList(t, ttype, nil, expected)
}

// Describe [T]! Func()Array<T> // equivalent to Promise<Array<T>>
func TestLists_NonNullListOfNullableFunc_ContainsValues(t *testing.T) {
	ttype := NewNonNull(NewList(Int))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return []interface{}{
			1, 2,
		}
	}
	expected := &Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": []interface{}{
					1, 2,
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NonNullListOfNullableFunc_ContainsNull(t *testing.T) {
	ttype := NewNonNull(NewList(Int))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return []interface{}{
			1, nil, 2,
		}
	}
	expected := &Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": []interface{}{
					1, nil, 2,
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NonNullListOfNullableFunc_ReturnsNull(t *testing.T) {
	ttype := NewNonNull(NewList(Int))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return nil
	}
	expected := &Result{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: "Cannot return null for non-nullable field DataType.test.",
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line:   1,
						Column: 10,
					},
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}

// Describe [T]! Array<Func()<T>> // equivalent to Array<Promise<T>>
func TestLists_NonNullListOfNullableArrayOfFunc_ContainsValues(t *testing.T) {
	ttype := NewNonNull(NewList(Int))

	// `data` is a slice of functions that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := []interface{}{
		func() interface{} {
			return 1
		},
		func() interface{} {
			return 2
		},
	}
	expected := &Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": []interface{}{
					1, 2,
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NonNullListOfNullableArrayOfFunc_ContainsNulls(t *testing.T) {
	ttype := NewNonNull(NewList(Int))

	// `data` is a slice of functions that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := []interface{}{
		func() interface{} {
			return 1
		},
		func() interface{} {
			return nil
		},
		func() interface{} {
			return 2
		},
	}
	expected := &Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": []interface{}{
					1, nil, 2,
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}

// Describe [T!] Array<T>
func TestLists_NullableListOfNonNullObjects_ContainsValues(t *testing.T) {
	ttype := NewList(NewNonNull(Int))
	data := []interface{}{
		1, 2,
	}
	expected := &Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": []interface{}{
					1, 2,
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NullableListOfNonNullObjects_ContainsNull(t *testing.T) {
	ttype := NewList(NewNonNull(Int))
	data := []interface{}{
		1, nil, 2,
	}
	expected := &Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": nil,
			},
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: "Cannot return null for non-nullable field DataType.test.",
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line:   1,
						Column: 10,
					},
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NullableListOfNonNullObjects_ReturnsNull(t *testing.T) {
	ttype := NewList(NewNonNull(Int))

	expected := &Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": nil,
			},
		},
	}
	checkList(t, ttype, nil, expected)
}

// Describe [T!] Func()Array<T> // equivalent to Promise<Array<T>>
func TestLists_NullableListOfNonNullFunc_ContainsValues(t *testing.T) {
	ttype := NewList(NewNonNull(Int))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return []interface{}{
			1, 2,
		}
	}
	expected := &Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": []interface{}{
					1, 2,
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NullableListOfNonNullFunc_ContainsNull(t *testing.T) {
	ttype := NewList(NewNonNull(Int))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return []interface{}{
			1, nil, 2,
		}
	}
	expected := &Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": nil,
			},
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: "Cannot return null for non-nullable field DataType.test.",
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line:   1,
						Column: 10,
					},
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NullableListOfNonNullFunc_ReturnsNull(t *testing.T) {
	ttype := NewList(NewNonNull(Int))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return nil
	}
	expected := &Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": nil,
			},
		},
	}
	checkList(t, ttype, data, expected)
}

// Describe [T!] Array<Func()<T>> // equivalent to Array<Promise<T>>
func TestLists_NullableListOfNonNullArrayOfFunc_ContainsValues(t *testing.T) {
	ttype := NewList(NewNonNull(Int))

	// `data` is a slice of functions that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := []interface{}{
		func() interface{} {
			return 1
		},
		func() interface{} {
			return 2
		},
	}
	expected := &Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": []interface{}{
					1, 2,
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NullableListOfNonNullArrayOfFunc_ContainsNulls(t *testing.T) {
	ttype := NewList(NewNonNull(Int))

	// `data` is a slice of functions that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := []interface{}{
		func() interface{} {
			return 1
		},
		func() interface{} {
			return nil
		},
		func() interface{} {
			return 2
		},
	}
	expected := &Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": []interface{}{
					1, nil, 2,
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}

// Describe [T!]! Array<T>
func TestLists_NonNullListOfNonNullObjects_ContainsValues(t *testing.T) {
	ttype := NewNonNull(NewList(NewNonNull(Int)))
	data := []interface{}{
		1, 2,
	}
	expected := &Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": []interface{}{
					1, 2,
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NonNullListOfNonNullObjects_ContainsNull(t *testing.T) {
	ttype := NewNonNull(NewList(NewNonNull(Int)))
	data := []interface{}{
		1, nil, 2,
	}
	expected := &Result{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: "Cannot return null for non-nullable field DataType.test.",
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line:   1,
						Column: 10,
					},
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NonNullListOfNonNullObjects_ReturnsNull(t *testing.T) {
	ttype := NewNonNull(NewList(NewNonNull(Int)))

	expected := &Result{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: "Cannot return null for non-nullable field DataType.test.",
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line:   1,
						Column: 10,
					},
				},
			},
		},
	}
	checkList(t, ttype, nil, expected)
}

// Describe [T!]! Func()Array<T> // equivalent to Promise<Array<T>>
func TestLists_NonNullListOfNonNullFunc_ContainsValues(t *testing.T) {
	ttype := NewNonNull(NewList(NewNonNull(Int)))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return []interface{}{
			1, 2,
		}
	}
	expected := &Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": []interface{}{
					1, 2,
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NonNullListOfNonNullFunc_ContainsNull(t *testing.T) {
	ttype := NewNonNull(NewList(NewNonNull(Int)))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return []interface{}{
			1, nil, 2,
		}
	}
	expected := &Result{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: "Cannot return null for non-nullable field DataType.test.",
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line:   1,
						Column: 10,
					},
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NonNullListOfNonNullFunc_ReturnsNull(t *testing.T) {
	ttype := NewNonNull(NewList(NewNonNull(Int)))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return nil
	}
	expected := &Result{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: "Cannot return null for non-nullable field DataType.test.",
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line:   1,
						Column: 10,
					},
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}

// Describe [T!]! Array<Func()<T>> // equivalent to Array<Promise<T>>
func TestLists_NonNullListOfNonNullArrayOfFunc_ContainsValues(t *testing.T) {
	ttype := NewNonNull(NewList(NewNonNull(Int)))

	// `data` is a slice of functions that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := []interface{}{
		func() interface{} {
			return 1
		},
		func() interface{} {
			return 2
		},
	}
	expected := &Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": []interface{}{
					1, 2,
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NonNullListOfNonNullArrayOfFunc_ContainsNulls(t *testing.T) {
	ttype := NewNonNull(NewList(NewNonNull(Int)))

	// `data` is a slice of functions that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := []interface{}{
		func() interface{} {
			return 1
		},
		func() interface{} {
			return nil
		},
		func() interface{} {
			return 2
		},
	}
	expected := &Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": []interface{}{
					1, nil, 2,
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
