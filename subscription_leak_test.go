package graphql_test

import (
	"context"
	"testing"
	"time"

	"github.com/graphql-go/graphql"
)

// TestProducerExitsOnContextCancel is a regression test: cancelling the
// context must let the producer exit and close resultChannel even with no
// active receiver.
func TestProducerExitsOnContextCancel(t *testing.T) {
	// unbuffered: once this send returns the producer has consumed the value
	// and is blocked sending on resultChannel.
	sourceCh := make(chan interface{})
	subscribeType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Subscription",
		Fields: graphql.Fields{
			"tick": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source, nil
				},
				Subscribe: func(p graphql.ResolveParams) (interface{}, error) {
					return sourceCh, nil
				},
			},
		},
	})
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{"ok": &graphql.Field{Type: graphql.Boolean}},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:        queryType,
		Subscription: subscribeType,
	})
	if err != nil {
		t.Fatalf("NewSchema: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	resultChannel := graphql.Subscribe(graphql.Params{
		Schema:        schema,
		Context:       ctx,
		RequestString: "subscription { tick }",
	})

	// make the producer block on sending the result.
	sourceCh <- 42

	// subscriber leaves: cancel the context.
	cancel()

	// give the producer time to exit and close the channel.
	time.Sleep(200 * time.Millisecond)

	select {
	case _, ok := <-resultChannel:
		if !ok {
			return // channel closed → producer exited → fix works
		}
		t.Fatalf("got a result instead of close: producer did not exit cleanly")
	case <-time.After(2 * time.Second):
		t.Fatalf("resultChannel still open after cancel - producer goroutine leaked")
	}
}
