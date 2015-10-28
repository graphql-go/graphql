package graphql

import (
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

var namedType = NewInterface(InterfaceConfig{
	Name: "Named",
	Fields: FieldConfigMap{
		"name": &FieldConfig{
			Type: String,
		},
	},
})
var dogType = NewObject(ObjectConfig{
	Name: "Dog",
	Interfaces: []*Interface{
		namedType,
	},
	Fields: FieldConfigMap{
		"name": &FieldConfig{
			Type: String,
		},
		"barks": &FieldConfig{
			Type: Boolean,
		},
	},
	IsTypeOf: func(value interface{}, info ResolveInfo) bool {
		_, ok := value.(*testDog2)
		return ok
	},
})
var catType = NewObject(ObjectConfig{
	Name: "Cat",
	Interfaces: []*Interface{
		namedType,
	},
	Fields: FieldConfigMap{
		"name": &FieldConfig{
			Type: String,
		},
		"meows": &FieldConfig{
			Type: Boolean,
		},
	},
	IsTypeOf: func(value interface{}, info ResolveInfo) bool {
		_, ok := value.(*testCat2)
		return ok
	},
})
var petType = NewUnion(UnionConfig{
	Name: "Pet",
	Types: []*Object{
		dogType, catType,
	},
	ResolveType: func(value interface{}, info ResolveInfo) *Object {
		if _, ok := value.(*testCat2); ok {
			return catType
		}
		if _, ok := value.(*testDog2); ok {
			return dogType
		}
		return nil
	},
})
var personType = NewObject(ObjectConfig{
	Name: "Person",
	Interfaces: []*Interface{
		namedType,
	},
	Fields: FieldConfigMap{
		"name": &FieldConfig{
			Type: String,
		},
		"pets": &FieldConfig{
			Type: NewList(petType),
		},
		"friends": &FieldConfig{
			Type: NewList(namedType),
		},
	},
	IsTypeOf: func(value interface{}, info ResolveInfo) bool {
		_, ok := value.(*testPerson)
		return ok
	},
})

var unionInterfaceTestSchema, _ = NewSchema(SchemaConfig{
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
	expected := &Result{
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
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: unionInterfaceTestSchema,
		AST:    ast,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
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
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: unionInterfaceTestSchema,
		AST:    ast,
		Root:   john,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
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
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: unionInterfaceTestSchema,
		AST:    ast,
		Root:   john,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
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
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: unionInterfaceTestSchema,
		AST:    ast,
		Root:   john,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
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
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: unionInterfaceTestSchema,
		AST:    ast,
		Root:   john,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
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
	expected := &Result{
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
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: unionInterfaceTestSchema,
		AST:    ast,
		Root:   john,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestUnionIntersectionTypes_GetsExecutionInfoInResolver(t *testing.T) {

	var encounteredSchema *Schema
	var encounteredRootValue interface{}

	var personType2 *Object

	namedType2 := NewInterface(InterfaceConfig{
		Name: "Named",
		Fields: FieldConfigMap{
			"name": &FieldConfig{
				Type: String,
			},
		},
		ResolveType: func(value interface{}, info ResolveInfo) *Object {
			encounteredSchema = &info.Schema
			encounteredRootValue = info.RootValue
			return personType2
		},
	})

	personType2 = NewObject(ObjectConfig{
		Name: "Person",
		Interfaces: []*Interface{
			namedType2,
		},
		Fields: FieldConfigMap{
			"name": &FieldConfig{
				Type: String,
			},
			"friends": &FieldConfig{
				Type: NewList(namedType2),
			},
		},
	})

	schema2, _ := NewSchema(SchemaConfig{
		Query: personType2,
	})

	john2 := &testPerson{
		Name: "John",
		Friends: []testNamedType{
			liz,
		},
	}

	doc := `{ name, friends { name } }`
	expected := &Result{
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
	ast := TestParse(t, doc)

	// execute
	ep := ExecuteParams{
		Schema: schema2,
		AST:    ast,
		Root:   john2,
	}
	result := TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
