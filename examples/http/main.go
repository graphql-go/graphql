package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/chris-ramon/graphql-go/types"
	"github.com/sogko/graphql-go-handler"
)

/*
   Create User object type with fields "id" and "name" by using GraphQLObjectTypeConfig:
       - Name: name of object type
       - Fields: a map of fields by using GraphQLFieldConfigMap
   Setup type of field use GraphQLFieldConfig
*/
var userType = types.NewGraphQLObjectType(
	types.GraphQLObjectTypeConfig{
		Name: "User",
		Fields: types.GraphQLFieldConfigMap{
			"id": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
			"name": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
		},
	},
)

/*
   Create Query object type with fields "user" has type [userType] by using GraphQLObjectTypeConfig:
       - Name: name of object type
       - Fields: a map of fields by using GraphQLFieldConfigMap
   Setup type of field use GraphQLFieldConfig to define:
       - Type: type of field
       - Args: arguments to query with current field
       - Resolve: function to query data using params from [Args] and return value with current type
*/
var queryType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
	Name: "Query",
	Fields: types.GraphQLFieldConfigMap{
		"user": &types.GraphQLFieldConfig{
			Type: userType,
			Args: types.GraphQLFieldConfigArgumentMap{
				"id": &types.GraphQLArgumentConfig{
					Type: types.GraphQLString,
				},
			},
			Resolve: func(p types.GQLFRParams) interface{} {
				idQuery, isOK := p.Args["id"].(string)
				if isOK {
					return data[idQuery]
				} else {
					return nil
				}
			},
		},
	},
})

var schema, _ = types.NewGraphQLSchema(
	types.GraphQLSchemaConfig{
		Query: queryType,
	},
)

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

//Helper function to import json from file to map
func importJsonDataFromFile(fileName string, result interface{}) (isOK bool) {
	isOK = true
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Print("Error:", err)
		isOK = false
	}
	err = json.Unmarshal(content, result)
	if err != nil {
		isOK = false
		fmt.Print("Error:", err)
	}
	return
}

var data map[string]User

func main() {
	_ = importJsonDataFromFile("data.json", &data)
	// create a graphl-go HTTP handler with our previously defined schema
	// and we also set it to return pretty JSON output
	h := gqlhandler.New(&gqlhandler.Config{
		Schema: &schema,
		Pretty: true,
	})

	// serve a GraphQL endpoint at `/graphql`
	http.Handle("/graphql", h)

	fmt.Println("Now server is running on port 8080")
	fmt.Println("Test with Post: curl -XPOST http://localhost:8080/graphql -H 'Content-Type: application/graphql'  -d 'query Root{ user(id:\"1\"){name}  }'")
	fmt.Println("Test with Get: http://localhost:8080/graphql?query={user(id:%221%22){name}}")
	// and serve!
	http.ListenAndServe(":8080", nil)

}
