package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/graphql-go/graphql"
)

type Todo struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

var TodoList []Todo
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func init() {
	todo1 := Todo{ID: "a", Text: "A todo not to forget", Done: false}
	todo2 := Todo{ID: "b", Text: "This is the most important", Done: false}
	todo3 := Todo{ID: "c", Text: "Please do this or else", Done: false}
	TodoList = append(TodoList, todo1, todo2, todo3)

	rand.Seed(time.Now().UnixNano())
}

// define custom GraphQL ObjectType `todoType` for our Golang struct `Todo`
// Note that
// - the fields in our todoType maps with the json tags for the fields in our struct
// - the field type matches the field type in our struct
var todoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Todo",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"text": &graphql.Field{
			Type: graphql.String,
		},
		"done": &graphql.Field{
			Type: graphql.Boolean,
		},
	},
})

// root mutation
var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		/*
			curl -g 'http://localhost:8080/graphql?query=mutation+_{createTodo(text:"My+new+todo"){id,text,done}}'
		*/
		"createTodo": &graphql.Field{
			Type:        todoType, // the return type for this field
			Description: "Create new todo",
			Args: graphql.FieldConfigArgument{
				"text": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				// marshall and cast the argument value
				text, _ := params.Args["text"].(string)

				// figure out new id
				newID := RandStringRunes(8)

				// perform mutation operation here
				// for e.g. create a Todo and save to DB.
				newTodo := Todo{
					ID:   newID,
					Text: text,
					Done: false,
				}

				TodoList = append(TodoList, newTodo)

				// return the new Todo object that we supposedly save to DB
				// Note here that
				// - we are returning a `Todo` struct instance here
				// - we previously specified the return Type to be `todoType`
				// - `Todo` struct maps to `todoType`, as defined in `todoType` ObjectConfig`
				return newTodo, nil
			},
		},
		/*
			curl -g 'http://localhost:8080/graphql?query=mutation+_{updateTodo(id:"a",done:true){id,text,done}}'
		*/
		"updateTodo": &graphql.Field{
			Type:        todoType, // the return type for this field
			Description: "Update existing todo, mark it done or not done",
			Args: graphql.FieldConfigArgument{
				"done": &graphql.ArgumentConfig{
					Type: graphql.Boolean,
				},
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				// marshall and cast the argument value
				done, _ := params.Args["done"].(bool)
				id, _ := params.Args["id"].(string)
				affectedTodo := Todo{}

				// Search list for todo with id and change the done variable
				for i := 0; i < len(TodoList); i++ {
					if TodoList[i].ID == id {
						TodoList[i].Done = done
						// Assign updated todo so we can return it
						affectedTodo = TodoList[i]
						break
					}
				}
				// Return affected todo
				return affectedTodo, nil
			},
		},
	},
})

// root query
// we just define a trivial example here, since root query is required.
// Test with curl
// curl -g 'http://localhost:8080/graphql?query={lastTodo{id,text,done}}'
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{

		/*
		   curl -g 'http://localhost:8080/graphql?query={todo(id:"b"){id,text,done}}'
		*/
		"todo": &graphql.Field{
			Type:        todoType,
			Description: "Get single todo",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				idQuery, isOK := params.Args["id"].(string)
				if isOK {
					// Search for el with id
					for _, todo := range TodoList {
						if todo.ID == idQuery {
							return todo, nil
						}
					}
				}

				return Todo{}, nil
			},
		},

		"lastTodo": &graphql.Field{
			Type:        todoType,
			Description: "Last todo added",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return TodoList[len(TodoList)-1], nil
			},
		},

		/*
		   curl -g 'http://localhost:8080/graphql?query={todoList{id,text,done}}'
		*/
		"todoList": &graphql.Field{
			Type:        graphql.NewList(todoType),
			Description: "List of todos",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return TodoList, nil
			},
		},
	},
})

// define schema, with our rootQuery and rootMutation
var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

func main() {
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})
	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	// Display some basic instructions
	fmt.Println("Now server is running on port 8080")
	fmt.Println("Get single todo: curl -g 'http://localhost:8080/graphql?query={todo(id:\"b\"){id,text,done}}'")
	fmt.Println("Create new todo: curl -g 'http://localhost:8080/graphql?query=mutation+_{createTodo(text:\"My+new+todo\"){id,text,done}}'")
	fmt.Println("Update todo: curl -g 'http://localhost:8080/graphql?query=mutation+_{updateTodo(id:\"a\",done:true){id,text,done}}'")
	fmt.Println("Load todo list: curl -g 'http://localhost:8080/graphql?query={todoList{id,text,done}}'")
	fmt.Println("Access the web app via browser at 'http://localhost:8080'")

	http.ListenAndServe(":8080", nil)
}
