package types_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/chris-ramon/graphql/testutil"
	"github.com/chris-ramon/graphql/types"
)

var blogImage = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
	Name: "Image",
	Fields: types.GraphQLFieldConfigMap{
		"url": &types.GraphQLFieldConfig{
			Type: types.GraphQLString,
		},
		"width": &types.GraphQLFieldConfig{
			Type: types.GraphQLInt,
		},
		"height": &types.GraphQLFieldConfig{
			Type: types.GraphQLInt,
		},
	},
})
var blogAuthor = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
	Name: "Author",
	Fields: types.GraphQLFieldConfigMap{
		"id": &types.GraphQLFieldConfig{
			Type: types.GraphQLString,
		},
		"name": &types.GraphQLFieldConfig{
			Type: types.GraphQLString,
		},
		"pic": &types.GraphQLFieldConfig{
			Type: blogImage,
			Args: types.GraphQLFieldConfigArgumentMap{
				"width": &types.GraphQLArgumentConfig{
					Type: types.GraphQLInt,
				},
				"height": &types.GraphQLArgumentConfig{
					Type: types.GraphQLInt,
				},
			},
		},
		"recentArticle": &types.GraphQLFieldConfig{},
	},
})
var blogArticle = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
	Name: "Article",
	Fields: types.GraphQLFieldConfigMap{
		"id": &types.GraphQLFieldConfig{
			Type: types.GraphQLString,
		},
		"isPublished": &types.GraphQLFieldConfig{
			Type: types.GraphQLBoolean,
		},
		"author": &types.GraphQLFieldConfig{
			Type: blogAuthor,
		},
		"title": &types.GraphQLFieldConfig{
			Type: types.GraphQLString,
		},
		"body": &types.GraphQLFieldConfig{
			Type: types.GraphQLString,
		},
	},
})
var blogQuery = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
	Name: "Query",
	Fields: types.GraphQLFieldConfigMap{
		"article": &types.GraphQLFieldConfig{
			Type: blogArticle,
			Args: types.GraphQLFieldConfigArgumentMap{
				"id": &types.GraphQLArgumentConfig{
					Type: types.GraphQLString,
				},
			},
		},
		"feed": &types.GraphQLFieldConfig{
			Type: types.NewGraphQLList(blogArticle),
		},
	},
})

var blogMutation = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
	Name: "Mutation",
	Fields: types.GraphQLFieldConfigMap{
		"writeArticle": &types.GraphQLFieldConfig{
			Type: blogArticle,
		},
	},
})

var objectType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
	Name: "Object",
	IsTypeOf: func(value interface{}, info types.GraphQLResolveInfo) bool {
		return true
	},
})
var interfaceType = types.NewGraphQLInterfaceType(types.GraphQLInterfaceTypeConfig{
	Name: "Interface",
})
var unionType = types.NewGraphQLUnionType(types.GraphQLUnionTypeConfig{
	Name: "Union",
	Types: []*types.GraphQLObjectType{
		objectType,
	},
})
var enumType = types.NewGraphQLEnumType(types.GraphQLEnumTypeConfig{
	Name: "Enum",
	Values: types.GraphQLEnumValueConfigMap{
		"foo": &types.GraphQLEnumValueConfig{},
	},
})
var inputObjectType = types.NewGraphQLInputObjectType(types.InputObjectConfig{
	Name: "InputObject",
})

func init() {
	blogAuthor.AddFieldConfig("recentArticle", &types.GraphQLFieldConfig{
		Type: blogArticle,
	})
}

