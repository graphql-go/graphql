package graphql_test

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/graphql-go/graphql"
)

// TestSubscribeProducerExitsWhenConsumerStops reproduces #758: if the
// consumer stops reading resultChannel and cancels the context, the producer
// goroutine must exit (the result channel must close) instead of blocking
// forever on an unbuffered send.
func TestSubscribeProducerExitsWhenConsumerStops(t *testing.T) {
	// secondResolved is closed once the second source value has been mapped,
	// i.e. the producer is about to send it on resultChannel.
	secondResolved := make(chan struct{})
	subType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Subscription",
		Fields: graphql.Fields{
			"tick": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if p.Source == 1 {
						close(secondResolved)
					}
					return p.Source, nil
				},
				Subscribe: func(p graphql.ResolveParams) (interface{}, error) {
					ch := make(chan interface{})
					go func() {
						defer close(ch)
						for i := 0; i < 100; i++ {
							select {
							case <-p.Context.Done():
								return
							case ch <- i:
							}
						}
					}()
					return ch, nil
				},
			},
		},
	})
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Query",
		Fields: graphql.Fields{"ok": &graphql.Field{Type: graphql.Boolean}},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType, Subscription: subType})
	if err != nil {
		t.Fatal(err)
	}

	before := runtime.NumGoroutine()
	ctx, cancel := context.WithCancel(context.Background())
	resultChannel := graphql.Subscribe(graphql.Params{
		Schema:        schema,
		Context:       ctx,
		RequestString: "subscription { tick }",
	})

	// Read one result, then stop reading. Wait until the producer has mapped
	// the second value, so it is about to block on the result send, and
	// cancel. Before the fix the producer stayed blocked on that send forever.
	select {
	case <-resultChannel:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for first subscription result")
	}
	select {
	case <-secondResolved:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for the second value to be resolved")
	}
	cancel()

	// Do not read resultChannel here: a read would unblock a leaked send.
	// Wait for the producer and source goroutines to exit instead.
	deadline := time.Now().Add(2 * time.Second)
	for runtime.NumGoroutine() > before {
		if time.Now().After(deadline) {
			t.Fatalf("subscription goroutines did not exit after cancel (#758): %d running, %d before", runtime.NumGoroutine(), before)
		}
		time.Sleep(10 * time.Millisecond)
	}

	if _, ok := <-resultChannel; ok {
		t.Fatal("result channel delivered a value after cancel instead of being closed (#758)")
	}
}
