package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/examples/todo/schema"
)

type postData struct {
	Query     string                 `json:"query"`
	Operation string                 `json:"operation"`
	Variables map[string]interface{} `json:"variables"`
}

func main() {
	http.HandleFunc("/graphql", func(w http.ResponseWriter, req *http.Request) {
		var p postData
		if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
			w.WriteHeader(400)
			return
		}
		result := graphql.Do(graphql.Params{
			Context:        req.Context(),
			Schema:         schema.TodoSchema,
			RequestString:  p.Query,
			VariableValues: p.Variables,
			OperationName:  p.Operation,
		})
		if err := json.NewEncoder(w).Encode(result); err != nil {
			fmt.Printf("could not write result to response: %s", err)
		}
	})

	fmt.Println("Now server is running on port 8080")

	fmt.Println("")

	fmt.Println(`Get single todo:
curl \
-X POST \
-H "Content-Type: application/json" \
--data '{ "query": "{ todo(id:\"b\") { id text done } }" }' \
http://localhost:8080/graphql`)

	fmt.Println("")

	fmt.Println(`Create new todo:
curl \
-X POST \
-H "Content-Type: application/json" \
--data '{ "query": "mutation { createTodo(text:\"My New todo\") { id text done } }" }' \
http://localhost:8080/graphql`)

	fmt.Println("")

	fmt.Println(`Update todo:
curl \
-X POST \
-H "Content-Type: application/json" \
--data '{ "query": "mutation { updateTodo(id:\"a\", done: true) { id text done } }" }' \
http://localhost:8080/graphql`)

	fmt.Println("")

	fmt.Println(`Load todo list:
curl \
-X POST \
-H "Content-Type: application/json" \
--data '{ "query": "{ todoList { id text done } }" }' \
http://localhost:8080/graphql`)

	http.ListenAndServe(":8080", nil)
}
