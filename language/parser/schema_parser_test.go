package parser_test

import (
	"reflect"
	"testing"

	"github.com/chris-ramon/graphql-go/errors"
	"github.com/chris-ramon/graphql-go/language/ast"
	"github.com/chris-ramon/graphql-go/language/location"
	"github.com/chris-ramon/graphql-go/language/parser"
	"github.com/chris-ramon/graphql-go/language/source"
)

func parse(t *testing.T, query string) *ast.Document {
	astDoc, err := parser.Parse(parser.ParseParams{
		Source: query,
		Options: parser.ParseOptions{
			NoLocation: false,
			NoSource:   true,
		},
	})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	return astDoc
}
func loc(start int, end int) *ast.Location {
	return &ast.Location{
		Start: start, End: end,
	}
}
func TestSchemaParser_SimpleType(t *testing.T) {

	body := `
type Hello {
  world: String
}`
	astDoc := parse(t, body)
	expected := ast.NewDocument(&ast.Document{
		Loc: loc(1, 31),
		Definitions: []ast.Node{
			ast.NewObjectTypeDefinition(&ast.ObjectTypeDefinition{
				Loc: loc(1, 31),
				Name: ast.NewName(&ast.Name{
					Value: "Hello",
					Loc:   loc(6, 11),
				}),
				Interfaces: []*ast.NamedType{},
				Fields: []*ast.FieldDefinition{
					ast.NewFieldDefinition(&ast.FieldDefinition{
						Loc: loc(16, 29),
						Name: ast.NewName(&ast.Name{
							Value: "world",
							Loc:   loc(16, 21),
						}),
						Arguments: []*ast.InputValueDefinition{},
						Type: ast.NewNamedType(&ast.NamedType{
							Loc: loc(23, 29),
							Name: ast.NewName(&ast.Name{
								Value: "String",
								Loc:   loc(23, 29),
							}),
						}),
					}),
				},
			}),
		},
	})
	if !reflect.DeepEqual(astDoc, expected) {
		t.Fatalf("unexpected document, expected: %v, got: %v", expected, astDoc)
	}
}

func TestSchemaParser_SimpleExtension(t *testing.T) {

	body := `
extend type Hello {
  world: String
}`
	astDoc := parse(t, body)
	expected := ast.NewDocument(&ast.Document{
		Loc: loc(1, 38),
		Definitions: []ast.Node{
			ast.NewTypeExtensionDefinition(&ast.TypeExtensionDefinition{
				Loc: loc(1, 38),
				Definition: ast.NewObjectTypeDefinition(&ast.ObjectTypeDefinition{
					Loc: loc(8, 38),
					Name: ast.NewName(&ast.Name{
						Value: "Hello",
						Loc:   loc(13, 18),
					}),
					Interfaces: []*ast.NamedType{},
					Fields: []*ast.FieldDefinition{
						ast.NewFieldDefinition(&ast.FieldDefinition{
							Loc: loc(23, 36),
							Name: ast.NewName(&ast.Name{
								Value: "world",
								Loc:   loc(23, 28),
							}),
							Arguments: []*ast.InputValueDefinition{},
							Type: ast.NewNamedType(&ast.NamedType{
								Loc: loc(30, 36),
								Name: ast.NewName(&ast.Name{
									Value: "String",
									Loc:   loc(30, 36),
								}),
							}),
						}),
					},
				}),
			}),
		},
	})
	if !reflect.DeepEqual(astDoc, expected) {
		t.Fatalf("unexpected document, expected: %v, got: %v", expected, astDoc)
	}
}

