package types_test

import (
	"testing"

	"github.com/chris-ramon/graphql/language/ast"
	"github.com/chris-ramon/graphql/types"
)

var someScalarType = types.NewGraphQLScalarType(types.GraphQLScalarTypeConfig{
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
var someObjectType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
	Name: "SomeObject",
	Fields: types.GraphQLFieldConfigMap{
		"f": &types.GraphQLFieldConfig{
			Type: types.GraphQLString,
		},
	},
})
var objectWithIsTypeOf = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
	Name: "ObjectWithIsTypeOf",
	IsTypeOf: func(value interface{}, info types.GraphQLResolveInfo) bool {
		return true
	},
	Fields: types.GraphQLFieldConfigMap{
		"f": &types.GraphQLFieldConfig{
			Type: types.GraphQLString,
		},
	},
})
var someUnionType = types.NewGraphQLUnionType(types.GraphQLUnionTypeConfig{
	Name: "SomeUnion",
	ResolveType: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
		return nil
	},
	Types: []*types.GraphQLObjectType{
		someObjectType,
	},
})
var someInterfaceType = types.NewGraphQLInterfaceType(types.GraphQLInterfaceTypeConfig{
	Name: "SomeInterface",
	ResolveType: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
		return nil
	},
	Fields: types.GraphQLFieldConfigMap{
		"f": &types.GraphQLFieldConfig{
			Type: types.GraphQLString,
		},
	},
})
var someEnumType = types.NewGraphQLEnumType(types.GraphQLEnumTypeConfig{
	Name: "SomeEnum",
	Values: types.GraphQLEnumValueConfigMap{
		"ONLY": &types.GraphQLEnumValueConfig{},
	},
})
var someInputObject = types.NewGraphQLInputObjectType(types.InputObjectConfig{
	Name: "SomeInputObject",
	Fields: types.InputObjectConfigFieldMap{
		"f": &types.InputObjectFieldConfig{
			Type:         types.GraphQLString,
			DefaultValue: "Hello",
		},
	},
})

func withModifiers(ttypes []types.GraphQLType) []types.GraphQLType {
	res := ttypes
	for _, ttype := range ttypes {
		res = append(res, types.NewGraphQLList(ttype))
	}
	for _, ttype := range ttypes {
		res = append(res, types.NewGraphQLNonNull(ttype))
	}
	for _, ttype := range ttypes {
		res = append(res, types.NewGraphQLNonNull(types.NewGraphQLList(ttype)))
	}
	return res
}

var outputTypes = withModifiers([]types.GraphQLType{
	types.GraphQLString,
	someScalarType,
	someEnumType,
	someObjectType,
	someUnionType,
	someInterfaceType,
})
var inputTypes = withModifiers([]types.GraphQLType{
	types.GraphQLString,
	someScalarType,
	someEnumType,
	someInputObject,
})

