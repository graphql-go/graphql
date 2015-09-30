package executor_test

import (
	"reflect"
	"testing"

	"github.com/chris-ramon/graphql-go/executor"
	"github.com/chris-ramon/graphql-go/gqltypes"
	"github.com/chris-ramon/graphql-go/testutil"
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

var namedType = gqltypes.NewGraphQLInterfaceType(gqltypes.GraphQLInterfaceTypeConfig{
	Name: "Named",
	Fields: gqltypes.GraphQLFieldConfigMap{
		"name": &gqltypes.GraphQLFieldConfig{
			Type: gqltypes.GraphQLString,
		},
	},
})
var dogType = gqltypes.NewGraphQLObjectType(gqltypes.GraphQLObjectTypeConfig{
	Name: "Dog",
	Interfaces: []*gqltypes.GraphQLInterfaceType{
		namedType,
	},
	Fields: gqltypes.GraphQLFieldConfigMap{
		"name": &gqltypes.GraphQLFieldConfig{
			Type: gqltypes.GraphQLString,
		},
		"barks": &gqltypes.GraphQLFieldConfig{
			Type: gqltypes.GraphQLBoolean,
		},
	},
	IsTypeOf: func(value interface{}, info gqltypes.GraphQLResolveInfo) bool {
		_, ok := value.(*testDog2)
		return ok
	},
})
var catType = gqltypes.NewGraphQLObjectType(gqltypes.GraphQLObjectTypeConfig{
	Name: "Cat",
	Interfaces: []*gqltypes.GraphQLInterfaceType{
		namedType,
	},
	Fields: gqltypes.GraphQLFieldConfigMap{
		"name": &gqltypes.GraphQLFieldConfig{
			Type: gqltypes.GraphQLString,
		},
		"meows": &gqltypes.GraphQLFieldConfig{
			Type: gqltypes.GraphQLBoolean,
		},
	},
	IsTypeOf: func(value interface{}, info gqltypes.GraphQLResolveInfo) bool {
		_, ok := value.(*testCat2)
		return ok
	},
})
var petType = gqltypes.NewGraphQLUnionType(gqltypes.GraphQLUnionTypeConfig{
	Name: "Pet",
	Types: []*gqltypes.GraphQLObjectType{
		dogType, catType,
	},
	ResolveType: func(value interface{}, info gqltypes.GraphQLResolveInfo) *gqltypes.GraphQLObjectType {
		if _, ok := value.(*testCat2); ok {
			return catType
		}
		if _, ok := value.(*testDog2); ok {
			return dogType
		}
		return nil
	},
})
var personType = gqltypes.NewGraphQLObjectType(gqltypes.GraphQLObjectTypeConfig{
	Name: "Person",
	Interfaces: []*gqltypes.GraphQLInterfaceType{
		namedType,
	},
	Fields: gqltypes.GraphQLFieldConfigMap{
		"name": &gqltypes.GraphQLFieldConfig{
			Type: gqltypes.GraphQLString,
		},
		"pets": &gqltypes.GraphQLFieldConfig{
			Type: gqltypes.NewGraphQLList(petType),
		},
		"friends": &gqltypes.GraphQLFieldConfig{
			Type: gqltypes.NewGraphQLList(namedType),
		},
	},
	IsTypeOf: func(value interface{}, info gqltypes.GraphQLResolveInfo) bool {
		_, ok := value.(*testPerson)
		return ok
	},
})

var unionInterfaceTestSchema, _ = gqltypes.NewGraphQLSchema(gqltypes.GraphQLSchemaConfig{
	Query: personType,
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
	expected := &gqltypes.GraphQLResult{
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
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: unionInterfaceTestSchema,
		AST:    ast,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
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
	expected := &gqltypes.GraphQLResult{
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
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: unionInterfaceTestSchema,
		AST:    ast,
		Root:   john,
	}
	result := testutil.Execute(t, ep)
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
	expected := &gqltypes.GraphQLResult{
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
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: unionInterfaceTestSchema,
		AST:    ast,
		Root:   john,
	}
	result := testutil.Execute(t, ep)
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
	expected := &gqltypes.GraphQLResult{
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
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: unionInterfaceTestSchema,
		AST:    ast,
		Root:   john,
	}
	result := testutil.Execute(t, ep)
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
	expected := &gqltypes.GraphQLResult{
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
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: unionInterfaceTestSchema,
		AST:    ast,
		Root:   john,
	}
	result := testutil.Execute(t, ep)
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
	expected := &gqltypes.GraphQLResult{
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
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: unionInterfaceTestSchema,
		AST:    ast,
		Root:   john,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestUnionIntersectionTypes_GetsExecutionInfoInResolver(t *testing.T) {

	var encounteredSchema *gqltypes.GraphQLSchema
	var encounteredRootValue interface{}

	var personType2 *gqltypes.GraphQLObjectType

	namedType2 := gqltypes.NewGraphQLInterfaceType(gqltypes.GraphQLInterfaceTypeConfig{
		Name: "Named",
		Fields: gqltypes.GraphQLFieldConfigMap{
			"name": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLString,
			},
		},
		ResolveType: func(value interface{}, info gqltypes.GraphQLResolveInfo) *gqltypes.GraphQLObjectType {
			encounteredSchema = &info.Schema
			encounteredRootValue = info.RootValue
			return personType2
		},
	})

	personType2 = gqltypes.NewGraphQLObjectType(gqltypes.GraphQLObjectTypeConfig{
		Name: "Person",
		Interfaces: []*gqltypes.GraphQLInterfaceType{
			namedType2,
		},
		Fields: gqltypes.GraphQLFieldConfigMap{
			"name": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLString,
			},
			"friends": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.NewGraphQLList(namedType2),
			},
		},
	})

	schema2, _ := gqltypes.NewGraphQLSchema(gqltypes.GraphQLSchemaConfig{
		Query: personType2,
	})

	john2 := &testPerson{
		Name: "John",
		Friends: []testNamedType{
			liz,
		},
	}

	doc := `{ name, friends { name } }`
	expected := &gqltypes.GraphQLResult{
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
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: schema2,
		AST:    ast,
		Root:   john2,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
