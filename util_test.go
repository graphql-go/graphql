package graphql

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

type Human struct {
	Alive  bool    `json:"alive"`
	Age    int     `json:"age"`
	Weight float64 `json:"weight"`
}

type person struct {
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

var personSource = person{
	Human: Human{
		Age:    24,
		Weight: 70.1,
		Alive:  true,
	},
	Name: "John Doe",
	Home: myaddress,
	Friends: []Friend{
		{"Arief", "palembang"},
		{"Al", "semarang"},
	},
}

var myaddress = Address{
	Street: "Jl. G1",
	City:   "Jakarta",
}

func TestBindFields(t *testing.T) {
	personObj := NewObject(ObjectConfig{
		Name:   "Person",
		Fields: BindFields(person{}),
	})
	fields := Fields{
		"person": &Field{
			Type: personObj,
			Resolve: func(p ResolveParams) (interface{}, error) {
				return personSource, nil
			},
		},
	}
	rootQuery := ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := SchemaConfig{Query: NewObject(rootQuery)}
	schema, err := NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// Query
	query := `
		{
			person{
				name,
				home{street},
				friends{name,address},
				age,
				weight,
				alive
			}
		}
	`
	params := Params{Schema: schema, RequestString: query}
	r := Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON)
}
