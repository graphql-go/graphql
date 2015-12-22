package main

/**
 * Use the cannonical example from Facebook: https://github.com/graphql/graphql-js/tree/master/src/__tests__
 * Most comments are taken directly from their source
 */

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	// external
	"github.com/graphql-go/graphql"
)

type Character struct {
	Id              string   `json:"id"`
	Name            string   `json:"name"`
	Friends         []string `json:"friends"`
	AppearsIn       []int    `json:"appearsIn"`
	HomePlanet      string   `json:"homePlanet"`
	PrimaryFunction string   `json:"primaryFunction"`
}

type Human struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	Friends    []string `json:"friends"`
	AppearsIn  []int    `json:"appearsIn"`
	HomePlanet string   `json:"homePlanet"`
}

type Droid struct {
	Id              string   `json:"id"`
	Name            string   `json:"name"`
	Friends         []string `json:"friends"`
	AppearsIn       []int    `json:"appearsIn"`
	PrimaryFunction string   `json:"primaryFunction"`
}

type Data struct {
	Droids map[string]json.RawMessage `json:"droids"`
	Humans map[string]json.RawMessage `json:"humans"`
}

var (
	episodeEnum        *graphql.Enum
	characterInterface *graphql.Interface
	droidType          *graphql.Object
	humanType          *graphql.Object
	StarWarsSchema     graphql.Schema
	Droids             map[string]Droid
	Humans             map[string]Human
	Characters         map[string]Character
)