func TestSchemaParser_SimpleNonNullType(t *testing.T) {

	body := `
type Hello {
  world: String!
}`
	astDoc := parse(t, body)
	expected := ast.NewDocument(&ast.Document{
		Loc: loc(1, 32),
		Definitions: []ast.Node{
			ast.NewObjectTypeDefinition(&ast.ObjectTypeDefinition{
				Loc: loc(1, 32),
				Name: ast.NewName(&ast.Name{
					Value: "Hello",
					Loc:   loc(6, 11),
				}),
				Interfaces: []*ast.NamedType{},
				Fields: []*ast.FieldDefinition{
					ast.NewFieldDefinition(&ast.FieldDefinition{
						Loc: loc(16, 30),
						Name: ast.NewName(&ast.Name{
							Value: "world",
							Loc:   loc(16, 21),
						}),
						Arguments: []*ast.InputValueDefinition{},
						Type: ast.NewNonNullType(&ast.NonNullType{
							Kind: "NonNullType",
							Loc:  loc(23, 30),
							Type: ast.NewNamedType(&ast.NamedType{
								Loc: loc(23, 29),
								Name: ast.NewName(&ast.Name{
									Value: "String",
									Loc:   loc(23, 29),
								}),
							}),
						}),
					}),
				},
			}),
		},
	})
	if !reflect.DeepEqual(astDoc, expected) {
		t.Fatalf("unexpected document, expected: %v, got: %v", expected, astDoc)
	}
}

func TestSchemaParser_SimpleTypeInheritingInterface(t *testing.T) {
	body := `type Hello implements World { }`
	astDoc := parse(t, body)
	expected := ast.NewDocument(&ast.Document{
		Loc: loc(0, 31),
		Definitions: []ast.Node{
			ast.NewObjectTypeDefinition(&ast.ObjectTypeDefinition{
				Loc: loc(0, 31),
				Name: ast.NewName(&ast.Name{
					Value: "Hello",
					Loc:   loc(5, 10),
				}),
				Interfaces: []*ast.NamedType{
					ast.NewNamedType(&ast.NamedType{
						Name: ast.NewName(&ast.Name{
							Value: "World",
							Loc:   loc(22, 27),
						}),
						Loc: loc(22, 27),
					}),
				},
				Fields: []*ast.FieldDefinition{},
			}),
		},
	})
	if !reflect.DeepEqual(astDoc, expected) {
		t.Fatalf("unexpected document, expected: %v, got: %v", expected, astDoc)
	}
}

func TestSchemaParser_SimpleTypeInheritingMultipleInterfaces(t *testing.T) {
	body := `type Hello implements Wo, rld { }`
	astDoc := parse(t, body)
	expected := ast.NewDocument(&ast.Document{
		Loc: loc(0, 33),
		Definitions: []ast.Node{
			ast.NewObjectTypeDefinition(&ast.ObjectTypeDefinition{
				Loc: loc(0, 33),
				Name: ast.NewName(&ast.Name{
					Value: "Hello",
					Loc:   loc(5, 10),
				}),
				Interfaces: []*ast.NamedType{
					ast.NewNamedType(&ast.NamedType{
						Name: ast.NewName(&ast.Name{
							Value: "Wo",
							Loc:   loc(22, 24),
						}),
						Loc: loc(22, 24),
					}),
					ast.NewNamedType(&ast.NamedType{
						Name: ast.NewName(&ast.Name{
							Value: "rld",
							Loc:   loc(26, 29),
						}),
						Loc: loc(26, 29),
					}),
				},
				Fields: []*ast.FieldDefinition{},
			}),
		},
	})
	if !reflect.DeepEqual(astDoc, expected) {
		t.Fatalf("unexpected document, expected: %v, got: %v", expected, astDoc)
	}
}

func TestSchemaParser_SingleValueEnum(t *testing.T) {
	body := `enum Hello { WORLD }`
	astDoc := parse(t, body)
	expected := ast.NewDocument(&ast.Document{
		Loc: loc(0, 20),
		Definitions: []ast.Node{
			ast.NewEnumTypeDefinition(&ast.EnumTypeDefinition{
				Loc: loc(0, 20),
				Name: ast.NewName(&ast.Name{
					Value: "Hello",
					Loc:   loc(5, 10),
				}),
				Values: []*ast.EnumValueDefinition{
					ast.NewEnumValueDefinition(&ast.EnumValueDefinition{
						Name: ast.NewName(&ast.Name{
							Value: "WORLD",
							Loc:   loc(13, 18),
						}),
						Loc: loc(13, 18),
					}),
				},
			}),
		},
	})
	if !reflect.DeepEqual(astDoc, expected) {
		t.Fatalf("unexpected document, expected: %v, got: %v", expected, astDoc)
	}
}

