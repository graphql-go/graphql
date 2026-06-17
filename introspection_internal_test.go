package graphql

import (
	"math"
	"testing"
)

func TestAstFromValue_NonNull(t *testing.T) {
	nonNull := NewNonNull(String)
	result := astFromValue("hello", nonNull)
	if result == nil {
		t.Fatal("expected non-nil result for NonNull type")
	}
}

func TestAstFromValue_Nullish(t *testing.T) {
	result := astFromValue(nil, String)
	if result != nil {
		t.Fatal("expected nil result for nullish value")
	}
}

func TestAstFromValue_NilPointer(t *testing.T) {
	var s *string
	result := astFromValue(s, String)
	if result != nil {
		t.Fatal("expected nil result for nil pointer")
	}
}

func TestAstFromValue_Pointer(t *testing.T) {
	val := "hello"
	result := astFromValue(&val, String)
	if result == nil {
		t.Fatal("expected non-nil result for pointer to string")
	}
}

func TestAstFromValue_ListWithSlice(t *testing.T) {
	listType := NewList(String)
	result := astFromValue([]interface{}{"a", "b"}, listType)
	if result == nil {
		t.Fatal("expected non-nil result for list with slice")
	}
}

func TestAstFromValue_ListWithSingleValue(t *testing.T) {
	listType := NewList(String)
	result := astFromValue("single", listType)
	if result == nil {
		t.Fatal("expected non-nil result for single value with list type")
	}
}

func TestAstFromValue_Map(t *testing.T) {
	result := astFromValue(map[string]interface{}{"key": "val"}, String)
	if result == nil {
		t.Fatal("expected non-nil result for map (falls through to fallback)")
	}
}

func TestAstFromValue_Bool(t *testing.T) {
	result := astFromValue(true, Boolean)
	if result == nil {
		t.Fatal("expected non-nil result for bool")
	}
}

func TestAstFromValue_Int(t *testing.T) {
	result := astFromValue(42, String)
	if result == nil {
		t.Fatal("expected non-nil result for int")
	}
}

func TestAstFromValue_IntWithFloatType(t *testing.T) {
	result := astFromValue(42, Float)
	if result == nil {
		t.Fatal("expected non-nil result for int with Float type")
	}
}

func TestAstFromValue_Float32(t *testing.T) {
	result := astFromValue(float32(3.14), Float)
	if result == nil {
		t.Fatal("expected non-nil result for float32")
	}
}

func TestAstFromValue_Float64(t *testing.T) {
	result := astFromValue(float64(3.14), Float)
	if result == nil {
		t.Fatal("expected non-nil result for float64")
	}
}

func TestAstFromValue_StringWithEnum(t *testing.T) {
	enumType := NewEnum(EnumConfig{
		Name: "TestEnum",
		Values: EnumValueConfigMap{
			"VALUE1": &EnumValueConfig{Value: "VALUE1"},
		},
	})
	result := astFromValue("VALUE1", enumType)
	if result == nil {
		t.Fatal("expected non-nil result for string with enum type")
	}
}

func TestAstFromValue_Fallback(t *testing.T) {
	type customType struct {
		name string
	}
	result := astFromValue(customType{name: "test"}, String)
	if result == nil {
		t.Fatal("expected non-nil result for fallback type")
	}
}

func TestIntrospection_TypeKindUnknownType(t *testing.T) {
	kindFieldDef := TypeType.Fields()["kind"]
	result, err := kindFieldDef.Resolve(ResolveParams{
		Source: "not a valid graphql type",
	})
	if err == nil {
		t.Fatal("expected error for unknown type kind")
	}
	if result != nil {
		t.Fatal("expected nil result for unknown type kind")
	}
}

