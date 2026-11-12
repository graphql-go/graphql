package graphql_test

import (
	"fmt"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

func TestScalarParseLiteral_WithParseLiteralFn(t *testing.T) {
	parseLiteralFn := func(valueAST ast.Value) interface{} {
		if enum, ok := valueAST.(*ast.EnumValue); ok {
			return "parsed:" + enum.Value
		}
		return nil
	}
	scalar := graphql.NewScalar(graphql.ScalarConfig{
		Name:         "CustomScalar",
		Serialize:    func(v interface{}) interface{} { return v },
		ParseValue:   func(v interface{}) interface{} { return v },
		ParseLiteral: parseLiteralFn,
	})
	result := scalar.ParseLiteral(&ast.EnumValue{Value: "FOO"})
	if result != "parsed:FOO" {
		t.Fatalf("expected parsed:FOO, got: %v", result)
	}
}

func TestNewObject_EmptyName(t *testing.T) {
	obj := graphql.NewObject(graphql.ObjectConfig{Name: ""})
	if obj.Error() == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestNewObject_InvalidName(t *testing.T) {
	obj := graphql.NewObject(graphql.ObjectConfig{Name: "123invalid"})
	if obj.Error() == nil {
		t.Fatal("expected error for invalid name")
	}
}

func TestObjectAddFieldConfig_EmptyFieldName(t *testing.T) {
	obj := graphql.NewObject(graphql.ObjectConfig{
		Name: "TestObject",
		Fields: graphql.Fields{
			"a": &graphql.Field{Type: graphql.String},
		},
	})
	obj.AddFieldConfig("", &graphql.Field{Type: graphql.Int})
	if _, ok := obj.Fields()[""]; ok {
		t.Fatal("expected no field added with empty name")
	}
}

func TestObjectAddFieldConfig_NilFieldConfig(t *testing.T) {
	obj := graphql.NewObject(graphql.ObjectConfig{
		Name: "TestObject",
		Fields: graphql.Fields{
			"a": &graphql.Field{Type: graphql.String},
		},
	})
	obj.AddFieldConfig("b", nil)
	if _, ok := obj.Fields()["b"]; ok {
		t.Fatal("expected no field added with nil config")
	}
}

func TestObjectDescription(t *testing.T) {
	obj := graphql.NewObject(graphql.ObjectConfig{
		Name:        "TestObject",
		Description: "an object",
	})
	if obj.Description() != "an object" {
		t.Fatalf("expected 'an object', got: %q", obj.Description())
	}
}

func TestObjectInterfaces_UnknownType(t *testing.T) {
	obj := graphql.NewObject(graphql.ObjectConfig{
		Name:       "TestObject",
		Interfaces: "not an interface type",
	})
	obj.Interfaces()
	if obj.Error() == nil {
		t.Fatal("expected error for unknown interfaces type")
	}
}

func TestDefineInterfaces_WithResolveType(t *testing.T) {
	iface := graphql.NewInterface(graphql.InterfaceConfig{
		Name:        "Named",
		ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object { return nil },
	})
	obj := graphql.NewObject(graphql.ObjectConfig{
		Name: "TestObject",
		Interfaces: []*graphql.Interface{
			iface,
		},
	})
	ifaces := obj.Interfaces()
	if len(ifaces) != 1 {
		t.Fatalf("expected 1 interface, got: %d", len(ifaces))
	}
}

func TestDefineFieldMap_NilField(t *testing.T) {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"a": nil,
			},
		}),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	queryType := schema.QueryType()
	fieldMap := queryType.Fields()
	if _, ok := fieldMap["a"]; ok {
		t.Fatal("expected nil field to be skipped")
	}
}

func TestDefineFieldMap_ArgTypeIsNil(t *testing.T) {
	query := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"field": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: nil,
					},
				},
			},
		},
	})
	_, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: query,
	})
	if err == nil {
		t.Fatal("expected error for nil argument type")
	}
}

