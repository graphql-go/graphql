package types_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/chris-ramon/graphql/testutil"
	"github.com/chris-ramon/graphql/types"
)

var blogImage = types.NewObject(types.ObjectConfig{
	Name: "Image",
	Fields: types.FieldConfigMap{
		"url": &types.FieldConfig{
			Type: types.String,
		},
		"width": &types.FieldConfig{
			Type: types.Int,
		},
		"height": &types.FieldConfig{
			Type: types.Int,
		},
	},
})
var blogAuthor = types.NewObject(types.ObjectConfig{
	Name: "Author",
	Fields: types.FieldConfigMap{
		"id": &types.FieldConfig{
			Type: types.String,
		},
		"name": &types.FieldConfig{
			Type: types.String,
		},
		"pic": &types.FieldConfig{
			Type: blogImage,
			Args: types.FieldConfigArgument{
				"width": &types.ArgumentConfig{
					Type: types.Int,
				},
				"height": &types.ArgumentConfig{
					Type: types.Int,
				},
			},
		},
		"recentArticle": &types.FieldConfig{},
	},
})
var blogArticle = types.NewObject(types.ObjectConfig{
	Name: "Article",
	Fields: types.FieldConfigMap{
		"id": &types.FieldConfig{
			Type: types.String,
		},
		"isPublished": &types.FieldConfig{
			Type: types.Boolean,
		},
		"author": &types.FieldConfig{
			Type: blogAuthor,
		},
		"title": &types.FieldConfig{
			Type: types.String,
		},
		"body": &types.FieldConfig{
			Type: types.String,
		},
	},
})
var blogQuery = types.NewObject(types.ObjectConfig{
	Name: "Query",
	Fields: types.FieldConfigMap{
		"article": &types.FieldConfig{
			Type: blogArticle,
			Args: types.FieldConfigArgument{
				"id": &types.ArgumentConfig{
					Type: types.String,
				},
			},
		},
		"feed": &types.FieldConfig{
			Type: types.NewList(blogArticle),
		},
	},
})

var blogMutation = types.NewObject(types.ObjectConfig{
	Name: "Mutation",
	Fields: types.FieldConfigMap{
		"writeArticle": &types.FieldConfig{
			Type: blogArticle,
		},
	},
})

var objectType = types.NewObject(types.ObjectConfig{
	Name: "Object",
	IsTypeOf: func(value interface{}, info types.ResolveInfo) bool {
		return true
	},
})
var interfaceType = types.NewInterface(types.InterfaceConfig{
	Name: "Interface",
})
var unionType = types.NewUnion(types.UnionConfig{
	Name: "Union",
	Types: []*types.Object{
		objectType,
	},
})
var enumType = types.NewEnum(types.EnumConfig{
	Name: "Enum",
	Values: types.EnumValueConfigMap{
		"foo": &types.EnumValueConfig{},
	},
})
var inputObjectType = types.NewInputObject(types.InputObjectConfig{
	Name: "InputObject",
})

func init() {
	blogAuthor.AddFieldConfig("recentArticle", &types.FieldConfig{
		Type: blogArticle,
	})
}

