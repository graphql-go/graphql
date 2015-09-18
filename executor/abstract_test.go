package executor_test

import (
	"github.com/chris-ramon/graphql-go"
	"github.com/chris-ramon/graphql-go/errors"
	"github.com/chris-ramon/graphql-go/language/location"
	"github.com/chris-ramon/graphql-go/testutil"
	"github.com/chris-ramon/graphql-go/types"
	"reflect"
	"testing"
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

	var dogType *types.GraphQLObjectType
	var catType *types.GraphQLObjectType

	petType := types.NewGraphQLInterfaceType(types.GraphQLInterfaceTypeConfig{
		Name: "Pet",
		Fields: types.GraphQLFieldConfigMap{
			"name": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
		},
		ResolveType: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
			if _, ok := value.(*testCat); ok {
				return catType
			}
			if _, ok := value.(*testDog); ok {
				return dogType
			}
			return nil
		},
	})

	dogType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "Dog",
		Interfaces: []*types.GraphQLInterfaceType{
			petType,
		},
		IsTypeOf: func(value interface{}, info types.GraphQLResolveInfo) bool {
			_, ok := value.(*testDog)
			return ok
		},
		Fields: types.GraphQLFieldConfigMap{
			"name": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Resolve: func(p types.GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Name
					}
					return nil
				},
			},
			"woofs": &types.GraphQLFieldConfig{
				Type: types.GraphQLBoolean,
				Resolve: func(p types.GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Woofs
					}
					return nil
				},
			},
		},
	})
	catType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "Cat",
		Interfaces: []*types.GraphQLInterfaceType{
			petType,
		},
		IsTypeOf: func(value interface{}, info types.GraphQLResolveInfo) bool {
			_, ok := value.(*testCat)
			return ok
		},
		Fields: types.GraphQLFieldConfigMap{
			"name": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Resolve: func(p types.GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Name
					}
					return nil
				},
			},
			"meows": &types.GraphQLFieldConfig{
				Type: types.GraphQLBoolean,
				Resolve: func(p types.GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Meows
					}
					return nil
				},
			},
		},
	})
	schema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "Query",
			Fields: types.GraphQLFieldConfigMap{
				"pets": &types.GraphQLFieldConfig{
					Type: types.NewGraphQLList(petType),
					Resolve: func(p types.GQLFRParams) interface{} {
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

	expected := &types.GraphQLResult{
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

	resultChannel := make(chan *types.GraphQLResult)

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

	dogType := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "Dog",
		IsTypeOf: func(value interface{}, info types.GraphQLResolveInfo) bool {
			_, ok := value.(*testDog)
			return ok
		},
		Fields: types.GraphQLFieldConfigMap{
			"name": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Resolve: func(p types.GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Name
					}
					return nil
				},
			},
			"woofs": &types.GraphQLFieldConfig{
				Type: types.GraphQLBoolean,
				Resolve: func(p types.GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Woofs
					}
					return nil
				},
			},
		},
	})
	catType := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "Cat",
		IsTypeOf: func(value interface{}, info types.GraphQLResolveInfo) bool {
			_, ok := value.(*testCat)
			return ok
		},
		Fields: types.GraphQLFieldConfigMap{
			"name": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Resolve: func(p types.GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Name
					}
					return nil
				},
			},
			"meows": &types.GraphQLFieldConfig{
				Type: types.GraphQLBoolean,
				Resolve: func(p types.GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Meows
					}
					return nil
				},
			},
		},
	})
	petType := types.NewGraphQLUnionType(types.GraphQLUnionTypeConfig{
		Name: "Pet",
		Types: []*types.GraphQLObjectType{
			dogType, catType,
		},
		ResolveType: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
			if _, ok := value.(*testCat); ok {
				return catType
			}
			if _, ok := value.(*testDog); ok {
				return dogType
			}
			return nil
		},
	})
	schema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "Query",
			Fields: types.GraphQLFieldConfigMap{
				"pets": &types.GraphQLFieldConfig{
					Type: types.NewGraphQLList(petType),
					Resolve: func(p types.GQLFRParams) interface{} {
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

	expected := &types.GraphQLResult{
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

	resultChannel := make(chan *types.GraphQLResult)

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

	var dogType *types.GraphQLObjectType
	var catType *types.GraphQLObjectType
	var humanType *types.GraphQLObjectType
	petType := types.NewGraphQLInterfaceType(types.GraphQLInterfaceTypeConfig{
		Name: "Pet",
		Fields: types.GraphQLFieldConfigMap{
			"name": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
		},
		ResolveType: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
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

	humanType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "Human",
		Fields: types.GraphQLFieldConfigMap{
			"name": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Resolve: func(p types.GQLFRParams) interface{} {
					if human, ok := p.Source.(*testHuman); ok {
						return human.Name
					}
					return nil
				},
			},
		},
	})
	dogType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "Dog",
		Interfaces: []*types.GraphQLInterfaceType{
			petType,
		},
		IsTypeOf: func(value interface{}, info types.GraphQLResolveInfo) bool {
			_, ok := value.(*testDog)
			return ok
		},
		Fields: types.GraphQLFieldConfigMap{
			"name": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Resolve: func(p types.GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Name
					}
					return nil
				},
			},
			"woofs": &types.GraphQLFieldConfig{
				Type: types.GraphQLBoolean,
				Resolve: func(p types.GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Woofs
					}
					return nil
				},
			},
		},
	})
	catType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "Cat",
		Interfaces: []*types.GraphQLInterfaceType{
			petType,
		},
		IsTypeOf: func(value interface{}, info types.GraphQLResolveInfo) bool {
			_, ok := value.(*testCat)
			return ok
		},
		Fields: types.GraphQLFieldConfigMap{
			"name": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Resolve: func(p types.GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Name
					}
					return nil
				},
			},
			"meows": &types.GraphQLFieldConfig{
				Type: types.GraphQLBoolean,
				Resolve: func(p types.GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Meows
					}
					return nil
				},
			},
		},
	})
	schema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "Query",
			Fields: types.GraphQLFieldConfigMap{
				"pets": &types.GraphQLFieldConfig{
					Type: types.NewGraphQLList(petType),
					Resolve: func(p types.GQLFRParams) interface{} {
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

	expected := &types.GraphQLResult{
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
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message:   `Runtime Object type "Human" is not a possible type for "Pet".`,
				Locations: []location.SourceLocation{},
			},
		},
	}

	resultChannel := make(chan *types.GraphQLResult)

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

	humanType := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "Human",
		Fields: types.GraphQLFieldConfigMap{
			"name": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Resolve: func(p types.GQLFRParams) interface{} {
					if human, ok := p.Source.(*testHuman); ok {
						return human.Name
					}
					return nil
				},
			},
		},
	})
	dogType := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "Dog",
		IsTypeOf: func(value interface{}, info types.GraphQLResolveInfo) bool {
			_, ok := value.(*testDog)
			return ok
		},
		Fields: types.GraphQLFieldConfigMap{
			"name": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Resolve: func(p types.GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Name
					}
					return nil
				},
			},
			"woofs": &types.GraphQLFieldConfig{
				Type: types.GraphQLBoolean,
				Resolve: func(p types.GQLFRParams) interface{} {
					if dog, ok := p.Source.(*testDog); ok {
						return dog.Woofs
					}
					return nil
				},
			},
		},
	})
	catType := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "Cat",
		IsTypeOf: func(value interface{}, info types.GraphQLResolveInfo) bool {
			_, ok := value.(*testCat)
			return ok
		},
		Fields: types.GraphQLFieldConfigMap{
			"name": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Resolve: func(p types.GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Name
					}
					return nil
				},
			},
			"meows": &types.GraphQLFieldConfig{
				Type: types.GraphQLBoolean,
				Resolve: func(p types.GQLFRParams) interface{} {
					if cat, ok := p.Source.(*testCat); ok {
						return cat.Meows
					}
					return nil
				},
			},
		},
	})
	petType := types.NewGraphQLUnionType(types.GraphQLUnionTypeConfig{
		Name: "Pet",
		Types: []*types.GraphQLObjectType{
			dogType, catType,
		},
		ResolveType: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
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
	schema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "Query",
			Fields: types.GraphQLFieldConfigMap{
				"pets": &types.GraphQLFieldConfig{
					Type: types.NewGraphQLList(petType),
					Resolve: func(p types.GQLFRParams) interface{} {
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

	expected := &types.GraphQLResult{
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
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message:   `Runtime Object type "Human" is not a possible type for "Pet".`,
				Locations: []location.SourceLocation{},
			},
		},
	}

	resultChannel := make(chan *types.GraphQLResult)

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