func TestArgumentDescriptionStringError(t *testing.T) {
	query := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"field": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "an input",
					},
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: query})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	field := schema.QueryType().Fields()["field"]
	arg := field.Args[0]
	if arg.Name() != "input" {
		t.Fatalf("expected 'input', got: %q", arg.Name())
	}
	if arg.Description() != "an input" {
		t.Fatalf("expected 'an input', got: %q", arg.Description())
	}
	if arg.String() != "input" {
		t.Fatalf("expected 'input', got: %q", arg.String())
	}
	if arg.Error() != nil {
		t.Fatalf("expected nil error, got: %v", arg.Error())
	}
}

func TestNewInterface_EmptyName(t *testing.T) {
	iface := graphql.NewInterface(graphql.InterfaceConfig{Name: ""})
	if iface.Error() == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestInterfaceAddFieldConfig_EmptyFieldName(t *testing.T) {
	iface := graphql.NewInterface(graphql.InterfaceConfig{
		Name: "TestInterface",
		Fields: graphql.Fields{
			"a": &graphql.Field{Type: graphql.String},
		},
	})
	iface.AddFieldConfig("", &graphql.Field{Type: graphql.Int})
	iface.Fields()
}

func TestInterfaceAddFieldConfig_NilFieldConfig(t *testing.T) {
	iface := graphql.NewInterface(graphql.InterfaceConfig{
		Name: "TestInterface",
		Fields: graphql.Fields{
			"a": &graphql.Field{Type: graphql.String},
		},
	})
	iface.AddFieldConfig("b", nil)
	iface.Fields()
}

func TestInterfaceDescription(t *testing.T) {
	iface := graphql.NewInterface(graphql.InterfaceConfig{
		Name:        "TestInterface",
		Description: "an interface",
	})
	if iface.Description() != "an interface" {
		t.Fatalf("expected 'an interface', got: %q", iface.Description())
	}
}

func TestNewUnion_EmptyName(t *testing.T) {
	u := graphql.NewUnion(graphql.UnionConfig{Name: ""})
	if u.Error() == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestUnionDescription(t *testing.T) {
	u := graphql.NewUnion(graphql.UnionConfig{
		Name:        "TestUnion",
		Description: "a union",
	})
	if u.Description() != "a union" {
		t.Fatalf("expected 'a union', got: %q", u.Description())
	}
}

func TestNewEnum_InvalidName(t *testing.T) {
	e := graphql.NewEnum(graphql.EnumConfig{
		Name:   "123invalid",
		Values: graphql.EnumValueConfigMap{"A": &graphql.EnumValueConfig{}},
	})
	if e.Error() == nil {
		t.Fatal("expected error for invalid name")
	}
}

func TestEnumDefineEnumValues_EmptyValueMap(t *testing.T) {
	e := graphql.NewEnum(graphql.EnumConfig{
		Name:   "TestEnum",
		Values: graphql.EnumValueConfigMap{},
	})
	if e.Error() == nil {
		t.Fatal("expected error for empty value map")
	}
}

func TestEnumDefineEnumValues_NilValueConfig(t *testing.T) {
	e := graphql.NewEnum(graphql.EnumConfig{
		Name: "TestEnum",
		Values: graphql.EnumValueConfigMap{
			"FOO": nil,
		},
	})
	if e.Error() == nil {
		t.Fatal("expected error for nil value config")
	}
}

func TestEnumSerialize_PointerValue(t *testing.T) {
	val := "internalFoo"
	e := graphql.NewEnum(graphql.EnumConfig{
		Name: "TestEnum",
		Values: graphql.EnumValueConfigMap{
			"FOO": &graphql.EnumValueConfig{Value: "internalFoo"},
		},
	})
	result := e.Serialize(&val)
	if result != "FOO" {
		t.Fatalf("expected 'FOO', got: %v", result)
	}
}

func TestEnumSerialize_NilPointer(t *testing.T) {
	var val *string
	e := graphql.NewEnum(graphql.EnumConfig{
		Name: "TestEnum",
		Values: graphql.EnumValueConfigMap{
			"FOO": &graphql.EnumValueConfig{Value: "internalFoo"},
		},
	})
	result := e.Serialize(val)
	if result != nil {
		t.Fatalf("expected nil, got: %v", result)
	}
}

func TestEnumParseValue_PointerString(t *testing.T) {
	val := "FOO"
	e := graphql.NewEnum(graphql.EnumConfig{
		Name: "TestEnum",
		Values: graphql.EnumValueConfigMap{
			"FOO": &graphql.EnumValueConfig{Value: "internalFoo"},
		},
	})
	result := e.ParseValue(&val)
	if result != "internalFoo" {
		t.Fatalf(`expected 'internalFoo', got: %v`, result)
	}
}

func TestEnumParseValue_NonStringType(t *testing.T) {
	e := graphql.NewEnum(graphql.EnumConfig{
		Name: "TestEnum",
		Values: graphql.EnumValueConfigMap{
			"FOO": &graphql.EnumValueConfig{Value: "internalFoo"},
		},
	})
	result := e.ParseValue(42)
	if result != nil {
		t.Fatalf("expected nil, got: %v", result)
	}
}

func TestEnumParseValue_UnknownName(t *testing.T) {
	e := graphql.NewEnum(graphql.EnumConfig{
		Name: "TestEnum",
		Values: graphql.EnumValueConfigMap{
			"FOO": &graphql.EnumValueConfig{Value: "internalFoo"},
		},
	})
	result := e.ParseValue("BAR")
	if result != nil {
		t.Fatalf("expected nil, got: %v", result)
	}
}

func TestEnumDescription(t *testing.T) {
	e := graphql.NewEnum(graphql.EnumConfig{
		Name:        "TestEnum",
		Description: "an enum",
		Values: graphql.EnumValueConfigMap{
			"A": &graphql.EnumValueConfig{},
		},
	})
	if e.Description() != "an enum" {
		t.Fatalf("expected 'an enum', got: %q", e.Description())
	}
}

func TestInputObjectFieldGetters(t *testing.T) {
	io := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "TestInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"field": &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "a field",
			},
		},
	})
	field := io.Fields()["field"]
	if field.Name() != "field" {
		t.Fatalf("expected 'field', got: %q", field.Name())
	}
	if field.Description() != "a field" {
		t.Fatalf("expected 'a field', got: %q", field.Description())
	}
	if field.String() != "field" {
		t.Fatalf("expected 'field', got: %q", field.String())
	}
	if field.Error() != nil {
		t.Fatalf("expected nil error, got: %v", field.Error())
	}
}

