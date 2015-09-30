package executor_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/chris-ramon/graphql-go/executor"
	"github.com/chris-ramon/graphql-go/gqltypes"
	"github.com/chris-ramon/graphql-go/testutil"
)

// TODO: have a separate package for other tests for eg `parser`
// maybe for:
// - tests that supposed to be black-boxed (no reason to access private identifiers)
// - tests that create internal tests structs, we might not want to pollute the package with too many test structs

type testPic struct {
	Url    string `json:"url"`
	Width  string `json:"width"`
	Height string `json:"height"`
}

type testPicFn func(width, height string) *testPic

type testAuthor struct {
	Id            int          `json:"id"`
	Name          string       `json:"name"`
	Pic           testPicFn    `json:"pic"`
	RecentArticle *testArticle `json:"recentArticle"`
}
type testArticle struct {
	Id          string        `json:"id"`
	IsPublished string        `json:"isPublished"`
	Author      *testAuthor   `json:"author"`
	Title       string        `json:"title"`
	Body        string        `json:"body"`
	Hidden      string        `json:"hidden"`
	Keywords    []interface{} `json:"keywords"`
}

func getPic(id int, width, height string) *testPic {
	return &testPic{
		Url:    fmt.Sprintf("cdn://%v", id),
		Width:  width,
		Height: height,
	}
}

var johnSmith *testAuthor

func article(id interface{}) *testArticle {
	return &testArticle{
		Id:          fmt.Sprintf("%v", id),
		IsPublished: "true",
		Author:      johnSmith,
		Title:       fmt.Sprintf("My Article %v", id),
		Body:        "This is a post",
		Hidden:      "This data is not exposed in the schema",
		Keywords: []interface{}{
			"foo", "bar", 1, true, nil,
		},
	}
}

func TestExecutesUsingAComplexSchema(t *testing.T) {

	johnSmith = &testAuthor{
		Id:   123,
		Name: "John Smith",
		Pic: func(width string, height string) *testPic {
			return getPic(123, width, height)
		},
		RecentArticle: article("1"),
	}

	blogImage := gqltypes.NewGraphQLObjectType(gqltypes.GraphQLObjectTypeConfig{
		Name: "Image",
		Fields: gqltypes.GraphQLFieldConfigMap{
			"url": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLString,
			},
			"width": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLInt,
			},
			"height": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLInt,
			},
		},
	})
	blogAuthor := gqltypes.NewGraphQLObjectType(gqltypes.GraphQLObjectTypeConfig{
		Name: "Author",
		Fields: gqltypes.GraphQLFieldConfigMap{
			"id": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLString,
			},
			"name": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLString,
			},
			"pic": &gqltypes.GraphQLFieldConfig{
				Type: blogImage,
				Args: gqltypes.GraphQLFieldConfigArgumentMap{
					"width": &gqltypes.GraphQLArgumentConfig{
						Type: gqltypes.GraphQLInt,
					},
					"height": &gqltypes.GraphQLArgumentConfig{
						Type: gqltypes.GraphQLInt,
					},
				},
				Resolve: func(p gqltypes.GQLFRParams) interface{} {
					if author, ok := p.Source.(*testAuthor); ok {
						width := fmt.Sprintf("%v", p.Args["width"])
						height := fmt.Sprintf("%v", p.Args["height"])
						return author.Pic(width, height)
					}
					return nil
				},
			},
			"recentArticle": &gqltypes.GraphQLFieldConfig{},
		},
	})
	blogArticle := gqltypes.NewGraphQLObjectType(gqltypes.GraphQLObjectTypeConfig{
		Name: "Article",
		Fields: gqltypes.GraphQLFieldConfigMap{
			"id": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.NewGraphQLNonNull(gqltypes.GraphQLString),
			},
			"isPublished": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLBoolean,
			},
			"author": &gqltypes.GraphQLFieldConfig{
				Type: blogAuthor,
			},
			"title": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLString,
			},
			"body": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.GraphQLString,
			},
			"keywords": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.NewGraphQLList(gqltypes.GraphQLString),
			},
		},
	})

	blogAuthor.AddFieldConfig("recentArticle", &gqltypes.GraphQLFieldConfig{
		Type: blogArticle,
	})

	blogQuery := gqltypes.NewGraphQLObjectType(gqltypes.GraphQLObjectTypeConfig{
		Name: "Query",
		Fields: gqltypes.GraphQLFieldConfigMap{
			"article": &gqltypes.GraphQLFieldConfig{
				Type: blogArticle,
				Args: gqltypes.GraphQLFieldConfigArgumentMap{
					"id": &gqltypes.GraphQLArgumentConfig{
						Type: gqltypes.GraphQLID,
					},
				},
				Resolve: func(p gqltypes.GQLFRParams) interface{} {
					id := p.Args["id"]
					return article(id)
				},
			},
			"feed": &gqltypes.GraphQLFieldConfig{
				Type: gqltypes.NewGraphQLList(blogArticle),
				Resolve: func(p gqltypes.GQLFRParams) interface{} {
					return []*testArticle{
						article(1),
						article(2),
						article(3),
						article(4),
						article(5),
						article(6),
						article(7),
						article(8),
						article(9),
						article(10),
					}
				},
			},
		},
	})

	blogSchema, err := gqltypes.NewGraphQLSchema(gqltypes.GraphQLSchemaConfig{
		Query: blogQuery,
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	request := `
      {
        feed {
          id,
          title
        },
        article(id: "1") {
          ...articleFields,
          author {
            id,
            name,
            pic(width: 640, height: 480) {
              url,
              width,
              height
            },
            recentArticle {
              ...articleFields,
              keywords
            }
          }
        }
      }

      fragment articleFields on Article {
        id,
        isPublished,
        title,
        body,
        hidden,
        notdefined
      }
	`

	expected := &gqltypes.GraphQLResult{
		Data: map[string]interface{}{
			"article": map[string]interface{}{
				"title": "My Article 1",
				"body":  "This is a post",
				"author": map[string]interface{}{
					"id":   "123",
					"name": "John Smith",
					"pic": map[string]interface{}{
						"url":    "cdn://123",
						"width":  640,
						"height": 480,
					},
					"recentArticle": map[string]interface{}{
						"id":          "1",
						"isPublished": bool(true),
						"title":       "My Article 1",
						"body":        "This is a post",
						"keywords": []interface{}{
							"foo",
							"bar",
							"1",
							"true",
							nil,
						},
					},
				},
				"id":          "1",
				"isPublished": bool(true),
			},
			"feed": []interface{}{
				map[string]interface{}{
					"id":    "1",
					"title": "My Article 1",
				},
				map[string]interface{}{
					"id":    "2",
					"title": "My Article 2",
				},
				map[string]interface{}{
					"id":    "3",
					"title": "My Article 3",
				},
				map[string]interface{}{
					"id":    "4",
					"title": "My Article 4",
				},
				map[string]interface{}{
					"id":    "5",
					"title": "My Article 5",
				},
				map[string]interface{}{
					"id":    "6",
					"title": "My Article 6",
				},
				map[string]interface{}{
					"id":    "7",
					"title": "My Article 7",
				},
				map[string]interface{}{
					"id":    "8",
					"title": "My Article 8",
				},
				map[string]interface{}{
					"id":    "9",
					"title": "My Article 9",
				},
				map[string]interface{}{
					"id":    "10",
					"title": "My Article 10",
				},
			},
		},
	}

	// parse query
	ast := testutil.Parse(t, request)

	// execute
	ep := executor.ExecuteParams{
		Schema: blogSchema,
		AST:    ast,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
