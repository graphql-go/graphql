package types_test

import (
	"testing"

	"github.com/chris-ramon/graphql/language/ast"
	"github.com/chris-ramon/graphql/types"
)

var someScalarType = types.NewScalar(types.ScalarConfig{
	Name: "SomeScalar",
	Serialize: func(value interface{}) interface{} {
		return nil
	},
	ParseValue: func(value interface{}) interface{} {
		return nil
	},
	ParseLiteral: func(valueAST ast.Value) interface{} {
		return nil
	},
})
var someObjectType = types.NewObject(types.ObjectConfig{
	Name: "SomeObject",
	Fields: types.FieldConfigMap{
		"f": &types.FieldConfig{
			Type: types.String,
		},
	},
})
var objectWithIsTypeOf = types.NewObject(types.ObjectConfig{
	Name: "ObjectWithIsTypeOf",
	IsTypeOf: func(value interface{}, info types.ResolveInfo) bool {
		return true
	},
	Fields: types.FieldConfigMap{
		"f": &types.FieldConfig{
			Type: types.String,
		},
	},
})
var someUnionType = types.NewUnion(types.UnionConfig{
	Name: "SomeUnion",
	ResolveType: func(value interface{}, info types.ResolveInfo) *types.Object {
		return nil
	},
	Types: []*types.Object{
		someObjectType,
	},
})
var someInterfaceType = types.NewInterface(types.InterfaceConfig{
	Name: "SomeInterface",
	ResolveType: func(value interface{}, info types.ResolveInfo) *types.Object {
		return nil
	},
	Fields: types.FieldConfigMap{
		"f": &types.FieldConfig{
			Type: types.String,
		},
	},
})
var someEnumType = types.NewEnum(types.EnumConfig{
	Name: "SomeEnum",
	Values: types.EnumValueConfigMap{
		"ONLY": &types.EnumValueConfig{},
	},
})
var someInputObject = types.NewInputObject(types.InputObjectConfig{
	Name: "SomeInputObject",
	Fields: types.InputObjectConfigFieldMap{
		"f": &types.InputObjectFieldConfig{
			Type:         types.String,
			DefaultValue: "Hello",
		},
	},
})

func withModifiers(ttypes []types.Type) []types.Type {
	res := ttypes
	for _, ttype := range ttypes {
		res = append(res, types.NewList(ttype))
	}
	for _, ttype := range ttypes {
		res = append(res, types.NewNonNull(ttype))
	}
	for _, ttype := range ttypes {
		res = append(res, types.NewNonNull(types.NewList(ttype)))
	}
	return res
}

var outputTypes = withModifiers([]types.Type{
	types.String,
	someScalarType,
	someEnumType,
	someObjectType,
	someUnionType,
	someInterfaceType,
})
var inputTypes = withModifiers([]types.Type{
	types.String,
	someScalarType,
	someEnumType,
	someInputObject,
})

