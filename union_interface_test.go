package graphql_test

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/testutil"
)

type testNamedType interface {
}
type testPet interface {
}
type testDog2 struct {
	Name  string `json:"name"`
	Barks bool   `json:"barks"`
}

type testCat2 struct {
	Name  string `json:"name"`
	Meows bool   `json:"meows"`
}

type testPerson struct {
	Name    string          `json:"name"`
	Pets    []testPet       `json:"pets"`
	Friends []testNamedType `json:"friends"`
}

var namedType = graphql.NewInterface(graphql.InterfaceConfig{
	Name: "Named",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
	},
})
var dogType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Dog",
	Interfaces: []*graphql.Interface{
		namedType,
	},
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"barks": &graphql.Field{
			Type: graphql.Boolean,
		},
	},
	IsTypeOf: func(p graphql.IsTypeOfParams) bool {
		_, ok := p.Value.(*testDog2)
		return ok
	},
})
var catType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Cat",
	Interfaces: []*graphql.Interface{
		namedType,
	},
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"meows": &graphql.Field{
			Type: graphql.Boolean,
		},
	},
	IsTypeOf: func(p graphql.IsTypeOfParams) bool {
		_, ok := p.Value.(*testCat2)
		return ok
	},
})
var petType = graphql.NewUnion(graphql.UnionConfig{
	Name: "Pet",
	Types: []*graphql.Object{
		dogType, catType,
	},
	ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
		if _, ok := p.Value.(*testCat2); ok {
			return catType
		}
		if _, ok := p.Value.(*testDog2); ok {
			return dogType
		}
		return nil
	},
})
var personType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Person",
	Interfaces: []*graphql.Interface{
		namedType,
	},
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"pets": &graphql.Field{
			Type: graphql.NewList(petType),
		},
		"friends": &graphql.Field{
			Type: graphql.NewList(namedType),
		},
	},
	IsTypeOf: func(p graphql.IsTypeOfParams) bool {
		_, ok := p.Value.(*testPerson)
		return ok
	},
})

var unionInterfaceTestSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: personType,
	Types: []graphql.Type{petType},
})

var garfield = &testCat2{"Garfield", false}
var odie = &testDog2{"Odie", true}
var liz = &testPerson{
	Name: "Liz",
}
var john = &testPerson{
	Name: "John",
	Pets: []testPet{
		garfield, odie,
	},
	Friends: []testNamedType{
		liz, odie,
	},
}

