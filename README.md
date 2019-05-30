# graphql [![Build Status](https://travis-ci.org/GannettDigital/graphql.svg?branch=master)](https://travis-ci.org/GannettDigital/graphql) [![GoDoc](https://godoc.org/github.com/GannettDigital/graphql?status.svg)](https://godoc.org/github.com/GannettDigital/graphql) [![Coverage Status](https://coveralls.io/repos/github/GannettDigital/graphql/badge.svg?branch=master)](https://coveralls.io/github/GannettDigital/graphql?branch=master)

An implementation of GraphQL in Go. Follows the official reference implementation [`graphql-js`](https://github.com/graphql/graphql-js).

Supports: queries, mutations & subscriptions.

This is a fork of the original repo [here](https://github.com/graphql-go/graphql). We maintain this fork mostly for speed of development. We've added features the original repo now has, like resolving fields in parallel. We've also added additional improvements, like query complexity costs.

### Documentation

godoc: https://godoc.org/github.com/GannettDigital/graphql

### Getting Started

To install the library, run:
```bash
go get github.com/GannettDigital/graphql
```

The following is a simple example which defines a schema with a single `hello` string-type field and a `Resolve` method which returns the string `world`. A GraphQL query is performed against this schema with the resulting output printed in JSON format.

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/GannettDigital/graphql"
)

func main() {
	// Schema
	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
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
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON) // {“data”:{“hello”:”world”}}
}
```
For more complex examples, refer to the [examples/](https://github.com/GannettDigital/graphql/tree/master/examples/) directory and [graphql_test.go](https://github.com/GannettDigital/graphql/blob/master/graphql_test.go).

### Third Party Libraries
| Name          | Author        | Description  |
|:-------------:|:-------------:|:------------:|
| [graphql-go-handler](https://github.com/graphql-go/graphql-go-handler) | [Hafiz Ismail](https://github.com/sogko) | Middleware to handle GraphQL queries through HTTP requests. |
| [graphql-relay-go](https://github.com/graphql-go/graphql-relay-go) | [Hafiz Ismail](https://github.com/sogko) | Lib to construct a graphql-go server supporting react-relay. |
| [golang-relay-starter-kit](https://github.com/sogko/golang-relay-starter-kit) | [Hafiz Ismail](https://github.com/sogko) | Barebones starting point for a Relay application with Golang GraphQL server. |
| [dataloader](https://github.com/nicksrandall/dataloader) | [Nick Randall](https://github.com/nicksrandall) | [DataLoader](https://github.com/facebook/dataloader) implementation in Go. |

### Blog Posts
- [Golang + GraphQL + Relay](http://wehavefaces.net/)