func schemaWithFieldType(ttype types.GraphQLOutputType) (types.GraphQLSchema, error) {
	return types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "Query",
			Fields: types.GraphQLFieldConfigMap{
				"f": &types.GraphQLFieldConfig{
					Type: ttype,
				},
			},
		}),
	})
}
func schemaWithInputObject(ttype types.GraphQLInputType) (types.GraphQLSchema, error) {
	return types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "Query",
			Fields: types.GraphQLFieldConfigMap{
				"f": &types.GraphQLFieldConfig{
					Type: types.GraphQLString,
					Args: types.GraphQLFieldConfigArgumentMap{
						"args": &types.GraphQLArgumentConfig{
							Type: ttype,
						},
					},
				},
			},
		}),
	})
}
func schemaWithObjectFieldOfType(fieldType types.GraphQLInputType) (types.GraphQLSchema, error) {

	badObjectType := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "BadObject",
		Fields: types.GraphQLFieldConfigMap{
			"badField": &types.GraphQLFieldConfig{
				Type: fieldType,
			},
		},
	})
	return types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "Query",
			Fields: types.GraphQLFieldConfigMap{
				"f": &types.GraphQLFieldConfig{
					Type: badObjectType,
				},
			},
		}),
	})
}
func schemaWithObjectImplementingType(implementedType *types.GraphQLInterfaceType) (types.GraphQLSchema, error) {

	badObjectType := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name:       "BadObject",
		Interfaces: []*types.GraphQLInterfaceType{implementedType},
		Fields: types.GraphQLFieldConfigMap{
			"f": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
		},
	})
	return types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "Query",
			Fields: types.GraphQLFieldConfigMap{
				"f": &types.GraphQLFieldConfig{
					Type: badObjectType,
				},
			},
		}),
	})
}
func schemaWithUnionOfType(ttype *types.GraphQLObjectType) (types.GraphQLSchema, error) {

	badObjectType := types.NewGraphQLUnionType(types.GraphQLUnionTypeConfig{
		Name: "BadUnion",
		ResolveType: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
			return nil
		},
		Types: []*types.GraphQLObjectType{ttype},
	})
	return types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "Query",
			Fields: types.GraphQLFieldConfigMap{
				"f": &types.GraphQLFieldConfig{
					Type: badObjectType,
				},
			},
		}),
	})
}
func schemaWithInterfaceFieldOfType(ttype types.GraphQLType) (types.GraphQLSchema, error) {

	badInterfaceType := types.NewGraphQLInterfaceType(types.GraphQLInterfaceTypeConfig{
		Name: "BadInterface",
		Fields: types.GraphQLFieldConfigMap{
			"badField": &types.GraphQLFieldConfig{
				Type: ttype,
			},
		},
	})
	return types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "Query",
			Fields: types.GraphQLFieldConfigMap{
				"f": &types.GraphQLFieldConfig{
					Type: badInterfaceType,
				},
			},
		}),
	})
}
func schemaWithArgOfType(ttype types.GraphQLType) (types.GraphQLSchema, error) {

	badObject := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "BadObject",
		Fields: types.GraphQLFieldConfigMap{
			"badField": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Args: types.GraphQLFieldConfigArgumentMap{
					"badArg": &types.GraphQLArgumentConfig{
						Type: ttype,
					},
				},
			},
		},
	})
	return types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "Query",
			Fields: types.GraphQLFieldConfigMap{
				"f": &types.GraphQLFieldConfig{
					Type: badObject,
				},
			},
		}),
	})
}
func schemaWithInputFieldOfType(ttype types.GraphQLType) (types.GraphQLSchema, error) {

	badInputObject := types.NewGraphQLInputObjectType(types.InputObjectConfig{
		Name: "BadInputObject",
		Fields: types.InputObjectConfigFieldMap{
			"badField": &types.InputObjectFieldConfig{
				Type: ttype,
			},
		},
	})
	return types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "Query",
			Fields: types.GraphQLFieldConfigMap{
				"f": &types.GraphQLFieldConfig{
					Type: types.GraphQLString,
					Args: types.GraphQLFieldConfigArgumentMap{
						"badArg": &types.GraphQLArgumentConfig{
							Type: badInputObject,
						},
					},
				},
			},
		}),
	})
}

func TestTypeSystem_SchemaMustHaveObjectRootTypes_AcceptsASchemaWhoseQueryTypeIsAnObjectType(t *testing.T) {
	_, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: someObjectType,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_SchemaMustHaveObjectRootTypes_AcceptsASchemaWhoseQueryAndMutationTypesAreObjectType(t *testing.T) {
	mutationObject := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "Mutation",
		Fields: types.GraphQLFieldConfigMap{
			"edit": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
		},
	})
	_, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query:    someObjectType,
		Mutation: mutationObject,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_SchemaMustHaveObjectRootTypes_RejectsASchemaWithoutAQueryType(t *testing.T) {
	_, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{})
	expectedError := "Schema query must be Object Type but got: nil."
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}

func TestTypeSystem_SchemaMustContainUniquelyNamedTypes_RejectsASchemaWhichRedefinesABuiltInType(t *testing.T) {

	fakeString := types.NewGraphQLScalarType(types.GraphQLScalarTypeConfig{
		Name: "String",
		Serialize: func(value interface{}) interface{} {
			return nil
		},
	})
	queryType := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "Query",
		Fields: types.GraphQLFieldConfigMap{
			"normal": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
			"fake": &types.GraphQLFieldConfig{
				Type: fakeString,
			},
		},
	})
	_, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: queryType,
	})
	expectedError := `Schema must contain unique named types but contains multiple types named "String".`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
