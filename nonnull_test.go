package graphql

import (
	"reflect"
	"sort"
	"testing"

	"github.com/chris-ramon/graphql/gqlerrors"
	"github.com/chris-ramon/graphql/language/location"
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

var dataType = NewObject(ObjectConfig{
	Name: "DataType",
	Fields: FieldConfigMap{
		"sync": &FieldConfig{
			Type: String,
		},
		"nonNullSync": &FieldConfig{
			Type: NewNonNull(String),
		},
		"promise": &FieldConfig{
			Type: String,
		},
		"nonNullPromise": &FieldConfig{
			Type: NewNonNull(String),
		},
	},
})

var nonNullTestSchema, _ = NewSchema(SchemaConfig{
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

	dataType.AddFieldConfig("nest", &FieldConfig{
		Type: dataType,
	})
	dataType.AddFieldConfig("nonNullNest", &FieldConfig{
		Type: NewNonNull(dataType),
	})
	dataType.AddFieldConfig("promiseNest", &FieldConfig{
		Type: dataType,
	})
	dataType.AddFieldConfig("nonNullPromiseNest", &FieldConfig{
		Type: NewNonNull(dataType),
	})
}

// nulls a nullable field that panics
func TestNonNull_NullsANullableFieldThatThrowsSynchronously(t *testing.T) {
	doc := `
      query Q {
        sync
      }
	`
	expected := &Result{
		Data: map[string]interface{}{
			"sync": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
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
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestNonNull_NullsANullableFieldThatThrowsInAPromise(t *testing.T) {
	doc := `
      query Q {
        promise
      }
	`
	expected := &Result{
		Data: map[string]interface{}{
			"promise": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
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
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
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
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
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
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: map[string]interface{}{
			"promiseNest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
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
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: map[string]interface{}{
			"promiseNest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
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
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
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
			gqlerrors.FormattedError{
				Message: syncError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 4, Column: 11},
				},
			},
			gqlerrors.FormattedError{
				Message: syncError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 7, Column: 13},
				},
			},
			gqlerrors.FormattedError{
				Message: syncError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 11, Column: 13},
				},
			},
			gqlerrors.FormattedError{
				Message: syncError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 16, Column: 11},
				},
			},
			gqlerrors.FormattedError{
				Message: syncError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 19, Column: 13},
				},
			},
			gqlerrors.FormattedError{
				Message: syncError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 23, Column: 13},
				},
			},
			gqlerrors.FormattedError{
				Message: promiseError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 5, Column: 11},
				},
			},
			gqlerrors.FormattedError{
				Message: promiseError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 8, Column: 13},
				},
			},
			gqlerrors.FormattedError{
				Message: promiseError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 12, Column: 13},
				},
			},
			gqlerrors.FormattedError{
				Message: promiseError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 17, Column: 11},
				},
			},
			gqlerrors.FormattedError{
				Message: promiseError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 20, Column: 13},
				},
			},
			gqlerrors.FormattedError{
				Message: promiseError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 24, Column: 13},
				},
			},
		},
	}
	// parse query
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected.Data, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected.Data, result.Data))
	}
	sort.Sort(gqlerrors.GQLFormattedErrorSlice(expected.Errors))
	sort.Sort(gqlerrors.GQLFormattedErrorSlice(result.Errors))
	if !reflect.DeepEqual(expected.Errors, result.Errors) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected.Errors, result.Errors))
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
	expected := &Result{
		Data: map[string]interface{}{
			"nest":               nil,
			"promiseNest":        nil,
			"anotherNest":        nil,
			"anotherPromiseNest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: nonNullSyncError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 8, Column: 19},
				},
			},
			gqlerrors.FormattedError{
				Message: nonNullSyncError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 19, Column: 19},
				},
			},
			gqlerrors.FormattedError{
				Message: nonNullPromiseError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 30, Column: 19},
				},
			},
			gqlerrors.FormattedError{
				Message: nonNullPromiseError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 41, Column: 19},
				},
			},
		},
	}
	// parse query
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected.Data, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected.Data, result.Data))
	}
	sort.Sort(gqlerrors.GQLFormattedErrorSlice(expected.Errors))
	sort.Sort(gqlerrors.GQLFormattedErrorSlice(result.Errors))
	if !reflect.DeepEqual(expected.Errors, result.Errors) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected.Errors, result.Errors))
	}

}
func TestNonNull_NullsANullableFieldThatSynchronouslyReturnsNull(t *testing.T) {
	doc := `
      query Q {
        sync
      }
	`
	expected := &Result{
		Data: map[string]interface{}{
			"sync": nil,
		},
	}
	// parse query
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected.Data, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected.Data, result.Data))
	}
	if !reflect.DeepEqual(expected.Errors, result.Errors) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected.Errors, result.Errors))
	}
}
func TestNonNull_NullsANullableFieldThatSynchronouslyReturnsNullInAPromise(t *testing.T) {
	doc := `
      query Q {
        promise
      }
	`
	expected := &Result{
		Data: map[string]interface{}{
			"promise": nil,
		},
	}
	// parse query
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected.Data, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected.Data, result.Data))
	}
	if !reflect.DeepEqual(expected.Errors, result.Errors) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected.Errors, result.Errors))
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
	expected := &Result{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: `Cannot return null for non-nullable field DataType.nonNullSync.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 4, Column: 11},
				},
			},
		},
	}
	// parse query
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: `Cannot return null for non-nullable field DataType.nonNullPromise.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 4, Column: 11},
				},
			},
		},
	}
	// parse query
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: map[string]interface{}{
			"promiseNest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: `Cannot return null for non-nullable field DataType.nonNullSync.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 4, Column: 11},
				},
			},
		},
	}
	// parse query
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
		Data: map[string]interface{}{
			"promiseNest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: `Cannot return null for non-nullable field DataType.nonNullPromise.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 4, Column: 11},
				},
			},
		},
	}
	// parse query
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
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
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected.Data, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected.Data, result.Data))
	}
	if !reflect.DeepEqual(expected.Errors, result.Errors) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected.Errors, result.Errors))
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
	expected := &Result{
		Data: map[string]interface{}{
			"nest":               nil,
			"promiseNest":        nil,
			"anotherNest":        nil,
			"anotherPromiseNest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: `Cannot return null for non-nullable field DataType.nonNullSync.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 8, Column: 19},
				},
			},
			gqlerrors.FormattedError{
				Message: `Cannot return null for non-nullable field DataType.nonNullSync.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 19, Column: 19},
				},
			},
			gqlerrors.FormattedError{
				Message: `Cannot return null for non-nullable field DataType.nonNullPromise.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 30, Column: 19},
				},
			},
			gqlerrors.FormattedError{
				Message: `Cannot return null for non-nullable field DataType.nonNullPromise.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 41, Column: 19},
				},
			},
		},
	}
	// parse query
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected.Data, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected.Data, result.Data))
	}
	sort.Sort(gqlerrors.GQLFormattedErrorSlice(expected.Errors))
	sort.Sort(gqlerrors.GQLFormattedErrorSlice(result.Errors))
	if !reflect.DeepEqual(expected.Errors, result.Errors) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected.Errors, result.Errors))
	}
}

func TestNonNull_NullsTheTopLevelIfSyncNonNullableFieldThrows(t *testing.T) {
	doc := `
      query Q { nonNullSync }
	`
	expected := &Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: nonNullSyncError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 2, Column: 17},
				},
			},
		},
	}
	// parse query
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestNonNull_NullsTheTopLevelIfSyncNonNullableFieldErrors(t *testing.T) {
	doc := `
      query Q { nonNullPromise }
	`
	expected := &Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: nonNullPromiseError,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 2, Column: 17},
				},
			},
		},
	}
	// parse query
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestNonNull_NullsTheTopLevelIfSyncNonNullableFieldReturnsNull(t *testing.T) {
	doc := `
      query Q { nonNullSync }
	`
	expected := &Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: `Cannot return null for non-nullable field DataType.nonNullSync.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 2, Column: 17},
				},
			},
		},
	}
	// parse query
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestNonNull_NullsTheTopLevelIfSyncNonNullableFieldResolvesNull(t *testing.T) {
	doc := `
      query Q { nonNullPromise }
	`
	expected := &Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message: `Cannot return null for non-nullable field DataType.nonNullPromise.`,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 2, Column: 17},
				},
			},
		},
	}
	// parse query
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
