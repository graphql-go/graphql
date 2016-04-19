package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type Todo struct {
	Name     string `json:"name"`
	Complete bool   `json:"complete"`
	Id       int    `json:"id"`
}

var (
	todos  map[int]Todo
	schema graphql.Schema
	err    error

	counter int = 0
)

func init() {
	counter++

	// setup a simple Todo
	todos = make(map[int]Todo)
	t := Todo{
		Name:     "Finish mutation example",
		Complete: false,
		Id:       counter,
	}
	todos[t.Id] = t

	todoType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Todo",
		Description: "The todo object",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"complete": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Boolean),
			},
		},
	})

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"newTodo": &graphql.Field{
				Type: todoType,
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if name, ok := p.Args["name"].(string); ok {
						counter++
						t := Todo{
							Name:     name,
							Id:       counter,
							Complete: false,
						}
						todos[t.Id] = t

						return t, nil
					}

					return nil, errors.New("could not get name from params")
				},
			},
			"updateTodo": &graphql.Field{
				Type: todoType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"complete": &graphql.ArgumentConfig{
						Type: graphql.Boolean,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if id, ok := p.Args["id"].(int); ok {
						if t, ok := todos[id]; ok {
							if name, ok := p.Args["name"].(string); ok {
								t.Name = name
							}

							if complete, ok := p.Args["complete"].(bool); ok {
								t.Complete = complete
							}

							todos[t.Id] = t
							return t, nil
						}

						return nil, errors.New("could not find todo with that ID")
					}

					return nil, errors.New("could not get id from params")
				},
			},
		},
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"todo": &graphql.Field{
				Type: todoType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if id, ok := p.Args["id"].(int); ok {
						if t, ok := todos[id]; ok {
							return t, nil
						}

						return nil, errors.New("could not find todo with that ID")
					}

					return nil, errors.New("could not parse ID from query")
				},
			},
			"todos": &graphql.Field{
				Type: graphql.NewList(todoType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var t []Todo

					for _, value := range todos {
						t = append(t, value)
					}

					return t, nil
				},
			},
		},
	})

	schema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})

	if err != nil {
		panic(err)
	}
}

func main() {
	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})

	http.Handle("/graphql", h)

	port := ":8080"
	log.Printf("Setting up server on localhost%s", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Failed to listen on port %s", port)
	}
}