func TestTypeSystem_DefinitionExample_DefinesAQueryOnlySchema(t *testing.T) {
	blogSchema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: blogQuery,
	})
	if err != nil {
		t.Fatalf("unexpected error, got: %v", err)
	}

	if blogSchema.GetQueryType() != blogQuery {
		t.Fatalf("expected blogSchema.GetQueryType() == blogQuery")
	}

	articleField, _ := blogQuery.GetFields()["article"]
	if articleField == nil {
		t.Fatalf("articleField is nil")
	}
	articleFieldType := articleField.Type
	if articleFieldType != blogArticle {
		t.Fatalf("articleFieldType expected to equal blogArticle, got: %v", articleField.Type)
	}
	if articleFieldType.GetName() != "Article" {
		t.Fatalf("articleFieldType.Name expected to equal `Article`, got: %v", articleField.Type.GetName())
	}
	if articleField.Name != "article" {
		t.Fatalf("articleField.Name expected to equal `article`, got: %v", articleField.Name)
	}
	articleFieldTypeObject, ok := articleFieldType.(*types.GraphQLObjectType)
	if !ok {
		t.Fatalf("expected articleFieldType to be *types.GraphQLObjectType`, got: %v", articleField)
	}

	// TODO: expose a GraphQLObjectType.GetField(key string), instead of this ghetto way of accessing a field map?
	titleField := articleFieldTypeObject.GetFields()["title"]
	if titleField == nil {
		t.Fatalf("titleField is nil")
	}
	if titleField.Name != "title" {
		t.Fatalf("titleField.Name expected to equal title, got: %v", titleField.Name)
	}
	if titleField.Type != types.GraphQLString {
		t.Fatalf("titleField.Type expected to equal types.GraphQLString, got: %v", titleField.Type)
	}
	if titleField.Type.GetName() != "String" {
		t.Fatalf("titleField.Type.GetName() expected to equal `String`, got: %v", titleField.Type.GetName())
	}

	authorField := articleFieldTypeObject.GetFields()["author"]
	if authorField == nil {
		t.Fatalf("authorField is nil")
	}
	authorFieldObject, ok := authorField.Type.(*types.GraphQLObjectType)
	if !ok {
		t.Fatalf("expected authorField.Type to be *types.GraphQLObjectType`, got: %v", authorField)
	}

	recentArticleField := authorFieldObject.GetFields()["recentArticle"]
	if recentArticleField == nil {
		t.Fatalf("recentArticleField is nil")
	}
	if recentArticleField.Type != blogArticle {
		t.Fatalf("recentArticleField.Type expected to equal blogArticle, got: %v", recentArticleField.Type)
	}

	feedField := blogQuery.GetFields()["feed"]
	feedFieldList, ok := feedField.Type.(*types.GraphQLList)
	if !ok {
		t.Fatalf("expected feedFieldList to be *types.GraphQLList`, got: %v", authorField)
	}
	if feedFieldList.OfType != blogArticle {
		t.Fatalf("feedFieldList.OfType expected to equal blogArticle, got: %v", feedFieldList.OfType)
	}
	if feedField.Name != "feed" {
		t.Fatalf("feedField.Name expected to equal `feed`, got: %v", feedField.Name)
	}
}
func TestTypeSystem_DefinitionExample_DefinesAMutationScheme(t *testing.T) {
	blogSchema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query:    blogQuery,
		Mutation: blogMutation,
	})
	if err != nil {
		t.Fatalf("unexpected error, got: %v", err)
	}

	if blogSchema.GetMutationType() != blogMutation {
		t.Fatalf("expected blogSchema.GetMutationType() == blogMutation")
	}

	writeMutation, _ := blogMutation.GetFields()["writeArticle"]
	if writeMutation == nil {
		t.Fatalf("writeMutation is nil")
	}
	writeMutationType := writeMutation.Type
	if writeMutationType != blogArticle {
		t.Fatalf("writeMutationType expected to equal blogArticle, got: %v", writeMutationType)
	}
	if writeMutationType.GetName() != "Article" {
		t.Fatalf("writeMutationType.Name expected to equal `Article`, got: %v", writeMutationType.GetName())
	}
	if writeMutation.Name != "writeArticle" {
		t.Fatalf("writeMutation.Name expected to equal `writeArticle`, got: %v", writeMutation.Name)
	}
}

func TestTypeSystem_DefinitionExample_IncludesNestedInputObjectsInTheMap(t *testing.T) {
	nestedInputObject := types.NewGraphQLInputObjectType(types.InputObjectConfig{
		Name: "NestedInputObject",
		Fields: types.InputObjectConfigFieldMap{
			"value": &types.InputObjectFieldConfig{
				Type: types.GraphQLString,
			},
		},
	})
	someInputObject := types.NewGraphQLInputObjectType(types.InputObjectConfig{
		Name: "SomeInputObject",
		Fields: types.InputObjectConfigFieldMap{
			"nested": &types.InputObjectFieldConfig{
				Type: nestedInputObject,
			},
		},
	})
	someMutation := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "SomeMutation",
		Fields: types.GraphQLFieldConfigMap{
			"mutateSomething": &types.GraphQLFieldConfig{
				Type: blogArticle,
				Args: types.GraphQLFieldConfigArgumentMap{
					"input": &types.GraphQLArgumentConfig{
						Type: someInputObject,
					},
				},
			},
		},
	})
	schema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query:    blogQuery,
		Mutation: someMutation,
	})
	if err != nil {
		t.Fatalf("unexpected error, got: %v", err)
	}
	if schema.GetType("NestedInputObject") != nestedInputObject {
		t.Fatalf(`schema.GetType("NestedInputObject") expected to equal nestedInputObject, got: %v`, schema.GetType("NestedInputObject"))
	}
}

