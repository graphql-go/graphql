package graphql

import (
	"container/list"
	"sync"
	"sync/atomic"

	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
)

// PlanCache is a bounded, schema-aware LRU of parsed + validated +
// planned query state. Drop-in: a server's hot loop becomes
//
//   pr := cache.Get(schema, queryString, opName)
//   if len(pr.Errors) > 0 { return errorResponse(pr.Errors) }
//   args := mergeArgs(requestArgs, pr.SynthArgs)
//   result := graphql.ExecutePlan(pr.Plan, graphql.ExecuteParams{
//       Schema: *schema, Args: args, Context: ctx,
//   })
//
// Each entry holds the *Plan plus any validation errors that arose
// at parse/validate/plan time. Entries are bound to the *Schema
// pointer they were planned against; on schema rebuild the *Schema
// pointer changes and stale entries fall out at lookup.
//
// The cache is safe for concurrent use.
type PlanCache struct {
	opts PlanCacheOptions

	mu      sync.Mutex
	entries map[string]*list.Element
	order   *list.List

	hits   atomic.Uint64
	misses atomic.Uint64
}

// PlanCacheOptions tunes the cache. Zero values get sensible
// defaults. MaxEntries bounds the entry count (LRU eviction past
// it); MaxQueryBytes bypasses the cache for over-cap queries (large
// queries plus unbounded entry retention is the easy way to OOM a
// gateway). Normalize triggers literal→variable rewriting before
// hashing, so two queries that differ only in literal values
// collapse to one cache entry.
type PlanCacheOptions struct {
	MaxEntries    int
	MaxQueryBytes int
	Normalize     bool
}

const (
	defaultPlanCacheMaxEntries    = 1024
	defaultPlanCacheMaxQueryBytes = 64 * 1024
)

// PlanResult is what PlanCache.Get returns. On a successful build,
// Plan is non-nil and Errors is empty; on failure, Plan is nil and
// Errors holds the parse/validate/plan errors in the order they
// were detected. SynthArgs is populated when normalization extracted
// literals into synthetic variables — callers must merge it into
// the request's Args before ExecutePlan.
type PlanResult struct {
	Plan      *Plan
	SynthArgs map[string]interface{}
	Errors    []gqlerrors.FormattedError
}

// internal cache entry layout. Stored as the value of a list.Element;
// the list owns LRU order, the map owns key-to-element lookup.
type planCacheEntry struct {
	schema *Schema
	result PlanResult
}

type planCacheItem struct {
	key string
	e   *planCacheEntry
}

// NewPlanCache returns a PlanCache with the given options. Pass
// PlanCacheOptions{} to take all defaults.
func NewPlanCache(opts PlanCacheOptions) *PlanCache {
	if opts.MaxEntries <= 0 {
		opts.MaxEntries = defaultPlanCacheMaxEntries
	}
	if opts.MaxQueryBytes <= 0 {
		opts.MaxQueryBytes = defaultPlanCacheMaxQueryBytes
	}
	return &PlanCache{
		opts:    opts,
		entries: make(map[string]*list.Element, opts.MaxEntries),
		order:   list.New(),
	}
}

