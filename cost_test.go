package graphql

import (
	"fmt"
	"testing"

	"github.com/GannettDigital/graphql/language/parser"
	"github.com/GannettDigital/graphql/language/source"
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
	interfaceFields := Fields{
		"b": &Field{
			Cost: 1,
			Type: String,
		},
	}
	deepDataInterface := NewInterface(InterfaceConfig{
		Name:   "deepD",
		Fields: interfaceFields,
	})

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
		Interfaces: []*Interface{deepDataInterface},
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
		Name:       "DeepDataType",
		Fields:     deepDataFields,
		Interfaces: []*Interface{deepDataInterface},
	})

	dataType.AddFieldConfig("deep", &Field{
		Cost: 25,
		Type: deepDataType,
	})
	dataType.AddFieldConfig("promise", &Field{
		Cost: 25,
		Type: dataType,
	})

	dataType.AddFieldConfig("iface", &Field{
		Cost: 50,
		Type: deepDataInterface,
	})

	queryCfg := ObjectConfig{
		Name: "query",
		Fields: Fields{
			"example": &Field{
				Type: dataType,
			},
		},
	}

	query := NewObject(queryCfg)
	schema, err := NewSchema(SchemaConfig{
		Query: query,
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	tests := []struct {
		description string
		query       string
		want        int
	}{
		// Note a test with an unused fragment isn't needed as that fails and a invalid query
		{
			description: "Simple Query",
			query: `{
					  example {
						a,
						b
					  }
					}`,
			want: 2,
		},
		{
			description: "Medium Complexity Query",
			query: `{
 						example {
							a,
							b,
							deep {
								a
								b
								c
							}
						}
					}`,
			want: 30,
		},
		{
			description: "Medium Complexity with inline fragment",
			query: `{
						example {
							a,
							b,
							...on DataType {
								promise {
									a
								}
							}
							deep {
								a
								b
								c
							}
						}
					}`,
			want: 56,
		},
		{
			description: "Complex Query",
			query: `query a($size: Int) {
						example {
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
					}

					fragment c on DataType {
						d
						e
					}`,
			want: 172,
		},
		{
			description: "Query with Interface",
			query: `{
						example {
							a,
							b,
							iface {
								b
							}
						}
					}`,
			want: 53,
		},
		{
			description: "Query with Interface and inline fragment",
			query: `{
						example {
							a,
							b,
							iface {
								... on DeepDataType {
									a
									b
									c
			                    }
							}
						}
					}`,
			want: 55,
		},
		{
			description: "Query with Interface and fragment on the interface",
			query: `{
						example {
							a,
							b,
							iface {
						    	... interfaceFrag
						    }
						}
					}
			        fragment interfaceFrag on deepD {
						... on DeepDataType {
							a
							b
							c
			            }
					}`,
			want: 55,
		},
		{
			description: "Query with Interface in fragment with inline fragment",
			query: `{
						example {
							... interfaceFrag
						}
					}
			        fragment interfaceFrag on DataType {
						a,
						b,
						iface {
							... on DeepDataType {
								a
								b
								c
							}
						}
					}`,
			want: 55,
		},
		{
			description: "Query with Interface and multiple inline fragment",
			query: `{
						example {
							a,
							b,
							iface {
								... on DeepDataType {
									a
									b
									c
                                }
								... on DataType {
									b
									e
									f
									pic(size: 1)
                                }
							}
						}
					}`,
			want: 68,
		},
		{
			description: "Query with Interface and multiple inline fragment in a fragment",
			query: `{
					  example {
			            ... interfaceFrag
			          }
			        }
                    fragment interfaceFrag on DataType {
						a,
						b,
						iface {
							... on DeepDataType {
								a
								b
								c
                            }
							... on DataType {
								b
								e
								f
								pic(size: 1)
                            }
						}
					}`,
			want: 68,
		},
		{
			description: "Medium Complexity Query as fragment",
			query: `{
			        	example {
			            	... mediumFrag
			            }
			        }
			        fragment mediumFrag on DataType {
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
			description: "Complex Query as fragment",
			query: `query a($size: Int){
					  example {
					   ... complexFrag
					  }
				    }
				    fragment complexFrag on DataType {
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
	}
	for _, test := range tests {
		s := source.NewSource(&source.Source{
			Body: []byte(test.query),
			Name: "GraphQL request",
		})
		astDoc, err := parser.Parse(parser.ParseParams{Source: s})
		if err != nil {
			t.Fatalf("Test %q - Parse failed: %v", test.description, err)
		}

		validationResult := ValidateDocument(&schema, astDoc, nil)
		if !validationResult.IsValid {
			t.Errorf("Test %q - failed validation: %v", test.description, validationResult)
			continue
		}

		ep := ExecuteParams{
			Schema: schema,
			Root:   data,
			AST:    astDoc,
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
