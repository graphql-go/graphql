package testutil

import (
	"encoding/json"
	"reflect"
	"strconv"
	"testing"

	"github.com/chris-ramon/graphql"
	"github.com/chris-ramon/graphql/language/ast"
	"github.com/chris-ramon/graphql/language/parser"
	"github.com/kr/pretty"
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
	StarWarsSchema graphql.Schema

	humanType *graphql.Object
	droidType *graphql.Object
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

	episodeEnum := graphql.NewEnum(graphql.EnumConfig{
		Name:        "Episode",
		Description: "One of the films in the Star Wars Trilogy",
		Values: graphql.EnumValueConfigMap{
			"NEWHOPE": &graphql.EnumValueConfig{
				Value:       4,
				Description: "Released in 1977.",
			},
			"EMPIRE": &graphql.EnumValueConfig{
				Value:       5,
				Description: "Released in 1980.",
			},
			"JEDI": &graphql.EnumValueConfig{
				Value:       6,
				Description: "Released in 1983.",
			},
		},
	})

	characterInterface := graphql.NewInterface(graphql.InterfaceConfig{
		Name:        "Character",
		Description: "A character in the Star Wars Trilogy",
		Fields: graphql.FieldConfigMap{
			"id": &graphql.FieldConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The id of the character.",
			},
			"name": &graphql.FieldConfig{
				Type:        graphql.String,
				Description: "The name of the character.",
			},
			"appearsIn": &graphql.FieldConfig{
				Type:        graphql.NewList(episodeEnum),
				Description: "Which movies they appear in.",
			},
		},
		ResolveType: func(value interface{}, info graphql.ResolveInfo) *graphql.Object {
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
	characterInterface.AddFieldConfig("friends", &graphql.FieldConfig{
		Type:        graphql.NewList(characterInterface),
		Description: "The friends of the character, or an empty list if they have none.",
	})

	humanType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Human",
		Description: "A humanoid creature in the Star Wars universe.",
		Fields: graphql.FieldConfigMap{
			"id": &graphql.FieldConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The id of the human.",
				Resolve: func(p graphql.GQLFRParams) interface{} {
					if human, ok := p.Source.(StarWarsChar); ok {
						return human.Id
					}
					return nil
				},
			},
			"name": &graphql.FieldConfig{
				Type:        graphql.String,
				Description: "The name of the human.",
				Resolve: func(p graphql.GQLFRParams) interface{} {
					if human, ok := p.Source.(StarWarsChar); ok {
						return human.Name
					}
					return nil
				},
			},
			"friends": &graphql.FieldConfig{
				Type:        graphql.NewList(characterInterface),
				Description: "The friends of the human, or an empty list if they have none.",
				Resolve: func(p graphql.GQLFRParams) interface{} {
					if human, ok := p.Source.(StarWarsChar); ok {
						return human.Friends
					}
					return []interface{}{}
				},
			},
			"appearsIn": &graphql.FieldConfig{
				Type:        graphql.NewList(episodeEnum),
				Description: "Which movies they appear in.",
				Resolve: func(p graphql.GQLFRParams) interface{} {
					if human, ok := p.Source.(StarWarsChar); ok {
						return human.AppearsIn
					}
					return nil
				},
			},
			"homePlanet": &graphql.FieldConfig{
				Type:        graphql.String,
				Description: "The home planet of the human, or null if unknown.",
				Resolve: func(p graphql.GQLFRParams) interface{} {
					if human, ok := p.Source.(StarWarsChar); ok {
						return human.HomePlanet
					}
					return nil
				},
			},
		},
		Interfaces: []*graphql.Interface{
			characterInterface,
		},
	})
	droidType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Droid",
		Description: "A mechanical creature in the Star Wars universe.",
		Fields: graphql.FieldConfigMap{
			"id": &graphql.FieldConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The id of the droid.",
				Resolve: func(p graphql.GQLFRParams) interface{} {
					if droid, ok := p.Source.(StarWarsChar); ok {
						return droid.Id
					}
					return nil
				},
			},
			"name": &graphql.FieldConfig{
				Type:        graphql.String,
				Description: "The name of the droid.",
				Resolve: func(p graphql.GQLFRParams) interface{} {
					if droid, ok := p.Source.(StarWarsChar); ok {
						return droid.Name
					}
					return nil
				},
			},
			"friends": &graphql.FieldConfig{
				Type:        graphql.NewList(characterInterface),
				Description: "The friends of the droid, or an empty list if they have none.",
				Resolve: func(p graphql.GQLFRParams) interface{} {
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
			"appearsIn": &graphql.FieldConfig{
				Type:        graphql.NewList(episodeEnum),
				Description: "Which movies they appear in.",
				Resolve: func(p graphql.GQLFRParams) interface{} {
					if droid, ok := p.Source.(StarWarsChar); ok {
						return droid.AppearsIn
					}
					return nil
				},
			},
			"primaryFunction": &graphql.FieldConfig{
				Type:        graphql.String,
				Description: "The primary function of the droid.",
				Resolve: func(p graphql.GQLFRParams) interface{} {
					if droid, ok := p.Source.(StarWarsChar); ok {
						return droid.PrimaryFunction
					}
					return nil
				},
			},
		},
		Interfaces: []*graphql.Interface{
			characterInterface,
		},
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.FieldConfigMap{
			"hero": &graphql.FieldConfig{
				Type: characterInterface,
				Args: graphql.FieldConfigArgument{
					"episode": &graphql.ArgumentConfig{
						Description: "If omitted, returns the hero of the whole saga. If " +
							"provided, returns the hero of that particular episode.",
						Type: episodeEnum,
					},
				},
				Resolve: func(p graphql.GQLFRParams) (r interface{}) {
					return GetHero(p.Args["episode"])
				},
			},
			"human": &graphql.FieldConfig{
				Type: humanType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Description: "id of the human",
						Type:        graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.GQLFRParams) (r interface{}) {
					return GetHuman(p.Args["id"].(int))
				},
			},
			"droid": &graphql.FieldConfig{
				Type: droidType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Description: "id of the droid",
						Type:        graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.GQLFRParams) (r interface{}) {
					return GetDroid(p.Args["id"].(int))
				},
			},
		},
	})
	StarWarsSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
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

