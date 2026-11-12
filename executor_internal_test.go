package graphql

import (
	"errors"
	"testing"

	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/ast"
)

// customUnknownASTType implements ast.Type with a non-standard kind,
// used to force typeFromAST's default error branch.
type customUnknownASTType struct {
	kind string
}

func (t *customUnknownASTType) GetKind() string { return t.kind }
func (t *customUnknownASTType) GetLoc() *ast.Location { return nil }
func (t *customUnknownASTType) String() string { return t.kind }

var _ ast.Type = (*customUnknownASTType)(nil)

func testSchemaWithType(extraTypes ...Type) (Schema, *Object) {
	q := NewObject(ObjectConfig{
		Name: "Query",
		Fields: Fields{
			"a": &Field{Type: String},
		},
	})
	s, err := NewSchema(SchemaConfig{
		Query: q,
		Types: extraTypes,
	})
	if err != nil {
		panic(err)
	}
	return s, q
}

// ---------------------------------------------------------------------------
// buildExecutionContext
// ---------------------------------------------------------------------------

func TestBuildExecutionContext_MultipleOperationsNoName(t *testing.T) {
	s, _ := testSchemaWithType()
	doc := ast.NewDocument(&ast.Document{
		Definitions: []ast.Node{
			ast.NewOperationDefinition(&ast.OperationDefinition{
				Operation: ast.OperationTypeQuery,
				SelectionSet: ast.NewSelectionSet(&ast.SelectionSet{
					Selections: []ast.Selection{
						ast.NewField(&ast.Field{Name: ast.NewName(&ast.Name{Value: "a"})}),
					},
				}),
			}),
			ast.NewOperationDefinition(&ast.OperationDefinition{
				Operation: ast.OperationTypeQuery,
				Name:      ast.NewName(&ast.Name{Value: "NamedOp"}),
				SelectionSet: ast.NewSelectionSet(&ast.SelectionSet{
					Selections: []ast.Selection{
						ast.NewField(&ast.Field{Name: ast.NewName(&ast.Name{Value: "a"})}),
					},
				}),
			}),
		},
	})
	_, err := buildExecutionContext(buildExecutionCtxParams{
		Schema: s,
		AST:    doc,
	})
	if err == nil || err.Error() != "Must provide operation name if query contains multiple operations." {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuildExecutionContext_NoOperation(t *testing.T) {
	s, _ := testSchemaWithType()
	doc := ast.NewDocument(&ast.Document{
		Definitions: []ast.Node{
			ast.NewFragmentDefinition(&ast.FragmentDefinition{
				Name: ast.NewName(&ast.Name{Value: "F"}),
				TypeCondition: ast.NewNamed(&ast.Named{
					Name: ast.NewName(&ast.Name{Value: "Query"}),
				}),
				SelectionSet: ast.NewSelectionSet(&ast.SelectionSet{
					Selections: []ast.Selection{
						ast.NewField(&ast.Field{Name: ast.NewName(&ast.Name{Value: "a"})}),
					},
				}),
			}),
		},
	})
	_, err := buildExecutionContext(buildExecutionCtxParams{
		Schema: s,
		AST:    doc,
	})
	if err == nil || err.Error() != "Must provide an operation." {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuildExecutionContext_VariableCoercionError(t *testing.T) {
	s, _ := testSchemaWithType()
	doc := ast.NewDocument(&ast.Document{
		Definitions: []ast.Node{
			ast.NewOperationDefinition(&ast.OperationDefinition{
				Operation: ast.OperationTypeQuery,
				VariableDefinitions: []*ast.VariableDefinition{
					{
						Variable: ast.NewVariable(&ast.Variable{
							Name: ast.NewName(&ast.Name{Value: "x"}),
						}),
						Type: ast.NewNonNull(&ast.NonNull{
							Type: ast.NewNamed(&ast.Named{
								Name: ast.NewName(&ast.Name{Value: "String"}),
							}),
						}),
					},
				},
				SelectionSet: ast.NewSelectionSet(&ast.SelectionSet{
					Selections: []ast.Selection{
						ast.NewField(&ast.Field{Name: ast.NewName(&ast.Name{Value: "a"})}),
					},
				}),
			}),
		},
	})
	_, err := buildExecutionContext(buildExecutionCtxParams{
		Schema: s,
		AST:    doc,
		Args:   map[string]interface{}{},
	})
	if err == nil {
		t.Fatal("expected variable coercion error, got nil")
	}
}

// ---------------------------------------------------------------------------
// getOperationRootType
// ---------------------------------------------------------------------------

func TestGetOperationRootType_NilOperation(t *testing.T) {
	s, _ := testSchemaWithType()
	_, err := getOperationRootType(s, nil)
	if err == nil {
		t.Fatal("expected error for nil operation")
	}
}

func TestGetOperationRootType_UnknownOperation(t *testing.T) {
	s, _ := testSchemaWithType()
	op := ast.NewOperationDefinition(&ast.OperationDefinition{
		Operation: "unknown_op_type",
		SelectionSet: ast.NewSelectionSet(&ast.SelectionSet{
			Selections: []ast.Selection{
				ast.NewField(&ast.Field{Name: ast.NewName(&ast.Name{Value: "a"})}),
			},
		}),
	})
	_, err := getOperationRootType(s, op)
	if err == nil {
		t.Fatal("expected error for unknown operation type")
	}
	if gqlErr, ok := err.(*gqlerrors.Error); ok {
		if gqlErr.Message != "Can only execute queries, mutations and subscription" {
			t.Fatalf("unexpected message: %v", gqlErr.Message)
		}
	} else {
		t.Fatalf("unexpected error type: %T", err)
	}
}

// ---------------------------------------------------------------------------
// collectFields
// ---------------------------------------------------------------------------

func TestCollectFields_NilSelectionSet(t *testing.T) {
	s, _ := testSchemaWithType()
	eCtx := &executionContext{Schema: s}
	result := collectFields(collectFieldsParams{
		ExeContext:   eCtx,
		SelectionSet: nil,
	})
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestCollectFields_InlineFragmentConditionFalse(t *testing.T) {
	// Create a type that exists in the schema but is different from the
	// runtime type, so doesFragmentConditionMatch returns false safely.
	otherType := NewObject(ObjectConfig{
		Name: "OtherType",
		Fields: Fields{
			"a": &Field{Type: String},
		},
	})
	s, q := testSchemaWithType(otherType)
	eCtx := &executionContext{
		Schema:         s,
		Fragments:      map[string]ast.Definition{},
		VariableValues: map[string]interface{}{},
	}
	selSet := ast.NewSelectionSet(&ast.SelectionSet{
		Selections: []ast.Selection{
			ast.NewInlineFragment(&ast.InlineFragment{
				TypeCondition: ast.NewNamed(&ast.Named{
					Name: ast.NewName(&ast.Name{Value: "OtherType"}),
				}),
				SelectionSet: ast.NewSelectionSet(&ast.SelectionSet{
					Selections: []ast.Selection{
						ast.NewField(&ast.Field{Name: ast.NewName(&ast.Name{Value: "a"})}),
					},
				}),
			}),
		},
	})
	fields := collectFields(collectFieldsParams{
		ExeContext:   eCtx,
		RuntimeType:  q,
		SelectionSet: selSet,
		Fields:       map[string][]*ast.Field{},
	})
	if len(fields) != 0 {
		t.Fatalf("expected empty fields, got %v", fields)
	}
}

func TestCollectFields_MissingFragmentSpread(t *testing.T) {
	s, q := testSchemaWithType()
	eCtx := &executionContext{
		Schema:         s,
		Fragments:      map[string]ast.Definition{},
		VariableValues: map[string]interface{}{},
	}
	selSet := ast.NewSelectionSet(&ast.SelectionSet{
		Selections: []ast.Selection{
			ast.NewFragmentSpread(&ast.FragmentSpread{
				Name: ast.NewName(&ast.Name{Value: "MissingFrag"}),
			}),
		},
	})
	fields := collectFields(collectFieldsParams{
		ExeContext:   eCtx,
		RuntimeType:  q,
		SelectionSet: selSet,
		Fields:       map[string][]*ast.Field{},
	})
	if len(fields) != 0 {
		t.Fatalf("expected empty fields, got %v", fields)
	}
}

func TestCollectFields_FragmentSpreadConditionFalse(t *testing.T) {
	// Create a type that exists but differs from runtime type
	otherType := NewObject(ObjectConfig{
		Name: "OtherType",
		Fields: Fields{
			"a": &Field{Type: String},
		},
	})
	s, q := testSchemaWithType(otherType)
	eCtx := &executionContext{
		Schema:         s,
		Fragments:      map[string]ast.Definition{},
		VariableValues: map[string]interface{}{},
	}
	// Add a fragment that won't match (different type from runtime type)
	fragDef := ast.NewFragmentDefinition(&ast.FragmentDefinition{
		Name: ast.NewName(&ast.Name{Value: "F"}),
		TypeCondition: ast.NewNamed(&ast.Named{
			Name: ast.NewName(&ast.Name{Value: "OtherType"}),
		}),
		SelectionSet: ast.NewSelectionSet(&ast.SelectionSet{
			Selections: []ast.Selection{
				ast.NewField(&ast.Field{Name: ast.NewName(&ast.Name{Value: "a"})}),
			},
		}),
	})
	eCtx.Fragments["F"] = fragDef

	selSet := ast.NewSelectionSet(&ast.SelectionSet{
		Selections: []ast.Selection{
			ast.NewFragmentSpread(&ast.FragmentSpread{
				Name: ast.NewName(&ast.Name{Value: "F"}),
			}),
		},
	})
	fields := collectFields(collectFieldsParams{
		ExeContext:   eCtx,
		RuntimeType:  q,
		SelectionSet: selSet,
		Fields:       map[string][]*ast.Field{},
	})
	if len(fields) != 0 {
		t.Fatalf("expected empty fields, got %v", fields)
	}
}

// ---------------------------------------------------------------------------
// shouldIncludeNode – nil directive / nil directive name
// ---------------------------------------------------------------------------

func TestShouldIncludeNode_NilDirective(t *testing.T) {
	eCtx := &executionContext{}
	result := shouldIncludeNode(eCtx, []*ast.Directive{
		ast.NewDirective(&ast.Directive{
			Name: nil,
		}),
	})
	if !result {
		t.Fatal("expected true for nil directive name")
	}
}

// ---------------------------------------------------------------------------
// doesFragmentConditionMatch
// ---------------------------------------------------------------------------

func TestDoesFragmentConditionMatch_FragmentDefinitionNilTypeCondition(t *testing.T) {
	s, q := testSchemaWithType()
	eCtx := &executionContext{Schema: s}
	frag := ast.NewFragmentDefinition(&ast.FragmentDefinition{
		TypeCondition: nil,
	})
	if !doesFragmentConditionMatch(eCtx, frag, q) {
		t.Fatal("expected true for nil type condition")
	}
}

func TestDoesFragmentConditionMatch_FragmentDefinitionRefEqual(t *testing.T) {
	obj := NewObject(ObjectConfig{
		Name: "RefType",
		Fields: Fields{
			"a": &Field{Type: String},
		},
	})
	s, _ := testSchemaWithType(obj)
	eCtx := &executionContext{Schema: s}
	frag := ast.NewFragmentDefinition(&ast.FragmentDefinition{
		TypeCondition: ast.NewNamed(&ast.Named{
			Name: ast.NewName(&ast.Name{Value: "RefType"}),
		}),
	})
	if !doesFragmentConditionMatch(eCtx, frag, obj) {
		t.Fatal("expected true for reference-equal type")
	}
}

func TestDoesFragmentConditionMatch_FragmentDefinitionNameMatch(t *testing.T) {
	schemaType := NewObject(ObjectConfig{
		Name: "SameName",
		Fields: Fields{
			"a": &Field{Type: String},
		},
	})
	runtimeType := NewObject(ObjectConfig{
		Name: "SameName",
		Fields: Fields{
			"a": &Field{Type: String},
		},
	})
	s, _ := testSchemaWithType(schemaType)
	eCtx := &executionContext{Schema: s}
	frag := ast.NewFragmentDefinition(&ast.FragmentDefinition{
		TypeCondition: ast.NewNamed(&ast.Named{
			Name: ast.NewName(&ast.Name{Value: "SameName"}),
		}),
	})
	if !doesFragmentConditionMatch(eCtx, frag, runtimeType) {
		t.Fatal("expected true for name-match")
	}
}

func TestDoesFragmentConditionMatch_FragmentDefinitionInterface(t *testing.T) {
	iface := NewInterface(InterfaceConfig{
		Name: "Char",
		Fields: Fields{
			"name": &Field{Type: String},
		},
	})
	human := NewObject(ObjectConfig{
		Name: "Human",
		Interfaces: []*Interface{iface},
		Fields: Fields{
			"name": &Field{Type: String},
		},
	})
	s, err := NewSchema(SchemaConfig{
		Query: NewObject(ObjectConfig{
			Name: "Query",
			Fields: Fields{
				"a": &Field{Type: String},
			},
		}),
		Types: []Type{human},
	})
	if err != nil {
		t.Fatal(err)
	}
	eCtx := &executionContext{Schema: s}
	frag := ast.NewFragmentDefinition(&ast.FragmentDefinition{
		TypeCondition: ast.NewNamed(&ast.Named{
			Name: ast.NewName(&ast.Name{Value: "Char"}),
		}),
	})
	if !doesFragmentConditionMatch(eCtx, frag, human) {
		t.Fatal("expected true for interface match")
	}
}

func TestDoesFragmentConditionMatch_FragmentDefinitionUnion(t *testing.T) {
	obj := NewObject(ObjectConfig{
		Name: "Result",
		Fields: Fields{
			"a": &Field{Type: String},
		},
		IsTypeOf: func(p IsTypeOfParams) bool { return true },
	})
	searchUnion := NewUnion(UnionConfig{
		Name:  "Search",
		Types: []*Object{obj},
	})
	s, err := NewSchema(SchemaConfig{
		Query: NewObject(ObjectConfig{
			Name: "Query",
			Fields: Fields{
				"a": &Field{Type: String},
			},
		}),
		Types: []Type{obj, searchUnion},
	})
	if err != nil {
		t.Fatal(err)
	}
	eCtx := &executionContext{Schema: s}
	frag := ast.NewFragmentDefinition(&ast.FragmentDefinition{
		TypeCondition: ast.NewNamed(&ast.Named{
			Name: ast.NewName(&ast.Name{Value: "Search"}),
		}),
	})
	if !doesFragmentConditionMatch(eCtx, frag, obj) {
		t.Fatal("expected true for union match")
	}
}

func TestDoesFragmentConditionMatch_FragmentDefinitionTypeNotFound(t *testing.T) {
	s, q := testSchemaWithType()
	eCtx := &executionContext{Schema: s}
	frag := ast.NewFragmentDefinition(&ast.FragmentDefinition{
		TypeCondition: ast.NewNamed(&ast.Named{
			Name: ast.NewName(&ast.Name{Value: "NonExistentType"}),
		}),
	})
	if doesFragmentConditionMatch(eCtx, frag, q) {
		t.Fatal("expected false when type not in schema")
	}
}

func TestDoesFragmentConditionMatch_InlineFragmentTypeNotFound(t *testing.T) {
	s, q := testSchemaWithType()
	eCtx := &executionContext{Schema: s}
	frag := ast.NewInlineFragment(&ast.InlineFragment{
		TypeCondition: ast.NewNamed(&ast.Named{
			Name: ast.NewName(&ast.Name{Value: "NonExistentType"}),
		}),
	})
	if doesFragmentConditionMatch(eCtx, frag, q) {
		t.Fatal("expected false when type not in schema")
	}
}

func TestDoesFragmentConditionMatch_InlineFragmentNilTypeCondition(t *testing.T) {
	s, q := testSchemaWithType()
	eCtx := &executionContext{Schema: s}
	frag := ast.NewInlineFragment(&ast.InlineFragment{
		TypeCondition: nil,
	})
	if !doesFragmentConditionMatch(eCtx, frag, q) {
		t.Fatal("expected true for nil type condition")
	}
}

func TestDoesFragmentConditionMatch_InlineFragmentRefEqual(t *testing.T) {
	obj := NewObject(ObjectConfig{
		Name: "RefType",
		Fields: Fields{
			"a": &Field{Type: String},
		},
	})
	s, _ := testSchemaWithType(obj)
	eCtx := &executionContext{Schema: s}
	frag := ast.NewInlineFragment(&ast.InlineFragment{
		TypeCondition: ast.NewNamed(&ast.Named{
			Name: ast.NewName(&ast.Name{Value: "RefType"}),
		}),
	})
	if !doesFragmentConditionMatch(eCtx, frag, obj) {
		t.Fatal("expected true for reference-equal type")
	}
}

func TestDoesFragmentConditionMatch_InlineFragmentNameMatch(t *testing.T) {
	schemaType := NewObject(ObjectConfig{
		Name: "SameName",
		Fields: Fields{
			"a": &Field{Type: String},
		},
	})
	runtimeType := NewObject(ObjectConfig{
		Name: "SameName",
		Fields: Fields{
			"a": &Field{Type: String},
		},
	})
	s, _ := testSchemaWithType(schemaType)
	eCtx := &executionContext{Schema: s}
	frag := ast.NewInlineFragment(&ast.InlineFragment{
		TypeCondition: ast.NewNamed(&ast.Named{
			Name: ast.NewName(&ast.Name{Value: "SameName"}),
		}),
	})
	if !doesFragmentConditionMatch(eCtx, frag, runtimeType) {
		t.Fatal("expected true for name-match")
	}
}

func TestDoesFragmentConditionMatch_InlineFragmentInterface(t *testing.T) {
	iface := NewInterface(InterfaceConfig{
		Name: "Char",
		Fields: Fields{
			"name": &Field{Type: String},
		},
	})
	human := NewObject(ObjectConfig{
		Name: "Human",
		Interfaces: []*Interface{iface},
		Fields: Fields{
			"name": &Field{Type: String},
		},
	})
	s, err := NewSchema(SchemaConfig{
		Query: NewObject(ObjectConfig{
			Name: "Query",
			Fields: Fields{
				"a": &Field{Type: String},
			},
		}),
		Types: []Type{human},
	})
	if err != nil {
		t.Fatal(err)
	}
	eCtx := &executionContext{Schema: s}
	frag := ast.NewInlineFragment(&ast.InlineFragment{
		TypeCondition: ast.NewNamed(&ast.Named{
			Name: ast.NewName(&ast.Name{Value: "Char"}),
		}),
	})
	if !doesFragmentConditionMatch(eCtx, frag, human) {
		t.Fatal("expected true for interface match")
	}
}

func TestDoesFragmentConditionMatch_InlineFragmentUnion(t *testing.T) {
	obj := NewObject(ObjectConfig{
		Name: "Result",
		Fields: Fields{
			"a": &Field{Type: String},
		},
		IsTypeOf: func(p IsTypeOfParams) bool { return true },
	})
	searchUnion := NewUnion(UnionConfig{
		Name:  "Search",
		Types: []*Object{obj},
	})
	s, err := NewSchema(SchemaConfig{
		Query: NewObject(ObjectConfig{
			Name: "Query",
			Fields: Fields{
				"a": &Field{Type: String},
			},
		}),
		Types: []Type{obj, searchUnion},
	})
	if err != nil {
		t.Fatal(err)
	}
	eCtx := &executionContext{Schema: s}
	frag := ast.NewInlineFragment(&ast.InlineFragment{
		TypeCondition: ast.NewNamed(&ast.Named{
			Name: ast.NewName(&ast.Name{Value: "Search"}),
		}),
	})
	if !doesFragmentConditionMatch(eCtx, frag, obj) {
		t.Fatal("expected true for union match")
	}
}

func TestDoesFragmentConditionMatch_ReturnFalse(t *testing.T) {
	// A fragment with a type condition that exists but the runtime type is
	// a different, unrelated type — should hit return false at line 406.
	objType := NewObject(ObjectConfig{
		Name: "ObjType",
		Fields: Fields{
			"a": &Field{Type: String},
		},
	})
	otherType := NewObject(ObjectConfig{
		Name: "OtherType",
		Fields: Fields{
			"a": &Field{Type: String},
		},
	})
	s, _ := testSchemaWithType(objType, otherType)
	eCtx := &executionContext{Schema: s}
	frag := ast.NewInlineFragment(&ast.InlineFragment{
		TypeCondition: ast.NewNamed(&ast.Named{
			Name: ast.NewName(&ast.Name{Value: "ObjType"}),
		}),
	})
	// Pass otherType as runtime — typeFromAST returns objType,
	// neither ref-equal, nor name-equal, nor Interface, nor Union
	if doesFragmentConditionMatch(eCtx, frag, otherType) {
		t.Fatal("expected false for unrelated types")
	}
}

// ---------------------------------------------------------------------------
// getFieldEntryKey – empty alias and name
// ---------------------------------------------------------------------------

func TestGetFieldEntryKey_Empty(t *testing.T) {
	// Both Alias and Name are nil → returns ""
	field := ast.NewField(&ast.Field{
		Alias: nil,
		Name:  nil,
	})
	if key := getFieldEntryKey(field); key != "" {
		t.Fatalf("expected empty key, got %q", key)
	}

	// Empty Alias, nil Name
	field2 := ast.NewField(&ast.Field{
		Alias: ast.NewName(&ast.Name{Value: ""}),
		Name:  nil,
	})
	if key := getFieldEntryKey(field2); key != "" {
		t.Fatalf("expected empty key, got %q", key)
	}

	// Non-nil Alias with value → returns alias
	field3 := ast.NewField(&ast.Field{
		Alias: ast.NewName(&ast.Name{Value: "alias"}),
		Name:  ast.NewName(&ast.Name{Value: "fieldName"}),
	})
	if key := getFieldEntryKey(field3); key != "alias" {
		t.Fatalf("expected 'alias', got %q", key)
	}

	// No Alias, Name with value → returns name
	field4 := ast.NewField(&ast.Field{
		Name: ast.NewName(&ast.Name{Value: "fieldName"}),
	})
	if key := getFieldEntryKey(field4); key != "fieldName" {
		t.Fatalf("expected 'fieldName', got %q", key)
	}
}

// ---------------------------------------------------------------------------
// defaultResolveTypeFn – no IsTypeOf on possible types
// ---------------------------------------------------------------------------

func TestDefaultResolveTypeFn_NoIsTypeOf(t *testing.T) {
	humanType := NewObject(ObjectConfig{
		Name: "Human",
		Fields: Fields{
			"name": &Field{Type: String},
		},
		// No IsTypeOf
	})
	s, _ := testSchemaWithType(humanType)
	searchUnion := NewUnion(UnionConfig{
		Name:  "Search",
		Types: []*Object{humanType},
	})
	rtParams := ResolveTypeParams{
		Value: map[string]interface{}{},
		Info: ResolveInfo{
			Schema: s,
		},
	}
	result := defaultResolveTypeFn(rtParams, searchUnion)
	if result != nil {
		t.Fatal("expected nil result when no IsTypeOf matches")
	}
}

// ---------------------------------------------------------------------------
// DefaultResolveFn – tag options empty
// ---------------------------------------------------------------------------

type emptyTagStruct struct {
	Val string `json:","`
}

func TestDefaultResolveFn_TagOptionsEmpty(t *testing.T) {
	// Struct with tag `json:","` — split(",") on empty string yields
	// [""] with len=1 > 0, but tOptions[0] != p.Info.FieldName → false.
	// We use a struct field whose json tag has empty first option.
	// The struct field name differs from the json tag, so the name
	// check also fails. This exercises lines 497-499.
	p := ResolveParams{
		Source: emptyTagStruct{Val: "value"},
		Info: ResolveInfo{
			FieldName: "Val",
		},
	}
	val, err := DefaultResolveFn(p)
	if err != nil {
		t.Fatal(err)
	}
	if val != "value" {
		t.Fatalf("expected 'value', got %v", val)
	}
}

// ---------------------------------------------------------------------------
// DefaultResolveFn – reflect.Map with func() interface{} value
// ---------------------------------------------------------------------------

func TestDefaultResolveFn_ReflectMapFunc(t *testing.T) {
	// Source is a typed map whose value type is exactly func() interface{},
	// so it does NOT match map[string]interface{} and instead enters the
	// reflect.Map branch (line 529). The func value exercises lines 533-538.
	source := map[string]func() interface{}{
		"fn": func() interface{} { return "called" },
	}
	p := ResolveParams{
		Source: source,
		Info: ResolveInfo{
			FieldName: "fn",
		},
	}
	val, err := DefaultResolveFn(p)
	if err != nil {
		t.Fatal(err)
	}
	if val != "called" {
		t.Fatalf("expected 'called', got %v", val)
	}
}

// ---------------------------------------------------------------------------
// getFieldDef – nil parentType
// ---------------------------------------------------------------------------

func TestGetFieldDef_NilParentType(t *testing.T) {
	// getFieldDef returns nil when parentType is nil (line 557-559)
	s, _ := testSchemaWithType()
	result := getFieldDef(s, nil, "a")
	if result != nil {
		t.Fatal("expected nil for nil parentType")
	}
}

// ---------------------------------------------------------------------------
// getFieldDef – parentType that is not the schema query type
// ---------------------------------------------------------------------------

func TestGetFieldDef_SchemaMetaFieldDefOnNonQueryType(t *testing.T) {
	// __schema should return SchemaMetaFieldDef only on Query type.
	// When parentType is not Query, it should fall through.
	otherType := NewObject(ObjectConfig{
		Name: "Other",
		Fields: Fields{
			"a": &Field{Type: String},
		},
	})
	s, _ := testSchemaWithType(otherType)
	// otherType is NOT the query type, so __schema should NOT match
	result := getFieldDef(s, otherType, "__schema")
	if result != nil {
		t.Fatalf("expected nil for __schema on non-query parent, got %v", result)
	}
}

// ---------------------------------------------------------------------------
// handleFieldError – NonNull panic (testing handleFieldError coverage)
// ---------------------------------------------------------------------------

func TestHandleFieldError_NonNullPanic(t *testing.T) {
	eCtx := &executionContext{
		Errors: []gqlerrors.FormattedError{},
	}
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for NonNull return type")
		}
	}()
	handleFieldError(
		errors.New("test error"),
		[]ast.Node{},
		&ResponsePath{},
		&NonNull{},
		eCtx,
	)
}

func TestHandleFieldError_Nullable(t *testing.T) {
	eCtx := &executionContext{
		Errors: []gqlerrors.FormattedError{},
	}
	// Should NOT panic for nullable type
	handleFieldError(
		errors.New("test error"),
		[]ast.Node{},
		&ResponsePath{},
		String,
		eCtx,
	)
	if len(eCtx.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(eCtx.Errors))
	}
}

// ---------------------------------------------------------------------------
// DefaultResolveFn – FieldResolver interface
// ---------------------------------------------------------------------------

type testFieldResolver struct {
	val string
}

func (r testFieldResolver) Resolve(p ResolveParams) (interface{}, error) {
	return r.val, nil
}

func TestDefaultResolveFn_FieldResolverInterface(t *testing.T) {
	p := ResolveParams{
		Source: testFieldResolver{val: "resolved"},
		Info: ResolveInfo{
			FieldName: "anything",
		},
	}
	val, err := DefaultResolveFn(p)
	if err != nil {
		t.Fatal(err)
	}
	if val != "resolved" {
		t.Fatalf("expected 'resolved', got %v", val)
	}
}
