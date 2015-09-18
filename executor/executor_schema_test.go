package executor_test

import (
	"fmt"
	"github.com/chris-ramon/graphql-go/executor"
	"github.com/chris-ramon/graphql-go/testutil"
	"github.com/chris-ramon/graphql-go/types"
	"github.com/kr/pretty"
	"reflect"
	"testing"
)

// TODO: have a separate package for other tests for eg `parser`
// maybe for:
// - tests that supposed to be black-boxed (no reason to access private identifiers)
// - tests that create internal tests structs, we might not want to pollute the package with too many test structs

type testPic struct {
	Url    string
	Width  string
	Height string
}

type testPicFn func(width, height string) *testPic

type testAuthor struct {
	Id            int
	Name          string
	Pic           testPicFn
	RecentArticle *testArticle
}
type testArticle struct {
	Id          string
	IsPublished string
	Author      *testAuthor
	Title       string
	Body        string
	Hidden      string
	Keywords    []interface{}
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

	blogImage := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
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
	blogAuthor := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "Author",
		Fields: types.GraphQLFieldConfigMap{
			"id": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Resolve: func(p types.GQLFRParams) interface{} {
					if author, ok := p.Source.(*testAuthor); ok {
						return author.Id
					}
					return nil
				},
			},
			"name": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Resolve: func(p types.GQLFRParams) interface{} {
					if author, ok := p.Source.(*testAuthor); ok {
						return author.Name
					}
					return nil
				},
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
				Resolve: func(p types.GQLFRParams) interface{} {
					if author, ok := p.Source.(*testAuthor); ok {
						width := fmt.Sprintf("%v", p.Args["width"])
						height := fmt.Sprintf("%v", p.Args["height"])
						return author.Pic(width, height)
					}
					return nil
				},
			},
			"recentArticle": &types.GraphQLFieldConfig{},
		},
	})
	blogArticle := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "Article",
		Fields: types.GraphQLFieldConfigMap{
			"id": &types.GraphQLFieldConfig{
				Type: types.NewGraphQLNonNull(types.GraphQLString),
				Resolve: func(p types.GQLFRParams) interface{} {
					if article, ok := p.Source.(*testArticle); ok {
						return article.Id
					}
					return nil
				},
			},
			"isPublished": &types.GraphQLFieldConfig{
				Type: types.GraphQLBoolean,
				Resolve: func(p types.GQLFRParams) interface{} {
					if article, ok := p.Source.(*testArticle); ok {
						return article.IsPublished
					}
					return false
				},
			},
			"author": &types.GraphQLFieldConfig{
				Type: blogAuthor,
				Resolve: func(p types.GQLFRParams) interface{} {
					if article, ok := p.Source.(*testArticle); ok {
						return article.Author
					}
					return nil
				},
			},
			"title": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Resolve: func(p types.GQLFRParams) interface{} {
					if article, ok := p.Source.(*testArticle); ok {
						return article.Title
					}
					return nil
				},
			},
			"body": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Resolve: func(p types.GQLFRParams) interface{} {
					if article, ok := p.Source.(*testArticle); ok {
						return article.Body
					}
					return nil
				},
			},
			"keywords": &types.GraphQLFieldConfig{
				Type: types.NewGraphQLList(types.GraphQLString),
				Resolve: func(p types.GQLFRParams) interface{} {
					if article, ok := p.Source.(*testArticle); ok {
						return article.Keywords
					}
					return nil
				},
			},
		},
	})

	blogAuthor.AddFieldConfig("recentArticle", &types.GraphQLFieldConfig{
		Type: blogArticle,
		Resolve: func(p types.GQLFRParams) interface{} {
			if author, ok := p.Source.(*testAuthor); ok {
				return author.RecentArticle
			}
			return nil
		},
	})

	blogQuery := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "Query",
		Fields: types.GraphQLFieldConfigMap{
			"article": &types.GraphQLFieldConfig{
				Type: blogArticle,
				Args: types.GraphQLFieldConfigArgumentMap{
					"id": &types.GraphQLArgumentConfig{
						Type: types.GraphQLID,
					},
				},
				Resolve: func(p types.GQLFRParams) interface{} {
					id := p.Args["id"]
					return article(id)
				},
			},
			"feed": &types.GraphQLFieldConfig{
				Type: types.NewGraphQLList(blogArticle),
				Resolve: func(p types.GQLFRParams) interface{} {
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

	blogSchema, err := types.NewGraphQLSchema(types.GraphQLSchemaConfig{
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

	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"article": map[string]interface{}{
				"title": "My Article 1",
				"body":  "This is a post",
				"author": map[string]interface{}{
					"id":   "123",
					"name": "John Smith",
					"pic": map[string]interface{}{
						"url":    "&{cdn://123 640 480}",
						"width":  int(0),
						"height": int(0),
					},
					"recentArticle": map[string]interface{}{
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
						"id": "1",
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
		t.Fatalf("Unexpected result, Diff: %v", pretty.Diff(expected, result))
	}
}
