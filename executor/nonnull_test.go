package executor_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/chris-ramon/graphql-go/errors"
	"github.com/chris-ramon/graphql-go/executor"
	"github.com/chris-ramon/graphql-go/language/location"
	"github.com/chris-ramon/graphql-go/testutil"
	"github.com/chris-ramon/graphql-go/types"
)

var syncError = "sync"
var nonNullSyncError = "nonNullSync"
var promiseError = "promise"
var nonNullPromiseError = "nonNullPromise"

var throwingData = map[string]interface{}{
	"sync": func() interface{} {
		panic(syncError)
	},
	"nonNullSync": func() interface{} {
		panic(nonNullSyncError)
	},
	"promise": func() interface{} {
		panic(promiseError)
	},
	"nonNullPromise": func() interface{} {
		panic(nonNullPromiseError)
	},
}

var nullingData = map[string]interface{}{
	"sync": func() interface{} {
		return nil
	},
	"nonNullSync": func() interface{} {
		return nil
	},
	"promise": func() interface{} {
		return nil
	},
	"nonNullPromise": func() interface{} {
		return nil
	},
}

var dataType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
	Name: "DataType",
	Fields: types.GraphQLFieldConfigMap{
		"sync": &types.GraphQLFieldConfig{
			Type: types.GraphQLString,
		},
		"nonNullSync": &types.GraphQLFieldConfig{
			Type: types.NewGraphQLNonNull(types.GraphQLString),
		},
		"promise": &types.GraphQLFieldConfig{
			Type: types.GraphQLString,
		},
		"nonNullPromise": &types.GraphQLFieldConfig{
			Type: types.NewGraphQLNonNull(types.GraphQLString),
		},
	},
})

var nonNullTestSchema, _ = types.NewGraphQLSchema(types.GraphQLSchemaConfig{
	Query: dataType,
})

func init() {
	throwingData["nest"] = func() interface{} {
		return throwingData
	}
	throwingData["nonNullNest"] = func() interface{} {
		return throwingData
	}
	throwingData["promiseNest"] = func() interface{} {
		return throwingData
	}
	throwingData["nonNullPromiseNest"] = func() interface{} {
		return throwingData
	}

	nullingData["nest"] = func() interface{} {
		return nullingData
	}
	nullingData["nonNullNest"] = func() interface{} {
		return nullingData
	}
	nullingData["promiseNest"] = func() interface{} {
		return nullingData
	}
	nullingData["nonNullPromiseNest"] = func() interface{} {
		return nullingData
	}

	dataType.AddFieldConfig("nest", &types.GraphQLFieldConfig{
		Type: dataType,
	})
	dataType.AddFieldConfig("nonNullNest", &types.GraphQLFieldConfig{
		Type: types.NewGraphQLNonNull(dataType),
	})
	dataType.AddFieldConfig("promiseNest", &types.GraphQLFieldConfig{
		Type: dataType,
	})
	dataType.AddFieldConfig("nonNullPromiseNest", &types.GraphQLFieldConfig{
		Type: types.NewGraphQLNonNull(dataType),
	})
}