func schemaWithFieldType(ttype types.Output) (types.Schema, error) {
	return types.NewSchema(types.SchemaConfig{
		Query: types.NewObject(types.ObjectConfig{
			Name: "Query",
			Fields: types.FieldConfigMap{
				"f": &types.FieldConfig{
					Type: ttype,
				},
			},
		}),
	})
}
func schemaWithInputObject(ttype types.Input) (types.Schema, error) {
	return types.NewSchema(types.SchemaConfig{
		Query: types.NewObject(types.ObjectConfig{
			Name: "Query",
			Fields: types.FieldConfigMap{
				"f": &types.FieldConfig{
					Type: types.String,
					Args: types.FieldConfigArgument{
						"args": &types.ArgumentConfig{
							Type: ttype,
						},
					},
				},
			},
		}),
	})
}
func schemaWithObjectFieldOfType(fieldType types.Input) (types.Schema, error) {

	badObjectType := types.NewObject(types.ObjectConfig{
		Name: "BadObject",
		Fields: types.FieldConfigMap{
			"badField": &types.FieldConfig{
				Type: fieldType,
			},
		},
	})
	return types.NewSchema(types.SchemaConfig{
		Query: types.NewObject(types.ObjectConfig{
			Name: "Query",
			Fields: types.FieldConfigMap{
				"f": &types.FieldConfig{
					Type: badObjectType,
				},
			},
		}),
	})
}
func schemaWithObjectImplementingType(implementedType *types.Interface) (types.Schema, error) {

	badObjectType := types.NewObject(types.ObjectConfig{
		Name:       "BadObject",
		Interfaces: []*types.Interface{implementedType},
		Fields: types.FieldConfigMap{
			"f": &types.FieldConfig{
				Type: types.String,
			},
		},
	})
	return types.NewSchema(types.SchemaConfig{
		Query: types.NewObject(types.ObjectConfig{
			Name: "Query",
			Fields: types.FieldConfigMap{
				"f": &types.FieldConfig{
					Type: badObjectType,
				},
			},
		}),
	})
}
func schemaWithUnionOfType(ttype *types.Object) (types.Schema, error) {

	badObjectType := types.NewUnion(types.UnionConfig{
		Name: "BadUnion",
		ResolveType: func(value interface{}, info types.ResolveInfo) *types.Object {
			return nil
		},
		Types: []*types.Object{ttype},
	})
	return types.NewSchema(types.SchemaConfig{
		Query: types.NewObject(types.ObjectConfig{
			Name: "Query",
			Fields: types.FieldConfigMap{
				"f": &types.FieldConfig{
					Type: badObjectType,
				},
			},
		}),
	})
}
func schemaWithInterfaceFieldOfType(ttype types.Type) (types.Schema, error) {

	badInterfaceType := types.NewInterface(types.InterfaceConfig{
		Name: "BadInterface",
		Fields: types.FieldConfigMap{
			"badField": &types.FieldConfig{
				Type: ttype,
			},
		},
	})
	return types.NewSchema(types.SchemaConfig{
		Query: types.NewObject(types.ObjectConfig{
			Name: "Query",
			Fields: types.FieldConfigMap{
				"f": &types.FieldConfig{
					Type: badInterfaceType,
				},
			},
		}),
	})
}
func schemaWithArgOfType(ttype types.Type) (types.Schema, error) {

	badObject := types.NewObject(types.ObjectConfig{
		Name: "BadObject",
		Fields: types.FieldConfigMap{
			"badField": &types.FieldConfig{
				Type: types.String,
				Args: types.FieldConfigArgument{
					"badArg": &types.ArgumentConfig{
						Type: ttype,
					},
				},
			},
		},
	})
	return types.NewSchema(types.SchemaConfig{
		Query: types.NewObject(types.ObjectConfig{
			Name: "Query",
			Fields: types.FieldConfigMap{
				"f": &types.FieldConfig{
					Type: badObject,
				},
			},
		}),
	})
}
func schemaWithInputFieldOfType(ttype types.Type) (types.Schema, error) {

	badInputObject := types.NewInputObject(types.InputObjectConfig{
		Name: "BadInputObject",
		Fields: types.InputObjectConfigFieldMap{
			"badField": &types.InputObjectFieldConfig{
				Type: ttype,
			},
		},
	})
	return types.NewSchema(types.SchemaConfig{
		Query: types.NewObject(types.ObjectConfig{
			Name: "Query",
			Fields: types.FieldConfigMap{
				"f": &types.FieldConfig{
					Type: types.String,
					Args: types.FieldConfigArgument{
						"badArg": &types.ArgumentConfig{
							Type: badInputObject,
						},
					},
				},
			},
		}),
	})
}

