package executor_test

import (
	"github.com/chris-ramon/graphql-go/errors"
	"github.com/chris-ramon/graphql-go/executor"
	"github.com/chris-ramon/graphql-go/language/location"
	"github.com/chris-ramon/graphql-go/types"
	"github.com/kr/pretty"
	"reflect"
	"testing"
)

func checkList(t *testing.T, testType types.GraphQLType, testData interface{}, expected *types.GraphQLResult) {
	data := map[string]interface{}{
		"test": testData,
	}

	dataType := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "DataType",
		Fields: types.GraphQLFieldConfigMap{
			"test": &types.GraphQLFieldConfig{
				Type: testType,
			},
		},
	})
	dataType.AddFieldConfig("nest", &types.GraphQLFieldConfig{
		Type: dataType,
		Resolve: func(p types.GQLFRParams) interface{} {
			return data
		},
	})

	schema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: dataType,
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := parse(`{ nest { test } }`, t)

	// execute
	ep := executor.ExecuteParams{
		Schema: schema,
		AST:    ast,
		Root:   data,
	}
	resultChannel := make(chan *types.GraphQLResult)
	go executor.Execute(ep, resultChannel)
	result := <-resultChannel
	if len(expected.Errors) != len(result.Errors) {
		t.Fatalf("wrong result, Diff: %v", pretty.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", pretty.Diff(expected, result))
	}

}

// Describe [T] Array<T>
func TestListsListOfNullableObjectsContainsValues(t *testing.T) {
	ttype := types.NewGraphQLList(types.GraphQLInt)
	data := []interface{}{
		1, 2,
	}
	expected := &types.GraphQLResult{
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
func TestListsListOfNullableObjectsContainsNull(t *testing.T) {
	ttype := types.NewGraphQLList(types.GraphQLInt)
	data := []interface{}{
		1, nil, 2,
	}
	expected := &types.GraphQLResult{
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
func TestListsListOfNullableObjectsReturnsNull(t *testing.T) {
	ttype := types.NewGraphQLList(types.GraphQLInt)
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": nil,
			},
		},
	}
	checkList(t, ttype, nil, expected)
}

// Describe [T] Func()Array<T> // equivalent to Promise<Array<T>>
func TestListsListOfNullableFuncContainsValues(t *testing.T) {
	ttype := types.NewGraphQLList(types.GraphQLInt)

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return []interface{}{
			1, 2,
		}
	}
	expected := &types.GraphQLResult{
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
func TestListsListOfNullableFuncContainsNull(t *testing.T) {
	ttype := types.NewGraphQLList(types.GraphQLInt)

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return []interface{}{
			1, nil, 2,
		}
	}
	expected := &types.GraphQLResult{
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
func TestListsListOfNullableFuncReturnsNull(t *testing.T) {
	ttype := types.NewGraphQLList(types.GraphQLInt)

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return nil
	}
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": nil,
			},
		},
	}
	checkList(t, ttype, data, expected)
}

