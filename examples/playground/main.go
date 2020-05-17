package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/testutil"
)

func main() {
	GRAPHQL_PORT := 8080
	PLAYGROUND_PORT := 8081
	query := `
  type Query {
    hello: String
  }
`
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		result := graphql.Do(graphql.Params{
			Schema:        testutil.StarWarsSchema,
			RequestString: query,
		})
		json.NewEncoder(w).Encode(result)
	})

	cmd := exec.Command("node", "index.js", query)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GRAPHQL_PORT=%d", GRAPHQL_PORT),       // GRAPHQL_PORT
		fmt.Sprintf("PLAYGROUND_PORT=%d", PLAYGROUND_PORT), // this value is used
	)
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("🚀 GraphQL Express playground server is running on: http://localhost:%d/graphql\n", PLAYGROUND_PORT)
	fmt.Println("Now server is running on port 8080")
	fmt.Println("Test with Get      : curl -g 'http://localhost:8080/graphql?query={hero{name}}'")

	http.ListenAndServe(fmt.Sprintf(":%d", GRAPHQL_PORT), nil)

	// if err := cmd.Run(); err != nil {
	// 	log.Fatal(err)
	// 	fmt.Printf("🚀 GraphQL Express playground server is running on: http://localhost:%d/graphql\n", PLAYGROUND_PORT)
	// }
	// 	// fmt.Printf("in all caps: %q\n", out.String())
}
