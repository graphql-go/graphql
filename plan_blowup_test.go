package graphql

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
)

func parseQuery(t testing.TB, query string) *ast.Document {
	t.Helper()
	doc, err := parser.Parse(parser.ParseParams{Source: query})
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	return doc
}

// buildAbstractFanSchema builds a schema with a single interface `Node`
// implemented by `numTypes` object types. The interface (and therefore
// every implementer) exposes a `next: Node` field, so the type graph is
// cyclic through the abstract type: planning a selection on `next` must
// consider every possible type, each of which again has a `next`.
//
// A naive planner that eagerly expands every possible type for every
// abstract field, at every nesting level, does O(numTypes ^ depth) work
// — even though a single request only ever resolves one concrete type
// per level. This is the shape that made a real, deeply-nested
// polymorphic query take minutes to plan.
func buildAbstractFanSchema(t testing.TB, numTypes int) Schema {
	t.Helper()

	nodeIface := NewInterface(InterfaceConfig{
		Name: "Node",
		Fields: Fields{
			"id": &Field{Type: String},
		},
	})
	// Self-reference: `next: Node` makes the type graph cyclic through
	// the abstract type. Added after construction to reference the
	// interface itself.
	nodeIface.AddFieldConfig("next", &Field{Type: nodeIface})

	objs := make([]*Object, numTypes)
	for i := 0; i < numTypes; i++ {
		objs[i] = NewObject(ObjectConfig{
			Name:       fmt.Sprintf("T%d", i),
			Interfaces: []*Interface{nodeIface},
			Fields: Fields{
				"id":   &Field{Type: String},
				"next": &Field{Type: nodeIface},
			},
		})
	}
	// Resolve every abstract value to T0 so the correctness test has a
	// deterministic concrete type.
	nodeIface.ResolveType = func(p ResolveTypeParams) *Object { return objs[0] }

	queryType := NewObject(ObjectConfig{
		Name: "Query",
		Fields: Fields{
			"start": &Field{
				Type: nodeIface,
				Resolve: func(p ResolveParams) (interface{}, error) {
					return map[string]interface{}{"id": "root"}, nil
				},
			},
		},
	})

	schema, err := NewSchema(SchemaConfig{
		Query: queryType,
		Types: objectsToTypes(objs),
	})
	if err != nil {
		t.Fatalf("build schema: %v", err)
	}
	return schema
}

func objectsToTypes(objs []*Object) []Type {
	types := make([]Type, len(objs))
	for i, o := range objs {
		types[i] = o
	}
	return types
}

// nestedNextQuery builds `{ start { id next { id next { ... } } } }`
// nested `depth` levels deep.
func nestedNextQuery(depth int) string {
	var b strings.Builder
	b.WriteString("{ start { id ")
	for i := 0; i < depth; i++ {
		b.WriteString("next { id ")
	}
	for i := 0; i < depth; i++ {
		b.WriteString("} ")
	}
	b.WriteString("} }")
	return b.String()
}

// TestPlanQueryAbstractFanDoesNotExplode is the regression guard: a
// deeply-nested abstract selection must plan in time linear in the
// query size, not exponential in the number of possible types.
//
// With numTypes=16 and depth=12, an eager-expansion planner would do
// ~16^12 (≈ 2.8e14) units of work and never return; the lazy planner
// returns essentially instantly.
func TestPlanQueryAbstractFanDoesNotExplode(t *testing.T) {
	schema := buildAbstractFanSchema(t, 16)
	doc := parseQuery(t, nestedNextQuery(12))

	done := make(chan error, 1)
	go func() {
		_, err := PlanQuery(&schema, doc, "")
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("PlanQuery returned error: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("PlanQuery did not complete within 10s: abstract-type " +
			"planning is expanding possible types exponentially with nesting depth")
	}
}

// TestPlanQueryAbstractFanCorrectness ensures the lazy abstract planning
// still resolves concrete types and returns correct data.
func TestPlanQueryAbstractFanCorrectness(t *testing.T) {
	schema := buildAbstractFanSchema(t, 4)
	result := Do(Params{
		Schema:        schema,
		RequestString: "{ start { id next { id } } }",
	})
	if len(result.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	data, _ := result.Data.(map[string]interface{})
	start, _ := data["start"].(map[string]interface{})
	if start == nil {
		t.Fatalf("expected start object, got: %#v", result.Data)
	}
	if start["id"] != "root" {
		t.Fatalf("expected start.id=root, got %#v", start["id"])
	}
	// next resolves to nil (no resolver returns a child) → null.
	if v, ok := start["next"]; ok && v != nil {
		t.Fatalf("expected start.next=nil, got %#v", v)
	}
}

func BenchmarkPlanQueryAbstractFan(b *testing.B) {
	for _, dims := range []struct{ numTypes, depth int }{
		{8, 6}, {16, 8}, {32, 10},
	} {
		schema := buildAbstractFanSchema(b, dims.numTypes)
		doc := parseQuery(b, nestedNextQuery(dims.depth))
		b.Run(fmt.Sprintf("types=%d/depth=%d", dims.numTypes, dims.depth), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if _, err := PlanQuery(&schema, doc, ""); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
