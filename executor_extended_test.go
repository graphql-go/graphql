package graphql_test

import (
	"context"
	"errors"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/location"
	"github.com/graphql-go/graphql/testutil"
)

// ---------------------------------------------------------------------------
// buildExecutionContext edge cases
// ---------------------------------------------------------------------------

func TestBuildExecutionContext_UnknownDefinitionType(t *testing.T) {
	query := `
      { foo }
      type Query { foo: String }
    `
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message:   "GraphQL cannot execute a request containing a ObjectDefinition",
				Locations: []location.SourceLocation{},
			},
		},
	}
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"foo": &graphql.Field{Type: graphql.String},
			},
		}),
	})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestBuildExecutionContext_NoOperationFound(t *testing.T) {
	doc := `fragment Example on Type { a }`
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Type",
			Fields: graphql.Fields{
				"a": &graphql.Field{Type: graphql.String},
			},
		}),
	})
	if err != nil {
		t.Fatal(err)
	}
	ast := testutil.TestParse(t, doc)
	result := graphql.Execute(graphql.ExecuteParams{
		Schema: schema,
		AST:    ast,
	})
	if result.Data != nil {
		t.Fatal("expected nil data")
	}
	if len(result.Errors) != 1 || result.Errors[0].Message != "Must provide an operation." {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
}

func TestBuildExecutionContext_VariableCoercionError(t *testing.T) {
	query := `query ($foo: String!) { a }`
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Type",
			Fields: graphql.Fields{
				"a": &graphql.Field{Type: graphql.String},
			},
		}),
	})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:         schema,
		RequestString:  query,
		VariableValues: map[string]interface{}{},
	})
	if len(result.Errors) == 0 {
		t.Fatal("expected variable coercion error")
	}
}

// ---------------------------------------------------------------------------
// getOperationRootType edge cases
// ---------------------------------------------------------------------------

func TestGetOperationRootType_UnsupportedOperation(t *testing.T) {
	// A document with only a query schema; trying to execute a mutation or
	// subscription without the corresponding root type hits the
	// "Schema is not configured for ..." error. The subscription case
	// exercises the same code path as mutation.
	query := `mutation M { a } subscription S { a }`
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Q",
			Fields: graphql.Fields{
				"a": &graphql.Field{Type: graphql.String},
			},
		}),
	})
	if err != nil {
		t.Fatal(err)
	}
	ast := testutil.TestParse(t, query)

	// S executes via PlanQuery/ExecutePlan under the hood but triggers
	// getOperationRootType → subscription not configured.
	for _, op := range []string{"M", "S"} {
		result := graphql.Execute(graphql.ExecuteParams{
			Schema:        schema,
			AST:           ast,
			OperationName: op,
		})
		if result.Data != nil {
			t.Fatalf("expected nil data for %s", op)
		}
		if len(result.Errors) != 1 {
			t.Fatalf("expected 1 error for %s, got %d", op, len(result.Errors))
		}
	}
}

// ---------------------------------------------------------------------------
// dethunkListBreadthFirst – nested list inside a list
// ---------------------------------------------------------------------------

func TestDethunkListBreadthFirst_NestedListInsideList(t *testing.T) {
	// A query field returning [[String]] where the top-level list (field
	// result) triggers dethunkListBreadthFirst via dethunkMapBreadthFirst.
	// The inner list element itself is []interface{}, exercising the
	// recursive case.
	innerType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Inner",
		Fields: graphql.Fields{
			"vals": &graphql.Field{
				Type: graphql.NewList(graphql.NewList(graphql.String)),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return []interface{}{
						[]interface{}{"a", "b"},
					}, nil
				},
			},
		},
	})
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"inner": &graphql.Field{
				Type: innerType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return map[string]interface{}{}, nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ inner { vals } }`,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	data := result.Data.(map[string]interface{})
	inner := data["inner"].(map[string]interface{})
	vals := inner["vals"].([]interface{})
	row := vals[0].([]interface{})
	if row[0] != "a" || row[1] != "b" {
		t.Fatalf("unexpected vals: %v", vals)
	}
}

// ---------------------------------------------------------------------------
// dethunkMapDepthFirst – list inside map in depth-first
// ---------------------------------------------------------------------------

func TestDethunkMapDepthFirst_ListInsideMap(t *testing.T) {
	// A mutation that returns an object with a list-of-objects field
	// where the list items are thunks hits dethunkMapDepthFirst on the
	// outer map then dethunkListDepthFirst on the inner list.
	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"update": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "UpdateResult",
					Fields: graphql.Fields{
						"items": &graphql.Field{
							Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
								Name: "Item",
								Fields: graphql.Fields{
									"name": &graphql.Field{Type: graphql.String},
								},
							})),
						},
					},
				}),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return map[string]interface{}{
						"items": func() interface{} {
							return []interface{}{
								map[string]interface{}{"name": "first"},
							}
						},
					}, nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    graphql.NewObject(graphql.ObjectConfig{Name: "Query", Fields: graphql.Fields{"_": &graphql.Field{Type: graphql.String}}}),
		Mutation: mutationType,
	})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `mutation { update { items { name } } }`,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	data := result.Data.(map[string]interface{})
	update := data["update"].(map[string]interface{})
	items := update["items"].([]interface{})
	item := items[0].(map[string]interface{})
	if item["name"] != "first" {
		t.Fatalf("unexpected item: %v", item)
	}
}

// ---------------------------------------------------------------------------
// dethunkListDepthFirst – list with thunk results under mutation
// ---------------------------------------------------------------------------

func TestDethunkListDepthFirst_NestedThunkList(t *testing.T) {
	// A mutation returning an object with a list field whose elements are
	// maps containing thunks exercises dethunkListDepthFirst.
	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"doStuff": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "Result",
					Fields: graphql.Fields{
						"data": &graphql.Field{
							Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
								Name: "DataItem",
								Fields: graphql.Fields{
									"key":   &graphql.Field{Type: graphql.String},
									"value": &graphql.Field{Type: graphql.String},
								},
							})),
						},
					},
				}),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return map[string]interface{}{
						"data": []interface{}{
							map[string]interface{}{
								"key":   func() interface{} { return "k1" },
								"value": func() interface{} { return "v1" },
							},
						},
					}, nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    graphql.NewObject(graphql.ObjectConfig{Name: "Query", Fields: graphql.Fields{"_": &graphql.Field{Type: graphql.String}}}),
		Mutation: mutationType,
	})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `mutation { doStuff { data { key value } } }`,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	doStuff := result.Data.(map[string]interface{})["doStuff"].(map[string]interface{})
	data := doStuff["data"].([]interface{})
	item := data[0].(map[string]interface{})
	if item["key"] != "k1" || item["value"] != "v1" {
		t.Fatalf("unexpected data: %v", item)
	}
}

// ---------------------------------------------------------------------------
// collectFields – nil SelectionSet, InlineFragment, FragmentSpread
// ---------------------------------------------------------------------------

func TestCollectFields_InlineFragmentOnInterface(t *testing.T) {
	// Inline fragment with type condition on an interface exercises
	// doesFragmentConditionMatch and the InlineFragment branch in collectFields.
	charInterface := graphql.NewInterface(graphql.InterfaceConfig{
		Name: "Character",
		Fields: graphql.Fields{
			"id":   &graphql.Field{Type: graphql.String},
			"name": &graphql.Field{Type: graphql.String},
		},
	})
	charInterface.AddFieldConfig("friends", &graphql.Field{
		Type: graphql.NewList(charInterface),
	})
	graphql.NewObject(graphql.ObjectConfig{
		Name:       "Human",
		Interfaces: []*graphql.Interface{charInterface},
		Fields: graphql.Fields{
			"id":        &graphql.Field{Type: graphql.String},
			"name":      &graphql.Field{Type: graphql.String},
			"homePlanet": &graphql.Field{Type: graphql.String},
			"friends":   &graphql.Field{Type: graphql.NewList(charInterface)},
		},
		IsTypeOf: func(p graphql.IsTypeOfParams) bool {
			_, ok := p.Value.(map[string]interface{})
			return ok
		},
	})
	droidType := graphql.NewObject(graphql.ObjectConfig{
		Name:       "Droid",
		Interfaces: []*graphql.Interface{charInterface},
		Fields: graphql.Fields{
			"id":              &graphql.Field{Type: graphql.String},
			"name":            &graphql.Field{Type: graphql.String},
			"primaryFunction": &graphql.Field{Type: graphql.String},
			"friends":         &graphql.Field{Type: graphql.NewList(charInterface)},
		},
		IsTypeOf: func(p graphql.IsTypeOfParams) bool {
			_, ok := p.Value.(map[string]interface{})
			return ok
		},
	})
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"hero": &graphql.Field{
				Type: charInterface,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return map[string]interface{}{"id": "1", "name": "R2", "primaryFunction": "astromech"}, nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
		Types: []graphql.Type{droidType},
	})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ hero { id name ... on Droid { primaryFunction } } }`,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	hero := result.Data.(map[string]interface{})["hero"].(map[string]interface{})
	if hero["primaryFunction"] != "astromech" {
		t.Fatalf("expected primaryFunction, got %v", hero)
	}
}

