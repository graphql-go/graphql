package graphql_test

import (
	"encoding/json"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/testutil"
	"log"
	"reflect"
	"testing"
)

type Human struct {
	Alive  bool    `json:"alive"`
	Age    int     `json:"age"`
	Weight float64 `json:"weight"`
}

type Person struct {
	Human
	Name    string   `json:"name"`
	Home    Address  `json:"home"`
	Friends []Friend `json:"friends"`
}

type Friend struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
}

var personSource = Person{
	Human: Human{
		Age:    24,
		Weight: 70.1,
		Alive:  true,
	},
	Name:    "John Doe",
	Home:    myaddress,
	Friends: friendSource,
}

var friendSource = []Friend{
	{"Arief", "palembang"},
	{"Al", "semarang"},
}

var myaddress = Address{
	Street: "Jl. G1",
	City:   "Jakarta",
}

func TestBindFields(t *testing.T) {
	personObj := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Person",
		Fields: graphql.BindFields(Person{}),
	})
	fields := graphql.Fields{
		"person": &graphql.Field{
			Type: personObj,
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
				alive
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