// nulls a nullable field that panics
func TestNonNull_NullsANullableFieldThatThrowsSynchronously(t *testing.T) {
	doc := `
      query Q {
        sync
      }
	`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"sync": nil,
		},
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: syncError,
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line: 3, Column: 9,
					},
				},
			},
		},
	}
	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsANullableFieldThatThrowsInAPromise(t *testing.T) {
	doc := `
      query Q {
        promise
      }
	`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"promise": nil,
		},
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: promiseError,
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line: 3, Column: 9,
					},
				},
			},
		},
	}
	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsASynchronouslyReturnedObjectThatContainsANullableFieldThatThrowsSynchronously(t *testing.T) {
	doc := `
      query Q {
        nest {
          nonNullSync,
        }
      }
	`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: nonNullSyncError,
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line: 4, Column: 11,
					},
				},
			},
		},
	}
	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsASynchronouslyReturnedObjectThatContainsANullableFieldThatThrowsInAPromise(t *testing.T) {
	t.Skipf("Promises not implemented")
	doc := `
      query Q {
        nest {
          nonNullPromise,
        }
      }
	`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: nonNullPromiseError,
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line: 4, Column: 11,
					},
				},
			},
		},
	}
	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsAnObjectReturnedInAPromiseThatContainsANonNulableFieldThatThrowsSynchronously(t *testing.T) {
	t.Skipf("Promises not implemented")
	doc := `
      query Q {
        promiseNest {
          nonNullSync,
        }
      }
	`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"promiseNest": nil,
		},
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: nonNullSyncError,
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line: 4, Column: 11,
					},
				},
			},
		},
	}
	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsAnObjectReturnedInAPromiseThatContainsANonNulableFieldThatThrowsInAPromise(t *testing.T) {
	t.Skipf("Promises not implemented")
	doc := `
      query Q {
        promiseNest {
          nonNullPromise,
        }
      }
	`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"promiseNest": nil,
		},
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: nonNullPromiseError,
				Locations: []location.SourceLocation{
					location.SourceLocation{
						Line: 4, Column: 11,
					},
				},
			},
		},
	}
	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestNonNull_NullsAComplexTreeOfNullableFieldsThatThrow(t *testing.T) {
	doc := `
      query Q {
        nest {
          sync
          promise
          nest {
            sync
            promise
          }
          promiseNest {
            sync
            promise
          }
        }
        promiseNest {
          sync
          promise
          nest {
            sync
            promise
          }
          promiseNest {
            sync
            promise
          }
        }
      }
	`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"sync":    nil,
				"promise": nil,
				"nest": map[string]interface{}{
					"sync":    nil,
					"promise": nil,
				},
				"promiseNest": map[string]interface{}{
					"sync":    nil,
					"promise": nil,
				},
			},
			"promiseNest": map[string]interface{}{
				"sync":    nil,
				"promise": nil,
				"nest": map[string]interface{}{
					"sync":    nil,
					"promise": nil,
				},
				"promiseNest": map[string]interface{}{
					"sync":    nil,
					"promise": nil,
				},
			},
		},
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: syncError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 4, Column: 11},
				},
			},
			graphqlerrors.GraphQLFormattedError{
				Message: syncError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 7, Column: 13},
				},
			},
			graphqlerrors.GraphQLFormattedError{
				Message: syncError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 11, Column: 13},
				},
			},
			graphqlerrors.GraphQLFormattedError{
				Message: syncError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 16, Column: 11},
				},
			},
			graphqlerrors.GraphQLFormattedError{
				Message: syncError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 19, Column: 13},
				},
			},
			graphqlerrors.GraphQLFormattedError{
				Message: syncError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 23, Column: 13},
				},
			},
			graphqlerrors.GraphQLFormattedError{
				Message: promiseError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 5, Column: 11},
				},
			},
			graphqlerrors.GraphQLFormattedError{
				Message: promiseError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 8, Column: 13},
				},
			},
			graphqlerrors.GraphQLFormattedError{
				Message: promiseError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 12, Column: 13},
				},
			},
			graphqlerrors.GraphQLFormattedError{
				Message: promiseError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 17, Column: 11},
				},
			},
			graphqlerrors.GraphQLFormattedError{
				Message: promiseError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 20, Column: 13},
				},
			},
			graphqlerrors.GraphQLFormattedError{
				Message: promiseError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 24, Column: 13},
				},
			},
		},
	}
	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected.Data, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Data, result.Data))
	}
	sort.Sort(graphqlerrors.GQLFormattedErrorSlice(expected.Errors))
	sort.Sort(graphqlerrors.GQLFormattedErrorSlice(result.Errors))
	if !reflect.DeepEqual(expected.Errors, result.Errors) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
}
func TestNonNull_NullsTheFirstNullableObjectAfterAFieldThrowsInALongChainOfFieldsThatAreNonNull(t *testing.T) {
	doc := `
      query Q {
        nest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullSync
                }
              }
            }
          }
        }
        promiseNest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullSync
                }
              }
            }
          }
        }
        anotherNest: nest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullPromise
                }
              }
            }
          }
        }
        anotherPromiseNest: promiseNest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullPromise
                }
              }
            }
          }
        }
      }
	`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"nest":               nil,
			"promiseNest":        nil,
			"anotherNest":        nil,
			"anotherPromiseNest": nil,
		},
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: `Cannot return null for non-nullable field DataType.nonNullSync.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 8, Column: 19},
				},
			},
			graphqlerrors.GraphQLFormattedError{
				Message: `Cannot return null for non-nullable field DataType.nonNullSync.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 19, Column: 19},
				},
			},
			graphqlerrors.GraphQLFormattedError{
				Message: `Cannot return null for non-nullable field DataType.nonNullPromise.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 30, Column: 19},
				},
			},
			graphqlerrors.GraphQLFormattedError{
				Message: `Cannot return null for non-nullable field DataType.nonNullPromise.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 41, Column: 19},
				},
			},
		},
	}
	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected.Data, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Data, result.Data))
	}
	t.Skipf("Testing equality for slice of errors in results")
	if !testutil.EqualSet(expected.Errors, result.Errors) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}

}
func TestNonNull_NullsANullableFieldThatSynchronouslyReturnsNull(t *testing.T) {
	doc := `
      query Q {
        sync
      }
	`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"sync": nil,
		},
	}
	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected.Data, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Data, result.Data))
	}
	if !testutil.EqualSet(expected.Errors, result.Errors) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
}
func TestNonNull_NullsANullableFieldThatSynchronouslyReturnsNullInAPromise(t *testing.T) {
	doc := `
      query Q {
        promise
      }
	`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"promise": nil,
		},
	}
	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected.Data, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Data, result.Data))
	}
	if !testutil.EqualSet(expected.Errors, result.Errors) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
}
func TestNonNull_NullsASynchronouslyReturnedObjectThatContainsANonNullableFieldThatReturnsNullSynchronously(t *testing.T) {
	doc := `
      query Q {
        nest {
          nonNullSync,
        }
      }
	`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: `Cannot return null for non-nullable field DataType.nonNullSync.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 4, Column: 11},
				},
			},
		},
	}
	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsASynchronouslyReturnedObjectThatContainsANonNullableFieldThatReturnsNullInAPromise(t *testing.T) {
	doc := `
      query Q {
        nest {
          nonNullPromise,
        }
      }
	`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: `Cannot return null for non-nullable field DataType.nonNullPromise.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 4, Column: 11},
				},
			},
		},
	}
	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestNonNull_NullsAnObjectReturnedInAPromiseThatContainsANonNullableFieldThatReturnsNullSynchronously(t *testing.T) {
	doc := `
      query Q {
        promiseNest {
          nonNullSync,
        }
      }
	`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"promiseNest": nil,
		},
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: `Cannot return null for non-nullable field DataType.nonNullSync.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 4, Column: 11},
				},
			},
		},
	}
	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsAnObjectReturnedInAPromiseThatContainsANonNullableFieldThatReturnsNullInAPromise(t *testing.T) {
	doc := `
      query Q {
        promiseNest {
          nonNullPromise,
        }
      }
	`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"promiseNest": nil,
		},
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: `Cannot return null for non-nullable field DataType.nonNullPromise.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 4, Column: 11},
				},
			},
		},
	}
	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsAComplexTreeOfNullableFieldsThatReturnNull(t *testing.T) {
	doc := `
      query Q {
        nest {
          sync
          promise
          nest {
            sync
            promise
          }
          promiseNest {
            sync
            promise
          }
        }
        promiseNest {
          sync
          promise
          nest {
            sync
            promise
          }
          promiseNest {
            sync
            promise
          }
        }
      }
	`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"sync":    nil,
				"promise": nil,
				"nest": map[string]interface{}{
					"sync":    nil,
					"promise": nil,
				},
				"promiseNest": map[string]interface{}{
					"sync":    nil,
					"promise": nil,
				},
			},
			"promiseNest": map[string]interface{}{
				"sync":    nil,
				"promise": nil,
				"nest": map[string]interface{}{
					"sync":    nil,
					"promise": nil,
				},
				"promiseNest": map[string]interface{}{
					"sync":    nil,
					"promise": nil,
				},
			},
		},
	}
	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected.Data, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Data, result.Data))
	}
	if !testutil.EqualSet(expected.Errors, result.Errors) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
}
func TestNonNull_NullsTheFirstNullableObjectAfterAFieldReturnsNullInALongChainOfFieldsThatAreNonNull(t *testing.T) {
	doc := `
      query Q {
        nest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullSync
                }
              }
            }
          }
        }
        promiseNest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullSync
                }
              }
            }
          }
        }
        anotherNest: nest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullPromise
                }
              }
            }
          }
        }
        anotherPromiseNest: promiseNest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullPromise
                }
              }
            }
          }
        }
      }
	`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"nest":               nil,
			"promiseNest":        nil,
			"anotherNest":        nil,
			"anotherPromiseNest": nil,
		},
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: `Cannot return null for non-nullable field DataType.nonNullSync.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 8, Column: 19},
				},
			},
			graphqlerrors.GraphQLFormattedError{
				Message: `Cannot return null for non-nullable field DataType.nonNullSync.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 19, Column: 19},
				},
			},
			graphqlerrors.GraphQLFormattedError{
				Message: `Cannot return null for non-nullable field DataType.nonNullPromise.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 30, Column: 19},
				},
			},
			graphqlerrors.GraphQLFormattedError{
				Message: `Cannot return null for non-nullable field DataType.nonNullPromise.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 41, Column: 19},
				},
			},
		},
	}
	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected.Data, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Data, result.Data))
	}
	t.Skipf("Testing equality for slice of errors in results")
	if !testutil.EqualSet(expected.Errors, result.Errors) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
}

