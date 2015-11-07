package graphql_test

import (
	"reflect"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/location"
	"github.com/graphql-go/graphql/testutil"
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

	petType := graphql.NewInterface(graphql.InterfaceConfig{
		Name: "Pet",
		Fields: graphql.FieldConfigMap{
			"name": &graphql.FieldConfig{
				Type: graphql.String,
			},
		},
	})

	// ie declare that Dog belongs to Pet interface
	_ = graphql.NewObject(graphql.ObjectConfig{
		Name: "Dog",
		Interfaces: []*graphql.Interface{
			petType,
		},
		IsTypeOf: func(value interface{}, info graphql.ResolveInfo) bool {
			_, ok := value.(*testDog)
			return ok
		},
		Fields: graphql.FieldConfigMap{
			"name": &graphql.FieldConfig{
				Type: graphql.String,
				Resolve: func(p graphql.GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Name
					}
					return nil
				},
			},
			"woofs": &graphql.FieldConfig{
				Type: graphql.Boolean,
				Resolve: func(p graphql.GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Woofs
					}
					return nil
				},
			},
		},
	})
	// ie declare that Cat belongs to Pet interface
	_ = graphql.NewObject(graphql.ObjectConfig{
		Name: "Cat",
		Interfaces: []*graphql.Interface{
			petType,
		},
		IsTypeOf: func(value interface{}, info graphql.ResolveInfo) bool {
			_, ok := value.(*testCat)
			return ok
		},
		Fields: graphql.FieldConfigMap{
			"name": &graphql.FieldConfig{
				Type: graphql.String,
				Resolve: func(p graphql.GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Name
					}
					return nil
				},
			},
			"meows": &graphql.FieldConfig{
				Type: graphql.Boolean,
				Resolve: func(p graphql.GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Meows
					}
					return nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.FieldConfigMap{
				"pets": &graphql.FieldConfig{
					Type: graphql.NewList(petType),
					Resolve: func(p graphql.GQLFRParams) interface{} {
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

	expected := &graphql.Result{
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

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestIsTypeOfUsedToResolveRuntimeTypeForUnion(t *testing.T) {

	dogType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Dog",
		IsTypeOf: func(value interface{}, info graphql.ResolveInfo) bool {
			_, ok := value.(*testDog)
			return ok
		},
		Fields: graphql.FieldConfigMap{
			"name": &graphql.FieldConfig{
				Type: graphql.String,
				Resolve: func(p graphql.GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Name
					}
					return nil
				},
			},
			"woofs": &graphql.FieldConfig{
				Type: graphql.Boolean,
				Resolve: func(p graphql.GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Woofs
					}
					return nil
				},
			},
		},
	})
	catType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Cat",
		IsTypeOf: func(value interface{}, info graphql.ResolveInfo) bool {
			_, ok := value.(*testCat)
			return ok
		},
		Fields: graphql.FieldConfigMap{
			"name": &graphql.FieldConfig{
				Type: graphql.String,
				Resolve: func(p graphql.GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Name
					}
					return nil
				},
			},
			"meows": &graphql.FieldConfig{
				Type: graphql.Boolean,
				Resolve: func(p graphql.GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Meows
					}
					return nil
				},
			},
		},
	})
	// ie declare Pet has Dot and Cat object types
	petType := graphql.NewUnion(graphql.UnionConfig{
		Name: "Pet",
		Types: []*graphql.Object{
			dogType, catType,
		},
		ResolveType: func(value interface{}, info graphql.ResolveInfo) *graphql.Object {
			if _, ok := value.(*testCat); ok {
				return catType
			}
			if _, ok := value.(*testDog); ok {
				return dogType
			}
			return nil
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.FieldConfigMap{
				"pets": &graphql.FieldConfig{
					Type: graphql.NewList(petType),
					Resolve: func(p graphql.GQLFRParams) interface{} {
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

	expected := &graphql.Result{
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

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})

	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestResolveTypeOnInterfaceYieldsUsefulError(t *testing.T) {

	var dogType *graphql.Object
	var catType *graphql.Object
	var humanType *graphql.Object
	petType := graphql.NewInterface(graphql.InterfaceConfig{
		Name: "Pet",
		Fields: graphql.FieldConfigMap{
			"name": &graphql.FieldConfig{
				Type: graphql.String,
			},
		},
		ResolveType: func(value interface{}, info graphql.ResolveInfo) *graphql.Object {
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

	humanType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Human",
		Fields: graphql.FieldConfigMap{
			"name": &graphql.FieldConfig{
				Type: graphql.String,
				Resolve: func(p graphql.GQLFRParams) interface{} {
					if human, ok := p.Source.(*testHuman); ok {
						return human.Name
					}
					return nil
				},
			},
		},
	})
	dogType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Dog",
		Interfaces: []*graphql.Interface{
			petType,
		},
		IsTypeOf: func(value interface{}, info graphql.ResolveInfo) bool {
			_, ok := value.(*testDog)
			return ok
		},
		Fields: graphql.FieldConfigMap{
			"name": &graphql.FieldConfig{
				Type: graphql.String,
				Resolve: func(p graphql.GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Name
					}
					return nil
				},
			},
			"woofs": &graphql.FieldConfig{
				Type: graphql.Boolean,
				Resolve: func(p graphql.GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Woofs
					}
					return nil
				},
			},
		},
	})
	catType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Cat",
		Interfaces: []*graphql.Interface{
			petType,
		},
		IsTypeOf: func(value interface{}, info graphql.ResolveInfo) bool {
			_, ok := value.(*testCat)
			return ok
		},
		Fields: graphql.FieldConfigMap{
			"name": &graphql.FieldConfig{
				Type: graphql.String,
				Resolve: func(p graphql.GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Name
					}
					return nil
				},
			},
			"meows": &graphql.FieldConfig{
				Type: graphql.Boolean,
				Resolve: func(p graphql.GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Meows
					}
					return nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.FieldConfigMap{
				"pets": &graphql.FieldConfig{
					Type: graphql.NewList(petType),
					Resolve: func(p graphql.GQLFRParams) interface{} {
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

	expected := &graphql.Result{
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

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) == 0 {
		t.Fatalf("wrong result, expected errors: %v, got: %v", len(expected.Errors), len(result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestResolveTypeOnUnionYieldsUsefulError(t *testing.T) {

	humanType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Human",
		Fields: graphql.FieldConfigMap{
			"name": &graphql.FieldConfig{
				Type: graphql.String,
				Resolve: func(p graphql.GQLFRParams) interface{} {
					if human, ok := p.Source.(*testHuman); ok {
						return human.Name
					}
					return nil
				},
			},
		},
	})
	dogType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Dog",
		IsTypeOf: func(value interface{}, info graphql.ResolveInfo) bool {
			_, ok := value.(*testDog)
			return ok
		},
		Fields: graphql.FieldConfigMap{
			"name": &graphql.FieldConfig{
				Type: graphql.String,
				Resolve: func(p graphql.GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Name
					}
					return nil
				},
			},
			"woofs": &graphql.FieldConfig{
				Type: graphql.Boolean,
				Resolve: func(p graphql.GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Woofs
					}
					return nil
				},
			},
		},
	})
	catType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Cat",
		IsTypeOf: func(value interface{}, info graphql.ResolveInfo) bool {
			_, ok := value.(*testCat)
			return ok
		},
		Fields: graphql.FieldConfigMap{
			"name": &graphql.FieldConfig{
				Type: graphql.String,
				Resolve: func(p graphql.GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Name
					}
					return nil
				},
			},
			"meows": &graphql.FieldConfig{
				Type: graphql.Boolean,
				Resolve: func(p graphql.GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Meows
					}
					return nil
				},
			},
		},
	})
	petType := graphql.NewUnion(graphql.UnionConfig{
		Name: "Pet",
		Types: []*graphql.Object{
			dogType, catType,
		},
		ResolveType: func(value interface{}, info graphql.ResolveInfo) *graphql.Object {
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
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.FieldConfigMap{
				"pets": &graphql.FieldConfig{
					Type: graphql.NewList(petType),
					Resolve: func(p graphql.GQLFRParams) interface{} {
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

	expected := &graphql.Result{
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

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) == 0 {
		t.Fatalf("wrong result, expected errors: %v, got: %v", len(expected.Errors), len(result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