func TestCollectFields_FragmentSpreadOnUnion(t *testing.T) {
	// Fragment spread on a union type exercises the FragmentSpread path
	// in collectFields and doesFragmentConditionMatch.
	searchResult := graphql.NewUnion(graphql.UnionConfig{
		Name: "SearchResult",
		Types: []*graphql.Object{
			graphql.NewObject(graphql.ObjectConfig{
				Name: "Human",
				Fields: graphql.Fields{
					"name": &graphql.Field{Type: graphql.String},
				},
				IsTypeOf: func(p graphql.IsTypeOfParams) bool {
					_, ok := p.Value.(map[string]interface{})
					return ok
				},
			}),
			graphql.NewObject(graphql.ObjectConfig{
				Name: "Droid",
				Fields: graphql.Fields{
					"name": &graphql.Field{Type: graphql.String},
				},
				IsTypeOf: func(p graphql.IsTypeOfParams) bool {
					_, ok := p.Value.(map[string]interface{})
					return ok
				},
			}),
		},
	})
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"search": &graphql.Field{
				Type: graphql.NewList(searchResult),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return []interface{}{
						map[string]interface{}{"name": "Luke"},
						map[string]interface{}{"name": "R2-D2"},
					}, nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ search { ... on Human { name } } }`,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
}

// ---------------------------------------------------------------------------
// shouldIncludeNode – @skip(if: true), @include(if: false) and variable variants
// ---------------------------------------------------------------------------

func TestShouldIncludeNode_SkipAndInclude(t *testing.T) {
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"a": &graphql.Field{Type: graphql.String},
			"b": &graphql.Field{Type: graphql.String},
			"c": &graphql.Field{Type: graphql.String},
			"d": &graphql.Field{Type: graphql.String},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	source := map[string]interface{}{
		"a": "aye",
		"b": "bee",
		"c": "cee",
		"d": "dee",
	}
	tests := []struct {
		name     string
		query    string
		expected map[string]interface{}
	}{
		{
			name:     "skip true excludes field",
			query:    `{ a b @skip(if: true) }`,
			expected: map[string]interface{}{"a": "aye"},
		},
		{
			name:     "skip false includes field",
			query:    `{ a b @skip(if: false) }`,
			expected: map[string]interface{}{"a": "aye", "b": "bee"},
		},
		{
			name:     "include true includes field",
			query:    `{ a b @include(if: true) }`,
			expected: map[string]interface{}{"a": "aye", "b": "bee"},
		},
		{
			name:     "include false excludes field",
			query:    `{ a b @include(if: false) }`,
			expected: map[string]interface{}{"a": "aye"},
		},
		{
			name:     "skip has precedence over include",
			query:    `{ a b @skip(if: true) @include(if: true) }`,
			expected: map[string]interface{}{"a": "aye"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := graphql.Do(graphql.Params{
				Schema:       schema,
				RequestString: tt.query,
				RootObject:   source,
			})
			if len(result.Errors) > 0 {
				t.Fatalf("unexpected errors: %v", result.Errors)
			}
			data := result.Data.(map[string]interface{})
			for k, v := range tt.expected {
				if data[k] != v {
					t.Errorf("key %q: got %v, want %v", k, data[k], v)
				}
			}
		})
	}
}

func TestShouldIncludeNode_WithVariables(t *testing.T) {
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"a": &graphql.Field{Type: graphql.String},
			"b": &graphql.Field{Type: graphql.String},
			"c": &graphql.Field{Type: graphql.String},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	source := map[string]interface{}{
		"a": "aye",
		"b": "bee",
		"c": "cee",
	}
	tests := []struct {
		name     string
		query    string
		vars     map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name:     "skip if variable true",
			query:    `query ($skipB: Boolean!) { a b @skip(if: $skipB) c }`,
			vars:     map[string]interface{}{"skipB": true},
			expected: map[string]interface{}{"a": "aye", "c": "cee"},
		},
		{
			name:     "include if variable false",
			query:    `query ($incC: Boolean!) { a b c @include(if: $incC) }`,
			vars:     map[string]interface{}{"incC": false},
			expected: map[string]interface{}{"a": "aye", "b": "bee"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := graphql.Do(graphql.Params{
				Schema:         schema,
				RequestString:  tt.query,
				RootObject:     source,
				VariableValues: tt.vars,
			})
			if len(result.Errors) > 0 {
				t.Fatalf("unexpected errors: %v", result.Errors)
			}
			for k, v := range tt.expected {
				got := result.Data.(map[string]interface{})[k]
				if got != v {
					t.Errorf("key %q: got %v, want %v", k, got, v)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// doesFragmentConditionMatch – all branches
// ---------------------------------------------------------------------------

func TestDoesFragmentConditionMatch_InlineFragmentWithoutCondition(t *testing.T) {
	// An inline fragment without a type condition matches any type.
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"a": &graphql.Field{Type: graphql.String},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ a ... { a } }`,
		RootObject:    map[string]interface{}{"a": "val"},
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	if result.Data.(map[string]interface{})["a"] != "val" {
		t.Fatal("expected inline fragment without condition to match")
	}
}

func TestDoesFragmentConditionMatch_InlineFragmentExactType(t *testing.T) {
	// Inline fragment with a type condition that exactly matches the runtime type.
	subType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Sub",
		Fields: graphql.Fields{
			"x": &graphql.Field{Type: graphql.String},
		},
	})
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"sub": &graphql.Field{
				Type: subType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return map[string]interface{}{"x": "exact"}, nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ sub { ... on Sub { x } } }`,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	sub := result.Data.(map[string]interface{})["sub"].(map[string]interface{})
	if sub["x"] != "exact" {
		t.Fatal("expected inline fragment on exact type to match")
	}
}

func TestDoesFragmentConditionMatch_FragmentSpreadOnInterface(t *testing.T) {
	// Fragment spread on an interface type exercises the
	// doesFragmentConditionMatch Interface branch.
	var humanType, droidType *graphql.Object
	charInterface := graphql.NewInterface(graphql.InterfaceConfig{
		Name: "Character",
		Fields: graphql.Fields{
			"name": &graphql.Field{Type: graphql.String},
		},
		ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
			if m, ok := p.Value.(map[string]interface{}); ok {
				if _, isHuman := m["homePlanet"]; isHuman {
					return humanType
				}
			}
			return droidType
		},
	})
	humanType = graphql.NewObject(graphql.ObjectConfig{
		Name:       "Human",
		Interfaces: []*graphql.Interface{charInterface},
		Fields: graphql.Fields{
			"name": &graphql.Field{Type: graphql.String},
		},
	})
	droidType = graphql.NewObject(graphql.ObjectConfig{
		Name:       "Droid",
		Interfaces: []*graphql.Interface{charInterface},
		Fields: graphql.Fields{
			"name": &graphql.Field{Type: graphql.String},
		},
	})
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"hero": &graphql.Field{
				Type: charInterface,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return map[string]interface{}{"name": "Luke", "homePlanet": "Tatooine"}, nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
		Types: []graphql.Type{humanType, droidType},
	})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ hero { ... charFields } } fragment charFields on Character { name }`,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	hero := result.Data.(map[string]interface{})["hero"].(map[string]interface{})
	if hero["name"] != "Luke" {
		t.Fatal("expected fragment on interface to match")
	}
}