func TestUnionIntersectionTypes_CanIntrospectOnUnionAndIntersectionTypes(t *testing.T) {
	doc := `
      {
        Named: __type(name: "Named") {
          kind
          name
          fields { name }
          interfaces { name }
          possibleTypes { name }
          enumValues { name }
          inputFields { name }
        }
        Pet: __type(name: "Pet") {
          kind
          name
          fields { name }
          interfaces { name }
          possibleTypes { name }
          enumValues { name }
          inputFields { name }
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"Named": map[string]interface{}{
				"kind": "INTERFACE",
				"name": "Named",
				"fields": []interface{}{
					map[string]interface{}{
						"name": "name",
					},
				},
				"interfaces": nil,
				"possibleTypes": []interface{}{
					map[string]interface{}{
						"name": "Dog",
					},
					map[string]interface{}{
						"name": "Cat",
					},
					map[string]interface{}{
						"name": "Person",
					},
				},
				"enumValues":  nil,
				"inputFields": nil,
			},
			"Pet": map[string]interface{}{
				"kind":       "UNION",
				"name":       "Pet",
				"fields":     nil,
				"interfaces": nil,
				"possibleTypes": []interface{}{
					map[string]interface{}{
						"name": "Dog",
					},
					map[string]interface{}{
						"name": "Cat",
					},
				},
				"enumValues":  nil,
				"inputFields": nil,
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: unionInterfaceTestSchema,
		AST:    ast,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !testutil.ContainSubset(expected.Data.(map[string]interface{}), result.Data.(map[string]interface{})) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Data, result.Data))
	}
}
func TestUnionIntersectionTypes_ExecutesUsingUnionTypes(t *testing.T) {
	// NOTE: This is an *invalid* query, but it should be an *executable* query.
	doc := `
      {
        __typename
        name
        pets {
          __typename
          name
          barks
          meows
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"__typename": "Person",
			"name":       "John",
			"pets": []interface{}{
				map[string]interface{}{
					"__typename": "Cat",
					"name":       "Garfield",
					"meows":      false,
				},
				map[string]interface{}{
					"__typename": "Dog",
					"name":       "Odie",
					"barks":      true,
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: unionInterfaceTestSchema,
		AST:    ast,
		Root:   john,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestUnionIntersectionTypes_ExecutesUnionTypesWithInlineFragments(t *testing.T) {
	// This is the valid version of the query in the above test.
	doc := `
      {
        __typename
        name
        pets {
          __typename
          ... on Dog {
            name
            barks
          }
          ... on Cat {
            name
            meows
          }
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"__typename": "Person",
			"name":       "John",
			"pets": []interface{}{
				map[string]interface{}{
					"__typename": "Cat",
					"name":       "Garfield",
					"meows":      false,
				},
				map[string]interface{}{
					"__typename": "Dog",
					"name":       "Odie",
					"barks":      true,
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: unionInterfaceTestSchema,
		AST:    ast,
		Root:   john,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestUnionIntersectionTypes_ExecutesUsingInterfaceTypes(t *testing.T) {

	// NOTE: This is an *invalid* query, but it should be an *executable* query.
	doc := `
      {
        __typename
        name
        friends {
          __typename
          name
          barks
          meows
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"__typename": "Person",
			"name":       "John",
			"friends": []interface{}{
				map[string]interface{}{
					"__typename": "Person",
					"name":       "Liz",
				},
				map[string]interface{}{
					"__typename": "Dog",
					"name":       "Odie",
					"barks":      true,
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: unionInterfaceTestSchema,
		AST:    ast,
		Root:   john,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestUnionIntersectionTypes_ExecutesInterfaceTypesWithInlineFragments(t *testing.T) {

	// This is the valid version of the query in the above test.
	doc := `
      {
        __typename
        name
        friends {
          __typename
          name
          ... on Dog {
            barks
          }
          ... on Cat {
            meows
          }
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"__typename": "Person",
			"name":       "John",
			"friends": []interface{}{
				map[string]interface{}{
					"__typename": "Person",
					"name":       "Liz",
				},
				map[string]interface{}{
					"__typename": "Dog",
					"name":       "Odie",
					"barks":      true,
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: unionInterfaceTestSchema,
		AST:    ast,
		Root:   john,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestUnionIntersectionTypes_AllowsFragmentConditionsToBeAbstractTypes(t *testing.T) {

	doc := `
      {
        __typename
        name
        pets { ...PetFields }
        friends { ...FriendFields }
      }

      fragment PetFields on Pet {
        __typename
        ... on Dog {
          name
          barks
        }
        ... on Cat {
          name
          meows
        }
      }

      fragment FriendFields on Named {
        __typename
        name
        ... on Dog {
          barks
        }
        ... on Cat {
          meows
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"__typename": "Person",
			"name":       "John",
			"friends": []interface{}{
				map[string]interface{}{
					"__typename": "Person",
					"name":       "Liz",
				},
				map[string]interface{}{
					"__typename": "Dog",
					"name":       "Odie",
					"barks":      true,
				},
			},
			"pets": []interface{}{
				map[string]interface{}{
					"__typename": "Cat",
					"name":       "Garfield",
					"meows":      false,
				},
				map[string]interface{}{
					"__typename": "Dog",
					"name":       "Odie",
					"barks":      true,
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: unionInterfaceTestSchema,
		AST:    ast,
		Root:   john,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestUnionIntersectionTypes_GetsExecutionInfoInResolver(t *testing.T) {

	var encounteredContextValue string
	var encounteredSchema graphql.Schema
	var encounteredRootValue string
	var personType2 *graphql.Object

	namedType2 := graphql.NewInterface(graphql.InterfaceConfig{
		Name: "Named",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
		},
		ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
			encounteredSchema = p.Info.Schema
			encounteredContextValue, _ = p.Context.Value("authToken").(string)
			encounteredRootValue = p.Info.RootValue.(*testPerson).Name
			return personType2
		},
	})

	personType2 = graphql.NewObject(graphql.ObjectConfig{
		Name: "Person",
		Interfaces: []*graphql.Interface{
			namedType2,
		},
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"friends": &graphql.Field{
				Type: graphql.NewList(namedType2),
			},
		},
	})

	schema2, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: personType2,
	})

	john2 := &testPerson{
		Name: "John",
		Friends: []testNamedType{
			liz,
		},
	}

	doc := `{ name, friends { name } }`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"name": "John",
			"friends": []interface{}{
				map[string]interface{}{
					"name": "Liz",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// create context
	ctx := context.Background()
	ctx = context.WithValue(ctx, "authToken", "contextStringValue123")

	// execute
	ep := graphql.ExecuteParams{
		Schema:  schema2,
		AST:     ast,
		Root:    john2,
		Context: ctx,
	}
	result := testutil.TestExecute(t, ep)

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
	if !reflect.DeepEqual("contextStringValue123", encounteredContextValue) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff("contextStringValue123", encounteredContextValue))
	}
	if !reflect.DeepEqual("John", encounteredRootValue) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff("John", encounteredRootValue))
	}
	if !reflect.DeepEqual(schema2, encounteredSchema) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(schema2, encounteredSchema))
	}
}

func TestUnionIntersectionTypes_ValueMayBeNilPointer(t *testing.T) {
	var unionInterfaceTestSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"query": &graphql.Field{
					Type: graphql.NewObject(graphql.ObjectConfig{
						Name: "query",
						Fields: graphql.Fields{
							"pet": &graphql.Field{
								Type: petType,
							},
							"named": &graphql.Field{
								Type: namedType,
							},
						},
					}),
					Resolve: func(_ graphql.ResolveParams) (interface{}, error) {
						return struct {
							Pet   *testCat2 `graphql:"pet"`
							Named *testCat2 `graphql:"named"`
						}{nil, nil}, nil
					},
				},
			},
		}),
	})
	query := `{
		query {
			pet {
				__typename
			}
			named {
				__typename
				name
			}
		}
	}`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"query": map[string]interface{}{
				"pet":   nil,
				"named": nil,
			}},
	}
	result := g(t, graphql.Params{
		Schema:        unionInterfaceTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestUnionIntersectionTypes_ConcurrentRequests(t *testing.T) {
	named := graphql.NewInterface(graphql.InterfaceConfig{
		Name:   "Named",
		Fields: graphql.Fields{"name": &graphql.Field{Type: graphql.String}},
	})
	hasFriends := graphql.NewInterface(graphql.InterfaceConfig{
		Name:   "HasFriends",
		Fields: graphql.Fields{"friends": &graphql.Field{Type: graphql.NewList(named)}},
	})
	dog := graphql.NewObject(graphql.ObjectConfig{
		Name:       "Dog",
		Interfaces: []*graphql.Interface{named, hasFriends},
		IsTypeOf:   func(p graphql.IsTypeOfParams) bool { return true },
		Fields:     graphql.Fields{"name": &graphql.Field{Type: graphql.String}},
	})
	// Dog's friends field is narrower than the HasFriends interface, in that
	// it's of type [Dog], not type [Named]. This schema pattern makes graphql-go
	// call `Schema.IsPossibleType` to validate the interface narrowing, which
	// *used to* populate a private mutable map of interface implementations,
	// which would then cause crashes during concurrent calls to Execute
	dog.AddFieldConfig("friends", &graphql.Field{Type: graphql.NewList(dog)})
	pet := graphql.NewUnion(graphql.UnionConfig{
		Name:  "Pet",
		Types: []*graphql.Object{dog},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				// Resolving the Pet union at request time checks that Dog is
				// one of Pet's possible types, which must be safe to do from
				// many concurrent requests.
				"pet": &graphql.Field{
					Type: pet,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return map[string]interface{}{}, nil
					},
				},
			},
		}),
	})
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := graphql.Do(graphql.Params{
				Schema:        schema,
				RequestString: `{ pet { __typename } }`,
			})
			if len(result.Errors) > 0 {
				t.Errorf("unexpected errors: %v", result.Errors)
			}
		}()
	}
	wg.Wait()
}
