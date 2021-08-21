package graphql_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/testutil"
)

func TestSchemaSubscribe(t *testing.T) {

	testutil.RunSubscribes(t, []*testutil.TestSubscription{
		{
			Name: "subscribe without resolver",
			Schema: makeSubscriptionSchema(t, graphql.ObjectConfig{
				Name: "Subscription",
				Fields: graphql.Fields{
					"sub_without_resolver": &graphql.Field{
						Type: graphql.String,
						Subscribe: makeSubscribeToMapFunction([]map[string]interface{}{
							{
								"sub_without_resolver": "a",
							},
							{
								"sub_without_resolver": "b",
							},
							{
								"sub_without_resolver": "c",
							},
						}),
					},
				},
			}),
			Query: `
				subscription {
					sub_without_resolver
				}
			`,
			ExpectedResults: []testutil.TestResponse{
				{Data: `{ "sub_without_resolver": "a" }`},
				{Data: `{ "sub_without_resolver": "b" }`},
				{Data: `{ "sub_without_resolver": "c" }`},
			},
		},
		{
			Name: "subscribe with resolver",
			Schema: makeSubscriptionSchema(t, graphql.ObjectConfig{
				Name: "Subscription",
				Fields: graphql.Fields{
					"sub_with_resolver": &graphql.Field{
						Type: graphql.String,
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							return p.Source, nil
						},
						Subscribe: makeSubscribeToStringFunction([]string{"a", "b", "c"}),
					},
				},
			}),
			Query: `
				subscription {
					sub_with_resolver
				}
			`,
			ExpectedResults: []testutil.TestResponse{
				{Data: `{ "sub_with_resolver": "a" }`},
				{Data: `{ "sub_with_resolver": "b" }`},
				{Data: `{ "sub_with_resolver": "c" }`},
			},
		},
		{
			Name: "receive query validation error",
			Schema: makeSubscriptionSchema(t, graphql.ObjectConfig{
				Name: "Subscription",
				Fields: graphql.Fields{
					"sub_without_resolver": &graphql.Field{
						Type:      graphql.String,
						Subscribe: makeSubscribeToStringFunction([]string{"a", "b", "c"}),
					},
				},
			}),
			Query: `
				subscription {
					sub_without_resolver
					xxx
				}
			`,
			ExpectedResults: []testutil.TestResponse{
				{Errors: []string{"Cannot query field \"xxx\" on type \"Subscription\"."}},
			},
		},
		{
			Name: "panic inside subscribe is recovered",
			Schema: makeSubscriptionSchema(t, graphql.ObjectConfig{
				Name: "Subscription",
				Fields: graphql.Fields{
					"should_error": &graphql.Field{
						Type: graphql.String,
						Subscribe: func(p graphql.ResolveParams) (interface{}, error) {
							panic(errors.New("got a panic error"))
						},
					},
				},
			}),
			Query: `
				subscription {
					should_error
				}
			`,
			ExpectedResults: []testutil.TestResponse{
				{Errors: []string{"got a panic error"}},
			},
		},
		{
			Name: "subscribe with resolver changes output",
			Schema: makeSubscriptionSchema(t, graphql.ObjectConfig{
				Name: "Subscription",
				Fields: graphql.Fields{
					"sub_with_resolver": &graphql.Field{
						Type:      graphql.String,
						Subscribe: makeSubscribeToStringFunction([]string{"a", "b", "c", "d"}),
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							return fmt.Sprintf("result=%v", p.Source), nil
						},
					},
				},
			}),
			Query: `
				subscription {
					sub_with_resolver
				}
			`,
			ExpectedResults: []testutil.TestResponse{
				{Data: `{ "sub_with_resolver": "result=a" }`},
				{Data: `{ "sub_with_resolver": "result=b" }`},
				{Data: `{ "sub_with_resolver": "result=c" }`},
				{Data: `{ "sub_with_resolver": "result=d" }`},
			},
		},
		{
			Name: "subscribe to a nested object",
			Schema: makeSubscriptionSchema(t, graphql.ObjectConfig{
				Name: "Subscription",
				Fields: graphql.Fields{
					"sub_with_object": &graphql.Field{
						Type: graphql.NewObject(graphql.ObjectConfig{
							Name: "Obj",
							Fields: graphql.Fields{
								"field": &graphql.Field{
									Type: graphql.String,
								},
							},
						}),
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							return p.Source, nil
						},
						Subscribe: makeSubscribeToMapFunction([]map[string]interface{}{
							{
								"field": "hello",
							},
							{
								"field": "bye",
							},
							{
								"field": nil,
							},
						}),
					},
				},
			}),
			Query: `
				subscription {
					sub_with_object {
						field
					}
				}
			`,
			ExpectedResults: []testutil.TestResponse{
				{Data: `{ "sub_with_object": { "field": "hello" } }`},
				{Data: `{ "sub_with_object": { "field": "bye" } }`},
				{Data: `{ "sub_with_object": { "field": null } }`},
			},
		},

		{
			Name: "subscription_resolver_can_error",
			Schema: makeSubscriptionSchema(t, graphql.ObjectConfig{
				Name: "Subscription",
				Fields: graphql.Fields{
					"should_error": &graphql.Field{
						Type: graphql.String,
						Subscribe: func(p graphql.ResolveParams) (interface{}, error) {
							return nil, errors.New("got a subscribe error")
						},
					},
				},
			}),
			Query: `
				subscription {
					should_error
				}
			`,
			ExpectedResults: []testutil.TestResponse{
				{
					Errors: []string{"got a subscribe error"},
				},
			},
		},
		{
			Name: "schema_without_subscribe_errors",
			Schema: makeSubscriptionSchema(t, graphql.ObjectConfig{
				Name: "Subscription",
				Fields: graphql.Fields{
					"should_error": &graphql.Field{
						Type: graphql.String,
					},
				},
			}),
			Query: `
				subscription {
					should_error
				}
			`,
			ExpectedResults: []testutil.TestResponse{
				{
					Errors: []string{"the subscription function \"should_error\" is not defined"},
				},
			},
		},
	})
}

func makeSubscribeToStringFunction(elements []string) func(p graphql.ResolveParams) (interface{}, error) {
	return func(p graphql.ResolveParams) (interface{}, error) {
		c := make(chan interface{})
		go func() {
			for _, r := range elements {
				select {
				case <-p.Context.Done():
					close(c)
					return
				case c <- r:
				}
			}
			close(c)
		}()
		return c, nil
	}
}

func makeSubscribeToMapFunction(elements []map[string]interface{}) func(p graphql.ResolveParams) (interface{}, error) {
	return func(p graphql.ResolveParams) (interface{}, error) {
		c := make(chan interface{})
		go func() {
			for _, r := range elements {
				select {
				case <-p.Context.Done():
					close(c)
					return
				case c <- r:
				}
			}
			close(c)
		}()
		return c, nil
	}
}

func makeSubscriptionSchema(t *testing.T, c graphql.ObjectConfig) graphql.Schema {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:        dummyQuery,
		Subscription: graphql.NewObject(c),
	})
	if err != nil {
		t.Errorf("failed to create schema: %v", err)
	}
	return schema
}

var dummyQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{

		"hello": &graphql.Field{Type: graphql.String},
	},
})