func TestSchemaParser_DoubleValueEnum(t *testing.T) {
	body := `enum Hello { WO, RLD }`
	astDoc := parse(t, body)
	expected := ast.NewDocument(&ast.Document{
		Loc: loc(0, 22),
		Definitions: []ast.Node{
			ast.NewEnumTypeDefinition(&ast.EnumTypeDefinition{
				Loc: loc(0, 22),
				Name: ast.NewName(&ast.Name{
					Value: "Hello",
					Loc:   loc(5, 10),
				}),
				Values: []*ast.EnumValueDefinition{
					ast.NewEnumValueDefinition(&ast.EnumValueDefinition{
						Name: ast.NewName(&ast.Name{
							Value: "WO",
							Loc:   loc(13, 15),
						}),
						Loc: loc(13, 15),
					}),
					ast.NewEnumValueDefinition(&ast.EnumValueDefinition{
						Name: ast.NewName(&ast.Name{
							Value: "RLD",
							Loc:   loc(17, 20),
						}),
						Loc: loc(17, 20),
					}),
				},
			}),
		},
	})
	if !reflect.DeepEqual(astDoc, expected) {
		t.Fatalf("unexpected document, expected: %v, got: %v", expected, astDoc)
	}
}

func TestSchemaParser_SimpleInterface(t *testing.T) {
	body := `
interface Hello {
  world: String
}`
	astDoc := parse(t, body)
	expected := ast.NewDocument(&ast.Document{
		Loc: loc(1, 36),
		Definitions: []ast.Node{
			ast.NewInterfaceTypeDefinition(&ast.InterfaceTypeDefinition{
				Loc: loc(1, 36),
				Name: ast.NewName(&ast.Name{
					Value: "Hello",
					Loc:   loc(11, 16),
				}),
				Fields: []*ast.FieldDefinition{
					ast.NewFieldDefinition(&ast.FieldDefinition{
						Loc: loc(21, 34),
						Name: ast.NewName(&ast.Name{
							Value: "world",
							Loc:   loc(21, 26),
						}),
						Arguments: []*ast.InputValueDefinition{},
						Type: ast.NewNamedType(&ast.NamedType{
							Loc: loc(28, 34),
							Name: ast.NewName(&ast.Name{
								Value: "String",
								Loc:   loc(28, 34),
							}),
						}),
					}),
				},
			}),
		},
	})
	if !reflect.DeepEqual(astDoc, expected) {
		t.Fatalf("unexpected document, expected: %v, got: %v", expected, astDoc)
	}
}

func TestSchemaParser_SimpleFieldWithArg(t *testing.T) {
	body := `
type Hello {
  world(flag: Boolean): String
}`
	astDoc := parse(t, body)
	expected := ast.NewDocument(&ast.Document{
		Loc: loc(1, 46),
		Definitions: []ast.Node{
			ast.NewObjectTypeDefinition(&ast.ObjectTypeDefinition{
				Loc: loc(1, 46),
				Name: ast.NewName(&ast.Name{
					Value: "Hello",
					Loc:   loc(6, 11),
				}),
				Interfaces: []*ast.NamedType{},
				Fields: []*ast.FieldDefinition{
					ast.NewFieldDefinition(&ast.FieldDefinition{
						Loc: loc(16, 44),
						Name: ast.NewName(&ast.Name{
							Value: "world",
							Loc:   loc(16, 21),
						}),
						Arguments: []*ast.InputValueDefinition{
							ast.NewInputValueDefinition(&ast.InputValueDefinition{
								Loc: loc(22, 35),
								Name: ast.NewName(&ast.Name{
									Value: "flag",
									Loc:   loc(22, 26),
								}),
								Type: ast.NewNamedType(&ast.NamedType{
									Loc: loc(28, 35),
									Name: ast.NewName(&ast.Name{
										Value: "Boolean",
										Loc:   loc(28, 35),
									}),
								}),
								DefaultValue: nil,
							}),
						},
						Type: ast.NewNamedType(&ast.NamedType{
							Loc: loc(38, 44),
							Name: ast.NewName(&ast.Name{
								Value: "String",
								Loc:   loc(38, 44),
							}),
						}),
					}),
				},
			}),
		},
	})
	if !reflect.DeepEqual(astDoc, expected) {
		t.Fatalf("unexpected document, expected: %v, got: %v", expected, astDoc)
	}
}

