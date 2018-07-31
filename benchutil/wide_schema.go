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
