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
	Name  string `json:"name"`
	Woofs bool   `json:"woofs"`
}

type testCat struct {
	Name  string `json:"name"`
	Meows bool   `json:"meows"`
}

type testHuman struct {
	Name string `json:"name"`
}

func TestIsTypeOfUsedToResolveRuntimeTypeForInterface(t *testing.T) {

	petType := graphql.NewInterface(graphql.InterfaceConfig{
		Name: "Pet",
		Fields: graphql.Fields{
			"name": &graphql.Field{
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
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Name, nil
					}
					return nil, nil
				},
			},
			"woofs": &graphql.Field{
				Type: graphql.Boolean,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Woofs, nil
					}
					return nil, nil
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
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Name, nil
					}
					return nil, nil
				},
			},
			"meows": &graphql.Field{
				Type: graphql.Boolean,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Meows, nil
					}
					return nil, nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"pets": &graphql.Field{
					Type: graphql.NewList(petType),
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return []interface{}{
							&testDog{"Odie", true},
							&testCat{"Garfield", false},
						}, nil
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
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"woofs": &graphql.Field{
				Type: graphql.Boolean,
			},
		},
	})
	catType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Cat",
		IsTypeOf: func(value interface{}, info graphql.ResolveInfo) bool {
			_, ok := value.(*testCat)
			return ok
		},
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"meows": &graphql.Field{
				Type: graphql.Boolean,
			},
		},
	})
	// ie declare Pet has Dot and Cat object types
	petType := graphql.NewUnion(graphql.UnionConfig{
		Name: "Pet",
		Types: []*graphql.Object{
			dogType, catType,
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"pets": &graphql.Field{
					Type: graphql.NewList(petType),
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return []interface{}{
							&testDog{"Odie", true},
							&testCat{"Garfield", false},
						}, nil
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
        ... on Dog {
          name
          woofs
        }
        ... on Cat {
          name
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
		Fields: graphql.Fields{
			"name": &graphql.Field{
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
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if human, ok := p.Source.(*testHuman); ok {
						return human.Name, nil
					}
					return nil, nil
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
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Name, nil
					}
					return nil, nil
				},
			},
			"woofs": &graphql.Field{
				Type: graphql.Boolean,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Woofs, nil
					}
					return nil, nil
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
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Name, nil
					}
					return nil, nil
				},
			},
			"meows": &graphql.Field{
				Type: graphql.Boolean,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Meows, nil
					}
					return nil, nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"pets": &graphql.Field{
					Type: graphql.NewList(petType),
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return []interface{}{
							&testDog{"Odie", true},
							&testCat{"Garfield", false},
							&testHuman{"Jon"},
						}, nil
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
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
	dogType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Dog",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"woofs": &graphql.Field{
				Type: graphql.Boolean,
			},
		},
	})
	catType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Cat",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"meows": &graphql.Field{
				Type: graphql.Boolean,
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
			Fields: graphql.Fields{
				"pets": &graphql.Field{
					Type: graphql.NewList(petType),
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return []interface{}{
							&testDog{"Odie", true},
							&testCat{"Garfield", false},
							&testHuman{"Jon"},
						}, nil
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
        ... on Dog {
          name
          woofs
        }
        ... on Cat {
          name
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
