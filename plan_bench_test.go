package graphql_test

import (
	"fmt"
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

// BenchmarkPlannedExecute_WideQuery_100_10_Varied demonstrates that
// a single cached *Plan handles arbitrary literal variations — the
// plan binds the field arg to a `$v` variable and the request's Args
// changes per call. No re-parse, no re-validate, no re-plan; just
// per-request arg coercion (the only inherently-dynamic work) plus
// the resolver loop. This is the canonical parametric-query path:
// every real client should look like this.
func BenchmarkPlannedExecute_WideQuery_100_10_Varied(b *testing.B) {
	schema := benchutil.WideArgedSchemaWithXFieldsAndYItems(100, 10)
	query := benchutil.WideArgedSchemaQueryWithVariable(100)

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
			Args:   map[string]interface{}{"v": fmt.Sprintf("v-%d", i)},
		})
		if len(result.Errors) > 0 {
			b.Fatalf("errors: %v", result.Errors)
		}
	}
}

// BenchmarkUncachedExecute_WideQuery_100_10_Varied is the comparison:
// same workload, but the caller bakes the literal into the query
// string itself and lets graphql.Do parse + validate + plan + execute
// every single call. This is what naive clients do when they don't
// use GraphQL variables.
func BenchmarkUncachedExecute_WideQuery_100_10_Varied(b *testing.B) {
	schema := benchutil.WideArgedSchemaWithXFieldsAndYItems(100, 10)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		query := benchutil.WideArgedSchemaQueryWithLiteral(100, fmt.Sprintf("v-%d", i))
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: query,
		})
		if len(result.Errors) > 0 {
			b.Fatalf("errors: %v", result.Errors)
		}
	}
}

// BenchmarkPlanCache_HotLoop_NativeVars: the canonical fast path —
// the client sends the query once with `$v` variables; the cache
// stores one plan; every iteration is `cache.Get + ExecutePlan` with
// varying Args. This is what real clients should look like.
func BenchmarkPlanCache_HotLoop_NativeVars(b *testing.B) {
	schema := benchutil.WideArgedSchemaWithXFieldsAndYItems(100, 10)
	query := benchutil.WideArgedSchemaQueryWithVariable(100)
	cache := graphql.NewPlanCache(graphql.PlanCacheOptions{})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pr := cache.Get(&schema, query, "")
		if len(pr.Errors) > 0 {
			b.Fatalf("get: %v", pr.Errors)
		}
		result := graphql.ExecutePlan(pr.Plan, graphql.ExecuteParams{
			Schema: schema,
			Args:   map[string]interface{}{"v": fmt.Sprintf("v-%d", i)},
		})
		if len(result.Errors) > 0 {
			b.Fatalf("execute: %v", result.Errors)
		}
	}
}

// BenchmarkPlanCache_HotLoop_Normalized: the salvage path — the
// client sends literal-baked queries that vary per call. With
// Normalize=true, every iteration parses + normalizes (cheap) and
// hits the same cached plan. Demonstrates the full effect of cache
// + normalization on the worst client behavior.
func BenchmarkPlanCache_HotLoop_Normalized(b *testing.B) {
	schema := benchutil.WideArgedSchemaWithXFieldsAndYItems(100, 10)
	cache := graphql.NewPlanCache(graphql.PlanCacheOptions{Normalize: true})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		query := benchutil.WideArgedSchemaQueryWithLiteral(100, fmt.Sprintf("v-%d", i))
		pr := cache.Get(&schema, query, "")
		if len(pr.Errors) > 0 {
			b.Fatalf("get: %v", pr.Errors)
		}
		result := graphql.ExecutePlan(pr.Plan, graphql.ExecuteParams{
			Schema: schema,
			Args:   pr.SynthArgs,
		})
		if len(result.Errors) > 0 {
			b.Fatalf("execute: %v", result.Errors)
		}
	}
}

// BenchmarkPlanCache_HotLoop_NoNorm: the worst case — literal-baked
// queries vary per call, but normalization is OFF. Every iteration
// is a fresh parse + validate + plan. The cache only hits when an
// LRU-replayed literal happens to match (vanishingly rare for real
// workloads). This measures what users get if they neither use
// variables nor turn on normalization.
func BenchmarkPlanCache_HotLoop_NoNorm(b *testing.B) {
	schema := benchutil.WideArgedSchemaWithXFieldsAndYItems(100, 10)
	cache := graphql.NewPlanCache(graphql.PlanCacheOptions{})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		query := benchutil.WideArgedSchemaQueryWithLiteral(100, fmt.Sprintf("v-%d", i))
		pr := cache.Get(&schema, query, "")
		if len(pr.Errors) > 0 {
			b.Fatalf("get: %v", pr.Errors)
		}
		result := graphql.ExecutePlan(pr.Plan, graphql.ExecuteParams{
			Schema: schema,
			Args:   pr.SynthArgs,
		})
		if len(result.Errors) > 0 {
			b.Fatalf("execute: %v", result.Errors)
		}
	}
}

// BenchmarkPlannedExecute_WideQuery_100_10_StaticArg is the static
// counterpart: same plan, same Args every call. Lets us see how
// close the Varied case (per-call arg coercion) gets to the
// fully-static ceiling.
func BenchmarkPlannedExecute_WideQuery_100_10_StaticArg(b *testing.B) {
	schema := benchutil.WideArgedSchemaWithXFieldsAndYItems(100, 10)
	query := benchutil.WideArgedSchemaQueryWithVariable(100)

	src := source.NewSource(&source.Source{Body: []byte(query), Name: "bench"})
	doc, err := parser.Parse(parser.ParseParams{Source: src})
	if err != nil {
		b.Fatalf("parse: %v", err)
	}
	plan, err := graphql.PlanQuery(&schema, doc, "")
	if err != nil {
		b.Fatalf("plan: %v", err)
	}
	args := map[string]interface{}{"v": "static"}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := graphql.ExecutePlan(plan, graphql.ExecuteParams{
			Schema: schema,
			AST:    doc,
			Args:   args,
		})
		if len(result.Errors) > 0 {
			b.Fatalf("errors: %v", result.Errors)
		}
	}
}
