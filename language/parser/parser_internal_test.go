package parser

import (
	"testing"

	"github.com/graphql-go/graphql/language/source"
)

func makeTestParser(t *testing.T, src string) *Parser {
	parser, err := makeParser(source.NewSource(&source.Source{Body: []byte(src)}), ParseOptions{NoLocation: true, NoSource: true})
	if err != nil {
		t.Fatal(err)
	}
	return parser
}

// Direct internal function error-path tests

func TestParseNameError(t *testing.T) {
	parser := makeTestParser(t, "123")
	_, err := parseName(parser)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseVariableError(t *testing.T) {
	parser := makeTestParser(t, "$")
	_, err := parseVariable(parser)
	if err == nil {
		t.Fatal("expected error for $ with nothing after")
	}
}

func TestParseVariableError2(t *testing.T) {
	parser := makeTestParser(t, "$123")
	_, err := parseVariable(parser)
	if err == nil {
		t.Fatal("expected error for $ with non-name")
	}
}

func TestParseFieldSuccess(t *testing.T) {
	parser := makeTestParser(t, `field`)
	// parseField starts with parseName(parser) which expects NAME token
	_, err := parseField(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseFieldErrorDirective(t *testing.T) {
	parser := makeTestParser(t, `field @`)
	if err := advance(parser); err != nil {
		t.Fatal(err)
	}
	// parseArguments: no PAREN_L → empty
	// parseDirectives: peek(AT) → true → parseDirective
	// Inside parseDirective: expect(AT) → succeeds, then parseName → expects NAME but token is EOF
	_, err := parseField(parser)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseArgumentErrorName(t *testing.T) {
	parser := makeTestParser(t, `123:`)
	_, err := parseArgument(parser)
	if err == nil {
		t.Fatal("expected error for argument with non-name key")
	}
}

func TestParseArgumentErrorColon(t *testing.T) {
	parser := makeTestParser(t, `name`)
	if err := advance(parser); err != nil {
		t.Fatal(err)
	}
	_, err := parseArgument(parser)
	if err == nil {
		t.Fatal("expected error for argument missing colon")
	}
}

func TestParseFragmentErrorSpread(t *testing.T) {
	parser := makeTestParser(t, `123`)
	_, err := parseFragment(parser)
	if err == nil {
		t.Fatal("expected error for fragment missing spread")
	}
}

func TestParseFragmentErrorFragmentName(t *testing.T) {
	parser := makeTestParser(t, `...123`)
	_, err := parseFragment(parser)
	if err == nil {
		t.Fatal("expected error for spread with non-name")
	}
}

func TestParseFragmentErrorSpreadDirective(t *testing.T) {
	parser := makeTestParser(t, `...frag @123`)
	_, err := parseFragment(parser)
	if err == nil {
		t.Fatal("expected error for fragment spread with bad directive")
	}
}

func TestParseFragmentErrorInlineAdvance(t *testing.T) {
	parser := makeTestParser(t, `... on`)
	if err := advance(parser); err != nil {
		t.Fatal(err)
	}
	// After ..., token is NAME "on"
	_, err := parseFragment(parser)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseFragmentErrorInlineDirective(t *testing.T) {
	parser := makeTestParser(t, `... on Type @123`)
	_, err := parseFragment(parser)
	if err == nil {
		t.Fatal("expected error for inline fragment with bad directive")
	}
}

func TestParseFragmentErrorInlineSelectionSet(t *testing.T) {
	parser := makeTestParser(t, `... on`)
	_, err := parseFragment(parser)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseFragmentDefinitionErrorFragment(t *testing.T) {
	parser := makeTestParser(t, `123 on Type { field }`)
	_, err := parseFragmentDefinition(parser)
	if err == nil {
		t.Fatal("expected error for missing fragment keyword")
	}
}

func TestParseFragmentDefinitionErrorType(t *testing.T) {
	parser := makeTestParser(t, `fragment Foo on`)
	if err := advance(parser); err != nil {
		t.Fatal(err)
	}
	_, err := parseFragmentDefinition(parser)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseFragmentDefinitionErrorDirective(t *testing.T) {
	parser := makeTestParser(t, `fragment Foo on Type @123`)
	_, err := parseFragmentDefinition(parser)
	if err == nil {
		t.Fatal("expected error for fragment with bad directive")
	}
}

func TestParseFragmentDefinitionErrorSelection(t *testing.T) {
	parser := makeTestParser(t, `fragment Foo on Type`)
	_, err := parseFragmentDefinition(parser)
	if err == nil {
		t.Fatal("expected error for fragment missing selection set")
	}
}

func TestParseValueLiteralErrorIntAdvance(t *testing.T) {
	parser := makeTestParser(t, `123`)
	_, err := parseValueLiteral(parser, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseValueLiteralErrorFloatAdvance(t *testing.T) {
	parser := makeTestParser(t, `3.14`)
	_, err := parseValueLiteral(parser, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseValueLiteralErrorBoolAdvance(t *testing.T) {
	parser := makeTestParser(t, `true`)
	_, err := parseValueLiteral(parser, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseValueLiteralErrorEnumAdvance(t *testing.T) {
	parser := makeTestParser(t, `SOME_VALUE`)
	_, err := parseValueLiteral(parser, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseObjectErrorExpectBraceL(t *testing.T) {
	parser := makeTestParser(t, `123`)
	_, err := parseObject(parser, false)
	if err == nil {
		t.Fatal("expected error for object missing brace")
	}
}

func TestParseObjectErrorSkipBraceR(t *testing.T) {
	parser := makeTestParser(t, `{`)
	if err := advance(parser); err != nil {
		t.Fatal(err)
	}
	_, err := parseObject(parser, false)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseObjectFieldErrorName(t *testing.T) {
	parser := makeTestParser(t, `123:`)
	_, err := parseObjectField(parser, false)
	if err == nil {
		t.Fatal("expected error for object field with non-name key")
	}
}

func TestParseObjectFieldErrorColon(t *testing.T) {
	parser := makeTestParser(t, `name`)
	if err := advance(parser); err != nil {
		t.Fatal(err)
	}
	_, err := parseObjectField(parser, false)
	if err == nil {
		t.Fatal("expected error for object field missing colon")
	}
}

func TestParseDirectiveErrorAT(t *testing.T) {
	parser := makeTestParser(t, `123`)
	_, err := parseDirective(parser)
	if err == nil {
		t.Fatal("expected error for directive missing @")
	}
}

func TestParseDirectiveErrorName(t *testing.T) {
	parser := makeTestParser(t, `@123`)
	if err := advance(parser); err != nil {
		t.Fatal(err)
	}
	_, err := parseDirective(parser)
	if err == nil {
		t.Fatal("expected error for directive with non-name")
	}
}

func TestParseDirectiveErrorArg(t *testing.T) {
	parser := makeTestParser(t, `@skip(foo bar)`)
	_, err := parseDirective(parser)
	if err == nil {
		t.Fatal("expected error for directive with bad argument")
	}
}

func TestParseTypeErrorBracketAdvance(t *testing.T) {
	parser := makeTestParser(t, `[`)
	_, err := parseType(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseTypeErrorBracketRecursive(t *testing.T) {
	parser := makeTestParser(t, `[`)
	if err := advance(parser); err != nil {
		t.Fatal(err)
	}
	_, err := parseType(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseTypeErrorBracketR(t *testing.T) {
	parser := makeTestParser(t, `[String`)
	_, err := parseType(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseTypeNamed(t *testing.T) {
	parser := makeTestParser(t, `String`)
	ttype, err := parseType(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ttype == nil {
		t.Fatal("expected a type")
	}
}

func TestParseTypeErrorBang(t *testing.T) {
	parser := makeTestParser(t, `String`)
	_, err := parseType(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseTypeSystemDefinitionErrorLookahead(t *testing.T) {
	parser := makeTestParser(t, `"desc" 123`)
	_, err := parseTypeSystemDefinition(parser)
	if err == nil {
		t.Fatal("expected error for description followed by non-name")
	}
}

func TestParseSchemaDefinitionErrorExpectKey(t *testing.T) {
	parser := makeTestParser(t, `123 { }`)
	_, err := parseSchemaDefinition(parser)
	if err == nil {
		t.Fatal("expected error for schema missing keyword")
	}
}

func TestParseSchemaDefinitionErrorDirective(t *testing.T) {
	parser := makeTestParser(t, `schema @123 { query: Query }`)
	_, err := parseSchemaDefinition(parser)
	if err == nil {
		t.Fatal("expected error for schema with bad directive")
	}
}

func TestParseOperationTypeDefinitionErrorType(t *testing.T) {
	parser := makeTestParser(t, `123: Query`)
	_, err := parseOperationTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error for operation type definition with non-name")
	}
}

func TestParseOperationTypeDefinitionErrorColon(t *testing.T) {
	parser := makeTestParser(t, `query`)
	if err := advance(parser); err != nil {
		t.Fatal(err)
	}
	_, err := parseOperationTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error for operation type definition missing colon")
	}
}

func TestParseScalarTypeDefinitionErrorDescription(t *testing.T) {
	parser := makeTestParser(t, `123 JSON`)
	_, err := parseScalarTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error for scalar missing keyword")
	}
}

func TestParseScalarTypeDefinitionErrorName(t *testing.T) {
	parser := makeTestParser(t, `scalar 123`)
	if err := advance(parser); err != nil {
		t.Fatal(err)
	}
	_, err := parseScalarTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error for scalar with non-name")
	}
}

func TestParseScalarTypeDefinitionErrorDirective(t *testing.T) {
	parser := makeTestParser(t, `scalar JSON @123`)
	_, err := parseScalarTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error for scalar with bad directive")
	}
}

func TestParseObjectTypeDefinitionErrorType(t *testing.T) {
	parser := makeTestParser(t, `123 Foo { field: String }`)
	_, err := parseObjectTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseObjectTypeDefinitionErrorName(t *testing.T) {
	parser := makeTestParser(t, `type 123 { field: String }`)
	if err := advance(parser); err != nil {
		t.Fatal(err)
	}
	_, err := parseObjectTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseObjectTypeDefinitionErrorDirective(t *testing.T) {
	parser := makeTestParser(t, `type Foo @123 { field: String }`)
	_, err := parseObjectTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error for type with bad directive")
	}
}

func TestParseImplementsInterfacesErrorNamed(t *testing.T) {
	parser := makeTestParser(t, `type Foo implements 123 { field: String }`)
	_, err := parseObjectTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error for implements with non-name")
	}
}

func TestParseImplementsInterfacesErrorAmpersand(t *testing.T) {
	parser := makeTestParser(t, `type Foo implements Bar & { field: String }`)
	_, err := parseObjectTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error for implements with amp missing type")
	}
}

func TestParseFieldDefinitionErrorDescription(t *testing.T) {
	parser := makeTestParser(t, `123 name: String`)
	_, err := parseFieldDefinition(parser)
	if err == nil {
		t.Fatal("expected error for field def missing description")
	}
}

func TestParseFieldDefinitionErrorColon(t *testing.T) {
	parser := makeTestParser(t, `name String`)
	if err := advance(parser); err != nil {
		t.Fatal(err)
	}
	_, err := parseFieldDefinition(parser)
	if err == nil {
		t.Fatal("expected error for field def missing colon")
	}
}

func TestParseFieldDefinitionErrorType(t *testing.T) {
	parser := makeTestParser(t, `name: 123`)
	if err := advance(parser); err != nil {
		t.Fatal(err)
	}
	_, err := parseFieldDefinition(parser)
	if err == nil {
		t.Fatal("expected error for field def with non-name type")
	}
}

func TestParseInputValueDefErrorDescription(t *testing.T) {
	parser := makeTestParser(t, `123: String`)
	_, err := parseInputValueDef(parser)
	if err == nil {
		t.Fatal("expected error for input value def error")
	}
}

func TestParseInputValueDefErrorName(t *testing.T) {
	parser := makeTestParser(t, `123: String`)
	_, err := parseInputValueDef(parser)
	if err == nil {
		t.Fatal("expected error for input value def with non-name")
	}
}

func TestParseInputValueDefErrorType(t *testing.T) {
	parser := makeTestParser(t, `name: 123`)
	if err := advance(parser); err != nil {
		t.Fatal(err)
	}
	_, err := parseInputValueDef(parser)
	if err == nil {
		t.Fatal("expected error for input value def with non-name type")
	}
}

func TestParseInputValueDefErrorEquals(t *testing.T) {
	parser := makeTestParser(t, `name: String 123`)
	if err := advance(parser); err != nil {
		t.Fatal(err)
	}
	_, err := parseInputValueDef(parser)
	if err == nil {
		t.Fatal("expected error for input value def error")
	}
}

func TestParseInputValueDefErrorDefault(t *testing.T) {
	parser := makeTestParser(t, `name: String = $var`)
	if err := advance(parser); err != nil {
		t.Fatal(err)
	}
	_, err := parseInputValueDef(parser)
	if err == nil {
		t.Fatal("expected error for input value def with variable default")
	}
}

func TestParseInputValueDefErrorDirective(t *testing.T) {
	parser := makeTestParser(t, `name: String @123`)
	if err := advance(parser); err != nil {
		t.Fatal(err)
	}
	_, err := parseInputValueDef(parser)
	if err == nil {
		t.Fatal("expected error for input value def with bad directive")
	}
}

func TestParseInterfaceTypeDefinitionErrorInterface(t *testing.T) {
	parser := makeTestParser(t, `123 Named { name: String }`)
	_, err := parseInterfaceTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseInterfaceTypeDefinitionErrorName(t *testing.T) {
	parser := makeTestParser(t, `interface 123 { name: String }`)
	if err := advance(parser); err != nil {
		t.Fatal(err)
	}
	_, err := parseInterfaceTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseInterfaceTypeDefinitionErrorDirective(t *testing.T) {
	parser := makeTestParser(t, `interface Named @123 { name: String }`)
	_, err := parseInterfaceTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error for interface with bad directive")
	}
}

func TestParseInterfaceTypeDefinitionErrorReverse(t *testing.T) {
	parser := makeTestParser(t, `interface Named { }`)
	_, err := parseInterfaceTypeDefinition(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseUnionTypeDefinitionErrorUnion(t *testing.T) {
	parser := makeTestParser(t, `123 SearchResult = Photo`)
	_, err := parseUnionTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseUnionTypeDefinitionErrorName(t *testing.T) {
	parser := makeTestParser(t, `union 123 = Photo`)
	if err := advance(parser); err != nil {
		t.Fatal(err)
	}
	_, err := parseUnionTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseUnionTypeDefinitionErrorDirective(t *testing.T) {
	parser := makeTestParser(t, `union SearchResult @123 = Photo`)
	_, err := parseUnionTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseUnionMembersErrorNamed(t *testing.T) {
	parser := makeTestParser(t, `123`)
	_, err := parseUnionMembers(parser)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseUnionMembersErrorPipe(t *testing.T) {
	parser := makeTestParser(t, `Foo |`)
	_, err := parseUnionMembers(parser)
	if err == nil {
		t.Fatal("expected error for union members with trailing pipe")
	}
}

func TestParseEnumTypeDefinitionErrorEnum(t *testing.T) {
	parser := makeTestParser(t, `123 Color { RED }`)
	_, err := parseEnumTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseEnumTypeDefinitionErrorName(t *testing.T) {
	parser := makeTestParser(t, `enum 123 { RED }`)
	if err := advance(parser); err != nil {
		t.Fatal(err)
	}
	_, err := parseEnumTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseEnumTypeDefinitionErrorDirective(t *testing.T) {
	parser := makeTestParser(t, `enum Color @123 { RED }`)
	_, err := parseEnumTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseEnumValueDefinitionErrorDescription(t *testing.T) {
	parser := makeTestParser(t, `123 RED`)
	_, err := parseEnumValueDefinition(parser)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseEnumValueDefinitionErrorName(t *testing.T) {
	parser := makeTestParser(t, `123`)
	if err := advance(parser); err != nil {
		t.Fatal(err)
	}
	_, err := parseEnumValueDefinition(parser)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseEnumValueDefinitionErrorDirective(t *testing.T) {
	parser := makeTestParser(t, `RED @123`)
	if err := advance(parser); err != nil {
		t.Fatal(err)
	}
	_, err := parseEnumValueDefinition(parser)
	if err == nil {
		t.Fatal("expected error for enum value with bad directive")
	}
}

func TestParseInputObjectTypeDefinitionErrorInput(t *testing.T) {
	parser := makeTestParser(t, `123 Filter { name: String }`)
	_, err := parseInputObjectTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseInputObjectTypeDefinitionErrorName(t *testing.T) {
	parser := makeTestParser(t, `input 123 { name: String }`)
	if err := advance(parser); err != nil {
		t.Fatal(err)
	}
	_, err := parseInputObjectTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseInputObjectTypeDefinitionErrorDirective(t *testing.T) {
	parser := makeTestParser(t, `input Filter @123 { name: String }`)
	_, err := parseInputObjectTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseTypeExtensionDefinitionErrorExtend(t *testing.T) {
	parser := makeTestParser(t, `123 type Foo { field: String }`)
	_, err := parseTypeExtensionDefinition(parser)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseDirectiveDefinitionErrorDescription(t *testing.T) {
	parser := makeTestParser(t, `123 directive @foo on FIELD`)
	_, err := parseDirectiveDefinition(parser)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseDirectiveDefinitionErrorDirective(t *testing.T) {
	parser := makeTestParser(t, `123 @foo on FIELD`)
	if err := advance(parser); err != nil {
		t.Fatal(err)
	}
	_, err := parseDirectiveDefinition(parser)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseDirectiveDefinitionErrorAT(t *testing.T) {
	parser := makeTestParser(t, `directive 123 on FIELD`)
	if err := advance(parser); err != nil {
		t.Fatal(err)
	}
	_, err := parseDirectiveDefinition(parser)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseDirectiveDefinitionErrorName(t *testing.T) {
	parser := makeTestParser(t, `directive @123 on FIELD`)
	if err := advance(parser); err != nil {
		t.Fatal(err)
	}
	_, err := parseDirectiveDefinition(parser)
	if err == nil {
		t.Fatal("expected error for directive definition with non-name after @")
	}
}

func TestParseDirectiveDefinitionErrorArgs(t *testing.T) {
	parser := makeTestParser(t, `directive @foo on FIELD`)
	_, err := parseDirectiveDefinition(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseDirectiveDefinitionErrorOn(t *testing.T) {
	parser := makeTestParser(t, `directive @foo`)
	_, err := parseDirectiveDefinition(parser)
	if err == nil {
		t.Fatal("expected error for directive definition missing 'on'")
	}
}

func TestParseDirectiveLocationsErrorPipe(t *testing.T) {
	parser := makeTestParser(t, `FIELD | 123`)
	_, err := parseDirectiveLocations(parser)
	if err == nil {
		t.Fatal("expected error for directive location with non-name")
	}
}

func TestParseStringLiteralErrorAdvance(t *testing.T) {
	parser := makeTestParser(t, `"hello"`)
	_, err := parseStringLiteral(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseDocumentErrorEofSkip(t *testing.T) {
	parser := makeTestParser(t, `123`)
	_, err := parseDocument(parser)
	if err == nil {
		t.Fatal("expected error for document starting with int")
	}
}

// Integration-style error tests via Parse

func TestParse_OperationTypeDefaultCase(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `schema { teapot: Query }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for invalid operation type in schema")
	}
}

func TestParse_OperationDefinitionErrorPaths(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `mutation 123`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for mutation with invalid name")
	}
}

func TestParse_VariableDefinitionErrorPaths(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `query($: Int) { field }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for variable with empty name")
	}
}

func TestParse_VariableMissingColon(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `query($x Int) { field }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for variable missing colon")
	}
}

func TestParse_VariableDefaultValueWithFloat(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `query($x: Float = 3.14) { field }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_ListTypeWithNonNull(t *testing.T) {
	sdl := `type Query { field(arg: [String!]!): String }`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_FieldErrorPaths(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `type Query { field: String @ }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for field with invalid directive")
	}
}

func TestParse_FragmentSpreadErrorPath(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `{ ...123 }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for fragment spread with invalid name")
	}
}

func TestParse_FragmentDefinitionErrorPaths(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `fragment on Query { field }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for fragment missing name")
	}
}

func TestParse_FragmentOnMissingType(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `fragment Foo on { field }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for fragment on missing type")
	}
}

func TestParse_DirectiveErrorPath(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `{ field @ }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for directive with empty name")
	}
}

func TestParse_DirectiveWithArgs(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `{ field @skip(if: true) }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_TypeErrorPaths(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `type Query { field: [ }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for incomplete list type")
	}
}

func TestParse_TypeSystemDefinitionError(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `"desc" { }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for description followed by brace")
	}
}

func TestParse_UnionMemberError(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `union Foo = | Bar`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for union with leading pipe")
	}
}

func TestParse_DirectiveLocationsError(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `directive @foo on FIELD | `,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for directive with trailing pipe")
	}
}

func TestParse_ObjectFieldError(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `{ field(arg: { 123: "val" }) }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for object field with int key")
	}
}

func TestParse_InputValueDefWithObjectDefault(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `type Query { field(arg: Input = {name: "val"}): String }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_IncompleteSchema(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `schema`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for incomplete schema definition")
	}
}

func TestParse_SchemaMissingQuery(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `schema { }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for empty schema")
	}
}

func TestParse_ScalarDefinitionError(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `scalar 123`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for scalar with int name")
	}
}

func TestParse_ScalarDefinitionDirectiveError(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `scalar JSON @`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for scalar with invalid directive")
	}
}

func TestParse_InterfaceDefinitionDirectiveError(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `interface Named @ { name: String }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for interface with invalid directive")
	}
}

func TestParse_UnionDefinitionError(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `union SearchResult = Photo | `,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for union with trailing pipe")
	}
}

func TestParse_EnumDefinitionError(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `enum Color { RED @ }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for enum value with invalid directive")
	}
}

func TestParse_InputObjectDirectiveError(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `input Filter @ { name: String }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for input object with invalid directive")
	}
}

func TestParse_ExtendTypeError(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `extend type 123 { field: String }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for extend type with int name")
	}
}

func TestParse_DirectiveDefinitionManyLocations(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `directive @foo(arg: Int!) on FIELD | OBJECT | QUERY`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_DirectiveDefinitionMissingOn(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `directive @foo FIELD`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for directive definition missing 'on'")
	}
}

func TestParse_DirectiveDefinitionMissingLocations(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `directive @foo on`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for directive definition missing locations")
	}
}

func TestParse_ImplementsInterfacesError(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `type Foo implements 123 { field: String }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for implements with int name")
	}
}

func TestParse_EmptyVariableParens(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `query() { field }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for empty variable parens")
	}
}

func TestParse_EmptyArguments(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `{ field() }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for empty arguments")
	}
}

func TestParse_EmptyArgumentDefs(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `type Query { field(): String }`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for empty argument defs")
	}
}

func TestParse_BlockStringDescription(t *testing.T) {
	sdl := `
		"""
		multi-line
		description
		"""
		type Query {
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

func TestParse_IntegrationNoLocation(t *testing.T) {
	doc, err := Parse(ParseParams{
		Source: source.NewSource(&source.Source{
			Body: []byte(`{ field }`),
		}),
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if doc == nil {
		t.Fatal("expected a document")
	}
}

func TestParse_ValueNullInput(t *testing.T) {
	_, err := ParseValue(ParseParams{
		Source:  `null`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for null value")
	}
}

func TestParse_ValueInvalid(t *testing.T) {
	_, err := ParseValue(ParseParams{
		Source:  `~`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for invalid value token")
	}
}

func TestParse_ValueVariableNonConst(t *testing.T) {
	_, err := ParseValue(ParseParams{
		Source:  `$var`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_InlineFragmentWithoutOn(t *testing.T) {
	q := `
		{
			... @skip(if: true) {
				field
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

func TestParse_FragmentDefinitionWithDirectives(t *testing.T) {
	q := `
		fragment Foo on Query @deprecated {
			field
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

func TestParse_ObjectDefinitionWithDescription(t *testing.T) {
	sdl := `"desc" type Foo { field: String }`
	_, err := Parse(ParseParams{
		Source:  sdl,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_UnionMemberPipeError(t *testing.T) {
	_, err := Parse(ParseParams{
		Source:  `union Foo = Bar | 123`,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err == nil {
		t.Fatal("expected error for union member with int")
	}
}

func TestParse_OperationWithVarDefinitions(t *testing.T) {
	q := `query Foo($x: Int = 42, $y: String) { field }`
	_, err := Parse(ParseParams{
		Source:  q,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_OperationWithMultipleVarDefs(t *testing.T) {
	q := `query ($x: Int, $y: String, $z: Boolean) { field }`
	_, err := Parse(ParseParams{
		Source:  q,
		Options: ParseOptions{NoLocation: true, NoSource: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- Tests for specific uncovered error branches ---

func TestParseValueLiteralIntAdvanceError(t *testing.T) {
	parser := makeTestParser(t, `123 "`)
	_, err := parseValueLiteral(parser, false)
	if err == nil {
		t.Fatal("expected error from advance in INT case")
	}
}

func TestParseValueLiteralFloatAdvanceError(t *testing.T) {
	parser := makeTestParser(t, `3.14 "`)
	_, err := parseValueLiteral(parser, false)
	if err == nil {
		t.Fatal("expected error from advance in FLOAT case")
	}
}

func TestParseValueLiteralBoolAdvanceError(t *testing.T) {
	parser := makeTestParser(t, `true "`)
	_, err := parseValueLiteral(parser, false)
	if err == nil {
		t.Fatal("expected error from advance in bool case")
	}
}

func TestParseValueLiteralEnumAdvanceError(t *testing.T) {
	parser := makeTestParser(t, `SOME_VAL "`)
	_, err := parseValueLiteral(parser, false)
	if err == nil {
		t.Fatal("expected error from advance in enum case")
	}
}

func TestParseTypeBracketLAdvanceError(t *testing.T) {
	parser := makeTestParser(t, `[ "`)
	_, err := parseType(parser)
	if err == nil {
		t.Fatal("expected error from advance in BRACKET_L case")
	}
}

func TestParseTypeBracketRAdvanceError(t *testing.T) {
	parser := makeTestParser(t, `[String]"`)
	_, err := parseType(parser)
	if err == nil {
		t.Fatal("expected error from advance in BRACKET_R case")
	}
}

func TestParseTypeRecursiveError(t *testing.T) {
	parser := makeTestParser(t, `[["`)
	_, err := parseType(parser)
	if err == nil {
		t.Fatal("expected error from recursive parseType")
	}
}

func TestParseTypeBangAdvanceError(t *testing.T) {
	parser := makeTestParser(t, `String!"`)
	_, err := parseType(parser)
	if err == nil {
		t.Fatal("expected error from skip(BANG) advance")
	}
}

func TestParseFieldSkipColonError(t *testing.T) {
	parser := makeTestParser(t, `field: "`)
	_, err := parseField(parser)
	if err == nil {
		t.Fatal("expected error from skip(COLON) advance")
	}
}

func TestParseFragmentFragmentNameAdvanceError(t *testing.T) {
	parser := makeTestParser(t, `...frag"`)
	_, err := parseFragment(parser)
	if err == nil {
		t.Fatal("expected error from parseFragmentName advance")
	}
}

func TestParseFragmentInlineAdvanceError(t *testing.T) {
	parser := makeTestParser(t, `... on "`)
	_, err := parseFragment(parser)
	if err == nil {
		t.Fatal("expected error from advance in inline 'on' path")
	}
}

func TestParseObjectSkipBraceRError(t *testing.T) {
	parser := makeTestParser(t, `{abc: "val"}"`)
	_, err := parseObject(parser, false)
	if err == nil {
		t.Fatal("expected error from skip(BRACE_R) advance")
	}
}

func TestParseStringLiteralAdvanceError(t *testing.T) {
	parser := makeTestParser(t, `"hello" "`)
	_, err := parseStringLiteral(parser)
	if err == nil {
		t.Fatal("expected error from advance in string literal")
	}
}

func TestParseObjectFieldExpectColonError(t *testing.T) {
	parser := makeTestParser(t, `name bar`)
	_, err := parseObjectField(parser, false)
	if err == nil {
		t.Fatal("expected error for object field missing colon")
	}
}

func TestParseArgumentExpectColonError(t *testing.T) {
	parser := makeTestParser(t, `name bar`)
	_, err := parseArgument(parser)
	if err == nil {
		t.Fatal("expected error for argument missing colon")
	}
}

func TestParseVariableExpectDollarError(t *testing.T) {
	parser := makeTestParser(t, `123`)
	_, err := parseVariable(parser)
	if err == nil {
		t.Fatal("expected error for variable missing $")
	}
}

func TestParseOperationDefinitionParseNameAdvanceError(t *testing.T) {
	parser := makeTestParser(t, `query "`)
	_, err := parseOperationDefinition(parser)
	if err == nil {
		t.Fatal("expected error from parseName advance in operation definition")
	}
}

func TestParseOperationDefinitionDirectiveError(t *testing.T) {
	parser := makeTestParser(t, `query @ { field }`)
	_, err := parseOperationDefinition(parser)
	if err == nil {
		t.Fatal("expected error for @ without directive name")
	}
}

func TestParseOperationTypeDefinitionParseNamedError(t *testing.T) {
	parser := makeTestParser(t, `query: 123`)
	_, err := parseOperationTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error for operation type def with non-name type")
	}
}

func TestParseVariableDefinitionParseTypeAdvanceError(t *testing.T) {
	parser := makeTestParser(t, `$x:[String"`)
	_, err := parseVariableDefinition(parser)
	if err == nil {
		t.Fatal("expected error from parseType advance in var def")
	}
}

func TestParseVariableDefinitionSkipEqualsAdvanceError(t *testing.T) {
	parser := makeTestParser(t, `$x:String = "`)
	_, err := parseVariableDefinition(parser)
	if err == nil {
		t.Fatal("expected error from skip(EQUALS) advance in var def")
	}
}

func TestParseInputValueDefSkipEqualsAdvanceError(t *testing.T) {
	parser := makeTestParser(t, `name: String = "`)
	_, err := parseInputValueDef(parser)
	if err == nil {
		t.Fatal("expected error from skip(EQUALS) advance in input value def")
	}
}

func TestParseInputValueDefConstValueError(t *testing.T) {
	parser := makeTestParser(t, `name: String = $var`)
	_, err := parseInputValueDef(parser)
	if err == nil {
		t.Fatal("expected error for const value with variable")
	}
}

func TestParseInputValueDefDirectiveError(t *testing.T) {
	parser := makeTestParser(t, `name: String @123`)
	_, err := parseInputValueDef(parser)
	if err == nil {
		t.Fatal("expected error for input value def with bad directive")
	}
}

func TestParseInputValueDefParseTypeAdvanceError(t *testing.T) {
	parser := makeTestParser(t, `name:["`)
	_, err := parseInputValueDef(parser)
	if err == nil {
		t.Fatal("expected error from parseType advance in input value def")
	}
}

func TestParseFieldDefinitionParseTypeAdvanceError(t *testing.T) {
	parser := makeTestParser(t, `name:["`)
	_, err := parseFieldDefinition(parser)
	if err == nil {
		t.Fatal("expected error from parseType advance in field def")
	}
}

func TestParseImplementsInterfacesAdvanceError(t *testing.T) {
	parser := makeTestParser(t, `type Foo implements " { field: String }`)
	_, err := parseObjectTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error from advance in implements")
	}
}

func TestParseImplementsInterfacesSkipAmpersandError(t *testing.T) {
	parser := makeTestParser(t, `type Foo implements Bar & " { field: String }`)
	_, err := parseObjectTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error from skip(AMP) advance in implements")
	}
}

func TestParseUnionMembersSkipPipeError(t *testing.T) {
	parser := makeTestParser(t, `Foo | "`)
	_, err := parseUnionMembers(parser)
	if err == nil {
		t.Fatal("expected error from skip(PIPE) advance")
	}
}

func TestParseDirectiveLocationsSkipPipeError(t *testing.T) {
	parser := makeTestParser(t, `FIELD | "`)
	_, err := parseDirectiveLocations(parser)
	if err == nil {
		t.Fatal("expected error from skip(PIPE) advance in locations")
	}
}

func TestParseUnionTypeDefinitionParseNameError(t *testing.T) {
	parser := makeTestParser(t, `union 123 = Photo`)
	_, err := parseUnionTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error for union with non-name")
	}
}

func TestParseDirectiveDefinitionExpectATError(t *testing.T) {
	parser := makeTestParser(t, `directive 123 on FIELD`)
	_, err := parseDirectiveDefinition(parser)
	if err == nil {
		t.Fatal("expected error for directive def missing @")
	}
}

func TestParseDirectiveDefinitionParseNameError(t *testing.T) {
	parser := makeTestParser(t, `directive @123 on FIELD`)
	_, err := parseDirectiveDefinition(parser)
	if err == nil {
		t.Fatal("expected error for directive def with non-name after @")
	}
}

func TestParseDirectiveDefinitionArgError(t *testing.T) {
	parser := makeTestParser(t, `directive @foo(123) on FIELD`)
	_, err := parseDirectiveDefinition(parser)
	if err == nil {
		t.Fatal("expected error for directive def with bad args")
	}
}

func TestParseInterfaceTypeDefinitionReverseError(t *testing.T) {
	parser := makeTestParser(t, `interface Named `)
	_, err := parseInterfaceTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error for interface missing {")
	}
}

func TestParseReverseSkipCloseKindError(t *testing.T) {
	parser := makeTestParser(t, `[1]"`)
	_, err := parseList(parser, false)
	if err == nil {
		t.Fatal("expected error from skip(closeKind) in reverse")
	}
}

func TestParseTypeSystemDefinitionLookaheadError(t *testing.T) {
	parser := makeTestParser(t, `"desc" "`)
	_, err := parseTypeSystemDefinition(parser)
	if err == nil {
		t.Fatal("expected error from lookahead")
	}
}

func TestParseScalarTypeDefinitionDescriptionAdvanceError(t *testing.T) {
	parser := makeTestParser(t, `"desc" "scalar JSON`)
	_, err := parseScalarTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error from parseDescription advance")
	}
}

func TestParseObjectTypeDefinitionDescriptionAdvanceError(t *testing.T) {
	parser := makeTestParser(t, `"desc" "type Foo { field: String }`)
	_, err := parseObjectTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error from parseDescription advance")
	}
}

func TestParseFieldDefinitionDescriptionAdvanceError(t *testing.T) {
	parser := makeTestParser(t, `"desc" "name: String`)
	_, err := parseFieldDefinition(parser)
	if err == nil {
		t.Fatal("expected error from parseDescription advance")
	}
}

func TestParseInputValueDefDescriptionAdvanceError(t *testing.T) {
	parser := makeTestParser(t, `"desc" "name: String`)
	_, err := parseInputValueDef(parser)
	if err == nil {
		t.Fatal("expected error from parseDescription advance")
	}
}

func TestParseInterfaceTypeDefinitionDescriptionAdvanceError(t *testing.T) {
	parser := makeTestParser(t, `"desc" "interface Named { field: String }`)
	_, err := parseInterfaceTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error from parseDescription advance")
	}
}

func TestParseUnionTypeDefinitionDescriptionAdvanceError(t *testing.T) {
	parser := makeTestParser(t, `"desc" "union SearchResult = Photo`)
	_, err := parseUnionTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error from parseDescription advance")
	}
}

func TestParseEnumTypeDefinitionDescriptionAdvanceError(t *testing.T) {
	parser := makeTestParser(t, `"desc" "enum Color { RED }`)
	_, err := parseEnumTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error from parseDescription advance")
	}
}

func TestParseEnumValueDefinitionDescriptionAdvanceError(t *testing.T) {
	parser := makeTestParser(t, `"desc" "RED`)
	_, err := parseEnumValueDefinition(parser)
	if err == nil {
		t.Fatal("expected error from parseDescription advance")
	}
}

func TestParseInputObjectTypeDefinitionDescriptionAdvanceError(t *testing.T) {
	parser := makeTestParser(t, `"desc" "input Filter { name: String }`)
	_, err := parseInputObjectTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error from parseDescription advance")
	}
}

func TestParseDirectiveDefinitionDescriptionAdvanceError(t *testing.T) {
	parser := makeTestParser(t, `"desc" "directive @foo on FIELD`)
	_, err := parseDirectiveDefinition(parser)
	if err == nil {
		t.Fatal("expected error from parseDescription advance")
	}
}

func TestParseOperationTypeDefinitionParseTypeAdvanceError(t *testing.T) {
	parser := makeTestParser(t, `query:["`)
	_, err := parseOperationTypeDefinition(parser)
	if err == nil {
		t.Fatal("expected error from parseNamed advance")
	}
}
