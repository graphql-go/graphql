package graphql

import (
	"fmt"
	"reflect"
	"testing"
)

var blogImage = NewObject(ObjectConfig{
	Name: "Image",
	Fields: FieldConfigMap{
		"url": &FieldConfig{
			Type: String,
		},
		"width": &FieldConfig{
			Type: Int,
		},
		"height": &FieldConfig{
			Type: Int,
		},
	},
})
var blogAuthor = NewObject(ObjectConfig{
	Name: "Author",
	Fields: FieldConfigMap{
		"id": &FieldConfig{
			Type: String,
		},
		"name": &FieldConfig{
			Type: String,
		},
		"pic": &FieldConfig{
			Type: blogImage,
			Args: FieldConfigArgument{
				"width": &ArgumentConfig{
					Type: Int,
				},
				"height": &ArgumentConfig{
					Type: Int,
				},
			},
		},
		"recentArticle": &FieldConfig{},
	},
})
var blogArticle = NewObject(ObjectConfig{
	Name: "Article",
	Fields: FieldConfigMap{
		"id": &FieldConfig{
			Type: String,
		},
		"isPublished": &FieldConfig{
			Type: Boolean,
		},
		"author": &FieldConfig{
			Type: blogAuthor,
		},
		"title": &FieldConfig{
			Type: String,
		},
		"body": &FieldConfig{
			Type: String,
		},
	},
})
var blogQuery = NewObject(ObjectConfig{
	Name: "Query",
	Fields: FieldConfigMap{
		"article": &FieldConfig{
			Type: blogArticle,
			Args: FieldConfigArgument{
				"id": &ArgumentConfig{
					Type: String,
				},
			},
		},
		"feed": &FieldConfig{
			Type: NewList(blogArticle),
		},
	},
})

var blogMutation = NewObject(ObjectConfig{
	Name: "Mutation",
	Fields: FieldConfigMap{
		"writeArticle": &FieldConfig{
			Type: blogArticle,
		},
	},
})

var objectType = NewObject(ObjectConfig{
	Name: "Object",
	IsTypeOf: func(value interface{}, info ResolveInfo) bool {
		return true
	},
})
var interfaceType = NewInterface(InterfaceConfig{
	Name: "Interface",
})
var unionType = NewUnion(UnionConfig{
	Name: "Union",
	Types: []*Object{
		objectType,
	},
})
var enumType = NewEnum(EnumConfig{
	Name: "Enum",
	Values: EnumValueConfigMap{
		"foo": &EnumValueConfig{},
	},
})
var inputObjectType = NewInputObject(InputObjectConfig{
	Name: "InputObject",
})

func init() {
	blogAuthor.AddFieldConfig("recentArticle", &FieldConfig{
		Type: blogArticle,
	})
}