// Describe [T] Array<Func()<T>> // equivalent to Array<Promise<T>>
func TestListsListOfNullableArrayOfFuncContainsValues(t *testing.T) {
	ttype := types.NewGraphQLList(types.GraphQLInt)

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
	expected := &types.GraphQLResult{
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
func TestListsListOfNullableArrayOfFuncContainsNulls(t *testing.T) {
	ttype := types.NewGraphQLList(types.GraphQLInt)

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
	expected := &types.GraphQLResult{
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
func TestListsNonNullListOfNullableObjectsContainsValues(t *testing.T) {
	ttype := types.NewGraphQLNonNull(types.NewGraphQLList(types.GraphQLInt))
	data := []interface{}{
		1, 2,
	}
	expected := &types.GraphQLResult{
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
func TestListsNonNullListOfNullableObjectsContainsNull(t *testing.T) {
	ttype := types.NewGraphQLNonNull(types.NewGraphQLList(types.GraphQLInt))
	data := []interface{}{
		1, nil, 2,
	}
	expected := &types.GraphQLResult{
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
func TestListsNonNullListOfNullableObjectsReturnsNull(t *testing.T) {
	ttype := types.NewGraphQLNonNull(types.NewGraphQLList(types.GraphQLInt))
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": nil,
			},
		},
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
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
func TestListsNonNullListOfNullableFuncContainsValues(t *testing.T) {
	ttype := types.NewGraphQLNonNull(types.NewGraphQLList(types.GraphQLInt))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return []interface{}{
			1, 2,
		}
	}
	expected := &types.GraphQLResult{
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
func TestListsNonNullListOfNullableFuncContainsNull(t *testing.T) {
	ttype := types.NewGraphQLNonNull(types.NewGraphQLList(types.GraphQLInt))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return []interface{}{
			1, nil, 2,
		}
	}
	expected := &types.GraphQLResult{
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
func TestListsNonNullListOfNullableFuncReturnsNull(t *testing.T) {
	ttype := types.NewGraphQLNonNull(types.NewGraphQLList(types.GraphQLInt))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return nil
	}
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": nil,
			},
		},
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
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
func TestListsNonNullListOfNullableArrayOfFuncContainsValues(t *testing.T) {
	ttype := types.NewGraphQLNonNull(types.NewGraphQLList(types.GraphQLInt))

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
	expected := &types.GraphQLResult{
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
func TestListsNonNullListOfNullableArrayOfFuncContainsNulls(t *testing.T) {
	ttype := types.NewGraphQLNonNull(types.NewGraphQLList(types.GraphQLInt))

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
	expected := &types.GraphQLResult{
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
func TestListsNullableListOfNonNullObjectsContainsValues(t *testing.T) {
	ttype := types.NewGraphQLList(types.NewGraphQLNonNull(types.GraphQLInt))
	data := []interface{}{
		1, 2,
	}
	expected := &types.GraphQLResult{
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
func TestListsNullableListOfNonNullObjectsContainsNull(t *testing.T) {
	ttype := types.NewGraphQLList(types.NewGraphQLNonNull(types.GraphQLInt))
	data := []interface{}{
		1, nil, 2,
	}
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": []interface{}{
					// if you're looking at this and wondering why "test" != nil like `graphql-js`,
					// it's because we don't throw errors and don't terminate in the middle of
					// finding a nil value for GraphQLNonNull
					1, 2,
				},
			},
		},
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
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
func TestListsNullableListOfNonNullObjectsReturnsNull(t *testing.T) {
	ttype := types.NewGraphQLList(types.NewGraphQLNonNull(types.GraphQLInt))

	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": nil,
			},
		},
	}
	checkList(t, ttype, nil, expected)
}

// Describe [T!] Func()Array<T> // equivalent to Promise<Array<T>>
func TestListsNullableListOfNonNullFuncContainsValues(t *testing.T) {
	ttype := types.NewGraphQLList(types.NewGraphQLNonNull(types.GraphQLInt))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return []interface{}{
			1, 2,
		}
	}
	expected := &types.GraphQLResult{
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
func TestListsNullableListOfNonNullFuncContainsNull(t *testing.T) {
	ttype := types.NewGraphQLList(types.NewGraphQLNonNull(types.GraphQLInt))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return []interface{}{
			1, nil, 2,
		}
	}
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": []interface{}{
					1, 2,
				},
			},
		},
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
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
func TestListsNullableListOfNonNullFuncReturnsNull(t *testing.T) {
	ttype := types.NewGraphQLList(types.NewGraphQLNonNull(types.GraphQLInt))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return nil
	}
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": nil,
			},
		},
	}
	checkList(t, ttype, data, expected)
}

// Describe [T!] Array<Func()<T>> // equivalent to Array<Promise<T>>
func TestListsNullableListOfNonNullArrayOfFuncContainsValues(t *testing.T) {
	ttype := types.NewGraphQLList(types.NewGraphQLNonNull(types.GraphQLInt))

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
	expected := &types.GraphQLResult{
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
func TestListsNullableListOfNonNullArrayOfFuncContainsNulls(t *testing.T) {
	ttype := types.NewGraphQLList(types.NewGraphQLNonNull(types.GraphQLInt))

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
	expected := &types.GraphQLResult{
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
func TestListsNonNullListOfNonNullObjectsContainsValues(t *testing.T) {
	ttype := types.NewGraphQLNonNull(types.NewGraphQLList(types.NewGraphQLNonNull(types.GraphQLInt)))
	data := []interface{}{
		1, 2,
	}
	expected := &types.GraphQLResult{
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
func TestListsNonNullListOfNonNullObjectsContainsNull(t *testing.T) {
	ttype := types.NewGraphQLNonNull(types.NewGraphQLList(types.NewGraphQLNonNull(types.GraphQLInt)))
	data := []interface{}{
		1, nil, 2,
	}
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": []interface{}{
					// if you're looking at this and wondering why "test" != nil like `graphql-js`,
					// it's because we don't throw errors and don't terminate in the middle of
					// finding a nil value for GraphQLNonNull
					1, 2,
				},
			},
		},
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
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
func TestListsNonNullListOfNonNullObjectsReturnsNull(t *testing.T) {
	t.Skip()
	ttype := types.NewGraphQLNonNull(types.NewGraphQLList(types.NewGraphQLNonNull(types.GraphQLInt)))

	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
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
func TestListsNonNullListOfNonNullFuncContainsValues(t *testing.T) {
	ttype := types.NewGraphQLNonNull(types.NewGraphQLList(types.NewGraphQLNonNull(types.GraphQLInt)))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return []interface{}{
			1, 2,
		}
	}
	expected := &types.GraphQLResult{
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
func TestListsNonNullListOfNonNullFuncContainsNull(t *testing.T) {
	ttype := types.NewGraphQLNonNull(types.NewGraphQLList(types.NewGraphQLNonNull(types.GraphQLInt)))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return []interface{}{
			1, nil, 2,
		}
	}
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"test": []interface{}{
					1, 2,
				},
			},
		},
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
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
func TestListsNonNullListOfNonNullFuncReturnsNull(t *testing.T) {
	t.Skip()
	ttype := types.NewGraphQLNonNull(types.NewGraphQLList(types.NewGraphQLNonNull(types.GraphQLInt)))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() interface{} {
		return nil
	}
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
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
func TestListsNonNullListOfNonNullArrayOfFuncContainsValues(t *testing.T) {
	ttype := types.NewGraphQLNonNull(types.NewGraphQLList(types.NewGraphQLNonNull(types.GraphQLInt)))

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
	expected := &types.GraphQLResult{
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
func TestListsNonNullListOfNonNullArrayOfFuncContainsNulls(t *testing.T) {
	ttype := types.NewGraphQLNonNull(types.NewGraphQLList(types.NewGraphQLNonNull(types.GraphQLInt)))

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
	expected := &types.GraphQLResult{
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