func TestSchemaParser_SimpleFieldWithArgWithDefaultValue(t *testing.T) {
	body := `
type Hello {
  world(flag: Boolean = true): String
}`
	astDoc := parse(t, body)
	expected := ast.NewDocument(&ast.Document{
		Loc: loc(1, 53),
		Definitions: []ast.Node{
			ast.NewObjectTypeDefinition(&ast.ObjectTypeDefinition{
				Loc: loc(1, 53),
				Name: ast.NewName(&ast.Name{
					Value: "Hello",
					Loc:   loc(6, 11),
				}),
				Interfaces: []*ast.NamedType{},
				Fields: []*ast.FieldDefinition{
					ast.NewFieldDefinition(&ast.FieldDefinition{
						Loc: loc(16, 51),
						Name: ast.NewName(&ast.Name{
							Value: "world",
							Loc:   loc(16, 21),
						}),
						Arguments: []*ast.InputValueDefinition{
							ast.NewInputValueDefinition(&ast.InputValueDefinition{
								Loc: loc(22, 42),
								Name: ast.NewName(&ast.Name{
									Value: "flag",
									Loc:   loc(22, 26),
								}),
								Type: ast.NewNamedType(&ast.NamedType{
									Loc: loc(28, 35),
									Name: ast.NewName(&ast.Name{
										Value: "Boolean",
										Loc:   loc(28, 35),
									}),
								}),
								DefaultValue: ast.NewBooleanValue(&ast.BooleanValue{
									Value: true,
									Loc:   loc(38, 42),
								}),
							}),
						},
						Type: ast.NewNamedType(&ast.NamedType{
							Loc: loc(45, 51),
							Name: ast.NewName(&ast.Name{
								Value: "String",
								Loc:   loc(45, 51),
							}),
						}),
					}),
				},
			}),
		},
	})
	if !reflect.DeepEqual(astDoc, expected) {
		t.Fatalf("unexpected document, expected: %v, got: %v", expected, astDoc)
	}
}

func TestSchemaParser_SimpleFieldWithListArg(t *testing.T) {
	body := `
type Hello {
  world(things: [String]): String
}`
	astDoc := parse(t, body)
	expected := ast.NewDocument(&ast.Document{
		Loc: loc(1, 49),
		Definitions: []ast.Node{
			ast.NewObjectTypeDefinition(&ast.ObjectTypeDefinition{
				Loc: loc(1, 49),
				Name: ast.NewName(&ast.Name{
					Value: "Hello",
					Loc:   loc(6, 11),
				}),
				Interfaces: []*ast.NamedType{},
				Fields: []*ast.FieldDefinition{
					ast.NewFieldDefinition(&ast.FieldDefinition{
						Loc: loc(16, 47),
						Name: ast.NewName(&ast.Name{
							Value: "world",
							Loc:   loc(16, 21),
						}),
						Arguments: []*ast.InputValueDefinition{
							ast.NewInputValueDefinition(&ast.InputValueDefinition{
								Loc: loc(22, 38),
								Name: ast.NewName(&ast.Name{
									Value: "things",
									Loc:   loc(22, 28),
								}),
								Type: ast.NewListType(&ast.ListType{
									Loc: loc(30, 38),
									Type: ast.NewNamedType(&ast.NamedType{
										Loc: loc(31, 37),
										Name: ast.NewName(&ast.Name{
											Value: "String",
											Loc:   loc(31, 37),
										}),
									}),
								}),
								DefaultValue: nil,
							}),
						},
						Type: ast.NewNamedType(&ast.NamedType{
							Loc: loc(41, 47),
							Name: ast.NewName(&ast.Name{
								Value: "String",
								Loc:   loc(41, 47),
							}),
						}),
					}),
				},
			}),
		},
	})
	if !reflect.DeepEqual(astDoc, expected) {
		t.Fatalf("unexpected document, expected: %v, got: %v", expected, astDoc)
	}
}

