package graphql_test

import (
	"fmt"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/testutil"
)

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

func TestSchemaSubscribe(t *testing.T) {

	testutil.RunSubscribes(t, []*testutil.TestSubscription{
		{
			Name: "subscribe without resolver",
			Schema: makeSubscriptionSchema(t, graphql.ObjectConfig{
				Name: "Subscription",
				Fields: graphql.Fields{
					"sub_without_resolver": &graphql.Field{
						Type:      graphql.String,
						Subscribe: makeSubscribeToStringFunction([]string{"a", "b", "c"}),
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							return p.Source, nil
						},
					},
				},
			}),
			Query: `
				subscription onHelloSaid {
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
			Name: "subscribe with resolver changes output",
			Schema: makeSubscriptionSchema(t, graphql.ObjectConfig{
				Name: "Subscription",
				Fields: graphql.Fields{
					"sub_with_resolver": &graphql.Field{
						Type:      graphql.String,
						Subscribe: makeSubscribeToStringFunction([]string{"a", "b", "c"}),
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							return fmt.Sprintf("result=%v", p.Source), nil
						},
					},
				},
			}),
			Query: `
				subscription onHelloSaid {
					sub_with_resolver
				}
			`,
			ExpectedResults: []testutil.TestResponse{
				{Data: `{ "sub_with_resolver": "result=a" }`},
				{Data: `{ "sub_with_resolver": "result=b" }`},
				{Data: `{ "sub_with_resolver": "result=c" }`},
			},
		},
		// {
		// 	Name: "subscribe to a nested object",
		// 	Schema: makeSubscriptionSchema(t, graphql.ObjectConfig{
		// 		Name: "Subscription",
		// 		Fields: graphql.Fields{
		// 			"sub_with_object": &graphql.Field{
		// 				Type: graphql.String,
		// 				Subscribe: makeSubscribeToMapFunction([]map[string]interface{}{
		// 					{
		// 						"field": "hello",
		// 						"obj": map[string]interface{}{
		// 							"field": "hello",
		// 						},
		// 					},
		// 					{
		// 						"field": "bye",
		// 						"obj": map[string]interface{}{
		// 							"field": "bye",
		// 						},
		// 					},
		// 				}),
		// 			},
		// 		},
		// 	}),
		// 	Query: `
		// 		subscription onHelloSaid {
		// 			sub_with_object {
		// 				field
		// 				obj {
		// 					field
		// 				}
		// 			}
		// 		}
		// 	`,
		// 	ExpectedResults: []testutil.TestResponse{
		// 		{Data: `{ "sub_with_object": { "field": "hello", "obj": { "field": "hello" } } }`},
		// 	},
		// },

		// {
		// 	Name:   "parse_errors",
		// 	Schema: schema,
		// 	Query:  `invalid graphQL query`,
		// 	ExpectedResults: []testutil.TestResponse{
		// 		{
		// 			Errors: []gqlerrors.FormattedError{{Message: ""}},
		// 		},
		// 	},
		// },
		// {
		// 	Name:   "subscribe_to_query_succeeds",
		// 	Schema: schema,
		// 	Query: `
		// 		query Hello {
		// 			hello
		// 		}
		// 	`,
		// 	ExpectedResults: []testutil.TestResponse{
		// 		{
		// 			Data: json.RawMessage(`
		// 				{
		// 					"hello": "Hello world!"
		// 				}
		// 			`),
		// 		},
		// 	},
		// },
		// {
		// 	Name:   "subscription_resolver_can_error",
		// 	Schema: schema,
		// 	Query: `
		// 		subscription onHelloSaid {
		// 			helloSaid {
		// 				msg
		// 			}
		// 		}
		// 	`,
		// 	ExpectedResults: []testutil.TestResponse{
		// 		{
		// 			Data: json.RawMessage(`
		// 				null
		// 			`),
		// 			Errors: []gqlerrors.FormattedError{{Message: ""}}},
		// 	},
		// },
		// {
		// 	Name:   "subscription_resolver_can_error_optional_msg",
		// 	Schema: schema,
		// 	Query: `
		// 		subscription onHelloSaid {
		// 			helloSaidNullable {
		// 				msg
		// 			}
		// 		}
		// 	`,
		// 	ExpectedResults: []testutil.TestResponse{
		// 		{
		// 			Data: json.RawMessage(`
		// 				{
		// 					"helloSaidNullable": {
		// 						"msg": null
		// 					}
		// 				}
		// 			`),
		// 			Errors: []gqlerrors.FormattedError{{Message: ""}}},
		// 	},
		// },
		// {
		// 	Name:   "subscription_resolver_can_error_optional_event",
		// 	Schema: schema,
		// 	Query: `
		// 		subscription onHelloSaid {
		// 			helloSaidNullable {
		// 				msg
		// 			}
		// 		}
		// 	`,
		// 	ExpectedResults: []testutil.TestResponse{
		// 		{
		// 			Data: json.RawMessage(`
		// 				{
		// 					"helloSaidNullable": null
		// 				}
		// 			`),
		// 			Errors: []gqlerrors.FormattedError{{Message: ""}}},
		// 	},
		// },
		// {
		// 	Name:   "schema_without_resolver_errors",
		// 	Schema: schema,
		// 	Query: `
		// 		subscription onHelloSaid {
		// 			helloSaid {
		// 				msg
		// 			}
		// 		}
		// 	`,
		// 	ExpectedErr: errors.New("schema created without resolver, can not subscribe"),
		// },
	})
}

// func TestRootOperations_invalidSubscriptionSchema(t *testing.T) {
// 	type args struct {
// 		Schema string
// 	}
// 	type want struct {
// 		Error string
// 	}
// 	testTable := map[string]struct {
// 		Args args
// 		Want want
// 	}{
// 		"Subscription as incorrect type": {
// 			Args: args{
// 				Schema: `
// 					schema {
// 						query: Query
// 						subscription: String
// 					}
// 					type Query {
// 						thing: String
// 					}
// 				`,
// 			},
// 			Want: want{Error: `root operation "subscription" must be an OBJECT`},
// 		},
// 		"Subscription declared by schema, but type not present": {
// 			Args: args{
// 				Schema: `
// 					schema {
// 						query: Query
// 						subscription: Subscription
// 					}
// 					type Query {
// 						hello: String!
// 					}
// 				`,
// 			},
// 			Want: want{Error: `graphql: type "Subscription" not found`},
// 		},
// 	}

// 	for name, tt := range testTable {
// 		tt := tt
// 		t.Run(name, func(t *testing.T) {
// 			t.Log(tt.Args.Schema) // TODO do something
// 		})
// 	}
// }

var dummyQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{

		"hello": &graphql.Field{Type: graphql.String},
	},
})
