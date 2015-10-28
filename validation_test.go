package graphql

import (
	"testing"
)

var someScalarType = NewScalar(ScalarConfig{
	Name: "SomeScalar",
	Serialize: func(value interface{}) interface{} {
		return nil
	},
	ParseValue: func(value interface{}) interface{} {
		return nil
	},
	ParseLiteral: func(valueAST Value) interface{} {
		return nil
	},
})
var someObjectType = NewObject(ObjectConfig{
	Name: "SomeObject",
	Fields: FieldConfigMap{
		"f": &FieldConfig{
			Type: String,
		},
	},
})
var objectWithIsTypeOf = NewObject(ObjectConfig{
	Name: "ObjectWithIsTypeOf",
	IsTypeOf: func(value interface{}, info ResolveInfo) bool {
		return true
	},
	Fields: FieldConfigMap{
		"f": &FieldConfig{
			Type: String,
		},
	},
})
var someUnionType = NewUnion(UnionConfig{
	Name: "SomeUnion",
	ResolveType: func(value interface{}, info ResolveInfo) *Object {
		return nil
	},
	Types: []*Object{
		someObjectType,
	},
})
var someInterfaceType = NewInterface(InterfaceConfig{
	Name: "SomeInterface",
	ResolveType: func(value interface{}, info ResolveInfo) *Object {
		return nil
	},
	Fields: FieldConfigMap{
		"f": &FieldConfig{
			Type: String,
		},
	},
})
var someEnumType = NewEnum(EnumConfig{
	Name: "SomeEnum",
	Values: EnumValueConfigMap{
		"ONLY": &EnumValueConfig{},
	},
})
var someInputObject = NewInputObject(InputObjectConfig{
	Name: "SomeInputObject",
	Fields: InputObjectConfigFieldMap{
		"f": &InputObjectFieldConfig{
			Type:         String,
			DefaultValue: "Hello",
		},
	},
})

func withModifiers(ttypes []Type) []Type {
	res := ttypes
	for _, ttype := range ttypes {
		res = append(res, NewList(ttype))
	}
	for _, ttype := range ttypes {
		res = append(res, NewNonNull(ttype))
	}
	for _, ttype := range ttypes {
		res = append(res, NewNonNull(NewList(ttype)))
	}
	return res
}

var outputTypes = withModifiers([]Type{
	String,
	someScalarType,
	someEnumType,
	someObjectType,
	someUnionType,
	someInterfaceType,
})
var inputTypes = withModifiers([]Type{
	String,
	someScalarType,
	someEnumType,
	someInputObject,
})

func schemaWithFieldType(ttype Output) (Schema, error) {
	return NewSchema(SchemaConfig{
		Query: NewObject(ObjectConfig{
			Name: "Query",
			Fields: FieldConfigMap{
				"f": &FieldConfig{
					Type: ttype,
				},
			},
		}),
	})
}
func schemaWithInputObject(ttype Input) (Schema, error) {
	return NewSchema(SchemaConfig{
		Query: NewObject(ObjectConfig{
			Name: "Query",
			Fields: FieldConfigMap{
				"f": &FieldConfig{
					Type: String,
					Args: FieldConfigArgument{
						"args": &ArgumentConfig{
							Type: ttype,
						},
					},
				},
			},
		}),
	})
}
func schemaWithObjectFieldOfType(fieldType Input) (Schema, error) {

	badObjectType := NewObject(ObjectConfig{
		Name: "BadObject",
		Fields: FieldConfigMap{
			"badField": &FieldConfig{
				Type: fieldType,
			},
		},
	})
	return NewSchema(SchemaConfig{
		Query: NewObject(ObjectConfig{
			Name: "Query",
			Fields: FieldConfigMap{
				"f": &FieldConfig{
					Type: badObjectType,
				},
			},
		}),
	})
}
func schemaWithObjectImplementingType(implementedType *Interface) (Schema, error) {

	badObjectType := NewObject(ObjectConfig{
		Name:       "BadObject",
		Interfaces: []*Interface{implementedType},
		Fields: FieldConfigMap{
			"f": &FieldConfig{
				Type: String,
			},
		},
	})
	return NewSchema(SchemaConfig{
		Query: NewObject(ObjectConfig{
			Name: "Query",
			Fields: FieldConfigMap{
				"f": &FieldConfig{
					Type: badObjectType,
				},
			},
		}),
	})
}
func schemaWithUnionOfType(ttype *Object) (Schema, error) {

	badObjectType := NewUnion(UnionConfig{
		Name: "BadUnion",
		ResolveType: func(value interface{}, info ResolveInfo) *Object {
			return nil
		},
		Types: []*Object{ttype},
	})
	return NewSchema(SchemaConfig{
		Query: NewObject(ObjectConfig{
			Name: "Query",
			Fields: FieldConfigMap{
				"f": &FieldConfig{
					Type: badObjectType,
				},
			},
		}),
	})
}
func schemaWithInterfaceFieldOfType(ttype Type) (Schema, error) {

	badInterfaceType := NewInterface(InterfaceConfig{
		Name: "BadInterface",
		Fields: FieldConfigMap{
			"badField": &FieldConfig{
				Type: ttype,
			},
		},
	})
	return NewSchema(SchemaConfig{
		Query: NewObject(ObjectConfig{
			Name: "Query",
			Fields: FieldConfigMap{
				"f": &FieldConfig{
					Type: badInterfaceType,
				},
			},
		}),
	})
}
func schemaWithArgOfType(ttype Type) (Schema, error) {

	badObject := NewObject(ObjectConfig{
		Name: "BadObject",
		Fields: FieldConfigMap{
			"badField": &FieldConfig{
				Type: String,
				Args: FieldConfigArgument{
					"badArg": &ArgumentConfig{
						Type: ttype,
					},
				},
			},
		},
	})
	return NewSchema(SchemaConfig{
		Query: NewObject(ObjectConfig{
			Name: "Query",
			Fields: FieldConfigMap{
				"f": &FieldConfig{
					Type: badObject,
				},
			},
		}),
	})
}
func schemaWithInputFieldOfType(ttype Type) (Schema, error) {

	badInputObject := NewInputObject(InputObjectConfig{
		Name: "BadInputObject",
		Fields: InputObjectConfigFieldMap{
			"badField": &InputObjectFieldConfig{
				Type: ttype,
			},
		},
	})
	return NewSchema(SchemaConfig{
		Query: NewObject(ObjectConfig{
			Name: "Query",
			Fields: FieldConfigMap{
				"f": &FieldConfig{
					Type: String,
					Args: FieldConfigArgument{
						"badArg": &ArgumentConfig{
							Type: badInputObject,
						},
					},
				},
			},
		}),
	})
}

