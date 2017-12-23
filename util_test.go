package graphql_test

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/testutil"
)

type Person struct {
	Human
	Name    string   `json:"name"`
	Home    Address  `json:"home"`
	Hobbies []string `json:"hobbies"`
	Friends []Friend `json:"friends"`
}

type Human struct {
	Alive  bool    `json:"alive,omitempty"`
	Age    int     `json:"age"`
	Weight float64 `json:"weight"`
}

type Friend struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
	Test   string `json:",omitempty"`
}

var personSource = Person{
	Human: Human{
		Age:    24,
		Weight: 70.1,
		Alive:  true,
	},
	Name: "John Doe",
	Home: Address{
		Street: "Jl. G1",
		City:   "Jakarta",
	},
	Friends: friendSource,
	Hobbies: []string{"eat", "sleep", "code"},
}

var friendSource = []Friend{
	{Name: "Arief", Address: "palembang"},
	{Name: "Al", Address: "semarang"},
}

func TestBindFields(t *testing.T) {
	// create person type based on Person struct
	personType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Person",
		// pass empty Person struct to bind all of it's fields
		Fields: graphql.BindFields(Person{}),
	})
	fields := graphql.Fields{
		"person": &graphql.Field{
			Type: personType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return personSource, nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// Query
	query := `
		{
			person{
				name,
				home{street,city},
				friends{name,address},
				age,
				weight,
				alive,
				hobbies
			}
		}
	`
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}

	rJSON, _ := json.Marshal(r)
	data := struct {
		Data struct {
			Person Person `json:"person"`
		} `json:"data"`
	}{}
	err = json.Unmarshal(rJSON, &data)
	if err != nil {
		log.Fatalf("failed to unmarshal. error: %v", err)
	}

	newPerson := data.Data.Person
	if !reflect.DeepEqual(newPerson, personSource) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(personSource, newPerson))
	}
}

func TestBindArg(t *testing.T) {
	var friendObj = graphql.NewObject(graphql.ObjectConfig{
		Name:   "friend",
		Fields: graphql.BindFields(Friend{}),
	})

	fields := graphql.Fields{
		"friend": &graphql.Field{
			Type: friendObj,
			//it can be added more than one since it's a slice
			Args: graphql.BindArg(Friend{}, "name"),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if name, ok := p.Args["name"].(string); ok {
					for _, friend := range friendSource {
						if friend.Name == name {
							return friend, nil
						}
					}
				}
				return nil, nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// Query
	query := `
		{
			friend(name:"Arief"){
				address
			}
		}
	`
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}

	rJSON, _ := json.Marshal(r)

	data := struct {
		Data struct {
			Friend Friend `json:"friend"`
		} `json:"data"`
	}{}
	err = json.Unmarshal(rJSON, &data)
	if err != nil {
		log.Fatalf("failed to unmarshal. error: %v", err)
	}

	expectedAddress := "palembang"
	newFriend := data.Data.Friend
	if newFriend.Address != expectedAddress {
		t.Fatalf("Unexpected result, expected address to be %s but got %s", expectedAddress, newFriend.Address)
	}
}
