package graphql_test

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/location"
	"github.com/graphql-go/graphql/testutil"
	"github.com/kr/pretty"
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

var dataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "DataType",
	Fields: graphql.Fields{
		"sync": &graphql.Field{
			Type: graphql.String,
		},
		"nonNullSync": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"promise": &graphql.Field{
			Type: graphql.String,
		},
		"nonNullPromise": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})

var nonNullTestSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
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

	dataType.AddFieldConfig("nest", &graphql.Field{
		Type: dataType,
	})
	dataType.AddFieldConfig("nonNullNest", &graphql.Field{
		Type: graphql.NewNonNull(dataType),
	})
	dataType.AddFieldConfig("promiseNest", &graphql.Field{
		Type: dataType,
	})
	dataType.AddFieldConfig("nonNullPromiseNest", &graphql.Field{
		Type: graphql.NewNonNull(dataType),
	})
}

// nulls a nullable field that panics
func TestNonNull_NullsANullableFieldThatThrowsSynchronously(t *testing.T) {
	doc := `
      query Q {
        sync
      }
	`
	originalError := errors.New(syncError)
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"sync": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormatError(gqlerrors.Error{
				Message: originalError.Error(),
				Locations: []location.SourceLocation{
					{
						Line: 3, Column: 9,
					},
				},
				Path: []interface{}{
					"sync",
				},
				OriginalError: originalError,
			}),
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	fmt.Printf("%v\n", pretty.Diff(expected, result))
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
	originalError := errors.New(promiseError)
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"promise": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormatError(gqlerrors.Error{
				Message: originalError.Error(),
				Locations: []location.SourceLocation{
					{
						Line: 3, Column: 9,
					},
				},
				Path: []interface{}{
					"promise",
				},
				OriginalError: originalError,
			}),
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
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
	originalError := errors.New(nonNullSyncError)
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormatError(gqlerrors.Error{
				Message: originalError.Error(),
				Locations: []location.SourceLocation{
					{
						Line: 4, Column: 11,
					},
				},
				Path: []interface{}{
					"nest",
					"nonNullSync",
				},
				OriginalError: originalError,
			}),
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsASynchronouslyReturnedObjectThatContainsANonNullableFieldThatThrowsInAPromise(t *testing.T) {
	doc := `
      query Q {
        nest {
          nonNullPromise,
        }
      }
	`
	originalError := errors.New(nonNullPromiseError)
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormatError(gqlerrors.Error{
				Message: originalError.Error(),
				Locations: []location.SourceLocation{
					{
						Line: 4, Column: 11,
					},
				},
				Path: []interface{}{
					"nest",
					"nonNullPromise",
				},
				OriginalError: originalError,
			}),
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsAnObjectReturnedInAPromiseThatContainsANonNullableFieldThatThrowsSynchronously(t *testing.T) {
	doc := `
      query Q {
        promiseNest {
          nonNullSync,
        }
      }
	`
	originalError := errors.New(nonNullSyncError)
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"promiseNest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormatError(gqlerrors.Error{
				Message: originalError.Error(),
				Locations: []location.SourceLocation{
					{
						Line: 4, Column: 11,
					},
				},
				Path: []interface{}{
					"promiseNest",
					"nonNullSync",
				},
				OriginalError: originalError,
			}),
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsAnObjectReturnedInAPromiseThatContainsANonNullableFieldThatThrowsInAPromise(t *testing.T) {
	doc := `
      query Q {
        promiseNest {
          nonNullPromise,
        }
      }
	`
	originalError := errors.New(nonNullPromiseError)
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"promiseNest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormatError(gqlerrors.Error{
				Message: originalError.Error(),
				Locations: []location.SourceLocation{
					{
						Line: 4, Column: 11,
					},
				},
				Path: []interface{}{
					"promiseNest",
					"nonNullPromise",
				},
				OriginalError: originalError,
			}),
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
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
	syncOriginalError := errors.New(syncError)
	promiseOriginalError := errors.New(promiseError)
	expected := &graphql.Result{
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
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormatError(gqlerrors.Error{
				Message: syncOriginalError.Error(),
				Locations: []location.SourceLocation{
					{Line: 4, Column: 11},
				},
				Path: []interface{}{
					"nest", "sync",
				},
				OriginalError: syncOriginalError,
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: syncOriginalError.Error(),
				Locations: []location.SourceLocation{
					{Line: 7, Column: 13},
				},
				Path: []interface{}{
					"nest", "nest", "sync",
				},
				OriginalError: syncOriginalError,
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: syncOriginalError.Error(),
				Locations: []location.SourceLocation{
					{Line: 11, Column: 13},
				},
				Path: []interface{}{
					"nest", "promiseNest", "sync",
				},
				OriginalError: syncOriginalError,
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: syncOriginalError.Error(),
				Locations: []location.SourceLocation{
					{Line: 16, Column: 11},
				},
				Path: []interface{}{
					"promiseNest", "sync",
				},
				OriginalError: syncOriginalError,
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: syncOriginalError.Error(),
				Locations: []location.SourceLocation{
					{Line: 19, Column: 13},
				},
				Path: []interface{}{
					"promiseNest", "nest", "sync",
				},
				OriginalError: syncOriginalError,
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: syncOriginalError.Error(),
				Locations: []location.SourceLocation{
					{Line: 23, Column: 13},
				},
				Path: []interface{}{
					"promiseNest", "promiseNest", "sync",
				},
				OriginalError: syncOriginalError,
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: promiseOriginalError.Error(),
				Locations: []location.SourceLocation{
					{Line: 5, Column: 11},
				},
				Path: []interface{}{
					"nest", "promise",
				},
				OriginalError: promiseOriginalError,
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: promiseOriginalError.Error(),
				Locations: []location.SourceLocation{
					{Line: 8, Column: 13},
				},
				Path: []interface{}{
					"nest", "nest", "promise",
				},
				OriginalError: promiseOriginalError,
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: promiseOriginalError.Error(),
				Locations: []location.SourceLocation{
					{Line: 12, Column: 13},
				},
				Path: []interface{}{
					"nest", "promiseNest", "promise",
				},
				OriginalError: promiseOriginalError,
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: promiseOriginalError.Error(),
				Locations: []location.SourceLocation{
					{Line: 17, Column: 11},
				},
				Path: []interface{}{
					"promiseNest", "promise",
				},
				OriginalError: promiseOriginalError,
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: promiseOriginalError.Error(),
				Locations: []location.SourceLocation{
					{Line: 20, Column: 13},
				},
				Path: []interface{}{
					"promiseNest", "nest", "promise",
				},
				OriginalError: promiseOriginalError,
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: promiseOriginalError.Error(),
				Locations: []location.SourceLocation{
					{Line: 24, Column: 13},
				},
				Path: []interface{}{
					"promiseNest", "promiseNest", "promise",
				},
				OriginalError: promiseOriginalError,
			}),
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected.Data, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Data, result.Data))
	}
	sort.Sort(gqlerrors.FormattedErrors(expected.Errors))
	sort.Sort(gqlerrors.FormattedErrors(result.Errors))
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
	nonNullSyncOriginalError := errors.New(nonNullSyncError)
	onNullPromiseOriginalError := errors.New(nonNullPromiseError)
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"nest":               nil,
			"promiseNest":        nil,
			"anotherNest":        nil,
			"anotherPromiseNest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormatError(gqlerrors.Error{
				Message: nonNullSyncOriginalError.Error(),
				Locations: []location.SourceLocation{
					{Line: 8, Column: 19},
				},
				Path: []interface{}{
					"nest", "nonNullNest", "nonNullPromiseNest", "nonNullNest",
					"nonNullPromiseNest", "nonNullSync",
				},
				OriginalError: nonNullSyncOriginalError,
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: nonNullSyncOriginalError.Error(),
				Locations: []location.SourceLocation{
					{Line: 19, Column: 19},
				},
				Path: []interface{}{
					"promiseNest", "nonNullNest", "nonNullPromiseNest", "nonNullNest",
					"nonNullPromiseNest", "nonNullSync",
				},
				OriginalError: nonNullSyncOriginalError,
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: onNullPromiseOriginalError.Error(),
				Locations: []location.SourceLocation{
					{Line: 30, Column: 19},
				},
				Path: []interface{}{
					"anotherNest", "nonNullNest", "nonNullPromiseNest", "nonNullNest",
					"nonNullPromiseNest", "nonNullPromise",
				},
				OriginalError: onNullPromiseOriginalError,
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: onNullPromiseOriginalError.Error(),
				Locations: []location.SourceLocation{
					{Line: 41, Column: 19},
				},
				Path: []interface{}{
					"anotherPromiseNest", "nonNullNest", "nonNullPromiseNest", "nonNullNest",
					"nonNullPromiseNest", "nonNullPromise",
				},
				OriginalError: onNullPromiseOriginalError,
			}),
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected.Data, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Data, result.Data))
	}
	sort.Sort(gqlerrors.FormattedErrors(expected.Errors))
	sort.Sort(gqlerrors.FormattedErrors(result.Errors))
	if !reflect.DeepEqual(expected.Errors, result.Errors) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}

}
func TestNonNull_NullsANullableFieldThatSynchronouslyReturnsNull(t *testing.T) {
	doc := `
      query Q {
        sync
      }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"sync": nil,
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected.Data, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Data, result.Data))
	}
	if !reflect.DeepEqual(expected.Errors, result.Errors) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
}
func TestNonNull_NullsANullableFieldThatSynchronouslyReturnsNullInAPromise(t *testing.T) {
	doc := `
      query Q {
        promise
      }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"promise": nil,
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected.Data, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Data, result.Data))
	}
	if !reflect.DeepEqual(expected.Errors, result.Errors) {
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
	rootError := errors.New(`Cannot return null for non-nullable field DataType.nonNullSync.`)
	originalError := gqlerrors.FormatError(gqlerrors.Error{
		Message: rootError.Error(),
		Locations: []location.SourceLocation{
			{Line: 4, Column: 11},
		},
		Path: []interface{}{
			"nest",
			"nonNullSync",
		},
		OriginalError: rootError,
	})
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormatError(gqlerrors.Error{
				Message:       originalError.Message,
				Locations:     originalError.Locations,
				Path:          originalError.Path,
				OriginalError: originalError,
			}),
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
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
	rootError := errors.New(`Cannot return null for non-nullable field DataType.nonNullPromise.`)
	originalError := gqlerrors.FormatError(gqlerrors.Error{
		Message: rootError.Error(),
		Locations: []location.SourceLocation{
			{Line: 4, Column: 11},
		},
		Path: []interface{}{
			"nest",
			"nonNullPromise",
		},
		OriginalError: rootError,
	})
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormatError(gqlerrors.Error{
				Message:       originalError.Message,
				Locations:     originalError.Locations,
				Path:          originalError.Path,
				OriginalError: originalError,
			}),
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
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
	rootError := errors.New(`Cannot return null for non-nullable field DataType.nonNullSync.`)
	originalError := gqlerrors.FormatError(gqlerrors.Error{
		Message: rootError.Error(),
		Locations: []location.SourceLocation{
			{Line: 4, Column: 11},
		},
		Path: []interface{}{
			"promiseNest",
			"nonNullSync",
		},
		OriginalError: rootError,
	})
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"promiseNest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormatError(gqlerrors.Error{
				Message:       originalError.Message,
				Locations:     originalError.Locations,
				Path:          originalError.Path,
				OriginalError: originalError,
			}),
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
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
	rootError := errors.New(`Cannot return null for non-nullable field DataType.nonNullPromise.`)
	originalError := gqlerrors.FormatError(gqlerrors.Error{
		Message: rootError.Error(),
		Locations: []location.SourceLocation{
			{Line: 4, Column: 11},
		},
		Path: []interface{}{
			"promiseNest",
			"nonNullPromise",
		},
		OriginalError: rootError,
	})
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"promiseNest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormatError(gqlerrors.Error{
				Message:       originalError.Message,
				Locations:     originalError.Locations,
				Path:          originalError.Path,
				OriginalError: originalError,
			}),
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
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
	expected := &graphql.Result{
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
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected.Data, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Data, result.Data))
	}
	if !reflect.DeepEqual(expected.Errors, result.Errors) {
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
	nonNullSyncRootError := errors.New(`Cannot return null for non-nullable field DataType.nonNullSync.`)
	nonNullSyncOriginalError := gqlerrors.FormatError(gqlerrors.Error{
		Message: nonNullSyncRootError.Error(),
		Locations: []location.SourceLocation{
			{Line: 8, Column: 19},
		},
		Path: []interface{}{
			"nest", "nonNullNest", "nonNullPromiseNest", "nonNullNest",
			"nonNullPromiseNest", "nonNullSync",
		},
		OriginalError: nonNullSyncRootError,
	})
	nonNullSyncOriginalError2 := gqlerrors.FormatError(gqlerrors.Error{
		Message: nonNullSyncRootError.Error(),
		Locations: []location.SourceLocation{
			{Line: 19, Column: 19},
		},
		Path: []interface{}{
			"promiseNest", "nonNullNest", "nonNullPromiseNest", "nonNullNest",
			"nonNullPromiseNest", "nonNullSync",
		},
		OriginalError: nonNullSyncRootError,
	})
	nonNullPromiseError := errors.New(`Cannot return null for non-nullable field DataType.nonNullPromise.`)
	nonNullPromiseOriginalError := gqlerrors.FormatError(gqlerrors.Error{
		Message: nonNullPromiseError.Error(),
		Locations: []location.SourceLocation{
			{Line: 30, Column: 19},
		},
		Path: []interface{}{
			"anotherNest", "nonNullNest", "nonNullPromiseNest", "nonNullNest",
			"nonNullPromiseNest", "nonNullPromise",
		},
		OriginalError: nonNullPromiseError,
	})
	nonNullPromiseOriginalError2 := gqlerrors.FormatError(gqlerrors.Error{
		Message: nonNullPromiseError.Error(),
		Locations: []location.SourceLocation{
			{Line: 41, Column: 19},
		},
		Path: []interface{}{
			"anotherPromiseNest", "nonNullNest", "nonNullPromiseNest", "nonNullNest",
			"nonNullPromiseNest", "nonNullPromise",
		},
		OriginalError: nonNullPromiseError,
	})
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"nest":               nil,
			"promiseNest":        nil,
			"anotherNest":        nil,
			"anotherPromiseNest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormatError(gqlerrors.Error{
				Message:       nonNullSyncOriginalError.Message,
				Locations:     nonNullSyncOriginalError.Locations,
				Path:          nonNullSyncOriginalError.Path,
				OriginalError: nonNullSyncOriginalError,
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message:       nonNullSyncOriginalError2.Message,
				Locations:     nonNullSyncOriginalError2.Locations,
				Path:          nonNullSyncOriginalError2.Path,
				OriginalError: nonNullSyncOriginalError2,
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message:       nonNullPromiseOriginalError.Message,
				Locations:     nonNullPromiseOriginalError.Locations,
				Path:          nonNullPromiseOriginalError.Path,
				OriginalError: nonNullPromiseOriginalError,
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message:       nonNullPromiseOriginalError2.Message,
				Locations:     nonNullPromiseOriginalError2.Locations,
				Path:          nonNullPromiseOriginalError2.Path,
				OriginalError: nonNullPromiseOriginalError2,
			}),
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected.Data, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Data, result.Data))
	}
	sort.Sort(gqlerrors.FormattedErrors(expected.Errors))
	sort.Sort(gqlerrors.FormattedErrors(result.Errors))
	if !reflect.DeepEqual(expected.Errors, result.Errors) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
}

