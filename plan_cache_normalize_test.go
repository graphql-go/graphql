package graphql_test

import (
	"testing"

	"github.com/graphql-go/graphql"
)

// buildNormalizeTestSchema builds a small schema with enough type
// variety to exercise the fingerprint walker's branches: scalar
// args (Int, Float, Bool, String, Enum), List arg, NonNull arg,
// InputObject arg, an inline-fragmentable interface, and a
// fragment-spreadable object.
func buildNormalizeTestSchema(t *testing.T) graphql.Schema {
	t.Helper()
	color := graphql.NewEnum(graphql.EnumConfig{
		Name: "Color",
		Values: graphql.EnumValueConfigMap{
			"RED":  &graphql.EnumValueConfig{Value: "RED"},
			"BLUE": &graphql.EnumValueConfig{Value: "BLUE"},
		},
	})
	point := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "Point",
		Fields: graphql.InputObjectConfigFieldMap{
			"x": &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"y": &graphql.InputObjectFieldConfig{Type: graphql.Int},
		},
	})
	thingIface := graphql.NewInterface(graphql.InterfaceConfig{
		Name: "Thing",
		Fields: graphql.Fields{
			"name": &graphql.Field{Type: graphql.String},
		},
		ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
			return nil
		},
	})
	widget := graphql.NewObject(graphql.ObjectConfig{
		Name:       "Widget",
		Interfaces: []*graphql.Interface{thingIface},
		Fields: graphql.Fields{
			"name":   &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) { return "w", nil }},
			"weight": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) { return 1, nil }},
		},
		IsTypeOf: func(p graphql.IsTypeOfParams) bool { return true },
	})
	q := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"echoInt": &graphql.Field{
				Type: graphql.Int,
				Args: graphql.FieldConfigArgument{
					"v": &graphql.ArgumentConfig{Type: graphql.Int},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) { return p.Args["v"], nil },
			},
			"echoFloat": &graphql.Field{
				Type: graphql.Float,
				Args: graphql.FieldConfigArgument{
					"v": &graphql.ArgumentConfig{Type: graphql.Float},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) { return p.Args["v"], nil },
			},
			"echoBool": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"v": &graphql.ArgumentConfig{Type: graphql.Boolean},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) { return p.Args["v"], nil },
			},
			"echoColor": &graphql.Field{
				Type: color,
				Args: graphql.FieldConfigArgument{
					"v": &graphql.ArgumentConfig{Type: color},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) { return p.Args["v"], nil },
			},
			"echoIntList": &graphql.Field{
				Type: graphql.NewList(graphql.Int),
				Args: graphql.FieldConfigArgument{
					"vs": &graphql.ArgumentConfig{Type: graphql.NewList(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) { return p.Args["vs"], nil },
			},
			"echoNonNull": &graphql.Field{
				Type: graphql.Int,
				Args: graphql.FieldConfigArgument{
					"v": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) { return p.Args["v"], nil },
			},
			"echoPoint": &graphql.Field{
				Type: graphql.Int,
				Args: graphql.FieldConfigArgument{
					"p": &graphql.ArgumentConfig{Type: point},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if m, ok := p.Args["p"].(map[string]interface{}); ok {
						return m["x"], nil
					}
					return 0, nil
				},
			},
			"thing": &graphql.Field{
				Type:    thingIface,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) { return struct{}{}, nil },
			},
			"widget": &graphql.Field{
				Type:    widget,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) { return struct{}{}, nil },
			},
		},
	})
	thingIface.ResolveType = func(p graphql.ResolveTypeParams) *graphql.Object { return widget }
	s, err := graphql.NewSchema(graphql.SchemaConfig{Query: q, Types: []graphql.Type{widget}})
	if err != nil {
		t.Fatalf("schema: %v", err)
	}
	return s
}

// TestPlanCacheNormalizeArgKinds drives every arm of the literal
// extractor + fingerprint writeValue: Int / Float / Bool / Enum /
// List / InputObject. Two identical-shape queries with different
// literals must collapse to the same plan; a different-shape query
// (different field) must miss.
func TestPlanCacheNormalizeArgKinds(t *testing.T) {
	schema := buildNormalizeTestSchema(t)
	cache := graphql.NewPlanCache(graphql.PlanCacheOptions{Normalize: true})

	cases := []struct {
		name string
		a, b string
	}{
		{"Int", `{ echoInt(v: 1) }`, `{ echoInt(v: 999) }`},
		{"Float", `{ echoFloat(v: 1.5) }`, `{ echoFloat(v: 9.5) }`},
		{"Bool", `{ echoBool(v: true) }`, `{ echoBool(v: false) }`},
		{"Enum", `{ echoColor(v: RED) }`, `{ echoColor(v: BLUE) }`},
		{"List", `{ echoIntList(vs: [1,2,3]) }`, `{ echoIntList(vs: [9,8]) }`},
		{"NonNull", `{ echoNonNull(v: 1) }`, `{ echoNonNull(v: 99) }`},
		{"InputObject", `{ echoPoint(p: {x: 1, y: 2}) }`, `{ echoPoint(p: {x: 9, y: 9}) }`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r1 := cache.Get(&schema, tc.a, "")
			if len(r1.Errors) > 0 || r1.Plan == nil {
				t.Fatalf("a: %v", r1.Errors)
			}
			r2 := cache.Get(&schema, tc.b, "")
			if len(r2.Errors) > 0 {
				t.Fatalf("b: %v", r2.Errors)
			}
			if r1.Plan != r2.Plan {
				t.Fatalf("normalized literal-only variations must share *Plan (kind=%s)", tc.name)
			}
		})
	}
}

