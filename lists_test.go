package graphql_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/location"
	"github.com/graphql-go/graphql/testutil"
)

func checkList(t *testing.T, testType graphql.Type, testData interface{}, expected *graphql.Result) {
	// TODO: uncomment t.Helper when support for go1.8 is dropped.
	//t.Helper()
	data := map[string]interface{}{
		"test": testData,
	}

	dataType := graphql.NewObject(graphql.ObjectConfig{
		Name: "DataType",
		Fields: graphql.Fields{
			"test": &graphql.Field{
				Type: testType,
			},
		},
	})
	dataType.AddFieldConfig("nest", &graphql.Field{
		Type: dataType,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return data, nil
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: dataType,
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.TestParse(t, `{ nest { test } }`)

	// execute
	ep := graphql.ExecuteParams{
		Schema: schema,
		AST:    ast,
		Root:   data,
	}
	result := testutil.TestExecute(t, ep)
	if len(expected.Errors) != len(result.Errors) {
		t.Fatalf("wrong result, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}

}

// Describe [T] Array<T>
func TestLists_ListOfNullableObjects_ContainsValues(t *testing.T) {
	ttype := graphql.NewList(graphql.Int)
	data := []interface{}{
		1, 2,
	}
	expected := &graphql.Result{
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
	ttype := graphql.NewList(graphql.Int)
	data := []interface{}{
		1, nil, 2,
	}
	expected := &graphql.Result{
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
	ttype := graphql.NewList(graphql.Int)
	expected := &graphql.Result{
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
	ttype := graphql.NewList(graphql.Int)

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return []interface{}{
			1, 2,
		}
	}
	expected := &graphql.Result{
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
	ttype := graphql.NewList(graphql.Int)

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return []interface{}{
			1, nil, 2,
		}
	}
	expected := &graphql.Result{
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
	ttype := graphql.NewList(graphql.Int)

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return nil
	}
	expected := &graphql.Result{
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
	ttype := graphql.NewList(graphql.Int)

	// `data` is a slice of functions that return values
	// Note that its uses the expected signature `func() (interface{}, error) {...}`
	data := []interface{}{
		func() (interface{}, error) {
			return 1, nil
		},
		func() (interface{}, error) {
			return 2, nil
		},
	}
	expected := &graphql.Result{
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
	ttype := graphql.NewList(graphql.Int)

	// `data` is a slice of functions that return values
	// Note that its uses the expected signature `func() (interface{}, error) {...}`
	data := []interface{}{
		func() (interface{}, error) {
			return 1, nil
		},
		func() (interface{}, error) {
			return nil, nil
		},
		func() (interface{}, error) {
			return 2, nil
		},
	}
	expected := &graphql.Result{
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
	ttype := graphql.NewNonNull(graphql.NewList(graphql.Int))
	data := []interface{}{
		1, 2,
	}
	expected := &graphql.Result{
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
	ttype := graphql.NewNonNull(graphql.NewList(graphql.Int))
	data := []interface{}{
		1, nil, 2,
	}
	expected := &graphql.Result{
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
	ttype := graphql.NewNonNull(graphql.NewList(graphql.Int))
	rootError := errors.New("Cannot return null for non-nullable field DataType.test.")
	originalError := gqlerrors.FormatError(gqlerrors.Error{
		Message: rootError.Error(),
		Locations: []location.SourceLocation{
			{
				Line:   1,
				Column: 10,
			},
		},
		Path: []interface{}{
			"nest",
			"test",
		},
		OriginalError: rootError,
	})
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{gqlerrors.FormatError(
			gqlerrors.Error{
				Message:       originalError.Message,
				Locations:     originalError.Locations,
				Path:          originalError.Path,
				OriginalError: originalError,
			}),
		},
	}
	checkList(t, ttype, nil, expected)
}

// Describe [T]! Func()Array<T> // equivalent to Promise<Array<T>>
func TestLists_NonNullListOfNullableFunc_ContainsValues(t *testing.T) {
	ttype := graphql.NewNonNull(graphql.NewList(graphql.Int))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return []interface{}{
			1, 2,
		}
	}
	expected := &graphql.Result{
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
	ttype := graphql.NewNonNull(graphql.NewList(graphql.Int))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return []interface{}{
			1, nil, 2,
		}
	}
	expected := &graphql.Result{
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
	ttype := graphql.NewNonNull(graphql.NewList(graphql.Int))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return nil
	}
	rootError := errors.New("Cannot return null for non-nullable field DataType.test.")
	originalError := gqlerrors.FormatError(gqlerrors.Error{
		Message: rootError.Error(),
		Locations: []location.SourceLocation{
			{
				Line:   1,
				Column: 10,
			},
		},
		Path: []interface{}{
			"nest",
			"test",
		},
		OriginalError: rootError,
	})
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{gqlerrors.FormatError(
			gqlerrors.Error{
				Message:       originalError.Message,
				Locations:     originalError.Locations,
				Path:          originalError.Path,
				OriginalError: originalError,
			}),
		},
	}
	checkList(t, ttype, data, expected)
}

// Describe [T]! Array<Func()<T>> // equivalent to Array<Promise<T>>
func TestLists_NonNullListOfNullableArrayOfFunc_ContainsValues(t *testing.T) {
	ttype := graphql.NewNonNull(graphql.NewList(graphql.Int))

	// `data` is a slice of functions that return values
	// Note that its uses the expected signature `func() (interface{}, error) {...}`
	data := []interface{}{
		func() (interface{}, error) {
			return 1, nil
		},
		func() (interface{}, error) {
			return 2, nil
		},
	}
	expected := &graphql.Result{
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
	ttype := graphql.NewNonNull(graphql.NewList(graphql.Int))

	// `data` is a slice of functions that return values
	// Note that its uses the expected signature `func() (interface{}, error) {...}`
	data := []interface{}{
		func() (interface{}, error) {
			return 1, nil
		},
		func() (interface{}, error) {
			return nil, nil
		},
		func() (interface{}, error) {
			return 2, nil
		},
	}
	expected := &graphql.Result{
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
	ttype := graphql.NewList(graphql.NewNonNull(graphql.Int))
	data := []interface{}{
		1, 2,
	}
	expected := &graphql.Result{
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
	ttype := graphql.NewList(graphql.NewNonNull(graphql.Int))
	data := []interface{}{
		1, nil, 2,
	}
	rootError := errors.New("Cannot return null for non-nullable field DataType.test.")
	originalError := gqlerrors.FormatError(gqlerrors.Error{
		Message: rootError.Error(),
		Locations: []location.SourceLocation{
			{
				Line:   1,
				Column: 10,
			},
		},
		Path: []interface{}{
			"nest",
			"test",
			1,
		},
		OriginalError: rootError,
	})
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": nil,
			},
		},
		Errors: []gqlerrors.FormattedError{gqlerrors.FormatError(
			gqlerrors.Error{
				Message:       originalError.Message,
				Locations:     originalError.Locations,
				Path:          originalError.Path,
				OriginalError: originalError,
			}),
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NullableListOfNonNullObjects_ReturnsNull(t *testing.T) {
	ttype := graphql.NewList(graphql.NewNonNull(graphql.Int))

	expected := &graphql.Result{
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
	ttype := graphql.NewList(graphql.NewNonNull(graphql.Int))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return []interface{}{
			1, 2,
		}
	}
	expected := &graphql.Result{
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
	ttype := graphql.NewList(graphql.NewNonNull(graphql.Int))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return []interface{}{
			1, nil, 2,
		}
	}
	rootError := errors.New("Cannot return null for non-nullable field DataType.test.")
	originalError := gqlerrors.FormatError(gqlerrors.Error{
		Message: rootError.Error(),
		Locations: []location.SourceLocation{
			{
				Line:   1,
				Column: 10,
			},
		},
		Path: []interface{}{
			"nest",
			"test",
			1,
		},
		OriginalError: rootError,
	})
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": nil,
			},
		},
		Errors: []gqlerrors.FormattedError{gqlerrors.FormatError(
			gqlerrors.Error{
				Message:       originalError.Message,
				Locations:     originalError.Locations,
				Path:          originalError.Path,
				OriginalError: originalError,
			}),
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NullableListOfNonNullFunc_ReturnsNull(t *testing.T) {
	ttype := graphql.NewList(graphql.NewNonNull(graphql.Int))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return nil
	}
	expected := &graphql.Result{
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
	ttype := graphql.NewList(graphql.NewNonNull(graphql.Int))

	// `data` is a slice of functions that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := []interface{}{
		func() (interface{}, error) {
			return 1, nil
		},
		func() (interface{}, error) {
			return 2, nil
		},
	}
	expected := &graphql.Result{
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
	ttype := graphql.NewList(graphql.NewNonNull(graphql.Int))

	// `data` is a slice of functions that return values
	// Note that its uses the expected signature `func() (interface{}, error){...}`
	data := []interface{}{
		func() (interface{}, error) {
			return 1, nil
		},
		func() (interface{}, error) {
			return nil, nil
		},
		func() (interface{}, error) {
			return 2, nil
		},
	}
	rootError := errors.New("Cannot return null for non-nullable field DataType.test.")
	originalError := gqlerrors.FormatError(gqlerrors.Error{
		Message: rootError.Error(),
		Locations: []location.SourceLocation{
			{
				Line:   1,
				Column: 10,
			},
		},
		Path: []interface{}{
			"nest",
			"test",
			1,
		},
		OriginalError: rootError,
	})
	expected := &graphql.Result{
		/*
			// TODO: Because thunks are called after the result map has been assembled,
			// we are not able to traverse up the tree until we find a nullable type,
			// so in this case the entire data is nil. Will need some significant code
			// restructure to restore this.
			Data: map[string]interface{}{
				"nest": map[string]interface{}{
					"test": nil,
				},
			},
		*/
		Data: nil,
		Errors: []gqlerrors.FormattedError{gqlerrors.FormatError(
			gqlerrors.Error{
				Message:       originalError.Message,
				Locations:     originalError.Locations,
				Path:          originalError.Path,
				OriginalError: originalError,
			}),
		},
	}
	checkList(t, ttype, data, expected)
}

// Describe [T!]! Array<T>
func TestLists_NonNullListOfNonNullObjects_ContainsValues(t *testing.T) {
	ttype := graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.Int)))
	data := []interface{}{
		1, 2,
	}
	expected := &graphql.Result{
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
	ttype := graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.Int)))
	data := []interface{}{
		1, nil, 2,
	}
	rootError := errors.New("Cannot return null for non-nullable field DataType.test.")
	originalError := gqlerrors.FormatError(gqlerrors.Error{
		Message: rootError.Error(),
		Locations: []location.SourceLocation{
			{
				Line:   1,
				Column: 10,
			},
		},
		Path: []interface{}{
			"nest",
			"test",
			1,
		},
		OriginalError: rootError,
	})
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{gqlerrors.FormatError(
			gqlerrors.Error{
				Message:       originalError.Message,
				Locations:     originalError.Locations,
				Path:          originalError.Path,
				OriginalError: originalError,
			}),
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NonNullListOfNonNullObjects_ReturnsNull(t *testing.T) {
	ttype := graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.Int)))
	rootError := errors.New("Cannot return null for non-nullable field DataType.test.")
	originalError := gqlerrors.FormatError(gqlerrors.Error{
		Message: rootError.Error(),
		Locations: []location.SourceLocation{
			{
				Line:   1,
				Column: 10,
			},
		},
		Path: []interface{}{
			"nest",
			"test",
		},
		OriginalError: rootError,
	})
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{gqlerrors.FormatError(
			gqlerrors.Error{
				Message:       originalError.Message,
				Locations:     originalError.Locations,
				Path:          originalError.Path,
				OriginalError: originalError,
			}),
		},
	}
	checkList(t, ttype, nil, expected)
}

// Describe [T!]! Func()Array<T> // equivalent to Promise<Array<T>>
func TestLists_NonNullListOfNonNullFunc_ContainsValues(t *testing.T) {
	ttype := graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.Int)))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return []interface{}{
			1, 2,
		}
	}
	expected := &graphql.Result{
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
	ttype := graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.Int)))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return []interface{}{
			1, nil, 2,
		}
	}
	rootError := errors.New("Cannot return null for non-nullable field DataType.test.")
	originalError := gqlerrors.FormatError(gqlerrors.Error{
		Message: rootError.Error(),
		Locations: []location.SourceLocation{
			{
				Line:   1,
				Column: 10,
			},
		},
		Path: []interface{}{
			"nest",
			"test",
			1,
		},
		OriginalError: rootError,
	})
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{gqlerrors.FormatError(
			gqlerrors.Error{
				Message:       originalError.Message,
				Locations:     originalError.Locations,
				Path:          originalError.Path,
				OriginalError: originalError,
			}),
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NonNullListOfNonNullFunc_ReturnsNull(t *testing.T) {
	ttype := graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.Int)))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return nil
	}
	rootError := errors.New("Cannot return null for non-nullable field DataType.test.")
	originalError := gqlerrors.FormatError(gqlerrors.Error{
		Message: rootError.Error(),
		Locations: []location.SourceLocation{
			{
				Line:   1,
				Column: 10,
			},
		},
		Path: []interface{}{
			"nest",
			"test",
		},
		OriginalError: rootError,
	})
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{gqlerrors.FormatError(
			gqlerrors.Error{
				Message:       originalError.Message,
				Locations:     originalError.Locations,
				Path:          originalError.Path,
				OriginalError: originalError,
			}),
		},
	}
	checkList(t, ttype, data, expected)
}