func TestTypeSystem_DefinitionExample_IncludesInterfacesSubTypesInTheTypeMap(t *testing.T) {

	someInterface := types.NewGraphQLInterfaceType(types.GraphQLInterfaceTypeConfig{
		Name: "SomeInterface",
		Fields: types.GraphQLFieldConfigMap{
			"f": &types.GraphQLFieldConfig{
				Type: types.GraphQLInt,
			},
		},
	})

	someSubType := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "SomeSubtype",
		Fields: types.GraphQLFieldConfigMap{
			"f": &types.GraphQLFieldConfig{
				Type: types.GraphQLInt,
			},
		},
		Interfaces: []*types.GraphQLInterfaceType{someInterface},
		IsTypeOf: func(value interface{}, info types.GraphQLResolveInfo) bool {
			return true
		},
	})
	schema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "Query",
			Fields: types.GraphQLFieldConfigMap{
				"iface": &types.GraphQLFieldConfig{
					Type: someInterface,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("unexpected error, got: %v", err)
	}
	if schema.GetType("SomeSubtype") != someSubType {
		t.Fatalf(`schema.GetType("SomeSubtype") expected to equal someSubType, got: %v`, schema.GetType("SomeSubtype"))
	}
}

func TestTypeSystem_DefinitionExample_IncludesInterfacesThunkSubtypesInTheTypeMap(t *testing.T) {

	someInterface := types.NewGraphQLInterfaceType(types.GraphQLInterfaceTypeConfig{
		Name: "SomeInterface",
		Fields: types.GraphQLFieldConfigMap{
			"f": &types.GraphQLFieldConfig{
				Type: types.GraphQLInt,
			},
		},
	})

	someSubType := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "SomeSubtype",
		Fields: types.GraphQLFieldConfigMap{
			"f": &types.GraphQLFieldConfig{
				Type: types.GraphQLInt,
			},
		},
		Interfaces: (types.GraphQLInterfacesThunk)(func() []*types.GraphQLInterfaceType {
			return []*types.GraphQLInterfaceType{someInterface}
		}),
		IsTypeOf: func(value interface{}, info types.GraphQLResolveInfo) bool {
			return true
		},
	})
	schema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
			Name: "Query",
			Fields: types.GraphQLFieldConfigMap{
				"iface": &types.GraphQLFieldConfig{
					Type: someInterface,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("unexpected error, got: %v", err)
	}
	if schema.GetType("SomeSubtype") != someSubType {
		t.Fatalf(`schema.GetType("SomeSubtype") expected to equal someSubType, got: %v`, schema.GetType("SomeSubtype"))
	}
}

func TestTypeSystem_DefinitionExample_StringifiesSimpleTypes(t *testing.T) {

	type Test struct {
		ttype    types.GraphQLType
		expected string
	}
	tests := []Test{
		Test{types.GraphQLInt, "Int"},
		Test{blogArticle, "Article"},
		Test{interfaceType, "Interface"},
		Test{unionType, "Union"},
		Test{enumType, "Enum"},
		Test{inputObjectType, "InputObject"},
		Test{types.NewGraphQLNonNull(types.GraphQLInt), "Int!"},
		Test{types.NewGraphQLList(types.GraphQLInt), "[Int]"},
		Test{types.NewGraphQLNonNull(types.NewGraphQLList(types.GraphQLInt)), "[Int]!"},
		Test{types.NewGraphQLList(types.NewGraphQLNonNull(types.GraphQLInt)), "[Int!]"},
		Test{types.NewGraphQLList(types.NewGraphQLList(types.GraphQLInt)), "[[Int]]"},
	}
	for _, test := range tests {
		ttypeStr := fmt.Sprintf("%v", test.ttype)
		if ttypeStr != test.expected {
			t.Fatalf(`expected %v , got: %v`, test.expected, ttypeStr)
		}
	}
}

func TestTypeSystem_DefinitionExample_IdentifiesInputTypes(t *testing.T) {
	type Test struct {
		ttype    types.GraphQLType
		expected bool
	}
	tests := []Test{
		Test{types.GraphQLInt, true},
		Test{objectType, false},
		Test{interfaceType, false},
		Test{unionType, false},
		Test{enumType, true},
		Test{inputObjectType, true},
	}
	for _, test := range tests {
		ttypeStr := fmt.Sprintf("%v", test.ttype)
		if types.IsInputType(test.ttype) != test.expected {
			t.Fatalf(`expected %v , got: %v`, test.expected, ttypeStr)
		}
		if types.IsInputType(types.NewGraphQLList(test.ttype)) != test.expected {
			t.Fatalf(`expected %v , got: %v`, test.expected, ttypeStr)
		}
		if types.IsInputType(types.NewGraphQLNonNull(test.ttype)) != test.expected {
			t.Fatalf(`expected %v , got: %v`, test.expected, ttypeStr)
		}
	}
}

