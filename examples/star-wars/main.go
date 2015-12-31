package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/testutil"
)

func main() {
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()["query"][0]
		result := graphql.Do(graphql.Params{
			Schema:        testutil.StarWarsSchema,
			RequestString: query,
		})
		err := json.NewEncoder(w).Encode(result)
		if err != nil {
			fmt.Printf("Error encoding result: %v", err)
		}
	})
	fmt.Println("Now server is running on port 8080")
	fmt.Println("Test with Get      : curl -g 'http://localhost:8080/graphql?query={hero{name}}'")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error opening HTTP server: %v", err)
	}
}