func TestTypeSystem_DefinitionExample_DefinesAQueryOnlySchema(t *testing.T) {
	blogSchema, err := types.NewSchema(types.SchemaConfig{
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
	articleFieldTypeObject, ok := articleFieldType.(*types.Object)
	if !ok {
		t.Fatalf("expected articleFieldType to be *types.Object`, got: %v", articleField)
	}

	// TODO: expose a Object.GetField(key string), instead of this ghetto way of accessing a field map?
	titleField := articleFieldTypeObject.GetFields()["title"]
	if titleField == nil {
		t.Fatalf("titleField is nil")
	}
	if titleField.Name != "title" {
		t.Fatalf("titleField.Name expected to equal title, got: %v", titleField.Name)
	}
	if titleField.Type != types.String {
		t.Fatalf("titleField.Type expected to equal types.String, got: %v", titleField.Type)
	}
	if titleField.Type.GetName() != "String" {
		t.Fatalf("titleField.Type.GetName() expected to equal `String`, got: %v", titleField.Type.GetName())
	}

	authorField := articleFieldTypeObject.GetFields()["author"]
	if authorField == nil {
		t.Fatalf("authorField is nil")
	}
	authorFieldObject, ok := authorField.Type.(*types.Object)
	if !ok {
		t.Fatalf("expected authorField.Type to be *types.Object`, got: %v", authorField)
	}

	recentArticleField := authorFieldObject.GetFields()["recentArticle"]
	if recentArticleField == nil {
		t.Fatalf("recentArticleField is nil")
	}
	if recentArticleField.Type != blogArticle {
		t.Fatalf("recentArticleField.Type expected to equal blogArticle, got: %v", recentArticleField.Type)
	}

	feedField := blogQuery.GetFields()["feed"]
	feedFieldList, ok := feedField.Type.(*types.List)
	if !ok {
		t.Fatalf("expected feedFieldList to be *types.List`, got: %v", authorField)
	}
	if feedFieldList.OfType != blogArticle {
		t.Fatalf("feedFieldList.OfType expected to equal blogArticle, got: %v", feedFieldList.OfType)
	}
	if feedField.Name != "feed" {
		t.Fatalf("feedField.Name expected to equal `feed`, got: %v", feedField.Name)
	}
}
func TestTypeSystem_DefinitionExample_DefinesAMutationScheme(t *testing.T) {
	blogSchema, err := types.NewSchema(types.SchemaConfig{
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
	nestedInputObject := types.NewInputObject(types.InputObjectConfig{
		Name: "NestedInputObject",
		Fields: types.InputObjectConfigFieldMap{
			"value": &types.InputObjectFieldConfig{
				Type: types.String,
			},
		},
	})
	someInputObject := types.NewInputObject(types.InputObjectConfig{
		Name: "SomeInputObject",
		Fields: types.InputObjectConfigFieldMap{
			"nested": &types.InputObjectFieldConfig{
				Type: nestedInputObject,
			},
		},
	})
	someMutation := types.NewObject(types.ObjectConfig{
		Name: "SomeMutation",
		Fields: types.FieldConfigMap{
			"mutateSomething": &types.FieldConfig{
				Type: blogArticle,
				Args: types.FieldConfigArgument{
					"input": &types.ArgumentConfig{
						Type: someInputObject,
					},
				},
			},
		},
	})
	schema, err := types.NewSchema(types.SchemaConfig{
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

	someInterface := types.NewInterface(types.InterfaceConfig{
		Name: "SomeInterface",
		Fields: types.FieldConfigMap{
			"f": &types.FieldConfig{
				Type: types.Int,
			},
		},
	})

	someSubType := types.NewObject(types.ObjectConfig{
		Name: "SomeSubtype",
		Fields: types.FieldConfigMap{
			"f": &types.FieldConfig{
				Type: types.Int,
			},
		},
		Interfaces: []*types.Interface{someInterface},
		IsTypeOf: func(value interface{}, info types.ResolveInfo) bool {
			return true
		},
	})
	schema, err := types.NewSchema(types.SchemaConfig{
		Query: types.NewObject(types.ObjectConfig{
			Name: "Query",
			Fields: types.FieldConfigMap{
				"iface": &types.FieldConfig{
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

	someInterface := types.NewInterface(types.InterfaceConfig{
		Name: "SomeInterface",
		Fields: types.FieldConfigMap{
			"f": &types.FieldConfig{
				Type: types.Int,
			},
		},
	})

	someSubType := types.NewObject(types.ObjectConfig{
		Name: "SomeSubtype",
		Fields: types.FieldConfigMap{
			"f": &types.FieldConfig{
				Type: types.Int,
			},
		},
		Interfaces: (types.InterfacesThunk)(func() []*types.Interface {
			return []*types.Interface{someInterface}
		}),
		IsTypeOf: func(value interface{}, info types.ResolveInfo) bool {
			return true
		},
	})
	schema, err := types.NewSchema(types.SchemaConfig{
		Query: types.NewObject(types.ObjectConfig{
			Name: "Query",
			Fields: types.FieldConfigMap{
				"iface": &types.FieldConfig{
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
		ttype    types.Type
		expected string
	}
	tests := []Test{
		Test{types.Int, "Int"},
		Test{blogArticle, "Article"},
		Test{interfaceType, "Interface"},
		Test{unionType, "Union"},
		Test{enumType, "Enum"},
		Test{inputObjectType, "InputObject"},
		Test{types.NewNonNull(types.Int), "Int!"},
		Test{types.NewList(types.Int), "[Int]"},
		Test{types.NewNonNull(types.NewList(types.Int)), "[Int]!"},
		Test{types.NewList(types.NewNonNull(types.Int)), "[Int!]"},
		Test{types.NewList(types.NewList(types.Int)), "[[Int]]"},
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
		ttype    types.Type
		expected bool
	}
	tests := []Test{
		Test{types.Int, true},
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
		if types.IsInputType(types.NewList(test.ttype)) != test.expected {
			t.Fatalf(`expected %v , got: %v`, test.expected, ttypeStr)
		}
		if types.IsInputType(types.NewNonNull(test.ttype)) != test.expected {
			t.Fatalf(`expected %v , got: %v`, test.expected, ttypeStr)
		}
	}
}

func TestTypeSystem_DefinitionExample_IdentifiesOutputTypes(t *testing.T) {
	type Test struct {
		ttype    types.Type
		expected bool
	}
	tests := []Test{
		Test{types.Int, true},
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
		if types.IsOutputType(types.NewList(test.ttype)) != test.expected {
			t.Fatalf(`expected %v , got: %v`, test.expected, ttypeStr)
		}
		if types.IsOutputType(types.NewNonNull(test.ttype)) != test.expected {
			t.Fatalf(`expected %v , got: %v`, test.expected, ttypeStr)
		}
	}
}

func TestTypeSystem_DefinitionExample_ProhibitsNestingNonNullInsideNonNull(t *testing.T) {
	ttype := types.NewNonNull(types.NewNonNull(types.Int))
	expected := `Can only create NonNull of a Nullable Type but got: Int!.`
	if ttype.GetError().Error() != expected {
		t.Fatalf(`expected %v , got: %v`, expected, ttype.GetError())
	}
}
func TestTypeSystem_DefinitionExample_ProhibitsNilInNonNull(t *testing.T) {
	ttype := types.NewNonNull(nil)
	expected := `Can only create NonNull of a Nullable Type but got: <nil>.`
	if ttype.GetError().Error() != expected {
		t.Fatalf(`expected %v , got: %v`, expected, ttype.GetError())
	}
}
func TestTypeSystem_DefinitionExample_ProhibitsNilTypeInUnions(t *testing.T) {
	ttype := types.NewUnion(types.UnionConfig{
		Name:  "BadUnion",
		Types: []*types.Object{nil},
	})
	expected := `BadUnion may only contain Object types, it cannot contain: <nil>.`
	if ttype.GetError().Error() != expected {
		t.Fatalf(`expected %v , got: %v`, expected, ttype.GetError())
	}
}
func TestTypeSystem_DefinitionExample_DoesNotMutatePassedFieldDefinitions(t *testing.T) {
	fields := types.FieldConfigMap{
		"field1": &types.FieldConfig{
			Type: types.String,
		},
		"field2": &types.FieldConfig{
			Type: types.String,
			Args: types.FieldConfigArgument{
				"id": &types.ArgumentConfig{
					Type: types.String,
				},
			},
		},
	}
	testObject1 := types.NewObject(types.ObjectConfig{
		Name:   "Test1",
		Fields: fields,
	})
	testObject2 := types.NewObject(types.ObjectConfig{
		Name:   "Test2",
		Fields: fields,
	})
	if !reflect.DeepEqual(testObject1.GetFields(), testObject2.GetFields()) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(testObject1.GetFields(), testObject2.GetFields()))
	}

	expectedFields := types.FieldConfigMap{
		"field1": &types.FieldConfig{
			Type: types.String,
		},
		"field2": &types.FieldConfig{
			Type: types.String,
			Args: types.FieldConfigArgument{
				"id": &types.ArgumentConfig{
					Type: types.String,
				},
			},
		},
	}
	if !reflect.DeepEqual(fields, expectedFields) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedFields, fields))
	}

	inputFields := types.InputObjectConfigFieldMap{
		"field1": &types.InputObjectFieldConfig{
			Type: types.String,
		},
		"field2": &types.InputObjectFieldConfig{
			Type: types.String,
		},
	}
	expectedInputFields := types.InputObjectConfigFieldMap{
		"field1": &types.InputObjectFieldConfig{
			Type: types.String,
		},
		"field2": &types.InputObjectFieldConfig{
			Type: types.String,
		},
	}
	testInputObject1 := types.NewInputObject(types.InputObjectConfig{
		Name:   "Test1",
		Fields: inputFields,
	})
	testInputObject2 := types.NewInputObject(types.InputObjectConfig{
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
