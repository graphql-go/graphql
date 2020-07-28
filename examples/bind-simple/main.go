package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
)

type GreetingInput struct {
	Name string `json:"name"`
}

func Greeting(input GreetingInput) string {
	return fmt.Sprintf("Hello %s", input.Name)
}

func main() {
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: graphql.Fields{
		"greeting": graphql.Bind(Greeting),
	}}

	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// Query
	query := `
	{
		greeting(name: "Alan")
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