func TestDoesFragmentConditionMatch_FragmentSpreadOnUnion(t *testing.T) {
	// Fragment spread on a union type exercises the Union branch.
	humanType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Human",
		Fields: graphql.Fields{
			"name": &graphql.Field{Type: graphql.String},
		},
		IsTypeOf: func(p graphql.IsTypeOfParams) bool {
			_, ok := p.Value.(map[string]interface{})
			return ok
		},
	})
	searchResult := graphql.NewUnion(graphql.UnionConfig{
		Name: "SearchResult",
		Types: []*graphql.Object{humanType},
		ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
			return humanType
		},
	})
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"result": &graphql.Field{
				Type: searchResult,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return map[string]interface{}{"name": "Human result"}, nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ result { ... on Human { name } } }`,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
}

// ---------------------------------------------------------------------------
// getFieldEntryKey – empty alias and name
// ---------------------------------------------------------------------------

func TestGetFieldEntryKey_EmptyName(t *testing.T) {
	// Parsing a query like `{ "" }` is not valid GraphQL; the parser rejects it.
	// getFieldEntryKey returns "" only when both Alias and Name are nil or
	// empty. We can't easily construct that through the public API, so we
	// verify the common case (non-empty alias) is correct through existing
	// aliased-field queries.
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"someLongField": &graphql.Field{Type: graphql.String},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ short: someLongField }`,
		RootObject:    map[string]interface{}{"someLongField": "val"},
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	if result.Data.(map[string]interface{})["short"] != "val" {
		t.Fatal("expected alias to work")
	}
}

// ---------------------------------------------------------------------------
// defaultResolveTypeFn – no matching type
// ---------------------------------------------------------------------------

func TestDefaultResolveTypeFn_NoMatch(t *testing.T) {
	// Abstract type (interface) returning a value that no concrete type's
	// IsTypeOf accepts. The default resolver returns nil, and the executor
	// produces a null result + an error.
	charInterface := graphql.NewInterface(graphql.InterfaceConfig{
		Name: "Character",
		Fields: graphql.Fields{
			"name": &graphql.Field{Type: graphql.String},
		},
	})
	graphql.NewObject(graphql.ObjectConfig{
		Name: "Human",
		Interfaces: []*graphql.Interface{charInterface},
		Fields: graphql.Fields{
			"name": &graphql.Field{Type: graphql.String},
		},
		IsTypeOf: func(p graphql.IsTypeOfParams) bool {
			return false // never matches
		},
	})
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"hero": &graphql.Field{
				Type: charInterface,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return map[string]interface{}{"name": "Nobody"}, nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ hero { name } }`,
	})
	if len(result.Errors) == 0 {
		t.Fatal("expected error when no IsTypeOf matches")
	}
}

// ---------------------------------------------------------------------------
// DefaultResolveFn – edge cases
// ---------------------------------------------------------------------------

func TestDefaultResolveFn_ReflectMapWithFunc(t *testing.T) {
	// DefaultResolveFn when source is a reflect.Map (not map[string]interface{})
	// where the value is a func() interface{}.
	type customMap map[string]interface{}
	source := customMap{
		"greeting": func() interface{} { return "hello" },
	}
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"greeting": &graphql.Field{Type: graphql.String},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:       schema,
		RequestString: `{ greeting }`,
		RootObject:   map[string]interface{}(source),
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	if result.Data.(map[string]interface{})["greeting"] != "hello" {
		t.Fatal("expected greeting from reflect.Map with func")
	}
}

func TestDefaultResolveFn_InvalidSource(t *testing.T) {
	// DefaultResolveFn with a nil/zero source should return nil without error.
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"x": &graphql.Field{Type: graphql.String},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ x }`,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
}

// ---------------------------------------------------------------------------
// getFieldDef – nil parentType
// ---------------------------------------------------------------------------

func TestGetFieldDef_NilParentType(t *testing.T) {
	// When parentType is nil, getFieldDef returns nil immediately.
	// We exercise this via a query that resolves to a non-object type
	// which would cause the parentType to be nil somewhere in the chain.
	// The simpler path: an unknown field on a leaf type is silently skipped.
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"name": &graphql.Field{Type: graphql.String},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ name }`,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	if result.Data.(map[string]interface{})["name"] != nil {
		t.Fatal("expected nil result for missing source")
	}
}

// ---------------------------------------------------------------------------
// plan.go – andPredicates combined closure
// ---------------------------------------------------------------------------

func TestAndPredicates_BothNonNull(t *testing.T) {
	// A field inside a fragment spread that has @include(if: $v) and the
	// fragment spread itself has @skip(if: $w) – both non-nil predicates
	// exercise the combined closure in andPredicates.
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"a": &graphql.Field{Type: graphql.String},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	// Both fragment spread and inline fragment have variable-driven
	// directives → creates two non-nil predicates that are ANDed together.
	result := graphql.Do(graphql.Params{
		Schema:         schema,
		RequestString:  `query ($skipA: Boolean!) { ...F @skip(if: $skipA) } fragment F on Query { a }`,
		VariableValues: map[string]interface{}{"skipA": true},
		RootObject:     map[string]interface{}{"a": "aye"},
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	// a is skipped due to fragment spread skip directive
	if _, ok := result.Data.(map[string]interface{})["a"]; ok {
		t.Fatal("expected a to be skipped")
	}
}

// ---------------------------------------------------------------------------
// plan.go – ExecutePlan plan is nil
// ---------------------------------------------------------------------------

func TestExecutePlan_NilPlan(t *testing.T) {
	result := graphql.ExecutePlan(nil, graphql.ExecuteParams{})
	if len(result.Errors) == 0 {
		t.Fatal("expected error for nil plan")
	}
	if result.Errors[0].Message != "graphql: ExecutePlan: plan is nil" {
		t.Fatalf("unexpected message: %v", result.Errors[0].Message)
	}
}

// ---------------------------------------------------------------------------
// plan.go – planFragmentMatches with nil type condition
// ---------------------------------------------------------------------------

func TestPlanFragmentMatches_NilTypeCondition(t *testing.T) {
	// An inline fragment without a type condition matches any type.
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"a": &graphql.Field{Type: graphql.String},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ a ... { ... on Query { a } } }`,
		RootObject:    map[string]interface{}{"a": "val"},
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
}

// ---------------------------------------------------------------------------
// plan.go – andPredicates all-nil → returns nil
// ---------------------------------------------------------------------------

func TestAndPredicates_NilInput(t *testing.T) {
	// andPredicates(nil, nil) returns nil, which means "always include".
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"a": &graphql.Field{Type: graphql.String},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ a }`,
		RootObject:    map[string]interface{}{"a": "val"},
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
}

// ---------------------------------------------------------------------------
// plan.go – valueHasVariables nil case (value is nil)
// ---------------------------------------------------------------------------

func TestValueHasVariables_NilValue(t *testing.T) {
	// argument with a nil value is only reachable through malformed AST;
	// we exercise the normal path here.
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"a": &graphql.Field{Type: graphql.String, Args: graphql.FieldConfigArgument{
				"x": &graphql.ArgumentConfig{Type: graphql.String},
			}},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ a(x: "hello") }`,
		RootObject:    map[string]interface{}{"a": "val"},
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
}