func TestIntrospection_InputValueDefaultValue_NullishArg(t *testing.T) {
	defaultValueFieldDef := InputValueType.Fields()["defaultValue"]
	source := &Argument{
		PrivateName:  "testArg",
		Type:         String,
		DefaultValue: math.NaN(),
	}
	result, err := defaultValueFieldDef.Resolve(ResolveParams{
		Source: source,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Fatal("expected nil result for nullish default value")
	}
}

func TestIntrospection_InputValueDefaultValue_NullishInputField(t *testing.T) {
	defaultValueFieldDef := InputValueType.Fields()["defaultValue"]
	source := &InputObjectField{
		PrivateName:  "testField",
		DefaultValue: math.NaN(),
	}
	_, err := defaultValueFieldDef.Resolve(ResolveParams{
		Source: source,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIntrospection_InputValueDefaultValue_NilInputField(t *testing.T) {
	defaultValueFieldDef := InputValueType.Fields()["defaultValue"]
	source := &InputObjectField{
		PrivateName:  "testField",
		DefaultValue: nil,
	}
	result, err := defaultValueFieldDef.Resolve(ResolveParams{
		Source: source,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Fatal("expected nil result for nil default value")
	}
}

func TestIntrospection_InputValueDefaultValue_UnknownSource(t *testing.T) {
	defaultValueFieldDef := InputValueType.Fields()["defaultValue"]
	result, err := defaultValueFieldDef.Resolve(ResolveParams{
		Source: "some unknown source",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Fatal("expected nil result for unknown source")
	}
}

func TestIntrospection_FieldArgs_NonFieldDefinition(t *testing.T) {
	argsFieldDef := FieldType.Fields()["args"]
	result, err := argsFieldDef.Resolve(ResolveParams{
		Source: "not a field definition",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	args, ok := result.([]interface{})
	if !ok {
		t.Fatal("expected []interface{} result")
	}
	if len(args) != 0 {
		t.Fatal("expected empty args")
	}
}

func TestIntrospection_FieldIsDeprecated_NonFieldDefinition(t *testing.T) {
	isDeprecatedFieldDef := FieldType.Fields()["isDeprecated"]
	result, err := isDeprecatedFieldDef.Resolve(ResolveParams{
		Source: "not a field definition",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != false {
		t.Fatal("expected false for isDeprecated on non-FieldDefinition")
	}
}

func TestIntrospection_FieldDeprecationReason_NonFieldDefinition(t *testing.T) {
	deprecationReasonFieldDef := FieldType.Fields()["deprecationReason"]
	result, err := deprecationReasonFieldDef.Resolve(ResolveParams{
		Source: "not a field definition",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Fatal("expected nil for deprecationReason on non-FieldDefinition")
	}
}

func TestIntrospection_DirectiveOnOperation_WithSubscription(t *testing.T) {
	onOperationFieldDef := DirectiveType.Fields()["onOperation"]
	source := NewDirective(DirectiveConfig{
		Name: "testDirective",
		Locations: []string{
			DirectiveLocationSubscription,
		},
	})
	result, err := onOperationFieldDef.Resolve(ResolveParams{
		Source: source,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != true {
		t.Fatal("expected true for onOperation with subscription location")
	}
}

func TestIntrospection_DirectiveResolvers_NonDirectiveSource(t *testing.T) {
	onOperationFieldDef := DirectiveType.Fields()["onOperation"]
	onFragmentFieldDef := DirectiveType.Fields()["onFragment"]
	onFieldFieldDef := DirectiveType.Fields()["onField"]

	result, err := onOperationFieldDef.Resolve(ResolveParams{
		Source: "not a directive",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != false {
		t.Fatal("expected false for onOperation with non-Directive source")
	}

	result, err = onFragmentFieldDef.Resolve(ResolveParams{
		Source: "not a directive",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != false {
		t.Fatal("expected false for onFragment with non-Directive source")
	}

	result, err = onFieldFieldDef.Resolve(ResolveParams{
		Source: "not a directive",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != false {
		t.Fatal("expected false for onField with non-Directive source")
	}
}

func TestIntrospection_SchemaTypes_NonSchemaSource(t *testing.T) {
	typesFieldDef := SchemaType.Fields()["types"]
	queryTypeFieldDef := SchemaType.Fields()["queryType"]
	directivesFieldDef := SchemaType.Fields()["directives"]

	result, err := typesFieldDef.Resolve(ResolveParams{
		Source: "not a schema",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	types, ok := result.([]Type)
	if !ok || len(types) != 0 {
		t.Fatal("expected empty []Type{} for non-Schema source")
	}

	result, err = queryTypeFieldDef.Resolve(ResolveParams{
		Source: "not a schema",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Fatal("expected nil for queryType with non-Schema source")
	}

	result, err = directivesFieldDef.Resolve(ResolveParams{
		Source: "not a schema",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Fatal("expected nil for directives with non-Schema source")
	}
}

func TestIntrospection_EnumValueIsDeprecated_NonEnumValueDef(t *testing.T) {
	isDeprecatedFieldDef := EnumValueType.Fields()["isDeprecated"]
	result, err := isDeprecatedFieldDef.Resolve(ResolveParams{
		Source: "not an enum value definition",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != false {
		t.Fatal("expected false for isDeprecated on non-EnumValueDefinition")
	}
}

func TestIntrospection_EnumValueDeprecationReason_NonEnumValueDef(t *testing.T) {
	deprecationReasonFieldDef := EnumValueType.Fields()["deprecationReason"]
	result, err := deprecationReasonFieldDef.Resolve(ResolveParams{
		Source: "not an enum value definition",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Fatal("expected nil for deprecationReason on non-EnumValueDefinition")
	}
}

func TestIntrospection_TypeFields_NilObjectSource(t *testing.T) {
	fieldsFieldDef := TypeType.Fields()["fields"]
	var nilObj *Object
	result, err := fieldsFieldDef.Resolve(ResolveParams{
		Source: nilObj,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Fatal("expected nil for nil Object source in fields resolver")
	}
}

func TestIntrospection_TypeFields_NilInterfaceSource(t *testing.T) {
	fieldsFieldDef := TypeType.Fields()["fields"]
	var nilInf *Interface
	result, err := fieldsFieldDef.Resolve(ResolveParams{
		Source: nilInf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Fatal("expected nil for nil Interface source in fields resolver")
	}
}

func TestIntrospection_TypeFields_DeprecatedInterfaceField(t *testing.T) {
	interfaceType := NewInterface(InterfaceConfig{
		Name: "TestInterface",
		Fields: Fields{
			"nonDeprecated": &Field{
				Type: String,
			},
			"deprecated": &Field{
				Type:              String,
				DeprecationReason: "Removed",
			},
		},
	})
	fieldsFieldDef := TypeType.Fields()["fields"]

	// Without includeDeprecated - should filter out deprecated field
	result, err := fieldsFieldDef.Resolve(ResolveParams{
		Source: interfaceType,
		Args: map[string]interface{}{
			"includeDeprecated": false,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fields, ok := result.([]*FieldDefinition)
	if !ok {
		t.Fatal("expected []*FieldDefinition result")
	}
	if len(fields) != 1 {
		t.Fatalf("expected 1 field (non-deprecated), got %d", len(fields))
	}
	if fields[0].Name != "nonDeprecated" {
		t.Fatal("expected nonDeprecated field")
	}
}

func TestIntrospection_TypeInterfaces_NonObjectSource(t *testing.T) {
	interfacesFieldDef := TypeType.Fields()["interfaces"]
	result, err := interfacesFieldDef.Resolve(ResolveParams{
		Source: "not an object",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Fatal("expected nil for non-Object source in interfaces resolver")
	}
}

func TestIntrospection_TypePossibleTypes_NonInterfaceUnionSource(t *testing.T) {
	possibleTypesFieldDef := TypeType.Fields()["possibleTypes"]
	result, err := possibleTypesFieldDef.Resolve(ResolveParams{
		Source: "neither interface nor union",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Fatal("expected nil for non-Interface/Union source in possibleTypes resolver")
	}
}

func TestIntrospection_TypeEnumValues_NonEnumSource(t *testing.T) {
	enumValuesFieldDef := TypeType.Fields()["enumValues"]
	result, err := enumValuesFieldDef.Resolve(ResolveParams{
		Source: "not an enum",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Fatal("expected nil for non-Enum source in enumValues resolver")
	}
}

func TestIntrospection_TypeInputFields_NonInputObjectSource(t *testing.T) {
	inputFieldsFieldDef := TypeType.Fields()["inputFields"]
	result, err := inputFieldsFieldDef.Resolve(ResolveParams{
		Source: "not an input object",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Fatal("expected nil for non-InputObject source in inputFields resolver")
	}
}

func TestIntrospection_TypeMetaFieldDef_NonStringName(t *testing.T) {
	result, err := TypeMetaFieldDef.Resolve(ResolveParams{
		Args: map[string]interface{}{
			"name": 123,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Fatal("expected nil for non-string name in __type resolver")
	}
}

func TestIntrospection_SchemaMutationType_WithMutation(t *testing.T) {
	mutationType := NewObject(ObjectConfig{
		Name: "MutationRoot",
		Fields: Fields{
			"doSomething": &Field{Type: String},
		},
	})
	schema, err := NewSchema(SchemaConfig{
		Query:    NewObject(ObjectConfig{Name: "QueryRoot", Fields: Fields{"f1": &Field{Type: String}}}),
		Mutation: mutationType,
	})
	if err != nil {
		t.Fatalf("unexpected error creating schema: %v", err)
	}
	mutationTypeFieldDef := SchemaType.Fields()["mutationType"]
	result, err := mutationTypeFieldDef.Resolve(ResolveParams{
		Source: schema,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil mutation type result")
	}
}

func TestIntrospection_SchemaSubscriptionType_WithSubscription(t *testing.T) {
	subscriptionType := NewObject(ObjectConfig{
		Name: "SubscriptionRoot",
		Fields: Fields{
			"onEvent": &Field{Type: String},
		},
	})
	schema, err := NewSchema(SchemaConfig{
		Query:        NewObject(ObjectConfig{Name: "QueryRoot", Fields: Fields{"f1": &Field{Type: String}}}),
		Subscription: subscriptionType,
	})
	if err != nil {
		t.Fatalf("unexpected error creating schema: %v", err)
	}
	subscriptionTypeFieldDef := SchemaType.Fields()["subscriptionType"]
	result, err := subscriptionTypeFieldDef.Resolve(ResolveParams{
		Source: schema,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil subscription type result")
	}
}