func TestTypeSystem_SchemaMustContainUniquelyNamedTypes_RejectsASchemaWhichDefinesAnObjectTypeTwice(t *testing.T) {

	a := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "SameName",
		Fields: types.GraphQLFieldConfigMap{
			"f": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
		},
	})
	b := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "SameName",
		Fields: types.GraphQLFieldConfigMap{
			"f": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
		},
	})
	queryType := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "Query",
		Fields: types.GraphQLFieldConfigMap{
			"a": &types.GraphQLFieldConfig{
				Type: a,
			},
			"b": &types.GraphQLFieldConfig{
				Type: b,
			},
		},
	})
	_, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: queryType,
	})
	expectedError := `Schema must contain unique named types but contains multiple types named "SameName".`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
func TestTypeSystem_SchemaMustContainUniquelyNamedTypes_RejectsASchemaWhichHaveSameNamedObjectsImplementingAnInterface(t *testing.T) {

	anotherInterface := types.NewGraphQLInterfaceType(types.GraphQLInterfaceTypeConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
			return nil
		},
		Fields: types.GraphQLFieldConfigMap{
			"f": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
		},
	})
	_ = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "BadObject",
		Interfaces: []*types.GraphQLInterfaceType{
			anotherInterface,
		},
		Fields: types.GraphQLFieldConfigMap{
			"f": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
		},
	})
	_ = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "BadObject",
		Interfaces: []*types.GraphQLInterfaceType{
			anotherInterface,
		},
		Fields: types.GraphQLFieldConfigMap{
			"f": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
		},
	})
	queryType := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "Query",
		Fields: types.GraphQLFieldConfigMap{
			"iface": &types.GraphQLFieldConfig{
				Type: anotherInterface,
			},
		},
	})
	_, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: queryType,
	})
	expectedError := `Schema must contain unique named types but contains multiple types named "BadObject".`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}

