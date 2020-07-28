package graphql_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/graphql-go/graphql"
)

type HelloOutput struct {
	Message string `json:"message"`
}

func Hello(ctx *context.Context) (output *HelloOutput, err error) {
	output = &HelloOutput{
		Message: "Hello World",
	}
	return output, nil
}

type GreetingInput struct {
	Name string `json:"name"`
}

type GreetingOutput struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func GreetingPtr(ctx *context.Context, input *GreetingInput) (output *GreetingOutput, err error) {
	return &GreetingOutput{
		Message:   fmt.Sprintf("Hello %s.", input.Name),
		Timestamp: time.Now(),
	}, nil
}

func Greeting(ctx context.Context, input GreetingInput) (output GreetingOutput, err error) {
	return GreetingOutput{
		Message:   fmt.Sprintf("Hello %s.", input.Name),
		Timestamp: time.Now(),
	}, nil
}

type FriendRecur struct {
	Name    string        `json:"name"`
	Friends []FriendRecur `json:"friends"`
}

func friends(ctx *context.Context) (output *FriendRecur) {
	recursiveFriendRecur := FriendRecur{
		Name: "Recursion",
	}
	recursiveFriendRecur.Friends = make([]FriendRecur, 2)
	recursiveFriendRecur.Friends[0] = recursiveFriendRecur
	recursiveFriendRecur.Friends[1] = recursiveFriendRecur

	return &FriendRecur{
		Name: "Alan",
		Friends: []FriendRecur{
			recursiveFriendRecur,
			{
				Name: "Samantha",
				Friends: []FriendRecur{
					{
						Name: "Olivia",
					},
					{
						Name: "Eric",
					},
				},
			},
			{
				Name: "Brian",
				Friends: []FriendRecur{
					{
						Name: "Windy",
					},
					{
						Name: "Kevin",
					},
				},
			},
			{
				Name: "Kevin",
				Friends: []FriendRecur{
					{
						Name: "Sergei",
					},
					{
						Name: "Michael",
					},
				},
			},
		},
	}
}

func TestBindHappyPath(t *testing.T) {
	// Schema
	fields := graphql.Fields{
		"hello":       graphql.Bind(Hello),
		"greeting":    graphql.Bind(Greeting),
		"greetingPtr": graphql.Bind(GreetingPtr),
		"friends":     graphql.Bind(friends),
		"string":      graphql.Bind("Hello World"),
		"number":      graphql.Bind(12345),
		"float":       graphql.Bind(123.45),
		"anonymous": graphql.Bind(struct {
			SomeField string `json:"someField"`
		}{
			SomeField: "Some Value",
		}),
		"simpleFunc": graphql.Bind(func() string {
			return "Hello World"
		}),
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
			hello {
				message
			}
			greeting(name:"Alan") {
				message
				timestamp
			}
			greetingPtr(name:"Alan") {
				message
				timestamp
			}
			friends {
				name
				friends {
					name
					friends {
						name
						friends {
							name
							friends {
								name
							}
						}
					}
				}
			}
			string
			number
			float
			anonymous {
				someField
			}
			simpleFunc
		}
	`
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		t.Errorf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
}

func TestBindPanicImproperInput(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected Bind to panic due to improper function signature")
		}
	}()
	graphql.Bind(func(a, b, c string) {})
}

func TestBindPanicImproperOutput(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected Bind to panic due to improper function signature")
		}
	}()
	graphql.Bind(func() (string, string) { return "Hello", "World" })
}

func TestBindWithRuntimeError(t *testing.T) {
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: graphql.Fields{
		"throwError": graphql.Bind(func() (string, error) {
			return "", errors.New("Some Error")
		}),
	}}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// Query
	query := `
	{
		throwError
	}
	`
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) == 0 {
		t.Error("Expected error")
	}
}