func TestNewInputObject_EmptyName(t *testing.T) {
	io := graphql.NewInputObject(graphql.InputObjectConfig{Name: ""})
	if io.Error() == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestInputObjectDefineFieldMap_NilFieldConfig(t *testing.T) {
	io := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "TestInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"a": nil,
		},
	})
	fields := io.Fields()
	if _, ok := fields["a"]; ok {
		t.Fatal("expected nil field config to be skipped")
	}
}

func TestInputObjectDefineFieldMap_InvalidFieldName(t *testing.T) {
	io := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "TestInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"123invalid": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"valid": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
		},
	})
	fields := io.Fields()
	if _, ok := fields["123invalid"]; ok {
		t.Fatal("expected invalid field name to be skipped")
	}
	if _, ok := fields["valid"]; !ok {
		t.Fatal("expected valid field to be present")
	}
}

func TestInputObjectAddFieldConfig_WithThunk(t *testing.T) {
	io := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "TestInput",
		Fields: (graphql.InputObjectConfigFieldMapThunk)(func() graphql.InputObjectConfigFieldMap {
			return graphql.InputObjectConfigFieldMap{
				"a": &graphql.InputObjectFieldConfig{Type: graphql.String},
			}
		}),
	})
	io.AddFieldConfig("b", &graphql.InputObjectFieldConfig{Type: graphql.Int})
	if io.Error() == nil {
		t.Fatal("expected error when adding field to thunk-based input object")
	}
}

func TestInputObjectDescription(t *testing.T) {
	io := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        "TestInput",
		Description: "an input",
	})
	if io.Description() != "an input" {
		t.Fatalf("expected 'an input', got: %q", io.Description())
	}
}

