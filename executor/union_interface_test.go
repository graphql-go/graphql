package executor_test

import (
	"github.com/chris-ramon/graphql/executor"
	"github.com/chris-ramon/graphql/testutil"
	"github.com/chris-ramon/graphql/types"
	"reflect"
	"testing"
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

var namedType = types.NewInterface(types.InterfaceConfig{
	Name: "Named",
	Fields: types.FieldConfigMap{
		"name": &types.FieldConfig{
			Type: types.String,
		},
	},
})
var dogType = types.NewObject(types.ObjectConfig{
	Name: "Dog",
	Interfaces: []*types.Interface{
		namedType,
	},
	Fields: types.FieldConfigMap{
		"name": &types.FieldConfig{
			Type: types.String,
		},
		"barks": &types.FieldConfig{
			Type: types.Boolean,
		},
	},
	IsTypeOf: func(value interface{}, info types.ResolveInfo) bool {
		_, ok := value.(*testDog2)
		return ok
	},
})
var catType = types.NewObject(types.ObjectConfig{
	Name: "Cat",
	Interfaces: []*types.Interface{
		namedType,
	},
	Fields: types.FieldConfigMap{
		"name": &types.FieldConfig{
			Type: types.String,
		},
		"meows": &types.FieldConfig{
			Type: types.Boolean,
		},
	},
	IsTypeOf: func(value interface{}, info types.ResolveInfo) bool {
		_, ok := value.(*testCat2)
		return ok
	},
})
var petType = types.NewUnion(types.UnionConfig{
	Name: "Pet",
	Types: []*types.Object{
		dogType, catType,
	},
	ResolveType: func(value interface{}, info types.ResolveInfo) *types.Object {
		if _, ok := value.(*testCat2); ok {
			return catType
		}
		if _, ok := value.(*testDog2); ok {
			return dogType
		}
		return nil
	},
})
var personType = types.NewObject(types.ObjectConfig{
	Name: "Person",
	Interfaces: []*types.Interface{
		namedType,
	},
	Fields: types.FieldConfigMap{
		"name": &types.FieldConfig{
			Type: types.String,
		},
		"pets": &types.FieldConfig{
			Type: types.NewList(petType),
		},
		"friends": &types.FieldConfig{
			Type: types.NewList(namedType),
		},
	},
	IsTypeOf: func(value interface{}, info types.ResolveInfo) bool {
		_, ok := value.(*testPerson)
		return ok
	},
})

var unionInterfaceTestSchema, _ = types.NewSchema(types.SchemaConfig{
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
	expected := &types.Result{
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
	expected := &types.Result{
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
	expected := &types.Result{
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
	expected := &types.Result{
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
	expected := &types.Result{
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
	expected := &types.Result{
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

	var encounteredSchema *types.Schema
	var encounteredRootValue interface{}

	var personType2 *types.Object

	namedType2 := types.NewInterface(types.InterfaceConfig{
		Name: "Named",
		Fields: types.FieldConfigMap{
			"name": &types.FieldConfig{
				Type: types.String,
			},
		},
		ResolveType: func(value interface{}, info types.ResolveInfo) *types.Object {
			encounteredSchema = &info.Schema
			encounteredRootValue = info.RootValue
			return personType2
		},
	})

	personType2 = types.NewObject(types.ObjectConfig{
		Name: "Person",
		Interfaces: []*types.Interface{
			namedType2,
		},
		Fields: types.FieldConfigMap{
			"name": &types.FieldConfig{
				Type: types.String,
			},
			"friends": &types.FieldConfig{
				Type: types.NewList(namedType2),
			},
		},
	})

	schema2, _ := types.NewSchema(types.SchemaConfig{
		Query: personType2,
	})

	john2 := &testPerson{
		Name: "John",
		Friends: []testNamedType{
			liz,
		},
	}

	doc := `{ name, friends { name } }`
	expected := &types.Result{
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
