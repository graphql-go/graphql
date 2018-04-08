package benchutil

import (
	"fmt"

	"github.com/graphql-go/graphql"
)

type color struct {
	Hex string
	R   int
	G   int
	B   int
}

func ListSchemaWithXItems(x int) graphql.Schema {

	list := generateXListItems(x)

	color := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Color",
		Description: "A color",
		Fields: graphql.Fields{
			"hex": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Hex color code.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if c, ok := p.Source.(color); ok {
						return c.Hex, nil
					}
					return nil, nil
				},
			},
			"r": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Red value.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if c, ok := p.Source.(color); ok {
						return c.R, nil
					}
					return nil, nil
				},
			},
			"g": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Green value.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if c, ok := p.Source.(color); ok {
						return c.G, nil
					}
					return nil, nil
				},
			},
			"b": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Blue value.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if c, ok := p.Source.(color); ok {
						return c.B, nil
					}
					return nil, nil
				},
			},
		},
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"colors": {
				Type: graphql.NewList(color),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return list, nil
				},
			},
		},
	})

	colorSchema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})

	return colorSchema
}

var colors []color

func init() {
	colors = make([]color, 0, 256*16*16)

	for r := 0; r < 256; r++ {
		for g := 0; g < 16; g++ {
			for b := 0; b < 16; b++ {
				colors = append(colors, color{
					Hex: fmt.Sprintf("#%x%x%x", r, g, b),
					R:   r,
					G:   g,
					B:   b,
				})
			}
		}
	}
}

func generateXListItems(x int) []color {
	if x > len(colors) {
		x = len(colors)
	}
	return colors[0:x]
}