func TestInputObjectError(t *testing.T) {
	io := graphql.NewInputObject(graphql.InputObjectConfig{Name: ""})
	if io.Error() == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestListDescription(t *testing.T) {
	l := graphql.NewList(graphql.String)
	if l.Description() != "" {
		t.Fatalf("expected empty description, got: %q", l.Description())
	}
}

func TestListString_NilOfType(t *testing.T) {
	l := &graphql.List{}
	str := fmt.Sprintf("%v", l)
	if str != "" {
		t.Fatalf("expected empty string, got: %q", str)
	}
}

func TestNonNullDescription(t *testing.T) {
	nn := graphql.NewNonNull(graphql.String)
	if nn.Description() != "" {
		t.Fatalf("expected empty description, got: %q", nn.Description())
	}
}

func TestNonNullString_NilOfType(t *testing.T) {
	nn := &graphql.NonNull{}
	str := fmt.Sprintf("%v", nn)
	if str != "" {
		t.Fatalf("expected empty string, got: %q", str)
	}
}

func TestEnumParseLiteral(t *testing.T) {
	e := graphql.NewEnum(graphql.EnumConfig{
		Name: "TestEnum",
		Values: graphql.EnumValueConfigMap{
			"FOO": &graphql.EnumValueConfig{Value: "internalFoo"},
		},
	})
	result := e.ParseLiteral(&ast.EnumValue{Value: "FOO"})
	if result != "internalFoo" {
		t.Fatalf("expected 'internalFoo', got: %v", result)
	}
	result = e.ParseLiteral(&ast.EnumValue{Value: "BAR"})
	if result != nil {
		t.Fatalf("expected nil for unknown value, got: %v", result)
	}
	result = e.ParseLiteral(&ast.StringValue{Value: "FOO"})
	if result != nil {
		t.Fatalf("expected nil for non-enum value, got: %v", result)
	}
}

func TestEnumValues(t *testing.T) {
	e := graphql.NewEnum(graphql.EnumConfig{
		Name: "TestEnum",
		Values: graphql.EnumValueConfigMap{
			"FOO": &graphql.EnumValueConfig{Value: "internalFoo", DeprecationReason: "old", Description: "a value"},
		},
	})
	vals := e.Values()
	if len(vals) != 1 {
		t.Fatalf("expected 1 value, got: %d", len(vals))
	}
	if vals[0].Name != "FOO" {
		t.Fatalf("expected 'FOO', got: %q", vals[0].Name)
	}
	if vals[0].Value != "internalFoo" {
		t.Fatalf("expected 'internalFoo', got: %v", vals[0].Value)
	}
	if vals[0].DeprecationReason != "old" {
		t.Fatalf("expected 'old', got: %q", vals[0].DeprecationReason)
	}
	if vals[0].Description != "a value" {
		t.Fatalf("expected 'a value', got: %q", vals[0].Description)
	}
}

func TestEnumDefaultValue(t *testing.T) {
	e := graphql.NewEnum(graphql.EnumConfig{
		Name: "TestEnum",
		Values: graphql.EnumValueConfigMap{
			"FOO": &graphql.EnumValueConfig{},
		},
	})
	vals := e.Values()
	if vals[0].Value != "FOO" {
		t.Fatalf("expected default value 'FOO', got: %v", vals[0].Value)
	}
}

func TestEnumSerialize_ValueLookupMiss(t *testing.T) {
	e := graphql.NewEnum(graphql.EnumConfig{
		Name: "TestEnum",
		Values: graphql.EnumValueConfigMap{
			"FOO": &graphql.EnumValueConfig{Value: "internalFoo"},
		},
	})
	result := e.Serialize("nonexistent")
	if result != nil {
		t.Fatalf("expected nil, got: %v", result)
	}
}

func TestObjectFieldsThunk(t *testing.T) {
	obj := graphql.NewObject(graphql.ObjectConfig{
		Name: "Test",
		Fields: (graphql.FieldsThunk)(func() graphql.Fields {
			return graphql.Fields{
				"a": &graphql.Field{Type: graphql.String},
			}
		}),
	})
	fields := obj.Fields()
	if len(fields) != 1 {
		t.Fatalf("expected 1 field, got: %d", len(fields))
	}
}

func TestObjectInterfacesThunk(t *testing.T) {
	iface := graphql.NewInterface(graphql.InterfaceConfig{Name: "Named"})
	obj := graphql.NewObject(graphql.ObjectConfig{
		Name: "TestObject",
		Interfaces: (graphql.InterfacesThunk)(func() []*graphql.Interface {
			return []*graphql.Interface{iface}
		}),
		IsTypeOf: func(p graphql.IsTypeOfParams) bool { return true },
	})
	ifaces := obj.Interfaces()
	if len(ifaces) != 1 {
		t.Fatalf("expected 1 interface, got: %d", len(ifaces))
	}
}

func TestUnionTypes_Thunk(t *testing.T) {
	obj := graphql.NewObject(graphql.ObjectConfig{Name: "SomeObject"})
	u := graphql.NewUnion(graphql.UnionConfig{
		Name: "TestUnion",
		Types: (graphql.UnionTypesThunk)(func() []*graphql.Object {
			return []*graphql.Object{obj}
		}),
		ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object { return nil },
	})
	types := u.Types()
	if len(types) != 1 {
		t.Fatalf("expected 1 type, got: %d", len(types))
	}
}

func TestUnionTypes_Nil(t *testing.T) {
	u := graphql.NewUnion(graphql.UnionConfig{
		Name: "TestUnion",
	})
	u.Types()
	if u.Error() == nil {
		t.Fatal("expected error for missing types")
	}
}

func TestDefineUnionTypes_WithResolveType(t *testing.T) {
	obj := graphql.NewObject(graphql.ObjectConfig{
		Name:    "SomeObject",
		IsTypeOf: func(p graphql.IsTypeOfParams) bool { return true },
	})
	u := graphql.NewUnion(graphql.UnionConfig{
		Name:        "TestUnion",
		Types:       []*graphql.Object{obj},
		ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object { return nil },
	})
	types := u.Types()
	if len(types) != 1 {
		t.Fatalf("expected 1 type, got: %d", len(types))
	}
}

func TestDefineUnionTypes_NoResolveTypeNoIsTypeOf(t *testing.T) {
	obj := graphql.NewObject(graphql.ObjectConfig{
		Name: "SomeObject",
	})
	u := graphql.NewUnion(graphql.UnionConfig{
		Name:  "TestUnion",
		Types: []*graphql.Object{obj},
	})
	u.Types()
	if u.Error() == nil {
		t.Fatal("expected error when neither ResolveType nor IsTypeOf is provided")
	}
}

func TestDefineFieldMap_NilFieldType(t *testing.T) {
	query := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"field": &graphql.Field{
				Type: nil,
			},
		},
	})
	_, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: query,
	})
	if err == nil {
		t.Fatal("expected error for nil field type")
	}
}