func TestTypeSystem_ObjectsMustHaveFields_AcceptsAnObjectTypeWithFieldsObject(t *testing.T) {
	_, err := schemaWithFieldType(types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "SomeObject",
		Fields: types.GraphQLFieldConfigMap{
			"f": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_ObjectsMustHaveFields_RejectsAnObjectTypeWithMissingFields(t *testing.T) {
	badObject := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "SomeObject",
	})
	_, err := schemaWithFieldType(badObject)
	expectedError := `SomeObject fields must be an object with field names as keys or a function which return such an object.`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
func TestTypeSystem_ObjectsMustHaveFields_RejectsAnObjectTypeWithIncorrectlyNamedFields(t *testing.T) {
	badObject := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "SomeObject",
		Fields: types.GraphQLFieldConfigMap{
			"bad-name-with-dashes": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
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
	badObject := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name:   "SomeObject",
		Fields: types.GraphQLFieldConfigMap{},
	})
	_, err := schemaWithFieldType(badObject)
	expectedError := `SomeObject fields must be an object with field names as keys or a function which return such an object.`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}

func TestTypeSystem_FieldsArgsMustBeProperlyNamed_AcceptsFieldArgsWithValidNames(t *testing.T) {
	_, err := schemaWithFieldType(types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "SomeObject",
		Fields: types.GraphQLFieldConfigMap{
			"goodField": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Args: types.GraphQLFieldConfigArgumentMap{
					"goodArgs": &types.GraphQLArgumentConfig{
						Type: types.GraphQLString,
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
	_, err := schemaWithFieldType(types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "SomeObject",
		Fields: types.GraphQLFieldConfigMap{
			"badField": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Args: types.GraphQLFieldConfigArgumentMap{
					"bad-name-with-dashes": &types.GraphQLArgumentConfig{
						Type: types.GraphQLString,
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
	_, err := schemaWithFieldType(types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "SomeObject",
		Fields: types.GraphQLFieldConfigMap{
			"goodField": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Args: types.GraphQLFieldConfigArgumentMap{
					"goodArgs": &types.GraphQLArgumentConfig{
						Type: types.GraphQLString,
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
	anotherInterfaceType := types.NewGraphQLInterfaceType(types.GraphQLInterfaceTypeConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
			return nil
		},
		Fields: types.GraphQLFieldConfigMap{
			"f": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
		},
	})
	_, err := schemaWithFieldType(types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "SomeObject",
		Interfaces: (types.GraphQLInterfacesThunk)(func() []*types.GraphQLInterfaceType {
			return []*types.GraphQLInterfaceType{anotherInterfaceType}
		}),
		Fields: types.GraphQLFieldConfigMap{
			"f": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTypeSystem_ObjectInterfacesMustBeArray_AcceptsAnObjectTypeWithInterfacesAsFunctionReturningAnArray(t *testing.T) {
	anotherInterfaceType := types.NewGraphQLInterfaceType(types.GraphQLInterfaceTypeConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
			return nil
		},
		Fields: types.GraphQLFieldConfigMap{
			"f": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
		},
	})
	_, err := schemaWithFieldType(types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name:       "SomeObject",
		Interfaces: []*types.GraphQLInterfaceType{anotherInterfaceType},
		Fields: types.GraphQLFieldConfigMap{
			"f": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTypeSystem_UnionTypesMustBeArray_AcceptsAUnionTypeWithArrayTypes(t *testing.T) {
	_, err := schemaWithFieldType(types.NewGraphQLUnionType(types.GraphQLUnionTypeConfig{
		Name: "SomeUnion",
		ResolveType: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
			return nil
		},
		Types: []*types.GraphQLObjectType{
			someObjectType,
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_UnionTypesMustBeArray_RejectsAUnionTypeWithoutTypes(t *testing.T) {
	_, err := schemaWithFieldType(types.NewGraphQLUnionType(types.GraphQLUnionTypeConfig{
		Name: "SomeUnion",
		ResolveType: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
			return nil
		},
	}))
	expectedError := "Must provide Array of types for Union SomeUnion."
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
func TestTypeSystem_UnionTypesMustBeArray_RejectsAUnionTypeWithEmptyTypes(t *testing.T) {
	_, err := schemaWithFieldType(types.NewGraphQLUnionType(types.GraphQLUnionTypeConfig{
		Name: "SomeUnion",
		ResolveType: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
			return nil
		},
		Types: []*types.GraphQLObjectType{},
	}))
	expectedError := "Must provide Array of types for Union SomeUnion."
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}

func TestTypeSystem_InputObjectsMustHaveFields_AcceptsAnInputObjectTypeWithFields(t *testing.T) {
	_, err := schemaWithInputObject(types.NewGraphQLInputObjectType(types.InputObjectConfig{
		Name: "SomeInputObject",
		Fields: types.InputObjectConfigFieldMap{
			"f": &types.InputObjectFieldConfig{
				Type: types.GraphQLString,
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTypeSystem_InputObjectsMustHaveFields_AcceptsAnInputObjectTypeWithAFieldFunction(t *testing.T) {
	_, err := schemaWithInputObject(types.NewGraphQLInputObjectType(types.InputObjectConfig{
		Name: "SomeInputObject",
		Fields: (types.InputObjectConfigFieldMapThunk)(func() types.InputObjectConfigFieldMap {
			return types.InputObjectConfigFieldMap{
				"f": &types.InputObjectFieldConfig{
					Type: types.GraphQLString,
				},
			}
		}),
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTypeSystem_InputObjectsMustHaveFields_RejectsAnInputObjectTypeWithMissingFields(t *testing.T) {
	_, err := schemaWithInputObject(types.NewGraphQLInputObjectType(types.InputObjectConfig{
		Name: "SomeInputObject",
	}))
	expectedError := "SomeInputObject fields must be an object with field names as keys or a function which return such an object."
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
func TestTypeSystem_InputObjectsMustHaveFields_RejectsAnInputObjectTypeWithEmptyFields(t *testing.T) {
	_, err := schemaWithInputObject(types.NewGraphQLInputObjectType(types.InputObjectConfig{
		Name:   "SomeInputObject",
		Fields: types.InputObjectConfigFieldMap{},
	}))
	expectedError := "SomeInputObject fields must be an object with field names as keys or a function which return such an object."
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}

func TestTypeSystem_ObjectTypesMustBeAssertable_AcceptsAnObjectTypeWithAnIsTypeOfFunction(t *testing.T) {
	_, err := schemaWithFieldType(types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "AnotherObject",
		IsTypeOf: func(value interface{}, info types.GraphQLResolveInfo) bool {
			return true
		},
		Fields: types.GraphQLFieldConfigMap{
			"f": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTypeSystem_InterfaceTypesMustBeResolvable_AcceptsAnInterfaceTypeDefiningResolveType(t *testing.T) {

	anotherInterfaceType := types.NewGraphQLInterfaceType(types.GraphQLInterfaceTypeConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
			return nil
		},
		Fields: types.GraphQLFieldConfigMap{
			"f": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
		},
	})
	_, err := schemaWithFieldType(types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name:       "SomeObject",
		Interfaces: []*types.GraphQLInterfaceType{anotherInterfaceType},
		Fields: types.GraphQLFieldConfigMap{
			"f": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_InterfaceTypesMustBeResolvable_AcceptsAnInterfaceWithImplementingTypeDefiningIsTypeOf(t *testing.T) {

	anotherInterfaceType := types.NewGraphQLInterfaceType(types.GraphQLInterfaceTypeConfig{
		Name: "AnotherInterface",
		Fields: types.GraphQLFieldConfigMap{
			"f": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
		},
	})
	_, err := schemaWithFieldType(types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name:       "SomeObject",
		Interfaces: []*types.GraphQLInterfaceType{anotherInterfaceType},
		IsTypeOf: func(value interface{}, info types.GraphQLResolveInfo) bool {
			return true
		},
		Fields: types.GraphQLFieldConfigMap{
			"f": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTypeSystem_InterfaceTypesMustBeResolvable_AcceptsAnInterfaceTypeDefiningResolveTypeWithImplementingTypeDefiningIsTypeOf(t *testing.T) {

	anotherInterfaceType := types.NewGraphQLInterfaceType(types.GraphQLInterfaceTypeConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
			return nil
		},
		Fields: types.GraphQLFieldConfigMap{
			"f": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
		},
	})
	_, err := schemaWithFieldType(types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name:       "SomeObject",
		Interfaces: []*types.GraphQLInterfaceType{anotherInterfaceType},
		IsTypeOf: func(value interface{}, info types.GraphQLResolveInfo) bool {
			return true
		},
		Fields: types.GraphQLFieldConfigMap{
			"f": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTypeSystem_UnionTypesMustBeResolvable_AcceptsAUnionTypeDefiningResolveType(t *testing.T) {

	_, err := schemaWithFieldType(types.NewGraphQLUnionType(types.GraphQLUnionTypeConfig{
		Name:  "SomeUnion",
		Types: []*types.GraphQLObjectType{someObjectType},
		ResolveType: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
			return nil
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_UnionTypesMustBeResolvable_AcceptsAUnionOfObjectTypesDefiningIsTypeOf(t *testing.T) {

	_, err := schemaWithFieldType(types.NewGraphQLUnionType(types.GraphQLUnionTypeConfig{
		Name:  "SomeUnion",
		Types: []*types.GraphQLObjectType{objectWithIsTypeOf},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_UnionTypesMustBeResolvable_AcceptsAUnionTypeDefiningResolveTypeOfObjectTypesDefiningIsTypeOf(t *testing.T) {

	_, err := schemaWithFieldType(types.NewGraphQLUnionType(types.GraphQLUnionTypeConfig{
		Name:  "SomeUnion",
		Types: []*types.GraphQLObjectType{objectWithIsTypeOf},
		ResolveType: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
			return nil
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_UnionTypesMustBeResolvable_RejectsAUnionTypeNotDefiningResolveTypeOfObjectTypesNotDefiningIsTypeOf(t *testing.T) {

	_, err := schemaWithFieldType(types.NewGraphQLUnionType(types.GraphQLUnionTypeConfig{
		Name:  "SomeUnion",
		Types: []*types.GraphQLObjectType{someObjectType},
	}))
	expectedError := `Union Type SomeUnion does not provide a "resolveType" function and ` +
		`possible Type SomeObject does not provide a "isTypeOf" function. ` +
		`There is no way to resolve this possible type during execution.`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}

func TestTypeSystem_ScalarTypesMustBeSerializable_AcceptsAScalarTypeDefiningSerialize(t *testing.T) {

	_, err := schemaWithFieldType(types.NewGraphQLScalarType(types.GraphQLScalarTypeConfig{
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

	_, err := schemaWithFieldType(types.NewGraphQLScalarType(types.GraphQLScalarTypeConfig{
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

	_, err := schemaWithFieldType(types.NewGraphQLScalarType(types.GraphQLScalarTypeConfig{
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

	_, err := schemaWithFieldType(types.NewGraphQLScalarType(types.GraphQLScalarTypeConfig{
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

	_, err := schemaWithFieldType(types.NewGraphQLScalarType(types.GraphQLScalarTypeConfig{
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

	_, err := schemaWithFieldType(types.NewGraphQLEnumType(types.GraphQLEnumTypeConfig{
		Name: "SomeEnum",
		Values: types.GraphQLEnumValueConfigMap{
			"FOO": &types.GraphQLEnumValueConfig{},
			"BAR": &types.GraphQLEnumValueConfig{},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_EnumTypesMustBeWellDefined_AcceptsAWellDefinedEnumTypeWithInternalValueDefinition(t *testing.T) {

	_, err := schemaWithFieldType(types.NewGraphQLEnumType(types.GraphQLEnumTypeConfig{
		Name: "SomeEnum",
		Values: types.GraphQLEnumValueConfigMap{
			"FOO": &types.GraphQLEnumValueConfig{
				Value: 10,
			},
			"BAR": &types.GraphQLEnumValueConfig{
				Value: 20,
			},
		},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestTypeSystem_EnumTypesMustBeWellDefined_RejectsAnEnumTypeWithoutValues(t *testing.T) {

	_, err := schemaWithFieldType(types.NewGraphQLEnumType(types.GraphQLEnumTypeConfig{
		Name: "SomeEnum",
	}))
	expectedError := `SomeEnum values must be an object with value names as keys.`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
func TestTypeSystem_EnumTypesMustBeWellDefined_RejectsAnEnumTypeWithEmptyValues(t *testing.T) {

	_, err := schemaWithFieldType(types.NewGraphQLEnumType(types.GraphQLEnumTypeConfig{
		Name:   "SomeEnum",
		Values: types.GraphQLEnumValueConfigMap{},
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
	anotherInterfaceType := types.NewGraphQLInterfaceType(types.GraphQLInterfaceTypeConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
			return nil
		},
		Fields: types.GraphQLFieldConfigMap{
			"f": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
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
	testTypes := withModifiers([]types.GraphQLType{
		types.GraphQLString,
		someScalarType,
		someEnumType,
		someObjectType,
		someUnionType,
		someInterfaceType,
	})
	for _, ttype := range testTypes {
		result := types.NewGraphQLList(ttype)
		if result.GetError() != nil {
			t.Fatalf(`unexpected error: %v for type "%v"`, result.GetError(), ttype)
		}
	}
}
func TestTypeSystem_ListMustAcceptGraphQLTypes_RejectsANilTypeAsItemTypeOfList(t *testing.T) {
	result := types.NewGraphQLList(nil)
	expectedError := `Can only create List of a GraphQLType but got: <nil>.`
	if result.GetError() == nil || result.GetError().Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, result.GetError())
	}
}

func TestTypeSystem_NonNullMustAcceptGraphQLTypes_AcceptsAnTypeAsNullableTypeOfNonNull(t *testing.T) {
	nullableTypes := []types.GraphQLType{
		types.GraphQLString,
		someScalarType,
		someObjectType,
		someUnionType,
		someInterfaceType,
		someEnumType,
		someInputObject,
		types.NewGraphQLList(types.GraphQLString),
		types.NewGraphQLList(types.NewGraphQLNonNull(types.GraphQLString)),
	}
	for _, ttype := range nullableTypes {
		result := types.NewGraphQLNonNull(ttype)
		if result.GetError() != nil {
			t.Fatalf(`unexpected error: %v for type "%v"`, result.GetError(), ttype)
		}
	}
}
func TestTypeSystem_NonNullMustAcceptGraphQLTypes_RejectsNilAsNonNullableType(t *testing.T) {
	result := types.NewGraphQLNonNull(nil)
	expectedError := `Can only create NonNull of a Nullable GraphQLType but got: <nil>.`
	if result.GetError() == nil || result.GetError().Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, result.GetError())
	}
}

func TestTypeSystem_ObjectsMustAdhereToInterfaceTheyImplement_AcceptsAnObjectWhichImplementsAnInterface(t *testing.T) {
	anotherInterface := types.NewGraphQLInterfaceType(types.GraphQLInterfaceTypeConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
			return nil
		},
		Fields: types.GraphQLFieldConfigMap{
			"field": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Args: types.GraphQLFieldConfigArgumentMap{
					"input": &types.GraphQLArgumentConfig{
						Type: types.GraphQLString,
					},
				},
			},
		},
	})
	anotherObject := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name:       "AnotherObject",
		Interfaces: []*types.GraphQLInterfaceType{anotherInterface},
		Fields: types.GraphQLFieldConfigMap{
			"field": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Args: types.GraphQLFieldConfigArgumentMap{
					"input": &types.GraphQLArgumentConfig{
						Type: types.GraphQLString,
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
	anotherInterface := types.NewGraphQLInterfaceType(types.GraphQLInterfaceTypeConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
			return nil
		},
		Fields: types.GraphQLFieldConfigMap{
			"field": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Args: types.GraphQLFieldConfigArgumentMap{
					"input": &types.GraphQLArgumentConfig{
						Type: types.GraphQLString,
					},
				},
			},
		},
	})
	anotherObject := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name:       "AnotherObject",
		Interfaces: []*types.GraphQLInterfaceType{anotherInterface},
		Fields: types.GraphQLFieldConfigMap{
			"field": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Args: types.GraphQLFieldConfigArgumentMap{
					"input": &types.GraphQLArgumentConfig{
						Type: types.GraphQLString,
					},
				},
			},
			"anotherfield": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
		},
	})
	_, err := schemaWithObjectFieldOfType(anotherObject)
	if err != nil {
		t.Fatalf(`unexpected error: %v for type "%v"`, err, anotherObject)
	}
}
func TestTypeSystem_ObjectsMustAdhereToInterfaceTheyImplement_RejectsAnObjectWhichImplementsAnInterfaceFieldAlongWithMoreArguments(t *testing.T) {
	anotherInterface := types.NewGraphQLInterfaceType(types.GraphQLInterfaceTypeConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
			return nil
		},
		Fields: types.GraphQLFieldConfigMap{
			"field": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Args: types.GraphQLFieldConfigArgumentMap{
					"input": &types.GraphQLArgumentConfig{
						Type: types.GraphQLString,
					},
				},
			},
		},
	})
	anotherObject := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name:       "AnotherObject",
		Interfaces: []*types.GraphQLInterfaceType{anotherInterface},
		Fields: types.GraphQLFieldConfigMap{
			"field": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Args: types.GraphQLFieldConfigArgumentMap{
					"input": &types.GraphQLArgumentConfig{
						Type: types.GraphQLString,
					},
					"anotherInput": &types.GraphQLArgumentConfig{
						Type: types.GraphQLString,
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
	anotherInterface := types.NewGraphQLInterfaceType(types.GraphQLInterfaceTypeConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
			return nil
		},
		Fields: types.GraphQLFieldConfigMap{
			"field": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Args: types.GraphQLFieldConfigArgumentMap{
					"input": &types.GraphQLArgumentConfig{
						Type: types.GraphQLString,
					},
				},
			},
		},
	})
	anotherObject := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name:       "AnotherObject",
		Interfaces: []*types.GraphQLInterfaceType{anotherInterface},
		Fields: types.GraphQLFieldConfigMap{
			"anotherfield": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
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
	anotherInterface := types.NewGraphQLInterfaceType(types.GraphQLInterfaceTypeConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
			return nil
		},
		Fields: types.GraphQLFieldConfigMap{
			"field": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Args: types.GraphQLFieldConfigArgumentMap{
					"input": &types.GraphQLArgumentConfig{
						Type: types.GraphQLString,
					},
				},
			},
		},
	})
	anotherObject := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name:       "AnotherObject",
		Interfaces: []*types.GraphQLInterfaceType{anotherInterface},
		Fields: types.GraphQLFieldConfigMap{
			"field": &types.GraphQLFieldConfig{
				Type: someScalarType,
				Args: types.GraphQLFieldConfigArgumentMap{
					"input": &types.GraphQLArgumentConfig{
						Type: types.GraphQLString,
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
	anotherInterface := types.NewGraphQLInterfaceType(types.GraphQLInterfaceTypeConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
			return nil
		},
		Fields: types.GraphQLFieldConfigMap{
			"field": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Args: types.GraphQLFieldConfigArgumentMap{
					"input": &types.GraphQLArgumentConfig{
						Type: types.GraphQLString,
					},
				},
			},
		},
	})
	anotherObject := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name:       "AnotherObject",
		Interfaces: []*types.GraphQLInterfaceType{anotherInterface},
		Fields: types.GraphQLFieldConfigMap{
			"field": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
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
	anotherInterface := types.NewGraphQLInterfaceType(types.GraphQLInterfaceTypeConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
			return nil
		},
		Fields: types.GraphQLFieldConfigMap{
			"field": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Args: types.GraphQLFieldConfigArgumentMap{
					"input": &types.GraphQLArgumentConfig{
						Type: types.GraphQLString,
					},
				},
			},
		},
	})
	anotherObject := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name:       "AnotherObject",
		Interfaces: []*types.GraphQLInterfaceType{anotherInterface},
		Fields: types.GraphQLFieldConfigMap{
			"field": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Args: types.GraphQLFieldConfigArgumentMap{
					"input": &types.GraphQLArgumentConfig{
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
	anotherInterface := types.NewGraphQLInterfaceType(types.GraphQLInterfaceTypeConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
			return nil
		},
		Fields: types.GraphQLFieldConfigMap{
			"field": &types.GraphQLFieldConfig{
				Type: types.NewGraphQLNonNull(types.NewGraphQLList(types.GraphQLString)),
			},
		},
	})
	anotherObject := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name:       "AnotherObject",
		Interfaces: []*types.GraphQLInterfaceType{anotherInterface},
		Fields: types.GraphQLFieldConfigMap{
			"field": &types.GraphQLFieldConfig{
				Type: types.NewGraphQLNonNull(types.NewGraphQLList(types.GraphQLString)),
			},
		},
	})
	_, err := schemaWithObjectFieldOfType(anotherObject)
	if err != nil {
		t.Fatalf(`unexpected error: %v for type "%v"`, err, anotherObject)
	}
}
func TestTypeSystem_ObjectsMustAdhereToInterfaceTheyImplement_RejectsAnObjectWithADifferentlyModifiedInterfaceFieldType(t *testing.T) {
	anotherInterface := types.NewGraphQLInterfaceType(types.GraphQLInterfaceTypeConfig{
		Name: "AnotherInterface",
		ResolveType: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
			return nil
		},
		Fields: types.GraphQLFieldConfigMap{
			"field": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
		},
	})
	anotherObject := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name:       "AnotherObject",
		Interfaces: []*types.GraphQLInterfaceType{anotherInterface},
		Fields: types.GraphQLFieldConfigMap{
			"field": &types.GraphQLFieldConfig{
				Type: types.NewGraphQLNonNull(types.GraphQLString),
			},
		},
	})
	_, err := schemaWithObjectFieldOfType(anotherObject)
	expectedError := `AnotherInterface.field expects type "String" but AnotherObject.field provides type "String!".`
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error: %v, got %v", expectedError, err)
	}
}
