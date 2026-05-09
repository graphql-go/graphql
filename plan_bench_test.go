package graphql_test

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/benchutil"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
)

// BenchmarkPlannedExecute_* compare a cached *Plan re-executed N times
// vs graphql.Do (which parses, validates, plans, and executes every
// call). The cached path skips parse + validate + plan; what's left
// is the work that's inherently per-request.

func BenchmarkPlannedExecute_WideQuery_100_10(b *testing.B) {
	schema := benchutil.WideSchemaWithXFieldsAndYItems(100, 10)
	query := benchutil.WideSchemaQuery(100)

	src := source.NewSource(&source.Source{Body: []byte(query), Name: "bench"})
	doc, err := parser.Parse(parser.ParseParams{Source: src})
	if err != nil {
		b.Fatalf("parse: %v", err)
	}
	plan, err := graphql.PlanQuery(&schema, doc, "")
	if err != nil {
		b.Fatalf("plan: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := graphql.ExecutePlan(plan, graphql.ExecuteParams{
			Schema: schema,
			AST:    doc,
		})
		if len(result.Errors) > 0 {
			b.Fatalf("errors: %v", result.Errors)
		}
	}
}

func BenchmarkUncachedExecute_WideQuery_100_10(b *testing.B) {
	schema := benchutil.WideSchemaWithXFieldsAndYItems(100, 10)
	query := benchutil.WideSchemaQuery(100)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: query,
		})
		if len(result.Errors) > 0 {
			b.Fatalf("errors: %v", result.Errors)
		}
	}
}

func BenchmarkPlannedExecute_ListQuery_1K(b *testing.B) {
	schema := benchutil.ListSchemaWithXItems(1000)
	query := `query { colors { hex r g b } }`

	src := source.NewSource(&source.Source{Body: []byte(query), Name: "bench"})
	doc, err := parser.Parse(parser.ParseParams{Source: src})
	if err != nil {
		b.Fatalf("parse: %v", err)
	}
	plan, err := graphql.PlanQuery(&schema, doc, "")
	if err != nil {
		b.Fatalf("plan: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := graphql.ExecutePlan(plan, graphql.ExecuteParams{
			Schema: schema,
			AST:    doc,
		})
		if len(result.Errors) > 0 {
			b.Fatalf("errors: %v", result.Errors)
		}
	}
}

func BenchmarkUncachedExecute_ListQuery_1K(b *testing.B) {
	schema := benchutil.ListSchemaWithXItems(1000)
	query := `query { colors { hex r g b } }`

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: query,
		})
		if len(result.Errors) > 0 {
			b.Fatalf("errors: %v", result.Errors)
		}
	}
}