func TestTypeSystem_SchemaMustHaveObjectRootTypes_AcceptsASchemaWhoseQueryTypeIsAnObjectType(t *testing.T) {
	_, err := types.NewSchema(types.SchemaConfig{
		Query: someObjectType,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_SchemaMustHaveObjectRootTypes_AcceptsASchemaWhoseQueryAndMutationTypesAreObjectType(t *testing.T) {
	mutationObject := types.NewObject(types.ObjectConfig{
		Name: "Mutation",
		Fields: types.FieldConfigMap{
			"edit": &types.FieldConfig{
				Type: types.String,
			},
		},
	})
	_, err := types.NewSchema(types.SchemaConfig{
		Query:    someObjectType,
		Mutation: mutationObject,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_SchemaMustHaveObjectRootTypes_RejectsASchemaWithoutAQueryType(t *testing.T) {
	_, err := types.NewSchema(types.SchemaConfig{})
	expectedError := "Schema query must be Object Type but got: nil."
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}

func TestTypeSystem_SchemaMustContainUniquelyNamedTypes_RejectsASchemaWhichRedefinesABuiltInType(t *testing.T) {

	fakeString := types.NewScalar(types.ScalarConfig{
		Name: "String",
		Serialize: func(value interface{}) interface{} {
			return nil
		},
	})
	queryType := types.NewObject(types.ObjectConfig{
		Name: "Query",
		Fields: types.FieldConfigMap{
			"normal": &types.FieldConfig{
				Type: types.String,
			},
			"fake": &types.FieldConfig{
				Type: fakeString,
			},
		},
	})
	_, err := types.NewSchema(types.SchemaConfig{
		Query: queryType,
	})
	expectedError := `Schema must contain unique named types but contains multiple types named "String".`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
func TestTypeSystem_SchemaMustContainUniquelyNamedTypes_RejectsASchemaWhichDefinesAnObjectTypeTwice(t *testing.T) {

	a := types.NewObject(types.ObjectConfig{
		Name: "SameName",
		Fields: types.FieldConfigMap{
			"f": &types.FieldConfig{
				Type: types.String,
			},
		},
	})
	b := types.NewObject(types.ObjectConfig{
		Name: "SameName",
		Fields: types.FieldConfigMap{
			"f": &types.FieldConfig{
				Type: types.String,
			},
		},
	})
	queryType := types.NewObject(types.ObjectConfig{
		Name: "Query",
		Fields: types.FieldConfigMap{
			"a": &types.FieldConfig{
				Type: a,
			},
			"b": &types.FieldConfig{
				Type: b,
			},
		},
	})
	_, err := types.NewSchema(types.SchemaConfig{
		Query: queryType,
	})
	expectedError := `Schema must contain unique named types but contains multiple types named "SameName".`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
func TestTypeSystem_SchemaMustContainUniquelyNamedTypes_RejectsASchemaWhichHaveSameNamedObjectsImplementingAnInterface(t *testing.T) {

	anotherInterface := types.NewInterface(types.InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.ResolveInfo) *types.Object {
			return nil
		},
		Fields: types.FieldConfigMap{
			"f": &types.FieldConfig{
				Type: types.String,
			},
		},
	})
	_ = types.NewObject(types.ObjectConfig{
		Name: "BadObject",
		Interfaces: []*types.Interface{
			anotherInterface,
		},
		Fields: types.FieldConfigMap{
			"f": &types.FieldConfig{
				Type: types.String,
			},
		},
	})
	_ = types.NewObject(types.ObjectConfig{
		Name: "BadObject",
		Interfaces: []*types.Interface{
			anotherInterface,
		},
		Fields: types.FieldConfigMap{
			"f": &types.FieldConfig{
				Type: types.String,
			},
		},
	})
	queryType := types.NewObject(types.ObjectConfig{
		Name: "Query",
		Fields: types.FieldConfigMap{
			"iface": &types.FieldConfig{
				Type: anotherInterface,
			},
		},
	})
	_, err := types.NewSchema(types.SchemaConfig{
		Query: queryType,
	})
	expectedError := `Schema must contain unique named types but contains multiple types named "BadObject".`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}

func TestTypeSystem_ObjectsMustHaveFields_AcceptsAnObjectTypeWithFieldsObject(t *testing.T) {
	_, err := schemaWithFieldType(types.NewObject(types.ObjectConfig{
		Name: "SomeObject",
		Fields: types.FieldConfigMap{
			"f": &types.FieldConfig{
				Type: types.String,
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_ObjectsMustHaveFields_RejectsAnObjectTypeWithMissingFields(t *testing.T) {
	badObject := types.NewObject(types.ObjectConfig{
		Name: "SomeObject",
	})
	_, err := schemaWithFieldType(badObject)
	expectedError := `SomeObject fields must be an object with field names as keys or a function which return such an object.`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
func TestTypeSystem_ObjectsMustHaveFields_RejectsAnObjectTypeWithIncorrectlyNamedFields(t *testing.T) {
	badObject := types.NewObject(types.ObjectConfig{
		Name: "SomeObject",
		Fields: types.FieldConfigMap{
			"bad-name-with-dashes": &types.FieldConfig{
				Type: types.String,
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
	badObject := types.NewObject(types.ObjectConfig{
		Name:   "SomeObject",
		Fields: types.FieldConfigMap{},
	})
	_, err := schemaWithFieldType(badObject)
	expectedError := `SomeObject fields must be an object with field names as keys or a function which return such an object.`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}

func TestTypeSystem_FieldsArgsMustBeProperlyNamed_AcceptsFieldArgsWithValidNames(t *testing.T) {
	_, err := schemaWithFieldType(types.NewObject(types.ObjectConfig{
		Name: "SomeObject",
		Fields: types.FieldConfigMap{
			"goodField": &types.FieldConfig{
				Type: types.String,
				Args: types.FieldConfigArgument{
					"goodArgs": &types.ArgumentConfig{
						Type: types.String,
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
	_, err := schemaWithFieldType(types.NewObject(types.ObjectConfig{
		Name: "SomeObject",
		Fields: types.FieldConfigMap{
			"badField": &types.FieldConfig{
				Type: types.String,
				Args: types.FieldConfigArgument{
					"bad-name-with-dashes": &types.ArgumentConfig{
						Type: types.String,
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
	_, err := schemaWithFieldType(types.NewObject(types.ObjectConfig{
		Name: "SomeObject",
		Fields: types.FieldConfigMap{
			"goodField": &types.FieldConfig{
				Type: types.String,
				Args: types.FieldConfigArgument{
					"goodArgs": &types.ArgumentConfig{
						Type: types.String,
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
	anotherInterfaceType := types.NewInterface(types.InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.ResolveInfo) *types.Object {
			return nil
		},
		Fields: types.FieldConfigMap{
			"f": &types.FieldConfig{
				Type: types.String,
			},
		},
	})
	_, err := schemaWithFieldType(types.NewObject(types.ObjectConfig{
		Name: "SomeObject",
		Interfaces: (types.InterfacesThunk)(func() []*types.Interface {
			return []*types.Interface{anotherInterfaceType}
		}),
		Fields: types.FieldConfigMap{
			"f": &types.FieldConfig{
				Type: types.String,
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTypeSystem_ObjectInterfacesMustBeArray_AcceptsAnObjectTypeWithInterfacesAsFunctionReturningAnArray(t *testing.T) {
	anotherInterfaceType := types.NewInterface(types.InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.ResolveInfo) *types.Object {
			return nil
		},
		Fields: types.FieldConfigMap{
			"f": &types.FieldConfig{
				Type: types.String,
			},
		},
	})
	_, err := schemaWithFieldType(types.NewObject(types.ObjectConfig{
		Name:       "SomeObject",
		Interfaces: []*types.Interface{anotherInterfaceType},
		Fields: types.FieldConfigMap{
			"f": &types.FieldConfig{
				Type: types.String,
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTypeSystem_UnionTypesMustBeArray_AcceptsAUnionTypeWithArrayTypes(t *testing.T) {
	_, err := schemaWithFieldType(types.NewUnion(types.UnionConfig{
		Name: "SomeUnion",
		ResolveType: func(value interface{}, info types.ResolveInfo) *types.Object {
			return nil
		},
		Types: []*types.Object{
			someObjectType,
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_UnionTypesMustBeArray_RejectsAUnionTypeWithoutTypes(t *testing.T) {
	_, err := schemaWithFieldType(types.NewUnion(types.UnionConfig{
		Name: "SomeUnion",
		ResolveType: func(value interface{}, info types.ResolveInfo) *types.Object {
			return nil
		},
	}))
	expectedError := "Must provide Array of types for Union SomeUnion."
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
func TestTypeSystem_UnionTypesMustBeArray_RejectsAUnionTypeWithEmptyTypes(t *testing.T) {
	_, err := schemaWithFieldType(types.NewUnion(types.UnionConfig{
		Name: "SomeUnion",
		ResolveType: func(value interface{}, info types.ResolveInfo) *types.Object {
			return nil
		},
		Types: []*types.Object{},
	}))
	expectedError := "Must provide Array of types for Union SomeUnion."
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}

func TestTypeSystem_InputObjectsMustHaveFields_AcceptsAnInputObjectTypeWithFields(t *testing.T) {
	_, err := schemaWithInputObject(types.NewInputObject(types.InputObjectConfig{
		Name: "SomeInputObject",
		Fields: types.InputObjectConfigFieldMap{
			"f": &types.InputObjectFieldConfig{
				Type: types.String,
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTypeSystem_InputObjectsMustHaveFields_AcceptsAnInputObjectTypeWithAFieldFunction(t *testing.T) {
	_, err := schemaWithInputObject(types.NewInputObject(types.InputObjectConfig{
		Name: "SomeInputObject",
		Fields: (types.InputObjectConfigFieldMapThunk)(func() types.InputObjectConfigFieldMap {
			return types.InputObjectConfigFieldMap{
				"f": &types.InputObjectFieldConfig{
					Type: types.String,
				},
			}
		}),
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTypeSystem_InputObjectsMustHaveFields_RejectsAnInputObjectTypeWithMissingFields(t *testing.T) {
	_, err := schemaWithInputObject(types.NewInputObject(types.InputObjectConfig{
		Name: "SomeInputObject",
	}))
	expectedError := "SomeInputObject fields must be an object with field names as keys or a function which return such an object."
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
func TestTypeSystem_InputObjectsMustHaveFields_RejectsAnInputObjectTypeWithEmptyFields(t *testing.T) {
	_, err := schemaWithInputObject(types.NewInputObject(types.InputObjectConfig{
		Name:   "SomeInputObject",
		Fields: types.InputObjectConfigFieldMap{},
	}))
	expectedError := "SomeInputObject fields must be an object with field names as keys or a function which return such an object."
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}

func TestTypeSystem_ObjectTypesMustBeAssertable_AcceptsAnObjectTypeWithAnIsTypeOfFunction(t *testing.T) {
	_, err := schemaWithFieldType(types.NewObject(types.ObjectConfig{
		Name: "AnotherObject",
		IsTypeOf: func(value interface{}, info types.ResolveInfo) bool {
			return true
		},
		Fields: types.FieldConfigMap{
			"f": &types.FieldConfig{
				Type: types.String,
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTypeSystem_InterfaceTypesMustBeResolvable_AcceptsAnInterfaceTypeDefiningResolveType(t *testing.T) {

	anotherInterfaceType := types.NewInterface(types.InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.ResolveInfo) *types.Object {
			return nil
		},
		Fields: types.FieldConfigMap{
			"f": &types.FieldConfig{
				Type: types.String,
			},
		},
	})
	_, err := schemaWithFieldType(types.NewObject(types.ObjectConfig{
		Name:       "SomeObject",
		Interfaces: []*types.Interface{anotherInterfaceType},
		Fields: types.FieldConfigMap{
			"f": &types.FieldConfig{
				Type: types.String,
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_InterfaceTypesMustBeResolvable_AcceptsAnInterfaceWithImplementingTypeDefiningIsTypeOf(t *testing.T) {

	anotherInterfaceType := types.NewInterface(types.InterfaceConfig{
		Name: "AnotherInterface",
		Fields: types.FieldConfigMap{
			"f": &types.FieldConfig{
				Type: types.String,
			},
		},
	})
	_, err := schemaWithFieldType(types.NewObject(types.ObjectConfig{
		Name:       "SomeObject",
		Interfaces: []*types.Interface{anotherInterfaceType},
		IsTypeOf: func(value interface{}, info types.ResolveInfo) bool {
			return true
		},
		Fields: types.FieldConfigMap{
			"f": &types.FieldConfig{
				Type: types.String,
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTypeSystem_InterfaceTypesMustBeResolvable_AcceptsAnInterfaceTypeDefiningResolveTypeWithImplementingTypeDefiningIsTypeOf(t *testing.T) {

	anotherInterfaceType := types.NewInterface(types.InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.ResolveInfo) *types.Object {
			return nil
		},
		Fields: types.FieldConfigMap{
			"f": &types.FieldConfig{
				Type: types.String,
			},
		},
	})
	_, err := schemaWithFieldType(types.NewObject(types.ObjectConfig{
		Name:       "SomeObject",
		Interfaces: []*types.Interface{anotherInterfaceType},
		IsTypeOf: func(value interface{}, info types.ResolveInfo) bool {
			return true
		},
		Fields: types.FieldConfigMap{
			"f": &types.FieldConfig{
				Type: types.String,
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTypeSystem_UnionTypesMustBeResolvable_AcceptsAUnionTypeDefiningResolveType(t *testing.T) {

	_, err := schemaWithFieldType(types.NewUnion(types.UnionConfig{
		Name:  "SomeUnion",
		Types: []*types.Object{someObjectType},
		ResolveType: func(value interface{}, info types.ResolveInfo) *types.Object {
			return nil
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_UnionTypesMustBeResolvable_AcceptsAUnionOfObjectTypesDefiningIsTypeOf(t *testing.T) {

	_, err := schemaWithFieldType(types.NewUnion(types.UnionConfig{
		Name:  "SomeUnion",
		Types: []*types.Object{objectWithIsTypeOf},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_UnionTypesMustBeResolvable_AcceptsAUnionTypeDefiningResolveTypeOfObjectTypesDefiningIsTypeOf(t *testing.T) {

	_, err := schemaWithFieldType(types.NewUnion(types.UnionConfig{
		Name:  "SomeUnion",
		Types: []*types.Object{objectWithIsTypeOf},
		ResolveType: func(value interface{}, info types.ResolveInfo) *types.Object {
			return nil
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_UnionTypesMustBeResolvable_RejectsAUnionTypeNotDefiningResolveTypeOfObjectTypesNotDefiningIsTypeOf(t *testing.T) {

	_, err := schemaWithFieldType(types.NewUnion(types.UnionConfig{
		Name:  "SomeUnion",
		Types: []*types.Object{someObjectType},
	}))
	expectedError := `Union Type SomeUnion does not provide a "resolveType" function and ` +
		`possible Type SomeObject does not provide a "isTypeOf" function. ` +
		`There is no way to resolve this possible type during execution.`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}

func TestTypeSystem_ScalarTypesMustBeSerializable_AcceptsAScalarTypeDefiningSerialize(t *testing.T) {

	_, err := schemaWithFieldType(types.NewScalar(types.ScalarConfig{
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

	_, err := schemaWithFieldType(types.NewScalar(types.ScalarConfig{
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

	_, err := schemaWithFieldType(types.NewScalar(types.ScalarConfig{
		Name: "SomeScalar",
		Serialize: func(value interface{}) interface{} {
			return nil
		},
		ParseValue: func(value interface{}) interface{} {
			return nil
		},
		ParseLiteral: func(valueAST ast.Value) interface{} {
			return nil
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_ScalarTypesMustBeSerializable_RejectsAScalarTypeDefiningParseValueButNotParseLiteral(t *testing.T) {

	_, err := schemaWithFieldType(types.NewScalar(types.ScalarConfig{
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

	_, err := schemaWithFieldType(types.NewScalar(types.ScalarConfig{
		Name: "SomeScalar",
		Serialize: func(value interface{}) interface{} {
			return nil
		},
		ParseLiteral: func(valueAST ast.Value) interface{} {
			return nil
		},
	}))
	expectedError := `SomeScalar must provide both "parseValue" and "parseLiteral" functions.`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}

func TestTypeSystem_EnumTypesMustBeWellDefined_AcceptsAWellDefinedEnumTypeWithEmptyValueDefinition(t *testing.T) {

	_, err := schemaWithFieldType(types.NewEnum(types.EnumConfig{
		Name: "SomeEnum",
		Values: types.EnumValueConfigMap{
			"FOO": &types.EnumValueConfig{},
			"BAR": &types.EnumValueConfig{},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_EnumTypesMustBeWellDefined_AcceptsAWellDefinedEnumTypeWithInternalValueDefinition(t *testing.T) {

	_, err := schemaWithFieldType(types.NewEnum(types.EnumConfig{
		Name: "SomeEnum",
		Values: types.EnumValueConfigMap{
			"FOO": &types.EnumValueConfig{
				Value: 10,
			},
			"BAR": &types.EnumValueConfig{
				Value: 20,
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_EnumTypesMustBeWellDefined_RejectsAnEnumTypeWithoutValues(t *testing.T) {

	_, err := schemaWithFieldType(types.NewEnum(types.EnumConfig{
		Name: "SomeEnum",
	}))
	expectedError := `SomeEnum values must be an object with value names as keys.`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
func TestTypeSystem_EnumTypesMustBeWellDefined_RejectsAnEnumTypeWithEmptyValues(t *testing.T) {

	_, err := schemaWithFieldType(types.NewEnum(types.EnumConfig{
		Name:   "SomeEnum",
		Values: types.EnumValueConfigMap{},
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
	anotherInterfaceType := types.NewInterface(types.InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.ResolveInfo) *types.Object {
			return nil
		},
		Fields: types.FieldConfigMap{
			"f": &types.FieldConfig{
				Type: types.String,
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
	testTypes := withModifiers([]types.Type{
		types.String,
		someScalarType,
		someEnumType,
		someObjectType,
		someUnionType,
		someInterfaceType,
	})
	for _, ttype := range testTypes {
		result := types.NewList(ttype)
		if result.GetError() != nil {
			t.Fatalf(`unexpected error: %v for type "%v"`, result.GetError(), ttype)
		}
	}
}
func TestTypeSystem_ListMustAcceptGraphQLTypes_RejectsANilTypeAsItemTypeOfList(t *testing.T) {
	result := types.NewList(nil)
	expectedError := `Can only create List of a Type but got: <nil>.`
	if result.GetError() == nil || result.GetError().Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, result.GetError())
	}
}

func TestTypeSystem_NonNullMustAcceptGraphQLTypes_AcceptsAnTypeAsNullableTypeOfNonNull(t *testing.T) {
	nullableTypes := []types.Type{
		types.String,
		someScalarType,
		someObjectType,
		someUnionType,
		someInterfaceType,
		someEnumType,
		someInputObject,
		types.NewList(types.String),
		types.NewList(types.NewNonNull(types.String)),
	}
	for _, ttype := range nullableTypes {
		result := types.NewNonNull(ttype)
		if result.GetError() != nil {
			t.Fatalf(`unexpected error: %v for type "%v"`, result.GetError(), ttype)
		}
	}
}
func TestTypeSystem_NonNullMustAcceptGraphQLTypes_RejectsNilAsNonNullableType(t *testing.T) {
	result := types.NewNonNull(nil)
	expectedError := `Can only create NonNull of a Nullable Type but got: <nil>.`
	if result.GetError() == nil || result.GetError().Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, result.GetError())
	}
}

func TestTypeSystem_ObjectsMustAdhereToInterfaceTheyImplement_AcceptsAnObjectWhichImplementsAnInterface(t *testing.T) {
	anotherInterface := types.NewInterface(types.InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.ResolveInfo) *types.Object {
			return nil
		},
		Fields: types.FieldConfigMap{
			"field": &types.FieldConfig{
				Type: types.String,
				Args: types.FieldConfigArgument{
					"input": &types.ArgumentConfig{
						Type: types.String,
					},
				},
			},
		},
	})
	anotherObject := types.NewObject(types.ObjectConfig{
		Name:       "AnotherObject",
		Interfaces: []*types.Interface{anotherInterface},
		Fields: types.FieldConfigMap{
			"field": &types.FieldConfig{
				Type: types.String,
				Args: types.FieldConfigArgument{
					"input": &types.ArgumentConfig{
						Type: types.String,
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
	anotherInterface := types.NewInterface(types.InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.ResolveInfo) *types.Object {
			return nil
		},
		Fields: types.FieldConfigMap{
			"field": &types.FieldConfig{
				Type: types.String,
				Args: types.FieldConfigArgument{
					"input": &types.ArgumentConfig{
						Type: types.String,
					},
				},
			},
		},
	})
	anotherObject := types.NewObject(types.ObjectConfig{
		Name:       "AnotherObject",
		Interfaces: []*types.Interface{anotherInterface},
		Fields: types.FieldConfigMap{
			"field": &types.FieldConfig{
				Type: types.String,
				Args: types.FieldConfigArgument{
					"input": &types.ArgumentConfig{
						Type: types.String,
					},
				},
			},
			"anotherfield": &types.FieldConfig{
				Type: types.String,
			},
		},
	})
	_, err := schemaWithObjectFieldOfType(anotherObject)
	if err != nil {
		t.Fatalf(`unexpected error: %v for type "%v"`, err, anotherObject)
	}
}
func TestTypeSystem_ObjectsMustAdhereToInterfaceTheyImplement_RejectsAnObjectWhichImplementsAnInterfaceFieldAlongWithMoreArguments(t *testing.T) {
	anotherInterface := types.NewInterface(types.InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.ResolveInfo) *types.Object {
			return nil
		},
		Fields: types.FieldConfigMap{
			"field": &types.FieldConfig{
				Type: types.String,
				Args: types.FieldConfigArgument{
					"input": &types.ArgumentConfig{
						Type: types.String,
					},
				},
			},
		},
	})
	anotherObject := types.NewObject(types.ObjectConfig{
		Name:       "AnotherObject",
		Interfaces: []*types.Interface{anotherInterface},
		Fields: types.FieldConfigMap{
			"field": &types.FieldConfig{
				Type: types.String,
				Args: types.FieldConfigArgument{
					"input": &types.ArgumentConfig{
						Type: types.String,
					},
					"anotherInput": &types.ArgumentConfig{
						Type: types.String,
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
	anotherInterface := types.NewInterface(types.InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.ResolveInfo) *types.Object {
			return nil
		},
		Fields: types.FieldConfigMap{
			"field": &types.FieldConfig{
				Type: types.String,
				Args: types.FieldConfigArgument{
					"input": &types.ArgumentConfig{
						Type: types.String,
					},
				},
			},
		},
	})
	anotherObject := types.NewObject(types.ObjectConfig{
		Name:       "AnotherObject",
		Interfaces: []*types.Interface{anotherInterface},
		Fields: types.FieldConfigMap{
			"anotherfield": &types.FieldConfig{
				Type: types.String,
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
	anotherInterface := types.NewInterface(types.InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.ResolveInfo) *types.Object {
			return nil
		},
		Fields: types.FieldConfigMap{
			"field": &types.FieldConfig{
				Type: types.String,
				Args: types.FieldConfigArgument{
					"input": &types.ArgumentConfig{
						Type: types.String,
					},
				},
			},
		},
	})
	anotherObject := types.NewObject(types.ObjectConfig{
		Name:       "AnotherObject",
		Interfaces: []*types.Interface{anotherInterface},
		Fields: types.FieldConfigMap{
			"field": &types.FieldConfig{
				Type: someScalarType,
				Args: types.FieldConfigArgument{
					"input": &types.ArgumentConfig{
						Type: types.String,
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
	anotherInterface := types.NewInterface(types.InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.ResolveInfo) *types.Object {
			return nil
		},
		Fields: types.FieldConfigMap{
			"field": &types.FieldConfig{
				Type: types.String,
				Args: types.FieldConfigArgument{
					"input": &types.ArgumentConfig{
						Type: types.String,
					},
				},
			},
		},
	})
	anotherObject := types.NewObject(types.ObjectConfig{
		Name:       "AnotherObject",
		Interfaces: []*types.Interface{anotherInterface},
		Fields: types.FieldConfigMap{
			"field": &types.FieldConfig{
				Type: types.String,
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
	anotherInterface := types.NewInterface(types.InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.ResolveInfo) *types.Object {
			return nil
		},
		Fields: types.FieldConfigMap{
			"field": &types.FieldConfig{
				Type: types.String,
				Args: types.FieldConfigArgument{
					"input": &types.ArgumentConfig{
						Type: types.String,
					},
				},
			},
		},
	})
	anotherObject := types.NewObject(types.ObjectConfig{
		Name:       "AnotherObject",
		Interfaces: []*types.Interface{anotherInterface},
		Fields: types.FieldConfigMap{
			"field": &types.FieldConfig{
				Type: types.String,
				Args: types.FieldConfigArgument{
					"input": &types.ArgumentConfig{
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
	anotherInterface := types.NewInterface(types.InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.ResolveInfo) *types.Object {
			return nil
		},
		Fields: types.FieldConfigMap{
			"field": &types.FieldConfig{
				Type: types.NewNonNull(types.NewList(types.String)),
			},
		},
	})
	anotherObject := types.NewObject(types.ObjectConfig{
		Name:       "AnotherObject",
		Interfaces: []*types.Interface{anotherInterface},
		Fields: types.FieldConfigMap{
			"field": &types.FieldConfig{
				Type: types.NewNonNull(types.NewList(types.String)),
			},
		},
	})
	_, err := schemaWithObjectFieldOfType(anotherObject)
	if err != nil {
		t.Fatalf(`unexpected error: %v for type "%v"`, err, anotherObject)
	}
}
func TestTypeSystem_ObjectsMustAdhereToInterfaceTheyImplement_RejectsAnObjectWithADifferentlyModifiedInterfaceFieldType(t *testing.T) {
	anotherInterface := types.NewInterface(types.InterfaceConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.ResolveInfo) *types.Object {
			return nil
		},
		Fields: types.FieldConfigMap{
			"field": &types.FieldConfig{
				Type: types.String,
			},
		},
	})
	anotherObject := types.NewObject(types.ObjectConfig{
		Name:       "AnotherObject",
		Interfaces: []*types.Interface{anotherInterface},
		Fields: types.FieldConfigMap{
			"field": &types.FieldConfig{
				Type: types.NewNonNull(types.String),
			},
		},
	})
	_, err := schemaWithObjectFieldOfType(anotherObject)
	expectedError := `AnotherInterface.field expects type "String" but AnotherObject.field provides type "String!".`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