func TestDefineFieldMap_FieldTypeWithError(t *testing.T) {
	badType := graphql.NewScalar(graphql.ScalarConfig{
		Name:      "",
		Serialize: func(v interface{}) interface{} { return v },
	})
	query := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"field": &graphql.Field{
				Type: badType,
			},
		},
	})
	_, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: query,
	})
	if err == nil {
		t.Fatal("expected error for field type with error")
	}
}

func TestDefineFieldMap_InvalidFieldName(t *testing.T) {
	query := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"123invalid": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
	_, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: query,
	})
	if err == nil {
		t.Fatal("expected error for invalid field name")
	}
}

func TestDefineFieldMap_InvalidArgName(t *testing.T) {
	query := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"field": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"123invalid": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
			},
		},
	})
	_, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: query,
	})
	if err == nil {
		t.Fatal("expected error for invalid argument name")
	}
}

func TestDefineFieldMap_NilArg(t *testing.T) {
	query := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"field": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"input": nil,
				},
			},
		},
	})
	_, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: query,
	})
	if err == nil {
		t.Fatal("expected error for nil argument")
	}
}

func TestObjectName(t *testing.T) {
	obj := graphql.NewObject(graphql.ObjectConfig{Name: "Foo"})
	if obj.Name() != "Foo" {
		t.Fatalf("expected 'Foo', got: %q", obj.Name())
	}
}

func TestObjectString(t *testing.T) {
	obj := graphql.NewObject(graphql.ObjectConfig{Name: "Foo"})
	if obj.String() != "Foo" {
		t.Fatalf("expected 'Foo', got: %q", obj.String())
	}
}