func TestNonNull_NullsTheTopLevelIfSyncNonNullableFieldThrows(t *testing.T) {
	doc := `
      query Q { nonNullSync }
	`
	originalError := errors.New(nonNullSyncError)
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormatError(gqlerrors.Error{
				Message: originalError.Error(),
				Locations: []location.SourceLocation{
					{Line: 2, Column: 17},
				},
				Path: []interface{}{
					"nonNullSync",
				},
				OriginalError: originalError,
			}),
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
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
	originalError := errors.New(nonNullPromiseError)
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormatError(gqlerrors.Error{
				Message: originalError.Error(),
				Locations: []location.SourceLocation{
					{Line: 2, Column: 17},
				},
				Path: []interface{}{
					"nonNullPromise",
				},
				OriginalError: originalError,
			}),
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
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
	rootError := errors.New(`Cannot return null for non-nullable field DataType.nonNullSync.`)
	originalError := gqlerrors.FormatError(gqlerrors.Error{
		Message: rootError.Error(),
		Locations: []location.SourceLocation{
			{Line: 2, Column: 17},
		},
		Path: []interface{}{
			"nonNullSync",
		},
		OriginalError: rootError,
	})
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormatError(gqlerrors.Error{
				Message:       originalError.Message,
				Locations:     originalError.Locations,
				Path:          originalError.Path,
				OriginalError: originalError,
			}),
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
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
	rootError := errors.New(`Cannot return null for non-nullable field DataType.nonNullPromise.`)
	originalError := gqlerrors.FormatError(gqlerrors.Error{
		Message: rootError.Error(),
		Locations: []location.SourceLocation{
			{Line: 2, Column: 17},
		},
		Path: []interface{}{
			"nonNullPromise",
		},
		OriginalError: rootError,
	})
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormatError(gqlerrors.Error{
				Message:       originalError.Message,
				Locations:     originalError.Locations,
				Path:          originalError.Path,
				OriginalError: originalError,
			}),
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