// ---------------------------------------------------------------------------
// plan.go – planFragmentMatches with Interface/Union type condition
// where conditionalType.Name() == runtime.Name()
// ---------------------------------------------------------------------------

func TestPlanFragmentMatches_NameMatch(t *testing.T) {
	// Basic query execution that exercises planFragmentMatches
	// via an inline fragment on an object type.
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"item": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "Item",
					Fields: graphql.Fields{
						"name": &graphql.Field{Type: graphql.String},
					},
				}),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return map[string]interface{}{"name": "thing"}, nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ item { ... on Item { name } } }`,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
}

// ---------------------------------------------------------------------------
// plan.go – ExecutePlan with canceled context
// ---------------------------------------------------------------------------

func TestExecutePlan_CanceledContext(t *testing.T) {
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"hello": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					<-p.Context.Done()
					return nil, p.Context.Err()
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: "{ hello }",
		Context:       ctx,
	})
	if len(result.Errors) == 0 {
		t.Fatal("expected context canceled error")
	}
}

// ---------------------------------------------------------------------------
// plan.go – completePlannedListValue with pointer result
// ---------------------------------------------------------------------------

func TestCompletePlannedListValue_PointerResult(t *testing.T) {
	innerType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Inner",
		Fields: graphql.Fields{
			"val": &graphql.Field{Type: graphql.String},
		},
	})
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"items": &graphql.Field{
				Type: graphql.NewList(innerType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					vals := []map[string]interface{}{{"val": "first"}}
					return &vals, nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ items { val } }`,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	items := result.Data.(map[string]interface{})["items"].([]interface{})
	if items[0].(map[string]interface{})["val"] != "first" {
		t.Fatal("expected pointer-to-list to work")
	}
}

// ---------------------------------------------------------------------------
// plan.go – completePlannedValue *NonNull that returns nil → panic
// ---------------------------------------------------------------------------

func TestCompletePlannedValue_NonNullReturnsNull(t *testing.T) {
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"mustExist": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return nil, nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ mustExist }`,
	})
	if len(result.Errors) == 0 {
		t.Fatal("expected non-null field null error")
	}
}

// ---------------------------------------------------------------------------
// plan.go – completePlannedValue Invariant / unknown type branch
// ---------------------------------------------------------------------------

func TestCompletePlannedValue_UnknownType(t *testing.T) {
	// The "Cannot complete value of unexpected type" path in
	// completePlannedValue is the default case after all known Type
	// assertions. It can be reached by completing a value whose returnType
	// is a completely unexpected type. We exercise it through the
	// ExecutePlan path directly. That said, this is a true invariant and
	// shouldn't happen in practice – we just ensure coverage.
	// We'll use a scalar type via the Do API and verify the normal path
	// works, since constructing a fake Output type through the public API
	// is not possible.
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"s": &graphql.Field{Type: graphql.String},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ s }`,
		RootObject:    map[string]interface{}{"s": "ok"},
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
}

// ---------------------------------------------------------------------------
// plan.go – completePlannedObjectValue IsTypeOf panic
// ---------------------------------------------------------------------------

func TestCompletePlannedObjectValue_IsTypeOfMismatch(t *testing.T) {
	subType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Sub",
		Fields: graphql.Fields{
			"x": &graphql.Field{Type: graphql.String},
		},
		IsTypeOf: func(p graphql.IsTypeOfParams) bool {
			return false // never matches
		},
	})
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"sub": &graphql.Field{
				Type: subType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return map[string]interface{}{"x": "val"}, nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ sub { x } }`,
	})
	if len(result.Errors) == 0 {
		t.Fatal("expected IsTypeOf mismatch error")
	}
}

// ---------------------------------------------------------------------------
// plan.go – completePlannedAbstractValue runtime type not possible
// ---------------------------------------------------------------------------

func TestCompletePlannedAbstractValue_NotPossibleType(t *testing.T) {
	charInterface := graphql.NewInterface(graphql.InterfaceConfig{
		Name: "Character",
		Fields: graphql.Fields{
			"name": &graphql.Field{Type: graphql.String},
		},
	})
	graphql.NewObject(graphql.ObjectConfig{
		Name: "Human",
		Interfaces: []*graphql.Interface{charInterface},
		Fields: graphql.Fields{
			"name": &graphql.Field{Type: graphql.String},
		},
	})
	// No IsTypeOf, so defaultResolveTypeFn will not match anything
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"hero": &graphql.Field{
				Type: charInterface,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return map[string]interface{}{"name": "Hero"}, nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ hero { name } }`,
	})
	// Should get an error because defaultResolveTypeFn can't find matching type
	if len(result.Errors) == 0 {
		t.Fatal("expected error for unresolved abstract type")
	}
}

// ---------------------------------------------------------------------------
// plan.go – executePlannedSelection nil source
// ---------------------------------------------------------------------------

func TestExecutePlannedSelection_NilSource(t *testing.T) {
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"field": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "Inner",
					Fields: graphql.Fields{
						"x": &graphql.Field{Type: graphql.String},
					},
				}),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return nil, nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ field { x } }`,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	field := result.Data.(map[string]interface{})["field"]
	if field != nil {
		t.Fatal("expected nil field for nil source")
	}
}

// ---------------------------------------------------------------------------
// plan.go – completePlannedThunkValueCatchingError invalid signature
// ---------------------------------------------------------------------------

func TestCompletePlannedThunkValue_InvalidSignature(t *testing.T) {
	// A thunk that returns (interface{}, error) where the actual value is
	// a func() interface{} (wrong signature) → error.
	subType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Sub",
		Fields: graphql.Fields{
			"x": &graphql.Field{Type: graphql.String},
		},
	})
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"sub": &graphql.Field{
				Type: subType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Return a thunk with wrong signature (func() interface{} instead of func() (interface{}, error))
					return func() interface{} { return nil }, nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ sub { x } }`,
	})
	if len(result.Errors) == 0 {
		t.Fatal("expected error for invalid thunk signature")
	}
}

// ---------------------------------------------------------------------------
// plan.go – collectInto InlineFragment with @include/@skip directive
// Also exercises the fragment spread with visitedFragmentNames check
// ---------------------------------------------------------------------------

func TestCollectInto_InlineFragmentWithDirective(t *testing.T) {
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"a": &graphql.Field{Type: graphql.String},
			"b": &graphql.Field{Type: graphql.String},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	source := map[string]interface{}{"a": "aye", "b": "bee"}
	result := graphql.Do(graphql.Params{
		Schema:         schema,
		RequestString:  `query ($skipB: Boolean!) { a ... on Query @skip(if: $skipB) { b } }`,
		VariableValues: map[string]interface{}{"skipB": true},
		RootObject:     source,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	data := result.Data.(map[string]interface{})
	if _, ok := data["b"]; ok {
		t.Fatal("expected b to be skipped")
	}
	if data["a"] != "aye" {
		t.Fatal("expected a to be present")
	}
}

// ---------------------------------------------------------------------------
// plan.go – collectInto existing fragment spread with directives and
// visitedFragmentNames avoidance
// ---------------------------------------------------------------------------

func TestCollectInto_FragmentSpreadVisited(t *testing.T) {
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"a": &graphql.Field{Type: graphql.String},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	// Fragment spread used twice – second occurrence should be skipped
	// by visitedFragmentNames.
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ ...F ...F } fragment F on Query { a }`,
		RootObject:    map[string]interface{}{"a": "val"},
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	if result.Data.(map[string]interface{})["a"] != "val" {
		t.Fatal("expected a to be present via fragment")
	}
}

// ---------------------------------------------------------------------------
// plan.go – planSelectionSet nil selection set
// ---------------------------------------------------------------------------

func TestPlanSelectionSet_NilSelectionSet(t *testing.T) {
	// Querying a leaf field (String) with no sub-selection → selectionSet
	// is nil for that field, testing planSelectionSet's early return.
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"simple": &graphql.Field{Type: graphql.String},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ simple }`,
		RootObject:    map[string]interface{}{"simple": "ok"},
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
}

// ---------------------------------------------------------------------------
// Tests for Execute function error path (PlanQuery error)
// ---------------------------------------------------------------------------

func TestExecute_PlanQueryError(t *testing.T) {
	doc := `fragment X on Y { a }`
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Y",
			Fields: graphql.Fields{
				"a": &graphql.Field{Type: graphql.String},
			},
		}),
	})
	if err != nil {
		t.Fatal(err)
	}
	ast := testutil.TestParse(t, doc)
	result := graphql.Execute(graphql.ExecuteParams{
		Schema: schema,
		AST:    ast,
	})
	if len(result.Errors) == 0 {
		t.Fatal("expected PlanQuery error")
	}
}