func TestTypeSystem_SchemaMustHaveObjectRootTypes_AcceptsASchemaWhoseQueryTypeIsAnObjectType(t *testing.T) {
	_, err := NewSchema(SchemaConfig{
		Query: someObjectType,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_SchemaMustHaveObjectRootTypes_AcceptsASchemaWhoseQueryAndMutationTypesAreObjectType(t *testing.T) {
	mutationObject := NewObject(ObjectConfig{
		Name: "Mutation",
		Fields: FieldConfigMap{
			"edit": &FieldConfig{
				Type: String,
			},
		},
	})
	_, err := NewSchema(SchemaConfig{
		Query:    someObjectType,
		Mutation: mutationObject,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_SchemaMustHaveObjectRootTypes_RejectsASchemaWithoutAQueryType(t *testing.T) {
	_, err := NewSchema(SchemaConfig{})
	expectedError := "Schema query must be Object Type but got: nil."
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}

func TestTypeSystem_SchemaMustContainUniquelyNamedTypes_RejectsASchemaWhichRedefinesABuiltInType(t *testing.T) {

	fakeString := NewScalar(ScalarConfig{
		Name: "String",
		Serialize: func(value interface{}) interface{} {
			return nil
		},
	})
	queryType := NewObject(ObjectConfig{
		Name: "Query",
		Fields: FieldConfigMap{
			"normal": &FieldConfig{
				Type: String,
			},
			"fake": &FieldConfig{
				Type: fakeString,
			},
		},
	})
	_, err := NewSchema(SchemaConfig{
		Query: queryType,
	})
	expectedError := `Schema must contain unique named types but contains multiple types named "String".`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
func TestTypeSystem_SchemaMustContainUniquelyNamedTypes_RejectsASchemaWhichDefinesAnObjectTypeTwice(t *testing.T) {

	a := NewObject(ObjectConfig{
		Name: "SameName",
		Fields: FieldConfigMap{
			"f": &FieldConfig{
				Type: String,
			},
		},
	})
	b := NewObject(ObjectConfig{
		Name: "SameName",
		Fields: FieldConfigMap{
			"f": &FieldConfig{
				Type: String,
			},
		},
	})
	queryType := NewObject(ObjectConfig{
		Name: "Query",
		Fields: FieldConfigMap{
			"a": &FieldConfig{
				Type: a,
			},
			"b": &FieldConfig{
				Type: b,
			},
		},
	})
	_, err := NewSchema(SchemaConfig{
		Query: queryType,
	})
	expectedError := `Schema must contain unique named types but contains multiple types named "SameName".`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
func TestTypeSystem_SchemaMustContainUniquelyNamedTypes_RejectsASchemaWhichHaveSameNamedObjectsImplementingAnInterface(t *testing.T) {

	anotherInterface := NewInterface(InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info ResolveInfo) *Object {
			return nil
		},
		Fields: FieldConfigMap{
			"f": &FieldConfig{
				Type: String,
			},
		},
	})
	_ = NewObject(ObjectConfig{
		Name: "BadObject",
		Interfaces: []*Interface{
			anotherInterface,
		},
		Fields: FieldConfigMap{
			"f": &FieldConfig{
				Type: String,
			},
		},
	})
	_ = NewObject(ObjectConfig{
		Name: "BadObject",
		Interfaces: []*Interface{
			anotherInterface,
		},
		Fields: FieldConfigMap{
			"f": &FieldConfig{
				Type: String,
			},
		},
	})
	queryType := NewObject(ObjectConfig{
		Name: "Query",
		Fields: FieldConfigMap{
			"iface": &FieldConfig{
				Type: anotherInterface,
			},
		},
	})
	_, err := NewSchema(SchemaConfig{
		Query: queryType,
	})
	expectedError := `Schema must contain unique named types but contains multiple types named "BadObject".`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}

func TestTypeSystem_ObjectsMustHaveFields_AcceptsAnObjectTypeWithFieldsObject(t *testing.T) {
	_, err := schemaWithFieldType(NewObject(ObjectConfig{
		Name: "SomeObject",
		Fields: FieldConfigMap{
			"f": &FieldConfig{
				Type: String,
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_ObjectsMustHaveFields_RejectsAnObjectTypeWithMissingFields(t *testing.T) {
	badObject := NewObject(ObjectConfig{
		Name: "SomeObject",
	})
	_, err := schemaWithFieldType(badObject)
	expectedError := `SomeObject fields must be an object with field names as keys or a function which return such an object.`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
func TestTypeSystem_ObjectsMustHaveFields_RejectsAnObjectTypeWithIncorrectlyNamedFields(t *testing.T) {
	badObject := NewObject(ObjectConfig{
		Name: "SomeObject",
		Fields: FieldConfigMap{
			"bad-name-with-dashes": &FieldConfig{
				Type: String,
			},
		},
	})
	_, err := schemaWithFieldType(badObject)
	expectedError := `Names must match /^[_a-zA-Z][_a-zA-Z0-9]*$/ but "bad-name-with-dashes" does not.`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
func TestTypeSystem_ObjectsMustHaveFields_RejectsAnObjectTypeWithEmptyFields(t *testing.T) {
	badObject := NewObject(ObjectConfig{
		Name:   "SomeObject",
		Fields: FieldConfigMap{},
	})
	_, err := schemaWithFieldType(badObject)
	expectedError := `SomeObject fields must be an object with field names as keys or a function which return such an object.`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}

func TestTypeSystem_FieldsArgsMustBeProperlyNamed_AcceptsFieldArgsWithValidNames(t *testing.T) {
	_, err := schemaWithFieldType(NewObject(ObjectConfig{
		Name: "SomeObject",
		Fields: FieldConfigMap{
			"goodField": &FieldConfig{
				Type: String,
				Args: FieldConfigArgument{
					"goodArgs": &ArgumentConfig{
						Type: String,
					},
				},
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_FieldsArgsMustBeProperlyNamed_RejectsFieldArgWithInvalidNames(t *testing.T) {
	_, err := schemaWithFieldType(NewObject(ObjectConfig{
		Name: "SomeObject",
		Fields: FieldConfigMap{
			"badField": &FieldConfig{
				Type: String,
				Args: FieldConfigArgument{
					"bad-name-with-dashes": &ArgumentConfig{
						Type: String,
					},
				},
			},
		},
	}))
	expectedError := `Names must match /^[_a-zA-Z][_a-zA-Z0-9]*$/ but "bad-name-with-dashes" does not.`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}

func TestTypeSystem_FieldsArgsMustBeObjects_AcceptsAnObjectTypeWithFieldArgs(t *testing.T) {
	_, err := schemaWithFieldType(NewObject(ObjectConfig{
		Name: "SomeObject",
		Fields: FieldConfigMap{
			"goodField": &FieldConfig{
				Type: String,
				Args: FieldConfigArgument{
					"goodArgs": &ArgumentConfig{
						Type: String,
					},
				},
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTypeSystem_ObjectInterfacesMustBeArray_AcceptsAnObjectTypeWithArrayInterfaces(t *testing.T) {
	anotherInterfaceType := NewInterface(InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info ResolveInfo) *Object {
			return nil
		},
		Fields: FieldConfigMap{
			"f": &FieldConfig{
				Type: String,
			},
		},
	})
	_, err := schemaWithFieldType(NewObject(ObjectConfig{
		Name: "SomeObject",
		Interfaces: (InterfacesThunk)(func() []*Interface {
			return []*Interface{anotherInterfaceType}
		}),
		Fields: FieldConfigMap{
			"f": &FieldConfig{
				Type: String,
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTypeSystem_ObjectInterfacesMustBeArray_AcceptsAnObjectTypeWithInterfacesAsFunctionReturningAnArray(t *testing.T) {
	anotherInterfaceType := NewInterface(InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info ResolveInfo) *Object {
			return nil
		},
		Fields: FieldConfigMap{
			"f": &FieldConfig{
				Type: String,
			},
		},
	})
	_, err := schemaWithFieldType(NewObject(ObjectConfig{
		Name:       "SomeObject",
		Interfaces: []*Interface{anotherInterfaceType},
		Fields: FieldConfigMap{
			"f": &FieldConfig{
				Type: String,
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTypeSystem_UnionTypesMustBeArray_AcceptsAUnionTypeWithArrayTypes(t *testing.T) {
	_, err := schemaWithFieldType(NewUnion(UnionConfig{
		Name: "SomeUnion",
		ResolveType: func(value interface{}, info ResolveInfo) *Object {
			return nil
		},
		Types: []*Object{
			someObjectType,
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_UnionTypesMustBeArray_RejectsAUnionTypeWithoutTypes(t *testing.T) {
	_, err := schemaWithFieldType(NewUnion(UnionConfig{
		Name: "SomeUnion",
		ResolveType: func(value interface{}, info ResolveInfo) *Object {
			return nil
		},
	}))
	expectedError := "Must provide Array of types for Union SomeUnion."
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
func TestTypeSystem_UnionTypesMustBeArray_RejectsAUnionTypeWithEmptyTypes(t *testing.T) {
	_, err := schemaWithFieldType(NewUnion(UnionConfig{
		Name: "SomeUnion",
		ResolveType: func(value interface{}, info ResolveInfo) *Object {
			return nil
		},
		Types: []*Object{},
	}))
	expectedError := "Must provide Array of types for Union SomeUnion."
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}

func TestTypeSystem_InputObjectsMustHaveFields_AcceptsAnInputObjectTypeWithFields(t *testing.T) {
	_, err := schemaWithInputObject(NewInputObject(InputObjectConfig{
		Name: "SomeInputObject",
		Fields: InputObjectConfigFieldMap{
			"f": &InputObjectFieldConfig{
				Type: String,
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTypeSystem_InputObjectsMustHaveFields_AcceptsAnInputObjectTypeWithAFieldFunction(t *testing.T) {
	_, err := schemaWithInputObject(NewInputObject(InputObjectConfig{
		Name: "SomeInputObject",
		Fields: (InputObjectConfigFieldMapThunk)(func() InputObjectConfigFieldMap {
			return InputObjectConfigFieldMap{
				"f": &InputObjectFieldConfig{
					Type: String,
				},
			}
		}),
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTypeSystem_InputObjectsMustHaveFields_RejectsAnInputObjectTypeWithMissingFields(t *testing.T) {
	_, err := schemaWithInputObject(NewInputObject(InputObjectConfig{
		Name: "SomeInputObject",
	}))
	expectedError := "SomeInputObject fields must be an object with field names as keys or a function which return such an object."
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
func TestTypeSystem_InputObjectsMustHaveFields_RejectsAnInputObjectTypeWithEmptyFields(t *testing.T) {
	_, err := schemaWithInputObject(NewInputObject(InputObjectConfig{
		Name:   "SomeInputObject",
		Fields: InputObjectConfigFieldMap{},
	}))
	expectedError := "SomeInputObject fields must be an object with field names as keys or a function which return such an object."
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}

func TestTypeSystem_ObjectTypesMustBeAssertable_AcceptsAnObjectTypeWithAnIsTypeOfFunction(t *testing.T) {
	_, err := schemaWithFieldType(NewObject(ObjectConfig{
		Name: "AnotherObject",
		IsTypeOf: func(value interface{}, info ResolveInfo) bool {
			return true
		},
		Fields: FieldConfigMap{
			"f": &FieldConfig{
				Type: String,
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTypeSystem_InterfaceTypesMustBeResolvable_AcceptsAnInterfaceTypeDefiningResolveType(t *testing.T) {

	anotherInterfaceType := NewInterface(InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info ResolveInfo) *Object {
			return nil
		},
		Fields: FieldConfigMap{
			"f": &FieldConfig{
				Type: String,
			},
		},
	})
	_, err := schemaWithFieldType(NewObject(ObjectConfig{
		Name:       "SomeObject",
		Interfaces: []*Interface{anotherInterfaceType},
		Fields: FieldConfigMap{
			"f": &FieldConfig{
				Type: String,
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_InterfaceTypesMustBeResolvable_AcceptsAnInterfaceWithImplementingTypeDefiningIsTypeOf(t *testing.T) {

	anotherInterfaceType := NewInterface(InterfaceConfig{
		Name: "AnotherInterface",
		Fields: FieldConfigMap{
			"f": &FieldConfig{
				Type: String,
			},
		},
	})
	_, err := schemaWithFieldType(NewObject(ObjectConfig{
		Name:       "SomeObject",
		Interfaces: []*Interface{anotherInterfaceType},
		IsTypeOf: func(value interface{}, info ResolveInfo) bool {
			return true
		},
		Fields: FieldConfigMap{
			"f": &FieldConfig{
				Type: String,
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTypeSystem_InterfaceTypesMustBeResolvable_AcceptsAnInterfaceTypeDefiningResolveTypeWithImplementingTypeDefiningIsTypeOf(t *testing.T) {

	anotherInterfaceType := NewInterface(InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info ResolveInfo) *Object {
			return nil
		},
		Fields: FieldConfigMap{
			"f": &FieldConfig{
				Type: String,
			},
		},
	})
	_, err := schemaWithFieldType(NewObject(ObjectConfig{
		Name:       "SomeObject",
		Interfaces: []*Interface{anotherInterfaceType},
		IsTypeOf: func(value interface{}, info ResolveInfo) bool {
			return true
		},
		Fields: FieldConfigMap{
			"f": &FieldConfig{
				Type: String,
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTypeSystem_UnionTypesMustBeResolvable_AcceptsAUnionTypeDefiningResolveType(t *testing.T) {

	_, err := schemaWithFieldType(NewUnion(UnionConfig{
		Name:  "SomeUnion",
		Types: []*Object{someObjectType},
		ResolveType: func(value interface{}, info ResolveInfo) *Object {
			return nil
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_UnionTypesMustBeResolvable_AcceptsAUnionOfObjectTypesDefiningIsTypeOf(t *testing.T) {

	_, err := schemaWithFieldType(NewUnion(UnionConfig{
		Name:  "SomeUnion",
		Types: []*Object{objectWithIsTypeOf},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_UnionTypesMustBeResolvable_AcceptsAUnionTypeDefiningResolveTypeOfObjectTypesDefiningIsTypeOf(t *testing.T) {

	_, err := schemaWithFieldType(NewUnion(UnionConfig{
		Name:  "SomeUnion",
		Types: []*Object{objectWithIsTypeOf},
		ResolveType: func(value interface{}, info ResolveInfo) *Object {
			return nil
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_UnionTypesMustBeResolvable_RejectsAUnionTypeNotDefiningResolveTypeOfObjectTypesNotDefiningIsTypeOf(t *testing.T) {

	_, err := schemaWithFieldType(NewUnion(UnionConfig{
		Name:  "SomeUnion",
		Types: []*Object{someObjectType},
	}))
	expectedError := `Union Type SomeUnion does not provide a "resolveType" function and ` +
		`possible Type SomeObject does not provide a "isTypeOf" function. ` +
		`There is no way to resolve this possible type during execution.`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}

func TestTypeSystem_ScalarTypesMustBeSerializable_AcceptsAScalarTypeDefiningSerialize(t *testing.T) {

	_, err := schemaWithFieldType(NewScalar(ScalarConfig{
		Name: "SomeScalar",
		Serialize: func(value interface{}) interface{} {
			return nil
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_ScalarTypesMustBeSerializable_RejectsAScalarTypeNotDefiningSerialize(t *testing.T) {

	_, err := schemaWithFieldType(NewScalar(ScalarConfig{
		Name: "SomeScalar",
	}))
	expectedError := `SomeScalar must provide "serialize" function. If this custom Scalar ` +
		`is also used as an input type, ensure "parseValue" and "parseLiteral" ` +
		`functions are also provided.`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
func TestTypeSystem_ScalarTypesMustBeSerializable_AcceptsAScalarTypeDefiningParseValueAndParseLiteral(t *testing.T) {

	_, err := schemaWithFieldType(NewScalar(ScalarConfig{
		Name: "SomeScalar",
		Serialize: func(value interface{}) interface{} {
			return nil
		},
		ParseValue: func(value interface{}) interface{} {
			return nil
		},
		ParseLiteral: func(valueAST Value) interface{} {
			return nil
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_ScalarTypesMustBeSerializable_RejectsAScalarTypeDefiningParseValueButNotParseLiteral(t *testing.T) {

	_, err := schemaWithFieldType(NewScalar(ScalarConfig{
		Name: "SomeScalar",
		Serialize: func(value interface{}) interface{} {
			return nil
		},
		ParseValue: func(value interface{}) interface{} {
			return nil
		},
	}))
	expectedError := `SomeScalar must provide both "parseValue" and "parseLiteral" functions.`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
func TestTypeSystem_ScalarTypesMustBeSerializable_RejectsAScalarTypeDefiningParseLiteralButNotParseValue(t *testing.T) {

	_, err := schemaWithFieldType(NewScalar(ScalarConfig{
		Name: "SomeScalar",
		Serialize: func(value interface{}) interface{} {
			return nil
		},
		ParseLiteral: func(valueAST Value) interface{} {
			return nil
		},
	}))
	expectedError := `SomeScalar must provide both "parseValue" and "parseLiteral" functions.`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}

func TestTypeSystem_EnumTypesMustBeWellDefined_AcceptsAWellDefinedEnumTypeWithEmptyValueDefinition(t *testing.T) {

	_, err := schemaWithFieldType(NewEnum(EnumConfig{
		Name: "SomeEnum",
		Values: EnumValueConfigMap{
			"FOO": &EnumValueConfig{},
			"BAR": &EnumValueConfig{},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_EnumTypesMustBeWellDefined_AcceptsAWellDefinedEnumTypeWithInternalValueDefinition(t *testing.T) {

	_, err := schemaWithFieldType(NewEnum(EnumConfig{
		Name: "SomeEnum",
		Values: EnumValueConfigMap{
			"FOO": &EnumValueConfig{
				Value: 10,
			},
			"BAR": &EnumValueConfig{
				Value: 20,
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_EnumTypesMustBeWellDefined_RejectsAnEnumTypeWithoutValues(t *testing.T) {

	_, err := schemaWithFieldType(NewEnum(EnumConfig{
		Name: "SomeEnum",
	}))
	expectedError := `SomeEnum values must be an object with value names as keys.`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
func TestTypeSystem_EnumTypesMustBeWellDefined_RejectsAnEnumTypeWithEmptyValues(t *testing.T) {

	_, err := schemaWithFieldType(NewEnum(EnumConfig{
		Name:   "SomeEnum",
		Values: EnumValueConfigMap{},
	}))
	expectedError := `SomeEnum values must be an object with value names as keys.`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}

func TestTypeSystem_ObjectFieldsMustHaveOutputTypes_AcceptAnOutputTypeAsAnObjectFieldType(t *testing.T) {
	for _, ttype := range outputTypes {
		_, err := schemaWithObjectFieldOfType(ttype)
		if err != nil {
			t.Fatalf(`unexpected error: %v for type "%v"`, err, ttype)
		}
	}
}
func TestTypeSystem_ObjectFieldsMustHaveOutputTypes_RejectsAnEmptyObjectFieldType(t *testing.T) {
	_, err := schemaWithObjectFieldOfType(nil)
	expectedError := `BadObject.badField field type must be Output Type but got: <nil>.`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}

func TestTypeSystem_ObjectsCanOnlyImplementInterfaces_AcceptsAnObjectImplementingAnInterface(t *testing.T) {
	anotherInterfaceType := NewInterface(InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info ResolveInfo) *Object {
			return nil
		},
		Fields: FieldConfigMap{
			"f": &FieldConfig{
				Type: String,
			},
		},
	})
	_, err := schemaWithObjectImplementingType(anotherInterfaceType)
	if err != nil {
		t.Fatalf(`unexpected error: %v"`, err)
	}
}
func TestTypeSystem_ObjectsCanOnlyImplementInterfaces_RejectsAnObjectImplementingANonInterfaceType(t *testing.T) {
	_, err := schemaWithObjectImplementingType(nil)
	expectedError := `BadObject may only implement Interface types, it cannot implement: <nil>.`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}

func TestTypeSystem_UnionsMustRepresentObjectTypes_AcceptsAUnionOfAnObjectType(t *testing.T) {
	_, err := schemaWithUnionOfType(someObjectType)
	if err != nil {
		t.Fatalf(`unexpected error: %v"`, err)
	}
}
func TestTypeSystem_UnionsMustRepresentObjectTypes_RejectsAUnionOfNonObjectTypes(t *testing.T) {
	_, err := schemaWithUnionOfType(nil)
	expectedError := `BadUnion may only contain Object types, it cannot contain: <nil>.`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}

func TestTypeSystem_InterfaceFieldsMustHaveOutputTypes_AcceptsAnOutputTypeAsAnInterfaceFieldType(t *testing.T) {
	for _, ttype := range outputTypes {
		_, err := schemaWithInterfaceFieldOfType(ttype)
		if err != nil {
			t.Fatalf(`unexpected error: %v for type "%v"`, err, ttype)
		}
	}
}
func TestTypeSystem_InterfaceFieldsMustHaveOutputTypes_RejectsAnEmptyInterfaceFieldType(t *testing.T) {
	_, err := schemaWithInterfaceFieldOfType(nil)
	expectedError := `BadInterface.badField field type must be Output Type but got: <nil>.`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}

func TestTypeSystem_FieldArgumentsMustHaveInputTypes_AcceptsAnInputTypeAsFieldArgType(t *testing.T) {
	for _, ttype := range inputTypes {
		_, err := schemaWithArgOfType(ttype)
		if err != nil {
			t.Fatalf(`unexpected error: %v for type "%v"`, err, ttype)
		}
	}
}
func TestTypeSystem_FieldArgumentsMustHaveInputTypes_RejectsAnEmptyFieldArgType(t *testing.T) {
	_, err := schemaWithArgOfType(nil)
	expectedError := `BadObject.badField(badArg:) argument type must be Input Type but got: <nil>.`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}

func TestTypeSystem_InputObjectFieldsMustHaveInputTypes_AcceptsAnInputTypeAsInputFieldType(t *testing.T) {
	for _, ttype := range inputTypes {
		_, err := schemaWithInputFieldOfType(ttype)
		if err != nil {
			t.Fatalf(`unexpected error: %v for type "%v"`, err, ttype)
		}
	}
}
func TestTypeSystem_InputObjectFieldsMustHaveInputTypes_RejectsAnEmptyInputFieldType(t *testing.T) {
	_, err := schemaWithInputFieldOfType(nil)
	expectedError := `BadInputObject.badField field type must be Input Type but got: <nil>.`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}

func TestTypeSystem_ListMustAcceptGraphQLTypes_AcceptsAnTypeAsItemTypeOfList(t *testing.T) {
	testTypes := withModifiers([]Type{
		String,
		someScalarType,
		someEnumType,
		someObjectType,
		someUnionType,
		someInterfaceType,
	})
	for _, ttype := range testTypes {
		result := NewList(ttype)
		if result.GetError() != nil {
			t.Fatalf(`unexpected error: %v for type "%v"`, result.GetError(), ttype)
		}
	}
}
func TestTypeSystem_ListMustAcceptGraphQLTypes_RejectsANilTypeAsItemTypeOfList(t *testing.T) {
	result := NewList(nil)
	expectedError := `Can only create List of a Type but got: <nil>.`
	if result.GetError() == nil || result.GetError().Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, result.GetError())
	}
}

func TestTypeSystem_NonNullMustAcceptGraphQLTypes_AcceptsAnTypeAsNullableTypeOfNonNull(t *testing.T) {
	nullableTypes := []Type{
		String,
		someScalarType,
		someObjectType,
		someUnionType,
		someInterfaceType,
		someEnumType,
		someInputObject,
		NewList(String),
		NewList(NewNonNull(String)),
	}
	for _, ttype := range nullableTypes {
		result := NewNonNull(ttype)
		if result.GetError() != nil {
			t.Fatalf(`unexpected error: %v for type "%v"`, result.GetError(), ttype)
		}
	}
}
func TestTypeSystem_NonNullMustAcceptGraphQLTypes_RejectsNilAsNonNullableType(t *testing.T) {
	result := NewNonNull(nil)
	expectedError := `Can only create NonNull of a Nullable Type but got: <nil>.`
	if result.GetError() == nil || result.GetError().Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, result.GetError())
	}
}

func TestTypeSystem_ObjectsMustAdhereToInterfaceTheyImplement_AcceptsAnObjectWhichImplementsAnInterface(t *testing.T) {
	anotherInterface := NewInterface(InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info ResolveInfo) *Object {
			return nil
		},
		Fields: FieldConfigMap{
			"field": &FieldConfig{
				Type: String,
				Args: FieldConfigArgument{
					"input": &ArgumentConfig{
						Type: String,
					},
				},
			},
		},
	})
	anotherObject := NewObject(ObjectConfig{
		Name:       "AnotherObject",
		Interfaces: []*Interface{anotherInterface},
		Fields: FieldConfigMap{
			"field": &FieldConfig{
				Type: String,
				Args: FieldConfigArgument{
					"input": &ArgumentConfig{
						Type: String,
					},
				},
			},
		},
	})
	_, err := schemaWithObjectFieldOfType(anotherObject)
	if err != nil {
		t.Fatalf(`unexpected error: %v for type "%v"`, err, anotherObject)
	}
}
func TestTypeSystem_ObjectsMustAdhereToInterfaceTheyImplement_AcceptsAnObjectWhichImplementsAnInterfaceAlongWithMoreFields(t *testing.T) {
	anotherInterface := NewInterface(InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info ResolveInfo) *Object {
			return nil
		},
		Fields: FieldConfigMap{
			"field": &FieldConfig{
				Type: String,
				Args: FieldConfigArgument{
					"input": &ArgumentConfig{
						Type: String,
					},
				},
			},
		},
	})
	anotherObject := NewObject(ObjectConfig{
		Name:       "AnotherObject",
		Interfaces: []*Interface{anotherInterface},
		Fields: FieldConfigMap{
			"field": &FieldConfig{
				Type: String,
				Args: FieldConfigArgument{
					"input": &ArgumentConfig{
						Type: String,
					},
				},
			},
			"anotherfield": &FieldConfig{
				Type: String,
			},
		},
	})
	_, err := schemaWithObjectFieldOfType(anotherObject)
	if err != nil {
		t.Fatalf(`unexpected error: %v for type "%v"`, err, anotherObject)
	}
}
func TestTypeSystem_ObjectsMustAdhereToInterfaceTheyImplement_RejectsAnObjectWhichImplementsAnInterfaceFieldAlongWithMoreArguments(t *testing.T) {
	anotherInterface := NewInterface(InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info ResolveInfo) *Object {
			return nil
		},
		Fields: FieldConfigMap{
			"field": &FieldConfig{
				Type: String,
				Args: FieldConfigArgument{
					"input": &ArgumentConfig{
						Type: String,
					},
				},
			},
		},
	})
	anotherObject := NewObject(ObjectConfig{
		Name:       "AnotherObject",
		Interfaces: []*Interface{anotherInterface},
		Fields: FieldConfigMap{
			"field": &FieldConfig{
				Type: String,
				Args: FieldConfigArgument{
					"input": &ArgumentConfig{
						Type: String,
					},
					"anotherInput": &ArgumentConfig{
						Type: String,
					},
				},
			},
		},
	})
	_, err := schemaWithObjectFieldOfType(anotherObject)
	expectedError := `AnotherInterface.field does not define argument "anotherInput" but AnotherObject.field provides it.`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
func TestTypeSystem_ObjectsMustAdhereToInterfaceTheyImplement_RejectsAnObjectMissingAnInterfaceField(t *testing.T) {
	anotherInterface := NewInterface(InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info ResolveInfo) *Object {
			return nil
		},
		Fields: FieldConfigMap{
			"field": &FieldConfig{
				Type: String,
				Args: FieldConfigArgument{
					"input": &ArgumentConfig{
						Type: String,
					},
				},
			},
		},
	})
	anotherObject := NewObject(ObjectConfig{
		Name:       "AnotherObject",
		Interfaces: []*Interface{anotherInterface},
		Fields: FieldConfigMap{
			"anotherfield": &FieldConfig{
				Type: String,
			},
		},
	})
	_, err := schemaWithObjectFieldOfType(anotherObject)
	expectedError := `"AnotherInterface" expects field "field" but "AnotherObject" does not provide it.`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
func TestTypeSystem_ObjectsMustAdhereToInterfaceTheyImplement_RejectsAnObjectWithAnIncorrectlyTypedInterfaceField(t *testing.T) {
	anotherInterface := NewInterface(InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info ResolveInfo) *Object {
			return nil
		},
		Fields: FieldConfigMap{
			"field": &FieldConfig{
				Type: String,
				Args: FieldConfigArgument{
					"input": &ArgumentConfig{
						Type: String,
					},
				},
			},
		},
	})
	anotherObject := NewObject(ObjectConfig{
		Name:       "AnotherObject",
		Interfaces: []*Interface{anotherInterface},
		Fields: FieldConfigMap{
			"field": &FieldConfig{
				Type: someScalarType,
				Args: FieldConfigArgument{
					"input": &ArgumentConfig{
						Type: String,
					},
				},
			},
		},
	})
	_, err := schemaWithObjectFieldOfType(anotherObject)
	expectedError := `AnotherInterface.field expects type "String" but AnotherObject.field provides type "SomeScalar".`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
func TestTypeSystem_ObjectsMustAdhereToInterfaceTheyImplement_RejectsAnObjectMissingAnInterfaceArgument(t *testing.T) {
	anotherInterface := NewInterface(InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info ResolveInfo) *Object {
			return nil
		},
		Fields: FieldConfigMap{
			"field": &FieldConfig{
				Type: String,
				Args: FieldConfigArgument{
					"input": &ArgumentConfig{
						Type: String,
					},
				},
			},
		},
	})
	anotherObject := NewObject(ObjectConfig{
		Name:       "AnotherObject",
		Interfaces: []*Interface{anotherInterface},
		Fields: FieldConfigMap{
			"field": &FieldConfig{
				Type: String,
			},
		},
	})
	_, err := schemaWithObjectFieldOfType(anotherObject)
	expectedError := `AnotherInterface.field expects argument "input" but AnotherObject.field does not provide it.`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
func TestTypeSystem_ObjectsMustAdhereToInterfaceTheyImplement_RejectsAnObjectWithAnIncorrectlyTypedInterfaceArgument(t *testing.T) {
	anotherInterface := NewInterface(InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info ResolveInfo) *Object {
			return nil
		},
		Fields: FieldConfigMap{
			"field": &FieldConfig{
				Type: String,
				Args: FieldConfigArgument{
					"input": &ArgumentConfig{
						Type: String,
					},
				},
			},
		},
	})
	anotherObject := NewObject(ObjectConfig{
		Name:       "AnotherObject",
		Interfaces: []*Interface{anotherInterface},
		Fields: FieldConfigMap{
			"field": &FieldConfig{
				Type: String,
				Args: FieldConfigArgument{
					"input": &ArgumentConfig{
						Type: someScalarType,
					},
				},
			},
		},
	})
	_, err := schemaWithObjectFieldOfType(anotherObject)
	expectedError := `AnotherInterface.field(input:) expects type "String" but AnotherObject.field(input:) provides type "SomeScalar".`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
func TestTypeSystem_ObjectsMustAdhereToInterfaceTheyImplement_AcceptsAnObjectWithAnEquivalentlyModifiedInterfaceField(t *testing.T) {
	anotherInterface := NewInterface(InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info ResolveInfo) *Object {
			return nil
		},
		Fields: FieldConfigMap{
			"field": &FieldConfig{
				Type: NewNonNull(NewList(String)),
			},
		},
	})
	anotherObject := NewObject(ObjectConfig{
		Name:       "AnotherObject",
		Interfaces: []*Interface{anotherInterface},
		Fields: FieldConfigMap{
			"field": &FieldConfig{
				Type: NewNonNull(NewList(String)),
			},
		},
	})
	_, err := schemaWithObjectFieldOfType(anotherObject)
	if err != nil {
		t.Fatalf(`unexpected error: %v for type "%v"`, err, anotherObject)
	}
}
func TestTypeSystem_ObjectsMustAdhereToInterfaceTheyImplement_RejectsAnObjectWithADifferentlyModifiedInterfaceFieldType(t *testing.T) {
	anotherInterface := NewInterface(InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info ResolveInfo) *Object {
			return nil
		},
		Fields: FieldConfigMap{
			"field": &FieldConfig{
				Type: String,
			},
		},
	})
	anotherObject := NewObject(ObjectConfig{
		Name:       "AnotherObject",
		Interfaces: []*Interface{anotherInterface},
		Fields: FieldConfigMap{
			"field": &FieldConfig{
				Type: NewNonNull(String),
			},
		},
	})
	_, err := schemaWithObjectFieldOfType(anotherObject)
	expectedError := `AnotherInterface.field expects type "String" but AnotherObject.field provides type "String!".`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
