package parser

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/chris-ramon/graphql-go/language/ast"
	"github.com/chris-ramon/graphql-go/language/source"
)

func printJSON(doc *ast.Document) (string, error) {
	b, err := json.Marshal(doc)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

type locFn func(start, end int) *ast.Location

func createLoc(body string) locFn {
	return func(start, end int) *ast.Location {
		return &ast.Location{
			Start:  start,
			End:    end,
			Source: source.NewSource(&source.Source{Body: body}),
		}
	}
}

func TestSimpleType(t *testing.T) {
	body := `
type Hello {
  world: String
}`
	params := ParseParams{
		Source: source.NewSource(&source.Source{Body: body}),
	}
	doc, err := Parse(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	docJSON, err := printJSON(doc)
	if err != nil {
		t.Fatalf("unexpected error, error: %v", err)
	}
	locFn := createLoc(body)
	oTypeDef := ast.NewObjectTypeDefinition(&ast.ObjectTypeDefinition{
		Name: ast.NewName(&ast.Name{
			Value: "Hello",
			Loc:   locFn(6, 11),
		}),
		Interfaces: []*ast.NamedType{},
		Fields: []*ast.FieldDefinition{
			ast.NewFieldDefinition(&ast.FieldDefinition{
				Loc: locFn(16, 29),
				Name: ast.NewName(&ast.Name{
					Loc:   locFn(16, 21),
					Value: "world",
				}),
				Arguments: []*ast.InputValueDefinition{},
				Type: ast.NewNamedType(&ast.NamedType{
					Name: ast.NewName(&ast.Name{
						Loc:   locFn(23, 29),
						Value: "String",
					}),
					Loc: locFn(23, 29),
				}),
			}),
		},
		Loc: locFn(1, 31),
	})
	expectedDocument := ast.NewDocument(&ast.Document{
		Loc:         locFn(1, 31),
		Definitions: []ast.Node{oTypeDef},
	})
	expectedJSON, err := printJSON(expectedDocument)
	if err != nil {
		t.Fatalf("unexpected error, error: %v", err)
	}
	if !reflect.DeepEqual(docJSON, expectedJSON) {
		t.Fatalf("unexpected document, \n\n expected: \n\n%v, \n\ngot: \n\n%v", expectedJSON, docJSON)
	}
}

func TestSimpleExtension(t *testing.T) {
	body := `
extend type Hello {
  world: String
}`
	params := ParseParams{
		Source: source.NewSource(&source.Source{Body: body}),
	}
	doc, err := Parse(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	docJSON, err := printJSON(doc)
	if err != nil {
		t.Fatalf("unexpected error, error: %v", err)
	}
	locFn := createLoc(body)
	typeExtDef := ast.NewTypeExtensionDefinition(&ast.TypeExtensionDefinition{
		Definition: ast.NewObjectTypeDefinition(&ast.ObjectTypeDefinition{
			Name: ast.NewName(&ast.Name{
				Value: "Hello",
				Loc:   locFn(13, 18),
			}),
			Interfaces: []*ast.NamedType{},
			Fields: []*ast.FieldDefinition{
				ast.NewFieldDefinition(&ast.FieldDefinition{
					Name: ast.NewName(&ast.Name{
						Loc:   locFn(23, 28),
						Value: "world",
					}),
					Arguments: []*ast.InputValueDefinition{},
					Type: ast.NewNamedType(&ast.NamedType{
						Name: ast.NewName(&ast.Name{
							Loc:   locFn(30, 36),
							Value: "String",
						}),
						Loc: locFn(30, 36),
					}),
					Loc: locFn(23, 36),
				}),
			},
			Loc: locFn(8, 38),
		}),
		Loc: locFn(1, 38),
	})
	expectedDocument := ast.NewDocument(&ast.Document{
		Loc:         locFn(1, 38),
		Definitions: []ast.Node{typeExtDef},
	})
	expectedJSON, err := printJSON(expectedDocument)
	if err != nil {
		t.Fatalf("unexpected error, error: %v", err)
	}
	if !reflect.DeepEqual(docJSON, expectedJSON) {
		t.Fatalf("unexpected document, \n\n expected: \n\n%v, \n\ngot: \n\n%v", expectedJSON, docJSON)
	}
}

func TestSimpleNonNullType(t *testing.T) {
	body := `
type Hello {
  world: String!
}`
	params := ParseParams{
		Source: source.NewSource(&source.Source{Body: body}),
	}
	doc, err := Parse(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	docJSON, err := printJSON(doc)
	if err != nil {
		t.Fatalf("unexpected error, error: %v", err)
	}
	locFn := createLoc(body)
	oTypeDef := ast.NewObjectTypeDefinition(&ast.ObjectTypeDefinition{
		Name: ast.NewName(&ast.Name{
			Value: "Hello",
			Loc:   locFn(6, 11),
		}),
		Interfaces: []*ast.NamedType{},
		Fields: []*ast.FieldDefinition{
			ast.NewFieldDefinition(&ast.FieldDefinition{
				Name: ast.NewName(&ast.Name{
					Loc:   locFn(16, 21),
					Value: "world",
				}),
				Arguments: []*ast.InputValueDefinition{},
				Type: ast.NewNonNullType(&ast.NonNullType{
					Loc: locFn(23, 30),
					Type: ast.NewNamedType(&ast.NamedType{
						Name: ast.NewName(&ast.Name{
							Loc:   locFn(23, 29),
							Value: "String",
						}),
						Loc: locFn(23, 29),
					}),
				}),
				Loc: locFn(16, 30),
			}),
		},
		Loc: locFn(1, 32),
	})
	expectedDocument := ast.NewDocument(&ast.Document{
		Loc:         locFn(1, 32),
		Definitions: []ast.Node{oTypeDef},
	})
	expectedJSON, err := printJSON(expectedDocument)
	if err != nil {
		t.Fatalf("unexpected error, error: %v", err)
	}
	if !reflect.DeepEqual(docJSON, expectedJSON) {
		t.Fatalf("unexpected document, \n\n expected: \n\n%v, \n\ngot: \n\n%v", expectedJSON, docJSON)
	}
}

func TestSimpleTypeInheritingInterface(t *testing.T) {
	body := `type Hello implements World { }`
	params := ParseParams{
		Source: source.NewSource(&source.Source{Body: body}),
	}
	doc, err := Parse(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	docJSON, err := printJSON(doc)
	if err != nil {
		t.Fatalf("unexpected error, error: %v", err)
	}
	locFn := createLoc(body)
	oTypeDef := ast.NewObjectTypeDefinition(&ast.ObjectTypeDefinition{
		Name: ast.NewName(&ast.Name{
			Value: "Hello",
			Loc:   locFn(5, 10),
		}),
		Interfaces: []*ast.NamedType{
			ast.NewNamedType(&ast.NamedType{
				Loc: locFn(22, 27),
				Name: ast.NewName(&ast.Name{
					Value: "World",
					Loc:   locFn(22, 27),
				}),
			}),
		},
		Fields: []*ast.FieldDefinition{},
		Loc:    locFn(0, 31),
	})
	expectedDocument := ast.NewDocument(&ast.Document{
		Loc:         locFn(0, 31),
		Definitions: []ast.Node{oTypeDef},
	})
	expectedJSON, err := printJSON(expectedDocument)
	if err != nil {
		t.Fatalf("unexpected error, error: %v", err)
	}
	if !reflect.DeepEqual(docJSON, expectedJSON) {
		t.Fatalf("unexpected document, \n\n expected: \n\n%v, \n\ngot: \n\n%v", expectedJSON, docJSON)
	}
}

func TestSimpleTypeInheritingMultipleInterfaces(t *testing.T) {
	body := `type Hello implements Wo, rld { }`
	params := ParseParams{
		Source: source.NewSource(&source.Source{Body: body}),
	}
	doc, err := Parse(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	docJSON, err := printJSON(doc)
	if err != nil {
		t.Fatalf("unexpected error, error: %v", err)
	}
	locFn := createLoc(body)
	oTypeDef := ast.NewObjectTypeDefinition(&ast.ObjectTypeDefinition{
		Name: ast.NewName(&ast.Name{
			Value: "Hello",
			Loc:   locFn(5, 10),
		}),
		Interfaces: []*ast.NamedType{
			ast.NewNamedType(&ast.NamedType{
				Loc: locFn(22, 24),
				Name: ast.NewName(&ast.Name{
					Value: "Wo",
					Loc:   locFn(22, 24),
				}),
			}),
			ast.NewNamedType(&ast.NamedType{
				Loc: locFn(26, 29),
				Name: ast.NewName(&ast.Name{
					Value: "rld",
					Loc:   locFn(26, 29),
				}),
			}),
		},
		Fields: []*ast.FieldDefinition{},
		Loc:    locFn(0, 33),
	})
	expectedDocument := ast.NewDocument(&ast.Document{
		Loc:         locFn(0, 33),
		Definitions: []ast.Node{oTypeDef},
	})
	expectedJSON, err := printJSON(expectedDocument)
	if err != nil {
		t.Fatalf("unexpected error, error: %v", err)
	}
	if !reflect.DeepEqual(docJSON, expectedJSON) {
		t.Fatalf("unexpected document, \n\n expected: \n\n%v, \n\ngot: \n\n%v", expectedJSON, docJSON)
	}
}

func TestSimpleValueEnum(t *testing.T) {
	body := `enum Hello { WORLD }`
	params := ParseParams{
		Source: source.NewSource(&source.Source{Body: body}),
	}
	doc, err := Parse(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	docJSON, err := printJSON(doc)
	if err != nil {
		t.Fatalf("unexpected error, error: %v", err)
	}
	locFn := createLoc(body)
	eTypeDef := ast.NewEnumTypeDefinition(&ast.EnumTypeDefinition{
		Name: ast.NewName(&ast.Name{
			Value: "Hello",
			Loc:   locFn(5, 10),
		}),
		Values: []*ast.EnumValueDefinition{
			ast.NewEnumValueDefinition(&ast.EnumValueDefinition{
				Loc: locFn(13, 18),
				Name: ast.NewName(&ast.Name{
					Value: "WORLD",
					Loc:   locFn(13, 18),
				}),
			}),
		},
		Loc: locFn(0, 20),
	})
	expectedDocument := ast.NewDocument(&ast.Document{
		Loc:         locFn(0, 20),
		Definitions: []ast.Node{eTypeDef},
	})
	expectedJSON, err := printJSON(expectedDocument)
	if err != nil {
		t.Fatalf("unexpected error, error: %v", err)
	}
	if !reflect.DeepEqual(docJSON, expectedJSON) {
		t.Fatalf("unexpected document, \n\n expected: \n\n%v, \n\ngot: \n\n%v", expectedJSON, docJSON)
	}
}

func TestDoubleValueEnum(t *testing.T) {
	body := `enum Hello { WO, RLD }`
	params := ParseParams{
		Source: source.NewSource(&source.Source{Body: body}),
	}
	doc, err := Parse(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	docJSON, err := printJSON(doc)
	if err != nil {
		t.Fatalf("unexpected error, error: %v", err)
	}
	locFn := createLoc(body)
	eTypeDef := ast.NewEnumTypeDefinition(&ast.EnumTypeDefinition{
		Name: ast.NewName(&ast.Name{
			Value: "Hello",
			Loc:   locFn(5, 10),
		}),
		Values: []*ast.EnumValueDefinition{
			ast.NewEnumValueDefinition(&ast.EnumValueDefinition{
				Loc: locFn(13, 15),
				Name: ast.NewName(&ast.Name{
					Value: "WO",
					Loc:   locFn(13, 15),
				}),
			}),
			ast.NewEnumValueDefinition(&ast.EnumValueDefinition{
				Loc: locFn(17, 20),
				Name: ast.NewName(&ast.Name{
					Value: "RLD",
					Loc:   locFn(17, 20),
				}),
			}),
		},
		Loc: locFn(0, 22),
	})
	expectedDocument := ast.NewDocument(&ast.Document{
		Loc:         locFn(0, 22),
		Definitions: []ast.Node{eTypeDef},
	})
	expectedJSON, err := printJSON(expectedDocument)
	if err != nil {
		t.Fatalf("unexpected error, error: %v", err)
	}
	if !reflect.DeepEqual(docJSON, expectedJSON) {
		t.Fatalf("unexpected document, \n\n expected: \n\n%v, \n\ngot: \n\n%v", expectedJSON, docJSON)
	}
}

func TestSimpleInterface(t *testing.T) {
	body := `
interface Hello {
  world: String
}`
	params := ParseParams{
		Source: source.NewSource(&source.Source{Body: body}),
	}
	doc, err := Parse(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	docJSON, err := printJSON(doc)
	if err != nil {
		t.Fatalf("unexpected error, error: %v", err)
	}
	locFn := createLoc(body)
	iTypeDef := ast.NewInterfaceTypeDefinition(&ast.InterfaceTypeDefinition{
		Name: ast.NewName(&ast.Name{
			Value: "Hello",
			Loc:   locFn(11, 16),
		}),
		Fields: []*ast.FieldDefinition{
			ast.NewFieldDefinition(&ast.FieldDefinition{
				Name: ast.NewName(&ast.Name{
					Loc:   locFn(21, 26),
					Value: "world",
				}),
				Arguments: []*ast.InputValueDefinition{},
				Type: ast.NewNamedType(&ast.NamedType{
					Name: ast.NewName(&ast.Name{
						Loc:   locFn(28, 34),
						Value: "String",
					}),
					Loc: locFn(28, 34),
				}),
				Loc: locFn(21, 34),
			}),
		},
		Loc: locFn(1, 36),
	})
	expectedDocument := ast.NewDocument(&ast.Document{
		Loc:         locFn(1, 36),
		Definitions: []ast.Node{iTypeDef},
	})
	expectedJSON, err := printJSON(expectedDocument)
	if err != nil {
		t.Fatalf("unexpected error, error: %v", err)
	}
	if !reflect.DeepEqual(docJSON, expectedJSON) {
		t.Fatalf("unexpected document, \n\n expected: \n\n%v, \n\ngot: \n\n%v", expectedJSON, docJSON)
	}
}

func TestSimpleFieldWithArg(t *testing.T) {
	body := `
type Hello {
  world(flag: Boolean): String
}`
	params := ParseParams{
		Source: source.NewSource(&source.Source{Body: body}),
	}
	doc, err := Parse(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	docJSON, err := printJSON(doc)
	if err != nil {
		t.Fatalf("unexpected error, error: %v", err)
	}
	locFn := createLoc(body)
	oTypeDef := ast.NewObjectTypeDefinition(&ast.ObjectTypeDefinition{
		Name: ast.NewName(&ast.Name{
			Value: "Hello",
			Loc:   locFn(6, 11),
		}),
		Interfaces: []*ast.NamedType{},
		Fields: []*ast.FieldDefinition{
			ast.NewFieldDefinition(&ast.FieldDefinition{
				Loc: locFn(16, 44),
				Name: ast.NewName(&ast.Name{
					Loc:   locFn(16, 21),
					Value: "world",
				}),
				Arguments: []*ast.InputValueDefinition{
					ast.NewInputValueDefinition(
						&ast.InputValueDefinition{
							Loc: locFn(22, 35),
							Name: ast.NewName(&ast.Name{
								Loc:   locFn(22, 26),
								Value: "flag",
							}),
							Type: ast.NewNamedType(&ast.NamedType{
								Name: ast.NewName(&ast.Name{
									Loc:   locFn(28, 35),
									Value: "Boolean",
								}),
								Loc: locFn(28, 35),
							}),
						}),
				},
				Type: ast.NewNamedType(&ast.NamedType{
					Name: ast.NewName(&ast.Name{
						Loc:   locFn(38, 44),
						Value: "String",
					}),
					Loc: locFn(38, 44),
				}),
			}),
		},
		Loc: locFn(1, 46),
	})
	expectedDocument := ast.NewDocument(&ast.Document{
		Loc:         locFn(1, 46),
		Definitions: []ast.Node{oTypeDef},
	})
	expectedJSON, err := printJSON(expectedDocument)
	if err != nil {
		t.Fatalf("unexpected error, error: %v", err)
	}
	if !reflect.DeepEqual(docJSON, expectedJSON) {
		t.Fatalf("unexpected document, \n\n expected: \n\n%v, \n\ngot: \n\n%v", expectedJSON, docJSON)
	}
}

func TestSimpleFieldWithArgWithDefaultValue(t *testing.T) {
	body := `
type Hello {
  world(flag: Boolean = true): String
}`
	params := ParseParams{
		Source: source.NewSource(&source.Source{Body: body}),
	}
	doc, err := Parse(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	docJSON, err := printJSON(doc)
	if err != nil {
		t.Fatalf("unexpected error, error: %v", err)
	}
	locFn := createLoc(body)
	oTypeDef := ast.NewObjectTypeDefinition(&ast.ObjectTypeDefinition{
		Name: ast.NewName(&ast.Name{
			Value: "Hello",
			Loc:   locFn(6, 11),
		}),
		Interfaces: []*ast.NamedType{},
		Fields: []*ast.FieldDefinition{
			ast.NewFieldDefinition(&ast.FieldDefinition{
				Loc: locFn(16, 51),
				Name: ast.NewName(&ast.Name{
					Loc:   locFn(16, 21),
					Value: "world",
				}),
				Arguments: []*ast.InputValueDefinition{
					ast.NewInputValueDefinition(
						&ast.InputValueDefinition{
							Loc: locFn(22, 42),
							Name: ast.NewName(&ast.Name{
								Loc:   locFn(22, 26),
								Value: "flag",
							}),
							Type: ast.NewNamedType(&ast.NamedType{
								Name: ast.NewName(&ast.Name{
									Value: "Boolean",
									Loc:   locFn(28, 35),
								}),
								Loc: locFn(28, 35),
							}),
							DefaultValue: ast.NewBooleanValue(&ast.BooleanValue{
								Value: true,
								Loc:   locFn(38, 42),
							}),
						}),
				},
				Type: ast.NewNamedType(&ast.NamedType{
					Name: ast.NewName(&ast.Name{
						Loc:   locFn(45, 51),
						Value: "String",
					}),
					Loc: locFn(45, 51),
				}),
			}),
		},
		Loc: locFn(1, 53),
	})
	expectedDocument := ast.NewDocument(&ast.Document{
		Loc:         locFn(1, 53),
		Definitions: []ast.Node{oTypeDef},
	})
	expectedJSON, err := printJSON(expectedDocument)
	if err != nil {
		t.Fatalf("unexpected error, error: %v", err)
	}
	if !reflect.DeepEqual(docJSON, expectedJSON) {
		t.Fatalf("unexpected document, \n\n expected: \n\n%v, \n\ngot: \n\n%v", expectedJSON, docJSON)
	}
}

func TestSimpleFieldWithListArg(t *testing.T) {
	body := `
type Hello {
  world(things: [String]): String
}`
	params := ParseParams{
		Source: source.NewSource(&source.Source{Body: body}),
	}
	doc, err := Parse(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	docJSON, err := printJSON(doc)
	if err != nil {
		t.Fatalf("unexpected error, error: %v", err)
	}
	locFn := createLoc(body)
	oTypeDef := ast.NewObjectTypeDefinition(&ast.ObjectTypeDefinition{
		Name: ast.NewName(&ast.Name{
			Value: "Hello",
			Loc:   locFn(6, 11),
		}),
		Interfaces: []*ast.NamedType{},
		Fields: []*ast.FieldDefinition{
			ast.NewFieldDefinition(&ast.FieldDefinition{
				Loc: locFn(16, 47),
				Name: ast.NewName(&ast.Name{
					Loc:   locFn(16, 21),
					Value: "world",
				}),
				Arguments: []*ast.InputValueDefinition{
					ast.NewInputValueDefinition(
						&ast.InputValueDefinition{
							Loc: locFn(22, 38),
							Name: ast.NewName(&ast.Name{
								Loc:   locFn(22, 28),
								Value: "things",
							}),
							Type: ast.NewListType(&ast.ListType{
								Loc: locFn(30, 38),
								Type: ast.NewNamedType(&ast.NamedType{
									Name: ast.NewName(&ast.Name{
										Loc:   locFn(31, 37),
										Value: "String",
									}),
									Loc: locFn(31, 37),
								}),
							}),
						}),
				},
				Type: ast.NewNamedType(&ast.NamedType{
					Name: ast.NewName(&ast.Name{
						Loc:   locFn(41, 47),
						Value: "String",
					}),
					Loc: locFn(41, 47),
				}),
			}),
		},
		Loc: locFn(1, 49),
	})
	expectedDocument := ast.NewDocument(&ast.Document{
		Loc:         locFn(1, 49),
		Definitions: []ast.Node{oTypeDef},
	})
	expectedJSON, err := printJSON(expectedDocument)
	if err != nil {
		t.Fatalf("unexpected error, error: %v", err)
	}
	if !reflect.DeepEqual(docJSON, expectedJSON) {
		t.Fatalf("unexpected document, \n\n expected: \n\n%v, \n\ngot: \n\n%v", expectedJSON, docJSON)
	}
}

func TestSimpleFieldWithTwoArgs(t *testing.T) {
	body := `
type Hello {
  world(argOne: Boolean, argTwo: Int): String
}`
	params := ParseParams{
		Source: source.NewSource(&source.Source{Body: body}),
	}
	doc, err := Parse(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	docJSON, err := printJSON(doc)
	if err != nil {
		t.Fatalf("unexpected error, error: %v", err)
	}
	locFn := createLoc(body)
	oTypeDef := ast.NewObjectTypeDefinition(&ast.ObjectTypeDefinition{
		Name: ast.NewName(&ast.Name{
			Value: "Hello",
			Loc:   locFn(6, 11),
		}),
		Interfaces: []*ast.NamedType{},
		Fields: []*ast.FieldDefinition{
			ast.NewFieldDefinition(&ast.FieldDefinition{
				Loc: locFn(16, 59),
				Name: ast.NewName(&ast.Name{
					Loc:   locFn(16, 21),
					Value: "world",
				}),
				Arguments: []*ast.InputValueDefinition{
					ast.NewInputValueDefinition(
						&ast.InputValueDefinition{
							Loc: locFn(22, 37),
							Name: ast.NewName(&ast.Name{
								Loc:   locFn(22, 28),
								Value: "argOne",
							}),
							Type: ast.NewNamedType(&ast.NamedType{
								Name: ast.NewName(&ast.Name{
									Loc:   locFn(30, 37),
									Value: "Boolean",
								}),
								Loc: locFn(30, 37),
							}),
						}),
					ast.NewInputValueDefinition(
						&ast.InputValueDefinition{
							Loc: locFn(39, 50),
							Name: ast.NewName(&ast.Name{
								Loc:   locFn(39, 45),
								Value: "argTwo",
							}),
							Type: ast.NewNamedType(&ast.NamedType{
								Name: ast.NewName(&ast.Name{
									Loc:   locFn(47, 50),
									Value: "Int",
								}),
								Loc: locFn(47, 50),
							}),
						}),
				},
				Type: ast.NewNamedType(&ast.NamedType{
					Name: ast.NewName(&ast.Name{
						Loc:   locFn(53, 59),
						Value: "String",
					}),
					Loc: locFn(53, 59),
				}),
			}),
		},
		Loc: locFn(1, 61),
	})
	expectedDocument := ast.NewDocument(&ast.Document{
		Loc:         locFn(1, 61),
		Definitions: []ast.Node{oTypeDef},
	})
	expectedJSON, err := printJSON(expectedDocument)
	if err != nil {
		t.Fatalf("unexpected error, error: %v", err)
	}
	if !reflect.DeepEqual(docJSON, expectedJSON) {
		t.Fatalf("unexpected document, \n\n expected: \n\n%v, \n\ngot: \n\n%v", expectedJSON, docJSON)
	}
}

func TestSimpleUnion(t *testing.T) {
	body := `union Hello = World`
	params := ParseParams{
		Source: source.NewSource(&source.Source{Body: body}),
	}
	doc, err := Parse(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	docJSON, err := printJSON(doc)
	if err != nil {
		t.Fatalf("unexpected error, error: %v", err)
	}
	locFn := createLoc(body)
	uTypeDef := ast.NewUnionTypeDefinition(&ast.UnionTypeDefinition{
		Name: ast.NewName(&ast.Name{
			Value: "Hello",
			Loc:   locFn(6, 11),
		}),
		Types: []*ast.NamedType{
			ast.NewNamedType(&ast.NamedType{
				Name: ast.NewName(&ast.Name{
					Loc:   locFn(14, 19),
					Value: "World",
				}),
				Loc: locFn(14, 19),
			}),
		},
		Loc: locFn(0, 19),
	})
	expectedDocument := ast.NewDocument(&ast.Document{
		Loc:         locFn(0, 19),
		Definitions: []ast.Node{uTypeDef},
	})
	expectedJSON, err := printJSON(expectedDocument)
	if err != nil {
		t.Fatalf("unexpected error, error: %v", err)
	}
	if !reflect.DeepEqual(docJSON, expectedJSON) {
		t.Fatalf("unexpected document, \n\n expected: \n\n%v, \n\ngot: \n\n%v", expectedJSON, docJSON)
	}
}

func TestUnionWithTwoTypes(t *testing.T) {
	body := `union Hello = Wo | Rld`
	params := ParseParams{
		Source: source.NewSource(&source.Source{Body: body}),
	}
	doc, err := Parse(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	docJSON, err := printJSON(doc)
	if err != nil {
		t.Fatalf("unexpected error, error: %v", err)
	}
	locFn := createLoc(body)
	uTypeDef := ast.NewUnionTypeDefinition(&ast.UnionTypeDefinition{
		Name: ast.NewName(&ast.Name{
			Value: "Hello",
			Loc:   locFn(6, 11),
		}),
		Types: []*ast.NamedType{
			ast.NewNamedType(&ast.NamedType{
				Name: ast.NewName(&ast.Name{
					Loc:   locFn(14, 16),
					Value: "Wo",
				}),
				Loc: locFn(14, 16),
			}),
			ast.NewNamedType(&ast.NamedType{
				Name: ast.NewName(&ast.Name{
					Loc:   locFn(19, 22),
					Value: "Rld",
				}),
				Loc: locFn(19, 22),
			}),
		},
		Loc: locFn(0, 22),
	})
	expectedDocument := ast.NewDocument(&ast.Document{
		Loc:         locFn(0, 22),
		Definitions: []ast.Node{uTypeDef},
	})
	expectedJSON, err := printJSON(expectedDocument)
	if err != nil {
		t.Fatalf("unexpected error, error: %v", err)
	}
	if !reflect.DeepEqual(docJSON, expectedJSON) {
		t.Fatalf("unexpected document, \n\n expected: \n\n%v, \n\ngot: \n\n%v", expectedJSON, docJSON)
	}
}