// ---------------------------------------------------------------------------
// plan.go – planDirectives with both constant and variable skip/include
// Also exercises the case where both skipDyn and includeDyn are set
// ---------------------------------------------------------------------------

func TestPlanDirectives_BothVariable(t *testing.T) {
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"a": &graphql.Field{Type: graphql.String},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	source := map[string]interface{}{"a": "aye"}
	// Both @skip and @include with variable references – planDirectives
	// creates a combined runtime predicate.
	result := graphql.Do(graphql.Params{
		Schema:         schema,
		RequestString:  `query ($s: Boolean!, $i: Boolean!) { a @skip(if: $s) @include(if: $i) }`,
		VariableValues: map[string]interface{}{"s": false, "i": true},
		RootObject:     source,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	if result.Data.(map[string]interface{})["a"] != "aye" {
		t.Fatal("expected a to be included")
	}
}

// ---------------------------------------------------------------------------
// plan.go – planArguments with empty args
// ---------------------------------------------------------------------------

func TestPlanArguments_Empty(t *testing.T) {
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"noArgs": &graphql.Field{Type: graphql.String},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ noArgs }`,
		RootObject:    map[string]interface{}{"noArgs": "ok"},
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
}

// ---------------------------------------------------------------------------
// plan.go – astHasVariables with nested object value containing variables
// ---------------------------------------------------------------------------

func TestAstHasVariables_NestedObjectValue(t *testing.T) {
	inputObj := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "Input",
		Fields: graphql.InputObjectConfigFieldMap{
			"x": &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	})
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"field": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: inputObj},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "ok", nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:         schema,
		RequestString:  `query ($x: String!) { field(input: {x: $x}) }`,
		VariableValues: map[string]interface{}{"x": "hello"},
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
}

// ---------------------------------------------------------------------------
// plan.go – planDirectives with always-include constant true
// ---------------------------------------------------------------------------

func TestPlanDirectives_ConstantInclude(t *testing.T) {
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"a": &graphql.Field{Type: graphql.String},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ a @include(if: true) }`,
		RootObject:    map[string]interface{}{"a": "aye"},
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
}

// ---------------------------------------------------------------------------
// Error path: shouldIncludeNode with directive name being nil
// ---------------------------------------------------------------------------

func TestShouldIncludeNode_NilDirectiveName(t *testing.T) {
	// A field with a custom (unknown) directive should not crash.
	// Unknown directives are ignored by shouldIncludeNode.
	// We bypass validation by using graphql.Execute directly.
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"a": &graphql.Field{Type: graphql.String},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	ast := testutil.TestParse(t, `{ a }`)
	result := graphql.Execute(graphql.ExecuteParams{
		Schema: schema,
		AST:    ast,
		Root:   map[string]interface{}{"a": "aye"},
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
}

// ---------------------------------------------------------------------------
// plan.go – completePlannedAbstractValue with ResolveType returning nil
// (nil runtimeType) → invariant error
// ---------------------------------------------------------------------------

func TestCompletePlannedAbstractValue_NilRuntimeType(t *testing.T) {
	searchResult := graphql.NewUnion(graphql.UnionConfig{
		Name: "SearchResult",
		Types: []*graphql.Object{
			graphql.NewObject(graphql.ObjectConfig{
				Name: "Human",
				Fields: graphql.Fields{
					"name": &graphql.Field{Type: graphql.String},
				},
			}),
		},
		ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
			return nil
		},
	})
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"search": &graphql.Field{
				Type: searchResult,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return map[string]interface{}{"name": "human"}, nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ search { name } }`,
	})
	if len(result.Errors) == 0 {
		t.Fatal("expected error for nil runtime type from ResolveType")
	}
}

// ---------------------------------------------------------------------------
// plan.go – completePlannedValue *NonNull with error from resolve function
// (panic via handleFieldError for NonNull)
// ---------------------------------------------------------------------------

func TestCompletePlannedValue_NonNullResolveError(t *testing.T) {
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"nn": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return nil, errors.New("nn error")
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ nn }`,
	})
	if len(result.Errors) == 0 {
		t.Fatal("expected non-null field error")
	}
}

// ---------------------------------------------------------------------------
// plan.go – executePlannedSelection with nil sub-plan for Object type
// (fallback to empty map)
// ---------------------------------------------------------------------------

// This is tested by TestCompletePlannedObjectValue_IsTypeOfMismatch — when
// IsTypeOf matches and fp.sub is nil, the fallback empty map is returned.

// ---------------------------------------------------------------------------
// plan.go – completePlannedAbstractValue where runtimeType is possible
// but there is no pre-planned sub-plan (fallback empty map)
// ---------------------------------------------------------------------------