func TestInterfaceName(t *testing.T) {
	iface := graphql.NewInterface(graphql.InterfaceConfig{Name: "Bar"})
	if iface.Name() != "Bar" {
		t.Fatalf("expected 'Bar', got: %q", iface.Name())
	}
}

func TestInterfaceString(t *testing.T) {
	iface := graphql.NewInterface(graphql.InterfaceConfig{Name: "Bar"})
	if iface.String() != "Bar" {
		t.Fatalf("expected 'Bar', got: %q", iface.String())
	}
}

func TestUnionName(t *testing.T) {
	u := graphql.NewUnion(graphql.UnionConfig{Name: "Baz"})
	if u.Name() != "Baz" {
		t.Fatalf("expected 'Baz', got: %q", u.Name())
	}
}

func TestUnionString(t *testing.T) {
	u := graphql.NewUnion(graphql.UnionConfig{Name: "Baz"})
	if u.String() != "Baz" {
		t.Fatalf("expected 'Baz', got: %q", u.String())
	}
}

func TestEnumString(t *testing.T) {
	e := graphql.NewEnum(graphql.EnumConfig{
		Name:   "Color",
		Values: graphql.EnumValueConfigMap{"RED": &graphql.EnumValueConfig{}},
	})
	if e.String() != "Color" {
		t.Fatalf("expected 'Color', got: %q", e.String())
	}
}

func TestInputObjectString(t *testing.T) {
	io := graphql.NewInputObject(graphql.InputObjectConfig{Name: "Input"})
	if io.String() != "Input" {
		t.Fatalf("expected 'Input', got: %q", io.String())
	}
}

func TestInputObjectName(t *testing.T) {
	io := graphql.NewInputObject(graphql.InputObjectConfig{Name: "Input"})
	if io.Name() != "Input" {
		t.Fatalf("expected 'Input', got: %q", io.Name())
	}
}

func TestListName(t *testing.T) {
	l := graphql.NewList(graphql.String)
	if l.Name() != "[String]" {
		t.Fatalf("expected '[String]', got: %q", l.Name())
	}
}

func TestListError(t *testing.T) {
	l := graphql.NewList(nil)
	if l.Error() == nil {
		t.Fatal("expected error for nil OfType")
	}
}

func TestNonNullName(t *testing.T) {
	nn := graphql.NewNonNull(graphql.String)
	if nn.Name() != "String!" {
		t.Fatalf("expected 'String!', got: %q", nn.Name())
	}
}

func TestNonNullError(t *testing.T) {
	nn := graphql.NewNonNull(nil)
	if nn.Error() == nil {
		t.Fatal("expected error for nil OfType")
	}
}

func TestInterfaceFields_Thunk(t *testing.T) {
	iface := graphql.NewInterface(graphql.InterfaceConfig{
		Name: "TestInterface",
		Fields: (graphql.FieldsThunk)(func() graphql.Fields {
			return graphql.Fields{
				"a": &graphql.Field{Type: graphql.String},
			}
		}),
	})
	fields := iface.Fields()
	if len(fields) != 1 {
		t.Fatalf("expected 1 field, got: %d", len(fields))
	}
}

func TestInputObjectDefineFieldMap_FieldTypeNil(t *testing.T) {
	io := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "TestInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"field": &graphql.InputObjectFieldConfig{
				Type: nil,
			},
		},
	})
	io.Fields()
	if io.Error() == nil {
		t.Fatal("expected error for nil field type")
	}
}

func TestInputObjectDefineFieldMap_EmptyFieldMap(t *testing.T) {
	io := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "TestInput",
		Fields: graphql.InputObjectConfigFieldMap{},
	})
	io.Fields()
	if io.Error() == nil {
		t.Fatal("expected error for empty field map")
	}
}

func TestInputObjectFields_WithThunk(t *testing.T) {
	io := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "TestInput",
		Fields: (graphql.InputObjectConfigFieldMapThunk)(func() graphql.InputObjectConfigFieldMap {
			return graphql.InputObjectConfigFieldMap{
				"a": &graphql.InputObjectFieldConfig{Type: graphql.String},
			}
		}),
	})
	fields := io.Fields()
	if len(fields) != 1 {
		t.Fatalf("expected 1 field, got: %d", len(fields))
	}
}

