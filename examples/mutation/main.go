package main

import (
	"errors"
	"fmt"

	"github.com/graphql-go/graphql"
)

type Todo struct {
	name     string
	complete bool
	id       int
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
		name:     "Finish mutation example",
		complete: false,
		id:       counter,
	}
	todos[t.id] = t

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
							name:     name,
							id:       counter,
							complete: false,
						}
						todos[t.id] = t

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
								t.name = name
							}

							if complete, ok := p.Args["complete"].(bool); ok {
								t.complete = complete
							}

							todos[t.id] = t
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

func printTodos() {
	for _, value := range todos {
		var checkbox string
		if value.complete {
			checkbox = "[x]"
		} else {
			checkbox = "[ ]"
		}
		fmt.Println(checkbox, value.name)
	}
}

func main() {
	fmt.Println("Before mutation:")
	printTodos()

	mutation := `
		mutation _ {
			newTodo: newTodo(name:"Submit pull request") {
				name
			}
		}
	`

	p := graphql.Params{
		Schema:        schema,
		RequestString: mutation,
	}

	_ = graphql.Do(p)

	fmt.Println("After mutation:")
	printTodos()

	mutation = `
		mutation _ {
			updateTodo: updateTodo(id:1, complete:true) {
				name
			}
		}
	`

	p = graphql.Params{
		Schema:        schema,
		RequestString: mutation,
	}

	_ = graphql.Do(p)

	fmt.Println("After mutation:")
	printTodos()
}
