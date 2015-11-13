package graphql_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/testutil"
)

var blogImage = graphql.NewObject(graphql.ObjectConfig{
	Name: "Image",
	Fields: graphql.FieldConfigMap{
		"url": &graphql.FieldConfig{
			Type: graphql.String,
		},
		"width": &graphql.FieldConfig{
			Type: graphql.Int,
		},
		"height": &graphql.FieldConfig{
			Type: graphql.Int,
		},
	},
})
var blogAuthor = graphql.NewObject(graphql.ObjectConfig{
	Name: "Author",
	Fields: graphql.FieldConfigMap{
		"id": &graphql.FieldConfig{
			Type: graphql.String,
		},
		"name": &graphql.FieldConfig{
			Type: graphql.String,
		},
		"pic": &graphql.FieldConfig{
			Type: blogImage,
			Args: graphql.FieldConfigArgument{
				"width": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"height": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
		},
		"recentArticle": &graphql.FieldConfig{},
	},
})
var blogArticle = graphql.NewObject(graphql.ObjectConfig{
	Name: "Article",
	Fields: graphql.FieldConfigMap{
		"id": &graphql.FieldConfig{
			Type: graphql.String,
		},
		"isPublished": &graphql.FieldConfig{
			Type: graphql.Boolean,
		},
		"author": &graphql.FieldConfig{
			Type: blogAuthor,
		},
		"title": &graphql.FieldConfig{
			Type: graphql.String,
		},
		"body": &graphql.FieldConfig{
			Type: graphql.String,
		},
	},
})
var blogQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.FieldConfigMap{
		"article": &graphql.FieldConfig{
			Type: blogArticle,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
		},
		"feed": &graphql.FieldConfig{
			Type: graphql.NewList(blogArticle),
		},
	},
})

var blogMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.FieldConfigMap{
		"writeArticle": &graphql.FieldConfig{
			Type: blogArticle,
		},
	},
})

var objectType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Object",
	IsTypeOf: func(value interface{}, info graphql.ResolveInfo) bool {
		return true
	},
})
var interfaceType = graphql.NewInterface(graphql.InterfaceConfig{
	Name: "Interface",
})
var unionType = graphql.NewUnion(graphql.UnionConfig{
	Name: "Union",
	Types: []*graphql.Object{
		objectType,
	},
})
var enumType = graphql.NewEnum(graphql.EnumConfig{
	Name: "Enum",
	Values: graphql.EnumValueConfigMap{
		"foo": &graphql.EnumValueConfig{},
	},
})
var inputObjectType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "InputObject",
})

func init() {
	blogAuthor.AddFieldConfig("recentArticle", &graphql.FieldConfig{
		Type: blogArticle,
	})
}