// Describe [T!]! Array<Func()<T>> // equivalent to Array<Promise<T>>
func TestLists_NonNullListOfNonNullArrayOfFunc_ContainsValues(t *testing.T) {
	ttype := graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.Int)))

	// `data` is a slice of functions that return values
	// Note that its uses the expected signature `func() (interface{}, error) {...}`
	data := []interface{}{
		func() (interface{}, error) {
			return 1, nil
		},
		func() (interface{}, error) {
			return 2, nil
		},
	}
	expected := &graphql.Result{
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
	ttype := graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.Int)))

	// `data` is a slice of functions that return values
	// Note that its uses the expected signature `func() (interface{}, error) {...}`
	data := []interface{}{
		func() (interface{}, error) {
			return 1, nil
		},
		func() (interface{}, error) {
			return nil, nil
		},
		func() (interface{}, error) {
			return 2, nil
		},
	}
	rootError := errors.New("Cannot return null for non-nullable field DataType.test.")
	originalError := gqlerrors.FormatError(gqlerrors.Error{
		Message: rootError.Error(),
		Locations: []location.SourceLocation{
			{
				Line:   1,
				Column: 10,
			},
		},
		Path: []interface{}{
			"nest",
			"test",
			1,
		},
		OriginalError: rootError,
	})
	expected := &graphql.Result{
		/*
			// TODO: Because thunks are called after the result map has been assembled,
			// we are not able to traverse up the tree until we find a nullable type,
			// so in this case the entire data is nil. Will need some significant code
			// restructure to restore this.
			Data: map[string]interface{}{
				"nest": nil,
			},
		*/
		Data: nil,
		Errors: []gqlerrors.FormattedError{gqlerrors.FormatError(
			gqlerrors.Error{
				Message:       originalError.Message,
				Locations:     originalError.Locations,
				Path:          originalError.Path,
				OriginalError: originalError,
			}),
		},
	}
	checkList(t, ttype, data, expected)
}

