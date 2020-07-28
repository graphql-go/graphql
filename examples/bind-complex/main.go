package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
)

var people = []Person{
	{
		Name: "Alan",
		Friends: []Person{
			{
				Name: "Nadeem",
				Friends: []Person{
					{
						Name: "Heidi",
					},
				},
			},
		},
	},
}

type Person struct {
	Name    string   `json:"name"`
	Friends []Person `json:"friends"`
}

type GetPersonInput struct {
	Name string `json:"name"`
}

type GetPersonOutput struct {
	Person
}

func GetPerson(ctx context.Context, input GetPersonInput) (*GetPersonOutput, error) {
	for _, person := range people {
		if person.Name == input.Name {
			return &GetPersonOutput{
				Person: person,
			}, nil
		}
	}
	return nil, errors.New("Could not find person.")
}

func main() {
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: graphql.Fields{
		"person": graphql.Bind(GetPerson),
	}}

	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// Query
	query := `
	{
		person(name: "Alan") {
			name
			friends {
				name
				friends {
					name
				}
			}
		}
	}
	`
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON)
}
