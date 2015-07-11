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
	Friends    []string
	AppearsIn  []int
	HomePlanet string
}

func init() {
	Luke = StarWarsChar{
		Id:         "1000",
		Name:       "Luke Skywalker",
		Friends:    []string{"1002", "1003", "2000", "2001"},
		AppearsIn:  []int{4, 5, 6},
		HomePlanet: "Tatooine",
	}
	Vader = StarWarsChar{
		Id:         "1001",
		Name:       "Darth Vader",
		Friends:    []string{"1004"},
		AppearsIn:  []int{4, 5, 6},
		HomePlanet: "Tatooine",
	}
	Han = StarWarsChar{
		Id:        "1002",
		Name:      "Han Solo",
		Friends:   []string{"1000", "1003", "2001"},
		AppearsIn: []int{4, 5, 6},
	}
	Leia = StarWarsChar{
		Id:         "1003",
		Name:       "Leia Organa",
		Friends:    []string{"1000", "1002", "2000", "2001"},
		AppearsIn:  []int{4, 5, 6},
		HomePlanet: "Alderaa",
	}
	Tarkin = StarWarsChar{
		Id:        "1004",
		Name:      "Wilhuff Tarkin",
		Friends:   []string{"1001"},
		AppearsIn: []int{4},
	}
	HumanData = map[int]StarWarsChar{
		1000: Luke,
		1001: Vader,
		1002: Han,
		1003: Leia,
		1004: Tarkin,
	}
	var fields types.GraphQLObjectTypeFields
	queryType := types.GraphQLObjectType{
		Name:   "Query",
		Fields: fields,
	}
	StarWarsSchema = types.GraphQLSchema{
		Query: queryType,
	}
}