func TestSchemaParser_SimpleFieldWithTwoArg(t *testing.T) {
	body := `
type Hello {
  world(argOne: Boolean, argTwo: Int): String
}`
	astDoc := parse(t, body)
	expected := ast.NewDocument(&ast.Document{
		Loc: loc(1, 61),
		Definitions: []ast.Node{
			ast.NewObjectTypeDefinition(&ast.ObjectTypeDefinition{
				Loc: loc(1, 61),
				Name: ast.NewName(&ast.Name{
					Value: "Hello",
					Loc:   loc(6, 11),
				}),
				Interfaces: []*ast.NamedType{},
				Fields: []*ast.FieldDefinition{
					ast.NewFieldDefinition(&ast.FieldDefinition{
						Loc: loc(16, 59),
						Name: ast.NewName(&ast.Name{
							Value: "world",
							Loc:   loc(16, 21),
						}),
						Arguments: []*ast.InputValueDefinition{
							ast.NewInputValueDefinition(&ast.InputValueDefinition{
								Loc: loc(22, 37),
								Name: ast.NewName(&ast.Name{
									Value: "argOne",
									Loc:   loc(22, 28),
								}),
								Type: ast.NewNamedType(&ast.NamedType{
									Loc: loc(30, 37),
									Name: ast.NewName(&ast.Name{
										Value: "Boolean",
										Loc:   loc(30, 37),
									}),
								}),
								DefaultValue: nil,
							}),
							ast.NewInputValueDefinition(&ast.InputValueDefinition{
								Loc: loc(39, 50),
								Name: ast.NewName(&ast.Name{
									Value: "argTwo",
									Loc:   loc(39, 45),
								}),
								Type: ast.NewNamedType(&ast.NamedType{
									Loc: loc(47, 50),
									Name: ast.NewName(&ast.Name{
										Value: "Int",
										Loc:   loc(47, 50),
									}),
								}),
								DefaultValue: nil,
							}),
						},
						Type: ast.NewNamedType(&ast.NamedType{
							Loc: loc(53, 59),
							Name: ast.NewName(&ast.Name{
								Value: "String",
								Loc:   loc(53, 59),
							}),
						}),
					}),
				},
			}),
		},
	})
	if !reflect.DeepEqual(astDoc, expected) {
		t.Fatalf("unexpected document, expected: %v, got: %v", expected, astDoc)
	}
}

func TestSchemaParser_SimpleUnion(t *testing.T) {
	body := `union Hello = World`
	astDoc := parse(t, body)
	expected := ast.NewDocument(&ast.Document{
		Loc: loc(0, 19),
		Definitions: []ast.Node{
			ast.NewUnionTypeDefinition(&ast.UnionTypeDefinition{
				Loc: loc(0, 19),
				Name: ast.NewName(&ast.Name{
					Value: "Hello",
					Loc:   loc(6, 11),
				}),
				Types: []*ast.NamedType{
					ast.NewNamedType(&ast.NamedType{
						Loc: loc(14, 19),
						Name: ast.NewName(&ast.Name{
							Value: "World",
							Loc:   loc(14, 19),
						}),
					}),
				},
			}),
		},
	})
	if !reflect.DeepEqual(astDoc, expected) {
		t.Fatalf("unexpected document, expected: %v, got: %v", expected, astDoc)
	}
}

