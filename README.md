# graphql [![Build Status](https://travis-ci.org/chris-ramon/graphql.svg)](https://travis-ci.org/chris-ramon/graphql) [![GoDoc](https://godoc.org/graphql.co/graphql?status.svg)](https://godoc.org/github.com/chris-ramon/graphql) [![Coverage Status](https://coveralls.io/repos/chris-ramon/graphql/badge.svg?branch=master&service=github)](https://coveralls.io/github/chris-ramon/graphql?branch=master) [![Join the chat at https://gitter.im/chris-ramon/graphql](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/chris-ramon/graphql?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)


A *work-in-progress* implementation of GraphQL for Go.

### Getting Started
Installation
```
go get github.com/chris-ramon/graphql
```

A simple example that defines a schema with a `hello` string field,
it’s resolve function returns a string `world`.
Then a graphql query is perform against that schema, finally
the result is printed as JSON:

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/chris-ramon/graphql"
)

func main() {
	// Schema
	fields := graphql.FieldConfigMap{
		"hello": &graphql.FieldConfig{
			Type: graphql.String,
			Resolve: func(p graphql.GQLFRParams) interface{} {
				return "world"
			},
		},
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
			hello
		}
	`
	params := graphql.Params{Schema: schema, RequestString: query}
	result := make(chan *graphql.Result)
	go graphql.Graphql(params, result)
	r := <-result
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON) // {“data”:{“hello”:”world”}}
}
```

For more complex examples see [examples](https://github.com/chris-ramon/graphql/tree/master/examples/) directory and [graphql tests](https://github.com/chris-ramon/graphql/blob/master/graphql_test.go).

### Origin and Current Direction

This project was originally a port of [v0.4.3](https://github.com/graphql/graphql-js/releases/tag/v0.4.3) of [graphql-js](https://github.com/graphql/graphql-js) (excluding the Validator), which was based on the July 2015 GraphQL specification. `graphql` is currently several versions behind `graphql-js`, however future efforts will be guided directly by the [latest formal GraphQL specification](https://github.com/facebook/graphql/releases) (currently: [October 2015](https://github.com/facebook/graphql/releases/tag/October2015)).

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
  - [ ] Basic Usage (see: [PR-#21](https://github.com/chris-ramon/graphql/pull/21)) 
  - [ ] React/Relay
- [ ] Alpha Release (v0.1)

The `Validator` is optional, per official GraphQL specification, but it would be a useful addition.
