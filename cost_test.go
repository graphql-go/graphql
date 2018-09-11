package graphql

import (
	"fmt"
	"testing"

	"github.com/GannettDigital/graphql/language/parser"
)

func TestQueryComplexity(t *testing.T) {
	// This is based off of TestExecutesArbitraryCode in executor_test.go

	deepData := map[string]interface{}{}
	data := map[string]interface{}{
		"a": func() interface{} { return "Apple" },
		"b": func() interface{} { return "Banana" },
		"c": func() interface{} { return "Cookie" },
		"d": func() interface{} { return "Donut" },
		"e": func() interface{} { return "Egg" },
		"f": "Fish",
		"pic": func(size int) string {
			return fmt.Sprintf("Pic of size: %v", size)
		},
		"deep": func() interface{} { return deepData },
	}
	data["promise"] = func() interface{} {
		return data
	}
	deepData = map[string]interface{}{
		"a":      func() interface{} { return "Already Been Done" },
		"b":      func() interface{} { return "Boring" },
		"c":      func() interface{} { return []string{"Contrived", "", "Confusing"} },
		"deeper": func() interface{} { return []interface{}{data, nil, data} },
	}

	// Schema Definitions
	picResolverFn := func(p ResolveParams) (interface{}, error) {
		// get and type assert ResolveFn for this field
		picResolver, ok := p.Source.(map[string]interface{})["pic"].(func(size int) string)
		if !ok {
			return nil, nil
		}
		// get and type assert argument
		sizeArg, ok := p.Args["size"].(int)
		if !ok {
			return nil, nil
		}
		return picResolver(sizeArg), nil
	}
	dataType := NewObject(ObjectConfig{
		Name: "DataType",
		Fields: Fields{
			"a": &Field{
				Cost: 1,
				Type: NewNonNull(String),
			},
			"b": &Field{
				Cost: 1,
				Type: String,
			},
			"c": &Field{
				Cost: 1,
				Type: String,
			},
			"d": &Field{
				Cost: 1,
				Type: String,
			},
			"e": &Field{
				Cost: 1,
				Type: String,
			},
			"f": &Field{
				Cost: 1,
				Type: String,
			},
			"pic": &Field{
				Cost: 10,
				Args: FieldConfigArgument{
					"size": &ArgumentConfig{
						Type: Int,
					},
				},
				Type:    String,
				Resolve: picResolverFn,
			},
		},
	})
	deepDataFields := Fields{
		"a": &Field{
			Cost: 1,
			Type: String,
		},
		"b": &Field{
			Cost: 1,
			Type: String,
		},
		"c": &Field{
			Cost: 1,
			Type: NewNonNull(NewList(String)),
		},
		"deeper": &Field{
			Cost: 100,
			Type: NewList(dataType),
		},
	}
	deepDataType := NewObject(ObjectConfig{
		Name:   "DeepDataType",
		Fields: deepDataFields,
	})

	dataType.AddFieldConfig("deep", &Field{
		Cost: 25,
		Type: deepDataType,
	})
	dataType.AddFieldConfig("promise", &Field{
		Cost: 25,
		Type: dataType,
	})

	deepDataInterface := NewInterface(InterfaceConfig{
		Name:   "deepD",
		Fields: deepDataFields,
		ResolveType: func(p ResolveTypeParams) *Object {
			return deepDataType
		},
	})

	dataType.AddFieldConfig("iface", &Field{
		Cost: 50,
		Type: deepDataInterface,
	})

	tests := []struct {
		description string
		query       string
		want        int
	}{
		{
			description: "Simple Query",
			query: `
				  query Example($size: Int) {
					a,
					b
				  }`,
			want: 2,
		},
		{
			description: "Medium Complexity Query",
			query: `
				  query Example($size: Int) {
					a,
					b,
					deep {
					  a
					  b
					  c
					}
				  }`,
			want: 30,
		},
		{
			description: "Complex Query",
			query: `
			  query Example($size: Int) {
				a,
				b,
				x: c
        		...c
				f
				...on DataType {
				  pic(size: $size)
				  promise {
					a
				  }
				}
				deep {
				  a
				  b
				  c
				  deeper {
					a
					b
				  }
				}
			  }

			  fragment c on DataType {
				d
				e
			  }`,
			want: 172,
		},
		{
			description: "Query with Interface",
			query: `
				  query Example($size: Int) {
					a,
					b,
					iface {
					  a
					  b
					  c
					}
				  }`,
			want: 55,
		},
	}
	schema, err := NewSchema(SchemaConfig{
		Query: dataType,
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	for _, test := range tests {

		// parse query
		astDoc, err := parser.Parse(parser.ParseParams{
			Source: test.query,
			Options: parser.ParseOptions{
				// include source, for error reporting
				NoSource: false,
			},
		})
		if err != nil {
			t.Fatalf("Test %q - Parse failed: %v", test.description, err)
		}

		operationName := "Example"
		ep := ExecuteParams{
			Schema:        schema,
			Root:          data,
			AST:           astDoc,
			OperationName: operationName,
		}
		got, err := QueryComplexity(ep)
		if err != nil {
			t.Errorf("Test %q - failed running query complexity: %v", test.description, err)
		}
		if got != test.want {
			t.Errorf("Test %q - got %d, want %d", test.description, got, test.want)
		}
	}
}
