package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
)

type Foo struct {
	Name string
}

var FieldFooType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Foo",
	Fields: graphql.Fields{
		"name": &graphql.Field{Type: graphql.String},
	},
})

type Bar struct {
	Name string
}

var FieldBarType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Bar",
	Fields: graphql.Fields{
		"name": &graphql.Field{Type: graphql.String},
	},
})

// QueryType fields: `concurrentFieldFoo` and `concurrentFieldBar` are resolved
// concurrently because they belong to the same field-level and their `Resolve`
// function returns a function (thunk).
var QueryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"concurrentFieldFoo": &graphql.Field{
			Type: FieldFooType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var foo = Foo{Name: "Foo's name"}
				return func() (interface{}, error) {
					return &foo, nil
				}, nil
			},
		},
		"concurrentFieldBar": &graphql.Field{
			Type: FieldBarType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				type result struct {
					data interface{}
					err  error
				}
				ch := make(chan *result, 1)
				go func() {
					defer close(ch)
					bar := &Bar{Name: "Bar's name"}
					// graphql.Do will finish immediately in the case concurrentFieldFoo returns an error. Therefore,
					// when using goroutines make sure to utilize a done channel to avoid leaking goroutines.
					select {
					case ch <- &result{data: bar, err: nil}:
					case <-p.Context.Done():
					}
				}()
				return func() (interface{}, error) {
					select {
					case r := <-ch:
						return r.data, r.err
					case <-p.Context.Done():
						return nil, nil
					}
				}, nil
			},
		},
	},
})

func main() {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: QueryType,
	})
	if err != nil {
		log.Fatal(err)
	}
	query := `
		query {
			concurrentFieldFoo {
				name
			}
			concurrentFieldBar {
				name
			}
		}
	`
	ctx, cancel := context.WithCancel(context.Background())
	result := graphql.Do(graphql.Params{
		RequestString: query,
		Schema:        schema,
		Context:       ctx,
	})
	cancel() // Notify via a done channel that the request is over.
	b, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", b)
	/*
		{
		  "data": {
		    "concurrentFieldBar": {
		      "name": "Bar's name"
		    },
		    "concurrentFieldFoo": {
		      "name": "Foo's name"
		    }
		  }
		}
	*/
}
