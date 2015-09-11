package testutil

import "github.com/chris-ramon/graphql-go/types"

var (
	Luke           StarWarsChar
	Vader          StarWarsChar
	Han            StarWarsChar
	Leia           StarWarsChar
	Tarkin         StarWarsChar
	HumanData      map[int]StarWarsChar
	StarWarsSchema types.GraphQLSchema
)

type StarWarsChar struct {
	Id         string
	Name       string
	Friends    []StarWarsChar
	AppearsIn  []int
	HomePlanet string
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
	Luke.Friends = append(Luke.Friends, []StarWarsChar{Han, Leia}...)
	Vader.Friends = append(Luke.Friends, []StarWarsChar{Tarkin}...)
	Han.Friends = append(Han.Friends, []StarWarsChar{Leia}...)
	Leia.Friends = append(Leia.Friends, []StarWarsChar{Han}...)
	Tarkin.Friends = append(Tarkin.Friends, []StarWarsChar{Vader}...)
	HumanData = map[int]StarWarsChar{
		1000: Luke,
		1001: Vader,
		1002: Han,
		1003: Leia,
		1004: Tarkin,
	}
	charIntFields := types.GraphQLFieldDefinitionMap{}
	characterInterface := types.GraphQLInterfaceType{
		Name:        "Character",
		Description: "A character in the Star Wars Trilogy",
		Fields:      charIntFields,
	}
	fields := types.GraphQLFieldDefinitionMap{}
	fields["hero"] = types.GraphQLFieldDefinition{
		Type: &characterInterface,
		Resolve: func(p types.GQLFRParams) (r interface{}) {
			return r
		},
	}
	queryType := types.GraphQLObjectType{
		Name:   "Query",
		Fields: fields,
	}
	StarWarsSchema = types.GraphQLSchema{
		Query: queryType,
	}
}