func TestSchemaParser_UnionWithTwoTypes(t *testing.T) {
	body := `union Hello = Wo | Rld`
	astDoc := parse(t, body)
	expected := ast.NewDocument(&ast.Document{
		Loc: loc(0, 22),
		Definitions: []ast.Node{
			ast.NewUnionTypeDefinition(&ast.UnionTypeDefinition{
				Loc: loc(0, 22),
				Name: ast.NewName(&ast.Name{
					Value: "Hello",
					Loc:   loc(6, 11),
				}),
				Types: []*ast.NamedType{
					ast.NewNamedType(&ast.NamedType{
						Loc: loc(14, 16),
						Name: ast.NewName(&ast.Name{
							Value: "Wo",
							Loc:   loc(14, 16),
						}),
					}),
					ast.NewNamedType(&ast.NamedType{
						Loc: loc(19, 22),
						Name: ast.NewName(&ast.Name{
							Value: "Rld",
							Loc:   loc(19, 22),
						}),
					}),
				},
			}),
		},
	})
	if !reflect.DeepEqual(astDoc, expected) {
		t.Fatalf("unexpected document, expected: %v, got: %v", expected, astDoc)
	}
}

func TestSchemaParser_Scalar(t *testing.T) {
	body := `scalar Hello`
	astDoc := parse(t, body)
	expected := ast.NewDocument(&ast.Document{
		Loc: loc(0, 12),
		Definitions: []ast.Node{
			ast.NewScalarTypeDefinition(&ast.ScalarTypeDefinition{
				Loc: loc(0, 12),
				Name: ast.NewName(&ast.Name{
					Value: "Hello",
					Loc:   loc(7, 12),
				}),
			}),
		},
	})
	if !reflect.DeepEqual(astDoc, expected) {
		t.Fatalf("unexpected document, expected: %v, got: %v", expected, astDoc)
	}
}

func TestSchemaParser_SimpleInputObject(t *testing.T) {
	body := `
input Hello {
  world: String
}`
	astDoc := parse(t, body)
	expected := ast.NewDocument(&ast.Document{
		Loc: loc(1, 32),
		Definitions: []ast.Node{
			ast.NewInputObjectTypeDefinition(&ast.InputObjectTypeDefinition{
				Loc: loc(1, 32),
				Name: ast.NewName(&ast.Name{
					Value: "Hello",
					Loc:   loc(7, 12),
				}),
				Fields: []*ast.InputValueDefinition{
					ast.NewInputValueDefinition(&ast.InputValueDefinition{
						Loc: loc(17, 30),
						Name: ast.NewName(&ast.Name{
							Value: "world",
							Loc:   loc(17, 22),
						}),
						Type: ast.NewNamedType(&ast.NamedType{
							Loc: loc(24, 30),
							Name: ast.NewName(&ast.Name{
								Value: "String",
								Loc:   loc(24, 30),
							}),
						}),
						DefaultValue: nil,
					}),
				},
			}),
		},
	})
	if !reflect.DeepEqual(astDoc, expected) {
		t.Fatalf("unexpected document, expected: %v, got: %v", expected, astDoc)
	}
}

func TestSchemaParser_SimpleInputObjectWithArgsShouldFail(t *testing.T) {
	body := `
input Hello {
  world(foo: Int): String
}`

	_, err := parser.Parse(parser.ParseParams{
		Source: body,
		Options: parser.ParseOptions{
			NoLocation: false,
			NoSource:   true,
		},
	})

	expectedError := &graphqlerrors.GraphQLError{
		Message: `Syntax Error GraphQL (3:8) Expected :, found (

2: input Hello {
3:   world(foo: Int): String
          ^
4: }
`,
		Stack: `Syntax Error GraphQL (3:8) Expected :, found (

2: input Hello {
3:   world(foo: Int): String
          ^
4: }
`,
		Nodes: []ast.Node{},
		Source: &source.Source{
			Body: `
input Hello {
  world(foo: Int): String
}`,
			Name: "GraphQL",
		},
		Positions: []int{22},
		Locations: []location.SourceLocation{
			{Line: 3, Column: 8},
		},
	}
	if err == nil {
		t.Fatalf("expected error, expected: %v, got: %v", expectedError, nil)
	}
	if !reflect.DeepEqual(expectedError, err) {
		t.Fatalf("unexpected document, expected: %v, got: %v", expectedError, err)
	}
}