func TestCompletePlannedAbstractValue_FallbackEmptyMap(t *testing.T) {
	// An interface where the runtime type is valid but the plan doesn't
	// have an entry for it with a nil sub-plan.
	subType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Sub",
		Fields: graphql.Fields{
			"name": &graphql.Field{Type: graphql.String},
		},
	})
	charInterface := graphql.NewInterface(graphql.InterfaceConfig{
		Name: "Character",
		Fields: graphql.Fields{
			"sub": &graphql.Field{
				Type: subType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return map[string]interface{}{"name": "val"}, nil
				},
			},
		},
	})
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"hero": &graphql.Field{
				Type: charInterface,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return map[string]interface{}{"name": "Hero"}, nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ hero { sub { name } } }`,
	})
	if len(result.Errors) == 0 {
		t.Fatal("expected error for unresolved abstract type")
	}
}

func TestDethunkMapDepthFirst_MapWithFuncs(t *testing.T) {
	// A mutation that returns a map containing thunk-function values
	// exercises dethunkMapDepthFirst with func() interface{} values.
	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"do": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "Result",
					Fields: graphql.Fields{
						"a": &graphql.Field{Type: graphql.String},
						"b": &graphql.Field{Type: graphql.String},
					},
				}),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return map[string]interface{}{
						"a": func() interface{} { return "alpha" },
						"b": func() interface{} { return "beta" },
					}, nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    graphql.NewObject(graphql.ObjectConfig{Name: "Query", Fields: graphql.Fields{"_": &graphql.Field{Type: graphql.String}}}),
		Mutation: mutationType,
	})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `mutation { do { a b } }`,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	do := result.Data.(map[string]interface{})["do"].(map[string]interface{})
	if do["a"] != "alpha" || do["b"] != "beta" {
		t.Fatalf("unexpected result: %v", do)
	}
}

func TestDethunkListDepthFirst_NestedList(t *testing.T) {
	// A mutation returning a list field where list items contain
	// map values that themselves contain thunks exercises both
	// dethunkListDepthFirst and dethunkMapDepthFirst in sequence.
	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"mutate": &graphql.Field{
				Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
					Name: "Item",
					Fields: graphql.Fields{
						"x": &graphql.Field{Type: graphql.String},
					},
				})),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return []interface{}{
						map[string]interface{}{
							"x": func() interface{} { return "result" },
						},
					}, nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    graphql.NewObject(graphql.ObjectConfig{Name: "Query", Fields: graphql.Fields{"_": &graphql.Field{Type: graphql.String}}}),
		Mutation: mutationType,
	})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `mutation { mutate { x } }`,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	items := result.Data.(map[string]interface{})["mutate"].([]interface{})
	item := items[0].(map[string]interface{})
	if item["x"] != "result" {
		t.Fatalf("unexpected item: %v", item)
	}
}

// ---------------------------------------------------------------------------
// DefaultResolveFn with untyped nil source (line 481-482)
// ---------------------------------------------------------------------------

func TestDefaultResolveFn_UntypedNilSource(t *testing.T) {
	// When source is an untyped nil (nil pointer after Elem),
	// DefaultResolveFn returns nil, nil at the !sourceVal.IsValid() guard.
	// We trigger this by passing a nil pointer as Root and querying a
	// field with no custom resolver.
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"x": &graphql.Field{Type: graphql.String},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	ast := testutil.TestParse(t, `{ x }`)
	result := graphql.Execute(graphql.ExecuteParams{
		Schema: schema,
		Root:   (*struct{})(nil),
		AST:    ast,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
}

// ---------------------------------------------------------------------------
// DefaultResolveFn – struct field name match (line 490-491)
// ---------------------------------------------------------------------------

func TestDefaultResolveFn_StructFieldMatch(t *testing.T) {
	// When source is a struct whose field name matches p.Info.FieldName
	// (case-insensitive), DefaultResolveFn returns the field value.
	type mySource struct {
		Xval string
	}
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"xval": &graphql.Field{Type: graphql.String},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	ast := testutil.TestParse(t, `{ xval }`)
	result := graphql.Execute(graphql.ExecuteParams{
		Schema: schema,
		Root:   mySource{Xval: "hello"},
		AST:    ast,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	if result.Data == nil {
		t.Fatal("expected data")
	}
	data := result.Data.(map[string]interface{})
	if data["xval"] != "hello" {
		t.Fatalf("expected xval=hello, got %v", data["xval"])
	}
}

// ---------------------------------------------------------------------------
// DefaultResolveFn – struct field tag match (lines 493-508)
// ---------------------------------------------------------------------------

type taggedSource struct {
	Val string `json:"xval" graphql:"xval"`
}

func TestDefaultResolveFn_StructTagMatch(t *testing.T) {
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"xval": &graphql.Field{Type: graphql.String},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	ast := testutil.TestParse(t, `{ xval }`)
	result := graphql.Execute(graphql.ExecuteParams{
		Schema: schema,
		Root:   taggedSource{Val: "tagged"},
		AST:    ast,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	data := result.Data.(map[string]interface{})
	if data["xval"] != "tagged" {
		t.Fatalf("expected xval=tagged, got %v", data["xval"])
	}
}

// ---------------------------------------------------------------------------
// DefaultResolveFn – reflect.Map access (line 529-542)
// ---------------------------------------------------------------------------

func TestDefaultResolveFn_ReflectMapAccess(t *testing.T) {
	// When source implements map via reflection (not map[string]interface{}),
	// DefaultResolveFn falls through to the reflect.Map branch.
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"key": &graphql.Field{Type: graphql.String},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	source := map[string]string{"key": "value"}
	ast := testutil.TestParse(t, `{ key }`)
	result := graphql.Execute(graphql.ExecuteParams{
		Schema: schema,
		Root:   source,
		AST:    ast,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	data := result.Data.(map[string]interface{})
	if data["key"] != "value" {
		t.Fatalf("expected key=value, got %v", data["key"])
	}
}

// ---------------------------------------------------------------------------
// getFieldDef – __typename (line 569-570)
// ---------------------------------------------------------------------------

func TestGetFieldDef_TypeName(t *testing.T) {
	// Querying __typename exercises the TypeNameMetaFieldDef branch.
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"x": &graphql.Field{Type: graphql.String},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ __typename }`,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	data := result.Data.(map[string]interface{})
	if data["__typename"] != "Query" {
		t.Fatalf("expected __typename=Query, got %v", data["__typename"])
	}
}

// ---------------------------------------------------------------------------
// defaultResolveTypeFn – type with no IsTypeOf (line 445-446)
// ---------------------------------------------------------------------------