func init() {
	/**
	 * Read our JSON file
	 **/
	json_data, err := ioutil.ReadFile("data.json")
	if err != nil {
		log.Fatal("Couldn't read JSON from data.json: ", err)
	}

	/**
	 * Init our data
	 **/
	d := Data{}
	err = json.Unmarshal(json_data, &d)
	if err != nil {
		log.Fatal("Couldn't unmarshal JSON data: ", err)
	}

	Droids = make(map[string]Droid)
	Humans = make(map[string]Human)
	Characters = make(map[string]Character)

	for id, json_data := range d.Droids {
		droid := Droid{}
		character := Character{}
		err = json.Unmarshal(json_data, &droid)
		err = json.Unmarshal(json_data, &character)
		if err != nil {
			log.Fatal("Couldn't unmarshal JSON data for droid '", id, "': ", err)
		}
		Droids[id] = droid
		Characters[id] = character
	}

	for id, json_data := range d.Humans {
		human := Human{}
		character := Character{}
		err = json.Unmarshal(json_data, &human)
		err = json.Unmarshal(json_data, &character)
		if err != nil {
			log.Fatal("Couldn't unmarshal JSON data for human '", id, "': ", err)
		}
		Humans[id] = human
		Characters[id] = character
	}

	/**
	 * This is designed to be an end-to-end test, demonstrating
	 * the full GraphQL stack.
	 *
	 * We will create a GraphQL schema that describes the major
	 * characters in the original Star Wars trilogy.
	 *
	 * NOTE: This may contain spoilers for the original Star
	 * Wars trilogy.
	 */

	/**
	 * Using our shorthand to describe type systems, the type system for our
	 * Star Wars example is:
	 *
	 * enum Episode { NEWHOPE, EMPIRE, JEDI }
	 *
	 * interface Character {
	 *   id: String!
	 *   name: String
	 *   friends: [Character]
	 *   appearsIn: [Episode]
	 * }
	 *
	 * type Human : Character {
	 *   id: String!
	 *   name: String
	 *   friends: [Character]
	 *   appearsIn: [Episode]
	 *   homePlanet: String
	 * }
	 *
	 * type Droid : Character {
	 *   id: String!
	 *   name: String
	 *   friends: [Character]
	 *   appearsIn: [Episode]
	 *   primaryFunction: String
	 * }
	 *
	 * type Query {
	 *   hero(episode: Episode): Character
	 *   human(id: String!): Human
	 *   droid(id: String!): Droid
	 * }
	 *
	 * We begin by setting up our schema.
	 */

	/**
	 * The original trilogy consists of three movies.
	 *
	 * This implements the following type system shorthand:
	 *   enum Episode { NEWHOPE, EMPIRE, JEDI }
	 */
	var episodeValues = &graphql.EnumValueConfigMap{
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
	}

	episodeEnum = graphql.NewEnum(
		graphql.EnumConfig{
			Name:        "Episode",
			Description: "One of the films in the Star Wars Trilogy",
			Values:      *episodeValues,
		},
	)

	/**
	 * Characters in the Star Wars trilogy are either humans or droids.
	 *
	 * This implements the following type system shorthand:
	 *   interface Character {
	 *     id: String!
	 *     name: String
	 *     friends: [Character]
	 *     appearsIn: [Episode]
	 *   }
	 */
	characterInterface = graphql.NewInterface(
		graphql.InterfaceConfig{
			Name:        "Character",
			Description: "A character in the Star Wars Trilogy",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.NewNonNull(graphql.String),
				},
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"appearsIn": &graphql.Field{
					Type: graphql.NewList(episodeEnum),
				},
			},
			ResolveType: func(value interface{}, info graphql.ResolveInfo) *graphql.Object {
				if character, ok := value.(Character); ok {
					human := FindHuman(character.Id)
					if human != nil {
						return humanType
					}
				} else if character, ok := value.(Human); ok {
					human := FindHuman(character.Id)
					if human != nil {
						return humanType
					}
				}
				return droidType
			},
		},
	)

	characterInterface.AddFieldConfig("friends", &graphql.Field{
		Type: graphql.NewList(characterInterface),
	})

	/**
	 * We define our human type, which implements the character interface.
	 *
	 * This implements the following type system shorthand:
	 *   type Human : Character {
	 *     id: String!
	 *     name: String
	 *     friends: [Character]
	 *     appearsIn: [Episode]
	 *   }
	 */
	humanType = graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Human",
			Description: "A humanoid creature in the Star Wars universe.",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.NewNonNull(graphql.String),
				},
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"friends": &graphql.Field{
					Type:        graphql.NewList(characterInterface),
					Description: "The firends of the human, or an empty list if they have none",
					Resolve: func(p graphql.ResolveParams) interface{} {
						// we have to try to resolve the incoming object as either
						// a Character or  Human
						if human, ok := p.Source.(Human); ok {
							friends := make([]interface{}, len(human.Friends))
							for index, friendId := range human.Friends {
								friends[index] = FindCharacter(friendId)
							}
							return friends
						} else if human, ok := p.Source.(Character); ok {
							friends := make([]interface{}, len(human.Friends))
							for index, friendId := range human.Friends {
								friends[index] = FindCharacter(friendId)
							}
							return friends
						}
						return []interface{}{}
					},
				},
				"appearsIn": &graphql.Field{
					Type: graphql.NewList(episodeEnum),
				},
				"homePlanet": &graphql.Field{
					Type: graphql.String,
				},
			},
			Interfaces: []*graphql.Interface{
				characterInterface,
			},
		},
	)

	/**
	 * The other type of character in Star Wars is a droid.
	 *
	 * This implements the following type system shorthand:
	 *   type Droid : Character {
	 *     id: String!
	 *     name: String
	 *     friends: [Character]
	 *     appearsIn: [Episode]
	 *     primaryFunction: String
	 *   }
	 */
	droidType = graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Droid",
			Description: "A mechanical creature in the Star Wars universe.",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.NewNonNull(graphql.String),
				},
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"friends": &graphql.Field{
					Type:        graphql.NewList(characterInterface),
					Description: "The firends of the droid, or an empty list if they have none",
					Resolve: func(p graphql.ResolveParams) interface{} {
						// we have to try to resolve the incoming object as
						// either a Droid or a Character
						if droid, ok := p.Source.(Droid); ok {
							friends := make([]interface{}, len(droid.Friends))
							for index, friendId := range droid.Friends {
								friends[index] = FindCharacter(friendId)
							}
							return friends
						} else if droid, ok := p.Source.(Character); ok {
							friends := make([]interface{}, len(droid.Friends))
							for index, friendId := range droid.Friends {
								friends[index] = FindCharacter(friendId)
							}
							return friends
						}
						return []interface{}{}
					},
				},
				"appearsIn": &graphql.Field{
					Type: graphql.NewList(episodeEnum),
				},
				"primaryFunction": &graphql.Field{
					Type: graphql.String,
				},
			},
			Interfaces: []*graphql.Interface{
				characterInterface,
			},
		},
	)

	/**
	 * This is the type that will be the root of our query, and the
	 * entry point into our schema. It gives us the ability to fetch
	 * objects by their IDs, as well as to fetch the undisputed hero
	 * of the Star Wars trilogy, R2-D2, directly.
	 *
	 * This implements the following type system shorthand:
	 *   type Query {
	 *     hero(episode: Episode): Character
	 *     human(id: String!): Human
	 *     droid(id: String!): Droid
	 *   }
	 *
	 */
	rootQuery := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"hero": &graphql.Field{
					Type: characterInterface,
					Args: graphql.FieldConfigArgument{
						"episode": &graphql.ArgumentConfig{
							Type: episodeEnum,
						},
					},
					Resolve: func(p graphql.ResolveParams) interface{} {
						if p.Args["episode"] == 5 {
							return FindHuman("1000")
						}
						return FindDroid("2001")
					},
				},
				"human": &graphql.Field{
					Type: humanType,
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: func(p graphql.ResolveParams) interface{} {
						if id, ok := p.Args["id"].(string); ok {
							return FindHuman(id)
						}
						return nil
					},
				},
				"droid": &graphql.Field{
					Type: droidType,
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: func(p graphql.ResolveParams) interface{} {
						if id, ok := p.Args["id"].(string); ok {
							return FindDroid(id)
						}
						return nil
					},
				},
			},
		},
	)

	StarWarsSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})
}

func FindHuman(id string) interface{} {
	human := Humans[id]
	if human.Name != "" {
		return interface{}(human)
	}
	return nil
}

func FindCharacter(id string) interface{} {
	character := Characters[id]
	if character.Name != "" {
		return interface{}(character)
	}
	return nil
}

func FindDroid(id string) interface{} {
	droid := Droids[id]
	if droid.Name != "" {
		return interface{}(droid)
	}
	return nil
}

func main() {
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()["query"][0]
		result := graphql.Do(graphql.Params{
			Schema:        StarWarsSchema,
			RequestString: query,
		})
		json.NewEncoder(w).Encode(result)
	})

	fmt.Println("Now server is running on port 8080")
	fmt.Println("Test with Get      : curl -g 'http://localhost:8080/graphql?query={hero{name}}")
	http.ListenAndServe(":8080", nil)
}
