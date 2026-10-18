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
	subType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Subscription",
		Fields: graphql.Fields{
			"tick": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
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

	// Read one result so the producer has entered the send loop, then stop
	// reading and cancel — the pre-fix producer would block forever.
	select {
	case <-resultChannel:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for first subscription result")
	}
	cancel()

	select {
	case _, ok := <-resultChannel:
		if ok {
			// Drain any in-flight result, then wait for close.
			select {
			case _, ok2 := <-resultChannel:
				if ok2 {
					// keep draining briefly
					deadline := time.After(2 * time.Second)
					for ok2 {
						select {
						case _, ok2 = <-resultChannel:
						case <-deadline:
							t.Fatal("result channel did not close after cancel (#758)")
						}
					}
				}
			case <-time.After(2 * time.Second):
				t.Fatal("result channel did not close after cancel (#758)")
			}
		}
	case <-time.After(2 * time.Second):
		t.Fatal("BUG: result channel NOT closed — producer goroutine leaked (#758)")
	}

	// Best-effort: goroutine count should not keep climbing; allow a little slack.
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if runtime.NumGoroutine() <= before+5 {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
}