func TestParse(t *testing.T, query string) *ast.Document {
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
func TestExecute(t *testing.T, ep graphql.ExecuteParams) *graphql.Result {
	resultChannel := make(chan *graphql.Result)
	go graphql.Execute(ep, resultChannel)
	result := <-resultChannel
	return result
}

func Diff(a, b interface{}) []string {
	return pretty.Diff(a, b)
}

func ASTToJSON(t *testing.T, a ast.Node) interface{} {
	b, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("Failed to marshal Node %v", err)
	}
	var f interface{}
	err = json.Unmarshal(b, &f)
	if err != nil {
		t.Fatalf("Failed to unmarshal Node %v", err)
	}
	return f
}

func ContainSubsetSlice(super []interface{}, sub []interface{}) bool {
	if len(sub) == 0 {
		return true
	}
subLoop:
	for _, subVal := range sub {
		found := false
	innerLoop:
		for _, superVal := range super {
			if subVal, ok := subVal.(map[string]interface{}); ok {
				if superVal, ok := superVal.(map[string]interface{}); ok {
					if ContainSubset(superVal, subVal) {
						found = true
						break innerLoop
					} else {
						continue
					}
				} else {
					return false
				}

			}
			if subVal, ok := subVal.([]interface{}); ok {
				if superVal, ok := superVal.([]interface{}); ok {
					if ContainSubsetSlice(superVal, subVal) {
						found = true
						break innerLoop
					} else {
						continue
					}
				} else {
					return false
				}
			}
			if reflect.DeepEqual(superVal, subVal) {
				found = true
				break innerLoop
			}
		}
		if !found {
			return false
		} else {
			continue subLoop
		}
	}
	return true
}

func ContainSubset(super map[string]interface{}, sub map[string]interface{}) bool {
	if len(sub) == 0 {
		return true
	}
	for subKey, subVal := range sub {
		if superVal, ok := super[subKey]; ok {
			switch superVal := superVal.(type) {
			case []interface{}:
				if subVal, ok := subVal.([]interface{}); ok {
					if !ContainSubsetSlice(superVal, subVal) {
						return false
					}
				} else {
					return false
				}
			case map[string]interface{}:
				if subVal, ok := subVal.(map[string]interface{}); ok {
					if !ContainSubset(superVal, subVal) {
						return false
					}
				} else {
					return false
				}
			default:
				if !reflect.DeepEqual(superVal, subVal) {
					return false
				}
			}
		} else {
			return false
		}
	}
	return true
}