// TestPlanCacheNormalizeFragmentSpreadFingerprint verifies that
// fragment definitions participate in the cache key — two queries
// whose only difference is the body of a spread fragment must miss.
func TestPlanCacheNormalizeFragmentSpreadFingerprint(t *testing.T) {
	schema := buildNormalizeTestSchema(t)
	cache := graphql.NewPlanCache(graphql.PlanCacheOptions{Normalize: true})

	q1 := `query { widget { ...F } } fragment F on Widget { name }`
	q2 := `query { widget { ...F } } fragment F on Widget { weight }`

	r1 := cache.Get(&schema, q1, "")
	if len(r1.Errors) > 0 || r1.Plan == nil {
		t.Fatalf("q1: %v", r1.Errors)
	}
	r2 := cache.Get(&schema, q2, "")
	if len(r2.Errors) > 0 || r2.Plan == nil {
		t.Fatalf("q2: %v", r2.Errors)
	}
	if r1.Plan == r2.Plan {
		t.Fatalf("queries with different fragment bodies must produce different *Plan")
	}
}

// Same fragment body, two queries that spread it identically must
// share a plan — exercises writeFragmentBody's visited-cache path.
func TestPlanCacheNormalizeFragmentSpreadShared(t *testing.T) {
	schema := buildNormalizeTestSchema(t)
	cache := graphql.NewPlanCache(graphql.PlanCacheOptions{Normalize: true})

	q := `query { widget { ...F ...F } } fragment F on Widget { name weight }`
	r1 := cache.Get(&schema, q, "")
	if len(r1.Errors) > 0 || r1.Plan == nil {
		t.Fatalf("q: %v", r1.Errors)
	}
	r2 := cache.Get(&schema, q, "")
	if r1.Plan != r2.Plan {
		t.Fatalf("identical query must hit cache")
	}
}

// TestPlanCacheNormalizeInlineFragmentPath drills the InlineFragment
// arm of normalizeSelectionSet — the planner must follow inline
// fragments' selection sets during literal extraction.
func TestPlanCacheNormalizeInlineFragmentPath(t *testing.T) {
	schema := buildNormalizeTestSchema(t)
	cache := graphql.NewPlanCache(graphql.PlanCacheOptions{Normalize: true})

	q1 := `{ thing { ... on Widget { name } } }`
	q2 := `{ thing { ... on Widget { weight } } }`

	r1 := cache.Get(&schema, q1, "")
	if len(r1.Errors) > 0 || r1.Plan == nil {
		t.Fatalf("q1: %v", r1.Errors)
	}
	r2 := cache.Get(&schema, q2, "")
	if r1.Plan == r2.Plan {
		t.Fatalf("inline fragments selecting different fields must miss the cache")
	}
}

// TestPlanCacheNormalizeUnextractableLiteralsFingerprinted exercises
// the literal-kind arms of writeValue. Values containing a nested
// variable are not extractable into synth vars, so the literal
// survives in the document and writeValue's Int/Float/Bool/Enum/
// String/List/ObjectValue branches must produce kind-distinct
// fingerprints.
func TestPlanCacheNormalizeUnextractableLiteralsFingerprinted(t *testing.T) {
	schema := buildNormalizeTestSchema(t)
	cache := graphql.NewPlanCache(graphql.PlanCacheOptions{Normalize: true})

	// List with nested variable can't be extracted as a whole.
	a := `query Q($x: Int) { echoIntList(vs: [1, $x, 3]) }`
	b := `query Q($x: Int) { echoIntList(vs: [9, $x, 9]) }`
	c := `query Q($x: Int) { echoIntList(vs: [1, $x, 9]) }`

	r1 := cache.Get(&schema, a, "Q")
	if len(r1.Errors) > 0 || r1.Plan == nil {
		t.Fatalf("a: %v", r1.Errors)
	}
	r2 := cache.Get(&schema, b, "Q")
	if len(r2.Errors) > 0 {
		t.Fatalf("b: %v", r2.Errors)
	}
	if r1.Plan == r2.Plan {
		t.Fatalf("un-extractable literal differences must not collapse to one plan")
	}
	r3 := cache.Get(&schema, c, "Q")
	if len(r3.Errors) > 0 {
		t.Fatalf("c: %v", r3.Errors)
	}
	if r1.Plan == r3.Plan || r2.Plan == r3.Plan {
		t.Fatalf("each distinct literal layout must produce its own plan")
	}
}

// TestPlanCacheNormalizeVariableTypes exercises writeType against
// NonNull and List wrappings in the operation's VariableDefinitions.
func TestPlanCacheNormalizeVariableTypes(t *testing.T) {
	schema := buildNormalizeTestSchema(t)
	cache := graphql.NewPlanCache(graphql.PlanCacheOptions{Normalize: true})

	queries := []string{
		`query Q($v: Int!) { echoNonNull(v: $v) }`,
		`query Q($v: Int) { echoInt(v: $v) }`,
		`query Q($vs: [Int]) { echoIntList(vs: $vs) }`,
	}
	plans := map[*graphql.Plan]bool{}
	for _, q := range queries {
		r := cache.Get(&schema, q, "Q")
		if len(r.Errors) > 0 || r.Plan == nil {
			t.Fatalf("%q: %v", q, r.Errors)
		}
		plans[r.Plan] = true
	}
	if len(plans) != len(queries) {
		t.Fatalf("expected %d distinct plans across NonNull/Nullable/List variable shapes, got %d", len(queries), len(plans))
	}
}
