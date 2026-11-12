package benchutil

import (
	"fmt"

	"github.com/graphql-go/graphql"
)

func WideSchemaWithXFieldsAndYItems(x int, y int) graphql.Schema {
	wide := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Wide",
		Description: "An object",
		Fields:      generateXWideFields(x),
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"wide": {
				Type: graphql.NewList(wide),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					out := make([]struct{}, 0, y)
					for i := 0; i < y; i++ {
						out = append(out, struct{}{})
					}
					return out, nil
				},
			},
		},
	})

	wideSchema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})

	return wideSchema
}

func generateXWideFields(x int) graphql.Fields {
	fields := graphql.Fields{}
	for i := 0; i < x; i++ {
		fields[generateFieldNameFromX(i)] = generateWideFieldFromX(i)
	}
	return fields
}

func generateWideFieldFromX(x int) *graphql.Field {
	return &graphql.Field{
		Type:    generateWideTypeFromX(x),
		Resolve: generateWideResolveFromX(x),
	}
}

func generateWideTypeFromX(x int) graphql.Type {
	switch x % 8 {
	case 0:
		return graphql.String
	case 1:
		return graphql.NewNonNull(graphql.String)
	case 2:
		return graphql.Int
	case 3:
		return graphql.NewNonNull(graphql.Int)
	case 4:
		return graphql.Float
	case 5:
		return graphql.NewNonNull(graphql.Float)
	case 6:
		return graphql.Boolean
	case 7:
		return graphql.NewNonNull(graphql.Boolean)
	}

	return nil
}

func generateFieldNameFromX(x int) string {
	var out string
	alphabet := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "z"}
	v := x
	for {
		r := v % 10
		out = alphabet[r] + out
		v /= 10
		if v == 0 {
			break
		}
	}
	return out
}

func generateWideResolveFromX(x int) func(p graphql.ResolveParams) (interface{}, error) {
	switch x % 8 {
	case 0:
		return func(p graphql.ResolveParams) (interface{}, error) {
			return fmt.Sprint(x), nil
		}
	case 1:
		return func(p graphql.ResolveParams) (interface{}, error) {
			return fmt.Sprint(x), nil
		}
	case 2:
		return func(p graphql.ResolveParams) (interface{}, error) {
			return x, nil
		}
	case 3:
		return func(p graphql.ResolveParams) (interface{}, error) {
			return x, nil
		}
	case 4:
		return func(p graphql.ResolveParams) (interface{}, error) {
			return float64(x), nil
		}
	case 5:
		return func(p graphql.ResolveParams) (interface{}, error) {
			return float64(x), nil
		}
	case 6:
		return func(p graphql.ResolveParams) (interface{}, error) {
			if x%2 == 0 {
				return false, nil
			}
			return true, nil
		}
	case 7:
		return func(p graphql.ResolveParams) (interface{}, error) {
			if x%2 == 0 {
				return false, nil
			}
			return true, nil
		}
	}

	return nil
}

func WideSchemaQuery(x int) string {
	var fields string
	for i := 0; i < x; i++ {
		fields = fields + generateFieldNameFromX(i) + " "
	}

	return fmt.Sprintf("query { wide { %s} }", fields)
}

// WideArgedSchemaWithXFieldsAndYItems is the same shape as
// WideSchemaWithXFieldsAndYItems but every field takes a `value:
// String` argument the resolver echoes (or, for non-string fields,
// ignores). Designed to exercise the parametric-query path: a single
// `value: $v` variable can fan out across all 100 fields, so a
// cached plan handles arbitrary literal variations via Args alone.
func WideArgedSchemaWithXFieldsAndYItems(x int, y int) graphql.Schema {
	wide := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Wide",
		Description: "An object",
		Fields:      generateXArgedWideFields(x),
	})
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"wide": {
				Type: graphql.NewList(wide),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					out := make([]struct{}, 0, y)
					for i := 0; i < y; i++ {
						out = append(out, struct{}{})
					}
					return out, nil
				},
			},
		},
	})
	s, _ := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	return s
}

func generateXArgedWideFields(x int) graphql.Fields {
	fields := graphql.Fields{}
	for i := 0; i < x; i++ {
		fields[generateFieldNameFromX(i)] = generateArgedWideFieldFromX(i)
	}
	return fields
}

func generateArgedWideFieldFromX(x int) *graphql.Field {
	return &graphql.Field{
		Type: generateWideTypeFromX(x),
		Args: graphql.FieldConfigArgument{
			"value": &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: generateArgedWideResolveFromX(x),
	}
}

// generateArgedWideResolveFromX echoes the arg for string fields and
// produces a deterministic fixed value for non-string fields. The
// argument variation alone is enough to exercise per-request arg
// coercion when present.
func generateArgedWideResolveFromX(x int) func(p graphql.ResolveParams) (interface{}, error) {
	switch x % 8 {
	case 0, 1:
		return func(p graphql.ResolveParams) (interface{}, error) {
			if v, ok := p.Args["value"].(string); ok {
				return v, nil
			}
			return fmt.Sprint(x), nil
		}
	case 2, 3:
		return func(p graphql.ResolveParams) (interface{}, error) { return x, nil }
	case 4, 5:
		return func(p graphql.ResolveParams) (interface{}, error) { return float64(x), nil }
	case 6, 7:
		return func(p graphql.ResolveParams) (interface{}, error) { return x%2 == 1, nil }
	}
	return nil
}

// WideArgedSchemaQueryWithVariable returns a parametric query: every
// field's `value` arg is bound to the same variable `$v`. The
// resulting *Plan can be cached once and reused for arbitrary literal
// variations passed as Args at execute time.
func WideArgedSchemaQueryWithVariable(x int) string {
	var fields string
	for i := 0; i < x; i++ {
		fields = fields + generateFieldNameFromX(i) + "(value: $v) "
	}
	return fmt.Sprintf("query Q($v: String) { wide { %s} }", fields)
}

// WideArgedSchemaQueryWithLiteral returns a query identical in shape
// to WideArgedSchemaQueryWithVariable except each `value` arg is a
// distinct literal string. Used to demonstrate the cost when clients
// don't use variables — every literal change forces a fresh
// parse+validate+plan.
func WideArgedSchemaQueryWithLiteral(x int, literal string) string {
	var fields string
	for i := 0; i < x; i++ {
		fields = fields + generateFieldNameFromX(i) + fmt.Sprintf("(value: %q) ", literal)
	}
	return fmt.Sprintf("query { wide { %s} }", fields)
}
