package testutil

import (
	"github.com/chris-ramon/graphql-go/executor"
	"github.com/chris-ramon/graphql-go/language/ast"
	"github.com/chris-ramon/graphql-go/language/parser"
	"github.com/chris-ramon/graphql-go/types"
	"github.com/deckarep/golang-set"
	"github.com/kr/pretty"
	"reflect"
	"strconv"
	"testing"
)

var (
	Luke           StarWarsChar
	Vader          StarWarsChar
	Han            StarWarsChar
	Leia           StarWarsChar
	Tarkin         StarWarsChar
	Threepio       StarWarsChar
	Artoo          StarWarsChar
	HumanData      map[int]StarWarsChar
	DroidData      map[int]StarWarsChar
	StarWarsSchema types.GraphQLSchema

	humanType *types.GraphQLObjectType
	droidType *types.GraphQLObjectType
)

type StarWarsChar struct {
	Id              string
	Name            string
	Friends         []StarWarsChar
	AppearsIn       []int
	HomePlanet      string
	PrimaryFunction string
}

func init() {
	Luke = StarWarsChar{
		Id:         "1000",
		Name:       "Luke Skywalker",
		AppearsIn:  []int{4, 5, 6},
		HomePlanet: "Tatooine",
	}
	Vader = StarWarsChar{
		Id:         "1001",
		Name:       "Darth Vader",
		AppearsIn:  []int{4, 5, 6},
		HomePlanet: "Tatooine",
	}
	Han = StarWarsChar{
		Id:        "1002",
		Name:      "Han Solo",
		AppearsIn: []int{4, 5, 6},
	}
	Leia = StarWarsChar{
		Id:         "1003",
		Name:       "Leia Organa",
		AppearsIn:  []int{4, 5, 6},
		HomePlanet: "Alderaa",
	}
	Tarkin = StarWarsChar{
		Id:        "1004",
		Name:      "Wilhuff Tarkin",
		AppearsIn: []int{4},
	}
	Threepio = StarWarsChar{
		Id:              "2000",
		Name:            "C-3PO",
		AppearsIn:       []int{4, 5, 6},
		PrimaryFunction: "Protocol",
	}
	Artoo = StarWarsChar{
		Id:              "2001",
		Name:            "R2-D2",
		AppearsIn:       []int{4, 5, 6},
		PrimaryFunction: "Astromech",
	}
	Luke.Friends = append(Luke.Friends, []StarWarsChar{Han, Leia, Threepio, Artoo}...)
	Vader.Friends = append(Luke.Friends, []StarWarsChar{Tarkin}...)
	Han.Friends = append(Han.Friends, []StarWarsChar{Luke, Leia, Artoo}...)
	Leia.Friends = append(Leia.Friends, []StarWarsChar{Luke, Han, Threepio, Artoo}...)
	Tarkin.Friends = append(Tarkin.Friends, []StarWarsChar{Vader}...)
	Threepio.Friends = append(Threepio.Friends, []StarWarsChar{Luke, Han, Leia, Artoo}...)
	Artoo.Friends = append(Artoo.Friends, []StarWarsChar{Luke, Han, Leia}...)
	HumanData = map[int]StarWarsChar{
		1000: Luke,
		1001: Vader,
		1002: Han,
		1003: Leia,
		1004: Tarkin,
	}
	DroidData = map[int]StarWarsChar{
		2000: Threepio,
		2001: Artoo,
	}

	episodeEnum := types.NewGraphQLEnumType(types.GraphQLEnumTypeConfig{
		Name:        "Episode",
		Description: "One of the films in the Star Wars Trilogy",
		Values: types.GraphQLEnumValueConfigMap{
			"NEWHOPE": types.GraphQLEnumValueConfig{
				Value:       4,
				Description: "Released in 1977.",
			},
			"EMPIRE": types.GraphQLEnumValueConfig{
				Value:       5,
				Description: "Released in 1980.",
			},
			"JEDI": types.GraphQLEnumValueConfig{
				Value:       6,
				Description: "Released in 1983.",
			},
		},
	})

	characterInterface := types.NewGraphQLInterfaceType(types.GraphQLInterfaceTypeConfig{
		Name:        "Character",
		Description: "A character in the Star Wars Trilogy",
		Fields: types.GraphQLFieldConfigMap{
			"id": &types.GraphQLFieldConfig{
				Type:        types.NewGraphQLNonNull(types.GraphQLString),
				Description: "The id of the character.",
			},
			"name": &types.GraphQLFieldConfig{
				Type:        types.GraphQLString,
				Description: "The name of the character.",
			},
			"appearsIn": &types.GraphQLFieldConfig{
				Type:        types.NewGraphQLList(episodeEnum),
				Description: "Which movies they appear in.",
			},
		},
		ResolveType: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
			if character, ok := value.(StarWarsChar); ok {
				id, _ := strconv.Atoi(character.Id)
				human := GetHuman(id)
				if human.Id != "" {
					return humanType
				}
			}
			return droidType
		},
	})
	characterInterface.AddFieldConfig("friends", &types.GraphQLFieldConfig{
		Type:        types.NewGraphQLList(characterInterface),
		Description: "The friends of the character, or an empty list if they have none.",
	})

	humanType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name:        "Human",
		Description: "A humanoid creature in the Star Wars universe.",
		Fields: types.GraphQLFieldConfigMap{
			"id": &types.GraphQLFieldConfig{
				Type:        types.NewGraphQLNonNull(types.GraphQLString),
				Description: "The id of the human.",
				Resolve: func(p types.GQLFRParams) interface{} {
					if human, ok := p.Source.(StarWarsChar); ok {
						return human.Id
					}
					return nil
				},
			},
			"name": &types.GraphQLFieldConfig{
				Type:        types.GraphQLString,
				Description: "The name of the human.",
				Resolve: func(p types.GQLFRParams) interface{} {
					if human, ok := p.Source.(StarWarsChar); ok {
						return human.Name
					}
					return nil
				},
			},
			"friends": &types.GraphQLFieldConfig{
				Type:        types.NewGraphQLList(characterInterface),
				Description: "The friends of the human, or an empty list if they have none.",
				Resolve: func(p types.GQLFRParams) interface{} {
					if human, ok := p.Source.(StarWarsChar); ok {
						return human.Friends
					}
					return []interface{}{}
				},
			},
			"appearsIn": &types.GraphQLFieldConfig{
				Type:        types.NewGraphQLList(episodeEnum),
				Description: "Which movies they appear in.",
				Resolve: func(p types.GQLFRParams) interface{} {
					if human, ok := p.Source.(StarWarsChar); ok {
						return human.AppearsIn
					}
					return nil
				},
			},
			"homePlanet": &types.GraphQLFieldConfig{
				Type:        types.GraphQLString,
				Description: "The home planet of the human, or null if unknown.",
				Resolve: func(p types.GQLFRParams) interface{} {
					if human, ok := p.Source.(StarWarsChar); ok {
						return human.HomePlanet
					}
					return nil
				},
			},
		},
		Interfaces: []*types.GraphQLInterfaceType{
			characterInterface,
		},
	})
	droidType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name:        "Droid",
		Description: "A mechanical creature in the Star Wars universe.",
		Fields: types.GraphQLFieldConfigMap{
			"id": &types.GraphQLFieldConfig{
				Type:        types.NewGraphQLNonNull(types.GraphQLString),
				Description: "The id of the droid.",
				Resolve: func(p types.GQLFRParams) interface{} {
					if droid, ok := p.Source.(StarWarsChar); ok {
						return droid.Id
					}
					return nil
				},
			},
			"name": &types.GraphQLFieldConfig{
				Type:        types.GraphQLString,
				Description: "The name of the droid.",
				Resolve: func(p types.GQLFRParams) interface{} {
					if droid, ok := p.Source.(StarWarsChar); ok {
						return droid.Name
					}
					return nil
				},
			},
			"friends": &types.GraphQLFieldConfig{
				Type:        types.NewGraphQLList(characterInterface),
				Description: "The friends of the droid, or an empty list if they have none.",
				Resolve: func(p types.GQLFRParams) interface{} {
					if droid, ok := p.Source.(StarWarsChar); ok {
						friends := []map[string]interface{}{}
						for _, friend := range droid.Friends {
							friends = append(friends, map[string]interface{}{
								"name": friend.Name,
								"id":   friend.Id,
							})
						}
						return droid.Friends
					}
					return []interface{}{}
				},
			},
			"appearsIn": &types.GraphQLFieldConfig{
				Type:        types.NewGraphQLList(episodeEnum),
				Description: "Which movies they appear in.",
				Resolve: func(p types.GQLFRParams) interface{} {
					if droid, ok := p.Source.(StarWarsChar); ok {
						return droid.AppearsIn
					}
					return nil
				},
			},
			"primaryFunction": &types.GraphQLFieldConfig{
				Type:        types.GraphQLString,
				Description: "The primary function of the droid.",
				Resolve: func(p types.GQLFRParams) interface{} {
					if droid, ok := p.Source.(StarWarsChar); ok {
						return droid.PrimaryFunction
					}
					return nil
				},
			},
		},
		Interfaces: []*types.GraphQLInterfaceType{
			characterInterface,
		},
	})

	queryType := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "Query",
		Fields: types.GraphQLFieldConfigMap{
			"hero": &types.GraphQLFieldConfig{
				Type: characterInterface,
				Args: types.GraphQLFieldConfigArgumentMap{
					"episode": &types.GraphQLArgumentConfig{
						Description: "If omitted, returns the hero of the whole saga. If " +
							"provided, returns the hero of that particular episode.",
						Type: episodeEnum,
					},
				},
				Resolve: func(p types.GQLFRParams) (r interface{}) {
					return GetHero(p.Args["episode"])
				},
			},
			"human": &types.GraphQLFieldConfig{
				Type: humanType,
				Args: types.GraphQLFieldConfigArgumentMap{
					"id": &types.GraphQLArgumentConfig{
						Description: "id of the human",
						Type:        types.NewGraphQLNonNull(types.GraphQLString),
					},
				},
				Resolve: func(p types.GQLFRParams) (r interface{}) {
					return GetHuman(p.Args["id"].(int))
				},
			},
			"droid": &types.GraphQLFieldConfig{
				Type: droidType,
				Args: types.GraphQLFieldConfigArgumentMap{
					"id": &types.GraphQLArgumentConfig{
						Description: "id of the droid",
						Type:        types.NewGraphQLNonNull(types.GraphQLString),
					},
				},
				Resolve: func(p types.GQLFRParams) (r interface{}) {
					return GetDroid(p.Args["id"].(int))
				},
			},
		},
	})
	StarWarsSchema, _ = types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: queryType,
	})
}

