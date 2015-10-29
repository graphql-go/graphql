package graphql

import (
	"reflect"
	"testing"

	"github.com/chris-ramon/graphql-go/gqlerrors"
	"github.com/chris-ramon/graphql-go/language/location"
)

type testDog struct {
	Name  string
	Woofs bool
}

type testCat struct {
	Name  string
	Meows bool
}

type testHuman struct {
	Name string
}

func TestIsTypeOfUsedToResolveRuntimeTypeForInterface(t *testing.T) {

	petType := NewInterface(InterfaceConfig{
		Name: "Pet",
		Fields: FieldConfigMap{
			"name": &FieldConfig{
				Type: String,
			},
		},
	})

	// ie declare that Dog belongs to Pet interface
	_ = NewObject(ObjectConfig{
		Name: "Dog",
		Interfaces: []*Interface{
			petType,
		},
		IsTypeOf: func(value interface{}, info ResolveInfo) bool {
			_, ok := value.(*testDog)
			return ok
		},
		Fields: FieldConfigMap{
			"name": &FieldConfig{
				Type: String,
				Resolve: func(p GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Name
					}
					return nil
				},
			},
			"woofs": &FieldConfig{
				Type: Boolean,
				Resolve: func(p GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Woofs
					}
					return nil
				},
			},
		},
	})
	// ie declare that Cat belongs to Pet interface
	_ = NewObject(ObjectConfig{
		Name: "Cat",
		Interfaces: []*Interface{
			petType,
		},
		IsTypeOf: func(value interface{}, info ResolveInfo) bool {
			_, ok := value.(*testCat)
			return ok
		},
		Fields: FieldConfigMap{
			"name": &FieldConfig{
				Type: String,
				Resolve: func(p GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Name
					}
					return nil
				},
			},
			"meows": &FieldConfig{
				Type: Boolean,
				Resolve: func(p GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Meows
					}
					return nil
				},
			},
		},
	})
	schema, err := NewSchema(SchemaConfig{
		Query: NewObject(ObjectConfig{
			Name: "Query",
			Fields: FieldConfigMap{
				"pets": &FieldConfig{
					Type: NewList(petType),
					Resolve: func(p GQLFRParams) interface{} {
						return []interface{}{
							&testDog{"Odie", true},
							&testCat{"Garfield", false},
						}
					},
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	query := `{
      pets {
        name
        ... on Dog {
          woofs
        }
        ... on Cat {
          meows
        }
      }
    }`

	expected := &Result{
		Data: map[string]interface{}{
			"pets": []interface{}{
				map[string]interface{}{
					"name":  "Odie",
					"woofs": bool(true),
				},
				map[string]interface{}{
					"name":  "Garfield",
					"meows": bool(false),
				},
			},
		},
		Errors: nil,
	}

	resultChannel := make(chan *Result)

	go Graphql(Params{
		Schema:        schema,
		RequestString: query,
	}, resultChannel)
	result := <-resultChannel
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}

func TestIsTypeOfUsedToResolveRuntimeTypeForUnion(t *testing.T) {

	dogType := NewObject(ObjectConfig{
		Name: "Dog",
		IsTypeOf: func(value interface{}, info ResolveInfo) bool {
			_, ok := value.(*testDog)
			return ok
		},
		Fields: FieldConfigMap{
			"name": &FieldConfig{
				Type: String,
				Resolve: func(p GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Name
					}
					return nil
				},
			},
			"woofs": &FieldConfig{
				Type: Boolean,
				Resolve: func(p GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Woofs
					}
					return nil
				},
			},
		},
	})
	catType := NewObject(ObjectConfig{
		Name: "Cat",
		IsTypeOf: func(value interface{}, info ResolveInfo) bool {
			_, ok := value.(*testCat)
			return ok
		},
		Fields: FieldConfigMap{
			"name": &FieldConfig{
				Type: String,
				Resolve: func(p GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Name
					}
					return nil
				},
			},
			"meows": &FieldConfig{
				Type: Boolean,
				Resolve: func(p GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Meows
					}
					return nil
				},
			},
		},
	})
	// ie declare Pet has Dot and Cat object types
	petType := NewUnion(UnionConfig{
		Name: "Pet",
		Types: []*Object{
			dogType, catType,
		},
		ResolveType: func(value interface{}, info ResolveInfo) *Object {
			if _, ok := value.(*testCat); ok {
				return catType
			}
			if _, ok := value.(*testDog); ok {
				return dogType
			}
			return nil
		},
	})
	schema, err := NewSchema(SchemaConfig{
		Query: NewObject(ObjectConfig{
			Name: "Query",
			Fields: FieldConfigMap{
				"pets": &FieldConfig{
					Type: NewList(petType),
					Resolve: func(p GQLFRParams) interface{} {
						return []interface{}{
							&testDog{"Odie", true},
							&testCat{"Garfield", false},
						}
					},
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	query := `{
      pets {
        name
        ... on Dog {
          woofs
        }
        ... on Cat {
          meows
        }
      }
    }`

	expected := &Result{
		Data: map[string]interface{}{
			"pets": []interface{}{
				map[string]interface{}{
					"name":  "Odie",
					"woofs": bool(true),
				},
				map[string]interface{}{
					"name":  "Garfield",
					"meows": bool(false),
				},
			},
		},
		Errors: nil,
	}

	resultChannel := make(chan *Result)

	go Graphql(Params{
		Schema:        schema,
		RequestString: query,
	}, resultChannel)
	result := <-resultChannel

	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}

func TestResolveTypeOnInterfaceYieldsUsefulError(t *testing.T) {

	var dogType *Object
	var catType *Object
	var humanType *Object
	petType := NewInterface(InterfaceConfig{
		Name: "Pet",
		Fields: FieldConfigMap{
			"name": &FieldConfig{
				Type: String,
			},
		},
		ResolveType: func(value interface{}, info ResolveInfo) *Object {
			if _, ok := value.(*testCat); ok {
				return catType
			}
			if _, ok := value.(*testDog); ok {
				return dogType
			}
			if _, ok := value.(*testHuman); ok {
				return humanType
			}
			return nil
		},
	})

	humanType = NewObject(ObjectConfig{
		Name: "Human",
		Fields: FieldConfigMap{
			"name": &FieldConfig{
				Type: String,
				Resolve: func(p GQLFRParams) interface{} {
					if human, ok := p.Source.(*testHuman); ok {
						return human.Name
					}
					return nil
				},
			},
		},
	})
	dogType = NewObject(ObjectConfig{
		Name: "Dog",
		Interfaces: []*Interface{
			petType,
		},
		IsTypeOf: func(value interface{}, info ResolveInfo) bool {
			_, ok := value.(*testDog)
			return ok
		},
		Fields: FieldConfigMap{
			"name": &FieldConfig{
				Type: String,
				Resolve: func(p GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Name
					}
					return nil
				},
			},
			"woofs": &FieldConfig{
				Type: Boolean,
				Resolve: func(p GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Woofs
					}
					return nil
				},
			},
		},
	})
	catType = NewObject(ObjectConfig{
		Name: "Cat",
		Interfaces: []*Interface{
			petType,
		},
		IsTypeOf: func(value interface{}, info ResolveInfo) bool {
			_, ok := value.(*testCat)
			return ok
		},
		Fields: FieldConfigMap{
			"name": &FieldConfig{
				Type: String,
				Resolve: func(p GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Name
					}
					return nil
				},
			},
			"meows": &FieldConfig{
				Type: Boolean,
				Resolve: func(p GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Meows
					}
					return nil
				},
			},
		},
	})
	schema, err := NewSchema(SchemaConfig{
		Query: NewObject(ObjectConfig{
			Name: "Query",
			Fields: FieldConfigMap{
				"pets": &FieldConfig{
					Type: NewList(petType),
					Resolve: func(p GQLFRParams) interface{} {
						return []interface{}{
							&testDog{"Odie", true},
							&testCat{"Garfield", false},
							&testHuman{"Jon"},
						}
					},
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	query := `{
      pets {
        name
        ... on Dog {
          woofs
        }
        ... on Cat {
          meows
        }
      }
    }`

	expected := &Result{
		Data: map[string]interface{}{
			"pets": []interface{}{
				map[string]interface{}{
					"name":  "Odie",
					"woofs": bool(true),
				},
				map[string]interface{}{
					"name":  "Garfield",
					"meows": bool(false),
				},
				nil,
			},
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message:   `Runtime Object type "Human" is not a possible type for "Pet".`,
				Locations: []location.SourceLocation{},
			},
		},
	}

	resultChannel := make(chan *Result)

	go Graphql(Params{
		Schema:        schema,
		RequestString: query,
	}, resultChannel)
	result := <-resultChannel
	if len(result.Errors) == 0 {
		t.Fatalf("wrong result, expected errors: %v, got: %v", len(expected.Errors), len(result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}

func TestResolveTypeOnUnionYieldsUsefulError(t *testing.T) {

	humanType := NewObject(ObjectConfig{
		Name: "Human",
		Fields: FieldConfigMap{
			"name": &FieldConfig{
				Type: String,
				Resolve: func(p GQLFRParams) interface{} {
					if human, ok := p.Source.(*testHuman); ok {
						return human.Name
					}
					return nil
				},
			},
		},
	})
	dogType := NewObject(ObjectConfig{
		Name: "Dog",
		IsTypeOf: func(value interface{}, info ResolveInfo) bool {
			_, ok := value.(*testDog)
			return ok
		},
		Fields: FieldConfigMap{
			"name": &FieldConfig{
				Type: String,
				Resolve: func(p GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Name
					}
					return nil
				},
			},
			"woofs": &FieldConfig{
				Type: Boolean,
				Resolve: func(p GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Woofs
					}
					return nil
				},
			},
		},
	})
	catType := NewObject(ObjectConfig{
		Name: "Cat",
		IsTypeOf: func(value interface{}, info ResolveInfo) bool {
			_, ok := value.(*testCat)
			return ok
		},
		Fields: FieldConfigMap{
			"name": &FieldConfig{
				Type: String,
				Resolve: func(p GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Name
					}
					return nil
				},
			},
			"meows": &FieldConfig{
				Type: Boolean,
				Resolve: func(p GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Meows
					}
					return nil
				},
			},
		},
	})
	petType := NewUnion(UnionConfig{
		Name: "Pet",
		Types: []*Object{
			dogType, catType,
		},
		ResolveType: func(value interface{}, info ResolveInfo) *Object {
			if _, ok := value.(*testCat); ok {
				return catType
			}
			if _, ok := value.(*testDog); ok {
				return dogType
			}
			if _, ok := value.(*testHuman); ok {
				return humanType
			}
			return nil
		},
	})
	schema, err := NewSchema(SchemaConfig{
		Query: NewObject(ObjectConfig{
			Name: "Query",
			Fields: FieldConfigMap{
				"pets": &FieldConfig{
					Type: NewList(petType),
					Resolve: func(p GQLFRParams) interface{} {
						return []interface{}{
							&testDog{"Odie", true},
							&testCat{"Garfield", false},
							&testHuman{"Jon"},
						}
					},
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	query := `{
      pets {
        name
        ... on Dog {
          woofs
        }
        ... on Cat {
          meows
        }
      }
    }`

	expected := &Result{
		Data: map[string]interface{}{
			"pets": []interface{}{
				map[string]interface{}{
					"name":  "Odie",
					"woofs": bool(true),
				},
				map[string]interface{}{
					"name":  "Garfield",
					"meows": bool(false),
				},
				nil,
			},
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormattedError{
				Message:   `Runtime Object type "Human" is not a possible type for "Pet".`,
				Locations: []location.SourceLocation{},
			},
		},
	}

	resultChannel := make(chan *Result)

	go Graphql(Params{
		Schema:        schema,
		RequestString: query,
	}, resultChannel)
	result := <-resultChannel
	if len(result.Errors) == 0 {
		t.Fatalf("wrong result, expected errors: %v, got: %v", len(expected.Errors), len(result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