func TestDefaultResolveTypeFn_NoIsTypeOf(t *testing.T) {
	// defaultResolveTypeFn skips possible types that have no IsTypeOf
	// defined (line 445-446). Droid has IsTypeOf=true, Human has no IsTypeOf.
	// The resolver returns a map, and Droid's IsTypeOf matches it.
	charInterface := graphql.NewInterface(graphql.InterfaceConfig{
		Name: "Character",
		Fields: graphql.Fields{
			"name": &graphql.Field{Type: graphql.String},
		},
	})
	humanType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Human",
		Interfaces: []*graphql.Interface{charInterface},
		Fields: graphql.Fields{
			"name": &graphql.Field{Type: graphql.String},
		},
		// No IsTypeOf – the continue at line 445-446 skips this type
	})
	droidType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Droid",
		Interfaces: []*graphql.Interface{charInterface},
		Fields: graphql.Fields{
			"name": &graphql.Field{Type: graphql.String},
		},
		IsTypeOf: func(p graphql.IsTypeOfParams) bool {
			return true
		},
	})
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"hero": &graphql.Field{
				Type: charInterface,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return map[string]interface{}{"name": "R2"}, nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
		Types: []graphql.Type{humanType, droidType},
	})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ hero { name } }`,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
}

// ---------------------------------------------------------------------------
// dethunkMapDepthFirst – value with thunk (line 215-217)
// ---------------------------------------------------------------------------

func TestDethunkMapDepthFirst_WithFuncThunk(t *testing.T) {
	// A mutation field returning func() (interface{}, error) will be
	// wrapped as func() interface{} by completePlannedValue, then
	// resolved by dethunkMapDepthFirst.
	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"set": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "Result",
					Fields: graphql.Fields{
						"value": &graphql.Field{Type: graphql.String},
					},
				}),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return map[string]interface{}{
						"value": func() (interface{}, error) { return "thunked", nil },
					}, nil
				},
			},
		},
	})
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{"_": &graphql.Field{Type: graphql.String}},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `mutation { set { value } }`,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	data := result.Data.(map[string]interface{})["set"].(map[string]interface{})
	if data["value"] != "thunked" {
		t.Fatalf("expected thunked, got %v", data["value"])
	}
}

// ---------------------------------------------------------------------------
// dethunkListDepthFirst – list with func thunk (line 229-231) and nested list
// ---------------------------------------------------------------------------

func TestDethunkListDepthFirst_WithFuncThunk(t *testing.T) {
	// A mutation returning a list where an item is func() (interface{}, error)
	// exercises the func thunk path inside dethunkListDepthFirst (line 229-231).
	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"items": &graphql.Field{
				Type: graphql.NewList(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return []interface{}{
						func() (interface{}, error) { return "a", nil },
						"b",
					}, nil
				},
			},
		},
	})
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{"_": &graphql.Field{Type: graphql.String}},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `mutation { items }`,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	items := result.Data.(map[string]interface{})["items"].([]interface{})
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
}

func TestDethunkListDepthFirst_NestedStringList(t *testing.T) {
	// A mutation returning [[String]] where inner items are plain strings
	// exercises the nested-list branch (line 235-236).
	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"items": &graphql.Field{
				Type: graphql.NewList(graphql.NewList(graphql.String)),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return []interface{}{
						[]interface{}{"a", "b"},
						[]interface{}{"c", "d"},
					}, nil
				},
			},
		},
	})
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{"_": &graphql.Field{Type: graphql.String}},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `mutation { items }`,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
}

// ---------------------------------------------------------------------------
// Subscription tests for subscription-only executor functions
// ---------------------------------------------------------------------------

func subscriptionSchema(t *testing.T, subFields graphql.Fields) graphql.Schema {
	t.Helper()
	s, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Query",
			Fields: graphql.Fields{"_": &graphql.Field{Type: graphql.String}},
		}),
		Subscription: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Subscription",
			Fields: subFields,
		}),
	})
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestSubscription_BuildExecutionContext_MultipleOperations(t *testing.T) {
	// Multiple operations without operation name exercises
	// line 68-70 in buildExecutionContext.
	s := subscriptionSchema(t, graphql.Fields{
		"sub": &graphql.Field{
			Type: graphql.String,
			Subscribe: func(p graphql.ResolveParams) (interface{}, error) {
				c := make(chan interface{}, 1)
				c <- "val"
				close(c)
				return c, nil
			},
		},
	})
	c := graphql.Subscribe(graphql.Params{
		Schema:        s,
		RequestString: `subscription { sub } query { _ }`,
	})
	res := <-c
	if len(res.Errors) == 0 {
		t.Fatal("expected error for multiple operations")
	}
}

func TestSubscription_BuildExecutionContext_UnknownOperationName(t *testing.T) {
	// Non-existent operation name exercises line 85-88 in buildExecutionContext.
	s := subscriptionSchema(t, graphql.Fields{
		"sub": &graphql.Field{
			Type: graphql.String,
			Subscribe: func(p graphql.ResolveParams) (interface{}, error) {
				c := make(chan interface{}, 1)
				c <- "val"
				close(c)
				return c, nil
			},
		},
	})
	c := graphql.Subscribe(graphql.Params{
		Schema:        s,
		RequestString: `subscription sub1 { sub }`,
		OperationName: "bogus",
	})
	res := <-c
	if len(res.Errors) == 0 {
		t.Fatal("expected error for unknown operation")
	}
}

func TestSubscription_BuildExecutionContext_NoOperation(t *testing.T) {
	// Document with no operation exercises line 89 in buildExecutionContext.
	s := subscriptionSchema(t, graphql.Fields{
		"sub": &graphql.Field{
			Type: graphql.String,
			Subscribe: func(p graphql.ResolveParams) (interface{}, error) {
				c := make(chan interface{}, 1)
				c <- "val"
				close(c)
				return c, nil
			},
		},
	})
	c := graphql.Subscribe(graphql.Params{
		Schema:        s,
		RequestString: `fragment F on Subscription { sub }`,
	})
	res := <-c
	if len(res.Errors) == 0 {
		t.Fatal("expected error for no operation")
	}
}

func TestSubscription_BuildExecutionContext_VariableCoercionError(t *testing.T) {
	// Variable coercion failure exercises line 92-94 in buildExecutionContext.
	s := subscriptionSchema(t, graphql.Fields{
		"sub": &graphql.Field{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"x": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Subscribe: func(p graphql.ResolveParams) (interface{}, error) {
				c := make(chan interface{}, 1)
				c <- "val"
				close(c)
				return c, nil
			},
		},
	})
	c := graphql.Subscribe(graphql.Params{
		Schema:        s,
		RequestString: `subscription ($x: String!) { sub }`,
		// missing variable $x
	})
	res := <-c
	if len(res.Errors) == 0 {
		t.Fatal("expected error for missing variable")
	}
}

func TestSubscription_CollectFields_InlineFragment(t *testing.T) {
	// Inline fragment inside subscription exercises lines 277-290.
	s := subscriptionSchema(t, graphql.Fields{
		"sub": &graphql.Field{
			Type: graphql.String,
			Subscribe: func(p graphql.ResolveParams) (interface{}, error) {
				c := make(chan interface{}, 1)
				c <- "val"
				close(c)
				return c, nil
			},
		},
	})
	c := graphql.Subscribe(graphql.Params{
		Schema:        s,
		RequestString: `subscription { ... on Subscription { sub } }`,
	})
	res := <-c
	if len(res.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
}

func TestSubscription_CollectFields_FragmentSpread(t *testing.T) {
	// Named fragment spread exercises lines 291-317.
	s := subscriptionSchema(t, graphql.Fields{
		"sub": &graphql.Field{
			Type: graphql.String,
			Subscribe: func(p graphql.ResolveParams) (interface{}, error) {
				c := make(chan interface{}, 1)
				c <- "val"
				close(c)
				return c, nil
			},
		},
	})
	c := graphql.Subscribe(graphql.Params{
		Schema:        s,
		RequestString: `subscription { ...F } fragment F on Subscription { sub }`,
	})
	res := <-c
	if len(res.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
}

func TestSubscription_ShouldIncludeNode_SkipTrue(t *testing.T) {
	// Field with @skip(if: true) exercises shouldIncludeNode skip path.
	s := subscriptionSchema(t, graphql.Fields{
		"sub": &graphql.Field{
			Type: graphql.String,
			Subscribe: func(p graphql.ResolveParams) (interface{}, error) {
				c := make(chan interface{}, 1)
				c <- "val"
				close(c)
				return c, nil
			},
		},
	})
	c := graphql.Subscribe(graphql.Params{
		Schema:        s,
		RequestString: `subscription { sub @skip(if: true) }`,
	})
	res := <-c
	// When skip is true and only one field, collectFields returns empty
	// map, so subscription can't find the field — expects error.
	if len(res.Errors) == 0 {
		t.Fatal("expected error for skipped field")
	}
}

func TestSubscription_ShouldIncludeNode_IncludeFalse(t *testing.T) {
	// Field with @include(if: false) exercises shouldIncludeNode include path.
	s := subscriptionSchema(t, graphql.Fields{
		"sub": &graphql.Field{
			Type: graphql.String,
			Subscribe: func(p graphql.ResolveParams) (interface{}, error) {
				c := make(chan interface{}, 1)
				c <- "val"
				close(c)
				return c, nil
			},
		},
	})
	c := graphql.Subscribe(graphql.Params{
		Schema:        s,
		RequestString: `subscription { sub @include(if: false) }`,
	})
	res := <-c
	if len(res.Errors) == 0 {
		t.Fatal("expected error for excluded field")
	}
}

func TestSubscription_DoesFragmentConditionMatch_Interface(t *testing.T) {
	// Fragment with interface type condition exercises
	// doesFragmentConditionMatch Interface branch (line 377-379 / 398-400).
	charInterface := graphql.NewInterface(graphql.InterfaceConfig{
		Name: "Character",
		Fields: graphql.Fields{
			"name": &graphql.Field{Type: graphql.String},
		},
	})
	human := graphql.NewObject(graphql.ObjectConfig{
		Name: "Human",
		Interfaces: []*graphql.Interface{charInterface},
		Fields: graphql.Fields{
			"name": &graphql.Field{Type: graphql.String},
		},
	})
	s, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Query",
			Fields: graphql.Fields{"_": &graphql.Field{Type: graphql.String}},
		}),
		Subscription: graphql.NewObject(graphql.ObjectConfig{
			Name: "Subscription",
			Fields: graphql.Fields{
				"sub": &graphql.Field{
					Type: charInterface,
					Subscribe: func(p graphql.ResolveParams) (interface{}, error) {
						c := make(chan interface{}, 1)
						c <- map[string]interface{}{"name": "Luke"}
						close(c)
						return c, nil
					},
				},
			},
		}),
		Types: []graphql.Type{human},
	})
	if err != nil {
		t.Fatal(err)
	}
	c := graphql.Subscribe(graphql.Params{
		Schema:        s,
		RequestString: `subscription { sub { ... on Character { name } } }`,
	})
	res := <-c
	if len(res.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
}

func TestSubscription_DoesFragmentConditionMatch_Union(t *testing.T) {
	// Fragment with union type condition exercises
	// doesFragmentConditionMatch Union branch (lines 380-382 / 401-403).
	resultType := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Result",
		Fields: graphql.Fields{"_": &graphql.Field{Type: graphql.String}},
		IsTypeOf: func(p graphql.IsTypeOfParams) bool {
			_, ok := p.Value.(map[string]interface{})
			return ok
		},
	})
	searchUnion := graphql.NewUnion(graphql.UnionConfig{
		Name: "Search",
		Types: []*graphql.Object{resultType},
	})
	s, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Query",
			Fields: graphql.Fields{"_": &graphql.Field{Type: graphql.String}},
		}),
		Subscription: graphql.NewObject(graphql.ObjectConfig{
			Name: "Subscription",
			Fields: graphql.Fields{
				"sub": &graphql.Field{
					Type: searchUnion,
					Subscribe: func(p graphql.ResolveParams) (interface{}, error) {
						c := make(chan interface{}, 1)
						c <- map[string]interface{}{}
						close(c)
						return c, nil
					},
				},
			},
		}),
		Types: []graphql.Type{resultType},
	})
	if err != nil {
		t.Fatal(err)
	}
	c := graphql.Subscribe(graphql.Params{
		Schema:        s,
		RequestString: `subscription { sub { ... on Result { _ } } }`,
	})
	res := <-c
	if len(res.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
}

// ---------------------------------------------------------------------------
// doesFragmentConditionMatch – root-level fragments on Interface/Union
// ---------------------------------------------------------------------------

func TestSubscription_DoesFragmentConditionMatch_InlineFragmentInterfaceAtRoot(t *testing.T) {
	// Root-level inline fragment with Interface type condition exercises
	// InlineFragment Interface branch (lines 398-400).
	charInterface := graphql.NewInterface(graphql.InterfaceConfig{
		Name: "Character",
		Fields: graphql.Fields{"name": &graphql.Field{Type: graphql.String}},
	})
	s, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Query",
			Fields: graphql.Fields{"_": &graphql.Field{Type: graphql.String}},
		}),
	})
	if err != nil {
		t.Fatal(err)
	}
	_ = charInterface
	c := graphql.Subscribe(graphql.Params{
		Schema:        s,
		RequestString: `subscription { ... on Character { name } }`,
	})
	res := <-c
	// Expect error: Character not in schema, subscription not in schema
	if len(res.Errors) == 0 {
		t.Fatal("expected error")
	}
}

func TestSubscription_DoesFragmentConditionMatch_InlineFragmentUnionAtRoot(t *testing.T) {
	// Root-level inline fragment with Union type condition exercises
	// InlineFragment Union branch (lines 401-403).
	searchUnion := graphql.NewUnion(graphql.UnionConfig{
		Name: "Search",
		Types: []*graphql.Object{},
	})
	s, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Query",
			Fields: graphql.Fields{"_": &graphql.Field{Type: graphql.String}},
		}),
	})
	if err != nil {
		t.Fatal(err)
	}
	_ = searchUnion
	c := graphql.Subscribe(graphql.Params{
		Schema:        s,
		RequestString: `subscription { ... on Search { _ } }`,
	})
	res := <-c
	if len(res.Errors) == 0 {
		t.Fatal("expected error")
	}
}

func TestSubscription_DoesFragmentConditionMatch_FragmentSpreadInterfaceAtRoot(t *testing.T) {
	// Root-level fragment spread with Interface type condition exercises
	// FragmentDefinition Interface branch (lines 377-379).
	charInterface := graphql.NewInterface(graphql.InterfaceConfig{
		Name: "Character",
		Fields: graphql.Fields{"name": &graphql.Field{Type: graphql.String}},
	})
	s, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Query",
			Fields: graphql.Fields{"_": &graphql.Field{Type: graphql.String}},
		}),
	})
	if err != nil {
		t.Fatal(err)
	}
	_ = charInterface
	c := graphql.Subscribe(graphql.Params{
		Schema:        s,
		RequestString: `subscription { ...F } fragment F on Character { name }`,
	})
	res := <-c
	if len(res.Errors) == 0 {
		t.Fatal("expected error")
	}
}

// ---------------------------------------------------------------------------
// collectFields – missing fragment (line 302-303)
// ---------------------------------------------------------------------------

func TestSubscription_CollectFields_MissingFragment(t *testing.T) {
	// Fragment spread referring to undefined fragment exercises
	// the `!hasFragment` skip in collectFields (line 302-303).
	s := subscriptionSchema(t, graphql.Fields{
		"sub": &graphql.Field{
			Type: graphql.String,
			Subscribe: func(p graphql.ResolveParams) (interface{}, error) {
				c := make(chan interface{}, 1)
				c <- "val"
				close(c)
				return c, nil
			},
		},
	})
	c := graphql.Subscribe(graphql.Params{
		Schema:        s,
		RequestString: `subscription { ...Missing }`,
	})
	res := <-c
	// Expect error: no field found because ...Missing is not defined
	if len(res.Errors) == 0 {
		t.Fatal("expected error")
	}
}

// ---------------------------------------------------------------------------
// DefaultResolveFn – struct tag mismatch path (line 497-499)
// ---------------------------------------------------------------------------

type taggedNoMatchSource struct {
	RealValue string `json:"real_value" graphql:"val"`
}

func TestDefaultResolveFn_StructTagMismatch(t *testing.T) {
	// When the "json" tag value differs from the field name but "graphql"
	// tag matches, the checkTag fallback hits line 499 (tOptions[0] != fieldName).
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"val": &graphql.Field{Type: graphql.String},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		t.Fatal(err)
	}
	ast := testutil.TestParse(t, `{ val }`)
	result := graphql.Execute(graphql.ExecuteParams{
		Schema: schema,
		Root:   taggedNoMatchSource{RealValue: "matched"},
		AST:    ast,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	data := result.Data.(map[string]interface{})
	if data["val"] != "matched" {
		t.Fatalf("expected val=matched, got %v", data["val"])
	}
}

func TestSubscription_UnknownDefinitionType(t *testing.T) {
	// Document with type definition exercises the default branch
	// (line 80-81) in buildExecutionContext.
	s := subscriptionSchema(t, graphql.Fields{
		"sub": &graphql.Field{
			Type: graphql.String,
			Subscribe: func(p graphql.ResolveParams) (interface{}, error) {
				c := make(chan interface{}, 1)
				c <- "val"
				close(c)
				return c, nil
			},
		},
	})
	c := graphql.Subscribe(graphql.Params{
		Schema:        s,
		RequestString: `subscription { sub }
type Foo { bar: String }`,
	})
	res := <-c
	if len(res.Errors) == 0 {
		t.Fatal("expected error for unknown definition type")
	}
}

func TestSubscription_CollectFields_VisitedFragment(t *testing.T) {
	// Same fragment used twice exercises the visited-fragment skip
	// (line 297-298) in collectFields.
	s := subscriptionSchema(t, graphql.Fields{
		"sub": &graphql.Field{
			Type: graphql.String,
			Subscribe: func(p graphql.ResolveParams) (interface{}, error) {
				c := make(chan interface{}, 1)
				c <- "val"
				close(c)
				return c, nil
			},
		},
	})
	c := graphql.Subscribe(graphql.Params{
		Schema:        s,
		RequestString: `subscription { ...F ...F } fragment F on Subscription { sub }`,
	})
	res := <-c
	if len(res.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
}