func TestTypeSystem_DefinitionExample_IdentifiesOutputTypes(t *testing.T) {
	type Test struct {
		ttype    types.GraphQLType
		expected bool
	}
	tests := []Test{
		Test{types.GraphQLInt, true},
		Test{objectType, true},
		Test{interfaceType, true},
		Test{unionType, true},
		Test{enumType, true},
		Test{inputObjectType, false},
	}
	for _, test := range tests {
		ttypeStr := fmt.Sprintf("%v", test.ttype)
		if types.IsOutputType(test.ttype) != test.expected {
			t.Fatalf(`expected %v , got: %v`, test.expected, ttypeStr)
		}
		if types.IsOutputType(types.NewGraphQLList(test.ttype)) != test.expected {
			t.Fatalf(`expected %v , got: %v`, test.expected, ttypeStr)
		}
		if types.IsOutputType(types.NewGraphQLNonNull(test.ttype)) != test.expected {
			t.Fatalf(`expected %v , got: %v`, test.expected, ttypeStr)
		}
	}
}

func TestTypeSystem_DefinitionExample_ProhibitsNestingNonNullInsideNonNull(t *testing.T) {
	ttype := types.NewGraphQLNonNull(types.NewGraphQLNonNull(types.GraphQLInt))
	expected := `Can only create NonNull of a Nullable GraphQLType but got: Int!.`
	if ttype.GetError().Error() != expected {
		t.Fatalf(`expected %v , got: %v`, expected, ttype.GetError())
	}
}
func TestTypeSystem_DefinitionExample_ProhibitsNilInNonNull(t *testing.T) {
	ttype := types.NewGraphQLNonNull(nil)
	expected := `Can only create NonNull of a Nullable GraphQLType but got: <nil>.`
	if ttype.GetError().Error() != expected {
		t.Fatalf(`expected %v , got: %v`, expected, ttype.GetError())
	}
}
func TestTypeSystem_DefinitionExample_ProhibitsNilTypeInUnions(t *testing.T) {
	ttype := types.NewGraphQLUnionType(types.GraphQLUnionTypeConfig{
		Name:  "BadUnion",
		Types: []*types.GraphQLObjectType{nil},
	})
	expected := `BadUnion may only contain Object types, it cannot contain: <nil>.`
	if ttype.GetError().Error() != expected {
		t.Fatalf(`expected %v , got: %v`, expected, ttype.GetError())
	}
}
func TestTypeSystem_DefinitionExample_DoesNotMutatePassedFieldDefinitions(t *testing.T) {
	fields := types.GraphQLFieldConfigMap{
		"field1": &types.GraphQLFieldConfig{
			Type: types.GraphQLString,
		},
		"field2": &types.GraphQLFieldConfig{
			Type: types.GraphQLString,
			Args: types.GraphQLFieldConfigArgumentMap{
				"id": &types.GraphQLArgumentConfig{
					Type: types.GraphQLString,
				},
			},
		},
	}
	testObject1 := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name:   "Test1",
		Fields: fields,
	})
	testObject2 := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name:   "Test2",
		Fields: fields,
	})
	if !reflect.DeepEqual(testObject1.GetFields(), testObject2.GetFields()) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(testObject1.GetFields(), testObject2.GetFields()))
	}

	expectedFields := types.GraphQLFieldConfigMap{
		"field1": &types.GraphQLFieldConfig{
			Type: types.GraphQLString,
		},
		"field2": &types.GraphQLFieldConfig{
			Type: types.GraphQLString,
			Args: types.GraphQLFieldConfigArgumentMap{
				"id": &types.GraphQLArgumentConfig{
					Type: types.GraphQLString,
				},
			},
		},
	}
	if !reflect.DeepEqual(fields, expectedFields) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedFields, fields))
	}

	inputFields := types.InputObjectConfigFieldMap{
		"field1": &types.InputObjectFieldConfig{
			Type: types.GraphQLString,
		},
		"field2": &types.InputObjectFieldConfig{
			Type: types.GraphQLString,
		},
	}
	expectedInputFields := types.InputObjectConfigFieldMap{
		"field1": &types.InputObjectFieldConfig{
			Type: types.GraphQLString,
		},
		"field2": &types.InputObjectFieldConfig{
			Type: types.GraphQLString,
		},
	}
	testInputObject1 := types.NewGraphQLInputObjectType(types.InputObjectConfig{
		Name:   "Test1",
		Fields: inputFields,
	})
	testInputObject2 := types.NewGraphQLInputObjectType(types.InputObjectConfig{
		Name:   "Test2",
		Fields: inputFields,
	})
	if !reflect.DeepEqual(testInputObject1.GetFields(), testInputObject2.GetFields()) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(testInputObject1.GetFields(), testInputObject2.GetFields()))
	}
	if !reflect.DeepEqual(inputFields, expectedInputFields) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedInputFields, fields))
	}

}
