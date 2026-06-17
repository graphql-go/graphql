package parser

import (
	"testing"

	"github.com/graphql-go/graphql/language/source"
)

func TestParseValue_Float(t *testing.T) {
	v, err := ParseValue(ParseParams{
		Source:  `3.14`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v == nil {
		t.Fatal("expected a value, got nil")
	}
}

func TestParseValue_StringInput(t *testing.T) {
	v, err := ParseValue(ParseParams{
		Source:  `"hello"`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v == nil {
		t.Fatal("expected a value, got nil")
	}
}

func TestParseValue_Int(t *testing.T) {
	v, err := ParseValue(ParseParams{
		Source:  `42`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v == nil {
		t.Fatal("expected a value, got nil")
	}
}

func TestParseValue_Bool(t *testing.T) {
	v, err := ParseValue(ParseParams{
		Source:  `true`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v == nil {
		t.Fatal("expected a value, got nil")
	}
}

func TestParseValue_Enum(t *testing.T) {
	v, err := ParseValue(ParseParams{
		Source:  `SOME_ENUM`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v == nil {
		t.Fatal("expected a value, got nil")
	}
}

func TestParseValue_List(t *testing.T) {
	v, err := ParseValue(ParseParams{
		Source:  `[1, 2, 3]`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v == nil {
		t.Fatal("expected a value, got nil")
	}
}

func TestParseValue_Object(t *testing.T) {
	v, err := ParseValue(ParseParams{
		Source:  `{key: "value"}`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v == nil {
		t.Fatal("expected a value, got nil")
	}
}

func TestParseValue_SourceStruct(t *testing.T) {
	v, err := ParseValue(ParseParams{
		Source: source.NewSource(&source.Source{
			Body: []byte(`"hello"`),
		}),
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v == nil {
		t.Fatal("expected a value, got nil")
	}
}

func TestParseValue_Error(t *testing.T) {
	_, err := ParseValue(ParseParams{
		Source:  `@bad`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestParseValue_BadToken(t *testing.T) {
	_, err := ParseValue(ParseParams{
		Source:  "`",
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestParse_SchemaDefinition(t *testing.T) {
	sdl := `
		schema {
			query: QueryRoot
			mutation: MutationRoot
		}
		type QueryRoot { hello: String }
		type MutationRoot { set: Boolean }
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_SchemaDefinitionWithDirectives(t *testing.T) {
	sdl := `
		schema @directive {
			query: QueryRoot
		}
		type QueryRoot { hello: String }
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_ScalarDefinition(t *testing.T) {
	sdl := `
		scalar JSON
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_ScalarDefinitionWithDirective(t *testing.T) {
	sdl := `
		scalar JSON @specifiedBy(url: "https://foo.com")
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_InterfaceDefinition(t *testing.T) {
	sdl := `
		interface NamedEntity {
			name: String
		}
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_InterfaceDefinitionWithDirectives(t *testing.T) {
	sdl := `
		interface NamedEntity @deprecated {
			name: String
		}
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_UnionDefinition(t *testing.T) {
	sdl := `
		union SearchResult = Photo | Person
		type Photo { url: String }
		type Person { name: String }
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_UnionDefinitionWithDirectives(t *testing.T) {
	sdl := `
		union SearchResult @deprecated = Photo | Person
		type Photo { url: String }
		type Person { name: String }
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_EnumDefinition(t *testing.T) {
	sdl := `
		enum Color {
			RED
			GREEN
			BLUE
		}
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_EnumDefinitionWithDirectives(t *testing.T) {
	sdl := `
		enum Color @deprecated {
			RED
			GREEN
			BLUE
		}
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_EnumValueWithDirective(t *testing.T) {
	sdl := `
		enum Color {
			RED @deprecated
			GREEN
		}
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_InputObjectDefinition(t *testing.T) {
	sdl := `
		input Filter {
			name: String
			age: Int
		}
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_InputObjectDefinitionWithDirectives(t *testing.T) {
	sdl := `
		input Filter @deprecated {
			name: String
		}
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_TypeExtension(t *testing.T) {
	sdl := `
		extend type Hello {
			world: String
		}
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_DirectiveDefinition(t *testing.T) {
	sdl := `
		directive @myDirective on FIELD
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_DirectiveDefinitionWithArgs(t *testing.T) {
	sdl := `
		directive @myDirective(arg: Int!) on FIELD | OBJECT
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_DirectiveDefinitionWithInputValueDefault(t *testing.T) {
	sdl := `
		directive @myDirective(arg: Int = 42) on FIELD
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_FloatDefaultValue(t *testing.T) {
	q := `query Foo($x: Float = 3.14) { field }`
	_, err := Parse(ParseParams{
		Source:  q,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_FloatArgumentValue(t *testing.T) {
	q := `{ field(arg: 3.14) }`
	_, err := Parse(ParseParams{
		Source:  q,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_NonNullListType(t *testing.T) {
	sdl := `
		type Query {
			field(arg: [String]!): [Int!]!
		}
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_ObjectTypeWithDirectivesAndInterfaces(t *testing.T) {
	sdl := `
		interface Named { name: String }
		type Person implements Named @deprecated {
			name: String
			age: Int
		}
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_ObjectTypeWithMultipleInterfaces(t *testing.T) {
	sdl := `
		interface A { a: String }
		interface B { b: Int }
		type C implements A & B {
			a: String
			b: Int
		}
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_FragmentSpreadWithDirectives(t *testing.T) {
	q := `
		query { ...myFrag @include(if: true) }
		fragment myFrag on Query { field }
	`
	_, err := Parse(ParseParams{
		Source:  q,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_InlineFragmentWithDirectives(t *testing.T) {
	q := `
		query {
			... on User @skip(if: false) {
				name
			}
		}
	`
	_, err := Parse(ParseParams{
		Source:  q,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_FieldWithAlias(t *testing.T) {
	q := `{ aliasName: realName }`
	_, err := Parse(ParseParams{
		Source:  q,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_FieldWithDirective(t *testing.T) {
	q := `{ field @deprecated }`
	_, err := Parse(ParseParams{
		Source:  q,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_NoLocation(t *testing.T) {
	doc, err := Parse(ParseParams{
		Source:  "{ field }",
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if doc == nil {
		t.Fatal("expected a document")
	}
}

func TestParse_FieldDefinitionWithArgs(t *testing.T) {
	sdl := `
		type Query {
			field(arg1: String, arg2: Int = 42): String @deprecated
		}
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_InputValueWithDefaultAndDirective(t *testing.T) {
	sdl := `
		type Query {
			field(arg: Int = 10 @deprecated): String
		}
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_OperationWithDirectives(t *testing.T) {
	q := `query Foo @deprecated { field }`
	_, err := Parse(ParseParams{
		Source:  q,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_DefaultValueObject(t *testing.T) {
	sdl := `
		input Filter { name: String }
		type Query {
			field(filter: Filter = {name: "foo"}): String
		}
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_DefaultValueList(t *testing.T) {
	sdl := `
		type Query {
			field(ids: [Int] = [1, 2, 3]): String
		}
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_DescriptionOnObjectType(t *testing.T) {
	sdl := `
		" A type with a description "
		type MyType {
			field: String
		}
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_DescriptionOnInterface(t *testing.T) {
	sdl := `
		" An interface "
		interface Named {
			name: String
		}
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_DescriptionOnUnion(t *testing.T) {
	sdl := `
		" A union type "
		union Result = A | B
		type A { a: String }
		type B { b: Int }
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_DescriptionOnEnum(t *testing.T) {
	sdl := `
		" Color enum "
		enum Color {
			RED
		}
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_DescriptionOnInput(t *testing.T) {
	sdl := `
		" Input type "
		input Filter {
			name: String
		}
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_DescriptionOnScalar(t *testing.T) {
	sdl := `
		" A custom scalar "
		scalar CustomScalar
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_DescriptionOnDirective(t *testing.T) {
	sdl := `
		"A directive"
		directive @custom on FIELD
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_DescriptionOnFieldDefinition(t *testing.T) {
	sdl := `
		type Query {
			"field doc"
			hello: String
		}
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_InvalidOperationType(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `teapot { field }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected an error for invalid operation type")
	}
}

func TestParse_EmptySchemaBody(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `schema { }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected an error for empty schema body")
	}
}

func TestParse_SchemaMissingColon(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `schema { query QueryRoot }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected an error for schema missing colon")
	}
}

func TestParse_InvalidDirectiveDefinition(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `directive @ on FIELD`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected an error for invalid directive definition")
	}
}

func TestParse_ObjectTypeMissingName(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `type { }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected an error for type missing name")
	}
}

func TestParse_EnumTypeMissingName(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `enum { }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected an error for enum missing name")
	}
}

func TestParse_InputTypeMissingName(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `input { }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected an error for input missing name")
	}
}

func TestParse_InterfaceTypeMissingName(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `interface { }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected an error for interface missing name")
	}
}

func TestParse_NullValue(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `{ field(arg: null) }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected an error for null value")
	}
}

func TestParse_UnionMissingEquals(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `union SearchResult Photo`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected an error for union missing =")
	}
}

func TestParse_ExtendTypeMissingName(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `extend type { }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected an error for extend type missing name")
	}
}

func TestParse_MutationOperationType(t *testing.T) {
	q := `mutation { doStuff }`
	_, err := Parse(ParseParams{
		Source:  q,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_SubscriptionOperationType(t *testing.T) {
	q := `subscription { watch }`
	_, err := Parse(ParseParams{
		Source:  q,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_NamedDirectiveDefinition(t *testing.T) {
	sdl := `
		directive @skip(if: Boolean!) on FIELD | FRAGMENT_SPREAD | INLINE_FRAGMENT
	`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_BadTokenInParse(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  "`",
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected an error for bad token")
	}
}
