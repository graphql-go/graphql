package dataloaderexample_test

import (
	"testing"

	"github.com/graphql-go/graphql/examples/dataloader"
)

func TestQuery(t *testing.T) {
	schema := dataloaderexample.CreateSchema()
	r := dataloaderexample.RunQuery(`{
        p1_0: post(id: "1") { id author { name }}
        p1_1: post(id: "1") { id author { name }}
        p1_2: post(id: "1") { id author { name }}
        p1_3: post(id: "1") { id author { name }}
        p1_4: post(id: "1") { id author { name }}
        p1_5: post(id: "1") { id author { name }}
        p2_1: post(id: "2") { id author { name }}
        p2_2: post(id: "2") { id author { name }}
        p2_3: post(id: "2") { id author { name }}
        p3_1: post(id: "3") { id author { name }}
        p3_2: post(id: "3") { id author { name }}
        p3_3: post(id: "3") { id author { name }}
    }`, schema)
	if len(r.Errors) != 0 {
		t.Error(r.Errors)
	}
	// The above query would produce log like this:
	// 2016/07/23 23:28:05 Load post 1
	// 2016/07/23 23:28:05 Load post 3
	// 2016/07/23 23:28:05 Load post 2
	// 2016/07/23 23:28:06 Load user 1
	// 2016/07/23 23:28:06 Load user 2
	// Notice the first level post loading is done concurrently without duplicate.
	// And the second level user loading is also done in the same fashion.
	// TODO: Make test actually verify that.
}