func TestScalarParseLiteral_NilFn(t *testing.T) {
	scalar := graphql.NewScalar(graphql.ScalarConfig{
		Name:      "Test",
		Serialize: func(v interface{}) interface{} { return v },
	})
	result := scalar.ParseLiteral(&ast.StringValue{Value: "hello"})
	if result != nil {
		t.Fatalf("expected nil, got: %v", result)
	}
}

func TestDefineFieldMap_EmptyFields(t *testing.T) {
	query := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Query",
		Fields: graphql.Fields{},
	})
	_, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: query,
	})
	if err == nil {
		t.Fatal("expected error for empty field map")
	}
}

func TestObjectError(t *testing.T) {
	obj := graphql.NewObject(graphql.ObjectConfig{Name: ""})
	if obj.Error() == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestInterfaceError(t *testing.T) {
	iface := graphql.NewInterface(graphql.InterfaceConfig{Name: ""})
	if iface.Error() == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestUnionError(t *testing.T) {
	u := graphql.NewUnion(graphql.UnionConfig{Name: ""})
	if u.Error() == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestEnumError(t *testing.T) {
	e := graphql.NewEnum(graphql.EnumConfig{
		Name:   "123bad",
		Values: graphql.EnumValueConfigMap{"A": &graphql.EnumValueConfig{}},
	})
	if e.Error() == nil {
		t.Fatal("expected error for invalid name")
	}
}

func TestNewEnum_EmptyName(t *testing.T) {
	e := graphql.NewEnum(graphql.EnumConfig{
		Name:   "",
		Values: graphql.EnumValueConfigMap{"A": &graphql.EnumValueConfig{}},
	})
	if e.Error() == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestNewInterface_InvalidName(t *testing.T) {
	iface := graphql.NewInterface(graphql.InterfaceConfig{Name: "123invalid"})
	if iface.Error() == nil {
		t.Fatal("expected error for invalid name")
	}
}

func TestNewUnion_InvalidName(t *testing.T) {
	u := graphql.NewUnion(graphql.UnionConfig{Name: "123invalid"})
	if u.Error() == nil {
		t.Fatal("expected error for invalid name")
	}
}

func TestEnumDefineEnumValues_InvalidValueName(t *testing.T) {
	e := graphql.NewEnum(graphql.EnumConfig{
		Name: "TestEnum",
		Values: graphql.EnumValueConfigMap{
			"123invalid": &graphql.EnumValueConfig{},
		},
	})
	if e.Error() == nil {
		t.Fatal("expected error for invalid value name")
	}
}

func TestInputObjectAddFieldConfig_EmptyName(t *testing.T) {
	io := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "TestInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"a": &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	})
	io.AddFieldConfig("", &graphql.InputObjectFieldConfig{Type: graphql.Int})
	if io.Error() != nil {
		t.Fatalf("unexpected error: %v", io.Error())
	}
}

func TestInputObjectAddFieldConfig_NilConfig(t *testing.T) {
	io := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "TestInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"a": &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	})
	io.AddFieldConfig("b", nil)
	if io.Error() != nil {
		t.Fatalf("unexpected error: %v", io.Error())
	}
}

func TestListNewList_NilType(t *testing.T) {
	l := graphql.NewList(nil)
	if l.Error() == nil {
		t.Fatal("expected error for nil OfType")
	}
}

func TestNonNullNewNonNull_NilType(t *testing.T) {
	nn := graphql.NewNonNull(nil)
	if nn.Error() == nil {
		t.Fatal("expected error for nil OfType")
	}
}

func TestNonNullNewNonNull_WrapsNonNull(t *testing.T) {
	nn := graphql.NewNonNull(graphql.NewNonNull(graphql.String))
	if nn.Error() == nil {
		t.Fatal("expected error for wrapping NonNull in NonNull")
	}
}

func TestScalarParseValue_WithParseValueFn(t *testing.T) {
	scalar := graphql.NewScalar(graphql.ScalarConfig{
		Name:       "Custom",
		Serialize:  func(v interface{}) interface{} { return v },
		ParseValue: func(v interface{}) interface{} { return "parsed:" + v.(string) },
		ParseLiteral: func(v ast.Value) interface{} { return nil },
	})
	result := scalar.ParseValue("hello")
	if result != "parsed:hello" {
		t.Fatalf("expected 'parsed:hello', got: %v", result)
	}
}

func TestScalarMissingSerialize(t *testing.T) {
	scalar := graphql.NewScalar(graphql.ScalarConfig{
		Name: "Test",
	})
	if scalar.Error() == nil {
		t.Fatal("expected error when Serialize is nil")
	}
}

func TestScalarMissingBothParseFn(t *testing.T) {
	scalar := graphql.NewScalar(graphql.ScalarConfig{
		Name:       "Test",
		Serialize:  func(v interface{}) interface{} { return v },
		ParseValue: func(v interface{}) interface{} { return v },
	})
	if scalar.Error() == nil {
		t.Fatal("expected error when only ParseValue is provided without ParseLiteral")
	}
}

func TestScalarSerialize(t *testing.T) {
	serializeFn := func(value interface{}) interface{} {
		if intVal, ok := value.(int); ok && intVal%2 == 1 {
			return intVal
		}
		return nil
	}

	tests := []struct {
		name     string
		config   graphql.ScalarConfig
		input    interface{}
		expected interface{}
	}{
		{
			name: "Serialize with custom function should call function",
			config: graphql.ScalarConfig{
				Name:      "OddInt",
				Serialize: serializeFn,
			},
			input:    3,
			expected: 3,
		},
		{
			name: "Serialize with custom function returns nil for even number",
			config: graphql.ScalarConfig{
				Name:      "OddInt",
				Serialize: serializeFn,
			},
			input:    4,
			expected: nil,
		},
		{
			name: "Serialize with nil function should return input value",
			config: graphql.ScalarConfig{
				Name:      "TestScalar",
				Serialize: nil,
			},
			input:    "test",
			expected: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scalar := graphql.NewScalar(tt.config)
			result := scalar.Serialize(tt.input)
			if result != tt.expected {
				t.Errorf("Serialize(%v) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestScalarParseValue(t *testing.T) {
	parseValueFn := func(value interface{}) interface{} {
		if strVal, ok := value.(string); ok {
			return strVal
		}
		return nil
	}

	tests := []struct {
		name     string
		config   graphql.ScalarConfig
		input    interface{}
		expected interface{}
	}{
		{
			name: "ParseValue with custom function should parse string",
			config: graphql.ScalarConfig{
				Name:       "StringScalar",
				Serialize:  func(v interface{}) interface{} { return v },
				ParseValue: parseValueFn,
			},
			input:    "hello",
			expected: "hello",
		},
		{
			name: "ParseValue with custom function returns given value",
			config: graphql.ScalarConfig{
				Name:       "StringScalar",
				Serialize:  func(v interface{}) interface{} { return v },
				ParseValue: parseValueFn,
			},
			input:    123,
			expected: 123,
		},
		{
			name: "ParseValue with nil function should return input value as-is",
			config: graphql.ScalarConfig{
				Name:      "TestScalar",
				Serialize: func(v interface{}) interface{} { return v },
			},
			input:    "test",
			expected: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scalar := graphql.NewScalar(tt.config)
			result := scalar.ParseValue(tt.input)
			if result != tt.expected {
				t.Errorf("ParseValue(%v) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestScalarDescription(t *testing.T) {
	tests := []struct {
		name        string
		description string
		expected    string
	}{
		{
			name:        "Description should return provided description",
			description: "A custom scalar type",
			expected:    "A custom scalar type",
		},
		{
			name:        "Description should return empty string when not provided",
			description: "",
			expected:    "",
		},
		{
			name:        "Description should return long description",
			description: "This is a longer description for a scalar type that can span multiple lines",
			expected:    "This is a longer description for a scalar type that can span multiple lines",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scalar := graphql.NewScalar(graphql.ScalarConfig{
				Name:        "TestScalar",
				Description: tt.description,
				Serialize:   func(v interface{}) interface{} { return v },
			})
			result := scalar.Description()
			if result != tt.expected {
				t.Errorf("Description() = %q; want %q", result, tt.expected)
			}
		})
	}
}
