package graphql_test

import (
	"reflect"
	"sync"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/benchutil"
)

// TestPlanCacheBasicHits exercises the un-normalized path: identical
// raw query strings hit the same cache entry; different ones miss.
func TestPlanCacheBasicHits(t *testing.T) {
	schema := benchutil.WideArgedSchemaWithXFieldsAndYItems(5, 1)
	cache := graphql.NewPlanCache(graphql.PlanCacheOptions{})

	q := `{ wide { a(value: "x") } }`
	r1 := cache.Get(&schema, q, "")
	if len(r1.Errors) > 0 || r1.Plan == nil {
		t.Fatalf("first call: %v", r1.Errors)
	}
	hits, misses := cache.HitsMisses()
	if hits != 0 || misses != 1 {
		t.Fatalf("after first miss: hits=%d misses=%d", hits, misses)
	}

	r2 := cache.Get(&schema, q, "")
	if len(r2.Errors) > 0 {
		t.Fatalf("second call: %v", r2.Errors)
	}
	if r1.Plan != r2.Plan {
		t.Fatalf("expected same *Plan pointer on cache hit")
	}
	hits, misses = cache.HitsMisses()
	if hits != 1 || misses != 1 {
		t.Fatalf("after hit: hits=%d misses=%d", hits, misses)
	}

	// Different literal → different raw key → miss (un-normalized).
	r3 := cache.Get(&schema, `{ wide { a(value: "y") } }`, "")
	if r1.Plan == r3.Plan {
		t.Fatalf("un-normalized: different literals must miss the cache")
	}
}

// TestPlanCacheNormalizeCollapsesLiterals exercises the normalize path:
// queries that differ only in literal values must hit the same plan.
func TestPlanCacheNormalizeCollapsesLiterals(t *testing.T) {
	schema := benchutil.WideArgedSchemaWithXFieldsAndYItems(5, 1)
	cache := graphql.NewPlanCache(graphql.PlanCacheOptions{Normalize: true})

	r1 := cache.Get(&schema, `{ wide { a(value: "asdf") } }`, "")
	if len(r1.Errors) > 0 || r1.Plan == nil {
		t.Fatalf("first call: %v", r1.Errors)
	}
	if got := r1.SynthArgs["__pcv0"]; got != "asdf" {
		t.Fatalf("synthArgs[__pcv0] = %v, want asdf", got)
	}

	r2 := cache.Get(&schema, `{ wide { a(value: "qwer") } }`, "")
	if len(r2.Errors) > 0 {
		t.Fatalf("second call: %v", r2.Errors)
	}
	if r1.Plan != r2.Plan {
		t.Fatalf("normalized: different literals must reuse the same *Plan pointer")
	}
	if got := r2.SynthArgs["__pcv0"]; got != "qwer" {
		t.Fatalf("synthArgs[__pcv0] = %v, want qwer", got)
	}
	hits, _ := cache.HitsMisses()
	if hits != 1 {
		t.Fatalf("expected 1 hit, got %d", hits)
	}
}

// TestPlanCacheExecutesNormalizedPlan verifies the full end-to-end
// path: cache hit + ExecutePlan + merged Args produces the right
// resolver output.
func TestPlanCacheExecutesNormalizedPlan(t *testing.T) {
	schema := benchutil.WideArgedSchemaWithXFieldsAndYItems(2, 1)
	cache := graphql.NewPlanCache(graphql.PlanCacheOptions{Normalize: true})

	for _, lit := range []string{"asdf", "qwer", "zxcv"} {
		query := `{ wide { a(value: "` + lit + `") } }`
		pr := cache.Get(&schema, query, "")
		if len(pr.Errors) > 0 {
			t.Fatalf("Get %q: %v", query, pr.Errors)
		}
		result := graphql.ExecutePlan(pr.Plan, graphql.ExecuteParams{
			Schema: schema,
			Args:   pr.SynthArgs,
		})
		if len(result.Errors) > 0 {
			t.Fatalf("ExecutePlan %q: %v", query, result.Errors)
		}
		// Each item's "a" field echoes the value arg.
		got := result.Data
		want := map[string]interface{}{"wide": []interface{}{map[string]interface{}{"a": lit}}}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("query %q: got %v, want %v", query, got, want)
		}
	}
}