// Get parses+validates+plans `query` against `schema`, or returns a
// cached entry if the query has been seen before with the same
// schema pointer.
//
// Behavior on miss:
//  1. Parse the query string. Parse errors → PlanResult{Errors}, no
//     cache (parse failures are usually client bugs that don't
//     repeat).
//  2. (When opts.Normalize is true) rewrite leaf literal arguments
//     into synth variables; the printed normalized doc becomes the
//     cache key, so two queries that differ only in literal values
//     collapse to one entry.
//  3. Validate the (normalized) document against the schema.
//     Validation errors → PlanResult{Errors}, cached so
//     repeat-bad-queries don't hammer the validator.
//  4. Plan the operation. Plan errors → cached as PlanResult{Errors}.
//
// On a hit with normalize=true, the synthetic arguments extracted at
// THIS call's parse+normalize step are returned in PlanResult.SynthArgs;
// the cached *Plan itself is shared across all calls.
func (c *PlanCache) Get(schema *Schema, query, operationName string) PlanResult {
	if c == nil {
		// Nil-receiver convenience: caller can pass nil and still
		// get a working (uncached) path. Useful for tests +
		// "off by default" deployments.
		return planAndValidate(schema, query, operationName)
	}
	if !c.shouldCache(len(query)) {
		// Over-cap queries bypass the cache entirely.
		return planAndValidate(schema, query, operationName)
	}

	if !c.opts.Normalize {
		// Plain cache: raw query string is the key.
		key := operationName + "\x00" + query
		if pr, ok := c.lookup(schema, key); ok {
			return pr
		}
		pr := planAndValidate(schema, query, operationName)
		c.store(schema, key, pr)
		return pr
	}

	// Normalized cache: parse first (the AST is the input to
	// normalization), then key the cache on the printed normalized
	// document. Each call carries its own synth-args even when the
	// underlying *Plan is shared.
	src := source.NewSource(&source.Source{Body: []byte(query), Name: "GraphQL request"})
	doc, parseErr := parser.Parse(parser.ParseParams{Source: src})
	if parseErr != nil {
		return PlanResult{Errors: gqlerrors.FormatErrors(parseErr)}
	}
	normDoc, synthArgs, normKey, normErr := normalizeDocument(schema, doc, operationName)
	if normErr != nil {
		return PlanResult{Errors: gqlerrors.FormatErrors(normErr)}
	}
	cacheKey := operationName + "\x00" + normKey
	if pr, ok := c.lookup(schema, cacheKey); ok {
		// Stash this call's synthArgs onto the returned result.
		// The cached PlanResult deliberately stores no synthArgs
		// — those are per-call, derived freshly from the literals
		// in the incoming query.
		pr.SynthArgs = synthArgs
		return pr
	}
	if vr := ValidateDocument(schema, normDoc, nil); !vr.IsValid {
		pr := PlanResult{Errors: vr.Errors}
		c.store(schema, cacheKey, pr)
		return pr
	}
	plan, err := PlanQuery(schema, normDoc, operationName)
	if err != nil {
		pr := PlanResult{Errors: gqlerrors.FormatErrors(err)}
		c.store(schema, cacheKey, pr)
		return pr
	}
	c.store(schema, cacheKey, PlanResult{Plan: plan})
	return PlanResult{Plan: plan, SynthArgs: synthArgs}
}

// HitsMisses returns the cumulative hit and miss counts. Useful for
// surfacing as Prometheus counters.
func (c *PlanCache) HitsMisses() (hits, misses uint64) {
	if c == nil {
		return 0, 0
	}
	return c.hits.Load(), c.misses.Load()
}

// Reset drops every entry. Operators rebuilding the schema can call
// this to reclaim memory immediately rather than waiting for the
// schema-pointer mismatch to evict entries one at a time.
func (c *PlanCache) Reset() {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]*list.Element, c.opts.MaxEntries)
	c.order = list.New()
}

func (c *PlanCache) shouldCache(querySize int) bool {
	return c.opts.MaxQueryBytes <= 0 || querySize <= c.opts.MaxQueryBytes
}

func (c *PlanCache) lookup(schema *Schema, key string) (PlanResult, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	el, ok := c.entries[key]
	if !ok {
		c.misses.Add(1)
		return PlanResult{}, false
	}
	item := el.Value.(*planCacheItem)
	if item.e.schema != schema {
		c.order.Remove(el)
		delete(c.entries, key)
		c.misses.Add(1)
		return PlanResult{}, false
	}
	c.order.MoveToFront(el)
	c.hits.Add(1)
	return item.e.result, true
}

func (c *PlanCache) store(schema *Schema, key string, pr PlanResult) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if el, ok := c.entries[key]; ok {
		item := el.Value.(*planCacheItem)
		item.e.schema = schema
		item.e.result = pr
		c.order.MoveToFront(el)
		return
	}
	item := &planCacheItem{key: key, e: &planCacheEntry{schema: schema, result: pr}}
	el := c.order.PushFront(item)
	c.entries[key] = el
	for c.order.Len() > c.opts.MaxEntries {
		oldest := c.order.Back()
		if oldest == nil {
			break
		}
		oi := oldest.Value.(*planCacheItem)
		c.order.Remove(oldest)
		delete(c.entries, oi.key)
	}
}

// planAndValidate is the cache-miss path: parse, validate, plan,
// returning a PlanResult ready to store. Exposed as a free function
// (no receiver) so the nil-receiver Get can fall through to it
// without re-implementing the pipeline.
func planAndValidate(schema *Schema, query, operationName string) PlanResult {
	src := source.NewSource(&source.Source{Body: []byte(query), Name: "GraphQL request"})
	doc, parseErr := parser.Parse(parser.ParseParams{Source: src})
	if parseErr != nil {
		return PlanResult{Errors: gqlerrors.FormatErrors(parseErr)}
	}
	if vr := ValidateDocument(schema, doc, nil); !vr.IsValid {
		return PlanResult{Errors: vr.Errors}
	}
	plan, err := PlanQuery(schema, doc, operationName)
	if err != nil {
		return PlanResult{Errors: gqlerrors.FormatErrors(err)}
	}
	return PlanResult{Plan: plan}
}
