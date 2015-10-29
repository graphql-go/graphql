Graphql Golang [![Build Status](https://travis-ci.org/chris-ramon/graphql-go.svg)](https://travis-ci.org/chris-ramon/graphql-go) [![GoDoc](https://godoc.org/graphql.co/graphql?status.svg)](https://godoc.org/github.com/chris-ramon/graphql-go) [![Coverage Status](https://coveralls.io/repos/chris-ramon/graphql-go/badge.svg?branch=master&service=github)](https://coveralls.io/github/chris-ramon/graphql-go?branch=master) [![Join the chat at https://gitter.im/chris-ramon/graphql-go](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/chris-ramon/graphql-go?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
=====

A *work-in-progress* implementation of GraphQL for Go. Right now a package renaming is in works.

Its very similar to the js implementation and aims to be somewhat API compatabile Server-side implementation of graphql   
[graphql-go](https://github.com/chris-ramon/graphql-go) == [graphql-js](https://github.com/graphql/graphql-js) 

### Origin and Current Direction

This project was originally a port of [v0.4.3](https://github.com/graphql/graphql-js/releases/tag/v0.4.3) of [graphql-js](https://github.com/graphql/graphql-js) (excluding the Validator), which was based on the July 2015 GraphQL specification. `graphql-go` is currently several versions behind `graphql-js`, however future efforts will be guided directly by the [latest formal GraphQL specification](https://github.com/facebook/graphql/releases) (currently: [October 2015](https://github.com/facebook/graphql/releases/tag/October2015)).

### Install
`go get https://github.com/chris-ramon/graphql-go`

### Example
```go
package main

import (
  "encoding/json"
  "github.com/chris-ramon/graphql-go"
  "github.com/chris-ramon/graphql-go/types"
)

// This is a basic grapqhl object type
var UserType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
	Name:        "User",
	Description: "A typical user",
	Fields: types.GraphQLFieldConfigMap{
		"id": &types.GraphQLFieldConfig{
			Description: "The id of the user",
			Type:        types.GraphQLString,
		},
		"name": &types.GraphQLFieldConfig{
			Description: "The name of the user",
			Type:        types.GraphQLString,
		},
		"email": &types.GraphQLFieldConfig{
			Description: "The full name of the user",
			Type:        types.GraphQLString,
		},
	},
})


// This is the type that will be the root of our query,
// and the entry point into our schema.
var RootQuery = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
	Name: "Query",
	Fields: types.GraphQLFieldConfigMap{
		"user": &types.GraphQLFieldConfig{
			Type: UserType,
			Args: types.GraphQLFieldConfigArgumentMap{
				"id": &types.GraphQLArgumentConfig{
					Type: types.GraphQLString,
				},
			},
			Resolve: func(p types.GQLFRParams) interface{} {
				return map[string]interface{}{
				  "id": "john_doe",
				  "name": "John Doe",
				  "email": "john_doe@abc.com",
				}
			},
		},
	},
})


func main() {
        // We create the schema first
    Schema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: RootQuery,
	})
	if err != nil {
		panic(err.Error())
	}
	
	// Then to execute a graphql request
    rootValue := map[string]interface{}{"property_1": "tester"}
	resultChannel := make(chan *types.GraphQLResult)
	params := gql.GraphqlParams{
		Schema:         Schema,
		RootObject:     rootValue,
		RequestString:  "query GetUser { user { id name } }",
		VariableValues: map[string]interface{}{},
		OperationName:  "",
	}
	go gql.Graphql(params, resultChannel)
	
	// Get the response
	result, err := json.Marshal(<-resultChannel)
	if err != nil {
	  panic(err.Error())
	}
	println(string(result))
}
```

### Other Libraries Related
* [graphql-go-handler](https://github.com/sogko/graphql-go-handler) == [express-graphql](https://github.com/graphql/express-graphql)  
Middleware to handle GraphQL queries through HTTP requests. It parses GET/POST params and passes them into Graphql(), which returns JSON response. You can choose not to use it, but you will end up writing similar code. As to whether it should be merged with graphql-go, I think it could possibly be. Or it could remain separate, just like express-graphql.

* [graphql-relay-go](https://github.com/sogko/graphql-relay-go) == [graphql-relay-js](https://github.com/graphql/graphql-relay-js)  
This is a library to construct Relay-compliant servers, which has additional specs for pagination, global IDs and those sort of things. Not needed if you chose to build a pure GraphQL server.

Some of the other projects that [@sogko](https://github.com/sogko) contributed do use graphql-go + graphql-go-handler + graphql-relay-go, but they are specifically Relay applications:

* [golang-relay-starter-kit](https://github.com/sogko/golang-relay-starter-kit)
* [todomvc-relay-go](https://github.com/sogko/todomvc-relay-go)

Another project that [@sogko](https://github.com/sogko) wrote used only graphql-go + graphql-go-handler, with graphiql for the front-end:

* [golang-graphql-playground](https://github.com/sogko/golang-graphql-playground) (Query only example, no mutations)

### Blog Posts that might be useful
A couple of posts written by [@sogko](https://github.com/sogko) on [Golang + GraphQL + Relay](http://wehavefaces.net/) but again, those are heading into the direction of Relay-specific details.

### Contributing

We actively welcome pull requests, learn how to contribute.

### Changelog

Changes are tracked as Github releases.(Todo)

### License
Todo


### Roadmap
- [x] Lexer
- [x] Parser
- [x] Schema Parser
- [x] Printer
- [x] Schema Printer
- [x] Visitor
- [x] Executor
- [ ] Validator
- [ ] Examples
  - [ ] Basic Usage (see: [PR-#21](https://github.com/chris-ramon/graphql-go/pull/21)) 
  - [ ] React/Relay
- [ ] Alpha Release (v0.1)

The `Validator` is optional, per official GraphQL specification, but it would be a useful addition.
