package graphql

import (
	"reflect"
	"testing"
)

func parsee(t *testing.T, query string) *AstDocument {
	astDoc, err := Parse(ParseParams{
		Source: query,
		Options: ParseOptions{
			NoLocation: false,
			NoSource:   true,
		},
	})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	return astDoc
}

func testLoc(start int, end int) *AstLocation {
	return &AstLocation{
		Start: start, End: end,
	}
}
func TestSchemaParser_SimpleType(t *testing.T) {

	body := `
type Hello {
  world: String
}`
	astDoc := parsee(t, body)
	expected := NewAstDocument(&AstDocument{
		Loc: testLoc(1, 31),
		Definitions: []Node{
			NewAstObjectDefinition(&AstObjectDefinition{
				Loc: testLoc(1, 31),
				Name: NewAstName(&AstName{
					Value: "Hello",
					Loc:   testLoc(6, 11),
				}),
				Interfaces: []*AstNamed{},
				Fields: []*AstFieldDefinition{
					NewAstFieldDefinition(&AstFieldDefinition{
						Loc: testLoc(16, 29),
						Name: NewAstName(&AstName{
							Value: "world",
							Loc:   testLoc(16, 21),
						}),
						Arguments: []*AstInputValueDefinition{},
						Type: NewAstNamed(&AstNamed{
							Loc: testLoc(23, 29),
							Name: NewAstName(&AstName{
								Value: "String",
								Loc:   testLoc(23, 29),
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
	astDoc := parsee(t, body)
	expected := NewAstDocument(&AstDocument{
		Loc: testLoc(1, 38),
		Definitions: []Node{
			NewAstTypeExtensionDefinition(&AstTypeExtensionDefinition{
				Loc: testLoc(1, 38),
				Definition: NewAstObjectDefinition(&AstObjectDefinition{
					Loc: testLoc(8, 38),
					Name: NewAstName(&AstName{
						Value: "Hello",
						Loc:   testLoc(13, 18),
					}),
					Interfaces: []*AstNamed{},
					Fields: []*AstFieldDefinition{
						NewAstFieldDefinition(&AstFieldDefinition{
							Loc: testLoc(23, 36),
							Name: NewAstName(&AstName{
								Value: "world",
								Loc:   testLoc(23, 28),
							}),
							Arguments: []*AstInputValueDefinition{},
							Type: NewAstNamed(&AstNamed{
								Loc: testLoc(30, 36),
								Name: NewAstName(&AstName{
									Value: "String",
									Loc:   testLoc(30, 36),
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
	astDoc := parsee(t, body)
	expected := NewAstDocument(&AstDocument{
		Loc: testLoc(1, 32),
		Definitions: []Node{
			NewAstObjectDefinition(&AstObjectDefinition{
				Loc: testLoc(1, 32),
				Name: NewAstName(&AstName{
					Value: "Hello",
					Loc:   testLoc(6, 11),
				}),
				Interfaces: []*AstNamed{},
				Fields: []*AstFieldDefinition{
					NewAstFieldDefinition(&AstFieldDefinition{
						Loc: testLoc(16, 30),
						Name: NewAstName(&AstName{
							Value: "world",
							Loc:   testLoc(16, 21),
						}),
						Arguments: []*AstInputValueDefinition{},
						Type: NewAstNonNull(&AstNonNull{
							Kind: "NonNullType",
							Loc:  testLoc(23, 30),
							Type: NewAstNamed(&AstNamed{
								Loc: testLoc(23, 29),
								Name: NewAstName(&AstName{
									Value: "String",
									Loc:   testLoc(23, 29),
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
	astDoc := parsee(t, body)
	expected := NewAstDocument(&AstDocument{
		Loc: testLoc(0, 31),
		Definitions: []Node{
			NewAstObjectDefinition(&AstObjectDefinition{
				Loc: testLoc(0, 31),
				Name: NewAstName(&AstName{
					Value: "Hello",
					Loc:   testLoc(5, 10),
				}),
				Interfaces: []*AstNamed{
					NewAstNamed(&AstNamed{
						Name: NewAstName(&AstName{
							Value: "World",
							Loc:   testLoc(22, 27),
						}),
						Loc: testLoc(22, 27),
					}),
				},
				Fields: []*AstFieldDefinition{},
			}),
		},
	})
	if !reflect.DeepEqual(astDoc, expected) {
		t.Fatalf("unexpected document, expected: %v, got: %v", expected, astDoc)
	}
}

func TestSchemaParser_SimpleTypeInheritingMultipleInterfaces(t *testing.T) {
	body := `type Hello implements Wo, rld { }`
	astDoc := parsee(t, body)
	expected := NewAstDocument(&AstDocument{
		Loc: testLoc(0, 33),
		Definitions: []Node{
			NewAstObjectDefinition(&AstObjectDefinition{
				Loc: testLoc(0, 33),
				Name: NewAstName(&AstName{
					Value: "Hello",
					Loc:   testLoc(5, 10),
				}),
				Interfaces: []*AstNamed{
					NewAstNamed(&AstNamed{
						Name: NewAstName(&AstName{
							Value: "Wo",
							Loc:   testLoc(22, 24),
						}),
						Loc: testLoc(22, 24),
					}),
					NewAstNamed(&AstNamed{
						Name: NewAstName(&AstName{
							Value: "rld",
							Loc:   testLoc(26, 29),
						}),
						Loc: testLoc(26, 29),
					}),
				},
				Fields: []*AstFieldDefinition{},
			}),
		},
	})
	if !reflect.DeepEqual(astDoc, expected) {
		t.Fatalf("unexpected document, expected: %v, got: %v", expected, astDoc)
	}
}

func TestSchemaParser_SingleValueEnum(t *testing.T) {
	body := `enum Hello { WORLD }`
	astDoc := parsee(t, body)
	expected := NewAstDocument(&AstDocument{
		Loc: testLoc(0, 20),
		Definitions: []Node{
			NewAstEnumDefinition(&AstEnumDefinition{
				Loc: testLoc(0, 20),
				Name: NewAstName(&AstName{
					Value: "Hello",
					Loc:   testLoc(5, 10),
				}),
				Values: []*AstEnumValueDefinition{
					NewAstEnumValueDefinition(&AstEnumValueDefinition{
						Name: NewAstName(&AstName{
							Value: "WORLD",
							Loc:   testLoc(13, 18),
						}),
						Loc: testLoc(13, 18),
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
	astDoc := parsee(t, body)
	expected := NewAstDocument(&AstDocument{
		Loc: testLoc(0, 22),
		Definitions: []Node{
			NewAstEnumDefinition(&AstEnumDefinition{
				Loc: testLoc(0, 22),
				Name: NewAstName(&AstName{
					Value: "Hello",
					Loc:   testLoc(5, 10),
				}),
				Values: []*AstEnumValueDefinition{
					NewAstEnumValueDefinition(&AstEnumValueDefinition{
						Name: NewAstName(&AstName{
							Value: "WO",
							Loc:   testLoc(13, 15),
						}),
						Loc: testLoc(13, 15),
					}),
					NewAstEnumValueDefinition(&AstEnumValueDefinition{
						Name: NewAstName(&AstName{
							Value: "RLD",
							Loc:   testLoc(17, 20),
						}),
						Loc: testLoc(17, 20),
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
	astDoc := parsee(t, body)
	expected := NewAstDocument(&AstDocument{
		Loc: testLoc(1, 36),
		Definitions: []Node{
			NewAstInterfaceDefinition(&AstInterfaceDefinition{
				Loc: testLoc(1, 36),
				Name: NewAstName(&AstName{
					Value: "Hello",
					Loc:   testLoc(11, 16),
				}),
				Fields: []*AstFieldDefinition{
					NewAstFieldDefinition(&AstFieldDefinition{
						Loc: testLoc(21, 34),
						Name: NewAstName(&AstName{
							Value: "world",
							Loc:   testLoc(21, 26),
						}),
						Arguments: []*AstInputValueDefinition{},
						Type: NewAstNamed(&AstNamed{
							Loc: testLoc(28, 34),
							Name: NewAstName(&AstName{
								Value: "String",
								Loc:   testLoc(28, 34),
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
	astDoc := parsee(t, body)
	expected := NewAstDocument(&AstDocument{
		Loc: testLoc(1, 46),
		Definitions: []Node{
			NewAstObjectDefinition(&AstObjectDefinition{
				Loc: testLoc(1, 46),
				Name: NewAstName(&AstName{
					Value: "Hello",
					Loc:   testLoc(6, 11),
				}),
				Interfaces: []*AstNamed{},
				Fields: []*AstFieldDefinition{
					NewAstFieldDefinition(&AstFieldDefinition{
						Loc: testLoc(16, 44),
						Name: NewAstName(&AstName{
							Value: "world",
							Loc:   testLoc(16, 21),
						}),
						Arguments: []*AstInputValueDefinition{
							NewAstInputValueDefinition(&AstInputValueDefinition{
								Loc: testLoc(22, 35),
								Name: NewAstName(&AstName{
									Value: "flag",
									Loc:   testLoc(22, 26),
								}),
								Type: NewAstNamed(&AstNamed{
									Loc: testLoc(28, 35),
									Name: NewAstName(&AstName{
										Value: "Boolean",
										Loc:   testLoc(28, 35),
									}),
								}),
								DefaultValue: nil,
							}),
						},
						Type: NewAstNamed(&AstNamed{
							Loc: testLoc(38, 44),
							Name: NewAstName(&AstName{
								Value: "String",
								Loc:   testLoc(38, 44),
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
	astDoc := parsee(t, body)
	expected := NewAstDocument(&AstDocument{
		Loc: testLoc(1, 53),
		Definitions: []Node{
			NewAstObjectDefinition(&AstObjectDefinition{
				Loc: testLoc(1, 53),
				Name: NewAstName(&AstName{
					Value: "Hello",
					Loc:   testLoc(6, 11),
				}),
				Interfaces: []*AstNamed{},
				Fields: []*AstFieldDefinition{
					NewAstFieldDefinition(&AstFieldDefinition{
						Loc: testLoc(16, 51),
						Name: NewAstName(&AstName{
							Value: "world",
							Loc:   testLoc(16, 21),
						}),
						Arguments: []*AstInputValueDefinition{
							NewAstInputValueDefinition(&AstInputValueDefinition{
								Loc: testLoc(22, 42),
								Name: NewAstName(&AstName{
									Value: "flag",
									Loc:   testLoc(22, 26),
								}),
								Type: NewAstNamed(&AstNamed{
									Loc: testLoc(28, 35),
									Name: NewAstName(&AstName{
										Value: "Boolean",
										Loc:   testLoc(28, 35),
									}),
								}),
								DefaultValue: NewAstBooleanValue(&AstBooleanValue{
									Value: true,
									Loc:   testLoc(38, 42),
								}),
							}),
						},
						Type: NewAstNamed(&AstNamed{
							Loc: testLoc(45, 51),
							Name: NewAstName(&AstName{
								Value: "String",
								Loc:   testLoc(45, 51),
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
	astDoc := parsee(t, body)
	expected := NewAstDocument(&AstDocument{
		Loc: testLoc(1, 49),
		Definitions: []Node{
			NewAstObjectDefinition(&AstObjectDefinition{
				Loc: testLoc(1, 49),
				Name: NewAstName(&AstName{
					Value: "Hello",
					Loc:   testLoc(6, 11),
				}),
				Interfaces: []*AstNamed{},
				Fields: []*AstFieldDefinition{
					NewAstFieldDefinition(&AstFieldDefinition{
						Loc: testLoc(16, 47),
						Name: NewAstName(&AstName{
							Value: "world",
							Loc:   testLoc(16, 21),
						}),
						Arguments: []*AstInputValueDefinition{
							NewAstInputValueDefinition(&AstInputValueDefinition{
								Loc: testLoc(22, 38),
								Name: NewAstName(&AstName{
									Value: "things",
									Loc:   testLoc(22, 28),
								}),
								Type: NewAstList(&AstList{
									Loc: testLoc(30, 38),
									Type: NewAstNamed(&AstNamed{
										Loc: testLoc(31, 37),
										Name: NewAstName(&AstName{
											Value: "String",
											Loc:   testLoc(31, 37),
										}),
									}),
								}),
								DefaultValue: nil,
							}),
						},
						Type: NewAstNamed(&AstNamed{
							Loc: testLoc(41, 47),
							Name: NewAstName(&AstName{
								Value: "String",
								Loc:   testLoc(41, 47),
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
	astDoc := parsee(t, body)
	expected := NewAstDocument(&AstDocument{
		Loc: testLoc(1, 61),
		Definitions: []Node{
			NewAstObjectDefinition(&AstObjectDefinition{
				Loc: testLoc(1, 61),
				Name: NewAstName(&AstName{
					Value: "Hello",
					Loc:   testLoc(6, 11),
				}),
				Interfaces: []*AstNamed{},
				Fields: []*AstFieldDefinition{
					NewAstFieldDefinition(&AstFieldDefinition{
						Loc: testLoc(16, 59),
						Name: NewAstName(&AstName{
							Value: "world",
							Loc:   testLoc(16, 21),
						}),
						Arguments: []*AstInputValueDefinition{
							NewAstInputValueDefinition(&AstInputValueDefinition{
								Loc: testLoc(22, 37),
								Name: NewAstName(&AstName{
									Value: "argOne",
									Loc:   testLoc(22, 28),
								}),
								Type: NewAstNamed(&AstNamed{
									Loc: testLoc(30, 37),
									Name: NewAstName(&AstName{
										Value: "Boolean",
										Loc:   testLoc(30, 37),
									}),
								}),
								DefaultValue: nil,
							}),
							NewAstInputValueDefinition(&AstInputValueDefinition{
								Loc: testLoc(39, 50),
								Name: NewAstName(&AstName{
									Value: "argTwo",
									Loc:   testLoc(39, 45),
								}),
								Type: NewAstNamed(&AstNamed{
									Loc: testLoc(47, 50),
									Name: NewAstName(&AstName{
										Value: "Int",
										Loc:   testLoc(47, 50),
									}),
								}),
								DefaultValue: nil,
							}),
						},
						Type: NewAstNamed(&AstNamed{
							Loc: testLoc(53, 59),
							Name: NewAstName(&AstName{
								Value: "String",
								Loc:   testLoc(53, 59),
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
	astDoc := parsee(t, body)
	expected := NewAstDocument(&AstDocument{
		Loc: testLoc(0, 19),
		Definitions: []Node{
			NewAstUnionDefinition(&AstUnionDefinition{
				Loc: testLoc(0, 19),
				Name: NewAstName(&AstName{
					Value: "Hello",
					Loc:   testLoc(6, 11),
				}),
				Types: []*AstNamed{
					NewAstNamed(&AstNamed{
						Loc: testLoc(14, 19),
						Name: NewAstName(&AstName{
							Value: "World",
							Loc:   testLoc(14, 19),
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
	astDoc := parsee(t, body)
	expected := NewAstDocument(&AstDocument{
		Loc: testLoc(0, 22),
		Definitions: []Node{
			NewAstUnionDefinition(&AstUnionDefinition{
				Loc: testLoc(0, 22),
				Name: NewAstName(&AstName{
					Value: "Hello",
					Loc:   testLoc(6, 11),
				}),
				Types: []*AstNamed{
					NewAstNamed(&AstNamed{
						Loc: testLoc(14, 16),
						Name: NewAstName(&AstName{
							Value: "Wo",
							Loc:   testLoc(14, 16),
						}),
					}),
					NewAstNamed(&AstNamed{
						Loc: testLoc(19, 22),
						Name: NewAstName(&AstName{
							Value: "Rld",
							Loc:   testLoc(19, 22),
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
	astDoc := parsee(t, body)
	expected := NewAstDocument(&AstDocument{
		Loc: testLoc(0, 12),
		Definitions: []Node{
			NewAstScalarDefinition(&AstScalarDefinition{
				Loc: testLoc(0, 12),
				Name: NewAstName(&AstName{
					Value: "Hello",
					Loc:   testLoc(7, 12),
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
	astDoc := parsee(t, body)
	expected := NewAstDocument(&AstDocument{
		Loc: testLoc(1, 32),
		Definitions: []Node{
			NewAstInputObjectDefinition(&AstInputObjectDefinition{
				Loc: testLoc(1, 32),
				Name: NewAstName(&AstName{
					Value: "Hello",
					Loc:   testLoc(7, 12),
				}),
				Fields: []*AstInputValueDefinition{
					NewAstInputValueDefinition(&AstInputValueDefinition{
						Loc: testLoc(17, 30),
						Name: NewAstName(&AstName{
							Value: "world",
							Loc:   testLoc(17, 22),
						}),
						Type: NewAstNamed(&AstNamed{
							Loc: testLoc(24, 30),
							Name: NewAstName(&AstName{
								Value: "String",
								Loc:   testLoc(24, 30),
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

	_, err := Parse(ParseParams{
		Source: body,
		Options: ParseOptions{
			NoLocation: false,
			NoSource:   true,
		},
	})

	expectedError := &Error{
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
		Nodes: []Node{},
		Source: &Source{
			Body: `
input Hello {
  world(foo: Int): String
}`,
			Name: "GraphQL",
		},
		Positions: []int{22},
		Locations: []SourceLocation{
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