func TestTypeSystem_DefinitionExample_DefinesAQueryOnlySchema(t *testing.T) {
	blogSchema, err := graphql.NewSchema(graphql.SchemaConfig{
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
	articleFieldTypeObject, ok := articleFieldType.(*graphql.Object)
	if !ok {
		t.Fatalf("expected articleFieldType to be graphql.Object`, got: %v", articleField)
	}

	// TODO: expose a Object.GetField(key string), instead of this ghetto way of accessing a field map?
	titleField := articleFieldTypeObject.GetFields()["title"]
	if titleField == nil {
		t.Fatalf("titleField is nil")
	}
	if titleField.Name != "title" {
		t.Fatalf("titleField.Name expected to equal title, got: %v", titleField.Name)
	}
	if titleField.Type != graphql.String {
		t.Fatalf("titleField.Type expected to equal graphql.String, got: %v", titleField.Type)
	}
	if titleField.Type.GetName() != "String" {
		t.Fatalf("titleField.Type.GetName() expected to equal `String`, got: %v", titleField.Type.GetName())
	}

	authorField := articleFieldTypeObject.GetFields()["author"]
	if authorField == nil {
		t.Fatalf("authorField is nil")
	}
	authorFieldObject, ok := authorField.Type.(*graphql.Object)
	if !ok {
		t.Fatalf("expected authorField.Type to be Object`, got: %v", authorField)
	}

	recentArticleField := authorFieldObject.GetFields()["recentArticle"]
	if recentArticleField == nil {
		t.Fatalf("recentArticleField is nil")
	}
	if recentArticleField.Type != blogArticle {
		t.Fatalf("recentArticleField.Type expected to equal blogArticle, got: %v", recentArticleField.Type)
	}

	feedField := blogQuery.GetFields()["feed"]
	feedFieldList, ok := feedField.Type.(*graphql.List)
	if !ok {
		t.Fatalf("expected feedFieldList to be List`, got: %v", authorField)
	}
	if feedFieldList.OfType != blogArticle {
		t.Fatalf("feedFieldList.OfType expected to equal blogArticle, got: %v", feedFieldList.OfType)
	}
	if feedField.Name != "feed" {
		t.Fatalf("feedField.Name expected to equal `feed`, got: %v", feedField.Name)
	}
}
func TestTypeSystem_DefinitionExample_DefinesAMutationScheme(t *testing.T) {
	blogSchema, err := graphql.NewSchema(graphql.SchemaConfig{
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
	nestedInputObject := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "NestedInputObject",
		Fields: graphql.InputObjectConfigFieldMap{
			"value": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
		},
	})
	someInputObject := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "SomeInputObject",
		Fields: graphql.InputObjectConfigFieldMap{
			"nested": &graphql.InputObjectFieldConfig{
				Type: nestedInputObject,
			},
		},
	})
	someMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "SomeMutation",
		Fields: graphql.FieldConfigMap{
			"mutateSomething": &graphql.FieldConfig{
				Type: blogArticle,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: someInputObject,
					},
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
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

	someInterface := graphql.NewInterface(graphql.InterfaceConfig{
		Name: "SomeInterface",
		Fields: graphql.FieldConfigMap{
			"f": &graphql.FieldConfig{
				Type: graphql.Int,
			},
		},
	})

	someSubType := graphql.NewObject(graphql.ObjectConfig{
		Name: "SomeSubtype",
		Fields: graphql.FieldConfigMap{
			"f": &graphql.FieldConfig{
				Type: graphql.Int,
			},
		},
		Interfaces: []*graphql.Interface{someInterface},
		IsTypeOf: func(value interface{}, info graphql.ResolveInfo) bool {
			return true
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.FieldConfigMap{
				"iface": &graphql.FieldConfig{
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

	someInterface := graphql.NewInterface(graphql.InterfaceConfig{
		Name: "SomeInterface",
		Fields: graphql.FieldConfigMap{
			"f": &graphql.FieldConfig{
				Type: graphql.Int,
			},
		},
	})

	someSubType := graphql.NewObject(graphql.ObjectConfig{
		Name: "SomeSubtype",
		Fields: graphql.FieldConfigMap{
			"f": &graphql.FieldConfig{
				Type: graphql.Int,
			},
		},
		Interfaces: (graphql.InterfacesThunk)(func() []*graphql.Interface {
			return []*graphql.Interface{someInterface}
		}),
		IsTypeOf: func(value interface{}, info graphql.ResolveInfo) bool {
			return true
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.FieldConfigMap{
				"iface": &graphql.FieldConfig{
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
		ttype    graphql.Type
		expected string
	}
	tests := []Test{
		Test{graphql.Int, "Int"},
		Test{blogArticle, "Article"},
		Test{interfaceType, "Interface"},
		Test{unionType, "Union"},
		Test{enumType, "Enum"},
		Test{inputObjectType, "InputObject"},
		Test{graphql.NewNonNull(graphql.Int), "Int!"},
		Test{graphql.NewList(graphql.Int), "[Int]"},
		Test{graphql.NewNonNull(graphql.NewList(graphql.Int)), "[Int]!"},
		Test{graphql.NewList(graphql.NewNonNull(graphql.Int)), "[Int!]"},
		Test{graphql.NewList(graphql.NewList(graphql.Int)), "[[Int]]"},
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
		ttype    graphql.Type
		expected bool
	}
	tests := []Test{
		Test{graphql.Int, true},
		Test{objectType, false},
		Test{interfaceType, false},
		Test{unionType, false},
		Test{enumType, true},
		Test{inputObjectType, true},
	}
	for _, test := range tests {
		ttypeStr := fmt.Sprintf("%v", test.ttype)
		if graphql.IsInputType(test.ttype) != test.expected {
			t.Fatalf(`expected %v , got: %v`, test.expected, ttypeStr)
		}
		if graphql.IsInputType(graphql.NewList(test.ttype)) != test.expected {
			t.Fatalf(`expected %v , got: %v`, test.expected, ttypeStr)
		}
		if graphql.IsInputType(graphql.NewNonNull(test.ttype)) != test.expected {
			t.Fatalf(`expected %v , got: %v`, test.expected, ttypeStr)
		}
	}
}

func TestTypeSystem_DefinitionExample_IdentifiesOutputTypes(t *testing.T) {
	type Test struct {
		ttype    graphql.Type
		expected bool
	}
	tests := []Test{
		Test{graphql.Int, true},
		Test{objectType, true},
		Test{interfaceType, true},
		Test{unionType, true},
		Test{enumType, true},
		Test{inputObjectType, false},
	}
	for _, test := range tests {
		ttypeStr := fmt.Sprintf("%v", test.ttype)
		if graphql.IsOutputType(test.ttype) != test.expected {
			t.Fatalf(`expected %v , got: %v`, test.expected, ttypeStr)
		}
		if graphql.IsOutputType(graphql.NewList(test.ttype)) != test.expected {
			t.Fatalf(`expected %v , got: %v`, test.expected, ttypeStr)
		}
		if graphql.IsOutputType(graphql.NewNonNull(test.ttype)) != test.expected {
			t.Fatalf(`expected %v , got: %v`, test.expected, ttypeStr)
		}
	}
}

func TestTypeSystem_DefinitionExample_ProhibitsNestingNonNullInsideNonNull(t *testing.T) {
	ttype := graphql.NewNonNull(graphql.NewNonNull(graphql.Int))
	expected := `Can only create NonNull of a Nullable Type but got: Int!.`
	if ttype.GetError().Error() != expected {
		t.Fatalf(`expected %v , got: %v`, expected, ttype.GetError())
	}
}
func TestTypeSystem_DefinitionExample_ProhibitsNilInNonNull(t *testing.T) {
	ttype := graphql.NewNonNull(nil)
	expected := `Can only create NonNull of a Nullable Type but got: <nil>.`
	if ttype.GetError().Error() != expected {
		t.Fatalf(`expected %v , got: %v`, expected, ttype.GetError())
	}
}
func TestTypeSystem_DefinitionExample_ProhibitsNilTypeInUnions(t *testing.T) {
	ttype := graphql.NewUnion(graphql.UnionConfig{
		Name:  "BadUnion",
		Types: []*graphql.Object{nil},
	})
	expected := `BadUnion may only contain Object types, it cannot contain: <nil>.`
	if ttype.GetError().Error() != expected {
		t.Fatalf(`expected %v , got: %v`, expected, ttype.GetError())
	}
}
func TestTypeSystem_DefinitionExample_DoesNotMutatePassedFieldDefinitions(t *testing.T) {
	fields := graphql.FieldConfigMap{
		"field1": &graphql.FieldConfig{
			Type: graphql.String,
		},
		"field2": &graphql.FieldConfig{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
		},
	}
	testObject1 := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Test1",
		Fields: fields,
	})
	testObject2 := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Test2",
		Fields: fields,
	})
	if !reflect.DeepEqual(testObject1.GetFields(), testObject2.GetFields()) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(testObject1.GetFields(), testObject2.GetFields()))
	}

	expectedFields := graphql.FieldConfigMap{
		"field1": &graphql.FieldConfig{
			Type: graphql.String,
		},
		"field2": &graphql.FieldConfig{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
		},
	}
	if !reflect.DeepEqual(fields, expectedFields) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedFields, fields))
	}

	inputFields := graphql.InputObjectConfigFieldMap{
		"field1": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"field2": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
	}
	expectedInputFields := graphql.InputObjectConfigFieldMap{
		"field1": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"field2": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
	}
	testInputObject1 := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:   "Test1",
		Fields: inputFields,
	})
	testInputObject2 := graphql.NewInputObject(graphql.InputObjectConfig{
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

func TestTypeSystem_DefinitionExampe_AllowsCyclicFieldTypes(t *testing.T) {
	personType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Person",
		Fields: (graphql.ObjectFieldMapThunk)(func() graphql.FieldConfigMap {
			return graphql.FieldConfigMap{
				"name": &graphql.FieldConfig{
					Type: graphql.String,
				},
				"bestFriend": &graphql.FieldConfig{
					Type: personType,
				},
			}
		}),
	})

	fieldMap := personType.GetFields()
	if !reflect.DeepEqual(fieldMap["name"].Type, graphql.String) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(fieldMap["bestFriend"].Type, personType))
	}

}