func TestTypeSystem_DefinitionExample_DefinesAQueryOnlySchema(t *testing.T) {
	blogSchema, err := NewSchema(SchemaConfig{
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
	articleFieldTypeObject, ok := articleFieldType.(*Object)
	if !ok {
		t.Fatalf("expected articleFieldType to be *Object`, got: %v", articleField)
	}

	// TODO: expose a Object.GetField(key string), instead of this ghetto way of accessing a field map?
	titleField := articleFieldTypeObject.GetFields()["title"]
	if titleField == nil {
		t.Fatalf("titleField is nil")
	}
	if titleField.Name != "title" {
		t.Fatalf("titleField.Name expected to equal title, got: %v", titleField.Name)
	}
	if titleField.Type != String {
		t.Fatalf("titleField.Type expected to equal String, got: %v", titleField.Type)
	}
	if titleField.Type.GetName() != "String" {
		t.Fatalf("titleField.Type.GetName() expected to equal `String`, got: %v", titleField.Type.GetName())
	}

	authorField := articleFieldTypeObject.GetFields()["author"]
	if authorField == nil {
		t.Fatalf("authorField is nil")
	}
	authorFieldObject, ok := authorField.Type.(*Object)
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
	feedFieldList, ok := feedField.Type.(*List)
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
	blogSchema, err := NewSchema(SchemaConfig{
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
	nestedInputObject := NewInputObject(InputObjectConfig{
		Name: "NestedInputObject",
		Fields: InputObjectConfigFieldMap{
			"value": &InputObjectFieldConfig{
				Type: String,
			},
		},
	})
	someInputObject := NewInputObject(InputObjectConfig{
		Name: "SomeInputObject",
		Fields: InputObjectConfigFieldMap{
			"nested": &InputObjectFieldConfig{
				Type: nestedInputObject,
			},
		},
	})
	someMutation := NewObject(ObjectConfig{
		Name: "SomeMutation",
		Fields: FieldConfigMap{
			"mutateSomething": &FieldConfig{
				Type: blogArticle,
				Args: FieldConfigArgument{
					"input": &ArgumentConfig{
						Type: someInputObject,
					},
				},
			},
		},
	})
	schema, err := NewSchema(SchemaConfig{
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

	someInterface := NewInterface(InterfaceConfig{
		Name: "SomeInterface",
		Fields: FieldConfigMap{
			"f": &FieldConfig{
				Type: Int,
			},
		},
	})

	someSubType := NewObject(ObjectConfig{
		Name: "SomeSubtype",
		Fields: FieldConfigMap{
			"f": &FieldConfig{
				Type: Int,
			},
		},
		Interfaces: []*Interface{someInterface},
		IsTypeOf: func(value interface{}, info ResolveInfo) bool {
			return true
		},
	})
	schema, err := NewSchema(SchemaConfig{
		Query: NewObject(ObjectConfig{
			Name: "Query",
			Fields: FieldConfigMap{
				"iface": &FieldConfig{
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

	someInterface := NewInterface(InterfaceConfig{
		Name: "SomeInterface",
		Fields: FieldConfigMap{
			"f": &FieldConfig{
				Type: Int,
			},
		},
	})

	someSubType := NewObject(ObjectConfig{
		Name: "SomeSubtype",
		Fields: FieldConfigMap{
			"f": &FieldConfig{
				Type: Int,
			},
		},
		Interfaces: (InterfacesThunk)(func() []*Interface {
			return []*Interface{someInterface}
		}),
		IsTypeOf: func(value interface{}, info ResolveInfo) bool {
			return true
		},
	})
	schema, err := NewSchema(SchemaConfig{
		Query: NewObject(ObjectConfig{
			Name: "Query",
			Fields: FieldConfigMap{
				"iface": &FieldConfig{
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
		ttype    Type
		expected string
	}
	tests := []Test{
		Test{Int, "Int"},
		Test{blogArticle, "Article"},
		Test{interfaceType, "Interface"},
		Test{unionType, "Union"},
		Test{enumType, "Enum"},
		Test{inputObjectType, "InputObject"},
		Test{NewNonNull(Int), "Int!"},
		Test{NewList(Int), "[Int]"},
		Test{NewNonNull(NewList(Int)), "[Int]!"},
		Test{NewList(NewNonNull(Int)), "[Int!]"},
		Test{NewList(NewList(Int)), "[[Int]]"},
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
		ttype    Type
		expected bool
	}
	tests := []Test{
		Test{Int, true},
		Test{objectType, false},
		Test{interfaceType, false},
		Test{unionType, false},
		Test{enumType, true},
		Test{inputObjectType, true},
	}
	for _, test := range tests {
		ttypeStr := fmt.Sprintf("%v", test.ttype)
		if IsInputType(test.ttype) != test.expected {
			t.Fatalf(`expected %v , got: %v`, test.expected, ttypeStr)
		}
		if IsInputType(NewList(test.ttype)) != test.expected {
			t.Fatalf(`expected %v , got: %v`, test.expected, ttypeStr)
		}
		if IsInputType(NewNonNull(test.ttype)) != test.expected {
			t.Fatalf(`expected %v , got: %v`, test.expected, ttypeStr)
		}
	}
}

func TestTypeSystem_DefinitionExample_IdentifiesOutputTypes(t *testing.T) {
	type Test struct {
		ttype    Type
		expected bool
	}
	tests := []Test{
		Test{Int, true},
		Test{objectType, true},
		Test{interfaceType, true},
		Test{unionType, true},
		Test{enumType, true},
		Test{inputObjectType, false},
	}
	for _, test := range tests {
		ttypeStr := fmt.Sprintf("%v", test.ttype)
		if IsOutputType(test.ttype) != test.expected {
			t.Fatalf(`expected %v , got: %v`, test.expected, ttypeStr)
		}
		if IsOutputType(NewList(test.ttype)) != test.expected {
			t.Fatalf(`expected %v , got: %v`, test.expected, ttypeStr)
		}
		if IsOutputType(NewNonNull(test.ttype)) != test.expected {
			t.Fatalf(`expected %v , got: %v`, test.expected, ttypeStr)
		}
	}
}

func TestTypeSystem_DefinitionExample_ProhibitsNestingNonNullInsideNonNull(t *testing.T) {
	ttype := NewNonNull(NewNonNull(Int))
	expected := `Can only create NonNull of a Nullable Type but got: Int!.`
	if ttype.GetError().Error() != expected {
		t.Fatalf(`expected %v , got: %v`, expected, ttype.GetError())
	}
}
func TestTypeSystem_DefinitionExample_ProhibitsNilInNonNull(t *testing.T) {
	ttype := NewNonNull(nil)
	expected := `Can only create NonNull of a Nullable Type but got: <nil>.`
	if ttype.GetError().Error() != expected {
		t.Fatalf(`expected %v , got: %v`, expected, ttype.GetError())
	}
}
func TestTypeSystem_DefinitionExample_ProhibitsNilTypeInUnions(t *testing.T) {
	ttype := NewUnion(UnionConfig{
		Name:  "BadUnion",
		Types: []*Object{nil},
	})
	expected := `BadUnion may only contain Object types, it cannot contain: <nil>.`
	if ttype.GetError().Error() != expected {
		t.Fatalf(`expected %v , got: %v`, expected, ttype.GetError())
	}
}
func TestTypeSystem_DefinitionExample_DoesNotMutatePassedFieldDefinitions(t *testing.T) {
	fields := FieldConfigMap{
		"field1": &FieldConfig{
			Type: String,
		},
		"field2": &FieldConfig{
			Type: String,
			Args: FieldConfigArgument{
				"id": &ArgumentConfig{
					Type: String,
				},
			},
		},
	}
	testObject1 := NewObject(ObjectConfig{
		Name:   "Test1",
		Fields: fields,
	})
	testObject2 := NewObject(ObjectConfig{
		Name:   "Test2",
		Fields: fields,
	})
	if !reflect.DeepEqual(testObject1.GetFields(), testObject2.GetFields()) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(testObject1.GetFields(), testObject2.GetFields()))
	}

	expectedFields := FieldConfigMap{
		"field1": &FieldConfig{
			Type: String,
		},
		"field2": &FieldConfig{
			Type: String,
			Args: FieldConfigArgument{
				"id": &ArgumentConfig{
					Type: String,
				},
			},
		},
	}
	if !reflect.DeepEqual(fields, expectedFields) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expectedFields, fields))
	}

	inputFields := InputObjectConfigFieldMap{
		"field1": &InputObjectFieldConfig{
			Type: String,
		},
		"field2": &InputObjectFieldConfig{
			Type: String,
		},
	}
	expectedInputFields := InputObjectConfigFieldMap{
		"field1": &InputObjectFieldConfig{
			Type: String,
		},
		"field2": &InputObjectFieldConfig{
			Type: String,
		},
	}
	testInputObject1 := NewInputObject(InputObjectConfig{
		Name:   "Test1",
		Fields: inputFields,
	})
	testInputObject2 := NewInputObject(InputObjectConfig{
		Name:   "Test2",
		Fields: inputFields,
	})
	if !reflect.DeepEqual(testInputObject1.GetFields(), testInputObject2.GetFields()) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(testInputObject1.GetFields(), testInputObject2.GetFields()))
	}
	if !reflect.DeepEqual(inputFields, expectedInputFields) {
		t.Fatalf("Unexpected result, Diff: %v", Diff(expectedInputFields, fields))
	}

}
