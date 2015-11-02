package graphql

import (
	"encoding/json"
	"reflect"
	"strconv"
	"testing"

	"github.com/chris-ramon/graphql-go/language/ast"
	"github.com/chris-ramon/graphql-go/language/parser"
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
	StarWarsSchema Schema

	humanType *Object
	droidType *Object
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

	episodeEnum := NewEnum(EnumConfig{
		Name:        "Episode",
		Description: "One of the films in the Star Wars Trilogy",
		Values: EnumValueConfigMap{
			"NEWHOPE": &EnumValueConfig{
				Value:       4,
				Description: "Released in 1977.",
			},
			"EMPIRE": &EnumValueConfig{
				Value:       5,
				Description: "Released in 1980.",
			},
			"JEDI": &EnumValueConfig{
				Value:       6,
				Description: "Released in 1983.",
			},
		},
	})

	characterInterface := NewInterface(InterfaceConfig{
		Name:        "Character",
		Description: "A character in the Star Wars Trilogy",
		Fields: FieldConfigMap{
			"id": &FieldConfig{
				Type:        NewNonNull(String),
				Description: "The id of the character.",
			},
			"name": &FieldConfig{
				Type:        String,
				Description: "The name of the character.",
			},
			"appearsIn": &FieldConfig{
				Type:        NewList(episodeEnum),
				Description: "Which movies they appear in.",
			},
		},
		ResolveType: func(value interface{}, info ResolveInfo) *Object {
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
	characterInterface.AddFieldConfig("friends", &FieldConfig{
		Type:        NewList(characterInterface),
		Description: "The friends of the character, or an empty list if they have none.",
	})

	humanType = NewObject(ObjectConfig{
		Name:        "Human",
		Description: "A humanoid creature in the Star Wars universe.",
		Fields: FieldConfigMap{
			"id": &FieldConfig{
				Type:        NewNonNull(String),
				Description: "The id of the human.",
				Resolve: func(p GQLFRParams) interface{} {
					if human, ok := p.Source.(StarWarsChar); ok {
						return human.Id
					}
					return nil
				},
			},
			"name": &FieldConfig{
				Type:        String,
				Description: "The name of the human.",
				Resolve: func(p GQLFRParams) interface{} {
					if human, ok := p.Source.(StarWarsChar); ok {
						return human.Name
					}
					return nil
				},
			},
			"friends": &FieldConfig{
				Type:        NewList(characterInterface),
				Description: "The friends of the human, or an empty list if they have none.",
				Resolve: func(p GQLFRParams) interface{} {
					if human, ok := p.Source.(StarWarsChar); ok {
						return human.Friends
					}
					return []interface{}{}
				},
			},
			"appearsIn": &FieldConfig{
				Type:        NewList(episodeEnum),
				Description: "Which movies they appear in.",
				Resolve: func(p GQLFRParams) interface{} {
					if human, ok := p.Source.(StarWarsChar); ok {
						return human.AppearsIn
					}
					return nil
				},
			},
			"homePlanet": &FieldConfig{
				Type:        String,
				Description: "The home planet of the human, or null if unknown.",
				Resolve: func(p GQLFRParams) interface{} {
					if human, ok := p.Source.(StarWarsChar); ok {
						return human.HomePlanet
					}
					return nil
				},
			},
		},
		Interfaces: []*Interface{
			characterInterface,
		},
	})
	droidType = NewObject(ObjectConfig{
		Name:        "Droid",
		Description: "A mechanical creature in the Star Wars universe.",
		Fields: FieldConfigMap{
			"id": &FieldConfig{
				Type:        NewNonNull(String),
				Description: "The id of the droid.",
				Resolve: func(p GQLFRParams) interface{} {
					if droid, ok := p.Source.(StarWarsChar); ok {
						return droid.Id
					}
					return nil
				},
			},
			"name": &FieldConfig{
				Type:        String,
				Description: "The name of the droid.",
				Resolve: func(p GQLFRParams) interface{} {
					if droid, ok := p.Source.(StarWarsChar); ok {
						return droid.Name
					}
					return nil
				},
			},
			"friends": &FieldConfig{
				Type:        NewList(characterInterface),
				Description: "The friends of the droid, or an empty list if they have none.",
				Resolve: func(p GQLFRParams) interface{} {
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
			"appearsIn": &FieldConfig{
				Type:        NewList(episodeEnum),
				Description: "Which movies they appear in.",
				Resolve: func(p GQLFRParams) interface{} {
					if droid, ok := p.Source.(StarWarsChar); ok {
						return droid.AppearsIn
					}
					return nil
				},
			},
			"primaryFunction": &FieldConfig{
				Type:        String,
				Description: "The primary function of the droid.",
				Resolve: func(p GQLFRParams) interface{} {
					if droid, ok := p.Source.(StarWarsChar); ok {
						return droid.PrimaryFunction
					}
					return nil
				},
			},
		},
		Interfaces: []*Interface{
			characterInterface,
		},
	})

	queryType := NewObject(ObjectConfig{
		Name: "Query",
		Fields: FieldConfigMap{
			"hero": &FieldConfig{
				Type: characterInterface,
				Args: FieldConfigArgument{
					"episode": &ArgumentConfig{
						Description: "If omitted, returns the hero of the whole saga. If " +
							"provided, returns the hero of that particular episode.",
						Type: episodeEnum,
					},
				},
				Resolve: func(p GQLFRParams) (r interface{}) {
					return GetHero(p.Args["episode"])
				},
			},
			"human": &FieldConfig{
				Type: humanType,
				Args: FieldConfigArgument{
					"id": &ArgumentConfig{
						Description: "id of the human",
						Type:        NewNonNull(String),
					},
				},
				Resolve: func(p GQLFRParams) (r interface{}) {
					return GetHuman(p.Args["id"].(int))
				},
			},
			"droid": &FieldConfig{
				Type: droidType,
				Args: FieldConfigArgument{
					"id": &ArgumentConfig{
						Description: "id of the droid",
						Type:        NewNonNull(String),
					},
				},
				Resolve: func(p GQLFRParams) (r interface{}) {
					return GetDroid(p.Args["id"].(int))
				},
			},
		},
	})
	StarWarsSchema, _ = NewSchema(SchemaConfig{
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
func TestExecute(t *testing.T, ep ExecuteParams) *Result {
	resultChannel := make(chan *Result)
	go Execute(ep, resultChannel)
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