func GetHuman(id int) StarWarsChar {
	if human, ok := HumanData[id]; ok {
		return human
	}
	return StarWarsChar{}
}
func GetDroid(id int) StarWarsChar {
	if droid, ok := DroidData[id]; ok {
		return droid
	}
	return StarWarsChar{}
}
func GetHero(episode interface{}) interface{} {
	if episode == 5 {
		return Luke
	}
	return Artoo
}

// Test helper functions

func Parse(t *testing.T, query string) *ast.Document {
	astDoc, err := parser.Parse(parser.ParseParams{
		Source: query,
		Options: parser.ParseOptions{
			// include source, for error reporting
			NoSource: false,
		},
	})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	return astDoc
}
func Execute(t *testing.T, ep executor.ExecuteParams) *types.GraphQLResult {
	resultChannel := make(chan *types.GraphQLResult)
	go executor.Execute(ep, resultChannel)
	result := <-resultChannel
	return result
}

func Diff(a, b interface{}) []string {
	return pretty.Diff(a, b)
}

// TODO: EqualSet
func EqualSet(a, b interface{}) bool {

	aa := []interface{}{}
	bb := []interface{}{}

	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)
	if aVal.Type().Kind() == reflect.Slice {
		for i := 0; i < aVal.Len(); i++ {
			val := aVal.Index(i).Interface()
			aa = append(aa, &val)
		}
	}
	if bVal.Type().Kind() == reflect.Slice {
		for i := 0; i < bVal.Len(); i++ {
			val := bVal.Index(i).Interface()
			bb = append(bb, &val)
		}
	}
	xx := mapset.NewSetFromSlice(aa)
	yy := mapset.NewSetFromSlice(bb)
	return xx.Equal(yy)
}