func TestLists_UserErrorExpectIterableButDidNotGetOne(t *testing.T) {
	ttype := graphql.NewList(graphql.Int)
	data := "Not an iterable"
	originalError := gqlerrors.NewFormattedError("User Error: expected iterable, but did not find one for field DataType.test.")
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": nil,
			},
		},
		Errors: []gqlerrors.FormattedError{gqlerrors.FormatError(
			gqlerrors.Error{
				Message: originalError.Message,
				Locations: []location.SourceLocation{
					{
						Line:   1,
						Column: 10,
					},
				},
				Path: []interface{}{
					"nest",
					"test",
				},
				OriginalError: originalError,
			}),
		},
	}
	checkList(t, ttype, data, expected)
}

func TestLists_ArrayOfNullableObjects_ContainsValues(t *testing.T) {
	ttype := graphql.NewList(graphql.Int)
	data := [2]interface{}{
		1, 2,
	}
	expected := &graphql.Result{
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

func TestLists_ValueMayBeNilPointer(t *testing.T) {
	var listTestSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"list": &graphql.Field{
					Type: graphql.NewList(graphql.Int),
					Resolve: func(_ graphql.ResolveParams) (interface{}, error) {
						return []int(nil), nil
					},
				},
			},
		}),
	})
	query := "{ list }"
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"list": []interface{}{},
		},
	}
	result := g(t, graphql.Params{
		Schema:        listTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestLists_NullableListOfInt_ReturnsNull(t *testing.T) {
	ttype := graphql.NewList(graphql.Int)
	type dataType *[]int
	var data dataType
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": nil,
			},
		},
	}
	checkList(t, ttype, data, expected)
}
