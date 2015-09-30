package executor_test

import (
	"reflect"
	"testing"

	"github.com/chris-ramon/graphql-go"
	"github.com/chris-ramon/graphql-go/errors"
	"github.com/chris-ramon/graphql-go/gqltypes"
	"github.com/chris-ramon/graphql-go/language/location"
	"github.com/chris-ramon/graphql-go/testutil"
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

	petType := gqltypes.NewGraphQLInterfaceType(gqltypes.GraphQLInterfaceTypeConfig{
		Name: "Pet",
		Fields: gqltypes.GraphQLFieldConfigMap{
			"name": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLString,
			},
		},
	})

	// ie declare that Dog belongs to Pet interface
	_ = gqltypes.NewGraphQLObjectType(gqltypes.GraphQLObjectTypeConfig{
		Name: "Dog",
		Interfaces: []*gqltypes.GraphQLInterfaceType{
			petType,
		},
		IsTypeOf: func(value interface{}, info gqltypes.GraphQLResolveInfo) bool {
			_, ok := value.(*testDog)
			return ok
		},
		Fields: gqltypes.GraphQLFieldConfigMap{
			"name": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLString,
				Resolve: func(p gqltypes.GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Name
					}
					return nil
				},
			},
			"woofs": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLBoolean,
				Resolve: func(p gqltypes.GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Woofs
					}
					return nil
				},
			},
		},
	})
	// ie declare that Cat belongs to Pet interface
	_ = gqltypes.NewGraphQLObjectType(gqltypes.GraphQLObjectTypeConfig{
		Name: "Cat",
		Interfaces: []*gqltypes.GraphQLInterfaceType{
			petType,
		},
		IsTypeOf: func(value interface{}, info gqltypes.GraphQLResolveInfo) bool {
			_, ok := value.(*testCat)
			return ok
		},
		Fields: gqltypes.GraphQLFieldConfigMap{
			"name": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLString,
				Resolve: func(p gqltypes.GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Name
					}
					return nil
				},
			},
			"meows": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLBoolean,
				Resolve: func(p gqltypes.GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Meows
					}
					return nil
				},
			},
		},
	})
	schema, err := gqltypes.NewGraphQLSchema(gqltypes.GraphQLSchemaConfig{
		Query: gqltypes.NewGraphQLObjectType(gqltypes.GraphQLObjectTypeConfig{
			Name: "Query",
			Fields: gqltypes.GraphQLFieldConfigMap{
				"pets": &gqltypes.GraphQLFieldConfig{
					Type: gqltypes.NewGraphQLList(petType),
					Resolve: func(p gqltypes.GQLFRParams) interface{} {
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

	expected := &gqltypes.GraphQLResult{
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

	resultChannel := make(chan *gqltypes.GraphQLResult)

	go gql.Graphql(gql.GraphqlParams{
		Schema:        schema,
		RequestString: query,
	}, resultChannel)
	result := <-resultChannel
	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestIsTypeOfUsedToResolveRuntimeTypeForUnion(t *testing.T) {

	dogType := gqltypes.NewGraphQLObjectType(gqltypes.GraphQLObjectTypeConfig{
		Name: "Dog",
		IsTypeOf: func(value interface{}, info gqltypes.GraphQLResolveInfo) bool {
			_, ok := value.(*testDog)
			return ok
		},
		Fields: gqltypes.GraphQLFieldConfigMap{
			"name": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLString,
				Resolve: func(p gqltypes.GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Name
					}
					return nil
				},
			},
			"woofs": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLBoolean,
				Resolve: func(p gqltypes.GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Woofs
					}
					return nil
				},
			},
		},
	})
	catType := gqltypes.NewGraphQLObjectType(gqltypes.GraphQLObjectTypeConfig{
		Name: "Cat",
		IsTypeOf: func(value interface{}, info gqltypes.GraphQLResolveInfo) bool {
			_, ok := value.(*testCat)
			return ok
		},
		Fields: gqltypes.GraphQLFieldConfigMap{
			"name": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLString,
				Resolve: func(p gqltypes.GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Name
					}
					return nil
				},
			},
			"meows": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLBoolean,
				Resolve: func(p gqltypes.GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Meows
					}
					return nil
				},
			},
		},
	})
	// ie declare Pet has Dot and Cat object types
	petType := gqltypes.NewGraphQLUnionType(gqltypes.GraphQLUnionTypeConfig{
		Name: "Pet",
		Types: []*gqltypes.GraphQLObjectType{
			dogType, catType,
		},
		ResolveType: func(value interface{}, info gqltypes.GraphQLResolveInfo) *gqltypes.GraphQLObjectType {
			if _, ok := value.(*testCat); ok {
				return catType
			}
			if _, ok := value.(*testDog); ok {
				return dogType
			}
			return nil
		},
	})
	schema, err := gqltypes.NewGraphQLSchema(gqltypes.GraphQLSchemaConfig{
		Query: gqltypes.NewGraphQLObjectType(gqltypes.GraphQLObjectTypeConfig{
			Name: "Query",
			Fields: gqltypes.GraphQLFieldConfigMap{
				"pets": &gqltypes.GraphQLFieldConfig{
					Type: gqltypes.NewGraphQLList(petType),
					Resolve: func(p gqltypes.GQLFRParams) interface{} {
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

	expected := &gqltypes.GraphQLResult{
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

	resultChannel := make(chan *gqltypes.GraphQLResult)

	go gql.Graphql(gql.GraphqlParams{
		Schema:        schema,
		RequestString: query,
	}, resultChannel)
	result := <-resultChannel

	if len(result.Errors) != 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestResolveTypeOnInterfaceYieldsUsefulError(t *testing.T) {

	var dogType *gqltypes.GraphQLObjectType
	var catType *gqltypes.GraphQLObjectType
	var humanType *gqltypes.GraphQLObjectType
	petType := gqltypes.NewGraphQLInterfaceType(gqltypes.GraphQLInterfaceTypeConfig{
		Name: "Pet",
		Fields: gqltypes.GraphQLFieldConfigMap{
			"name": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLString,
			},
		},
		ResolveType: func(value interface{}, info gqltypes.GraphQLResolveInfo) *gqltypes.GraphQLObjectType {
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

	humanType = gqltypes.NewGraphQLObjectType(gqltypes.GraphQLObjectTypeConfig{
		Name: "Human",
		Fields: gqltypes.GraphQLFieldConfigMap{
			"name": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLString,
				Resolve: func(p gqltypes.GQLFRParams) interface{} {
					if human, ok := p.Source.(*testHuman); ok {
						return human.Name
					}
					return nil
				},
			},
		},
	})
	dogType = gqltypes.NewGraphQLObjectType(gqltypes.GraphQLObjectTypeConfig{
		Name: "Dog",
		Interfaces: []*gqltypes.GraphQLInterfaceType{
			petType,
		},
		IsTypeOf: func(value interface{}, info gqltypes.GraphQLResolveInfo) bool {
			_, ok := value.(*testDog)
			return ok
		},
		Fields: gqltypes.GraphQLFieldConfigMap{
			"name": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLString,
				Resolve: func(p gqltypes.GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Name
					}
					return nil
				},
			},
			"woofs": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLBoolean,
				Resolve: func(p gqltypes.GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Woofs
					}
					return nil
				},
			},
		},
	})
	catType = gqltypes.NewGraphQLObjectType(gqltypes.GraphQLObjectTypeConfig{
		Name: "Cat",
		Interfaces: []*gqltypes.GraphQLInterfaceType{
			petType,
		},
		IsTypeOf: func(value interface{}, info gqltypes.GraphQLResolveInfo) bool {
			_, ok := value.(*testCat)
			return ok
		},
		Fields: gqltypes.GraphQLFieldConfigMap{
			"name": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLString,
				Resolve: func(p gqltypes.GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Name
					}
					return nil
				},
			},
			"meows": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLBoolean,
				Resolve: func(p gqltypes.GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Meows
					}
					return nil
				},
			},
		},
	})
	schema, err := gqltypes.NewGraphQLSchema(gqltypes.GraphQLSchemaConfig{
		Query: gqltypes.NewGraphQLObjectType(gqltypes.GraphQLObjectTypeConfig{
			Name: "Query",
			Fields: gqltypes.GraphQLFieldConfigMap{
				"pets": &gqltypes.GraphQLFieldConfig{
					Type: gqltypes.NewGraphQLList(petType),
					Resolve: func(p gqltypes.GQLFRParams) interface{} {
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

	expected := &gqltypes.GraphQLResult{
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
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message:   `Runtime Object type "Human" is not a possible type for "Pet".`,
				Locations: []location.SourceLocation{},
			},
		},
	}

	resultChannel := make(chan *gqltypes.GraphQLResult)

	go gql.Graphql(gql.GraphqlParams{
		Schema:        schema,
		RequestString: query,
	}, resultChannel)
	result := <-resultChannel
	if len(result.Errors) == 0 {
		t.Fatalf("wrong result, expected errors: %v, got: %v", len(expected.Errors), len(result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestResolveTypeOnUnionYieldsUsefulError(t *testing.T) {

	humanType := gqltypes.NewGraphQLObjectType(gqltypes.GraphQLObjectTypeConfig{
		Name: "Human",
		Fields: gqltypes.GraphQLFieldConfigMap{
			"name": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLString,
				Resolve: func(p gqltypes.GQLFRParams) interface{} {
					if human, ok := p.Source.(*testHuman); ok {
						return human.Name
					}
					return nil
				},
			},
		},
	})
	dogType := gqltypes.NewGraphQLObjectType(gqltypes.GraphQLObjectTypeConfig{
		Name: "Dog",
		IsTypeOf: func(value interface{}, info gqltypes.GraphQLResolveInfo) bool {
			_, ok := value.(*testDog)
			return ok
		},
		Fields: gqltypes.GraphQLFieldConfigMap{
			"name": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLString,
				Resolve: func(p gqltypes.GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Name
					}
					return nil
				},
			},
			"woofs": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLBoolean,
				Resolve: func(p gqltypes.GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Woofs
					}
					return nil
				},
			},
		},
	})
	catType := gqltypes.NewGraphQLObjectType(gqltypes.GraphQLObjectTypeConfig{
		Name: "Cat",
		IsTypeOf: func(value interface{}, info gqltypes.GraphQLResolveInfo) bool {
			_, ok := value.(*testCat)
			return ok
		},
		Fields: gqltypes.GraphQLFieldConfigMap{
			"name": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLString,
				Resolve: func(p gqltypes.GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Name
					}
					return nil
				},
			},
			"meows": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLBoolean,
				Resolve: func(p gqltypes.GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Meows
					}
					return nil
				},
			},
		},
	})
	petType := gqltypes.NewGraphQLUnionType(gqltypes.GraphQLUnionTypeConfig{
		Name: "Pet",
		Types: []*gqltypes.GraphQLObjectType{
			dogType, catType,
		},
		ResolveType: func(value interface{}, info gqltypes.GraphQLResolveInfo) *gqltypes.GraphQLObjectType {
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
	schema, err := gqltypes.NewGraphQLSchema(gqltypes.GraphQLSchemaConfig{
		Query: gqltypes.NewGraphQLObjectType(gqltypes.GraphQLObjectTypeConfig{
			Name: "Query",
			Fields: gqltypes.GraphQLFieldConfigMap{
				"pets": &gqltypes.GraphQLFieldConfig{
					Type: gqltypes.NewGraphQLList(petType),
					Resolve: func(p gqltypes.GQLFRParams) interface{} {
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

	expected := &gqltypes.GraphQLResult{
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
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message:   `Runtime Object type "Human" is not a possible type for "Pet".`,
				Locations: []location.SourceLocation{},
			},
		},
	}

	resultChannel := make(chan *gqltypes.GraphQLResult)

	go gql.Graphql(gql.GraphqlParams{
		Schema:        schema,
		RequestString: query,
	}, resultChannel)
	result := <-resultChannel
	if len(result.Errors) == 0 {
		t.Fatalf("wrong result, expected errors: %v, got: %v", len(expected.Errors), len(result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
