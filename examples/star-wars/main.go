package main

import (
	"encoding/json"
	"fmt"
	"log"
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
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Printf("failed to encode gql result: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	fmt.Println("Now server is running on port 8080")
	fmt.Println("Test with Get      : curl -g 'http://localhost:8080/graphql?query={hero{name}}'")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