func TestNonNull_NullsTheTopLevelIfSyncNonNullableFieldThrows(t *testing.T) {
	doc := `
      query Q { nonNullSync }
	`
	expected := &types.GraphQLResult{
		Data: nil,
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: nonNullSyncError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 2, Column: 17},
				},
			},
		},
	}
	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsTheTopLevelIfSyncNonNullableFieldErrors(t *testing.T) {
	doc := `
      query Q { nonNullPromise }
	`
	expected := &types.GraphQLResult{
		Data: nil,
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: nonNullPromiseError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 2, Column: 17},
				},
			},
		},
	}
	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsTheTopLevelIfSyncNonNullableFieldReturnsNull(t *testing.T) {
	doc := `
      query Q { nonNullSync }
	`
	expected := &types.GraphQLResult{
		Data: nil,
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: `Cannot return null for non-nullable field DataType.nonNullSync.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 2, Column: 17},
				},
			},
		},
	}
	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsTheTopLevelIfSyncNonNullableFieldResolvesNull(t *testing.T) {
	doc := `
      query Q { nonNullPromise }
	`
	expected := &types.GraphQLResult{
		Data: nil,
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: `Cannot return null for non-nullable field DataType.nonNullPromise.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 2, Column: 17},
				},
			},
		},
	}
	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
