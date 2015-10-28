package graphql

import (
	"reflect"
	"testing"
)

func g(t *testing.T, p Params) *Result {
	resultChannel := make(chan *Result)
	go Graphql(p, resultChannel)
	result := <-resultChannel
	return result
}

func TestIntrospection_ExecutesAnIntrospectionQuery(t *testing.T) {
	emptySchema, err := NewSchema(SchemaConfig{
		Query: NewObject(ObjectConfig{
			Name: "QueryRoot",
			Fields: FieldConfigMap{
				"onlyField": &FieldConfig{
					Type: String,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error creating Schema: %v", err.Error())
	}
	expectedDataSubSet := map[string]interface{}{
		"__schema": map[string]interface{}{
			"mutationType": nil,
			"queryType": map[string]interface{}{
				"name": "QueryRoot",
			},
			"types": []interface{}{
				map[string]interface{}{
					"kind":          "OBJECT",
					"name":          "QueryRoot",
					"inputFields":   nil,
					"interfaces":    []interface{}{},
					"enumValues":    nil,
					"possibleTypes": nil,
				},
				map[string]interface{}{
					"kind": "OBJECT",
					"name": "__Schema",
					"fields": []interface{}{
						map[string]interface{}{
							"name": "types",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]interface{}{
									"kind": "LIST",
									"name": nil,
									"ofType": map[string]interface{}{
										"kind": "NON_NULL",
										"name": nil,
										"ofType": map[string]interface{}{
											"kind": "OBJECT",
											"name": "__Type",
										},
									},
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name": "queryType",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]interface{}{
									"kind": "OBJECT",
									"name": "__Type",
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name": "mutationType",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind": "OBJECT",
								"name": "__Type",
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name": "directives",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]interface{}{
									"kind": "LIST",
									"name": nil,
									"ofType": map[string]interface{}{
										"kind": "NON_NULL",
										"name": nil,
										"ofType": map[string]interface{}{
											"kind": "OBJECT",
											"name": "__Directive",
										},
									},
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
					},
					"inputFields":   nil,
					"interfaces":    []interface{}{},
					"enumValues":    nil,
					"possibleTypes": nil,
				},
				map[string]interface{}{
					"kind": "OBJECT",
					"name": "__Type",
					"fields": []interface{}{
						map[string]interface{}{
							"name": "kind",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]interface{}{
									"kind":   "ENUM",
									"name":   "__TypeKind",
									"ofType": nil,
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name": "name",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind":   "SCALAR",
								"name":   "String",
								"ofType": nil,
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name": "description",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind":   "SCALAR",
								"name":   "String",
								"ofType": nil,
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name": "fields",
							"args": []interface{}{
								map[string]interface{}{
									"name": "includeDeprecated",
									"type": map[string]interface{}{
										"kind":   "SCALAR",
										"name":   "Boolean",
										"ofType": nil,
									},
									"defaultValue": "false",
								},
							},
							"type": map[string]interface{}{
								"kind": "LIST",
								"name": nil,
								"ofType": map[string]interface{}{
									"kind": "NON_NULL",
									"name": nil,
									"ofType": map[string]interface{}{
										"kind":   "OBJECT",
										"name":   "__Field",
										"ofType": nil,
									},
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name": "interfaces",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind": "LIST",
								"name": nil,
								"ofType": map[string]interface{}{
									"kind": "NON_NULL",
									"name": nil,
									"ofType": map[string]interface{}{
										"kind":   "OBJECT",
										"name":   "__Type",
										"ofType": nil,
									},
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name": "possibleTypes",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind": "LIST",
								"name": nil,
								"ofType": map[string]interface{}{
									"kind": "NON_NULL",
									"name": nil,
									"ofType": map[string]interface{}{
										"kind":   "OBJECT",
										"name":   "__Type",
										"ofType": nil,
									},
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name": "enumValues",
							"args": []interface{}{
								map[string]interface{}{
									"name": "includeDeprecated",
									"type": map[string]interface{}{
										"kind":   "SCALAR",
										"name":   "Boolean",
										"ofType": nil,
									},
									"defaultValue": "false",
								},
							},
							"type": map[string]interface{}{
								"kind": "LIST",
								"name": nil,
								"ofType": map[string]interface{}{
									"kind": "NON_NULL",
									"name": nil,
									"ofType": map[string]interface{}{
										"kind":   "OBJECT",
										"name":   "__EnumValue",
										"ofType": nil,
									},
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name": "inputFields",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind": "LIST",
								"name": nil,
								"ofType": map[string]interface{}{
									"kind": "NON_NULL",
									"name": nil,
									"ofType": map[string]interface{}{
										"kind":   "OBJECT",
										"name":   "__InputValue",
										"ofType": nil,
									},
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name": "ofType",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind":   "OBJECT",
								"name":   "__Type",
								"ofType": nil,
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
					},
					"inputFields":   nil,
					"interfaces":    []interface{}{},
					"enumValues":    nil,
					"possibleTypes": nil,
				},
				map[string]interface{}{
					"kind":        "ENUM",
					"name":        "__TypeKind",
					"fields":      nil,
					"inputFields": nil,
					"interfaces":  nil,
					"enumValues": []interface{}{
						map[string]interface{}{
							"name":              "SCALAR",
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name":              "OBJECT",
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name":              "INTERFACE",
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name":              "UNION",
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name":              "ENUM",
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name":              "INPUT_OBJECT",
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name":              "LIST",
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name":              "NON_NULL",
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
					},
					"possibleTypes": nil,
				},
				map[string]interface{}{
					"kind":          "SCALAR",
					"name":          "String",
					"fields":        nil,
					"inputFields":   nil,
					"interfaces":    nil,
					"enumValues":    nil,
					"possibleTypes": nil,
				},
				map[string]interface{}{
					"kind":          "SCALAR",
					"name":          "Boolean",
					"fields":        nil,
					"inputFields":   nil,
					"interfaces":    nil,
					"enumValues":    nil,
					"possibleTypes": nil,
				},
				map[string]interface{}{
					"kind": "OBJECT",
					"name": "__Field",
					"fields": []interface{}{
						map[string]interface{}{
							"name": "name",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]interface{}{
									"kind":   "SCALAR",
									"name":   "String",
									"ofType": nil,
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name": "description",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind":   "SCALAR",
								"name":   "String",
								"ofType": nil,
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name": "args",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]interface{}{
									"kind": "LIST",
									"name": nil,
									"ofType": map[string]interface{}{
										"kind": "NON_NULL",
										"name": nil,
										"ofType": map[string]interface{}{
											"kind": "OBJECT",
											"name": "__InputValue",
										},
									},
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name": "type",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]interface{}{
									"kind":   "OBJECT",
									"name":   "__Type",
									"ofType": nil,
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name": "isDeprecated",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]interface{}{
									"kind":   "SCALAR",
									"name":   "Boolean",
									"ofType": nil,
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name": "deprecationReason",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind":   "SCALAR",
								"name":   "String",
								"ofType": nil,
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
					},
					"inputFields":   nil,
					"interfaces":    []interface{}{},
					"enumValues":    nil,
					"possibleTypes": nil,
				},
				map[string]interface{}{
					"kind": "OBJECT",
					"name": "__InputValue",
					"fields": []interface{}{
						map[string]interface{}{
							"name": "name",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]interface{}{
									"kind":   "SCALAR",
									"name":   "String",
									"ofType": nil,
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name": "description",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind":   "SCALAR",
								"name":   "String",
								"ofType": nil,
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name": "type",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]interface{}{
									"kind":   "OBJECT",
									"name":   "__Type",
									"ofType": nil,
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name": "defaultValue",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind":   "SCALAR",
								"name":   "String",
								"ofType": nil,
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
					},
					"inputFields":   nil,
					"interfaces":    []interface{}{},
					"enumValues":    nil,
					"possibleTypes": nil,
				},
				map[string]interface{}{
					"kind": "OBJECT",
					"name": "__EnumValue",
					"fields": []interface{}{
						map[string]interface{}{
							"name": "name",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]interface{}{
									"kind":   "SCALAR",
									"name":   "String",
									"ofType": nil,
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name": "description",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind":   "SCALAR",
								"name":   "String",
								"ofType": nil,
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name": "isDeprecated",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]interface{}{
									"kind":   "SCALAR",
									"name":   "Boolean",
									"ofType": nil,
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name": "deprecationReason",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind":   "SCALAR",
								"name":   "String",
								"ofType": nil,
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
					},
					"inputFields":   nil,
					"interfaces":    []interface{}{},
					"enumValues":    nil,
					"possibleTypes": nil,
				},
				map[string]interface{}{
					"kind": "OBJECT",
					"name": "__Directive",
					"fields": []interface{}{
						map[string]interface{}{
							"name": "name",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]interface{}{
									"kind":   "SCALAR",
									"name":   "String",
									"ofType": nil,
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name": "description",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind":   "SCALAR",
								"name":   "String",
								"ofType": nil,
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name": "args",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]interface{}{
									"kind": "LIST",
									"name": nil,
									"ofType": map[string]interface{}{
										"kind": "NON_NULL",
										"name": nil,
										"ofType": map[string]interface{}{
											"kind": "OBJECT",
											"name": "__InputValue",
										},
									},
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name": "onOperation",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]interface{}{
									"kind":   "SCALAR",
									"name":   "Boolean",
									"ofType": nil,
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name": "onFragment",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]interface{}{
									"kind":   "SCALAR",
									"name":   "Boolean",
									"ofType": nil,
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]interface{}{
							"name": "onField",
							"args": []interface{}{},
							"type": map[string]interface{}{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]interface{}{
									"kind":   "SCALAR",
									"name":   "Boolean",
									"ofType": nil,
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
					},
					"inputFields":   nil,
					"interfaces":    []interface{}{},
					"enumValues":    nil,
					"possibleTypes": nil,
				},
			},
			"directives": []interface{}{
				map[string]interface{}{
					"name": "include",
					"args": []interface{}{
						map[string]interface{}{
							"defaultValue": nil,
							"name":         "if",
							"type": map[string]interface{}{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]interface{}{
									"kind":   "SCALAR",
									"name":   "Boolean",
									"ofType": nil,
								},
							},
						},
					},
					"onOperation": false,
					"onFragment":  true,
					"onField":     true,
				},
				map[string]interface{}{
					"name": "skip",
					"args": []interface{}{
						map[string]interface{}{
							"defaultValue": nil,
							"name":         "if",
							"type": map[string]interface{}{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]interface{}{
									"kind":   "SCALAR",
									"name":   "Boolean",
									"ofType": nil,
								},
							},
						},
					},
					"onOperation": false,
					"onFragment":  true,
					"onField":     true,
				},
			},
		},
	}
	result := g(t, Params{
		Schema:        emptySchema,
		RequestString: IntrospectionQuery,
	})
	if !ContainSubset(result.Data.(map[string]interface{}), expectedDataSubSet) {
		t.Fatalf("unexpected, result does not contain subset of expected data")
	}
}

func TestIntrospection_ExecutesAnInputObject(t *testing.T) {

	testInputObject := NewInputObject(InputObjectConfig{
		Name: "TestInputObject",
		Fields: InputObjectConfigFieldMap{
			"a": &InputObjectFieldConfig{
				Type:         String,
				DefaultValue: "foo",
			},
			"b": &InputObjectFieldConfig{
				Type: NewList(String),
			},
		},
	})
	testType := NewObject(ObjectConfig{
		Name: "TestType",
		Fields: FieldConfigMap{
			"field": &FieldConfig{
				Type: String,
				Args: FieldConfigArgument{
					"complex": &ArgumentConfig{
						Type: testInputObject,
					},
				},
				Resolve: func(p GQLFRParams) interface{} {
					return p.Args["complex"]
				},
			},
		},
	})
	schema, err := NewSchema(SchemaConfig{
		Query: testType,
	})
	if err != nil {
		t.Fatalf("Error creating Schema: %v", err.Error())
	}
	query := `
      {
        __schema {
          types {
            kind
            name
            inputFields {
              name
              type { ...TypeRef }
              defaultValue
            }
          }
        }
      }

      fragment TypeRef on __Type {
        kind
        name
        ofType {
          kind
          name
          ofType {
            kind
            name
            ofType {
              kind
              name
            }
          }
        }
      }
    `
	expectedDataSubSet := map[string]interface{}{
		"__schema": map[string]interface{}{
			"types": []interface{}{
				map[string]interface{}{
					"kind": "INPUT_OBJECT",
					"name": "TestInputObject",
					"inputFields": []interface{}{
						map[string]interface{}{
							"name": "a",
							"type": map[string]interface{}{
								"kind":   "SCALAR",
								"name":   "String",
								"ofType": nil,
							},
							"defaultValue": `"foo"`,
						},
						map[string]interface{}{
							"name": "b",
							"type": map[string]interface{}{
								"kind": "LIST",
								"name": nil,
								"ofType": map[string]interface{}{
									"kind":   "SCALAR",
									"name":   "String",
									"ofType": nil,
								},
							},
							"defaultValue": nil,
						},
					},
				},
			},
		},
	}

	result := g(t, Params{
		Schema:        schema,
		RequestString: query,
	})
	if !ContainSubset(result.Data.(map[string]interface{}), expectedDataSubSet) {
		t.Fatalf("unexpected, result does not contain subset of expected data")
	}
}

func TestIntrospection_SupportsThe__TypeRootField(t *testing.T) {

	testType := NewObject(ObjectConfig{
		Name: "TestType",
		Fields: FieldConfigMap{
			"testField": &FieldConfig{
				Type: String,
			},
		},
	})
	schema, err := NewSchema(SchemaConfig{
		Query: testType,
	})
	if err != nil {
		t.Fatalf("Error creating Schema: %v", err.Error())
	}
	query := `
      {
        __type(name: "TestType") {
          name
        }
      }
    `
	expected := &Result{
		Data: map[string]interface{}{
			"__type": map[string]interface{}{
				"name": "TestType",
			},
		},
	}
	result := g(t, Params{
		Schema:        schema,
		RequestString: query,
	})
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestIntrospection_IdentifiesDeprecatedFields(t *testing.T) {

	testType := NewObject(ObjectConfig{
		Name: "TestType",
		Fields: FieldConfigMap{
			"nonDeprecated": &FieldConfig{
				Type: String,
			},
			"deprecated": &FieldConfig{
				Type:              String,
				DeprecationReason: "Removed in 1.0",
			},
		},
	})
	schema, err := NewSchema(SchemaConfig{
		Query: testType,
	})
	if err != nil {
		t.Fatalf("Error creating Schema: %v", err.Error())
	}
	query := `
      {
        __type(name: "TestType") {
          name
          fields(includeDeprecated: true) {
            name
            isDeprecated,
            deprecationReason
          }
        }
      }
    `
	expected := &Result{
		Data: map[string]interface{}{
			"__type": map[string]interface{}{
				"name": "TestType",
				"fields": []interface{}{
					map[string]interface{}{
						"name":              "nonDeprecated",
						"isDeprecated":      false,
						"deprecationReason": nil,
					},
					map[string]interface{}{
						"name":              "deprecated",
						"isDeprecated":      true,
						"deprecationReason": "Removed in 1.0",
					},
				},
			},
		},
	}
	result := g(t, Params{
		Schema:        schema,
		RequestString: query,
	})
	if !ContainSubset(result.Data.(map[string]interface{}), expected.Data.(map[string]interface{})) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestIntrospection_RespectsTheIncludeDeprecatedParameterForFields(t *testing.T) {

	testType := NewObject(ObjectConfig{
		Name: "TestType",
		Fields: FieldConfigMap{
			"nonDeprecated": &FieldConfig{
				Type: String,
			},
			"deprecated": &FieldConfig{
				Type:              String,
				DeprecationReason: "Removed in 1.0",
			},
		},
	})
	schema, err := NewSchema(SchemaConfig{
		Query: testType,
	})
	if err != nil {
		t.Fatalf("Error creating Schema: %v", err.Error())
	}
	query := `
      {
        __type(name: "TestType") {
          name
          trueFields: fields(includeDeprecated: true) {
            name
          }
          falseFields: fields(includeDeprecated: false) {
            name
          }
          omittedFields: fields {
            name
          }
        }
      }
    `
	expected := &Result{
		Data: map[string]interface{}{
			"__type": map[string]interface{}{
				"name": "TestType",
				"trueFields": []interface{}{
					map[string]interface{}{
						"name": "nonDeprecated",
					},
					map[string]interface{}{
						"name": "deprecated",
					},
				},
				"falseFields": []interface{}{
					map[string]interface{}{
						"name": "nonDeprecated",
					},
				},
				"omittedFields": []interface{}{
					map[string]interface{}{
						"name": "nonDeprecated",
					},
				},
			},
		},
	}
	result := g(t, Params{
		Schema:        schema,
		RequestString: query,
	})
	if !ContainSubset(result.Data.(map[string]interface{}), expected.Data.(map[string]interface{})) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestIntrospection_IdentifiesDeprecatedEnumValues(t *testing.T) {

	testEnum := NewEnum(EnumConfig{
		Name: "TestEnum",
		Values: EnumValueConfigMap{
			"NONDEPRECATED": &EnumValueConfig{
				Value: 0,
			},
			"DEPRECATED": &EnumValueConfig{
				Value:             1,
				DeprecationReason: "Removed in 1.0",
			},
			"ALSONONDEPRECATED": &EnumValueConfig{
				Value: 2,
			},
		},
	})
	testType := NewObject(ObjectConfig{
		Name: "TestType",
		Fields: FieldConfigMap{
			"testEnum": &FieldConfig{
				Type: testEnum,
			},
		},
	})
	schema, err := NewSchema(SchemaConfig{
		Query: testType,
	})
	if err != nil {
		t.Fatalf("Error creating Schema: %v", err.Error())
	}
	query := `
      {
        __type(name: "TestEnum") {
          name
          enumValues(includeDeprecated: true) {
            name
            isDeprecated,
            deprecationReason
          }
        }
      }
    `
	expected := &Result{
		Data: map[string]interface{}{
			"__type": map[string]interface{}{
				"name": "TestEnum",
				"enumValues": []interface{}{
					map[string]interface{}{
						"name":              "NONDEPRECATED",
						"isDeprecated":      false,
						"deprecationReason": nil,
					},
					map[string]interface{}{
						"name":              "DEPRECATED",
						"isDeprecated":      true,
						"deprecationReason": "Removed in 1.0",
					},
					map[string]interface{}{
						"name":              "ALSONONDEPRECATED",
						"isDeprecated":      false,
						"deprecationReason": nil,
					},
				},
			},
		},
	}
	result := g(t, Params{
		Schema:        schema,
		RequestString: query,
	})
	if !ContainSubset(result.Data.(map[string]interface{}), expected.Data.(map[string]interface{})) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestIntrospection_RespectsTheIncludeDeprecatedParameterForEnumValues(t *testing.T) {

	testEnum := NewEnum(EnumConfig{
		Name: "TestEnum",
		Values: EnumValueConfigMap{
			"NONDEPRECATED": &EnumValueConfig{
				Value: 0,
			},
			"DEPRECATED": &EnumValueConfig{
				Value:             1,
				DeprecationReason: "Removed in 1.0",
			},
			"ALSONONDEPRECATED": &EnumValueConfig{
				Value: 2,
			},
		},
	})
	testType := NewObject(ObjectConfig{
		Name: "TestType",
		Fields: FieldConfigMap{
			"testEnum": &FieldConfig{
				Type: testEnum,
			},
		},
	})
	schema, err := NewSchema(SchemaConfig{
		Query: testType,
	})
	if err != nil {
		t.Fatalf("Error creating Schema: %v", err.Error())
	}
	query := `
      {
        __type(name: "TestEnum") {
          name
          trueValues: enumValues(includeDeprecated: true) {
            name
          }
          falseValues: enumValues(includeDeprecated: false) {
            name
          }
          omittedValues: enumValues {
            name
          }
        }
      }
    `
	expected := &Result{
		Data: map[string]interface{}{
			"__type": map[string]interface{}{
				"name": "TestEnum",
				"trueValues": []interface{}{
					map[string]interface{}{
						"name": "NONDEPRECATED",
					},
					map[string]interface{}{
						"name": "DEPRECATED",
					},
					map[string]interface{}{
						"name": "ALSONONDEPRECATED",
					},
				},
				"falseValues": []interface{}{
					map[string]interface{}{
						"name": "NONDEPRECATED",
					},
					map[string]interface{}{
						"name": "ALSONONDEPRECATED",
					},
				},
				"omittedValues": []interface{}{
					map[string]interface{}{
						"name": "NONDEPRECATED",
					},
					map[string]interface{}{
						"name": "ALSONONDEPRECATED",
					},
				},
			},
		},
	}
	result := g(t, Params{
		Schema:        schema,
		RequestString: query,
	})
	if !ContainSubset(result.Data.(map[string]interface{}), expected.Data.(map[string]interface{})) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestIntrospection_FailsAsExpectedOnThe__TypeRootFieldWithoutAnArg(t *testing.T) {

	testType := NewObject(ObjectConfig{
		Name: "TestType",
		Fields: FieldConfigMap{
			"testField": &FieldConfig{
				Type: String,
			},
		},
	})
	schema, err := NewSchema(SchemaConfig{
		Query: testType,
	})
	if err != nil {
		t.Fatalf("Error creating Schema: %v", err.Error())
	}
	query := `
      {
        __type {
          name
        }
      }
    `
	expected := &Result{
		Errors: []FormattedError{
			FormattedError{
				Message: `Field "__type" argument "name" of type "String!" ` +
					`is required but not provided.`,
				Locations: []SourceLocation{
					SourceLocation{Line: 3, Column: 9},
				},
			},
		},
	}
	result := g(t, Params{
		Schema:        schema,
		RequestString: query,
	})
	t.Skipf("Pending `validator` implementation")
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}

func TestIntrospection_ExposesDescriptionsOnTypesAndFields(t *testing.T) {

	queryRoot := NewObject(ObjectConfig{
		Name: "QueryRoot",
		Fields: FieldConfigMap{
			"onlyField": &FieldConfig{
				Type: String,
			},
		},
	})
	schema, err := NewSchema(SchemaConfig{
		Query: queryRoot,
	})
	if err != nil {
		t.Fatalf("Error creating Schema: %v", err.Error())
	}
	query := `
      {
        schemaType: __type(name: "__Schema") {
          name,
          description,
          fields {
            name,
            description
          }
        }
      }
    `
	expected := &Result{
		Data: map[string]interface{}{
			"schemaType": map[string]interface{}{
				"name": "__Schema",
				"description": `A GraphQL Schema defines the capabilities of a GraphQL
server. It exposes all available types and directives on
the server, as well as the entry points for query and
mutation operations.`,
				"fields": []interface{}{
					map[string]interface{}{
						"name":        "types",
						"description": "A list of all types supported by this server.",
					},
					map[string]interface{}{
						"name":        "queryType",
						"description": "The type that query operations will be rooted at.",
					},
					map[string]interface{}{
						"name": "mutationType",
						"description": "If this server supports mutation, the type that " +
							"mutation operations will be rooted at.",
					},
					map[string]interface{}{
						"name":        "directives",
						"description": "A list of all directives supported by this server.",
					},
				},
			},
		},
	}
	result := g(t, Params{
		Schema:        schema,
		RequestString: query,
	})
	if !ContainSubset(result.Data.(map[string]interface{}), expected.Data.(map[string]interface{})) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
func TestIntrospection_ExposesDescriptionsOnEnums(t *testing.T) {

	queryRoot := NewObject(ObjectConfig{
		Name: "QueryRoot",
		Fields: FieldConfigMap{
			"onlyField": &FieldConfig{
				Type: String,
			},
		},
	})
	schema, err := NewSchema(SchemaConfig{
		Query: queryRoot,
	})
	if err != nil {
		t.Fatalf("Error creating Schema: %v", err.Error())
	}
	query := `
      {
        typeKindType: __type(name: "__TypeKind") {
          name,
          description,
          enumValues {
            name,
            description
          }
        }
      }
    `
	expected := &Result{
		Data: map[string]interface{}{
			"typeKindType": map[string]interface{}{
				"name":        "__TypeKind",
				"description": `An enum describing what kind of type a given __Type is`,
				"enumValues": []interface{}{
					map[string]interface{}{
						"name":        "SCALAR",
						"description": "Indicates this type is a scalar.",
					},
					map[string]interface{}{
						"name":        "OBJECT",
						"description": "Indicates this type is an object. `fields` and `interfaces` are valid fields.",
					},
					map[string]interface{}{
						"name":        "INTERFACE",
						"description": "Indicates this type is an interface. `fields` and `possibleTypes` are valid fields.",
					},
					map[string]interface{}{
						"name":        "UNION",
						"description": "Indicates this type is a union. `possibleTypes` is a valid field.",
					},
					map[string]interface{}{
						"name":        "ENUM",
						"description": "Indicates this type is an enum. `enumValues` is a valid field.",
					},
					map[string]interface{}{
						"name":        "INPUT_OBJECT",
						"description": "Indicates this type is an input object. `inputFields` is a valid field.",
					},
					map[string]interface{}{
						"name":        "LIST",
						"description": "Indicates this type is a list. `ofType` is a valid field.",
					},
					map[string]interface{}{
						"name":        "NON_NULL",
						"description": "Indicates this type is a non-null. `ofType` is a valid field.",
					},
				},
			},
		},
	}
	result := g(t, Params{
		Schema:        schema,
		RequestString: query,
	})
	if !ContainSubset(result.Data.(map[string]interface{}), expected.Data.(map[string]interface{})) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expected, result))
	}
}