// TestPlanCacheNormalizeUnnormalizableNoCollision exercises the
// fallback for documents that normalizeDocument can't analyse
// (ambiguous multi-op without operationName, named op not found in
// the document). Two such queries must not collide on the cache key
// — without the raw-doc fingerprint fallback, both would land in the
// same "" slot and the second Get would incorrectly hit the first's
// cached errors.
func TestPlanCacheNormalizeUnnormalizableNoCollision(t *testing.T) {
	schema := benchutil.WideArgedSchemaWithXFieldsAndYItems(2, 1)
	cache := graphql.NewPlanCache(graphql.PlanCacheOptions{Normalize: true})

	q1 := `query A { wide { a(value: "x") } } query B { wide { a(value: "y") } }`
	q2 := `query A { wide { a(value: "x") } } query C { wide { a(value: "z") } }`

	if r := cache.Get(&schema, q1, ""); len(r.Errors) == 0 {
		t.Fatalf("expected errors on ambiguous multi-op q1")
	}
	if r := cache.Get(&schema, q2, ""); len(r.Errors) == 0 {
		t.Fatalf("expected errors on ambiguous multi-op q2")
	}
	hits, misses := cache.HitsMisses()
	if hits != 0 || misses != 2 {
		t.Fatalf("distinct unnormalizable queries collided: hits=%d misses=%d (want 0/2)", hits, misses)
	}

	// Re-fetch both: each should hit its own slot.
	_ = cache.Get(&schema, q1, "")
	_ = cache.Get(&schema, q2, "")
	hits, misses = cache.HitsMisses()
	if hits != 2 || misses != 2 {
		t.Fatalf("re-fetch didn't hit own slots: hits=%d misses=%d (want 2/2)", hits, misses)
	}
}

// TestPlanCacheSchemaRebuildInvalidation confirms the schema-pointer
// guard: an entry planned against one *Schema falls out of the cache
// when looked up against a freshly built one, even when the schemas
// are structurally identical.
func TestPlanCacheSchemaRebuildInvalidation(t *testing.T) {
	s1 := benchutil.WideArgedSchemaWithXFieldsAndYItems(2, 1)
	cache := graphql.NewPlanCache(graphql.PlanCacheOptions{})

	q := `{ wide { a(value: "x") } }`
	r1 := cache.Get(&s1, q, "")
	if len(r1.Errors) > 0 || r1.Plan == nil {
		t.Fatalf("first call against s1: %v", r1.Errors)
	}

	s2 := benchutil.WideArgedSchemaWithXFieldsAndYItems(2, 1)
	r2 := cache.Get(&s2, q, "")
	if len(r2.Errors) > 0 || r2.Plan == nil {
		t.Fatalf("first call against s2: %v", r2.Errors)
	}
	if r1.Plan == r2.Plan {
		t.Fatalf("rebuilt schema reused stale *Plan pointer")
	}

	// s1 is now invalidated for the same key (cache replaced its slot
	// with the s2-bound plan), so a re-lookup against s1 must miss
	// and re-plan.
	r1b := cache.Get(&s1, q, "")
	if r1b.Plan == r2.Plan {
		t.Fatalf("post-rebuild lookup against s1 returned the s2-bound plan")
	}
}

// TestPlanCacheConcurrentGets exercises the LRU/mutex behavior under
// concurrent access. Many goroutines hammer the same key; we don't
// guarantee they share a single planning pass (that would need
// singleflight) but we do guarantee no panics, no data races (use
// `go test -race`), and consistent results.
func TestPlanCacheConcurrentGets(t *testing.T) {
	schema := benchutil.WideArgedSchemaWithXFieldsAndYItems(2, 1)
	cache := graphql.NewPlanCache(graphql.PlanCacheOptions{})
	q := `{ wide { a(value: "x") } }`

	const goroutines = 16
	const perG = 32

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < perG; j++ {
				r := cache.Get(&schema, q, "")
				if len(r.Errors) > 0 || r.Plan == nil {
					t.Errorf("concurrent Get: %v", r.Errors)
					return
				}
			}
		}()
	}
	wg.Wait()

	hits, misses := cache.HitsMisses()
	if hits+misses != goroutines*perG {
		t.Fatalf("hits+misses=%d, want %d", hits+misses, goroutines*perG)
	}
	if misses == 0 {
		t.Fatalf("expected at least one miss")
	}
}

// TestPlanCacheReset clears the cache so subsequent lookups miss.
func TestPlanCacheReset(t *testing.T) {
	schema := benchutil.WideArgedSchemaWithXFieldsAndYItems(2, 1)
	cache := graphql.NewPlanCache(graphql.PlanCacheOptions{})
	q := `{ wide { a(value: "x") } }`
	_ = cache.Get(&schema, q, "")
	cache.Reset()
	_, _ = cache.HitsMisses()
	r := cache.Get(&schema, q, "")
	_, misses := cache.HitsMisses()
	if misses != 2 {
		t.Fatalf("expected 2 misses after Reset, got %d", misses)
	}
	if r.Plan == nil {
		t.Fatalf("expected new plan after Reset, got nil")
	}
}
